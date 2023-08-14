import {
  base64ToBytes,
  bigIntToBytes,
  bytesToBase64,
  bytesToBigInt,
} from '../src/bytes-helpers';

const TEST_CASES: [bigint, string, Uint8Array][] = [
  [BigInt('0'), 'Ag==', Uint8Array.from([0x02])],
  [BigInt('-0'), 'Ag==', Uint8Array.from([0x02])],
  [BigInt('1'), 'AgE=', Uint8Array.from([0x02, 0x01])],
  [BigInt('-1'), 'AwE=', Uint8Array.from([0x03, 0x01])],
  [BigInt('255'), 'Av8=', Uint8Array.from([0x02, 0xFF])],
  [BigInt('-255'), 'A/8=', Uint8Array.from([0x03, 0xFF])],
  [BigInt('256'), 'AgEA', Uint8Array.from([0x02, 0x01, 0x00])],
  [BigInt('-256'), 'AwEA', Uint8Array.from([0x03, 0x01, 0x00])],
  [BigInt('123456789'), 'AgdbzRU=', Uint8Array.from([0x02, 0x07, 0x5b, 0xcd, 0x15])],
  [BigInt('-123456789'), 'AwdbzRU=', Uint8Array.from([0x03, 0x07, 0x5b, 0xcd, 0x15])],
  [BigInt('123456789123456789'), 'AgG2m0us0F8V', Uint8Array.from([0x02, 0x01, 0xb6, 0x9b, 0x4b, 0xac, 0xd0, 0x5f, 0x15])],
  [BigInt('-123456789123456789'), 'AwG2m0us0F8V', Uint8Array.from([0x03, 0x01, 0xb6, 0x9b, 0x4b, 0xac, 0xd0, 0x5f, 0x15])],
  [BigInt('123456789123456789123456789'), 'AmYe/fLjsZ98BF8V', Uint8Array.from([0x02, 0x66, 0x1e, 0xfd, 0xf2, 0xe3, 0xb1, 0x9f, 0x7c, 0x04, 0x5f, 0x15])],
  [BigInt('-123456789123456789123456789'), 'A2Ye/fLjsZ98BF8V', Uint8Array.from([0x03, 0x66, 0x1e, 0xfd, 0xf2, 0xe3, 0xb1, 0x9f, 0x7c, 0x04, 0x5f, 0x15])],
];

describe('bigIntToBytes', () => {
  it('Success', () => {
    TEST_CASES.forEach(([i, _, b]: [bigint, string, Uint8Array]) => {
      expect(bigIntToBytes(i)).toEqual(b);
    });
  });
});

describe('bytesToBigInt', () => {
  it('Success', () => {
    TEST_CASES.forEach(([i, _, b]: [bigint, string, Uint8Array]) => {
      expect(bytesToBigInt(b)).toEqual(i);
    });
  });
});

describe('base64ToBytes', () => {
  it('Success', () => {
    TEST_CASES.forEach(([_, i, b]: [bigint, string, Uint8Array]) => {
      expect(base64ToBytes(i)).toEqual(b);
    });
  });
});

describe('bytesToBase64', () => {
  it('Success', () => {
    TEST_CASES.forEach(([_, i, b]: [bigint, string, Uint8Array]) => {
      expect(bytesToBase64(b)).toEqual(i);
    });
  });
});
