package rabbitmq

import (
	"context"

	"github.com/embrace-century/basin-warden/internal/brokers"
)

type deliveryMode struct{}
type routingKey struct{}
type exchangeName struct{}
type queueName struct{}
type routingKeys struct{}

// type exchangeName struct{}
// type exchangeType struct{}

func RoutingKey(value string) brokers.PublishOption {
	return setPublishOption(routingKey{}, value)
}

func DeliveryMode(value uint8) brokers.PublishOption {
	return setPublishOption(deliveryMode{}, value)
}

func QueueName(value string) brokers.SubscribeOption {
	return setSubscribeOption(queueName{}, value)
}

func RoutingKeys(value []string) brokers.SubscribeOption {
	return setSubscribeOption(routingKeys{}, value)
}

func setPublishOption(k, v interface{}) brokers.PublishOption {
	return func(opt *brokers.PublishOptions) {
		if opt.Context == nil {
			opt.Context = context.Background()
		}
		opt.Context = context.WithValue(opt.Context, k, v)
	}
}

func setSubscribeOption(k, v interface{}) brokers.SubscribeOption {
	return func(opt *brokers.SubscribeOptions) {
		if opt.Context == nil {
			opt.Context = context.Background()
		}
		opt.Context = context.WithValue(opt.Context, k, v)
	}
}
