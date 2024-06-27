const DataError = require('../../../../errors/DataError');
const { getSqlStateClass } = require('../../../../utils/sqlState');

const ERROR_CODES = [
  'ER_DATA_TOO_LONG',
  'ER_TRUNCATED_WRONG_VALUE',
  'ER_TRUNCATED_WRONG_VALUE_FOR_FIELD'
];

module.exports = {
  error: DataError,

  parse: (err) => {
    // MySQL mainly uses the SQLSTATE codes, but some errors don't have
    // an SQLSTATE equivalent.
    if (getSqlStateClass(err.sqlState) === '22' || ERROR_CODES.includes(err.code)) {
      return {};
    } else {
      return null;
    }
  },

  subclassParsers: []
};