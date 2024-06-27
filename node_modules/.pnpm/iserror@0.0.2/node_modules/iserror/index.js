
/**
 * Expose `isError`.
 */

module.exports = isError;

/**
 * Test whether `value` is error object.
 *
 * @param {*} value
 * @returns {boolean}
 */

function isError(value) {
  switch (Object.prototype.toString.call(value)) {
    case '[object Error]': return true;
    case '[object Exception]': return true;
    case '[object DOMException]': return true;
    default: return value instanceof Error;
  }
}
