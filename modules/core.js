const axios = require("axios")
const { setTransactionCheckpoint, getTransactionCheckpoint } = require("./sync.js")
const { parallelize, getMax } = require("./utils.js")
const log = require('./logger.js')
const VError = require("verror")

const options = require("config").get("core")

class Tenant {

  constructor(tenantName) {
    if (!tenantName) {
      throw Error("When creating Tenant you have to provide his name")
    }

    this._tenantName = tenantName
  }

  get apiUrl() {
    return `${options.url}/v1/${this._tenantName}/core`
  }

  async _accountExists(accountNumber) {
    try {
      await axios.get(`${this.apiUrl}/account/${accountNumber}`)
      return true
    } catch (err) {
      if (err.response && err.response.status === 404) {
        return false
      } else {
        throw new VError(err, "Request to core api failed")
      }
    }
  }

  async createMissingAccounts(accounts) {
    await parallelize(accounts, options.accountsParallelismSize,
      async account => {
        if (await this._accountExists(account.accountNumber)) {
          log.debug(`Account ${account.accountNumber} already exists`)
        } else {
          try {
            await axios.post(`${this.apiUrl}/account`, account)
            log.debug(`Created account ${account.accountNumber}`)
          } catch (err) {
            throw new VError(err, "Request to core api failed")
          }
        }
      }
    )

    log.info(`Created missing accounts for tenant ${this._tenantName}`)
  }

  async createTransactions(transactions, accountNumber) {
    await parallelize(transactions, options.transactionsParallelismSize,
      async (transaction, index) => {
        const transferId = transaction.transfers.reduce((maxTransferId, transfer) => {
          const newMaxTransferId = getMax(maxTransferId, transfer.id)
          delete(transfer.id)
          return newMaxTransferId
        }, null)

        await axios.put(`${this.apiUrl}/transaction`, transaction)
        log.debug(`Created transaction ID ${transaction.id}`)
        return transferId
      },
      async (transferIds) => {
        const max = transferIds.reduce(getMax)
        await setTransactionCheckpoint(options.db, this._tenantName, accountNumber, max)
        log.debug(`Max ID ${max}`)
      }
    )

    log.info(`Created transactions for tenant ${this._tenantName}`)
  }

  async getCheckpoint(accountNumber) {
    const transactionCheckpoint = await getTransactionCheckpoint(options.db, this._tenantName, accountNumber)
    log.info(`Checkpoint for tenant/account ${this._tenantName}/${accountNumber}:${transactionCheckpoint}`)
    return transactionCheckpoint
  }
}

module.exports = {
  Tenant
}
