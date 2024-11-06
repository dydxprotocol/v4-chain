require('dotenv-flow').config({
    node_env: 'test',
    silent: true,
  });
  
  module.exports = {
    roots: [
      '<rootDir>/build/__tests__',
    ],
    testRegex: 'build/__tests__\\/.*\\.test\\.js$',
    moduleFileExtensions: [
      'js',
      'json',
      'node',
    ],
    moduleNameMapper: {
      '@klyraprotocol/v4-client-js': '<rootDir>/../v4-client-js/build/v4-client-js/src',
      '@klyraprotocol/v4-proto/(.*)': '<rootDir>/../v4-client-js/build/v4-proto-js/$1'
    },
    moduleDirectories: [
      'node_modules',
      '../v4-client-js/build/v4-client-js/src',
      '../v4-client-js/build/v4-proto-js'
    ],
    resetMocks: true,
    testEnvironment: 'node',
    testTimeout: 30000,
    transform: {
      '^.+\\.js$': 'babel-jest'
    },
    coveragePathIgnorePatterns: ['src/codegen/'],
    setupFiles: ['<rootDir>/jest.setup.js']
  };