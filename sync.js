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

async function getTransactionCheckpoint(fn, accountNumber) {
  try {
    let checkpoints = await jsonfile.readFile(fn);
    if (checkpoints && checkpoints[accountNumber]) {
      return checkpoints[accountNumber].idTransactionTo;
    } else {
      return null;
    }
  } catch (error) {
    if (error.code === "ENOENT") return null;
    else throw error;
  }
}

module.exports = {
  updateTransactionCheckpoint,
  getTransactionCheckpoint
};