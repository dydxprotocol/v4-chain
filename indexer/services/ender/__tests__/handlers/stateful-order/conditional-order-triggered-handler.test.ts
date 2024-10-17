import {
  dbHelpers,
  OrderFromDatabase,
  OrderStatus,
  OrderTable,
  orderTranslations,
  perpetualMarketRefresher,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import {
  IndexerOrder,
  IndexerOrderId,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  OffChainUpdateV1,
  OrderPlaceV1_OrderPlacementStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../../src/lib/types';
import {
  defaultDateTime,
  defaultHeight,
  defaultOrderId, defaultPreviousHeight, defaultTime, defaultTxHash,
} from '../../helpers/constants';
import { createKafkaMessageFromStatefulOrderEvent } from '../../helpers/kafka-helpers';
import { updateBlockCache } from '../../../src/caches/block-cache';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectVulcanKafkaMessage,
} from '../../helpers/indexer-proto-helpers';
import { STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE } from '../../../src/constants';
import { producer } from '@dydxprotocol-indexer/kafka';
import { ORDER_FLAG_CONDITIONAL } from '@dydxprotocol-indexer/v4-proto-parser';
import { ConditionalOrderTriggeredHandler } from '../../../src/handlers/stateful-order/conditional-order-triggered-handler';
import { defaultPerpetualMarket } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';
import config from '../../../src/config';

describe('conditionalOrderTriggeredHandler', () => {
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

  const conditionalOrderId: IndexerOrderId = {
    ...defaultOrderId,
    orderFlags: ORDER_FLAG_CONDITIONAL,
  };
  const defaultStatefulOrderEvent: StatefulOrderEventV1 = {
    conditionalOrderTriggered: {
      triggeredOrderId: conditionalOrderId,
    },
  };
  const orderId: string = OrderTable.orderIdToUuid(conditionalOrderId);
  let producerSendMock: jest.SpyInstance;

  describe('getParallelizationIds', () => {
    it('returns the correct parallelization ids', () => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.STATEFUL_ORDER,
        StatefulOrderEventV1.encode(defaultStatefulOrderEvent).finish(),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: ConditionalOrderTriggeredHandler = new ConditionalOrderTriggeredHandler(
        block,
        0,
        indexerTendermintEvent,
        0,
        defaultStatefulOrderEvent,
      );

      expect(handler.getParallelizationIds()).toEqual([
        `${handler.eventType}_${orderId}`,
        `${STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE}_${orderId}`,
      ]);
    });
  });

  it.each([
    ['transaction event', 0],
    ['block event', -1],
  ])('successfully triggers order and sends to vulcan (as %s)', async (
    _name: string,
    transactionIndex: number,
  ) => {
    await OrderTable.create({
      ...testConstants.defaultOrderGoodTilBlockTime,
      orderFlags: conditionalOrderId.orderFlags.toString(),
      status: OrderStatus.UNTRIGGERED,
      triggerPrice: '1000',
      clientId: '0',
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulOrderEvent,
      transactionIndex,
    );

    await onMessage(kafkaMessage);
    const order: OrderFromDatabase | undefined = await OrderTable.findById(orderId);
    const indexerOrder: IndexerOrder = await orderTranslations.convertToIndexerOrder(
      order!,
      defaultPerpetualMarket,
    );

    expect(order).toBeDefined();
    expect(order).toEqual(expect.objectContaining({
      status: OrderStatus.OPEN,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    }));

    const expectedOffchainUpdate: OffChainUpdateV1 = {
      orderPlace: {
        order: indexerOrder,
        placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
      },
    };
    expectVulcanKafkaMessage({
      producerSendMock,
      orderId: conditionalOrderId,
      offchainUpdate: expectedOffchainUpdate,
      headers: { message_received_timestamp: kafkaMessage.timestamp, event_type: 'ConditionalOrderTriggered' },
    });
  });

  it('throws error when attempting to trigger an order that does not exist', async () => {
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulOrderEvent,
    );

    await expect(onMessage(kafkaMessage)).rejects.toThrowError(
      `Unable to update order status with orderId: ${orderId}`,
    );
  });

  it.each([
    ['transaction event', 0],
    ['block event', -1],
  ])('successfully skips order trigger event (as %s)', async (
    _name: string,
    transactionIndex: number,
  ) => {
    config.SKIP_STATEFUL_ORDER_UUIDS = OrderTable.uuid(
      testConstants.defaultOrderGoodTilBlockTime.subaccountId,
      '0',
      testConstants.defaultOrderGoodTilBlockTime.clobPairId,
      testConstants.defaultOrderGoodTilBlockTime.orderFlags,
    );
    await OrderTable.create({
      ...testConstants.defaultOrderGoodTilBlockTime,
      orderFlags: conditionalOrderId.orderFlags.toString(),
      status: OrderStatus.UNTRIGGERED,
      triggerPrice: '1000',
      clientId: '0',
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulOrderEvent,
      transactionIndex,
    );

    await onMessage(kafkaMessage);
    const order: OrderFromDatabase | undefined = await OrderTable.findById(orderId);

    expect(order).toBeDefined();
    expect(order).toEqual(expect.objectContaining({
      status: OrderStatus.OPEN,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    }));
  });
});
