let fioAccountStatement = require("./fio-response-example.json");
let fio = require("./fio.js");
let core = require("./core.js");
let sync = require("./sync.js");

async function main() {
  let coreAccountStatement = fio.extractCoreAccountStatement(fioAccountStatement);
  let accounts = fio.extractUniqueCoreAccounts(fioAccountStatement);

  let tenantJohny = new core.CoreTenant("http://localhost:8080", "johny");
  await tenantJohny.createMissingAccounts(accounts);
  await tenantJohny.createTransactions(coreAccountStatement.transactions);

  await sync.updateTransactionCheckpoint("./db.json", coreAccountStatement.accountNumber, coreAccountStatement.idTransactionTo);
}


main()
  .then(() => console.log("FINISH - Everything went well"))
  .catch(error => {
    console.log(error);
    console.log(error.message);
    console.log("There were some unexpected error see above");
  });
