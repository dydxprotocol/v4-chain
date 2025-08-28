// This file runs before each test file.
const { performance } = require('perf_hooks');

// eslint-disable-next-line
const originalDescriptor = Object.getOwnPropertyDescriptor(global, 'performance');

// Remove the existing property
delete global.performance;

// Re-define it as writable
Object.defineProperty(global, 'performance', {
  value: performance,
  writable: true,
  enumerable: true,
  configurable: true,
});
