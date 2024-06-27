const ConstraintViolationError = require('../../../../errors/ConstraintViolationError');
const { getSqlStateClass } = require('../../../../utils/sqlState');

module.exports = {
  error: ConstraintViolationError,

  parse: (err) => {
    if (getSqlStateClass(err.code) === '23') {
      return {};
    } else {
      return null;
    }
  },

  subclassParsers: [
    require('./UniqueViolationError/parser'),
    require('./NotNullViolationError/parser'),
    require('./ForeignKeyViolationError/parser'),
    require('./CheckViolationError/parser')
  ]
};