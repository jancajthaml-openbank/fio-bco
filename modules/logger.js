const log = require('winston')

log.level = require("config").get("logger").level

module.exports = log