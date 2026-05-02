package pubsub

import "context"

type SubscriptionEvent struct {
	NotificationType int
	SubscriptionId   string
	PurchaseToken    string
}

type SubscriptionEventHandler interface {
	HandleSubscriptionEvent(ctx context.Context, event SubscriptionEvent) error
}
