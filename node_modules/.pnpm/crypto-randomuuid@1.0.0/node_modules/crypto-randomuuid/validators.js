'use strict';

// Polyfill relevant Node.js core validators and error types

const assert = require('assert');
const util = require('util');

// Sorted by a rough estimate on most frequently used entries.
const kTypes = [
  'string',
  'function',
  'number',
  'object',
  // Accept 'Function' and 'Object' as alternative to the lower cased version.
  'Function',
  'Object',
  'boolean',
  'bigint',
  'symbol',
];

class ERR_OPERATION_FAILED extends TypeError {
  constructor(message) {
    super(`Operation failed: ${message}`);
    this.code = this.constructor.name;

    Object.defineProperties(this, {
      toString: {
        value () {
          return `${this.name} [${this.code}]: ${this.message}`;
        },
        enumerable: false,
        writable: true,
        configurable: true,
      },
    });
  }
}

class ERR_INVALID_ARG_TYPE extends TypeError {
  constructor(name, expected, actual) {
    super();

    assert(typeof name === 'string', "'name' must be a string");
    if (!Array.isArray(expected)) {
      expected = [expected];
    }

    let msg = 'The ';
    if (name.endsWith(' argument')) {
      // For cases like 'first argument'
      msg += `${name} `;
    } else {
      const type = name.includes('.') ? 'property' : 'argument';
      msg += `"${name}" ${type} `;
    }
    msg += 'must be ';

    const types = [];
    const instances = [];
    const other = [];

    for (const value of expected) {
      assert(typeof value === 'string',
        'All expected entries have to be of type string');
      if (kTypes.includes(value)) {
        types.push(value.toLowerCase());
      } else if (classRegExp.test(value)) {
        instances.push(value);
      } else {
        assert(value !== 'object',
          'The value "object" should be written as "Object"');
        other.push(value);
      }
    }

    // Special handle `object` in case other instances are allowed to outline
    // the differences between each other.
    if (instances.length > 0) {
      const pos = types.indexOf('object');
      if (pos !== -1) {
        types.splice(pos, 1);
        instances.push('Object');
      }
    }

    if (types.length > 0) {
      if (types.length > 2) {
        const last = types.pop();
        msg += `one of type ${types.join(', ')}, or ${last}`;
      } else if (types.length === 2) {
        msg += `one of type ${types[0]} or ${types[1]}`;
      } else {
        msg += `of type ${types[0]}`;
      }
      if (instances.length > 0 || other.length > 0)
        msg += ' or ';
    }

    if (instances.length > 0) {
      if (instances.length > 2) {
        const last = instances.pop();
        msg +=
          `an instance of ${instances.join(', ')}, or ${last}`;
      } else {
        msg += `an instance of ${instances[0]}`;
        if (instances.length === 2) {
          msg += ` or ${instances[1]}`;
        }
      }
      if (other.length > 0)
        msg += ' or ';
    }

    if (other.length > 0) {
      if (other.length > 2) {
        const last = other.pop();
        msg += `one of ${other.join(', ')}, or ${last}`;
      } else if (other.length === 2) {
        msg += `one of ${other[0]} or ${other[1]}`;
      } else {
        if (other[0].toLowerCase() !== other[0])
          msg += 'an ';
        msg += `${other[0]}`;
      }
    }

    if (actual == null) {
      msg += `. Received ${actual}`;
    } else if (typeof actual === 'function' && actual.name) {
      msg += `. Received function ${actual.name}`;
    } else if (typeof actual === 'object') {
      if (actual.constructor && actual.constructor.name) {
        msg += `. Received an instance of ${actual.constructor.name}`;
      } else {
        const inspected = util.inspect(actual, { depth: -1 });
        msg += `. Received ${inspected}`;
      }
    } else {
      let inspected = util.inspect(actual, { colors: false });
      if (inspected.length > 25) {
        inspected = `${inspected.slice(0, 25)}...`;
      }
      msg += `. Received type ${typeof actual} (${inspected})`;
    }

    this.code = this.constructor.name;

    Object.defineProperties(this, {
      message: {
        value: msg,
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

function validateBoolean(value, name) {
  if (typeof value !== 'boolean')
    throw new ERR_INVALID_ARG_TYPE(name, 'boolean', value);
}

function validateObject(value, name, {
  nullable = false,
  allowArray = false,
  allowFunction = false,
} = {}) {
  if ((!nullable && value === null) ||
    (!allowArray && Array.isArray(value)) ||
    (typeof value !== 'object' && (
      !allowFunction || typeof value !== 'function'
    ))) {
    throw new ERR_INVALID_ARG_TYPE(name, 'Object', value);
  }
};

module.exports = {
  validateBoolean,
  validateObject,
  codes: {
    ERR_OPERATION_FAILED
  }
};
