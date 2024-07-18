import {
  dbHelpers,
  SubaccountTable,
  testMocks,
  SubaccountUsernamesTable,
} from '@dydxprotocol-indexer/postgres';

import { RequestMethod } from '../../../../src/types';
import { sendRequest } from '../../../helpers/helpers';

describe('social-trading-controller', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await testMocks.seedData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  it('successfuly fetches subaccount info by address', async () => {
    const subaccounts = await SubaccountTable.findAll({}, [], { readReplica: true });
    const subaccount = subaccounts[0];

    const subaccountUsernames = await SubaccountUsernamesTable.create({
      subaccountId: subaccount.id,
      username: 'test_username',
    });

    const response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/trader/search?searchParam=${subaccount.address}`,
    });

    expect(response.status).toEqual(200);
    expect(response.body).toEqual({
      result: {
        address: subaccount.address,
        subaccountNumber: subaccount.subaccountNumber,
        username: subaccountUsernames.username,
        subaccountId: subaccount.id,
      },
    });
  });

  it('successfuly fetches subaccount info by username', async () => {

    const subaccounts = await SubaccountTable.findAll({}, [], { readReplica: true });
    const subaccount = subaccounts[0];
    const subaccountUsernames = await SubaccountUsernamesTable.create({
      subaccountId: subaccount.id,
      username: 'test_username',
    });

    const response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/trader/search?searchParam=${subaccountUsernames.username}`,
    });

    expect(response.status).toEqual(200);
    expect(response.body).toEqual({
      result: {
        address: subaccount.address,
        subaccountNumber: subaccount.subaccountNumber,
        username: subaccountUsernames.username,
        subaccountId: subaccount.id,
      },
    });
  });

});
