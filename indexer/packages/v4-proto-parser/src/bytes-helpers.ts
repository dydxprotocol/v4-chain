/**
 * Converts a byte array (representing an arbitrary-size signed integer) into a bigint.
 * @param u Array of bytes represented as a Uint8Array.
 */
export function bytesToBigInt(
  u: Uint8Array,
): bigint {
  if (u.length <= 1) {
    return BigInt(0);
  }
  // eslint-disable-next-line no-bitwise
  const negated: boolean = (u[0] & 1) === 1;
  const hex: string = Buffer.from(u.slice(1)).toString('hex');
  const abs: bigint = BigInt(`0x${hex}`);
  return negated
    ? -abs
    : abs;
}

/**
 * Converts a bigint to a byte array.
 * @param b bigint value that must be translated.
 */
export function bigIntToBytes(
  b: bigint,
): Uint8Array {
  // Special-case zero.
  if (b === BigInt(0)) {
    return Uint8Array.from([0x02]);
  }

  const negated: boolean = b < 0;
  const abs: bigint = negated ? -b : b;

  // Generate the hex string and have it be even-length (prepended with a zero if needed).
  const hex: string = abs.toString(16);
  const hexPadded: string = hex.length % 2 === 0
    ? hex
    : `0${hex}`;

  // Add the sign+version byte.
  const hexWithSign: string = `${negated ? '03' : '02'}${hexPadded}`;

  return Uint8Array.from(Buffer.from(hexWithSign, 'hex'));
}

/**
 * Converts a byte array (representing an arbitrary-size signed integer) into a base64 string.
 * This is useful for JSON encoding.
 * @param u Array of bytes represented as a Uint8Array.
 */
export function bytesToBase64(
  u: Uint8Array,
): string {
  return Buffer.from(u).toString('base64');
}

/**
 * Converts a base64 string into a byte array.
 * This is useful for JSON decoding.
 * @param u string value in base64.
 */
export function base64ToBytes(
  s: string,
): Uint8Array {
  return Uint8Array.from(
    Buffer.from(s, 'base64'),
  );
}
