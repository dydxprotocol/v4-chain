import express from 'express';

import { BlockedCode, RequestMethod } from '../../src/types';
import Server from '../../src/request-helpers/server';
import { sendRequestToApp } from '../helpers/helpers';
import { complianceCheck } from '../../src/lib/compliance-check';
import { handleValidationErrors } from '../../src/request-helpers/error-handler';
import { checkSchema } from 'express-validator';
import {
  ComplianceTable, dbHelpers, testConstants, testMocks,
} from '@dydxprotocol-indexer/postgres';
import { blockedComplianceData, nonBlockedComplianceData } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import request from 'supertest';
import { INDEXER_COMPLIANCE_BLOCKED_PAYLOAD } from '@dydxprotocol-indexer/compliance';

// Create a router to test the middleware with
const router: express.Router = express.Router();

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
  complianceCheck,
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
  complianceCheck,
  (req: express.Request, res: express.Response) => {
    res.sendStatus(200);
  },
);

export const complianceCheckApp = Server(router);

describe('compliance-check', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await testMocks.seedData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  it('does not return 403 if no address in request', async () => {
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
  ])('does not return 403 if address in request is not blocked (%s)', async (
    _name: string,
    path: string,
  ) => {
    await ComplianceTable.create(nonBlockedComplianceData);
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
    await ComplianceTable.create(blockedComplianceData);
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
