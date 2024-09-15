import {
  dbHelpers,
  testConstants,
  testMocks,
  WalletTable,
  SubaccountTable,
  PersistentCacheTable,
  FillTable,
  OrderTable,
} from '@dydxprotocol-indexer/postgres';
import walletTotalVolumeUpdateTask from '../../src/tasks/update-wallet-total-volume';
import { DateTime } from 'luxon';

describe('update-wallet-total-volume', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await OrderTable.create(testConstants.defaultOrder);
    await OrderTable.create(testConstants.isolatedMarketOrder);
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  it('Successfully updates totalVolume multiple times', async () => {
    const defaultSubaccountId = await SubaccountTable.findAll(
      { subaccountNumber: testConstants.defaultSubaccount.subaccountNumber },
      [],
      {},
    );
    // Set persistent cache totalVolumeUpdateTime so walletTotalVolumeUpdateTask() does not attempt
    // to backfill
    await PersistentCacheTable.create({
      key: 'totalVolumeUpdateTime',
      value: DateTime.utc().toISO(),
    });

    // First task run: one new fill
    await FillTable.create({
      ...testConstants.defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: DateTime.utc().toISO(),
      eventId: testConstants.defaultTendermintEventId,
      price: '1',
      size: '1',
    });
    await walletTotalVolumeUpdateTask();
    let wallet = await WalletTable.findById(testConstants.defaultWallet.address);
    expect(wallet).toEqual(expect.objectContaining({
      ...testConstants.defaultWallet,
      totalVolume: '1',
    }));

    // Second task run: no new fills
    await walletTotalVolumeUpdateTask();
    wallet = await WalletTable.findById(testConstants.defaultWallet.address);
    expect(wallet).toEqual(expect.objectContaining({
      ...testConstants.defaultWallet,
      totalVolume: '1',
    }));

    // Third task run: one new fill
    await FillTable.create({
      ...testConstants.defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: DateTime.utc().toISO(),
      eventId: testConstants.defaultTendermintEventId2,
      price: '1',
      size: '1',
    });
    await walletTotalVolumeUpdateTask();
    wallet = await WalletTable.findById(testConstants.defaultWallet.address);
    expect(wallet).toEqual(expect.objectContaining({
      ...testConstants.defaultWallet,
      totalVolume: '2',
    }));
  });

  it('Successfully updates totalVolumeUpdateTime in persistent cache table', async () => {
    // Set persistent cache totalVolumeUpdateTime so walletTotalVolumeUpdateTask() does not attempt
    // to backfill
    await PersistentCacheTable.create({
      key: 'totalVolumeUpdateTime',
      value: DateTime.utc().toISO(),
    });

    await walletTotalVolumeUpdateTask();
    const lastUpdateTime1 = await getTotalVolumeUpdateTime();

    await walletTotalVolumeUpdateTask();
    const lastUpdateTime2 = await getTotalVolumeUpdateTime();

    expect(lastUpdateTime1).not.toBeUndefined();
    expect(lastUpdateTime2).not.toBeUndefined();
    if (lastUpdateTime1?.toMillis() !== undefined && lastUpdateTime2?.toMillis() !== undefined) {
      expect(lastUpdateTime2.toMillis())
        .toBeGreaterThan(lastUpdateTime1.toMillis());
    }
  });

  it('Successfully backfills from past date', async () => {
    const currentDt: DateTime = DateTime.utc();
    const defaultSubaccountId = await SubaccountTable.findAll(
      { subaccountNumber: testConstants.defaultSubaccount.subaccountNumber },
      [],
      {},
    );

    // Create 3 fills spanning 2 weeks in the past
    await FillTable.create({
      ...testConstants.defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: currentDt.toISO(),
      eventId: testConstants.defaultTendermintEventId,
      price: '1',
      size: '1',
    });
    await FillTable.create({
      ...testConstants.defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: currentDt.minus({ weeks: 1 }).toISO(),
      eventId: testConstants.defaultTendermintEventId2,
      price: '2',
      size: '2',
    });
    await FillTable.create({
      ...testConstants.defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: currentDt.minus({ weeks: 2 }).toISO(),
      eventId: testConstants.defaultTendermintEventId3,
      price: '3',
      size: '3',
    });

    // Set persistent cache totalVolumeUpdateTime to 3 weeks ago to emulate backfill from 3 weeks.
    await PersistentCacheTable.create({
      key: 'totalVolumeUpdateTime',
      value: currentDt.minus({ weeks: 3 }).toISO(),
    });

    let backfillTime = await getTotalVolumeUpdateTime();
    while (backfillTime !== undefined && DateTime.fromISO(backfillTime.toISO()) < currentDt) {
      await walletTotalVolumeUpdateTask();
      backfillTime = await getTotalVolumeUpdateTime();
    }

    const wallet = await WalletTable.findById(testConstants.defaultWallet.address);
    expect(wallet).toEqual(expect.objectContaining({
      ...testConstants.defaultWallet,
      totalVolume: '14', // 1 + 4 + 9
    }));
  });

  it('Successfully backfills on first run', async () => {
    const defaultSubaccountId = await SubaccountTable.findAll(
      { subaccountNumber: testConstants.defaultSubaccount.subaccountNumber },
      [],
      {},
    );

    // Leave persistent cache totalVolumeUpdateTime empty and create fills around
    // `defaultLastUpdateTime` value to emulate backfilling from very beginning
    expect(await getTotalVolumeUpdateTime()).toBeUndefined();

    const referenceDt = DateTime.fromISO('2020-01-01T00:00:00Z');

    await FillTable.create({
      ...testConstants.defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: referenceDt.plus({ days: 1 }).toISO(),
      eventId: testConstants.defaultTendermintEventId,
      price: '1',
      size: '1',
    });
    await FillTable.create({
      ...testConstants.defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: referenceDt.plus({ days: 2 }).toISO(),
      eventId: testConstants.defaultTendermintEventId2,
      price: '2',
      size: '2',
    });
    await FillTable.create({
      ...testConstants.defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: referenceDt.plus({ days: 3 }).toISO(),
      eventId: testConstants.defaultTendermintEventId3,
      price: '3',
      size: '3',
    });

    // Emulate 10 roundtable runs (this should backfill all the fills)
    for (let i = 0; i < 10; i++) {
      await walletTotalVolumeUpdateTask();
    }

    const wallet = await WalletTable.findById(testConstants.defaultWallet.address);
    expect(wallet).toEqual(expect.objectContaining({
      ...testConstants.defaultWallet,
      totalVolume: '14', // 1 + 4 + 9
    }));
  });
});

async function getTotalVolumeUpdateTime(): Promise<DateTime | undefined> {
  const persistentCache = await PersistentCacheTable.findById('totalVolumeUpdateTime');
  const lastUpdateTime1 = persistentCache?.value
    ? DateTime.fromISO(persistentCache.value)
    : undefined;
  return lastUpdateTime1;
}
