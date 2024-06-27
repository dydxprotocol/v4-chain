const UniqueViolationError = require('../../../../../errors/UniqueViolationError');

const UNIQUE_COLUMNS_REGEX = /Key \((.+)\)=\(.+\) already exists/;

module.exports = {
  error: UniqueViolationError,

  parse: (err) => {
    if (err.code !== '23505') {
      return null;
    }

    const colsMatch = UNIQUE_COLUMNS_REGEX.exec(err.detail);

    if (!colsMatch) {
      return null;
    }

    const cols = colsMatch[1]
      .split(',')
      .map(it => it.trim())
      .map(it => stripQuotes(it));

    return {
      table: err.table,
      columns: cols,
      constraint: err.constraint
    };
  },

  subclassParsers: []
};

function stripQuotes(it) {
  if (it[0] === '"' && it[it.length - 1] === '"') {
    return it.substr(1, it.length - 2);
  } else {
    return it;
  }
}