'use strict';

const DBError = require('../../../errors/DBError');
const errorCodes = require('../../../errorCodes').mssql;
const { getCode } = require('./util');

module.exports = {
  error: DBError,

  parse: (err) => {
    if (err.originalError && err.code === 'EREQUEST' && errorCodes.has(getCode(err))) {
      return {
        nativeError: err.originalError,
        client: 'mssql'
      };
    }

    return null;
  },

  subclassParsers: [
    require('./CheckViolationError/parser'),
    require('./DataError/parser'),
    require('./ForeignKeyViolationError/parser'),
    require('./NotNullViolationError/parser'),
    require('./UniqueViolationError/parser'),
  ]
};