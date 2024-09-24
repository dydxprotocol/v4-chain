import { RequestMethod } from '../../../../src/types';
import { sendRequest } from '../../../helpers/helpers';
import {
  ComplianceDataFromDatabase,
  ComplianceTable,
  dbHelpers,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { stats } from '@dydxprotocol-indexer/base';
import { complianceProvider } from '../../../../src/helpers/compliance/compliance-clients';
import {
  ComplianceClientResponse,
  INDEXER_COMPLIANCE_BLOCKED_PAYLOAD,
  NOT_IN_BLOCKCHAIN_RISK_SCORE,
} from '@dydxprotocol-indexer/compliance';
import { ratelimitRedis } from '../../../../src/caches/rate-limiters';
import { redis } from '@dydxprotocol-indexer/redis';
import { DateTime } from 'luxon';
import config from '../../../../src/config';
import { getIpAddr } from '../../../../src/lib/utils';

jest.mock('../../../../src/lib/utils', () => ({
  ...jest.requireActual('../../../../src/lib/utils'),
  getIpAddr: jest.fn(),
}));

describe('compliance-controller#V4', () => {
  const riskScore: string = '10.00';
  const blocked: boolean = false;
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
      jest.spyOn(complianceProvider.client, 'getComplianceResponse').mockImplementation(
        (address: string): Promise<ComplianceClientResponse> => {
          return Promise.resolve({
            address,
            blocked,
            riskScore,
          });
        },
      );
      ipAddrMock.mockReturnValue(ipAddr);
      await testMocks.seedData();
    });

    afterEach(async () => {
      await redis.deleteAllAsync(ratelimitRedis.client);
      await dbHelpers.clearData();
    });

    it('Get /screen with new address gets and stores compliance data from provider', async () => {
      let data: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {});
      expect(data).toHaveLength(0);

      const response: any = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/screen?address=${testConstants.defaultAddress}`,
      });

      expect(response.body.restricted).toEqual(false);
      expect(response.reason).toBeUndefined();
      expect(stats.timing).toHaveBeenCalledTimes(1);
      expect(stats.increment).toHaveBeenCalledWith(
        'comlink.compliance-controller.compliance_data_cache_miss',
        { provider: complianceProvider.provider },
      );
      expect(complianceProvider.client.getComplianceResponse).toHaveBeenCalledTimes(1);

      data = await ComplianceTable.findAll({}, [], {});
      expect(data).toHaveLength(1);
      expect(data[0]).toEqual(expect.objectContaining({
        address: testConstants.defaultAddress,
        provider: complianceProvider.provider,
        blocked,
        riskScore,
      }));
    });

    it(
      'Get /screen with existing non-blocked address retrieves compliance data from database',
      async () => {
      // Seed the database with a non-blocked compliance record
        await ComplianceTable.create(testConstants.nonBlockedComplianceData);
        let data: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {});
        expect(data).toHaveLength(1);

        const response: any = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/screen?address=${testConstants.defaultAddress}`,
        });

        expect(response.body.restricted).toEqual(false);
        expect(response.reason).toBeUndefined();
        expect(stats.timing).toHaveBeenCalledTimes(1);
        expect(stats.increment).toHaveBeenCalledWith(
          'comlink.compliance-controller.compliance_data_cache_hit',
          { provider: complianceProvider.provider },
        );
        expect(complianceProvider.client.getComplianceResponse).toHaveBeenCalledTimes(0);

        data = await ComplianceTable.findAll({}, [], {});
        expect(data).toHaveLength(1);
        expect(data[0]).toEqual(testConstants.nonBlockedComplianceData);
      });

    it(
      'Get /screen with existing blocked address retrieves compliance data from database',
      async () => {
      // Seed the database with a blocked compliance record
        await ComplianceTable.create(testConstants.blockedComplianceData);
        let data: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {});
        expect(data).toHaveLength(1);

        const response: any = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/screen?address=${testConstants.blockedAddress}`,
        });

        expect(response.body.restricted).toEqual(true);
        expect(response.body.reason).toEqual(INDEXER_COMPLIANCE_BLOCKED_PAYLOAD);
        expect(stats.timing).toHaveBeenCalledTimes(1);
        expect(stats.increment).toHaveBeenCalledWith(
          'comlink.compliance-controller.compliance_data_cache_hit',
          { provider: complianceProvider.provider },
        );
        expect(complianceProvider.client.getComplianceResponse).toHaveBeenCalledTimes(0);

        data = await ComplianceTable.findAll({}, [], {});
        expect(data).toHaveLength(1);
        expect(data[0]).toEqual(testConstants.blockedComplianceData);
      });

    it(
      'Get /screen with old existing non-blocked address refreshes compliance data',
      async () => {
      // Seed the database with an old non-blocked compliance record
        await ComplianceTable.create({
          ...testConstants.nonBlockedComplianceData,
          updatedAt: DateTime.utc().minus({
            seconds: config.MAX_AGE_SCREENED_ADDRESS_COMPLIANCE_DATA_SECONDS * 2,
          }).toISO(),
        });
        let data: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {});
        expect(data).toHaveLength(1);

        const response: any = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/screen?address=${testConstants.defaultAddress}`,
        });

        expect(response.body.restricted).toEqual(false);
        expect(response.body.reason).toBeUndefined();
        expect(stats.timing).toHaveBeenCalledTimes(1);
        expect(stats.increment).toHaveBeenCalledWith(
          'comlink.compliance-controller.compliance_data_cache_hit',
          { provider: complianceProvider.provider },
        );
        expect(stats.increment).toHaveBeenCalledWith(
          'comlink.compliance-controller.refresh_compliance_data_cache',
          { provider: complianceProvider.provider },
        );
        expect(complianceProvider.client.getComplianceResponse).toHaveBeenCalledTimes(1);

        data = await ComplianceTable.findAll({}, [], {});
        expect(data).toHaveLength(1);
        expect(data[0]).not.toEqual({
          address: testConstants.defaultAddress,
          provider: complianceProvider.provider,
          blocked,
          riskScore,
        });
      });

    it(
      'Get /screen with old existing blocked addressdoes not refresh compliance data',
      async () => {
      // Seed the database with an old blocked compliance record
        const oldUpdatedAt: string = DateTime.utc().minus({
          seconds: config.MAX_AGE_SCREENED_ADDRESS_COMPLIANCE_DATA_SECONDS * 2,
        }).toISO();
        await ComplianceTable.create({
          ...testConstants.blockedComplianceData,
          updatedAt: oldUpdatedAt,
        });
        let data: ComplianceDataFromDatabase[] = await ComplianceTable.findAll({}, [], {});
        expect(data).toHaveLength(1);

        const response: any = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/screen?address=${testConstants.blockedAddress}`,
        });

        expect(response.body.restricted).toEqual(true);
        expect(response.body.reason).toEqual(INDEXER_COMPLIANCE_BLOCKED_PAYLOAD);
        expect(stats.timing).toHaveBeenCalledTimes(1);
        expect(stats.increment).toHaveBeenCalledWith(
          'comlink.compliance-controller.compliance_data_cache_hit',
          { provider: complianceProvider.provider },
        );
        expect(complianceProvider.client.getComplianceResponse).toHaveBeenCalledTimes(0);

        data = await ComplianceTable.findAll({}, [], {});
        expect(data).toHaveLength(1);
        expect(data[0]).toEqual({
          ...testConstants.blockedComplianceData,
          updatedAt: oldUpdatedAt,
        });
      });

    it('Get /screen with multiple new address from same IP gets rate-limited', async () => {
      for (let i: number = 0; i < config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_POINTS; i++) {
        await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/screen?address=${i}`,
        });
      }

      await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/screen?address=${testConstants.defaultAddress}`,
        errorMsg: 'Too many requests',
        expectedStatus: 429,
      });
      expect(stats.increment).toHaveBeenCalledWith(
        'comlink.compliance-controller.compliance_screen_rate_limited_attempts',
        { provider: complianceProvider.provider },
      );
    });

    it('Get /screen with multiple new address globally gets rate-limited', async () => {
      ipAddrMock.mockImplementation(() => Math.random().toString());
      for (let i: number = 0; i < config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_POINTS; i++) {
        await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/screen?address=${i}`,
        });
      }

      await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/screen?address=${testConstants.defaultAddress}`,
        errorMsg: 'Too many requests',
        expectedStatus: 429,
      });
      expect(stats.increment).toHaveBeenCalledWith(
        'comlink.compliance-controller.compliance_screen_rate_limited_attempts',
        { provider: complianceProvider.provider },
      );
    });

    it('GET /screen for invalid address does not upsert compliance data', async () => {
      const invalidAddress: string = 'invalidAddress';
      const notInBlockchainRiskScore: string = NOT_IN_BLOCKCHAIN_RISK_SCORE.toString();

      jest.spyOn(complianceProvider.client, 'getComplianceResponse').mockImplementation(
        (address: string): Promise<ComplianceClientResponse> => {
          return Promise.resolve({
            address,
            blocked,
            riskScore: notInBlockchainRiskScore,
          });
        },
      );

      const response: any = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/screen?address=${invalidAddress}`,
      });

      expect(response.body).toEqual({
        restricted: false,
        reason: undefined,
      });

      const data = await ComplianceTable.findAll({}, [], {});
      expect(data).toHaveLength(0);
    });
  });
});
