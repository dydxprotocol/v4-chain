const DataError = require('../../../../errors/DataError');
const { getSqlStateClass } = require('../../../../utils/sqlState');

module.exports = {
  error: DataError,

  parse: (err) => {
    if (getSqlStateClass(err.code) === '22') {
      return {};
    } else {
      return null;
    }
  },

  subclassParsers: []
};