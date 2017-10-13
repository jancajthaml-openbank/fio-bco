//let fioAccountStatement = require("./fio-response-example.json");
let fio = require("./fio.js");
let core = require("./core.js");

async function main(argv) {
  if (!argv || !argv.tenantName || !argv.accountNumber || !argv.token) {
    console.log("Run program using npm start <tenant_name> <tenant_accountIban> <fio_token> [wait]");
    return;
  }

  let tenant = new core.Tenant(argv.tenantName);
  let transactionCheckpoint = await tenant.getTransactionCheckpoint(argv.accountNumber);

  let fioAccountStatement = await fio.getFioAccountStatement(argv.token, transactionCheckpoint, argv.wait);
  let coreAccountStatement = fio.extractCoreAccountStatement(fioAccountStatement);
  let accounts = fio.extractUniqueCoreAccounts(fioAccountStatement);

  await tenant.createMissingAccounts(accounts);
  await tenant.createTransactions(coreAccountStatement.transactions, coreAccountStatement.accountNumber);
}

main({
  "tenantName": process.argv[2],
  "accountNumber": process.argv[3],
  "token": process.argv[4],
  "wait": process.argv[5] && process.argv[5] === "wait"
}).catch(error => {
    console.log(error);
    console.log(error.message);
    console.log("There were some unexpected error see above");
  });
