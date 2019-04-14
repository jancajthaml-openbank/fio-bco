require 'json'
require 'thread'
require_relative '../shims/harden_webrick'
require_relative './vault_mock'

class LedgerTransactionHandler < WEBrick::HTTPServlet::AbstractServlet

  def do_POST(request, response)
    status, content_type, body = create_transaction(request)

    response.status = status
    response['Content-Type'] = content_type
    response.body = body
  end

  def create_transaction(request)
    begin
      body = JSON.parse(request.body)

      raise JSON::ParserError if body["transfers"].nil? || body["transfers"].empty?

      if VaultMock.create_transaction(body["id"], body["transfers"])
        return 200, "application/json", "{}"
      else
        return 409, "application/json", "{}"
      end
    rescue JSON::ParserError
      return 400, "application/json", "{}"
    rescue Exception => _
      return 500, "application/json", "{}"
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
      raise "Failed to allocate server binding! #{err}"
    end

    self.server.mount "/transaction", LedgerTransactionHandler

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
