package rabbitmq

import (
	"sync"

	"github.com/embrace-century/basin-warden/internal/brokers"
	"github.com/streadway/amqp"
)

type RabbitBroker struct {
	conn           *amqp.Connection
	channelManager *ChannelManager // 用来管理信道的
	// publishLock 是为了保证调用  Publish(exchangeName string, message *brokers.Message, opts ...brokers.PublishOption) error 方法的顺序
	// 和publish到rabbitmq的顺序保持一致
	publishLock sync.Mutex
}

func Dial(dsn string) (brokers.Broker, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	rabbitMQBroker := &RabbitBroker{
		conn: conn,
		channelManager: &ChannelManager{
			conn:     conn,
			channels: make(chan *rabbitmqChannel, 5),
			exchangeManager: &exchangeManager{
				exchanges: make(map[string]*exchange),
			},
		},
	}
	return rabbitMQBroker, nil
}

func (rb *RabbitBroker) exchangeManager(em *exchangeManager) {
	em = rb.channelManager.exchangeManager
}

func (rb *RabbitBroker) Close() {
	for {
		ch, ok := <-rb.channelManager.channels
		if !ok {
			break
		}
		ch.close()
	}
	_ = rb.conn.Close()
}

func (rb *RabbitBroker) Publish(exchangeName string, message *brokers.Message, opts ...brokers.PublishOption) error {
	exchange := rb.channelManager.obtainExchange(exchangeName)
	return rb.channelManager.publish(exchange, message, opts...)
}

func (rb *RabbitBroker) Subscribe(exchangeName string, handler brokers.MessageHandler, opts ...brokers.SubscribeOption) error {
	options := brokers.SubscribeOptions{}
	for _, oFun := range opts {
		oFun(&options)
	}

	return rb.channelManager.subscribe(exchangeName, handler, options)
}
