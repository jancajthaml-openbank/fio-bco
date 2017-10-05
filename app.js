let fioAccountStatement = require("./fio-response-example.json");
let fio = require("./fio.js");
let core = require("./core.js");
let sync = require("./sync");

async function main() {
  let coreAccountStatement = fio.extractCoreAccountStatement(fioAccountStatement);
  let accounts = fio.extractUniqueCoreAccounts(fioAccountStatement);

  await core.createMissingAccounts(accounts);
  await core.createTransactions(coreAccountStatement.transactions);

  await sync.updateTransactionCheckpoint("./db.json", coreAccountStatement.accountNumber, coreAccountStatement.idTransactionTo);
}


main()
  .then(() => console.log("FINISH - Everything went well"))
  .catch(error => {
    console.log(error);
    console.log("There were some unexpected error see above");
  });
