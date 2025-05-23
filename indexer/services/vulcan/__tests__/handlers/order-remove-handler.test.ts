import {
  logger,
  ParseMessageError,
  stats,
  STATS_FUNCTION_NAME,
  wrapBackgroundTask,
} from '@dydxprotocol-indexer/base';
import {
  defaultTime,
} from '../helpers/constants';
import { synchronizeWrapBackgroundTask } from '@dydxprotocol-indexer/dev';
import {
  ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
  producer,
  SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import {
  dbHelpers,
  BlockTable,
  OrderCreateObject,
  OrderbookMessageContents,
  OrderFromDatabase,
  OrderSide,
  OrderStatus,
  OrderTable,
  OrderType,
  perpetualMarketRefresher,
  protocolTranslations,
  SubaccountMessageContents,
  testConstants,
  testMocks,
  apiTranslations,
  TimeInForce, blockHeightRefresher,
} from '@dydxprotocol-indexer/postgres';
import {
  OrderbookLevelsCache,
  OrderData,
  OrdersCache,
  OrdersDataCache,
  placeOrder,
  redis,
  redisTestConstants,
  StateFilledQuantumsCache,
  SubaccountOrderIdsCache,
  updateOrder,
  CanceledOrderStatus,
} from '@dydxprotocol-indexer/redis';
import {
  OffChainUpdateV1,
  OrderbookMessage,
  IndexerOrderId,
  OrderRemoveV1,
  OrderRemovalReason,
  OrderRemoveV1_OrderRemovalStatus,
  RedisOrder,
  SubaccountMessage,
  protoTimestampToDate,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import { IHeaders, ProducerRecord } from 'kafkajs';
import { DateTime } from 'luxon';

import { OrderRemoveHandler } from '../../src/handlers/order-remove-handler';
import { OrderbookSide } from '../../src/lib/types';
import { redisClient } from '../../src/helpers/redis/redis-controller';

import {
  expectCanceledOrderStatus,
  expectOrderbookLevelCache,
  handleOrderUpdate,
} from '../helpers/helpers';
import { expectWebsocketOrderbookMessage, expectWebsocketSubaccountMessage } from '../helpers/websocket-helpers';
import { ORDER_FLAG_LONG_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import Long from 'long';
import config from '../../src/config';

jest.mock('@dydxprotocol-indexer/base', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/base'),
  wrapBackgroundTask: jest.fn(),
}));

describe('OrderRemoveHandler', () => {
  beforeAll(async () => {
    jest.useFakeTimers();
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await Promise.all([
      perpetualMarketRefresher.updatePerpetualMarkets(),
      blockHeightRefresher.updateBlockHeight(),
    ]);
    jest.spyOn(stats, 'timing');
    jest.spyOn(stats, 'increment');
    jest.spyOn(logger, 'info');
    jest.spyOn(logger, 'error');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    await redis.deleteAllAsync(redisClient);
    jest.resetAllMocks();
    config.SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_CANCELS_MISSING_ORDERS = false;
  });

  afterAll(async () => {
    jest.useRealTimers();
    await dbHelpers.teardown();
  });

  const defaultOrderRemove: OrderRemoveV1 = {
    removedOrderId: redisTestConstants.defaultOrderId,
    reason: OrderRemovalReason.ORDER_REMOVAL_REASON_EXPIRED,
    removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED,
  };

  const statefulCancelationOrderRemove: OrderRemoveV1 = {
    removedOrderId: redisTestConstants.defaultOrderIdGoodTilBlockTime,
    reason: OrderRemovalReason.ORDER_REMOVAL_REASON_USER_CANCELED,
    removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED,
  };

  const defaultQuantums = redisTestConstants.defaultOrder.quantums;
  const defaultSubticks = redisTestConstants.defaultOrder.subticks;
  const defaultPrice = protocolTranslations.subticksToPrice(
    defaultSubticks.toString(),
    testConstants.defaultPerpetualMarket,
  );
  const defaultSize = protocolTranslations.quantumsToHumanFixedString(
    defaultQuantums.toString(),
    testConstants.defaultPerpetualMarket.atomicResolution,
  );

  const dbOrderFok: OrderCreateObject = {
    ...testConstants.defaultOrder,
    timeInForce: TimeInForce.FOK,
  };
  const dbOrderIoc: OrderCreateObject = {
    ...testConstants.defaultOrder,
    timeInForce: TimeInForce.IOC,
  };

  const defaultKafkaHeaders: IHeaders = {
    message_received_timestamp: String(protoTimestampToDate(defaultTime)),
  };

  it.each([
    [
      {
        ...defaultOrderRemove,
        removedOrderId: undefined,
      },
      'OrderRemove must contain a removedOrderId',
    ],
    [
      {
        ...defaultOrderRemove,
        removedOrderId: {
          ...redisTestConstants.defaultOrderId,
          subaccountId: undefined,
        },
      },
      'OrderRemove must contain a removedOrderId.subaccountId',
    ],
    [
      {
        ...defaultOrderRemove,
        removedOrderId: {
          ...redisTestConstants.defaultOrderId,
          clientId: undefined,
        },
      },
      'OrderRemove must contain a removedOrderId.clientId',
    ],
    [
      {
        ...defaultOrderRemove,
        removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_UNSPECIFIED,
      },
      'OrderRemove removalStatus cannot be unspecified',
    ],
    [
      {
        ...defaultOrderRemove,
        reason: OrderRemovalReason.ORDER_REMOVAL_REASON_UNSPECIFIED,
      },
      'OrderRemove reason cannot be unspecified',
    ],
  ])('fails when update is invalid', async (orderRemoveJson: any, errorMessage: string) => {
    const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate(orderRemoveJson);

    const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
    await expect(orderRemoveHandler.handleUpdate(
      offChainUpdate,
      defaultKafkaHeaders,
    )).rejects.toThrow(
      new ParseMessageError(errorMessage),
    );
    expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
      at: 'OrderRemoveHandler#logAndThrowParseMessageError',
      message: errorMessage,
      orderRemove: orderRemoveJson,
    }));
  });

  describe('Order Remove Message - not a Stateful Cancelation', () => {
    it('successfully returns early if unable to find order in redis', async () => {
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate(defaultOrderRemove);

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      expect(logger.info).toHaveBeenCalledWith(expect.objectContaining({
        at: 'orderRemoveHandler#handleOrderRemoval',
        message: 'Unable to find order',
        orderId: defaultOrderRemove.removedOrderId,
      }));
      expect(logger.error).not.toHaveBeenCalled();
      expectTimingStats();
    });

    it('successfully sends subaccount websocket message and returns if unable to find order in redis', async () => {
      config.SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_CANCELS_MISSING_ORDERS = true;
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate(defaultOrderRemove);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      expect(logger.info).toHaveBeenCalledWith(expect.objectContaining({
        at: 'orderRemoveHandler#handleOrderRemoval',
        message: 'Unable to find order',
        orderId: defaultOrderRemove.removedOrderId,
      }));

      // Subaccounts message is sent
      const subaccountContents: SubaccountMessageContents = {
        orders: [
          {
            id: OrderTable.orderIdToUuid(redisTestConstants.defaultOrderId),
            subaccountId: testConstants.defaultSubaccountId,
            clientId: redisTestConstants.defaultOrderId.clientId.toString(),
            clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
            status: OrderStatus.CANCELED,
            orderFlags: redisTestConstants.defaultOrderId.orderFlags.toString(),
            ticker: redisTestConstants.defaultRedisOrder.ticker,
            removalReason: OrderRemovalReason[defaultOrderRemove.reason],
          },
        ],
        blockHeight: blockHeightRefresher.getLatestBlockHeight(),
      };
      expectWebsocketMessagesSent(
        producerSendSpy,
        SubaccountMessage.fromPartial({
          contents: JSON.stringify(subaccountContents),
          subaccountId: redisTestConstants.defaultSubaccountId,
          version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
        }),
      );

      expect(logger.error).not.toHaveBeenCalled();
      expectTimingStats();
    });

    it('successfully sends subaccount websocket message with db order fields if unable to find order in redis',
      async () => {
        config.SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_CANCELS_MISSING_ORDERS = true;
        await OrderTable.create(testConstants.defaultOrder);
        const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate(defaultOrderRemove);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

        const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
        await orderRemoveHandler.handleUpdate(
          offChainUpdate,
          defaultKafkaHeaders,
        );

        expect(logger.info).toHaveBeenCalledWith(expect.objectContaining({
          at: 'orderRemoveHandler#handleOrderRemoval',
          message: 'Unable to find order',
          orderId: defaultOrderRemove.removedOrderId,
        }));

        // Subaccounts message is sent
        const subaccountContents: SubaccountMessageContents = {
          orders: [
            {
              id: OrderTable.orderIdToUuid(redisTestConstants.defaultOrderId),
              subaccountId: testConstants.defaultSubaccountId,
              clientId: redisTestConstants.defaultOrderId.clientId.toString(),
              clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
              status: OrderStatus.CANCELED,
              orderFlags: redisTestConstants.defaultOrderId.orderFlags.toString(),
              ticker: redisTestConstants.defaultRedisOrder.ticker,
              removalReason: OrderRemovalReason[defaultOrderRemove.reason],
              updatedAt: testConstants.defaultOrder.updatedAt,
              updatedAtHeight: testConstants.defaultOrder.updatedAtHeight,
              price: testConstants.defaultOrder.price,
              size: testConstants.defaultOrder.size,
              clientMetadata: testConstants.defaultOrder.clientMetadata,
              side: testConstants.defaultOrder.side,
              timeInForce: apiTranslations.orderTIFToAPITIF(testConstants.defaultOrder.timeInForce),
              totalFilled: testConstants.defaultOrder.totalFilled,
              goodTilBlock: testConstants.defaultOrder.goodTilBlock,
              type: testConstants.defaultOrder.type,
            },
          ],
          blockHeight: blockHeightRefresher.getLatestBlockHeight(),
        };
        expectWebsocketMessagesSent(
          producerSendSpy,
          SubaccountMessage.fromPartial({
            contents: JSON.stringify(subaccountContents),
            subaccountId: redisTestConstants.defaultSubaccountId,
            version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
          }),
        );

        expect(logger.error).not.toHaveBeenCalled();
        expectTimingStats();
      });

    it('successfully returns early if unable to find perpetualMarket', async () => {
      await Promise.all([
        dbHelpers.clearData(),
        placeOrder({ redisOrder: redisTestConstants.defaultRedisOrder, client: redisClient }),
      ]);

      await perpetualMarketRefresher.updatePerpetualMarkets();

      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate(defaultOrderRemove);
      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      const clobPairId: string = testConstants.defaultPerpetualMarket.clobPairId;
      expect(logger.error).toHaveBeenCalledWith({
        at: 'orderRemoveHandler#handle',
        message: `Unable to find perpetual market with clobPairId: ${clobPairId}`,
      });
      expectTimingStats();
    });

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrderId,
        testConstants.defaultOrder,
        redisTestConstants.defaultRedisOrder,
        redisTestConstants.defaultOrderUuid,
        true,
        undefined,
      ],
      [
        'goodTilBlockTime',
        redisTestConstants.defaultOrderIdGoodTilBlockTime,
        testConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        true,
        undefined,
      ],
      [
        'conditional',
        redisTestConstants.defaultOrderIdConditional,
        testConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
        true,
        testConstants.defaultConditionalOrder.triggerPrice,
      ],
      [
        'Fill-or-Kill',
        redisTestConstants.defaultOrderId,
        dbOrderFok,
        redisTestConstants.defaultRedisOrderFok,
        redisTestConstants.defaultOrderUuid,
        false,
        undefined,
      ],
      [
        'Immediate-or-Cancel',
        redisTestConstants.defaultOrderId,
        dbOrderIoc,
        redisTestConstants.defaultRedisOrderIoc,
        redisTestConstants.defaultOrderUuid,
        false,
        undefined,
      ],
    ])('successfully removes order (with %s)', async (
      _name: string,
      removedOrderId: IndexerOrderId,
      removedOrder: OrderCreateObject,
      removedRedisOrder: RedisOrder,
      expectedOrderUuid: string,
      expectOrderbookUpdate: boolean,
      triggerPrice?: string,
    ) => {
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
        ...defaultOrderRemove,
        removedOrderId,
      });
      const orderbookLevel: string = Big(defaultQuantums.toString()).times(2).toString();

      await Promise.all([
        placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
        OrderbookLevelsCache.updatePriceLevel(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          orderbookLevel,
          redisClient,
        ),
      ]);

      await Promise.all([
        // Must be done separately so that the subaccount and perpetualMarket have been created
        // before the order
        OrderTable.create(removedOrder),
        // Must be done after adding orders to all caches to overwrite the ordersDataCache
        setOrderToRestingOnOrderbook(removedRedisOrder),
      ]);

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      // orderbook level reduced by defaultQuantums
      const remainingOrderbookLevel: string = Big(
        orderbookLevel,
      ).minus(
        defaultQuantums.toString(),
      ).toString();
      await Promise.all([
        expectOrderStatus(expectedOrderUuid, OrderStatus.CANCELED),
        expectOrderbookLevelCache(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          expectOrderbookUpdate ? remainingOrderbookLevel : orderbookLevel,
        ),
        expectOrdersCacheEmpty(expectedOrderUuid),
        expectOrdersDataCacheEmpty(removedOrderId),
        expectSubaccountsOrderIdsCacheEmpty(redisTestConstants.defaultSubaccountUuid),
        expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.CANCELED),
      ]);

      // Subaccounts message is sent first followed by orderbooks message
      const subaccountContents: SubaccountMessageContents = {
        orders: [
          {
            id: expectedOrderUuid,
            subaccountId: testConstants.defaultSubaccountId,
            clientId: removedOrderId.clientId.toString(),
            clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
            side: OrderSide.BUY,
            size: defaultSize,
            totalOptimisticFilled: '0',
            price: defaultPrice,
            type: protocolTranslations.protocolConditionTypeToOrderType(
              removedRedisOrder.order!.conditionType,
            ),
            status: OrderStatus.CANCELED,
            timeInForce: apiTranslations.orderTIFToAPITIF(
              protocolTranslations.protocolOrderTIFToTIF(removedRedisOrder.order!.timeInForce),
            ),
            postOnly: apiTranslations.isOrderTIFPostOnly(
              protocolTranslations.protocolOrderTIFToTIF(removedRedisOrder.order!.timeInForce),
            ),
            reduceOnly: removedRedisOrder.order!.reduceOnly,
            orderFlags: removedRedisOrder.order!.orderId!.orderFlags.toString(),
            goodTilBlock: protocolTranslations.getGoodTilBlock(removedRedisOrder.order!)
              ?.toString(),
            goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(removedRedisOrder.order!),
            ticker: redisTestConstants.defaultRedisOrder.ticker,
            removalReason: OrderRemovalReason[defaultOrderRemove.reason],
            createdAtHeight: removedOrder.createdAtHeight,
            updatedAt: removedOrder.updatedAt,
            updatedAtHeight: removedOrder.updatedAtHeight,
            clientMetadata: removedRedisOrder.order!.clientMetadata.toString(),
            triggerPrice,
          },
        ],
        blockHeight: blockHeightRefresher.getLatestBlockHeight(),
      };
      const orderbookContents: OrderbookMessageContents = {
        [OrderbookSide.BIDS]: [[
          defaultPrice,
          protocolTranslations.quantumsToHuman(
            remainingOrderbookLevel,
            testConstants.defaultPerpetualMarket.atomicResolution,
          ).toString(),
        ]],
      };
      expectWebsocketMessagesSent(
        producerSendSpy,
        SubaccountMessage.fromPartial({
          contents: JSON.stringify(subaccountContents),
          subaccountId: redisTestConstants.defaultSubaccountId,
          version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
        }),
        expectOrderbookUpdate
          ? OrderbookMessage.fromPartial({
            contents: JSON.stringify(orderbookContents),
            clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
            version: ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
          }) : undefined,
      );
      expect(logger.error).not.toHaveBeenCalled();
      expectTimingStats(true, true);
    });

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrderId,
        testConstants.defaultOrder,
        redisTestConstants.defaultRedisOrder,
        redisTestConstants.defaultOrderUuid,
        undefined,
      ],
      [
        'goodTilBlockTime',
        redisTestConstants.defaultOrderIdGoodTilBlockTime,
        testConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        undefined,
      ],
      [
        'conditional',
        redisTestConstants.defaultOrderIdConditional,
        testConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
        testConstants.defaultConditionalOrder.triggerPrice,
      ],
    ])('successfully removes order (with %s) and can set reason to BEST_EFFORT_CANCELED', async (
      _name: string,
      removedOrderId: IndexerOrderId,
      removedOrder: OrderCreateObject,
      removedRedisOrder: RedisOrder,
      expectedOrderUuid: string,
      triggerPrice?: string,
    ) => {
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
        ...defaultOrderRemove,
        removedOrderId,
        removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
      });

      await Promise.all([
        placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
        OrderbookLevelsCache.updatePriceLevel(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          defaultQuantums.toString(),
          redisClient,
        ),
      ]);

      await Promise.all([
        // Must be done separately so that the subaccount and perpetualMarket have been created
        // before the order
        OrderTable.create(removedOrder),
        // Must be done after adding orders to all caches to overwrite the ordersDataCache
        setOrderToRestingOnOrderbook(removedRedisOrder),
      ]);

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      await Promise.all([
        expectOrderStatus(expectedOrderUuid, OrderStatus.BEST_EFFORT_CANCELED),
        // default quantums - default quantums = 0
        expectOrderbookLevelCache(removedRedisOrder.ticker, OrderSide.BUY, defaultPrice, '0'),
        expectOrdersCacheEmpty(expectedOrderUuid),
        expectOrdersDataCacheEmpty(removedOrderId),
        expectSubaccountsOrderIdsCacheEmpty(redisTestConstants.defaultSubaccountUuid),
        expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.BEST_EFFORT_CANCELED),
      ]);

      // Subaccounts message is sent first followed by orderbooks message
      const subaccountContents: SubaccountMessageContents = {
        orders: [
          {
            id: expectedOrderUuid,
            subaccountId: testConstants.defaultSubaccountId,
            clientId: removedOrderId.clientId.toString(),
            clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
            side: OrderSide.BUY,
            size: defaultSize,
            totalOptimisticFilled: '0',
            price: defaultPrice,
            type: protocolTranslations.protocolConditionTypeToOrderType(
              removedRedisOrder.order!.conditionType,
            ),
            status: OrderStatus.BEST_EFFORT_CANCELED,
            timeInForce: apiTranslations.orderTIFToAPITIF(
              protocolTranslations.protocolOrderTIFToTIF(removedRedisOrder.order!.timeInForce),
            ),
            postOnly: apiTranslations.isOrderTIFPostOnly(
              protocolTranslations.protocolOrderTIFToTIF(removedRedisOrder.order!.timeInForce),
            ),
            reduceOnly: removedRedisOrder.order!.reduceOnly,
            orderFlags: removedRedisOrder.order!.orderId!.orderFlags.toString(),
            goodTilBlock: protocolTranslations.getGoodTilBlock(removedRedisOrder.order!)
              ?.toString(),
            goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(removedRedisOrder.order!),
            ticker: redisTestConstants.defaultRedisOrder.ticker,
            removalReason: OrderRemovalReason[defaultOrderRemove.reason],
            createdAtHeight: removedOrder.createdAtHeight,
            updatedAt: removedOrder.updatedAt,
            updatedAtHeight: removedOrder.updatedAtHeight,
            clientMetadata: removedRedisOrder.order!.clientMetadata.toString(),
            triggerPrice,
          },
        ],
        blockHeight: blockHeightRefresher.getLatestBlockHeight(),
      };

      const orderbookContents: OrderbookMessageContents = {
        [OrderbookSide.BIDS]: [[
          defaultPrice,
          '0',
        ]],
      };
      expectWebsocketMessagesSent(
        producerSendSpy,
        SubaccountMessage.fromPartial({
          contents: JSON.stringify(subaccountContents),
          subaccountId: redisTestConstants.defaultSubaccountId,
          version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
        }),
        OrderbookMessage.fromPartial({
          contents: JSON.stringify(orderbookContents),
          clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
          version: ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
        }),
      );
      expect(logger.error).not.toHaveBeenCalled();
      expectTimingStats(true, true);
    });

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrderId,
        testConstants.defaultOrder,
        redisTestConstants.defaultRedisOrder,
        redisTestConstants.defaultOrderUuid,
        undefined,
      ],
      [
        'goodTilBlockTime',
        redisTestConstants.defaultOrderIdGoodTilBlockTime,
        testConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        undefined,
      ],
      [
        'conditional',
        redisTestConstants.defaultOrderIdConditional,
        testConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
        testConstants.defaultConditionalOrder.triggerPrice,
      ],
    ])(
      'successfully removes order (with %s) and does not change orderbookLevelsCache when order is on book',
      async (
        _name: string,
        removedOrderId: IndexerOrderId,
        removedOrder: OrderCreateObject,
        removedRedisOrder: RedisOrder,
        expectedOrderUuid: string,
        triggerPrice?: string,
      ) => {
        const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
          ...defaultOrderRemove,
          removedOrderId,
        });

        await Promise.all([
          placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
          OrderbookLevelsCache.updatePriceLevel(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            defaultPrice,
            defaultQuantums.toString(),
            redisClient,
          ),
        ]);

        // Must be done separately so that the subaccount and perpetualMarket have been created
        // before the order
        await Promise.all([
          OrderTable.create(removedOrder),
        ]);

        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

        const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
        await orderRemoveHandler.handleUpdate(
          offChainUpdate,
          defaultKafkaHeaders,
        );

        await Promise.all([
          expectOrderStatus(expectedOrderUuid, OrderStatus.CANCELED),
          // orderbook should not be affected, so it will be set to defaultQuantums
          expectOrderbookLevelCache(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            redisTestConstants.defaultPrice,
            defaultQuantums.toString(),
          ),
          expectOrdersCacheEmpty(expectedOrderUuid),
          expectOrdersDataCacheEmpty(removedOrderId),
          expectSubaccountsOrderIdsCacheEmpty(redisTestConstants.defaultSubaccountUuid),
          expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.CANCELED),
        ]);

        // Subaccounts message is sent first followed by orderbooks message
        const subaccountContents: SubaccountMessageContents = {
          orders: [{
            id: expectedOrderUuid,
            subaccountId: testConstants.defaultSubaccountId,
            clientId: removedOrderId.clientId.toString(),
            clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
            side: OrderSide.BUY,
            size: defaultSize,
            totalOptimisticFilled: '0',
            price: defaultPrice,
            type: protocolTranslations.protocolConditionTypeToOrderType(
              removedRedisOrder.order!.conditionType,
            ),
            status: OrderStatus.CANCELED,
            timeInForce: apiTranslations.orderTIFToAPITIF(
              protocolTranslations.protocolOrderTIFToTIF(removedRedisOrder.order!.timeInForce),
            ),
            postOnly: apiTranslations.isOrderTIFPostOnly(
              protocolTranslations.protocolOrderTIFToTIF(removedRedisOrder.order!.timeInForce),
            ),
            reduceOnly: removedRedisOrder.order!.reduceOnly,
            orderFlags: removedRedisOrder.order!.orderId!.orderFlags.toString(),
            goodTilBlock: protocolTranslations.getGoodTilBlock(removedRedisOrder.order!)
              ?.toString(),
            goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(removedRedisOrder.order!),
            ticker: redisTestConstants.defaultRedisOrder.ticker,
            removalReason: OrderRemovalReason[defaultOrderRemove.reason],
            createdAtHeight: removedOrder.createdAtHeight,
            updatedAt: removedOrder.updatedAt,
            updatedAtHeight: removedOrder.updatedAtHeight,
            clientMetadata: removedRedisOrder.order!.clientMetadata.toString(),
            triggerPrice,
          }],
          blockHeight: blockHeightRefresher.getLatestBlockHeight(),
        };
        expectWebsocketMessagesSent(
          producerSendSpy,
          SubaccountMessage.fromPartial({
            contents: JSON.stringify(subaccountContents),
            subaccountId: redisTestConstants.defaultSubaccountId,
            version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
          }),
          undefined,
        );

        expect(logger.error).not.toHaveBeenCalled();
        expectTimingStats(true, true);
      },
    );

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrderId,
        testConstants.defaultOrder,
        redisTestConstants.defaultRedisOrder,
        redisTestConstants.defaultOrderUuid,
        undefined,
      ],
      [
        'goodTilBlockTime',
        redisTestConstants.defaultOrderIdGoodTilBlockTime,
        testConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        undefined,
      ],
      [
        'conditional',
        redisTestConstants.defaultOrderIdConditional,
        testConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
        testConstants.defaultConditionalOrder.triggerPrice,
      ],
    ])(
      'does not increase orderbook level if total filled > quantums of order (with %s)',
      async (
        _name: string,
        removedOrderId: IndexerOrderId,
        removedOrder: OrderCreateObject,
        removedRedisOrder: RedisOrder,
        expectedOrderUuid: string,
        triggerPrice?: string,
      ) => {
        const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
          ...defaultOrderRemove,
          removedOrderId,
        });

        await Promise.all([
          placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
          OrderbookLevelsCache.updatePriceLevel(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            defaultPrice,
            defaultQuantums.toString(),
            redisClient,
          ),
        ]);

        const exceedsFilledUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage = {
          orderPlace: undefined,
          orderRemove: undefined,
          orderUpdate: {
            orderId: removedRedisOrder.order!.orderId!,
            totalFilledQuantums: removedRedisOrder.order!.quantums.add(Long.fromValue(100, true)),
          },
        };
        await handleOrderUpdate(exceedsFilledUpdate);

        // Must be done separately so that the subaccount and perpetualMarket have been created
        // before the order
        await Promise.all([
          OrderTable.create(removedOrder),
        ]);

        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

        const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
        await orderRemoveHandler.handleUpdate(
          offChainUpdate,
          defaultKafkaHeaders,
        );

        await Promise.all([
          expectOrderStatus(expectedOrderUuid, OrderStatus.CANCELED),
          // orderbook should not be affected, so it will be set to defaultQuantums
          expectOrderbookLevelCache(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            redisTestConstants.defaultPrice,
            defaultQuantums.toString(),
          ),
          expectOrdersCacheEmpty(expectedOrderUuid),
          expectOrdersDataCacheEmpty(removedOrderId),
          expectSubaccountsOrderIdsCacheEmpty(redisTestConstants.defaultSubaccountUuid),
          expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.CANCELED),
        ]);

        // Subaccounts message is sent first followed by orderbooks message
        const subaccountContents: SubaccountMessageContents = {
          orders: [{
            id: expectedOrderUuid,
            subaccountId: testConstants.defaultSubaccountId,
            clientId: removedOrderId.clientId.toString(),
            clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
            side: OrderSide.BUY,
            size: defaultSize,
            // Check that the total filled was > than quantums
            totalOptimisticFilled: '0.00010001',
            price: defaultPrice,
            type: protocolTranslations.protocolConditionTypeToOrderType(
              removedRedisOrder.order!.conditionType,
            ),
            status: OrderStatus.CANCELED,
            timeInForce: apiTranslations.orderTIFToAPITIF(
              protocolTranslations.protocolOrderTIFToTIF(removedRedisOrder.order!.timeInForce),
            ),
            postOnly: apiTranslations.isOrderTIFPostOnly(
              protocolTranslations.protocolOrderTIFToTIF(removedRedisOrder.order!.timeInForce),
            ),
            reduceOnly: removedRedisOrder.order!.reduceOnly,
            orderFlags: removedRedisOrder.order!.orderId!.orderFlags.toString(),
            goodTilBlock: protocolTranslations.getGoodTilBlock(removedRedisOrder.order!)
              ?.toString(),
            goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(removedRedisOrder.order!),
            ticker: redisTestConstants.defaultRedisOrder.ticker,
            removalReason: OrderRemovalReason[defaultOrderRemove.reason],
            createdAtHeight: removedOrder.createdAtHeight,
            updatedAt: removedOrder.updatedAt,
            updatedAtHeight: removedOrder.updatedAtHeight,
            clientMetadata: removedRedisOrder.order!.clientMetadata.toString(),
            triggerPrice,
          }],
          blockHeight: blockHeightRefresher.getLatestBlockHeight(),
        };
        expectWebsocketMessagesSent(
          producerSendSpy,
          SubaccountMessage.fromPartial({
            contents: JSON.stringify(subaccountContents),
            subaccountId: redisTestConstants.defaultSubaccountId,
            version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
          }),
          // no orderbook message because no change in orderbook levels
          undefined,
        );

        expect(logger.error).not.toHaveBeenCalled();
        expectTimingStats(true, true);
      },
    );

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrderId,
        {
          ...testConstants.defaultOrder,
          status: OrderStatus.FILLED,
        },
        redisTestConstants.defaultRedisOrder,
        redisTestConstants.defaultOrderUuid,
      ],
      [
        'goodTilBlockTime',
        redisTestConstants.defaultOrderIdGoodTilBlockTime,
        {
          ...testConstants.defaultOrderGoodTilBlockTime,
          status: OrderStatus.FILLED,
        },
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
      ],
      [
        'conditional',
        redisTestConstants.defaultOrderIdConditional,
        {
          ...testConstants.defaultConditionalOrder,
          status: OrderStatus.FILLED,
        },
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
      ],
    ])(
      'does not send subaccount message for orders fully-filled in state for best effort ' +
      'user cancel (with %s)',
      async (
        _name: string,
        removedOrderId: IndexerOrderId,
        removedOrder: OrderCreateObject,
        removedRedisOrder: RedisOrder,
        expectedOrderUuid: string,
      ) => {
        const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
          ...defaultOrderRemove,
          removedOrderId,
          removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
          reason: OrderRemovalReason.ORDER_REMOVAL_REASON_USER_CANCELED,
        });

        await Promise.all([
          placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
          OrderbookLevelsCache.updatePriceLevel(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            defaultPrice,
            defaultQuantums.toString(),
            redisClient,
          ),
          StateFilledQuantumsCache.updateStateFilledQuantums(
            expectedOrderUuid,
            removedRedisOrder.order!.quantums.toString(),
            redisClient,
          ),
        ]);

        const fullyFilledUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage = {
          orderPlace: undefined,
          orderRemove: undefined,
          orderUpdate: {
            orderId: removedRedisOrder.order!.orderId!,
            totalFilledQuantums: removedRedisOrder.order!.quantums,
          },
        };
        await handleOrderUpdate(fullyFilledUpdate);

        // Must be done separately so that the subaccount and perpetualMarket have been created
        // before the order
        await Promise.all([
          OrderTable.create(removedOrder),
        ]);

        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

        const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
        await orderRemoveHandler.handleUpdate(
          offChainUpdate,
          defaultKafkaHeaders,
        );

        await Promise.all([
          expectOrderStatus(expectedOrderUuid, removedOrder.status),
          // orderbook should not be affected, so it will be set to defaultQuantums
          expectOrderbookLevelCache(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            redisTestConstants.defaultPrice,
            defaultQuantums.toString(),
          ),
          expectOrdersCacheEmpty(expectedOrderUuid),
          expectOrdersDataCacheEmpty(removedOrderId),
          expectSubaccountsOrderIdsCacheEmpty(redisTestConstants.defaultSubaccountUuid),
          expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.BEST_EFFORT_CANCELED),
        ]);

        // no orderbook message because no change in orderbook levels
        expectNoWebsocketMessagesSent(producerSendSpy);
        expect(logger.error).not.toHaveBeenCalled();
        expectTimingStats(true, false);
      },
    );

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrderId,
        testConstants.defaultOrder,
        redisTestConstants.defaultRedisOrder,
        redisTestConstants.defaultOrderUuid,
      ],
      [
        'goodTilBlockTime',
        redisTestConstants.defaultOrderIdGoodTilBlockTime,
        testConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
      ],
      [
        'conditional',
        redisTestConstants.defaultOrderIdConditional,
        testConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
      ],
    ])(
      'does not send subaccount message for removals with fully-filled reason (with %s)',
      async (
        _name: string,
        removedOrderId: IndexerOrderId,
        removedOrder: OrderCreateObject,
        removedRedisOrder: RedisOrder,
        expectedOrderUuid: string,
      ) => {
        const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
          ...defaultOrderRemove,
          removedOrderId,
          removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_FILLED,
          reason: OrderRemovalReason.ORDER_REMOVAL_REASON_FULLY_FILLED,
        });

        await Promise.all([
          placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
          OrderbookLevelsCache.updatePriceLevel(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            defaultPrice,
            defaultQuantums.toString(),
            redisClient,
          ),
        ]);

        // Must be done separately so that the subaccount and perpetualMarket have been created
        // before the order
        await Promise.all([
          OrderTable.create(removedOrder),
        ]);

        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

        const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
        await orderRemoveHandler.handleUpdate(
          offChainUpdate,
          defaultKafkaHeaders,
        );

        await Promise.all([
          expectOrderStatus(expectedOrderUuid, OrderStatus.FILLED),
          // orderbook should not be affected, so it will be set to defaultQuantums
          expectOrderbookLevelCache(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            redisTestConstants.defaultPrice,
            defaultQuantums.toString(),
          ),
          expectOrdersCacheEmpty(expectedOrderUuid),
          expectOrdersDataCacheEmpty(removedOrderId),
          expectSubaccountsOrderIdsCacheEmpty(redisTestConstants.defaultSubaccountUuid),
          expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.NOT_CANCELED),
        ]);

        // no orderbook message because no change in orderbook levels
        expectNoWebsocketMessagesSent(producerSendSpy);
        expect(logger.error).not.toHaveBeenCalled();
        expectTimingStats(true, true);
      },
    );
  });

  describe('Order Remove Message - Stateful Order Cancelation', () => {
    it('logs an error if the order cannot be found in Postgres', async () => {
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
        ...statefulCancelationOrderRemove,
        removedOrderId: redisTestConstants.defaultOrderIdGoodTilBlockTime,
      });

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      const statefulOrderCancelSpy = jest.spyOn(orderRemoveHandler as any, 'handleStatefulOrderCancelation');
      const orderRemovalSpy = jest.spyOn(orderRemoveHandler as any, 'handleOrderRemoval');
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      expect(producerSendSpy).not.toHaveBeenCalled();
      expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
        at: 'orderRemoveHandler#handleStatefulOrderCancelation',
        message: expect.stringContaining('Could not find order for stateful order cancelation'),
        orderRemove: statefulCancelationOrderRemove,
      }));
      expect(statefulOrderCancelSpy).toHaveBeenCalled();
      expect(orderRemovalSpy).not.toHaveBeenCalled();
    });

    it('calls order removal instead of stateful order cancellation for vault orders', async () => {
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
        ...statefulCancelationOrderRemove,
        removedOrderId: redisTestConstants.defaultOrderIdVault,
      });

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      const statefulOrderCancelSpy = jest.spyOn(orderRemoveHandler as any, 'handleStatefulOrderCancelation');
      const orderRemovalSpy = jest.spyOn(orderRemoveHandler as any, 'handleOrderRemoval');
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      expect(statefulOrderCancelSpy).not.toHaveBeenCalled();
      expect(orderRemovalSpy).toHaveBeenCalled();
    });

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrderIdGoodTilBlockTime,
        testConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        undefined,
      ],
      [
        'conditional',
        redisTestConstants.defaultOrderIdConditional,
        testConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
        testConstants.defaultConditionalOrder.triggerPrice,
      ],
    ])('sends subaccount websocket message if order is not redis (with %s)', async (
      _name: string,
      removedOrderId: IndexerOrderId,
      removedOrder: OrderCreateObject,
      removedRedisOrder: RedisOrder,
      expectedOrderUuid: string,
      triggerPrice?: string,
    ) => {
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
        ...statefulCancelationOrderRemove,
        removedOrderId,
      });

      await Promise.all([
        OrderTable.create(removedOrder),
      ]);

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      // Subaccounts message is sent first followed by orderbooks message
      const subaccountContents: SubaccountMessageContents = {
        orders: [{
          id: expectedOrderUuid,
          subaccountId: testConstants.defaultSubaccountId,
          clientId: removedOrderId.clientId.toString(),
          clobPairId: testConstants.defaultOrderGoodTilBlockTime.clobPairId,
          side: OrderSide.BUY,
          size: removedOrder.size,
          totalFilled: '0',
          price: removedOrder.price,
          type: protocolTranslations.protocolConditionTypeToOrderType(
            removedRedisOrder.order!.conditionType,
          ),
          status: OrderStatus.CANCELED,
          timeInForce: apiTranslations.orderTIFToAPITIF(removedOrder.timeInForce),
          postOnly: apiTranslations.isOrderTIFPostOnly(removedOrder.timeInForce),
          reduceOnly: removedOrder.reduceOnly,
          orderFlags: removedOrder.orderFlags,
          goodTilBlockTime: removedOrder.goodTilBlockTime,
          ticker: removedRedisOrder.ticker,
          removalReason: OrderRemovalReason[statefulCancelationOrderRemove.reason],
          createdAtHeight: removedOrder.createdAtHeight,
          updatedAt: removedOrder.updatedAt,
          updatedAtHeight: removedOrder.updatedAtHeight,
          clientMetadata: removedOrder.clientMetadata.toString(),
          triggerPrice,
        }],
        blockHeight: blockHeightRefresher.getLatestBlockHeight(),
      };
      expectWebsocketMessagesSent(
        producerSendSpy,
        SubaccountMessage.fromPartial({
          contents: JSON.stringify(subaccountContents),
          subaccountId: redisTestConstants.defaultSubaccountId,
          version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
        }),
        undefined,
      );

      expect(logger.error).not.toHaveBeenCalled();
      expectTimingStats(true, false, false, true);
    });

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrderIdGoodTilBlockTime,
        testConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        undefined,
      ],
      [
        'conditional',
        redisTestConstants.defaultOrderIdConditional,
        testConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
        testConstants.defaultConditionalOrder.triggerPrice,
      ],
    ])('successfully removes stateful order, not resting on book', async (
      _name: string,
      removedOrderId: IndexerOrderId,
      removedOrder: OrderCreateObject,
      removedRedisOrder: RedisOrder,
      expectedOrderUuid: string,
      triggerPrice?: string,
    ) => {
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
        ...statefulCancelationOrderRemove,
        removedOrderId,
      });

      await Promise.all([
        placeOrder(
          {
            redisOrder: removedRedisOrder,
            client: redisClient,
          }),
        OrderbookLevelsCache.updatePriceLevel(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          defaultQuantums.toString(),
          redisClient,
        ),
      ]);

      // Must be done separately so that the subaccount and perpetualMarket have been created
      // before the order
      await Promise.all([
        OrderTable.create(removedOrder),
      ]);

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      await Promise.all([
        // orderbook should not be affected, so it will be set to defaultQuantums
        expectOrderbookLevelCache(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          redisTestConstants.defaultPrice,
          defaultQuantums.toString(),
        ),
        expectOrdersCacheEmpty(expectedOrderUuid),
        expectOrdersDataCacheEmpty(removedOrderId),
        expectSubaccountsOrderIdsCacheEmpty(redisTestConstants.defaultSubaccountUuid),
        expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.NOT_CANCELED),
      ]);

      // Subaccounts message is sent first followed by orderbooks message
      const subaccountContents: SubaccountMessageContents = {
        orders: [{
          id: expectedOrderUuid,
          subaccountId: testConstants.defaultSubaccountId,
          clientId: removedOrderId.clientId.toString(),
          clobPairId: testConstants.defaultOrderGoodTilBlockTime.clobPairId,
          side: OrderSide.BUY,
          size: removedOrder.size,
          totalFilled: '0',
          price: removedOrder.price,
          type: protocolTranslations.protocolConditionTypeToOrderType(
            removedRedisOrder.order!.conditionType,
          ),
          status: OrderStatus.CANCELED,
          timeInForce: apiTranslations.orderTIFToAPITIF(removedOrder.timeInForce),
          postOnly: apiTranslations.isOrderTIFPostOnly(removedOrder.timeInForce),
          reduceOnly: removedOrder.reduceOnly,
          orderFlags: removedOrder.orderFlags,
          goodTilBlockTime: removedOrder.goodTilBlockTime,
          ticker: removedRedisOrder.ticker,
          removalReason: OrderRemovalReason[statefulCancelationOrderRemove.reason],
          createdAtHeight: removedOrder.createdAtHeight,
          updatedAt: removedOrder.updatedAt,
          updatedAtHeight: removedOrder.updatedAtHeight,
          clientMetadata: removedOrder.clientMetadata.toString(),
          triggerPrice,
        }],
        blockHeight: blockHeightRefresher.getLatestBlockHeight(),
      };
      expectWebsocketMessagesSent(
        producerSendSpy,
        SubaccountMessage.fromPartial({
          contents: JSON.stringify(subaccountContents),
          subaccountId: redisTestConstants.defaultSubaccountId,
          version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
        }),
        undefined,
      );

      expect(logger.error).not.toHaveBeenCalled();
      expectTimingStats(true, false, false, true);
    });

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrderIdGoodTilBlockTime,
        testConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        undefined,
      ],
      [
        'conditional',
        redisTestConstants.defaultOrderIdConditional,
        testConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
        testConstants.defaultConditionalOrder.triggerPrice,
      ],
    ])('successfully removes stateful order, resting on book', async (
      _name: string,
      removedOrderId: IndexerOrderId,
      removedOrder: OrderCreateObject,
      removedRedisOrder: RedisOrder,
      expectedOrderUuid: string,
      triggerPrice?: string,
    ) => {
      const orderbookLevel: string = Big(defaultQuantums.toString()).times(2).toString();
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
        ...statefulCancelationOrderRemove,
        removedOrderId,
      });

      await Promise.all([
        placeOrder(
          {
            redisOrder: removedRedisOrder,
            client: redisClient,
          }),
        OrderbookLevelsCache.updatePriceLevel(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          orderbookLevel,
          redisClient,
        ),
      ]);

      // Must be done separately so that the subaccount and perpetualMarket have been created
      // before the order
      await Promise.all([
        OrderTable.create(removedOrder),
        // Must be done after adding orders to all caches to overwrite the ordersDataCache
        setOrderToRestingOnOrderbook(removedRedisOrder),
      ]);

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      // orderbook level reduced by defaultQuantums
      const remainingOrderbookLevel: string = Big(
        orderbookLevel,
      ).minus(
        defaultQuantums.toString(),
      ).toString();
      await Promise.all([
        expectOrderbookLevelCache(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          redisTestConstants.defaultPrice,
          remainingOrderbookLevel,
        ),
        expectOrdersCacheEmpty(expectedOrderUuid),
        expectOrdersDataCacheEmpty(removedOrderId),
        expectSubaccountsOrderIdsCacheEmpty(redisTestConstants.defaultSubaccountUuid),
        expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.NOT_CANCELED),
      ]);

      // Subaccounts message is sent first followed by orderbooks message
      const subaccountContents: SubaccountMessageContents = {
        orders: [{
          id: expectedOrderUuid,
          subaccountId: testConstants.defaultSubaccountId,
          clientId: removedOrderId.clientId.toString(),
          clobPairId: testConstants.defaultOrderGoodTilBlockTime.clobPairId,
          side: OrderSide.BUY,
          size: removedOrder.size,
          totalFilled: '0',
          price: removedOrder.price,
          type: protocolTranslations.protocolConditionTypeToOrderType(
            removedRedisOrder.order!.conditionType,
          ),
          status: OrderStatus.CANCELED,
          timeInForce: apiTranslations.orderTIFToAPITIF(removedOrder.timeInForce),
          postOnly: apiTranslations.isOrderTIFPostOnly(removedOrder.timeInForce),
          reduceOnly: removedOrder.reduceOnly,
          orderFlags: removedOrder.orderFlags,
          goodTilBlockTime: removedOrder.goodTilBlockTime,
          ticker: removedRedisOrder.ticker,
          removalReason: OrderRemovalReason[statefulCancelationOrderRemove.reason],
          createdAtHeight: removedOrder.createdAtHeight,
          updatedAt: removedOrder.updatedAt,
          updatedAtHeight: removedOrder.updatedAtHeight,
          clientMetadata: removedOrder.clientMetadata.toString(),
          triggerPrice,
        }],
        blockHeight: blockHeightRefresher.getLatestBlockHeight(),
      };

      const orderbookContents: OrderbookMessageContents = {
        [OrderbookSide.BIDS]: [[
          defaultPrice,
          protocolTranslations.quantumsToHuman(
            remainingOrderbookLevel,
            testConstants.defaultPerpetualMarket.atomicResolution,
          ).toString(),
        ]],
      };
      expectWebsocketMessagesSent(
        producerSendSpy,
        SubaccountMessage.fromPartial({
          contents: JSON.stringify(subaccountContents),
          subaccountId: redisTestConstants.defaultSubaccountId,
          version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
        }),
        OrderbookMessage.fromPartial({
          contents: JSON.stringify(orderbookContents),
          clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
          version: ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
        }),
      );

      expect(logger.error).not.toHaveBeenCalled();
      expectTimingStats(true, false, true, true);
    });

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrderIdGoodTilBlockTime,
        testConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        undefined,
      ],
      [
        'conditional',
        redisTestConstants.defaultOrderIdConditional,
        testConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
        testConstants.defaultConditionalOrder.triggerPrice,
      ],
    ])('does not increase orderbook level if total filled > quantums', async (
      _name: string,
      removedOrderId: IndexerOrderId,
      removedOrder: OrderCreateObject,
      removedRedisOrder: RedisOrder,
      expectedOrderUuid: string,
      triggerPrice?: string,
    ) => {
      const orderbookLevel: string = Big(defaultQuantums.toString()).times(2).toString();
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
        ...statefulCancelationOrderRemove,
        removedOrderId,
      });

      await Promise.all([
        placeOrder(
          {
            redisOrder: removedRedisOrder,
            client: redisClient,
          }),
        OrderbookLevelsCache.updatePriceLevel(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          orderbookLevel,
          redisClient,
        ),
      ]);

      const exceedsFilledUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage = {
        orderPlace: undefined,
        orderRemove: undefined,
        orderUpdate: {
          orderId: removedRedisOrder.order!.orderId!,
          totalFilledQuantums: defaultQuantums.add(Long.fromValue(100, true)),
        },
      };
      await handleOrderUpdate(exceedsFilledUpdate);

      // Must be done separately so that the subaccount and perpetualMarket have been created
      // before the order
      await Promise.all([
        OrderTable.create(removedOrder),
      ]);

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      await Promise.all([
        expectOrderbookLevelCache(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          redisTestConstants.defaultPrice,
          orderbookLevel,
        ),
        expectOrdersCacheEmpty(expectedOrderUuid),
        expectOrdersDataCacheEmpty(removedOrderId),
        expectSubaccountsOrderIdsCacheEmpty(redisTestConstants.defaultSubaccountUuid),
        expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.NOT_CANCELED),
      ]);

      // Subaccounts message is sent first followed by orderbooks message
      const subaccountContents: SubaccountMessageContents = {
        orders: [{
          id: expectedOrderUuid,
          subaccountId: testConstants.defaultSubaccountId,
          clientId: removedOrderId.clientId.toString(),
          clobPairId: testConstants.defaultOrderGoodTilBlockTime.clobPairId,
          side: OrderSide.BUY,
          size: removedOrder.size,
          totalFilled: '0',
          price: removedOrder.price,
          type: protocolTranslations.protocolConditionTypeToOrderType(
            removedRedisOrder.order!.conditionType,
          ),
          status: OrderStatus.CANCELED,
          timeInForce: apiTranslations.orderTIFToAPITIF(removedOrder.timeInForce),
          postOnly: apiTranslations.isOrderTIFPostOnly(removedOrder.timeInForce),
          reduceOnly: removedOrder.reduceOnly,
          orderFlags: removedOrder.orderFlags,
          goodTilBlockTime: removedOrder.goodTilBlockTime,
          ticker: removedRedisOrder.ticker,
          removalReason: OrderRemovalReason[statefulCancelationOrderRemove.reason],
          createdAtHeight: removedOrder.createdAtHeight,
          updatedAt: removedOrder.updatedAt,
          updatedAtHeight: removedOrder.updatedAtHeight,
          clientMetadata: removedOrder.clientMetadata.toString(),
          triggerPrice,
        }],
        blockHeight: blockHeightRefresher.getLatestBlockHeight(),
      };

      const orderbookContents: OrderbookMessageContents = {
        [OrderbookSide.BIDS]: [[
          defaultPrice,
          protocolTranslations.quantumsToHuman(
            orderbookLevel,
            testConstants.defaultPerpetualMarket.atomicResolution,
          ).toString(),
        ]],
      };
      expectWebsocketMessagesSent(
        producerSendSpy,
        SubaccountMessage.fromPartial({
          contents: JSON.stringify(subaccountContents),
          subaccountId: redisTestConstants.defaultSubaccountId,
          version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
        }),
        OrderbookMessage.fromPartial({
          contents: JSON.stringify(orderbookContents),
          clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
          version: ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
        }),
      );

      expect(logger.error).not.toHaveBeenCalled();
      expectTimingStats(true, false, true, true);
    });
  });

  describe('Order Remove Message - Indexer-expired', () => {
    const indexerExpiredOrderRemoved: OrderRemoveV1 = {
      removedOrderId: redisTestConstants.defaultOrderId,
      reason: OrderRemovalReason.ORDER_REMOVAL_REASON_INDEXER_EXPIRED,
      removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED,
    };

    const indexerExpiredDefaultOrder: OrderCreateObject = {
      ...testConstants.defaultOrder,
      goodTilBlock: redisTestConstants.defaultOrder.goodTilBlock!.toString(),
    };

    it('successfully removes expired order', async () => {
      const removedOrderId: IndexerOrderId = redisTestConstants.defaultOrderId;
      const removedOrder: OrderCreateObject = indexerExpiredDefaultOrder;
      const removedRedisOrder: RedisOrder = redisTestConstants.defaultRedisOrder;
      const expectedOrderUuid: string = redisTestConstants.defaultOrderUuid;

      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
        ...indexerExpiredOrderRemoved,
        removedOrderId,
      });
      const orderbookLevel: string = Big(defaultQuantums.toString()).times(2).toString();

      await Promise.all([
        // testConstants.defaultOrder has a goodTilBlock of 1150
        BlockTable.create({ blockHeight: '1151', time: DateTime.utc(2022, 6, 1).toISO() }),
        placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
        OrderbookLevelsCache.updatePriceLevel(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          orderbookLevel,
          redisClient,
        ),
      ]);

      await Promise.all([
        // Must be done separately so that the subaccount and perpetualMarket have been created
        // before the order
        OrderTable.create(removedOrder),
        // Must be done after adding orders to all caches to overwrite the ordersDataCache
        setOrderToRestingOnOrderbook(removedRedisOrder),
      ]);

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      // orderbook level reduced by defaultQuantums
      const remainingOrderbookLevel: string = Big(
        orderbookLevel,
      ).minus(
        defaultQuantums.toString(),
      ).toString();
      await Promise.all([
        expectOrderStatus(expectedOrderUuid, OrderStatus.CANCELED),
        expectOrderbookLevelCache(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          remainingOrderbookLevel,
        ),
        expectOrdersCacheEmpty(expectedOrderUuid),
        expectOrdersDataCacheEmpty(removedOrderId),
        expectSubaccountsOrderIdsCacheEmpty(redisTestConstants.defaultSubaccountUuid),
        expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.CANCELED),
      ]);

      // Subaccounts message is sent first followed by orderbooks message
      const subaccountContents: SubaccountMessageContents = {
        orders: [
          {
            id: expectedOrderUuid,
            subaccountId: testConstants.defaultSubaccountId,
            clientId: removedOrderId.clientId.toString(),
            clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
            side: OrderSide.BUY,
            size: defaultSize,
            totalOptimisticFilled: '0',
            price: defaultPrice,
            type: OrderType.LIMIT,
            status: OrderStatus.CANCELED,
            timeInForce: apiTranslations.orderTIFToAPITIF(
              protocolTranslations.protocolOrderTIFToTIF(removedRedisOrder.order!.timeInForce),
            ),
            postOnly: apiTranslations.isOrderTIFPostOnly(
              protocolTranslations.protocolOrderTIFToTIF(removedRedisOrder.order!.timeInForce),
            ),
            reduceOnly: removedRedisOrder.order!.reduceOnly,
            orderFlags: removedRedisOrder.order!.orderId!.orderFlags.toString(),
            goodTilBlock: protocolTranslations.getGoodTilBlock(removedRedisOrder.order!)
              ?.toString(),
            goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(removedRedisOrder.order!),
            ticker: redisTestConstants.defaultRedisOrder.ticker,
            removalReason: OrderRemovalReason[indexerExpiredOrderRemoved.reason],
            createdAtHeight: removedOrder.createdAtHeight,
            updatedAt: removedOrder.updatedAt,
            updatedAtHeight: removedOrder.updatedAtHeight,
            clientMetadata: testConstants.defaultOrderGoodTilBlockTime.clientMetadata.toString(),
          },
        ],
        blockHeight: blockHeightRefresher.getLatestBlockHeight(),
      };
      const orderbookContents: OrderbookMessageContents = {
        [OrderbookSide.BIDS]: [[
          defaultPrice,
          protocolTranslations.quantumsToHuman(
            remainingOrderbookLevel,
            testConstants.defaultPerpetualMarket.atomicResolution,
          ).toString(),
        ]],
      };
      expectWebsocketMessagesSent(
        producerSendSpy,
        SubaccountMessage.fromPartial({
          contents: JSON.stringify(subaccountContents),
          subaccountId: redisTestConstants.defaultSubaccountId,
          version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
        }),
        OrderbookMessage.fromPartial({
          contents: JSON.stringify(orderbookContents),
          clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
          version: ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
        }),
      );
      expectTimingStats(true, true);
    });

    it('successfully removes fully filled expired order and does not send websocket message', async () => {
      const removedOrderId: IndexerOrderId = redisTestConstants.defaultOrderId;
      const removedOrder: OrderCreateObject = {
        ...indexerExpiredDefaultOrder,
        status: OrderStatus.FILLED,
      };
      const removedRedisOrder: RedisOrder = redisTestConstants.defaultRedisOrder;
      const expectedOrderUuid: string = redisTestConstants.defaultOrderUuid;

      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate({
        ...indexerExpiredOrderRemoved,
        removedOrderId,
      });
      const orderbookLevel: string = Big(defaultQuantums.toString()).times(2).toString();

      await Promise.all([
        // testConstants.defaultOrder has a goodTilBlock of 1150
        BlockTable.create({ blockHeight: '1151', time: DateTime.utc(2022, 6, 1).toISO() }),
        placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
        OrderbookLevelsCache.updatePriceLevel(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          orderbookLevel,
          redisClient,
        ),
        StateFilledQuantumsCache.updateStateFilledQuantums(
          expectedOrderUuid,
          removedRedisOrder.order!.quantums.toString(),
          redisClient,
        ),
      ]);

      await Promise.all([
        // Must be done separately so that the subaccount and perpetualMarket have been created
        // before the order
        OrderTable.create(removedOrder),
        // Must be done after adding orders to all caches to overwrite the ordersDataCache
        setOrderToRestingOnOrderbook(removedRedisOrder),
        updateOrder({
          updatedOrderId: removedRedisOrder.order!.orderId!,
          newTotalFilledQuantums: removedRedisOrder.order!.quantums.toNumber(),
          client: redisClient,
        }),
      ]);

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      // orderbook level should not be reduced
      const remainingOrderbookLevel: string = Big(
        orderbookLevel,
      ).toString();
      await Promise.all([
        expectOrderStatus(expectedOrderUuid, removedOrder.status),
        expectOrderbookLevelCache(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          remainingOrderbookLevel,
        ),
        expectOrdersCacheEmpty(expectedOrderUuid),
        expectOrdersDataCacheEmpty(removedOrderId),
        expectSubaccountsOrderIdsCacheEmpty(redisTestConstants.defaultSubaccountUuid),
        expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.CANCELED),
      ]);
      expectNoWebsocketMessagesSent(producerSendSpy);
      expectTimingStats(true, false);
    });

    it('error: when latest block not found, log and exit', async () => {
      const removedOrderId: IndexerOrderId = redisTestConstants.defaultOrderId;
      const removedOrder: OrderCreateObject = indexerExpiredDefaultOrder;
      const removedRedisOrder: RedisOrder = redisTestConstants.defaultRedisOrder;
      const expectedOrderUuid: string = redisTestConstants.defaultOrderUuid;

      // eslint-disable-next-line @typescript-eslint/require-await
      const tableSpy = jest.spyOn(BlockTable, 'getLatest').mockImplementation(async () => {
        throw new Error();
      });

      try {
        const orderRemoveJson: OrderRemoveV1 = { ...indexerExpiredOrderRemoved, removedOrderId };
        const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate(orderRemoveJson);
        const orderbookLevel: string = Big(defaultQuantums.toString()).times(2).toString();

        await Promise.all([
          placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
          OrderbookLevelsCache.updatePriceLevel(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            defaultPrice,
            orderbookLevel,
            redisClient,
          ),
        ]);

        await Promise.all([
          // Must be done separately so that the subaccount and perpetualMarket have been created
          // before the order
          OrderTable.create(removedOrder),
          // Must be done after adding orders to all caches to overwrite the ordersDataCache
          setOrderToRestingOnOrderbook(removedRedisOrder),
        ]);

        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

        const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
        await orderRemoveHandler.handleUpdate(
          offChainUpdate,
          defaultKafkaHeaders,
        );

        expect(producerSendSpy).not.toHaveBeenCalled();
        expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
          at: 'orderRemoveHandler#isOrderExpired',
          message: expect.stringContaining('Unable to find latest block'),
          orderRemove: orderRemoveJson,
        }));

        await Promise.all([
          expectOrderStatus(expectedOrderUuid, OrderStatus.OPEN),
          expectOrderbookLevelCache(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            defaultPrice,
            orderbookLevel,
          ),
          expectOrdersCacheFound(expectedOrderUuid),
          expectOrdersDataCacheFound(removedOrderId),
          expectSubaccountsOrderIdsCacheFound(redisTestConstants.defaultSubaccountUuid),
          expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.NOT_CANCELED),
        ]);

        expectTimingStats(false, false, false, false, true);
      } finally {
        tableSpy.mockRestore();
      }
    });

    it('error: when order not found, log and exit', async () => {
      const removedOrderId: IndexerOrderId = redisTestConstants.defaultOrderId;
      const removedOrder: OrderCreateObject = indexerExpiredDefaultOrder;
      const removedRedisOrder: RedisOrder = redisTestConstants.defaultRedisOrder;
      const expectedOrderUuid: string = redisTestConstants.defaultOrderUuid;

      const tableSpy = jest.spyOn(OrdersCache, 'getOrder').mockResolvedValueOnce(null);

      try {
        const orderRemoveJson: OrderRemoveV1 = { ...indexerExpiredOrderRemoved, removedOrderId };
        const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate(orderRemoveJson);
        const orderbookLevel: string = Big(defaultQuantums.toString()).times(2).toString();

        await Promise.all([
          // testConstants.defaultOrder has a goodTilBlock of 1150
          BlockTable.create({ blockHeight: '1151', time: DateTime.utc(2022, 6, 1).toISO() }),
          placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
          OrderbookLevelsCache.updatePriceLevel(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            defaultPrice,
            orderbookLevel,
            redisClient,
          ),
        ]);

        await Promise.all([
          // Must be done separately so that the subaccount and perpetualMarket have been created
          // before the order
          OrderTable.create(removedOrder),
          // Must be done after adding orders to all caches to overwrite the ordersDataCache
          setOrderToRestingOnOrderbook(removedRedisOrder),
        ]);

        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

        const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
        await orderRemoveHandler.handleUpdate(
          offChainUpdate,
          defaultKafkaHeaders,
        );

        expect(producerSendSpy).not.toHaveBeenCalled();
        expect(stats.increment).toHaveBeenCalledWith('vulcan.indexer_expired_order_not_found', 1, { instance: '' });

        await Promise.all([
          expectOrderStatus(expectedOrderUuid, OrderStatus.OPEN),
          expectOrderbookLevelCache(
            removedRedisOrder.ticker,
            OrderSide.BUY,
            defaultPrice,
            orderbookLevel,
          ),
          expectOrdersCacheFound(expectedOrderUuid),
          expectOrdersDataCacheFound(removedOrderId),
          expectSubaccountsOrderIdsCacheFound(redisTestConstants.defaultSubaccountUuid),
          expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.NOT_CANCELED),
        ]);

        expectTimingStats(false, false, false, false, true, true);
      } finally {
        tableSpy.mockRestore();
      }
    });

    it('error: when order found is not short-term, log and exit', async () => {
      const removedOrderId: IndexerOrderId = {
        ...redisTestConstants.defaultOrderId,
        orderFlags: ORDER_FLAG_LONG_TERM,
      };
      const removedOrder: OrderCreateObject = {
        ...testConstants.defaultOrder,
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
      };
      const expectedOrderUuid: string = OrderTable.orderIdToUuid(removedOrderId);
      const removedRedisOrder: RedisOrder = {
        ...redisTestConstants.defaultRedisOrder,
        order: {
          ...redisTestConstants.defaultOrder,
          orderId: removedOrderId,
        },
        id: expectedOrderUuid,
      };

      const orderRemoveJson: OrderRemoveV1 = { ...indexerExpiredOrderRemoved, removedOrderId };
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate(orderRemoveJson);
      const orderbookLevel: string = Big(defaultQuantums.toString()).times(2).toString();

      await Promise.all([
        // testConstants.defaultOrder has a goodTilBlock of 1150
        BlockTable.create({ blockHeight: '1151', time: DateTime.utc(2022, 6, 1).toISO() }),
        placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
        OrderbookLevelsCache.updatePriceLevel(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          orderbookLevel,
          redisClient,
        ),
      ]);

      await Promise.all([
        // Must be done separately so that the subaccount and perpetualMarket have been created
        // before the order
        OrderTable.create(removedOrder),
        // Must be done after adding orders to all caches to overwrite the ordersDataCache
        setOrderToRestingOnOrderbook(removedRedisOrder),
      ]);

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      expect(producerSendSpy).not.toHaveBeenCalled();
      expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
        at: 'orderRemoveHandler#isOrderExpired',
        message: expect.stringContaining(
          'Long-term order retrieved during Indexer-expired expiry verification',
        ),
        orderRemove: orderRemoveJson,
        redisOrder: removedRedisOrder,
      }));

      await Promise.all([
        expectOrderStatus(expectedOrderUuid, OrderStatus.OPEN),
        expectOrderbookLevelCache(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          orderbookLevel,
        ),
        expectOrdersCacheFound(expectedOrderUuid),
        expectOrdersDataCacheFound(removedOrderId),
        expectSubaccountsOrderIdsCacheFound(redisTestConstants.defaultSubaccountUuid),
        expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.NOT_CANCELED),
      ]);

      expectTimingStats(false, false, false, false, true, true);
    });

    it('error: when order is not expired, log and exit', async () => {
      const removedOrderId: IndexerOrderId = redisTestConstants.defaultOrderId;
      const removedOrder: OrderCreateObject = indexerExpiredDefaultOrder;
      const removedRedisOrder: RedisOrder = redisTestConstants.defaultRedisOrder;
      const expectedOrderUuid: string = redisTestConstants.defaultOrderUuid;

      const orderRemoveJson: OrderRemoveV1 = { ...indexerExpiredOrderRemoved, removedOrderId };
      const offChainUpdate: OffChainUpdateV1 = orderRemoveToOffChainUpdate(orderRemoveJson);
      const orderbookLevel: string = Big(defaultQuantums.toString()).times(2).toString();

      await Promise.all([
        placeOrder({ redisOrder: removedRedisOrder, client: redisClient }),
        OrderbookLevelsCache.updatePriceLevel(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          orderbookLevel,
          redisClient,
        ),
      ]);

      await Promise.all([
        // Must be done separately so that the subaccount and perpetualMarket have been created
        // before the order
        OrderTable.create(removedOrder),
        // Must be done after adding orders to all caches to overwrite the ordersDataCache
        setOrderToRestingOnOrderbook(removedRedisOrder),
      ]);

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      const orderRemoveHandler: OrderRemoveHandler = new OrderRemoveHandler();
      await orderRemoveHandler.handleUpdate(
        offChainUpdate,
        defaultKafkaHeaders,
      );

      expect(producerSendSpy).not.toHaveBeenCalled();
      expect(
        stats.increment,
      ).toHaveBeenCalledWith('vulcan.indexer_expired_order_is_not_expired', 1, { instance: '' });

      await Promise.all([
        expectOrderStatus(expectedOrderUuid, OrderStatus.OPEN),
        expectOrderbookLevelCache(
          removedRedisOrder.ticker,
          OrderSide.BUY,
          defaultPrice,
          orderbookLevel,
        ),
        expectOrdersCacheFound(expectedOrderUuid),
        expectOrdersDataCacheFound(removedOrderId),
        expectSubaccountsOrderIdsCacheFound(redisTestConstants.defaultSubaccountUuid),
        expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.NOT_CANCELED),
      ]);

      expectTimingStats(false, false, false, false, true, true);
    });
  });

  async function expectOrderStatus(orderId: string, status: OrderStatus): Promise<void> {
    const order: OrderFromDatabase | undefined = await OrderTable.findById(orderId);
    expect(order).toBeDefined();
    expect(order!.status).toEqual(status);
  }

  async function expectOrdersCacheFound(
    orderId: string,
  ): Promise<void> {
    const order: RedisOrder | null = await OrdersCache.getOrder(orderId, redisClient);
    expect(order).not.toBeNull();
  }

  async function expectOrdersDataCacheFound(
    orderId: IndexerOrderId,
  ): Promise<void> {
    const orderData: OrderData | null = await OrdersDataCache.getOrderData(orderId, redisClient);
    expect(orderData).not.toBeNull();
  }

  async function expectSubaccountsOrderIdsCacheFound(
    subaccountUuid: string,
  ): Promise<void> {
    const orderIds: string[] = await SubaccountOrderIdsCache.getOrderIdsForSubaccount(
      subaccountUuid,
      redisClient,
    );
    expect(orderIds.length).not.toEqual(0);
  }

  async function expectOrdersCacheEmpty(
    orderId: string,
  ): Promise<void> {
    const order: RedisOrder | null = await OrdersCache.getOrder(orderId, redisClient);
    expect(order).toBeNull();
  }

  async function expectOrdersDataCacheEmpty(
    orderId: IndexerOrderId,
  ): Promise<void> {
    const orderData: OrderData | null = await OrdersDataCache.getOrderData(orderId, redisClient);
    expect(orderData).toBeNull();
  }

  async function expectSubaccountsOrderIdsCacheEmpty(
    subaccountUuid: string,
  ): Promise<void> {
    const orderIds: string[] = await SubaccountOrderIdsCache.getOrderIdsForSubaccount(
      subaccountUuid,
      redisClient,
    );
    expect(orderIds.length).toEqual(0);
  }

  function expectNoWebsocketMessagesSent(
    producerSendSpy: jest.SpyInstance,
  ): void {
    jest.runOnlyPendingTimers();
    expect(producerSendSpy).not.toHaveBeenCalled();
  }

  function expectWebsocketMessagesSent(
    producerSendSpy: jest.SpyInstance,
    expectedSubaccountMessage?: SubaccountMessage,
    expectedOrderbookMessage?: OrderbookMessage,
  ): void {
    jest.runOnlyPendingTimers();
    let numMessages: number = 0;
    if (expectedSubaccountMessage !== undefined) {
      numMessages += 1;
    }
    if (expectedOrderbookMessage !== undefined) {
      numMessages += 1;
    }
    // expect one call for subaccount and one for orderbook if expectedOrderbookMessage is defined
    expect(producerSendSpy).toHaveBeenCalledTimes(numMessages);

    if (expectedSubaccountMessage !== undefined) {
      const subaccountProducerRecord: ProducerRecord = producerSendSpy.mock.calls[0][0];
      expectWebsocketSubaccountMessage(
        subaccountProducerRecord,
        expectedSubaccountMessage,
        defaultKafkaHeaders,
      );
    }

    if (expectedOrderbookMessage !== undefined) {
      // If a subaccount message was not sent, the orderbook message should be the first call
      const callIndex: number = expectedSubaccountMessage !== undefined ? 1 : 0;
      const orderbookProducerRecord: ProducerRecord = producerSendSpy.mock.calls[callIndex][0];
      expectWebsocketOrderbookMessage(orderbookProducerRecord, expectedOrderbookMessage);
    }
  }
});

function orderRemoveToOffChainUpdate(
  orderRemoveJson: any,
): OffChainUpdateV1 {
  return {
    orderRemove: orderRemoveJson,
  };
}

async function setOrderToRestingOnOrderbook(
  redisOrder: RedisOrder,
): Promise<void> {
  await redis.setAsync({
    key: OrdersDataCache.getOrderDataCacheKey(redisOrder.order!.orderId!),
    // [good-til-block or good-til-blocktime]_[totalFilled]_[resting on book]
    value: '5_0_true',
  }, redisClient);
}

function expectTimingStats(
  shouldRemoveOrder: boolean = true,
  shouldCancelOrderInPostgres: boolean = false,
  shouldUpdatePriceLevels: boolean = false,
  shouldFindOrderInPostgrers: boolean = false,
  shouldGetLatestBlockForIndexerExpiredExpiryVerification: boolean = false,
  shouldFindOrderForIndexerExpiredExpiryVerification: boolean = false,
) {
  if (shouldRemoveOrder) {
    expectTimingStat('remove_order');
  }
  if (shouldCancelOrderInPostgres) {
    expectTimingStat('cancel_order_in_postgres');
  }
  if (shouldFindOrderInPostgrers) {
    expectTimingStat('find_order_for_stateful_cancelation');
  }
  if (shouldUpdatePriceLevels) {
    expectTimingStat('update_price_level_cache');
  }
  if (shouldGetLatestBlockForIndexerExpiredExpiryVerification) {
    expectTimingStat('get_latest_block_for_indexer_expired_expiry_verification');
  }
  if (shouldFindOrderForIndexerExpiredExpiryVerification) {
    expectTimingStat('find_order_for_indexer_expired_expiry_verification');
  }
}

function expectTimingStat(fnName: string) {
  expect(stats.timing).toHaveBeenCalledWith(
    `vulcan.${STATS_FUNCTION_NAME}.timing`,
    expect.any(Number),
    { className: 'OrderRemoveHandler', fnName, instance: '' },
  );
}
