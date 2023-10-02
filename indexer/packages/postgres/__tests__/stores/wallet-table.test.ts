import { WalletFromDatabase } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { blockedAddress, defaultAddress, defaultWallet } from '../helpers/constants';
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
    await WalletTable.create(defaultWallet);
  });

  it('Successfully upserts a Wallet multiple times', async () => {
    await WalletTable.upsert(defaultWallet);
    let wallet: WalletFromDatabase | undefined = await WalletTable.findById(
      defaultAddress,
    );

    expect(wallet).toEqual(expect.objectContaining(defaultWallet));
    await WalletTable.upsert(defaultWallet);
    wallet = await WalletTable.findById(
      defaultAddress,
    );

    expect(wallet).toEqual(expect.objectContaining(defaultWallet));
  });

  it('Successfully finds all Wallets', async () => {
    await Promise.all([
      WalletTable.create(defaultWallet),
      WalletTable.create({
        address: 'fake_address',
      }),
    ]);

    const wallets: WalletFromDatabase[] = await WalletTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(wallets.length).toEqual(2);
    expect(wallets[0]).toEqual(expect.objectContaining(defaultWallet));
    expect(wallets[1]).toEqual(expect.objectContaining({
      address: 'fake_address',
    }));
  });

  it('Successfully finds a Wallet', async () => {
    await WalletTable.create(defaultWallet);

    const wallet: WalletFromDatabase | undefined = await WalletTable.findById(
      defaultAddress,
    );

    expect(wallet).toEqual(expect.objectContaining(defaultWallet));
  });

  describe('bulkCreate', () => {
    it('Successfully creates multiple wallets', async () => {
      const createdWallets:
      WalletFromDatabase[] = await WalletTable.bulkCreate([
        defaultWallet,
        {
          address: blockedAddress,
        },
      ]);

      expect(createdWallets).toHaveLength(2);
      for (let i = 0; i < createdWallets.length; i += 1) {
        const wallet: WalletFromDatabase = createdWallets[i];
        expect(
          await WalletTable.findById(wallet.address),
        ).toEqual(wallet);
      }
    });
  });
});
