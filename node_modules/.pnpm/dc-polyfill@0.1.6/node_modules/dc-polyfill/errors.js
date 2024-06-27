const util = require('util');

const { ObjectDefineProperties } = require('./primordials.js');

// Port from node core lib/internal/errors.js
class ERR_INVALID_ARG_TYPE extends TypeError {
  constructor(message, actual) {
    super();

    if (actual == null) {
      message += `. Received ${actual}`;
    } else if (typeof actual === 'function' && actual.name) {
      message += `. Received function ${actual.name}`;
    } else if (typeof actual === 'object') {
      if (actual.constructor && actual.constructor.name) {
        message += `. Received an instance of ${actual.constructor.name}`;
      } else {
        const inspected = util.inspect(actual, { depth: -1 });
        message += `. Received ${inspected}`;
      }
    } else {
      let inspected = util.inspect(actual, { colors: false });
      if (inspected.length > 25) {
        inspected = `${inspected.slice(0, 25)}...`;
      }
      message += `. Received type ${typeof actual} (${inspected})`;
    }

    this.code = this.constructor.name;

    ObjectDefineProperties(this, {
      message: {
        value: message,
        enumerable: false,
        writable: true,
        configurable: true,
      },
      toString: {
        value() {
          return `${this.name} [${this.code}]: ${this.message}`;
        },
        enumerable: false,
        writable: true,
        configurable: true,
      },
    });
  }
}

module.exports = {
  ERR_INVALID_ARG_TYPE,
};
