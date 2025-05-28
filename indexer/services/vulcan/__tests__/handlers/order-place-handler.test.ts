import {
  logger, stats, STATS_FUNCTION_NAME, wrapBackgroundTask,
} from '@dydxprotocol-indexer/base';
import { synchronizeWrapBackgroundTask } from '@dydxprotocol-indexer/dev';
import {
  createKafkaMessage,
  getTriggerPrice,
  producer,
  SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import {
  APIOrderStatus,
  APIOrderStatusEnum,
  apiTranslations,
  blockHeightRefresher,
  BlockTable,
  dbHelpers,
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
  CanceledOrdersCache,
  CanceledOrderStatus,
  OrderbookLevels,
  OrderbookLevelsCache,
  OrdersCache,
  PriceLevel,
  redis,
  redisTestConstants,
  StatefulOrderUpdatesCache,
  SubaccountOrderIdsCache,
  updateOrder,
} from '@dydxprotocol-indexer/redis';
import {
  IndexerOrder,
  OffChainUpdateV1,
  OrderbookMessage,
  OrderPlaceV1_OrderPlacementStatus,
  OrderUpdateV1,
  RedisOrder,
  SubaccountId,
  SubaccountMessage,
} from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import Long from 'long';
import { redisClient, redisClient as client } from '../../src/helpers/redis/redis-controller';
import { onMessage } from '../../src/lib/on-message';
import { expectCanceledOrderStatus, handleInitialOrderPlace } from '../helpers/helpers';
import {
  expectOffchainUpdateMessage,
  expectWebsocketOrderbookMessage,
  expectWebsocketSubaccountMessage,
} from '../helpers/websocket-helpers';
import { getOrderIdHash, isStatefulOrder } from '@dydxprotocol-indexer/v4-proto-parser';
import { defaultKafkaHeaders } from '../helpers/constants';

jest.mock('@dydxprotocol-indexer/base', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/base'),
  wrapBackgroundTask: jest.fn(),
}));

interface OffchainUpdateRecord {
  key: Buffer,
  offchainUpdate: OffChainUpdateV1,
}

describe('order-place-handler', () => {
  beforeAll(async () => {
    await BlockTable.create(testConstants.defaultBlock);
    await blockHeightRefresher.updateBlockHeight();
    jest.useFakeTimers();
  });

  afterAll(() => {
    jest.useRealTimers();
  });

  describe('handle', () => {
    const replacementOrder: IndexerOrder = {
      ...redisTestConstants.defaultOrder,
      goodTilBlock: 1160,
      goodTilBlockTime: undefined,
      quantums: Long.fromValue(500_000, true),
      subticks: Long.fromValue(1_000_000, true),
    };
    const replacementOrderGoodTilBlockTime: IndexerOrder = {
      ...replacementOrder,
      orderId: redisTestConstants.defaultOrderIdGoodTilBlockTime,
      goodTilBlock: undefined,
      goodTilBlockTime: 1_300_000_000,
    };
    const replacementOrderConditional: IndexerOrder = {
      ...redisTestConstants.defaultConditionalOrder,
      quantums: Long.fromValue(500_000, true),
      subticks: Long.fromValue(1_000_000, true),
      goodTilBlock: undefined,
      goodTilBlockTime: 1_300_000_000,
    };
    const replacementOrderFok: IndexerOrder = {
      ...redisTestConstants.defaultOrderFok,
      goodTilBlock: 1160,
      goodTilBlockTime: undefined,
      quantums: Long.fromValue(500_000, true),
      subticks: Long.fromValue(1_000_000, true),
    };
    const replacementOrderIoc: IndexerOrder = {
      ...redisTestConstants.defaultOrderIoc,
      goodTilBlock: 1160,
      goodTilBlockTime: undefined,
      quantums: Long.fromValue(500_000, true),
      subticks: Long.fromValue(1_000_000, true),
    };
    const replacedOrder: RedisOrder = redisPackage.convertToRedisOrder(
      replacementOrder,
      testConstants.defaultPerpetualMarket,
    );
    const replacedOrderGoodTilBlockTime: RedisOrder = redisPackage.convertToRedisOrder(
      replacementOrderGoodTilBlockTime,
      testConstants.defaultPerpetualMarket,
    );
    const replacedOrderConditional: RedisOrder = redisPackage.convertToRedisOrder(
      replacementOrderConditional,
      testConstants.defaultPerpetualMarket,
    );
    const replacedOrderFok: RedisOrder = redisPackage.convertToRedisOrder(
      replacementOrderFok,
      testConstants.defaultPerpetualMarket,
    );
    const replacedOrderIoc: RedisOrder = redisPackage.convertToRedisOrder(
      replacementOrderIoc,
      testConstants.defaultPerpetualMarket,
    );
    const replacementUpdate: OffChainUpdateV1 = {
      orderPlace: {
        order: replacementOrder,
        placementStatus:
        OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
      },
    };
    const replacementUpdateGoodTilBlockTime: OffChainUpdateV1 = {
      orderPlace: {
        order: replacementOrderGoodTilBlockTime,
        placementStatus:
        OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
      },
    };
    const replacementUpdateConditional: OffChainUpdateV1 = {
      orderPlace: {
        order: replacementOrderConditional,
        placementStatus:
        OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
      },
    };
    const replacementUpdateFok: OffChainUpdateV1 = {
      orderPlace: {
        order: replacementOrderFok,
        placementStatus:
        OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
      },
    };
    const replacementUpdateIoc: OffChainUpdateV1 = {
      orderPlace: {
        order: replacementOrderIoc,
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
    const replacementMessageConditional: KafkaMessage = createKafkaMessage(
      Buffer.from(Uint8Array.from(
        OffChainUpdateV1.encode(replacementUpdateConditional).finish(),
      )),
    );
    const replacementMessageFok: KafkaMessage = createKafkaMessage(
      Buffer.from(Uint8Array.from(OffChainUpdateV1.encode(replacementUpdateFok).finish())),
    );
    const replacementMessageIoc: KafkaMessage = createKafkaMessage(
      Buffer.from(Uint8Array.from(OffChainUpdateV1.encode(replacementUpdateIoc).finish())),
    );
    [replacementMessage, replacementMessageGoodTilBlockTime, replacementMessageConditional,
      replacementMessageFok, replacementMessageIoc].forEach((message) => {
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
    const dbConditionalOrder: OrderFromDatabase = {
      ...testConstants.defaultConditionalOrder,
      id: testConstants.defaultConditionalOrderId,
      createdAtHeight: '3',
    };
    const dbDefaultOrderFok: OrderFromDatabase = {
      ...testConstants.defaultOrder,
      id: testConstants.defaultOrderId,
      timeInForce: TimeInForce.FOK,
    };
    const dbDefaultOrderIoc: OrderFromDatabase = {
      ...testConstants.defaultOrder,
      id: testConstants.defaultOrderId,
      timeInForce: TimeInForce.IOC,
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
        OrderTable.create(dbConditionalOrder),
      ]);
      jest.spyOn(stats, 'timing');
      jest.spyOn(OrderbookLevelsCache, 'updatePriceLevel');
      jest.spyOn(CanceledOrdersCache, 'removeOrderFromCaches');
      jest.spyOn(StatefulOrderUpdatesCache, 'removeStatefulOrderUpdate');
      jest.spyOn(OrderTable, 'findById');
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
        'good-til-block',
        redisTestConstants.defaultOrder,
        redisTestConstants.defaultRedisOrder,
        redisTestConstants.defaultOrderUuid,
        true,
      ],
      [
        'good-til-block-time',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        // Subaccount messages should not be sent for best-effort-opened for long-term orders
        false,
      ],
      [
        'conditional',
        redisTestConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
        // Subaccount messages should not be sent for best-effort-opened for conditional orders
        false,
      ],
    ])('handles order place for new order (with %s)', async (
      _name: string,
      orderToPlace: IndexerOrder,
      expectedRedisOrder: RedisOrder,
      expectedOrderUuid: string,
      expectSubaccountMessageSent: boolean,
    ) => {
      const expectedOrder: RedisOrder = redisPackage.convertToRedisOrder(
        orderToPlace,
        testConstants.defaultPerpetualMarket,
      );
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      await handleInitialOrderPlace({
        ...redisTestConstants.orderPlace,
        orderPlace: {
          order: orderToPlace,
          placementStatus:
          OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
        },
      });

      expectWebsocketMessagesSent(
        producerSendSpy,
        expectedRedisOrder,
        dbDefaultOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.BEST_EFFORT_OPENED,
        expectSubaccountMessageSent,
      );

      await checkOrderPlace(
        expectedOrderUuid,
        redisTestConstants.defaultSubaccountUuid,
        expectedOrder,
      );
      expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();

      expect(logger.error).not.toHaveBeenCalled();
      expectStats();
    });

    // TODO(IND-68): Remove this test once order replacement logic does not change price levels as
    // orders are removed before being re-placed.
    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrder,
        replacementMessage,
        redisTestConstants.defaultRedisOrder,
        dbDefaultOrder,
        redisTestConstants.defaultOrderUuid,
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
        replacedOrderGoodTilBlockTime,
        false,
        false,
      ],
      [
        'conditional',
        redisTestConstants.defaultConditionalOrder,
        replacementMessageConditional,
        redisTestConstants.defaultRedisOrderConditional,
        dbConditionalOrder,
        redisTestConstants.defaultOrderUuidConditional,
        replacedOrderConditional,
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
        replacedOrderGoodTilBlockTime,
        false,
        true,
      ],
      [
        'conditional and canceled order',
        redisTestConstants.defaultConditionalOrder,
        replacementMessageConditional,
        redisTestConstants.defaultRedisOrderConditional,
        dbConditionalOrder,
        redisTestConstants.defaultOrderUuidConditional,
        replacedOrderConditional,
        false,
        true,
      ],
    ])('handles order place for replacing order (with %s), not resting on book', async (
      _name: string,
      initialOrderToPlace: IndexerOrder,
      orderReplacementMessage: KafkaMessage,
      expectedRedisOrder: RedisOrder,
      dbOrder: OrderFromDatabase,
      expectedOrderUuid: string,
      expectedReplacedOrder: RedisOrder,
      expectSubaccountMessage: boolean,
      hasCanceledOrderId: boolean,
    ) => {
      if (hasCanceledOrderId) {
        await redisPackage.CanceledOrdersCache.addCanceledOrderId(
          expectedOrderUuid,
          Date.now(),
          redisClient,
        );
      }
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      // Handle the order place off-chain update for the initial order
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

      // Handle the order place off-chain update with the replacement order
      await onMessage(orderReplacementMessage);

      await checkOrderPlace(
        expectedOrderUuid,
        redisTestConstants.defaultSubaccountUuid,
        expectedReplacedOrder,
      );
      expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();
      if (hasCanceledOrderId) {
        expect(CanceledOrdersCache.removeOrderFromCaches).toHaveBeenCalled();
      }
      await expectCanceledOrderStatus(expectedOrderUuid, CanceledOrderStatus.NOT_CANCELED);

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

    // TODO(IND-68): Remove this test once order replacement logic does not change price levels as
    // orders are removed before being re-placed.
    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrder,
        replacementMessage,
        redisTestConstants.defaultRedisOrder,
        dbDefaultOrder,
        redisTestConstants.defaultOrderUuid,
        replacedOrder,
        true,
        true,
      ],
      [
        'goodTilBlockTime',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        replacementMessageGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        dbOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        replacedOrderGoodTilBlockTime,
        false,
        true,
      ],
      [
        'conditional',
        redisTestConstants.defaultConditionalOrder,
        replacementMessageConditional,
        redisTestConstants.defaultRedisOrderConditional,
        dbConditionalOrder,
        redisTestConstants.defaultOrderUuidConditional,
        replacedOrderConditional,
        false,
        true,
      ],
      [
        'Fill-or-Kill',
        redisTestConstants.defaultOrderFok,
        replacementMessageFok,
        redisTestConstants.defaultRedisOrderFok,
        dbDefaultOrderFok,
        redisTestConstants.defaultOrderUuid,
        replacedOrderFok,
        true,
        false,
      ],
      [
        'Immediate-or-Cancel',
        redisTestConstants.defaultOrderIoc,
        replacementMessageIoc,
        redisTestConstants.defaultRedisOrderIoc,
        dbDefaultOrderIoc,
        redisTestConstants.defaultOrderUuid,
        replacedOrderIoc,
        true,
        false,
      ],
    ])('handles order place for replacing order (with %s), resting on book', async (
      _name: string,
      initialOrderToPlace: IndexerOrder,
      orderReplacementMessage: KafkaMessage,
      expectedRedisOrder: RedisOrder,
      dbOrder: OrderFromDatabase,
      expectedOrderUuid: string,
      expectedReplacedOrder: RedisOrder,
      expectSubaccountMessage: boolean,
      expectOrderBookUpdate: boolean,
    ) => {
      const oldOrderTotalFilled: number = 10;
      const oldPriceLevelInitialQuantums: number = Number(initialOrderToPlace.quantums) * 2;
      const expectedPriceLevel: PriceLevel = {
        humanPrice: expectedRedisOrder.price,
        quantums: oldPriceLevelInitialQuantums.toString(),
        lastUpdated: expect.stringMatching(/^[0-9]{10}$/),
      };

      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      // Handle the order place event for the initial order
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
      await OrderbookLevelsCache.updatePriceLevel(
        testConstants.defaultPerpetualMarket.ticker,
        protocolTranslations.protocolOrderSideToOrderSide(
          initialOrderToPlace.side,
        ),
        expectedRedisOrder.price,
        oldPriceLevelInitialQuantums.toString(),
        client,
      );

      // Handle the order place off-chain update with the replacement order
      await onMessage(orderReplacementMessage);

      await checkOrderPlace(
        expectedOrderUuid,
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

    // TODO(IND-68): Remove this test once order replacement logic does not change price levels as
    // orders are removed before being re-placed.
    it.each([
      [
        'goodTilBlock',
        redisTestConstants.defaultOrder,
        replacementMessage,
        redisTestConstants.defaultRedisOrder,
        dbDefaultOrder,
        redisTestConstants.defaultOrderUuid,
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
        replacedOrderGoodTilBlockTime,
        false,
      ],
      [
        'conditional',
        redisTestConstants.defaultConditionalOrder,
        replacementMessageConditional,
        redisTestConstants.defaultRedisOrderConditional,
        dbConditionalOrder,
        redisTestConstants.defaultOrderUuidConditional,
        replacedOrderConditional,
        false,
      ],
    ])('handles order place for replacing order (with %s), resting on book, 0 remaining quantums',
      async (
        _name: string,
        initialOrderToPlace: IndexerOrder,
        orderReplacementMessage: KafkaMessage,
        expectedRedisOrder: RedisOrder,
        dbOrder: OrderFromDatabase,
        expectedOrderUuid: string,
        expectedReplacedOrder: RedisOrder,
        expectSubaccountMessage: boolean,
      ) => {
        synchronizeWrapBackgroundTask(wrapBackgroundTask);
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
        // Handle the order place event for the initial order
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
        // Handle the order place off-chain update with the replacement order
        await onMessage(orderReplacementMessage);

        await checkOrderPlace(
          expectedOrderUuid,
          redisTestConstants.defaultSubaccountUuid,
          expectedReplacedOrder,
        );
        expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();

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
        'good-til-block-time and no cached update',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        dbOrderGoodTilBlockTime,
        undefined,
      ],
      [
        'conditional and no cached update',
        redisTestConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
        dbConditionalOrder,
        undefined,
      ],
      [
        'good-til-block-time and cached update',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        dbOrderGoodTilBlockTime,
        {
          ...redisTestConstants.orderUpdate.orderUpdate,
          orderId: redisTestConstants.defaultOrderIdGoodTilBlockTime,
        },
      ],
      [
        'conditional and cached update',
        redisTestConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        redisTestConstants.defaultOrderUuidConditional,
        dbConditionalOrder,
        {
          ...redisTestConstants.orderUpdate.orderUpdate,
          orderId: redisTestConstants.defaultOrderIdConditional,
        },
      ],
    ])('handles order place with OPEN placement status, does not exist initially (with %s)', async (
      _name: string,
      orderToPlace: IndexerOrder,
      expectedRedisOrder: RedisOrder,
      expectedOrderUuid: string,
      placedOrder: OrderFromDatabase,
      cachedOrderUpdate: OrderUpdateV1 | undefined,
    ) => {
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

      let expectedOffchainUpdate: OffchainUpdateRecord | undefined;
      if (cachedOrderUpdate !== undefined) {
        await StatefulOrderUpdatesCache.addStatefulOrderUpdate(
          expectedOrderUuid,
          cachedOrderUpdate,
          Date.now(),
          redisClient,
        );
        expectedOffchainUpdate = {
          key: getOrderIdHash(orderToPlace.orderId!),
          offchainUpdate: {
            orderUpdate: cachedOrderUpdate,
          },
        };
      }

      await handleInitialOrderPlace({
        ...redisTestConstants.orderPlace,
        orderPlace: {
          order: orderToPlace,
          placementStatus:
          OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
        },
      });

      expectWebsocketMessagesSent(
        producerSendSpy,
        expectedRedisOrder,
        placedOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.OPEN,
        // Subaccount message should be sent for stateful order if status is OPEN
        true,
        undefined,
        expectedOffchainUpdate,
      );

      await checkOrderPlace(
        expectedOrderUuid,
        redisTestConstants.defaultSubaccountUuid,
        expectedRedisOrder,
      );

      expect(logger.error).not.toHaveBeenCalled();
      expectStats();
    });

    it.each([
      [
        'good-til-block-time',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        dbOrderGoodTilBlockTime,
      ],
      [
        'conditional',
        redisTestConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        dbConditionalOrder,
      ],
      [
        'good-til-block-time',
        redisTestConstants.defaultOrderGoodTilBlockTime,
        redisTestConstants.defaultRedisOrderGoodTilBlockTime,
        dbOrderGoodTilBlockTime,
      ],
      [
        'conditional',
        redisTestConstants.defaultConditionalOrder,
        redisTestConstants.defaultRedisOrderConditional,
        dbConditionalOrder,
      ],
    ])('handles order place with OPEN placement status, exists initially (with %s)', async (
      _name: string,
      orderToPlace: IndexerOrder,
      expectedRedisOrder: RedisOrder,
      placedOrder: OrderFromDatabase,
    ) => {
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      // Handle the order place event for the initial order with BEST_EFFORT_OPENED
      await handleInitialOrderPlace({
        ...redisTestConstants.orderPlace,
        orderPlace: {
          order: orderToPlace,
          placementStatus:
          OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED,
        },
      });
      expectWebsocketMessagesSent(
        producerSendSpy,
        expectedRedisOrder,
        placedOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.BEST_EFFORT_OPENED,
        // Subaccount message should not be sent for BEST_EFFORT_OPEN
        false,
      );
      expectStats();
      // clear mocks
      jest.clearAllMocks();

      // Handle the order place with OPEN placement status
      await handleInitialOrderPlace({
        ...redisTestConstants.orderPlace,
        orderPlace: {
          order: orderToPlace,
          placementStatus:
          OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
        },
      });
      expectWebsocketMessagesSent(
        producerSendSpy,
        expectedRedisOrder,
        placedOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.OPEN,
        // Subaccount messages should be sent for stateful order with OPEN status
        true,
      );

      expect(logger.error).not.toHaveBeenCalled();
      expectStats();
    });

    it('handles unplaced and unreplaced order place with BEST_EFFORT_OPENED placement status', async () => {
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      // Handle the order place event for the initial order with BEST_EFFORT_OPENED
      await handleInitialOrderPlace(redisTestConstants.orderPlace);
      expectWebsocketMessagesSent(
        producerSendSpy,
        redisTestConstants.defaultRedisOrder,
        dbDefaultOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.BEST_EFFORT_OPENED,
        true,
      );
      expectStats();
      // clear mocks
      jest.clearAllMocks();

      // Handle the order place with OPEN placement status
      await handleInitialOrderPlace(redisTestConstants.orderPlace);
      expectWebsocketMessagesSent(
        producerSendSpy,
        redisTestConstants.defaultRedisOrder,
        dbDefaultOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.BEST_EFFORT_OPENED,
        // Subaccount messages should be sent for stateful order with OPEN status
        false,
      );

      expect(logger.error).not.toHaveBeenCalled();
      expectStats();
    });

    it.each([
      [
        'missing order',
        {
          ...redisTestConstants.orderPlace,
          orderPlace: { ...redisTestConstants.orderPlace.orderPlace, order: undefined },
        },
        'Invalid OrderPlace, order is undefined',
      ],
      [
        'missing order id',
        {
          orderPlace: {
            ...redisTestConstants.orderPlace.orderPlace,
            order: {
              ...redisTestConstants.defaultOrder,
              orderId: undefined,
            },
          },
        },
        'Invalid OrderPlace, order id is undefined',
      ],
      [
        'missing order id',
        {
          orderPlace: {
            ...redisTestConstants.orderPlace.orderPlace,
            order: {
              ...redisTestConstants.defaultOrder,
              orderId: {
                ...redisTestConstants.defaultOrderId,
                subaccountId: undefined,
              },
            },
          },
        },
        'Invalid OrderPlace, subaccount id is undefined',
      ],
      [
        'unspecified placement status',
        {
          ...redisTestConstants.orderPlace,
          orderPlace: {
            ...redisTestConstants.orderPlace.orderPlace,
            placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_UNSPECIFIED,
          },
        },
        'Invalid OrderPlace, placement status is UNSPECIFIED',
      ],
    ])('logs error and does not update caches on invalid order place off-chain update: %s', async (
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
        at: 'OrderPlaceHandler#logAndThrowParseMessageError',
        message: errorMsg,
      }));
    });

    it('logs error and does not update caches if order has invalid clobPairId', async () => {
      const update: OffChainUpdateV1 = {
        orderPlace: {
          ...redisTestConstants.orderPlace.orderPlace,
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
        at: 'OrderPlaceHandler#logAndThrowParseMessageError',
        message: 'Order in OrderPlace has invalid clobPairId',
      }));
    });

    it('logs error if replaced order had total filled quantums > quantums', async () => {
      synchronizeWrapBackgroundTask(wrapBackgroundTask);
      const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
      // Handle the order place event for the initial order
      await handleInitialOrderPlace(redisTestConstants.orderPlace);
      expectWebsocketMessagesSent(
        producerSendSpy,
        redisTestConstants.defaultRedisOrder,
        dbDefaultOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.BEST_EFFORT_OPENED,
        true,
      );
      // clear mocks
      jest.clearAllMocks();

      // This should never happen, testing that an error is logged in this case
      // Update the order to set it to be resting on the book, with total filled greater than
      // quantums of order
      await updateOrder({
        updatedOrderId: redisTestConstants.defaultOrderId,
        newTotalFilledQuantums: Number(
          redisTestConstants.defaultOrder.quantums.multiply(
            Long.fromValue(2),
          ),
        ),
        client,
      });
      // Handle the order place off-chain update with the replacement order
      await onMessage(replacementMessage);

      // Order should still have been replaced
      await checkOrderPlace(
        redisTestConstants.defaultOrderUuid,
        redisTestConstants.defaultSubaccountUuid,
        replacedOrder,
      );
      // Order book price levels should not have been updated
      expect(OrderbookLevelsCache.updatePriceLevel).not.toHaveBeenCalled();
      expectWebsocketMessagesSent(
        producerSendSpy,
        replacedOrder,
        dbDefaultOrder,
        testConstants.defaultPerpetualMarket,
        APIOrderStatusEnum.BEST_EFFORT_OPENED,
        true,
      );
    });

    describe('stateful order update cache', () => {
      it('looks up stateful order in db and send cached order update', async () => {
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();
        const fakeSystemTime = 123456789;
        jest.setSystemTime(fakeSystemTime);
        await StatefulOrderUpdatesCache.addStatefulOrderUpdate(
          redisTestConstants.defaultOrderUuidGoodTilBlockTime,
          {
            ...redisTestConstants.orderUpdate.orderUpdate,
            orderId: redisTestConstants.defaultOrderIdGoodTilBlockTime,
          },
          Date.now(),
          redisClient,
        );

        await handleInitialOrderPlace({
          ...redisTestConstants.orderPlace,
          orderPlace: {
            order: redisTestConstants.defaultOrderGoodTilBlockTime,
            placementStatus:
            OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
          },
        });

        expect(OrderTable.findById).toHaveBeenCalledWith(
          redisTestConstants.defaultOrderUuidGoodTilBlockTime,
        );
        expect(StatefulOrderUpdatesCache.removeStatefulOrderUpdate).toHaveBeenCalledWith(
          redisTestConstants.defaultOrderUuidGoodTilBlockTime,
          fakeSystemTime,
          client,
        );
        expectWebsocketMessagesSent(
          producerSendSpy,
          redisTestConstants.defaultRedisOrderGoodTilBlockTime,
          dbOrderGoodTilBlockTime,
          testConstants.defaultPerpetualMarket,
          APIOrderStatusEnum.OPEN,
          true,
          undefined,
          {
            key: getOrderIdHash(redisTestConstants.defaultOrderIdGoodTilBlockTime),
            offchainUpdate: {
              orderUpdate: {
                orderId: redisTestConstants.defaultOrderIdGoodTilBlockTime,
                totalFilledQuantums: redisTestConstants.orderUpdate.orderUpdate.totalFilledQuantums,
              },
            },
          },
        );
      });

      it('does not look up vault order in db but sends cached order update', async () => {
        const producerSendSpy: jest.SpyInstance = jest.spyOn(producer, 'send').mockReturnThis();

        await handleInitialOrderPlace({
          ...redisTestConstants.orderPlace,
          orderPlace: {
            order: redisTestConstants.defaultOrderVault,
            placementStatus:
            OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
          },
        });

        expect(OrderTable.findById).not.toHaveBeenCalled();
        expect(StatefulOrderUpdatesCache.removeStatefulOrderUpdate).toHaveBeenCalled();
        expectWebsocketMessagesSent(
          producerSendSpy,
          redisTestConstants.defaultRedisOrderVault,
          undefined,
          testConstants.defaultPerpetualMarket,
          APIOrderStatusEnum.OPEN,
          true,
          undefined,
          undefined,
        );
      });
    });
  });
});

async function checkOrderPlace(
  placedOrderId: string,
  placedSubaccountId: string,
  expectedOrder: RedisOrder,
): Promise<void> {
  const redisOrder: RedisOrder | null = await OrdersCache.getOrder(placedOrderId, client);
  const orderIdsForSubaccount: string[] = await SubaccountOrderIdsCache.getOrderIdsForSubaccount(
    placedSubaccountId,
    client,
  );

  expect(redisOrder).toEqual(expectedOrder);
  expect(orderIdsForSubaccount).toEqual([placedOrderId]);
}

function expectStats(orderWasReplaced: boolean = false): void {
  expect(stats.timing).toHaveBeenCalledWith(
    `vulcan.${STATS_FUNCTION_NAME}.timing`,
    expect.any(Number),
    { className: 'OrderPlaceHandler', fnName: 'place_order_cache_update', instance: '' },
  );

  if (orderWasReplaced) {
    expect(stats.increment).toHaveBeenCalledWith('vulcan.place_order_handler.replaced_order', 1, { instance: '' });
  } else {
    expect(stats.increment).not.toHaveBeenCalledWith(
      'vulcan.place_order_handler.replaced_order',
      expect.any(Number),
      { instance: '' },
    );
  }
}

function expectWebsocketMessagesSent(
  producerSendSpy: jest.SpyInstance,
  redisOrder: RedisOrder,
  dbOrder: OrderFromDatabase | undefined,
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
          ...(isStateful && dbOrder && { createdAtHeight: dbOrder.createdAtHeight }),
          ...(isStateful && dbOrder && { updatedAt: dbOrder.updatedAt }),
          ...(isStateful && dbOrder && { updatedAtHeight: dbOrder.updatedAtHeight }),
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
