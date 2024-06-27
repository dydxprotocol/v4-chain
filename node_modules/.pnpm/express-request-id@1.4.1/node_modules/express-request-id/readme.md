# express-request-id [![NPM version][npm-image]][npm-url] [![Build Status][travis-image]][travis-url]

Generate UUID for request and add it to `X-Request-Id` header. In case request contains `X-Request-Id` header, uses its value instead.

```js

var app = require('express')();
var addRequestId = require('express-request-id')();

app.use(addRequestId);

app.get('/', function (req, res, next) {
    res.send(req.id);
    next();
});

app.listen(3000, function() {
    console.log('Listening on port %d', server.address().port);
});

// curl localhost:3000
// d7c32387-3feb-452b-8df1-2d8338b3ea22
```

# API

### express-request-id([options])

Returns middleware function, that appends request id to req object.

#### options

 * `uuidVersion` - version of uuid to use (defaults to `v4`). Can be one of methods from [node-uuid](https://github.com/broofa/node-uuid).
 * `setHeader` - boolean, indicates that header should be added to response (defaults to `true`).
 * `headerName` - string, indicates the header name to use (defaults to `X-Request-Id`).
 * `attributeName` - string, indicates the attribute name used for the identifier on the request object (defaults to `id`)

This options fields are passed to node-uuid functions directly:

 * Whole `options` object, that can contain fields like: `node`, `clockseq`, `msecs`, `nsecs`.
 * `options.buffer` and `options.offset` to uuid function as second and third parameters.

# License

MIT (c) 2014 Vsevolod Strukchinsky (floatdrop@gmail.com)

[npm-url]: https://npmjs.org/package/express-request-id
[npm-image]: http://img.shields.io/npm/v/express-request-id.svg

[travis-url]: https://travis-ci.org/floatdrop/express-request-id
[travis-image]: http://img.shields.io/travis/floatdrop/express-request-id.svg
