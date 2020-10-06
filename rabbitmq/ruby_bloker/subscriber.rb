class Rabbitmq::Subscriber
  def initialize(route_keys = [], queue_config = {}, exchange_config = {})
    @route_keys =route_keys
    @queue_config = queue_config
    @exchange_config = exchange_config
  end

  def run(&block)
    channel = Rabbitmq::Channel.get_suscribe_channel
    exchange_name = @exchange_config[:exchange_name] || "test_topic_exchange"
    exchange_durable = @exchange_config[:durable].nil? ? true : @exchange_config[:durable]
    exchange_auto_delete = @exchange_config[:auto_delete].nil? ? false : @exchange_config[:auto_delete]
    exchange_internal = @exchange_config[:internal].nil? ? false : @exchange_config[:internal]
    exchange = channel.topic(exchange_name, {
        durable: exchange_durable,
        auto_delete: exchange_auto_delete,
        internal: exchange_internal,
    })

    #Durable queues that are shared by many consumers and have an independent existence: i.e.
    # they will continue to exist and collect messages whether or not there are consumers to receive them.
    queue_durable = @queue_config[:durable].nil? ? true : @queue_config[:durable]
    queue_auto_delete = @queue_config[:auto_delete].nil? ? false : @queue_config[:auto_delete] # Whether the queue is auto-deleted when no longer used
    queue = channel.queue('',
                          durable: queue_durable,
                          auto_delete: queue_auto_delete
    )
    @route_keys.each do |route_name|
      queue.bind(exchange, routing_key: route_name)
    end
    queue.subscribe(manual_ack: true, block: true) do |delivery_info, properties, body|
      puts " [x] #{delivery_info.routing_key} : #{body}"
      yield delivery_info, properties, body
      channel.ack(delivery_info.delivery_tag)
    end
  end

end