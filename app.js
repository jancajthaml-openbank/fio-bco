let fioAccountStatement = require("./fio-response-example.json");
let fio = require("./fio.js");
let core = require("./core.js");

async function main() {
  let transactions = fio.extractCoreTransactions(fioAccountStatement);
  let accounts = fio.extractUniqueCoreAccounts(fioAccountStatement);

  await core.createMissingAccounts(accounts);
  await core.createTransactions(transactions);
}

main()
  .then(() => console.log("FINISH - Everything went well"))
  .catch(error => {
    console.log(error);
    console.log("There were some unexpected errors see above")
  });

