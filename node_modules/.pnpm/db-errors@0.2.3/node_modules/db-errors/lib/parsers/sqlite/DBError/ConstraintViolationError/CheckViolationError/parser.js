const CheckViolationError = require('../../../../../errors/CheckViolationError');

const REGEX = /SQLITE_CONSTRAINT: CHECK constraint failed/;

module.exports = {
  error: CheckViolationError,

  parse: (err) => {
    const match = REGEX.exec(err.message);

    if (!match) {
      return null;
    }

    // No way to extract anything reliably.
    return {};
  },

  subclassParsers: []
};