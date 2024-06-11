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
  expectOrderSubaccountKafkaMessage,
  expectVulcanKafkaMessage,
} from '../../helpers/indexer-proto-helpers';
import { getPrice, getSize } from '../../../src/lib/helper';
import { producer } from '@dydxprotocol-indexer/kafka';
import { ORDER_FLAG_LONG_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';
import config from '../../../src/config';

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
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
    config.SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_STATEFUL_ORDERS = false;
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

  const defaultStatefulOrderReplacementEvent: StatefulOrderEventV1 = {
    orderReplacement: {
      oldOrderId: defaultOldOrder.orderId!,
      order: defaultNewOrder,
    },
  };

  const oldOrderUuid: string = OrderTable.orderIdToUuid(defaultOldOrder.orderId!);
  const newOrderUuid: string = OrderTable.orderIdToUuid(defaultNewOrder.orderId!);
  let producerSendMock: jest.SpyInstance;

  it.each([
    ['stateful order placement as txn event', defaultStatefulOrderReplacementEvent, false, 0],
    ['stateful order placement as txn event', defaultStatefulOrderReplacementEvent, true, 0],
    ['stateful order placement as block event', defaultStatefulOrderReplacementEvent, false, -1],
    ['stateful order placement as block event', defaultStatefulOrderReplacementEvent, true, -1],
  ])('successfully places order with %s (emit subaccount websocket msg: %s)', async (
    _name: string,
    statefulOrderEvent: StatefulOrderEventV1,
    emitSubaccountMessage: boolean,
    transactionIndex: number,
  ) => {
    await OrderTable.create({
      ...testConstants.defaultOrder,
      clientId: '0',
      orderFlags: ORDER_FLAG_LONG_TERM.toString(),
    });
    config.SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_STATEFUL_ORDERS = emitSubaccountMessage;
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
    if (emitSubaccountMessage) {
      expectOrderSubaccountKafkaMessage(
        producerSendMock,
        defaultNewOrder.orderId!.subaccountId!,
        newOrder!,
        defaultHeight.toString(),
        transactionIndex,
      );
    }
  });

  it.each([
    ['stateful order replacement', defaultStatefulOrderReplacementEvent],
  ])('successfully upserts order with %s', async (
    _name: string,
    statefulOrderReplacementEvent: StatefulOrderEventV1,
  ) => {
    const subaccountId: string = SubaccountTable.subaccountIdToUuid(
      defaultNewOrder.orderId!.subaccountId!,
    );
    const clientId: string = defaultNewOrder.orderId!.clientId.toString();
    const clobPairId: string = defaultNewOrder.orderId!.clobPairId.toString();
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
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(defaultNewOrder),
      createdAtHeight: '1',
      clientMetadata: '0',
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: '0',
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      statefulOrderReplacementEvent,
    );

    await onMessage(kafkaMessage);
    const order: OrderFromDatabase | undefined = await OrderTable.findById(newOrderUuid);
    expect(order).toEqual({
      id: newOrderUuid,
      subaccountId,
      clientId,
      clobPairId,
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
