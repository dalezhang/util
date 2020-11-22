str = "{daag (dag )[ddd]} {ddd}"

class Test
  def initialize(str)
    @match = true
    @str = str
  end

  def slid_str(str)
    match1 = /\{/.match(str)
    match2 = /\[/.match(str)
    match3 = /\(/.match(str)
    if match1 == nil && match2 == nil && match3 == nil
      return ""
    end
    index1 = match1 ? match1.pre_match.size : nil
    index2 = match2 ? match2.pre_match.size : nil
    index3 = match3 ? match3.pre_match.size : nil
    arr = [index1, index2, index3] - [nil]
    case arr.sort.first
    when index1
      str = match1.post_match
      next_step_match(/\}/.match(str))
    when index2
      puts "index 2 ===="
      str = match2.post_match
      next_step_match(/\]/.match(str))
    when index3
      str = match3.post_match
      next_step_match(/\)/.match(str))
    end
    str
  end

  def validate_breaks
    @match = true
    slid_str(@str)
    @match
  end
  def next_step_match(match)
    if match == nil
      @match = false
      @str = ""
      return
    end
    slid_str(match.post_match)
    slid_str(match.pre_match)
  end
end
t = Test.new(str)
puts t.validate_breaks
