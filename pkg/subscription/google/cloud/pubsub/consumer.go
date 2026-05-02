package pubsub

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"cloud.google.com/go/pubsub/v2"
	"github.com/mephistolie/chefbook-backend-common/log"
	"google.golang.org/api/option"
)

type SubscriptionEventConsumer struct {
	subscriber *pubsub.Subscriber
	cancelFunc *context.CancelFunc
}

func NewSubscriptionEventConsumer(ctx context.Context, projectId, subscriptionId string, credentialsJson []byte) (*SubscriptionEventConsumer, error) {
	creds := option.WithCredentialsJSON(credentialsJson)
	client, err := pubsub.NewClient(ctx, projectId, creds)
	if err != nil {
		return nil, err
	}

	return &SubscriptionEventConsumer{subscriber: client.Subscriber(subscriptionId)}, nil
}

func (c *SubscriptionEventConsumer) Subscribe(ctx context.Context, handler SubscriptionEventHandler) error {
	c.Stop()
	ctx, cancel := context.WithCancel(ctx)
	c.cancelFunc = &cancel

	err := c.subscriber.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		data, err := base64.StdEncoding.DecodeString(string(m.Data))
		if err != nil {
			log.Errorf("unable to decode message %s: %s", m.ID, err)
			m.Nack()
			return
		}

		var notification DeveloperNotification
		if err := json.Unmarshal(data, &notification); err != nil {
			log.Errorf("unable to unmarshal message %s with body %s: %s", m.ID, data, err)
			m.Nack()
			return
		}

		event := SubscriptionEvent{
			NotificationType: notification.SubscriptionNotification.NotificationType,
			SubscriptionId:   notification.SubscriptionNotification.SubscriptionId,
			PurchaseToken:    notification.SubscriptionNotification.PurchaseToken,
		}
		if err := handler.HandleSubscriptionEvent(ctx, event); err != nil {
			m.Nack()
			return
		}

		m.Ack()
	})

	return err
}

func (c *SubscriptionEventConsumer) Stop() {
	if c.cancelFunc != nil {
		(*c.cancelFunc)()
	}
}
