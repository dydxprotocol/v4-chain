import { stats } from '@dydxprotocol-indexer/base';
import {
  dbHelpers,
  testConstants,
  testMocks,
  PersistentCacheTable,
  FillTable,
  OrderTable,
  PersistentCacheKeys,
  AffiliateReferredUsersTable,
  AffiliateInfoFromDatabase,
  AffiliateInfoTable,
  Liquidity,
  PersistentCacheFromDatabase,
  BlockTable,
} from '@dydxprotocol-indexer/postgres';
import affiliateInfoUpdateTask from '../../src/tasks/update-affiliate-info';
import { DateTime } from 'luxon';

describe('update-affiliate-info', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await OrderTable.create(testConstants.defaultOrder);
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  it('Successfully updates affiliate info and persistent cache multiple times', async () => {
    // Set persistent cache affiliateInfoUpdateTime to slightly in past so task does not backfill
    await PersistentCacheTable.create({
      key: PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME,
      value: DateTime.utc().toISO(),
    });

    const updatedDt1: DateTime = DateTime.utc();
    await Promise.all([
      // First task run: add referral w/o any fills
      // defaultWallet2 will be affiliate and defaultWallet will be referee
      AffiliateReferredUsersTable.create({
        affiliateAddress: testConstants.defaultWallet2.address,
        refereeAddress: testConstants.defaultWallet.address,
        referredAtBlock: '1',
      }),

      // Create block to simulate time passing
      BlockTable.create({
        blockHeight: '3',
        time: updatedDt1.toISO(),
      }),
    ]);

    // Run task
    await affiliateInfoUpdateTask();

    const updatedInfo1: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
      testConstants.defaultWallet2.address,
    );
    const expectedAffiliateInfo1: AffiliateInfoFromDatabase = {
      address: testConstants.defaultWallet2.address,
      affiliateEarnings: '0',
      referredMakerTrades: 0,
      referredTakerTrades: 0,
      totalReferredMakerFees: '0',
      totalReferredTakerFees: '0',
      totalReferredMakerRebates: '0',
      totalReferredUsers: 1,
      firstReferralBlockHeight: '1',
      referredTotalVolume: '0',
    };
    expect(updatedInfo1).toEqual(expectedAffiliateInfo1);

    // Check that persistent cache updated
    const lastUpdateTime1: DateTime | undefined = await getAffiliateInfoUpdateTime();
    if (lastUpdateTime1 !== undefined) {
      expect(lastUpdateTime1.toMillis()).toEqual(updatedDt1.toMillis());
    }

    // Second task run: one new fill and one new referral
    await Promise.all([
      FillTable.create({
        ...testConstants.defaultFill,
        liquidity: Liquidity.TAKER,
        createdAt: DateTime.utc().toISO(),
        eventId: testConstants.defaultTendermintEventId,
        price: '1',
        size: '1',
        fee: '1000',
        affiliateRevShare: '500',
      }),
      AffiliateReferredUsersTable.create({
        affiliateAddress: testConstants.defaultWallet2.address,
        refereeAddress: testConstants.defaultWallet3.address,
        referredAtBlock: '2',
      }),
    ]);

    const updatedDt2: DateTime = DateTime.utc();
    await BlockTable.create({
      blockHeight: '4',
      time: updatedDt2.toISO(),
    });

    await affiliateInfoUpdateTask();

    const updatedInfo2: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
      testConstants.defaultWallet2.address,
    );
    const expectedAffiliateInfo2: AffiliateInfoFromDatabase = {
      address: testConstants.defaultWallet2.address,
      affiliateEarnings: '500',
      referredMakerTrades: 0,
      referredTakerTrades: 1,
      totalReferredMakerFees: '0',
      totalReferredTakerFees: '1000',
      totalReferredMakerRebates: '0',
      totalReferredUsers: 2,
      firstReferralBlockHeight: '1',
      referredTotalVolume: '1',
    };
    expect(updatedInfo2).toEqual(expectedAffiliateInfo2);
    const lastUpdateTime2: DateTime | undefined = await getAffiliateInfoUpdateTime();
    if (lastUpdateTime2 !== undefined) {
      expect(lastUpdateTime2.toMillis()).toEqual(updatedDt2.toMillis());
    }
  });

  it('Successfully backfills from past date', async () => {
    const currentDt: DateTime = DateTime.utc();

    await Promise.all([
      // Set persistent cache to 3 weeks ago to emulate backfill from 3 weeks.
      PersistentCacheTable.create({
        key: PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME,
        value: currentDt.minus({ weeks: 3 }).toISO(),
      }),
      // defaultWallet2 will be affiliate and defaultWallet will be referee
      AffiliateReferredUsersTable.create({
        affiliateAddress: testConstants.defaultWallet2.address,
        refereeAddress: testConstants.defaultWallet.address,
        referredAtBlock: '1',
      }),
      // Fills spannings 2 weeks
      FillTable.create({
        ...testConstants.defaultFill,
        liquidity: Liquidity.TAKER,
        createdAt: currentDt.minus({ weeks: 1 }).toISO(),
        eventId: testConstants.defaultTendermintEventId,
        price: '1',
        size: '1',
        fee: '1000',
        affiliateRevShare: '500',
      }),
      FillTable.create({
        ...testConstants.defaultFill,
        liquidity: Liquidity.TAKER,
        createdAt: currentDt.minus({ weeks: 2 }).toISO(),
        eventId: testConstants.defaultTendermintEventId2,
        price: '1',
        size: '1',
        fee: '1000',
        affiliateRevShare: '500',
      }),
      // Create block at current time
      BlockTable.create({
        blockHeight: '3',
        time: DateTime.utc().toISO(),
      }),
    ]);

    // Simulate backfill
    let backfillTime: DateTime | undefined = await getAffiliateInfoUpdateTime();
    while (backfillTime !== undefined && DateTime.fromISO(backfillTime.toISO()) < currentDt) {
      await affiliateInfoUpdateTask();
      backfillTime = await getAffiliateInfoUpdateTime();
    }

    const expectedAffiliateInfo: AffiliateInfoFromDatabase = {
      address: testConstants.defaultWallet2.address,
      affiliateEarnings: '1000',
      referredMakerTrades: 0,
      referredTakerTrades: 2,
      totalReferredMakerFees: '0',
      totalReferredTakerFees: '2000',
      totalReferredMakerRebates: '0',
      totalReferredUsers: 1,
      firstReferralBlockHeight: '1',
      referredTotalVolume: '2',
    };
    const updatedInfo: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable
      .findById(testConstants.defaultWallet2.address);
    expect(updatedInfo).toEqual(expectedAffiliateInfo);
  });

  it('Successfully backfills on first run', async () => {
    // We will simulate a 1 week backfill from the beginning time of
    // `defaultLastUpdateTime`=2024-09-16T00:00:00Z. We do this by leaving persistent cache
    // affiliateInfoUpdateTime empty and create fills around `defaultLastUpdateTime`. Then we run
    // the backfill 7 times.
    expect(await getAffiliateInfoUpdateTime()).toBeUndefined();

    const referenceDt: DateTime = DateTime.fromISO('2024-09-16T00:00:00Z');

    // defaultWallet2 will be affiliate and defaultWallet will be referee
    await AffiliateReferredUsersTable.create({
      affiliateAddress: testConstants.defaultWallet2.address,
      refereeAddress: testConstants.defaultWallet.address,
      referredAtBlock: '1',
    });

    await Promise.all([
      // Fills spannings 7 days after referenceDt
      FillTable.create({
        ...testConstants.defaultFill,
        liquidity: Liquidity.TAKER,
        createdAt: referenceDt.plus({ days: 1 }).toISO(),
        eventId: testConstants.defaultTendermintEventId,
        price: '1',
        size: '1',
        fee: '1000',
        affiliateRevShare: '500',
      }),
      FillTable.create({
        ...testConstants.defaultFill,
        liquidity: Liquidity.TAKER,
        createdAt: referenceDt.plus({ days: 7 }).toISO(),
        eventId: testConstants.defaultTendermintEventId2,
        price: '1',
        size: '1',
        fee: '1000',
        affiliateRevShare: '500',
      }),
      // Create block in the future relative to referenceDt
      BlockTable.create({
        blockHeight: '3',
        time: referenceDt.plus({ days: 7 }).toISO(),
      }),
    ]);

    // Simulate roundtable runs
    for (let i = 0; i < 7; i++) {
      await affiliateInfoUpdateTask();
    }

    const expectedAffiliateInfo: AffiliateInfoFromDatabase = {
      address: testConstants.defaultWallet2.address,
      affiliateEarnings: '1000',
      referredMakerTrades: 0,
      referredTakerTrades: 2,
      totalReferredMakerFees: '0',
      totalReferredTakerFees: '2000',
      totalReferredMakerRebates: '0',
      totalReferredUsers: 1,
      firstReferralBlockHeight: '1',
      referredTotalVolume: '2',
    };
    const updatedInfo = await AffiliateInfoTable.findById(testConstants.defaultWallet2.address);
    expect(updatedInfo).toEqual(expectedAffiliateInfo);
  });

  it('Successfully records metrics', async () => {
    jest.spyOn(stats, 'gauge');

    await PersistentCacheTable.create({
      key: PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME,
      value: DateTime.utc().toISO(),
    });

    await affiliateInfoUpdateTask();

    expect(stats.gauge).toHaveBeenCalledWith(
      `roundtable.persistent_cache_${PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME}_lag_seconds`,
      expect.any(Number),
      { cache: PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME },
    );
  });
});

async function getAffiliateInfoUpdateTime(): Promise<DateTime | undefined> {
  const persistentCache: PersistentCacheFromDatabase | undefined = await PersistentCacheTable
    .findById(
      PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME,
    );
  const lastUpdateTime: DateTime | undefined = persistentCache?.value
    ? DateTime.fromISO(persistentCache.value)
    : undefined;
  return lastUpdateTime;
}
