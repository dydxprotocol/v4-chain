import {
  OrderbookLevelsCache,
  OrderData,
  OrdersDataCache,
  redis,
  redisTestConstants,
  StatefulOrderUpdatesCache,
} from '@dydxprotocol-indexer/redis';
import {
  expectOrderbookLevelCache,
  handleInitialOrderPlace,
  handleOrderUpdate,
} from '../helpers/helpers';
import { redisClient as client } from '../../src/helpers/redis/redis-controller';
import {
  blockHeightRefresher,
  dbHelpers,
  OrderbookMessageContents,
  perpetualMarketRefresher,
  protocolTranslations,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import {
  logger,
  stats,
  STATS_FUNCTION_NAME,
  wrapBackgroundTask,
} from '@dydxprotocol-indexer/base';
import { synchronizeWrapBackgroundTask } from '@dydxprotocol-indexer/dev';
import {
  IndexerOrder,
  OrderbookMessage,
  IndexerOrderId,
  OrderPlaceV1_OrderPlacementStatus,
  RedisOrder,
  OrderUpdateV1,
} from '@dydxprotocol-indexer/v4-protos';
import * as redisPackage from '@dydxprotocol-indexer/redis';
import {
  ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
  producer,
} from '@dydxprotocol-indexer/kafka';
import { ProducerRecord } from 'kafkajs';
import { expectWebsocketOrderbookMessage } from '../helpers/websocket-helpers';
import { OrderbookSide } from '../../src/lib/types';
import Long from 'long';

jest.mock('@dydxprotocol-indexer/base', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/base'),
  wrapBackgroundTask: jest.fn(),
}));

describe('OrderUpdateHandler', () => {
  describe('handle', () => {
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
      jest.spyOn(OrderbookLevelsCache, 'updatePriceLevel');
      jest.spyOn(redisPackage, 'updateOrder');
      jest.spyOn(logger, 'error');
      jest.spyOn(logger, 'info');
      jest.spyOn(OrderbookLevelsCache, 'updatePriceLevel');
    });

    afterEach(async () => {
      await dbHelpers.clearData();
      await redis.deleteAllAsync(client);
      jest.restoreAllMocks();
    });

    afterAll(async () => {
      jest.useRealTimers();
      await dbHelpers.teardown();
    });

    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrder,
        redisTestConstants.defaultOrderId,
        redisTestConstants.defaultRedisOrder,
      ],
      [
        'goodTilBlockTime',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderIdGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
      ],
    ])(
      'updates new order (with %s) (not resting on book) total filled and orderbook price level',
      async (
        _name: string,
        initialOrderToPlace: IndexerOrder,
        updatedOrderId: IndexerOrderId,
        updatedRedisOrder: RedisOrder,
      ) => {
        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
        const orderUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage = {
          ...redisTestConstants.orderUpdate,
          orderUpdate: {
            ...redisTestConstants.orderUpdate.orderUpdate,
            orderId: updatedOrderId,
          },
        };

        // Create a new order by handling an order place message
        await handleInitialOrderPlace({
          ...redisTestConstants.orderPlace,
          orderPlace: {
            order: initialOrderToPlace,
            placementStatus:
              OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
          },
        });
        jest.runOnlyPendingTimers();
        jest.clearAllMocks();

        const expectedPriceLevelSize: string = (
          initialOrderToPlace.quantums.sub(orderUpdate.orderUpdate.totalFilledQuantums)
        ).toString();
        await handleOrderUpdate(orderUpdate);

        expectWebsocketMessagesSent(
          producerSendSpy,
          orderbookMessageFromBidContents([
            updatedRedisOrder.price,
            protocolTranslations.quantumsToHumanFixedString(
              (initialOrderToPlace.quantums.sub(
                orderUpdate.orderUpdate.totalFilledQuantums)).toString(),
              testConstants.defaultPerpetualMarket.atomicResolution,
            ),
          ]),
        );
        await expectOrdersDataCache(orderUpdate);
        await expectOrderbookLevelCache(
          testConstants.defaultPerpetualMarket.ticker,
          protocolTranslations.protocolOrderSideToOrderSide(initialOrderToPlace.side),
          updatedRedisOrder.price,
          expectedPriceLevelSize,
        );
        expectTimingStats();
      },
    );

    it(
      'updates existing order (resting on book) total filled and orderbook price level',
      async () => {
        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
        // Create a new order by handling an order place message
        await handleInitialOrderPlace(redisTestConstants.orderPlace);
        jest.runOnlyPendingTimers();
        jest.clearAllMocks();
        // Handle the first update to the order, updating the order to be resting on the book
        // and updating the price level for the order
        await handleOrderUpdate(redisTestConstants.orderUpdate);

        expectWebsocketMessagesSent(
          producerSendSpy,
          orderbookMessageFromBidContents([
            redisTestConstants.defaultRedisOrder.price,
            protocolTranslations.quantumsToHumanFixedString(
              (redisTestConstants.defaultOrder.quantums.sub(
                redisTestConstants.orderUpdate.orderUpdate.totalFilledQuantums,
              )).toNumber().toString(),
              testConstants.defaultPerpetualMarket.atomicResolution,
            ),
          ]),
        );

        jest.clearAllMocks();
        const secondTotalFilledQuantums: Long = Long.fromValue(500_350, true);
        const secondUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage = {
          ...redisTestConstants.orderUpdate,
          orderUpdate: {
            ...redisTestConstants.orderUpdate.orderUpdate,
            totalFilledQuantums: secondTotalFilledQuantums,
          },
        };
        const expectedPriceLevelSize: string = (
          redisTestConstants.defaultOrder.quantums.sub(secondUpdate.orderUpdate.totalFilledQuantums)
        ).toString();
        await handleOrderUpdate(secondUpdate);
        expectWebsocketMessagesSent(
          producerSendSpy,
          orderbookMessageFromBidContents([
            redisTestConstants.defaultRedisOrder.price,
            protocolTranslations.quantumsToHumanFixedString(
              (redisTestConstants.defaultOrder.quantums.sub(secondTotalFilledQuantums))
                .toNumber().toString(),
              testConstants.defaultPerpetualMarket.atomicResolution,
            ),
          ]),
        );
        await expectOrdersDataCache(secondUpdate);
        await expectOrderbookLevelCache(
          testConstants.defaultPerpetualMarket.ticker,
          protocolTranslations.protocolOrderSideToOrderSide(redisTestConstants.defaultOrder.side),
          redisTestConstants.defaultRedisOrder.price,
          expectedPriceLevelSize,
        );
        expectTimingStats();

        const thirdTotalFilledQuantums: Long = redisTestConstants.defaultOrder.quantums;
        const thirdUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage = {
          ...redisTestConstants.orderUpdate,
          orderUpdate: {
            ...redisTestConstants.orderUpdate.orderUpdate,
            totalFilledQuantums: thirdTotalFilledQuantums,
          },
        };
        await handleOrderUpdate(thirdUpdate);
      },
    );

    it.each([
      [
        'Fill-or-Kill',
        redisTestConstants.defaultOrderFok,
        redisTestConstants.defaultOrderId,
        redisTestConstants.defaultRedisOrderFok,
      ],
      [
        'Immediate-or-Cancel',
        redisTestConstants.defaultOrderIoc,
        redisTestConstants.defaultOrderId,
        redisTestConstants.defaultRedisOrderIoc,
      ],
    ])(
      'updates new order (with %s) (not resting on book) total filled but does not update order ' +
      'book for IOC order',
      async (
        _name: string,
        initialOrderToPlace: IndexerOrder,
        updatedOrderId: IndexerOrderId,
        updatedRedisOrder: RedisOrder,
      ) => {
        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
        const orderUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage = {
          ...redisTestConstants.orderUpdate,
          orderUpdate: {
            ...redisTestConstants.orderUpdate.orderUpdate,
            orderId: updatedOrderId,
          },
        };

        // Create a new order by handling an order place message
        await handleInitialOrderPlace({
          ...redisTestConstants.orderPlace,
          orderPlace: {
            order: initialOrderToPlace,
            placementStatus:
              OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
          },
        });
        jest.runOnlyPendingTimers();
        jest.clearAllMocks();
        await handleOrderUpdate(orderUpdate);

        // No order book update websocket messages should be sent
        expectWebsocketMessagesNotSent(producerSendSpy);
        await expectOrdersDataCache(orderUpdate);
        // Price-level should be 0 as placing the order and updating the order should not have
        // changed the order book
        await expectOrderbookLevelCache(
          testConstants.defaultPerpetualMarket.ticker,
          protocolTranslations.protocolOrderSideToOrderSide(initialOrderToPlace.side),
          updatedRedisOrder.price,
          '0',
        );
      },
    );

    it(
      'logs error, and caps price level update if new total filled quantums exceed order size',
      async () => {
        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
        // Create a new order by handling an order place message
        await handleInitialOrderPlace(redisTestConstants.orderPlace);
        jest.runOnlyPendingTimers();
        jest.clearAllMocks();
        // Handle the first update to the order, updating the order to be resting on the book
        // and updating the price level for the order
        await handleOrderUpdate(redisTestConstants.orderUpdate);

        expectWebsocketMessagesSent(
          producerSendSpy,
          orderbookMessageFromBidContents([
            redisTestConstants.defaultRedisOrder.price,
            protocolTranslations.quantumsToHumanFixedString(
              (redisTestConstants.defaultOrder.quantums.sub(
                redisTestConstants.orderUpdate.orderUpdate.totalFilledQuantums)).toString(),
              testConstants.defaultPerpetualMarket.atomicResolution,
            ),
          ]),
        );
        jest.clearAllMocks();
        const exceedsFilledUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage = {
          ...redisTestConstants.orderUpdate,
          orderUpdate: {
            ...redisTestConstants.orderUpdate.orderUpdate,
            totalFilledQuantums: redisTestConstants.defaultOrder.quantums.add(
              Long.fromValue(100, true),
            ),
          },
        };
        await handleOrderUpdate(exceedsFilledUpdate);

        expectWebsocketMessagesSent(
          producerSendSpy,
          orderbookMessageFromBidContents([
            redisTestConstants.defaultRedisOrder.price,
            protocolTranslations.quantumsToHumanFixedString(
              '0',
              testConstants.defaultPerpetualMarket.atomicResolution,
            ),
          ]),
        );

        await expectOrdersDataCache(exceedsFilledUpdate);
        await expectOrderbookLevelCache(
          testConstants.defaultPerpetualMarket.ticker,
          protocolTranslations.protocolOrderSideToOrderSide(redisTestConstants.defaultOrder.side),
          redisTestConstants.defaultRedisOrder.price,
          '0', // Should be 0 since total filled quantums > order size in quantums
        );

        expectTimingStats();
        expect(stats.increment).toHaveBeenCalledWith(
          'vulcan.order_update_total_filled_exceeds_size',
          1,
          { instance: '' },
        );
        expect(logger.info).toHaveBeenCalledWith(expect.objectContaining({
          at: 'OrderUpdateHandler#getCappedNewTotalFilledQuantums',
          message: 'New total filled quantums of order exceeds order size in quantums.',
        }));
      },
    );

    it(
      'Caps price level update if old total filled quantums exceed order size',
      async () => {
        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
        // Create a new order by handling an order place message
        await handleInitialOrderPlace(redisTestConstants.orderPlace);
        jest.runOnlyPendingTimers();
        jest.clearAllMocks();
        // Handle the first update to the order, updating the order to be resting on the book
        // and updating the price level for the order with an update the exceeds the size of the
        // order
        const exceedsFilledUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage = {
          ...redisTestConstants.orderUpdate,
          orderUpdate: {
            ...redisTestConstants.orderUpdate.orderUpdate,
            totalFilledQuantums: redisTestConstants.defaultOrder.quantums.add(
              Long.fromValue(100, true),
            ),
          },
        };
        await handleOrderUpdate(exceedsFilledUpdate);
        // Size delta will be zero as the order went from unfilled (0 quantums) to fully-filled
        // and no update should have been made to the orderbook levels cache
        expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();
        jest.runOnlyPendingTimers();
        // No websocket update messages should be sent for an update with a size delta of 0
        expect(producerSendSpy).not.toHaveBeenCalled();
        jest.clearAllMocks();

        const expectedPriceLevelSize: string = (
          redisTestConstants.defaultOrder.quantums.sub(
            redisTestConstants.orderUpdate.orderUpdate.totalFilledQuantums)
        ).toString();
        await handleOrderUpdate(redisTestConstants.orderUpdate);
        expectWebsocketMessagesSent(
          producerSendSpy,
          orderbookMessageFromBidContents([
            redisTestConstants.defaultRedisOrder.price,
            protocolTranslations.quantumsToHumanFixedString(
              (redisTestConstants.defaultOrder.quantums.sub(
                redisTestConstants.orderUpdate.orderUpdate.totalFilledQuantums)).toString(),
              testConstants.defaultPerpetualMarket.atomicResolution,
            ),
          ]),
        );

        await expectOrdersDataCache(redisTestConstants.orderUpdate);
        await expectOrderbookLevelCache(
          testConstants.defaultPerpetualMarket.ticker,
          protocolTranslations.protocolOrderSideToOrderSide(redisTestConstants.defaultOrder.side),
          redisTestConstants.defaultRedisOrder.price,
          expectedPriceLevelSize,
        );

        expect(logger.info).toHaveBeenCalledWith(expect.objectContaining({
          at: 'OrderUpdateHandler#getCappedOldTotalFilledQuantums',
          message: 'Old total filled quantums of order exceeds order size in quantums.',
        }));
        expect(stats.increment).toHaveBeenCalledWith(
          'vulcan.order_update_old_total_filled_exceeds_size',
          1,
          { instance: '' },
        );
        expectTimingStats();
      },
    );

    it('Does not update orderbook cache or send messages if size delta is 0', async () => {
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      // Create a new order by handling an order place message
      await handleInitialOrderPlace(redisTestConstants.orderPlace);
      jest.runOnlyPendingTimers();
      jest.clearAllMocks();

      const expectedPriceLevelSize: string = (
        redisTestConstants.defaultOrder.quantums.sub(
          redisTestConstants.orderUpdate.orderUpdate.totalFilledQuantums)
      ).toString();
      await handleOrderUpdate(redisTestConstants.orderUpdate);
      expectWebsocketMessagesSent(
        producerSendSpy,
        orderbookMessageFromBidContents([
          redisTestConstants.defaultRedisOrder.price,
          protocolTranslations.quantumsToHumanFixedString(
            (redisTestConstants.defaultOrder.quantums.sub(
              redisTestConstants.orderUpdate.orderUpdate.totalFilledQuantums)).toString(),
            testConstants.defaultPerpetualMarket.atomicResolution,
          ),
        ]),
      );
      jest.clearAllMocks();

      await expectOrdersDataCache(redisTestConstants.orderUpdate);
      await expectOrderbookLevelCache(
        testConstants.defaultPerpetualMarket.ticker,
        protocolTranslations.protocolOrderSideToOrderSide(redisTestConstants.defaultOrder.side),
        redisTestConstants.defaultRedisOrder.price,
        expectedPriceLevelSize,
      );

      // Handle the same update again, size delta should be zero
      await handleOrderUpdate(redisTestConstants.orderUpdate);

      // No update should have been made to the orderbook levels cache since the size delta is 0
      expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();
      jest.runOnlyPendingTimers();
      // No websocket update messages should be sent for an update with a size delta of 0
      expect(producerSendSpy).not.toHaveBeenCalled();

      await expectOrdersDataCache(redisTestConstants.orderUpdate);
      await expectOrderbookLevelCache(
        testConstants.defaultPerpetualMarket.ticker,
        protocolTranslations.protocolOrderSideToOrderSide(redisTestConstants.defaultOrder.side),
        redisTestConstants.defaultRedisOrder.price,
        expectedPriceLevelSize,
      );
      expectTimingStat('update_order_cache_update');
    });

    it.each([
      [
        'missing orderId',
        {
          ...redisTestConstants.orderUpdate,
          orderUpdate: {
            ...redisTestConstants.orderUpdate.orderUpdate,
            orderId: undefined,
          },
        },
        'Invalid OrderUpdate, order id is undefined',
      ],
      [
        'missing subaccountId',
        {
          ...redisTestConstants.orderUpdate,
          orderUpdate: {
            ...redisTestConstants.orderUpdate.orderUpdate,
            orderId: {
              ...redisTestConstants.defaultOrderId,
              subaccountId: undefined,
            },
          },
        },
        'Invalid OrderUpdate, subaccount id is undefined',
      ],
    ])('logs error and does not update order on invalid OrderUpdate: %s', async (
      _name: string,
      updateMessage: any,
      errorMsg: string,
    ) => {
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      // Create a new order by handling an order place message
      await handleInitialOrderPlace(redisTestConstants.orderPlace);
      jest.clearAllMocks();

      await handleOrderUpdate(updateMessage);
      expectWebsocketMessagesNotSent(producerSendSpy);

      expect(redisPackage.updateOrder).not.toHaveBeenCalled();
      expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();
      expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
        at: 'OrderUpdateHandler#logAndThrowParseMessageError',
        message: errorMsg,
      }));
    });

    it('logs error and does not update OrderbookLevels if short-term order not found', async () => {
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      await handleOrderUpdate(redisTestConstants.orderUpdate);

      expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();
      expect(logger.info).toHaveBeenCalledWith(expect.objectContaining({
        at: 'OrderUpdateHandler#handle',
        message: expect.stringMatching('Received order update for order that does not exist, order id '),
      }));
      expectWebsocketMessagesNotSent(producerSendSpy);
      expect(stats.increment).toHaveBeenCalledWith(
        'vulcan.order_update_order_does_not_exist',
        1,
        {
          orderFlags: String(redisTestConstants.orderUpdate.orderUpdate.orderId!.orderFlags),
          instance: '',
        },
      );
    });

    it('adds order update to stateful order update cache if stateful order not found', async () => {
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      const statefulOrderUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage = {
        ...redisTestConstants.orderUpdate,
        orderUpdate: {
          ...redisTestConstants.orderUpdate.orderUpdate,
          orderId: redisTestConstants.defaultOrderIdGoodTilBlockTime,
        },
      };
      await handleOrderUpdate(statefulOrderUpdate);

      const cachedOrderUpdate: OrderUpdateV1 | undefined = await StatefulOrderUpdatesCache
        .removeStatefulOrderUpdate(
          redisTestConstants.defaultOrderUuidGoodTilBlockTime,
          Date.now(),
          client,
        );
      expect(cachedOrderUpdate).toBeDefined();
      expect(cachedOrderUpdate).toEqual(statefulOrderUpdate.orderUpdate);

      expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();
      expect(logger.info).toHaveBeenCalledWith(expect.objectContaining({
        at: 'OrderUpdateHandler#handle',
        message: expect.stringMatching('Received order update for order that does not exist, order id '),
      }));
      expectWebsocketMessagesNotSent(producerSendSpy);
      expect(stats.increment).toHaveBeenCalledWith(
        'vulcan.order_update_order_does_not_exist',
        1,
        {
          orderFlags: String(statefulOrderUpdate.orderUpdate.orderId!.orderFlags),
          instance: '',
        },
      );
    });

    it('adds order update to stateful order update cache if vault order not found', async () => {
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const vaultOrderUpdate: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage = {
        ...redisTestConstants.orderUpdate,
        orderUpdate: {
          ...redisTestConstants.orderUpdate.orderUpdate,
          orderId: redisTestConstants.defaultOrderIdVault,
        },
      };
      await handleOrderUpdate(vaultOrderUpdate);

      const cachedOrderUpdate: OrderUpdateV1 | undefined = await StatefulOrderUpdatesCache
        .removeStatefulOrderUpdate(
          redisTestConstants.defaultOrderUuidVault,
          Date.now(),
          client,
        );
      expect(cachedOrderUpdate).toBeDefined();
    });
  });
});

async function expectOrdersDataCache(
  updateMessage: redisTestConstants.OffChainUpdateOrderUpdateUpdateMessage,
): Promise<void> {
  const orderData: OrderData | null = await OrdersDataCache.getOrderData(
    updateMessage.orderUpdate.orderId!,
    client,
  );
  expect(orderData).toBeDefined();
  expect(orderData!.totalFilledQuantums).toEqual(
    updateMessage.orderUpdate.totalFilledQuantums.toString(),
  );
  expect(orderData!.restingOnBook).toBe(true);
}

function expectTimingStats(): void {
  expectTimingStat('update_order_cache_update');
  expectTimingStat('update_price_level');
}

function expectTimingStat(fnName: string) {
  expect(stats.timing).toHaveBeenCalledWith(
    `vulcan.${STATS_FUNCTION_NAME}.timing`,
    expect.any(Number),
    { className: 'OrderUpdateHandler', fnName, instance: '' },
  );
}

function orderbookMessageFromBidContents(
  bidContents: any,
): OrderbookMessage {
  const contents: OrderbookMessageContents = {
    [OrderbookSide.BIDS]: [bidContents],
  };
  return OrderbookMessage.fromPartial({
    contents: JSON.stringify(contents),
    clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
    version: ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
  });
}

function expectWebsocketMessagesSent(
  producerSendSpy: jest.SpyInstance,
  expectedOrderbookMessage: OrderbookMessage,
): void {
  jest.runOnlyPendingTimers();
  // expect one call for subaccount and one for orderbook
  expect(producerSendSpy).toHaveBeenCalledTimes(1);

  const orderbookProducerRecord: ProducerRecord = producerSendSpy.mock.calls[0][0];
  expectWebsocketOrderbookMessage(orderbookProducerRecord, expectedOrderbookMessage);
}

function expectWebsocketMessagesNotSent(
  producerSendSpy: jest.SpyInstance,
): void {
  expect(producerSendSpy).not.toBeCalled();
}
