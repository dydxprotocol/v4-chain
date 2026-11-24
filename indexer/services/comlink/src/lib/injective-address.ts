import { decode, encode, fromWords } from 'bech32';
import { getAddress, pad } from 'viem';

const INJECTIVE_PREFIX: string = 'inj';
const ACCOUNT_ADDRESS_LENGTH: number = 20;

/**
 * Deposits are bridged to an Injective account that the user names, and the user sends the funds
 * on to dYdX from it. It is the account their app holds the key for.
 *
 * This service cannot work that account out. The app derives its Injective key from the same
 * secret as the dYdX key, but as a different key, and nothing known here leads from one to the
 * other. So the user gives the address when they approve the deposit policy, where it becomes
 * part of what they sign (see `policy-engine.ts`), and it is stored with the user from then on.
 *
 * The forms below are the same account: as an Injective address, as the EVM-style address Skip
 * and CCTP use for Injective, and as the 32-byte CCTP mint recipient.
 */
function getAddressBytes(injectiveAddress: string): Buffer {
  const { prefix, words } = decode(injectiveAddress);
  if (prefix !== INJECTIVE_PREFIX) {
    throw new Error(`Expected an Injective address, got prefix "${prefix}"`);
  }
  const bytes: Buffer = Buffer.from(fromWords(words));
  if (bytes.length !== ACCOUNT_ADDRESS_LENGTH) {
    throw new Error(`Expected a ${ACCOUNT_ADDRESS_LENGTH}-byte address, got ${bytes.length} bytes`);
  }
  return bytes;
}

/**
 * The one way of writing an Injective address that is stored and compared: bech32 allows an
 * address to be written in upper case as well.
 *
 * @throws if `injectiveAddress` is not the address of an Injective account.
 */
export function normalizeInjectiveAddress(injectiveAddress: string): string {
  getAddressBytes(injectiveAddress);
  return encode(INJECTIVE_PREFIX, decode(injectiveAddress).words);
}

export function getInjectiveEvmAddress(injectiveAddress: string): `0x${string}` {
  return getAddress(`0x${getAddressBytes(injectiveAddress).toString('hex')}`);
}

/**
 * The CCTP mint recipient: the address left-padded with zeros to 32 bytes, in lower-case hex.
 */
export function getInjectiveMintRecipient(injectiveAddress: string): `0x${string}` {
  return pad(`0x${getAddressBytes(injectiveAddress).toString('hex')}`, { size: 32 });
}
