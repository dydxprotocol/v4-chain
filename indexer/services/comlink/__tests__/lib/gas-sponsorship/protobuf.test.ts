import { toHex } from '@cosmjs/encoding';

import { ProtoWriter } from '../../../src/lib/gas-sponsorship/protobuf';

function hex(writer: ProtoWriter): string {
  return toHex(writer.finish());
}

describe('ProtoWriter', () => {
  describe('uint', () => {
    // The encodings of 1, 150 and 300 are the worked examples in the protobuf encoding guide.
    it.each([
      [1, '0801'],
      [127, '087f'],
      [128, '088001'],
      [150, '089601'],
      [300, '08ac02'],
      [BigInt('18446744073709551615'), '08ffffffffffffffffff01'],
    ])('encodes %s as a varint', (value: number | bigint, expected: string) => {
      expect(hex(new ProtoWriter().uint(1, value))).toEqual(expected);
    });

    it('leaves out zero', () => {
      expect(hex(new ProtoWriter().uint(1, 0))).toEqual('');
      expect(hex(new ProtoWriter().uint(1, BigInt(0)))).toEqual('');
    });

    it('rejects values outside the unsigned 64-bit range', () => {
      expect(() => new ProtoWriter().uint(1, -1)).toThrow('varint out of range');
      expect(() => new ProtoWriter().uint(1, BigInt('18446744073709551616'))).toThrow(
        'varint out of range',
      );
    });
  });

  describe('string', () => {
    it('encodes the worked example from the protobuf encoding guide', () => {
      expect(hex(new ProtoWriter().string(2, 'testing'))).toEqual('120774657374696e67');
    });

    it('encodes the length in bytes, not characters', () => {
      expect(hex(new ProtoWriter().string(1, 'é'))).toEqual('0a02c3a9');
    });

    it('leaves out the empty string', () => {
      expect(hex(new ProtoWriter().string(1, ''))).toEqual('');
    });
  });

  describe('bytes', () => {
    it('encodes a length of more than one varint byte', () => {
      const encoded: string = hex(new ProtoWriter().bytes(1, new Uint8Array(200).fill(7)));

      expect(encoded.slice(0, 6)).toEqual('0ac801');
      expect(encoded.slice(6)).toEqual('07'.repeat(200));
    });

    it('leaves out empty bytes', () => {
      expect(hex(new ProtoWriter().bytes(1, new Uint8Array(0)))).toEqual('');
    });
  });

  describe('message', () => {
    it('writes an empty message, unlike empty bytes', () => {
      expect(hex(new ProtoWriter().message(3, new Uint8Array(0)))).toEqual('1a00');
    });

    it('nests an encoded message', () => {
      const inner: Uint8Array = new ProtoWriter().uint(1, 150).finish();

      expect(hex(new ProtoWriter().message(3, inner))).toEqual('1a03089601');
    });
  });

  describe('repeatedString', () => {
    it('writes every element, including an empty one', () => {
      expect(hex(new ProtoWriter().repeatedString(2, ['a', '', 'b']))).toEqual(
        ['120161', '1200', '120162'].join(''),
      );
    });

    it('writes nothing for an empty list', () => {
      expect(hex(new ProtoWriter().repeatedString(2, []))).toEqual('');
    });
  });

  it('writes fields in the order they are added', () => {
    expect(hex(new ProtoWriter().uint(1, 1).string(2, 'a').uint(3, 2))).toEqual(
      ['0801', '120161', '1802'].join(''),
    );
  });

  it('encodes field numbers above 15 with a multi-byte tag', () => {
    expect(hex(new ProtoWriter().uint(16, 1))).toEqual('800101');
  });

  it('rejects an invalid field number', () => {
    expect(() => new ProtoWriter().uint(0, 1)).toThrow('invalid field number');
  });
});
