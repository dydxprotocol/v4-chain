# eslint-plugin-no-only-tests

[![Version](https://img.shields.io/npm/v/eslint-plugin-no-only-tests.svg)](https://www.npmjs.com/package/eslint-plugin-no-only-tests) [![Downloads](https://img.shields.io/npm/dm/eslint-plugin-no-only-tests.svg)](https://npmcharts.com/compare/eslint-plugin-no-only-tests?minimal=true) [![CircleCI](https://circleci.com/gh/levibuzolic/eslint-plugin-no-only-tests.svg?style=shield)](https://circleci.com/gh/levibuzolic/eslint-plugin-no-only-tests)

ESLint rule for `.only` tests in [mocha](https://mochajs.org/) and other JS testing libraries.

Currently matches the following test blocks by default: `describe`, `it`, `context`, `tape`, `test`, `fixture`, `serial`.

Designed to prevent you from committing `.only` tests to CI, disabling tests for your whole team.

As of v2.3 you can now override the test blocks and focus functions.

## Installation

First you'll need to install [ESLint](http://eslint.org) and the plugin:

```bash
npm install --save-dev eslint eslint-plugin-no-only-tests
# OR
yarn add --dev eslint eslint-plugin-no-only-tests
```

**Note:** If you installed ESLint globally (using the `-g` flag) then you must also install `eslint-plugin-no-only-tests` globally.

## Usage

Add `no-only-tests` to the plugins section of your `.eslintrc` configuration file. You can omit the `eslint-plugin-` prefix:

```json
{
  "plugins": [
    "no-only-tests"
  ]
}
```

Then configure the rules you want to use under the rules section.

```json
{
  "rules": {
    "no-only-tests/no-only-tests": "error"
  }
}
```

If you use a testing framework that uses an unsupported block name, or a different way of focusing test (something other than `.only`) you can specify an array of blocks and focus methods to match in the options.

```json
{
  "rules": {
    "no-only-tests/no-only-tests": ["error", {"block": ["test", "it", "assert"], "focus": ["only", "focus"]}]
  }
}
```

The above example will catch any uses of `test.only`, `test.focus`, `it.only`, `it.focus`, `assert.only` and `assert.focus`.
