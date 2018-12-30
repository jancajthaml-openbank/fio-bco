class Object
  def blank?
    respond_to?(:empty?) ? !!empty? : !self
  end

  def present?
    !blank?
  end
end

class Hash

  def deep_diff(b, exceptions = [])
    a = self
    unless b.instance_of?(Hash)
      puts a, b
    end
    (a.keys | b.keys).inject({}) do |diff, k|
      if a[k] != b[k] && !exceptions.any? { |key| k.include?(key) }
        if a[k].respond_to?(:deep_diff) && b[k].respond_to?(:deep_diff)
          diff[k] = a[k].deep_diff(b[k])
        else
          if a[k].instance_of?(Array) && b[k].instance_of?(Array)
            a[k].sort_by! { |h| h }
            b[k].sort_by! { |h| h }
            if a[k].present? && a[k].first.instance_of?(Hash)
              a[k].each_with_index do |hash, index|
                if (b[k][index]).present?
                  diff[k] = hash.deep_diff(b[k][index])
                else
                  diff[k] = [a[k], b[k]]
                end
              end
            else
              delta = (a[k] | b[k]) - (a[k] & b[k])
              diff[k] = [a[k], b[k]] unless delta.empty?
            end
          elsif a[k] != b[k]
            diff[k] = [a[k], b[k]]
          end
        end
      end
      diff.delete_blank
    end
  end

  def delete_blank
    delete_if { |_, v| v.empty? or v.instance_of?(Hash) && v.delete_blank.empty? }
  end

end
