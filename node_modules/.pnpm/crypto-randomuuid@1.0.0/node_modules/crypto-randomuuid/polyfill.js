'use strict';

const crypto = require('crypto');

const {
  validateBoolean,
  validateObject,
  codes: {
    ERR_OPERATION_FAILED
  }
} = require('./validators');

const { randomFillSync } = crypto;

// This is a non-cryptographically secure replacement for the native version
// of the `secureBuffer` function used in Node.js core. This means `randomUUID`
// should not be used where cryptographically secure uuids are important.
//
// Node.js core uses a native version which uses `OPENSSL_secure_malloc`
// rather than `randomFillSync`.
function secureBuffer (size) {
  const buf = Buffer.alloc(size);
  return randomFillSync(buf);
}

// Implements an RFC 4122 version 4 random UUID.
// To improve performance, random data is generated in batches
// large enough to cover kBatchSize UUID's at a time. The uuidData
// buffer is reused. Each call to randomUUID() consumes 16 bytes
// from the buffer.

const kBatchSize = 128;
let uuidData;
let uuidNotBuffered;
let uuidBatch = 0;

let hexBytesCache;
function getHexBytes () {
  if (hexBytesCache === undefined) {
    hexBytesCache = new Array(256);
    for (let i = 0; i < hexBytesCache.length; i++) {
      const hex = i.toString(16);
      hexBytesCache[i] = hex.padStart(2, '0');
    }
  }
  return hexBytesCache;
}

function serializeUUID (buf, offset = 0) {
  const kHexBytes = getHexBytes();
  // xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  return kHexBytes[buf[offset]] +
    kHexBytes[buf[offset + 1]] +
    kHexBytes[buf[offset + 2]] +
    kHexBytes[buf[offset + 3]] +
    '-' +
    kHexBytes[buf[offset + 4]] +
    kHexBytes[buf[offset + 5]] +
    '-' +
    kHexBytes[(buf[offset + 6] & 0x0f) | 0x40] +
    kHexBytes[buf[offset + 7]] +
    '-' +
    kHexBytes[(buf[offset + 8] & 0x3f) | 0x80] +
    kHexBytes[buf[offset + 9]] +
    '-' +
    kHexBytes[buf[offset + 10]] +
    kHexBytes[buf[offset + 11]] +
    kHexBytes[buf[offset + 12]] +
    kHexBytes[buf[offset + 13]] +
    kHexBytes[buf[offset + 14]] +
    kHexBytes[buf[offset + 15]];
}

function getBufferedUUID () {
  if (!uuidData) uuidData = secureBuffer(16 * kBatchSize);
  if (uuidData === undefined)
    throw new ERR_OPERATION_FAILED('Out of memory');

  if (uuidBatch === 0) randomFillSync(uuidData);
  uuidBatch = (uuidBatch + 1) % kBatchSize;
  return serializeUUID(uuidData, uuidBatch * 16);
}

function getUnbufferedUUID () {
  if (!uuidNotBuffered) uuidNotBuffered = secureBuffer(16 * kBatchSize);
  if (uuidNotBuffered === undefined)
    throw new ERR_OPERATION_FAILED('Out of memory');
  randomFillSync(uuidNotBuffered);
  return serializeUUID(uuidNotBuffered);
}

function randomUUID (options) {
  if (options !== undefined)
    validateObject(options, 'options');
  const {
    disableEntropyCache = false,
  } = options || {};

  validateBoolean(disableEntropyCache, 'options.disableEntropyCache');

  return disableEntropyCache ? getUnbufferedUUID() : getBufferedUUID();
}

module.exports = randomUUID;
