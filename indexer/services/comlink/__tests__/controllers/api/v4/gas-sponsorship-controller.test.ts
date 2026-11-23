import {
  INDEXER_COMPLIANCE_BLOCKED_PAYLOAD,
  INDEXER_GEOBLOCKED_PAYLOAD,
  isRestrictedCountryHeaders,
  isWhitelistedAddress,
} from '@dydxprotocol-indexer/compliance';
import {
  ComplianceReason,
  ComplianceStatus,
  ComplianceStatusTable,
  dbHelpers,
} from '@dydxprotocol-indexer/postgres';
import { redis } from '@dydxprotocol-indexer/redis';
import request from 'supertest';

import { ratelimitRedis } from '../../../../src/caches/rate-limiters';
import config, { configSchema } from '../../../../src/config';
import IndexV4 from '../../../../src/controllers/api/index-v4';
import {
  GasSponsorshipError,
  GasSponsorshipErrorCode,
  GasSponsorshipGrant,
  GasSponsorshipGrantStatus,
  getGasSponsorshipService,
} from '../../../../src/lib/gas-sponsorship';
import Server from '../../../../src/request-helpers/server';
import { BlockedCode, RequestMethod } from '../../../../src/types';
import { sendRequest } from '../../../helpers/helpers';

jest.mock('@dydxprotocol-indexer/compliance');
jest.mock('../../../../src/lib/gas-sponsorship', () => ({
  ...jest.requireActual('../../../../src/lib/gas-sponsorship'),
  getGasSponsorshipService: jest.fn(),
}));

const PATH: string = '/v4/gasSponsorship';
const DYDX_ADDRESS: string = 'dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4';
const INJECTIVE_ADDRESS: string = 'inj1frwsqymyg3wnfhy4mqrq9cghwwzst0tgsmfal7';
const BODY: { address: string, injectiveAddress: string } = {
  address: DYDX_ADDRESS,
  injectiveAddress: INJECTIVE_ADDRESS,
};
const GRANT: GasSponsorshipGrant = {
  chainId: 'injective-1',
  granter: 'inj1xdrkxgke8qzey9ze4cx48heav6lzkkszuc224r',
  grantee: INJECTIVE_ADDRESS,
  feeDenom: 'inj',
  spendLimit: '324000000000000',
  expiresAt: '2026-11-19T12:00:00.000Z',
  allowedMessages: ['/ibc.applications.transfer.v1.MsgTransfer'],
  status: GasSponsorshipGrantStatus.CONFIRMED,
  transactionHash: 'C554B3ACB4F98161CCFE7938399A73C4D452FBD9EF92F88C4F79E68896905BE4',
  gasLimit: '450000',
  gasPrice: '240000000',
};

describe('gas-sponsorship-controller#V4', () => {
  let sponsor: jest.Mock;

  const post = (
    body: object,
    expectedStatus: number,
  ): Promise<request.Response> => sendRequest({
    type: RequestMethod.POST,
    path: PATH,
    body,
    expectedStatus,
  });

  // The shared helper sends a request that has headers before it attaches the body.
  const postFrom = (ip: string): Promise<request.Response> => request(Server(IndexV4))
    .post(PATH)
    .set('cf-connecting-ip', ip)
    .send(BODY);

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(() => {
    config.GAS_SPONSORSHIP_ENABLED = true;
    sponsor = jest.fn().mockResolvedValue(GRANT);
    (getGasSponsorshipService as jest.Mock).mockReturnValue({ sponsor });
    (isWhitelistedAddress as jest.Mock).mockReturnValue(false);
    (isRestrictedCountryHeaders as unknown as jest.Mock).mockReturnValue(false);
  });

  afterEach(async () => {
    config.GAS_SPONSORSHIP_ENABLED = false;
    config.RATE_LIMIT_ENABLED = false;
    await redis.deleteAllAsync(ratelimitRedis.client);
    await dbHelpers.clearData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('while sponsorship is switched off', () => {
    beforeEach(() => {
      config.GAS_SPONSORSHIP_ENABLED = false;
    });

    it('answers that it is unavailable, without sponsoring', async () => {
      const response: request.Response = await post(BODY, 503);

      expect(response.body).toEqual({
        errors: [{
          msg: 'Gas sponsorship is not available',
          code: GasSponsorshipErrorCode.DISABLED,
        }],
      });
      expect(getGasSponsorshipService).not.toHaveBeenCalled();
      expect(sponsor).not.toHaveBeenCalled();
    });

    it('answers the same for a request it would otherwise reject', async () => {
      const response: request.Response = await post({ address: 'nonsense' }, 503);

      expect(response.body.errors[0].code).toEqual(GasSponsorshipErrorCode.DISABLED);
    });
  });

  it('is switched off unless configured otherwise', () => {
    // The value the config schema gives when the environment does not set the flag.
    const previous: string | undefined = process.env.GAS_SPONSORSHIP_ENABLED;
    delete process.env.GAS_SPONSORSHIP_ENABLED;
    try {
      expect(configSchema.GAS_SPONSORSHIP_ENABLED('GAS_SPONSORSHIP_ENABLED')).toEqual(false);
    } finally {
      if (previous !== undefined) {
        process.env.GAS_SPONSORSHIP_ENABLED = previous;
      }
    }
  });

  it('sponsors the Injective address and returns the grant', async () => {
    const response: request.Response = await post(BODY, 200);

    expect(response.body).toEqual(GRANT);
    expect(sponsor).toHaveBeenCalledTimes(1);
    expect(sponsor).toHaveBeenCalledWith({
      granteeAddress: INJECTIVE_ADDRESS,
      dydxAddress: DYDX_ADDRESS,
    });
  });

  it('ignores fields it does not know', async () => {
    await post({ ...BODY, spendLimit: '999999999999999999999', granter: 'inj1attacker' }, 200);

    expect(sponsor).toHaveBeenCalledWith({
      granteeAddress: INJECTIVE_ADDRESS,
      dydxAddress: DYDX_ADDRESS,
    });
  });

  describe('validation', () => {
    it.each([
      ['a missing dYdX address', { injectiveAddress: INJECTIVE_ADDRESS }, 'address must be a valid dydx address'],
      [
        'a dYdX address with a bad checksum',
        { ...BODY, address: `${DYDX_ADDRESS.slice(0, -1)}q` },
        'address must be a valid dydx address',
      ],
      [
        'an Injective address given as the dYdX address',
        { ...BODY, address: INJECTIVE_ADDRESS },
        'address must be a valid dydx address',
      ],
      [
        'a dYdX address that is not a string',
        { ...BODY, address: { $ne: '' } },
        'address must be a valid dydx address',
      ],
      [
        'a missing Injective address',
        { address: DYDX_ADDRESS },
        'injectiveAddress must be a valid Injective address',
      ],
      [
        'an Injective address with a bad checksum',
        { ...BODY, injectiveAddress: `${INJECTIVE_ADDRESS.slice(0, -1)}q` },
        'injectiveAddress must be a valid Injective address',
      ],
      [
        'a dYdX address given as the Injective address',
        { ...BODY, injectiveAddress: DYDX_ADDRESS },
        'injectiveAddress must be a valid Injective address',
      ],
      [
        'an Ethereum address given as the Injective address',
        { ...BODY, injectiveAddress: '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266' },
        'injectiveAddress must be a valid Injective address',
      ],
      [
        'a list of Injective addresses',
        { ...BODY, injectiveAddress: [INJECTIVE_ADDRESS, INJECTIVE_ADDRESS] },
        'injectiveAddress must be a valid Injective address',
      ],
    ])('rejects %s', async (_name: string, body: object, message: string) => {
      const response: request.Response = await post(body, 400);

      expect(response.body.errors).toEqual([expect.objectContaining({ msg: message })]);
      expect(sponsor).not.toHaveBeenCalled();
    });
  });

  describe('compliance', () => {
    it('refuses a dYdX address that is blocked', async () => {
      await ComplianceStatusTable.create({
        address: DYDX_ADDRESS,
        status: ComplianceStatus.BLOCKED,
        reason: ComplianceReason.COMPLIANCE_PROVIDER,
        createdAt: '2026-11-18T00:00:00.000Z',
        updatedAt: '2026-11-18T00:00:00.000Z',
      });

      const response: request.Response = await post(BODY, 403);

      expect(response.body).toEqual({
        errors: [{
          msg: INDEXER_COMPLIANCE_BLOCKED_PAYLOAD,
          code: BlockedCode.COMPLIANCE_BLOCKED,
        }],
      });
      expect(sponsor).not.toHaveBeenCalled();
    });

    it('refuses a request from a restricted location', async () => {
      (isRestrictedCountryHeaders as unknown as jest.Mock).mockReturnValue(true);

      const response: request.Response = await post(BODY, 403);

      expect(response.body).toEqual({
        errors: [{
          msg: INDEXER_GEOBLOCKED_PAYLOAD,
          code: BlockedCode.GEOBLOCKED,
        }],
      });
      expect(sponsor).not.toHaveBeenCalled();
    });

    it('sponsors a dYdX address that is compliant', async () => {
      await ComplianceStatusTable.create({
        address: DYDX_ADDRESS,
        status: ComplianceStatus.COMPLIANT,
        createdAt: '2026-11-18T00:00:00.000Z',
        updatedAt: '2026-11-18T00:00:00.000Z',
      });

      await post(BODY, 200);

      expect(sponsor).toHaveBeenCalledTimes(1);
    });
  });

  describe('rate limiting', () => {
    it('limits requests per IP address', async () => {
      config.RATE_LIMIT_ENABLED = true;

      for (let i: number = 0; i < config.RATE_LIMIT_GAS_SPONSORSHIP_POINTS; i++) {
        expect((await postFrom('203.0.113.7')).status).toEqual(200);
      }
      const response: request.Response = await postFrom('203.0.113.7');

      expect(response.status).toEqual(429);
      expect(response.body).toEqual({ errors: [{ msg: 'Too many requests' }] });
      expect(sponsor).toHaveBeenCalledTimes(config.RATE_LIMIT_GAS_SPONSORSHIP_POINTS);

      // Another IP address has its own allowance of requests.
      expect((await postFrom('203.0.113.8')).status).toEqual(200);
    });

    it('refuses while the rate-limit store is not ready, without sponsoring', async () => {
      config.RATE_LIMIT_ENABLED = true;
      const client: any = ratelimitRedis.client;
      const ready: boolean = client.ready;
      client.ready = false;
      try {
        const response: request.Response = await postFrom('203.0.113.9');
        expect(response.status).toEqual(503);
        expect(response.body).toEqual({ errors: [{ msg: 'Service temporarily unavailable' }] });
      } finally {
        client.ready = ready;
      }
      expect(sponsor).not.toHaveBeenCalled();
    });
  });

  describe('when the service refuses or fails', () => {
    it.each([
      [GasSponsorshipErrorCode.INVALID_ADDRESS, 400],
      [GasSponsorshipErrorCode.INSUFFICIENT_USDC_BALANCE, 403],
      [GasSponsorshipErrorCode.ALREADY_HAS_GAS, 403],
      [GasSponsorshipErrorCode.COOLDOWN, 429],
      [GasSponsorshipErrorCode.GRANT_FAILED, 502],
      [GasSponsorshipErrorCode.DISABLED, 503],
      [GasSponsorshipErrorCode.BUDGET_EXHAUSTED, 503],
      [GasSponsorshipErrorCode.BASE_FEE_TOO_HIGH, 503],
      [GasSponsorshipErrorCode.DEPOSITS_CLOSED, 503],
      [GasSponsorshipErrorCode.SPONSOR_UNDERFUNDED, 503],
      [GasSponsorshipErrorCode.BUSY, 503],
      [GasSponsorshipErrorCode.CHAIN_UNAVAILABLE, 503],
    ])('answers %s with status %d and the code', async (
      code: GasSponsorshipErrorCode,
      status: number,
    ) => {
      sponsor.mockRejectedValue(new GasSponsorshipError(code, 'the reason'));

      const response: request.Response = await post(BODY, status);

      expect(response.body).toEqual({ errors: [{ msg: 'the reason', code }] });
    });

    it('has a status for every refusal code', () => {
      Object.values(GasSponsorshipErrorCode).forEach((code: GasSponsorshipErrorCode) => {
        expect(new GasSponsorshipError(code, '').httpStatus).toEqual(expect.any(Number));
      });
    });

    it('answers an unexpected failure with an internal error that reveals nothing', async () => {
      sponsor.mockRejectedValue(new Error('turnkey credentials abc123 rejected'));

      const response: request.Response = await post(BODY, 500);

      expect(response.body).toEqual({ errors: [{ msg: 'Internal Server Error' }] });
    });

    it('answers that it is unavailable when the service cannot be built', async () => {
      (getGasSponsorshipService as jest.Mock).mockImplementation(() => {
        throw new GasSponsorshipError(
          GasSponsorshipErrorCode.DISABLED,
          'Gas sponsorship is not available',
        );
      });

      const response: request.Response = await post(BODY, 503);

      expect(response.body.errors[0].code).toEqual(GasSponsorshipErrorCode.DISABLED);
    });
  });
});
