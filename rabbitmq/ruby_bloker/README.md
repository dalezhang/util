# 发布者
## 目标
* 做一个topic模式的exchange
* 每次放入1条消息（厕所模式）
* 当项目rabbitmq clash或背kill后，重启后消息不能丢失
## 手段
* 创建topic模式的exchange
```ruby
    exchange = channel.topic("durable")
```
* channel prefetch 设为1
```ruby
    channel.prefetch(1) # worker每次只处理1条消息，在上一条消息没有ack前不会发送其他消息
```
* persistent 让消息入磁盘
```ruby
 presistant = message_config[:presistant] || true 
```
## 使用方法
```ruby
      # exchange
      exchange_config = {
          exchange_name: "test_topic_exchange",
          durable: true,
          auto_delete: false,
          internal: false
      }
      p = Rabbitmq::Publisher.new(exchange_config)

      # messsage
      message_config = {
          routing_key: "dale.test",
          content_type: "application/text",
          content_encoding: "utf8",
          mandatory: false, # false：找不到订阅者不报错， true：找不到订阅者报错
          reply_to: "icm-scripts-scheduler",
          correlation_id: "",
          app_id: "icm-scripts-scheduler"
      }
      (0..10).each do |i|
        message_config[:message_id] = "test_#{i}"
        p.publish("test message [#{i}]", message_config)
      end
```

# 订阅者
## 目标
* 做一个topic模式的exchange
* 找到创建一个queue用于订阅消息
* 当项目rabbitmq clash或背kill后，queue不能消失
* 手动ack
## 手段
* 创建topic模式的exchange
```ruby
    exchange = channel.topic("durable")
```
* 找到创建一个queue用于订阅消息 & 当项目rabbitmq clash或背kill后，queue不能消失
```ruby
    queue = channel.queue('', durable: ture)
```
* 手动ack
```ruby
    queue.subscribe(manual_ack: true, block: true) do |delivery_info, properties, body|
      puts " [x] #{delivery_info.routing_key} : #{body}"
      yield delivery_info, properties, body
      channel.ack(delivery_info.delivery_tag)
    end
```
## 使用方法
```ruby
    s = Rabbitmq::Subscriber.new(["#"],{}, {exchange_name: "test_topic_exchange"})
    s.run do |delivery_info, properties, body|
      puts "body: #{body}\n properties: #{properties} \n delivery_info: #{delivery_info}"
    end
```
