import { toBase64 } from '@cosmjs/encoding';
import fetch, { Response } from 'node-fetch';

import {
  ALLOWED_MSG_ALLOWANCE_TYPE_URL,
  BASIC_ALLOWANCE_TYPE_URL,
} from './injective-tx';

// The Cosmos SDK's error for a transaction signed with the wrong account sequence.
const SDK_CODESPACE: string = 'sdk';
const SDK_CODE_WRONG_SEQUENCE: number = 32;
// gRPC status codes, as reported in the body of a failed REST query.
const GRPC_CODE_NOT_FOUND: number = 5;

export class InjectiveClientError extends Error {
  public readonly status: number | undefined;

  constructor(message: string, status?: number) {
    super(message);
    this.name = 'InjectiveClientError';
    this.status = status;
  }
}

export interface InjectiveAccount {
  accountNumber: bigint,
  sequence: bigint,
}

export interface FeeAllowance {
  granter: string,
  grantee: string,
  // False when the allowance is of a kind this module does not issue or understand. The
  // remaining fields are then unset and say nothing about what the allowance permits.
  recognized: boolean,
  // What is left to spend, by denom. Undefined means there is no limit.
  spendLimit?: { [denom: string]: bigint },
  // Undefined means the allowance does not expire.
  expiration?: Date,
  // Undefined means the allowance covers every message type.
  allowedMessages?: string[],
}

export interface TxResult {
  txHash: string,
  // Zero means success.
  code: number,
  codespace: string,
  rawLog: string,
}

export interface IncludedTx extends TxResult {
  height: bigint,
}

interface RestResponse {
  status: number,
  ok: boolean,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  body: any,
}

function parseUint(value: unknown, what: string): bigint {
  if (typeof value !== 'string' || !/^\d+$/.test(value)) {
    throw new InjectiveClientError(`unexpected ${what}: ${JSON.stringify(value)}`);
  }
  return BigInt(value);
}

/**
 * Parses a non-negative decimal such as "160000000.000000000000000000", rounding up to a whole
 * number.
 */
export function parseDecimalRoundingUp(value: unknown, what: string): bigint {
  if (typeof value !== 'string' || !/^\d+(\.\d+)?$/.test(value)) {
    throw new InjectiveClientError(`unexpected ${what}: ${JSON.stringify(value)}`);
  }
  const [whole, fraction]: string[] = value.split('.');
  const hasFraction: boolean = fraction !== undefined && /[1-9]/.test(fraction);
  return BigInt(whole) + (hasFraction ? BigInt(1) : BigInt(0));
}

/**
 * If `result` is the rejection of a transaction signed with the wrong sequence, returns the
 * sequence the chain expected.
 */
export function getExpectedSequence(result: TxResult): bigint | undefined {
  if (result.codespace !== SDK_CODESPACE || result.code !== SDK_CODE_WRONG_SEQUENCE) {
    return undefined;
  }
  const match: RegExpMatchArray | null = result.rawLog.match(/expected (\d+), got \d+/);
  return match === null ? undefined : BigInt(match[1]);
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function parseSpendLimit(coins: any): { [denom: string]: bigint } | undefined {
  // An absent or empty spend limit means the allowance is unlimited.
  if (!Array.isArray(coins) || coins.length === 0) {
    return undefined;
  }
  const spendLimit: { [denom: string]: bigint } = {};
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  coins.forEach((coin: any) => {
    spendLimit[String(coin?.denom)] = parseUint(coin?.amount, 'allowance amount');
  });
  return spendLimit;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function parseExpiration(expiration: any): Date | undefined {
  if (expiration === undefined || expiration === null) {
    return undefined;
  }
  const date: Date = new Date(String(expiration));
  if (Number.isNaN(date.getTime())) {
    throw new InjectiveClientError(`unexpected allowance expiration: ${JSON.stringify(expiration)}`);
  }
  return date;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function parseFeeAllowance(grant: any): FeeAllowance {
  const base: { granter: string, grantee: string } = {
    granter: String(grant?.granter),
    grantee: String(grant?.grantee),
  };
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let allowance: any = grant?.allowance;
  let allowedMessages: string[] | undefined;

  if (allowance?.['@type'] === ALLOWED_MSG_ALLOWANCE_TYPE_URL) {
    if (!Array.isArray(allowance.allowed_messages)) {
      return { ...base, recognized: false };
    }
    allowedMessages = allowance.allowed_messages.map(String);
    allowance = allowance.allowance;
  }
  if (allowance?.['@type'] !== BASIC_ALLOWANCE_TYPE_URL) {
    return { ...base, recognized: false };
  }

  return {
    ...base,
    recognized: true,
    spendLimit: parseSpendLimit(allowance.spend_limit),
    expiration: parseExpiration(allowance.expiration),
    allowedMessages,
  };
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function parseTxResponse(txResponse: any): TxResult {
  if (typeof txResponse?.txhash !== 'string' || typeof txResponse?.code !== 'number') {
    throw new InjectiveClientError(`unexpected tx response: ${JSON.stringify(txResponse)}`);
  }
  return {
    txHash: txResponse.txhash.toUpperCase(),
    code: txResponse.code,
    codespace: typeof txResponse.codespace === 'string' ? txResponse.codespace : '',
    rawLog: typeof txResponse.raw_log === 'string' ? txResponse.raw_log : '',
  };
}

/**
 * Reads from and broadcasts to Injective through a node's REST (LCD) API.
 */
export class InjectiveClient {
  private readonly restUrl: string;

  private readonly timeoutMs: number;

  constructor(restUrl: string, timeoutMs: number) {
    this.restUrl = restUrl.replace(/\/+$/, '');
    this.timeoutMs = timeoutMs;
  }

  /**
   * Returns the account's number and sequence, or undefined if the account does not exist on
   * chain yet.
   */
  public async getAccount(address: string): Promise<InjectiveAccount | undefined> {
    const response: RestResponse = await this.request(
      `/cosmos/auth/v1beta1/accounts/${encodeURIComponent(address)}`,
    );
    if (response.status === 404) {
      return undefined;
    }
    this.assertOk(response, 'account query');
    // Injective's own account type nests the standard account; plain accounts do not.
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const account: any = response.body?.account?.base_account ?? response.body?.account;
    return {
      accountNumber: parseUint(account?.account_number, 'account number'),
      sequence: parseUint(account?.sequence, 'account sequence'),
    };
  }

  public async getBalance(address: string, denom: string): Promise<bigint> {
    const response: RestResponse = await this.request(
      `/cosmos/bank/v1beta1/balances/${encodeURIComponent(address)}/by_denom` +
      `?denom=${encodeURIComponent(denom)}`,
    );
    this.assertOk(response, 'balance query');
    return parseUint(response.body?.balance?.amount, 'balance');
  }

  /**
   * Returns the fee allowance `granter` has given `grantee`, or undefined if there is none.
   */
  public async getAllowance(granter: string, grantee: string): Promise<FeeAllowance | undefined> {
    const response: RestResponse = await this.request(
      `/cosmos/feegrant/v1beta1/allowance/${encodeURIComponent(granter)}` +
      `/${encodeURIComponent(grantee)}`,
    );
    // A missing allowance is reported as an internal error whose message says "not found",
    // rather than with a not-found status.
    if (
      response.status === 404 ||
      response.body?.code === GRPC_CODE_NOT_FOUND ||
      (!response.ok && /not found/i.test(String(response.body?.message)))
    ) {
      return undefined;
    }
    this.assertOk(response, 'allowance query');
    return parseFeeAllowance(response.body?.allowance);
  }

  /**
   * Returns the current base fee per unit of gas, in the smallest unit of the fee denom, rounded
   * up to a whole number.
   */
  public async getBaseFee(): Promise<bigint> {
    const response: RestResponse = await this.request('/injective/txfees/v1beta1/cur_eip_base_fee');
    this.assertOk(response, 'base fee query');
    return parseDecimalRoundingUp(response.body?.base_fee?.base_fee, 'base fee');
  }

  /**
   * Submits a signed transaction and returns the node's verdict on admitting it to the mempool.
   * A zero code means it was admitted, not that it has been included in a block.
   */
  public async broadcastTx(txBytes: Uint8Array): Promise<TxResult> {
    const response: RestResponse = await this.request('/cosmos/tx/v1beta1/txs', {
      tx_bytes: toBase64(txBytes),
      mode: 'BROADCAST_MODE_SYNC',
    });
    this.assertOk(response, 'broadcast');
    return parseTxResponse(response.body?.tx_response);
  }

  /**
   * Returns the transaction if it has been included in a block, or undefined if the node does
   * not know of it as included.
   */
  public async getTx(txHash: string): Promise<IncludedTx | undefined> {
    const response: RestResponse = await this.request(
      `/cosmos/tx/v1beta1/txs/${encodeURIComponent(txHash)}`,
    );
    if (response.status === 404 || response.body?.code === GRPC_CODE_NOT_FOUND) {
      return undefined;
    }
    this.assertOk(response, 'tx query');
    return {
      ...parseTxResponse(response.body?.tx_response),
      height: parseUint(response.body?.tx_response?.height, 'tx height'),
    };
  }

  private assertOk(response: RestResponse, what: string): void {
    if (!response.ok) {
      throw new InjectiveClientError(
        `${what} failed with status ${response.status}: ${JSON.stringify(response.body)}`,
        response.status,
      );
    }
  }

  private async request(path: string, postBody?: object): Promise<RestResponse> {
    let response: Response;
    try {
      response = await fetch(`${this.restUrl}${path}`, {
        method: postBody === undefined ? 'GET' : 'POST',
        headers: {
          accept: 'application/json',
          ...(postBody === undefined ? {} : { 'content-type': 'application/json' }),
        },
        body: postBody === undefined ? undefined : JSON.stringify(postBody),
        timeout: this.timeoutMs,
      });
    } catch (error) {
      throw new InjectiveClientError(`request to ${path} failed: ${error.message}`);
    }

    const text: string = await response.text();
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    let body: any;
    try {
      body = JSON.parse(text);
    } catch (error) {
      if (response.ok) {
        throw new InjectiveClientError(`response from ${path} is not JSON`, response.status);
      }
      body = { message: text.slice(0, 200) };
    }
    return { status: response.status, ok: response.ok, body };
  }
}

/**
 * Whether the dYdX chain is currently accepting USDC deposits from Injective, according to the
 * node at `dydxRestUrl`. A deposit sent while it is not would be refunded on Injective, having
 * spent the gas for nothing.
 */
export async function isInjectiveDepositOpen(
  dydxRestUrl: string,
  timeoutMs: number,
): Promise<boolean> {
  let response: Response;
  try {
    response = await fetch(
      `${dydxRestUrl.replace(/\/+$/, '')}/dydxprotocol/v4/canonicalusdc/state`,
      { headers: { accept: 'application/json' }, timeout: timeoutMs },
    );
  } catch (error) {
    throw new InjectiveClientError(`canonical USDC state query failed: ${error.message}`);
  }
  if (!response.ok) {
    throw new InjectiveClientError(
      `canonical USDC state query failed with status ${response.status}`,
      response.status,
    );
  }
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const body: any = await response.json();
  if (typeof body?.controls?.mode !== 'string') {
    throw new InjectiveClientError('unexpected canonical USDC state');
  }
  // Injective deposits are only credited as USDC while the module is in its gradual mode, and
  // can be switched off separately within it.
  return body.controls.mode === 'MODE_GRADUAL' &&
    body.controls.injective_deposits_disabled !== true;
}
