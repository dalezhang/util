gem "minitest"

require "minitest/pride"
require "minitest/autorun"

$VERBOSE = 1

require_relative "../rabbitmq/channel.rb"
require_relative "../rabbitmq/publisher.rb"
require_relative "../rabbitmq/connection.rb"
require_relative "../rabbitmq/subscriber.rb"

module Rabbitmq
  class TestRabbitmq < Minitest::Test
    def test_publish
      # exchange
      exchange_config = {
          exchange_name: "test_topic_exchange",
          durable: true,
          auto_delete: false,
          internal: false
      }
      p = Rabbitmq::Publisher.new(exchange_config)

      # message
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
      # assert_equal((1..generations).cycle(pool_size).sort, result.sort)
      #
      # assert_operator(finish - start, :>, generations * NetworkConnection::SLEEP_TIME)

     # test_subscribe
      s = Rabbitmq::Subscriber.new(["#"],{}, {exchange_name: "test_topic_exchange"})
      s.run do |delivery_info, properties, body|
        puts "body: #{body}\n properties: #{properties} \n delivery_info: #{delivery_info}"
      end
    end
  end

end

