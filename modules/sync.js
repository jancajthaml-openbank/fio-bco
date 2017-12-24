/**
 * Sync is used as simple DB that stores synchronization checkpoints of tenant's
 * accounts in core. The structure of the
 * DB is:
 *
 * {
 *   "tenantName": {
 *     "accountNumber1": {
 *       "idTransactionTo": #####
 *     },
 *     "accountNumber2": {
 *       "idTransactionTo": #####
 *     }
 *   },
 *   "anotherTenant": ...
 * }
 *
 * where idTransactionTo is ID of last transaction that was synced for given
 * tenant and his account.
 */

const jsonfile = require("jsonfile-promised")
const log = require("./logger.js")
const VError = require("verror")

const NoSuchFileException = 'ENOENT'

async function setTransactionCheckpoint(db, tenantName, accountNumber, token, idTransferTo) {
  let checkpoints
  // Todo: add check for parameters, token is optional

  try {
    checkpoints = await jsonfile.readFile(db)
  } catch (error) {
    if (error.code === NoSuchFileException) {
      log.info(`Database ${db} will be created for the first time`)
      checkpoints = {}
    } else {
      // fixme: add verror
      throw error
    }
  }

  if (checkpoints[tenantName]) {
    checkpoints[tenantName][accountNumber] = {
      idTransferTo,
      token
    }
  } else {
    checkpoints[tenantName] = {
      [accountNumber]: {
        idTransferTo,
        token
      }
    }
  }

  await jsonfile.writeFile(db, checkpoints)
}

async function getTransactionCheckpoint(db, tenantName, accountNumber) {
  return await getCheckpoint(db, (checkpoints) =>
    (checkpoints[tenantName] && checkpoints[tenantName][accountNumber])
    ? checkpoints[tenantName][accountNumber].idTransferTo
    : null
  )
}

async function getCheckpoint(db, searchCb) {
  try {
    const checkpoints = await jsonfile.readFile(db)
    return searchCb(checkpoints)
  } catch (err) {
    if (err.code === NoSuchFileException) {
      return null
    }
    throw new VError(err, `Error when reading DB file ${db}`)
  }
}

async function getTransactionCheckpointByToken(db, tenantName, token) {
  return await getCheckpoint(db, (checkpoints) => {
    if (checkpoints[tenantName]) {
      const result = Object.values(checkpoints[tenantName]).find((account) =>
        account.token && (account.token === token)
      )
      return result && result.idTransferTo || null
    }
    return null
  })
}

module.exports = {
  setTransactionCheckpoint,
  getTransactionCheckpoint,
  getTransactionCheckpointByToken
}
