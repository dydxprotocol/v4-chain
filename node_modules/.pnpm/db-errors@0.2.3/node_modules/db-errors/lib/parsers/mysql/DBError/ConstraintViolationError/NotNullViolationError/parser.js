const NotNullViolationError = require('../../../../../errors/NotNullViolationError');

const COLUMN_REGEX_BAD_NULL_ERROR = /Column '(.+)' cannot be null/;
const COLUMN_REGEX_NO_DEFAULT_FOR_FIELD = /Field '(.+)' doesn't have a default value/;

module.exports = {
  error: NotNullViolationError,

  parse: (err) => {
    let columnMatch = null;

    if (err.code === 'ER_BAD_NULL_ERROR') {
      columnMatch = COLUMN_REGEX_BAD_NULL_ERROR.exec(err.sqlMessage);
    } else if (err.code === 'ER_NO_DEFAULT_FOR_FIELD') {
      columnMatch = COLUMN_REGEX_NO_DEFAULT_FOR_FIELD.exec(err.sqlMessage);
    }

    if (!columnMatch) {
      return null;
    }

    // No way to reliably get `table` from mysql error.
    return {
      column: columnMatch[1]
    };
  },

  subclassParsers: []
};
