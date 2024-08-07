// // Use the base configuration as-is.
// module.exports = require('./node_modules/@dydxprotocol-indexer/dev/jest.config');

// Use the base configuration as-is.
const baseConfig = require('./node_modules/@dydxprotocol-indexer/dev/jest.config');

module.exports = {
  ...baseConfig,
  testSequencer: './customSequencer.js',
};