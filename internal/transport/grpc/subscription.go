package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	api "github.com/mephistolie/chefbook-backend-subscription/api/proto/implementation/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SubscriptionServer) GetProfileSubscriptions(_ context.Context, req *api.GetProfileSubscriptionsRequest) (*api.GetProfileSubscriptionsResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	subscriptions := s.service.GetProfileSubscriptions(userId)
	dtos := make([]*api.Subscription, len(subscriptions))

	for i, sub := range subscriptions {
		var expirationDate *timestamppb.Timestamp
		if sub.Expiration != nil {
			expirationDate = timestamppb.New(*sub.Expiration)
		}
		dtos[i] = &api.Subscription{
			Plan:           sub.Plan,
			Source:         sub.Source,
			ExpirationDate: expirationDate,
			AutoRenew:      sub.AutoRenew,
		}
	}

	return &api.GetProfileSubscriptionsResponse{Subscriptions: dtos}, nil
}

func (s *SubscriptionServer) GetProfileCurrentSubscription(_ context.Context, req *api.GetProfileCurrentSubscriptionRequest) (*api.GetProfileCurrentSubscriptionResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	sub := s.service.GetProfileCurrentSubscription(userId)
	var expirationDate *timestamppb.Timestamp
	if sub.Expiration != nil {
		expirationDate = timestamppb.New(*sub.Expiration)
	}

	return &api.GetProfileCurrentSubscriptionResponse{
		Plan:           sub.Plan,
		ExpirationDate: expirationDate,
	}, nil
}
