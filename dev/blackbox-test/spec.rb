require 'turnip/rspec'

RSpec.configure do |config|
  config.raise_error_for_unimplemented_steps = true
  config.color = true

  Dir.glob("./**/*_helper.rb") { |f| load f }
  config.include EventuallyHelper, :type => :feature
  Dir.glob("./**/*_steps.rb") { |f| load f, true }

  config.after(:suite) do |suite|
    FileUtils.rm_rf("/data/.", secure: true)
  end

  puts "starting"
end


class Hash

  def deep_diff(b)
    a = self
    (a.keys | b.keys).inject({}) do |diff, k|
      if a[k] != b[k]
        if a[k].respond_to?(:deep_diff) && b[k].respond_to?(:deep_diff)
          diff[k] = a[k].deep_diff(b[k])
        else
          diff[k] = [a[k], b[k]]
        end
      end
      diff
    end
  end

end
