function getDebitAccount(transaction, accountIban) {
  if (transaction.column1.value > 0) {
    return (transaction.column2 && transaction.column2.value) || "FIO internal account";
  } else {
    return accountIban;
  }
}

function getCreditAccount(transaction, accountIban) {
  if (transaction.column1.value < 0) {
    return (transaction.column2 && transaction.column2.value) || "FIO internal account";
  } else {
    return accountIban;
  }
}

function fioTransaction2CoreTransfer(fioTransction, accountIban) {
  return {
    "transactionId": fioTransction.column17.value,
    "transferId": fioTransction.column22.value,
    "amount": fioTransction.column1.value,
    "debit": getDebitAccount(fioTransction, accountIban),
    "credit": getCreditAccount(fioTransction, accountIban),
    "valueDate": fioTransction.column0.value
  };
}

function accumulateCoreTransfers2CoreTransaction(transactions, transfer) {
  if (transactions[transfer.transactionId]) {
    transactions[transfer.transactionId].push(transfer);
  } else {
    transactions[transfer.transactionId] = [transfer];
  }
  return transactions;
}

// TODO: rethink, maybe extracting core transactions is enough
function normalizeAccountStatement(fioAccountStatement) {
  let transactions = fioAccountStatement.accountStatement.transactionList.transaction
    .map((fioTransaction) => fioTransaction2CoreTransfer(fioTransaction, fioAccountStatement.accountStatement.info.iban))
    .reduce(accumulateCoreTransfers2CoreTransaction, {});
  let transactionsArray = Object.keys(transactions).map((transactionId) => transactions[transactionId]);

  return {
    "iban": fioAccountStatement.accountStatement.info.iban,
    "currency": fioAccountStatement.accountStatement.info.currency,
    "transactions": transactionsArray
  };
}

// TODO: Add method for extracting unique core accounts

module.exports = {
  normalizeAccountStatement: normalizeAccountStatement
};