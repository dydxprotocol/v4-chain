'use strict';

const ConstraintViolationError = require('./ConstraintViolationError');

class ForeignKeyViolationError extends ConstraintViolationError {

  constructor(args) {
    super(args);

    this.table = args.table;
    this.constraint = args.constraint;
    this.schema = args.schema;
  }
}

module.exports = ForeignKeyViolationError;