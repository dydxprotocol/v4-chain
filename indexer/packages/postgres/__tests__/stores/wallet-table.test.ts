import { WalletFromDatabase } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { defaultWallet2, defaultWallet3 } from '../helpers/constants';
import * as WalletTable from '../../src/stores/wallet-table';

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
      totalVolume: '100',
    });
    wallet = await WalletTable.findById(defaultWallet2.address);

    expect(wallet).toEqual(expect.objectContaining({
      ...defaultWallet2,
      isWhitelistAffiliate: true,
      totalVolume: '100',
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
      WalletTable.create(defaultWallet2)
    ])

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
});
