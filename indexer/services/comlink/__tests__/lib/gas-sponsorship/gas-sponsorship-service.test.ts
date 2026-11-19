import { fromBase64, toBech32 } from '@cosmjs/encoding';
import {
  dbHelpers,
  GasSponsorshipFromDatabase,
  GasSponsorshipStatus,
  GasSponsorshipTable,
} from '@dydxprotocol-indexer/postgres';
import { redis } from '@dydxprotocol-indexer/redis';

import { redisClient } from '../../../src/helpers/redis/redis-controller';
import {
  GasSponsorshipService,
  getSponsorLockKey,
} from '../../../src/lib/gas-sponsorship/gas-sponsorship-service';
import { InjectiveClientError } from '../../../src/lib/gas-sponsorship/injective-client';
import { SponsorSigner } from '../../../src/lib/gas-sponsorship/sponsor-signer';
import {
  GasSponsorshipError,
  GasSponsorshipErrorCode,
  GasSponsorshipGrant,
  GasSponsorshipGrantStatus,
  GasSponsorshipSettings,
} from '../../../src/lib/gas-sponsorship/types';
import {
  DecodedTx,
  decodeTx,
  FakeInjectiveChain,
  LocalKeySigner,
} from '../../helpers/gas-sponsorship-helpers';

const USDC_DENOM: string = 'erc20:0xa00C59fF5a080D2b954d0c75e46E22a0c371235a';
const MSG_TRANSFER: string = '/ibc.applications.transfer.v1.MsgTransfer';
const DYDX_ADDRESS: string = 'dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4';

// 150% of the base fee of the fake chain, which is 160,000,000.
const GAS_PRICE: bigint = BigInt(240_000_000);
// 450,000 gas at the gas price.
const FEE_PER_ATTEMPT: bigint = BigInt('108000000000000');
// Three attempts.
const SPEND_LIMIT: bigint = BigInt('324000000000000');
// 250,000 gas at the gas price.
const GRANT_FEE: bigint = BigInt('60000000000000');
const ONE_INJ: bigint = BigInt('1000000000000000000');
const INITIAL_SEQUENCE: bigint = BigInt(41);

const START_TIME: string = '2026-11-18T12:00:00.000Z';
const EXPIRY_OF_GRANT_AT_START: string = '2026-11-19T12:00:00.000Z';
const ONE_MINUTE_MS: number = 60 * 1000;
const ONE_HOUR_MS: number = 60 * ONE_MINUTE_MS;

const DEFAULT_SETTINGS: GasSponsorshipSettings = {
  enabled: true,
  chainId: 'injective-1',
  addressPrefix: 'inj',
  feeDenom: 'inj',
  usdcDenom: USDC_DENOM,
  allowedMessages: [MSG_TRANSFER],
  transferGasLimit: BigInt(450_000),
  transferAttempts: BigInt(3),
  gasPriceBufferPercent: BigInt(150),
  maxBaseFee: BigInt(1_600_000_000),
  grantGasLimit: BigInt(250_000),
  grantExpirySeconds: 86_400,
  minUsdcBalance: BigInt(1_000_000),
  cooldownSeconds: 3_600,
  dailyBudget: ONE_INJ,
  lockExpiryMs: 10_000,
  lockWaitMs: 5_000,
  confirmationTimeoutMs: 60,
  confirmationPollMs: 5,
};

function addressOf(byte: number): string {
  return toBech32('inj', new Uint8Array(20).fill(byte));
}

const GRANTEE: string = addressOf(1);
const OTHER_GRANTEE: string = addressOf(2);

describe('GasSponsorshipService', () => {
  let chain: FakeInjectiveChain;
  let signer: LocalKeySigner;
  let granter: string;
  let currentTime: Date;
  let depositOpen: boolean;

  const makeService = (
    overrides: Partial<GasSponsorshipSettings> = {},
    sponsorSigner: SponsorSigner = signer,
  ): GasSponsorshipService => new GasSponsorshipService(
    { ...DEFAULT_SETTINGS, ...overrides },
    chain,
    sponsorSigner,
    redisClient,
    () => Promise.resolve(depositOpen),
    () => new Date(currentTime),
  );

  const advanceTime = (ms: number): void => {
    currentTime = new Date(currentTime.getTime() + ms);
  };

  const fundGrantee = (address: string): void => {
    chain.setBalance(address, USDC_DENOM, BigInt(25_000_000));
  };

  const sponsor = (
    service: GasSponsorshipService,
    granteeAddress: string = GRANTEE,
  ): Promise<GasSponsorshipGrant> => service.sponsor({ granteeAddress, dydxAddress: DYDX_ADDRESS });

  const expectRefusal = async (
    promise: Promise<GasSponsorshipGrant>,
    code: GasSponsorshipErrorCode,
  ): Promise<void> => {
    const error: GasSponsorshipError = await promise.then(
      () => { throw new Error(`expected a refusal with code ${code}, but it was sponsored`); },
      (e: GasSponsorshipError) => e,
    );
    expect(error).toBeInstanceOf(GasSponsorshipError);
    expect(error.code).toEqual(code);
  };

  const allSponsorships = (): Promise<GasSponsorshipFromDatabase[]> => {
    return GasSponsorshipTable.findAll();
  };

  beforeAll(async () => {
    await dbHelpers.migrate();
    signer = await LocalKeySigner.create('11'.repeat(32));
  });

  beforeEach(async () => {
    await redis.deleteAllAsync(redisClient);
    chain = new FakeInjectiveChain();
    signer.signCount = 0;
    granter = chain.addAccount(signer.publicKey, INITIAL_SEQUENCE);
    chain.setBalance(granter, 'inj', ONE_INJ);
    fundGrantee(GRANTEE);
    fundGrantee(OTHER_GRANTEE);
    currentTime = new Date(START_TIME);
    depositOpen = true;
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    await redis.deleteAllAsync(redisClient);
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('the fake chain used by these tests', () => {
    it('reads a fee-grant transaction from Injective mainnet correctly', () => {
      // Transaction C554B3ACB4F98161CCFE7938399A73C4D452FBD9EF92F88C4F79E68896905BE4.
      const tx: DecodedTx = decodeTx(fromBase64(
        'CuoDCocBCisvY29zbW9zLmZlZWdyYW50LnYxYmV0YTEuTXNnUmV2b2tlQWxsb3dhbmNlElgKKmluajF4ZHJr' +
        'eGdrZThxemV5OXplNGN4NDhoZWF2Nmx6a2tzenVjMjI0chIqaW5qMWZyd3NxeW15ZzN3bmZoeTRtcXJxOWNn' +
        'aHd3enN0MHRnc21mYWw3Ct0CCiovY29zbW9zLmZlZWdyYW50LnYxYmV0YTEuTXNnR3JhbnRBbGxvd2FuY2US' +
        'rgIKKmluajF4ZHJreGdrZThxemV5OXplNGN4NDhoZWF2Nmx6a2tzenVjMjI0chIqaW5qMWZyd3NxeW15ZzN3' +
        'bmZoeTRtcXJxOWNnaHd3enN0MHRnc21mYWw3GtMBCiwvY29zbW9zLmZlZWdyYW50LnYxYmV0YTEuQWxsb3dl' +
        'ZE1zZ0FsbG93YW5jZRKiAQpTCicvY29zbW9zLmZlZWdyYW50LnYxYmV0YTEuQmFzaWNBbGxvd2FuY2USKAoY' +
        'CgNpbmoSETM1MjAwMDAwMDAwMDAwMDAwEgwIod/81wYQu/jztgISHS9jb3Ntb3MuYXV0aHoudjFiZXRhMS5N' +
        'c2dFeGVjEiwvaW5qZWN0aXZlLndhc214LnYxLk1zZ0V4ZWN1dGVDb250cmFjdENvbXBhdBJ/CmAKVAotL2lu' +
        'amVjdGl2ZS5jcnlwdG8udjFiZXRhMS5ldGhzZWNwMjU2azEuUHViS2V5EiMKIQKvJ53oWgb9UF1LQcvlSqmF' +
        'OZtqy98KgFqgueuoCjKuKxIECgIIARiCsgMSGwoVCgNpbmoSDjk2MDAwMDAwMDAwMDAwEMDPJBpBilgYbrEY' +
        'eAk2LvFoK7l5BrZkgAHItMeM3oQCCkWGpZMHz7ExpX2GcBYZyMb1KYvSRlpL/xw4KHAX0EhlLOCMOgA=',
      ));
      const mainnetGranter: string = 'inj1xdrkxgke8qzey9ze4cx48heav6lzkkszuc224r';
      const mainnetGrantee: string = 'inj1frwsqymyg3wnfhy4mqrq9cghwwzst0tgsmfal7';

      expect(tx.txHash).toEqual('C554B3ACB4F98161CCFE7938399A73C4D452FBD9EF92F88C4F79E68896905BE4');
      expect(tx.messages).toEqual([
        { type: 'revoke', granter: mainnetGranter, grantee: mainnetGrantee },
        {
          type: 'grant',
          granter: mainnetGranter,
          grantee: mainnetGrantee,
          spendLimit: { inj: BigInt('35200000000000000') },
          // The chain has 2026-11-19T17:19:29.652016699Z; a Date keeps the milliseconds.
          expiration: new Date('2026-11-19T17:19:29.652Z'),
          allowedMessages: [
            '/cosmos.authz.v1beta1.MsgExec',
            '/injective.wasmx.v1.MsgExecuteContractCompat',
          ],
        },
      ]);
      expect(tx.sequence).toEqual(BigInt(55554));
      expect(tx.fee).toEqual({ inj: BigInt('96000000000000') });
      expect(tx.gasLimit).toEqual(BigInt(600000));
      expect(tx.signerPublicKey).toEqual(
        Uint8Array.from(fromBase64('Aq8nnehaBv1QXUtBy+VKqYU5m2rL3wqAWqC566gKMq4r')),
      );
    });
  });

  describe('granting', () => {
    it('grants an allowance sized from the base fee, and reports it once it is on chain', async () => {
      const grant: GasSponsorshipGrant = await sponsor(makeService());

      expect(chain.submitted).toHaveLength(1);
      const tx: DecodedTx = chain.submitted[0];
      expect(grant).toEqual({
        chainId: 'injective-1',
        granter,
        grantee: GRANTEE,
        feeDenom: 'inj',
        spendLimit: SPEND_LIMIT.toString(),
        expiresAt: EXPIRY_OF_GRANT_AT_START,
        allowedMessages: [MSG_TRANSFER],
        status: GasSponsorshipGrantStatus.CONFIRMED,
        transactionHash: tx.txHash,
        gasLimit: '450000',
        gasPrice: GAS_PRICE.toString(),
      });
    });

    it('sends a transaction that grants exactly the allowance, and nothing else', async () => {
      await sponsor(makeService());

      const tx: DecodedTx = chain.submitted[0];
      expect(tx.messages).toEqual([{
        type: 'grant',
        granter,
        grantee: GRANTEE,
        spendLimit: { inj: SPEND_LIMIT },
        expiration: new Date(EXPIRY_OF_GRANT_AT_START),
        allowedMessages: [MSG_TRANSFER],
      }]);
      expect(tx.sequence).toEqual(INITIAL_SEQUENCE);
      expect(tx.fee).toEqual({ inj: GRANT_FEE });
      expect(tx.gasLimit).toEqual(BigInt(250_000));
      expect(tx.memo).toEqual('');
      expect(tx.feePayer).toEqual('');
      expect(tx.feeGranter).toEqual('');
      expect(tx.signerPublicKey).toEqual(signer.publicKey);
    });

    it('leaves the allowance on chain, having paid the grant fee from the sponsor', async () => {
      await sponsor(makeService());

      expect(await chain.getAllowance(granter, GRANTEE)).toEqual({
        granter,
        grantee: GRANTEE,
        recognized: true,
        spendLimit: { inj: SPEND_LIMIT },
        expiration: new Date(EXPIRY_OF_GRANT_AT_START),
        allowedMessages: [MSG_TRANSFER],
      });
      expect(await chain.getBalance(granter, 'inj')).toEqual(ONE_INJ - GRANT_FEE);
      expect((await chain.getAccount(granter))?.sequence).toEqual(INITIAL_SEQUENCE + BigInt(1));
    });

    it('records the sponsorship', async () => {
      const grant: GasSponsorshipGrant = await sponsor(makeService());

      expect(await allSponsorships()).toEqual([expect.objectContaining({
        chain_id: 'injective-1',
        granter_address: granter,
        grantee_address: GRANTEE,
        dydx_address: DYDX_ADDRESS,
        spend_limit: SPEND_LIMIT.toString(),
        fee_amount: GRANT_FEE.toString(),
        status: GasSponsorshipStatus.CONFIRMED,
        transaction_hash: grant.transactionHash,
        expires_at: EXPIRY_OF_GRANT_AT_START,
        created_at: START_TIME,
      })]);
    });

    it('releases the lock', async () => {
      await sponsor(makeService());

      expect(await redis.getAsync(getSponsorLockKey('injective-1', granter), redisClient)).toBeNull();
    });

    it('rounds the gas price up', async () => {
      chain.baseFee = BigInt(160_000_001);

      const grant: GasSponsorshipGrant = await sponsor(makeService());

      // 150% of 160,000,001 is 240,000,001.5.
      expect(grant.gasPrice).toEqual('240000002');
      expect(grant.spendLimit).toEqual(
        (BigInt(240_000_002) * BigInt(450_000) * BigInt(3)).toString(),
      );
    });

    it('sponsors at exactly the maximum base fee', async () => {
      chain.baseFee = DEFAULT_SETTINGS.maxBaseFee;

      const grant: GasSponsorshipGrant = await sponsor(makeService());

      expect(grant.status).toEqual(GasSponsorshipGrantStatus.CONFIRMED);
      expect(grant.spendLimit).toEqual((SPEND_LIMIT * BigInt(10)).toString());
    });

    it('sponsors an address holding a little of the fee denom, but not enough for a transfer', async () => {
      chain.setBalance(GRANTEE, 'inj', FEE_PER_ATTEMPT - BigInt(1));

      expect((await sponsor(makeService())).status).toEqual(GasSponsorshipGrantStatus.CONFIRMED);
    });

    it('sponsors an address holding exactly the minimum USDC', async () => {
      chain.setBalance(GRANTEE, USDC_DENOM, DEFAULT_SETTINGS.minUsdcBalance);

      expect((await sponsor(makeService())).status).toEqual(GasSponsorshipGrantStatus.CONFIRMED);
    });
  });

  describe('refusing', () => {
    afterEach(async () => {
      // None of these may sign, send or record anything.
      expect(signer.signCount).toEqual(0);
      expect(chain.submitted).toHaveLength(0);
      expect(await allSponsorships()).toHaveLength(0);
    });

    it('does nothing at all while disabled', async () => {
      const getBaseFee: jest.SpyInstance = jest.spyOn(chain, 'getBaseFee');
      const getBalance: jest.SpyInstance = jest.spyOn(chain, 'getBalance');
      const getAllowance: jest.SpyInstance = jest.spyOn(chain, 'getAllowance');

      await expectRefusal(
        sponsor(makeService({ enabled: false })),
        GasSponsorshipErrorCode.DISABLED,
      );
      expect(getBaseFee).not.toHaveBeenCalled();
      expect(getBalance).not.toHaveBeenCalled();
      expect(getAllowance).not.toHaveBeenCalled();
    });

    it.each([
      ['an address of another chain', DYDX_ADDRESS],
      ['an address with a bad checksum', `${addressOf(1).slice(0, -1)}q`],
      ['an address of the wrong length', toBech32('inj', new Uint8Array(32).fill(1))],
      ['something that is not an address', 'not-an-address'],
      ['an empty address', ''],
    ])('refuses %s', async (_name: string, address: string) => {
      await expectRefusal(
        sponsor(makeService(), address),
        GasSponsorshipErrorCode.INVALID_ADDRESS,
      );
    });

    it('refuses to sponsor the sponsor itself', async () => {
      chain.setBalance(granter, USDC_DENOM, BigInt(25_000_000));

      await expectRefusal(
        sponsor(makeService(), granter),
        GasSponsorshipErrorCode.INVALID_ADDRESS,
      );
    });

    it('refuses while dYdX is not accepting the deposits', async () => {
      depositOpen = false;

      await expectRefusal(sponsor(makeService()), GasSponsorshipErrorCode.DEPOSITS_CLOSED);
    });

    it('refuses while the base fee is above the maximum', async () => {
      chain.baseFee = DEFAULT_SETTINGS.maxBaseFee + BigInt(1);

      await expectRefusal(sponsor(makeService()), GasSponsorshipErrorCode.BASE_FEE_TOO_HIGH);
    });

    it('refuses an address holding less than the minimum USDC', async () => {
      chain.setBalance(GRANTEE, USDC_DENOM, DEFAULT_SETTINGS.minUsdcBalance - BigInt(1));

      await expectRefusal(
        sponsor(makeService()),
        GasSponsorshipErrorCode.INSUFFICIENT_USDC_BALANCE,
      );
    });

    it('refuses an address that can pay for a transfer itself', async () => {
      chain.setBalance(GRANTEE, 'inj', FEE_PER_ATTEMPT);

      await expectRefusal(sponsor(makeService()), GasSponsorshipErrorCode.ALREADY_HAS_GAS);
    });

    it('refuses when granting would take the last 24 hours over the daily budget', async () => {
      await expectRefusal(
        sponsor(makeService({ dailyBudget: SPEND_LIMIT + GRANT_FEE - BigInt(1) })),
        GasSponsorshipErrorCode.BUDGET_EXHAUSTED,
      );
    });

    it('refuses when the sponsor cannot cover the allowance and the grant fee', async () => {
      chain.setBalance(granter, 'inj', SPEND_LIMIT + GRANT_FEE - BigInt(1));

      await expectRefusal(sponsor(makeService()), GasSponsorshipErrorCode.SPONSOR_UNDERFUNDED);
    });

    it('refuses when the sponsor account does not exist on chain', async () => {
      const unknownSigner: LocalKeySigner = await LocalKeySigner.create('22'.repeat(32));
      const service: GasSponsorshipService = makeService({}, unknownSigner);
      chain.setBalance(service.granter, 'inj', ONE_INJ);

      await expectRefusal(sponsor(service), GasSponsorshipErrorCode.SPONSOR_UNDERFUNDED);
      expect(unknownSigner.signCount).toEqual(0);
    });

    it.each([
      'getBaseFee',
      'getBalance',
      'getAllowance',
      'getAccount',
    ])('reports the chain as unavailable when %s fails', async (method: string) => {
      chain.failNext[method] = true;

      await expectRefusal(sponsor(makeService()), GasSponsorshipErrorCode.CHAIN_UNAVAILABLE);
    });

    it('reports the chain as unavailable when the dYdX deposit check fails', async () => {
      const service: GasSponsorshipService = new GasSponsorshipService(
        DEFAULT_SETTINGS,
        chain,
        signer,
        redisClient,
        () => Promise.reject(new InjectiveClientError('canonical USDC state query failed')),
      );

      await expectRefusal(sponsor(service), GasSponsorshipErrorCode.CHAIN_UNAVAILABLE);
    });

    it('refuses as busy when the lock cannot be taken in time', async () => {
      await redis.lockWithExpiry(
        redisClient,
        getSponsorLockKey('injective-1', granter),
        'held-by-another-request',
        10_000,
      );

      await expectRefusal(
        sponsor(makeService({ lockWaitMs: 150 })),
        GasSponsorshipErrorCode.BUSY,
      );
      // The lock belongs to the other request and must be left alone.
      expect(
        await redis.getAsync(getSponsorLockKey('injective-1', granter), redisClient),
      ).toEqual('held-by-another-request');
    });
  });

  describe('when the signer does not hold the configured key', () => {
    it('fails before sending or recording anything', async () => {
      const otherKey: LocalKeySigner = await LocalKeySigner.create('22'.repeat(32));
      const mismatched: SponsorSigner = {
        publicKey: signer.publicKey,
        signDigest: (digest: Uint8Array) => otherKey.signDigest(digest),
      };

      await expect(sponsor(makeService({}, mismatched))).rejects.toThrow(
        'The sponsor signature does not match the configured sponsor public key',
      );
      expect(chain.submitted).toHaveLength(0);
      expect(await allSponsorships()).toHaveLength(0);
      expect(await redis.getAsync(getSponsorLockKey('injective-1', granter), redisClient)).toBeNull();
    });
  });

  describe('an address that already has an allowance', () => {
    it('returns the allowance again without granting another', async () => {
      const service: GasSponsorshipService = makeService();
      const first: GasSponsorshipGrant = await sponsor(service);

      const second: GasSponsorshipGrant = await sponsor(service);

      expect(second).toEqual(first);
      expect(chain.submitted).toHaveLength(1);
      expect(signer.signCount).toEqual(1);
      expect(await allSponsorships()).toHaveLength(1);
    });

    it('reports what is left of a partly used allowance', async () => {
      const service: GasSponsorshipService = makeService();
      const first: GasSponsorshipGrant = await sponsor(service);
      // One transfer's worth has been spent.
      chain.setAllowance({
        ...(await chain.getAllowance(granter, GRANTEE))!,
        spendLimit: { inj: SPEND_LIMIT - FEE_PER_ATTEMPT },
      });

      const second: GasSponsorshipGrant = await sponsor(service);

      expect(second).toEqual({
        ...first,
        spendLimit: (SPEND_LIMIT - FEE_PER_ATTEMPT).toString(),
      });
      expect(chain.submitted).toHaveLength(1);
    });

    it('returns the allowance even if the address no longer meets the balance checks', async () => {
      const service: GasSponsorshipService = makeService();
      await sponsor(service);
      chain.setBalance(GRANTEE, USDC_DENOM, BigInt(0));

      expect((await sponsor(service)).status).toEqual(GasSponsorshipGrantStatus.CONFIRMED);
      expect(chain.submitted).toHaveLength(1);
    });

    it.each([
      [
        'has too little left for one transfer',
        { spendLimit: { inj: FEE_PER_ATTEMPT - BigInt(1) } },
      ],
      [
        'expires in under five minutes',
        { expiration: new Date(new Date(START_TIME).getTime() + 5 * ONE_MINUTE_MS - 1) },
      ],
      [
        'does not cover the transfer message',
        { allowedMessages: ['/cosmos.bank.v1beta1.MsgSend'] },
      ],
      [
        'is not limited in what it can spend',
        { spendLimit: undefined },
      ],
      [
        'is of a kind this module does not issue',
        {
          recognized: false,
          spendLimit: undefined,
          expiration: undefined,
          allowedMessages: undefined,
        },
      ],
    ])('replaces an allowance that %s', async (_name: string, change: object) => {
      chain.setAllowance({
        granter,
        grantee: GRANTEE,
        recognized: true,
        spendLimit: { inj: SPEND_LIMIT },
        expiration: new Date(EXPIRY_OF_GRANT_AT_START),
        allowedMessages: [MSG_TRANSFER],
        ...change,
      });

      const grant: GasSponsorshipGrant = await sponsor(makeService());

      expect(chain.submitted).toHaveLength(1);
      expect(chain.submitted[0].messages.map((message) => message.type)).toEqual(
        ['revoke', 'grant'],
      );
      expect(grant.status).toEqual(GasSponsorshipGrantStatus.CONFIRMED);
      expect(grant.spendLimit).toEqual(SPEND_LIMIT.toString());
      expect(await chain.getAllowance(granter, GRANTEE)).toEqual({
        granter,
        grantee: GRANTEE,
        recognized: true,
        spendLimit: { inj: SPEND_LIMIT },
        expiration: new Date(EXPIRY_OF_GRANT_AT_START),
        allowedMessages: [MSG_TRANSFER],
      });
    });

    it('keeps an allowance with exactly enough left for one transfer and five minutes', async () => {
      chain.setAllowance({
        granter,
        grantee: GRANTEE,
        recognized: true,
        spendLimit: { inj: FEE_PER_ATTEMPT },
        expiration: new Date(new Date(START_TIME).getTime() + 5 * ONE_MINUTE_MS),
        allowedMessages: [MSG_TRANSFER, '/cosmos.bank.v1beta1.MsgSend'],
      });

      const grant: GasSponsorshipGrant = await sponsor(makeService());

      expect(chain.submitted).toHaveLength(0);
      expect(grant).toEqual(expect.objectContaining({
        status: GasSponsorshipGrantStatus.CONFIRMED,
        spendLimit: FEE_PER_ATTEMPT.toString(),
        transactionHash: undefined,
      }));
    });
  });

  describe('cooldown', () => {
    it('refuses another allowance within the cooldown once the first is used up', async () => {
      const service: GasSponsorshipService = makeService();
      await sponsor(service);
      chain.removeAllowance(granter, GRANTEE);
      advanceTime(ONE_HOUR_MS - 1000);

      await expectRefusal(sponsor(service), GasSponsorshipErrorCode.COOLDOWN);
      expect(chain.submitted).toHaveLength(1);
      expect(await allSponsorships()).toHaveLength(1);
    });

    it('grants again once the cooldown has passed', async () => {
      const service: GasSponsorshipService = makeService();
      await sponsor(service);
      chain.removeAllowance(granter, GRANTEE);
      advanceTime(ONE_HOUR_MS + 1000);

      const grant: GasSponsorshipGrant = await sponsor(service);

      expect(grant.status).toEqual(GasSponsorshipGrantStatus.CONFIRMED);
      expect(chain.submitted).toHaveLength(2);
      expect(chain.submitted[1].sequence).toEqual(INITIAL_SEQUENCE + BigInt(1));
    });

    it('applies to each address separately', async () => {
      const service: GasSponsorshipService = makeService();
      await sponsor(service);

      expect((await sponsor(service, OTHER_GRANTEE)).status).toEqual(
        GasSponsorshipGrantStatus.CONFIRMED,
      );
    });
  });

  describe('daily budget', () => {
    const budgetForOneGrant: Partial<GasSponsorshipSettings> = {
      dailyBudget: SPEND_LIMIT + GRANT_FEE,
    };

    it('counts each allowance in full plus its grant fee', async () => {
      const service: GasSponsorshipService = makeService(budgetForOneGrant);
      await sponsor(service);

      await expectRefusal(
        sponsor(service, OTHER_GRANTEE),
        GasSponsorshipErrorCode.BUDGET_EXHAUSTED,
      );
      expect(chain.submitted).toHaveLength(1);
    });

    it('frees up what was committed more than 24 hours ago', async () => {
      const service: GasSponsorshipService = makeService(budgetForOneGrant);
      await sponsor(service);
      advanceTime(24 * ONE_HOUR_MS + 1000);

      expect((await sponsor(service, OTHER_GRANTEE)).status).toEqual(
        GasSponsorshipGrantStatus.CONFIRMED,
      );
    });

    it('does not count a grant that failed', async () => {
      const service: GasSponsorshipService = makeService(budgetForOneGrant);
      chain.rejectNextBroadcast = { code: 13, rawLog: 'insufficient fees' };
      await expectRefusal(sponsor(service), GasSponsorshipErrorCode.GRANT_FAILED);

      expect((await sponsor(service, OTHER_GRANTEE)).status).toEqual(
        GasSponsorshipGrantStatus.CONFIRMED,
      );
    });
  });

  describe('when the grant transaction fails', () => {
    it('reports a transaction the node rejects, and lets the address try again', async () => {
      const service: GasSponsorshipService = makeService();
      chain.rejectNextBroadcast = { code: 13, rawLog: 'insufficient fees' };

      await expectRefusal(sponsor(service), GasSponsorshipErrorCode.GRANT_FAILED);

      // The rejected transaction does not exist on the chain, so its hash is not kept.
      expect(await allSponsorships()).toEqual([
        expect.objectContaining({ status: GasSponsorshipStatus.FAILED, transaction_hash: null }),
      ]);
      expect(await chain.getAllowance(granter, GRANTEE)).toBeUndefined();

      // A failed grant gave nothing, so it does not start the cooldown. The clock has not
      // moved, so the retry signs the very same transaction.
      const retry: GasSponsorshipGrant = await sponsor(service);
      expect(retry.status).toEqual(GasSponsorshipGrantStatus.CONFIRMED);
      expect(retry.transactionHash).toEqual(chain.submitted[0].txHash);
      expect(chain.submitted.map((tx: DecodedTx) => tx.sequence)).toEqual(
        [INITIAL_SEQUENCE, INITIAL_SEQUENCE],
      );
    });

    it('reports a transaction that is included but fails', async () => {
      const service: GasSponsorshipService = makeService();
      // Someone else grants the address an allowance from the sponsor in the same block, ahead
      // of this transaction, so its grant message fails.
      chain.beforeBlock = () => chain.setAllowance({
        granter,
        grantee: GRANTEE,
        recognized: true,
        spendLimit: { inj: BigInt(1) },
        expiration: new Date(EXPIRY_OF_GRANT_AT_START),
        allowedMessages: [MSG_TRANSFER],
      });

      await expectRefusal(sponsor(service), GasSponsorshipErrorCode.GRANT_FAILED);

      expect(await allSponsorships()).toEqual([
        expect.objectContaining({ status: GasSponsorshipStatus.FAILED }),
      ]);
      expect(await redis.getAsync(getSponsorLockKey('injective-1', granter), redisClient)).toBeNull();
    });
  });

  describe('when the sequence read from the chain is behind', () => {
    it('signs again with the sequence the node expects', async () => {
      const service: GasSponsorshipService = makeService();
      // The first grant stays in the mempool, so the chain still reports the old sequence.
      chain.autoInclude = false;
      const first: GasSponsorshipGrant = await sponsor(service);
      expect(first.status).toEqual(GasSponsorshipGrantStatus.PENDING);

      const second: GasSponsorshipGrant = await sponsor(service, OTHER_GRANTEE);

      expect(chain.submitted.map((tx: DecodedTx) => tx.sequence)).toEqual([
        INITIAL_SEQUENCE,
        // Rejected: the node expects the next sequence.
        INITIAL_SEQUENCE,
        INITIAL_SEQUENCE + BigInt(1),
      ]);
      expect(chain.getMempoolSize()).toEqual(2);
      expect(second.transactionHash).toEqual(chain.submitted[2].txHash);

      // One record per address, the second pointing at the transaction that was accepted.
      const sponsorships: GasSponsorshipFromDatabase[] = await allSponsorships();
      expect(sponsorships).toHaveLength(2);
      expect(sponsorships.map((s) => s.transaction_hash).sort()).toEqual(
        [first.transactionHash, second.transactionHash].sort(),
      );

      chain.includePending();
      expect(await chain.getAllowance(granter, GRANTEE)).toBeDefined();
      expect(await chain.getAllowance(granter, OTHER_GRANTEE)).toBeDefined();
    });

    it('corrects the sequence only once per grant', async () => {
      const service: GasSponsorshipService = makeService();
      chain.autoInclude = false;
      await sponsor(service);
      // The node reports an expected sequence that it then rejects as well.
      const broadcastTx = chain.broadcastTx.bind(chain);
      let calls: number = 0;
      jest.spyOn(chain, 'broadcastTx').mockImplementation((txBytes: Uint8Array) => {
        calls += 1;
        if (calls === 1) {
          chain.rejectNextBroadcast = {
            code: 32,
            rawLog: 'account sequence mismatch, expected 99, got 41: incorrect account sequence',
          };
        }
        return broadcastTx(txBytes);
      });

      await expectRefusal(
        sponsor(service, OTHER_GRANTEE),
        GasSponsorshipErrorCode.GRANT_FAILED,
      );

      expect(calls).toEqual(2);
      expect(chain.submitted.map((tx: DecodedTx) => tx.sequence)).toEqual(
        [INITIAL_SEQUENCE, INITIAL_SEQUENCE, BigInt(99)],
      );
      expect(chain.getMempoolSize()).toEqual(1);
      expect(await GasSponsorshipTable.findAll({ grantee_address: OTHER_GRANTEE })).toEqual([
        expect.objectContaining({ status: GasSponsorshipStatus.FAILED }),
      ]);
    });
  });

  describe('when the grant is not confirmed in time', () => {
    beforeEach(() => {
      chain.autoInclude = false;
    });

    it('reports the grant as pending, with its transaction hash', async () => {
      const grant: GasSponsorshipGrant = await sponsor(makeService());

      expect(grant).toEqual(expect.objectContaining({
        status: GasSponsorshipGrantStatus.PENDING,
        transactionHash: chain.submitted[0].txHash,
        spendLimit: SPEND_LIMIT.toString(),
        expiresAt: EXPIRY_OF_GRANT_AT_START,
      }));
      expect(await allSponsorships()).toEqual([
        expect.objectContaining({ status: GasSponsorshipStatus.PENDING }),
      ]);
    });

    it('returns the same pending grant when asked again, without sending another', async () => {
      const service: GasSponsorshipService = makeService();
      const first: GasSponsorshipGrant = await sponsor(service);
      advanceTime(ONE_MINUTE_MS);

      const second: GasSponsorshipGrant = await sponsor(service);

      expect(second).toEqual(first);
      expect(chain.submitted).toHaveLength(1);
      expect(signer.signCount).toEqual(1);
    });

    it('reports the grant as confirmed once the transaction is included', async () => {
      const service: GasSponsorshipService = makeService();
      const first: GasSponsorshipGrant = await sponsor(service);
      chain.includePending();

      const second: GasSponsorshipGrant = await sponsor(service);

      expect(second).toEqual({ ...first, status: GasSponsorshipGrantStatus.CONFIRMED });
      expect(chain.submitted).toHaveLength(1);
      expect(await allSponsorships()).toEqual([
        expect.objectContaining({ status: GasSponsorshipStatus.CONFIRMED }),
      ]);
    });

    it('grants again once a transaction that never made it into a block is given up on', async () => {
      const service: GasSponsorshipService = makeService();
      const first: GasSponsorshipGrant = await sponsor(service);
      chain.dropPending();
      chain.autoInclude = true;
      advanceTime(2 * ONE_MINUTE_MS + 1000);

      const second: GasSponsorshipGrant = await sponsor(service);

      expect(second.status).toEqual(GasSponsorshipGrantStatus.CONFIRMED);
      expect(second.transactionHash).not.toEqual(first.transactionHash);
      expect(chain.submitted).toHaveLength(2);
      expect((await allSponsorships()).map((s) => s.status).sort()).toEqual(
        [GasSponsorshipStatus.CONFIRMED, GasSponsorshipStatus.FAILED],
      );
    });

    it('settles a pending grant whose transaction was included but failed, then grants', async () => {
      const service: GasSponsorshipService = makeService();
      await sponsor(service);
      chain.beforeBlock = () => chain.setAllowance({
        granter,
        grantee: GRANTEE,
        recognized: true,
        spendLimit: { inj: BigInt(1) },
        expiration: new Date(EXPIRY_OF_GRANT_AT_START),
        allowedMessages: [MSG_TRANSFER],
      });
      chain.includePending();
      chain.beforeBlock = undefined;
      chain.autoInclude = true;

      const second: GasSponsorshipGrant = await sponsor(service);

      expect(second.status).toEqual(GasSponsorshipGrantStatus.CONFIRMED);
      expect(chain.submitted[1].messages.map((message) => message.type)).toEqual(
        ['revoke', 'grant'],
      );
      expect((await allSponsorships()).map((s) => s.status).sort()).toEqual(
        [GasSponsorshipStatus.CONFIRMED, GasSponsorshipStatus.FAILED],
      );
    });

    it('counts a pending grant against the budget', async () => {
      const service: GasSponsorshipService = makeService({
        dailyBudget: SPEND_LIMIT + GRANT_FEE,
      });
      await sponsor(service);

      await expectRefusal(
        sponsor(service, OTHER_GRANTEE),
        GasSponsorshipErrorCode.BUDGET_EXHAUSTED,
      );
    });
  });

  describe('when the outcome of broadcasting is unknown', () => {
    it('reports the grant as pending, and confirmed once it is seen on chain', async () => {
      const service: GasSponsorshipService = makeService();
      // The node takes the transaction, but the response is lost.
      chain.loseNextBroadcastResponse = true;

      const first: GasSponsorshipGrant = await sponsor(service);
      expect(first.status).toEqual(GasSponsorshipGrantStatus.PENDING);
      expect(first.transactionHash).toEqual(chain.submitted[0].txHash);

      const second: GasSponsorshipGrant = await sponsor(service);
      expect(second).toEqual({ ...first, status: GasSponsorshipGrantStatus.CONFIRMED });
      expect(chain.submitted).toHaveLength(1);
    });

    it('does not send a second grant while the first may still land', async () => {
      const service: GasSponsorshipService = makeService();
      // The request never reaches the node.
      chain.failNext.broadcastTx = true;

      const first: GasSponsorshipGrant = await sponsor(service);
      expect(first.status).toEqual(GasSponsorshipGrantStatus.PENDING);
      expect(signer.signCount).toEqual(1);

      advanceTime(ONE_MINUTE_MS);
      expect(await sponsor(service)).toEqual(first);
      expect(signer.signCount).toEqual(1);

      advanceTime(ONE_MINUTE_MS + 1000);
      const third: GasSponsorshipGrant = await sponsor(service);
      expect(third.status).toEqual(GasSponsorshipGrantStatus.CONFIRMED);
      expect(signer.signCount).toEqual(2);
    });

    it('keeps waiting for inclusion when looking the transaction up fails', async () => {
      const service: GasSponsorshipService = makeService();
      chain.failGetTx = true;

      const first: GasSponsorshipGrant = await sponsor(service);
      expect(first.status).toEqual(GasSponsorshipGrantStatus.PENDING);

      chain.failGetTx = false;
      expect(await sponsor(service)).toEqual({
        ...first,
        status: GasSponsorshipGrantStatus.CONFIRMED,
      });
    });
  });

  describe('concurrent requests', () => {
    it('grants one allowance when the same address asks twice at once', async () => {
      const service: GasSponsorshipService = makeService();

      const [first, second]: GasSponsorshipGrant[] = await Promise.all([
        sponsor(service),
        sponsor(service),
      ]);

      expect(first).toEqual(second);
      expect(first.status).toEqual(GasSponsorshipGrantStatus.CONFIRMED);
      expect(chain.submitted).toHaveLength(1);
      expect(await allSponsorships()).toHaveLength(1);
    });

    it('gives each of many addresses its own allowance, one sequence after another', async () => {
      // Every request waits its turn, however slowly the machine running the test gets to it.
      const service: GasSponsorshipService = makeService({ lockWaitMs: 25_000 });
      const grantees: string[] = [11, 12, 13, 14, 15].map(addressOf);
      grantees.forEach(fundGrantee);

      const grants: GasSponsorshipGrant[] = await Promise.all(
        grantees.map((grantee: string) => sponsor(service, grantee)),
      );

      expect(grants.map((grant) => grant.status)).toEqual(
        grantees.map(() => GasSponsorshipGrantStatus.CONFIRMED),
      );
      // No transaction was rejected for its sequence: each was signed with the next one.
      expect(chain.submitted.map((tx: DecodedTx) => tx.sequence)).toEqual(
        [0, 1, 2, 3, 4].map((offset: number) => INITIAL_SEQUENCE + BigInt(offset)),
      );
      for (const grantee of grantees) {
        expect(await chain.getAllowance(granter, grantee)).toBeDefined();
      }
      expect(await chain.getBalance(granter, 'inj')).toEqual(ONE_INJ - GRANT_FEE * BigInt(5));
    });

    it('holds concurrent requests to the daily budget', async () => {
      const service: GasSponsorshipService = makeService({
        dailyBudget: (SPEND_LIMIT + GRANT_FEE) * BigInt(2),
        lockWaitMs: 25_000,
      });
      const grantees: string[] = [21, 22, 23, 24].map(addressOf);
      grantees.forEach(fundGrantee);

      const outcomes: string[] = await Promise.all(grantees.map(
        (grantee: string) => sponsor(service, grantee).then(
          () => 'granted',
          (error: GasSponsorshipError) => error.code,
        ),
      ));

      expect(outcomes.sort()).toEqual([
        GasSponsorshipErrorCode.BUDGET_EXHAUSTED,
        GasSponsorshipErrorCode.BUDGET_EXHAUSTED,
        'granted',
        'granted',
      ]);
      expect(chain.submitted).toHaveLength(2);
    });
  });
});
