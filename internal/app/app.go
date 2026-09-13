package app

import (
	"context"
	"fmt"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/shutdown"
	subscriptionpb "github.com/mephistolie/chefbook-backend-subscription/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-subscription/internal/config"
	"github.com/mephistolie/chefbook-backend-subscription/internal/logging"
	grpcRepo "github.com/mephistolie/chefbook-backend-subscription/internal/repository/grpc"
	"github.com/mephistolie/chefbook-backend-subscription/internal/repository/postgres"
	"github.com/mephistolie/chefbook-backend-subscription/internal/transport/dependencies/service"
	subscription "github.com/mephistolie/chefbook-backend-subscription/internal/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"time"
)

func Run(cfg *config.Config) {
	log.InitWithService("subscription", *cfg.LogsPath, *cfg.Environment == config.EnvDev)
	ctx := context.Background()
	cfg.Print(ctx)

	db, err := postgres.Connect(ctx, cfg.Database)
	if err != nil {
		logging.Events{}.StartupFailed(ctx, "connect_postgres")
		return
	}

	repository := postgres.NewRepository(db)

	grpcRepository, err := grpcRepo.NewRepository(cfg)
	if err != nil {
		logging.Events{}.StartupFailed(ctx, "initialize_grpc_clients")
		return
	}

	subscriptionService, err := service.New(ctx, repository, grpcRepository, cfg)
	if err != nil {
		logging.Events{}.StartupFailed(ctx, "initialize_service")
		return
	}

	mqSubscriber, err := NewMqConsumer(ctx, cfg.Amqp, subscriptionService.MQ)
	if err != nil {
		logging.Events{}.StartupFailed(ctx, "initialize_mq_consumer")
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *cfg.Port))
	if err != nil {
		logging.Events{}.StartupFailed(ctx, "listen_grpc")
		return
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			log.UnaryServerInterceptor(),
		),
	)

	healthServer := health.NewServer()
	subscriptionServer := subscription.NewServer(subscriptionService.Subscription)

	go monitorHealthChecking(ctx, db, healthServer)

	subscriptionpb.RegisterSubscriptionServiceServer(grpcServer, subscriptionServer)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	logging.Events{}.GRPCServerStarted(ctx)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logging.Events{}.GRPCServerFailed(ctx, err)
		}
	}()

	wait := shutdown.Graceful(ctx, 5*time.Second, map[string]shutdown.Operation{
		"grpc-server": func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
		"database": func(ctx context.Context) error {
			return db.Close()
		},
		"services": func(ctx context.Context) error {
			return grpcRepository.Stop()
		},
		"mq": func(ctx context.Context) error {
			if mqSubscriber != nil {
				_ = mqSubscriber.Stop()
			}
			return nil
		},
	})
	<-wait
}
