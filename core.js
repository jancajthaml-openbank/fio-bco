let axios = require("axios");

function CoreTenant (coreHost, tenant) {
  if (!coreHost || !tenant) {
    throw Error("When creating CoreTenant you have to provide both parameters - coreHost and tenant");
  }
  this._coreHost = coreHost;
  this._tenant = tenant;
}

CoreTenant.prototype._getApiUrl = function () {
  return this._coreHost + "/v1/" + this._tenant + "/core";
};

CoreTenant.prototype.createMissingAccounts = async function (accounts) {
  console.log(this._getApiUrl());
  await Promise.all(accounts.map(
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
  ));
};

CoreTenant.prototype.createTransactions = async function (transactions) {
  await Promise.all(transactions.map(
    async (transaction, index) => {
      await axios.put(this._getApiUrl() + "/transaction", transaction);
      console.log("Created " + (index+1) + ". transaction");
    }
  ));
};

module.exports = {
  CoreTenant
};