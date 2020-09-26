require "bunny"
module Rabbitmq
  class Connection
    @@connection = nil

    def self.get_connection
      if @@connection.nil?
        # To disable automatic connection recovery, pass :automatic_recovery => false to Bunny.new.
        @@connection = ::Bunny.new(self.connection_opts)
        @@connection.start
      elsif @@connection.status != :open
        if @@connection.start
          return @@connection
        else
          raise "connect failure"
        end
      else
        @@connection
      end
    end

    private

    def self.connection_opts
      @@connection_opts ||= {
          :host => "127.0.0.1",
          :port => 5672,
          :ssl => false,
          :vhost => "/",
          :user => "guest",
          :pass => "guest",
          :heartbeat => :server, # will use RabbitMQ setting
          :frame_max => 131072,
          :auth_mechanism => "PLAIN"
      }
    end
  end
end