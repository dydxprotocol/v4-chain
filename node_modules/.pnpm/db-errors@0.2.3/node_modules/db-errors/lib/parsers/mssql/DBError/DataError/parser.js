const DataError = require('../../../../errors/DataError');
const { isCode } = require('../util');

module.exports = {
  error: DataError,

  parse: (err) => {
    // 241 - Conversion failed when converting date and/or time from character string.
    // 242 - The conversion of a nvarchar data type to a datetime data type resulted in an out-of-range value.
    // 245 - Conversion failed when converting the nvarchar value 'lol' to data type int.
    // 8152 - String or binary data would be truncated.
    if (isCode(err, 16, 241) || isCode(err, 16, 242) || isCode(err, 16, 245) || isCode(err, 16, 8152)) {
      return {};
    } else {
      return null;
    }
  },

  subclassParsers: []
};