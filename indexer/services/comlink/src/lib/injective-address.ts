import { decode, encode, fromWords } from 'bech32';
import { getAddress, pad } from 'viem';

const DYDX_PREFIX: string = 'dydx';
const INJECTIVE_PREFIX: string = 'inj';
const ACCOUNT_ADDRESS_LENGTH: number = 20;

/**
 * Deposits are bridged to the Injective account that shares its 20 address bytes with the user's
 * dYdX account. Injective accepts the dYdX key type and derives that key's address the same way
 * dYdX does, so the key that controls the dYdX account also controls this one. The user then
 * sends the funds on to dYdX from it.
 *
 * The three forms below are the same account: as an Injective address, as the EVM-style address
 * Skip and CCTP use for Injective, and as the 32-byte CCTP mint recipient.
 */
function getAddressBytes(dydxAddress: string): Buffer {
  const { prefix, words } = decode(dydxAddress);
  if (prefix !== DYDX_PREFIX) {
    throw new Error(`Expected a dydx address, got prefix "${prefix}"`);
  }
  const bytes: Buffer = Buffer.from(fromWords(words));
  if (bytes.length !== ACCOUNT_ADDRESS_LENGTH) {
    throw new Error(`Expected a ${ACCOUNT_ADDRESS_LENGTH}-byte address, got ${bytes.length} bytes`);
  }
  return bytes;
}

export function getInjectiveAddress(dydxAddress: string): string {
  getAddressBytes(dydxAddress);
  return encode(INJECTIVE_PREFIX, decode(dydxAddress).words);
}

export function getInjectiveEvmAddress(dydxAddress: string): `0x${string}` {
  return getAddress(`0x${getAddressBytes(dydxAddress).toString('hex')}`);
}

/**
 * The CCTP mint recipient: the address left-padded with zeros to 32 bytes, in lower-case hex.
 */
export function getInjectiveMintRecipient(dydxAddress: string): `0x${string}` {
  return pad(`0x${getAddressBytes(dydxAddress).toString('hex')}`, { size: 32 });
}
