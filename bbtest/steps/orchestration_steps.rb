require_relative 'placeholders'

step "fio-bco is restarted" do ||
  ids = %x(systemctl -t service --no-legend | awk '{ print $1 }')
  expect($?).to be_success, ids

  ids = ids.split("\n").map(&:strip).reject { |x|
    x.empty? || !x.start_with?("fio-bco@")
  }.map { |x| x.chomp(".service") }

  expect(ids).not_to be_empty

  ids.each { |e|
    %x(systemctl restart #{e} 2>&1)
  }

  eventually() {
    ids.each { |e|
      out = %x(systemctl show -p SubState #{e} 2>&1 | sed 's/SubState=//g')
      expect(out.strip).to eq("running")
    }
  }
end

step "tenant :tenant is offboarded" do |tenant|
  eventually() {
    %x(journalctl -o short-precise -u fio-bco@#{tenant}.service --no-pager > /reports/fio-bco@#{tenant}.log 2>&1)
    %x(systemctl stop fio-bco@#{tenant} 2>&1)
    %x(systemctl disable fio-bco@#{tenant} 2>&1)
    %x(journalctl -o short-precise -u fio-bco@#{tenant}.service --no-pager > /reports/fio-bco@#{tenant}.log 2>&1)
  }
end

step "tenant :tenant is onbdoarded" do |tenant|
  params = [
    "FIO_BCO_STORAGE=/data",
    "FIO_BCO_LOG_LEVEL=DEBUG",
    "FIO_BCO_FIO_GATEWAY=https://localhost:4000",
    "FIO_BCO_SYNC_RATE=1h",
    "FIO_BCO_WALL_GATEWAY=https://localhost:3000",
    "FIO_BCO_METRICS_OUTPUT=/reports/metrics.json",
    "FIO_BCO_LAKE_HOSTNAME=localhost",
    "FIO_BCO_METRICS_REFRESHRATE=1h",
    "FIO_BCO_HTTP_PORT=443",
    "FIO_BCO_SECRETS=/opt/fio-bco/secrets"
  ].join("\n").inspect.delete('\"')

  %x(mkdir -p /etc/init)
  %x(echo '#{params}' > /etc/init/fio-bco.conf)

  %x(systemctl enable fio-bco@#{tenant} 2>&1)
  %x(systemctl start fio-bco@#{tenant} 2>&1)

  eventually() {
    out = %x(systemctl show -p SubState fio-bco@#{tenant} 2>&1 | sed 's/SubState=//g')
    expect(out.strip).to eq("running")
  }
end

step "fio-bco is reconfigured with" do |configuration|
  params = Hash[configuration.split("\n").map(&:strip).reject(&:empty?).map {|el| el.split '='}]
  defaults = {
    "STORAGE" => "/data",
    "LOG_LEVEL" => "DEBUG",
    "FIO_GATEWAY" => "https://localhost:4000",
    "SYNC_RATE" => "1h",
    "WALL_GATEWAY" => "https://localhost:3000",
    "METRICS_OUTPUT" => "/reports/metrics.json",
    "LAKE_HOSTNAME" => "localhost",
    "METRICS_REFRESHRATE" => "1h",
    "HTTP_PORT" => "443",
    "SECRETS" => "/opt/fio-bco/secrets"
  }

  config = Array[defaults.merge(params).map {|k,v| "FIO_BCO_#{k}=#{v}"}]
  config = config.join("\n").inspect.delete('\"')

  %x(mkdir -p /etc/init)
  %x(echo '#{config}' > /etc/init/fio-bco.conf)

  ids = %x(systemctl list-units | awk '{ print $1 }')
  expect($?).to be_success, ids

  ids = ids.split("\n").map(&:strip).reject { |x|
    x.empty? || !x.start_with?("fio-bco")
  }.map { |x| x.chomp(".service") }

  expect(ids).not_to be_empty

  ids.each { |e|
    %x(systemctl restart #{e} 2>&1)
  }

  eventually() {
    ids.each { |e|
      out = %x(systemctl show -p SubState #{e} 2>&1 | sed 's/SubState=//g')
      expect(out.strip).to eq("running")
    }
  }
end
