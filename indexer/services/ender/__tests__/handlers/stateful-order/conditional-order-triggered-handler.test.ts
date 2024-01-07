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
  OffChainUpdateV1,
  OrderPlaceV1_OrderPlacementStatus,
  StatefulOrderEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../../src/lib/on-message';
import {
  defaultDateTime,
  defaultHeight,
  defaultOrderId,
  defaultPreviousHeight,
} from '../../helpers/constants';
import { createKafkaMessageFromStatefulOrderEvent } from '../../helpers/kafka-helpers';
import { updateBlockCache } from '../../../src/caches/block-cache';
import {
  expectVulcanKafkaMessage,
} from '../../helpers/indexer-proto-helpers';
import { producer } from '@dydxprotocol-indexer/kafka';
import { ORDER_FLAG_CONDITIONAL } from '@dydxprotocol-indexer/v4-proto-parser';
import { defaultPerpetualMarket } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';

describe('conditionalOrderTriggeredHandler', () => {
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

  it('successfully triggers order and sends to vulcan', async () => {
    await OrderTable.create({
      ...testConstants.defaultOrderGoodTilBlockTime,
      orderFlags: conditionalOrderId.orderFlags.toString(),
      status: OrderStatus.UNTRIGGERED,
      triggerPrice: '1000',
      clientId: '0',
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromStatefulOrderEvent(
      defaultStatefulOrderEvent,
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
});
