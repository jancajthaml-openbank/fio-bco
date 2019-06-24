require 'bigdecimal'

module LedgerMock

  class << self
    attr_accessor :tenants
  end

  self.tenants = Hash.new()

  def self.reset()
    self.tenants = Hash.new()
  end

  def self.get_transactions(tenant)
    return nil unless self.tenants.has_key?(tenant)
    return self.tenants[tenant].keys
  end

  def self.get_transaction(tenant, id)
    return {} unless self.tenants.has_key?(tenant)
    return {} unless self.tenants[tenant].has_key?(id)
    return self.tenants[tenant][id]
  end

  def self.create_transaction(tenant, id, transfers, status)
    return if self.tenants.has_key?(tenant) && self.tenants[tenant].has_key?(id)
    self.tenants[tenant] = Hash.new() unless self.tenants.has_key?(tenant)
    self.tenants[tenant][id] = {
      :id => id,
      :status => (if status.nil? then "committed" else status end),
      :transfers => transfers,
    }
    return
  end
end
