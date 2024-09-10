import { WalletFromDatabase, PersistentCacheKeys } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { DateTime } from 'luxon';
import {
  defaultFill,
  defaultOrder,
  defaultSubaccount,
  defaultTendermintEventId,
  defaultTendermintEventId2,
  defaultTendermintEventId3,
  defaultTendermintEventId4,
  defaultWallet,
  defaultWallet2,
  defaultWallet3,
  isolatedMarketOrder,
  isolatedSubaccount,
} from '../helpers/constants';
import * as FillTable from '../../src/stores/fill-table';
import * as OrderTable from '../../src/stores/order-table';
import * as WalletTable from '../../src/stores/wallet-table';
import * as SubaccountTable from '../../src/stores/subaccount-table';
import * as PersistentCacheTable from '../../src/stores/persistent-cache-table';
import { seedData } from '../helpers/mock-generators';

describe('Wallet store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a Wallet', async () => {
    await WalletTable.create(defaultWallet2);
  });

  it('Successfully upserts a Wallet multiple times', async () => {
    await WalletTable.upsert(defaultWallet2);
    let wallet: WalletFromDatabase | undefined = await WalletTable.findById(
      defaultWallet2.address,
    );

    expect(wallet).toEqual(expect.objectContaining(defaultWallet2));
    await WalletTable.upsert({
      ...defaultWallet2,
      isWhitelistAffiliate: true,
      totalVolume: '100.1',
    });
    wallet = await WalletTable.findById(defaultWallet2.address);

    expect(wallet).toEqual(expect.objectContaining({
      ...defaultWallet2,
      isWhitelistAffiliate: true,
      totalVolume: '100.1',
    }));
  });

  it('Successfully finds all Wallets', async () => {
    await Promise.all([
      WalletTable.create(defaultWallet2),
      WalletTable.create({
        address: 'fake_address',
        totalTradingRewards: '0',
        totalVolume: '0',
        isWhitelistAffiliate: false,
      }),
    ]);

    const wallets: WalletFromDatabase[] = await WalletTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(wallets.length).toEqual(2);
    expect(wallets[0]).toEqual(expect.objectContaining(defaultWallet2));
    expect(wallets[1]).toEqual(expect.objectContaining({
      address: 'fake_address',
    }));
  });

  it('Successfully finds a Wallet', async () => {
    await WalletTable.create(defaultWallet2);

    const wallet: WalletFromDatabase | undefined = await WalletTable.findById(
      defaultWallet2.address,
    );

    expect(wallet).toEqual(expect.objectContaining(defaultWallet2));
  });

  it('Successfully finds wallets by whitelist flag', async () => {
    await Promise.all([
      WalletTable.create(defaultWallet3),
      WalletTable.create(defaultWallet2),
    ]);

    const wallets: WalletFromDatabase[] = await WalletTable.findAll(
      {
        isWhitelistAffiliate: true,
      },
      [],
      { readReplica: true },
    );

    expect(wallets.length).toEqual(1);
    expect(wallets[0]).toEqual(expect.objectContaining(defaultWallet3));
  });

  describe('Wallet .updateTotalVolume()', () => {
    it('Successfully updates totalVolume for time window multiple times', async () => {
      const firstFillTime = await populateWalletSubaccountFill();

      // Update totalVolume for a time window that covers all fills
      await WalletTable.updateTotalVolume(
        firstFillTime.minus({ hours: 1 }).toISO(), // need to minus because left bound is exclusive
        firstFillTime.plus({ hours: 1 }).toISO(),
      );
      let wallet = await WalletTable.findById(defaultWallet.address);
      expect(wallet).toEqual(expect.objectContaining({
        ...defaultWallet,
        totalVolume: '103',
      }));

      // Update totalVolume for a time window that excludes some fills
      // For convenience, we will reuse the existing fills data. The total volume calculated in this
      // window should be added to the total volume above.
      await WalletTable.updateTotalVolume(
        firstFillTime.toISO(), // exclusive -> filters out first fill from each subaccount
        firstFillTime.plus({ minutes: 2 }).toISO(),
      );
      wallet = await WalletTable.findById(defaultWallet.address);
      expect(wallet).toEqual(expect.objectContaining({
        ...defaultWallet,
        totalVolume: '105', // 103 + 2
      }));
    });

    it('Successfully upserts persistent cache', async () => {
      const referenceDt = DateTime.utc();

      // Sets initial persistent cache value
      let leftBound = referenceDt.minus({ hours: 2 });
      let rightBound = referenceDt.minus({ hours: 1 });

      await WalletTable.updateTotalVolume(leftBound.toISO(), rightBound.toISO());

      let persistentCache = await PersistentCacheTable.findById(
        PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME,
      );
      let lastUpdateTime = persistentCache?.value;
      expect(lastUpdateTime).not.toBeUndefined();
      if (lastUpdateTime !== undefined) {
        expect(lastUpdateTime).toEqual(rightBound.toISO());
      }

      // Updates persistent cache value
      leftBound = referenceDt.minus({ hours: 1 });
      rightBound = referenceDt;

      await WalletTable.updateTotalVolume(leftBound.toISO(), rightBound.toISO());

      persistentCache = await PersistentCacheTable.findById(
        PersistentCacheKeys.TOTAL_VOLUME_UPDATE_TIME,
      );
      lastUpdateTime = persistentCache?.value;
      expect(lastUpdateTime).not.toBeUndefined();
      if (lastUpdateTime !== undefined) {
        expect(lastUpdateTime).toEqual(rightBound.toISO());
      }
    });
  });
});

/**
 * Helper function to add entries into wallet, subaccount, fill tables.
 * Create a wallet with 2 subaccounts; one subaccount has 3 fills and the other has 1 fill.
 * The fills are at t=0,1,2 and t=1 for the subaccounts respectively.
 * This setup allows us to test that the totalVolume is correctly calculated for a time window.
 * @returns first fill time in ISO format
 */
async function populateWalletSubaccountFill(): Promise<DateTime> {
  await seedData();
  await OrderTable.create(defaultOrder);
  await OrderTable.create(isolatedMarketOrder);

  // seedData() creates defaultWallet with defaultSubaccount and isolatedSubaccount
  const defaultSubaccountId = await SubaccountTable.findAll(
    { subaccountNumber: defaultSubaccount.subaccountNumber },
    [],
    {},
  );
  const isolatedSubaccountId = await SubaccountTable.findAll(
    { subaccountNumber: isolatedSubaccount.subaccountNumber },
    [],
    {},
  );

  const referenceDt = DateTime.utc().minus({ hours: 1 });
  const eventIds = [
    defaultTendermintEventId,
    defaultTendermintEventId2,
    defaultTendermintEventId3,
    defaultTendermintEventId4,
  ];
  let eventIdx = 0;

  // Create 3 fills with 1 min increments for defaultSubaccount
  for (let i = 0; i < 3; i++) {
    await FillTable.create({
      ...defaultFill,
      subaccountId: defaultSubaccountId[0].id,
      createdAt: referenceDt.plus({ minutes: i }).toISO(),
      eventId: eventIds[eventIdx],
      price: '1',
      size: '1',
    });
    eventIdx += 1;
  }
  // Create 1 fill at referenceDt for isolatedSubaccount
  await FillTable.create({
    ...defaultFill,
    subaccountId: isolatedSubaccountId[0].id,
    createdAt: referenceDt.toISO(),
    eventId: eventIds[eventIdx],
    price: '10',
    size: '10',
  });

  return referenceDt;
}
