package rabbitmq

type QueueConfig struct {
	Name       string // 队列名
	Durable    bool   // 消息代理重启后，队列依旧存在
	Exclusive  bool   //（只被一个连接（connection）使用，而且当连接关闭后队列即被删除）
	AutoDelete bool   //（当最后一个消费者退订后即被删除）
	Arguments  string //（一些消息代理用他来完成类似与 TTL 的某些额外功能）
	NoWait     bool
}
type QueueConfigOption func(qc *QueueConfig)

func NewQueueConfig(qcs ...func(q *QueueConfig)) (qc QueueConfig) {
	defaultQC := QueueConfig{
		Name:       "",
		Durable:    true,
		Exclusive:  false,
		AutoDelete: false,
		NoWait:     false,
	}
	for _, f := range qcs {
		f(&defaultQC)
	}
	return defaultQC
}

/*
 usage:
	qc = NewQueueConfig(func(q *QueueConfig) {
 		q.Durable = false
	}
*/
