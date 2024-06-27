const ForeignKeyViolationError = require('../../../../../errors/ForeignKeyViolationError');

const NO_REFERENCED_REGEX = /a foreign key constraint fails \(`(.+)`\.`(.+)`, CONSTRAINT `(.+)` FOREIGN KEY \(`(.+)`\) REFERENCES `(.+)` \(`(.+)`\)/;
const ROW_IS_REFERENCED_REGEX = /Cannot delete or update a parent row: a foreign key constraint fails \(`(.+)`\.`(.+)`, CONSTRAINT `(.+)` FOREIGN KEY \(`(.+)`\) REFERENCES `(.+)` \(`(.+)`\)/;

module.exports = {
  error: ForeignKeyViolationError,

  parse: (err) => {
    let match = null;

    if (
      err.code === 'ER_NO_REFERENCED_ROW' ||
      err.code === 'ER_ROW_IS_REFERENCED'
    ) {
      // For these variants, there is no table or constraint information available.
      // These seem to be mostly thrown on mysql 8 when the db user doesn't have
      // priviledges to the parent table. There seems to be a bug however that
      // causes these generic errors to be thrown even in some cases where the
      // user DOES have priviledges.
      return {
        table: null,
        constraint:null
      };
    }

    if (err.code === 'ER_NO_REFERENCED_ROW_2') {
      match = NO_REFERENCED_REGEX.exec(err.sqlMessage);
    } else if (err.code === 'ER_ROW_IS_REFERENCED_2') {
      match = ROW_IS_REFERENCED_REGEX.exec(err.sqlMessage);
    }

    if (!match) {
      return null;
    }

    return {
      table: match[2],
      constraint: match[3]
    };
  },

  subclassParsers: []
};
