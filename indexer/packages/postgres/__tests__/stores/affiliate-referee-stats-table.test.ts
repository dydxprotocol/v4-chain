import {
  AffiliateRefereeStatsFromDatabase, Liquidity, FillType,
} from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import {
  defaultOrder,
  defaultAddress4,
  defaultFill,
  defaultAddress2,
  affiliateStatDefaultAddrReferredByAddr2,
  affiliateStatAddr3ReferredByAddr2,
  affiliateStatAddr4ReferredByAddr,
  defaultTendermintEventId,
  defaultTendermintEventId2,
  defaultTendermintEventId3,
  defaultTendermintEventId4,
  // vaultAddress,
} from '../helpers/constants';
import * as AffiliateRefereeStatsTable from '../../src/stores/affiliate-referee-stats-table';
import * as OrderTable from '../../src/stores/order-table';
import * as AffiliateReferredUsersTable from '../../src/stores/affiliate-referred-users-table';
import * as FillTable from '../../src/stores/fill-table';
import { seedData } from '../helpers/mock-generators';
import { DateTime } from 'luxon';

describe('Affiliate info store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates affiliate stats', async () => {
    await AffiliateRefereeStatsTable.create(affiliateStatDefaultAddrReferredByAddr2);
  });

  it('Cannot create duplicate stats for referee', async () => {
    await AffiliateRefereeStatsTable.create(affiliateStatDefaultAddrReferredByAddr2);
    await expect(AffiliateRefereeStatsTable.create(
      affiliateStatDefaultAddrReferredByAddr2)).rejects.toThrowError();
  });

  it('Can upsert referee stats multiple times', async () => {
    await AffiliateRefereeStatsTable.upsert(affiliateStatDefaultAddrReferredByAddr2);
    let info:
    AffiliateRefereeStatsFromDatabase | undefined = await AffiliateRefereeStatsTable.findById(
      affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
    );
    expect(info).toEqual(expect.objectContaining(affiliateStatDefaultAddrReferredByAddr2));

    await AffiliateRefereeStatsTable.upsert(affiliateStatAddr3ReferredByAddr2);
    info = await AffiliateRefereeStatsTable.findById(
      affiliateStatAddr3ReferredByAddr2.refereeAddress,
    );
    expect(info).toEqual(expect.objectContaining(affiliateStatAddr3ReferredByAddr2));
  });

  it('Successfully finds all referee stats', async () => {
    await Promise.all([
      AffiliateRefereeStatsTable.create(affiliateStatDefaultAddrReferredByAddr2),
      AffiliateRefereeStatsTable.create(affiliateStatAddr3ReferredByAddr2),
    ]);

    const infos: AffiliateRefereeStatsFromDatabase[] = await AffiliateRefereeStatsTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(infos.length).toEqual(2);
    expect(infos).toEqual(expect.arrayContaining([
      expect.objectContaining(affiliateStatDefaultAddrReferredByAddr2),
      expect.objectContaining(affiliateStatAddr3ReferredByAddr2),
    ]));
  });

  it('Successfully finds referee stat by id', async () => {
    await AffiliateRefereeStatsTable.create(affiliateStatDefaultAddrReferredByAddr2);

    const info:
    AffiliateRefereeStatsFromDatabase | undefined = await AffiliateRefereeStatsTable.findById(
      affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
    );
    expect(info).toEqual(expect.objectContaining(affiliateStatDefaultAddrReferredByAddr2));
  });

  it('Successfully finds referee stats refereed by one affiliate user', async () => {
    await Promise.all([
      AffiliateRefereeStatsTable.create(affiliateStatDefaultAddrReferredByAddr2),
      AffiliateRefereeStatsTable.create(affiliateStatAddr3ReferredByAddr2),
      AffiliateRefereeStatsTable.create(affiliateStatAddr4ReferredByAddr),
    ]);

    const stats:
    AffiliateRefereeStatsFromDatabase[] = await AffiliateRefereeStatsTable.findAll(
      {
        affiliateAddress: affiliateStatDefaultAddrReferredByAddr2.affiliateAddress,
      },
      [],
      { readReplica: true },
    );

    expect(stats).toHaveLength(2);
    expect(stats).toEqual(expect.arrayContaining([
      expect.objectContaining(affiliateStatDefaultAddrReferredByAddr2),
      expect.objectContaining(affiliateStatAddr3ReferredByAddr2),
    ]));
  });

  it('Returns undefined if affiliate stats not found by pair', async () => {
    await AffiliateRefereeStatsTable.create(affiliateStatDefaultAddrReferredByAddr2);

    const ret = await AffiliateRefereeStatsTable.findById(
      'non_existent_referee_address',
    );
    expect(ret).toBeUndefined();
  });

  describe('updateStats', () => {
    it('Successfully creates new affiliate stats', async () => {
      // Get affiliate info (wallet2 is affiliate)
      const referenceDt: DateTime = await populateFillsAndReferrals();

      const oldStats
      : AffiliateRefereeStatsFromDatabase | undefined = await AffiliateRefereeStatsTable.findById(
        affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
      );

      // No stats yet in DB.
      expect(oldStats).toBeUndefined();

      // Perform update
      await AffiliateRefereeStatsTable.updateStats(
        // exclusive - only includes all fills > 2 minutes ago
        referenceDt.minus({ minutes: 2 }).toISO(),
        referenceDt.toISO(),
      );

      const updatedStats
      : AffiliateRefereeStatsFromDatabase | undefined = await AffiliateRefereeStatsTable.findById(
        affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
      );

      const expectedStats: AffiliateRefereeStatsFromDatabase = {
        affiliateAddress: affiliateStatDefaultAddrReferredByAddr2.affiliateAddress,
        refereeAddress: affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
        affiliateEarnings: '600.6',
        referredMakerTrades: 1,
        referredTakerTrades: 1,
        referredMakerFees: '0',
        referredTakerFees: '888.8',
        referredMakerRebates: '-22.2',
        referredLiquidationFees: '0',
        referralBlockHeight: '1',
        referredTotalVolume: '2.2',
      };

      expect(updatedStats).toEqual(expectedStats);
    });

    it('Successfully updates/increments affiliate info for stats and new referrals', async () => {
      const referenceDt: DateTime = await populateFillsAndReferrals();

      // Perform update: covers first 4 fills during period (-3min, -2min].
      await AffiliateRefereeStatsTable.updateStats(
        referenceDt.minus({ minutes: 3 }).toISO(),
        referenceDt.minus({ minutes: 2 }).toISO(),
      );

      let updatedDefaultAddrStat
      : AffiliateRefereeStatsFromDatabase | undefined = await AffiliateRefereeStatsTable.findById(
        affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
      );
      const expectedDefaultAddrStat1: AffiliateRefereeStatsFromDatabase = {
        refereeAddress: affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
        affiliateAddress: affiliateStatDefaultAddrReferredByAddr2.affiliateAddress,
        affiliateEarnings: '500.5',
        referredMakerTrades: 3,
        referredTakerTrades: 1,
        referredMakerFees: '203.2',
        referredTakerFees: '0',
        referredMakerRebates: '-33.3',
        referredLiquidationFees: '1200.2',
        referralBlockHeight: '1',
        referredTotalVolume: '4.4',
      };
      expect(updatedDefaultAddrStat).toEqual(expectedDefaultAddrStat1);

      // Perform update: covers next 2 fills during period (-2min, -1min].
      await AffiliateRefereeStatsTable.updateStats(
        referenceDt.minus({ minutes: 2 }).toISO(),
        referenceDt.minus({ minutes: 1 }).toISO(),
      );

      updatedDefaultAddrStat = await AffiliateRefereeStatsTable.findById(
        affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
      );
      const expectedDefaultAddrStat2: AffiliateRefereeStatsFromDatabase = {
        refereeAddress: affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
        affiliateAddress: affiliateStatDefaultAddrReferredByAddr2.affiliateAddress,
        affiliateEarnings: '1101.1', // 600.6 + 500.5
        referredMakerTrades: 4,
        referredTakerTrades: 2,
        referredMakerFees: '203.2',
        referredTakerFees: '888.8',
        referredMakerRebates: '-55.5', // -33.3 - 22.2
        referredLiquidationFees: '1200.2',
        referralBlockHeight: '1',
        referredTotalVolume: '6.6',
      };
      expect(updatedDefaultAddrStat).toEqual(expectedDefaultAddrStat2);

      // Perform update: catches no fills but new affiliate referral
      await AffiliateReferredUsersTable.create({
        affiliateAddress: defaultAddress2,
        refereeAddress: defaultAddress4,
        referredAtBlock: '3',
      });
      await AffiliateRefereeStatsTable.updateStats(
        referenceDt.minus({ minutes: 1 }).toISO(),
        referenceDt.toISO(),
      );
      // TODO: update query to find by referee only?
      const updatedDefaultAddr4Stat = await AffiliateRefereeStatsTable.findById(
        defaultAddress4,
      );
      const expectedDefaultAddr4Stat: AffiliateRefereeStatsFromDatabase = {
        refereeAddress: defaultAddress4,
        affiliateAddress: defaultAddress2,
        affiliateEarnings: '0',
        referredMakerTrades: 0,
        referredTakerTrades: 0,
        referredMakerFees: '0',
        referredTakerFees: '0',
        referredMakerRebates: '0',
        referredLiquidationFees: '0',
        referralBlockHeight: '3',
        referredTotalVolume: '0',
      };
      expect(updatedDefaultAddr4Stat).toEqual(expectedDefaultAddr4Stat);
    });

    it('Does not use fills from before referal block height', async () => {
      const referenceDt: DateTime = DateTime.utc();

      await seedData();
      await OrderTable.create(defaultOrder);

      // Referal at block 2 but fill is at block 1
      await AffiliateReferredUsersTable.create({
        affiliateAddress: affiliateStatDefaultAddrReferredByAddr2.affiliateAddress,
        refereeAddress: affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
        referredAtBlock: '2',
      });
      await FillTable.create({
        ...defaultFill,
        liquidity: Liquidity.TAKER,
        subaccountId: defaultOrder.subaccountId,
        createdAt: referenceDt.toISO(),
        createdAtHeight: '1',
        eventId: defaultTendermintEventId,
        price: '1',
        size: '1',
        fee: '1000',
        affiliateRevShare: '500',
      });

      await AffiliateRefereeStatsTable.updateStats(
        referenceDt.minus({ minutes: 1 }).toISO(),
        referenceDt.toISO(),
      );

      const updatedStats
      : AffiliateRefereeStatsFromDatabase | undefined = await AffiliateRefereeStatsTable.findById(
        affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
      );
      // expect one referred user but no fill stats
      const expectedRefereeStats: AffiliateRefereeStatsFromDatabase = {
        refereeAddress: affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
        affiliateAddress: affiliateStatDefaultAddrReferredByAddr2.affiliateAddress,
        affiliateEarnings: '0',
        referredMakerTrades: 0,
        referredTakerTrades: 0,
        referredLiquidationFees: '0',
        referredMakerFees: '0',
        referredTakerFees: '0',
        referredMakerRebates: '0',
        referralBlockHeight: '2',
        referredTotalVolume: '0',
      };
      expect(updatedStats).toEqual(expectedRefereeStats);
    });
  });

  // TODO(CT-1341): Add paginated query tests similar to `affiliate-info-table.test.ts`
});

async function populateFillsAndReferrals(): Promise<DateTime> {
  const referenceDt = DateTime.utc();
  const referredAtBlock = 1;

  await seedData();

  // default address 2 refers default address
  await AffiliateReferredUsersTable.create({
    affiliateAddress: affiliateStatDefaultAddrReferredByAddr2.affiliateAddress,
    refereeAddress: affiliateStatDefaultAddrReferredByAddr2.refereeAddress,
    referredAtBlock: `${referredAtBlock}`,
  });

  await AffiliateReferredUsersTable.create({
    affiliateAddress: affiliateStatAddr3ReferredByAddr2.affiliateAddress,
    refereeAddress: affiliateStatAddr3ReferredByAddr2.refereeAddress,
    referredAtBlock: `${referredAtBlock}`,
  });
  // Create order and fils for defaultWallet (referee)
  await OrderTable.create(defaultOrder);

  await Promise.all([
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.TAKER,
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 1 }).toISO(),
      createdAtHeight: `${referredAtBlock}`,
      eventId: defaultTendermintEventId,
      price: '1',
      size: '1.1',
      fee: '888.8',
      affiliateRevShare: '600.6',
    }),
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.MAKER,
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 1 }).toISO(),
      createdAtHeight: `${referredAtBlock}`,
      eventId: defaultTendermintEventId2,
      price: '1',
      size: '1.1',
      fee: '-22.2',
      affiliateRevShare: '0',
    }),
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.MAKER, // use uneven number of maker/taker
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 2 }).toISO(),
      createdAtHeight: `${referredAtBlock + 1}`,
      eventId: defaultTendermintEventId3,
      price: '1',
      size: '1.1',
      fee: '99.9',
      affiliateRevShare: '0',
    }),
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.MAKER,
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 2 }).toISO(),
      createdAtHeight: `${referredAtBlock + 1}`,
      eventId: defaultTendermintEventId4,
      price: '1',
      size: '1.1',
      fee: '103.3',
      affiliateRevShare: '0',
    }),
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.TAKER,
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 2 }).toISO(),
      createdAtHeight: `${referredAtBlock + 1}`,
      eventId: defaultTendermintEventId4,
      price: '1',
      size: '1.1',
      fee: '1200.2',
      affiliateRevShare: '500.5',
      type: FillType.LIQUIDATED,
    }),
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.MAKER,
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 2 }).toISO(),
      createdAtHeight: `${referredAtBlock + 1}`,
      eventId: defaultTendermintEventId,
      price: '1',
      size: '1.1',
      fee: '-33.3',
      affiliateRevShare: '0',
      type: FillType.LIQUIDATION,
    }),
  ]);

  return referenceDt;
}
