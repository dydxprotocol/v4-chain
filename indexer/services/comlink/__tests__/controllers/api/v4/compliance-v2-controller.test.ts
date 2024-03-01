import {
  ComplianceStatus,
  dbHelpers,
  testMocks,
  testConstants,
  ComplianceReason,
  ComplianceStatusTable,
  ComplianceStatusFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { getIpAddr } from '../../../../src/lib/utils';
import { sendRequest } from '../../../helpers/helpers';
import { RequestMethod } from '../../../../src/types';
import { stats } from '@dydxprotocol-indexer/base';
import { redis } from '@dydxprotocol-indexer/redis';
import { ratelimitRedis } from '../../../../src/caches/rate-limiters';
import { ComplianceControllerHelper } from '../../../../src/controllers/api/v4/compliance-controller';
import config from '../../../../src/config';
import { DateTime } from 'luxon';
import { ComplianceAction } from '../../../../src/controllers/api/v4/compliance-v2-controller';
import { ExtendedSecp256k1Signature, Secp256k1, sha256 } from '@cosmjs/crypto';

jest.mock('../../../../src/lib/utils', () => ({
  ...jest.requireActual('../../../../src/lib/utils'),
  getIpAddr: jest.fn(),
}));
jest.mock('@cosmjs/crypto', () => ({
  ...jest.requireActual('@cosmjs/crypto'),
  Secp256k1: {
    verifySignature: jest.fn(),
  },
  ExtendedSecp256k1Signature: {
    fromFixedLength: jest.fn(),
  },
}));

describe('ComplianceV2Controller', () => {
  const ipAddr: string = '192.168.1.1';

  const ipAddrMock: jest.Mock = (getIpAddr as unknown as jest.Mock);

  beforeAll(async () => {
    await dbHelpers.migrate();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET', () => {
    beforeEach(async () => {
      ipAddrMock.mockReturnValue(ipAddr);
      await testMocks.seedData();
    });

    afterEach(async () => {
      await redis.deleteAllAsync(ratelimitRedis.client);
      await dbHelpers.clearData();
      jest.clearAllMocks();
    });

    it('should return COMPLIANT for a non-restricted, non-dydx address', async () => {
      jest.spyOn(ComplianceControllerHelper.prototype, 'screen').mockImplementation(() => {
        return Promise.resolve({
          restricted: false,
        });
      });

      const response: any = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/compliance/screen/0x123',
      });
      expect(response.body.status).toEqual(ComplianceStatus.COMPLIANT);
    });

    it('should return BLOCKED/COMPLIANCE_PROVIDER for a restricted, non-dydx address', async () => {
      jest.spyOn(ComplianceControllerHelper.prototype, 'screen').mockImplementation(() => {
        return Promise.resolve({
          restricted: true,
        });
      });

      const response: any = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/compliance/screen/0x123',
      });
      expect(response.body.status).toEqual(ComplianceStatus.BLOCKED);
      expect(response.body.reason).toEqual(ComplianceReason.COMPLIANCE_PROVIDER);
    });

    it('should return BLOCKED & upsert for a restricted, dydx address without existing compliance status',
      async () => {
        jest.spyOn(ComplianceControllerHelper.prototype, 'screen').mockImplementation(() => {
          return Promise.resolve({
            restricted: true,
          });
        });

        let data: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll({}, [], {});
        expect(data).toHaveLength(0);

        const response: any = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/compliance/screen/${testConstants.defaultAddress}`,
        });
        expect(response.body.status).toEqual(ComplianceStatus.BLOCKED);
        expect(response.body.reason).toEqual(ComplianceReason.COMPLIANCE_PROVIDER);
        data = await ComplianceStatusTable.findAll({}, [], {});
        expect(data).toHaveLength(1);
        expect(data[0]).toEqual(expect.objectContaining({
          address: testConstants.defaultAddress,
          status: ComplianceStatus.BLOCKED,
          reason: ComplianceReason.COMPLIANCE_PROVIDER,
        }));
      });

    it('should return CLOSE_ONLY & update for a restricted, dydx address with existing compliance status',
      async () => {
        jest.spyOn(ComplianceControllerHelper.prototype, 'screen').mockImplementation(() => {
          return Promise.resolve({
            restricted: true,
          });
        });

        await ComplianceStatusTable.create({
          address: testConstants.defaultAddress,
          status: ComplianceStatus.FIRST_STRIKE,
        });
        let data: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll({}, [], {});
        expect(data).toHaveLength(1);

        const response: any = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/compliance/screen/${testConstants.defaultAddress}`,
        });
        expect(response.body.status).toEqual(ComplianceStatus.CLOSE_ONLY);
        expect(response.body.reason).toEqual(ComplianceReason.COMPLIANCE_PROVIDER);
        data = await ComplianceStatusTable.findAll({}, [], {});
        expect(data).toHaveLength(1);
        expect(data[0]).toEqual(expect.objectContaining({
          address: testConstants.defaultAddress,
          status: ComplianceStatus.CLOSE_ONLY,
          reason: ComplianceReason.COMPLIANCE_PROVIDER,
        }));
      });

    it('should return COMPLIANT for a non-restricted, dydx address', async () => {
      jest.spyOn(ComplianceControllerHelper.prototype, 'screen').mockImplementation(() => {
        return Promise.resolve({
          restricted: false,
        });
      });

      const response: any = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/compliance/screen/${testConstants.defaultAddress}`,
      });
      expect(response.body.status).toEqual(ComplianceStatus.COMPLIANT);
    });
  });

  describe('POST /setStatus', () => {
    beforeEach(async () => {
      ipAddrMock.mockReturnValue(ipAddr);
      await testMocks.seedData();
    });

    afterEach(async () => {
      await redis.deleteAllAsync(ratelimitRedis.client);
      await dbHelpers.clearData();
      jest.clearAllMocks();
    });

    it('should return 400 for non-dydx address', async () => {
      config.EXPOSE_SET_COMPLIANCE_ENDPOINT = true;
      await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/compliance/setStatus',
        body: {
          address: '0x123',
          status: ComplianceStatus.COMPLIANT,
        },
        expectedStatus: 400,
      });
    });

    it('should upsert db row for dydx address', async () => {
      let data: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll({}, [], {});
      expect(data).toHaveLength(0);
      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/compliance/setStatus',
        body: {
          address: testConstants.defaultAddress,
          status: ComplianceStatus.COMPLIANT,
        },
      });
      expect(response.body.status).toEqual(ComplianceStatus.COMPLIANT);
      data = await ComplianceStatusTable.findAll({}, [], {});
      expect(data).toHaveLength(1);
      expect(data[0]).toEqual(expect.objectContaining({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.COMPLIANT,
      }));
    });

    it('should update exisitng db row for dydx address', async () => {
      await ComplianceStatusTable.create({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.FIRST_STRIKE,
      });
      let data: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll({}, [], {});
      expect(data).toHaveLength(1);
      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/compliance/setStatus',
        body: {
          address: testConstants.defaultAddress,
          status: ComplianceStatus.COMPLIANT,
        },
      });
      expect(response.body.status).toEqual(ComplianceStatus.COMPLIANT);
      data = await ComplianceStatusTable.findAll({}, [], {});
      expect(data).toHaveLength(1);
      expect(data[0]).toEqual(expect.objectContaining({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.COMPLIANT,
      }));
    });
  });

  describe('POST /geoblock', () => {
    beforeEach(async () => {
      ipAddrMock.mockReturnValue(ipAddr);
      await testMocks.seedData();
      jest.mock('@cosmjs/crypto', () => ({
        Secp256k1: {
          verifySignature: jest.fn().mockResolvedValue(true),
        },
        ExtendedSecp256k1Signature: {
          fromFixedLength: jest.fn().mockResolvedValue({} as ExtendedSecp256k1Signature),
        },
      }));
      jest.spyOn(DateTime, 'now').mockReturnValue(DateTime.fromSeconds(1620000000)); // Mock current time
    });

    afterEach(async () => {
      await redis.deleteAllAsync(ratelimitRedis.client);
      await dbHelpers.clearData();
      jest.clearAllMocks();
      jest.restoreAllMocks();
    });

    it('should return 400 for non-dYdX address', async () => {
      await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/compliance/geoblock',
        body: {
          address: '0x123', // Non-dYdX address
          message: 'Test message',
          action: ComplianceAction.ONBOARD,
          signedMessage: sha256(Buffer.from('msg')),
          pubkey: new Uint8Array([/* public key bytes */]),
          timestamp: 1620000000,
        },
        expectedStatus: 400,
      });
    });

    it('should return 400 for invalid timestamp', async () => {
      await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/compliance/geoblock',
        body: {
          address: testConstants.defaultAddress,
          message: 'Test message',
          action: ComplianceAction.ONBOARD,
          signedMessage: sha256(Buffer.from('msg')),
          pubkey: new Uint8Array([/* public key bytes */]),
          timestamp: 1619996600, // More than 30 seconds difference
        },
        expectedStatus: 400,
      });
    });

    it('should return 400 for invalid signature', async () => {
      // Mock verifySignature to return false for this test
      (Secp256k1.verifySignature as jest.Mock).mockResolvedValueOnce(false);

      await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/compliance/geoblock',
        body: {
          address: testConstants.defaultAddress,
          message: 'Test message',
          action: ComplianceAction.ONBOARD,
          signedMessage: sha256(Buffer.from('msg')),
          pubkey: new Uint8Array([/* public key bytes */]),
          timestamp: 1620000000,
        },
        expectedStatus: 400,
      });
    });

    it('should process valid request', async () => {
      (Secp256k1.verifySignature as jest.Mock).mockResolvedValueOnce(true);

      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: '/v4/compliance/geoblock',
        body: {
          address: testConstants.defaultAddress,
          message: 'Test message',
          action: ComplianceAction.ONBOARD,
          signedMessage: sha256(Buffer.from('msg')),
          pubkey: new Uint8Array([/* public key bytes */]),
          timestamp: 1620000000, // Valid timestamp
        },
      });

      expect(response.status).toEqual(200);
      expect(response.body.status).toEqual(ComplianceStatus.COMPLIANT);
    });
  });
});
