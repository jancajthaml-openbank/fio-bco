let axios = require("axios");

let defaultOptions = {
  "accountsParallelismSize": 1,
  "transactionsParallelismSize": 2
} ;

function CoreTenant (coreHost, tenant) {
  if (!coreHost || !tenant) {
    throw Error("When creating CoreTenant you have to provide both parameters - coreHost and tenant");
  }
  this._coreHost = coreHost;
  this._tenant = tenant;
  this._options = defaultOptions;
}

CoreTenant.prototype._runInParallel = async function (items, parallelismSize, processItem) {
  for (let i = 0; i < items.length; i += parallelismSize) {
    let bulk = items.slice(i, Math.min(i + parallelismSize, items.length));
    await Promise.all(bulk.map((item, index) => processItem(item, index)));
    console.log("Finished " + (Math.floor(i/parallelismSize)+1) + ". bulk");
  }
};

CoreTenant.prototype._getApiUrl = function () {
  return this._coreHost + "/v1/" + this._tenant + "/core";
};

CoreTenant.prototype.createMissingAccounts = async function (accounts) {
  await this._runInParallel(accounts, this._options.accountsParallelismSize,
    async account => {
      try {
        await axios.get(this._getApiUrl() + "/account/" + account.accountNumber);
        console.log("Account " + account.accountNumber + " already exists");
      } catch (error) {
        if (error.response && error.response.status === 404) {
          await axios.put(this._getApiUrl() + "/account/", account);
          console.log("Created account " + account.accountNumber);
        } else {
          throw error;
        }
      }
    }
  );
};

CoreTenant.prototype.createTransactions = async function (transactions) {
  await this._runInParallel(transactions, this._options.transactionsParallelismSize,
    async (transaction, index) => {
      await axios.put(this._getApiUrl() + "/transaction", transaction);
      console.log("Created " + (index+1) + ". transaction");
    }
  );
};

module.exports = {
  CoreTenant
};