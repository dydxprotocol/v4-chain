const baseConfig = require('./node_modules/@dydxprotocol-indexer/dev/.eslintrc');

module.exports = {
  ...baseConfig,

  // Override the base configuration to set the correct tsconfigRootDir.
  parserOptions: {
    ...baseConfig.parserOptions,
    tsconfigRootDir: __dirname,
  },
  // ignore generated ts files
  ignorePatterns: baseConfig.ignorePatterns.concat([
    '**/generated/**/*.ts',
  ]),
};
