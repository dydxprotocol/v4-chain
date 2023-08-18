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
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  OffChainUpdateV1,
  IndexerOrder,
  OrderPlaceV1_OrderPlacementStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../../src/lib/types';
import {
  defaultMakerOrder,
  defaultOrderId,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
} from '../../helpers/constants';
import { createKafkaMessageFromStatefulOrderEvent } from '../../helpers/kafka-helpers';
import { updateBlockCache } from '../../../src/caches/block-cache';
import {
  binaryToBase64String,
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectVulcanKafkaMessage,
} from '../../helpers/indexer-proto-helpers';
import { StatefulOrderPlacementHandler } from '../../../src/handlers/stateful-order/stateful-order-placement-handler';
import { getPrice, getSize } from '../../../src/lib/helper';
import { stats, STATS_FUNCTION_NAME } from '@dydxprotocol-indexer/base';
import { STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE } from '../../../src/constants';
import { producer } from '@dydxprotocol-indexer/kafka';

describe('statefulOrderPlacementHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
    jest.spyOn(stats, 'gauge');
  });

  beforeEach(async () => {
    await testMocks.seedData();
    updateBlockCache(defaultPreviousHeight);
    await perpetualMarketRefresher.updatePerpetualMarkets();
    producerSendMock = jest.spyOn(producer, 'send');
  });

  afterEach(async () => {
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
    goodTilBlock: undefined,
    goodTilBlockTime,
  };
  const defaultStatefulOrderEvent: StatefulOrderEventV1 = {
    orderPlace: {
      order: defaultOrder,
    },
  };
  const orderId: string = OrderTable.orderIdToUuid(defaultOrderId);
  let producerSendMock: jest.SpyInstance;

  describe('getParallelizationIds', () => {
    it('returns the correct parallelization ids', () => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.STATEFUL_ORDER,
        binaryToBase64String(
          StatefulOrderEventV1.encode(defaultStatefulOrderEvent).finish(),
        ),
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
        indexerTendermintEvent,
        0,
        defaultStatefulOrderEvent,
      );

      const orderUuid: string = OrderTable.orderIdToUuid(defaultOrder.orderId!);
      expect(handler.getParallelizationIds()).toEqual([
        `${handler.eventType}_${orderId}`,
        `${STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE}_${orderUuid}`,
      ]);
    });
  });

  it('successfully places order', async () => {
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulOrderEvent,
    );

    await onMessage(kafkaMessage);
    const order: OrderFromDatabase | undefined = await OrderTable.findById(orderId);
    expect(order).toEqual({
      id: orderId,
      subaccountId: SubaccountTable.subaccountIdToUuid(defaultOrder.orderId!.subaccountId!),
      clientId: defaultOrder.orderId!.clientId.toString(),
      clobPairId: defaultOrder.orderId!.clobPairId.toString(),
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
    });
    expectTimingStats();

    const expectedOffchainUpdate: OffChainUpdateV1 = {
      orderPlace: {
        order: defaultOrder,
        placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
      },
    };
    expectVulcanKafkaMessage({
      producerSendMock,
      orderId: defaultOrder.orderId!,
      offchainUpdate: expectedOffchainUpdate,
    });
  });

  it('successfully upserts order', async () => {
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
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulOrderEvent,
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
    });
    expectTimingStats();
    // TODO[IND-20]: Add tests for vulcan messages
  });
});

function expectTimingStats() {
  expectTimingStat('upsert_order');
}

function expectTimingStat(fnName: string) {
  expect(stats.timing).toHaveBeenCalledWith(
    `ender.${STATS_FUNCTION_NAME}.timing`,
    expect.any(Number),
    { className: 'StatefulOrderPlacementHandler', eventType: 'StatefulOrderEvent', fnName },
  );
}
