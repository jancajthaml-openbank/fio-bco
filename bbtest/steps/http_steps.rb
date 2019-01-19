require_relative '../shims/deep_diff'

require 'json'
require 'net/http'
require 'time'

step "I request curl :http_method :url" do |http_method, url, body = nil|
  cmd = ["curl --insecure"]
  cmd << ["-X #{http_method.upcase}"] unless http_method.upcase == "GET"
  cmd << ["#{url} -sw \"%{http_code}\""]
  cmd << ["-d \'#{JSON.parse(body).to_json}\'"] unless body.nil? or http_method.upcase == "GET"

  @http_req = cmd.join(" ")
end

step "curl responds with :http_status" do |http_status, body = nil|
  raise if @http_req.nil?

  @resp = Hash.new
  resp = %x(#{@http_req})

  @resp[:code] = resp[resp.length-3...resp.length].to_i
  @resp[:body] = resp[0...resp.length-3] unless resp.nil?

  expect(@resp[:code]).to eq(http_status)

  return if body.nil?

  expected_body = { wrap: JSON.parse(body) }

  begin
    resp_body = { wrap: JSON.parse(@resp[:body]) }
    resp_body.deep_diff(expected_body).each do |key, array|
      (have, want) = array
      raise "unexpected attribute \"#{key}\" in response \"#{@resp[:body]}\" expected \"#{expected_body.to_json}\"" if want.nil?
      raise "\"#{key}\" expected \"#{want}\" but got \"#{have}\" instead"
    end
  rescue JSON::ParserError
    raise "invalid response got \"#{@resp[:body].strip}\", expected \"#{expected_body.to_json}\""
  end

end


