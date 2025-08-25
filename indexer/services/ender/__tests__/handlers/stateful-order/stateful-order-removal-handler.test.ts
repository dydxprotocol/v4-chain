import { producer } from '@dydxprotocol-indexer/kafka';
import {
  dbHelpers,
  OrderCreateObject,
  OrderFromDatabase,
  OrderStatus,
  OrderTable,
  perpetualMarketRefresher,
  SubaccountTable,
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
import { updateBlockCache } from '../../../src/caches/block-cache';
import config from '../../../src/config';
import { STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE } from '../../../src/constants';
import { StatefulOrderRemovalHandler } from '../../../src/handlers/stateful-order/stateful-order-removal-handler';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';
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
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectVulcanKafkaMessage,
} from '../../helpers/indexer-proto-helpers';
import { createKafkaMessageFromStatefulOrderEvent } from '../../helpers/kafka-helpers';
import { ORDER_FLAG_TWAP, ORDER_FLAG_TWAP_SUBORDER } from '@dydxprotocol-indexer/v4-proto-parser';

describe('statefulOrderRemovalHandler', () => {
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

  const reason: OrderRemovalReason = OrderRemovalReason.ORDER_REMOVAL_REASON_REPLACED;
  const defaultStatefulOrderEvent: StatefulOrderEventV1 = {
    orderRemoval: {
      removedOrderId: defaultOrderId,
      reason,
    },
  };

  const defaultStatefulTwapOrderEvent: StatefulOrderEventV1 = {
    orderRemoval: {
      removedOrderId: {
        ...defaultOrderId,
        orderFlags: ORDER_FLAG_TWAP,
      },
      reason,
    },
  };

  const defaultStatefulTwapSuborderEvent: StatefulOrderEventV1 = {
    orderRemoval: {
      removedOrderId: {
        ...defaultOrderId,
        orderFlags: ORDER_FLAG_TWAP_SUBORDER,
      },
      reason,
    },
  };
  const defaultStatefulVaultOrderEvent: StatefulOrderEventV1 = {
    orderRemoval: {
      removedOrderId: {
        ...defaultOrderId,
        subaccountId: {
          owner: testConstants.defaultVaultAddress,
          number: 0,
        },
      },
      reason,
    },
  };
  const orderId: string = OrderTable.orderIdToUuid(defaultOrderId);
  const twapOrderId: string = OrderTable.orderIdToUuid({
    ...defaultOrderId,
    orderFlags: ORDER_FLAG_TWAP,
  });
  const vaultOrderId: string = OrderTable.orderIdToUuid(
    defaultStatefulVaultOrderEvent.orderRemoval!.removedOrderId!,
  );
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

  it.each([
    ['transaction event', 0],
    ['block event', -1],
  ])('successfully cancels and removes order (as %s)', async (
    _name: string,
    transactionIndex: number,
  ) => {
    await OrderTable.create({
      ...testConstants.defaultOrder,
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

  it.each([
    ['transaction event', 0],
    ['block event', -1],
  ])('successfully does not cancel parent twap on suborder removal (as %s)', async (
    _name: string,
    transactionIndex: number,
  ) => {
    // Create Parent TWAP Order
    await OrderTable.create({
      ...testConstants.defaultOrder,
      orderFlags: '128',
      clientId: '0',
    });

    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulTwapSuborderEvent,
      transactionIndex,
    );

    await onMessage(kafkaMessage);
    const twapOrder: OrderFromDatabase | undefined = await OrderTable.findById(twapOrderId);
    expect(twapOrder).toBeDefined();
    expect(twapOrder).toEqual(expect.objectContaining({
      status: OrderStatus.OPEN,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    }));

    const expectedOffchainUpdate: OffChainUpdateV1 = {
      orderRemove: {
        removedOrderId: defaultStatefulTwapSuborderEvent.orderRemoval!.removedOrderId!,
        reason,
        removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED,
      },
    };
    expectVulcanKafkaMessage({
      producerSendMock,
      orderId: defaultStatefulTwapSuborderEvent.orderRemoval!.removedOrderId!,
      offchainUpdate: expectedOffchainUpdate,
      headers: { message_received_timestamp: kafkaMessage.timestamp, event_type: 'StatefulOrderRemoval' },
    });
  });

  it.each([
    ['transaction event', 0],
    ['block event', -1],
  ])('successfully cancels and removes twap order (as %s)', async (
    _name: string,
    transactionIndex: number,
  ) => {
    // Create Parent TWAP Order
    await OrderTable.create({
      ...testConstants.defaultOrder,
      orderFlags: '128',
      clientId: '0',
    });

    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulTwapOrderEvent,
      transactionIndex,
    );

    await onMessage(kafkaMessage);
    const twapOrder: OrderFromDatabase | undefined = await OrderTable.findById(twapOrderId);
    expect(twapOrder).toBeDefined();
    expect(twapOrder).toEqual(expect.objectContaining({
      status: OrderStatus.CANCELED,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    }));

    const expectedOffchainUpdate: OffChainUpdateV1 = {
      orderRemove: {
        removedOrderId: defaultStatefulTwapOrderEvent.orderRemoval!.removedOrderId!,
        reason,
        removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED,
      },
    };
    expectVulcanKafkaMessage({
      producerSendMock,
      orderId: defaultStatefulTwapOrderEvent.orderRemoval!.removedOrderId!,
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

  it.each([
    ['transaction event', 0],
    ['block event', -1],
  ])('successfully skips order removal event (as %s)', async (
    _name: string,
    transactionIndex: number,
  ) => {
    config.SKIP_STATEFUL_ORDER_UUIDS = OrderTable.uuid(
      testConstants.defaultOrder.subaccountId,
      '0',
      testConstants.defaultOrder.clobPairId,
      testConstants.defaultOrder.orderFlags,
    );
    await OrderTable.create({
      ...testConstants.defaultOrder,
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
      ...testConstants.defaultOrder,
      clientId: '0',
    }));
  });

  it.each([
    ['transaction event', 0],
    ['block event', -1],
  ])('successfully skips vault order removals (as %s)', async (
    _name: string,
    transactionIndex: number,
  ) => {
    const vaultOrderCreateEvent: OrderCreateObject = {
      ...testConstants.defaultOrder,
      subaccountId: SubaccountTable.uuid(testConstants.defaultVaultAddress, 0),
      clientId: '0',
    };
    await SubaccountTable.create({
      ...testConstants.defaultSubaccount,
      address: testConstants.defaultVaultAddress,
      subaccountNumber: 0,
    });
    await OrderTable.create(vaultOrderCreateEvent);
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulVaultOrderEvent,
      transactionIndex,
    );

    await onMessage(kafkaMessage);
    const order: OrderFromDatabase | undefined = await OrderTable.findById(vaultOrderId);
    expect(order).toBeDefined();
    expect(order).toEqual(expect.objectContaining({
      ...vaultOrderCreateEvent,
      clientId: '0',
    }));
    const expectedOffchainUpdate: OffChainUpdateV1 = {
      orderRemove: {
        removedOrderId: defaultStatefulVaultOrderEvent.orderRemoval!.removedOrderId!,
        reason,
        removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED,
      },
    };
    expectVulcanKafkaMessage({
      producerSendMock,
      orderId: defaultStatefulVaultOrderEvent.orderRemoval!.removedOrderId!,
      offchainUpdate: expectedOffchainUpdate,
      headers: { message_received_timestamp: kafkaMessage.timestamp, event_type: 'StatefulOrderRemoval' },
    });
  });
});
