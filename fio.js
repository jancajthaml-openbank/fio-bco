let options = require("config").get("fio");
let axios = require("axios");

function getAccount(fioTransaction) {
  return (fioTransaction.column2 && fioTransaction.column2.value) || "FIO";
}

function getDebitAccount(transaction, accountIban) {
  if (transaction.column1.value > 0) {
    return (transaction.column2 && transaction.column2.value) || "FIO";
  } else {
    return accountIban;
  }
}

function getCreditAccount(transaction, accountIban) {
  if (transaction.column1.value < 0) {
    return (transaction.column2 && transaction.column2.value) || "FIO";
  } else {
    return accountIban;
  }
}

function fioTransaction2CoreTransfer(fioTransaction, accountIban, accountCurrency) {
  return {
    "id": fioTransaction.column22.value.toString(),
    "credit": getCreditAccount(fioTransaction, accountIban),
    "debit": getDebitAccount(fioTransaction, accountIban),
    "amount": fioTransaction.column1.value.toString(),
    "currency": accountCurrency
  };
}

function fioTransactions2CoreTransactions(transactions, fioTransaction, accountIban, accountCurrency) {
  let transactionId = fioTransaction.column17.value;
  if (transactions[transactionId]) {
    transactions[transactionId].transfers.push(fioTransaction2CoreTransfer(fioTransaction, accountIban, accountCurrency));
  } else {
    transactions[transactionId] = {
      "blame": "fio-sync",
      "transfers": [fioTransaction2CoreTransfer(fioTransaction, accountIban, accountCurrency)]
    };
  }
  return transactions;
}

function extractCoreAccountStatement(fioAccountStatement) {
  let transactions = fioAccountStatement.accountStatement.transactionList.transaction
    .reduce((transactions, fioTransaction) => fioTransactions2CoreTransactions(
      transactions,
      fioTransaction,
      fioAccountStatement.accountStatement.info.iban,
      fioAccountStatement.accountStatement.info.currency), {});

  return {
    "accountNumber": fioAccountStatement.accountStatement.info.iban,
    "idTransactionTo": fioAccountStatement.accountStatement.info.idTo,
    "transactions": Object.keys(transactions).map((transactionId) => {
      let transaction = transactions[transactionId];
      transaction.id = transactionId;
      return transaction;
    })
  }
}

function extractUniqueCoreAccounts(fioAccountStatement) {
  let coreAccounts = fioAccountStatement.accountStatement.transactionList.transaction
    .map((fioTransaction) => {
      return {
        "accountNumber": getAccount(fioTransaction),
        "currency": fioAccountStatement.accountStatement.info.currency,
        "isBalanceCheck": false
      };
    })
    .filter((coreAccount, index, coreAccounts) => {
      let findIndex = coreAccounts.findIndex((acc, accIndex) => {
        return acc.accountNumber === coreAccount.accountNumber && acc.currency === coreAccount.currency;
      });
      return index === findIndex;
    });
  coreAccounts
    .push({
      "accountNumber": fioAccountStatement.accountStatement.info.iban,
      "currency": fioAccountStatement.accountStatement.info.currency,
      "isBalanceCheck": false
    });
  return coreAccounts;
}

let sleep = ms => new Promise(resolve => setTimeout(resolve, ms));

async function getFioAccountStatement(token, idTransactionFrom, wait) {
  if (!idTransactionFrom) {
    await axios.get(options.apiUrl + "/set-last-date/" + token +  "/1900-01-01/");
  } else {
    await axios.get(options.apiUrl + "/set-last-id/" + token +  "/" + idTransactionFrom + "/");
  }
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
  extractCoreAccountStatement,
  extractUniqueCoreAccounts,
  getFioAccountStatement
};