// Use the base configuration as-is.
exports = require('./node_modules/@dydxprotocol/node-service-base-dev/jest.config');

module.exports = {
  ...exports,
  moduleNameMapper: {
    '^axios$': require.resolve('axios'),
  },
  coveragePathIgnorePatterns: ['src/codegen/'],
};
