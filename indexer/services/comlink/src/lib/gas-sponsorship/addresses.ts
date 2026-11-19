import { fromBech32 } from '@cosmjs/encoding';

export const INJECTIVE_ADDRESS_PREFIX: string = 'inj';
export const DYDX_ADDRESS_PREFIX: string = 'dydx';
// Length in bytes of the address of an account controlled by a key.
export const ACCOUNT_ADDRESS_LENGTH: number = 20;

/**
 * Whether `address` is a well-formed address, with a valid checksum, of an account controlled
 * by a key on the chain with the given bech32 prefix.
 */
export function isValidAccountAddress(address: unknown, prefix: string): boolean {
  if (typeof address !== 'string') {
    return false;
  }
  try {
    const decoded: { prefix: string, data: Uint8Array } = fromBech32(address);
    return decoded.prefix === prefix && decoded.data.length === ACCOUNT_ADDRESS_LENGTH;
  } catch (error) {
    return false;
  }
}
