var fioAccountStatement = require("./fio-response-example.json");
var fio = require("./fio.js");

var accountStatement = fio.normalizeAccountStatement(fioAccountStatement);

console.log(JSON.stringify(accountStatement, null, 2));