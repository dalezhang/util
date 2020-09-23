package brokers

type Broker interface {
	AddMessages(string) error
	Close() error
	Publish() error
	Subscribe() error
}

type Message interface {
	MsgContent() string
}

type ExchanteConfig struct {
}

type ReceivedMessage interface {
}
