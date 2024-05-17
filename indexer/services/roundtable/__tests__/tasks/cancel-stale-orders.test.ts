import {
  BlockTable,
  OrderFromDatabase,
  OrderTable,
  OrderStatus,
  dbHelpers,
  testConstants,
  testMocks,
  PaginationFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import cancelStaleOrdersTask from '../../src/tasks/cancel-stale-orders';
import { defaultOrderGoodTilBlockTime } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import { stats } from '@dydxprotocol-indexer/base';
import _ from 'lodash';
import config from '../../src/config';

describe('cancel-stale-orders', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await Promise.all([
      BlockTable.create({
        ...testConstants.defaultBlock,
        blockHeight: '3',
      }),
      BlockTable.create({
        ...testConstants.defaultBlock,
        blockHeight: '4',
      }),
      BlockTable.create({
        ...testConstants.defaultBlock,
        blockHeight: '5',
      }),
    ]);
    jest.spyOn(stats, 'gauge');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.resetAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  it('succeeds with no orders', async () => {
    await cancelStaleOrdersTask();
  });

  it('updates stale OPEN orders', async () => {
    const createdOrders: OrderFromDatabase[] = await Promise.all([
      // will be updated
      OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '2',
        goodTilBlock: '1',
      }),
      // will not be updated, status is not OPEN
      OrderTable.create({
        ...testConstants.defaultOrder,
        totalFilled: testConstants.defaultOrder.size,
        clientId: '3',
        goodTilBlock: '3',
        status: OrderStatus.FILLED,
      }),
      // will be updated
      OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '4',
        goodTilBlock: '4',
      }),
      // will not be updated, is not stale, goodTilBlock >= latestBlock
      OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '5',
        goodTilBlock: '5',
      }),
      // will not be updated, is stateful order
      OrderTable.create(defaultOrderGoodTilBlockTime),
    ]);
    const createdOrderIds: string[] = createdOrders.map((order: OrderFromDatabase) => order.id);
    const expectedOrders: OrderFromDatabase[] = _.cloneDeep(createdOrders);
    expectedOrders[0].status = OrderStatus.CANCELED;
    expectedOrders[2].status = OrderStatus.CANCELED;

    await cancelStaleOrdersTask();

    const { results: ordersAfterTask }: PaginationFromDatabase<OrderFromDatabase> = await OrderTable
      .findAll(
        {
          id: createdOrderIds,
        },
        [],
        {},
      );
    expect(_.sortBy(ordersAfterTask, ['id'])).toEqual(_.sortBy(expectedOrders, ['id']));
    expect(stats.gauge).toHaveBeenCalledWith('roundtable.num_stale_orders.count', 2);
    expect(stats.gauge).toHaveBeenCalledWith('roundtable.num_stale_orders_canceled.count', 2);
  });

  it('updates up to the batch size of stale OPEN orders', async () => {
    const oldLimit: number = config.CANCEL_STALE_ORDERS_QUERY_BATCH_SIZE;
    config.CANCEL_STALE_ORDERS_QUERY_BATCH_SIZE = 1;

    const createdOrders: OrderFromDatabase[] = await Promise.all([
      // will be updated
      OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '2',
        goodTilBlock: '1',
      }),
      // will not be updated, status is not OPEN
      OrderTable.create({
        ...testConstants.defaultOrder,
        totalFilled: testConstants.defaultOrder.size,
        clientId: '3',
        goodTilBlock: '3',
        status: OrderStatus.FILLED,
      }),
      // will be updated
      OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '4',
        goodTilBlock: '4',
      }),
      // will not be updated, is not stale, goodTilBlock >= latestBlock
      OrderTable.create({
        ...testConstants.defaultOrder,
        clientId: '5',
        goodTilBlock: '5',
      }),
      // will not be updated, is stateful order
      OrderTable.create(defaultOrderGoodTilBlockTime),
    ]);
    const createdOrderIds: string[] = createdOrders.map((order: OrderFromDatabase) => order.id);
    const expectedOrders: OrderFromDatabase[] = _.cloneDeep(createdOrders);
    expectedOrders[0].status = OrderStatus.CANCELED;
    expectedOrders[2].status = OrderStatus.CANCELED;

    // Requires 2 runs of the task to update all the stale OPEN orders
    await cancelStaleOrdersTask();

    expect(stats.gauge).toHaveBeenCalledWith('roundtable.num_stale_orders.count', 1);
    expect(stats.gauge).toHaveBeenCalledWith('roundtable.num_stale_orders_canceled.count', 1);

    await cancelStaleOrdersTask();

    expect(stats.gauge).toHaveBeenCalledWith('roundtable.num_stale_orders.count', 1);
    expect(stats.gauge).toHaveBeenCalledWith('roundtable.num_stale_orders_canceled.count', 1);

    const { results: ordersAfterTask }: PaginationFromDatabase<OrderFromDatabase> = await OrderTable
      .findAll(
        {
          id: createdOrderIds,
        },
        [],
        {},
      );
    expect(_.sortBy(ordersAfterTask, ['id'])).toEqual(_.sortBy(expectedOrders, ['id']));

    config.CANCEL_STALE_ORDERS_QUERY_BATCH_SIZE = oldLimit;
  });
});
