class Biz::MarkLogger
  def initialize(type = "request", json_data = "")
    if json_data.present?
      # 使用异步队列时一定要传json_data，不然peform会把Biz::MarkLogger实例变成字符串，导致报错！！！
      hash = JSON.parse json_data
      @mark_arr = hash["mark_arr"]
      @uuid = hash["uuid"]
      @type = type["type"]
    else
      @mark_arr = []
      @uuid = UUID.new.generate
      @type = type
    end
  end

  def mark(append_str)
    @mark_arr << append_str.to_s.strip
  end

  def format_log(message)
    hash = {
      uuid: @uuid,
      type: @type,
      mark_arr: @mark_arr,
      message: message,
    }
    "\e[m\n#{JSON.pretty_generate hash}"
  end

  def log_info(message)
    Rails.mark_logger.info(format_log(message))
    mark(message)
  end

  def error(error)
    new_message = ""
    if error.message =~ /#{@uuid}/
      new_message = error.message
    else
      new_message = format_log(error.message)
    end
    new_error = StandardError.new(new_message)
    new_error.set_backtrace(error.backtrace)
    new_error
  end
end
