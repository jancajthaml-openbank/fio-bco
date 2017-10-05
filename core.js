let axios = require("axios");

async function createMissingAccounts(accounts) {
  await Promise.all(accounts.map(
    async account => {
      try {
        await axios.get("http://localhost:8080/v1/johny/core/account/" + account.accountNumber);
        console.log("Account " + account.accountNumber + " already exists");
      } catch (error) {
        if (error.response && error.response.status === 404) {
          await axios.put("http://localhost:8080/v1/johny/core/account/", account);
          console.log("Created account " + account.accountNumber);
        } else {
          throw error;
        }
      }
    }
  ));
}

async function createTransactions(transactions) {
  await Promise.all(transactions.map(
    async (transaction, index) => {
      await axios.put("http://localhost:8080/v1/johny/core/transaction", transaction);
      console.log("Created " + (index+1) + ". transaction");
    }
  ));
}

module.exports = {
  createMissingAccounts,
  createTransactions
}