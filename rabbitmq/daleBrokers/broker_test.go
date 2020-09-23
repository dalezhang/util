package brokers

import (
	"flag"
	"fmt"
	"testing"

	"github.com/streadway/amqp"
)

var (
	uri = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
)

func init() {
	// fmt.Println("init========================")
	// // err := Dial(*uri)
	// err := Dial("amqp://guest:guest@localhost:5672/")
	// // defer Close()
	// if err != nil {
	// 	panic(err)
	// }
	// queueExchange := &QueueExchange{
	// 	"test",
	// 	"rabbit.routeone",
	// 	"test.exchangeone",
	// 	amqp.ExchangeTopic,
	// }
	// rbbroker, err = NewRabbitmq(queueExchange)
	// if err != nil {
	// 	panic(err)
	// }
}

type TestMessage struct {
	A string
}

func (t *TestMessage) MsgContent() string {
	fmt.Println(t.A)
	return t.A
}

type AAA interface {
	MsgContent() string
}

func BBB(a AAA) {
	fmt.Println("========BBB: ", a.MsgContent())
}

var AAAs []AAA

func AppendAAA(a AAA) {
	AAAs = append(AAAs, a)
}

func TestBrokerClose(t *testing.T) {
	m := new(TestMessage)
	m.A = "11111"
	BBB(m)
	AppendAAA(m)
	fmt.Println("====AAAs: ", AAAs)
}

func TestBrokerDial(t *testing.T) {
	fmt.Println("init========================")
	// err := Dial(*uri)
	conn, err := Dial(*uri)
	defer conn.Close()
	if err != nil {
		fmt.Println("Dial err================")
		panic(err)
	}
	queueExchange := &QueueExchange{
		"test",
		"routeone",
		"logs_topic",
		amqp.ExchangeTopic,
	}
	ch, err := GetChannel()
	defer ch.Close()
	// 发送者

	if err != nil {
		panic(err)
	}
	productBroker, err := NewRabbitmq(ch, queueExchange)
	if err != nil {
		panic(err)
	}
	// 接收者
	consumeBroker, err := NewRabbitmq(ch, queueExchange)
	if err != nil {
		panic(err)
	}

	productBroker.AddMessages("11111111111")
	productBroker.AddMessages("22222222222")
	productBroker.AddMessages("333333333333")

	err = productBroker.Publish()
	if err != nil {
		panic(err)
	}

	consumeBroker.Subscribe()
	if err != nil {
		panic(err)
	}
	for i, m := range consumeBroker.ReceivedMessages {
		fmt.Printf("ReceivedMessages %d, body: %s", i, string(m.Body))
	}

	// if result.Data.TopicableType != "Series" {
	// 	t.Errorf("result.Data: %+v \n", result.Data)
	// }
	// // "{\"timestamp\":\"2019-08-08 19:52:09\",\"type\":\"topic\",\"id\":8038}
	// if result.Encrypted.ID != 8038 {
	// 	t.Errorf("result.Encrypted: %+v \n", result.Encrypted)
	// }
}
