const axios = require("axios")
const sync = require("./sync.js")
const { elapsedTime, parallelize, getMax } = require("./utils.js")
const log = require("./logger.js")
const VError = require("verror")

const options = require("config").get("core")

class Tenant {

  constructor(tenant) {
    if (!tenant) {
      throw Error("When creating Tenant you have to provide his name")
    }

    this._tenant = tenant
  }

  async createMissingAccounts(accounts) {
    log.info(`Asserting ${accounts.length} accounts for tenant ${this._tenant}`)

    let t0 = process.hrtime()
    await parallelize(accounts,
      async (account) => {
        if (await this._accountExists(account.accountNumber)) {
          log.debug(`Account ${account.accountNumber} already exists`)
        } else {
          try {
            await axios.post(`${this._baseUrl}/account/${this._tenant}`, account)
            log.debug(`Created account ${account.accountNumber}`)
          } catch (err) {
            throw new VError(err, "Account core api")
          }
        }
      }
    )

    log.info(`Creation of ${accounts.length} accounts took ${elapsedTime(t0)}.`)
  }

  async createTransactions(transactions, accountNumber, token) {
    if (transactions.length > 0) {
      log.info(`Creating ${transactions.length} new transactions for tenant ${this._tenant}`)

      let t0 = process.hrtime()
      await parallelize(transactions,
        async (transaction, index) => {
          //const transferId = transaction.transfers.map((transfer) => transfer.id).reduce(getMax)

          try {
            await axios.post(`${this._baseUrl}/transaction/${this._tenant}`, transaction)
            log.debug(`Created transaction ID ${transaction.id}`)
          } catch (err) {
            // Returned by core when transaction already exists but with different data
            if (err.response && err.response.status === 406) {
              log.warn(`Transaction with ID ${transaction.id} already exits in core but has different data.
                        Source data:
                        ` + JSON.stringify(transaction, null, 2))
            } else {
              throw new VError(err, "Transaction core api")
            }
          }
          return transaction.transfers.map((transfer) => transfer.id).reduce(getMax)
        },
        async (transferIds) => {
          const max = transferIds.reduce(getMax)
          await sync.setTransactionCheckpoint(options.db, this._tenant, accountNumber, token, max)
          log.debug(`Max TransferId ${max}`)
          log.info(`Creation of ${transactions.length} transactions took ${elapsedTime(t0)}.`)
        }
      )
    }
  }

  async getCheckpointByAccountNumber(accountNumber) {
    const transactionCheckpoint = await sync.getTransactionCheckpoint(options.db, this._tenant, accountNumber)
    log.info(`Checkpoint (by account ${accountNumber}) for tenant ${this._tenant}:${transactionCheckpoint}`)
    return transactionCheckpoint
  }

  async getCheckpointByToken(token) {
    const transactionCheckpoint = await sync.getTransactionCheckpointByToken(options.db, this._tenant, token)
    log.info(`Checkpoint (by token) for tenant ${this._tenant}:${transactionCheckpoint}`)
    return transactionCheckpoint
  }

  get _baseUrl() {
    return `${options.url}/v1/sparrow`
  }

  async _accountExists(accountNumber) {
    try {
      await axios.get(`${this._baseUrl}/account/${this._tenant}/${accountNumber}`)
      return true
    } catch (err) {
      if (err.response && err.response.status === 404) {
        return false
      } else {
        throw new VError(err, "Account core api")
      }
    }
  }
}

module.exports = {
  Tenant
}
