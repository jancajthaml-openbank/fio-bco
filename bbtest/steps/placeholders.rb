
placeholder :slice do
  match(/([^\s]{1,100},?+)$/) do |slice|
    slice.split(",")
  end
end

placeholder :amount do
  match(/-?\d{1,100}\.\d{1,100}|-?\d{1,100}/) do |amount|
    amount
  end
end

placeholder :count do
  match(/\d{1,100}/) do |count|
    count.to_i
  end
end

placeholder :path do
  match(/((?:\/[a-z0-9]+[a-z0-9(\/)(\-)]{1,100}[\w,\s-]+\.?[A-Za-z0-9_-]{0,100})+)/) do |path|
    path
  end
end

placeholder :http_method do
  match(/(GET|get|POST|post|PATCH|patch|DELETE|delete)/) do |http_method|
    http_method.downcase
  end
end

placeholder :http_status do
  match(/\d{3}/) do |http_status|
    http_status.to_i
  end
end

placeholder :url do
  match(/https?:\/\/[\S]+/) do |url|
    url
  end
end

placeholder :slash_pair do
  match(/[\S]+\/[\S]+/) do |slash_pair|
    slash_pair
  end
end
