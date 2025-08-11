import {
  dbHelpers,
  testMocks,
  PersistentCacheTable,
  PersistentCacheKeys,
  BlockTable,
  FillTable,
  OrderTable,
  FundingPaymentsTable,
  OrderSide,
  FundingIndexUpdatesTable,
  PositionSide,
  FillType,
  Liquidity,
  Ordering,
} from '@dydxprotocol-indexer/postgres';

import updatePnlTicksTask from '../../src/tasks/update-pnl-ticks';

import {
  createdDateTime,
  defaultFundingPayment,
  defaultFundingPayment2,
} from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

describe('update-pnl-ticks', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await dbHelpers.clearData();
    await testMocks.seedData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  it('Maintains persistent cache value across multiple runs', async () => {
    await FundingPaymentsTable.create(defaultFundingPayment);
    await FundingPaymentsTable.create(defaultFundingPayment2);

    await updatePnlTicksTask();
    const persistentCache = await PersistentCacheTable.findById(
      PersistentCacheKeys.PNL_TICKS_LAST_PROCESSED_HEIGHT,
    );
    expect(persistentCache?.value).toEqual('2');
  });
});