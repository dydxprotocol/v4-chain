'use strict';

const DBError = require('../../../errors/DBError');
const errorCodes = require('../../../errorCodes').postgres;

module.exports = {
  error: DBError,

  parse: (err) => {
    if (typeof err.code === 'string'
      && err.code.length === 5
      && errorCodes.has(err.code)
      && 'internalQuery' in err
      && 'table' in err) {

      return {
        nativeError: err,
        client: 'postgres'
      };
    }

    return null;
  },

  subclassParsers: [
    require('./ConstraintViolationError/parser'),
    require('./DataError/parser')
  ]
};