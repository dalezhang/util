class Rabbitmq::Subscriber
  def initialize(route_names = [], queue_config = {})
    @route_names = route_names
    @queue_config = queue_config
  end

  def run(&block)
    channel = Rabbitmq::Channel.get_channel
    exchange = channel.topic("durable")

    #Durable queues that are shared by many consumers and have an independent existence: i.e.
    # they will continue to exist and collect messages whether or not there are consumers to receive them.
    durable = @queue_config[:durable] || true
    auto_delete = @queue_config[:auto_delete] || false # Whether the queue is auto-deleted when no longer used
    queue = channel.queue('',
                          durable: durable,
                          auto_delete: auto_delete
    )
    @route_names.each do |route_name|
      queue.bind(exchange, routing_key: route_name)
    end
    queue.subscribe(manual_ack: true, block: true) do |delivery_info, properties, body|
      puts " [x] #{delivery_info.routing_key} : #{body}"
      yield delivery_info, properties, body
      channel.ack(delivery_info.delivery_tag)
    end
  end

end