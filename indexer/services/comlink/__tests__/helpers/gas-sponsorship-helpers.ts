import { Secp256k1 } from '@cosmjs/crypto';

import {
  FeeAllowance,
  IncludedTx,
  InjectiveAccount,
  InjectiveClientError,
  TxResult,
} from '../../src/lib/gas-sponsorship/injective-client';
import {
  ALLOWED_MSG_ALLOWANCE_TYPE_URL,
  BASIC_ALLOWANCE_TYPE_URL,
  ETH_SECP256K1_PUBKEY_TYPE_URL,
  getSignDigest,
  getTxHash,
  MSG_GRANT_ALLOWANCE_TYPE_URL,
  MSG_REVOKE_ALLOWANCE_TYPE_URL,
} from '../../src/lib/gas-sponsorship/injective-tx';
import {
  injectiveAddressFromPublicKey,
  SponsorSigner,
  verifyDigestSignature,
} from '../../src/lib/gas-sponsorship/sponsor-signer';
import { SponsorshipChainClient } from '../../src/lib/gas-sponsorship/types';

/**
 * Signs with a private key held in memory. For tests only: the service is never given a raw
 * private key outside of them.
 */
export class LocalKeySigner implements SponsorSigner {
  public signCount: number = 0;

  private constructor(
    public readonly publicKey: Uint8Array,
    private readonly privateKey: Uint8Array,
  ) {}

  public static async create(privateKeyHex: string): Promise<LocalKeySigner> {
    const privateKey: Uint8Array = Uint8Array.from(Buffer.from(privateKeyHex, 'hex'));
    const keypair: { pubkey: Uint8Array } = await Secp256k1.makeKeypair(privateKey);
    return new LocalKeySigner(Secp256k1.compressPubkey(keypair.pubkey), privateKey);
  }

  public async signDigest(digest: Uint8Array): Promise<Uint8Array> {
    this.signCount += 1;
    const signature: { toFixedLength(): Uint8Array } = await Secp256k1.createSignature(
      digest,
      this.privateKey,
    );
    // Drop the trailing recovery byte.
    return signature.toFixedLength().slice(0, 64);
  }
}

/* ------- A protobuf reader, so that the fake chain reads transactions for itself ------- */

interface ProtoField {
  fieldNumber: number,
  // Set for varint fields.
  varint?: bigint,
  // Set for length-delimited fields.
  bytes?: Uint8Array,
}

function readProto(buffer: Uint8Array): ProtoField[] {
  const fields: ProtoField[] = [];
  let offset: number = 0;

  const readVarint = (): bigint => {
    let result: bigint = BigInt(0);
    let shift: bigint = BigInt(0);
    for (;;) {
      if (offset >= buffer.length) {
        throw new Error('truncated varint');
      }
      const byte: number = buffer[offset];
      offset += 1;
      // eslint-disable-next-line no-bitwise
      result |= BigInt(byte & 0x7f) << shift;
      // eslint-disable-next-line no-bitwise
      if ((byte & 0x80) === 0) {
        return result;
      }
      shift += BigInt(7);
    }
  };

  while (offset < buffer.length) {
    const tag: number = Number(readVarint());
    // eslint-disable-next-line no-bitwise
    const fieldNumber: number = tag >> 3;
    // eslint-disable-next-line no-bitwise
    const wireType: number = tag & 0x7;
    if (wireType === 0) {
      fields.push({ fieldNumber, varint: readVarint() });
    } else if (wireType === 2) {
      const length: number = Number(readVarint());
      if (offset + length > buffer.length) {
        throw new Error('truncated field');
      }
      fields.push({ fieldNumber, bytes: buffer.slice(offset, offset + length) });
      offset += length;
    } else {
      throw new Error(`unsupported wire type ${wireType}`);
    }
  }
  return fields;
}

function allBytes(fields: ProtoField[], fieldNumber: number): Uint8Array[] {
  return fields
    .filter((field: ProtoField) => field.fieldNumber === fieldNumber && field.bytes !== undefined)
    .map((field: ProtoField) => field.bytes!);
}

function oneBytes(fields: ProtoField[], fieldNumber: number): Uint8Array {
  const found: Uint8Array[] = allBytes(fields, fieldNumber);
  if (found.length > 1) {
    throw new Error(`field ${fieldNumber} appears ${found.length} times`);
  }
  return found[0] ?? new Uint8Array(0);
}

function oneString(fields: ProtoField[], fieldNumber: number): string {
  return Buffer.from(oneBytes(fields, fieldNumber)).toString('utf8');
}

function oneVarint(fields: ProtoField[], fieldNumber: number): bigint {
  const found: ProtoField | undefined = fields.find(
    (field: ProtoField) => field.fieldNumber === fieldNumber && field.varint !== undefined,
  );
  return found?.varint ?? BigInt(0);
}

function readAny(buffer: Uint8Array): { typeUrl: string, value: Uint8Array } {
  const fields: ProtoField[] = readProto(buffer);
  return { typeUrl: oneString(fields, 1), value: oneBytes(fields, 2) };
}

function readCoins(buffers: Uint8Array[]): { [denom: string]: bigint } {
  const coins: { [denom: string]: bigint } = {};
  buffers.forEach((buffer: Uint8Array) => {
    const fields: ProtoField[] = readProto(buffer);
    coins[oneString(fields, 1)] = BigInt(oneString(fields, 2));
  });
  return coins;
}

/* ------- What the fake chain reads out of a transaction ------- */

export interface DecodedGrant {
  granter: string,
  grantee: string,
  spendLimit: { [denom: string]: bigint },
  expiration: Date,
  allowedMessages: string[],
}

export type DecodedMessage =
  { type: 'revoke', granter: string, grantee: string } |
  ({ type: 'grant' } & DecodedGrant);

export interface DecodedTx {
  txHash: string,
  messages: DecodedMessage[],
  memo: string,
  signerPublicKey: Uint8Array,
  sequence: bigint,
  fee: { [denom: string]: bigint },
  gasLimit: bigint,
  feePayer: string,
  feeGranter: string,
  bodyBytes: Uint8Array,
  authInfoBytes: Uint8Array,
  signature: Uint8Array,
}

function decodeMessage(packed: { typeUrl: string, value: Uint8Array }): DecodedMessage {
  const fields: ProtoField[] = readProto(packed.value);
  if (packed.typeUrl === MSG_REVOKE_ALLOWANCE_TYPE_URL) {
    return { type: 'revoke', granter: oneString(fields, 1), grantee: oneString(fields, 2) };
  }
  if (packed.typeUrl !== MSG_GRANT_ALLOWANCE_TYPE_URL) {
    throw new Error(`unexpected message type ${packed.typeUrl}`);
  }

  const allowedMsgAllowance: { typeUrl: string, value: Uint8Array } = readAny(oneBytes(fields, 3));
  if (allowedMsgAllowance.typeUrl !== ALLOWED_MSG_ALLOWANCE_TYPE_URL) {
    throw new Error(`unexpected allowance type ${allowedMsgAllowance.typeUrl}`);
  }
  const allowedFields: ProtoField[] = readProto(allowedMsgAllowance.value);
  const basicAllowance: { typeUrl: string, value: Uint8Array } = readAny(
    oneBytes(allowedFields, 1),
  );
  if (basicAllowance.typeUrl !== BASIC_ALLOWANCE_TYPE_URL) {
    throw new Error(`unexpected inner allowance type ${basicAllowance.typeUrl}`);
  }
  const basicFields: ProtoField[] = readProto(basicAllowance.value);
  const timestamp: ProtoField[] = readProto(oneBytes(basicFields, 2));

  return {
    type: 'grant',
    granter: oneString(fields, 1),
    grantee: oneString(fields, 2),
    spendLimit: readCoins(allBytes(basicFields, 1)),
    expiration: new Date(
      Number(oneVarint(timestamp, 1)) * 1000 + Number(oneVarint(timestamp, 2)) / 1_000_000,
    ),
    allowedMessages: allBytes(allowedFields, 2).map(
      (message: Uint8Array) => Buffer.from(message).toString('utf8'),
    ),
  };
}

export function decodeTx(txBytes: Uint8Array): DecodedTx {
  const raw: ProtoField[] = readProto(txBytes);
  const bodyBytes: Uint8Array = oneBytes(raw, 1);
  const authInfoBytes: Uint8Array = oneBytes(raw, 2);
  const body: ProtoField[] = readProto(bodyBytes);
  const authInfo: ProtoField[] = readProto(authInfoBytes);
  const signerInfo: ProtoField[] = readProto(oneBytes(authInfo, 1));
  const fee: ProtoField[] = readProto(oneBytes(authInfo, 2));

  const publicKey: { typeUrl: string, value: Uint8Array } = readAny(oneBytes(signerInfo, 1));
  if (publicKey.typeUrl !== ETH_SECP256K1_PUBKEY_TYPE_URL) {
    throw new Error(`unexpected public key type ${publicKey.typeUrl}`);
  }
  const signMode: bigint = oneVarint(readProto(oneBytes(readProto(oneBytes(signerInfo, 2)), 1)), 1);
  if (signMode !== BigInt(1)) {
    throw new Error(`unexpected sign mode ${signMode}`);
  }

  return {
    txHash: getTxHash(txBytes),
    messages: allBytes(body, 1).map((message: Uint8Array) => decodeMessage(readAny(message))),
    memo: oneString(body, 2),
    signerPublicKey: oneBytes(readProto(publicKey.value), 1),
    sequence: oneVarint(signerInfo, 3),
    fee: readCoins(allBytes(fee, 1)),
    gasLimit: oneVarint(fee, 2),
    feePayer: oneString(fee, 3),
    feeGranter: oneString(fee, 4),
    bodyBytes,
    authInfoBytes,
    signature: oneBytes(raw, 3),
  };
}

/* ------- The fake chain ------- */

const FEE_DENOM: string = 'inj';

/**
 * An in-memory stand-in for an Injective node, implementing the behaviour the sponsorship
 * service depends on: signature and sequence checks when a transaction is submitted, a mempool,
 * and the fee-grant module's rules when a block is made.
 *
 * It reads submitted transactions with its own decoder, so what the tests observe is what the
 * service actually put on the wire.
 */
export class FakeInjectiveChain implements SponsorshipChainClient {
  public readonly chainId: string = 'injective-1';

  public baseFee: bigint = BigInt(160_000_000);

  // When true, a submitted transaction is included in a block straight away.
  public autoInclude: boolean = true;

  // Every transaction submitted, accepted or not, in order.
  public readonly submitted: DecodedTx[] = [];

  // Set to make the next call of the named method fail as an unreachable node would.
  public failNext: { [method: string]: boolean } = {};

  // Set to make every transaction lookup fail as an unreachable node would.
  public failGetTx: boolean = false;

  // Set to make the next submission be taken into the mempool but report an error to the
  // caller, as when the response is lost.
  public loseNextBroadcastResponse: boolean = false;

  // Set to make the node refuse the next submission with this result.
  public rejectNextBroadcast: { code: number, rawLog: string } | undefined;

  // Called just before each block is made, to change state between a transaction being
  // submitted and being executed.
  public beforeBlock: (() => void) | undefined;

  private readonly accounts: Map<string, InjectiveAccount & { publicKey?: Uint8Array }> = new Map();

  private readonly balances: Map<string, bigint> = new Map();

  private readonly allowances: Map<string, FeeAllowance> = new Map();

  private readonly included: Map<string, IncludedTx> = new Map();

  private mempool: DecodedTx[] = [];

  private height: bigint = BigInt(1000);

  private nextAccountNumber: bigint = BigInt(7);

  public addAccount(publicKey: Uint8Array, sequence: bigint = BigInt(0)): string {
    const address: string = injectiveAddressFromPublicKey(publicKey, 'inj');
    this.accounts.set(address, { accountNumber: this.nextAccountNumber, sequence, publicKey });
    this.nextAccountNumber += BigInt(1);
    return address;
  }

  public setBalance(address: string, denom: string, amount: bigint): void {
    this.balances.set(`${address}/${denom}`, amount);
  }

  public setAllowance(allowance: FeeAllowance): void {
    this.allowances.set(`${allowance.granter}/${allowance.grantee}`, allowance);
  }

  public removeAllowance(granter: string, grantee: string): void {
    this.allowances.delete(`${granter}/${grantee}`);
  }

  public getMempoolSize(): number {
    return this.mempool.length;
  }

  /**
   * Makes a block out of everything in the mempool.
   */
  public includePending(): void {
    if (this.beforeBlock !== undefined) {
      this.beforeBlock();
    }
    this.height += BigInt(1);
    const txs: DecodedTx[] = this.mempool;
    this.mempool = [];
    txs.forEach((tx: DecodedTx) => this.deliver(tx));
  }

  /**
   * Drops everything in the mempool, as a node does with transactions it gives up on.
   */
  public dropPending(): void {
    this.mempool = [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async getAccount(address: string): Promise<InjectiveAccount | undefined> {
    this.maybeFail('getAccount');
    const account: InjectiveAccount | undefined = this.accounts.get(address);
    return account === undefined
      ? undefined
      : { accountNumber: account.accountNumber, sequence: account.sequence };
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async getBalance(address: string, denom: string): Promise<bigint> {
    this.maybeFail('getBalance');
    return this.balances.get(`${address}/${denom}`) ?? BigInt(0);
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async getAllowance(granter: string, grantee: string): Promise<FeeAllowance | undefined> {
    this.maybeFail('getAllowance');
    return this.allowances.get(`${granter}/${grantee}`);
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async getBaseFee(): Promise<bigint> {
    this.maybeFail('getBaseFee');
    return this.baseFee;
  }

  public async broadcastTx(txBytes: Uint8Array): Promise<TxResult> {
    this.maybeFail('broadcastTx');
    const tx: DecodedTx = decodeTx(txBytes);
    this.submitted.push(tx);

    if (this.rejectNextBroadcast !== undefined) {
      const { code, rawLog }: { code: number, rawLog: string } = this.rejectNextBroadcast;
      this.rejectNextBroadcast = undefined;
      return {
        txHash: tx.txHash, code, codespace: 'sdk', rawLog,
      };
    }
    const rejection: TxResult | undefined = await this.check(tx);
    if (rejection !== undefined) {
      return rejection;
    }
    this.mempool.push(tx);
    if (this.autoInclude) {
      this.includePending();
    }
    if (this.loseNextBroadcastResponse) {
      this.loseNextBroadcastResponse = false;
      throw new InjectiveClientError('request to /cosmos/tx/v1beta1/txs failed: socket hang up');
    }
    return {
      txHash: tx.txHash, code: 0, codespace: '', rawLog: '',
    };
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async getTx(txHash: string): Promise<IncludedTx | undefined> {
    if (this.failGetTx) {
      throw new InjectiveClientError('tx query failed with status 502', 502);
    }
    return this.included.get(txHash);
  }

  /**
   * The checks a node makes before taking a transaction into its mempool.
   */
  private async check(tx: DecodedTx): Promise<TxResult | undefined> {
    const reject = (code: number, rawLog: string): TxResult => ({
      txHash: tx.txHash, code, codespace: 'sdk', rawLog,
    });

    const signerAddress: string = injectiveAddressFromPublicKey(tx.signerPublicKey, 'inj');
    const account: InjectiveAccount | undefined = this.accounts.get(signerAddress);
    if (account === undefined) {
      return reject(9, `account ${signerAddress} not found: unknown address`);
    }
    const pending: number = this.mempool.filter(
      (queued: DecodedTx) => Buffer.from(queued.signerPublicKey).equals(
        Buffer.from(tx.signerPublicKey),
      ),
    ).length;
    const expectedSequence: bigint = account.sequence + BigInt(pending);
    if (tx.sequence !== expectedSequence) {
      return reject(
        32,
        `account sequence mismatch, expected ${expectedSequence}, got ${tx.sequence}: ` +
        'incorrect account sequence',
      );
    }
    const digest: Uint8Array = getSignDigest(
      tx.bodyBytes,
      tx.authInfoBytes,
      this.chainId,
      account.accountNumber,
    );
    if (!await verifyDigestSignature(tx.signature, digest, tx.signerPublicKey)) {
      return reject(4, 'signature verification failed; please verify account number: unauthorized');
    }
    if ((tx.fee[FEE_DENOM] ?? BigInt(0)) < tx.gasLimit * this.baseFee) {
      return reject(13, 'insufficient fees: insufficient fee');
    }
    if ((this.balances.get(`${signerAddress}/${FEE_DENOM}`) ?? BigInt(0)) < tx.fee[FEE_DENOM]) {
      return reject(5, 'insufficient funds');
    }
    return undefined;
  }

  /**
   * Executes a transaction in a block. The fee and the sequence are consumed whether or not the
   * messages succeed, and the messages take effect together or not at all.
   */
  private deliver(tx: DecodedTx): void {
    const signerAddress: string = injectiveAddressFromPublicKey(tx.signerPublicKey, 'inj');
    const account: InjectiveAccount = this.accounts.get(signerAddress)!;
    account.sequence += BigInt(1);
    const balanceKey: string = `${signerAddress}/${FEE_DENOM}`;
    this.balances.set(balanceKey, this.balances.get(balanceKey)! - tx.fee[FEE_DENOM]);

    const staged: Map<string, FeeAllowance | undefined> = new Map();
    const current = (key: string): FeeAllowance | undefined => (
      staged.has(key) ? staged.get(key) : this.allowances.get(key)
    );
    let failure: string | undefined;

    for (const message of tx.messages) {
      const key: string = `${message.granter}/${message.grantee}`;
      if (message.granter !== signerAddress) {
        failure = 'pubKey does not match signer address: invalid pubkey';
      } else if (message.type === 'revoke') {
        if (current(key) === undefined) {
          failure = 'fee-grant not found: not found';
        }
        staged.set(key, undefined);
      } else {
        if (current(key) !== undefined) {
          failure = 'fee allowance already exists';
        }
        staged.set(key, {
          granter: message.granter,
          grantee: message.grantee,
          recognized: true,
          spendLimit: message.spendLimit,
          expiration: message.expiration,
          allowedMessages: message.allowedMessages,
        });
      }
      if (failure !== undefined) {
        break;
      }
    }

    if (failure === undefined) {
      staged.forEach((allowance: FeeAllowance | undefined, key: string) => {
        if (allowance === undefined) {
          this.allowances.delete(key);
        } else {
          this.allowances.set(key, allowance);
        }
      });
    }
    this.included.set(tx.txHash, {
      txHash: tx.txHash,
      code: failure === undefined ? 0 : 2,
      codespace: failure === undefined ? '' : 'feegrant',
      rawLog: failure ?? '',
      height: this.height,
    });
  }

  private maybeFail(method: string): void {
    if (this.failNext[method]) {
      this.failNext[method] = false;
      throw new InjectiveClientError(`${method} failed with status 502`, 502);
    }
  }
}
