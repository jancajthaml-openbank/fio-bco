const axios = require("axios");
const sync = require("./sync.js");
const options = require("config").get("core");
const log = require("winston");

function Tenant (tenantName) {
  if (!tenantName) {
    throw Error("When creating Tenant you have to provide his name");
  }
  this._tenantName = tenantName;
}

Tenant.prototype._runInParallel = async function (items, parallelismSize, processItem, afterBatch) {
  for (let i = 0; i < items.length; i += parallelismSize) {
    const bulk = items.slice(i, Math.min(i + parallelismSize, items.length));
    const batchResult = await Promise.all(bulk.map((item, index) => processItem(item, index)));
    if (afterBatch) await afterBatch(batchResult);
    log.debug("Finished " + (Math.floor(i/parallelismSize)+1) + ". bulk");
  }
};

Tenant.prototype._getApiUrl = function () {
  return options.url + "/v1/" + this._tenantName + "/core";
};

Tenant.prototype.createMissingAccounts = async function (accounts) {
  await this._runInParallel(accounts, options.accountsParallelismSize,
    async account => {
      try {
        await axios.get(this._getApiUrl() + "/account/" + account.accountNumber);
        log.debug("Account " + account.accountNumber + " already exists");
      } catch (error) {
        if (error.response && error.response.status === 404) {
          await axios.post(this._getApiUrl() + "/account/", account);
          log.debug("Created account " + account.accountNumber);
        } else {
          throw error;
        }
      }
    }
  );
  log.info("Created missing accounts for tenant " + this._tenantName);
};

Tenant.prototype.createTransactions = async function (transactions, accountNumber) {
  await this._runInParallel(transactions, options.transactionsParallelismSize,
    async (transaction, index) => {
      const transferId = transaction.transfers.reduce((maxTransferId, transfer) => {
        const newMaxTransferId = Math.max(maxTransferId, transfer.id);
        delete(transfer.id);
        return newMaxTransferId;
      }, null);
      await axios.put(this._getApiUrl() + "/transaction", transaction);
      log.debug("Created transaction ID " + transaction.id);
      return transferId;
    },
    async (transferIds) => {
      const max = transferIds.reduce((max, transactionId) => {
        return Math.max(max, transactionId);
      });
      await sync.setTransactionCheckpoint(options.db, this._tenantName, accountNumber, max);
      log.debug("Max ID " + max);
    }
  );
  log.info("Created transactions for tenant " + this._tenantName);
};

Tenant.prototype.getTransactionCheckpoint = async function (accountNumber) {
  transactionCheckpoint = await sync.getTransactionCheckpoint(options.db, this._tenantName, accountNumber);
  log.info("Checkpoint for tenant/account " + this._tenantName + "/" + accountNumber + ": " + transactionCheckpoint);
  return transactionCheckpoint;
};

module.exports = {
  Tenant
};