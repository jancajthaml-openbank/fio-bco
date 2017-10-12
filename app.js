//let fioAccountStatement = require("./fio-response-example.json");
let fio = require("./fio.js");
let core = require("./core.js");

async function main(argv) {
  if (!argv || !argv.tenantName || !argv.accountNumber || !argv.token) {
    console.log("Run program using npm start <tenant_name> <tenant_accountIban> <fio_token>");
    return;
  }

  let tenantJohny = new core.Tenant(argv.tenantName);
  let transactionCheckpoint = await tenantJohny.getTransactionCheckpoint(argv.accountNumber);

  let fioAccountStatement = await fio.getFioAccountStatement(argv.token, transactionCheckpoint, true);
  let coreAccountStatement = fio.extractCoreAccountStatement(fioAccountStatement);
  let accounts = fio.extractUniqueCoreAccounts(fioAccountStatement);

  await tenantJohny.createMissingAccounts(accounts);
  await tenantJohny.createTransactions(coreAccountStatement.transactions, coreAccountStatement.accountNumber);
}

main({
  "tenantName": process.argv[2],
  "accountNumber": process.argv[3],
  "token": process.argv[4]
}).catch(error => {
    console.log(error);
    console.log(error.message);
    console.log("There were some unexpected error see above");
  });
