let fioAccountStatement = require("./fio-response-example.json");
let fio = require("./fio.js");
let core = require("./core.js");

async function main(args) {
  let tenantJohny = new core.Tenant(args.coreTenant);
  let transactionCheckpoint = tenantJohny.getTransactionCheckpoint(args.iban);
  console.log("Starting from " + transactionCheckpoint);

  let coreAccountStatement = fio.extractCoreAccountStatement(fioAccountStatement);
  let accounts = fio.extractUniqueCoreAccounts(fioAccountStatement);

  await tenantJohny.createMissingAccounts(accounts);
  await tenantJohny.createTransactions(coreAccountStatement.transactions, coreAccountStatement.accountNumber);
}

main({
  "iban": "CZ7920100000002400222233",
  "token": "xxx",
  "coreTenant": "johny"
}).then(() => console.log("FINISH - Everything went well"))
  .catch(error => {
    console.log(error);
    console.log(error.message);
    console.log("There were some unexpected error see above");
  });
