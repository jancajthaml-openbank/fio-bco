require 'json'
require 'thread'
require_relative '../shims/harden_webrick'
require_relative './vault_mock'

class WallAccountHandler < WEBrick::HTTPServlet::AbstractServlet

  def do_POST(request, response)
    status, content_type, body = create_account(request)

    response.status = status
    response['Content-Type'] = content_type
    response.body = body
  end

  def create_account(request)
    begin
      body = JSON.parse(request.body)

      raise JSON::ParserError if body["accountNumber"].nil? || body["accountNumber"].empty?
      raise JSON::ParserError if body["currency"].nil? || body["currency"].empty?
      raise JSON::ParserError if body["isBalanceCheck"].nil? || body["isBalanceCheck"].empty?

      if VaultMock.create_account(body["accountNumber"], body["currency"], body["isBalanceCheck"] != "false")
        puts "created account #{body}"
        return 200, "application/json", "{}"
      else
        #puts "already exists #{body}"
        return 409, "application/json", "{}"
      end
    rescue JSON::ParserError
      return 400, "application/json", "{}"
    rescue Exception => _
      return 500, "application/json", "{}"
    end

  end
end

class WallTransactionHandler < WEBrick::HTTPServlet::AbstractServlet

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
        puts "created transaction #{body}"
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

module WallHelper

  def self.start
    self.server = nil

    begin
      self.server = WEBrick::HTTPServer.new(
        Port: 3000,
        Logger: WEBrick::Log.new("/dev/null"),
        AccessLog: [],
        SSLEnable: true
      )
    rescue Exception => err
      raise "Failed to allocate server binding! #{err}"
    end

    self.server.mount "/account", WallAccountHandler
    self.server.mount "/transaction", WallTransactionHandler

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