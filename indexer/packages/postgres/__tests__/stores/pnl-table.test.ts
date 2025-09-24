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

  it('Successfully retrieves daily PNL records (every 24th record)', async () => {
  // Create 50 PNL records for the same subaccount with hourly intervals
    const hourlyRecords = [];
    const baseDate = new Date('2023-01-01T00:00:00.000Z');
    for (let i = 0; i < 50; i++) {
      const date = new Date(baseDate);
      date.setUTCHours(baseDate.getUTCHours() + i);

      hourlyRecords.push({
        ...defaultPnl,
        createdAt: date.toISOString(),
        createdAtHeight: (1000 + i).toString(), // Incrementing heights
        equity: (1000 + i).toString(), // Different equity values to verify correct records
      });
    }

    // Insert all records
    await Promise.all(
      hourlyRecords.map((record) => PnlTable.create(record)),
    );

    // Test 1: Get all daily records without pagination
    const allDailyResults = await PnlTable.findAllDailyPnl(
      { subaccountId: [defaultSubaccountId] },
      [],
      {},
    );

    // We should get the first record (latest) and then one every 24 hours
    // Expected: Latest + 2 days = 3 records
    expect(allDailyResults.results.length).toBeLessThanOrEqual(3);
    expect(allDailyResults.results.length).toBeGreaterThanOrEqual(2);

    // Verify we got the correct records (the latest and then every 24th)
    // The first record should be the latest one (highest height)
    expect(allDailyResults.results[0].createdAtHeight).toBe('1049');

    // Test 2: Test with pagination - first page
    const dailyPage1 = await PnlTable.findAllDailyPnl(
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
    // The total should match the total number of daily records
    expect(dailyPage1.total).toBeLessThanOrEqual(3);
    expect(dailyPage1.total).toBeGreaterThanOrEqual(2);

    // Test 3: Test with pagination - second page
    const dailyPage2 = await PnlTable.findAllDailyPnl(
      {
        subaccountId: [defaultSubaccountId],
        page: 2,
        limit: 2,
      },
      [],
      {},
    );

    // The second page should contain any remaining daily records
    const expectedPage2Length = Math.max(0, allDailyResults.results.length - 2);
    expect(dailyPage2.results.length).toBe(expectedPage2Length);
    expect(dailyPage2.limit).toBe(2);
    expect(dailyPage2.offset).toBe(2);
    expect(dailyPage2.total).toBe(allDailyResults.results.length);

    // Test 4: Test with limit only (no pagination)
    const dailyWithLimit = await PnlTable.findAllDailyPnl(
      {
        subaccountId: [defaultSubaccountId],
        limit: 1,
      },
      [],
      {},
    );

    expect(dailyWithLimit.results.length).toBe(1);
    // Should be the latest record
    expect(dailyWithLimit.results[0].createdAtHeight).toBe('1049');

    // Test 5: Test with date range filter
    const middleDate = new Date(baseDate);
    middleDate.setHours(baseDate.getHours() + 25); // Just after the first day

    const dailyWithDateFilter = await PnlTable.findAllDailyPnl(
      {
        subaccountId: [defaultSubaccountId],
        createdBeforeOrAt: middleDate.toISOString(),
      },
      [],
      {},
    );

    // We should only get daily records before or at the middle date
    // This should be at most 2 records (latest within range + one 24 hours before)
    expect(dailyWithDateFilter.results.length).toBeLessThanOrEqual(2);

    // The highest height should be the record at or just before our date filter
    expect(parseInt(dailyWithDateFilter.results[0].createdAtHeight, 10))
      .toBeLessThanOrEqual(1000 + 25);
  });
});
