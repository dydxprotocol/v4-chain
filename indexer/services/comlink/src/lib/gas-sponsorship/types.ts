import {
  FeeAllowance, IncludedTx, InjectiveAccount, TxResult,
} from './injective-client';

export interface GasSponsorshipSettings {
  // The kill switch. Nothing is read, signed or sent while this is false.
  enabled: boolean,
  chainId: string,
  // bech32 prefix of account addresses on the chain.
  addressPrefix: string,
  // Denom fees are paid in, and so the denom of the allowance.
  feeDenom: string,
  // Denom of the USDC an address must hold to be sponsored.
  usdcDenom: string,
  // Message type URLs the allowance may pay for.
  allowedMessages: string[],
  // Gas limit one sponsored transaction is expected to set.
  transferGasLimit: bigint,
  // How many sponsored transactions the allowance is sized for.
  transferAttempts: bigint,
  // Gas price the allowance and the grant transaction are priced at, as a percentage of the
  // current base fee. Above 100 leaves room for the base fee to rise.
  gasPriceBufferPercent: bigint,
  // No allowance is given while the base fee per unit of gas is above this. This bounds the size
  // of an allowance.
  maxBaseFee: bigint,
  // Gas limit of the grant transaction itself.
  grantGasLimit: bigint,
  grantExpirySeconds: number,
  // The least USDC, in its smallest unit, an address must hold to be sponsored.
  minUsdcBalance: bigint,
  // How long after being given an allowance an address must wait for another.
  cooldownSeconds: number,
  // The most that may be committed in any 24 hours, in the smallest unit of the fee denom,
  // counting each allowance in full plus the fee paid to grant it.
  dailyBudget: bigint,
  // How long the lock that serializes the sponsor's transactions is held before it lapses. Must
  // exceed the longest a single grant can take.
  lockExpiryMs: number,
  // How long a request waits for the lock before giving up.
  lockWaitMs: number,
  // How long to wait for the grant transaction to be included before reporting it as pending.
  confirmationTimeoutMs: number,
  confirmationPollMs: number,
}

/**
 * The chain reads and writes the service needs. Implemented by InjectiveClient.
 */
export interface SponsorshipChainClient {
  getAccount(address: string): Promise<InjectiveAccount | undefined>,
  getBalance(address: string, denom: string): Promise<bigint>,
  getAllowance(granter: string, grantee: string): Promise<FeeAllowance | undefined>,
  getBaseFee(): Promise<bigint>,
  broadcastTx(txBytes: Uint8Array): Promise<TxResult>,
  getTx(txHash: string): Promise<IncludedTx | undefined>,
}

export interface GasSponsorshipRequest {
  // The address on the sponsored chain that will pay fees from the allowance.
  granteeAddress: string,
  // The dYdX address the deposit is for. Recorded with the sponsorship.
  dydxAddress: string,
}

export enum GasSponsorshipGrantStatus {
  // The grant transaction has been sent but is not yet known to be in a block. The allowance
  // cannot be used yet.
  PENDING = 'PENDING',
  // The allowance is on chain and can be used.
  CONFIRMED = 'CONFIRMED',
}

export interface GasSponsorshipGrant {
  chainId: string,
  // The address to name as the fee granter of the sponsored transaction.
  granter: string,
  grantee: string,
  feeDenom: string,
  // What the grantee may still spend on fees, in the smallest unit of the fee denom.
  spendLimit: string,
  expiresAt: string,
  allowedMessages: string[],
  status: GasSponsorshipGrantStatus,
  // Hash of the grant transaction, when known.
  transactionHash?: string,
  // The gas limit and gas price the allowance is sized for. A sponsored transaction whose fee
  // is at most gasLimit * gasPrice fits within the allowance.
  gasLimit: string,
  gasPrice: string,
}

export enum GasSponsorshipErrorCode {
  DISABLED = 'GAS_SPONSORSHIP_DISABLED',
  INVALID_ADDRESS = 'GAS_SPONSORSHIP_INVALID_ADDRESS',
  INSUFFICIENT_USDC_BALANCE = 'GAS_SPONSORSHIP_INSUFFICIENT_USDC_BALANCE',
  ALREADY_HAS_GAS = 'GAS_SPONSORSHIP_ALREADY_HAS_GAS',
  COOLDOWN = 'GAS_SPONSORSHIP_COOLDOWN',
  BUDGET_EXHAUSTED = 'GAS_SPONSORSHIP_BUDGET_EXHAUSTED',
  BASE_FEE_TOO_HIGH = 'GAS_SPONSORSHIP_BASE_FEE_TOO_HIGH',
  DEPOSITS_CLOSED = 'GAS_SPONSORSHIP_DEPOSITS_CLOSED',
  SPONSOR_UNDERFUNDED = 'GAS_SPONSORSHIP_SPONSOR_UNDERFUNDED',
  BUSY = 'GAS_SPONSORSHIP_BUSY',
  CHAIN_UNAVAILABLE = 'GAS_SPONSORSHIP_CHAIN_UNAVAILABLE',
  GRANT_FAILED = 'GAS_SPONSORSHIP_GRANT_FAILED',
}

const HTTP_STATUS_BY_ERROR_CODE: Record<GasSponsorshipErrorCode, number> = {
  [GasSponsorshipErrorCode.DISABLED]: 503,
  [GasSponsorshipErrorCode.INVALID_ADDRESS]: 400,
  [GasSponsorshipErrorCode.INSUFFICIENT_USDC_BALANCE]: 403,
  [GasSponsorshipErrorCode.ALREADY_HAS_GAS]: 403,
  [GasSponsorshipErrorCode.COOLDOWN]: 429,
  [GasSponsorshipErrorCode.BUDGET_EXHAUSTED]: 503,
  [GasSponsorshipErrorCode.BASE_FEE_TOO_HIGH]: 503,
  [GasSponsorshipErrorCode.DEPOSITS_CLOSED]: 503,
  [GasSponsorshipErrorCode.SPONSOR_UNDERFUNDED]: 503,
  [GasSponsorshipErrorCode.BUSY]: 503,
  [GasSponsorshipErrorCode.CHAIN_UNAVAILABLE]: 503,
  [GasSponsorshipErrorCode.GRANT_FAILED]: 502,
};

/**
 * A refusal to sponsor, or a failure to, that is reported to the caller with its code.
 */
export class GasSponsorshipError extends Error {
  public readonly code: GasSponsorshipErrorCode;

  public readonly httpStatus: number;

  constructor(code: GasSponsorshipErrorCode, message: string) {
    super(message);
    this.name = 'GasSponsorshipError';
    this.code = code;
    this.httpStatus = HTTP_STATUS_BY_ERROR_CODE[code];
  }
}
