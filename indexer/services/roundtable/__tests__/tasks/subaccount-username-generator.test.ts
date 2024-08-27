import {
  SubaccountUsernamesTable,
  SubaccountTable,
  QueryableField,

  testMocks,
  dbHelpers,
  SubaccountUsernamesFromDatabase,
  SubaccountFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import subaccountUsernameGenerator from '../../src/tasks/subaccount-username-generator';
import { generateUsername } from '../../src/helpers/usernames-helper';

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
    const subaccounts: SubaccountFromDatabase[] = await
    SubaccountTable.findAll({
      subaccountNumber: 0,
    }, [QueryableField.SUBACCOUNT_NUMBER], {});

    const subaccountsLength: number = subaccounts.length;
    const subaccountsWithUsernames: SubaccountUsernamesFromDatabase[] = await
    SubaccountUsernamesTable.findAll(
      {}, [], {});
    expect(subaccountsWithUsernames.length).toEqual(0);

    await subaccountUsernameGenerator();
    const subaccountsWithUsernamesAfter: SubaccountUsernamesFromDatabase[] = await
    SubaccountUsernamesTable.findAll(
      {}, [], {});

    expect(subaccountsWithUsernamesAfter.length).toEqual(subaccountsLength);
  });
});
