import {
  dbHelpers,
  testConstants,
  testMocks,
  OrderTable,
  perpetualMarketRefresher,
  OrderSide,
  OrderCreateObject,
  protocolTranslations,
  OrderType,
  APIOrderStatusEnum,
} from '@dydxprotocol-indexer/postgres';
import { RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { getQueryString, sendRequest } from '../../../helpers/helpers';
import {
  placeOrder,
  redis,
  redisTestConstants,
} from '@dydxprotocol-indexer/redis';
import { redisClient } from '../../../../src/helpers/redis/redis-controller';
import {
  postgresAndRedisOrderToResponseObject,
  postgresOrderToResponseObject,
  redisOrderToResponseObject,
} from '../../../../src/request-helpers/request-transformer';
import {
  IndexerOrder,
  IndexerOrderId,
  IndexerOrder_Side,
  RedisOrder,
} from '@dydxprotocol-indexer/v4-protos';
import {
  ORDER_FLAG_CONDITIONAL,
} from '@dydxprotocol-indexer/v4-proto-parser';

describe('orders-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  beforeEach(async () => {
    await redis.deleteAllAsync(redisClient);
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterEach(async () => {
    await redis.deleteAllAsync(redisClient);
    await dbHelpers.clearData();
  });

  describe('/orders/:orderId', () => {
    it('Get /:orderId returns 400 if orderId is not found', async () => {
      await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/${testConstants.defaultOrderId}`,
        expectedStatus: 404,
      });
    });

    it('Get /:orderId gets order in postgres', async () => {
      await OrderTable.create(testConstants.defaultOrder);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/${testConstants.defaultOrderId}`,
      });

      expect(response.body).toEqual(postgresOrderToResponseObject(
        {
          ...testConstants.defaultOrder,
          id: testConstants.defaultOrderId,
        },
      ));
    });

    it('Get /:orderId gets order in redis', async () => {
      await placeOrder({
        redisOrder: redisTestConstants.defaultRedisOrder,
        client: redisClient,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/${testConstants.defaultOrderId}`,
      });

      expect(response.body).toEqual(
        redisOrderToResponseObject(
          redisTestConstants.defaultRedisOrder,
        ),
      );
    });

    it('Get /:orderId gets order in postgres and redis', async () => {
      await Promise.all([
        OrderTable.create(testConstants.defaultOrder),
        placeOrder({
          redisOrder: redisTestConstants.defaultRedisOrder,
          client: redisClient,
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/${testConstants.defaultOrderId}`,
      });

      expect(response.body).toEqual(
        postgresAndRedisOrderToResponseObject(
          {
            ...testConstants.defaultOrder,
            id: testConstants.defaultOrderId,
          },
          redisTestConstants.defaultRedisOrder,
        ),
      );
    });

    it('Get /:orderId errors when parameter is not a uuid', async () => {
      await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/orders/1',
        expectedStatus: 400,
      });
    });
  });

  describe('List orders', () => {
    const defaultQueryParams = {
      address: testConstants.defaultSubaccount.address,
      subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
    };
    const orderIdWithDifferentClobPairId: IndexerOrderId = {
      ...redisTestConstants.defaultOrderId,
      clobPairId: 2,
    };
    const newerOrderGoodTilBlockTimeId: IndexerOrderId = {
      ...redisTestConstants.defaultOrderIdGoodTilBlockTime,
      clientId: 4,
    };
    const orderWithDifferentClobPairId: IndexerOrder = {
      ...redisTestConstants.defaultOrder,
      orderId: orderIdWithDifferentClobPairId,
      goodTilBlock: 1200,
    };
    const newerOrderGoodTilBlockTime: IndexerOrder = {
      ...redisTestConstants.defaultOrderGoodTilBlockTime,
      orderId: newerOrderGoodTilBlockTimeId,
      goodTilBlockTime: 1_600_000_000,
    };
    const redisOrderWithDifferentMarket: RedisOrder = {
      ...redisTestConstants.defaultRedisOrder,
      order: orderWithDifferentClobPairId,
      id: OrderTable.orderIdToUuid(orderIdWithDifferentClobPairId),
      ticker: testConstants.defaultPerpetualMarket2.ticker,
    };
    const newerRedisOrderGoodTilBlockTime: RedisOrder = {
      ...redisTestConstants.defaultRedisOrderGoodTilBlockTime,
      order: newerOrderGoodTilBlockTime,
      id: OrderTable.orderIdToUuid(newerOrderGoodTilBlockTimeId),
    };
    const secondOrder: OrderCreateObject = {
      ...testConstants.defaultOrder,
      clientId: '2',
      goodTilBlock: '1250',
    };
    const secondOrderGoodTilBlockTime: OrderCreateObject = {
      ...testConstants.defaultOrderGoodTilBlockTime,
      clientId: '5',
      goodTilBlockTime: '2023-01-13T00:00:00.000Z',
    };
    const filledOrder: OrderCreateObject = {
      ...testConstants.defaultOrder,
      clientId: '3',
      goodTilBlock: '1251',
      status: APIOrderStatusEnum.FILLED,
    };
    const untriggeredOrder: OrderCreateObject = {
      ...testConstants.defaultOrder,
      clientId: '4',
      orderFlags: ORDER_FLAG_CONDITIONAL.toString(),
      status: APIOrderStatusEnum.UNTRIGGERED,
      triggerPrice: '1000',
    };

    it('Successfully gets multiple redis orders', async () => {
      await Promise.all([
        placeOrder({
          redisOrder: redisTestConstants.defaultRedisOrder,
          client: redisClient,
        }),
        placeOrder({
          redisOrder: redisOrderWithDifferentMarket,
          client: redisClient,
        }),
        placeOrder({
          redisOrder: redisTestConstants.defaultRedisOrderGoodTilBlockTime,
          client: redisClient,
        }),
        placeOrder({
          redisOrder: newerRedisOrderGoodTilBlockTime,
          client: redisClient,
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString(defaultQueryParams)}`,
      });

      expect(response.body).toEqual([ // by default sort by desc goodTilBlock
        redisOrderToResponseObject(newerRedisOrderGoodTilBlockTime),
        redisOrderToResponseObject(redisTestConstants.defaultRedisOrderGoodTilBlockTime),
        redisOrderToResponseObject(
          redisOrderWithDifferentMarket,
        ),
        redisOrderToResponseObject(
          redisTestConstants.defaultRedisOrder,
        ),
      ]);

      const response2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString({
          ...defaultQueryParams,
          returnLatestOrders: 'false',
        })}`,
      });

      expect(response2.body).toEqual([
        redisOrderToResponseObject(
          redisTestConstants.defaultRedisOrder,
        ),
        redisOrderToResponseObject(
          redisOrderWithDifferentMarket,
        ),
        redisOrderToResponseObject(redisTestConstants.defaultRedisOrderGoodTilBlockTime),
        redisOrderToResponseObject(newerRedisOrderGoodTilBlockTime),
      ]);
    });

    it.each([
      [
        'ticker',
        [redisTestConstants.defaultRedisOrder, redisOrderWithDifferentMarket],
        {
          ...defaultQueryParams,
          ticker: testConstants.defaultPerpetualMarket.ticker,
        },
        redisTestConstants.defaultRedisOrder,
      ],
      [
        'goodTilBlock',
        [
          redisTestConstants.defaultRedisOrder,
          redisOrderWithDifferentMarket,
          redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        ],
        {
          ...defaultQueryParams,
          goodTilBlockBeforeOrAt: protocolTranslations.getGoodTilBlock(
            redisTestConstants.defaultRedisOrder.order!,
          ),
        },
        redisTestConstants.defaultRedisOrder,
      ],
      [
        'goodTilBlockTime',
        [
          redisTestConstants.defaultRedisOrder,
          redisTestConstants.defaultRedisOrderGoodTilBlockTime,
          newerRedisOrderGoodTilBlockTime,
        ],
        {
          ...defaultQueryParams,
          goodTilBlockTimeBeforeOrAt: protocolTranslations.getGoodTilBlockTime(
            redisTestConstants.defaultRedisOrderGoodTilBlockTime.order!,
          ),
        },
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
      ],
    ])('Successfully filters redis order by %s', async (
      _testName: string,
      redisOrders: RedisOrder[],
      queryParams: any,
      expectedRedisOrder: RedisOrder,
    ) => {
      await Promise.all([
        placeOrder({
          redisOrder: redisTestConstants.defaultRedisOrder,
          client: redisClient,
        }),
        redisOrders.map(
          (order: RedisOrder) => placeOrder({
            redisOrder: order,
            client: redisClient,
          }),
        ),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString(queryParams)}`,
      });

      expect(response.body).toEqual([redisOrderToResponseObject(expectedRedisOrder)]);
    });

    it.each([
      [
        'BUY',
        redisTestConstants.defaultRedisOrder,
        {
          ...defaultQueryParams,
          side: OrderSide.BUY,
        },
      ],
      [
        'SELL',
        {
          ...redisOrderWithDifferentMarket,
          order: {
            ...orderWithDifferentClobPairId,
            side: IndexerOrder_Side.SIDE_SELL,
          },
        },
        {
          ...defaultQueryParams,
          side: OrderSide.SELL,
        },
      ],
    ])('Successfully filters redis order by side: %s', async (
      _testName: string,
      expectedOrder: RedisOrder,
      queryParams: any,
    ) => {
      await Promise.all([
        placeOrder({
          redisOrder: redisTestConstants.defaultRedisOrder,
          client: redisClient,
        }),
        placeOrder({
          redisOrder: {
            ...redisOrderWithDifferentMarket,
            order: {
              ...orderWithDifferentClobPairId,
              side: IndexerOrder_Side.SIDE_SELL,
            },
          },
          client: redisClient,
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString(queryParams)}`,
      });

      expect(response.body).toEqual([
        redisOrderToResponseObject(expectedOrder),
      ]);
    });

    it('Successfully gets multiple postgres orders', async () => {
      await Promise.all([
        OrderTable.create(testConstants.defaultOrder),
        OrderTable.create(secondOrder),
        OrderTable.create(testConstants.defaultOrderGoodTilBlockTime),
        OrderTable.create(secondOrderGoodTilBlockTime),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString(defaultQueryParams)}`,
      });

      expect(response.body).toEqual([
        postgresOrderToResponseObject({
          ...testConstants.defaultOrderGoodTilBlockTime,
          id: getUuidForTest(testConstants.defaultOrderGoodTilBlockTime),
        }),
        postgresOrderToResponseObject({
          ...secondOrderGoodTilBlockTime,
          id: getUuidForTest(secondOrderGoodTilBlockTime),
        }),
        postgresOrderToResponseObject({
          ...secondOrder,
          id: getUuidForTest(secondOrder),
        }),
        postgresOrderToResponseObject({
          ...testConstants.defaultOrder,
          id: testConstants.defaultOrderId,
        }),
      ]);

      const response2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString({
          ...defaultQueryParams,
          returnLatestOrders: 'false',
        })}`,
      });

      expect(response2.body).toEqual([
        postgresOrderToResponseObject({
          ...testConstants.defaultOrder,
          id: testConstants.defaultOrderId,
        }),
        postgresOrderToResponseObject({
          ...secondOrder,
          id: getUuidForTest(secondOrder),
        }),
        postgresOrderToResponseObject({
          ...secondOrderGoodTilBlockTime,
          id: getUuidForTest(secondOrderGoodTilBlockTime),
        }),
        postgresOrderToResponseObject({
          ...testConstants.defaultOrderGoodTilBlockTime,
          id: getUuidForTest(testConstants.defaultOrderGoodTilBlockTime),
        }),
      ]);
    });

    it('Successfully filters orders by status', async () => {
      await Promise.all([
        OrderTable.create(testConstants.defaultOrder),
        OrderTable.create(filledOrder),
        OrderTable.create(untriggeredOrder),
        placeOrder({
          redisOrder: redisTestConstants.defaultRedisOrder,
          client: redisClient,
        }),
        placeOrder({
          redisOrder: redisOrderWithDifferentMarket,
          client: redisClient,
        }),
      ]);

      let response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString({
          ...defaultQueryParams,
          status: APIOrderStatusEnum.OPEN,
        })}`,
      });

      // Filled order should not be in response.
      expect(response.body).toEqual([
        postgresAndRedisOrderToResponseObject(
          {
            ...testConstants.defaultOrder,
            id: testConstants.defaultOrderId,
          },
          redisTestConstants.defaultRedisOrder,
        ),
      ]);

      response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString({
          ...defaultQueryParams,
          status: APIOrderStatusEnum.FILLED,
        })}`,
      });

      // Filled order should be only order in response.
      expect(response.body).toEqual([
        postgresOrderToResponseObject({
          ...filledOrder,
          id: OrderTable.uuid(
            filledOrder.subaccountId,
            filledOrder.clientId,
            filledOrder.clobPairId,
            filledOrder.orderFlags,
          ),
        }),
      ]);

      response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString({
          ...defaultQueryParams,
          status: APIOrderStatusEnum.BEST_EFFORT_OPENED,
        })}`,
      });

      // Best effort opened order should be only order in response.
      expect(response.body).toEqual([
        redisOrderToResponseObject(redisOrderWithDifferentMarket),
      ]);

      response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString({
          ...defaultQueryParams,
          status: APIOrderStatusEnum.UNTRIGGERED,
        })}`,
      });

      // Untriggered order should be only order in response.
      expect(response.body).toEqual([
        postgresOrderToResponseObject({
          ...untriggeredOrder,
          id: OrderTable.uuid(
            untriggeredOrder.subaccountId,
            untriggeredOrder.clientId,
            untriggeredOrder.clobPairId,
            untriggeredOrder.orderFlags,
          ),
        }),
      ]);

      response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString({
          ...defaultQueryParams,
          status: [APIOrderStatusEnum.UNTRIGGERED, APIOrderStatusEnum.OPEN],
        })}`,
      });

      // Untriggered order and open order should be in response.
      expect(response.body).toEqual([
        postgresAndRedisOrderToResponseObject(
          {
            ...testConstants.defaultOrder,
            id: testConstants.defaultOrderId,
          },
          redisTestConstants.defaultRedisOrder,
        ),
        postgresOrderToResponseObject({
          ...untriggeredOrder,
          id: OrderTable.uuid(
            untriggeredOrder.subaccountId,
            untriggeredOrder.clientId,
            untriggeredOrder.clobPairId,
            untriggeredOrder.orderFlags,
          ),
        }),
      ]);
    });

    it('Successfully pulls both redis and postgres orders', async () => {
      await Promise.all([
        OrderTable.create(testConstants.defaultOrder),
        OrderTable.create(secondOrder),
        OrderTable.create(testConstants.defaultOrderGoodTilBlockTime),
        OrderTable.create(secondOrderGoodTilBlockTime),
        placeOrder({
          redisOrder: redisTestConstants.defaultRedisOrder,
          client: redisClient,
        }),
        placeOrder({
          redisOrder: redisOrderWithDifferentMarket,
          client: redisClient,
        }),
        placeOrder({
          redisOrder: redisTestConstants.defaultRedisOrderGoodTilBlockTime,
          client: redisClient,
        }),
        placeOrder({
          redisOrder: newerRedisOrderGoodTilBlockTime,
          client: redisClient,
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString(defaultQueryParams)}`,
      });

      expect(response.body).toEqual([
        postgresOrderToResponseObject({
          ...secondOrderGoodTilBlockTime,
          id: getUuidForTest(secondOrderGoodTilBlockTime),
        }),
        redisOrderToResponseObject(newerRedisOrderGoodTilBlockTime),
        postgresAndRedisOrderToResponseObject(
          {
            ...testConstants.defaultOrderGoodTilBlockTime,
            id: getUuidForTest(testConstants.defaultOrderGoodTilBlockTime),
          },
          redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        ),
        postgresOrderToResponseObject({
          ...secondOrder,
          id: getUuidForTest(secondOrder),
        }),
        redisOrderToResponseObject(redisOrderWithDifferentMarket),
        postgresAndRedisOrderToResponseObject(
          {
            ...testConstants.defaultOrder,
            id: testConstants.defaultOrderId,
          },
          redisTestConstants.defaultRedisOrder,
        ),
      ]);
    });

    it.each([
      [
        'goodTilBlock passed in with goodTilBlockTime',
        {
          ...defaultQueryParams,
          goodTilBlock: '50',
          goodTilBlockTime: '2023-01-13T00:00:00.000Z',
        },
        'goodTilBlock',
        'Cannot provide both goodTilBlock and goodTilBlockTime',
      ],
      [
        'invalid side',
        {
          ...defaultQueryParams,
          side: 'INVALID',
        },
        'side',
        `side must be one of ${Object.values(OrderSide)}`,
      ],
      [
        'invalid type',
        {
          ...defaultQueryParams,
          type: 'INVALID',
        },
        'type',
        `type must be one of ${Object.values(OrderType)}`,
      ],
      [
        'invalid status',
        {
          ...defaultQueryParams,
          status: 'INVALID',
        },
        'status',
        `status must be one of ${Object.values(APIOrderStatusEnum)}`,
      ],
    ])('Returns 400 when validation fails: %s', async (
      _reason: string,
      queryParams: {
        address?: string,
        subaccountNumber?: number,
        goodTilBlock?: string,
        goodTilBlockTime?: string,
        side?: string,
      },
      fieldWithError: string,
      expectedErrorMsg: string,
    ) => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString(queryParams)}`,
        expectedStatus: 400,
      });

      expect(response.body).toEqual(expect.objectContaining({
        errors: expect.arrayContaining([
          expect.objectContaining({
            param: fieldWithError,
            msg: expectedErrorMsg,
          }),
        ]),
      }));
    });
  });
});

function getUuidForTest(orderCreate: OrderCreateObject): string {
  return OrderTable.uuid(
    orderCreate.subaccountId,
    orderCreate.clientId,
    orderCreate.clobPairId,
    orderCreate.orderFlags,
  );
}
