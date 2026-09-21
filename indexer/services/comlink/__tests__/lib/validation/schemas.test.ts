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

  describe('CheckPaginationSchema offset ceiling', () => {
    it('allows a page/limit combination exactly at MAX_PAGINATION_OFFSET', async () => {
      const limit: number = 1000;
      const page: number = config.MAX_PAGINATION_OFFSET / limit + 1;

      await sendRequestToApp({
        type: RequestMethod.GET,
        path: `/v4/check-pagination-schema?${getQueryString({ page, limit })}`,
        expressApp: schemaTestApp,
        expectedStatus: 200,
      });
    });

    it('rejects a page/limit combination past MAX_PAGINATION_OFFSET, without leaking the offset or threshold values', async () => {
      const limit: number = 1000;
      const page: number = config.MAX_PAGINATION_OFFSET / limit + 2;
      const expectedOffset: number = (page - 1) * limit;

      const response: request.Response = await sendRequestToApp({
        type: RequestMethod.GET,
        path: `/v4/check-pagination-schema?${getQueryString({ page, limit })}`,
        expressApp: schemaTestApp,
        expectedStatus: 400,
      });

      expect(response.body).toEqual(expect.objectContaining({
        errors: expect.arrayContaining([
          expect.objectContaining({
            param: 'page',
            // The message should name the check (a debugging clue for us) but must not disclose
            // the computed offset or the numeric MAX_PAGINATION_OFFSET threshold to the caller.
            msg: expect.stringContaining('exceeds MAX_PAGINATION_OFFSET'),
          }),
        ]),
      }));
      const msg: string = response.body.errors.find(
        (error: { param: string }) => error.param === 'page',
      ).msg;
      expect(msg).not.toContain(String(expectedOffset));
      expect(msg).not.toContain(String(config.MAX_PAGINATION_OFFSET));
    });

    it('rejects a deep page using the default limit when no limit is provided', async () => {
      const page: number = Math.floor(config.MAX_PAGINATION_OFFSET / config.API_LIMIT_V4) + 2;

      const response: request.Response = await sendRequestToApp({
        type: RequestMethod.GET,
        path: `/v4/check-pagination-schema?${getQueryString({ page })}`,
        expressApp: schemaTestApp,
        expectedStatus: 400,
      });

      expect(response.body).toEqual(expect.objectContaining({
        errors: expect.arrayContaining([
          expect.objectContaining({
            param: 'page',
          }),
        ]),
      }));
    });

    it('does not report an offset error for shallow pagination', async () => {
      await sendRequestToApp({
        type: RequestMethod.GET,
        path: `/v4/check-pagination-schema?${getQueryString({ page: 2, limit: 100 })}`,
        expressApp: schemaTestApp,
        expectedStatus: 200,
      });
    });

    it('rejects a duplicate `page` query key instead of silently passing (HTTP parameter pollution)', async () => {
      // Express/qs parses a repeated query key into an array (req.query.page = ['1', '999999999']).
      // `isInt` only validates the array's first element ('1', a valid page), so without an
      // explicit non-scalar check this would otherwise slip through both `isInt` and the offset
      // check below (Number(['1','999999999']) is NaN, which the offset guard used to treat as
      // "some other validator will catch it").
      const response: request.Response = await sendRequestToApp({
        type: RequestMethod.GET,
        path: '/v4/check-pagination-schema?page=1&page=999999999&limit=1000',
        expressApp: schemaTestApp,
        expectedStatus: 400,
      });

      expect(response.body).toEqual(expect.objectContaining({
        errors: expect.arrayContaining([
          expect.objectContaining({
            param: 'page',
            msg: 'page must be a single integer value, not a list or object',
          }),
        ]),
      }));
    });
  });
});
