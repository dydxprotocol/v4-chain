const Sequencer = require('@jest/test-sequencer').default;

class CustomSequencer extends Sequencer {
  sort(tests) {
    // Sort tests in reverse alphabetical order by file path
    return tests.sort((a, b) => (a.path < b.path ? 1 : -1));
  }
}

module.exports = CustomSequencer;