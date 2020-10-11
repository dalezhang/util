package rabbitmq

import (
	"github.com/embrace-century/basin-warden/internal/brokers"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type ChannelManager struct {
	conn            *amqp.Connection
	channels        chan *rabbitmqChannel
	channelCount    int
	exchangeManager *exchangeManager
}

// obtainChannel 同时最多有5个channel，如果有空闲的channel就弹出一个，如果没有就阻塞等待其他go routing返还channel
// 使用后请调用 returnChannel(ch *rabbitmqChannel) 返还
func (cm *ChannelManager) obtainChannel() (ch *rabbitmqChannel, err error) {
	if len(cm.channels) < 1 && cm.channelCount < 5 {
		if ch, err := cm.conn.Channel(); err != nil {
			return nil, err
		} else {
			cm.channelCount = cm.channelCount + 1
			cm.channels <- &rabbitmqChannel{channel: ch, uuid: uuid.New().String()}
		}
	}
	for ch = range cm.channels {
		break
	}
	return ch, nil
}

// returnChannel 返还  obtainChannel() (ch *rabbitmqChannel, err error) 生成或弹出的channel
// 如果ChannelManager里的channels数量大于5个，则关闭并丢弃该channel
func (cm *ChannelManager) returnChannel(ch *rabbitmqChannel) {
	if len(cm.channels) > 5 {
		ch.close()
	} else {
		cm.channels <- ch
	}
}

func (cm *ChannelManager) publish(exchange *exchange, message *brokers.Message, opts ...brokers.PublishOption) error {
	ch, err := cm.obtainChannel()
	if err != nil {
		return err
	}
	defer func() {
		cm.returnChannel(ch)
	}()

	key := defaultRoutingKey

	m := amqp.Publishing{
		ContentType: message.ContentType,
		Body:        message.Body,
		AppId:       "basin-warden",
	}

	options := brokers.PublishOptions{}
	for _, oFun := range opts {
		oFun(&options)
	}

	if options.Context != nil {
		if value, ok := options.Context.Value(routingKey{}).(string); ok {
			key = value
		}
		if value, ok := options.Context.Value(deliveryMode{}).(uint8); ok {
			m.DeliveryMode = value
		}
	}
	return ch.channel.Publish(
		exchange.name,
		key,
		false,
		false,
		m,
	)
}

func (cm *ChannelManager) newExchange(name string) (*exchange, error) {
	ch, err := cm.obtainChannel()
	defer func() {
		cm.returnChannel(ch)
	}()

	if err != nil {
		return nil, err
	}

	if err := ch.channel.ExchangeDeclare(name, string(topic), true, false, false, false, nil); err != nil {
		return nil, err
	}
	exchange := &exchange{
		name: name,
	}
	cm.exchangeManager.exchanges[name] = exchange
	return exchange, nil
}

func (cm *ChannelManager) obtainExchange(name string) *exchange {
	if exchange, ok := cm.exchangeManager.exchanges[name]; ok {
		return exchange
	} else {
		exchange, _ := cm.newExchange(name)
		return exchange
	}
}

func (cm *ChannelManager) obtainSubscriber(exchangeName string, handler brokers.MessageHandler, opts brokers.SubscribeOptions) (s *subscriber, err error) {
	if opts.Context != nil {
		if value, ok := opts.Context.Value(queueName{}).(string); ok {
			opts.Queue = value
		}
	}
	exchange := cm.obtainExchange(exchangeName)
	ch, err := cm.obtainChannel()

	return &subscriber{channelManager: cm, channel: ch, handler: handler, opts: opts, exchange: exchange}, nil
}

func (cm *ChannelManager) subscribe(exchangeName string, handler brokers.MessageHandler, opts brokers.SubscribeOptions) error {
	sub, err := cm.obtainSubscriber(exchangeName, handler, opts)
	if err != nil {
		return err
	}

	return sub.run()
}
