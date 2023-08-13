package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/mephistolie/chefbook-backend-common/log"
	"google.golang.org/api/option"
)
import _ "cloud.google.com/go/pubsub"

type SubscriptionEventConsumer struct {
	subscription *pubsub.Subscription
	cancelFunc   *context.CancelFunc
}

func NewSubscriptionEventConsumer(projectId, subscriptionId string, credentialsJson []byte) (*SubscriptionEventConsumer, error) {
	creds := option.WithCredentialsJSON(credentialsJson)
	client, err := pubsub.NewClient(context.Background(), projectId, creds)
	if err != nil {
		return nil, err
	}

	subscription := client.SubscriptionInProject(subscriptionId, projectId)
	return &SubscriptionEventConsumer{subscription: subscription}, nil
}

func (c *SubscriptionEventConsumer) Subscribe(handler SubscriptionEventHandler) error {
	c.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	c.cancelFunc = &cancel

	err := c.subscription.Receive(ctx, func(_ context.Context, m *pubsub.Message) {
		var data []byte
		if _, err := base64.StdEncoding.Decode(data, m.Data); err != nil {
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
		if err := handler.Handle(event); err != nil {
			m.Nack()
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
