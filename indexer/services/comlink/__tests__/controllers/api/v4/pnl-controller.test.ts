import {
  dbHelpers,
  testConstants,
  testMocks,
  PnlCreateObject,
  PnlTable,
} from '@dydxprotocol-indexer/postgres';
import { PnlResponseObject, RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { getQueryString, pnlCreateObjectToResponseObject, sendRequest } from '../../../helpers/helpers';

describe('pnl-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET', () => {
    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('Get /pnl', async () => {
      await testMocks.seedData();
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const createdAtHeight: string = '3';
      const pnl2: PnlCreateObject = {
        ...testConstants.defaultPnl,
        createdAt,
        createdAtHeight,
      };

      await Promise.all([
        PnlTable.create(testConstants.defaultPnl),
        PnlTable.create(pnl2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        })}`,
      });

      const expectedPnlResponse: PnlResponseObject = pnlCreateObjectToResponseObject({
        ...testConstants.defaultPnl,
      });

      const expectedPnl2Response: PnlResponseObject = pnlCreateObjectToResponseObject({
        ...testConstants.defaultPnl,
        createdAt,
        createdAtHeight,
      });

      expect(response.body.pnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnl2Response,
          }),
          expect.objectContaining({
            ...expectedPnlResponse,
          }),
        ]),
      );
    });

    it('Get /pnl respects pagination', async () => {
      await testMocks.seedData();
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const createdAtHeight: string = '1';
      const pnl2: PnlCreateObject = {
        ...testConstants.defaultPnl,
        createdAt,
        createdAtHeight,
      };

      await Promise.all([
        PnlTable.create(testConstants.defaultPnl),
        PnlTable.create(pnl2),
      ]);

      const responsePage1: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          page: 1,
          limit: 1,
        })}`,
      });

      const responsePage2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          page: 2,
          limit: 1,
        })}`,
      });

      const expectedPnlResponse: PnlResponseObject = pnlCreateObjectToResponseObject({
        ...testConstants.defaultPnl,
      });

      const expectedPnl2Response: PnlResponseObject = pnlCreateObjectToResponseObject({
        ...testConstants.defaultPnl,
        createdAt,
        createdAtHeight,
      });

      expect(responsePage1.body.pageSize).toStrictEqual(1);
      expect(responsePage1.body.offset).toStrictEqual(0);
      expect(responsePage1.body.totalResults).toStrictEqual(2);
      expect(responsePage1.body.pnl).toHaveLength(1);
      expect(responsePage1.body.pnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnlResponse,
          }),
        ]),
      );

      expect(responsePage2.body.pageSize).toStrictEqual(1);
      expect(responsePage2.body.offset).toStrictEqual(1);
      expect(responsePage2.body.totalResults).toStrictEqual(2);
      expect(responsePage2.body.pnl).toHaveLength(1);
      expect(responsePage2.body.pnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnl2Response,
          }),
        ]),
      );
    });

    it('Get /pnl respects createdBeforeOrAt and createdBeforeOrAtHeight field', async () => {
      await testMocks.seedData();
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const createdAtHeight: string = '1';
      const pnl2: PnlCreateObject = {
        ...testConstants.defaultPnl,
        createdAt,
        createdAtHeight,
      };

      await Promise.all([
        PnlTable.create(testConstants.defaultPnl),
        PnlTable.create(pnl2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          createdBeforeOrAt: createdAt,
          createdBeforeOrAtHeight: createdAtHeight,
        })}`,
      });

      const expectedPnl2Response: PnlResponseObject = pnlCreateObjectToResponseObject({
        ...testConstants.defaultPnl,
        createdAt,
        createdAtHeight,
      });

      expect(response.body.pnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnl2Response,
          }),
        ]),
      );
    });

    it('Get /pnl respects createdOnOrAfter and createdOnOrAfterHeight field', async () => {
      await testMocks.seedData();
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const createdAtHeight: string = '1';
      const pnl2: PnlCreateObject = {
        ...testConstants.defaultPnl,
        createdAt,
        createdAtHeight,
      };

      await Promise.all([
        PnlTable.create(testConstants.defaultPnl),
        PnlTable.create(pnl2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          createdOnOrAfter: testConstants.defaultPnl.createdAt,
          createdOnOrAfterHeight: testConstants.defaultPnl.createdAtHeight,
        })}`,
      });

      const expectedPnlResponse: PnlResponseObject = pnlCreateObjectToResponseObject({
        ...testConstants.defaultPnl,
      });

      expect(response.body.pnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnlResponse,
          }),
        ]),
      );
    });

    it('Get /pnl respects all created fields', async () => {
      await testMocks.seedData();
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const createdAtHeight: string = '1';
      const pnl2: PnlCreateObject = {
        ...testConstants.defaultPnl,
        createdAt,
        createdAtHeight,
      };

      await Promise.all([
        PnlTable.create(testConstants.defaultPnl),
        PnlTable.create(pnl2),
      ]);

      const createdOnOrAfter: string = '2001-05-25T00:00:00.000Z';
      const createdBeforeOrAt: string = '2005-05-25T00:00:00.000Z';

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          createdOnOrAfter,
          createdBeforeOrAt,
        })}`,
      });

      expect(response.body.pnl).toHaveLength(0);
    });

    it('Get /pnl returns empty when there are no pnl records', async () => {
      await testMocks.seedData();

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        })}`,
      });

      expect(response.body.pnl).toHaveLength(0);
    });

    it('Get /pnl with non-existent address and subaccount number returns 404', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: 'invalidaddress',
          subaccountNumber: 100,
        })}`,
        expectedStatus: 404,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: 'No subaccount found with address invalidaddress and subaccountNumber 100',
          },
        ],
      });
    });

    it('Get /pnl with daily=true returns daily records', async () => {
      await testMocks.seedData();

      // Create hourly records spanning 3 days
      const baseDate = new Date('2023-01-01T00:00:00.000Z');
      const hourlyRecords = [];

      // Create 72 hourly records (3 days)
      for (let i = 0; i < 72; i++) {
        const date = new Date(baseDate);
        date.setHours(baseDate.getHours() + i);

        hourlyRecords.push({
          ...testConstants.defaultPnl,
          createdAt: date.toISOString(),
          createdAtHeight: (1000 + i).toString(), // Incrementing heights
          equity: (1000 + i).toString(), // Different equity values to verify correct records
        });
      }

      // Insert all records
      await Promise.all(
        hourlyRecords.map((record) => PnlTable.create(record)),
      );

      // Test regular request (should return all records with pagination)
      const regularResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          limit: 10, // Just get the first 10
        })}`,
      });

      // Should have 10 records (due to limit)
      expect(regularResponse.body.pnl).toHaveLength(10);

      // Test with daily=true
      const dailyResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          daily: 'true',
        })}`,
      });

      expect(dailyResponse.body.pnl.length).toEqual(3);

      // Verify that records are spaced at least a day apart
      if (dailyResponse.body.pnl.length >= 2) {
        const timestamps = dailyResponse.body.pnl.map(
          (record: { createdAt: string | number | Date }) => new Date(record.createdAt).getTime());

        // Check time gaps between consecutive records (should be ~24h = 86400000ms)
        for (let i = 0; i < timestamps.length - 1; i++) {
          const gap = timestamps[i] - timestamps[i + 1];
          // Allow for some flexibility in the gap (at least 20 hours)
          expect(gap).toBeGreaterThanOrEqual(20 * 60 * 60 * 1000);
        }
      }

      // Test daily with pagination
      const dailyPageResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          daily: 'true',
          page: 1,
          limit: 2,
        })}`,
      });

      // Should have 2 records (due to pagination limit)
      expect(dailyPageResponse.body.pnl).toHaveLength(2);
      expect(dailyPageResponse.body.pageSize).toBe(2);
      expect(dailyPageResponse.body.offset).toBe(0);
    });
  });

});
