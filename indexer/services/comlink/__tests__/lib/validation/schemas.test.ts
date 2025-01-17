import { RequestMethod } from '../../../src/types';
import { getQueryString, sendRequestToApp } from '../../helpers/helpers';
import { schemaTestApp } from './helpers';
import request from 'supertest';
import config from '../../../src/config';
import {
  testConstants,
  MAX_PARENT_SUBACCOUNTS,
  CHILD_SUBACCOUNT_MULTIPLIER,
} from '@dydxprotocol-indexer/postgres';

describe('schemas', () => {
  const positiveNonInteger: number = 3.2;
  const negativeInteger: number = -1;
  const zeroInteger: number = 0;
  const defaultSubaccountNumber: number = testConstants.defaultSubaccount.subaccountNumber;
  const defaultAddress: string = testConstants.defaultSubaccount.address;
  describe('CheckSubaccountSchema', () => {
    it.each([
      [
        'missingaddress',
        { subaccountNumber: defaultSubaccountNumber },
        'address',
        'address must be a valid dydx address',
      ],
      [
        'missing subaccountNumber',
        { address: defaultAddress },
        'subaccountNumber',
        'subaccountNumber must be a non-negative integer less than 128001',
      ],
      [
        'non-integer subaccountNumber',
        { address: defaultAddress, subaccountNumber: positiveNonInteger },
        'subaccountNumber',
        'subaccountNumber must be a non-negative integer less than 128001',
      ],
      [
        'negative subaccountNumber',
        { address: defaultAddress, subaccountNumber: negativeInteger },
        'subaccountNumber',
        'subaccountNumber must be a non-negative integer less than 128001',
      ],
      [
        'subaccountNumber greater than maximum subaccount number',
        {
          address: defaultAddress,
          subaccountNumber: MAX_PARENT_SUBACCOUNTS * CHILD_SUBACCOUNT_MULTIPLIER + 1,
        },
        'subaccountNumber',
        'subaccountNumber must be a non-negative integer less than 128001',
      ],
    ])('Returns 400 when validation fails: %s', async (
      _reason: string,
      queryParams: {
        address?: string,
        subaccountNumber?: number,
      },
      fieldWithError: string,
      expectedErrorMsg: string,
    ) => {
      const response: request.Response = await sendRequestToApp({
        type: RequestMethod.GET,
        path: `/v4/check-subaccount-schema?${getQueryString(queryParams)}`,
        expressApp: schemaTestApp,
        expectedStatus: 400,
      });

      expect(response.body).toEqual(expect.objectContaining({
        errors: expect.arrayContaining([
          expect.objectContaining({
            param: fieldWithError,
            msg: expectedErrorMsg,
          }),
        ]),
      }));
    });
  });

  describe('CheckLimitAndCreatedBeforeSchema', () => {
    it.each([
      [
        'non-integer limit',
        {
          limit: positiveNonInteger,
        },
        'limit',
        `limit must be a positive integer that is not greater than max: ${config.API_LIMIT_V4}`,
      ],
      [
        'limit equals 0',
        {
          limit: zeroInteger,
        },
        'limit',
        `limit must be a positive integer that is not greater than max: ${config.API_LIMIT_V4}`,
      ],
      [
        'negative limit',
        {
          limit: negativeInteger,
        },
        'limit',
        `limit must be a positive integer that is not greater than max: ${config.API_LIMIT_V4}`,
      ],
      [
        'limit > API LIMIT',
        {
          limit: config.API_LIMIT_V4 + 1,
        },
        'limit',
        `limit must be a positive integer that is not greater than max: ${config.API_LIMIT_V4}`,
      ],
      [
        'createdBeforeOrAt is not an ISO8601 formatted string',
        {
          createdBeforeOrAt: '0:0:0:0',
        },
        'createdBeforeOrAt',
        'Invalid value',
      ],
      [
        'negative createdBeforeOrAtHeight',
        {
          createdBeforeOrAtHeight: negativeInteger,
        },
        'createdBeforeOrAtHeight',
        'createdBeforeOrAtHeight must be a non-negative integer',
      ],
    ])('Returns 400 when validation fails: %s', async (
      _reason: string,
      queryParams: {
        limit?: number,
        createdBeforeOrAt?: string,
        createdBeforeOrAtHeight?: number,
      },
      fieldWithError: string,
      expectedErrorMsg: string,
    ) => {
      const response: request.Response = await sendRequestToApp({
        type: RequestMethod.GET,
        path: `/v4/check-limit-and-created-before-schema?${getQueryString(queryParams)}`,
        expressApp: schemaTestApp,
        expectedStatus: 400,
      });

      expect(response.body).toEqual(expect.objectContaining({
        errors: expect.arrayContaining([
          expect.objectContaining({
            param: fieldWithError,
            msg: expectedErrorMsg,
          }),
        ]),
      }));
    });
  });
});
