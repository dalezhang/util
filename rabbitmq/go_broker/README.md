## 连接rabbitmq
```
conn, err := Dial(*uri)
```
* 连接成功后这个连接回放入全局的Conn，再次调用回先去Conn中取，如果Conn可用则直接返回，不可用则再次创建连接
* 如果原连接不可用，会关闭conn和它的channel，然后删除

## 获取channel
```
ch, err := GetChannel()
```
* 从Conn中获取一个channel
## 创建Rabbitmq Broker
#### 配置队列参数
```
type QueueExchange struct {
	QuName string // 队列名称
	RtKey  string // key值
	ExName string // 交换机名称
	ExType string // 交换机类型
}
queueExchange := &QueueExchange{
    "test",
    "routeone",
    "logs_topic",
    amqp.ExchangeTopic,
}
```
#### 创建broker
```
productBroker, err := NewRabbitmq(ch, queueExchange)
```

```
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
```