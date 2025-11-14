import {
  Ordering,
  PnlColumns,
} from '../../src/types';
import * as PnlTable from '../../src/stores/pnl-table';
import * as BlockTable from '../../src/stores/block-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import {
  defaultBlock,
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultPnl,
  defaultPnl2,
<<<<<<< HEAD
=======
  defaultBlock,
  isolatedSubaccountId,
>>>>>>> 3f8d74c1 ([ENG-1369] Fix daily PnL aggregation to exclude child subaccounts created mid-day (#3248))
} from '../helpers/constants';

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

  it('Successfully retrieves daily PNL records with latest for current day and earliest for previous days', async () => {
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
    const dailyResults = await PnlTable.findAllDailyPnl(
      { subaccountId: [defaultSubaccountId] },
      [],
      {},
    );

    // We should get exactly 3 records (one for each day)
    expect(dailyResults.results.length).toBe(3);

    // The first record should be the latest one from day 3 (18:00)
    expect(dailyResults.results[0].createdAtHeight).toBe('3018');
    expect(dailyResults.results[0].createdAt).toBe('2023-01-03T18:00:00.000Z');

    // The second record should be the earliest one from day 2 (00:00)
    expect(dailyResults.results[1].createdAtHeight).toBe('2000');
    expect(dailyResults.results[1].createdAt).toBe('2023-01-02T00:00:00.000Z');

    // The third record should be the earliest one from day 1 (00:00)
    expect(dailyResults.results[2].createdAtHeight).toBe('1000');
    expect(dailyResults.results[2].createdAt).toBe('2023-01-01T00:00:00.000Z');

<<<<<<< HEAD
    // Test with pagination - first page
    const dailyPage1 = await PnlTable.findAllDailyPnl(
=======
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
    expect(day2Record?.equity).toBe('4100');
    expect(day2Record?.totalPnl).toBe('410');
    expect(day2Record?.netTransfers).toBe('1250');

    // Day 1 aggregated values
    // Equity sum: 1000 + 1100 = 2100
    // TotalPnl sum: 100 + 110 = 210
    // NetTransfers sum: 500 + 550 = 1050
    expect(day1Record?.equity).toBe('2100');
    expect(day1Record?.totalPnl).toBe('210');
    expect(day1Record?.netTransfers).toBe('1050');
  });

  it('Successfully handles child subaccounts created mid-day by excluding them from daily aggregation', async () => {
    const records = [];

    // Parent subaccount (0) and child subaccount (128)
    const parentSubaccountId = defaultSubaccountId; // subaccount 0
    const childSubaccountId = isolatedSubaccountId; // subaccount 128 (child of 0)

    // Day 1 (Jan 1): Only parent subaccount exists from 00:00
    records.push({
      ...defaultPnl,
      subaccountId: parentSubaccountId,
      createdAt: '2023-01-01T00:00:00.000Z',
      createdAtHeight: '1000',
      equity: '1000',
      totalPnl: '100',
      netTransfers: '500',
    });

    // Day 2 (Jan 2): Parent subaccount has records starting from 00:00
    records.push({
      ...defaultPnl,
      subaccountId: parentSubaccountId,
      createdAt: '2023-01-02T00:00:00.000Z',
      createdAtHeight: '2000',
      equity: '2000',
      totalPnl: '200',
      netTransfers: '600',
    });

    // Add more records throughout the day for parent subaccount
    records.push({
      ...defaultPnl,
      subaccountId: parentSubaccountId,
      createdAt: '2023-01-02T06:00:00.000Z',
      createdAtHeight: '2060',
      equity: '2100',
      totalPnl: '210',
      netTransfers: '610',
    });

    // Day 2 (Jan 2): Child subaccount 128 was created at 05:00
    // This subaccount has NO 00:00 record because it didn't exist yet
    records.push({
      ...defaultPnl,
      subaccountId: childSubaccountId,
      createdAt: '2023-01-02T05:00:00.000Z',
      createdAtHeight: '2050',
      equity: '500', // First equity value for this new child subaccount
      totalPnl: '50',
      netTransfers: '100',
    });

    // Add more records for the newly created child subaccount
    records.push({
      ...defaultPnl,
      subaccountId: childSubaccountId,
      createdAt: '2023-01-02T06:00:00.000Z',
      createdAtHeight: '2061',
      equity: '550',
      totalPnl: '55',
      netTransfers: '105',
    });

    records.push({
      ...defaultPnl,
      subaccountId: childSubaccountId,
      createdAt: '2023-01-02T12:00:00.000Z',
      createdAtHeight: '2120',
      equity: '600',
      totalPnl: '60',
      netTransfers: '110',
    });

    // Day 3 (Jan 3): Both parent and child subaccounts have records at 00:00
    records.push({
      ...defaultPnl,
      subaccountId: parentSubaccountId,
      createdAt: '2023-01-03T00:00:00.000Z',
      createdAtHeight: '3000',
      equity: '3000',
      totalPnl: '300',
      netTransfers: '700',
    });

    records.push({
      ...defaultPnl,
      subaccountId: childSubaccountId,
      createdAt: '2023-01-03T00:00:00.000Z',
      createdAtHeight: '3001',
      equity: '800',
      totalPnl: '80',
      netTransfers: '150',
    });

    // Add more records throughout Day 3
    records.push({
      ...defaultPnl,
      subaccountId: parentSubaccountId,
      createdAt: '2023-01-03T06:00:00.000Z',
      createdAtHeight: '3060',
      equity: '3100',
      totalPnl: '310',
      netTransfers: '710',
    });

    records.push({
      ...defaultPnl,
      subaccountId: childSubaccountId,
      createdAt: '2023-01-03T06:00:00.000Z',
      createdAtHeight: '3061',
      equity: '850',
      totalPnl: '85',
      netTransfers: '155',
    });

    // Insert all records
    await Promise.all(records.map((record) => PnlTable.create(record)));

    // Get aggregated daily records for parent subaccount and its children
    // This simulates the parentSubaccountNumber endpoint
    const dailyResults = await PnlTable.findAllDailyAggregate(
      { subaccountId: [parentSubaccountId, childSubaccountId] },
      [],
      {},
    );

    // We should get exactly 3 records (one for each day)
    expect(dailyResults.results.length).toBe(3);

    // Verify Day 1 (only parent subaccount existed)
    const day1Record = dailyResults.results.find((r) => r.createdAt.includes('2023-01-01'));
    expect(day1Record).toBeDefined();
    expect(day1Record?.equity).toBe('1000'); // Only parent
    expect(day1Record?.totalPnl).toBe('100');
    expect(day1Record?.netTransfers).toBe('500');
    expect(day1Record?.createdAt).toBe('2023-01-01T00:00:00.000Z');

    // Verify Day 2 (child created at 05:00, should be excluded)
    const day2Record = dailyResults.results.find((r) => r.createdAt.includes('2023-01-02'));
    expect(day2Record).toBeDefined();

    // Critical: The aggregation should ONLY include records at 00:00:00
    // Since childSubaccountId (128) has NO record at 00:00:00, it should be EXCLUDED
    // Expected: equity = 2000 (only from parent subaccount 0 at 00:00)
    // NOT: equity = 2000 + 500 = 2500 (mixing timestamps - this was the BUG)
    expect(day2Record?.equity).toBe('2000');
    expect(day2Record?.totalPnl).toBe('200');
    expect(day2Record?.netTransfers).toBe('600');
    expect(day2Record?.createdAt).toBe('2023-01-02T00:00:00.000Z');

    // The height should be from 00:00, not from 05:00
    expect(day2Record?.createdAtHeight).toBe('2000'); // MAX height from 00:00 records

    // Verify Day 3 (both parent and child have 00:00 records - should aggregate both)
    const day3Record = dailyResults.results.find((r) => r.createdAt.includes('2023-01-03'));
    expect(day3Record).toBeDefined();

    // Now both subaccounts should be included in the aggregation
    // Equity: 3000 (parent) + 800 (child) = 3800
    // TotalPnl: 300 (parent) + 80 (child) = 380
    // NetTransfers: 700 (parent) + 150 (child) = 850
    expect(day3Record?.equity).toBe('3800');
    expect(day3Record?.totalPnl).toBe('380');
    expect(day3Record?.netTransfers).toBe('850');
    expect(day3Record?.createdAt).toBe('2023-01-03T00:00:00.000Z');
    expect(day3Record?.createdAtHeight).toBe('3001'); // MAX height from 00:00 records
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
>>>>>>> 3f8d74c1 ([ENG-1369] Fix daily PnL aggregation to exclude child subaccounts created mid-day (#3248))
      {
        subaccountId: [defaultSubaccountId],
        page: 1,
        limit: 2,
      },
      [],
      {},
    );

    expect(dailyPage1.results.length).toBe(2);
    expect(dailyPage1.limit).toBe(2);
    expect(dailyPage1.offset).toBe(0);
    expect(dailyPage1.total).toBe(3);

    // First page should have day 3 (latest) and day 2 (earliest)
    expect(dailyPage1.results[0].createdAtHeight).toBe('3018');
    expect(dailyPage1.results[1].createdAtHeight).toBe('2000');

    // Test with pagination - second page
    const dailyPage2 = await PnlTable.findAllDailyPnl(
      {
        subaccountId: [defaultSubaccountId],
        page: 2,
        limit: 2,
      },
      [],
      {},
    );

    // The second page should have only day 1
    expect(dailyPage2.results.length).toBe(1);
    expect(dailyPage2.limit).toBe(2);
    expect(dailyPage2.offset).toBe(2);
    expect(dailyPage2.total).toBe(3);
    expect(dailyPage2.results[0].createdAtHeight).toBe('1000');

    // Test with date range filter
    const cutoffDate = new Date('2023-01-02T12:00:00.000Z');

    const dailyWithDateFilter = await PnlTable.findAllDailyPnl(
      {
        subaccountId: [defaultSubaccountId],
        createdBeforeOrAt: cutoffDate.toISOString(),
      },
      [],
      {},
    );

    // We should get 2 records: day 1 (earliest) and day 2 (records up to 12:00)
    expect(dailyWithDateFilter.results.length).toBe(2);

    // Day 2 should be represented by the latest record before our cutoff (12:00)
    expect(dailyWithDateFilter.results[0].createdAtHeight).toBe('2012');
    expect(dailyWithDateFilter.results[0].createdAt).toBe('2023-01-02T12:00:00.000Z');

    // Day 1 should still be the earliest record
    expect(dailyWithDateFilter.results[1].createdAtHeight).toBe('1000');
    expect(dailyWithDateFilter.results[1].createdAt).toBe('2023-01-01T00:00:00.000Z');
  });

  it('Successfully handles case where latest record is at midnight (00:00)', async () => {
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

    // Day 3 (Jan 3): Create ONLY a record at 00:00 (to test the case where latest is at midnight)
    // Give this record the highest height to ensure it's the latest
    records.push({
      ...defaultPnl,
      createdAt: '2023-01-03T00:00:00.000Z',
      createdAtHeight: '3500', // Highest height to ensure it's the latest
      equity: '3500',
    });

    // Insert all records
    await Promise.all(records.map((record) => PnlTable.create(record)));

    // Get daily records
    const dailyResults = await PnlTable.findAllDailyPnl(
      { subaccountId: [defaultSubaccountId] },
      [],
      {},
    );

    // We should get exactly 3 records (one for each day)
    expect(dailyResults.results.length).toBe(3);

    // The first record should be the latest one from day 3 (00:00)
    expect(dailyResults.results[0].createdAtHeight).toBe('3500');
    expect(dailyResults.results[0].createdAt).toBe('2023-01-03T00:00:00.000Z');

    // The second record should be the earliest one from day 2 (00:00)
    expect(dailyResults.results[1].createdAtHeight).toBe('2000');
    expect(dailyResults.results[1].createdAt).toBe('2023-01-02T00:00:00.000Z');

    // The third record should be the earliest one from day 1 (00:00)
    expect(dailyResults.results[2].createdAtHeight).toBe('1000');
    expect(dailyResults.results[2].createdAt).toBe('2023-01-01T00:00:00.000Z');
  });
});
