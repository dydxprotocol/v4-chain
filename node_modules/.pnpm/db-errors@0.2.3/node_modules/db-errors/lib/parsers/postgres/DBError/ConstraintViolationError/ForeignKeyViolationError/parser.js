const ForeignKeyViolationError = require('../../../../../errors/ForeignKeyViolationError');

module.exports = {
  error: ForeignKeyViolationError,

  parse: (err) => {
    if (err.code === '23503') {
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