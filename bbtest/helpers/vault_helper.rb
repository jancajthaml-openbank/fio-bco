require 'json'
require 'date'
require 'json-schema'
require 'thread'
require_relative '../shims/harden_webrick'
require_relative './vault_mock'

class VaultAccountPartial < WEBrick::HTTPServlet::AbstractServlet

  def do_GET(request, response)
    status, body = process_get(request)

    response.status = status
    response.body = body
  end

  def do_POST(request, response)
    status, content_type, body = process_post(request)

    response.status = status
    response.body = body
  end

  def process_post(request)
    path = request.path.split("/").map(&:strip).reject(&:empty?)

    case path.length
    when 3
      begin
        body = JSON.parse(request.body)

        raise JSON::ParserError if body["name"].nil? || body["name"].empty?
        raise JSON::ParserError if body["format"].nil? || body["format"].empty?
        raise JSON::ParserError if body["currency"].nil? || body["currency"].empty?
        raise JSON::ParserError if body["isBalanceCheck"].nil?

        if VaultMock.create_account(path[1], body["name"], body["format"], body["currency"], body["isBalanceCheck"] != "false")
          return 200, "{}"
        else
          return 409, "{}"
        end
      rescue JSON::ParserError
        return 400, "{}"
      rescue Exception => err
        puts err
        return 500, "{}"
      end
    else
      return 404, "{}"
    end
  end

  def process_get(request)
    path = request.path.split("/").map(&:strip).reject(&:empty?)

    case path.length
    when 2
      meta = VaultMock.get_accounts(path[1])
      return 404, "{}" if meta.nil?
      return 200, meta.to_json
    when 3
      meta = VaultMock.get_acount(path[1], path[2])
      return 404, "{}" if meta.empty?
      return 200, meta.to_json
    else
      return 404, "{}"
    end
  end
end

module VaultHelper

  def self.start
    self.server = nil

    begin
      self.server = WEBrick::HTTPServer.new(
        Port: 4400,
        Logger: WEBrick::Log.new("/dev/null"),
        AccessLog: [],
        SSLEnable: true
      )

    rescue Exception => err
      raise err
      raise "Failed to allocate server binding! #{err}"
    end

    self.server.mount "/account", VaultAccountPartial

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
