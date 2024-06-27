'use strict';

const ConstraintViolationError = require('./ConstraintViolationError');

class CheckViolationError extends ConstraintViolationError {

  constructor(args) {
    super(args);

    this.table = args.table;
    this.constraint = args.constraint;
  }
}

module.exports = CheckViolationError;