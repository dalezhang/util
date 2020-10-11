

require 'bunny'

connection = Bunny.new
connection.start

channel  = connection.create_channel
exchange = channel.topic('test_topic_exchange', {
  durable: true,
  auto_delete: false,
  internal: false,
})
severity = ARGV.shift || "zeng.yi.chen"
message = ARGV.empty? ? 'Hello World!' : ARGV.join(' ')

exchange.publish(message, routing_key: severity)
puts " [x] Sent #{severity}:#{message}"

connection.close

# 接收 第一个单词是 kern。后面至少有任意一个 单词
# ruby internal/brokers/rabbitmq/receive_demo.rb "kern.*"