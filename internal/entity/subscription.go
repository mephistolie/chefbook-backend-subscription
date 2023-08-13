package entity

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	Plan       string
	Source     *string
	AutoRenew  bool
	Expiration *time.Time
}

type SubscriptionInput struct {
	UserId     uuid.UUID
	Plan       string
	Source     string
	AutoRenew  bool
	Start      time.Time
	Expiration time.Time
}

type ExpiringSubscription struct {
	UserId uuid.UUID
	Plan   string
	Source *string
}
