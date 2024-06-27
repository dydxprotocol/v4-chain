const ConstraintViolationError = require('../../../../errors/ConstraintViolationError');
const { getSqlStateClass } = require('../../../../utils/sqlState');

module.exports = {
  error: ConstraintViolationError,

  parse: (err) => {
    if (getSqlStateClass(err.sqlState) === '23' || err.code === 'ER_NO_DEFAULT_FOR_FIELD') {
      return {};
    } else {
      return null;
    }
  },

  subclassParsers: [
    require('./UniqueViolationError/parser'),
    require('./NotNullViolationError/parser'),
    require('./ForeignKeyViolationError/parser')
  ]
};
