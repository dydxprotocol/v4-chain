const NotNullViolationError = require('../../../../../errors/NotNullViolationError');

module.exports = {
  error: NotNullViolationError,

  parse: (err) => {
    if (err.code === '23502') {
      return {
        table: err.table,
        column: err.column
      };
    } else {
      return null;
    }
  },

  subclassParsers: []
};