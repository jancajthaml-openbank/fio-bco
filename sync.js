let jsonfile = require("jsonfile-promised");

async function setTransactionCheckpoint(fn, tenantName, accountNumber, idTransactionTo) {
  let checkpoints;
  try {
    checkpoints = await jsonfile.readFile(fn);
  } catch (error) {
    if (error.code === "ENOENT") {
      console.log("Database " + fn + " will be created for the first time");
      checkpoints = {};
    }
    else throw error;
  }

  if (checkpoints[tenantName]) {
    checkpoints[tenantName][accountNumber] = { idTransactionTo };
  } else {
    checkpoints[tenantName] = {
      [accountNumber] : {
        idTransactionTo
      }
    }
  }

  await jsonfile.writeFile(fn, checkpoints);
}

async function getTransactionCheckpoint(fn, tenantName, accountNumber) {
  try {
    let checkpoints = await jsonfile.readFile(fn);
    if (checkpoints && checkpoints[tenantName] && checkpoints[tenantName][accountNumber]) {
      return checkpoints[tenantName][accountNumber].idTransactionTo;
    } else {
      return null;
    }
  } catch (error) {
    if (error.code === "ENOENT") return null;
    else throw error;
  }
}

module.exports = {
  setTransactionCheckpoint,
  getTransactionCheckpoint
};