module EventuallyHelper

  def self.eventually(timeout: 10, delay: 2, &_blk)
    wait_until = Time.now + timeout
    begin
      yield
    rescue => e
      raise e if Time.now >= wait_until
      sleep delay
      retry
    end
  end

  def eventually(*args, &blk)
    EventuallyHelper.eventually(*args, &blk)
  end
end
