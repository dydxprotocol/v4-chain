import { keccak256, Secp256k1, Secp256k1Signature } from '@cosmjs/crypto';
import { fromHex, toBech32, toHex } from '@cosmjs/encoding';
import { getAddress } from 'viem';

import { ACCOUNT_ADDRESS_LENGTH } from './addresses';

const COMPRESSED_PUBLIC_KEY_LENGTH: number = 33;
const SIGNATURE_COMPONENT_LENGTH: number = 32;
const SECP256K1_ORDER: bigint = BigInt(
  '0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141',
);
const SECP256K1_HALF_ORDER: bigint = SECP256K1_ORDER / BigInt(2);

/**
 * Signs for the account that sponsors gas. Implementations hold or reach the key; nothing else
 * in the module does.
 */
export interface SponsorSigner {
  // Compressed secp256k1 public key of the sponsor account.
  readonly publicKey: Uint8Array,

  /**
   * Signs a 32-byte digest, returning the 64-byte signature r || s.
   */
  signDigest(digest: Uint8Array): Promise<Uint8Array>,
}

/**
 * The part of the Turnkey API client the signer uses.
 */
export interface TurnkeyRawSigner {
  signRawPayload(input: {
    signWith: string,
    payload: string,
    encoding: 'PAYLOAD_ENCODING_HEXADECIMAL',
    hashFunction: 'HASH_FUNCTION_NO_OP',
  }): Promise<{ r: string, s: string }>,
}

function stripHexPrefix(hex: string): string {
  return hex.startsWith('0x') || hex.startsWith('0X') ? hex.slice(2) : hex;
}

export function parseCompressedPublicKey(hex: string): Uint8Array {
  let publicKey: Uint8Array;
  try {
    publicKey = fromHex(stripHexPrefix(hex));
  } catch (error) {
    throw new Error('sponsor public key is not valid hex');
  }
  if (
    publicKey.length !== COMPRESSED_PUBLIC_KEY_LENGTH ||
    (publicKey[0] !== 0x02 && publicKey[0] !== 0x03)
  ) {
    throw new Error('sponsor public key must be a 33-byte compressed secp256k1 key');
  }
  // Throws if the bytes are not a point on the curve.
  Secp256k1.uncompressPubkey(publicKey);
  return publicKey;
}

/**
 * Injective derives addresses the way Ethereum does: the last 20 bytes of the keccak256 hash of
 * the uncompressed public key, without its leading 0x04 byte.
 */
function addressBytesFromPublicKey(compressedPublicKey: Uint8Array): Uint8Array {
  const uncompressed: Uint8Array = Secp256k1.uncompressPubkey(compressedPublicKey);
  return keccak256(uncompressed.slice(1)).slice(-ACCOUNT_ADDRESS_LENGTH);
}

/**
 * The key's Ethereum address, in checksum case.
 */
export function ethereumAddressFromPublicKey(compressedPublicKey: Uint8Array): string {
  return getAddress(`0x${toHex(addressBytesFromPublicKey(compressedPublicKey))}`);
}

export function injectiveAddressFromPublicKey(
  compressedPublicKey: Uint8Array,
  prefix: string,
): string {
  return toBech32(prefix, addressBytesFromPublicKey(compressedPublicKey));
}

function componentToBigInt(hex: string, name: string): bigint {
  const stripped: string = stripHexPrefix(hex);
  if (
    stripped.length === 0 ||
    stripped.length > SIGNATURE_COMPONENT_LENGTH * 2 ||
    !/^[0-9a-fA-F]+$/.test(stripped)
  ) {
    throw new Error(`signature component ${name} is not a 32-byte hex value`);
  }
  return BigInt(`0x${stripped}`);
}

function bigIntToComponent(value: bigint): Uint8Array {
  return fromHex(value.toString(16).padStart(SIGNATURE_COMPONENT_LENGTH * 2, '0'));
}

/**
 * Builds the 64-byte signature r || s from hex components, moving s into the lower half of the
 * curve order. Injective rejects signatures with a high s, and both forms are valid for the
 * same digest and key.
 */
export function toCanonicalSignature(rHex: string, sHex: string): Uint8Array {
  const r: bigint = componentToBigInt(rHex, 'r');
  let s: bigint = componentToBigInt(sHex, 's');
  if (r === BigInt(0) || r >= SECP256K1_ORDER || s === BigInt(0) || s >= SECP256K1_ORDER) {
    throw new Error('signature component out of range');
  }
  if (s > SECP256K1_HALF_ORDER) {
    s = SECP256K1_ORDER - s;
  }
  const signature: Uint8Array = new Uint8Array(SIGNATURE_COMPONENT_LENGTH * 2);
  signature.set(bigIntToComponent(r), 0);
  signature.set(bigIntToComponent(s), SIGNATURE_COMPONENT_LENGTH);
  return signature;
}

/**
 * Whether `signature` (64 bytes, r || s) is a signature of `digest` by `compressedPublicKey`.
 */
export async function verifyDigestSignature(
  signature: Uint8Array,
  digest: Uint8Array,
  compressedPublicKey: Uint8Array,
): Promise<boolean> {
  if (signature.length !== SIGNATURE_COMPONENT_LENGTH * 2) {
    return false;
  }
  try {
    return await Secp256k1.verifySignature(
      Secp256k1Signature.fromFixedLength(signature),
      digest,
      compressedPublicKey,
    );
  } catch (error) {
    return false;
  }
}

/**
 * Signs with a secp256k1 key held in Turnkey. The digest is computed here and sent to Turnkey
 * to sign as-is, so the key never leaves Turnkey.
 */
export class TurnkeySponsorSigner implements SponsorSigner {
  public readonly publicKey: Uint8Array;

  private readonly client: TurnkeyRawSigner;

  private readonly signWith: string;

  /**
   * @param client Turnkey API client for the organisation that holds the key.
   * @param publicKey Compressed public key of the sponsor account.
   * @param signWith How Turnkey identifies the key: a wallet account address, a private key
   * address, or a private key id. Defaults to the key's Ethereum address.
   */
  constructor(client: TurnkeyRawSigner, publicKey: Uint8Array, signWith?: string) {
    this.client = client;
    this.publicKey = publicKey;
    this.signWith = signWith !== undefined && signWith !== ''
      ? signWith
      : ethereumAddressFromPublicKey(publicKey);
  }

  public async signDigest(digest: Uint8Array): Promise<Uint8Array> {
    if (digest.length !== 32) {
      throw new Error(`expected a 32-byte digest, got ${digest.length} bytes`);
    }
    const result: { r: string, s: string } = await this.client.signRawPayload({
      signWith: this.signWith,
      payload: toHex(digest),
      encoding: 'PAYLOAD_ENCODING_HEXADECIMAL',
      hashFunction: 'HASH_FUNCTION_NO_OP',
    });
    return toCanonicalSignature(result.r, result.s);
  }
}
