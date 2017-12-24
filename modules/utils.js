const getMax = (a, b) => +a > +b ? a : b

// FIXME differentiate between success and failure, pass only successfull to andThen
const parallelize = (items, processItem, andThen) => andThen
  ? Promise.all(items.map(processItem)).then(andThen)
  : Promise.all(items.map(processItem))

const sleep = (ms) => new Promise((andThen) => setTimeout(andThen, ms))

const elapsedTime = (start) => {
  var elapsed = process.hrtime(start)[1] / 1000000;
  return process.hrtime(start)[0] + " s, " + elapsed.toFixed(3) + " ms"
}

const parseDate = (input) => {
  let idx = input.indexOf("+")

  return new Date((idx === -1)
    ? `${input}T00:00:00+0000`
    : `${input.substring(0, idx)}T00:00:00${input.substring(idx)}`)
}

module.exports = {
  getMax,
  parallelize,
  sleep,
  parseDate,
  elapsedTime
}
