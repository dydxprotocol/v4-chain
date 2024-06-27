'use strict';

const DBError = require('./DBError');

class ConstraintViolationError extends DBError {}

module.exports = ConstraintViolationError;