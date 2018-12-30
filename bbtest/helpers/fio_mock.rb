require 'bigdecimal'

module FioMock

  class << self
    attr_accessor :info, :transfers
  end

  def self.set_statements(statements)
    self.transfers = statements["accountStatement"]["transactionList"]["transaction"]
    self.info = statements["accountStatement"]["info"]
  end

  def self.get_statements()
    transfers = self.transfers unless self.transfers.nil?
    lastId = self.info["idLastDownload"] unless self.info.nil?

    transfers.reject! { |x| x["column22"]["value"].to_i <= lastId } unless lastId.nil?

    {
      "accountStatement" => {
        "info" => self.info,
        "transactionList" => {
          "transaction" => transfers,
        }
      }
    }
  end

  def self.set_confirmed_transfer_pivot_id(transferId)
    begin
      self.info["idLastDownload"] = transferId.to_i
    rescue
    end
  end

  def self.get_last_transfer_id()
    return nil if self.transfers.nil?
    self.transfers.map { |x| x["column22"]["value"].to_i }.to_a.max
  end

  def self.set_confirmed_transfer_pivot_date(date)
    puts "not implemented last transferId by date"

    #self.info[:idLastDownload] = transferId
    #https://www.fio.cz/ib_api/rest/set-last-id/{token}/{transferId}/
  end

  def get_statements()
    FioMock.get_statements()
  end

  def get_last_transfer_id()
    FioMock.get_last_transfer_id()
  end

  def set_statements(statements)
    FioMock.set_statements(statements)
  end
end
