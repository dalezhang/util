# 发布者
## 目标
* 做一个topic模式的exchange
* 每次放入1条消息（厕所模式）
* 当项目rabbitmq clash或背kill后，重启后消息不能丢失
## 手段
* 创建topic模式的exchange
```
    exchange = channel.topic("durable")
```
* channel prefetch 设为1
```
    channel.prefetch(1) # worker每次只处理1条消息，在上一条消息没有ack前不会发送其他消息
```
* persistent 让消息入磁盘
```
 presistant = message_config[:presistant] || true 
```
## 使用方法
```
      p = Rabbitmq::Publisher.new("1111111111111111111",
                                  routing_key: "dale.test",
                                  content_type: "text",
                                  content_encoding: "utf8",
                                  message_id: "test_1",
                                  reply_to: "",
                                  correlation_id: "",
                                  app_id: nil
      )
      p.run
```

# 订阅者
## 目标
* 做一个topic模式的exchange
* 找到创建一个queue用于订阅消息
* 当项目rabbitmq clash或背kill后，queue不能消失
* 手动ack
## 手段
* 创建topic模式的exchange
```
    exchange = channel.topic("durable")
```
* 找到创建一个queue用于订阅消息
```
    queue = channel.queue('', durable: ture)
```
* 当项目rabbitmq clash或背kill后，queue不能消失
## 使用方法
```
    s = Rabbitmq::Subscriber.new(["#"])
    s.run do |delivery_info, properties, body|
      puts "body: #{body}\n properties: #{properties} \n delivery_info: #{delivery_info}"
    end
```
