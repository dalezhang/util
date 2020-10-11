package rabbitmq

/*
我们这边把所有的 Exchange 的定义写死

| Durable | AutoDelete |                             说明                             	|
| :-----: | :--------: | :----------------------------------------------------------: 	|
|  True   |   False    |             重启之后Exchange会重新定义（默认的）             		|
|  True   |    True    | 当 durable 的队列需要绑定在 auto delete 的 exchange 的时候使用 	|
|  False  |   False    |          即使没有绑定也不会删除、但是重启不会重定义          		|
|  False  |    True    |             当没有绑定的时候会删除、重启不会保留             		|
*/
type exchangeType string

const (
	topic  exchangeType = "topic"
	fanout              = "fanout"
	direct              = "direct"
)

type exchange struct {
	name string
}

type exchangeManager struct {
	exchanges map[string]*exchange
}
