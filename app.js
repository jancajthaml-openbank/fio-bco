let fioAccountStatement = require("./fio-response-example.json");
let fio = require("./fio.js");
let core = require("./core.js");
let sync = require("./sync.js");

async function main(args) {
  let transactionCheckpoint = await sync.getTransactionCheckpoint(args.db, args.iban);
  console.log("Starting from " + transactionCheckpoint);
  let coreAccountStatement = fio.extractCoreAccountStatement(fioAccountStatement);
  let accounts = fio.extractUniqueCoreAccounts(fioAccountStatement);

  let tenantJohny = new core.CoreTenant(args.coreHost, args.coreTenant);
  await tenantJohny.createMissingAccounts(accounts);
  await tenantJohny.createTransactions(coreAccountStatement.transactions);

  await sync.setTransactionCheckpoint(args.db, coreAccountStatement.accountNumber, coreAccountStatement.idTransactionTo);
}

main({
  "iban": "CZ7920100000002400222233",
  "token": "xxx",
  "fioApiUrl": "xxx",
  "db": "./db.json",
  "coreHost": "http://localhost:8080",
  "coreTenant": "johny"
}).then(() => console.log("FINISH - Everything went well"))
  .catch(error => {
    console.log(error);
    console.log(error.message);
    console.log("There were some unexpected error see above");
  });
