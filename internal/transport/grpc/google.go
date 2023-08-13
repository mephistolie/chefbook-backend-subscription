package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	api "github.com/mephistolie/chefbook-backend-subscription/api/proto/implementation/v1"
)

func (s *SubscriptionServer) ConfirmGoogleSubscription(_ context.Context, req *api.ConfirmGoogleSubscriptionRequest) (*api.ConfirmGoogleSubscriptionResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	if err = s.service.ConfirmGoogleSubscription(userId, req.SubscriptionId, req.PurchaseToken); err != nil {
		return nil, err
	}

	return &api.ConfirmGoogleSubscriptionResponse{Message: "Google subscription activated"}, nil
}
