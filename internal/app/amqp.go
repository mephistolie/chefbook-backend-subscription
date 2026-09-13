package app

import (
	"context"
	auth "github.com/mephistolie/chefbook-backend-auth/api/mq"
	"github.com/mephistolie/chefbook-backend-common/mq/config"
	mqConsumer "github.com/mephistolie/chefbook-backend-common/mq/consumer"
	mqApi "github.com/mephistolie/chefbook-backend-common/mq/dependencies"
	"github.com/mephistolie/chefbook-backend-subscription/internal/logging"
	amqp "github.com/wagslane/go-rabbitmq"
)

const queueProfiles = "subscription.profiles"

var supportedMsgTypes = []string{
	auth.MsgTypeProfileFirebaseImport,
	auth.MsgTypeProfileDeleted,
}

func NewMqConsumer(
	ctx context.Context,
	cfg config.Amqp,
	service mqApi.Inbox,
) (*mqConsumer.Consumer, error) {
	var consumer *mqConsumer.Consumer = nil
	var err error

	if len(*cfg.Host) > 0 {

		consumer, err = mqConsumer.New(cfg, service, supportedMsgTypes)
		if err != nil {
			return nil, err
		}
		if err = consumer.Start(
			mqConsumer.Params{
				QueueName: queueProfiles,
				Options: []func(*amqp.ConsumerOptions){
					amqp.WithConsumerOptionsQueueQuorum,
					amqp.WithConsumerOptionsQueueDurable,
					amqp.WithConsumerOptionsExchangeName(auth.ExchangeProfiles),
					amqp.WithConsumerOptionsExchangeKind("fanout"),
					amqp.WithConsumerOptionsExchangeDurable,
					amqp.WithConsumerOptionsExchangeDeclare,
					amqp.WithConsumerOptionsRoutingKey(""),
				},
			},
		); err != nil {
			return nil, err
		}

		logging.Events{}.MQConsumerInitialized(ctx)
	}

	return consumer, nil
}
