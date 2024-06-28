import {
  dbHelpers,
  OrderFromDatabase,
  OrderStatus,
  OrderTable,
  perpetualMarketRefresher,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  OffChainUpdateV1,
  OrderRemovalReason,
  OrderRemoveV1_OrderRemovalStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../../src/lib/types';
import {
  defaultDateTime,
  defaultHeight,
  defaultOrderId,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
} from '../../helpers/constants';
import { createKafkaMessageFromStatefulOrderEvent } from '../../helpers/kafka-helpers';
import { updateBlockCache } from '../../../src/caches/block-cache';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectVulcanKafkaMessage,
} from '../../helpers/indexer-proto-helpers';
import { StatefulOrderRemovalHandler } from '../../../src/handlers/stateful-order/stateful-order-removal-handler';
import { STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE } from '../../../src/constants';
import { producer } from '@dydxprotocol-indexer/kafka';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';

describe('statefulOrderRemovalHandler', () => {
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
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  const reason: OrderRemovalReason = OrderRemovalReason.ORDER_REMOVAL_REASON_REPLACED;
  const defaultStatefulOrderEvent: StatefulOrderEventV1 = {
    orderRemoval: {
      removedOrderId: defaultOrderId,
      reason,
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

      const handler: StatefulOrderRemovalHandler = new StatefulOrderRemovalHandler(
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

  it('successfully cancels and removes order', async () => {
    await OrderTable.create({
      ...testConstants.defaultOrder,
      clientId: '0',
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulOrderEvent,
    );

    await onMessage(kafkaMessage);
    const order: OrderFromDatabase | undefined = await OrderTable.findById(orderId);
    expect(order).toBeDefined();
    expect(order).toEqual(expect.objectContaining({
      status: OrderStatus.CANCELED,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    }));

    const expectedOffchainUpdate: OffChainUpdateV1 = {
      orderRemove: {
        removedOrderId: defaultOrderId,
        reason,
        removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED,
      },
    };
    expectVulcanKafkaMessage({
      producerSendMock,
      orderId: defaultOrderId,
      offchainUpdate: expectedOffchainUpdate,
      headers: { message_received_timestamp: kafkaMessage.timestamp, event_type: 'StatefulOrderRemoval' },
    });
  });

  it('throws error when attempting to cancel an order that does not exist', async () => {
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulOrderEvent,
    );

    await expect(onMessage(kafkaMessage)).rejects.toThrowError(
      `Unable to update order status with orderId: ${orderId}`,
    );
  });
});
