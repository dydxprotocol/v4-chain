import BigNumber from 'bignumber.js';
import Long from 'long';

import { encodeJson } from '../src/lib/helpers';

const longValue = Long.fromInt(1, true);
const objWithLong = {
  longValue,
};

const jsonedWithLong = encodeJson(objWithLong);
console.log(jsonedWithLong);

const bigNumberValue = BigNumber(1);
const objWithBigNumber = {
  bigNumberValue,
};
const jsonedWithBigNumber = encodeJson(objWithBigNumber);
console.log(jsonedWithBigNumber);

const buffer = Buffer.from('this is a tést');
const objWithBuffer = {
  buffer,
};

// Object.keys(objWithBuffer).forEach((k) => {
//   console.log(k, objWithBuffer[k]); // ERROR!!
// });
const isBuffer2 = objWithBuffer.buffer instanceof Uint8Array;
console.log(isBuffer2);

const jsonedWithBuffer = encodeJson(objWithBuffer);
console.log(jsonedWithBuffer);
