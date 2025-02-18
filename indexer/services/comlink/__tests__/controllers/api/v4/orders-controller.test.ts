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
        testConstants.defaultSubaccount.subaccountNumber,
      ));
    });

    it('Get /:orderId gets isolated position order in postgres', async () => {
      await OrderTable.create(testConstants.isolatedMarketOrder);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/${testConstants.isolatedMarketOrderId}`,
      });

      expect(response.body).toEqual(postgresOrderToResponseObject(
        {
          ...testConstants.isolatedMarketOrder,
          id: testConstants.isolatedMarketOrderId,
        },
        testConstants.isolatedSubaccount.subaccountNumber,
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

    it('Get /:orderId gets isolated position order in redis', async () => {
      await placeOrder({
        redisOrder: redisTestConstants.isolatedMarketRedisOrder,
        client: redisClient,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/${testConstants.isolatedMarketOrderId}`,
      });

      expect(response.body).toEqual(
        redisOrderToResponseObject(
          redisTestConstants.isolatedMarketRedisOrder,
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
          {
            [testConstants.defaultSubaccountId]:
              testConstants.defaultSubaccount.subaccountNumber,
          },
          redisTestConstants.defaultRedisOrder,
        ),
      );
    });

    it('Get /:orderId gets isolated market order in postgres and redis', async () => {
      await Promise.all([
        OrderTable.create(testConstants.isolatedMarketOrder),
        placeOrder({
          redisOrder: redisTestConstants.isolatedMarketRedisOrder,
          client: redisClient,
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/${testConstants.isolatedMarketOrderId}`,
      });

      expect(response.body).toEqual(
        postgresAndRedisOrderToResponseObject(
          {
            ...testConstants.isolatedMarketOrder,
            id: testConstants.isolatedMarketOrderId,
          },
          {
            [testConstants.isolatedSubaccountId]:
                testConstants.isolatedSubaccount.subaccountNumber,
          },
          redisTestConstants.isolatedMarketRedisOrder,
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
    const isolatedOrderIdWithDiffClientId: IndexerOrderId = {
      ...redisTestConstants.isolatedMarketOrderId,
      clientId: 2,
    };
    const isolatedOrderWithDiffClientId: IndexerOrder = {
      ...redisTestConstants.isolatedMarketOrder,
      orderId: isolatedOrderIdWithDiffClientId,
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
      goodTilBlock: '99',
    };
    const untriggeredOrderId: string = OrderTable.uuid(
      untriggeredOrder.subaccountId,
      untriggeredOrder.clientId,
      untriggeredOrder.clobPairId,
      untriggeredOrder.orderFlags,
    );
    const isolatedRedisOrder: RedisOrder = {
      ...redisTestConstants.isolatedMarketRedisOrder,
      order: {
        ...redisTestConstants.isolatedMarketOrder,
        goodTilBlock: 1200,
      },
    };
    const isolatedRedisOrderWithDiffClientId: RedisOrder = {
      ...redisTestConstants.isolatedMarketRedisOrder,
      order: isolatedOrderWithDiffClientId,
      id: OrderTable.orderIdToUuid(isolatedOrderIdWithDiffClientId),
      ticker: testConstants.isolatedPerpetualMarket.ticker,
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

    it('Successfully gets multiple redis orders for parent subaccount', async () => {
      await Promise.all([
        placeOrder({
          redisOrder: redisTestConstants.defaultRedisOrder,
          client: redisClient,
        }),
        placeOrder({
          redisOrder: isolatedRedisOrder,
          client: redisClient,
        }),
      ]);

      const parentSubaccountNumber: number = 0;
      const queryParams = {
        address: testConstants.defaultSubaccount.address,
        parentSubaccountNumber,
      };

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/parentSubaccountNumber?${getQueryString(queryParams)}`,
      });

      expect(response.body).toEqual([ // by default sort by desc goodTilBlock
        redisOrderToResponseObject(isolatedRedisOrder),
        redisOrderToResponseObject(
          redisTestConstants.defaultRedisOrder,
        ),
      ]);

      const response2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/parentSubaccountNumber?${getQueryString({
          ...queryParams,
          returnLatestOrders: 'false',
        })}`,
      });

      expect(response2.body).toEqual([ // by default sort by desc goodTilBlock
        redisOrderToResponseObject(
          redisTestConstants.defaultRedisOrder,
        ),
        redisOrderToResponseObject(isolatedRedisOrder),
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
        'tickers across parent subaccount',
        [redisTestConstants.defaultRedisOrder, redisTestConstants.isolatedMarketRedisOrder],
        {
          ...defaultQueryParams,
          ticker: testConstants.defaultPerpetualMarket.ticker,
        },
        redisTestConstants.defaultRedisOrder,
      ],
      [
        'goodTilBlockBeforeOrAt',
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
        'goodTilBlockBeforeOrAt with isolated market',
        [
          redisTestConstants.defaultRedisOrder,
          {
            ...redisTestConstants.isolatedMarketRedisOrder,
            order: {
              ...redisTestConstants.isolatedMarketOrder,
              goodTilBlock: 1200,
            },
          },
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
        'goodTilBlockAfter',
        [
          redisTestConstants.defaultRedisOrder,
          redisOrderWithDifferentMarket,
          redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        ],
        {
          ...defaultQueryParams,
          goodTilBlockAfter: protocolTranslations.getGoodTilBlock(
            redisOrderWithDifferentMarket.order!,
          )! - 1,
        },
        redisOrderWithDifferentMarket,
      ],
      [
        'goodTilBlockAfter with isolated market',
        [
          redisTestConstants.defaultRedisOrder,
          {
            ...redisTestConstants.isolatedMarketRedisOrder,
            order: {
              ...redisTestConstants.isolatedMarketOrder,
              goodTilBlock: 1200,
            },
          },
        ],
        {
          ...defaultQueryParams,
          goodTilBlockAfter: protocolTranslations.getGoodTilBlock(
            redisTestConstants.defaultRedisOrder.order!,
          )! - 1,
        },
        redisTestConstants.defaultRedisOrder,
      ],
      [
        'goodTilBlockTimeBeforeOrAt',
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
      [
        'goodTilBlockTimeAfter',
        [
          redisTestConstants.defaultRedisOrder,
          redisTestConstants.defaultRedisOrderGoodTilBlockTime,
          newerRedisOrderGoodTilBlockTime,
        ],
        {
          ...defaultQueryParams,
          goodTilBlockTimeAfter: '2020-09-13T12:26:39.000Z',
        },
        newerRedisOrderGoodTilBlockTime,
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
        }, testConstants.defaultSubaccount.subaccountNumber),
        postgresOrderToResponseObject({
          ...secondOrderGoodTilBlockTime,
          id: getUuidForTest(secondOrderGoodTilBlockTime),
        },
        testConstants.defaultSubaccount.subaccountNumber),
        postgresOrderToResponseObject({
          ...secondOrder,
          id: getUuidForTest(secondOrder),
        }, testConstants.defaultSubaccount.subaccountNumber),
        postgresOrderToResponseObject({
          ...testConstants.defaultOrder,
          id: testConstants.defaultOrderId,
        }, testConstants.defaultSubaccount.subaccountNumber),
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
        }, testConstants.defaultSubaccount.subaccountNumber),
        postgresOrderToResponseObject({
          ...secondOrder,
          id: getUuidForTest(secondOrder),
        }, testConstants.defaultSubaccount.subaccountNumber),
        postgresOrderToResponseObject({
          ...secondOrderGoodTilBlockTime,
          id: getUuidForTest(secondOrderGoodTilBlockTime),
        }, testConstants.defaultSubaccount.subaccountNumber),
        postgresOrderToResponseObject({
          ...testConstants.defaultOrderGoodTilBlockTime,
          id: getUuidForTest(testConstants.defaultOrderGoodTilBlockTime),
        }, testConstants.defaultSubaccount.subaccountNumber),
      ]);
    });

    it('Successfully gets multiple postgres orders for parent subaccount', async () => {
      await Promise.all([
        OrderTable.create(testConstants.defaultOrder),
        OrderTable.create({
          ...testConstants.isolatedMarketOrder,
          goodTilBlock: '1000',
        }),
      ]);
      const parentSubaccountNumber: number = 0;
      const queryParams = {
        address: testConstants.defaultSubaccount.address,
        parentSubaccountNumber,
      };

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/parentSubaccountNumber?${getQueryString(queryParams)}`,
      });

      expect(response.body).toEqual([
        postgresOrderToResponseObject({
          ...testConstants.isolatedMarketOrder,
          id: testConstants.isolatedMarketOrderId,
          goodTilBlock: '1000',
        }, testConstants.isolatedSubaccount.subaccountNumber),
        postgresOrderToResponseObject({
          ...testConstants.defaultOrder,
          id: testConstants.defaultOrderId,
        }, testConstants.defaultSubaccount.subaccountNumber),
      ]);

      const response2: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/parentSubaccountNumber?${getQueryString({
          ...queryParams,
          returnLatestOrders: 'false',
        })}`,
      });

      expect(response2.body).toEqual([
        postgresOrderToResponseObject({
          ...testConstants.defaultOrder,
          id: testConstants.defaultOrderId,
        }, testConstants.defaultSubaccount.subaccountNumber),
        postgresOrderToResponseObject({
          ...testConstants.isolatedMarketOrder,
          id: testConstants.isolatedMarketOrderId,
          goodTilBlock: '1000',
        }, testConstants.isolatedSubaccount.subaccountNumber),
      ]);
    });

    it('Successfully returns filtered order when > limit orders exist', async () => {
      await Promise.all([
        OrderTable.create(testConstants.defaultOrder),
        OrderTable.create(untriggeredOrder),
      ]);

      const response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString(defaultQueryParams)}`,
      });

      expect(response.body).toEqual([
        postgresOrderToResponseObject({
          ...testConstants.defaultOrder,
          id: testConstants.defaultOrderId,
        }, testConstants.defaultSubaccount.subaccountNumber),
        postgresOrderToResponseObject({
          ...untriggeredOrder,
          id: untriggeredOrderId,
        }, testConstants.defaultSubaccount.subaccountNumber),
      ]);

      const response2 = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders?${getQueryString({
          ...defaultQueryParams,
          status: APIOrderStatusEnum.UNTRIGGERED,
          limit: 1,
        })}`,
      });

      expect(response2.body).toEqual([
        postgresOrderToResponseObject({
          ...untriggeredOrder,
          id: untriggeredOrderId,
        }, testConstants.defaultSubaccount.subaccountNumber),
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
          {
            [testConstants.defaultSubaccountId]:
              testConstants.defaultSubaccount.subaccountNumber,
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
        }, testConstants.defaultSubaccount.subaccountNumber),
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
          id: untriggeredOrderId,
        }, testConstants.defaultSubaccount.subaccountNumber),
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
          {
            [testConstants.defaultSubaccountId]:
              testConstants.defaultSubaccount.subaccountNumber,
          },
          redisTestConstants.defaultRedisOrder,
        ),
        postgresOrderToResponseObject({
          ...untriggeredOrder,
          id: untriggeredOrderId,
        }, testConstants.defaultSubaccount.subaccountNumber),
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
        }, testConstants.defaultSubaccount.subaccountNumber),
        redisOrderToResponseObject(newerRedisOrderGoodTilBlockTime),
        postgresAndRedisOrderToResponseObject(
          {
            ...testConstants.defaultOrderGoodTilBlockTime,
            id: getUuidForTest(testConstants.defaultOrderGoodTilBlockTime),
          },
          {
            [testConstants.defaultSubaccountId]:
              testConstants.defaultSubaccount.subaccountNumber,
          },
          redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        ),
        postgresOrderToResponseObject({
          ...secondOrder,
          id: getUuidForTest(secondOrder),
        }, testConstants.defaultSubaccount.subaccountNumber),
        redisOrderToResponseObject(redisOrderWithDifferentMarket),
        postgresAndRedisOrderToResponseObject(
          {
            ...testConstants.defaultOrder,
            id: testConstants.defaultOrderId,
          },
          {
            [testConstants.defaultSubaccountId]:
              testConstants.defaultSubaccount.subaccountNumber,
          },
          redisTestConstants.defaultRedisOrder,
        ),
      ]);
    });

    it('Successfully pulls both redis and postgres orders for parent subaccount', async () => {
      await Promise.all([
        OrderTable.create(testConstants.defaultOrder),
        OrderTable.create(secondOrder),
        OrderTable.create(testConstants.isolatedMarketOrder),
        placeOrder({
          redisOrder: redisTestConstants.defaultRedisOrder,
          client: redisClient,
        }),
        placeOrder({
          redisOrder: redisTestConstants.isolatedMarketRedisOrder,
          client: redisClient,
        }),
        placeOrder({
          redisOrder: isolatedRedisOrderWithDiffClientId,
          client: redisClient,
        }),
      ]);

      const parentSubaccountNumber: number = 0;
      const queryParams = {
        address: testConstants.defaultSubaccount.address,
        parentSubaccountNumber,
      };

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/orders/parentSubaccountNumber?${getQueryString(queryParams)}`,
      });

      expect(response.body).toEqual(
        expect.arrayContaining([
          postgresOrderToResponseObject({
            ...secondOrder,
            id: getUuidForTest(secondOrder),
          }, testConstants.defaultSubaccount.subaccountNumber),
          redisOrderToResponseObject(isolatedRedisOrderWithDiffClientId),
          postgresAndRedisOrderToResponseObject(
            {
              ...testConstants.isolatedMarketOrder,
              id: testConstants.isolatedMarketOrderId,
            },
            {
              [testConstants.isolatedSubaccountId]:
                  testConstants.isolatedSubaccount.subaccountNumber,
            },
            redisTestConstants.isolatedMarketRedisOrder,
          ),
          postgresAndRedisOrderToResponseObject(
            {
              ...testConstants.defaultOrder,
              id: testConstants.defaultOrderId,
            },
            {
              [testConstants.defaultSubaccountId]:
                  testConstants.defaultSubaccount.subaccountNumber,
            },
            redisTestConstants.defaultRedisOrder,
          ),
        ]),
      );
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
