require 'bigdecimal'

module VaultMock

  class << self
    attr_accessor :accounts
  end

  self.accounts = {}

  def self.reset()
    self.accounts = {}
  end

  def self.get_acount(id)
    return self.accounts[id]
  end

  def self.account_exist?(id)
    return self.accounts.has_key?(id)
  end

  def self.create_account(id, currency, is_balance_check)
    return false if self.accounts.has_key?(id)

    self.accounts[id] = {
      :currency => currency,
      :is_balance_check => is_balance_check,
      :balance => BigDecimal.new("0"),
      :blocking => BigDecimal.new("0"),
      :promised => {}
    }
    return true
  end

  def self.create_transaction(id, transfers)
    #return false if self.accounts.has_key?(id)

    #self.accounts[id] = {
      #:currency => currency,
      #:is_balance_check => is_balance_check,
      #:balance => BigDecimal.new("0"),
      #:blocking => BigDecimal.new("0"),
      #:promised => {}
    #}
    return true
  end

end
