package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-subscription/internal/entity"
)

type Subscription interface {
	GetProfileSubscriptions(ctx context.Context, userId uuid.UUID) []entity.Subscription
	ClaimProfileSubscription(ctx context.Context, input entity.SubscriptionInput) error
	UpdateProfileSubscription(ctx context.Context, input entity.SubscriptionInput) error
	SetProfileSubscriptionAutoRenewStatus(ctx context.Context, input entity.SubscriptionInput) error
	EndProfileSubscription(ctx context.Context, userId uuid.UUID, plan string, source string) error

	GetUserIdByGooglePurchaseToken(ctx context.Context, purchaseToken string) (*uuid.UUID, error)
}

type MQ interface {
	GetExpiringSubscriptions(ctx context.Context) []entity.ExpiringSubscription
	ImportPremiumVersion(ctx context.Context, userId, messageId uuid.UUID) error
	DeleteProfile(ctx context.Context, userId, messageId uuid.UUID) error
}
