import { stats } from '@dydxprotocol-indexer/base';
import {
  dbHelpers,
  testConstants,
  testMocks,
  WalletTable,
  PersistentCacheTable,
  FillTable,
  OrderTable,
  PersistentCacheKeys,
  PersistentCacheFromDatabase,
  WalletFromDatabase,
  BlockTable,
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

  it('Successfully updates totalVolume and persistent cache multiple times', async () => {
    // Set persistent cache totalVolumeUpdateTime to now so task does not backfill
    await PersistentCacheTable.create({
      key: PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME,
      value: DateTime.utc().toISO(),
    });

    // First task run: one new fill
    await FillTable.create({
      ...testConstants.defaultFill,
      createdAt: DateTime.utc().toISO(),
      eventId: testConstants.defaultTendermintEventId,
      price: '1',
      size: '1',
    });

    // Create block to simulate time passing
    const updatedDt1: DateTime = DateTime.utc();
    await BlockTable.create({
      blockHeight: '3',
      time: updatedDt1.toISO(),
    });

    // Run task
    await walletTotalVolumeUpdateTask();

    // Check that wallet updated correctly
    const wallet1: WalletFromDatabase | undefined = await WalletTable
      .findById(testConstants.defaultWallet.address);
    expect(wallet1).toEqual(expect.objectContaining({
      ...testConstants.defaultWallet,
      totalVolume: '1',
    }));

    // Check that persistent cache updated
    const lastUpdateTime1: DateTime | undefined = await getTotalVolumeUpdateTime();
    if (lastUpdateTime1 !== undefined) {
      expect(lastUpdateTime1.toMillis()).toEqual(updatedDt1.toMillis());
    }

    // Second task run: no new fills
    const updatedDt2: DateTime = DateTime.utc();
    await BlockTable.create({
      blockHeight: '4',
      time: updatedDt2.toISO(),
    });
    await walletTotalVolumeUpdateTask();

    const wallet2: WalletFromDatabase | undefined = await WalletTable
      .findById(testConstants.defaultWallet.address);
    expect(wallet2).toEqual(expect.objectContaining({
      ...testConstants.defaultWallet,
      totalVolume: '1',
    }));
    const lastUpdateTime2: DateTime | undefined = await getTotalVolumeUpdateTime();
    if (lastUpdateTime2 !== undefined) {
      expect(lastUpdateTime2.toMillis()).toEqual(updatedDt2.toMillis());
    }

    // Third task run: one new fill
    await FillTable.create({
      ...testConstants.defaultFill,
      createdAt: DateTime.utc().toISO(),
      eventId: testConstants.defaultTendermintEventId2,
      price: '1',
      size: '1',
    });
    const updatedDt3: DateTime = DateTime.utc();
    await BlockTable.create({
      blockHeight: '5',
      time: updatedDt3.toISO(),
    });
    await walletTotalVolumeUpdateTask();

    const wallet3: WalletFromDatabase | undefined = await WalletTable
      .findById(testConstants.defaultWallet.address);
    expect(wallet3).toEqual(expect.objectContaining({
      ...testConstants.defaultWallet,
      totalVolume: '2',
    }));
    const lastUpdateTime3: DateTime | undefined = await getTotalVolumeUpdateTime();
    if (lastUpdateTime3 !== undefined) {
      expect(lastUpdateTime3.toMillis()).toEqual(updatedDt3.toMillis());
    }
  });

  it('Successfully backfills from past date', async () => {
    const currentDt: DateTime = DateTime.utc();

    await Promise.all([
      // Create 3 fills spanning 2 weeks in the past
      FillTable.create({
        ...testConstants.defaultFill,
        createdAt: currentDt.toISO(),
        eventId: testConstants.defaultTendermintEventId,
        price: '1',
        size: '1',
      }),
      FillTable.create({
        ...testConstants.defaultFill,
        createdAt: currentDt.minus({ weeks: 1 }).toISO(),
        eventId: testConstants.defaultTendermintEventId2,
        price: '2',
        size: '2',
      }),
      FillTable.create({
        ...testConstants.defaultFill,
        createdAt: currentDt.minus({ weeks: 2 }).toISO(),
        eventId: testConstants.defaultTendermintEventId3,
        price: '3',
        size: '3',
      }),
      // Set persistent cache totalVolumeUpdateTime to 3 weeks ago to emulate backfill from 3 weeks.
      PersistentCacheTable.create({
        key: PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME,
        value: currentDt.minus({ weeks: 3 }).toISO(),
      }),
      // Create block at current time
      BlockTable.create({
        blockHeight: '3',
        time: DateTime.utc().toISO(),
      }),
    ]);

    // Simulate backfill
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
    // We will simulate a 1 week backfill from the beginning time of
    // `defaultLastUpdateTime`=2023-10-26T00:00:00Z. We do this by leaving persistent cache
    // totalVolumeUpdateTime empty and create fills around `defaultLastUpdateTime`. Then we run
    // the backfill 7 times.
    expect(await getTotalVolumeUpdateTime()).toBeUndefined();

    const referenceDt = DateTime.fromISO('2023-10-26T00:00:00Z');

    await Promise.all([
      FillTable.create({
        ...testConstants.defaultFill,
        createdAt: referenceDt.plus({ days: 1 }).toISO(),
        eventId: testConstants.defaultTendermintEventId,
        price: '1',
        size: '1',
      }),
      FillTable.create({
        ...testConstants.defaultFill,
        createdAt: referenceDt.plus({ days: 4 }).toISO(),
        eventId: testConstants.defaultTendermintEventId2,
        price: '2',
        size: '2',
      }),
      FillTable.create({
        ...testConstants.defaultFill,
        createdAt: referenceDt.plus({ days: 7 }).toISO(),
        eventId: testConstants.defaultTendermintEventId3,
        price: '3',
        size: '3',
      }),
      // Create block in the future relative to referenceDt
      BlockTable.create({
        blockHeight: '3',
        time: referenceDt.plus({ days: 7 }).toISO(),
      }),
    ]);

    // Simulate roundtable runs
    for (let i = 0; i < 7; i++) {
      await walletTotalVolumeUpdateTask();
    }

    const wallet = await WalletTable.findById(testConstants.defaultWallet.address);
    expect(wallet).toEqual(expect.objectContaining({
      ...testConstants.defaultWallet,
      totalVolume: '14', // 1 + 4 + 9
    }));
  });

  it('Successfully records metrics', async () => {
    jest.spyOn(stats, 'gauge');

    await PersistentCacheTable.create({
      key: PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME,
      value: DateTime.utc().toISO(),
    });

    await walletTotalVolumeUpdateTask();

    expect(stats.gauge).toHaveBeenCalledWith(
      `roundtable.persistent_cache_${PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME}_lag_seconds`,
      expect.any(Number),
      { cache: PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME },
    );
  });
});

async function getTotalVolumeUpdateTime(): Promise<DateTime | undefined> {
  const persistentCache: PersistentCacheFromDatabase | undefined = await PersistentCacheTable
    .findById(
      PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME,
    );
  const lastUpdateTime = persistentCache?.value
    ? DateTime.fromISO(persistentCache.value)
    : undefined;
  return lastUpdateTime;
}
