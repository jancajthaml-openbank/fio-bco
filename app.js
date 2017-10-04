var fioAccountStatement = require("./fio-response-example.json");
var fio = require("./fio.js");

console.log(JSON.stringify(fio.extractCoreTransactions(fioAccountStatement), null, 2));