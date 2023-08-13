package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-subscription/internal/entity"
)

type Subscription interface {
	GetProfileSubscriptions(userId uuid.UUID) []entity.Subscription
	ClaimProfileSubscription(input entity.SubscriptionInput) error
	UpdateProfileSubscription(input entity.SubscriptionInput) error
	SetProfileSubscriptionAutoRenewStatus(input entity.SubscriptionInput) error
	EndProfileSubscription(userId uuid.UUID, plan string, source string) error

	GetUserIdByGooglePurchaseToken(purchaseToken string) (*uuid.UUID, error)
}

type MQ interface {
	GetExpiringSubscriptions() []entity.ExpiringSubscription
	ImportPremiumVersion(userId, messageId uuid.UUID) error
	DeleteProfile(userId, messageId uuid.UUID) error
}
