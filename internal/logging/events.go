package logging

import (
	"context"
	"errors"

	"github.com/mephistolie/chefbook-backend-common/log"
)

type Events struct{}

var (
	errGooglePlayRequest         = errors.New("google play request failed")
	errGoogleSubscriptionMapping = errors.New("google subscription mapping failed")
	errGoogleTimestampFormat     = errors.New("invalid Google subscription timestamp")
)

type SubscriptionData struct {
	UserID    string
	Plan      string
	Source    string
	AutoRenew bool
}

type MessageData struct {
	MessageID string
	UserID    string
}

func (Events) ProfileSubscriptionsQueryFailed(ctx context.Context, userID string, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.storage.profile_subscriptions.query_failed",
		Message:   "failed to query profile subscriptions",
		Component: log.ComponentPostgres,
		UserID:    userID,
		Operation: "query_profile_subscriptions",
	}, err)
}

func (Events) ProfileSubscriptionScanFailed(ctx context.Context, userID string, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.storage.profile_subscriptions.scan_failed",
		Message:   "failed to scan profile subscription",
		Component: log.ComponentPostgres,
		UserID:    userID,
		Operation: "scan_profile_subscription",
	}, err)
}

func (Events) ProfileSubscriptionsIterationFailed(ctx context.Context, userID string, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.storage.profile_subscriptions.iteration_failed",
		Message:   "failed while iterating profile subscriptions",
		Component: log.ComponentPostgres,
		UserID:    userID,
		Operation: "iterate_profile_subscriptions",
	}, err)
}

func (Events) SubscriptionCreateFailed(ctx context.Context, data SubscriptionData, err error) {
	log.LogError(ctx, subscriptionEvent(
		"subscription.storage.subscription.create_failed",
		"failed to create subscription",
		"create_subscription",
		data,
	), err)
}

func (Events) SubscriptionUpdateFailed(ctx context.Context, data SubscriptionData, err error) {
	log.LogError(ctx, subscriptionEvent(
		"subscription.storage.subscription.update_failed",
		"failed to update subscription",
		"update_subscription",
		data,
	), err)
}

func (Events) SubscriptionAutoRenewUpdateFailed(ctx context.Context, data SubscriptionData, err error) {
	log.LogError(ctx, subscriptionEvent(
		"subscription.storage.subscription.auto_renew_update_failed",
		"failed to update subscription auto-renew status",
		"update_subscription_auto_renew",
		data,
	), err)
}

func (Events) SubscriptionEndFailed(ctx context.Context, data SubscriptionData, err error) {
	log.LogError(ctx, subscriptionEvent(
		"subscription.storage.subscription.end_failed",
		"failed to end subscription",
		"end_subscription",
		data,
	), err)
}

func (Events) ExpiringSubscriptionsQueryFailed(ctx context.Context, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.storage.expiring_subscriptions.query_failed",
		Message:   "failed to query expiring subscriptions",
		Component: log.ComponentPostgres,
		Operation: "query_expiring_subscriptions",
	}, err)
}

func (Events) ExpiringSubscriptionScanFailed(ctx context.Context, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.storage.expiring_subscriptions.scan_failed",
		Message:   "failed to scan expiring subscription",
		Component: log.ComponentPostgres,
		Operation: "scan_expiring_subscription",
	}, err)
}

func (Events) ExpiringSubscriptionsIterationFailed(ctx context.Context, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.storage.expiring_subscriptions.iteration_failed",
		Message:   "failed while iterating expiring subscriptions",
		Component: log.ComponentPostgres,
		Operation: "iterate_expiring_subscriptions",
	}, err)
}

func (Events) PremiumImportFailed(ctx context.Context, data MessageData, err error) {
	log.LogWarnError(ctx, log.Event{
		Event:     "subscription.storage.premium_import.failed",
		Message:   "failed to import premium subscription",
		Component: log.ComponentPostgres,
		MessageID: data.MessageID,
		UserID:    data.UserID,
		Operation: "import_premium_subscription",
		Payload: map[string]any{
			"source": "firebase",
		},
	}, err)
}

func (Events) ProfileSubscriptionsDeleteFailed(ctx context.Context, data MessageData, err error) {
	log.LogWarnError(ctx, log.Event{
		Event:     "subscription.storage.profile_subscriptions.delete_failed",
		Message:   "failed to delete profile subscriptions",
		Component: log.ComponentPostgres,
		MessageID: data.MessageID,
		UserID:    data.UserID,
		Operation: "delete_profile_subscriptions",
	}, err)
}

func (Events) InboxMessageStoreFailed(ctx context.Context, messageID string, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.storage.inbox.store_failed",
		Message:   "failed to store inbox message",
		Component: log.ComponentPostgres,
		MessageID: messageID,
		Operation: "store_inbox_message",
	}, err)
}

func (Events) TransactionBeginFailed(ctx context.Context, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.storage.transaction.begin_failed",
		Message:   "failed to begin database transaction",
		Component: log.ComponentPostgres,
		Operation: "begin_transaction",
	}, err)
}

func (Events) TransactionCommitFailed(ctx context.Context, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.storage.transaction.commit_failed",
		Message:   "failed to commit database transaction",
		Component: log.ComponentPostgres,
		Operation: "commit_transaction",
	}, err)
}

func (Events) GooglePurchaseLookupFailed(ctx context.Context, operation string, err error) {
	log.LogWarnError(ctx, log.Event{
		Event:     "subscription.google.purchase.lookup_failed",
		Message:   "failed to look up Google purchase owner",
		Component: log.ComponentPostgres,
		Operation: operation,
		Payload: map[string]any{
			"provider": "google_play",
		},
	}, err)
}

func (Events) GooglePurchaseAcknowledgeFailed(ctx context.Context, userID string) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.google.purchase.acknowledge_failed",
		Message:   "failed to acknowledge Google purchase",
		Component: "google_play",
		UserID:    userID,
		Operation: "acknowledge_purchase",
		Payload: map[string]any{
			"provider": "google_play",
			"status":   "request_failed",
		},
		ErrorType: "provider_request_failed",
	}, errGooglePlayRequest)
}

func (Events) GoogleSubscriptionMappingFailed(ctx context.Context, notificationType int) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.google.subscription.mapping_failed",
		Message:   "failed to map Google subscription",
		Component: "google_play",
		Operation: "map_subscription",
		Payload: map[string]any{
			"notification_type": notificationType,
			"provider":          "google_play",
		},
		ErrorType: "unknown_subscription",
	}, errGoogleSubscriptionMapping)
}

func (Events) GoogleSubscriptionsDisabled(ctx context.Context) {
	log.LogWarn(ctx, log.Event{
		Event:     "subscription.google.disabled",
		Message:   "Google subscriptions are disabled",
		Component: "google_play",
		Operation: "validate_purchase",
		Payload: map[string]any{
			"provider": "google_play",
		},
	})
}

func (Events) GooglePurchaseValidationFailed(ctx context.Context) {
	log.LogWarnError(ctx, log.Event{
		Event:     "subscription.google.purchase.validation_failed",
		Message:   "Google purchase validation failed",
		Component: "google_play",
		Operation: "validate_purchase",
		Payload: map[string]any{
			"provider": "google_play",
			"status":   "request_failed",
		},
		ErrorType: "provider_request_failed",
	}, errGooglePlayRequest)
}

func (Events) GoogleSubscriptionTimestampParseFailed(ctx context.Context, field string) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.google.subscription.timestamp_parse_failed",
		Message:   "failed to parse Google subscription timestamp",
		Component: "google_play",
		Operation: "parse_subscription_timestamp",
		Payload: map[string]any{
			"field":    field,
			"provider": "google_play",
			"status":   "invalid_format",
		},
		ErrorType: "invalid_provider_response",
	}, errGoogleTimestampFormat)
}

func (Events) SubscriptionMailSending(ctx context.Context, plan string) {
	log.Log(ctx, log.Event{
		Event:     "subscription.mail.plan_change.sending",
		Message:   "sending subscription plan change mail",
		Component: "mail",
		Operation: "send_plan_change_mail",
		Payload: map[string]any{
			"plan": plan,
		},
	})
}

func (Events) SubscriptionMailRenderFailed(ctx context.Context, plan string, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.mail.plan_change.render_failed",
		Message:   "failed to render subscription plan change mail",
		Component: "mail",
		Operation: "render_plan_change_mail",
		Payload: map[string]any{
			"plan": plan,
		},
	}, err)
}

func (Events) SubscriptionMailSendFailed(ctx context.Context, plan string, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.mail.plan_change.send_failed",
		Message:   "failed to send subscription plan change mail",
		Component: "mail",
		Operation: "send_plan_change_mail",
		Payload: map[string]any{
			"plan": plan,
		},
	}, err)
}

func (Events) DependencyInitialized(ctx context.Context, dependency string) {
	log.Log(ctx, log.Event{
		Event:     "subscription.dependency.initialized",
		Message:   "service dependency initialized",
		Component: "dependencies",
		Operation: "initialize_dependency",
		Payload: map[string]any{
			"dependency": dependency,
		},
	})
}

func (Events) GoogleNotificationDecodeFailed(ctx context.Context, messageID string, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.google.notification.decode_failed",
		Message:   "failed to decode Google subscription notification",
		Component: "google_pubsub",
		MessageID: messageID,
		Operation: "decode_notification",
		Payload: map[string]any{
			"provider": "google_play",
		},
	}, err)
}

func (Events) GoogleNotificationUnmarshalFailed(ctx context.Context, messageID string, err error) {
	log.LogError(ctx, log.Event{
		Event:     "subscription.google.notification.unmarshal_failed",
		Message:   "failed to unmarshal Google subscription notification",
		Component: "google_pubsub",
		MessageID: messageID,
		Operation: "unmarshal_notification",
		Payload: map[string]any{
			"provider": "google_play",
		},
	}, err)
}

func subscriptionEvent(event, message, operation string, data SubscriptionData) log.Event {
	return log.Event{
		Event:     event,
		Message:   message,
		Component: log.ComponentPostgres,
		UserID:    data.UserID,
		Operation: operation,
		Payload: map[string]any{
			"auto_renew": data.AutoRenew,
			"plan":       data.Plan,
			"source":     data.Source,
		},
	}
}
