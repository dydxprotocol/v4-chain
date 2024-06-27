# Unique validation for Objection.js

[![npm](https://img.shields.io/npm/v/objection-unique.svg?style=flat-square)](https://npmjs.org/package/objection-unique)
![node](https://img.shields.io/node/v/objection-unique.svg?style=flat-square)
[![Build Status](https://img.shields.io/travis/seegno/objection-unique/master.svg?style=flat-square)](https://travis-ci.org/seegno/objection-unique)
[![Coverage Status](https://img.shields.io/coveralls/seegno/objection-unique/master.svg?style=flat-square)](https://coveralls.io/github/seegno/objection-unique?branch=master)

This plugin adds a unique validation for [Objection.js](https://github.com/Vincit/objection.js/) models.

**NOTE:** Unique validation at update only works with `$query` methods.

## Installation

### NPM

```sh
npm i objection-unique --save
```

### Yarn

```sh
yarn add objection-unique
```

## Usage

### Mixin the plugin

```js
// Import objection model.
const Model = require('objection').Model;

// Import the plugin.
const unique = require('objection-unique')({
  fields: ['email', 'username', ['phone_prefix','phone_number']],
  identifiers: ['id']
});

// Mixin the plugin.
class User extends unique(Model) {
  static get tableName() {
    return 'User';
  }
}
```

### Validate insert

```js
/**
 * Insert.
 */

// Insert one user.
await User.query().insert({ email: 'foo', username: 'bar' });

try {
  // Try to insert another user with the same data.
  await User.query().insert({ email: 'foo', username: 'bar' });
} catch (e) {
    // Exception with the invalid unique fields
    //
    // {
    //   email: [{
    //     keyword: 'unique',
    //     message: 'email already in use.'
    //   }],
    //   username: [{
    //     keyword: 'unique',
    //     message: 'username already in use.'
    //   }
    // }
}
```

### Validate update/patch

```js
/**
 * Update/Patch.
 */

// Insert one user.
await User.query().insert({ email: 'foo', username: 'bar' });

// Insert the user that we want to update.
const user = await User.query().insertAndFetch({ email: 'biz', username: 'buz' });

try {
  user.$query().update({ email: 'foo', username: 'buz' });
  // user.$query().patch({ email: 'foo' });
} catch (e) {
  // Exception with the invalid unique fields
  //
  // {
  //   email: [{
  //     keyword: 'unique',
  //     message: 'email already in use.'
  //   }]
  // }
}
```

## Options

**fields:** The unique fields. Compound fields can be specified as an array

**identifiers:** The fields that identifies the model. (Default: ['id'])

These options can be provided when instantiating the plugin:

```js
const unique = require('objection-unique')({
  fields: ['email', 'username', ['phone_prefix', 'phone_number']],
  identifiers: ['id']
});
```

## Tests

Run the tests from the root directory:

```sh
npm test
```

## Contributing & Development

### Contributing

Found a bug or want to suggest something? Take a look first on the current and closed [issues](https://github.com/seegno/objection-unique/issues). If it is something new, please [submit an issue](https://github.com/seegno/objection-unique/issues/new).

### Develop

It will be awesome if you can help us evolve `objection-unique`. Want to help?

1. [Fork it](https://github.com/seegno/objection-unique).
2. `npm install`.
3. Hack away.
4. Run the tests: `npm test`.
5. Create a [Pull Request](https://github.com/seegno/objection-unique/compare).
