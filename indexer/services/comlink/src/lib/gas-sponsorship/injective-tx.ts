import { keccak256, sha256 } from '@cosmjs/crypto';
import { toHex } from '@cosmjs/encoding';

import { ProtoWriter } from './protobuf';

export const MSG_GRANT_ALLOWANCE_TYPE_URL: string = '/cosmos.feegrant.v1beta1.MsgGrantAllowance';
export const MSG_REVOKE_ALLOWANCE_TYPE_URL: string = '/cosmos.feegrant.v1beta1.MsgRevokeAllowance';
export const BASIC_ALLOWANCE_TYPE_URL: string = '/cosmos.feegrant.v1beta1.BasicAllowance';
export const ALLOWED_MSG_ALLOWANCE_TYPE_URL: string = '/cosmos.feegrant.v1beta1.AllowedMsgAllowance';
// Injective accounts sign with secp256k1 over a keccak256 digest, and name the key type
// accordingly.
export const ETH_SECP256K1_PUBKEY_TYPE_URL: string = '/injective.crypto.v1beta1.ethsecp256k1.PubKey';

// cosmos.tx.signing.v1beta1.SignMode.SIGN_MODE_DIRECT
const SIGN_MODE_DIRECT: number = 1;
const NANOS_PER_MILLISECOND: number = 1_000_000;

export interface Coin {
  denom: string,
  amount: bigint,
}

export interface ProtoAny {
  typeUrl: string,
  value: Uint8Array,
}

export interface ProtoTimestamp {
  seconds: bigint,
  nanos: number,
}

export interface GrantTxParams {
  chainId: string,
  granter: string,
  grantee: string,
  // Compressed secp256k1 public key of the granter.
  granterPublicKey: Uint8Array,
  accountNumber: bigint,
  sequence: bigint,
  // The most the grantee may spend on fees in total.
  spendLimit: Coin,
  expiration: ProtoTimestamp,
  // Message type URLs the grantee may pay for with the allowance.
  allowedMessages: string[],
  // The fee the granter pays for this transaction.
  fee: Coin,
  gasLimit: bigint,
  // Revoke the granter's current allowance for the grantee first. Needed when one exists,
  // because granting on top of an existing allowance is an error.
  revokeExisting: boolean,
}

export interface UnsignedTx {
  bodyBytes: Uint8Array,
  authInfoBytes: Uint8Array,
  // keccak256 of the SIGN_MODE_DIRECT sign doc. This is what the granter signs.
  signDigest: Uint8Array,
}

export function timestampFromDate(date: Date): ProtoTimestamp {
  const millis: number = date.getTime();
  if (Number.isNaN(millis)) {
    throw new Error('invalid date');
  }
  const seconds: number = Math.floor(millis / 1000);
  return {
    seconds: BigInt(seconds),
    nanos: (millis - seconds * 1000) * NANOS_PER_MILLISECOND,
  };
}

export function encodeAny(message: ProtoAny): Uint8Array {
  return new ProtoWriter()
    .string(1, message.typeUrl)
    .bytes(2, message.value)
    .finish();
}

export function encodeCoin(coin: Coin): Uint8Array {
  return new ProtoWriter()
    .string(1, coin.denom)
    .string(2, coin.amount.toString())
    .finish();
}

export function encodeTimestamp(timestamp: ProtoTimestamp): Uint8Array {
  return new ProtoWriter()
    .uint(1, timestamp.seconds)
    .uint(2, timestamp.nanos)
    .finish();
}

export function encodeBasicAllowance(
  spendLimit: Coin[],
  expiration: ProtoTimestamp,
): Uint8Array {
  const writer: ProtoWriter = new ProtoWriter();
  spendLimit.forEach((coin: Coin) => writer.message(1, encodeCoin(coin)));
  return writer
    .message(2, encodeTimestamp(expiration))
    .finish();
}

export function encodeAllowedMsgAllowance(
  allowance: ProtoAny,
  allowedMessages: string[],
): Uint8Array {
  return new ProtoWriter()
    .message(1, encodeAny(allowance))
    .repeatedString(2, allowedMessages)
    .finish();
}

export function encodeMsgGrantAllowance(
  granter: string,
  grantee: string,
  allowance: ProtoAny,
): Uint8Array {
  return new ProtoWriter()
    .string(1, granter)
    .string(2, grantee)
    .message(3, encodeAny(allowance))
    .finish();
}

export function encodeMsgRevokeAllowance(granter: string, grantee: string): Uint8Array {
  return new ProtoWriter()
    .string(1, granter)
    .string(2, grantee)
    .finish();
}

export function encodeTxBody(messages: ProtoAny[], memo: string = ''): Uint8Array {
  const writer: ProtoWriter = new ProtoWriter();
  messages.forEach((message: ProtoAny) => writer.message(1, encodeAny(message)));
  return writer
    .string(2, memo)
    .finish();
}

/**
 * Encodes the AuthInfo for a transaction with a single SIGN_MODE_DIRECT signer who pays the fee.
 */
export function encodeAuthInfo(
  publicKey: ProtoAny,
  sequence: bigint,
  feeAmount: Coin[],
  gasLimit: bigint,
): Uint8Array {
  const modeInfo: Uint8Array = new ProtoWriter()
    .message(1, new ProtoWriter().uint(1, SIGN_MODE_DIRECT).finish())
    .finish();
  const signerInfo: Uint8Array = new ProtoWriter()
    .message(1, encodeAny(publicKey))
    .message(2, modeInfo)
    .uint(3, sequence)
    .finish();

  const feeWriter: ProtoWriter = new ProtoWriter();
  feeAmount.forEach((coin: Coin) => feeWriter.message(1, encodeCoin(coin)));
  const fee: Uint8Array = feeWriter
    .uint(2, gasLimit)
    .finish();

  return new ProtoWriter()
    .message(1, signerInfo)
    .message(2, fee)
    .finish();
}

export function encodeSignDoc(
  bodyBytes: Uint8Array,
  authInfoBytes: Uint8Array,
  chainId: string,
  accountNumber: bigint,
): Uint8Array {
  return new ProtoWriter()
    .bytes(1, bodyBytes)
    .bytes(2, authInfoBytes)
    .string(3, chainId)
    .uint(4, accountNumber)
    .finish();
}

export function encodeTxRaw(
  bodyBytes: Uint8Array,
  authInfoBytes: Uint8Array,
  signatures: Uint8Array[],
): Uint8Array {
  const writer: ProtoWriter = new ProtoWriter()
    .bytes(1, bodyBytes)
    .bytes(2, authInfoBytes);
  signatures.forEach((signature: Uint8Array) => writer.message(3, signature));
  return writer.finish();
}

export function encodeEthSecp256k1PubKey(compressedPublicKey: Uint8Array): ProtoAny {
  if (compressedPublicKey.length !== 33) {
    throw new Error(
      `expected a 33-byte compressed public key, got ${compressedPublicKey.length} bytes`,
    );
  }
  return {
    typeUrl: ETH_SECP256K1_PUBKEY_TYPE_URL,
    value: new ProtoWriter().bytes(1, compressedPublicKey).finish(),
  };
}

/**
 * The digest an Injective account signs for a SIGN_MODE_DIRECT transaction: keccak256 of the
 * sign doc, where other Cosmos chains use sha256.
 */
export function getSignDigest(
  bodyBytes: Uint8Array,
  authInfoBytes: Uint8Array,
  chainId: string,
  accountNumber: bigint,
): Uint8Array {
  return keccak256(encodeSignDoc(bodyBytes, authInfoBytes, chainId, accountNumber));
}

/**
 * Builds the transaction that gives `grantee` a fee allowance from `granter`, limited to
 * `spendLimit` in total, to `allowedMessages`, and until `expiration`.
 */
export function buildGrantTx(params: GrantTxParams): UnsignedTx {
  const basicAllowance: ProtoAny = {
    typeUrl: BASIC_ALLOWANCE_TYPE_URL,
    value: encodeBasicAllowance([params.spendLimit], params.expiration),
  };
  const allowance: ProtoAny = {
    typeUrl: ALLOWED_MSG_ALLOWANCE_TYPE_URL,
    value: encodeAllowedMsgAllowance(basicAllowance, params.allowedMessages),
  };

  const messages: ProtoAny[] = [];
  if (params.revokeExisting) {
    messages.push({
      typeUrl: MSG_REVOKE_ALLOWANCE_TYPE_URL,
      value: encodeMsgRevokeAllowance(params.granter, params.grantee),
    });
  }
  messages.push({
    typeUrl: MSG_GRANT_ALLOWANCE_TYPE_URL,
    value: encodeMsgGrantAllowance(params.granter, params.grantee, allowance),
  });

  const bodyBytes: Uint8Array = encodeTxBody(messages);
  const authInfoBytes: Uint8Array = encodeAuthInfo(
    encodeEthSecp256k1PubKey(params.granterPublicKey),
    params.sequence,
    [params.fee],
    params.gasLimit,
  );

  return {
    bodyBytes,
    authInfoBytes,
    signDigest: getSignDigest(bodyBytes, authInfoBytes, params.chainId, params.accountNumber),
  };
}

/**
 * Combines an unsigned transaction with its 64-byte signature into the bytes to broadcast.
 */
export function assembleSignedTx(unsignedTx: UnsignedTx, signature: Uint8Array): Uint8Array {
  return encodeTxRaw(unsignedTx.bodyBytes, unsignedTx.authInfoBytes, [signature]);
}

/**
 * The hash a transaction is known by on chain, computable before it is broadcast.
 */
export function getTxHash(txBytes: Uint8Array): string {
  return toHex(sha256(txBytes)).toUpperCase();
}
