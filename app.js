const fio = require("./modules/fio.js")
const core = require("./modules/core.js")
const log = require("./modules/logger")

async function main(argv) {
  if (!argv || !argv.tenantName || !argv.token) {
    log.error("Run program using npm start <tenant> <fio_token> [wait]")
    return
  }

  log.info(`Running synchronization for tenant ${argv.tenantName}`)
  const tenant = new core.Tenant(argv.tenantName, argv.token)

  const transactionCheckpoint = await tenant.getCheckpointByToken(argv.token)

  const fioAccountStatement = await fio.getFioAccountStatement(argv.token, transactionCheckpoint, argv.wait)
  const coreAccountStatement = fio.toCoreAccountStatement(fioAccountStatement)
  const accounts = fio.extractUniqueCoreAccounts(fioAccountStatement)

  try {
    await tenant.createMissingAccounts(accounts)
  } catch (err) {
    log.error(`Account creation ended with error: ${err}`)
  }

  try {
    await tenant.createTransactions(coreAccountStatement.transactions, coreAccountStatement.accountNumber, argv.token)
  } catch (err) {
    log.error(`Transaction creation ended with error: ${err}`)
  }
}

main({
  "tenantName": process.argv[2],
  "token": process.argv[3],
  "wait": process.argv[4] && process.argv[5] === "wait"
}).catch((error) => {
  log.error("Synchronization failed, exception:\n" + error.stack)
})
