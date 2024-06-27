# crypto-randomuuid

This is a polyfill for the `crypto.randomUUID` method in Node.js. It will use
the built-in version, if present. There are plenty of other uuid modules, but
this one aims to be as functionally identical as possible to the Node.js core
function.

This uses a pure JavaScript replacement of the `secureBuffer` function using
`randomFillSync` rather than the native version using `OPENSSL_secure_malloc`
in Node.js core. This may have security implications, so I'd recommend against
using this anywhere that cryptographically secure uuids are important.

## Install

```sh
npm install crypto-randomuuid
```

## Usage

https://nodejs.org/api/crypto.html#crypto_crypto_randomuuid_options

## License

This is all copy/pasted from Node.js core, so see the license there: https://github.com/nodejs/node/blob/master/LICENSE
