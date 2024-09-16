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
    const startDt = DateTime.utc();

    // Set persistent cache affiliateInfoUpdateTIme so task does not use backfill windows
    await PersistentCacheTable.create({
      key: PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME,
      value: startDt.toISO(),
    });

    // First task run: add refereal w/o any fills
    // defaultWallet2 will be affiliate and defaultWallet will be referee
    await AffiliateReferredUsersTable.create({
      affiliateAddress: testConstants.defaultWallet2.address,
      refereeAddress: testConstants.defaultWallet.address,
      referredAtBlock: '1',
    });
    await affiliateInfoUpdateTask();

    let updatedInfo: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
      testConstants.defaultWallet2.address,
    );
    let expectedAffiliateInfo: AffiliateInfoFromDatabase = {
      address: testConstants.defaultWallet2.address,
      affiliateEarnings: '0',
      referredMakerTrades: 0,
      referredTakerTrades: 0,
      totalReferredFees: '0',
      totalReferredUsers: 1,
      referredNetProtocolEarnings: '0',
      firstReferralBlockHeight: '1',
      referredTotalVolume: '0',
    };
    expect(updatedInfo).toEqual(expect.objectContaining(expectedAffiliateInfo));

    // Check that persistent cache updated
    const lastUpdateTime1 = await getAffiliateInfoUpdateTime();
    if (lastUpdateTime1 !== undefined) {
      expect(lastUpdateTime1.toMillis())
        .toBeGreaterThan(startDt.toMillis());
    }

    // Second task run: one new fill and one new referral
    await FillTable.create({
      ...testConstants.defaultFill,
      liquidity: Liquidity.TAKER,
      createdAt: DateTime.utc().toISO(),
      eventId: testConstants.defaultTendermintEventId,
      price: '1',
      size: '1',
      fee: '1000',
      affiliateRevShare: '500',
    });
    await AffiliateReferredUsersTable.create({
      affiliateAddress: testConstants.defaultWallet2.address,
      refereeAddress: testConstants.defaultWallet3.address,
      referredAtBlock: '2',
    });

    await affiliateInfoUpdateTask();

    updatedInfo = await AffiliateInfoTable.findById(
      testConstants.defaultWallet2.address,
    );
    expectedAffiliateInfo = {
      address: testConstants.defaultWallet2.address,
      affiliateEarnings: '500',
      referredMakerTrades: 0,
      referredTakerTrades: 1,
      totalReferredFees: '1000',
      totalReferredUsers: 2,
      referredNetProtocolEarnings: '500',
      firstReferralBlockHeight: '1',
      referredTotalVolume: '1',
    };
    expect(updatedInfo).toEqual(expectedAffiliateInfo);
    const lastUpdateTime2 = await getAffiliateInfoUpdateTime();
    if (lastUpdateTime2 !== undefined && lastUpdateTime1 !== undefined) {
      expect(lastUpdateTime2.toMillis())
        .toBeGreaterThan(lastUpdateTime1.toMillis());
    }
  });

  it('Successfully backfills from past date', async () => {
    const currentDt = DateTime.utc();

    // Set persistent cache to 3 weeks ago to emulate backfill from 3 weeks.
    await PersistentCacheTable.create({
      key: PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME,
      value: currentDt.minus({ weeks: 3 }).toISO(),
    });

    // defaultWallet2 will be affiliate and defaultWallet will be referee
    await AffiliateReferredUsersTable.create({
      affiliateAddress: testConstants.defaultWallet2.address,
      refereeAddress: testConstants.defaultWallet.address,
      referredAtBlock: '1',
    });

    // Fills spannings 2 weeks
    await FillTable.create({
      ...testConstants.defaultFill,
      liquidity: Liquidity.TAKER,
      createdAt: currentDt.minus({ weeks: 1 }).toISO(),
      eventId: testConstants.defaultTendermintEventId,
      price: '1',
      size: '1',
      fee: '1000',
      affiliateRevShare: '500',
    });
    await FillTable.create({
      ...testConstants.defaultFill,
      liquidity: Liquidity.TAKER,
      createdAt: currentDt.minus({ weeks: 2 }).toISO(),
      eventId: testConstants.defaultTendermintEventId2,
      price: '1',
      size: '1',
      fee: '1000',
      affiliateRevShare: '500',
    });

    // Simulate backfill
    let backfillTime = await getAffiliateInfoUpdateTime();
    while (backfillTime !== undefined && DateTime.fromISO(backfillTime.toISO()) < currentDt) {
      await affiliateInfoUpdateTask();
      backfillTime = await getAffiliateInfoUpdateTime();
    }

    const expectedAffiliateInfo = {
      address: testConstants.defaultWallet2.address,
      affiliateEarnings: '1000',
      referredMakerTrades: 0,
      referredTakerTrades: 2,
      totalReferredFees: '2000',
      totalReferredUsers: 1,
      referredNetProtocolEarnings: '1000',
      firstReferralBlockHeight: '1',
      totalReferredVolume: '2',
    };
    const updatedInfo = await AffiliateInfoTable.findById(testConstants.defaultWallet2.address);
    expect(updatedInfo).toEqual(expectedAffiliateInfo);
  });

  it('Successfully backfills on first run', async () => {
    // Leave persistent cache affiliateInfoUpdateTime empty and create fills around
    // `defaultLastUpdateTime` value to emulate backfilling from very beginning
    expect(await getAffiliateInfoUpdateTime()).toBeUndefined();

    const referenceDt = DateTime.fromISO('2020-01-01T00:00:00Z');

    // defaultWallet2 will be affiliate and defaultWallet will be referee
    await AffiliateReferredUsersTable.create({
      affiliateAddress: testConstants.defaultWallet2.address,
      refereeAddress: testConstants.defaultWallet.address,
      referredAtBlock: '1',
    });

    // Fills spannings 2 weeks after referenceDt
    await FillTable.create({
      ...testConstants.defaultFill,
      liquidity: Liquidity.TAKER,
      createdAt: referenceDt.plus({ weeks: 1 }).toISO(),
      eventId: testConstants.defaultTendermintEventId,
      price: '1',
      size: '1',
      fee: '1000',
      affiliateRevShare: '500',
    });
    await FillTable.create({
      ...testConstants.defaultFill,
      liquidity: Liquidity.TAKER,
      createdAt: referenceDt.plus({ weeks: 2 }).toISO(),
      eventId: testConstants.defaultTendermintEventId2,
      price: '1',
      size: '1',
      fee: '1000',
      affiliateRevShare: '500',
    });

    // Simulate 20 roundtable runs (this is enough to backfill all the fills)
    for (let i = 0; i < 20; i++) {
      await affiliateInfoUpdateTask();
    }

    const expectedAffiliateInfo = {
      address: testConstants.defaultWallet2.address,
      affiliateEarnings: '1000',
      referredMakerTrades: 0,
      referredTakerTrades: 2,
      totalReferredFees: '2000',
      totalReferredUsers: 1,
      referredNetProtocolEarnings: '1000',
      firstReferralBlockHeight: '1',
      totalReferredVolume: '2',
    };
    const updatedInfo = await AffiliateInfoTable.findById(testConstants.defaultWallet2.address);
    expect(updatedInfo).toEqual(expectedAffiliateInfo);
  });
});

async function getAffiliateInfoUpdateTime(): Promise<DateTime | undefined> {
  const persistentCache = await PersistentCacheTable.findById(
    PersistentCacheKeys.AFFILIATE_INFO_UPDATE_TIME,
  );
  const lastUpdateTime = persistentCache?.value
    ? DateTime.fromISO(persistentCache.value)
    : undefined;
  return lastUpdateTime;
}
