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
        date.setUTCHours(baseDate.getUTCHours() + i);

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

      // Verify the record structure:
      // 1. Latest record from day 3 (hour 71)
      // 2. Earliest record from day 2 (hour 24)
      // 3. Earliest record from day 1 (hour 0)

      // First record should be latest from day 3 (hour 71)
      expect(dailyResponse.body.pnl[0].createdAtHeight).toBe('1071');
      expect(dailyResponse.body.pnl[0].createdAt).toBe('2023-01-03T23:00:00.000Z');

      // Second record should be earliest from day 2 (hour 24)
      expect(dailyResponse.body.pnl[1].createdAtHeight).toBe('1024');
      expect(dailyResponse.body.pnl[1].createdAt).toBe('2023-01-02T00:00:00.000Z');

      // Third record should be earliest from day 1 (hour 0)
      expect(dailyResponse.body.pnl[2].createdAtHeight).toBe('1000');
      expect(dailyResponse.body.pnl[2].createdAt).toBe('2023-01-01T00:00:00.000Z');

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

      // Should contain the first two daily records
      expect(dailyPageResponse.body.pnl[0].createdAtHeight).toBe('1071');
      expect(dailyPageResponse.body.pnl[1].createdAtHeight).toBe('1024');
    });

    it('Get /pnl with daily=false returns regular hourly records', async () => {
      await testMocks.seedData();

      // Create hourly records spanning 3 days
      const baseDate = new Date('2023-01-01T00:00:00.000Z');
      const hourlyRecords = [];

      // Create 72 hourly records (3 days)
      for (let i = 0; i < 72; i++) {
        const date = new Date(baseDate);
        date.setUTCHours(baseDate.getUTCHours() + i);

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

      // Test with daily=false explicitly
      const regularResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          daily: 'false',
          limit: 10,
        })}`,
      });

      // Should have 10 records (due to limit)
      expect(regularResponse.body.pnl).toHaveLength(10);

      // Verify these are hourly records by checking timestamp spacing
      if (regularResponse.body.pnl.length >= 2) {
        const timestamps = regularResponse.body.pnl.map(
          (record: { createdAt: string | number | Date }) => new Date(record.createdAt).getTime(),
        );

        // Check time gaps between consecutive records (should be ~1h = 3600000ms)
        for (let i = 0; i < Math.min(3, timestamps.length - 1); i++) {
          const gap = timestamps[i] - timestamps[i + 1];
          // Expect gap to be around 1 hour (with some flexibility)
          expect(gap).toBeGreaterThanOrEqual(3500000); // ~58 minutes
          expect(gap).toBeLessThanOrEqual(3700000);   // ~62 minutes
        }
      }
    });

    it('Get /pnl with invalid daily parameter returns 400 error', async () => {
      await testMocks.seedData();

      // Create some test records
      await PnlTable.create(testConstants.defaultPnl);

      // Test with invalid daily parameter
      const invalidResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          daily: 'invalid',
        })}`,
        expectedStatus: 400,
      });

      // Check that we get the expected validation error
      expect(invalidResponse.body.errors).toBeDefined();
      expect(invalidResponse.body.errors.length).toBeGreaterThan(0);
      expect(invalidResponse.body.errors[0].param).toBe('daily');
    });
  });

  describe('Get /pnl/parentSubaccountNumber', () => {
    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('Get /pnl/parentSubaccountNumber returns aggregated PNL across child subaccounts', async () => {
      await testMocks.seedData();

      // Create PNL records for two different subaccounts (children of the same parent)
      const pnl1: PnlCreateObject = {
        ...testConstants.defaultPnl,
        equity: '1000.00',
        totalPnl: '100.00',
        netTransfers: '900.00',
      };

      const pnl2: PnlCreateObject = {
        ...testConstants.defaultPnl,
        subaccountId: testConstants.isolatedSubaccountId,
        equity: '2000.00',
        totalPnl: '200.00',
        netTransfers: '1800.00',
      };

      await Promise.all([
        PnlTable.create(pnl1),
        PnlTable.create(pnl2),
      ]);

      const parentSubaccountNumber: number = 0;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber,
        })}`,
      });

      // Check for the aggregated values but with the correct format (without decimal places)
      expect(response.body.pnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            equity: '3000',
            totalPnl: '300',
            netTransfers: '2700',
            createdAt: testConstants.defaultPnl.createdAt,
          }),
        ]),
      );
    });

    it('Get /pnl/parentSubaccountNumber with daily=true returns daily aggregated PNL', async () => {
      await testMocks.seedData();

      // Create hourly PNL records spanning 2 days for two different subaccounts
      const baseDate = new Date('2023-01-01T00:00:00.000Z');
      const hourlyRecords = [];

      // Create 48 hourly records (2 days) for first subaccount
      for (let i = 0; i < 48; i++) {
        const date = new Date(baseDate);
        date.setHours(baseDate.getHours() + i);

        hourlyRecords.push({
          ...testConstants.defaultPnl,
          createdAt: date.toISOString(),
          createdAtHeight: (1000 + i).toString(),
          equity: (1000 + i).toString(),
          totalPnl: (100 + i).toString(),
          netTransfers: (900 + i).toString(),
        });
      }

      // Create 48 hourly records (2 days) for second subaccount
      for (let i = 0; i < 48; i++) {
        const date = new Date(baseDate);
        date.setHours(baseDate.getHours() + i);

        hourlyRecords.push({
          ...testConstants.defaultPnl,
          subaccountId: testConstants.isolatedSubaccountId,
          createdAt: date.toISOString(),
          createdAtHeight: (2000 + i).toString(),
          equity: (2000 + i).toString(),
          totalPnl: (200 + i).toString(),
          netTransfers: (1800 + i).toString(),
        });
      }

      // Insert all records
      await Promise.all(
        hourlyRecords.map((record) => PnlTable.create(record)),
      );

      // Test with daily=true
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          daily: 'true',
        })}`,
      });

      // Should have 2 records (one for each day)
      expect(response.body.pnl.length).toEqual(2);

      // Get the records for further verification
      const day1Record = response.body.pnl.find(
        (r: any) => r.createdAt.startsWith('2023-01-01'),
      );
      const day2Record = response.body.pnl.find(
        (r: any) => r.createdAt.startsWith('2023-01-02'),
      );

      // Verify both day records exist
      expect(day1Record).toBeDefined();
      expect(day2Record).toBeDefined();

      // Expected values for day 1 (the earliest/first hour of day 1)
      // Earliest record for day 1 for first subaccount is i=0:
      //    equity=1000, totalPnl=100, netTransfers=900
      // Earliest record for day 1 for second subaccount is i=0:
      //    equity=2000, totalPnl=200, netTransfers=1800
      // Total: equity=3000, totalPnl=300, netTransfers=2700
      expect(day1Record.createdAt).toBe('2023-01-01T00:00:00.000Z');
      expect(Number(day1Record.equity)).toEqual(3000);
      expect(Number(day1Record.totalPnl)).toEqual(300);
      expect(Number(day1Record.netTransfers)).toEqual(2700);

      // Expected values for day 2 (the last hour of day 2)
      // Last record for day 2 for first subaccount is i=47:
      //    equity=1047, totalPnl=147, netTransfers=947
      // Last record for day 2 for second subaccount is i=47:
      //    equity=2047, totalPnl=247, netTransfers=1847
      // Total: equity=3094, totalPnl=394, netTransfers=2794

      // Verify day 2 values
      expect(day2Record.createdAt).toBe('2023-01-02T23:00:00.000Z');
      expect(Number(day2Record.equity)).toEqual(3094);
      expect(Number(day2Record.totalPnl)).toEqual(394);
      expect(Number(day2Record.netTransfers)).toEqual(2794);

      // Verify the records are in descending order by date (day 2 should come before day 1)
      const timestamps = response.body.pnl.map(
        (record: any) => new Date(record.createdAt).getTime());
      expect(timestamps[0]).toBeGreaterThan(timestamps[1]);
    });

    it('Get /pnl/parentSubaccountNumber respects filtering parameters', async () => {
      await testMocks.seedData();

      // Create PNL records with different timestamps and heights
      const pnl1: PnlCreateObject = {
        ...testConstants.defaultPnl,
        createdAt: '2023-01-01T00:00:00.000Z',
        createdAtHeight: '1000',
      };

      const pnl2: PnlCreateObject = {
        ...testConstants.defaultPnl,
        subaccountId: testConstants.isolatedSubaccountId,
        createdAt: '2023-01-02T00:00:00.000Z',
        createdAtHeight: '2000',
      };

      const pnl3: PnlCreateObject = {
        ...testConstants.defaultPnl,
        createdAt: '2023-01-03T00:00:00.000Z',
        createdAtHeight: '3000',
      };

      await Promise.all([
        PnlTable.create(pnl1),
        PnlTable.create(pnl2),
        PnlTable.create(pnl3),
      ]);

      // Test with createdBeforeOrAt filter
      const responseBefore: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          createdBeforeOrAt: '2023-01-02T12:00:00.000Z',
        })}`,
      });

      // Should only include records before or at the specified time
      expect(responseBefore.body.pnl.length).toBeGreaterThanOrEqual(1);
      expect(responseBefore.body.pnl.every(
        (record: any) => new Date(record.createdAt) <= new Date('2023-01-02T12:00:00.000Z'),
      )).toBe(true);

      // Test with createdOnOrAfter filter
      const responseAfter: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          createdOnOrAfter: '2023-01-02T00:00:00.000Z',
        })}`,
      });

      // Should only include records on or after the specified time
      expect(responseAfter.body.pnl.length).toBeGreaterThanOrEqual(1);
      expect(responseAfter.body.pnl.every(
        (record: any) => new Date(record.createdAt) >= new Date('2023-01-02T00:00:00.000Z'),
      )).toBe(true);

      // Test with height filters
      const responseHeight: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          createdBeforeOrAtHeight: 2000,
          createdOnOrAfterHeight: 1000,
        })}`,
      });

      // Should only include records within the specified height range
      expect(responseHeight.body.pnl.length).toBeGreaterThanOrEqual(1);
      expect(responseHeight.body.pnl.every(
        (record: any) => parseInt(record.createdAtHeight, 10) >= 1000 &&
        parseInt(record.createdAtHeight, 10) <= 2000,
      )).toBe(true);
    });

    it('Get /pnl/parentSubaccountNumber with limit parameter respects the limit', async () => {
      await testMocks.seedData();

      // Create multiple PNL records
      const records = [];
      for (let i = 0; i < 5; i++) {
        records.push({
          ...testConstants.defaultPnl,
          createdAt: new Date(Date.now() - i * 3600000).toISOString(), // 1 hour apart
          createdAtHeight: (1000 + i).toString(),
        });
      }

      await Promise.all(records.map((record) => PnlTable.create(record)));

      // Test with limit parameter
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          limit: 3,
        })}`,
      });

      // Should respect the limit
      expect(response.body.pnl.length).toBeLessThanOrEqual(3);
    });

    it('Get /pnl/parentSubaccountNumber with invalid parentSubaccountNumber returns error', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 128,
        })}`,
        expectedStatus: 400,
      });

      expect(response.body).toEqual({
        errors: [
          {
            location: 'query',
            msg: 'parentSubaccountNumber must be a non-negative integer less than 128',
            param: 'parentSubaccountNumber',
            value: '128',
          },
        ],
      });
    });

    it('Get /pnl/parentSubaccountNumber with non-existent address returns 404', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: 'nonexistentaddress',
          parentSubaccountNumber: 0,
        })}`,
        expectedStatus: 404,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: 'No subaccounts found with address nonexistentaddress and parentSubaccountNumber 0',
          },
        ],
      });
    });

    it('Get /pnl/parentSubaccountNumber with daily=false returns hourly aggregated PNL', async () => {
      await testMocks.seedData();

      // Create hourly PNL records for two different subaccounts
      const baseDate = new Date('2023-01-01T00:00:00.000Z');
      const hourlyRecords = [];

      // Create 5 hourly records for first subaccount
      for (let i = 0; i < 5; i++) {
        const date = new Date(baseDate);
        date.setUTCHours(baseDate.getUTCHours() + i);

        hourlyRecords.push({
          ...testConstants.defaultPnl,
          createdAt: date.toISOString(),
          createdAtHeight: (1000 + i).toString(),
          equity: (1000 + i).toString(),
          totalPnl: (100 + i).toString(),
          netTransfers: (900 + i).toString(),
        });
      }

      // Create 5 hourly records for second subaccount (same hours)
      for (let i = 0; i < 5; i++) {
        const date = new Date(baseDate);
        date.setUTCHours(baseDate.getUTCHours() + i);

        hourlyRecords.push({
          ...testConstants.defaultPnl,
          subaccountId: testConstants.isolatedSubaccountId,
          createdAt: date.toISOString(),
          createdAtHeight: (2000 + i).toString(),
          equity: (2000 + i).toString(),
          totalPnl: (200 + i).toString(),
          netTransfers: (1800 + i).toString(),
        });
      }

      // Insert all records
      await Promise.all(
        hourlyRecords.map((record) => PnlTable.create(record)),
      );

      // Test with daily=false explicitly
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          daily: 'false',
        })}`,
      });

      // Should have 5 records (one for each hour)
      expect(response.body.pnl.length).toEqual(5);

      // Verify hourly timestamps
      const timestamps = response.body.pnl.map((r: any) => new Date(r.createdAt));

      // Check timestamps are in descending order
      for (let i = 0; i < timestamps.length - 1; i++) {
        expect(timestamps[i] > timestamps[i + 1]).toBe(true);
      }

      // Check that records are 1 hour apart
      for (let i = 0; i < timestamps.length - 1; i++) {
        const diffHours = (timestamps[i].getTime() - timestamps[i + 1].getTime()) /
        (1000 * 60 * 60);
        expect(diffHours).toBeCloseTo(1, 1); // Should be close to 1 hour apart
      }

      // Verify the aggregated values for each hour
      // For hour 0 (first hour): 1000 + 2000 = 3000
      expect(Number(response.body.pnl[4].equity)).toEqual(3000); // index 4 since descending
      expect(Number(response.body.pnl[4].totalPnl)).toEqual(300);
      expect(Number(response.body.pnl[4].netTransfers)).toEqual(2700);

      // For hour 4 (last hour): 1004 + 2004 = 3008
      expect(Number(response.body.pnl[0].equity)).toEqual(3008); // index 0 since descending
      expect(Number(response.body.pnl[0].totalPnl)).toEqual(308);
      expect(Number(response.body.pnl[0].netTransfers)).toEqual(2708);

      // Test with omitted daily parameter (should default to false)
      const defaultResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
        })}`,
      });

      // Should also have 5 records when daily is omitted (default behavior)
      expect(defaultResponse.body.pnl.length).toEqual(5);

      // First record should match in both responses (same aggregation logic)
      expect(defaultResponse.body.pnl[0].equity).toEqual(response.body.pnl[0].equity);
      expect(defaultResponse.body.pnl[0].totalPnl).toEqual(response.body.pnl[0].totalPnl);
      expect(defaultResponse.body.pnl[0].netTransfers).toEqual(response.body.pnl[0].netTransfers);
    });
  });
});
