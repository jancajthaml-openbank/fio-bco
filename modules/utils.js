const getMax = (a, b) => a > b ? a : b

const partition = (list, chunkSize) => {
  let groups = [], i = 0, j = chunkSize, len = list.length

  while (i < len) {
    groups.push(list.slice(i, j))

    i = j
    j += chunkSize
  }

  return groups
}

const parallelize = (items, partitionSize, processItem, andThen) => {

  let withThen = () => partition(items, partitionSize)
    .map((bulk) => Promise.all(bulk.map(processItem)).then(andThen))

  let withPass = () => partition(items, partitionSize)
    .map((bulk) => bulk.map(processItem))
    .reduce((a, b) => a.concat(b), [])

  return Promise.all(andThen ? withThen() : withPass())
}

const sleep = (ms) => new Promise((andThen) => setTimeout(andThen, ms))

const parseDate = (input) => {
  let idx = input.indexOf('+')

  return new Date((idx === -1)
    ? `${input}T00:00:00+0000`
    : `${input.substring(0, idx)}T00:00:00${input.substring(idx)}`)  
}

module.exports = {
  getMax,
  parallelize,
  sleep,
  parseDate
}
