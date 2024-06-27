const DBError = require('../../../errors/DBError');

module.exports = {
  error: DBError,

  parse: (err) => {
    if (typeof err.code === 'string'
      && err.code.startsWith('SQLITE_')
      && typeof err.errno === 'number') {

      return {
        nativeError: err,
        client: 'sqlite'
      };
    }

    return null;
  },

  subclassParsers: [
    require('./ConstraintViolationError/parser')
  ]
};