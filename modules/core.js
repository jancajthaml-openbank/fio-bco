const axios = require("axios")
const sync = require("./sync.js")
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

  async createMissingAccounts(accounts) {
    await parallelize(accounts, options.accountsParallelismSize,
      async account => {
        if (await this._accountExists(account.accountNumber)) {
          log.debug(`Account ${account.accountNumber} already exists`)
        } else {
          try {
            await axios.post(`${this._apiUrl}/account`, account)
            log.debug(`Created account ${account.accountNumber}`)
          } catch (err) {
            throw new VError(err, "Request to core api failed")
          }
        }
      }
    )

    log.info(`Created missing accounts for tenant ${this._tenantName}`)
  }

  async createTransactions(transactions, accountNumber, token) {
    await parallelize(transactions, options.transactionsParallelismSize,
      async (transaction, index) => {
        // TODO: should be part of fioTransfersToCoreTransactions, here we just pass core transaction to core
        const transferId = transaction.transfers.reduce((maxTransferId, transfer) => {
          const newMaxTransferId = getMax(maxTransferId, transfer.id)
          delete(transfer.id)
          return newMaxTransferId
        }, null)

        try {
          await axios.put(`${this._apiUrl}/transaction`, transaction)
          log.debug(`Created transaction ID ${transaction.id}`)
        } catch (err) {
          // Returned by core when transaction already exists but with different data
          if (err.response && err.response.status === 406) {
            log.warn(`Transaction with ID ${transaction.id} already exits in core but has different data.
                      Source data: 
                      ` + JSON.stringify(transaction, null, 2))
          } else {
            throw new VError(err, "Unable to create transaction in core")
          }
        }
        return transferId
      },
      async (transferIds) => {
        const max = transferIds.reduce(getMax)
        await sync.setTransactionCheckpoint(options.db, this._tenantName, accountNumber, token, max)
        log.debug(`Max ID ${max}`)
      }
    )

    log.info(`Created transactions for tenant ${this._tenantName}`)
  }

  async getCheckpointByAccountNumber(accountNumber) {
    const transactionCheckpoint = await sync.getTransactionCheckpoint(options.db, this._tenantName, accountNumber)
    log.info(`Checkpoint for tenant/account ${this._tenantName}/${accountNumber}:${transactionCheckpoint}`)
    return transactionCheckpoint
  }

  async getCheckpointByToken(token) {
    const transactionCheckpoint = await sync.getTransactionCheckpointByToken(options.db, this._tenantName, token)
    log.info(`Checkpoint (by token) for tenant ${this._tenantName}:${transactionCheckpoint}`)
    return transactionCheckpoint
  }

  get _apiUrl() {
    return `${options.url}/v1/${this._tenantName}/core`
  }

  async _accountExists(accountNumber) {
    try {
      await axios.get(`${this._apiUrl}/account/${accountNumber}`)
      return true
    } catch (err) {
      if (err.response && err.response.status === 404) {
        return false
      } else {
        throw new VError(err, "Request to core api failed")
      }
    }
  }
}

module.exports = {
  Tenant
}
