module Rabbitmq
  class Channel
    @@publish_channel = nil

    def self.get_publish_channel
        connection = Rabbitmq::Connection.get_connection
        @@publish_channel = connection.create_channel
        return @@publish_channel
      elsif @@publish_channel.status != :open
        if @@publish_channel.open
          return @@publish_channel
        else
          connection = Rabbitmq::Connection.get_connection
          @@publish_channel = connection.create_channel
          return @@publish_channel
        end
      else
        @@publish_channel
      end
    end

    def self.get_suscribe_channel
        connection = Rabbitmq::Connection.get_connection
        channel = connection.create_channel
        channel.prefetch(1) # worker每次只处理1条消息，在上一条消息没有ack前不会发送其他消息
        return channel
    end
  end
end
