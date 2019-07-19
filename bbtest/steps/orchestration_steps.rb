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
    %x(journalctl -o short-precise -u fio-bco-import@#{tenant}.service --no-pager > /tmp/reports/fio-bco@#{tenant}.log 2>&1)
    %x(systemctl stop fio-bco-import@#{tenant} 2>&1)
    %x(systemctl disable fio-bco-import@#{tenant} 2>&1)
    %x(journalctl -o short-precise -u fio-bco-import@#{tenant}.service --no-pager > /tmp/reports/fio-bco@#{tenant}.log 2>&1)
  }
end

step "tenant :tenant is onbdoarded" do |tenant|
  config = Array[UnitHelper.default_config.map {|k,v| "FIO_BCO_#{k}=#{v}"}]
  config = config.join("\n").inspect.delete('\"')

  %x(mkdir -p /etc/init)
  %x(echo '#{config}' > /etc/init/fio-bco.conf)

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
  config = Array[UnitHelper.default_config.merge(params).map {|k,v| "FIO_BCO_#{k}=#{v}"}]
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
