package pubsub

const (
	NotificationTypeSubscriptionRenewed       = 2
	NotificationTypeSubscriptionCanceled      = 3
	NotificationTypeSubscriptionOnHold        = 5
	NotificationTypeSubscriptionInGracePeriod = 6
	NotificationTypeSubscriptionRestarted     = 7
	NotificationTypeSubscriptionRevoked       = 12
	NotificationTypeSubscriptionExpired       = 13
)

type DeveloperNotification struct {
	Version                  string                   `json:"version"`
	PackageName              string                   `json:"package_name"`
	EventTimeMillis          int64                    `json:"eventTimeMillis"`
	SubscriptionNotification SubscriptionNotification `json:"subscriptionNotification"`
}

type SubscriptionNotification struct {
	Version          string `json:"version"`
	NotificationType int    `json:"notificationType"`
	PurchaseToken    string `json:"purchaseToken"`
	SubscriptionId   string `json:"subscriptionId"`
}
