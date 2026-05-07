package app

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-backend-common/log"
	subscriptionpb "github.com/mephistolie/chefbook-backend-subscription/api/proto/implementation/v1"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

func monitorHealthChecking(db *sqlx.DB, healthServer *health.Server) {
	for {
		status := healthpb.HealthCheckResponse_SERVING
		if db.Ping() != nil {
			status = healthpb.HealthCheckResponse_NOT_SERVING
			log.LogWarn(context.Background(), log.Event{
				Event:     "postgres.health_check.failed",
				Message:   "database is unavailable",
				Component: log.ComponentPostgres,
			})
		}
		setHealthStatus(healthServer, status)
		time.Sleep(1 * time.Minute)
	}
}

func setHealthStatus(healthServer *health.Server, status healthpb.HealthCheckResponse_ServingStatus) {
	healthServer.SetServingStatus("", status)
	healthServer.SetServingStatus(subscriptionpb.SubscriptionService_ServiceDesc.ServiceName, status)
}
