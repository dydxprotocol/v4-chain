import { producer } from '@dydxprotocol-indexer/kafka';
import {
  dbHelpers,
  OrderFromDatabase,
  OrderSide,
  OrderStatus,
  OrderTable,
  OrderType,
  perpetualMarketRefresher,
  protocolTranslations,
  SubaccountTable,
  testConstants,
  testMocks,
  TimeInForce,
} from '@dydxprotocol-indexer/postgres';
import { ORDER_FLAG_LONG_TERM, ORDER_FLAG_TWAP } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  OffChainUpdateV1,
  OrderPlaceV1_OrderPlacementStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import { updateBlockCache } from '../../../src/caches/block-cache';
import config from '../../../src/config';
import { STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE } from '../../../src/constants';
import { StatefulOrderPlacementHandler } from '../../../src/handlers/stateful-order/stateful-order-placement-handler';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';
import { getPrice, getSize } from '../../../src/lib/helper';
import { onMessage } from '../../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../../src/lib/types';
import {
  defaultDateTime,
  defaultHeight,
  defaultMakerOrder,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
  defaultVaultOrder,
  defaultVaultOrderPlacementEvent,
} from '../../helpers/constants';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectVulcanKafkaMessage,
} from '../../helpers/indexer-proto-helpers';
import { createKafkaMessageFromStatefulOrderEvent } from '../../helpers/kafka-helpers';

describe('statefulOrderPlacementHandler', () => {
  const prevSkippedOrderUUIDs: string = config.SKIP_STATEFUL_ORDER_UUIDS;

  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    updateBlockCache(defaultPreviousHeight);
    await perpetualMarketRefresher.updatePerpetualMarkets();
    producerSendMock = jest.spyOn(producer, 'send');
  });

  afterEach(async () => {
    config.SKIP_STATEFUL_ORDER_UUIDS = prevSkippedOrderUUIDs;
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  const goodTilBlockTime: number = 123;
  const defaultOrder: IndexerOrder = {
    ...defaultMakerOrder,
    orderId: {
      ...defaultMakerOrder.orderId!,
      orderFlags: ORDER_FLAG_LONG_TERM,
    },
    goodTilBlock: undefined,
    goodTilBlockTime,
    builderCodeParams: {
      builderAddress: 'dydx123',
      feePpm: 1000,
    },
  };
  const defaultTwapOrder: IndexerOrder = {
    ...defaultMakerOrder,
    orderId: {
      ...defaultMakerOrder.orderId!,
      orderFlags: ORDER_FLAG_TWAP,
    },
    goodTilBlock: undefined,
    goodTilBlockTime,
    twapParameters: {
      duration: 300,
      interval: 30,
      priceTolerance: 0,
    },
  };
  const defaultStatefulOrderLongTermEvent: StatefulOrderEventV1 = {
    longTermOrderPlacement: {
      order: defaultOrder,
    },
  };
  // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
  const defaultStatefulOrderEvent: StatefulOrderEventV1 = {
    orderPlace: {
      order: defaultOrder,
    },
  };

  const defaultStatefulOrderEventWithORRS: StatefulOrderEventV1 = {
    orderPlace: {
      order: {
        ...defaultOrder,
        orderRouterAddress: testConstants.defaultAddress,
      },
    },
  };

  const defaultStatefulTwapOrderEvent: StatefulOrderEventV1 = {
    twapOrderPlacement: {
      order: defaultTwapOrder,
    },
  };
  const orderId: string = OrderTable.orderIdToUuid(defaultOrder.orderId!);
  const twapOrderId: string = OrderTable.orderIdToUuid(defaultTwapOrder.orderId!);
  let producerSendMock: jest.SpyInstance;

  describe('getParallelizationIds', () => {
    it.each([
      // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
      ['stateful order placement', defaultStatefulOrderEvent],
      ['stateful twap order placement', defaultStatefulTwapOrderEvent],
      ['stateful long term order placement', defaultStatefulOrderLongTermEvent],
    ])('returns the correct parallelization ids for %s', (
      _name: string,
      statefulOrderEvent: StatefulOrderEventV1,
    ) => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.STATEFUL_ORDER,
        StatefulOrderEventV1.encode(statefulOrderEvent).finish(),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: StatefulOrderPlacementHandler = new StatefulOrderPlacementHandler(
        block,
        0,
        indexerTendermintEvent,
        0,
        statefulOrderEvent,
      );
      const order = statefulOrderEvent.twapOrderPlacement ? defaultTwapOrder : defaultOrder;

      const orderUuid: string = OrderTable.orderIdToUuid(order.orderId!);
      expect(handler.getParallelizationIds()).toEqual([
        `${handler.eventType}_${orderUuid}`,
        `${STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE}_${orderUuid}`,
      ]);
    });
  });

  it.each([
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    ['stateful order placement as txn event', defaultStatefulOrderEvent, 0],
    ['stateful long term order placement as txn event', defaultStatefulOrderLongTermEvent, 0],
    ['stateful order placement as block event', defaultStatefulOrderEvent, -1],
    ['stateful long term order placement as block event', defaultStatefulOrderLongTermEvent, -1],
    ['stateful twap order placement as txn event', defaultStatefulTwapOrderEvent, 0],
    ['stateful twap order placement as block event', defaultStatefulTwapOrderEvent, -1],
  ])('successfully places order with %s', async (
    _name: string,
    statefulOrderEvent: StatefulOrderEventV1,
    transactionIndex: number,
  ) => {
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      statefulOrderEvent,
      transactionIndex,
    );

    await onMessage(kafkaMessage);
    const orderId_ = statefulOrderEvent.twapOrderPlacement ? twapOrderId : orderId;
    const order: OrderFromDatabase | undefined = await OrderTable.findById(orderId_);

    const testOrder = statefulOrderEvent.twapOrderPlacement ? defaultTwapOrder : defaultOrder;
    expect(order).toEqual({
      id: orderId_,
      subaccountId: SubaccountTable.subaccountIdToUuid(testOrder.orderId!.subaccountId!),
      clientId: testOrder.orderId!.clientId.toString(),
      clobPairId: testOrder.orderId!.clobPairId.toString(),
      side: OrderSide.BUY,
      size: getSize(testOrder, testConstants.defaultPerpetualMarket),
      totalFilled: '0',
      price: getPrice(testOrder, testConstants.defaultPerpetualMarket),
      type: testOrder.twapParameters ? OrderType.TWAP : OrderType.LIMIT,
      status: OrderStatus.OPEN,
      timeInForce: protocolTranslations.protocolOrderTIFToTIF(testOrder.timeInForce),
      reduceOnly: testOrder.reduceOnly,
      orderFlags: testOrder.orderId!.orderFlags.toString(),
      goodTilBlock: null,
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(testOrder),
      createdAtHeight: '3',
      clientMetadata: '0',
      triggerPrice: null,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
      builderAddress: testOrder.builderCodeParams?.builderAddress ?? null,
      feePpm: testOrder.builderCodeParams?.feePpm.toString() ?? null,
      orderRouterAddress: testOrder.orderRouterAddress,
      duration: testOrder.twapParameters?.duration.toString() ?? null,
      interval: testOrder.twapParameters?.interval.toString() ?? null,
      priceTolerance: testOrder.twapParameters?.priceTolerance.toString() ?? null,
    });

    const expectedOffchainUpdate: OffChainUpdateV1 = {
      orderPlace: {
        order: testOrder,
        placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
      },
    };
    expectVulcanKafkaMessage({
      producerSendMock,
      orderId: testOrder.orderId!,
      offchainUpdate: expectedOffchainUpdate,
      headers: { message_received_timestamp: kafkaMessage.timestamp, event_type: 'StatefulOrderPlacement' },
    });
  });

  it.each([
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    ['stateful order placement', defaultStatefulOrderEvent],
    ['stateful twap order placement', defaultStatefulTwapOrderEvent],
    ['stateful long term order placement', defaultStatefulOrderLongTermEvent],
  ])('successfully upserts order with %s', async (
    _name: string,
    statefulOrderEvent: StatefulOrderEventV1,
  ) => {
    const testOrder = statefulOrderEvent.twapOrderPlacement ? defaultTwapOrder : defaultOrder;
    const orderId_ = statefulOrderEvent.twapOrderPlacement ? twapOrderId : orderId;
    const subaccountId: string = SubaccountTable.subaccountIdToUuid(
      testOrder.orderId!.subaccountId!,
    );
    const clientId: string = testOrder.orderId!.clientId.toString();
    const clobPairId: string = testOrder.orderId!.clobPairId.toString();
    await OrderTable.create({
      subaccountId,
      clientId,
      clobPairId,
      side: OrderSide.SELL,
      size: '100',
      totalFilled: '0',
      price: '200',
      type: OrderType.LIMIT,
      status: OrderStatus.CANCELED,
      timeInForce: TimeInForce.GTT,
      reduceOnly: true,
      orderFlags: '0',
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(testOrder),
      createdAtHeight: '1',
      clientMetadata: '0',
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: '0',
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      statefulOrderEvent,
    );

    await onMessage(kafkaMessage);
    const order: OrderFromDatabase | undefined = await OrderTable.findById(orderId_);
    expect(order).toEqual({
      id: orderId_,
      subaccountId,
      clientId,
      clobPairId,
      side: OrderSide.BUY,
      size: getSize(testOrder, testConstants.defaultPerpetualMarket),
      totalFilled: '0',
      price: getPrice(testOrder, testConstants.defaultPerpetualMarket),
      type: testOrder.twapParameters ? OrderType.TWAP : OrderType.LIMIT,
      status: OrderStatus.OPEN,
      timeInForce: protocolTranslations.protocolOrderTIFToTIF(testOrder.timeInForce),
      reduceOnly: testOrder.reduceOnly,
      orderFlags: testOrder.orderId!.orderFlags.toString(),
      goodTilBlock: null,
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(testOrder),
      createdAtHeight: '3',
      clientMetadata: '0',
      triggerPrice: null,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
      builderAddress: testOrder.builderCodeParams?.builderAddress ?? null,
      feePpm: testOrder.builderCodeParams?.feePpm.toString() ?? null,
      orderRouterAddress: '',
      duration: testOrder.twapParameters?.duration.toString() ?? null,
      interval: testOrder.twapParameters?.interval.toString() ?? null,
      priceTolerance: testOrder.twapParameters?.priceTolerance.toString() ?? null,
    });
    // TODO[IND-20]: Add tests for vulcan messages
  });

  it.each([
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    ['stateful order placement with ORRS', defaultStatefulOrderEventWithORRS],
  ])('successfully upserts order with %s', async (
    _name: string,
    statefulOrderEvent: StatefulOrderEventV1,
  ) => {
    const subaccountId: string = SubaccountTable.subaccountIdToUuid(
      defaultOrder.orderId!.subaccountId!,
    );
    const clientId: string = defaultOrder.orderId!.clientId.toString();
    const clobPairId: string = defaultOrder.orderId!.clobPairId.toString();
    await OrderTable.create({
      subaccountId,
      clientId,
      clobPairId,
      side: OrderSide.SELL,
      size: '100',
      totalFilled: '0',
      price: '200',
      type: OrderType.LIMIT,
      status: OrderStatus.CANCELED,
      timeInForce: TimeInForce.GTT,
      reduceOnly: true,
      orderFlags: '0',
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(defaultOrder),
      createdAtHeight: '1',
      clientMetadata: '0',
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: '0',
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      statefulOrderEvent,
    );

    await onMessage(kafkaMessage);
    const order: OrderFromDatabase | undefined = await OrderTable.findById(orderId);
    expect(order).toEqual({
      id: orderId,
      subaccountId,
      clientId,
      clobPairId,
      side: OrderSide.BUY,
      size: getSize(defaultOrder, testConstants.defaultPerpetualMarket),
      totalFilled: '0',
      price: getPrice(defaultOrder, testConstants.defaultPerpetualMarket),
      type: OrderType.LIMIT, // TODO: Add additional order types once we support
      status: OrderStatus.OPEN,
      timeInForce: protocolTranslations.protocolOrderTIFToTIF(defaultOrder.timeInForce),
      reduceOnly: defaultOrder.reduceOnly,
      orderFlags: defaultOrder.orderId!.orderFlags.toString(),
      goodTilBlock: null,
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(defaultOrder),
      createdAtHeight: '3',
      clientMetadata: '0',
      triggerPrice: null,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
      builderAddress: defaultOrder.builderCodeParams?.builderAddress,
      feePpm: defaultOrder.builderCodeParams?.feePpm.toString(),
      orderRouterAddress: testConstants.defaultAddress,
      duration: defaultOrder.twapParameters?.duration.toString() ?? null,
      interval: defaultOrder.twapParameters?.interval.toString() ?? null,
      priceTolerance: defaultOrder.twapParameters?.priceTolerance.toString() ?? null,
    });
    // TODO[IND-20]: Add tests for vulcan messages
  });

  it.each([
    // TODO(IND-334): Remove after deprecating StatefulOrderPlacementEvent
    ['stateful order placement as txn event', defaultStatefulOrderEvent, 0],
    ['stateful long term order placement as txn event', defaultStatefulOrderLongTermEvent, 0],
    ['stateful order placement as block event', defaultStatefulOrderEvent, -1],
    ['stateful long term order placement as block event', defaultStatefulOrderLongTermEvent, -1],
    ['stateful twap order placement as txn event', defaultStatefulTwapOrderEvent, 0],
    ['stateful twap order placement as block event', defaultStatefulTwapOrderEvent, -1],
  ])('successfully skips order with %s', async (
    _name: string,
    statefulOrderEvent: StatefulOrderEventV1,
    transactionIndex: number,
  ) => {
    const testOrder = statefulOrderEvent.twapOrderPlacement ? defaultTwapOrder : defaultOrder;
    const orderId_ = statefulOrderEvent.twapOrderPlacement ? twapOrderId : orderId;
    config.SKIP_STATEFUL_ORDER_UUIDS = OrderTable.orderIdToUuid(testOrder.orderId!);
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      statefulOrderEvent,
      transactionIndex,
    );

    await onMessage(kafkaMessage);
    const order: OrderFromDatabase | undefined = await OrderTable.findById(orderId_);
    expect(order).toBeUndefined();
  });

  it.each([
    ['txn event', defaultVaultOrderPlacementEvent, 0],
    ['block event', defaultVaultOrderPlacementEvent, -1],
  ])('successfully skips vault order placements (as %s)', async (
    _name: string,
    statefulOrderEvent: StatefulOrderEventV1,
    transactionIndex: number,
  ) => {
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      statefulOrderEvent,
      transactionIndex,
    );

    await onMessage(kafkaMessage);
    const order: OrderFromDatabase | undefined = await OrderTable.findById(orderId);
    expect(order).toBeUndefined();
    const expectedOffchainUpdate: OffChainUpdateV1 = {
      orderPlace: {
        order: defaultVaultOrder,
        placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
      },
    };
    expectVulcanKafkaMessage({
      producerSendMock,
      orderId: defaultVaultOrder.orderId!,
      offchainUpdate: expectedOffchainUpdate,
      headers: { message_received_timestamp: kafkaMessage.timestamp, event_type: 'StatefulOrderPlacement' },
    });
  });
});
