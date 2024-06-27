const ConstraintViolationError = require('../../../../errors/ConstraintViolationError');

module.exports = {
  error: ConstraintViolationError,

  parse: (err) => {
    if (err.code === 'SQLITE_CONSTRAINT') {
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