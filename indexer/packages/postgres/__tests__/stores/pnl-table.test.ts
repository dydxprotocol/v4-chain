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

    // Test with pagination - first page
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
