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
  globalSetup: './jest.globalSetup.js',
  resetMocks: true,
  setupFilesAfterEnv: ['./jest.setup.js'],
  testEnvironment: 'node',
  testTimeout: 30000,
};
