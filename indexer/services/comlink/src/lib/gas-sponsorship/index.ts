import { logger } from '@dydxprotocol-indexer/base';
import { Turnkey } from '@turnkey/sdk-server';

import config from '../../config';
import { redisClient } from '../../helpers/redis/redis-controller';
import { INJECTIVE_ADDRESS_PREFIX } from './addresses';
import { GasSponsorshipService } from './gas-sponsorship-service';
import { InjectiveClient, isInjectiveDepositOpen } from './injective-client';
import { parseCompressedPublicKey, TurnkeySponsorSigner } from './sponsor-signer';
import {
  GasSponsorshipError,
  GasSponsorshipErrorCode,
  GasSponsorshipSettings,
} from './types';

export * from './types';
export { GasSponsorshipService } from './gas-sponsorship-service';

export function getGasSponsorshipSettings(): GasSponsorshipSettings {
  return {
    enabled: config.GAS_SPONSORSHIP_ENABLED,
    chainId: config.GAS_SPONSORSHIP_INJECTIVE_CHAIN_ID,
    addressPrefix: INJECTIVE_ADDRESS_PREFIX,
    feeDenom: config.GAS_SPONSORSHIP_FEE_DENOM,
    usdcDenom: config.GAS_SPONSORSHIP_USDC_DENOM,
    allowedMessages: config.GAS_SPONSORSHIP_ALLOWED_MESSAGES
      .split(',')
      .map((message: string) => message.trim())
      .filter((message: string) => message.length > 0),
    transferGasLimit: BigInt(config.GAS_SPONSORSHIP_TRANSFER_GAS_LIMIT),
    transferAttempts: BigInt(config.GAS_SPONSORSHIP_TRANSFER_ATTEMPTS),
    gasPriceBufferPercent: BigInt(config.GAS_SPONSORSHIP_GAS_PRICE_BUFFER_PERCENT),
    maxBaseFee: config.GAS_SPONSORSHIP_MAX_BASE_FEE,
    grantGasLimit: BigInt(config.GAS_SPONSORSHIP_GRANT_GAS_LIMIT),
    grantExpirySeconds: config.GAS_SPONSORSHIP_GRANT_EXPIRY_SECONDS,
    minUsdcBalance: config.GAS_SPONSORSHIP_MIN_USDC_BALANCE,
    cooldownSeconds: config.GAS_SPONSORSHIP_COOLDOWN_SECONDS,
    dailyBudget: config.GAS_SPONSORSHIP_DAILY_BUDGET,
    lockExpiryMs: config.GAS_SPONSORSHIP_LOCK_EXPIRY_MS,
    lockWaitMs: config.GAS_SPONSORSHIP_LOCK_WAIT_MS,
    confirmationTimeoutMs: config.GAS_SPONSORSHIP_CONFIRMATION_TIMEOUT_MS,
    confirmationPollMs: config.GAS_SPONSORSHIP_CONFIRMATION_POLL_MS,
  };
}

/**
 * Builds the service from the config, with the sponsor key held in the Turnkey parent
 * organization.
 *
 * @throws GasSponsorshipError if sponsorship is switched off, or switched on without the
 * settings it needs.
 */
export function createGasSponsorshipService(): GasSponsorshipService {
  const settings: GasSponsorshipSettings = getGasSponsorshipSettings();
  if (!settings.enabled) {
    throw new GasSponsorshipError(
      GasSponsorshipErrorCode.DISABLED,
      'Gas sponsorship is not available',
    );
  }

  const missing: string[] = [
    ['GAS_SPONSORSHIP_INJECTIVE_REST_URL', config.GAS_SPONSORSHIP_INJECTIVE_REST_URL],
    ['GAS_SPONSORSHIP_SPONSOR_PUBLIC_KEY', config.GAS_SPONSORSHIP_SPONSOR_PUBLIC_KEY],
    ['TURNKEY_API_PRIVATE_KEY', config.TURNKEY_API_PRIVATE_KEY],
    ['TURNKEY_API_PUBLIC_KEY', config.TURNKEY_API_PUBLIC_KEY],
    ['TURNKEY_ORGANIZATION_ID', config.TURNKEY_ORGANIZATION_ID],
    ['GAS_SPONSORSHIP_ALLOWED_MESSAGES', settings.allowedMessages.join(',')],
  ].filter(([, value]: string[]) => value === '').map(([name]: string[]) => name);
  if (missing.length > 0) {
    logger.error({
      at: 'gas-sponsorship#createGasSponsorshipService',
      message: 'Gas sponsorship is enabled but not fully configured',
      missing,
    });
    throw new GasSponsorshipError(
      GasSponsorshipErrorCode.DISABLED,
      'Gas sponsorship is not available',
    );
  }

  const signer: TurnkeySponsorSigner = new TurnkeySponsorSigner(
    new Turnkey({
      apiBaseUrl: config.TURNKEY_API_BASE_URL,
      apiPublicKey: config.TURNKEY_API_PUBLIC_KEY,
      apiPrivateKey: config.TURNKEY_API_PRIVATE_KEY,
      defaultOrganizationId: config.TURNKEY_ORGANIZATION_ID,
    }).apiClient(),
    parseCompressedPublicKey(config.GAS_SPONSORSHIP_SPONSOR_PUBLIC_KEY),
    config.GAS_SPONSORSHIP_TURNKEY_SIGN_WITH,
  );
  const dydxRestUrl: string = config.GAS_SPONSORSHIP_DYDX_REST_URL;

  return new GasSponsorshipService(
    settings,
    new InjectiveClient(
      config.GAS_SPONSORSHIP_INJECTIVE_REST_URL,
      config.GAS_SPONSORSHIP_REQUEST_TIMEOUT_MS,
    ),
    signer,
    redisClient,
    dydxRestUrl === ''
      ? undefined
      : () => isInjectiveDepositOpen(dydxRestUrl, config.GAS_SPONSORSHIP_REQUEST_TIMEOUT_MS),
  );
}

let gasSponsorshipService: GasSponsorshipService | undefined;

/**
 * Returns the service shared by everything in this process that sponsors gas, building it on
 * first use.
 */
export function getGasSponsorshipService(): GasSponsorshipService {
  if (gasSponsorshipService === undefined) {
    gasSponsorshipService = createGasSponsorshipService();
  }
  return gasSponsorshipService;
}
