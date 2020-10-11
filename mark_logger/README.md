#### 在config/initializers里添加文件
```
# logger.rb
require 'logger'
module Rails
  def self.mark_logger
    @mark_logger ||= Logger.new Rails.root.join('log', 'mark_logger.log')
  end
end

```
#### 在引入的类里加上
```
class A
    def initialize( logger = nil)
      logger ||= MarkLogger.new
      @logger = logger
      @logger.mark("A#initialize")
    end
end
```

#### 当你需要在日志链里加入临时记录时
```
@logger.mark("message: #{message}")
```

#### 当你需要保存日志时
```
 @logger.log_info("message")
```

#### 当你需要在异常信息里带上mark_logger时
```
def abc
...
rescue => e
   raise @logger.error(e)
end

```
#### 当你需要在吧mark_logger传递给异步任务时
```
json = @logger.to_json
class A
    def initialize( logger = nil)
      logger ||= MarkLogger.new("", logger)
      @logger = logger
      @logger.mark("A#initialize")
    end
end

```