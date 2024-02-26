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

jest.mock('../../../../src/lib/utils', () => ({
  ...jest.requireActual('../../../../src/lib/utils'),
  getIpAddr: jest.fn(),
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
});
