package app

import (
	"context"
	"github.com/jmoiron/sqlx"
	subscriptionpb "github.com/mephistolie/chefbook-backend-subscription/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-subscription/internal/logging"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

func monitorHealthChecking(ctx context.Context, db *sqlx.DB, healthServer *health.Server) {
	for {
		status := healthpb.HealthCheckResponse_SERVING
		if err := db.PingContext(ctx); err != nil {
			status = healthpb.HealthCheckResponse_NOT_SERVING
			logging.Events{}.DatabaseHealthCheckFailed(ctx, err)
		}
		setHealthStatus(healthServer, status)
		time.Sleep(1 * time.Minute)
	}
}

func setHealthStatus(healthServer *health.Server, status healthpb.HealthCheckResponse_ServingStatus) {
	healthServer.SetServingStatus("", status)
	healthServer.SetServingStatus(subscriptionpb.SubscriptionService_ServiceDesc.ServiceName, status)
}
