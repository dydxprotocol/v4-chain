import {
  ComplianceReason,
  ComplianceStatus,
  ComplianceStatusFromDatabase,
  ComplianceStatusTable,
  dbHelpers,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { getIpAddr } from '../../../../src/lib/utils';
import { sendRequest } from '../../../helpers/helpers';
import { BlockedCode, RequestMethod } from '../../../../src/types';
import { stats } from '@dydxprotocol-indexer/base';
import { redis } from '@dydxprotocol-indexer/redis';
import { ratelimitRedis } from '../../../../src/caches/rate-limiters';
import { ComplianceControllerHelper } from '../../../../src/controllers/api/v4/compliance-controller';
import config from '../../../../src/config';
import { DateTime } from 'luxon';
import { ExtendedSecp256k1Signature, Secp256k1 } from '@cosmjs/crypto';
import { verifyADR36Amino } from '@keplr-wallet/cosmos';
import { getGeoComplianceReason, ComplianceAction } from '../../../../src/helpers/compliance/compliance-utils';
import {
  INDEXER_GEOBLOCKED_PAYLOAD,
  isRestrictedCountryHeaders,
  isWhitelistedAddress,
} from '@dydxprotocol-indexer/compliance';
import { toBech32 } from '@cosmjs/encoding';

jest.mock('@dydxprotocol-indexer/compliance');
jest.mock('../../../../src/helpers/compliance/compliance-utils', () => {
  const actualModule = jest.requireActual('../../../../src/helpers/compliance/compliance-utils');
  return {
    ...actualModule,
    getGeoComplianceReason: jest.fn(),
  };
});

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

jest.mock('@keplr-wallet/cosmos', () => ({
  ...jest.requireActual('@keplr-wallet/cosmos'),
  verifyADR36Amino: jest.fn(),
}));

jest.mock('@cosmjs/encoding', () => ({
  toBech32: jest.fn(),
}));

describe('ComplianceV2Controller', () => {
  const ipAddr: string = '192.168.1.1';

  const verifySignatureMock = Secp256k1.verifySignature as jest.Mock;
  const fromFixedLengthMock = ExtendedSecp256k1Signature.fromFixedLength as jest.Mock;
  const verifyADR36AminoMock = verifyADR36Amino as jest.Mock;
  const ipAddrMock = getIpAddr as jest.Mock;
  const toBech32Mock = toBech32 as jest.Mock;

  beforeAll(async () => {
    await dbHelpers.migrate();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET', () => {
    let isWhitelistedAddressSpy: jest.SpyInstance;

    beforeEach(async () => {
      ipAddrMock.mockReturnValue(ipAddr);
      isWhitelistedAddressSpy = isWhitelistedAddress as unknown as jest.Mock;
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
        expect(response.body.updatedAt).toBeDefined();
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
        expect(response.body.updatedAt).toBeDefined();
        data = await ComplianceStatusTable.findAll({}, [], {});
        expect(data).toHaveLength(1);
        expect(data[0]).toEqual(expect.objectContaining({
          address: testConstants.defaultAddress,
          status: ComplianceStatus.CLOSE_ONLY,
          reason: ComplianceReason.COMPLIANCE_PROVIDER,
        }));
      });

    it('should return COMPLIANT for a restricted, dydx address with existing CLOSE_ONLY compliance status', async () => {
      jest.spyOn(ComplianceControllerHelper.prototype, 'screen').mockImplementation(() => {
        return Promise.resolve({
          restricted: true,
        });
      });

      const createdAt: string = DateTime.utc().minus({ days: 1 }).toISO();
      await ComplianceStatusTable.create({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.CLOSE_ONLY,
        createdAt,
        updatedAt: createdAt,
      });
      const data: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll({}, [], {});
      expect(data).toHaveLength(1);

      isWhitelistedAddressSpy.mockReturnValue(true);
      const response: any = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/compliance/screen/${testConstants.defaultAddress}`,
      });
      expect(response.body.status).toEqual(ComplianceStatus.COMPLIANT);
    });

    it('should return CLOSE_ONLY & not update for a restricted, dydx address with existing CLOSE_ONLY compliance status',
      async () => {
        jest.spyOn(ComplianceControllerHelper.prototype, 'screen').mockImplementation(() => {
          return Promise.resolve({
            restricted: true,
          });
        });

        const createdAt: string = DateTime.utc().minus({ days: 1 }).toISO();
        await ComplianceStatusTable.create({
          address: testConstants.defaultAddress,
          status: ComplianceStatus.CLOSE_ONLY,
          createdAt,
          updatedAt: createdAt,
        });
        let data: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll({}, [], {});
        expect(data).toHaveLength(1);

        const response: any = await sendRequest({
          type: RequestMethod.GET,
          path: `/v4/compliance/screen/${testConstants.defaultAddress}`,
        });
        expect(response.body.status).toEqual(ComplianceStatus.CLOSE_ONLY);
        expect(response.body.updatedAt).toEqual(createdAt);
        data = await ComplianceStatusTable.findAll({}, [], {});
        expect(data).toHaveLength(1);
        expect(data[0]).toEqual(expect.objectContaining({
          address: testConstants.defaultAddress,
          status: ComplianceStatus.CLOSE_ONLY,
          createdAt,
          updatedAt: createdAt,
        }));
      },
    );

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

    it('should return existing compliance data for a non-restricted, dydx address', async () => {
      await ComplianceStatusTable.create({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.FIRST_STRIKE,
      });
      jest.spyOn(ComplianceControllerHelper.prototype, 'screen').mockImplementation(() => {
        return Promise.resolve({
          restricted: false,
        });
      });

      const response: any = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/compliance/screen/${testConstants.defaultAddress}`,
      });
      expect(response.body.status).toEqual(ComplianceStatus.FIRST_STRIKE);
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

    it('should update existing db row for dydx address', async () => {
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

  const geoblockEndpoint = '/v4/compliance/geoblock';
  const geoblockKeplerEndpoint = '/v4/compliance/geoblock-keplr';
  const geoblockBody = {
    address: testConstants.defaultAddress,
    message: 'Test message',
    action: ComplianceAction.CONNECT,
    signedMessage: 'signedmessage123',
    pubkey: 'asdfasdf',
    timestamp: 1620000000,
  };
  const geoblockKeplrBody = {
    address: testConstants.defaultAddress,
    message: 'Test message',
    action: ComplianceAction.CONNECT,
    signedMessage: 'signedmessage123',
    pubkey: 'asdfasdf',
  };
  const endpoints = [
    { endpoint: geoblockEndpoint, description: 'POST /geoblock', body: geoblockBody },
    { endpoint: geoblockKeplerEndpoint, description: 'POST /geoblock-keplr', body: geoblockKeplrBody },
  ];

  describe.each(endpoints)('$description endpoint', ({ endpoint, body }) => {
    let getGeoComplianceReasonSpy: jest.SpyInstance;
    let isRestrictedCountryHeadersSpy: jest.SpyInstance;
    let isWhitelistedAddressSpy: jest.SpyInstance;

    beforeEach(async () => {
      getGeoComplianceReasonSpy = getGeoComplianceReason as unknown as jest.Mock;
      isRestrictedCountryHeadersSpy = isRestrictedCountryHeaders as unknown as jest.Mock;
      isWhitelistedAddressSpy = isWhitelistedAddress as unknown as jest.Mock;
      ipAddrMock.mockReturnValue(ipAddr);
      await testMocks.seedData();
      // Mock verification to true to reduce mocking within individual tests
      verifySignatureMock.mockResolvedValue(true);
      fromFixedLengthMock.mockResolvedValue({} as ExtendedSecp256k1Signature);
      verifyADR36AminoMock.mockReturnValue(true);

      toBech32Mock.mockReturnValue(testConstants.defaultAddress);
      jest.spyOn(DateTime, 'now').mockReturnValue(DateTime.fromSeconds(1620000000)); // Mock current time
      jest.spyOn(stats, 'increment');
    });

    afterEach(async () => {
      await redis.deleteAllAsync(ratelimitRedis.client);
      await dbHelpers.clearData();
      jest.clearAllMocks();
      jest.restoreAllMocks();
    });

    it('should return 400 for non-dYdX address', async () => {
      if (endpoint === geoblockEndpoint) {
        await sendRequest({
          type: RequestMethod.POST,
          path: endpoint,
          body: {
            ...body,
            address: '0x123', // Non-dYdX address
          },
          expectedStatus: 400,
        });
      }
    });

    it('should return 400 for invalid timestamp', async () => {
      if (endpoint === geoblockEndpoint) {
        await sendRequest({
          type: RequestMethod.POST,
          path: endpoint,
          body: {
            ...body,
            timestamp: 1619996600, // More than 30 seconds difference
          },
          expectedStatus: 400,
        });
      }
    });

    it('should return 400 for invalid signature', async () => {
      if (endpoint === geoblockEndpoint) {
        // Mock verifySignature to return false for this test
        (Secp256k1.verifySignature as jest.Mock).mockResolvedValueOnce(false);
        await sendRequest({
          type: RequestMethod.POST,
          path: endpoint,
          body,
          expectedStatus: 400,
        });
      }
    });

    it('should return 400 for incorrect address', async () => {
      if (endpoint === geoblockEndpoint) {
        toBech32Mock.mockResolvedValueOnce('invalid_address');
        await sendRequest({
          type: RequestMethod.POST,
          path: endpoint,
          body,
          expectedStatus: 400,
        });
      }
    });

    it('should return 400 for failed keplr validation', async () => {
      if (endpoint === geoblockKeplerEndpoint) {
        (verifyADR36Amino as jest.Mock).mockReturnValueOnce(false);
        await sendRequest({
          type: RequestMethod.POST,
          path: endpoint,
          body,
          expectedStatus: 400,
        });
      }
    });

    it('should process valid request', async () => {
      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: endpoint,
        body,
      });

      expect(response.status).toEqual(200);
      expect(response.body.status).toEqual(ComplianceStatus.COMPLIANT);
    });

    it('should return COMPLIANT from a restricted country when whitelisted', async () => {
      getGeoComplianceReasonSpy.mockReturnValueOnce(ComplianceReason.US_GEO);
      isRestrictedCountryHeadersSpy.mockReturnValue(true);
      await dbHelpers.clearData();

      const data2: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll({}, [], {});
      expect(data2).toHaveLength(0);

      isWhitelistedAddressSpy.mockReturnValue(true);
      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: endpoint,
        body,
        expectedStatus: 200,
      });

      // expect compliance status to be empty
      const data: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll({}, [], {});
      expect(data).toHaveLength(0);

      // expect valid response
      expect(response.body.updatedAt).toBeDefined();
      expect(response.body.status).toEqual(ComplianceStatus.COMPLIANT);
    });

    it('should return 403 for CONNECT action from a restricted country with no existing compliance status', async () => {
      getGeoComplianceReasonSpy.mockReturnValueOnce(ComplianceReason.US_GEO);
      isRestrictedCountryHeadersSpy.mockReturnValue(true);
      await dbHelpers.clearData();

      expect(await ComplianceStatusTable.findAll({}, [], {})).toHaveLength(0);

      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: endpoint,
        body: {
          ...body,
          action: ComplianceAction.CONNECT,
        },
        expectedStatus: 403,
      });

      // Verify no database changes occurred
      expect(await ComplianceStatusTable.findAll({}, [], {})).toHaveLength(0);

      // Verify error response
      expect(response.body.errors).toBeDefined();
      expect(response.body.errors[0].msg).toEqual(INDEXER_GEOBLOCKED_PAYLOAD);
      expect(response.body.errors[0].code).toEqual(BlockedCode.GEOBLOCKED);

      // Verify no stats were incremented for status changes
      expect(stats.increment).not.toHaveBeenCalledWith(
        expect.stringContaining('compliance_status_changed'),
        expect.any(Object),
      );
    });

    it('should return 403 for CONNECT action from a restricted country with existing COMPLIANT status', async () => {
      await ComplianceStatusTable.create({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.COMPLIANT,
      });
      getGeoComplianceReasonSpy.mockReturnValueOnce(ComplianceReason.US_GEO);
      isRestrictedCountryHeadersSpy.mockReturnValue(true);

      const dataBefore: ComplianceStatusFromDatabase[] = await
      ComplianceStatusTable.findAll({}, [], {});

      expect(dataBefore).toHaveLength(1);

      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: endpoint,
        body: {
          ...body,
          action: ComplianceAction.CONNECT,
        },
        expectedStatus: 403,
      });

      // Verify no database changes occurred
      const dataAfter: ComplianceStatusFromDatabase[] = await
      ComplianceStatusTable.findAll({}, [], {});
      expect(dataAfter).toHaveLength(1);
      expect(dataAfter).toEqual(dataBefore);

      // Verify error response
      expect(response.body.errors).toBeDefined();
      expect(response.body.errors[0].msg).toEqual(INDEXER_GEOBLOCKED_PAYLOAD);
      expect(response.body.errors[0].code).toEqual(BlockedCode.GEOBLOCKED);
    });

    it('should return 403 for CONNECT action from a restricted country with existing FIRST_STRIKE status', async () => {
      await ComplianceStatusTable.create({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.FIRST_STRIKE,
        reason: ComplianceReason.US_GEO,
      });
      getGeoComplianceReasonSpy.mockReturnValueOnce(ComplianceReason.US_GEO);
      isRestrictedCountryHeadersSpy.mockReturnValue(true);

      const dataBefore: ComplianceStatusFromDatabase[] = await
      ComplianceStatusTable.findAll({}, [], {});

      expect(dataBefore).toHaveLength(1);
      const originalStatus = dataBefore[0].status;
      const originalUpdatedAt = dataBefore[0].updatedAt;

      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: endpoint,
        body: {
          ...body,
          action: ComplianceAction.CONNECT,
        },
        expectedStatus: 403,
      });

      // Verify no database changes occurred
      const dataAfter: ComplianceStatusFromDatabase[] = await
      ComplianceStatusTable.findAll({}, [], {});

      expect(dataAfter).toHaveLength(1);
      expect(dataAfter[0].status).toEqual(originalStatus);
      expect(dataAfter[0].updatedAt).toEqual(originalUpdatedAt);

      // Verify error response
      expect(response.body.errors).toBeDefined();
      expect(response.body.errors[0].msg).toEqual(INDEXER_GEOBLOCKED_PAYLOAD);
      expect(response.body.errors[0].code).toEqual(BlockedCode.GEOBLOCKED);
    });

    it('should return 403 for CONNECT action from a restricted country with existing CLOSE_ONLY status', async () => {
      const createdAt: string = DateTime.utc().minus({ days: 1 }).toISO();
      await ComplianceStatusTable.create({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.CLOSE_ONLY,
        reason: ComplianceReason.US_GEO,
        updatedAt: createdAt,
      });
      getGeoComplianceReasonSpy.mockReturnValueOnce(ComplianceReason.US_GEO);
      isRestrictedCountryHeadersSpy.mockReturnValue(true);

      expect(await ComplianceStatusTable.findAll({}, [], {})).toHaveLength(1);

      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: endpoint,
        body: {
          ...body,
          action: ComplianceAction.CONNECT,
        },
        expectedStatus: 403,
      });

      // Verify no database changes occurred
      const dataAfter: ComplianceStatusFromDatabase[] = await
      ComplianceStatusTable.findAll({}, [], {});

      expect(dataAfter).toHaveLength(1);
      expect(dataAfter[0].status).toEqual(ComplianceStatus.CLOSE_ONLY);
      expect(dataAfter[0].updatedAt).toEqual(createdAt);

      // Verify error response
      expect(response.body.errors).toBeDefined();
      expect(response.body.errors[0].msg).toEqual(INDEXER_GEOBLOCKED_PAYLOAD);
      expect(response.body.errors[0].code).toEqual(BlockedCode.GEOBLOCKED);
    });

    it('should return 403 for INVALID_SURVEY action from a restricted country', async () => {
      await ComplianceStatusTable.create({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.FIRST_STRIKE_CLOSE_ONLY,
        reason: ComplianceReason.US_GEO,
      });
      getGeoComplianceReasonSpy.mockReturnValueOnce(ComplianceReason.US_GEO);
      isRestrictedCountryHeadersSpy.mockReturnValue(true);

      const dataBefore: ComplianceStatusFromDatabase[] = await
      ComplianceStatusTable.findAll({}, [], {});

      expect(dataBefore).toHaveLength(1);
      const originalStatus = dataBefore[0].status;
      const originalUpdatedAt = dataBefore[0].updatedAt;

      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: endpoint,
        body: {
          ...body,
          action: ComplianceAction.INVALID_SURVEY,
        },
        expectedStatus: 403,
      });

      // Verify no database changes occurred
      const dataAfter: ComplianceStatusFromDatabase[] = await
      ComplianceStatusTable.findAll({}, [], {});

      expect(dataAfter).toHaveLength(1);
      expect(dataAfter[0].status).toEqual(originalStatus);
      expect(dataAfter[0].updatedAt).toEqual(originalUpdatedAt);

      // Verify error response
      expect(response.body.errors).toBeDefined();
      expect(response.body.errors[0].msg).toEqual(INDEXER_GEOBLOCKED_PAYLOAD);
      expect(response.body.errors[0].code).toEqual(BlockedCode.GEOBLOCKED);
    });

    it('should return 403 for VALID_SURVEY action from a restricted country', async () => {
      await ComplianceStatusTable.create({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.FIRST_STRIKE_CLOSE_ONLY,
        reason: ComplianceReason.US_GEO,
      });
      getGeoComplianceReasonSpy.mockReturnValueOnce(ComplianceReason.US_GEO);
      isRestrictedCountryHeadersSpy.mockReturnValue(true);

      const dataBefore: ComplianceStatusFromDatabase[] = await
      ComplianceStatusTable.findAll({}, [], {});

      expect(dataBefore).toHaveLength(1);
      const originalStatus = dataBefore[0].status;
      const originalUpdatedAt = dataBefore[0].updatedAt;

      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: endpoint,
        body: {
          ...body,
          action: ComplianceAction.VALID_SURVEY,
        },
        expectedStatus: 403,
      });

      // Verify no database changes occurred
      const dataAfter: ComplianceStatusFromDatabase[] = await
      ComplianceStatusTable.findAll({}, [], {});

      expect(dataAfter).toHaveLength(1);
      expect(dataAfter[0].status).toEqual(originalStatus);
      expect(dataAfter[0].updatedAt).toEqual(originalUpdatedAt);

      // Verify error response
      expect(response.body.errors).toBeDefined();
      expect(response.body.errors[0].msg).toEqual(INDEXER_GEOBLOCKED_PAYLOAD);
      expect(response.body.errors[0].code).toEqual(BlockedCode.GEOBLOCKED);
    });

    it.each([
      ComplianceStatus.FIRST_STRIKE,
      ComplianceStatus.FIRST_STRIKE_CLOSE_ONLY,
    ])('should not mutate status for any action from a non-restricted country with existing compliance status FIRST_STRIKE or FIRST_STRIKE_CLOSE_ONLY', async (status: ComplianceStatus) => {
      isRestrictedCountryHeadersSpy.mockReturnValue(false);

      const complianceStatus: ComplianceStatusFromDatabase = {
        address: testConstants.defaultAddress,
        createdAt: DateTime.utc().toISO(),
        updatedAt: DateTime.utc().toISO(),
        status,
      };

      await ComplianceStatusTable.create(complianceStatus);

      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: endpoint,
        body,
        expectedStatus: 200,
      });

      expect(response.body.updatedAt).toBeDefined();
      expect(response.body.status).toEqual(status);

      const data: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll({}, [], {});

      expect(data).toHaveLength(1);
      expect(data[0]).toEqual(expect.objectContaining(complianceStatus));
    });

    it('should leave status at BLOCKED for any action from a non-restricted country with existing compliance status BLOCKED', async () => {
      isRestrictedCountryHeadersSpy.mockReturnValue(false);

      await ComplianceStatusTable.create({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.BLOCKED,
      });

      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: endpoint,
        body,
        expectedStatus: 200,
      });

      expect(response.body.updatedAt).toBeDefined();
      expect(response.body.status).toEqual(ComplianceStatus.BLOCKED);

      const data: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll({}, [], {});

      expect(data).toHaveLength(1);
      expect(data[0]).toEqual(expect.objectContaining({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.BLOCKED,
      }));
    });

    it('should set status to COMPLIANT for any action from a non-restricted country with no existing compliance status', async () => {
      isRestrictedCountryHeadersSpy.mockReturnValue(false);

      const response: any = await sendRequest({
        type: RequestMethod.POST,
        path: endpoint,
        body,
        expectedStatus: 200,
      });

      const data: ComplianceStatusFromDatabase[] = await ComplianceStatusTable.findAll({}, [], {});
      expect(data).toHaveLength(1);
      expect(data[0]).toEqual(expect.objectContaining({
        address: testConstants.defaultAddress,
        status: ComplianceStatus.COMPLIANT,
      }));

      expect(response.body.status).toEqual(ComplianceStatus.COMPLIANT);
    });
  });
});
