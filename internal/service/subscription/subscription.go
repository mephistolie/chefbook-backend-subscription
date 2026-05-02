package subscription

import (
	"context"

	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/subscription"
	"github.com/mephistolie/chefbook-backend-subscription/internal/entity"
)

var freeSubscription = entity.Subscription{Plan: subscription.PlanFree}

func (s *Service) GetProfileSubscriptions(ctx context.Context, userId uuid.UUID) []entity.Subscription {
	subscriptions := s.repo.GetProfileSubscriptions(ctx, userId)
	subscriptions = append(subscriptions, freeSubscription)
	return subscriptions
}

func (s *Service) GetProfileCurrentSubscription(ctx context.Context, userId uuid.UUID) entity.Subscription {
	subscriptions := s.repo.GetProfileSubscriptions(ctx, userId)
	if len(subscriptions) > 0 {
		return subscriptions[0]
	}

	return entity.Subscription{Plan: subscription.PlanFree}
}
