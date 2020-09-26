module Rabbitmq
  class Channel
    @@chennel = nil

    def self.get_channel
      if @@chennel.nil?
        connection = Rabbitmq::Connection.get_connection
        @@channel = connection.create_channel
        @@channel.prefetch(1) # worker每次只处理1条消息，在上一条消息没有ack前不会发送其他消息
        return @@channel
      elsif @@channel.status != :open
        if @@channel.open
          return @@channel
        else
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
