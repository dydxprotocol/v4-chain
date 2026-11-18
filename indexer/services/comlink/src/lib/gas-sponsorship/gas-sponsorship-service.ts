import { randomUUID } from 'crypto';

import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  GasSponsorshipFromDatabase,
  GasSponsorshipStatus,
  GasSponsorshipTable,
} from '@dydxprotocol-indexer/postgres';
import { redis } from '@dydxprotocol-indexer/redis';
import { RedisClient } from 'redis';

import config from '../../config';
import { isValidAccountAddress } from './addresses';
import {
  FeeAllowance,
  getExpectedSequence,
  IncludedTx,
  InjectiveAccount,
  InjectiveClientError,
  TxResult,
} from './injective-client';
import {
  assembleSignedTx,
  buildGrantTx,
  getTxHash,
  timestampFromDate,
  UnsignedTx,
} from './injective-tx';
import {
  injectiveAddressFromPublicKey,
  SponsorSigner,
  verifyDigestSignature,
} from './sponsor-signer';
import {
  GasSponsorshipError,
  GasSponsorshipErrorCode,
  GasSponsorshipGrant,
  GasSponsorshipGrantStatus,
  GasSponsorshipRequest,
  GasSponsorshipSettings,
  SponsorshipChainClient,
} from './types';

const ONE_DAY_MS: number = 24 * 60 * 60 * 1000;
// An allowance with less than this long left is replaced instead of being handed back, so that
// the caller is not given one that lapses before it can be used.
const MIN_REMAINING_ALLOWANCE_MS: number = 5 * 60 * 1000;
// A grant transaction that is still not in a block this long after it was signed is taken to
// have been dropped.
const PENDING_GIVE_UP_MS: number = 2 * 60 * 1000;
const LOCK_RETRY_MS: number = 100;
// The sequence is corrected and the transaction re-signed at most this many times per grant.
const MAX_SEQUENCE_RETRIES: number = 1;
// How many recent sponsorships of one address are looked at when applying the cooldown.
const RECENT_SPONSORSHIPS_LIMIT: number = 10;

interface GrantTerms {
  gasPrice: bigint,
  // The fee of one sponsored transaction at `gasPrice`.
  feePerAttempt: bigint,
  spendLimit: bigint,
  // The fee the sponsor pays for the grant transaction.
  grantFee: bigint,
}

/**
 * The Redis key of the lock that serializes the transactions of one sponsor account.
 */
export function getSponsorLockKey(chainId: string, granter: string): string {
  return `gas-sponsorship/lock/${chainId}/${granter}`;
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * Gives addresses on Injective a fee allowance from a sponsor account, so that they can send
 * USDC to the dYdX chain without holding INJ.
 *
 * The allowance lets the address pay fees out of the sponsor's balance, for the allowed message
 * types only, up to a spend limit and until an expiry. The address then signs and sends its own
 * transaction, naming the sponsor as the fee granter.
 *
 * All of the sponsor's transactions are signed and sent while holding a Redis lock, because each
 * must carry the account's next sequence number.
 */
export class GasSponsorshipService {
  public readonly granter: string;

  private readonly settings: GasSponsorshipSettings;

  private readonly client: SponsorshipChainClient;

  private readonly signer: SponsorSigner;

  private readonly redisClient: RedisClient;

  private readonly isDepositOpen: () => Promise<boolean>;

  private readonly now: () => Date;

  /**
   * @param isDepositOpen Whether the dYdX chain is accepting the deposits being sponsored.
   * @param now The clock used for expiries, the cooldown and the budget window.
   */
  constructor(
    settings: GasSponsorshipSettings,
    client: SponsorshipChainClient,
    signer: SponsorSigner,
    redisClient: RedisClient,
    isDepositOpen: () => Promise<boolean> = () => Promise.resolve(true),
    now: () => Date = () => new Date(),
  ) {
    this.settings = settings;
    this.client = client;
    this.signer = signer;
    this.redisClient = redisClient;
    this.isDepositOpen = isDepositOpen;
    this.now = now;
    this.granter = injectiveAddressFromPublicKey(signer.publicKey, settings.addressPrefix);
  }

  /**
   * Returns a fee allowance for the address, granting one if it is eligible and does not have a
   * usable one already.
   *
   * Safe to call repeatedly: while the address has a usable allowance, or a grant for it is
   * still pending, that is returned and nothing new is granted.
   *
   * @throws GasSponsorshipError if the address is not eligible, sponsorship is unavailable, or
   * the grant failed.
   */
  public async sponsor(request: GasSponsorshipRequest): Promise<GasSponsorshipGrant> {
    try {
      return await this.sponsorOrThrow(request);
    } catch (error) {
      if (error instanceof InjectiveClientError) {
        // A node could not be read, so eligibility cannot be established. Nothing was sent.
        logger.error({
          at: 'GasSponsorshipService#sponsor',
          message: 'Gas sponsorship could not read from the chain',
          error,
          grantee: request.granteeAddress,
        });
        throw new GasSponsorshipError(
          GasSponsorshipErrorCode.CHAIN_UNAVAILABLE,
          'Gas sponsorship is temporarily unavailable',
        );
      }
      throw error;
    }
  }

  private async sponsorOrThrow(request: GasSponsorshipRequest): Promise<GasSponsorshipGrant> {
    if (!this.settings.enabled) {
      throw new GasSponsorshipError(
        GasSponsorshipErrorCode.DISABLED,
        'Gas sponsorship is not available',
      );
    }
    const grantee: string = request.granteeAddress;
    if (!isValidAccountAddress(grantee, this.settings.addressPrefix)) {
      throw new GasSponsorshipError(
        GasSponsorshipErrorCode.INVALID_ADDRESS,
        `Address must be a valid ${this.settings.addressPrefix} address`,
      );
    }
    if (grantee === this.granter) {
      throw new GasSponsorshipError(
        GasSponsorshipErrorCode.INVALID_ADDRESS,
        'Address must not be the sponsor',
      );
    }

    const [depositOpen, baseFee, allowance, usdcBalance, feeBalance]: [
      boolean, bigint, FeeAllowance | undefined, bigint, bigint,
    ] = await Promise.all([
      this.isDepositOpen(),
      this.client.getBaseFee(),
      this.client.getAllowance(this.granter, grantee),
      this.client.getBalance(grantee, this.settings.usdcDenom),
      this.client.getBalance(grantee, this.settings.feeDenom),
    ]);

    if (!depositOpen) {
      throw new GasSponsorshipError(
        GasSponsorshipErrorCode.DEPOSITS_CLOSED,
        'Deposits to the dYdX chain are not open',
      );
    }
    if (baseFee > this.settings.maxBaseFee) {
      throw new GasSponsorshipError(
        GasSponsorshipErrorCode.BASE_FEE_TOO_HIGH,
        'Network fees are too high to sponsor right now',
      );
    }
    const terms: GrantTerms = this.getTerms(baseFee);

    if (allowance !== undefined && this.isUsable(allowance, terms)) {
      return this.describeExistingAllowance(allowance, terms);
    }
    if (usdcBalance < this.settings.minUsdcBalance) {
      throw new GasSponsorshipError(
        GasSponsorshipErrorCode.INSUFFICIENT_USDC_BALANCE,
        'Address does not hold enough USDC to be sponsored',
      );
    }
    if (feeBalance >= terms.feePerAttempt) {
      throw new GasSponsorshipError(
        GasSponsorshipErrorCode.ALREADY_HAS_GAS,
        'Address can already pay for gas',
      );
    }

    return this.withLock(() => this.grant(request, terms));
  }

  /**
   * Grants an allowance. Must be called while holding the lock.
   */
  private async grant(
    request: GasSponsorshipRequest,
    terms: GrantTerms,
  ): Promise<GasSponsorshipGrant> {
    const grantee: string = request.granteeAddress;

    // Read the allowance again now that no other grant can be in flight: a concurrent request
    // for this address may have granted one while this request waited for the lock.
    const allowance: FeeAllowance | undefined = await this.client.getAllowance(
      this.granter,
      grantee,
    );
    if (allowance !== undefined && this.isUsable(allowance, terms)) {
      return this.describeExistingAllowance(allowance, terms);
    }

    const pending: GasSponsorshipFromDatabase | undefined = await this.applyCooldown(grantee);
    if (pending !== undefined) {
      return this.describeSponsorship(pending, terms);
    }

    const now: Date = this.now();
    const committed: bigint = BigInt(await GasSponsorshipTable.getTotalCommittedSince(
      this.settings.chainId,
      new Date(now.getTime() - ONE_DAY_MS).toISOString(),
    ));
    if (committed + terms.spendLimit + terms.grantFee > this.settings.dailyBudget) {
      logger.error({
        at: 'GasSponsorshipService#grant',
        message: 'Gas sponsorship daily budget is exhausted',
        committed: committed.toString(),
        dailyBudget: this.settings.dailyBudget.toString(),
      });
      throw new GasSponsorshipError(
        GasSponsorshipErrorCode.BUDGET_EXHAUSTED,
        'Gas sponsorship is temporarily unavailable',
      );
    }

    const [account, sponsorBalance]: [InjectiveAccount | undefined, bigint] = await Promise.all([
      this.client.getAccount(this.granter),
      this.client.getBalance(this.granter, this.settings.feeDenom),
    ]);
    stats.gauge(`${config.SERVICE_NAME}.gas_sponsorship.sponsor_balance`, Number(sponsorBalance));
    if (account === undefined || sponsorBalance < terms.spendLimit + terms.grantFee) {
      logger.error({
        at: 'GasSponsorshipService#grant',
        message: 'Gas sponsor account is missing or cannot cover a grant',
        granter: this.granter,
        accountExists: account !== undefined,
        sponsorBalance: sponsorBalance.toString(),
      });
      throw new GasSponsorshipError(
        GasSponsorshipErrorCode.SPONSOR_UNDERFUNDED,
        'Gas sponsorship is temporarily unavailable',
      );
    }

    const expiration: Date = new Date(now.getTime() + this.settings.grantExpirySeconds * 1000);
    let sequence: bigint = account.sequence;
    let sponsorship: GasSponsorshipFromDatabase | undefined;

    for (let attempt: number = 0; ; attempt++) {
      const unsignedTx: UnsignedTx = buildGrantTx({
        chainId: this.settings.chainId,
        granter: this.granter,
        grantee,
        granterPublicKey: this.signer.publicKey,
        accountNumber: account.accountNumber,
        sequence,
        spendLimit: { denom: this.settings.feeDenom, amount: terms.spendLimit },
        expiration: timestampFromDate(expiration),
        allowedMessages: this.settings.allowedMessages,
        fee: { denom: this.settings.feeDenom, amount: terms.grantFee },
        gasLimit: this.settings.grantGasLimit,
        // Granting on top of an allowance is an error, so one that is not usable is replaced.
        revokeExisting: allowance !== undefined,
      });
      let txBytes: Uint8Array;
      try {
        txBytes = await this.sign(unsignedTx);
      } catch (error) {
        if (sponsorship !== undefined) {
          // Its recorded transaction is the one the node rejected for its sequence.
          await this.markFailed(sponsorship, `signing failed: ${error.message}`, false);
        }
        throw error;
      }
      const txHash: string = getTxHash(txBytes);

      // Recorded before broadcasting, so that the allowance counts against the budget and the
      // cooldown even if this process stops before it learns the outcome.
      sponsorship = sponsorship === undefined
        ? await GasSponsorshipTable.create({
          chain_id: this.settings.chainId,
          granter_address: this.granter,
          grantee_address: grantee,
          dydx_address: request.dydxAddress,
          spend_limit: terms.spendLimit.toString(),
          fee_amount: terms.grantFee.toString(),
          status: GasSponsorshipStatus.PENDING,
          transaction_hash: txHash,
          expires_at: expiration.toISOString(),
          created_at: now.toISOString(),
        })
        : (await GasSponsorshipTable.update(sponsorship.id, { transaction_hash: txHash }))!;

      let result: TxResult;
      try {
        result = await this.client.broadcastTx(txBytes);
      } catch (error) {
        // The node may or may not have taken the transaction. The sponsorship stays pending and
        // is settled by looking the transaction up on a later request.
        logger.error({
          at: 'GasSponsorshipService#grant',
          message: 'Broadcasting the grant transaction failed with an unknown outcome',
          error,
          txHash,
          grantee,
        });
        return this.describeSponsorship(sponsorship, terms);
      }

      if (result.code === 0) {
        break;
      }
      const expectedSequence: bigint | undefined = getExpectedSequence(result);
      if (
        expectedSequence !== undefined &&
        expectedSequence !== sequence &&
        attempt < MAX_SEQUENCE_RETRIES
      ) {
        // An earlier grant transaction is still in the mempool, so the sequence read from the
        // chain was behind. Sign again with the sequence the node expects.
        sequence = expectedSequence;
        continue;
      }
      await this.markFailed(sponsorship, `rejected on broadcast: ${result.rawLog}`, false);
      throw new GasSponsorshipError(
        GasSponsorshipErrorCode.GRANT_FAILED,
        'The grant transaction was rejected',
      );
    }

    stats.increment(`${config.SERVICE_NAME}.gas_sponsorship.grant_broadcast`);
    return this.describeSponsorship(await this.awaitInclusion(sponsorship), terms);
  }

  private async sign(unsignedTx: UnsignedTx): Promise<Uint8Array> {
    const signature: Uint8Array = await this.signer.signDigest(unsignedTx.signDigest);
    // Catches a signer that holds a different key from the configured one before anything is
    // sent, where the chain would only report a failed signature check.
    if (!await verifyDigestSignature(signature, unsignedTx.signDigest, this.signer.publicKey)) {
      throw new Error('The sponsor signature does not match the configured sponsor public key');
    }
    return assembleSignedTx(unsignedTx, signature);
  }

  /**
   * Applies the per-address cooldown. Returns a sponsorship of the address that is still
   * pending, if there is one, in which case that is the answer to the request and nothing new
   * is granted.
   *
   * @throws GasSponsorshipError if the address was given an allowance within the cooldown.
   */
  private async applyCooldown(
    grantee: string,
  ): Promise<GasSponsorshipFromDatabase | undefined> {
    const recent: GasSponsorshipFromDatabase[] = await GasSponsorshipTable.findAll(
      {
        chain_id: this.settings.chainId,
        grantee_address: grantee,
        statuses: [GasSponsorshipStatus.PENDING, GasSponsorshipStatus.CONFIRMED],
        createdOnOrAfter: new Date(
          this.now().getTime() - this.settings.cooldownSeconds * 1000,
        ).toISOString(),
      },
      RECENT_SPONSORSHIPS_LIMIT,
    );

    let confirmed: boolean = false;
    for (const sponsorship of recent) {
      const settled: GasSponsorshipFromDatabase = await this.settle(sponsorship);
      if (settled.status === GasSponsorshipStatus.PENDING) {
        return settled;
      }
      confirmed = confirmed || settled.status === GasSponsorshipStatus.CONFIRMED;
    }
    if (confirmed) {
      throw new GasSponsorshipError(
        GasSponsorshipErrorCode.COOLDOWN,
        'Address was sponsored recently',
      );
    }
    return undefined;
  }

  /**
   * Brings a pending sponsorship up to date with the chain, returning it as it now stands.
   */
  private async settle(
    sponsorship: GasSponsorshipFromDatabase,
  ): Promise<GasSponsorshipFromDatabase> {
    if (sponsorship.status !== GasSponsorshipStatus.PENDING) {
      return sponsorship;
    }
    const tx: IncludedTx | undefined = sponsorship.transaction_hash
      ? await this.client.getTx(sponsorship.transaction_hash)
      : undefined;
    if (tx !== undefined) {
      return tx.code === 0
        ? this.markConfirmed(sponsorship)
        : this.markFailed(sponsorship, `failed in block ${tx.height}: ${tx.rawLog}`);
    }
    const age: number = this.now().getTime() - new Date(sponsorship.created_at).getTime();
    if (age > PENDING_GIVE_UP_MS) {
      return this.markFailed(sponsorship, 'not included in a block in time');
    }
    return sponsorship;
  }

  /**
   * Waits for the grant transaction to be included in a block. The allowance cannot be used
   * before then. Returns the sponsorship as confirmed, or still pending if the wait timed out.
   *
   * @throws GasSponsorshipError if the transaction was included but failed.
   */
  private async awaitInclusion(
    sponsorship: GasSponsorshipFromDatabase,
  ): Promise<GasSponsorshipFromDatabase> {
    const deadline: number = Date.now() + this.settings.confirmationTimeoutMs;
    for (;;) {
      let tx: IncludedTx | undefined;
      try {
        tx = await this.client.getTx(sponsorship.transaction_hash!);
      } catch (error) {
        // Keep polling: one failed lookup says nothing about the transaction.
        logger.info({
          at: 'GasSponsorshipService#awaitInclusion',
          message: 'Looking up the grant transaction failed',
          error,
          txHash: sponsorship.transaction_hash,
        });
      }
      if (tx !== undefined) {
        if (tx.code === 0) {
          return this.markConfirmed(sponsorship);
        }
        await this.markFailed(sponsorship, `failed in block ${tx.height}: ${tx.rawLog}`);
        throw new GasSponsorshipError(
          GasSponsorshipErrorCode.GRANT_FAILED,
          'The grant transaction failed',
        );
      }
      if (Date.now() >= deadline) {
        stats.increment(`${config.SERVICE_NAME}.gas_sponsorship.grant_unconfirmed`);
        return sponsorship;
      }
      await sleep(this.settings.confirmationPollMs);
    }
  }

  private async markConfirmed(
    sponsorship: GasSponsorshipFromDatabase,
  ): Promise<GasSponsorshipFromDatabase> {
    stats.increment(`${config.SERVICE_NAME}.gas_sponsorship.grant_confirmed`);
    return (await GasSponsorshipTable.update(
      sponsorship.id,
      { status: GasSponsorshipStatus.CONFIRMED },
    ))!;
  }

  /**
   * Records that a grant gave no allowance. `transactionExists` is false when no node accepted
   * the transaction: its hash is then dropped from the record, since nothing on the chain has it
   * and a later grant signed over the same content would have the same one.
   */
  private async markFailed(
    sponsorship: GasSponsorshipFromDatabase,
    reason: string,
    transactionExists: boolean = true,
  ): Promise<GasSponsorshipFromDatabase> {
    stats.increment(`${config.SERVICE_NAME}.gas_sponsorship.grant_failed`);
    logger.error({
      at: 'GasSponsorshipService#markFailed',
      message: 'Gas sponsorship grant failed',
      reason,
      sponsorshipId: sponsorship.id,
      txHash: sponsorship.transaction_hash,
      grantee: sponsorship.grantee_address,
    });
    return (await GasSponsorshipTable.update(
      sponsorship.id,
      {
        status: GasSponsorshipStatus.FAILED,
        ...(transactionExists ? {} : { transaction_hash: null }),
      },
    ))!;
  }

  private getTerms(baseFee: bigint): GrantTerms {
    // Rounded up, so that the gas price is never below the buffered base fee.
    const gasPrice: bigint = (
      baseFee * this.settings.gasPriceBufferPercent + BigInt(99)
    ) / BigInt(100);
    const feePerAttempt: bigint = this.settings.transferGasLimit * gasPrice;
    return {
      gasPrice,
      feePerAttempt,
      spendLimit: feePerAttempt * this.settings.transferAttempts,
      grantFee: this.settings.grantGasLimit * gasPrice,
    };
  }

  /**
   * Whether an allowance from the sponsor can pay for at least one more sponsored transaction.
   * Only allowances of the shape this service grants count; anything else is replaced.
   */
  private isUsable(allowance: FeeAllowance, terms: GrantTerms): boolean {
    if (
      !allowance.recognized ||
      allowance.spendLimit === undefined ||
      allowance.expiration === undefined ||
      allowance.allowedMessages === undefined
    ) {
      return false;
    }
    const remaining: bigint = allowance.spendLimit[this.settings.feeDenom] ?? BigInt(0);
    const allowedMessages: string[] = allowance.allowedMessages;
    return remaining >= terms.feePerAttempt &&
      allowance.expiration.getTime() - this.now().getTime() >= MIN_REMAINING_ALLOWANCE_MS &&
      this.settings.allowedMessages.every(
        (message: string) => allowedMessages.includes(message),
      );
  }

  /**
   * Describes a usable allowance that is already on chain.
   */
  private async describeExistingAllowance(
    allowance: FeeAllowance,
    terms: GrantTerms,
  ): Promise<GasSponsorshipGrant> {
    // The record of the grant that created this allowance, if there is one, is the latest for
    // the address and has the same expiry. It supplies the transaction hash, and if it is still
    // marked pending this is the evidence that it went through.
    const latest: GasSponsorshipFromDatabase | undefined = await GasSponsorshipTable.findLatest({
      chain_id: this.settings.chainId,
      grantee_address: allowance.grantee,
      statuses: [GasSponsorshipStatus.PENDING, GasSponsorshipStatus.CONFIRMED],
    });
    const isSameGrant: boolean = latest !== undefined &&
      new Date(latest.expires_at).getTime() === allowance.expiration!.getTime();
    if (isSameGrant && latest!.status === GasSponsorshipStatus.PENDING) {
      await this.markConfirmed(latest!);
    }

    return {
      chainId: this.settings.chainId,
      granter: this.granter,
      grantee: allowance.grantee,
      feeDenom: this.settings.feeDenom,
      spendLimit: allowance.spendLimit![this.settings.feeDenom].toString(),
      expiresAt: allowance.expiration!.toISOString(),
      allowedMessages: allowance.allowedMessages!,
      status: GasSponsorshipGrantStatus.CONFIRMED,
      transactionHash: (isSameGrant && latest!.transaction_hash) || undefined,
      gasLimit: this.settings.transferGasLimit.toString(),
      gasPrice: terms.gasPrice.toString(),
    };
  }

  private describeSponsorship(
    sponsorship: GasSponsorshipFromDatabase,
    terms: GrantTerms,
  ): GasSponsorshipGrant {
    return {
      chainId: sponsorship.chain_id,
      granter: sponsorship.granter_address,
      grantee: sponsorship.grantee_address,
      feeDenom: this.settings.feeDenom,
      spendLimit: sponsorship.spend_limit,
      expiresAt: new Date(sponsorship.expires_at).toISOString(),
      allowedMessages: this.settings.allowedMessages,
      status: sponsorship.status === GasSponsorshipStatus.CONFIRMED
        ? GasSponsorshipGrantStatus.CONFIRMED
        : GasSponsorshipGrantStatus.PENDING,
      transactionHash: sponsorship.transaction_hash ?? undefined,
      gasLimit: this.settings.transferGasLimit.toString(),
      gasPrice: terms.gasPrice.toString(),
    };
  }

  private async withLock<T>(fn: () => Promise<T>): Promise<T> {
    const lockKey: string = getSponsorLockKey(this.settings.chainId, this.granter);
    const lockValue: string = randomUUID();
    const deadline: number = Date.now() + this.settings.lockWaitMs;

    while (!await redis.lockWithExpiry(
      this.redisClient,
      lockKey,
      lockValue,
      this.settings.lockExpiryMs,
    )) {
      if (Date.now() >= deadline) {
        stats.increment(`${config.SERVICE_NAME}.gas_sponsorship.lock_timeout`);
        throw new GasSponsorshipError(
          GasSponsorshipErrorCode.BUSY,
          'Gas sponsorship is busy, try again shortly',
        );
      }
      await sleep(LOCK_RETRY_MS);
    }

    try {
      return await fn();
    } finally {
      try {
        // Only releases the lock if it is still this request's: had it lapsed, another request
        // may hold it by now.
        await redis.unlock(this.redisClient, lockKey, lockValue);
      } catch (error) {
        logger.error({
          at: 'GasSponsorshipService#withLock',
          message: 'Failed to release the gas sponsorship lock',
          error,
        });
      }
    }
  }
}
