/**
 * A minimal protobuf encoder covering the wire types the grant transaction needs: varints and
 * length-delimited fields. Fields are written in the order the methods are called, so callers
 * write them in ascending field-number order to match what the Cosmos SDK produces.
 *
 * comlink has no protobuf runtime or Cosmos transaction types among its dependencies, and the
 * handful of messages needed here are small and stable, so they are encoded directly. The
 * encoding is checked byte-for-byte against a transaction from Injective mainnet in the tests.
 */

const WIRE_TYPE_VARINT: number = 0;
const WIRE_TYPE_LENGTH_DELIMITED: number = 2;
const MAX_UINT64: bigint = BigInt('18446744073709551615');

function encodeVarint(value: bigint): number[] {
  if (value < BigInt(0) || value > MAX_UINT64) {
    throw new Error(`varint out of range: ${value}`);
  }
  const bytes: number[] = [];
  let remaining: bigint = value;
  while (remaining >= BigInt(0x80)) {
    // eslint-disable-next-line no-bitwise
    bytes.push(Number(remaining & BigInt(0x7f)) | 0x80);
    // eslint-disable-next-line no-bitwise
    remaining >>= BigInt(7);
  }
  bytes.push(Number(remaining));
  return bytes;
}

export class ProtoWriter {
  private readonly chunks: Uint8Array[] = [];

  /**
   * Writes an unsigned integer field. Zero is the default value and is left out, as proto3
   * requires.
   */
  public uint(fieldNumber: number, value: bigint | number): ProtoWriter {
    const asBigInt: bigint = BigInt(value);
    if (asBigInt === BigInt(0)) {
      return this;
    }
    this.tag(fieldNumber, WIRE_TYPE_VARINT);
    this.chunks.push(Uint8Array.from(encodeVarint(asBigInt)));
    return this;
  }

  /**
   * Writes a string field. The empty string is the default value and is left out.
   */
  public string(fieldNumber: number, value: string): ProtoWriter {
    if (value.length === 0) {
      return this;
    }
    return this.lengthDelimited(fieldNumber, Buffer.from(value, 'utf8'));
  }

  /**
   * Writes a bytes field. Empty bytes are the default value and are left out.
   */
  public bytes(fieldNumber: number, value: Uint8Array): ProtoWriter {
    if (value.length === 0) {
      return this;
    }
    return this.lengthDelimited(fieldNumber, value);
  }

  /**
   * Writes an embedded message, or one element of a repeated field. Always written, even when
   * the encoded message is empty, because presence is meaningful for these.
   */
  public message(fieldNumber: number, encoded: Uint8Array): ProtoWriter {
    return this.lengthDelimited(fieldNumber, encoded);
  }

  /**
   * Writes a repeated string field. Every element is written, including empty strings.
   */
  public repeatedString(fieldNumber: number, values: string[]): ProtoWriter {
    values.forEach((value: string) => {
      this.lengthDelimited(fieldNumber, Buffer.from(value, 'utf8'));
    });
    return this;
  }

  public finish(): Uint8Array {
    return Uint8Array.from(Buffer.concat(this.chunks));
  }

  private lengthDelimited(fieldNumber: number, value: Uint8Array): ProtoWriter {
    this.tag(fieldNumber, WIRE_TYPE_LENGTH_DELIMITED);
    this.chunks.push(Uint8Array.from(encodeVarint(BigInt(value.length))));
    this.chunks.push(value);
    return this;
  }

  private tag(fieldNumber: number, wireType: number): void {
    if (!Number.isInteger(fieldNumber) || fieldNumber < 1) {
      throw new Error(`invalid field number: ${fieldNumber}`);
    }
    // eslint-disable-next-line no-bitwise
    const tag: bigint = (BigInt(fieldNumber) << BigInt(3)) | BigInt(wireType);
    this.chunks.push(Uint8Array.from(encodeVarint(tag)));
  }
}
