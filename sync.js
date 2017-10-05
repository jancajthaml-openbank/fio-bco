let jsonfile = require("jsonfile-promised");

async function updateTransactionCheckpoint(fn, accountNumber, idTransactionTo) {
  let checkpoints = {};
  try {
    checkpoints = await jsonfile.readFile(fn);
  } catch (error) {
    if (error.code === "ENOENT") console.log("Database " + fn + " will be created for the first time");
    else throw error;
  }
  checkpoints[accountNumber] = { idTransactionTo };
  await jsonfile.writeFile(fn, checkpoints);
}

module.exports = {
  updateTransactionCheckpoint: updateTransactionCheckpoint
};