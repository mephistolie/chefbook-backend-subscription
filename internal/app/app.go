package app

import (
	"context"
	"fmt"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/shutdown"
	subscriptionpb "github.com/mephistolie/chefbook-backend-subscription/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-subscription/internal/config"
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
	log.Init(*cfg.LogsPath, *cfg.Environment == config.EnvDev)
	cfg.Print()

	db, err := postgres.Connect(cfg.Database)
	if err != nil {
		log.Fatal(err)
		return
	}

	repository := postgres.NewRepository(db)

	grpcRepository, err := grpcRepo.NewRepository(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	subscriptionService, err := service.New(repository, grpcRepository, cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	mqSubscriber, err := NewMqConsumer(cfg.Amqp, subscriptionService.MQ)
	if err != nil {
		log.Fatal(err)
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *cfg.Port))
	if err != nil {
		log.Fatal(err)
		return
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			log.UnaryServerInterceptor(),
		),
	)

	healthServer := health.NewServer()
	subscriptionServer := subscription.NewServer(subscriptionService.Subscription)

	go monitorHealthChecking(db, healthServer)

	subscriptionpb.RegisterSubscriptionServiceServer(grpcServer, subscriptionServer)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Errorf("error occurred while running http server: %s\n", err.Error())
		} else {
			log.Info("gRPC server started")
		}
	}()

	wait := shutdown.Graceful(context.Background(), 5*time.Second, map[string]shutdown.Operation{
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
