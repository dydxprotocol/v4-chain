# @bugsnag/safe-json-stringify
[![Build Status](https://travis-ci.org/bugsnag/safe-json-stringify.svg?branch=master)](https://travis-ci.org/bugsnag/safe-json-stringify)
[![NPM Version](https://img.shields.io/npm/v/@bugsnag/safe-json-stringify.svg)](https://www.npmjs.com/package/@bugsnag/safe-json-stringify)

This is a fork of [`safe-json-stringify`](https://github.com/debitoor/safe-json-stringify) with some opinionated changes, specifically for [Bugsnag's](https://bugsnag.com) use-case:

- IE8 support. Our [JS notifier](https://github.com/bugsnag/bugsnag-js) supports IE8, so all of its dependencies must use IE8 compatible APIs.
- A hard limit on object depth/breadth. Given a deep/wide enough data structure, the original version would run out of memory or stack allocations. Values are replaced with the string `...` in such circumstances.
- Doesn't replace repeated, but non-circular references (implemented by [@MikeRalphson](https://github.com/MikeRalphson) in [their fork](https://github.com/MikeRalphson/safe-json-stringify/tree/circular))
- Adds the ability to apply key-based filters to specific subtrees

## Installation

```sh
npm install @bugsnag/safe-json-stringify
```

## Usage

The API is the same as `JSON.stringify`, with an additional `options` object:

```js
stringify(obj, [optional replacer], [optional spaces], [options])
```

- `options.redactedKeys`: a list of keys whose value should be replaced with the string `[REDACTED]`. Keys can be strings for case insensitive matches, or regexes for partial or pattern matches. The array can contain a mixture of both.
- `options.redactedPaths`: a list of paths where the `redactedKeys` option will be applied. The format of these strings are key names separated by `.` and if the property is an array, it is represented with `[]`. For example: `events.[].metaData`.

### Example

```js
var stringify = require('@bugsnag/safe-json-stringify')
stringify({
  api_key: 'd145b8e5afb56516423bc4d605e45442',
  events: [
    {
      errorMessage: 'Failed load tickets',
      errorClass: 'CheckoutError',
      user: {
        name: 'Jim Bug',
        email: 'jim@bugsnag.com',
        api_key: '245b39ebd3cd3992e85bffc81c045924'
      }
    }
  ]
}, null, 2, {
  redactedKeys: [ 'api_key' ],
  redactedPaths: [ 'events.[].user' ]
})

// yields the following json:
// {
//   "api_key": "d145b8e5afb56516423bc4d605e45442",
//   "events": [
//     {
//       "errorMessage": "Failed load tickets",
//       "errorClass": "CheckoutError",
//       "user": {
//         "name": "Jim Bug",
//         "email": "jim@bugsnag.com",
//         "api_key": "[REDACTED]"
//       }
//     }
//   ]
// }
```


## License

The original code is licensed with MIT, and so are the modifications.

The MIT License (MIT)

Copyright (c) 2014-2018 [Debitoor](https://debitoor.com/), [Bugsnag](https://bugsnag.com/)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
