import express from 'express';

import { BlockedCode, RequestMethod } from '../../src/types';
import Server from '../../src/request-helpers/server';
import { sendRequestToApp } from '../helpers/helpers';
import { complianceAndGeoCheck } from '../../src/lib/compliance-and-geo-check';
import { handleValidationErrors } from '../../src/request-helpers/error-handler';
import { checkSchema } from 'express-validator';
import {
  ComplianceStatus,
  ComplianceStatusTable,
  dbHelpers,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import request from 'supertest';
import {
  INDEXER_COMPLIANCE_BLOCKED_PAYLOAD,
  INDEXER_GEOBLOCKED_PAYLOAD,
  isRestrictedCountryHeaders,
  isWhitelistedAddress,
} from '@dydxprotocol-indexer/compliance';
import config from '../../src/config';

jest.mock('@dydxprotocol-indexer/compliance');

// Create a router to test the middleware with
const router: express.Router = express.Router();

const restrictedHeaders = {
  'cf-ipcountry': 'US',
};

const nonRestrictedHeaders = {
  'cf-ipcountry': 'SA',
};

router.get(
  '/check-compliance-query',
  checkSchema({
    address: {
      in: ['query'],
      isString: true,
      optional: true,
    },
  }),
  handleValidationErrors,
  complianceAndGeoCheck,
  (req: express.Request, res: express.Response) => {
    res.sendStatus(200);
  },
);

router.get(
  '/check-compliance-param/:address',
  checkSchema({
    address: {
      in: ['params'],
      isString: true,
    },
  }),
  handleValidationErrors,
  complianceAndGeoCheck,
  (req: express.Request, res: express.Response) => {
    res.sendStatus(200);
  },
);

export const complianceCheckApp = Server(router);

describe('compliance-check', () => {
  let isRestrictedCountrySpy: jest.SpyInstance;
  let isWhitelistedAddressSpy: jest.SpyInstance;

  beforeAll(async () => {
    config.INDEXER_LEVEL_GEOBLOCKING_ENABLED = true;
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    isRestrictedCountrySpy = isRestrictedCountryHeaders as unknown as jest.Mock;
    isWhitelistedAddressSpy = isWhitelistedAddress as jest.Mock;
    isWhitelistedAddressSpy.mockReturnValue(false);
    await testMocks.seedData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  afterEach(async () => {
    jest.restoreAllMocks();
    config.WHITELISTED_ADDRESSES = '';
    await dbHelpers.clearData();
  });

  it('does not return 403 if no address in request', async () => {
    isRestrictedCountrySpy.mockReturnValueOnce(false);
    await sendRequestToApp({
      type: RequestMethod.GET,
      path: '/v4/check-compliance-query',
      expressApp: complianceCheckApp,
      expectedStatus: 200,
    });
  });

  it.each([
    ['query', '/v4/check-compliance-query?address=random'],
    ['param', '/v4/check-compliance-param/random'],
  ])('does not return 403 if address in request is not in database (%s)', async (
    _name: string,
    path: string,
  ) => {
    isRestrictedCountrySpy.mockReturnValueOnce(false);
    await sendRequestToApp({
      type: RequestMethod.GET,
      path,
      expressApp: complianceCheckApp,
      expectedStatus: 200,
    });
  });

  it.each([
    ['query', '/v4/check-compliance-query?address=random'],
    ['param', '/v4/check-compliance-param/random'],
  ])('does not return 403 if address in request is not in database (%s) and non-restricted country', async (
    _name: string,
    path: string,
  ) => {
    isRestrictedCountrySpy.mockReturnValueOnce(false);
    await sendRequestToApp({
      type: RequestMethod.GET,
      path,
      expressApp: complianceCheckApp,
      expectedStatus: 200,
      headers: nonRestrictedHeaders,
    });
  });

  it.each([
    ['query', `/v4/check-compliance-query?address=${testConstants.defaultAddress}`],
    ['param', `/v4/check-compliance-param/${testConstants.defaultAddress}`],
  ])('does not return 403 if address in request is not blocked (%s)', async (
    _name: string,
    path: string,
  ) => {
    isRestrictedCountrySpy.mockReturnValueOnce(false);
    await ComplianceStatusTable.create(testConstants.compliantStatusData);
    await sendRequestToApp({
      type: RequestMethod.GET,
      path,
      expressApp: complianceCheckApp,
      expectedStatus: 200,
    });
  });

  it.each([
    ['query', `/v4/check-compliance-query?address=${testConstants.defaultAddress}`],
    ['param', `/v4/check-compliance-param/${testConstants.defaultAddress}`],
  ])('does not return 403 if address in request is in CLOSE_ONLY (%s)', async (
    _name: string,
    path: string,
  ) => {
    isRestrictedCountrySpy.mockReturnValueOnce(false);
    await ComplianceStatusTable.create({
      ...testConstants.compliantStatusData,
      status: ComplianceStatus.CLOSE_ONLY,
    });
    await sendRequestToApp({
      type: RequestMethod.GET,
      path,
      expressApp: complianceCheckApp,
      expectedStatus: 200,
    });
  });

  it.each([
    ['query', `/v4/check-compliance-query?address=${testConstants.defaultAddress}`],
    ['param', `/v4/check-compliance-param/${testConstants.defaultAddress}`],
  ])('does not return 403 if address in request is in CLOSE_ONLY and from restricted country (%s)', async (
    _name: string,
    path: string,
  ) => {
    isRestrictedCountrySpy.mockReturnValueOnce(true);
    await ComplianceStatusTable.create({
      ...testConstants.compliantStatusData,
      status: ComplianceStatus.CLOSE_ONLY,
    });
    await sendRequestToApp({
      type: RequestMethod.GET,
      path,
      expressApp: complianceCheckApp,
      expectedStatus: 200,
    });
  });

  it.each([
    ['query', `/v4/check-compliance-query?address=${testConstants.defaultAddress}`],
    ['param', `/v4/check-compliance-param/${testConstants.defaultAddress}`],
  ])('does not return 403 if address in request is in FIRST_STRIKE_CLOSE_ONLY and from restricted country (%s)', async (
    _name: string,
    path: string,
  ) => {
    isRestrictedCountrySpy.mockReturnValueOnce(true);
    await ComplianceStatusTable.create({
      ...testConstants.compliantStatusData,
      status: ComplianceStatus.FIRST_STRIKE_CLOSE_ONLY,
    });
    await sendRequestToApp({
      type: RequestMethod.GET,
      path,
      expressApp: complianceCheckApp,
      expectedStatus: 200,
    });
  });

  it.each([
    ['query', `/v4/check-compliance-query?address=${testConstants.defaultAddress}`],
    ['param', `/v4/check-compliance-param/${testConstants.defaultAddress}`],
  ])('does return 403 if request is from restricted country (%s)', async (
    _name: string,
    path: string,
  ) => {
    isRestrictedCountrySpy.mockReturnValueOnce(true);
    const response: request.Response = await sendRequestToApp({
      type: RequestMethod.GET,
      path,
      expressApp: complianceCheckApp,
      expectedStatus: 403,
      headers: restrictedHeaders,
    });

    expect(response.body).toEqual(expect.objectContaining({
      errors: expect.arrayContaining([{
        msg: INDEXER_GEOBLOCKED_PAYLOAD,
        code: BlockedCode.GEOBLOCKED,
      }]),
    }));
  });

  it.each([
    ['query', `/v4/check-compliance-query?address=${testConstants.defaultAddress}`],
    ['param', `/v4/check-compliance-param/${testConstants.defaultAddress}`],
  ])('does not return 403 if address is whitelisted and request is from restricted country (%s)', async (
    _name: string,
    path: string,
  ) => {
    isWhitelistedAddressSpy.mockReturnValue(true);
    isRestrictedCountrySpy.mockReturnValueOnce(true);
    await sendRequestToApp({
      type: RequestMethod.GET,
      path,
      expressApp: complianceCheckApp,
      expectedStatus: 200,
    });
  });

  it.each([
    ['query', `/v4/check-compliance-query?address=${testConstants.blockedAddress}`],
    ['param', `/v4/check-compliance-param/${testConstants.blockedAddress}`],
  ])('does return 403 if address in request is blocked (%s)', async (
    _name: string,
    path: string,
  ) => {
    isRestrictedCountrySpy.mockReturnValueOnce(false);
    await ComplianceStatusTable.create(testConstants.noncompliantStatusData);
    const response: request.Response = await sendRequestToApp({
      type: RequestMethod.GET,
      path,
      expressApp: complianceCheckApp,
      expectedStatus: 403,
    });

    expect(response.body).toEqual(expect.objectContaining({
      errors: expect.arrayContaining([{
        msg: INDEXER_COMPLIANCE_BLOCKED_PAYLOAD,
        code: BlockedCode.COMPLIANCE_BLOCKED,
      }]),
    }));
  });
});
