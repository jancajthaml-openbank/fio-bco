require 'json'
require 'date'
require 'json-schema'
require 'thread'
require_relative '../shims/harden_webrick'
require_relative './ledger_mock'

class LedgerTransactionsPartial < WEBrick::HTTPServlet::AbstractServlet

  def do_GET(request, response)
    status, body = process(request)

    response.status = status
    response.body = body
  end

  def process(request)
    path = request.path.split("/").map(&:strip).reject(&:empty?)

    case path.length
    when 2
      meta = LedgerMock.get_transactions(path[1])
      return 404, "{}" if meta.nil?
      return 200, meta.to_json
    when 3
      meta = LedgerMock.get_transaction(path[1], path[2])
      return 404, "{}" if meta.empty?
      return 200, meta.to_json
    else
      return 404, "{}"
    end
  end
end

module LedgerHelper

  def self.start
    self.server = nil

    begin
      self.server = WEBrick::HTTPServer.new(
        Port: 4401,
        Logger: WEBrick::Log.new("/dev/null"),
        AccessLog: [],
        SSLEnable: true
      )

    rescue Exception => err
      raise err
      raise "Failed to allocate server binding! #{err}"
    end

    self.server.mount "/transaction", LedgerTransactionsPartial

    self.server_daemon = Thread.new do
      self.server.start()
    end
  end

  def self.stop
    self.server.shutdown() unless self.server.nil?
    begin
      self.server_daemon.join() unless self.server_daemon.nil?
    rescue
    ensure
      self.server_daemon = nil
      self.server = nil
    end
  end

  class << self
    attr_accessor :server_daemon, :server
  end

end
