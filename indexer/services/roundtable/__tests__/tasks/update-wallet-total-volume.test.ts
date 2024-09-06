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

  it('Succeeds in populating historical totalVolume on first run', async () => {
    const referenceDt: DateTime = DateTime.fromISO('2020-01-01T00:00:00Z');
    const defaultSubaccountId = await SubaccountTable.findAll(
      { subaccountNumber: testConstants.defaultSubaccount.subaccountNumber },
      [],
      {},
    );

    await FillTable.create({
      ...testConstants.defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: referenceDt.toISO(),
      eventId: testConstants.defaultTendermintEventId,
      price: '1',
      size: '1',
    });
    await FillTable.create({
      ...testConstants.defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: referenceDt.plus({ years: 1 }).toISO(),
      eventId: testConstants.defaultTendermintEventId2,
      price: '2',
      size: '2',
    });
    await FillTable.create({
      ...testConstants.defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: referenceDt.plus({ years: 2 }).toISO(),
      eventId: testConstants.defaultTendermintEventId3,
      price: '3',
      size: '3',
    });

    await walletTotalVolumeUpdateTask();

    const wallet = await WalletTable.findById(testConstants.defaultWallet.address);
    expect(wallet).toEqual(expect.objectContaining({
      ...testConstants.defaultWallet,
      totalVolume: '14', // 1 + 4 + 9
    }));
  });

  it('Succeeds in incremental totalVolume updates', async () => {
    const defaultSubaccountId = await SubaccountTable.findAll(
      { subaccountNumber: testConstants.defaultSubaccount.subaccountNumber },
      [],
      {},
    );

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
    await walletTotalVolumeUpdateTask();
    const lastUpdateTime1 = await getTotalVolumeUpdateTime();

    // Sleep for 1s
    await new Promise((resolve) => setTimeout(resolve, 1000));

    await walletTotalVolumeUpdateTask();
    const lastUpdateTime2 = await getTotalVolumeUpdateTime();

    expect(lastUpdateTime1).not.toBeUndefined();
    expect(lastUpdateTime2).not.toBeUndefined();
    if (lastUpdateTime1?.toMillis() !== undefined && lastUpdateTime2?.toMillis() !== undefined) {
      expect(lastUpdateTime2.toMillis())
        .toBeGreaterThan(lastUpdateTime1.plus({ seconds: 1 }).toMillis());
    }
  });
});

async function getTotalVolumeUpdateTime(): Promise<DateTime | undefined> {
  const persistentCache = await PersistentCacheTable.findById('totalVolumeUpdateTime');
  const lastUpdateTime1 = persistentCache?.value
    ? DateTime.fromISO(persistentCache.value)
    : undefined;
  return lastUpdateTime1;
}
