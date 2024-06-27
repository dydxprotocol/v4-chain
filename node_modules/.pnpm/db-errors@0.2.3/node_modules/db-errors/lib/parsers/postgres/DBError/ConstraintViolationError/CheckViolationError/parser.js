const CheckViolationError = require('../../../../../errors/CheckViolationError');

module.exports = {
  error: CheckViolationError,

  parse: (err) => {
    if (err.code === '23514') {
      return {
        table: err.table,
        constraint: err.constraint
      };
    } else {
      return null;
    }
  },

  subclassParsers: []
};