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
} from '@dydxprotocol-indexer/postgres';
import {
  OffChainUpdateV1,
  IndexerOrder,
  OrderPlaceV1_OrderPlacementStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../../src/lib/on-message';
import {
  defaultDateTime,
  defaultHeight,
  defaultMakerOrder,
  defaultOrderId2,
  defaultPreviousHeight,
} from '../../helpers/constants';
import { createKafkaMessageFromStatefulOrderEvent } from '../../helpers/kafka-helpers';
import { updateBlockCache } from '../../../src/caches/block-cache';
import {
  expectVulcanKafkaMessage,
} from '../../helpers/indexer-proto-helpers';
import { getPrice, getSize } from '../../../src/lib/helper';
import { producer } from '@dydxprotocol-indexer/kafka';
import { ORDER_FLAG_LONG_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';
import { logger } from '@dydxprotocol-indexer/base';

describe('stateful-order-replacement-handler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    updateBlockCache(defaultPreviousHeight);
    await perpetualMarketRefresher.updatePerpetualMarkets();
    producerSendMock = jest.spyOn(producer, 'send');
    jest.spyOn(logger, 'error');
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
  const defaultOldOrder: IndexerOrder = {
    ...defaultMakerOrder,
    orderId: {
      ...defaultMakerOrder.orderId!,
      orderFlags: ORDER_FLAG_LONG_TERM,
    },
    goodTilBlock: undefined,
    goodTilBlockTime,
  };
  const defaultNewOrder: IndexerOrder = {
    ...defaultMakerOrder,
    orderId: defaultOrderId2,
    quantums: defaultOldOrder.quantums.mul(2),
    goodTilBlock: undefined,
    goodTilBlockTime,
  };

  // replacing order with a different order ID
  const defaultStatefulOrderReplacementEvent: StatefulOrderEventV1 = {
    orderReplacement: {
      oldOrderId: defaultOldOrder.orderId!,
      order: defaultNewOrder,
    },
  };

  // replacing order with the same order ID
  const statefulOrderReplacementEventSameId: StatefulOrderEventV1 = {
    orderReplacement: {
      oldOrderId: defaultOldOrder.orderId!,
      order: {
        ...defaultNewOrder,
        orderId: defaultOldOrder.orderId,
      },
    },
  };

  const oldOrderUuid: string = OrderTable.orderIdToUuid(defaultOldOrder.orderId!);
  const newOrderUuid: string = OrderTable.orderIdToUuid(defaultNewOrder.orderId!);
  let producerSendMock: jest.SpyInstance;

  it.each([
    ['stateful order replacement as txn event', defaultStatefulOrderReplacementEvent, 0],
    ['stateful order replacement as txn event', defaultStatefulOrderReplacementEvent, 0],
    ['stateful order replacement as block event', defaultStatefulOrderReplacementEvent, -1],
    ['stateful order replacement as block event', defaultStatefulOrderReplacementEvent, -1],
  ])('successfully replaces order with %s', async (
    _name: string,
    statefulOrderEvent: StatefulOrderEventV1,
    transactionIndex: number,
  ) => {
    await OrderTable.create({
      ...testConstants.defaultOrder,
      clientId: '0',
      orderFlags: ORDER_FLAG_LONG_TERM.toString(),
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      statefulOrderEvent,
      transactionIndex,
    );

    await onMessage(kafkaMessage);

    const oldOrder: OrderFromDatabase | undefined = await OrderTable.findById(oldOrderUuid);
    expect(oldOrder).toBeDefined();
    expect(oldOrder).toEqual(expect.objectContaining({
      status: OrderStatus.CANCELED,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    }));

    const newOrder: OrderFromDatabase | undefined = await OrderTable.findById(newOrderUuid);
    expect(newOrder).toEqual({
      id: newOrderUuid,
      subaccountId: SubaccountTable.subaccountIdToUuid(defaultNewOrder.orderId!.subaccountId!),
      clientId: defaultNewOrder.orderId!.clientId.toString(),
      clobPairId: defaultNewOrder.orderId!.clobPairId.toString(),
      side: OrderSide.BUY,
      size: getSize(defaultNewOrder, testConstants.defaultPerpetualMarket),
      totalFilled: '0',
      price: getPrice(defaultNewOrder, testConstants.defaultPerpetualMarket),
      type: OrderType.LIMIT, // TODO: Add additional order types once we support
      status: OrderStatus.OPEN,
      timeInForce: protocolTranslations.protocolOrderTIFToTIF(defaultNewOrder.timeInForce),
      reduceOnly: defaultNewOrder.reduceOnly,
      orderFlags: defaultNewOrder.orderId!.orderFlags.toString(),
      goodTilBlock: null,
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(defaultNewOrder),
      createdAtHeight: '3',
      clientMetadata: '0',
      triggerPrice: null,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    });

    const expectedOffchainUpdate: OffChainUpdateV1 = {
      orderReplace: {
        oldOrderId: defaultOldOrder.orderId!,
        order: defaultNewOrder,
        placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
      },
    };
    expectVulcanKafkaMessage({
      producerSendMock,
      orderId: defaultNewOrder.orderId!,
      offchainUpdate: expectedOffchainUpdate,
      headers: { message_received_timestamp: kafkaMessage.timestamp, event_type: 'StatefulOrderReplacement' },
    });
  });

  it('successfully replaces order where old order ID is the same as new order ID', async () => {
    // create existing order with the same ID as the one we will cancel and place again
    await OrderTable.create({
      ...testConstants.defaultOrderGoodTilBlockTime,
      clientId: '0',
    });

    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      statefulOrderReplacementEventSameId,
    );

    await onMessage(kafkaMessage);
    const order: OrderFromDatabase | undefined = await OrderTable.findById(oldOrderUuid);
    expect(order).toEqual({
      id: oldOrderUuid,
      subaccountId: SubaccountTable.subaccountIdToUuid(defaultOldOrder.orderId!.subaccountId!),
      clientId: defaultNewOrder.orderId!.clientId.toString(),
      clobPairId: defaultNewOrder.orderId!.clobPairId.toString(),
      side: OrderSide.BUY,
      size: getSize(defaultNewOrder, testConstants.defaultPerpetualMarket),
      totalFilled: '0',
      price: getPrice(defaultNewOrder, testConstants.defaultPerpetualMarket),
      type: OrderType.LIMIT, // TODO: Add additional order types once we support
      status: OrderStatus.OPEN,
      timeInForce: protocolTranslations.protocolOrderTIFToTIF(defaultNewOrder.timeInForce),
      reduceOnly: defaultNewOrder.reduceOnly,
      orderFlags: defaultNewOrder.orderId!.orderFlags.toString(),
      goodTilBlock: null,
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(defaultNewOrder),
      createdAtHeight: '3',
      clientMetadata: '0',
      triggerPrice: null,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    });
  });

  it('logs error if old order ID does not exist in DB', async () => {
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulOrderReplacementEvent,
    );

    await onMessage(kafkaMessage);

    expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
      at: 'StatefulOrderReplacementHandler#handleOrderReplacement',
      message: 'Unable to cancel replaced order because orderId not found',
      orderId: defaultStatefulOrderReplacementEvent.orderReplacement!.oldOrderId,
    }));

    // We still expect new order to be created
    const newOrder: OrderFromDatabase | undefined = await OrderTable.findById(newOrderUuid);
    expect(newOrder).toEqual({
      id: newOrderUuid,
      subaccountId: SubaccountTable.subaccountIdToUuid(defaultNewOrder.orderId!.subaccountId!),
      clientId: defaultNewOrder.orderId!.clientId.toString(),
      clobPairId: defaultNewOrder.orderId!.clobPairId.toString(),
      side: OrderSide.BUY,
      size: getSize(defaultNewOrder, testConstants.defaultPerpetualMarket),
      totalFilled: '0',
      price: getPrice(defaultNewOrder, testConstants.defaultPerpetualMarket),
      type: OrderType.LIMIT, // TODO: Add additional order types once we support
      status: OrderStatus.OPEN,
      timeInForce: protocolTranslations.protocolOrderTIFToTIF(defaultNewOrder.timeInForce),
      reduceOnly: defaultNewOrder.reduceOnly,
      orderFlags: defaultNewOrder.orderId!.orderFlags.toString(),
      goodTilBlock: null,
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(defaultNewOrder),
      createdAtHeight: '3',
      clientMetadata: '0',
      triggerPrice: null,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    });
  });
});
