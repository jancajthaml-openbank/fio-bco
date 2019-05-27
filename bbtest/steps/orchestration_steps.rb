require_relative 'placeholders'

step "fio-bco is restarted" do ||
  ids = %x(systemctl -t service --no-legend | awk '{ print $1 }')
  expect($?).to be_success, ids

  ids = ids.split("\n").map(&:strip).reject { |x|
    x.empty? || !x.start_with?("fio-bco-import@")
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

step "fio-bco is running" do ||
  ids = %x(systemctl -t service --no-legend | awk '{ print $1 }')
  expect($?).to be_success, ids

  ids = ids.split("\n").map(&:strip).reject { |x|
    x.empty? || !x.start_with?("fio-bco-import@")
  }.map { |x| x.chomp(".service") }

  ids << "fio-bco-rest"

  eventually() {
    ids.each { |e|
      out = %x(systemctl show -p SubState #{e} 2>&1 | sed 's/SubState=//g')
      expect(out.strip).to eq("running")
    }
  }
end

step "tenant :tenant is offboarded" do |tenant|
  eventually() {
    %x(journalctl -o short-precise -u fio-bco-import@#{tenant}.service --no-pager > /reports/fio-bco@#{tenant}.log 2>&1)
    %x(systemctl stop fio-bco-import@#{tenant} 2>&1)
    %x(systemctl disable fio-bco-import@#{tenant} 2>&1)
    %x(journalctl -o short-precise -u fio-bco-import@#{tenant}.service --no-pager > /reports/fio-bco@#{tenant}.log 2>&1)
  }
end

step "tenant :tenant is onbdoarded" do |tenant|
  params = [
    "FIO_BCO_STORAGE=/data",
    "FIO_BCO_LOG_LEVEL=DEBUG",
    "FIO_BCO_FIO_GATEWAY=https://127.0.0.1:4000",
    "FIO_BCO_SYNC_RATE=1h",
    "FIO_BCO_VAULT_GATEWAY=https://127.0.0.1:4400",
    "FIO_BCO_LEDGER_GATEWAY=https://127.0.0.1:4401",
    "FIO_BCO_METRICS_OUTPUT=/reports",
    "FIO_BCO_LAKE_HOSTNAME=127.0.0.1",
    "FIO_BCO_METRICS_REFRESHRATE=1h",
    "FIO_BCO_HTTP_PORT=4002",
    "FIO_BCO_SECRETS=/opt/fio-bco/secrets",
    "FIO_BCO_ENCRYPTION_KEY=/opt/fio-bco/secrets/fs_encryption.key"
  ].join("\n").inspect.delete('\"')

  %x(mkdir -p /etc/init)
  %x(echo '#{params}' > /etc/init/fio-bco.conf)

  %x(systemctl enable fio-bco-import@#{tenant} 2>&1)
  %x(systemctl start fio-bco-import@#{tenant} 2>&1)

  ids = %x(systemctl list-units | awk '{ print $1 }')
  expect($?).to be_success, ids

  ids = ids.split("\n").map(&:strip).reject { |x|
    x.empty? || !x.start_with?("fio-bco-")
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

step "fio-bco is reconfigured with" do |configuration|
  params = Hash[configuration.split("\n").map(&:strip).reject(&:empty?).map { |el| el.split '=' }]
  defaults = {
    "STORAGE" => "/data",
    "LOG_LEVEL" => "DEBUG",
    "FIO_GATEWAY" => "https://127.0.0.1:4000",
    "SYNC_RATE" => "1h",
    "VAULT_GATEWAY" => "https://127.0.0.1:4400",
    "LEDGER_GATEWAY" => "https://127.0.0.1:4401",
    "METRICS_OUTPUT" => "/reports",
    "LAKE_HOSTNAME" => "127.0.0.1",
    "METRICS_REFRESHRATE" => "1h",
    "HTTP_PORT" => "4002",
    "SECRETS" => "/opt/fio-bco/secrets",
    "ENCRYPTION_KEY" => "/opt/fio-bco/secrets/fs_encryption.key"
  }

  config = Array[defaults.merge(params).map {|k,v| "FIO_BCO_#{k}=#{v}"}]
  config = config.join("\n").inspect.delete('\"')

  %x(mkdir -p /etc/init)
  %x(echo '#{config}' > /etc/init/fio-bco.conf)

  ids = %x(systemctl list-units | awk '{ print $1 }')
  expect($?).to be_success, ids

  ids = ids.split("\n").map(&:strip).reject { |x|
    x.empty? || !x.start_with?("fio-bco-")
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
