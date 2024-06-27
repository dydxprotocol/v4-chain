'use strict';

const crypto = require('crypto');

module.exports = typeof crypto.randomUUID === 'function'
  ? crypto.randomUUID
  : require('./polyfill');
