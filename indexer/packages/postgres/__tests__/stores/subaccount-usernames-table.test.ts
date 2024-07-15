import { SubaccountUsernamesFromDatabase } from '../../src/types';
import * as SubaccountUsernamesTable from '../../src/stores/subaccount-usernames-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import {
  defaultSubaccountUsername,
  defaultSubaccountUsername2,
  duplicatedSubaccountUsername,
} from '../helpers/constants';
import { seedData } from '../helpers/mock-generators';

describe('SubaccountUsernames store', () => {
  beforeEach(async () => {
    await seedData();
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

  it('Successfully creates a SubaccountUsername', async () => {
    await SubaccountUsernamesTable.create(defaultSubaccountUsername);
  });

  it('Successfully finds all SubaccountUsernames', async () => {
    await Promise.all([
      SubaccountUsernamesTable.create(defaultSubaccountUsername),
      SubaccountUsernamesTable.create(defaultSubaccountUsername2),
    ]);

    const subaccountUsernames:
    SubaccountUsernamesFromDatabase[] = await SubaccountUsernamesTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(subaccountUsernames.length).toEqual(2);
    expect(subaccountUsernames[0]).toEqual(expect.objectContaining(defaultSubaccountUsername));
    expect(subaccountUsernames[1]).toEqual(expect.objectContaining(defaultSubaccountUsername2));
  });

  it('Successfully finds SubaccountUsername with subaccountId', async () => {
    await Promise.all([
      SubaccountUsernamesTable.create(defaultSubaccountUsername),
      SubaccountUsernamesTable.create(defaultSubaccountUsername2),
    ]);

    const subaccountUsername:
    SubaccountUsernamesFromDatabase | undefined = await SubaccountUsernamesTable.findByUsername(
      defaultSubaccountUsername.username,
    );
    expect(subaccountUsername).toEqual(expect.objectContaining(defaultSubaccountUsername));
  });

  it('Duplicate SubaccountUsername creation fails', async () => {
    await SubaccountUsernamesTable.create(defaultSubaccountUsername);
    await expect(SubaccountUsernamesTable.create(duplicatedSubaccountUsername)).rejects.toThrow();
  });
});
