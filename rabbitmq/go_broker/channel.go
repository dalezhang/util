package rabbitmq

import (
	"github.com/streadway/amqp"
)

var (
	defaultRoutingKey = "embrace-century-basin.warden.default-routing-key"
)

type rabbitmqChannel struct {
	uuid    string
	channel *amqp.Channel
	//lock    mutx.Lock
}

func (rc *rabbitmqChannel) close() {
	rc.channel.Close()
}
