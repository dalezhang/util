#!/usr/bin/env ruby

require 'bunny'

connection = Bunny.new
connection.start

channel  = connection.create_channel
exchange = channel.topic('test_topic_exchange', {
  durable: true,
  auto_delete: false,
  internal: false,
})
queue    = channel.queue('',   durable: false, auto_delete: true)

ARGV.each do |severity|
  puts severity
  queue.bind(exchange, routing_key: severity)
end

puts ' [*] Waiting for logs. To exit press CTRL+C'

begin
  queue.subscribe(block: true) do |delivery_info, _properties, body|
    puts " [x] #{delivery_info.routing_key}:#{body}"
  end
rescue Interrupt => _
  channel.close
  connection.close

  exit(0)
end

# 接收 第一个单词是 kern。后面至少有任意一个 单词
# ruby internal/brokers/rabbitmq/receive_demo.rb "kern.*"