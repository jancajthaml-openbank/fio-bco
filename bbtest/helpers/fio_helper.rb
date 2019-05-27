require 'json'
require 'thread'
require_relative '../shims/harden_webrick'
require_relative './fio_mock'

# get card transactions
#https://www.fio.cz/ib_api/rest/merchant/xxx/2017-01-01/2018-01-01/transactions.xml

# get account transactions
#https://www.fio.cz/ib_api/rest/last/xxx/transactions.xml
#https://www.fio.cz/ib_api/rest/last/xxx/transactions.json
#https://www.fio.cz/ib_api/rest/last/xxx/transactions.csv
#https://www.fio.cz/ib_api/rest/last/xxx/transactions.oxf

class FioDownloadNewStatemetsHandler < WEBrick::HTTPServlet::AbstractServlet

  def do_GET(request, response)
    status, content_type, body = process(request)

    response.status = status
    response['Content-Type'] = content_type
    response.body = body
  end

  def process(request)
    params = request.path_info.split("/").map(&:strip).reject(&:empty?)

    return 404, "application/json", "{}" if params.length < 2

    token = params[0]

    if params[1].include? "transactions.json"
      #
    elsif params[1].include? "transactions.xml"
      #
    end

    statements = FioMock.get_statements()

    FioMock.set_confirmed_transfer_pivot_id(FioMock.get_last_transfer_id())

    return 200, "application/json", statements.to_json
  end
end

class FioSetLastStatemetPivotIdHandler < WEBrick::HTTPServlet::AbstractServlet

  def do_GET(request, response)
    status, content_type, body = process(request)

    response.status = status
    response['Content-Type'] = content_type
    response.body = body
  end

  def process(request)
    return 404, "application/json", "{}" unless request.request_uri.to_s.end_with?("/")

    params = request.path_info.split("/").map(&:strip).reject(&:empty?)

    return 404, "application/json", "{}" if params.length < 2

    token = params[0]
    transferId = params[1]

    FioMock.set_confirmed_transfer_pivot_id(transferId)

    if request.accept.include? "application/json"
      return 200, "application/json", "{}"
    else
      puts "unknown accept #{request.accept}"
      return 200, "application/json", "{}"
    end
  end
end

class FioSetLastStatemetPivotDateHandler < WEBrick::HTTPServlet::AbstractServlet

  def do_GET(request, response)
    # FIXME check if url ends with slash (be robust againts real fio gateway)
    status, content_type, body = process(request)

    response.status = status
    response['Content-Type'] = content_type
    response.body = body
  end

  def process(request)
    return 404, "application/json", "{}" unless request.request_uri.to_s.end_with?("/")

    params = request.path_info.split("/").map(&:strip).reject(&:empty?)

    return 404, "application/json", "{}" if params.length < 2

    token = params[0]
    date = params[1]

    if request.accept.include? "application/json"
      return 200, "application/json", "{}"
    else
      puts "unknown acept #{request.accept}"
      return 200, "application/json", "{}"
    end
  end
end

module FioHelper

  def self.start
    self.server = nil

    begin
      self.server = WEBrick::HTTPServer.new(
        Port: 4000,
        Logger: WEBrick::Log.new("/dev/null"),
        AccessLog: [],
        SSLEnable: true
      )

    rescue Exception => err
      raise err
      raise "Failed to allocate server binding! #{err}"
    end

    self.server.mount "/ib_api/rest/last", FioDownloadNewStatemetsHandler
    self.server.mount "/ib_api/rest/set-last-id", FioSetLastStatemetPivotIdHandler
    self.server.mount "/ib_api/rest/set-last-date", FioSetLastStatemetPivotDateHandler

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
