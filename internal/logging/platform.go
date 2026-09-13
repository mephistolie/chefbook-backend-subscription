package logging

import (
	"context"
	"errors"

	"github.com/mephistolie/chefbook-backend-common/log"
)

func (Events) ConfigLoaded(ctx context.Context) {
	log.Log(ctx, log.Event{
		Event:     "config.loaded",
		Message:   "service configuration loaded",
		Component: "config",
	})
}

func (Events) StartupFailed(ctx context.Context, stage string) {
	log.LogFatal(ctx, log.Event{
		Event:     "app.startup.failed",
		Message:   "service startup failed",
		Component: "app",
		Operation: stage,
		ErrorType: "startup_stage_failed",
	}, errors.New("service startup stage failed"))
}

func (Events) GRPCServerStarted(ctx context.Context) {
	log.Log(ctx, log.Event{
		Event:     "grpc.server.started",
		Message:   "grpc server started",
		Component: log.ComponentGRPC,
	})
}

func (Events) GRPCServerFailed(ctx context.Context, err error) {
	log.LogError(ctx, log.Event{
		Event:     "grpc.server.failed",
		Message:   "error occurred while running grpc server",
		Component: log.ComponentGRPC,
	}, err)
}

func (Events) DatabaseHealthCheckFailed(ctx context.Context, err error) {
	log.LogWarnError(ctx, log.Event{
		Event:     "postgres.health_check.failed",
		Message:   "database is unavailable",
		Component: log.ComponentPostgres,
	}, err)
}

func (Events) MQConsumerInitialized(ctx context.Context) {
	log.Log(ctx, log.Event{
		Event:     "mq.consumer.initialized",
		Message:   "mq consumer initialized",
		Component: log.ComponentAMQP,
	})
}

func (Events) MQMessageProcessing(ctx context.Context, messageID, messageType string) {
	log.LogDebug(ctx, log.Event{
		Event:     "mq.message.processing",
		Message:   "processing message",
		Component: log.ComponentAMQP,
		MessageID: messageID,
		Payload: map[string]any{
			"message_type": messageType,
		},
	})
}

func (Events) MQMessageTypeUnsupported(ctx context.Context, messageID, messageType string) {
	log.LogWarn(ctx, log.Event{
		Event:     "mq.message.unsupported_type",
		Message:   "got unsupported message type",
		Component: log.ComponentAMQP,
		MessageID: messageID,
		Payload: map[string]any{
			"message_type": messageType,
		},
	})
}

func (Events) FirebaseProfileLoadFailed(ctx context.Context, userID string) {
	log.LogWarnError(ctx, log.Event{
		Event:     "firebase.profile.load_failed",
		Message:   "unable to get firebase profile",
		Component: log.ComponentFirebase,
		UserID:    userID,
		ErrorType: "provider_request_failed",
	}, errors.New("firebase request failed"))
}

func (Events) ProfileDeletedMessageProcessing(ctx context.Context, messageID, userID string) {
	log.LogDebug(ctx, log.Event{
		Event:     "profile.deleted.message.processing",
		Message:   "processing profile deleted message",
		Component: log.ComponentAMQP,
		MessageID: messageID,
		UserID:    userID,
	})
}
