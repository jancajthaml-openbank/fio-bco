/**
 * The main purpose of this module is to transform fio account statement to core account statement. The account
 * statement contains all transactions that was realized in given time period. To clarify terminology we have
 * transaction and transfer. Transaction contains 1 - N transfers, each transfer belongs under one particular
 * transaction. In fio terminology transaction is called "Pokyn" and transfer is called "Pohyb".
 */

let options = require("config").get("fio");
let axios = require("axios");

function extractCounterPartAccountNumber(fioTransfer) {
  return (fioTransfer.column2 && fioTransfer.column2.value) || "FIO";
}

function extractAmount(fioTransfer) {
  return fioTransfer.column1.value;
}

function extractDebitAccountNumber(fioTransfer, mainAccountNumber) {
  if (extractAmount(fioTransfer) > 0) {
    return extractCounterPartAccountNumber(fioTransfer);
  } else {
    return mainAccountNumber;
  }
}

function extractCreditAccountNumber(fioTransfer, mainAccountNumber) {
  if (extractAmount(fioTransfer) < 0) {
    return extractCounterPartAccountNumber(fioTransfer);
  } else {
    return mainAccountNumber;
  }
}

function extractTransferValueDate(fioTransfer) {
  let part = fioTransfer.column0.value.substring(0, fioTransfer.column0.value.indexOf('+'));
  let stringDate = part + "T00:00:00" + fioTransfer.column0.value.substring(fioTransfer.column0.value.indexOf('+'));
  return new Date(stringDate);
}

function extractTransferId(fioTransfer) {
  return fioTransfer.column22.value;
}

function extractTransactionId(fioTransfer) {
  return fioTransfer.column17.value;
}

function extractMainAccountNumber(fioAccountStatement) {
  return fioAccountStatement.accountStatement.info.iban;
}

function extractMainAccountCurrency(fioAccountStatement) {
  return fioAccountStatement.accountStatement.info.currency;
}

function fioTransferToCoreTransfer(fioTransfer, mainAccountNumber, mainAccountCurrency) {
  return {
    "id": extractTransferId(fioTransfer).toString(),
    "valueDate": extractTransferValueDate(fioTransfer).toISOString(),
    "credit": extractCreditAccountNumber(fioTransfer, mainAccountNumber),
    "debit": extractDebitAccountNumber(fioTransfer, mainAccountNumber),
    "amount": Math.abs(extractAmount(fioTransfer)).toString(),
    "currency": mainAccountCurrency
  };
}

function fioTransfersToCoreTransactions(fioTransfers, mainAccountNumber, mainAccountCurrency) {
  let result = fioTransfers
    .reduce((coreTransactions, fioTransfer) => {
      let transactionId = extractTransactionId(fioTransfer);
      if (coreTransactions[transactionId]) {
        coreTransactions[transactionId].transfers.push(fioTransferToCoreTransfer(fioTransfer, mainAccountNumber, mainAccountCurrency));
      } else {
        coreTransactions[transactionId] = {
          "blame": "fio-sync",
          "transfers": [fioTransferToCoreTransfer(fioTransfer, mainAccountNumber, mainAccountCurrency)]
        };
      }
      return coreTransactions;
    }, {});

  // Return as array
  return Object.keys(result).map((transactionId) => {
    let transaction = result[transactionId];
    transaction.id = transactionId;
    return transaction;
  });
}

function toCoreAccountStatement(fioAccountStatement) {
  return {
    "accountNumber": extractMainAccountNumber(fioAccountStatement),
    "transactions": fioTransfersToCoreTransactions(
      fioAccountStatement.accountStatement.transactionList.transaction,
      extractMainAccountNumber(fioAccountStatement),
      extractMainAccountCurrency(fioAccountStatement)
    )
  }
}

function extractUniqueCoreAccounts(fioAccountStatement) {
  let coreAccounts = fioAccountStatement.accountStatement.transactionList.transaction
    .filter((fioTransfer, currIndex, fioTransfers) => {
      let currAccountNumber = extractCounterPartAccountNumber(fioTransfer);
      let foundIndex = fioTransfers.findIndex((fioTransfer) => {
        return currAccountNumber === extractCounterPartAccountNumber(fioTransfer);
      });
      return foundIndex === currIndex;
    })
    .map((fioTransfer) => {
      return {
        "accountNumber": extractCounterPartAccountNumber(fioTransfer),
        "currency": extractMainAccountCurrency(fioAccountStatement),
        "isBalanceCheck": false
      };
    });

  // Add main account
  coreAccounts
    .push({
      "accountNumber": extractMainAccountNumber(fioAccountStatement),
      "currency": extractMainAccountCurrency(fioAccountStatement),
      "isBalanceCheck": false
    });

  return coreAccounts;
}

let sleep = ms => new Promise(resolve => setTimeout(resolve, ms));

async function setLastTransaction(token, idLastTransaction) {
  if (!idLastTransaction) {
    await axios.get(options.apiUrl + "/set-last-date/" + token +  "/1900-01-01/");
  } else {
    await axios.get(options.apiUrl + "/set-last-id/" + token +  "/" + idLastTransaction + "/");
  }
}

async function getFioAccountStatement(token, idTransactionFrom, wait) {
  await setLastTransaction(token, idTransactionFrom);
  try {
    let response = await axios.get(options.apiUrl + "/last/" + token + "/transactions.json");
    return response.data;
  } catch (error) {
    if (error.response.status === 409) {
      if (wait) {
        console.log("Request to FIO for transactions is too early - waiting 20 seconds ...");
        await sleep(1000 * 20);
        let response = await axios.get(options.apiUrl + "/last/" + token + "/transactions.json");
        return response.data;
      } else {
        console.log("Request to FIO for transactions is too early");
        throw error;
      }
    }
    throw error;
  }
}

module.exports = {
  extractCoreAccountStatement: toCoreAccountStatement,
  extractUniqueCoreAccounts,
  getFioAccountStatement
};