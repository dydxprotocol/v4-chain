'use strict';

const DBError = require('./DBError');

class DataError extends DBError {}

module.exports = DataError;