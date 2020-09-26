package brokers

import (
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

var (
	Conn     *amqp.Connection
	Channels []*amqp.Channel
	Dsn      string
)

type RabbitMQBroker struct {
	Channel          *amqp.Channel
	QueueName        string // 队列名称
	QueueConfig      QueueConfig
	RoutingKey       string // key名称
	ExchangeName     string // 交换机名称
	ExchangeType     string // 交换机类型 direct|fanout|topic|x-custom
	Messages         []Message
	ReceivedMessages []amqp.Delivery
	mu               sync.RWMutex
}

func Dial(dsn string) (conn *amqp.Connection, err error) {
	Dsn = dsn
	if Conn == nil || Conn.IsClosed() {
		if Conn.IsClosed() {
			Conn.Close()
			for _, c := range Channels {
				c.Close()
			}
			Channels = Channels[:0]
		}
		conn, err = amqp.Dial(dsn)
		if err != nil {
			fmt.Println("Dial connect err", err)
			return nil, err
		}
		Conn = conn
	}

	return Conn, err
}
func GetChannel() (mqChan *amqp.Channel, err error) {
	//TODO 根据状态选取channel
	Conn, err = Dial(Dsn)
	if err != nil {
		fmt.Println("Dial err================", err)
		return nil, err
	}

	fmt.Println("====len(Channels)", len(Channels))
	if len(Channels) == 0 {
		mqChan, err = Conn.Channel()
		if err != nil {
			fmt.Println("bulid channel err: ", err)
			return nil, err
		}
		Channels = []*amqp.Channel{mqChan}
	} else {
		mqChan = Channels[0]
	}
	fmt.Println("get channel", mqChan)
	return mqChan, err
}

func Close() error {
	// 先关闭管道,再关闭链接
	for _, c := range Channels {
		err := c.Close()
		if err != nil {
			fmt.Printf("MQ管道关闭失败:%s \n", err)
			return err
		}
	}

	err := Conn.Close()
	if err != nil {
		fmt.Printf("MQ链接关闭失败:%s \n", err)
		return err
	}
	return nil
}

// 定义队列交换机对象
type QueueExchange struct {
	QuName string // 队列名称
	RtKey  string // key值
	ExName string // 交换机名称
	ExType string // 交换机类型
}

// 创建一个新的操作对象
func NewRabbitmq(mqChan *amqp.Channel, q *QueueExchange) (b *RabbitMQBroker, err error) {
	qc := NewQueueConfig(func(q *QueueConfig) {
		q.NoWait = true
	})
	if err != nil {
		fmt.Printf("MQ获取channel失败:%s \n", err)
		return nil, err
	}
	return &RabbitMQBroker{
		Channel:      mqChan,
		QueueName:    q.QuName,
		QueueConfig:  qc,
		RoutingKey:   q.RtKey,
		ExchangeName: q.ExName,
		ExchangeType: q.ExType,
	}, nil
}
func (b *RabbitMQBroker) Close() {
	b.Channel.Close()
}

func (rb *RabbitMQBroker) AddMessages(msg string) error {
	m := new(RabbitmqMessage)
	m.body = []byte(msg)
	rb.Messages = append(rb.Messages, m)
	return nil
}

func (rb *RabbitMQBroker) Publish() error {

	qc := rb.QueueConfig
	_, err := rb.Channel.QueueDeclarePassive(rb.QueueName, qc.Durable, qc.AutoDelete, qc.Exclusive, qc.NoWait, nil)
	if err != nil {
		fmt.Printf("MQ注册队列失败:%s \n", err)
		return err
	}
	// 队列绑定
	err = rb.Channel.QueueBind(rb.QueueName, rb.RoutingKey, rb.ExchangeName, qc.NoWait, nil)
	if err != nil {
		fmt.Printf("MQ绑定队列失败:%s \n", err)
		return err
	}
	// 用于检查交换机是否存在,已经存在不需要重复声明
	err = rb.Channel.ExchangeDeclarePassive(
		rb.ExchangeName, // name
		rb.ExchangeType, // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // noWait
		nil,             // arguments
	)

	if err != nil {
		fmt.Printf("====MQ注册交换机失败1:%s, rb.ExchangeName: %s, rb.ExchangeType:%s \n", err, rb.ExchangeName, rb.ExchangeType)
		// 注册交换机
		// name:交换机名称,kind:交换机类型,durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;
		// noWait:是否非阻塞, true为是,不等待RMQ返回信息;args:参数,传nil即可; internal:是否为内部
		err = rb.Channel.ExchangeDeclare(
			rb.ExchangeName, // name
			rb.ExchangeType, // type
			true,            // durable
			false,           // auto-deleted
			false,           // internal
			false,           // noWait
			nil,             // arguments
		)
		if err != nil {
			fmt.Printf("======MQ注册交换机失败2: %s, rb.ExchangeName: %s, rb.ExchangeType:%s \n", err, rb.ExchangeName, rb.ExchangeType)
			return err
		}
	}

	failureMessages := rb.Messages[:0]
	for _, m := range rb.Messages {
		if err = rb.Channel.Publish(
			rb.ExchangeName, // publish to an exchange
			rb.RoutingKey,   // routing to 0 or more queues
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "text/plain",
				ContentEncoding: "",
				Body:            []byte(m.MsgContent()),
				DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
				Priority:        0,              // 0-9
				// a bunch of application/implementation-specific fields
			},
		); err != nil {
			failureMessages = append(failureMessages, m)
			fmt.Errorf("Exchange Publish: %s", err)
		}

	}
	rb.Messages = failureMessages

	return nil
}

func (rb *RabbitMQBroker) Subscribe() (err error) {
	// 用于检查交换机是否存在,已经存在不需要重复声明
	err = rb.Channel.ExchangeDeclarePassive(
		rb.ExchangeName, // name
		rb.ExchangeType, // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // noWait
		nil,             // arguments
	)
	if err != nil {
		// 注册交换机
		// name:交换机名称,kind:交换机类型,durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;
		// noWait:是否非阻塞, true为是,不等待RMQ返回信息;args:参数,传nil即可; internal:是否为内部
		err = rb.Channel.ExchangeDeclare(
			rb.ExchangeName, // name
			rb.ExchangeType, // type
			true,            // durable
			false,           // auto-deleted
			false,           // internal
			false,           // noWait
			nil,             // arguments
		)
		if err != nil {
			fmt.Printf("MQ注册交换机失败:%s \n", err)
			return err
		}
	}

	qc := rb.QueueConfig
	_, err = rb.Channel.QueueDeclare(rb.QueueName, qc.Durable, qc.AutoDelete, qc.Exclusive, qc.NoWait, nil)
	if err != nil {
		fmt.Printf("MQ注册队列失败:%s \n", err)
		return err
	}
	// 队列绑定
	err = rb.Channel.QueueBind(rb.QueueName, rb.RoutingKey, rb.ExchangeName, true, nil)
	if err != nil {
		fmt.Printf("MQ绑定队列失败:%s \n", err)
		return err
	}

	msgs, err := rb.Channel.Consume(rb.QueueName, "", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err, "Failed to register a consumer")
	}

	go func() {
		for msg := range msgs {
			rb.ReceivedMessages = append(rb.ReceivedMessages, msg)
			log.Printf("================Received a message: %s", msg.Body)
		}
	}()

	return nil
}
