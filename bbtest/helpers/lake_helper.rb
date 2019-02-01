require 'ffi-rzmq'
require 'thread'
require 'timeout'

# FIXME problem is within pake because it does not actually relay message from rest to unit

class LakeMessage

  attr_accessor :tenantId, :raw

  def initialize(tenantId)
    self.tenantId = tenantId
    raise "raw LakeMessage cannot be initialized"
  end

  def ===(other)
    if defined? other.tenantId
      unless self.tenantId.nil? or other.tenantId.nil?
        return false unless self.tenantId? == other.tenantId
      end

      other.raw === self.raw
    else
      other === self.raw
    end
  end

  def to_s
   self.raw
  end

end

class LakeMessageError < LakeMessage

  def initialize(tenantId, reqId)
    self.tenantId = tenantId
    self.raw = "Error[#{reqId}]"
  end

end

class LakeMessageTokenCreated < LakeMessage

  def initialize(tenantId, reqId)
    self.tenantId = tenantId
    self.raw = "TokenCreated[#{reqId}]"
  end

end

class LakeMessageTokenDeleted < LakeMessage

  def initialize(tenantId, reqId)
    self.tenantId = tenantId
    self.raw = "TokenDeleted[#{reqId}]"
  end

end

module LakeMock

  def self.parse_message(msg)

    if groups = msg.match(/^Wall\/bbtest FioUnit\/([^\s]{1,100}) ([^\s]{1,100}) token TN$/i)
      tenantId, reqId = groups.captures
      return LakeMessageTokenCreated.new(tenantId, reqId)

    elsif groups = msg.match(/^([^\s]{1,100}) token TN$/i)
      reqId, _ = groups.captures
      return LakeMessageTokenCreated.new(nil, reqId)

    elsif groups = msg.match(/^Wall\/bbtest FioUnit\/([^\s]{1,100}) ([^\s]{1,100}) token TD$/i)
      tenantId, reqId = groups.captures
      return LakeMessageTokenDeleted.new(tenantId, reqId)

    elsif groups = msg.match(/^([^\s]{1,100}) token TD$/i)
      reqId, _ = groups.captures
      return LakeMessageTokenDeleted.new(nil, reqId)

    elsif groups = msg.match(/^Wall\/bbtest FioUnit\/([^\s]{1,100}) ([^\s]{1,100}) token EE$/i)
      tenantId, reqId = groups.captures
      return LakeMessageError.new(tenantId, reqId)

    elsif groups = msg.match(/^([^\s]{1,100}) token EE$/i)
      reqId, _ = groups.captures
      return LakeMessageError.new(nil, reqId)

    else
      raise "lake unknown event \"#{msg}\""
    end
  end

  def self.start
    raise "cannot start when shutting down" if self.poisonPill
    self.poisonPill = false

    begin
      ctx = ZMQ::Context.new
      pull_channel = ctx.socket(ZMQ::PULL)
      raise "unable to bind PULL" unless pull_channel.bind("tcp://*:5562") >= 0
      pub_channel = ctx.socket(ZMQ::PUB)
      raise "unable to bind PUB" unless pub_channel.bind("tcp://*:5561") >= 0
    rescue ContextError => _
      raise "Failed to allocate context or socket!"
    end

    self.ctx = ctx
    self.pull_channel = pull_channel
    self.pub_channel = pub_channel

    self.pull_daemon = Thread.new do
      loop do
        break if self.poisonPill or self.pull_channel.nil?
        data = ""
        begin
          Timeout.timeout(1) do
            self.pull_channel.recv_string(data, 0)
          end
        rescue Timeout::Error => _
          break if self.poisonPill or self.pull_channel.nil?
          next
        end
        next if data.empty?

        if data.end_with?("]")
          self.pub_channel.send_string(data)
          self.pub_channel.send_string(data)
          next
        end

        unless data.start_with?("Wall/bbtest")
          self.send(data)
          next
        end
        self.mutex.synchronize do
          self.recv_backlog << data
        end
      end
    end
  end

  def self.stop
    self.poisonPill = true
    begin
      self.pull_daemon.join() unless self.pull_daemon.nil?
      self.pub_channel.close() unless self.pub_channel.nil?
      self.pull_channel.close() unless self.pull_channel.nil?
      self.ctx.terminate() unless self.ctx.nil?
    rescue
    ensure
      self.pull_daemon = nil
      self.ctx = nil
      self.pull_channel = nil
      self.pub_channel = nil
    end
    self.poisonPill = false
  end

  def ack(data)
    LakeMock.ack(data)
  end

  def parsed_mailbox()
    LakeMock.parsed_mailbox()
  end

  def mailbox()
    LakeMock.mailbox()
  end

  def send(data)
    LakeMock.send(data)
  end

  def pulled_message?(expected)
    LakeMock.pulled_message?(expected)
  end

  class << self
    attr_accessor :ctx,
                  :pull_channel,
                  :pub_channel,
                  :pull_daemon,
                  :mutex,
                  :recv_backlog,
                  :poisonPill
  end

  self.recv_backlog = []

  self.mutex = Mutex.new
  self.poisonPill = false


  def self.parsed_mailbox()
    return self.recv_backlog.map { |item| self.parse_message(item) }
  end

  def self.mailbox()
    return self.recv_backlog
  end

  def self.pulled_message?(expected)
    copy = self.recv_backlog.dup
    copy.each { |item|
      return true if self.parse_message(item) === expected
    }
    return false
  end

  def self.send(data)
    self.pub_channel.send_string(data) unless self.pub_channel.nil?
  end

  def self.ack(data)
    self.mutex.synchronize do
      self.recv_backlog.reject! { |v| self.parse_message(v) === data }
    end
  end

end
