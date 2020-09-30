module Rabbitmq
  class Channel
    @@channel = nil

    def self.get_channel
      if @@channel.nil?
        connection = Rabbitmq::Connection.get_connection
        @@channel = connection.create_channel
        @@channel.prefetch(1) # worker每次只处理1条消息，在上一条消息没有ack前不会发送其他消息
        return @@channel
      elsif @@channel.status != :open
        if @@channel.open
          return @@channel
        else
          connection = Rabbitmq::Connection.get_connection
          @@channel = connection.create_channel
          @@channel.prefetch(1) # worker每次只处理1条消息，在上一条消息没有ack前不会发送其他消息
          return @@channel
        end
      else
        @@channel
      end
    end
  end
end
