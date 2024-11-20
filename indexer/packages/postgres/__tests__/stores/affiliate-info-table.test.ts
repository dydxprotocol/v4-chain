import {
  AffiliateInfoFromDatabase, Liquidity, FillType,
} from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import {
  defaultOrder,
  defaultWallet,
  defaultFill,
  defaultWallet2,
  defaultAffiliateInfo,
  defaultAffiliateInfo2,
  defaultTendermintEventId,
  defaultTendermintEventId2,
  defaultTendermintEventId3,
  defaultTendermintEventId4,
  vaultAddress,
} from '../helpers/constants';
import * as AffiliateInfoTable from '../../src/stores/affiliate-info-table';
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

  it('Successfully creates affiliate info', async () => {
    await AffiliateInfoTable.create(defaultAffiliateInfo);
  });

  it('Cannot create duplicate info for duplicate address', async () => {
    await AffiliateInfoTable.create(defaultAffiliateInfo);
    await expect(AffiliateInfoTable.create(defaultAffiliateInfo)).rejects.toThrowError();
  });

  it('Can upsert affiliate info multiple times', async () => {
    await AffiliateInfoTable.upsert(defaultAffiliateInfo);
    let info: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
      defaultAffiliateInfo.address,
    );
    expect(info).toEqual(expect.objectContaining(defaultAffiliateInfo));

    await AffiliateInfoTable.upsert(defaultAffiliateInfo2);
    info = await AffiliateInfoTable.findById(defaultAffiliateInfo2.address);
    expect(info).toEqual(expect.objectContaining(defaultAffiliateInfo2));
  });

  it('Successfully finds all affiliate infos', async () => {
    await Promise.all([
      AffiliateInfoTable.create(defaultAffiliateInfo),
      AffiliateInfoTable.create(defaultAffiliateInfo2),
    ]);

    const infos: AffiliateInfoFromDatabase[] = await AffiliateInfoTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(infos.length).toEqual(2);
    expect(infos).toEqual(expect.arrayContaining([
      expect.objectContaining(defaultAffiliateInfo),
      expect.objectContaining(defaultAffiliateInfo2),
    ]));
  });

  it('Successfully finds affiliate info by Id', async () => {
    await AffiliateInfoTable.create(defaultAffiliateInfo);

    const info: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
      defaultAffiliateInfo.address,
    );
    expect(info).toEqual(expect.objectContaining(defaultAffiliateInfo));
  });

  it('Returns undefined if affiliate info not found by Id', async () => {
    await AffiliateInfoTable.create(defaultAffiliateInfo);

    const info: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
      'non_existent_address',
    );
    expect(info).toBeUndefined();
  });

  describe('updateInfo', () => {
    it('Successfully creates new affiliate info', async () => {
      const referenceDt: DateTime = await populateFillsAndReferrals();

      // Perform update
      await AffiliateInfoTable.updateInfo(
        referenceDt.minus({ minutes: 2 }).toISO(),
        referenceDt.toISO(),
      );

      // Get affiliate info (wallet2 is affiliate)
      const updatedInfo: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
        defaultWallet2.address,
      );

      const expectedAffiliateInfo: AffiliateInfoFromDatabase = {
        address: defaultWallet2.address,
        affiliateEarnings: '1000',
        referredMakerTrades: 1,
        referredTakerTrades: 1,
        totalReferredMakerFees: '0',
        totalReferredTakerFees: '1000',
        totalReferredMakerRebates: '-1000',
        totalReferredUsers: 1,
        firstReferralBlockHeight: '1',
        referredTotalVolume: '2',
      };

      expect(updatedInfo).toEqual(expectedAffiliateInfo);
    });

    it('Successfully updates/increments affiliate info for stats and new referrals', async () => {
      const referenceDt: DateTime = await populateFillsAndReferrals();

      // Perform update: catches first 2 fills
      await AffiliateInfoTable.updateInfo(
        referenceDt.minus({ minutes: 3 }).toISO(),
        referenceDt.minus({ minutes: 2 }).toISO(),
      );

      const updatedInfo1: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
        defaultWallet2.address,
      );
      const expectedAffiliateInfo1: AffiliateInfoFromDatabase = {
        address: defaultWallet2.address,
        affiliateEarnings: '1005',
        referredMakerTrades: 3,
        referredTakerTrades: 1,
        totalReferredMakerFees: '2100',
        totalReferredTakerFees: '0',
        totalReferredMakerRebates: '0',
        totalReferredUsers: 1,
        firstReferralBlockHeight: '1',
        referredTotalVolume: '4',
      };
      expect(updatedInfo1).toEqual(expectedAffiliateInfo1);

      // Perform update: catches next 2 fills
      await AffiliateInfoTable.updateInfo(
        referenceDt.minus({ minutes: 2 }).toISO(),
        referenceDt.minus({ minutes: 1 }).toISO(),
      );

      const updatedInfo2 = await AffiliateInfoTable.findById(
        defaultWallet2.address,
      );
      const expectedAffiliateInfo2: AffiliateInfoFromDatabase = {
        address: defaultWallet2.address,
        affiliateEarnings: '2005',
        referredMakerTrades: 4,
        referredTakerTrades: 2,
        totalReferredMakerFees: '2100',
        totalReferredTakerFees: '1000',
        totalReferredMakerRebates: '-1000',
        totalReferredUsers: 1,
        firstReferralBlockHeight: '1',
        referredTotalVolume: '6',
      };
      expect(updatedInfo2).toEqual(expectedAffiliateInfo2);

      // Perform update: catches no fills but new affiliate referral
      await AffiliateReferredUsersTable.create({
        affiliateAddress: defaultWallet2.address,
        refereeAddress: vaultAddress,
        referredAtBlock: '2',
      });
      await AffiliateInfoTable.updateInfo(
        referenceDt.minus({ minutes: 1 }).toISO(),
        referenceDt.toISO(),
      );
      const updatedInfo3 = await AffiliateInfoTable.findById(
        defaultWallet2.address,
      );
      const expectedAffiliateInfo3: AffiliateInfoFromDatabase = {
        address: defaultWallet2.address,
        affiliateEarnings: '2005',
        referredMakerTrades: 4,
        referredTakerTrades: 2,
        totalReferredMakerFees: '2100',
        totalReferredTakerFees: '1000',
        totalReferredMakerRebates: '-1000',
        totalReferredUsers: 2,
        firstReferralBlockHeight: '1',
        referredTotalVolume: '6',
      };
      expect(updatedInfo3).toEqual(expectedAffiliateInfo3);
    });

    it('Does not use fills from before referal block height', async () => {
      const referenceDt: DateTime = DateTime.utc();

      await seedData();
      await OrderTable.create(defaultOrder);

      // Referal at block 2 but fill is at block 1
      await AffiliateReferredUsersTable.create({
        affiliateAddress: defaultWallet2.address,
        refereeAddress: defaultWallet.address,
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

      await AffiliateInfoTable.updateInfo(
        referenceDt.minus({ minutes: 1 }).toISO(),
        referenceDt.toISO(),
      );

      const updatedInfo: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
        defaultWallet2.address,
      );
      // expect one referred user but no fill stats
      const expectedAffiliateInfo: AffiliateInfoFromDatabase = {
        address: defaultWallet2.address,
        affiliateEarnings: '0',
        referredMakerTrades: 0,
        referredTakerTrades: 0,
        totalReferredMakerFees: '0',
        totalReferredTakerFees: '0',
        totalReferredMakerRebates: '0',
        totalReferredUsers: 1,
        firstReferralBlockHeight: '2',
        referredTotalVolume: '0',
      };
      expect(updatedInfo).toEqual(expectedAffiliateInfo);
    });
  });

  describe('paginatedFindWithAddressFilter', () => {
    beforeEach(async () => {
      await migrate();
      await Promise.all(
        Array.from({ length: 10 }, (_, i) => AffiliateInfoTable.create({
          ...defaultAffiliateInfo,
          address: `address_${i}`,
          affiliateEarnings: i.toString(),
        }),
        ),
      );
    });

    it('Successfully filters by address', async () => {
      const infos: AffiliateInfoFromDatabase[] = await AffiliateInfoTable
        .paginatedFindWithAddressFilter(
          ['address_0'],
          0,
          10,
          false,
        );
      expect(infos).toBeDefined();
      expect(infos!.length).toEqual(1);
      expect(infos![0]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_0',
        affiliateEarnings: '0',
      }));
    });

    it('Successfully sorts by affiliate earning', async () => {
      const infos: AffiliateInfoFromDatabase[] = await AffiliateInfoTable
        .paginatedFindWithAddressFilter(
          [],
          0,
          10,
          true,
        );
      expect(infos).toBeDefined();
      expect(infos!.length).toEqual(10);
      expect(infos![0]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_9',
        affiliateEarnings: '9',
      }));
      expect(infos![9]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_0',
        affiliateEarnings: '0',
      }));
    });

    it('Successfully uses offset (default to sorted) and limit', async () => {
      const infos: AffiliateInfoFromDatabase[] = await AffiliateInfoTable
        .paginatedFindWithAddressFilter(
          [],
          5,
          2,
          false,
        );
      expect(infos).toBeDefined();
      expect(infos!.length).toEqual(2);
      expect(infos![0]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_4',
        // affiliateEarnings in DB: 9, 8, 7, 6, 5, 4, ...
        // so we get 4 with offset = 5.
        affiliateEarnings: '4',
      }));
      expect(infos![1]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_3',
        affiliateEarnings: '3',
      }));
    });

    it('Successfully filters, sorts, offsets, and limits', async () => {
      const infos: AffiliateInfoFromDatabase[] = await AffiliateInfoTable
        .paginatedFindWithAddressFilter(
          [],
          3,
          2,
          true,
        );
      expect(infos).toBeDefined();
      expect(infos!.length).toEqual(2);
      expect(infos![0]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_6',
        affiliateEarnings: '6',
      }));
      expect(infos![1]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_5',
        affiliateEarnings: '5',
      }));
    });

    it('Returns empty array if no results', async () => {
      const infos: AffiliateInfoFromDatabase[] = await AffiliateInfoTable
        .paginatedFindWithAddressFilter(
          ['address_11'],
          0,
          10,
          false,
        );
      expect(infos).toBeDefined();
      expect(infos!.length).toEqual(0);
    });

    it('Successfully use sorted - equal earnings between affiliates', async () => {
      await AffiliateInfoTable.create({
        ...defaultAffiliateInfo,
        address: 'address_10',
        affiliateEarnings: '9', // same as address_9
      });
      const infos: AffiliateInfoFromDatabase[] = await AffiliateInfoTable
        .paginatedFindWithAddressFilter(
          [],
          0,
          100,
          true,
        );
      expect(infos).toBeDefined();
      expect(infos!.length).toEqual(11);
      expect(infos![0]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_10', // '10' < '9' in lexicographical order
        affiliateEarnings: '9',
      }));
      expect(infos![1]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_9',
        affiliateEarnings: '9',
      }));
    });

  });
});

async function populateFillsAndReferrals(): Promise<DateTime> {
  const referenceDt = DateTime.utc();

  await seedData();

  // defaultWallet2 will be affiliate and defaultWallet will be referee
  await AffiliateReferredUsersTable.create({
    affiliateAddress: defaultWallet2.address,
    refereeAddress: defaultWallet.address,
    referredAtBlock: '1',
  });

  // Create order and fils for defaultWallet (referee)
  await OrderTable.create(defaultOrder);

  await Promise.all([
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.TAKER,
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 1 }).toISO(),
      eventId: defaultTendermintEventId,
      price: '1',
      size: '1',
      fee: '1000',
      affiliateRevShare: '500',
    }),
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.MAKER,
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 1 }).toISO(),
      eventId: defaultTendermintEventId2,
      price: '1',
      size: '1',
      fee: '-1000',
      affiliateRevShare: '500',
    }),
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.MAKER, // use uneven number of maker/taker
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 2 }).toISO(),
      eventId: defaultTendermintEventId3,
      price: '1',
      size: '1',
      fee: '1000',
      affiliateRevShare: '500',
    }),
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.MAKER,
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 2 }).toISO(),
      eventId: defaultTendermintEventId4,
      price: '1',
      size: '1',
      fee: '1000',
      affiliateRevShare: '500',
    }),
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.TAKER,
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 2 }).toISO(),
      eventId: defaultTendermintEventId4,
      price: '1',
      size: '1',
      fee: '1000',
      affiliateRevShare: '0',
      type: FillType.LIQUIDATED,
    }),
    FillTable.create({
      ...defaultFill,
      liquidity: Liquidity.MAKER,
      subaccountId: defaultOrder.subaccountId,
      createdAt: referenceDt.minus({ minutes: 2 }).toISO(),
      eventId: defaultTendermintEventId,
      price: '1',
      size: '1',
      fee: '100',
      affiliateRevShare: '5',
      type: FillType.LIQUIDATION,
    }),
  ]);

  return referenceDt;
}
