import { logger } from '@dydxprotocol-indexer/base';

import config from '../../../src/config';
import {
  createGasSponsorshipService,
  GasSponsorshipError,
  GasSponsorshipErrorCode,
  GasSponsorshipService,
  getGasSponsorshipSettings,
} from '../../../src/lib/gas-sponsorship';

// The well-known key of the first account of the Hardhat and Anvil development chains.
const PUBLIC_KEY: string = '038318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed75';
const INJECTIVE_ADDRESS_OF_PUBLIC_KEY: string = 'inj17w0adeg64ky0daxwd2ugyuneellmjgnxf5vkec';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const mutableConfig: any = config;

function expectDisabled(fn: () => unknown): void {
  let thrown: unknown;
  try {
    fn();
  } catch (error) {
    thrown = error;
  }
  expect(thrown).toBeInstanceOf(GasSponsorshipError);
  expect((thrown as GasSponsorshipError).code).toEqual(GasSponsorshipErrorCode.DISABLED);
}

describe('gas-sponsorship', () => {
  const original: { [name: string]: unknown } = {};
  const settingsUnderTest: string[] = [
    'GAS_SPONSORSHIP_ENABLED',
    'GAS_SPONSORSHIP_INJECTIVE_REST_URL',
    'GAS_SPONSORSHIP_SPONSOR_PUBLIC_KEY',
    'GAS_SPONSORSHIP_ALLOWED_MESSAGES',
    'TURNKEY_API_PRIVATE_KEY',
    'TURNKEY_API_PUBLIC_KEY',
    'TURNKEY_ORGANIZATION_ID',
  ];

  const configure = (): void => {
    mutableConfig.GAS_SPONSORSHIP_ENABLED = true;
    mutableConfig.GAS_SPONSORSHIP_INJECTIVE_REST_URL = 'https://injective.example';
    mutableConfig.GAS_SPONSORSHIP_SPONSOR_PUBLIC_KEY = PUBLIC_KEY;
    mutableConfig.TURNKEY_API_PRIVATE_KEY = 'test-private-key';
    mutableConfig.TURNKEY_API_PUBLIC_KEY = 'test-public-key';
    mutableConfig.TURNKEY_ORGANIZATION_ID = 'test-organization';
  };

  beforeEach(() => {
    settingsUnderTest.forEach((name: string) => {
      original[name] = mutableConfig[name];
    });
  });

  afterEach(() => {
    settingsUnderTest.forEach((name: string) => {
      mutableConfig[name] = original[name];
    });
  });

  describe('getGasSponsorshipSettings', () => {
    it('is switched off, and bounded, with no configuration at all', () => {
      expect(getGasSponsorshipSettings()).toEqual({
        enabled: false,
        chainId: 'injective-1',
        addressPrefix: 'inj',
        feeDenom: 'inj',
        usdcDenom: 'erc20:0xa00C59fF5a080D2b954d0c75e46E22a0c371235a',
        allowedMessages: ['/ibc.applications.transfer.v1.MsgTransfer'],
        transferGasLimit: BigInt(450_000),
        transferAttempts: BigInt(3),
        gasPriceBufferPercent: BigInt(150),
        maxBaseFee: BigInt(1_600_000_000),
        grantGasLimit: BigInt(250_000),
        grantExpirySeconds: 86_400,
        minUsdcBalance: BigInt(1_000_000),
        cooldownSeconds: 3_600,
        dailyBudget: BigInt('1000000000000000000'),
        lockExpiryMs: 60_000,
        lockWaitMs: 10_000,
        confirmationTimeoutMs: 15_000,
        confirmationPollMs: 500,
      });
    });

    it('holds the lock for longer than a grant takes to talk to the chain by default', () => {
      // Under the lock, a grant makes at most these requests to a node one after another: the
      // allowance, a lookup of the one sponsorship of the address that can be pending, the
      // account and balance together, and two broadcasts. It then waits for confirmation, and
      // the last lookup of that wait may start just before the timeout.
      const requestsBeforeConfirmation: number = 5;
      const worstCaseMs: number = (
        requestsBeforeConfirmation * config.GAS_SPONSORSHIP_REQUEST_TIMEOUT_MS +
        config.GAS_SPONSORSHIP_CONFIRMATION_TIMEOUT_MS +
        config.GAS_SPONSORSHIP_REQUEST_TIMEOUT_MS
      );

      expect(worstCaseMs).toEqual(45_000);
      expect(config.GAS_SPONSORSHIP_LOCK_EXPIRY_MS).toBeGreaterThan(worstCaseMs);
    });

    it('reads a comma-separated list of allowed messages', () => {
      mutableConfig.GAS_SPONSORSHIP_ALLOWED_MESSAGES = (
        '/ibc.applications.transfer.v1.MsgTransfer, /cosmos.bank.v1beta1.MsgSend,'
      );

      expect(getGasSponsorshipSettings().allowedMessages).toEqual([
        '/ibc.applications.transfer.v1.MsgTransfer',
        '/cosmos.bank.v1beta1.MsgSend',
      ]);
    });
  });

  describe('createGasSponsorshipService', () => {
    it('refuses while sponsorship is switched off, even if fully configured', () => {
      configure();
      mutableConfig.GAS_SPONSORSHIP_ENABLED = false;

      expectDisabled(createGasSponsorshipService);
    });

    it.each([
      'GAS_SPONSORSHIP_INJECTIVE_REST_URL',
      'GAS_SPONSORSHIP_SPONSOR_PUBLIC_KEY',
      'GAS_SPONSORSHIP_ALLOWED_MESSAGES',
      'TURNKEY_API_PRIVATE_KEY',
      'TURNKEY_API_PUBLIC_KEY',
      'TURNKEY_ORGANIZATION_ID',
    ])('refuses, and logs what is missing, when switched on without %s', (name: string) => {
      const logError: jest.SpyInstance = jest.spyOn(logger, 'error').mockImplementation();
      configure();
      mutableConfig[name] = '';

      expectDisabled(createGasSponsorshipService);
      expect(logError).toHaveBeenCalledWith(expect.objectContaining({
        message: 'Gas sponsorship is enabled but not fully configured',
        missing: [name],
      }));
    });

    it('fails on a sponsor public key that is not a key', () => {
      configure();
      mutableConfig.GAS_SPONSORSHIP_SPONSOR_PUBLIC_KEY = 'not-a-key';

      expect(createGasSponsorshipService).toThrow('sponsor public key is not valid hex');
    });

    it('builds a service that sponsors from the address of the configured public key', () => {
      configure();

      const service: GasSponsorshipService = createGasSponsorshipService();

      expect(service.granter).toEqual(INJECTIVE_ADDRESS_OF_PUBLIC_KEY);
    });
  });
});
