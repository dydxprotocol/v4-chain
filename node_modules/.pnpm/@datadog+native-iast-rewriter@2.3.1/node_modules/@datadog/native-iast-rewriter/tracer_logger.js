'use strict'

let defaultLogger = console || {}
function setLogger (logger) {
  if (logger) {
    defaultLogger = logger
  }
}

function log (level, msg) {
  const logFn = defaultLogger[level.toLowerCase()]
  if (logFn) {
    logFn(msg)
  }
}

module.exports = {
  setLogger,
  log
}
