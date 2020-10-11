package rabbitmq

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/embrace-century/basin-warden/internal/brokers"
	"github.com/sirupsen/logrus"
)

var (
	err    error
	broker brokers.Broker
	dsn    = "amqp://guest:guest@localhost:5672/"
)

func init() {
	broker, err = Dial(dsn)
	if err != nil {
		panic(err)
	}
}

type testMessageStruct struct {
	Number int
	lock   sync.Mutex
}

func TestPublish(t *testing.T) {

	routingOpts := RoutingKey("zeng.yi.chen")
	deliveryMode := DeliveryMode(1)

	var wg sync.WaitGroup
	m := testMessageStruct{Number: 0}
	wg.Add(100)
	for {
		m.lock.Lock()
		if m.Number < 100 {
			go func() {
				defer m.lock.Unlock()
				m.Number = m.Number + 1
				message := &brokers.Message{
					ContentType: "application/json",
					Body:        []byte(fmt.Sprintf("My Tao: [%d]", m.Number)),
				}
				fmt.Println("m.Number === ", m.Number)
				err = broker.Publish("test_topic_exchange", message, routingOpts, deliveryMode)
				if err != nil {
					logrus.Fatal(err)
				}

				wg.Done()
			}()

		} else {
			break
		}

	}
	wg.Wait()
}

type testResponseMessage struct {
	ContentType string
	Body        []byte
	MessageID   string
	Timestamp   time.Time
}

func (r testResponseMessage) Handler(message *brokers.ReceivedMessage) error {
	r.ContentType = message.ContentType
	r.Body = message.Body
	r.MessageID = message.MessageID
	r.Timestamp = message.Timestamp
	fmt.Printf("message: %s \nrespons_detail: %+v", string(r.Body), r)

	return nil
}
func TestSubscribe(t *testing.T) {
	qn := QueueName("")
	rkeys := RoutingKeys([]string{"#"})

	r := new(testResponseMessage)
	err := broker.Subscribe("test_topic_exchange", r, qn, rkeys)
	if err != nil {
		logrus.Fatal(err)
	}

	log.Printf(" [*] defer close")
	defer broker.Close()
	forever := make(chan bool)
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
