# iserror [![Build Status](https://travis-ci.org/yefremov/iserror.svg?branch=master)](https://travis-ci.org/yefremov/iserror) [![Coverage Status](https://coveralls.io/repos/github/yefremov/iserror/badge.svg?branch=master)](https://coveralls.io/github/yefremov/iserror?branch=master) [![npm version](https://badge.fury.io/js/iserror.svg)](https://badge.fury.io/js/iserror)

Test whether value is error object.

## Installation

```bash
$ npm install iserror
```

## API

```js
const isError = require('iserror');

isError(new Error('Some error'));
// => true

isError({});
// => false
```

## Running tests

```bash
$ npm test
```

## License

[MIT](LICENSE)
