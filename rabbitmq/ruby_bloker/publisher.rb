module Rabbitmq
  class Publisher

    def initialize(message, message_config = {})
      @message = message
      @message_config = message_config
    end

    def run
      channel = Rabbitmq::Channel.get_channel

      # create or find a exchange
      exchange = channel.topic("durable")

      # do message settings
      routing_key = @message_config[:routing_key] || "" # Used for routing messages depending on the exchange type and configuration.
      presistant = @message_config[:presistant] || true # When set to true, RabbitMQ will persist message to disk.

      # This flag tells the server how to react if the message cannot be routed to a queue.
      # If this flag is set to true, the server will return an unroutable message to the producer with a `basic.return` AMQP method.
      # If this flag is set to false, the server silently drops the message.
      mandatory = @message_config[:mandatory] || false
      content_type = @message_config[:content_type] # MIME content type of message payload. Has the same purpose/semantics as HTTP Content-Type header.
      content_encoding = @message_config[:content_encoding] # MIME content type of message payload. Has the same purpose/semantics as HTTP Content-Type header.

      # Message identifier as a string. If applications need to identify messages, it is recommended that they use this attribute instead of putting it into the message payload.
      message_id = @message_config[:message_id]

      #Commonly used to name a reply queue (or any other identifier that helps a consumer application to direct its response).
      # Applications are encouraged to use this attribute instead of putting this information into the message payload.
      reply_to = @message_config[:reply_to]

      # ID of the message that this message is a reply to. Applications are encouraged to use this attribute instead of putting this information into the message payload.
      correlation_id = @message_config[:correlation_id]
      app_id = @message_config[:app_id] || "icm-scripts-scheduler" # Application identifier string, for example, "eventoverse" or "webcrawler"
      puts " [x] #{@message}"
      exchange.publish(
          @message,
          routing_key: routing_key,
          persistent: presistant,
          mandatory: mandatory,
          content_type: content_type,
          content_encoding: content_encoding,
          message_id: message_id,
          reply_to: reply_to,
          correlation_id: correlation_id,
          app_id: app_id
      )
    end
  end
end
