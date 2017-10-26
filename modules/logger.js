const options = require("config").get("logger");
const log = require("winston");

log.level = options.level;

module.exports = log;