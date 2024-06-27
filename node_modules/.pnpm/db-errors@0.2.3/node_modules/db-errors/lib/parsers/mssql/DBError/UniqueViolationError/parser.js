const UniqueViolationError = require('../../../../errors/UniqueViolationError');
const { isCode } = require('../util');

const UNIQUE_INDEX_REGEX = /Cannot insert duplicate key row in object '(.+)\.(.+)' with unique index '(.+)'. The duplicate key value is (.+)./;
const UNIQUE_CONSTRAINT_REGEX = /Violation of UNIQUE KEY constraint '(.+)'. Cannot insert duplicate key in object '(.+)\.(.+)'. The duplicate key value is \((.+)\)/;

// 2601 - Violation in unique index
// 2627 - Violation in unique constraint (although it is implemented using unique index)

module.exports = {
  error: UniqueViolationError,

  parse: (err) => {
    if (isCode(err, 14, 2627)) {
      const constraintMatch = UNIQUE_CONSTRAINT_REGEX.exec(err.originalError.message);

      if (!constraintMatch) {
        return null;
      }

      return {
        table: constraintMatch[3],
        schema: constraintMatch[2],
        constraint: constraintMatch[1]
      };
    }

    // TODO: this case is missing a test
    if (isCode(err, 14, 2601)) {
      const indexMatch = UNIQUE_INDEX_REGEX.exec(err.originalError.message);

      if (!indexMatch) {
        return null;
      }

      return {
        table: indexMatch[2],
        constraint: indexMatch[3],
        schema: indexMatch[1],
      };
    }

    return null;
  },

  subclassParsers: []
};