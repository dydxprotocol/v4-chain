const CheckViolationError = require('../../../../errors/CheckViolationError');
const { isCode } = require('../util');

const REGEX = /The (INSERT|UPDATE) statement conflicted with the CHECK constraint "(.+)". The conflict occurred in database "(.+)", table "(.+)", column '(.+)'./

module.exports = {
  error: CheckViolationError,

  parse: (err) => {
    if (isCode(err, 16, 547)) {
      const match = REGEX.exec(err.originalError.message);

      if (!match) {
        return null;
      }

      return {
        table: match[3],
        constraint: match[1]
      };
    } else {
      return null;
    }
  },

  subclassParsers: []
};