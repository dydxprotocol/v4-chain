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
    const subaccountsWithUsernames = await SubaccountUsernamesTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    expect(subaccountsWithUsernames.length).toEqual(2);

    await subaccountUsernameGenerator();
    const subaccountsWithUsernamesAfter = await SubaccountUsernamesTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    const expectedUsernames = [
      'BubblyEarH5Y', // dydx1n88uc38xhjgxzw9nwre4ep2c8ga4fjxc575lnf
      'GreenSnowWTT', // dydx1n88uc38xhjgxzw9nwre4ep2c8ga4fjxc565lnf
      'LunarMatFK5', // dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4
    ];
    expect(subaccountsWithUsernamesAfter.length).toEqual(subaccountsLength);
    for (let i = 0; i < expectedUsernames.length; i++) {
      expect(subaccountsWithUsernamesAfter[i].username).toEqual(
        expectedUsernames[i],
      );
    }
  });
});
