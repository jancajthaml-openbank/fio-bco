
step "container :image should be running" do |container_name|
  container_id = %x(docker ps -aqf "name=#{container_name}" 2>/dev/null)
  expect($?).to eq(0), "error running `docker ps -aq`: err:\n #{container_name}"

  eventually(timeout: 10) {
    container_state = %x(docker inspect -f {{.State.Running}} #{container_id} 2>/dev/null)
    expect($?).to eq(0), "error running `docker inspect -f {{.State.Running}}`: err:\n #{container_id}"

    expect(container_state.strip).to eq("true")
  }
end

step ":host is listening on :port" do |host, port|
  eventually(timeout: 10) {
    `nc -z #{host} #{port}`
    expect($?).to be_success
  }
end

step "server is healthy" do ||
  $http_client.server_service.health_check()
end
