require 'turnip/rspec'
require 'json'
require 'thread'
require 'openssl'

Thread.abort_on_exception = true

RSpec.configure do |config|
  config.raise_error_for_unimplemented_steps = true
  config.color = true
  config.fail_fast = true

  Dir.glob("./helpers/*_helper.rb") { |f| load f }
  config.include EventuallyHelper, :type => :feature
  Dir.glob("./steps/*_steps.rb") { |f| load f, true }

  config.register_ordering(:global) do |items|
    (install, others) = items.partition { |spec| spec.metadata[:install] }
    (uninstall, others) = others.partition { |spec| spec.metadata[:uninstall] }

    install + others.shuffle + uninstall
  end

  $unit = UnitHelper.new()

  config.before(:suite) do |_|
    print "[ suite starting ]\n"

    LakeMock.start()
    FioHelper.start()
    VaultHelper.start()
    LedgerHelper.start()

    ["/reports"].each { |folder|
      FileUtils.mkdir_p folder
      %x(rm -rf #{folder}/*)
    }

    print "[ downloading unit ]\n"

    $unit.download()

    print "[ suite started    ]\n"
  end

  config.after(:type => :feature) do
    LakeMock.reset()

    $unit.cleanup()
  end

  config.after(:suite) do
    print "\n[ suite ending   ]\n"

    $unit.teardown()

    VaultHelper.stop()
    LedgerHelper.stop()
    FioHelper.stop()
    LakeMock.stop()

    print "[ suite cleaning ]\n"

    ["/data"].each { |folder|
      %x(rm -rf #{folder}/*)
    }

    print "[ suite ended    ]"
  end


end
