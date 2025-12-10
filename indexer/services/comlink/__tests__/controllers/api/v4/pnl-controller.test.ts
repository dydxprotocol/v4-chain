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

      await Promise.all([
        PnlTable.create({
          ...testConstants.defaultPnl,
          createdAt: '2023-01-01T10:00:00.000Z',
          createdAtHeight: '1000',
          equity: '5000.00',
          totalPnl: '500.00',
        }),

        PnlTable.create({
          ...testConstants.defaultPnl,
          createdAt: '2023-01-01T11:00:00.000Z',
          createdAtHeight: '1100',
          equity: '5100.00',
          totalPnl: '600.00',
        }),
      ]);

      // Test first page with limit 1
      const responsePage1: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          page: 1,
          limit: 1,
        })}`,
      });

      // Verify first page - should have the more recent record (11 AM)
      expect(responsePage1.body.pageSize).toStrictEqual(1);
      expect(responsePage1.body.offset).toStrictEqual(0);
      expect(responsePage1.body.totalResults).toStrictEqual(2);
      expect(responsePage1.body.pnl).toHaveLength(1);
      expect(responsePage1.body.pnl[0].createdAtHeight).toEqual('1100');
      expect(responsePage1.body.pnl[0].createdAt).toEqual('2023-01-01T11:00:00.000Z');

      // Test second page with limit 1
      const responsePage2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl?${getQueryString({
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          page: 2,
          limit: 1,
        })}`,
      });

      // Verify second page - should have the earlier record (10 AM)
      expect(responsePage2.body.pageSize).toStrictEqual(1);
      expect(responsePage2.body.offset).toStrictEqual(1);
      expect(responsePage2.body.totalResults).toStrictEqual(2);
      expect(responsePage2.body.pnl).toHaveLength(1);
      expect(responsePage2.body.pnl[0].createdAtHeight).toEqual('1000');
      expect(responsePage2.body.pnl[0].createdAt).toEqual('2023-01-01T10:00:00.000Z');
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
          createdAtHeight: (1000 + i).toString(),
          equity: (1000 + i).toString(),
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

      // First record should be earliest from day 3 (hour 48)
      expect(dailyResponse.body.pnl[0].createdAtHeight).toBe('1048');
      expect(dailyResponse.body.pnl[0].createdAt).toBe('2023-01-03T00:00:00.000Z');

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
      expect(dailyPageResponse.body.pnl[0].createdAtHeight).toBe('1048');
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
            equity: '3000.00',
            totalPnl: '300.00',
            netTransfers: '2700.00',
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

      // Create 48 hourly records (2 days) for second subaccount
      for (let i = 0; i < 48; i++) {
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
      // First record for day 1 for first subaccount is i=0:
      //    equity=1000, totalPnl=100, netTransfers=900
      // First record for day 1 for second subaccount is i=0:
      //    equity=2000, totalPnl=200, netTransfers=1800
      // Total: equity=3000, totalPnl=300, netTransfers=2700
      expect(day1Record.createdAt).toBe('2023-01-01T00:00:00.000Z');
      expect(Number(day1Record.equity)).toEqual(3000);
      expect(Number(day1Record.totalPnl)).toEqual(300);
      expect(Number(day1Record.netTransfers)).toEqual(2700);

      // Expected values for day 2 (the first hour of day 2)
      // First record for day 2 for first subaccount is i=24:
      //    equity=1024, totalPnl=124, netTransfers=924
      // First record for day 2 for second subaccount is i=24:
      //    equity=2024, totalPnl=224, netTransfers=1824
      // Total: equity=3048, totalPnl=348, netTransfers=2748
      expect(day2Record.createdAt).toBe('2023-01-02T00:00:00.000Z');
      expect(Number(day2Record.equity)).toEqual(3048);
      expect(Number(day2Record.totalPnl)).toEqual(348);
      expect(Number(day2Record.netTransfers)).toEqual(2748);

      // Verify the records are in descending order by date (day 2 should come before day 1)
      const timestamps = response.body.pnl.map(
        (record: any) => new Date(record.createdAt).getTime());
      expect(timestamps[0]).toBeGreaterThan(timestamps[1]);
    });

    it('Get /pnl/parentSubaccountNumber with daily=true correctly handles child subaccounts created mid-day', async () => {
      await testMocks.seedData();

      const records = [];

      // Day 1 (Jan 1): Only parent subaccount (0) exists from 00:00
      records.push({
        ...testConstants.defaultPnl,
        subaccountId: testConstants.defaultSubaccountId,
        createdAt: '2023-01-01T00:00:00.000Z',
        createdAtHeight: '1000',
        equity: '1000',
        totalPnl: '100',
        netTransfers: '500',
      });

      // Day 2 (Jan 2): Parent subaccount has records starting from 00:00
      records.push({
        ...testConstants.defaultPnl,
        subaccountId: testConstants.defaultSubaccountId,
        createdAt: '2023-01-02T00:00:00.000Z',
        createdAtHeight: '2000',
        equity: '2000',
        totalPnl: '200',
        netTransfers: '600',
      });

      records.push({
        ...testConstants.defaultPnl,
        subaccountId: testConstants.defaultSubaccountId,
        createdAt: '2023-01-02T06:00:00.000Z',
        createdAtHeight: '2060',
        equity: '2100',
        totalPnl: '210',
        netTransfers: '610',
      });

      // Day 2 (Jan 2): Child subaccount (128) was created at 05:00
      // This subaccount has NO 00:00 record because it didn't exist yet
      records.push({
        ...testConstants.defaultPnl,
        subaccountId: testConstants.isolatedSubaccountId,
        createdAt: '2023-01-02T05:00:00.000Z',
        createdAtHeight: '2050',
        equity: '500',
        totalPnl: '50',
        netTransfers: '100',
      });

      records.push({
        ...testConstants.defaultPnl,
        subaccountId: testConstants.isolatedSubaccountId,
        createdAt: '2023-01-02T06:00:00.000Z',
        createdAtHeight: '2061',
        equity: '550',
        totalPnl: '55',
        netTransfers: '105',
      });

      records.push({
        ...testConstants.defaultPnl,
        subaccountId: testConstants.isolatedSubaccountId,
        createdAt: '2023-01-02T12:00:00.000Z',
        createdAtHeight: '2120',
        equity: '600',
        totalPnl: '60',
        netTransfers: '110',
      });

      // Day 3 (Jan 3): Both parent and child subaccounts have records at 00:00
      records.push({
        ...testConstants.defaultPnl,
        subaccountId: testConstants.defaultSubaccountId,
        createdAt: '2023-01-03T00:00:00.000Z',
        createdAtHeight: '3000',
        equity: '3000',
        totalPnl: '300',
        netTransfers: '700',
      });

      records.push({
        ...testConstants.defaultPnl,
        subaccountId: testConstants.isolatedSubaccountId,
        createdAt: '2023-01-03T00:00:00.000Z',
        createdAtHeight: '3001',
        equity: '800',
        totalPnl: '80',
        netTransfers: '150',
      });

      records.push({
        ...testConstants.defaultPnl,
        subaccountId: testConstants.defaultSubaccountId,
        createdAt: '2023-01-03T06:00:00.000Z',
        createdAtHeight: '3060',
        equity: '3100',
        totalPnl: '310',
        netTransfers: '710',
      });

      records.push({
        ...testConstants.defaultPnl,
        subaccountId: testConstants.isolatedSubaccountId,
        createdAt: '2023-01-03T06:00:00.000Z',
        createdAtHeight: '3061',
        equity: '850',
        totalPnl: '85',
        netTransfers: '155',
      });

      // Insert all records
      await Promise.all(records.map((record) => PnlTable.create(record)));

      // Test with daily=true
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          daily: 'true',
        })}`,
      });

      // Should have 3 records (one for each day)
      expect(response.body.pnl.length).toEqual(3);

      // Find records by day
      const day1Record = response.body.pnl.find(
        (r: any) => r.createdAt.startsWith('2023-01-01'),
      );
      const day2Record = response.body.pnl.find(
        (r: any) => r.createdAt.startsWith('2023-01-02'),
      );
      const day3Record = response.body.pnl.find(
        (r: any) => r.createdAt.startsWith('2023-01-03'),
      );

      // Verify all day records exist
      expect(day1Record).toBeDefined();
      expect(day2Record).toBeDefined();
      expect(day3Record).toBeDefined();

      // Day 1: Only parent subaccount existed
      expect(day1Record.createdAt).toBe('2023-01-01T00:00:00.000Z');
      expect(Number(day1Record.equity)).toEqual(1000);
      expect(Number(day1Record.totalPnl)).toEqual(100);
      expect(Number(day1Record.netTransfers)).toEqual(500);

      // Day 2: Child created at 05:00, should be EXCLUDED from daily aggregate
      // Critical: Should only include parent's 00:00 record
      // NOT: 2000 + 500 = 2500 (mixing timestamps - this was the BUG)
      expect(day2Record.createdAt).toBe('2023-01-02T00:00:00.000Z');
      expect(Number(day2Record.equity)).toEqual(2000);
      expect(Number(day2Record.totalPnl)).toEqual(200);
      expect(Number(day2Record.netTransfers)).toEqual(600);

      // Day 3: Both parent and child have 00:00 records - should aggregate both
      expect(day3Record.createdAt).toBe('2023-01-03T00:00:00.000Z');
      expect(Number(day3Record.equity)).toEqual(3800); // 3000 + 800
      expect(Number(day3Record.totalPnl)).toEqual(380); // 300 + 80
      expect(Number(day3Record.netTransfers)).toEqual(850); // 700 + 150
    });

    it('Get /pnl/parentSubaccountNumber daily values match hourly values at 00:00 timestamps', async () => {
      await testMocks.seedData();

      const records = [];

      // Create 3 days of data with multiple hourly records
      for (let day = 1; day <= 3; day++) {
        const hours = [0, 3, 6, 12, 18];

        for (const hour of hours) {
          const date = new Date(`2023-01-0${day}T${hour.toString().padStart(2, '0')}:00:00.000Z`);

          // Parent subaccount (0)
          records.push({
            ...testConstants.defaultPnl,
            subaccountId: testConstants.defaultSubaccountId,
            createdAt: date.toISOString(),
            createdAtHeight: (day * 10000 + hour * 100).toString(),
            equity: (day * 1000 + hour * 10).toString(),
            totalPnl: (day * 100 + hour).toString(),
            netTransfers: (day * 200 + hour * 5).toString(),
          });

          // Child subaccount (128)
          records.push({
            ...testConstants.defaultPnl,
            subaccountId: testConstants.isolatedSubaccountId,
            createdAt: date.toISOString(),
            createdAtHeight: (day * 10000 + hour * 100 + 1).toString(),
            equity: (day * 1000 + hour * 10 + 5).toString(),
            totalPnl: (day * 100 + hour + 2).toString(),
            netTransfers: (day * 200 + hour * 5 + 3).toString(),
          });
        }
      }

      // Insert all records
      await Promise.all(records.map((record) => PnlTable.create(record)));

      // Get daily aggregated records
      const dailyResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          daily: 'true',
        })}`,
      });

      // Get hourly aggregated records
      const hourlyResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
        })}`,
      });

      // Should have 3 daily records
      expect(dailyResponse.body.pnl.length).toEqual(3);

      // For each day, verify daily record matches hourly record at 00:00
      for (let day = 1; day <= 3; day++) {
        const expectedTimestamp = `2023-01-0${day}T00:00:00.000Z`;

        const dailyRecord = dailyResponse.body.pnl.find(
          (r: any) => r.createdAt === expectedTimestamp,
        );

        const hourlyRecord = hourlyResponse.body.pnl.find(
          (r: any) => r.createdAt === expectedTimestamp,
        );

        // Both should exist
        expect(dailyRecord).toBeDefined();
        expect(hourlyRecord).toBeDefined();

        // Critical: Daily and hourly values should be IDENTICAL at 00:00
        expect(dailyRecord.equity).toBe(hourlyRecord.equity);
        expect(dailyRecord.totalPnl).toBe(hourlyRecord.totalPnl);
        expect(dailyRecord.netTransfers).toBe(hourlyRecord.netTransfers);
        expect(dailyRecord.createdAtHeight).toBe(hourlyRecord.createdAtHeight);

        // Verify the expected aggregated values for 00:00 (hour = 0)
        // Parent: equity = day * 1000 + 0 * 10 = day * 1000
        // Child:  equity = day * 1000 + 0 * 10 + 5 = day * 1000 + 5
        // Total:  equity = day * 2000 + 5
        const expectedEquity = (day * 2000 + 5).toString();
        expect(dailyRecord.equity).toBe(expectedEquity);
        expect(hourlyRecord.equity).toBe(expectedEquity);
      }
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
            msg: 'No PnL data found for address nonexistentaddress and parentSubaccountNumber 0',
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

      // Check specific timestamps for a few records
      // Since records are in descending order, latest hour (4) is at index 0
      expect(response.body.pnl[0].createdAt).toEqual('2023-01-01T04:00:00.000Z'); // Hour 4
      expect(response.body.pnl[2].createdAt).toEqual('2023-01-01T02:00:00.000Z'); // Hour 2
      expect(response.body.pnl[4].createdAt).toEqual('2023-01-01T00:00:00.000Z'); // Hour 0

      // Verify the aggregated values for each hour
      // For hour 0 (first hour): 1000 + 2000 = 3000
      expect(Number(response.body.pnl[4].equity)).toEqual(3000); // index 4 since descending
      expect(Number(response.body.pnl[4].totalPnl)).toEqual(300);
      expect(Number(response.body.pnl[4].netTransfers)).toEqual(2700);

      // For hour 2 (middle hour): 1002 + 2002 = 3004
      expect(Number(response.body.pnl[2].equity)).toEqual(3004); // index 2 is hour 2
      expect(Number(response.body.pnl[2].totalPnl)).toEqual(304);
      expect(Number(response.body.pnl[2].netTransfers)).toEqual(2704);

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

      // Check that the same timestamps are returned with the default parameter
      expect(defaultResponse.body.pnl[0].createdAt).toEqual('2023-01-01T04:00:00.000Z');
      expect(defaultResponse.body.pnl[4].createdAt).toEqual('2023-01-01T00:00:00.000Z');

      // First record should match in both responses (same aggregation logic)
      expect(defaultResponse.body.pnl[0].equity).toEqual(response.body.pnl[0].equity);
      expect(defaultResponse.body.pnl[0].totalPnl).toEqual(response.body.pnl[0].totalPnl);
      expect(defaultResponse.body.pnl[0].netTransfers).toEqual(response.body.pnl[0].netTransfers);
    });

    it('Get /pnl/parentSubaccountNumber performs efficiently with large datasets', async () => {
      await testMocks.seedData();

      // Create a large, realistic dataset: 3 years * 365 days * 24 hours
      //   = 26,280 records per subaccount
      // With 3 child subaccounts = 78,840 total records
      const baseDate = new Date('2021-01-01T00:00:00.000Z');
      const yearsToCreate = 3;
      const hoursPerYear = 365 * 24;
      const totalHours = yearsToCreate * hoursPerYear;
      const apiLimit = 1000; // API has a default limit of 1000 records

      // All three are child subaccounts of parent subaccount 0
      const subaccountIds = [
        testConstants.defaultSubaccountId, // subaccount 0
        testConstants.isolatedSubaccountId, // subaccount 128
        testConstants.isolatedSubaccountId2, // subaccount 256
      ];

      const largeDataset = [];

      for (let hour = 0; hour < totalHours; hour++) {
        const date = new Date(baseDate);
        date.setUTCHours(baseDate.getUTCHours() + hour);

        for (let i = 0; i < subaccountIds.length; i++) {
          largeDataset.push({
            ...testConstants.defaultPnl,
            subaccountId: subaccountIds[i],
            createdAt: date.toISOString(),
            createdAtHeight: (10000 + hour + (i * 100000)).toString(),
            equity: (1000 + hour + (i * 100)).toString(),
            totalPnl: (100 + hour + (i * 10)).toString(),
            netTransfers: (900 + hour + (i * 90)).toString(),
          });
        }
      }

      // Batch insert for better performance
      // Keep per-batch concurrency modest to avoid exhausting the pool on CI
      const batchSize = 200;
      for (let i = 0; i < largeDataset.length; i += batchSize) {
        const batch = largeDataset.slice(i, i + batchSize);
        // Insert in sub-batches to cap peak concurrency (~50)
        for (let j = 0; j < batch.length; j += 50) {
          const subBatch = batch.slice(j, j + 50);
          await Promise.all(subBatch.map((record) => PnlTable.create(record)));
        }
      }

      // Test 1: Hourly aggregation performance (most expensive operation)
      // Returns up to 1000 hourly records due to API limit
      const startHourly = Date.now();
      const hourlyResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          daily: 'false',
          limit: apiLimit,
        })}`,
      });
      const hourlyTime = Date.now() - startHourly;

      // Should complete within 3 seconds for hourly aggregation even with large dataset
      expect(hourlyTime).toBeLessThan(3000);
      expect(hourlyResponse.body.pnl.length).toEqual(apiLimit);

      // Verify we're getting the most recent records (descending order)
      const firstHourlyRecord = hourlyResponse.body.pnl[0];
      const lastHourlyRecord = hourlyResponse.body.pnl[hourlyResponse.body.pnl.length - 1];
      expect(new Date(firstHourlyRecord.createdAt).getTime()).toBeGreaterThan(
        new Date(lastHourlyRecord.createdAt).getTime(),
      );

      // Test 2: Daily aggregation performance
      // Returns up to 1000 daily records due to API limit
      const startDaily = Date.now();
      const dailyResponse: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/pnl/parentSubaccountNumber?${getQueryString({
          address: testConstants.defaultAddress,
          parentSubaccountNumber: 0,
          daily: 'true',
          limit: apiLimit,
        })}`,
      });
      const dailyTime = Date.now() - startDaily;

      // Daily aggregation should complete within 3 seconds
      const expectedDays = yearsToCreate * 365;
      expect(dailyTime).toBeLessThan(3000);
      // Should return 1000 days (API limit) even though we have more
      expect(dailyResponse.body.pnl.length).toEqual(Math.min(apiLimit, expectedDays));

      // Verify we're getting the most recent records (descending order)
      const firstDailyRecord = dailyResponse.body.pnl[0];
      const lastDailyRecord = dailyResponse.body.pnl[dailyResponse.body.pnl.length - 1];
      expect(new Date(firstDailyRecord.createdAt).getTime()).toBeGreaterThan(
        new Date(lastDailyRecord.createdAt).getTime(),
      );

      // Verify data integrity - check that aggregation is correct for the oldest record returned
      // Since we have 1095 days but only get 1000, the oldest returned is day 95 (1095 - 1000)
      // This would be 2021-01-01 + 95 days = 2021-04-06
      const daysToSkip = expectedDays - apiLimit;
      const oldestReturnedDate = new Date(baseDate.getTime() + daysToSkip * 24 * 60 * 60 * 1000);

      expect(lastDailyRecord.createdAt).toBe(oldestReturnedDate.toISOString());

      // Verify aggregation is correct for this day
      // Hour offset for this day: (expectedDays - apiLimit) * 24
      const hourOffset = daysToSkip * 24;
      // S0: equity=1000+hourOffset, totalPnl=100+hourOffset, netTransfers=900+hourOffset
      // S128: equity=1100+hourOffset, totalPnl=110+hourOffset, netTransfers=990+hourOffset
      // S256: equity=1200+hourOffset, totalPnl=120+hourOffset, netTransfers=1080+hourOffset
      const expectedEquity = 3300 + (hourOffset * 3);
      const expectedTotalPnl = 330 + (hourOffset * 3);
      const expectedNetTransfers = 2970 + (hourOffset * 3);

      expect(Number(lastDailyRecord.equity)).toEqual(expectedEquity);
      expect(Number(lastDailyRecord.totalPnl)).toEqual(expectedTotalPnl);
      expect(Number(lastDailyRecord.netTransfers)).toEqual(expectedNetTransfers);
    }, 120000);
  });
});
