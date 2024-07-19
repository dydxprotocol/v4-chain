import { KafkaTopics, producer } from '@dydxprotocol-indexer/kafka';
import {
  dbHelpers,
  OrderFromDatabase,
  OrderTable,
  orderTranslations,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  SubaccountFromDatabase,
  SubaccountTable,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';

import { getOrderIdHash, ORDER_FLAG_LONG_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import { sendStatefulOrderMessages } from '../src/vulcan-helpers';
import {
  IndexerOrder,
  IndexerOrderId,
  OffChainUpdateV1,
  OrderPlaceV1_OrderPlacementStatus,
} from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';

describe('vulcan-helpers', () => {
  let producerSendMock: jest.SpyInstance;

  beforeAll(async () => {
    producerSendMock = jest.spyOn(producer, 'send');
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterEach(async () => {
    jest.clearAllMocks();
    await dbHelpers.clearData();
  });

  afterAll(async () => {
    jest.resetAllMocks();
    await dbHelpers.teardown();
  });

  it('sendStatefulOrderMessages without fills', async () => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [_unusedOrder, longTermOrder]: [OrderFromDatabase, OrderFromDatabase] = await
    Promise.all([
      OrderTable.create(testConstants.defaultOrder),
      OrderTable.create({
        ...testConstants.defaultOrder,
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
      }),
    ]);
    await sendStatefulOrderMessages();

    expect(producerSendMock.mock.calls).toHaveLength(1);

    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(
      testConstants.defaultOrder.subaccountId,
    );

    const market: PerpetualMarketFromDatabase = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(testConstants.defaultOrder.clobPairId)!;
    const order: IndexerOrder = await
    orderTranslations.convertToIndexerOrder(longTermOrder, market);
    expect(producerSendMock.mock.lastCall).toEqual([
      {
        topic: KafkaTopics.TO_VULCAN,
        messages: [
          {
            key: getOrderIdHash(
              {
                subaccountId: {
                  owner: subaccount?.address!,
                  number: subaccount?.subaccountNumber!,
                },
                clientId: Number(testConstants.defaultOrder.clientId),
                clobPairId: Number(testConstants.defaultOrder.clobPairId),
                orderFlags: ORDER_FLAG_LONG_TERM,
              },
            ),
            value: Buffer.from(
              Uint8Array.from(
                OffChainUpdateV1.encode(
                  OffChainUpdateV1.fromPartial({
                    orderPlace: {
                      order,
                      // eslint-disable-next-line max-len
                      placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
                    },
                  }),
                ).finish(),
              ),
            ),
          },
        ],
      },
    ]);
  });

  it('sendStatefulOrderMessages with fills', async () => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [_, longTermOrder]: [OrderFromDatabase, OrderFromDatabase] = await
    Promise.all([
      OrderTable.create(testConstants.defaultOrder),
      OrderTable.create({
        ...testConstants.defaultOrderGoodTilBlockTime,
        totalFilled: '1000',
      }),
    ]);
    await sendStatefulOrderMessages();

    expect(producerSendMock.mock.calls).toHaveLength(1);

    const subaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(
      testConstants.defaultOrderGoodTilBlockTime.subaccountId,
    );

    const market: PerpetualMarketFromDatabase = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(testConstants.defaultOrderGoodTilBlockTime.clobPairId)!;
    const order: IndexerOrder = await
    orderTranslations.convertToIndexerOrder(longTermOrder, market);
    const expectedOrderId: IndexerOrderId = {
      subaccountId: {
        owner: subaccount?.address!,
        number: subaccount?.subaccountNumber!,
      },
      clientId: Number(testConstants.defaultOrderGoodTilBlockTime.clientId),
      clobPairId: Number(testConstants.defaultOrderGoodTilBlockTime.clobPairId),
      orderFlags: ORDER_FLAG_LONG_TERM,
    };
    expect(producerSendMock.mock.lastCall[0].messages).toHaveLength(2);
    expect(producerSendMock.mock.lastCall).toEqual([
      {
        topic: KafkaTopics.TO_VULCAN,
        messages: [
          {
            key: getOrderIdHash(expectedOrderId),
            value: Buffer.from(
              Uint8Array.from(
                OffChainUpdateV1.encode(
                  OffChainUpdateV1.fromPartial({
                    orderPlace: {
                      order,
                      // eslint-disable-next-line max-len
                      placementStatus: OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED,
                    },
                  }),
                ).finish(),
              ),
            ),
          },
          {
            key: getOrderIdHash(expectedOrderId),
            value: Buffer.from(
              Uint8Array.from(
                OffChainUpdateV1.encode(
                  OffChainUpdateV1.fromPartial({
                    orderUpdate: {
                      orderId: expectedOrderId,
                      // eslint-disable-next-line max-len
                      totalFilledQuantums: Long.fromValue(10_000_000_000_000, true),  // 1e3 / 1e-10 = 1e13
                    },
                  }),
                ).finish(),
              ),
            ),
          },
        ],
      },
    ]);
  });
});
