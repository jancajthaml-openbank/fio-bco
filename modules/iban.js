const log = require("./logger.js")

const mod97 = (digitString) => {
  let m = 0
  for (var i = 0; i < digitString.length; ++i) {
    m = (((m * 10) + parseInt(digitString.charAt(i), 10)) % 97)
  }
  return m
}

const calculateCzech = (bankCode, accountId) => {
  if (!bankCode) {
    return !accountId ? undefined : accountId
  }

  // canonise input
  let account = accountId.replace(/-/g,'')
  // accountNumber of length 16
  let number = `0000000000000000${account}`.slice(-16)
  // bankCode of length 4
  let code = `0000${bankCode}`.slice(-4)
  // country code for "Czech Republic"
  let country = "CZ"
  // country code converted to digits
  let cc = "123500"
  // bban checksum mod 97
  let checksum = `00${(98 - mod97(code + number + cc))}`.slice(-2)

  return `${country}${checksum}${code}${number}`
}

module.exports = {
  calculateCzech
}
