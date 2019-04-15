
step "token :slash_pair is created" do |slash_pair|
  @tokens ||= {}

  expect(@tokens).not_to have_key(slash_pair), "expected not to find \"#{slash_pair}\" in #{@tokens.keys}"

  (tenant, token) = slash_pair.split('/')

  send "I request curl :http_method :url", "POST", "https://localhost/token/#{tenant}", {
    "username": "X",
    "password": "Y"
  }.to_json

  @resp = { :code => 0 }

  eventually(timeout: 60, backoff: 2) {
    resp = %x(#{@http_req})
    @resp[:code] = resp[resp.length-3...resp.length].to_i

    if @resp[:code] === 0
      raise "endpoint #{@http_req} is unreachable"
    end

    @resp[:body] = resp[0...resp.length-3] unless resp.nil?
  }

  resp = JSON.parse(@resp[:body])

  @tokens[tenant+"/"+token] = resp["value"] if @resp[:code] == 200
end

step "token :slash_pair should exist" do |slash_pair|
  @tokens ||= {}

  expect(@tokens).to have_key(slash_pair), "expected to find \"#{slash_pair}\" in #{@tokens.keys}"

  (tenant, token) = slash_pair.split('/')
  token_value = @tokens[slash_pair]

  send "I request curl :http_method :url", "GET", "https://localhost/token/#{tenant}"

  @resp = { :code => 0 }

  eventually(timeout: 60, backoff: 2) {
    resp = %x(#{@http_req})
    @resp[:code] = resp[resp.length-3...resp.length].to_i

    if @resp[:code] === 0
      raise "endpoint #{@http_req} is unreachable"
    end

    @resp[:body] = resp[0...resp.length-3] unless resp.nil?
  }

  send "request should succeed"

  actual_tokens = JSON.parse(@resp[:body]).map { |item| tenant + "/" + item }

  expect(actual_tokens).to include(tenant+"/"+token_value), "expected to find \"#{slash_pair}\" in #{actual_tokens}"
end

step "request should succeed" do ||
  expect(@resp[:code]).to eq(200), "actual response #{@resp}"
end

step "request should fail" do ||
  expect(@resp[:code]).to_not eq(200), "actual response #{@resp}"
end
