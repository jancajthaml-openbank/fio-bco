/**
 * The main purpose of this module is to transform fio account statement to
 * core account statement. The account statement contains all transactions that
 * was realized in given time period. To clarify terminology we have transaction
 * and transfer. Transaction contains 1 - N transfers, each transfer belongs
 * under one particular transaction. In Fio terminology transaction is called
 * "Pokyn" and transfer is called "Pohyb".
 */

const axios = require('axios')
const log = require('./logger')
const VError = require('verror')

const options = require("config").get("fio")

const sleep = ms => new Promise(andThen => setTimeout(andThen, ms))

const extractCounterPartAccountNumber = _ =>
  ((_.column2 && _.column2.value) || "FIO")

const extractAmount = _ =>
  Number(_.column1.value)

const extractAbsAmount = _ =>
  Math.abs(Number(_.column1.value))

const extractDebitAccountNumber = (fioTransfer, mainAccountNumber) =>
  (extractAmount(fioTransfer) > 0)
    ? extractCounterPartAccountNumber(fioTransfer)
    : mainAccountNumber

const extractCreditAccountNumber = (fioTransfer, mainAccountNumber) => 
  (extractAmount(fioTransfer) < 0)
    ? extractCounterPartAccountNumber(fioTransfer)
    : mainAccountNumber

const extractTransferValueDate = _ => {
  let sDate = _.column0.value,
      idx = sDate.indexOf('+'),
      date = new Date((idx == -1)
        ? `${sDate}T00:00:00+0000`
        : `${sDate.substring(0, idx)}T00:00:00${sDate.substring(idx)}`)

  return date.toISOString()
}

const extractTransferId = _ =>
  String(_.column22.value)

const extractTransactionId = _ =>
  _.column17.value

const extractMainAccountNumber = _ =>
  _.accountStatement.info.iban

const extractMainAccountCurrency = _ =>
  _.accountStatement.info.currency

const fioTransferToCoreTransfer = (fioTransfer, mainAccountNumber, mainAccountCurrency) => ({
  "id": extractTransferId(fioTransfer),
  "valueDate": extractTransferValueDate(fioTransfer),
  "credit": extractCreditAccountNumber(fioTransfer, mainAccountNumber),
  "debit": extractDebitAccountNumber(fioTransfer, mainAccountNumber),
  "amount": String(extractAbsAmount(fioTransfer)),
  "currency": mainAccountCurrency
})

function fioTransfersToCoreTransactions(fioTransfers, mainAccountNumber, mainAccountCurrency) {
  const result = fioTransfers
    .reduce((coreTransactions, fioTransfer) => {
      const transactionId = extractTransactionId(fioTransfer)

      if (coreTransactions[transactionId]) {
        coreTransactions[transactionId].transfers.push(fioTransferToCoreTransfer(fioTransfer, mainAccountNumber, mainAccountCurrency))
      } else {
        coreTransactions[transactionId] = {
          "blame": "fio-bco", // FIXME to config
          "transfers": [fioTransferToCoreTransfer(fioTransfer, mainAccountNumber, mainAccountCurrency)]
        }
      }
      return coreTransactions
    }, {})

  // Return as array
  return Object.keys(result).map(transactionId => {
    // FIXME possible undefined
    const transaction = result[transactionId]
    transaction.id = transactionId
    return transaction
  })
}

const toCoreAccountStatement = fioAccountStatement => ({
  "accountNumber": extractMainAccountNumber(fioAccountStatement),
  "transactions": fioTransfersToCoreTransactions(
    fioAccountStatement.accountStatement.transactionList.transaction,
    extractMainAccountNumber(fioAccountStatement),
    extractMainAccountCurrency(fioAccountStatement)
  )
})

function extractUniqueCoreAccounts(fioAccountStatement) {
  const mainCurrency = extractMainAccountCurrency(fioAccountStatement)

  const coreAccounts = fioAccountStatement.accountStatement.transactionList.transaction
    .filter((fioTransfer, currIndex, fioTransfers) => {
      const currAccountNumber = extractCounterPartAccountNumber(fioTransfer)
      const foundIndex = fioTransfers.findIndex(fioTransferCmp =>
        currAccountNumber === extractCounterPartAccountNumber(fioTransferCmp)
      )
      return foundIndex === currIndex
    })
    .map(extractCounterPartAccountNumber)

  // Add main account
  coreAccounts.push(extractMainAccountNumber(fioAccountStatement))

  return coreAccounts.map(accountNumber => ({
    "accountNumber": accountNumber,
    "currency": mainCurrency,
    "isBalanceCheck": false
  }))
}

async function setLastTransaction(token, idLastTransaction) {
  try {
    return await axios.get(idLastTransaction
      ? `${options.apiUrl}/set-last-id/${token}/${idLastTransaction}/`
      : `${options.apiUrl}/set-last-date/${token}/1900-01-01/`
    )
  } catch (err) {
    throw new VError(err, "Request to FIO api failed")
  }
}

async function getLastTransactions(token, retry) {
  try {
    return await axios.get(`${options.apiUrl}/last/${token}/transactions.json`)
  } catch (err) {
    if (err.response && err.response.status === 409) {
      if (retry) {
        // FIXME 20 secs to constant
        log.warn("Request to FIO for transactions is too early - waiting 20 seconds ...")
        await sleep(20 * 1000) // FIXME to config
        return await getLastTransactions(token, false)
      } else {
        throw new VError(err, "FIO transaction api unavailable, you have to wait 20 seconds between calls")
      }
    } else {
      throw new VError(err, "Request to FIO api failed")
    }
  }
}

async function getFioAccountStatement(token, idTransactionFrom, wait) {
  await setLastTransaction(token, idTransactionFrom)
  const response = await getLastTransactions(token, wait)
  log.info(`Loaded FIO account statement for account ${response.data.accountStatement.info.iban}`)
  return response.data
}

module.exports = {
  toCoreAccountStatement,
  extractUniqueCoreAccounts,
  getFioAccountStatement
}
