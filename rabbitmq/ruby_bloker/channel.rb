module Rabbitmq
  class Channel
    @@publish_channel = nil
    @@suscribe_channel = nil

    def self.get_publish_channel
      if @@publish_channel.nil?
        connection = Rabbitmq::Connection.get_connection
        @@publish_channel = connection.create_channel
        @@publish_channel.prefetch(1) # worker每次只处理1条消息，在上一条消息没有ack前不会发送其他消息
        return @@publish_channel
      elsif @@publish_channel.status != :open
        if @@publish_channel.open
          return @@publish_channel
        else
          connection = Rabbitmq::Connection.get_connection
          @@publish_channel = connection.create_channel
          @@publish_channel.prefetch(1) # worker每次只处理1条消息，在上一条消息没有ack前不会发送其他消息
          return @@publish_channel
        end
      else
        @@publish_channel
      end
    end

    def self.get_suscribe_channel
      if @@suscribe_channel.nil?
        connection = Rabbitmq::Connection.get_connection
        @@suscribe_channel = connection.create_channel
        @@suscribe_channel.prefetch(1) # worker每次只处理1条消息，在上一条消息没有ack前不会发送其他消息
        return @@suscribe_channel
      elsif @@suscribe_channel.status != :open
        if @@suscribe_channel.open
          return @@suscribe_channel
        else
          connection = Rabbitmq::Connection.get_connection
          @@suscribe_channel = connection.create_channel
          @@suscribe_channel.prefetch(1) # worker每次只处理1条消息，在上一条消息没有ack前不会发送其他消息
          return @@suscribe_channel
        end
      else
        @@suscribe_channel
      end
    end
  end
end
