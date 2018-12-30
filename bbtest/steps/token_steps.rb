
step "token :slash_pair is created" do |slash_pair|
  @tokens ||= {}

  expect(@tokens).not_to have_key(slash_pair), "expected not to find \"#{slash_pair}\" in #{@tokens.keys}"

  (tenant, token) = slash_pair.split('/')

  send "I request curl :http_method :url", "POST", "https://localhost/token/#{tenant}", {
    "value" => token
  }.to_json

  @resp = Hash.new
  resp = %x(#{@http_req})

  @resp[:code] = resp[resp.length-3...resp.length].to_i
  @resp[:body] = resp[0...resp.length-3] unless resp.nil?
  @tokens[tenant+"/"+token] = true if @resp[:code] == 200
end

step "token :slash_pair should exist" do |slash_pair|
  @tokens ||= {}

  expect(@tokens).to have_key(slash_pair), "expected to find \"#{slash_pair}\" in #{@tokens.keys}"

  (tenant, token) = slash_pair.split('/')

  send "I request curl :http_method :url", "GET", "https://localhost/tokens/#{tenant}"

  @resp = Hash.new
  resp = %x(#{@http_req})

  @resp[:code] = resp[resp.length-3...resp.length].to_i
  @resp[:body] = resp[0...resp.length-3] unless resp.nil?
  send "request should succeed"

  actual_tokens = JSON.parse(@resp[:body]).map { |item| tenant + "/" + item["value"] }

  expect(actual_tokens).to include(slash_pair), "expected to find \"#{slash_pair}\" in #{actual_tokens}"
end

step "request should succeed" do ||
  expect(@resp[:code]).to eq(200), "actual response #{@resp}"
end

step "request should fail" do ||
  expect(@resp[:code]).to_not eq(200), "actual response #{@resp}"
end
