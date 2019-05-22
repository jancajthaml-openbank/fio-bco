
step "token :slash_pair is created" do |slash_pair|
  @tokens ||= {}

  expect(@tokens).not_to have_key(slash_pair), "expected not to find \"#{slash_pair}\" in #{@tokens.keys}"

  (tenant, token) = slash_pair.split('/')

  uri = "https://127.0.0.1:4002/token/#{tenant}"

  payload = {
    value: "X",
  }.to_json

  send "I request curl :http_method :url", "POST", uri, payload
  send "curl responds with :http_status", 200

  token_value = JSON.parse(HTTPHelper.response[:body])["value"]
  @tokens[tenant+"/"+token] = token_value
end

step "token :slash_pair should exist" do |slash_pair|
  @tokens ||= {}

  expect(@tokens).to have_key(slash_pair), "expected to find \"#{slash_pair}\" in #{@tokens.keys}"

  (tenant, token) = slash_pair.split('/')
  token_value = @tokens[slash_pair]

  uri = "https://127.0.0.1:4002/token/#{tenant}"

  send "I request curl :http_method :url", "GET", uri
  send "curl responds with :http_status", 200

  actual_tokens = JSON.parse(HTTPHelper.response[:body]).map { |item| tenant + "/" + item["value"] }

  puts HTTPHelper.response[:body]

  expect(token_value).not_to be_empty

  expect(actual_tokens).to include(tenant+"/"+token_value), "expected to find \"#{slash_pair}\" in #{actual_tokens}"
end

step "request should succeed" do ||
  expect(HTTPHelper.response[:code]).to eq(200), "expected 200 got\n#{HTTPHelper.response[:raw]}"
end

step "request should fail" do ||
  expect(HTTPHelper.response[:code]).to_not eq(200), "expected non 200 got\n#{HTTPHelper.response[:raw]}"
end
