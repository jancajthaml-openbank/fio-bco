
step "fio gateway contains following statements" do |statements|
  statements = JSON.parse(statements)

  FioMock.set_statements(statements)
end
