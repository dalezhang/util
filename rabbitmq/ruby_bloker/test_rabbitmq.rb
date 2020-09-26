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
      s = Rabbitmq::Subscriber.new(["#"])
      s.run do |delivery_info, properties, body|
        puts "body: #{body}\n properties: #{properties} \n delivery_info: #{delivery_info}"
      end
      # assert_equal((1..generations).cycle(pool_size).sort, result.sort)
      #
      # assert_operator(finish - start, :>, generations * NetworkConnection::SLEEP_TIME)
    end
  end

end

