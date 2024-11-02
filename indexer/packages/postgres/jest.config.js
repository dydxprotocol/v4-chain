// Use the base configuration as-is.
const baseConfig = require('./node_modules/@klyraprotocol-indexer/dev/jest.config');

module.exports = {
  ...baseConfig,
  testSequencer: './customSequencer.js',
};
