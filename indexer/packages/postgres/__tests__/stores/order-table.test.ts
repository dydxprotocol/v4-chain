import {
  OrderColumns,
  OrderCreateObject,
  OrderFromDatabase,
  Ordering,
  OrderStatus,
  PaginationFromDatabase,
  TimeInForce,
} from '../../src/types';
import * as OrderTable from '../../src/stores/order-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import { defaultOrder, defaultOrderGoodTilBlockTime, defaultOrderId } from '../helpers/constants';
import { CheckViolationError } from 'objection';
import { ORDER_FLAG_CONDITIONAL, ORDER_FLAG_LONG_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import Transaction from '../../src/helpers/transaction';

describe('Order store', () => {
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

  it('isLongTermOrConditionalOrder', () => {
    expect(OrderTable.isLongTermOrConditionalOrder('64')).toEqual(true);
    expect(OrderTable.isLongTermOrConditionalOrder('65')).toEqual(true);
    expect(OrderTable.isLongTermOrConditionalOrder('32')).toEqual(true);
    expect(OrderTable.isLongTermOrConditionalOrder('33')).toEqual(true);
    expect(OrderTable.isLongTermOrConditionalOrder('128')).toEqual(false);
    expect(OrderTable.isLongTermOrConditionalOrder('0')).toEqual(false);
    expect(OrderTable.isLongTermOrConditionalOrder('191')).toEqual(true);
    expect(OrderTable.isLongTermOrConditionalOrder('159')).toEqual(false);
  });

  it('Successfully creates a Order', async () => {
    await OrderTable.create(defaultOrder);
  });

  it('Successfully creates an Order with goodTilBlockTime', async () => {
    await OrderTable.create(defaultOrderGoodTilBlockTime);

    const {
      results: orders,
    }: PaginationFromDatabase<OrderFromDatabase> = await OrderTable.findAll({}, [], {});

    expect(orders).toHaveLength(1);
    expect(orders[0]).toEqual(expect.objectContaining({
      ...defaultOrderGoodTilBlockTime,
      goodTilBlock: null,
    }));
  });

  it('Successfully finds all Orders', async () => {
    await Promise.all([
      OrderTable.create(defaultOrder),
      OrderTable.create({
        ...defaultOrder,
        clientId: '2',
      }),
    ]);

    const {
      results: orders,
    }: PaginationFromDatabase<OrderFromDatabase> = await OrderTable.findAll({}, [], {
      orderBy: [[OrderColumns.clientId, Ordering.ASC]],
    });

    expect(orders.length).toEqual(2);
    expect(orders[0]).toEqual(expect.objectContaining(defaultOrder));
    expect(orders[1]).toEqual(expect.objectContaining({
      ...defaultOrder,
      clientId: '2',
    }));
  });

  it('Successfully finds all Orders using pagination', async () => {
    await Promise.all([
      OrderTable.create(defaultOrder),
      OrderTable.create({
        ...defaultOrder,
        clientId: '2',
      }),
    ]);

    const responsePageOne: PaginationFromDatabase<OrderFromDatabase> = await OrderTable.findAll({
      page: 1,
      limit: 1,
    },
    [],
    {
      orderBy: [[OrderColumns.clientId, Ordering.ASC]],
    });

    expect(responsePageOne.results.length).toEqual(1);
    expect(responsePageOne.results[0]).toEqual(expect.objectContaining(defaultOrder));
    expect(responsePageOne.offset).toEqual(0);
    expect(responsePageOne.total).toEqual(2);

    const responsePageTwo: PaginationFromDatabase<OrderFromDatabase> = await OrderTable.findAll({
      page: 2,
      limit: 1,
    },
    [],
    {
      orderBy: [[OrderColumns.clientId, Ordering.ASC]],
    });

    expect(responsePageTwo.results.length).toEqual(1);
    expect(responsePageTwo.results[0]).toEqual(expect.objectContaining({
      ...defaultOrder,
      clientId: '2',
    }));
    expect(responsePageTwo.offset).toEqual(1);
    expect(responsePageTwo.total).toEqual(2);

    const responsePageAllPages: PaginationFromDatabase<OrderFromDatabase> = await OrderTable
      .findAll({
        page: 1,
        limit: 2,
      },
      [],
      {
        orderBy: [[OrderColumns.clientId, Ordering.ASC]],
      });

    expect(responsePageAllPages.results.length).toEqual(2);
    expect(responsePageAllPages.results[0]).toEqual(expect.objectContaining(defaultOrder));
    expect(responsePageAllPages.results[1]).toEqual(expect.objectContaining({
      ...defaultOrder,
      clientId: '2',
    }));
    expect(responsePageAllPages.offset).toEqual(0);
    expect(responsePageAllPages.total).toEqual(2);
  });

  it('findOpenLongTermOrConditionalOrders', async () => {
    await Promise.all([
      OrderTable.create(defaultOrder),
      OrderTable.create({
        ...defaultOrder,
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
      }),
      OrderTable.create({
        ...defaultOrder,
        orderFlags: ORDER_FLAG_CONDITIONAL.toString(),
      }),
    ]);

    const orders: OrderFromDatabase[] = await OrderTable.findOpenLongTermOrConditionalOrders();
    expect(orders.length).toEqual(2);
  });

  // TODO: Add a bunch of tests for different search parameters
  it('Successfully finds Order with clientId', async () => {
    await Promise.all([
      OrderTable.create(defaultOrder),
      OrderTable.create({
        ...defaultOrder,
        clientId: '2',
      }),
    ]);

    const { results: orders }: PaginationFromDatabase<OrderFromDatabase> = await OrderTable.findAll(
      {
        clientId: '1',
      },
      [],
      { readReplica: true },
    );

    expect(orders.length).toEqual(1);
    expect(orders[0]).toEqual(expect.objectContaining(defaultOrder));
  });

  it('Successfully finds a Order', async () => {
    await OrderTable.create(defaultOrder);

    const order: OrderFromDatabase | undefined = await OrderTable.findById(
      OrderTable.uuid(
        defaultOrder.subaccountId,
        defaultOrder.clientId,
        defaultOrder.clobPairId,
        defaultOrder.orderFlags,
      ),
    );

    expect(order).toEqual(expect.objectContaining(defaultOrder));
  });

  it('Successfully finds all orders by subaccount/clob pair id', async () => {
    await Promise.all([
      OrderTable.create(defaultOrder),
      OrderTable.create({
        ...defaultOrder,
        clientId: '2',
      }),
    ]);
    const orders: OrderFromDatabase[] = await OrderTable.findBySubaccountIdAndClobPair(
      defaultOrder.subaccountId, defaultOrder.clobPairId,
    );

    expect(orders).toHaveLength(2);
  });

  it('Successfully finds all orders by subaccount/clob pair id after height', async () => {
    await Promise.all([
      OrderTable.create(defaultOrder),
      OrderTable.create({
        ...defaultOrder,
        clientId: '2',
        createdAtHeight: '5',
      }),
    ]);
    const orders: OrderFromDatabase[] = await OrderTable.findBySubaccountIdAndClobPairAfterHeight(
      defaultOrder.subaccountId,
      defaultOrder.clobPairId,
      4,
    );

    expect(orders).toHaveLength(1);
  });

  it.each([
    [
      'goodTilBlockBeforeOrAt',
      'goodTilBlock',
      {
        goodTilBlockBeforeOrAt: '101',
      },
      {
        ...defaultOrder,
        goodTilBlockTime: null,
      },
    ],
    [
      'goodTilBlockTimeBeforeOrAt',
      'goodTilBlockTime',
      {
        goodTilBlockTimeBeforeOrAt: '2023-01-23T00:00:00.000Z',
      },
      {
        ...defaultOrderGoodTilBlockTime,
        goodTilBlock: null,
      },
    ],
    [
      'goodTilBlockAfter',
      'goodTilBlock',
      {
        goodTilBlockAfter: '1',
      },
      {
        ...defaultOrder,
        goodTilBlockTime: null,
      },
    ],
    [
      'goodTilBlockTimeAfter',
      'goodTilBlockTime',
      {
        goodTilBlockTimeAfter: '2022-02-01T00:00:00.000Z',
      },
      {
        ...defaultOrderGoodTilBlockTime,
        goodTilBlock: null,
      },
    ],
  ])(
    'Successfully finds all orders by %s, excludes orders with null %s',
    async (
      _filterParam: string,
      _nullColumn: string,
      filter: Object,
      expectedOrder: Object,
    ) => {
      await Promise.all([
        OrderTable.create(defaultOrder),
        OrderTable.create(defaultOrderGoodTilBlockTime),
      ]);

      const { results: orders }: PaginationFromDatabase<OrderFromDatabase> = await OrderTable
        .findAll(
          filter,
          [],
          { readReplica: true },
        );

      expect(orders).toHaveLength(1);
      expect(orders[0]).toEqual(expect.objectContaining(expectedOrder));
    },
  );

  it('Successfully updates an Order', async () => {
    await OrderTable.create(defaultOrder);

    const order: OrderFromDatabase | undefined = await OrderTable.update({
      id: defaultOrderId,
      size: '32.50',
    });

    expect(order).toEqual(expect.objectContaining({
      ...defaultOrder,
      size: '32.50',
    }));

  });

  it('Successfully upserts a new Order', async () => {
    const createOrder: OrderCreateObject = {
      ...defaultOrder,
      createdAtHeight: '2',
    };
    const order: OrderFromDatabase = await OrderTable.upsert(createOrder);

    expect(order).toEqual(expect.objectContaining(createOrder));
  });

  it('Successfully upserts a new Order with UNTRIGGERED', async () => {
    const createOrder: OrderCreateObject = {
      ...defaultOrder,
      status: OrderStatus.UNTRIGGERED,
      createdAtHeight: '2',
    };
    const order: OrderFromDatabase = await OrderTable.upsert(createOrder);

    expect(order).toEqual(expect.objectContaining(createOrder));
  });

  it('Successfully upserts a new Order and updates status', async () => {
    const createOrder: OrderCreateObject = {
      ...defaultOrder,
      createdAtHeight: '2',
      totalFilled: '50',
    };
    const order: OrderFromDatabase = await OrderTable.upsert(createOrder);

    expect(order).toEqual(expect.objectContaining({
      ...createOrder,
      status: OrderStatus.FILLED,
    }));
  });

  it('Successfully upserts an existing Order, changing status to filled', async () => {
    await OrderTable.create(defaultOrder);

    const upsertOrder: OrderCreateObject = {
      ...defaultOrder,
      totalFilled: defaultOrder.size,
    };
    const order: OrderFromDatabase = await OrderTable.upsert(upsertOrder);

    expect(order).toEqual(expect.objectContaining({
      ...upsertOrder,
      status: OrderStatus.FILLED,
    }));
  });

  it('Successfully upserts an existing Order', async () => {
    await OrderTable.create(defaultOrder);

    const upsertOrder: OrderCreateObject = {
      ...defaultOrder,
      totalFilled: '10.65',
      timeInForce: TimeInForce.FOK,
      reduceOnly: true,
    };
    const order: OrderFromDatabase = await OrderTable.upsert(upsertOrder);

    expect(order).toEqual(expect.objectContaining(upsertOrder));
  });

  it('Successfully upserts an existing Order multiple times', async () => {
    await OrderTable.create(defaultOrder);

    const upsertOrder: OrderCreateObject = {
      ...defaultOrder,
      totalFilled: '10.65',
    };
    let order: OrderFromDatabase = await OrderTable.upsert(upsertOrder);
    order = await OrderTable.upsert(upsertOrder);

    expect(order).toEqual(expect.objectContaining(upsertOrder));
  });

  it('Successfully upserts an existing Order with OPEN status', async () => {
    await OrderTable.create(defaultOrder);

    const upsertOrder: OrderCreateObject = {
      ...defaultOrder,
      totalFilled: '1',
    };
    const order: OrderFromDatabase = await OrderTable.upsert(upsertOrder);

    expect(order).toEqual(expect.objectContaining({
      ...upsertOrder,
      status: OrderStatus.OPEN,
    }));
  });

  it('Successfully upserts an existing Order, created within the same transaction', async () => {
    const txId: number = await Transaction.start();
    await OrderTable.create(defaultOrder, { txId });

    const upsertOrder: OrderCreateObject = {
      ...defaultOrder,
      totalFilled: '10.65',
    };
    await OrderTable.upsert(upsertOrder, { txId });
    await Transaction.commit(txId);

    // Find order after committing transaction
    const order: OrderFromDatabase | undefined = await OrderTable.findById(
      OrderTable.uuid(
        defaultOrder.subaccountId,
        defaultOrder.clientId,
        defaultOrder.clobPairId,
        defaultOrder.orderFlags,
      ),
    );

    expect(order).toEqual(expect.objectContaining({
      ...upsertOrder,
    }));
  });

  it('Successfully upserts an existing Order, respects existing BEST_EFFORT_CANCELED status',
    async () => {
      await OrderTable.create(defaultOrder);

      const upsertOrder: OrderCreateObject = {
        ...defaultOrder,
        totalFilled: defaultOrder.size,
        status: OrderStatus.BEST_EFFORT_CANCELED,
      };
      const order: OrderFromDatabase = await OrderTable.upsert(upsertOrder);

      expect(order).toEqual(expect.objectContaining({
        ...upsertOrder,
        status: OrderStatus.BEST_EFFORT_CANCELED,
      }));
    });

  it('Successfully upserts an existing Order with UNTRIGGRERED status', async () => {
    await OrderTable.create(defaultOrder);

    const upsertOrder: OrderCreateObject = {
      ...defaultOrder,
      status: OrderStatus.UNTRIGGERED,
    };
    const order: OrderFromDatabase = await OrderTable.upsert(upsertOrder);

    expect(order).toEqual(expect.objectContaining({
      ...upsertOrder,
      status: OrderStatus.UNTRIGGERED,
    }));
  });

  it('Successfully upserts an existing Order, with fixed-decimal notation', async () => {
    await OrderTable.create(defaultOrder);

    const upsertOrder: OrderCreateObject = {
      ...defaultOrder,
      totalFilled: '0.00000001', // should not be converted to exponential notation
    };
    const order: OrderFromDatabase = await OrderTable.upsert(upsertOrder);

    expect(order).toEqual(expect.objectContaining({
      ...upsertOrder,
    }));
  });

  it('Fails to create invalid order with both goodTilBlock and goodTilBlockTime set', async () => {
    const invalidOrder: OrderCreateObject = {
      ...defaultOrder,
      goodTilBlockTime: defaultOrderGoodTilBlockTime.goodTilBlockTime,
    };

    await expect(OrderTable.create(invalidOrder)).rejects.toBeInstanceOf(CheckViolationError);
  });

  it('Successfully updates stale order status by id', async () => {
    const createdOrders: OrderFromDatabase[] = await Promise.all([
      // will be updated
      OrderTable.create(
        defaultOrder,
      ),
      // will not be updated as status doesn't match old status in update
      OrderTable.create({
        ...defaultOrder,
        clientId: '2',
        totalFilled: defaultOrder.size,
        status: OrderStatus.FILLED,
      }),
      // will be updated
      OrderTable.create({
        ...defaultOrder,
        clientId: '3',
        goodTilBlock: '120',
      }),
      // will not be updated as goodTilBlock >= latestBlock in update
      OrderTable.create({
        ...defaultOrder,
        clientId: '4',
        goodTilBlock: '150',
      }),
      // will not be updated as goodTilBlock is null
      OrderTable.create(defaultOrderGoodTilBlockTime),
    ]);

    const updatedOrders: OrderFromDatabase[] = await OrderTable.updateStaleOrderStatusByIds(
      OrderStatus.OPEN,
      OrderStatus.CANCELED,
      '135',
      createdOrders.map((order: OrderFromDatabase) => order.id),
    );

    const expectedOrders: OrderFromDatabase[] = [createdOrders[0], createdOrders[2]].map(
      (order: OrderFromDatabase) => {
        return {
          ...order,
          status: OrderStatus.CANCELED,
        };
      },
    );
    expect(updatedOrders).toEqual(expect.arrayContaining(expectedOrders));
  });
});
