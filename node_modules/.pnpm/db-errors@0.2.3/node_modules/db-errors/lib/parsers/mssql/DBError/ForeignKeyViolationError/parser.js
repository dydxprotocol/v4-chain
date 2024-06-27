const ForeignKeyViolationError = require('../../../../errors/ForeignKeyViolationError');
const { isCode } = require('../util');

const INSERT_UPDATE_REGEX = /The (?:INSERT|UPDATE) statement conflicted with the FOREIGN KEY constraint "(.+)". The conflict occurred in database "(.+)", table "(.+)\.(.+)", column '(.+)'./;
const DELETE_REGEX = /The DELETE statement conflicted with the REFERENCE constraint "(.+)". The conflict occurred in database "(.+)", table "(.+)\.(.+)", column '(.+)'./;

module.exports = {
  error: ForeignKeyViolationError,

  parse: (err) => {
    if (isCode(err, 16, 547)) {
      const insertUpdateMatch = INSERT_UPDATE_REGEX.exec(err.originalError.message);

      if (insertUpdateMatch) {
        return {
          table: insertUpdateMatch[4],
          schema: insertUpdateMatch[3],
          constraint: insertUpdateMatch[1],
        };
      }

      const deleteMatch = DELETE_REGEX.exec(err.originalError.message);

      if (deleteMatch) {
        return {
          table: deleteMatch[4],
          schema: deleteMatch[3],
          constraint: deleteMatch[1],
        };
      }
    }

    return null;
  },

  subclassParsers: []
};