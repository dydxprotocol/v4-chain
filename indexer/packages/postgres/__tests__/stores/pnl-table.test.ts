import {
  Ordering,
  PnlColumns,
} from '../../src/types';
import * as PnlTable from '../../src/stores/pnl-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import {
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultPnl,
  defaultPnl2,
  defaultBlock,
} from '../helpers/constants';
import { BlockTable } from 'packages/postgres/src';

describe('Pnl store', () => {
  beforeEach(async () => {
    await seedData();
  });

  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a Pnl record', async () => {
    await PnlTable.create(defaultPnl);
    const fetched = await PnlTable.findById(defaultPnl.subaccountId, defaultPnl.createdAt);
    expect(fetched).toEqual(expect.objectContaining(defaultPnl));
  });

  it('Successfully creates multiple Pnl records', async () => {
    await BlockTable.create({
      ...defaultBlock,
      blockHeight: '5',
    });
    await Promise.all([
      PnlTable.create(defaultPnl),
      PnlTable.create(defaultPnl2),
    ]);

    const { results: pnls } = await PnlTable.findAll({}, [], {
      orderBy: [[PnlColumns.createdAtHeight, Ordering.ASC]],
    });

    expect(pnls.length).toEqual(2);
    expect(pnls[0]).toEqual(expect.objectContaining(defaultPnl));
    expect(pnls[1]).toEqual(expect.objectContaining(defaultPnl2));
  });

  it('Successfully finds Pnl records with subaccountId', async () => {
    await Promise.all([
      PnlTable.create(defaultPnl),
      PnlTable.create({
        ...defaultPnl,
        createdAt: '2022-06-01T01:00:00.000Z',
      }),
      PnlTable.create({
        ...defaultPnl,
        subaccountId: defaultSubaccountId2,
        createdAt: '2022-06-01T00:00:00.000Z',
      }),
    ]);

    const { results: pnls } = await PnlTable.findAll(
      {
        subaccountId: [defaultSubaccountId],
      },
      [],
      {},
    );

    expect(pnls.length).toEqual(2);
  });

  it('Successfully finds Pnl records using pagination', async () => {
    await Promise.all([
      PnlTable.create(defaultPnl),
      PnlTable.create({
        ...defaultPnl,
        createdAt: '2020-01-01T00:00:00.000Z',
        createdAtHeight: '1000',
      }),
    ]);

    const responsePageOne = await PnlTable.findAll({
      page: 1,
      limit: 1,
    }, [], {
      orderBy: [[PnlColumns.createdAtHeight, Ordering.DESC]],
    });

    expect(responsePageOne.results.length).toEqual(1);
    expect(responsePageOne.results[0]).toEqual(expect.objectContaining({
      ...defaultPnl,
      createdAt: '2020-01-01T00:00:00.000Z',
      createdAtHeight: '1000',
    }));
    expect(responsePageOne.offset).toEqual(0);
    expect(responsePageOne.total).toEqual(2);

    const responsePageTwo = await PnlTable.findAll({
      page: 2,
      limit: 1,
    }, [], {
      orderBy: [[PnlColumns.createdAtHeight, Ordering.DESC]],
    });

    expect(responsePageTwo.results.length).toEqual(1);
    expect(responsePageTwo.results[0]).toEqual(expect.objectContaining(defaultPnl));
    expect(responsePageTwo.offset).toEqual(1);
    expect(responsePageTwo.total).toEqual(2);

    const responsePageAllPages = await PnlTable.findAll({
      page: 1,
      limit: 2,
    }, [], {
      orderBy: [[PnlColumns.createdAtHeight, Ordering.DESC]],
    });

    expect(responsePageAllPages.results.length).toEqual(2);
    expect(responsePageAllPages.results[0]).toEqual(expect.objectContaining({
      ...defaultPnl,
      createdAt: '2020-01-01T00:00:00.000Z',
      createdAtHeight: '1000',
    }));
    expect(responsePageAllPages.results[1]).toEqual(expect.objectContaining(defaultPnl));
    expect(responsePageAllPages.offset).toEqual(0);
    expect(responsePageAllPages.total).toEqual(2);
  });

  it('Successfully retrieves hourly aggregated PNL records for single subaccount', async () => {
    const records = [];

    // Day 1: Create records at 00:00, 06:00, 12:00, 18:00
    const day1Date = new Date('2023-01-01T00:00:00.000Z');
    for (let i = 0; i < 24; i += 6) {
      const date = new Date(day1Date);
      date.setUTCHours(i);
      records.push({
        ...defaultPnl,
        createdAt: date.toISOString(),
        createdAtHeight: (1000 + i).toString(),
        equity: (1000 + i).toString(),
      });
    }

    // Day 2: Create records at 00:00, 06:00, 12:00, 18:00
    const day2Date = new Date('2023-01-02T00:00:00.000Z');
    for (let i = 0; i < 24; i += 6) {
      const date = new Date(day2Date);
      date.setUTCHours(i);
      records.push({
        ...defaultPnl,
        createdAt: date.toISOString(),
        createdAtHeight: (2000 + i).toString(),
        equity: (2000 + i).toString(),
      });
    }

    // Insert all records
    await Promise.all(records.map((record) => PnlTable.create(record)));

    // Get hourly aggregated records
    const hourlyResults = await PnlTable.findAllHourlyAggregate(
      { subaccountId: [defaultSubaccountId] },
      [],
      {},
    );

    // We should get exactly 8 records (one for each hour we created)
    expect(hourlyResults.results.length).toBe(8);

    // Check a few specific records
    const day1FirstHour = hourlyResults.results.find((r) => r.createdAtHeight === '1000');
    expect(day1FirstHour).toBeDefined();
    expect(day1FirstHour?.equity).toBe('1000');

    const day2LastHour = hourlyResults.results.find((r) => r.createdAtHeight === '2018');
    expect(day2LastHour).toBeDefined();
    expect(day2LastHour?.equity).toBe('2018');
  });

  it('Successfully aggregates hourly PNL records across multiple subaccounts', async () => {
    const records = [];
    const subaccountIds = [defaultSubaccountId, defaultSubaccountId2];

    // Create records for specific hours, one for each subaccount
    const hours = [
      '2023-01-01T00:00:00.000Z',
      '2023-01-01T06:00:00.000Z',
      '2023-01-01T12:00:00.000Z',
      '2023-01-01T18:00:00.000Z',
    ];

    for (let hourIndex = 0; hourIndex < hours.length; hourIndex++) {
      for (let i = 0; i < subaccountIds.length; i++) {
        records.push({
          ...defaultPnl,
          subaccountId: subaccountIds[i],
          createdAt: hours[hourIndex],
          createdAtHeight: (1000 + hourIndex * 100 + i).toString(),
          equity: (1000 + (i * 100)).toString(),
          totalPnl: (100 + (i * 10)).toString(),
          netTransfers: (500 + (i * 50)).toString(),
        });
      }
    }

    // Insert all records
    await Promise.all(records.map((record) => PnlTable.create(record)));

    // Get aggregated hourly records for all subaccounts
    const hourlyResults = await PnlTable.findAllHourlyAggregate(
      { subaccountId: subaccountIds },
      [],
      {},
    );

    // We should get exactly 4 records (one for each hour, aggregated across subaccounts)
    expect(hourlyResults.results.length).toBe(4);

    // Sort results by time for consistent checking
    const sortedResults = [...hourlyResults.results].sort((a, b) => new Date(a.createdAt)
      .getTime() - new Date(b.createdAt).getTime());

    // Check the actual timestamp values - expecting exact UTC times
    expect(sortedResults[0].createdAt).toBe('2023-01-01T00:00:00.000Z');
    expect(sortedResults[1].createdAt).toBe('2023-01-01T06:00:00.000Z');
    expect(sortedResults[2].createdAt).toBe('2023-01-01T12:00:00.000Z');
    expect(sortedResults[3].createdAt).toBe('2023-01-01T18:00:00.000Z');

    // Verify aggregation for each hour
    for (const result of sortedResults) {
    // For each hour, the values should be aggregated across both subaccounts
    // Equity: 1000 + 1100 = 2100
    // TotalPnl: 100 + 110 = 210
    // NetTransfers: 500 + 550 = 1050
      expect(Number(result.equity)).toBeCloseTo(2100, 0);
      expect(Number(result.totalPnl)).toBeCloseTo(210, 0);
      expect(Number(result.netTransfers)).toBeCloseTo(1050, 0);
    }
  });

  it('Successfully paginates hourly aggregated PNL records', async () => {
    const records = [];
    const subaccountIds = [defaultSubaccountId, defaultSubaccountId2];

    // Create records for 8 hours (across 2 days)
    for (let hour = 0; hour < 8; hour++) {
      const day = Math.floor(hour / 4) + 1; // Day 1 for hours 0-3, Day 2 for hours 4-7
      const hourOfDay = (hour % 4) * 6; // Hours 0, 6, 12, 18 of each day

      const date = new Date(`2023-01-0${day}T${hourOfDay.toString().padStart(2, '0')}:00:00.000Z`);

      // For each hour, create records for both subaccounts
      for (let i = 0; i < subaccountIds.length; i++) {
        records.push({
          ...defaultPnl,
          subaccountId: subaccountIds[i],
          createdAt: date.toISOString(),
          createdAtHeight: (day * 1000 + hourOfDay * 10 + i).toString(),
          equity: (day * 1000 + hourOfDay * 10 + (i * 100)).toString(),
        });
      }
    }

    // Insert all records
    await Promise.all(records.map((record) => PnlTable.create(record)));

    // Test with pagination - first page (3 records)
    const hourlyPage1 = await PnlTable.findAllHourlyAggregate(
      {
        subaccountId: subaccountIds,
        page: 1,
        limit: 3,
      },
      [],
      {},
    );

    // Basic checks
    expect(hourlyPage1.results.length).toBe(3);
    expect(hourlyPage1.limit).toBe(3);
    expect(hourlyPage1.offset).toBe(0);

    // Test with pagination - second page (3 records)
    const hourlyPage2 = await PnlTable.findAllHourlyAggregate(
      {
        subaccountId: subaccountIds,
        page: 2,
        limit: 3,
      },
      [],
      {},
    );

    expect(hourlyPage2.results.length).toBe(3);
    expect(hourlyPage2.offset).toBe(3);

    // Test with pagination - third page (remaining records)
    const hourlyPage3 = await PnlTable.findAllHourlyAggregate(
      {
        subaccountId: subaccountIds,
        page: 3,
        limit: 3,
      },
      [],
      {},
    );

    expect(hourlyPage3.results.length).toBe(2);
    expect(hourlyPage3.offset).toBe(6);
  });

  it('Successfully filters hourly PNL records by time range', async () => {
    const records = [];
    const subaccountIds = [defaultSubaccountId, defaultSubaccountId2];

    // Create records for 12 hours (covering 2 days)
    for (let hour = 0; hour < 12; hour++) {
      const day = Math.floor(hour / 6) + 1; // Day 1 for hours 0-5, Day 2 for hours 6-11
      // Fix the mixed operators error by using parentheses
      const hourOfDay = (hour % 6) * 4; // Hours 0, 4, 8, 12, 16, 20 of each day

      const date = new Date(`2023-01-0${day}T${hourOfDay.toString().padStart(2, '0')}:00:00.000Z`);

      // For each hour, create records for both subaccounts
      for (let i = 0; i < subaccountIds.length; i++) {
        records.push({
          ...defaultPnl,
          subaccountId: subaccountIds[i],
          createdAt: date.toISOString(),
          createdAtHeight: (hour * 100 + i).toString(), // Unique height for each record
          equity: (hour * 1000 + (i * 100)).toString(),
        });
      }
    }

    // Insert all records
    await Promise.all(records.map((record) => PnlTable.create(record)));

    // Test with time range filter - createdOnOrAfter
    const startDate = new Date('2023-01-01T12:00:00.000Z'); // Start from hour 3 (noon on day 1)
    const hourlyWithStartDate = await PnlTable.findAllHourlyAggregate(
      {
        subaccountId: subaccountIds,
        createdOnOrAfter: startDate.toISOString(),
      },
      [],
      { orderBy: [[PnlColumns.createdAt, Ordering.ASC]] },
    );

    // Should include hours 3-11 (9 hours total)
    expect(hourlyWithStartDate.results.length).toBe(9);

    // Test with height filters
    const hourlyWithHeightRange = await PnlTable.findAllHourlyAggregate(
      {
        subaccountId: subaccountIds,
        createdOnOrAfterHeight: '300', // Start from hour 3
        createdBeforeOrAtHeight: '900', // End at hour 9
      },
      [],
      { orderBy: [[PnlColumns.createdAtHeight, Ordering.ASC]] },
    );

    // Should include hours 3-9 (7 hours total)
    expect(hourlyWithHeightRange.results.length).toBe(7);

    // Check the height range
    const heights = hourlyWithHeightRange.results.map((r) => Number(r.createdAtHeight));
    const minHeight = Math.min(...heights);
    const maxHeight = Math.max(...heights);
    expect(minHeight).toBeGreaterThanOrEqual(300);
    expect(maxHeight).toBeLessThanOrEqual(900);
  });

  it('Successfully retrieves daily PNL records with first of each day for single subaccount', async () => {
    const records = [];
    // Day 1 (Jan 1): Create records at 00:00, 06:00, 12:00, 18:00
    const day1Date = new Date('2023-01-01T00:00:00.000Z');
    for (let i = 0; i < 24; i += 6) {
      const date = new Date(day1Date);
      date.setUTCHours(i);
      records.push({
        ...defaultPnl,
        createdAt: date.toISOString(),
        createdAtHeight: (1000 + i).toString(),
        equity: (1000 + i).toString(),
      });
    }

    // Day 2 (Jan 2): Create records at 00:00, 06:00, 12:00, 18:00
    const day2Date = new Date('2023-01-02T00:00:00.000Z');
    for (let i = 0; i < 24; i += 6) {
      const date = new Date(day2Date);
      date.setUTCHours(i);
      records.push({
        ...defaultPnl,
        createdAt: date.toISOString(),
        createdAtHeight: (2000 + i).toString(),
        equity: (2000 + i).toString(),
      });
    }

    // Day 3 (Jan 3): Create records at 00:00, 06:00, 12:00, 18:00
    const day3Date = new Date('2023-01-03T00:00:00.000Z');
    for (let i = 0; i < 24; i += 6) {
      const date = new Date(day3Date);
      date.setUTCHours(i);
      records.push({
        ...defaultPnl,
        createdAt: date.toISOString(),
        createdAtHeight: (3000 + i).toString(),
        equity: (3000 + i).toString(),
      });
    }

    // Insert all records
    await Promise.all(records.map((record) => PnlTable.create(record)));

    // Get daily records
    const dailyResults = await PnlTable.findAllDailyAggregate(
      { subaccountId: [defaultSubaccountId] },
      [],
      {},
    );

    // We should get exactly 3 records (one for each day)
    expect(dailyResults.results.length).toBe(3);

    // Check the record heights - they should correspond to the first record of each day
    expect(dailyResults.results[0]).toEqual(expect.objectContaining({
      createdAtHeight: '3000',  // The first record of day 3
    }));

    expect(dailyResults.results[1]).toEqual(expect.objectContaining({
      createdAtHeight: '2000',  // The first record of day 2
    }));

    expect(dailyResults.results[2]).toEqual(expect.objectContaining({
      createdAtHeight: '1000',  // The first record of day 1
    }));

    // Check the actual timestamp values
    expect(dailyResults.results[0].createdAt).toBe('2023-01-03T00:00:00.000Z');
    expect(dailyResults.results[1].createdAt).toBe('2023-01-02T00:00:00.000Z');
    expect(dailyResults.results[2].createdAt).toBe('2023-01-01T00:00:00.000Z');
  });

  it('Successfully aggregates daily PNL records across multiple subaccounts', async () => {
    const records = [];
    const subaccountIds = [defaultSubaccountId, defaultSubaccountId2];

    // Day 1 (Jan 1): Create earliest record (00:00) for each subaccount
    for (let i = 0; i < subaccountIds.length; i++) {
      records.push({
        ...defaultPnl,
        subaccountId: subaccountIds[i],
        createdAt: '2023-01-01T00:00:00.000Z',
        createdAtHeight: (1000 + i).toString(),
        equity: (1000 + (i * 100)).toString(),
        totalPnl: (100 + (i * 10)).toString(),
        netTransfers: (500 + (i * 50)).toString(),
      });
    }

    // Add a later record on Day 1 to verify we take earliest
    records.push({
      ...defaultPnl,
      subaccountId: defaultSubaccountId,
      createdAt: '2023-01-01T12:00:00.000Z',
      createdAtHeight: '1500',
      equity: '9999', // Should be ignored for daily aggregation
      totalPnl: '999',
      netTransfers: '999',
    });

    // Day 2 (Jan 2): Create earliest record (00:00) for each subaccount
    for (let i = 0; i < subaccountIds.length; i++) {
      records.push({
        ...defaultPnl,
        subaccountId: subaccountIds[i],
        createdAt: '2023-01-02T00:00:00.000Z',
        createdAtHeight: (2000 + i).toString(),
        equity: (2000 + (i * 100)).toString(),
        totalPnl: (200 + (i * 10)).toString(),
        netTransfers: (600 + (i * 50)).toString(),
      });
    }

    // Insert all records
    await Promise.all(records.map((record) => PnlTable.create(record)));

    // Get aggregated daily records for all subaccounts
    const dailyResults = await PnlTable.findAllDailyAggregate(
      { subaccountId: subaccountIds },
      [],
      {},
    );

    // We should get exactly 2 records (one for each day, aggregated across subaccounts)
    expect(dailyResults.results.length).toBe(2);

    // Check the specific timestamps
    const day1Record = dailyResults.results.find((r) => r.createdAt.includes('2023-01-01'));
    const day2Record = dailyResults.results.find((r) => r.createdAt.includes('2023-01-02'));

    // Verify both day records exist
    expect(day1Record).toBeDefined();
    expect(day2Record).toBeDefined();
    expect(day1Record?.createdAt).toBe('2023-01-01T00:00:00.000Z');
    expect(day2Record?.createdAt).toBe('2023-01-02T00:00:00.000Z');

    // Day 2 aggregated values
    // Equity sum: 2000 + 2100 = 4100
    // TotalPnl sum: 200 + 210 = 410
    // NetTransfers sum: 600 + 650 = 1250
    expect(Number(day2Record?.equity)).toBeCloseTo(4100, 0);
    expect(Number(day2Record?.totalPnl)).toBeCloseTo(410, 0);
    expect(Number(day2Record?.netTransfers)).toBeCloseTo(1250, 0);

    // Day 1 aggregated values
    // Equity sum: 1000 + 1100 = 2100
    // TotalPnl sum: 100 + 110 = 210
    // NetTransfers sum: 500 + 550 = 1050
    expect(Number(day1Record?.equity)).toBeCloseTo(2100, 0);
    expect(Number(day1Record?.totalPnl)).toBeCloseTo(210, 0);
    expect(Number(day1Record?.netTransfers)).toBeCloseTo(1050, 0);
  });

  it('Successfully paginates daily aggregated PNL records', async () => {
    const records = [];
    const subaccountIds = [defaultSubaccountId, defaultSubaccountId2];

    // Create records for 5 days
    for (let day = 1; day <= 5; day++) {
      const date = new Date(`2023-01-0${day}T00:00:00.000Z`);

      // For each day, create records for both subaccounts
      for (let i = 0; i < subaccountIds.length; i++) {
        records.push({
          ...defaultPnl,
          subaccountId: subaccountIds[i],
          createdAt: date.toISOString(),
          createdAtHeight: (day * 1000 + i).toString(),
          equity: (day * 1000 + (i * 100)).toString(),
          totalPnl: (day * 100 + (i * 10)).toString(),
          netTransfers: (day * 200 + (i * 50)).toString(),
        });
      }
    }

    // Insert all records
    await Promise.all(records.map((record) => PnlTable.create(record)));

    // Test with pagination - first page (2 records)
    const dailyPage1 = await PnlTable.findAllDailyAggregate(
      {
        subaccountId: subaccountIds,
        page: 1,
        limit: 2,
      },
      [],
      {},
    );

    // Basic checks
    expect(dailyPage1.results.length).toBe(2);
    expect(dailyPage1.limit).toBe(2);
    expect(dailyPage1.offset).toBe(0);

    // Test with pagination - second page (2 records)
    const dailyPage2 = await PnlTable.findAllDailyAggregate(
      {
        subaccountId: subaccountIds,
        page: 2,
        limit: 2,
      },
      [],
      {},
    );

    expect(dailyPage2.results.length).toBe(2);
    expect(dailyPage2.offset).toBe(2);

    // Test with pagination - third page (should have the oldest data)
    const dailyPage3 = await PnlTable.findAllDailyAggregate(
      {
        subaccountId: subaccountIds,
        page: 3,
        limit: 2,
      },
      [],
      {},
    );

    expect(dailyPage3.results.length).toBeGreaterThan(0);  // At least some results
    expect(dailyPage3.offset).toBe(4);
  });

  it('Successfully applies time range filters to daily aggregated PNL', async () => {
    const records = [];
    const subaccountIds = [defaultSubaccountId, defaultSubaccountId2];

    // Create records for 5 days
    for (let day = 1; day <= 5; day++) {
      const date = new Date(`2023-01-0${day}T00:00:00.000Z`);

      // For each day, create records for both subaccounts
      for (let i = 0; i < subaccountIds.length; i++) {
        records.push({
          ...defaultPnl,
          subaccountId: subaccountIds[i],
          createdAt: date.toISOString(),
          createdAtHeight: (day * 1000 + i).toString(),
          equity: (day * 1000 + (i * 100)).toString(),
        });
      }
    }

    // Insert all records
    await Promise.all(records.map((record) => PnlTable.create(record)));

    // Test with date range filter - createdOnOrAfter
    const startDate = new Date('2023-01-03T00:00:00.000Z');
    const dailyWithStartDate = await PnlTable.findAllDailyAggregate(
      {
        subaccountId: subaccountIds,
        createdOnOrAfter: startDate.toISOString(),
      },
      [],
      {},
    );

    // Should include days 3, 4, and 5 (3 records)
    expect(dailyWithStartDate.results.length).toBe(3);

    // Test with height filters
    const dailyWithHeightRange = await PnlTable.findAllDailyAggregate(
      {
        subaccountId: subaccountIds,
        createdOnOrAfterHeight: '2000',
        createdBeforeOrAtHeight: '4000',
      },
      [],
      {},
    );

    // Should include days 2, 3, and 4 (3 records)
    expect(dailyWithHeightRange.results.length).toBe(3);
    // Check the height range
    const heights = dailyWithHeightRange.results.map((r) => Number(r.createdAtHeight));
    const minHeight = Math.min(...heights);
    const maxHeight = Math.max(...heights);
    expect(minHeight).toBeGreaterThanOrEqual(2000);
    expect(maxHeight).toBeLessThanOrEqual(4000);
  });
});
