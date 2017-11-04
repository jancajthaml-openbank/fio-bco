const jsonfile = require('jsonfile-promised')
const log = require('./logger.js')

async function setTransactionCheckpoint(db, tenantName, accountNumber, idTransactionTo) {
  let checkpoints

  try {
    checkpoints = await jsonfile.readFile(db)
  } catch (error) {
    if (error.code === "ENOENT") {
      log.info(`Database ${db} will be created for the first time`)
      checkpoints = {}
    } else {
      throw error
    }
  }

  if (checkpoints[tenantName]) {
    checkpoints[tenantName][accountNumber] = { idTransactionTo }
  } else {
    checkpoints[tenantName] = {
      [accountNumber]: {
        idTransactionTo
      }
    }
  }

  await jsonfile.writeFile(db, checkpoints)
}

async function getTransactionCheckpoint(db, tenantName, accountNumber) {
  try {
    let checkpoints = await jsonfile.readFile(db)

    return (checkpoints && checkpoints[tenantName] && checkpoints[tenantName][accountNumber])
      ? checkpoints[tenantName][accountNumber].idTransactionTo
      : null
  } catch (error) {
    if (error.code === "ENOENT") {
      return null
    } else {
      throw error
    }
  }
}

module.exports = {
  setTransactionCheckpoint,
  getTransactionCheckpoint
}
