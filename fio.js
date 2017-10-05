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

function extractCoreTransactions(fioAccountStatement) {
  let transactions = fioAccountStatement.accountStatement.transactionList.transaction
    .reduce((transactions, fioTransaction) => fioTransactions2CoreTransactions(
      transactions,
      fioTransaction,
      fioAccountStatement.accountStatement.info.iban,
      fioAccountStatement.accountStatement.info.currency), {});
  return Object.keys(transactions).map((transactionId) => transactions[transactionId]);
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

module.exports = {
  extractCoreTransactions,
  extractUniqueCoreAccounts
};