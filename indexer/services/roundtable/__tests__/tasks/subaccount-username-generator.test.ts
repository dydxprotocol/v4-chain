import {
  SubaccountUsernamesTable,
  SubaccountTable,
  QueryableField,
  testMocks,
  dbHelpers,
  SubaccountFromDatabase,
  testConstants,
} from '@dydxprotocol-indexer/postgres';
import subaccountUsernameGenerator from '../../src/tasks/subaccount-username-generator';

describe('subaccount-username-generator', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await testMocks.seedAdditionalSubaccounts();
    // delete all usernames that were seeded
    await SubaccountUsernamesTable.deleteBySubaccountId(
      testConstants.defaultSubaccountId,
    );
    await SubaccountUsernamesTable.deleteBySubaccountId(
      testConstants.defaultSubaccountId2,
    );
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  it('Successfully creates a username for all subaccount', async () => {
    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        subaccountNumber: 0,
      },
      [QueryableField.SUBACCOUNT_NUMBER],
      {},
    );

    const subaccountsLength: number = subaccounts.length;
    const before = await SubaccountUsernamesTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(before.length).toEqual(0);

    await subaccountUsernameGenerator();
    const after = await SubaccountUsernamesTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    const expectedUsernames = [
      'BubblyEarH5Y', // dydx1n88uc38xhjgxzw9nwre4ep2c8ga4fjxc575lnf
      'GreenSnowWTT', // dydx1n88uc38xhjgxzw9nwre4ep2c8ga4fjxc565lnf
      'LunarMatFK5', // dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4
    ];
    expect(after.length).toEqual(subaccountsLength);
    for (let i = 0; i < expectedUsernames.length; i++) {
      expect(after[i].username).toEqual(expectedUsernames[i]);
    }
  });

  it('Falls back to a second username when there is a conflict on the first attempt', async () => {
    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        subaccountNumber: 0,
      },
      [QueryableField.SUBACCOUNT_NUMBER],
      {},
    );
    const targetSubaccount = subaccounts[0];
    const otherSubaccount = subaccounts[1];

    const { generateUsernameForSubaccount } = require('../../src/helpers/usernames-helper');

    const usernameAttempt0 = generateUsernameForSubaccount(
      targetSubaccount.address,
      0,
      0,
    );
    await SubaccountUsernamesTable.create({
      username: usernameAttempt0,
      subaccountId: otherSubaccount.id,
    });

    const afterPreInsert = await SubaccountUsernamesTable.findAll(
      { subaccountId: [targetSubaccount.id] }, [QueryableField.SUBACCOUNT_ID], {},
    );
    expect(afterPreInsert.length).toBe(0);

    await subaccountUsernameGenerator();

    const created = await SubaccountUsernamesTable.findAll(
      { subaccountId: [targetSubaccount.id] }, [QueryableField.SUBACCOUNT_ID], {},
    );
    expect(created.length).toBe(1);

    const fallbackUsername = generateUsernameForSubaccount(
      targetSubaccount.address,
      0,
      1,
    );
    expect(created[0].username).toEqual(fallbackUsername);

    const conflict = await SubaccountUsernamesTable.findAll(
      { username: [usernameAttempt0] }, [QueryableField.USERNAME], {},
    );
    expect(conflict.length).toBe(1);
    expect(conflict[0].subaccountId).not.toEqual(targetSubaccount.id);
  });

  it('Handles batch where one username succeeds and the other needs a fallback', async () => {
    const subaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      {
        subaccountNumber: 0,
      },
      [QueryableField.SUBACCOUNT_NUMBER],
      {},
    );
    expect(subaccounts.length).toBeGreaterThanOrEqual(2);

    const sub0 = subaccounts[0];
    const sub1 = subaccounts[1];

    const { generateUsernameForSubaccount } = require('../../src/helpers/usernames-helper');

    const sub0Attempt0 = generateUsernameForSubaccount(sub0.address, 0, 0);
    await SubaccountUsernamesTable.create({
      username: sub0Attempt0,
      subaccountId: sub1.id,
    });

    // pre-run checks
    const preUsernames = await SubaccountUsernamesTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(
      preUsernames.find((u: any) => u.username === sub0Attempt0),
    ).toBeDefined();
    expect(preUsernames.filter((u: any) => u.subaccountId === sub0.id).length).toBe(0);
    expect(preUsernames.filter(
      (u: any) => u.subaccountId === sub1.id && u.username !== sub0Attempt0,
    ).length).toBe(0);

    // run generator
    await subaccountUsernameGenerator();

    // fetch results
    const after = await SubaccountUsernamesTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    // sub0 should have fallback username
    const sub0UsernameRow = after.find((u: any) => u.subaccountId === sub0.id);
    const sub0ExpectedFallback = generateUsernameForSubaccount(sub0.address, 0, 1);
    expect(sub0UsernameRow).toBeDefined();
    if (sub0UsernameRow) {
      expect(sub0UsernameRow.username).toEqual(sub0ExpectedFallback);
    }
    const sub1UsernameRow = after.find(
      (u: any) => u.subaccountId === sub1.id && u.username !== sub0Attempt0);
    if (sub1UsernameRow) {
      expect(sub1UsernameRow).toBeDefined();
    }

    // There should not be two usernames with the same value
    const usernameCounts = after.reduce((acc: Record<string, number>, u: any) => {
      acc[u.username] = (acc[u.username] || 0) + 1;
      return acc;
    }, {});
    for (const count of Object.values(usernameCounts)) {
      expect(count).toBe(1);
    }
  });
});
