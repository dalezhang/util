work_q = Queue.new
workers = (0...10).map do
    Thread.new do
    begin
      while page = work_q.pop(true)      
        begin
          puts 1
        rescue => e
          puts "problem on page #{page}"
          puts e.inspect
        end
      end # while
      puts ""
    rescue ThreadError
    end
  end
end
work_q << 1

class A
  include B
  @@a = 1
  def self.a
    @@a += 1
    puts @@a

  end
end
a = A.new
a.b


