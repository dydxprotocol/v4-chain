import { SubaccountFromDatabase, SubaccountUsernamesFromDatabase, SubaccountsWithoutUsernamesResult } from '../../src/types';
import * as SubaccountUsernamesTable from '../../src/stores/subaccount-usernames-table';
import * as WalletTable from '../../src/stores/wallet-table';
import * as SubaccountsTable from '../../src/stores/subaccount-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import {
  defaultSubaccountUsername,
  defaultSubaccountUsername2,
  defaultWallet,
  defaultWallet2,
  duplicatedSubaccountUsername,
  subaccountUsernameWithAlternativeAddress,
} from '../helpers/constants';
import { seedData, seedAdditionalSubaccounts } from '../helpers/mock-generators';

describe('SubaccountUsernames store', () => {
  beforeEach(async () => {
    await seedData();
    await seedAdditionalSubaccounts();
  });

  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully finds all SubaccountUsernames', async () => {
    const subaccountUsernames:
    SubaccountUsernamesFromDatabase[] = await SubaccountUsernamesTable.findAll(
      {},
      [],
      {},
    );

    expect(subaccountUsernames.length).toEqual(2);
    expect(subaccountUsernames[0]).toEqual(expect.objectContaining(defaultSubaccountUsername));
    expect(subaccountUsernames[1]).toEqual(expect.objectContaining(defaultSubaccountUsername2));
  });

  it('Successfully finds SubaccountUsername with subaccountId', async () => {
    const subaccountUsername:
    SubaccountUsernamesFromDatabase | undefined = await SubaccountUsernamesTable.findByUsername(
      defaultSubaccountUsername.username,
    );
    expect(subaccountUsername).toEqual(expect.objectContaining(defaultSubaccountUsername));
  });

  it('Duplicate SubaccountUsername creation fails', async () => {
    await expect(SubaccountUsernamesTable.create(duplicatedSubaccountUsername)).rejects.toThrow();
  });

  it('Creation of row without subaccountId fails', async () => {
    await expect(SubaccountUsernamesTable.create({ ...defaultSubaccountUsername, subaccountId: '' })).rejects.toThrow();
  });

  it('Get subaccount ids which arent in the subaccount usernames table', async () => {
    const subaccounts: SubaccountFromDatabase[] = await SubaccountsTable.findAll({
      subaccountNumber: 0,
    }, [], {});
    const subaccountLength = subaccounts.length;
    const subaccountIds: SubaccountsWithoutUsernamesResult[] = await
    SubaccountUsernamesTable.getSubaccountZerosWithoutUsernames(1000);
    expect(subaccountIds.length).toEqual(subaccountLength - 1);
  });

  it('Get username using address', async () => {
    await Promise.all([
      // Add one username for alternativeWallet
      WalletTable.create(defaultWallet2),
      SubaccountUsernamesTable.create(subaccountUsernameWithAlternativeAddress),
    ]);

    // Should only get username for defaultWallet's subaccount 0
    const usernames = await SubaccountUsernamesTable.findByAddress([defaultWallet.address]);
    expect(usernames.length).toEqual(1);
    expect(usernames[0]).toEqual(expect.objectContaining({
      address: defaultWallet.address,
      username: defaultSubaccountUsername.username,
    }));
  });
});
