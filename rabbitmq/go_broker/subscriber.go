package rabbitmq

import (
	"errors"
	"fmt"

	"github.com/embrace-century/basin-warden/internal/brokers"
	"github.com/streadway/amqp"
)

type subscriber struct {
	channelManager *ChannelManager
	channel        *rabbitmqChannel
	handler        brokers.MessageHandler
	opts           brokers.SubscribeOptions
	exchange       *exchange
}

func (sub *subscriber) obtainQueue(qc QueueConfig) (q *amqp.Queue, err error) {
	if queue, err := sub.channel.channel.QueueDeclare(qc.Name, qc.Durable, qc.AutoDelete, qc.Exclusive, qc.NoWait, nil); err != nil {
		if errors.As(fmt.Errorf("channel/connection is not open"), &err) {
			ch, err := sub.channelManager.obtainChannel()
			if err != nil {
				return nil, err
			}
			sub.channel = ch
			return sub.obtainQueue(qc)
		} else {
			return nil, fmt.Errorf("obtainQueue err: [%w]", err)
		}
	} else {
		return &queue, nil
	}
}

func (sub *subscriber) run() error {
	var rKeys []string
	if sub.opts.Context != nil {
		if value, ok := sub.opts.Context.Value(routingKeys{}).([]string); ok {
			rKeys = value
		}
	}
	qc := NewQueueConfig(func(q *QueueConfig) { q.Name = sub.opts.Queue })

	q, err := sub.obtainQueue(qc)
	if err != nil {
		return fmt.Errorf("MQ注册队列失败: [%w]", err)
	}
	// 队列绑定
	for _, k := range rKeys {
		err = sub.channel.channel.QueueBind(q.Name, k, sub.exchange.name, true, nil)
		if err != nil {
			return fmt.Errorf("MQ绑定队列失败: [%w]", err)
		}
	}

	// 设置 QoS：
	// Qos controls how many messages or how many bytes the server will try to keep on the network for consumers before receiving delivery acks.
	// The intent of Qos is to make sure the network buffers stay full between the server and client.
	err = sub.channel.channel.Qos(
		1, // prefetch count
		0, // prefetch size
		// When global is true, these Qos settings apply to all existing and future consumers on all channels on the same connection.
		// When false, the Channel.Qos settings will apply to all existing and future consumers on this channel.
		false, // global
	)
	if err != nil {
		return err
	}

	// 消费
	delivery, err := sub.channel.channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			sub.channelManager.returnChannel(sub.channel)
		}()
		for d := range delivery {
			message := brokers.ReceivedMessage{
				ContentType: d.ContentType,
				Body:        d.Body,
				MessageID:   d.MessageId,
				Timestamp:   d.Timestamp,
			}
			err = sub.handler.Handler(&message)
			if err != nil {
				_ = d.Nack(false, false)
			} else {
				_ = d.Ack(false) // When multiple is true, this delivery and all prior unacknowledged deliveries on the same channel will be acknowledged.
			}
		}

	}()
	return nil
}
