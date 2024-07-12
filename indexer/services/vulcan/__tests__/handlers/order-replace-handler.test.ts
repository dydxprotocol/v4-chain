import {
  logger,
  stats,
  STATS_FUNCTION_NAME,
  wrapBackgroundTask,
} from '@dydxprotocol-indexer/base';
import { synchronizeWrapBackgroundTask } from '@dydxprotocol-indexer/dev';
import {
  createKafkaMessage,
  producer,
  SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
  getTriggerPrice,
  ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import {
  APIOrderStatus,
  APIOrderStatusEnum,
  apiTranslations,
  blockHeightRefresher,
  BlockTable,
  dbHelpers,
  OrderbookMessageContents,
  OrderFromDatabase,
  OrderTable,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  protocolTranslations,
  SubaccountMessageContents,
  SubaccountTable,
  testConstants,
  testMocks,
  TimeInForce,
} from '@dydxprotocol-indexer/postgres';
import * as redisPackage from '@dydxprotocol-indexer/redis';
import {
  PriceLevel,
  OrdersCache,
  OrderbookLevels,
  OrderbookLevelsCache,
  redis,
  redisTestConstants,
  SubaccountOrderIdsCache,
  CanceledOrdersCache,
  updateOrder,
  CanceledOrderStatus,
} from '@dydxprotocol-indexer/redis';
import {
  OffChainUpdateV1,
  IndexerOrder,
  OrderbookMessage,
  OrderPlaceV1_OrderPlacementStatus,
  RedisOrder,
  SubaccountId,
  SubaccountMessage,
} from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import { redisClient, redisClient as client } from '../../src/helpers/redis/redis-controller';
import { onMessage } from '../../src/lib/on-message';
import { expectCanceledOrderStatus, handleInitialOrderPlace } from '../helpers/helpers';
import { expectOffchainUpdateMessage, expectWebsocketOrderbookMessage, expectWebsocketSubaccountMessage } from '../helpers/websocket-helpers';
import { isStatefulOrder } from '@dydxprotocol-indexer/v4-proto-parser';
import { defaultKafkaHeaders } from '../helpers/constants';
import config from '../../src/config';
import { OrderbookSide } from '../../src/lib/types';

jest.mock('@dydxprotocol-indexer/base', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/base'),
  wrapBackgroundTask: jest.fn(),
}));

interface OffchainUpdateRecord {
  key: Buffer,
  offchainUpdate: OffChainUpdateV1
}

describe('order-replace-handler', () => {
  beforeAll(async () => {
    await BlockTable.create(testConstants.defaultBlock);
    await blockHeightRefresher.updateBlockHeight();
    jest.useFakeTimers();
  });

  afterAll(() => {
    jest.useRealTimers();
  });

  afterEach(() => {
    config.SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_STATEFUL_ORDERS = true;
  });

  describe('handle', () => {
    const replacementOrder: IndexerOrder = redisTestConstants.defaultReplacementOrder;
    // eslint-disable-next-line max-len
    const replacementOrderGoodTilBlockTime: IndexerOrder = redisTestConstants.defaultReplacementOrderGTBT;
    const replacementOrderDifferentPrice: IndexerOrder = {
      ...replacementOrderGoodTilBlockTime,
      subticks: replacementOrderGoodTilBlockTime.subticks.mul(2),
    };

    const replacedOrder: RedisOrder = redisPackage.convertToRedisOrder(
      replacementOrder,
      testConstants.defaultPerpetualMarket,
    );
    const replacedOrderGoodTilBlockTime: RedisOrder = redisPackage.convertToRedisOrder(
      replacementOrderGoodTilBlockTime,
      testConstants.defaultPerpetualMarket,
    );
    const replacedOrderDifferentPrice: RedisOrder = redisPackage.convertToRedisOrder(
      replacementOrderDifferentPrice,
      testConstants.defaultPerpetualMarket,
    );

    const replacementUpdate: OffChainUpdateV1 = {
      orderReplace: {
        oldOrderId: redisTestConstants.defaultOrderId,
        order: replacementOrder,
        placementStatus:
            OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
      },
    };
    const replacementUpdateGoodTilBlockTime: OffChainUpdateV1 = {
      orderReplace: {
        oldOrderId: redisTestConstants.defaultOrderIdGoodTilBlockTime,
        order: replacementOrderGoodTilBlockTime,
        placementStatus:
            OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
      },
    };
    const replacementUpdateDifferentPrice: OffChainUpdateV1 = {
      orderReplace: {
        oldOrderId: redisTestConstants.defaultOrderIdGoodTilBlockTime,
        order: replacementOrderDifferentPrice,
        placementStatus:
            OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
      },
    };

    const replacementMessage: KafkaMessage = createKafkaMessage(
      Buffer.from(Uint8Array.from(OffChainUpdateV1.encode(replacementUpdate).finish())),
    );
    const replacementMessageGoodTilBlockTime: KafkaMessage = createKafkaMessage(
      Buffer.from(Uint8Array.from(
        OffChainUpdateV1.encode(replacementUpdateGoodTilBlockTime).finish(),
      )),
    );
    const replacementMessageDifferentPrice: KafkaMessage = createKafkaMessage(
      Buffer.from(Uint8Array.from(
        OffChainUpdateV1.encode(replacementUpdateDifferentPrice).finish(),
      )),
    );

    [replacementMessage, replacementMessageGoodTilBlockTime,
      replacementMessageDifferentPrice].forEach((message) => {
      // eslint-disable-next-line no-param-reassign
      message.headers = defaultKafkaHeaders;
    });

    const dbDefaultOrder: OrderFromDatabase = {
      ...testConstants.defaultOrder,
      id: testConstants.defaultOrderId,
    };
    const dbOrderGoodTilBlockTime: OrderFromDatabase = {
      ...testConstants.defaultOrderGoodTilBlockTime,
      id: testConstants.defaultOrderGoodTilBlockTimeId,
      createdAtHeight: '2',
    };

    beforeAll(async () => {
      await dbHelpers.migrate();
    });

    beforeEach(async () => {
      await dbHelpers.clearData();
      await testMocks.seedData();
      await Promise.all([
        perpetualMarketRefresher.updatePerpetualMarkets(),
        blockHeightRefresher.updateBlockHeight(),
      ]);
      await Promise.all([
        OrderTable.create(dbDefaultOrder),
        OrderTable.create(dbOrderGoodTilBlockTime),
      ]);
      jest.spyOn(stats, 'timing');
      jest.spyOn(OrderbookLevelsCache, 'updatePriceLevel');
      jest.spyOn(CanceledOrdersCache, 'removeOrderFromCaches');
      jest.spyOn(stats, 'increment');
      jest.spyOn(redisPackage, 'placeOrder');
      jest.spyOn(logger, 'error');
      jest.spyOn(logger, 'info');
    });

    afterEach(async () => {
      await redis.deleteAllAsync(client);
      await dbHelpers.clearData();
      jest.restoreAllMocks();
    });

    afterAll(async () => {
      await dbHelpers.teardown();
    });
    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrder,
        replacementMessage,
        redisTestConstants.defaultRedisOrder,
        dbDefaultOrder,
        redisTestConstants.defaultOrderUuid,
        redisTestConstants.defaultReplacementOrderUuid,
        replacedOrder,
        true,
        false,
      ],
      [
        'goodTilBlockTime',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        replacementMessageGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        dbOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        redisTestConstants.defaultReplacementOrderUuidGTBT,
        replacedOrderGoodTilBlockTime,
        false,
        false,
      ],
      [
        'goodTilBlock and canceled order',
        redisTestConstants.defaultOrder,
        replacementMessage,
        redisTestConstants.defaultRedisOrder,
        dbDefaultOrder,
        redisTestConstants.defaultOrderUuid,
        redisTestConstants.defaultReplacementOrderUuid,
        replacedOrder,
        true,
        true,
      ],
      [
        'goodTilBlockTime and canceled order',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        replacementMessageGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        dbOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        redisTestConstants.defaultReplacementOrderUuidGTBT,
        replacedOrderGoodTilBlockTime,
        false,
        true,
      ],
    ])('handles replacement (with %s), not resting on book', async (
      _name: string,
      initialOrderToPlace: IndexerOrder,
      orderReplacementMessage: KafkaMessage,
      expectedRedisOrder: RedisOrder,
      dbOrder: OrderFromDatabase,
      expectedOldOrderUuid: string,
      expectedNewOrderUuid: string,
      expectedReplacedOrder: RedisOrder,
      expectSubaccountMessage: boolean,
      hasCanceledOrderId: boolean,
    ) => {
      if (hasCanceledOrderId) {
        await redisPackage.CanceledOrdersCache.addCanceledOrderId(
          expectedOldOrderUuid,
          Date.now(),
          redisClient,
        );
      }
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      // Handle the order place event for the initial order that will get replaced
      await handleInitialOrderPlace({
        ...redisTestConstants.orderPlace,
        orderPlace: {
          order: initialOrderToPlace,
          placementStatus:
              OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
        },
      });
      expectWebsocketMessagesSent(
        producerSendSpy,
        expectedRedisOrder,
        dbOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.BEST_EFFORT_OPENED,
        expectSubaccountMessage,
      );
      expectStats();
      // clear mocks
      jest.clearAllMocks();

      // Handle the order replacement off-chain update with the replacement order
      await onMessage(orderReplacementMessage);

      await checkOrderReplace(
        expectedOldOrderUuid,
        expectedNewOrderUuid,
        redisTestConstants.defaultSubaccountUuid,
        expectedReplacedOrder,
      );
      expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();
      if (hasCanceledOrderId) {
        expect(CanceledOrdersCache.removeOrderFromCaches).toHaveBeenCalled();
      }
      await expectCanceledOrderStatus(expectedOldOrderUuid, CanceledOrderStatus.NOT_CANCELED);

      expect(logger.error).not.toHaveBeenCalled();
      expectWebsocketMessagesSent(
        producerSendSpy,
        expectedReplacedOrder,
        dbOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.BEST_EFFORT_OPENED,
        expectSubaccountMessage,
      );
      expectStats(true);
    });

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrder,
        replacementMessage,
        redisTestConstants.defaultRedisOrder,
        dbDefaultOrder,
        redisTestConstants.defaultOrderUuid,
        redisTestConstants.defaultReplacementOrderUuid,
        replacedOrder,
        true,
        true,
        false,
      ],
      [
        'goodTilBlockTime',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        replacementMessageGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        dbOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        redisTestConstants.defaultReplacementOrderUuidGTBT,
        replacedOrderGoodTilBlockTime,
        false,
        true,
        false,
      ],
      [
        'goodTilBlockTime different price',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        replacementMessageDifferentPrice,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        dbOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        redisTestConstants.defaultReplacementOrderUuidGTBT,
        replacedOrderDifferentPrice,
        false,
        true,
        true,
      ],
    ])('handles order replace (with %s), resting on book', async (
      _name: string,
      initialOrderToPlace: IndexerOrder,
      orderReplacementMessage: KafkaMessage,
      expectedRedisOrder: RedisOrder,
      dbOrder: OrderFromDatabase,
      expectedOldOrderUuid: string,
      expectedNewOrderUuid: string,
      expectedReplacedOrder: RedisOrder,
      expectSubaccountMessage: boolean,
      expectOrderBookUpdate: boolean,
      expectOrderBookMessage: boolean,
    ) => {
      const oldOrderTotalFilled: number = 10;
      const oldPriceLevelInitialQuantums: number = Number(initialOrderToPlace.quantums) * 2;
      // After replacing the order the quantums at the price level of the old order should be:
      // initial quantums - (old order quantums - old order total filled)
      const expectedPriceLevelQuantums: number = (
        oldPriceLevelInitialQuantums - (Number(initialOrderToPlace.quantums) - oldOrderTotalFilled)
      );
      const expectedPriceLevel: PriceLevel = {
        humanPrice: expectedRedisOrder.price,
        quantums: expectedPriceLevelQuantums.toString(),
        lastUpdated: expect.stringMatching(/^[0-9]{10}$/),
      };
      const expectedPriceLevelSize: string = protocolTranslations.quantumsToHumanFixedString(
        expectedPriceLevelQuantums.toString(),
        testConstants.defaultPerpetualMarket.atomicResolution,
      );

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      // Handle the order place event for the initial order that will get replaced
      await handleInitialOrderPlace({
        ...redisTestConstants.orderPlace,
        orderPlace: {
          order: initialOrderToPlace,
          placementStatus:
              OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
        },
      });
      expectWebsocketMessagesSent(
        producerSendSpy,
        expectedRedisOrder,
        dbOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.BEST_EFFORT_OPENED,
        expectSubaccountMessage,
      );
      // clear mocks
      jest.clearAllMocks();

      // Update the order to set it to be resting on the book
      await updateOrder({
        updatedOrderId: initialOrderToPlace.orderId!,
        newTotalFilledQuantums: oldOrderTotalFilled,
        client,
      });

      // Update the price level in the order book to a value larger than the quantums of the order
      await OrderbookLevelsCache.updatePriceLevel({
        ticker: testConstants.defaultPerpetualMarket.ticker,
        side: protocolTranslations.protocolOrderSideToOrderSide(
          initialOrderToPlace.side,
        ),
        humanPrice: expectedRedisOrder.price,
        sizeDeltaInQuantums: oldPriceLevelInitialQuantums.toString(),
        client,
      });

      // Handle the order replacement off-chain update with the replacement order
      await onMessage(orderReplacementMessage);

      await checkOrderReplace(
        expectedOldOrderUuid,
        expectedNewOrderUuid,
        redisTestConstants.defaultSubaccountUuid,
        expectedReplacedOrder,
      );
      const orderbook: OrderbookLevels = await OrderbookLevelsCache.getOrderBookLevels(
        testConstants.defaultPerpetualMarket.ticker,
        client,
      );

      // Check the order book levels were updated
      if (expectOrderBookUpdate) {
        expect(orderbook.bids).toHaveLength(1);
        expect(orderbook.asks).toHaveLength(0);
        expect(orderbook.bids).toContainEqual(expectedPriceLevel);
      }

      expect(logger.error).not.toHaveBeenCalled();
      const orderbookContents: OrderbookMessageContents = {
        [OrderbookSide.BIDS]: [[
          redisTestConstants.defaultPrice,
          expectedPriceLevelSize,
        ]],
      };
      expectWebsocketMessagesSent(
        producerSendSpy,
        expectedReplacedOrder,
        dbOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.BEST_EFFORT_OPENED,
        expectSubaccountMessage,
        expectOrderBookMessage
          ? OrderbookMessage.fromPartial({
            contents: JSON.stringify(orderbookContents),
            clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
            version: ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
          }) : undefined,
      );
      expectStats(true);
    });

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrder,
        replacementMessage,
        redisTestConstants.defaultRedisOrder,
        dbDefaultOrder,
        redisTestConstants.defaultOrderUuid,
        redisTestConstants.defaultReplacementOrderUuid,
        replacedOrder,
        true,
      ],
      [
        'goodTilBlockTime',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        replacementMessageGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        dbOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        redisTestConstants.defaultReplacementOrderUuidGTBT,
        replacedOrderGoodTilBlockTime,
        false,
      ],
    ])('handles order replacement (with %s), resting on book, 0 remaining quantums',
      async (
        _name: string,
        initialOrderToPlace: IndexerOrder,
        orderReplacementMessage: KafkaMessage,
        expectedRedisOrder: RedisOrder,
        dbOrder: OrderFromDatabase,
        expectedOldOrderUuid: string,
        expectedNewOrderUuid: string,
        expectedReplacedOrder: RedisOrder,
        expectSubaccountMessage: boolean,
      ) => {
        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
        // Handle the order place event for the initial order that will get replaced
        await handleInitialOrderPlace({
          ...redisTestConstants.orderPlace,
          orderPlace: {
            order: initialOrderToPlace,
            placementStatus:
                OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
          },
        });
        expectWebsocketMessagesSent(
          producerSendSpy,
          expectedRedisOrder,
          dbOrder,
          testConstants.defaultPerpetualMarket,
          APIOrderStatusEnum.BEST_EFFORT_OPENED,
          expectSubaccountMessage,
        );
        expectStats();
        // clear mocks
        jest.clearAllMocks();

        // Update the order to set it to be resting on the book
        await updateOrder({
          updatedOrderId: initialOrderToPlace.orderId!,
          newTotalFilledQuantums: Number(initialOrderToPlace.quantums),
          client,
        });
        // Handle the order replacement off-chain update with the replacement order
        await onMessage(orderReplacementMessage);

        await checkOrderReplace(
          expectedOldOrderUuid,
          expectedNewOrderUuid,
          redisTestConstants.defaultSubaccountUuid,
          expectedReplacedOrder,
        );
        // expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();

        expect(logger.error).not.toHaveBeenCalled();
        expectWebsocketMessagesSent(
          producerSendSpy,
          expectedReplacedOrder,
          dbOrder,
          testConstants.defaultPerpetualMarket,
          APIOrderStatusEnum.BEST_EFFORT_OPENED,
          expectSubaccountMessage,
        );
        expectStats(true);
      },
    );

    it.each([
      [
        'missing order',
        {
          ...redisTestConstants.orderReplace,
          orderReplace: { ...redisTestConstants.orderReplace.orderReplace, order: undefined },
        },
        'Invalid OrderReplace, order is undefined',
      ],
      [
        'missing order id',
        {
          orderReplace: {
            ...redisTestConstants.orderReplace.orderReplace,
            order: {
              ...redisTestConstants.defaultOrder,
              orderId: undefined,
            },
          },
        },
        'Invalid OrderReplace, order id is undefined',
      ],
      [
        'missing order id',
        {
          orderReplace: {
            ...redisTestConstants.orderReplace.orderReplace,
            order: {
              ...redisTestConstants.defaultOrder,
              orderId: {
                ...redisTestConstants.defaultOrderId,
                subaccountId: undefined,
              },
            },
          },
        },
        'Invalid OrderReplace, subaccount id is undefined',
      ],
      [
        'unspecified placement status',
        {
          ...redisTestConstants.orderReplace,
          orderReplace: {
            ...redisTestConstants.orderReplace.orderReplace,
            placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_UNSPECIFIED,
          },
        },
        'Invalid OrderReplace, placement status is UNSPECIFIED',
      ],
    ])('logs error and does not update caches on invalid order replacement off-chain update: %s', async (
      _name: string,
      updateMessage: any,
      errorMsg: string,
    ) => {
      const update: OffChainUpdateV1 = {
        ...updateMessage,
      };
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(OffChainUpdateV1.encode(update).finish())),
      );

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      await onMessage(message);

      expect(redisPackage.placeOrder).not.toHaveBeenCalled();
      expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();
      expectWebsocketMessagesNotSent(producerSendSpy);

      expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
        at: 'OrderReplaceHandler#logAndThrowParseMessageError',
        message: errorMsg,
      }));
    });

    it('logs error and does not update caches if order has invalid clobPairId', async () => {
      const update: OffChainUpdateV1 = {
        orderReplace: {
          ...redisTestConstants.orderReplace.orderReplace,
          order: {
            ...redisTestConstants.defaultOrder,
            orderId: {
              ...redisTestConstants.defaultOrderId,
              clobPairId: 34,
            },
          },
          placementStatus:
              OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
        },
      };
      const message: KafkaMessage = createKafkaMessage(
        Buffer.from(Uint8Array.from(OffChainUpdateV1.encode(update).finish())),
      );

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      await onMessage(message);

      expect(redisPackage.placeOrder).not.toHaveBeenCalled();
      expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();
      expectWebsocketMessagesNotSent(producerSendSpy);

      expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
        at: 'OrderReplaceHandler#logAndThrowParseMessageError',
        message: 'Order in OrderReplace has invalid clobPairId',
      }));
    });
  });
});

async function checkOrderReplace(
  oldOrderId: string,
  placedOrderId: string,
  placedSubaccountId: string,
  expectedOrder: RedisOrder,
): Promise<void> {
  const oldRedisOrder: RedisOrder | null = await OrdersCache.getOrder(oldOrderId, client);
  expect(oldRedisOrder).toBeNull();

  const newRedisOrder: RedisOrder | null = await OrdersCache.getOrder(placedOrderId, client);
  const orderIdsForSubaccount: string[] = await SubaccountOrderIdsCache.getOrderIdsForSubaccount(
    placedSubaccountId,
    client,
  );

  expect(newRedisOrder).toEqual(expectedOrder);
  expect(orderIdsForSubaccount).toEqual([placedOrderId]);
}

function expectStats(orderWasReplaced: boolean = false): void {
  let className: string = 'OrderPlaceHandler';

  if (orderWasReplaced) {
    className = 'OrderReplaceHandler';
  }
  expect(stats.timing).toHaveBeenCalledWith(
    `vulcan.${STATS_FUNCTION_NAME}.timing`,
    expect.any(Number),
    { className, fnName: 'place_order_cache_update' },
  );

  expect(stats.increment).not.toHaveBeenCalledWith(
    'vulcan.place_order_handler.replaced_order',
    expect.any(Number),
  );
}

function expectWebsocketMessagesSent(
  producerSendSpy: jest.SpyInstance,
  redisOrder: RedisOrder,
  dbOrder: OrderFromDatabase,
  perpetualMarket: PerpetualMarketFromDatabase,
  placementStatus: APIOrderStatus,
  expectSubaccountMessage: boolean,
  expectedOrderbookMessage?: OrderbookMessage,
  expectedOffchainUpdate?: OffchainUpdateRecord,
): void {
  jest.runOnlyPendingTimers();
  // expect one subaccount update message being sent
  let numMessages: number = 0;
  if (expectSubaccountMessage) {
    numMessages += 1;
  }
  if (expectedOrderbookMessage !== undefined) {
    numMessages += 1;
  }
  if (expectedOffchainUpdate !== undefined) {
    numMessages += 1;
  }
  expect(producerSendSpy).toHaveBeenCalledTimes(numMessages);

  let callIndex: number = 0;

  if (expectedOffchainUpdate) {
    expectOffchainUpdateMessage(
      producerSendSpy.mock.calls[callIndex][0],
      expectedOffchainUpdate.key,
      expectedOffchainUpdate.offchainUpdate,
    );
    callIndex += 1;
  }

  if (expectSubaccountMessage) {
    const orderTIF: TimeInForce = protocolTranslations.protocolOrderTIFToTIF(
      redisOrder.order!.timeInForce,
    );
    const isStateful: boolean = isStatefulOrder(redisOrder.order!.orderId!.orderFlags);
    const contents: SubaccountMessageContents = {
      orders: [
        {
          id: OrderTable.orderIdToUuid(
            redisOrder.order!.orderId!,
          ),
          subaccountId: SubaccountTable.subaccountIdToUuid(
            redisOrder.order!.orderId!.subaccountId!,
          ),
          clientId: redisOrder.order!.orderId!.clientId.toString(),
          clobPairId: perpetualMarket.clobPairId,
          side: protocolTranslations.protocolOrderSideToOrderSide(redisOrder.order!.side),
          size: redisOrder.size,
          price: redisOrder.price,
          status: placementStatus,
          type: protocolTranslations.protocolConditionTypeToOrderType(
            redisOrder.order!.conditionType,
          ),
          timeInForce: apiTranslations.orderTIFToAPITIF(orderTIF),
          postOnly: apiTranslations.isOrderTIFPostOnly(orderTIF),
          reduceOnly: redisOrder.order!.reduceOnly,
          orderFlags: redisOrder.order!.orderId!.orderFlags.toString(),
          goodTilBlock: protocolTranslations.getGoodTilBlock(redisOrder.order!)
            ?.toString(),
          goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(redisOrder.order!),
          ticker: redisOrder.ticker,
          ...(isStateful && { createdAtHeight: dbOrder.createdAtHeight }),
          ...(isStateful && { updatedAt: dbOrder.updatedAt }),
          ...(isStateful && { updatedAtHeight: dbOrder.updatedAtHeight }),
          clientMetadata: redisOrder.order!.clientMetadata.toString(),
          triggerPrice: getTriggerPrice(redisOrder.order!, perpetualMarket),
        },
      ],
      blockHeight: blockHeightRefresher.getLatestBlockHeight(),
    };
    const subaccountMessage: SubaccountMessage = SubaccountMessage.fromPartial({
      contents: JSON.stringify(contents),
      subaccountId: SubaccountId.fromPartial(redisOrder.order!.orderId!.subaccountId!),
      version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
    });

    expectWebsocketSubaccountMessage(
      producerSendSpy.mock.calls[callIndex][0],
      subaccountMessage,
      defaultKafkaHeaders,
    );
    callIndex += 1;
  }

  if (expectedOrderbookMessage !== undefined) {
    expectWebsocketOrderbookMessage(
      producerSendSpy.mock.calls[callIndex][0],
      expectedOrderbookMessage,
    );
  }
}

function expectWebsocketMessagesNotSent(
  producerSendSpy: jest.SpyInstance,
): void {
  expect(producerSendSpy).not.toBeCalled();
}
