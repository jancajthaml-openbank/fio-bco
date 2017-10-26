const fio = require("./modules/fio.js");
const core = require("./modules/core.js");
const log = require("./modules/logger");

async function main(argv) {
  if (!argv || !argv.tenantName || !argv.accountNumber || !argv.token) {
    log.error("Run program using npm start <tenant_name> <tenant_accountIban> <fio_token> [wait]");
    return;
  }

  log.info("Running synchronization for tenant/account" + argv.tenantName + "/" + argv.accountNumber);
  const tenant = new core.Tenant(argv.tenantName);
  const transactionCheckpoint = await tenant.getTransactionCheckpoint(argv.accountNumber);

  const fioAccountStatement = await fio.getFioAccountStatement(argv.token, transactionCheckpoint, argv.wait);
  const coreAccountStatement = fio.toCoreAccountStatement(fioAccountStatement);
  const accounts = fio.extractUniqueCoreAccounts(fioAccountStatement);

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
