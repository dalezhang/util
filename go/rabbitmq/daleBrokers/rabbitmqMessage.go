package brokers

import (
	"fmt"
)

type RabbitmqMessage struct {
	body []byte
}

func (m *RabbitmqMessage) MsgContent() string {
	return fmt.Sprintf("%s", m.body)
}
