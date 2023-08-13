package pubsub

type SubscriptionEvent struct {
	NotificationType int
	SubscriptionId   string
	PurchaseToken    string
}

type SubscriptionEventHandler interface {
	HandleSubscriptionEvent(event SubscriptionEvent) error
}
