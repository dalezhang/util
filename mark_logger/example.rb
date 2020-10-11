module Rabbitmq
  class Handler
    attr_accessor :body

    def initialize(body, logger = nil)
      logger ||= MarkLogger.new
      @logger = logger
      @logger.mark("Rabbitmq::Handler#initialize")
      @body = body
    end

    def handle
      @logger.mark("Rabbitmq::Handler#handle")
      @logger.mark("parse: #{parse}")
      obj = parse[:object]
      raise "object is not specify" unless obj.present?
      action = parse[:action]
      raise "action is not present" unless action.present?
      opts = parse[:opts]
      raise "opts is not present" unless opts.present?
      db_path = opts[:db_path]
      raise "db_path is not present" unless db_path.present?
      @logger.mark("klass: #{klass.to_s}")
      klass.new(parse)
    rescue => e
      @logger.log_info(e.message)
      raise @logger.error(e)
    end

    private

    def klass
      "icm::_#{parse[:action]}_#{parse[:object]}_service".camelize.constantize
    end

    def parse
      @parse ||= JSON.parse(@body).symbolize_keys
    end

  end
end
