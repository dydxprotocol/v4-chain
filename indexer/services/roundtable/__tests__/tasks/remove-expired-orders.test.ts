import {
  logger,
  stats,
} from '@dydxprotocol-indexer/base';
import { KafkaTopics, ProducerMessage, producer } from '@dydxprotocol-indexer/kafka';
import {
  dbHelpers,
  perpetualMarketRefresher,
  testMocks,
  BlockTable,
  OrderTable,
} from '@dydxprotocol-indexer/postgres';
import {
  IndexerOrder,
  IndexerOrderId,
  IndexerOrder_TimeInForce,
  RedisOrder,
} from '@dydxprotocol-indexer/v4-protos';
import {
  placeOrder,
  redis,
  redisTestConstants,
  OrdersCache,
  OrdersDataCache,
  OrderData,
  OrderExpiryCache,
} from '@dydxprotocol-indexer/redis';
import { ProducerRecord } from 'kafkajs';
import _ from 'lodash';
import { DateTime } from 'luxon';

import config from '../../src/config';
import { redisClient } from '../../src/helpers/redis';
import { getExpiredOffChainUpdateMessage } from '../../src/helpers/websocket';
import removeExpiredOrdersTask from '../../src/tasks/remove-expired-orders';
import { getOrderIdHash } from '@dydxprotocol-indexer/v4-proto-parser';

describe('remove-expired-orders', () => {
  let producerSendMock: jest.SpyInstance;

  function expectProducerMessages(expectedOrderIds: IndexerOrderId[]) {
    const sentMessages: ProducerMessage[] = [];
    for (const call of producerSendMock.mock.calls) {
      const message: ProducerRecord = call[0];
      expect(message.topic).toEqual(KafkaTopics.TO_VULCAN);
      sentMessages.push(...message.messages as ProducerMessage[]);
    }
    const expectedMessages: ProducerMessage[] = _.map(
      expectedOrderIds,
      (orderId: IndexerOrderId): ProducerMessage => {
        return {
          key: getOrderIdHash(orderId),
          value: getExpiredOffChainUpdateMessage(orderId),
        };
      },
    );
    expect(sentMessages).toHaveLength(expectedMessages.length);
    expect(sentMessages).toEqual(expect.arrayContaining(expectedMessages));
  }

  function expectStats(
    numExpectedProducerMessages: number,
    extraIncrementCalls?: [string, number][],
  ): void {
    const expectedCalls: [string, number][] = [...(extraIncrementCalls ?? [])];
    for (let i: number = 0; i < numExpectedProducerMessages; ++i) {
      expectedCalls.push([`${config.SERVICE_NAME}.expiry_message_sent`, 1]);
    }

    expect((stats.increment as jest.MockedFunction<typeof stats.increment>).mock.calls)
      .toEqual(expect.arrayContaining(expectedCalls));
    expect(stats.increment).toHaveBeenCalledTimes(expectedCalls.length);

    expect(stats.timing).toBeCalledWith(
      `${config.SERVICE_NAME}.remove_expired_orders.timing`,
      expect.any(Number),
    );
  }

  const ORDER_CLIENT1_GTB8: IndexerOrder = {
    ...redisTestConstants.defaultOrder,
    timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
    goodTilBlock: 8,
  };
  const ORDERID_CLIENT1: IndexerOrderId = redisTestConstants.defaultOrderId;
  const ORDERID_CLIENT2: IndexerOrderId = { ...ORDERID_CLIENT1, clientId: 2 };
  const ORDERID_CLIENT3: IndexerOrderId = { ...ORDERID_CLIENT1, clientId: 3 };
  const ORDERID_CLIENT4: IndexerOrderId = { ...ORDERID_CLIENT1, clientId: 4 };
  const REDISORDER_CLIENT1_GTB8: RedisOrder = {
    ...redisTestConstants.defaultRedisOrder,
    order: ORDER_CLIENT1_GTB8,
  };
  const REDISORDER_CLIENT2_GTB9: RedisOrder = {
    ...redisTestConstants.defaultRedisOrder,
    id: OrderTable.orderIdToUuid(ORDERID_CLIENT2),
    order: {
      ...ORDER_CLIENT1_GTB8,
      orderId: ORDERID_CLIENT2,
      goodTilBlock: 9,
    },
  };
  const REDISORDER_CLIENT3_GTB10: RedisOrder = {
    ...redisTestConstants.defaultRedisOrder,
    id: OrderTable.orderIdToUuid(ORDERID_CLIENT3),
    order: {
      ...ORDER_CLIENT1_GTB8,
      orderId: ORDERID_CLIENT3,
      goodTilBlock: 10,
    },
  };
  const REDISORDER_CLIENT4_GTB11: RedisOrder = {
    ...redisTestConstants.defaultRedisOrder,
    id: OrderTable.orderIdToUuid(ORDERID_CLIENT4),
    order: {
      ...ORDER_CLIENT1_GTB8,
      orderId: ORDERID_CLIENT4,
      goodTilBlock: 11,
    },
  };
  const defaultOrderData: OrderData = {
    totalFilledQuantums: '0',
    restingOnBook: true,
    goodTilBlock: ORDER_CLIENT1_GTB8.goodTilBlock!.toString(),
  };

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await Promise.all([
      BlockTable.create({ blockHeight: '30', time: DateTime.utc(2022, 6, 1).toISO() }),
      perpetualMarketRefresher.updatePerpetualMarkets(),
      placeOrder({ redisOrder: REDISORDER_CLIENT1_GTB8, client: redisClient }),
      placeOrder({ redisOrder: REDISORDER_CLIENT2_GTB9, client: redisClient }),
      placeOrder({ redisOrder: REDISORDER_CLIENT3_GTB10, client: redisClient }),
      placeOrder({ redisOrder: REDISORDER_CLIENT4_GTB11, client: redisClient }),
    ]);

    // mock out producer
    producerSendMock = jest.spyOn(producer, 'send');
    // mock out stats
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
    // mock out logger
    jest.spyOn(logger, 'info');
    jest.spyOn(logger, 'error');

    // If this isn't after the logger mock, we get an error:
    // [KafkaJSError: The producer is disconnected]
    jest.resetAllMocks();
  });

  afterEach(async () => {
    // Calling mock.restoreMock() at the end of the tests where OrdersCache.getOrder &
    // OrdersDataCache.getOrderDataWithUUID caused errors in later tests
    jest.restoreAllMocks();
    await Promise.all([
      dbHelpers.clearData(),
      redis.deleteAllAsync(redisClient),
    ]);
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  it('removes expired orders with GTB <= currentBlock - 20', async () => {
    const expectedMessages: IndexerOrderId[] = [ORDERID_CLIENT1, ORDERID_CLIENT2, ORDERID_CLIENT3];
    await removeExpiredOrdersTask();

    expect(logger.error).not.toHaveBeenCalled();
    expectProducerMessages(expectedMessages);
    expectStats(expectedMessages.length);
  });

  it('when latest block is updated, retrieved expired orders are updated', async () => {
    const expectedMessages: IndexerOrderId[] = [
      ORDERID_CLIENT1,
      ORDERID_CLIENT2,
      ORDERID_CLIENT3,
      ORDERID_CLIENT4,
    ];
    await BlockTable.create({ blockHeight: '31', time: DateTime.utc(2022, 6, 2).toISO() });

    await removeExpiredOrdersTask();

    expect(logger.error).not.toHaveBeenCalled();
    expectProducerMessages(expectedMessages);
    expectStats(expectedMessages.length);
  });

  it('error: when orders-cache cannot find order, continue the remainder', async () => {
    const expectedMessages: IndexerOrderId[] = [ORDERID_CLIENT1, ORDERID_CLIENT3];
    jest.spyOn(OrdersCache, 'getOrder')
      .mockResolvedValue(REDISORDER_CLIENT4_GTB11)
      .mockResolvedValueOnce(REDISORDER_CLIENT1_GTB8)
      .mockResolvedValueOnce(null)
      .mockResolvedValueOnce(REDISORDER_CLIENT3_GTB10);

    await removeExpiredOrdersTask();

    expectProducerMessages(expectedMessages);
    expectStats(
      expectedMessages.length,
      [['roundtable.expired_order_data_not_found', 1]],
    );
  });

  it('error: when orders-data-cache cannot find order, continue the remainder', async () => {
    const expectedMessages: IndexerOrderId[] = [ORDERID_CLIENT2, ORDERID_CLIENT3];
    jest.spyOn(OrdersDataCache, 'getOrderDataWithUUID')
      .mockResolvedValue(defaultOrderData)
      .mockResolvedValueOnce(null);

    await removeExpiredOrdersTask();

    expectProducerMessages(expectedMessages);
    expectStats(
      expectedMessages.length,
      [['roundtable.expired_order_data_not_found', 1]],
    );
  });

  it('error: excludes orders with expiry > (blockHeight - buffer)', async () => {
    const expectedMessages: IndexerOrderId[] = [ORDERID_CLIENT2, ORDERID_CLIENT3];
    const mockOrder: RedisOrder = {
      ...REDISORDER_CLIENT1_GTB8,
      order: {
        ...ORDER_CLIENT1_GTB8,
        goodTilBlock: 11,
      },
    };
    jest.spyOn(OrdersCache, 'getOrder')
      .mockResolvedValue(REDISORDER_CLIENT4_GTB11)
      .mockResolvedValueOnce(mockOrder)
      .mockResolvedValueOnce(REDISORDER_CLIENT2_GTB9)
      .mockResolvedValueOnce(REDISORDER_CLIENT3_GTB10);

    await removeExpiredOrdersTask();

    expectProducerMessages(expectedMessages);
    expectStats(
      expectedMessages.length,
      [['roundtable.indexer_expired_order_has_newer_expiry', 1]],
    );
  });

  it('increments the stat for fully-filled order expiries & continues as normal', async () => {
    const expectedMessages: IndexerOrderId[] = [ORDERID_CLIENT1, ORDERID_CLIENT2, ORDERID_CLIENT3];
    jest.spyOn(OrdersDataCache, 'getOrderDataWithUUID')
      .mockResolvedValue(defaultOrderData)
      .mockResolvedValueOnce(defaultOrderData)
      .mockResolvedValueOnce({
        ...defaultOrderData,
        totalFilledQuantums: ORDER_CLIENT1_GTB8.quantums.toString(),
      });

    await removeExpiredOrdersTask();

    expectProducerMessages(expectedMessages);
    expect(logger.error).not.toHaveBeenCalled();
    expectStats(expectedMessages.length, [
      [`${config.SERVICE_NAME}.fully_filled_orders_expired_by_roundtable`, 1],
    ]);
  });

  it('error: catches critical error and logs', async () => {
    const expectedMessages: IndexerOrderId[] = [];
    jest.spyOn(OrderExpiryCache, 'getOrderExpiries').mockImplementation(
      () => { throw new Error('TEST'); },
    );

    await removeExpiredOrdersTask();

    expect(logger.error).toHaveBeenCalledWith({
      at: 'remove-expired-orders#runTask',
      message: 'Error occurred in task to remove expired orders',
      error: new Error('TEST'),
    });
    expect(logger.error).toBeCalledTimes(1);
    expectProducerMessages(expectedMessages);
    expectStats(expectedMessages.length);
  });

  it('error: catches rejected promise, logs, and continues as normal', async () => {
    const expectedMessages: IndexerOrderId[] = [ORDERID_CLIENT2, ORDERID_CLIENT3];
    jest.spyOn(OrdersDataCache, 'getOrderDataWithUUID')
      .mockResolvedValue(defaultOrderData)
      .mockRejectedValueOnce(new Error('testing'));

    await removeExpiredOrdersTask();

    expect(logger.error).toHaveBeenCalledWith({
      at: 'remove-expired-orders#runTask',
      message: 'Encountered error expiring order',
      orderUuid: OrderTable.orderIdToUuid(ORDERID_CLIENT1),
      expiryCutoff: 10,
      error: new Error('testing'),
    });
    expect(logger.error).toBeCalledTimes(1);
    expectProducerMessages(expectedMessages);
    expectStats(expectedMessages.length);
  });

});
