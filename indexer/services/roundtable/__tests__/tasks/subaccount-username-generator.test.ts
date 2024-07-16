import {
  SubaccountUsernamesTable,
  SubaccountTable,
  QueryableField,

  testMocks,
  dbHelpers,
} from '@dydxprotocol-indexer/postgres';
import subaccountUsernameGenerator from '../../src/tasks/subaccount-username-generator';

describe('subaccount-username-generator', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
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
    const subaccounts = await SubaccountTable.findAll({
      subaccountNumber: 0,
    }, [QueryableField.SUBACCOUNT_NUMBER], {});

    const subaccountsLength = subaccounts.length;
    const subaccountsWithUsernames = await SubaccountUsernamesTable.findAll(
      {}, [], { readReplica: true });
    expect(subaccountsWithUsernames.length).toEqual(0);

    await subaccountUsernameGenerator();
    const subaccountsWithUsernamesAfter = await SubaccountUsernamesTable.findAll(
      {}, [], { readReplica: true });

    expect(subaccountsWithUsernamesAfter.length).toEqual(subaccountsLength);
  });
});
