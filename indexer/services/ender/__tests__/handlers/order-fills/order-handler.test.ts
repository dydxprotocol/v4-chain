import { logger } from '@dydxprotocol-indexer/base';
import { producer } from '@dydxprotocol-indexer/kafka';
import {
  assetRefresher,
  CandleFromDatabase,
  CandleTable,
  dbHelpers,
  FillTable,
  FillType,
  Liquidity,
  OrderFromDatabase,
  OrderSide,
  OrderStatus,
  OrderTable,
  OrderType,
  perpetualMarketRefresher,
  PerpetualPositionCreateObject,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  PositionSide,
  protocolTranslations,
  SubaccountTable,
  TendermintEventTable,
  testConstants,
  testMocks,
  TimeInForce,
} from '@dydxprotocol-indexer/postgres';
import { CanceledOrdersCache, redis } from '@dydxprotocol-indexer/redis';
import {
  ORDER_FLAG_LONG_TERM,
  ORDER_FLAG_SHORT_TERM,
  ORDER_FLAG_TWAP,
  ORDER_FLAG_TWAP_SUBORDER,
} from '@dydxprotocol-indexer/v4-proto-parser';
import {
  IndexerOrder,
  IndexerOrder_Side,
  IndexerOrder_TimeInForce,
  IndexerOrderId,
  IndexerSubaccountId,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  OffChainUpdateV1,
  OrderFillEventV1,
  OrderRemovalReason,
  OrderRemoveV1_OrderRemovalStatus,
  Timestamp,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import { KafkaMessage } from 'kafkajs';
import Long from 'long';
import { DateTime } from 'luxon';
import { updateBlockCache } from '../../../src/caches/block-cache';
import { clearCandlesMap } from '../../../src/caches/candle-cache';
import {
  MILLIS_IN_NANOS,
  SECONDS_IN_MILLIS,
  STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE,
  SUBACCOUNT_ORDER_FILL_EVENT_TYPE,
} from '../../../src/constants';
import { OrderHandler } from '../../../src/handlers/order-fills/order-handler';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';
import { redisClient } from '../../../src/helpers/redis/redis-controller';
import { getWeightedAverage } from '../../../src/lib/helper';
import { onMessage } from '../../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../../src/lib/types';
import {
  defaultOrder, defaultOrderEvent, defaultPreviousHeight, defaultTakerOrder,
} from '../../helpers/constants';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  createKafkaMessageFromOrderFillEvent,
  createOrder,
  expectDefaultTradeKafkaMessageFromTakerFillId,
  expectFillInDatabase,
  expectOrderFillAndPositionSubaccountKafkaMessageFromIds,
  expectOrderInDatabase,
  expectPerpetualPosition,
  expectVulcanKafkaMessage,
} from '../../helpers/indexer-proto-helpers';
import { expectStateFilledQuantums } from '../../helpers/redis-helpers';

const defaultClobPairId: string = testConstants.defaultPerpetualMarket.clobPairId;
const defaultMakerFeeQuantum: number = 1_000_000;
const defaultTakerFeeQuantum: number = 2_000_000;
const defaultAffiliateRevShareQuantum: number = 3_000_000;
const defaultMakerFee: string = protocolTranslations.quantumsToHumanFixedString(
  defaultMakerFeeQuantum.toString(),
  testConstants.defaultAsset.atomicResolution,
);
const defaultTakerFee: string = protocolTranslations.quantumsToHumanFixedString(
  defaultTakerFeeQuantum.toString(),
  testConstants.defaultAsset.atomicResolution,
);
const defaultAffiliateRevShare: string = protocolTranslations.quantumsToHumanFixedString(
  defaultAffiliateRevShareQuantum.toString(),
  testConstants.defaultAsset.atomicResolution,
);

describe('OrderHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
    await assetRefresher.updateAssets();
    updateBlockCache(defaultPreviousHeight);
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
    clearCandlesMap();
    await redis.deleteAllAsync(redisClient);
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  const defaultHeight: string = '3';
  const defaultDateTime: DateTime = DateTime.utc(2022, 6, 1, 12, 1, 1, 2);
  const defaultTime: Timestamp = {
    seconds: Long.fromValue(Math.floor(defaultDateTime.toSeconds()), true),
    nanos: (defaultDateTime.toMillis() % SECONDS_IN_MILLIS) * MILLIS_IN_NANOS,
  };
  const defaultTxHash: string = '0x32343534306431622d306461302d343831322d613730372d3965613162336162';
  const defaultSubaccountId: IndexerSubaccountId = IndexerSubaccountId.fromPartial({
    owner: testConstants.defaultSubaccount.address,
    number: testConstants.defaultSubaccount.subaccountNumber,
  });
  const defaultSubaccountId2: IndexerSubaccountId = IndexerSubaccountId.fromPartial({
    owner: testConstants.defaultSubaccount2.address,
    number: testConstants.defaultSubaccount2.subaccountNumber,
  });
  const defaultPerpetualPosition: PerpetualPositionCreateObject = {
    subaccountId: testConstants.defaultSubaccountId,
    perpetualId: testConstants.defaultPerpetualMarket.id,
    side: PositionSide.LONG,
    status: PerpetualPositionStatus.OPEN,
    size: '10',
    maxSize: '25',
    sumOpen: '10',
    entryPrice: '15000',
    createdAt: DateTime.utc().toISO(),
    createdAtHeight: '10',
    openEventId: testConstants.defaultTendermintEventId4,
    lastEventId: testConstants.defaultTendermintEventId4,
    settledFunding: '200000',
  };

  describe('getParallelizationIds', () => {
    it.each([
      [
        'maker',
        Liquidity.MAKER,
        defaultOrderEvent.makerOrder!.orderId!,
      ],
      [
        'taker',
        Liquidity.TAKER,
        defaultTakerOrder.orderId!,
      ],
    ])('returns the correct %s parallelization ids', (
      _name: string,
      liquidity: Liquidity,
      orderId: IndexerOrderId,
    ) => {
      const subaccountId: IndexerSubaccountId = orderId.subaccountId!;
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.ORDER_FILL,
        OrderFillEventV1.encode(defaultOrderEvent).finish(),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: OrderHandler = new OrderHandler(
        block,
        0,
        indexerTendermintEvent,
        0,
        {
          ...defaultOrder,
          liquidity,
        },
      );

      const orderUuid: string = OrderTable.orderIdToUuid(orderId);
      expect(handler.getParallelizationIds()).toEqual([
        `${handler.eventType}_${SubaccountTable.subaccountIdToUuid(subaccountId)}_${defaultOrderEvent.makerOrder!.orderId!.clobPairId}`,
        `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${SubaccountTable.subaccountIdToUuid(subaccountId)}`,
        `${STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE}_${orderUuid}`,
      ]);
    });
  });

  it.each([
    [
      'goodTilBlock',
      {
        goodTilBlock: 10,
      },
      {
        goodTilBlock: 15,
      },
    ],
    [
      'goodTilBlockTime',
      {
        goodTilBlockTime: 1_000_000_000,
      },
      {
        goodTilBlockTime: 1_000_005_000,
      },
    ],
  ])(
    'creates fills and orders (with %s), sends vulcan messages for order updates and order ' +
    'removal for maker order fully filled, and updates perpetualPosition',
    async (
      _name: string,
      makerGoodTilOneof: Partial<IndexerOrder>,
      takerGoodTilOneof: Partial<IndexerOrder>,
    ) => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;
      const makerQuantums: number = 1_000_000;
      const makerSubticks: number = 100_000_000;

      const makerOrderProto: IndexerOrder = createOrder({
        subaccountId: defaultSubaccountId,
        clientId: 0,
        side: IndexerOrder_Side.SIDE_BUY,
        quantums: makerQuantums,
        subticks: makerSubticks,
        goodTilOneof: makerGoodTilOneof,
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL,
        reduceOnly: false,
        clientMetadata: 0,
      });

      const takerSubticks: number = 15_000_000;
      const takerQuantums: number = 10_000_000;
      const takerOrderProto: IndexerOrder = createOrder({
        subaccountId: defaultSubaccountId2,
        clientId: 0,
        side: IndexerOrder_Side.SIDE_SELL,
        quantums: takerQuantums,
        subticks: takerSubticks,
        goodTilOneof: takerGoodTilOneof,
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
        reduceOnly: true,
        clientMetadata: 0,
      });

      const fillAmount: number = 1_000_000;
      const orderFillEvent: OrderFillEventV1 = createOrderFillEvent(
        makerOrderProto,
        takerOrderProto,
        fillAmount,
        fillAmount,
        fillAmount,
      );
      const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
        orderFillEvent,
        transactionIndex,
        eventIndex,
        height: parseInt(defaultHeight, 10),
        time: defaultTime,
        txHash: defaultTxHash,
      });

      // create PerpetualPositions
      await Promise.all([
        PerpetualPositionTable.create(defaultPerpetualPosition),
        PerpetualPositionTable.create({
          ...defaultPerpetualPosition,
          subaccountId: testConstants.defaultSubaccountId2,
        }),
        // older perpetual position to ensure that the correct perpetual position is being updated
        PerpetualPositionTable.create({
          ...defaultPerpetualPosition,
          openEventId: testConstants.defaultTendermintEventId2,
        }),
      ]);

      const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
      await onMessage(kafkaMessage);

      const makerOrderSize: string = '0.0001'; // quantums in human = 1e6 * 1e-10 = 1e-4
      const makerPrice: string = '10000'; // quote currency / base currency = 1e8 * 1e-8 * 1e-6 / 1e-10 = 1e4
      const takerPrice: string = '1500'; // quote currency / base currency = 1.5e7 * 1e-8 * 1e-6 / 1e-10 = 1.5e3
      const totalFilled: string = '0.0001'; // fillAmount in human = 1e6 * 1e-10 = 1e-4
      await expectOrderInDatabase({
        subaccountId: testConstants.defaultSubaccountId,
        clientId: '0',
        size: makerOrderSize,
        totalFilled,
        price: makerPrice,
        status: OrderStatus.FILLED, // orderSize == totalFilled so status is filled
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
        orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
        timeInForce: TimeInForce.FOK,
        reduceOnly: false,
        goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
        updatedAt: defaultDateTime.toISO(),
        updatedAtHeight: defaultHeight.toString(),
      });

      const takerOrderSize: string = '0.001'; // quantums in human = 1e7 * 1e-10 = 1e-3
      await expectOrderInDatabase({
        subaccountId: testConstants.defaultSubaccountId2,
        clientId: '0',
        size: takerOrderSize,
        totalFilled,
        price: takerPrice,
        status: OrderStatus.OPEN, // orderSize > totalFilled so status is open
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
        orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
        timeInForce: TimeInForce.GTT,
        reduceOnly: true,
        goodTilBlock: protocolTranslations.getGoodTilBlock(takerOrderProto)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(takerOrderProto),
        clientMetadata: takerOrderProto.clientMetadata.toString(),
        updatedAt: defaultDateTime.toISO(),
        updatedAtHeight: defaultHeight.toString(),
      });

      const eventId: Buffer = TendermintEventTable.createEventId(
        defaultHeight,
        transactionIndex,
        eventIndex,
      );
      const quoteAmount: string = '1'; // quote amount is price * fillAmount = 1e4 * 1e-4 = 1
      await expectFillInDatabase({
        subaccountId: testConstants.defaultSubaccountId,
        clientId: '0',
        liquidity: Liquidity.MAKER,
        size: totalFilled,
        price: makerPrice,
        quoteAmount,
        eventId,
        transactionHash: defaultTxHash,
        createdAt: defaultDateTime.toISO(),
        createdAtHeight: defaultHeight,
        type: FillType.LIMIT,
        clobPairId: makerOrderProto.orderId!.clobPairId.toString(),
        side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
        orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
        fee: defaultMakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
        positionSideBefore: 'LONG',
        entryPriceBefore: '15000',
        positionSizeBefore: '10',
      });
      await expectFillInDatabase({
        subaccountId: testConstants.defaultSubaccountId2,
        clientId: '0',
        liquidity: Liquidity.TAKER,
        size: totalFilled,
        price: makerPrice,
        quoteAmount,
        eventId,
        transactionHash: defaultTxHash,
        createdAt: defaultDateTime.toISO(),
        createdAtHeight: defaultHeight,
        type: FillType.LIMIT,
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
        orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
        clientMetadata: takerOrderProto.clientMetadata.toString(),
        fee: defaultTakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
        positionSideBefore: 'LONG',
        entryPriceBefore: '15000',
        positionSizeBefore: '10',
      });

      const expectedMakerOffchainUpdate: OffChainUpdateV1 = {
        orderUpdate: {
          orderId: makerOrderProto.orderId,
          totalFilledQuantums: orderFillEvent.totalFilledMaker,
        },
      };
      const expectedMakerRemoveOffchainUpdate: OffChainUpdateV1 = {
        orderRemove: {
          removedOrderId: makerOrderProto.orderId,
          reason: OrderRemovalReason.ORDER_REMOVAL_REASON_FULLY_FILLED,
          removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_FILLED,
        },
      };
      const expectedTakerOffchainUpdate: OffChainUpdateV1 = {
        orderUpdate: {
          orderId: takerOrderProto.orderId,
          totalFilledQuantums: orderFillEvent.totalFilledTaker,
        },
      };

      await Promise.all([
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: makerOrderProto.orderId!,
          offchainUpdate: expectedMakerOffchainUpdate,
        }),
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: takerOrderProto.orderId!,
          offchainUpdate: expectedTakerOffchainUpdate,
        }),
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: makerOrderProto.orderId!,
          offchainUpdate: expectedMakerRemoveOffchainUpdate,
        }),
        expectDefaultOrderAndFillSubaccountKafkaMessages(
          producerSendMock,
          eventId,
          ORDER_FLAG_LONG_TERM,
          ORDER_FLAG_SHORT_TERM,
        ),
        expectDefaultTradeKafkaMessageFromTakerFillId(
          producerSendMock,
          eventId,
        ),
        expectPerpetualPosition(
          PerpetualPositionTable.uuid(
            testConstants.defaultSubaccountId,
            defaultPerpetualPosition.openEventId,
          ),
          {
            sumOpen: Big(defaultPerpetualPosition.size).plus(totalFilled).toFixed(),
            entryPrice: getWeightedAverage(
              defaultPerpetualPosition.entryPrice!,
              defaultPerpetualPosition.size,
              makerPrice,
              totalFilled,
            ),
            totalRealizedPnl: '-1',
          },
        ),
        expectPerpetualPosition(
          PerpetualPositionTable.uuid(
            testConstants.defaultSubaccountId2,
            defaultPerpetualPosition.openEventId,
          ),
          {
            sumClose: totalFilled,
            exitPrice: makerPrice,
            totalRealizedPnl: '-2.5000',
          },
        ),
        expectCandlesUpdated(),
        expectStateFilledQuantums(
          OrderTable.orderIdToUuid(makerOrderProto.orderId!),
          orderFillEvent.totalFilledMaker.toString(),
        ),
        expectStateFilledQuantums(
          OrderTable.orderIdToUuid(takerOrderProto.orderId!),
          orderFillEvent.totalFilledTaker.toString(),
        ),
      ]);
    });

  it.each([
    [
      'goodTilBlock',
      {
        goodTilBlock: 10,
      },
      {
        goodTilBlock: 15,
      },
      false,
      '5',
      undefined,
    ],
    [
      'goodTilBlockTime',
      {
        goodTilBlockTime: 1_000_000_000,
      },
      {
        goodTilBlockTime: 1_000_005_000,
      },
      false,
      undefined,
      '1970-01-11T13:46:40.000Z',
    ],
    [
      'goodTilBlockTime w/ cancelled maker order',
      {
        goodTilBlockTime: 1_000_000_000,
      },
      {
        goodTilBlockTime: 1_000_005_000,
      },
      true,
      undefined,
      '1970-01-11T13:46:40.000Z',
    ],
  ])(
    'updates existing orders (with %s), sends vulcan messages for order updates and order ' +
    'removal for taker order fully-filled',
    async (
      _name: string,
      makerGoodTilOneof: Partial<IndexerOrder>,
      takerGoodTilOneof: Partial<IndexerOrder>,
      isOrderCanceled: boolean,
      existingGoodTilBlock?: string,
      existingGoodTilBlockTime?: string,
    ) => {
      if (isOrderCanceled) {
        await CanceledOrdersCache.addCanceledOrderId(
          OrderTable.uuid(
            testConstants.defaultSubaccountId,
            '0',
            defaultClobPairId,
            ORDER_FLAG_SHORT_TERM.toString(),
          ),
          Date.now(),
          redisClient,
        );
      }
      // create initial orders
      await Promise.all([
        // maker order
        OrderTable.create({
          subaccountId: testConstants.defaultSubaccountId,
          clientId: '0',
          clobPairId: defaultClobPairId,
          side: OrderSide.BUY,
          size: '1',
          totalFilled: '0.1',
          price: '10000',
          type: OrderType.LIMIT,
          status: OrderStatus.OPEN,
          timeInForce: TimeInForce.GTT,
          reduceOnly: false,
          goodTilBlock: existingGoodTilBlock,
          goodTilBlockTime: existingGoodTilBlockTime,
          orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
          clientMetadata: '0',
          updatedAt: DateTime.fromMillis(0).toISO(),
          updatedAtHeight: '0',
        }),
        // taker order
        OrderTable.create({
          subaccountId: testConstants.defaultSubaccountId2,
          clientId: '0',
          clobPairId: defaultClobPairId,
          side: OrderSide.SELL,
          size: '1',
          totalFilled: '0.1',
          price: '10000',
          type: OrderType.LIMIT,
          status: OrderStatus.OPEN,
          timeInForce: TimeInForce.GTT,
          reduceOnly: false,
          goodTilBlock: existingGoodTilBlock,
          goodTilBlockTime: existingGoodTilBlockTime,
          orderFlags: ORDER_FLAG_LONG_TERM.toString(),
          clientMetadata: '0',
          updatedAt: DateTime.fromMillis(0).toISO(),
          updatedAtHeight: '0',
        }),
      ]);

      // create initial PerpetualPositions
      await Promise.all([
        PerpetualPositionTable.create(defaultPerpetualPosition),
        PerpetualPositionTable.create({
          ...defaultPerpetualPosition,
          subaccountId: testConstants.defaultSubaccountId2,
        }),
      ]);

      const transactionIndex: number = 0;
      const eventIndex: number = 0;
      const makerQuantums: number = 11_001_000_000;
      const subticks: number = 100_000_000;
      const takerQuantums: number = 1_002_000_000;

      const makerOrderProto: IndexerOrder = createOrder({
        subaccountId: defaultSubaccountId,
        clientId: 0,
        side: IndexerOrder_Side.SIDE_BUY,
        quantums: makerQuantums,
        subticks,
        goodTilOneof: makerGoodTilOneof,
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
        reduceOnly: true,
        clientMetadata: 0,
      });
      const takerOrderProto: IndexerOrder = createOrder({
        subaccountId: defaultSubaccountId2,
        clientId: 0,
        side: IndexerOrder_Side.SIDE_SELL,
        quantums: takerQuantums,
        subticks,
        goodTilOneof: takerGoodTilOneof,
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY,
        reduceOnly: false,
        clientMetadata: 0,
      });

      const fillAmount: number = 1_000_000;
      const totalMakerFilled: number = 1_001_000_000;
      const totalTakerFilled: number = 1_002_000_000;
      const orderFillEvent: OrderFillEventV1 = createOrderFillEvent(
        makerOrderProto,
        takerOrderProto,
        fillAmount,
        totalMakerFilled,
        totalTakerFilled,
      );
      const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
        orderFillEvent,
        transactionIndex,
        eventIndex,
        height: parseInt(defaultHeight, 10),
        time: defaultTime,
        txHash: defaultTxHash,
      });
      const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
      await onMessage(kafkaMessage);

      const makerOrderSize: string = '1.1001'; // quantums in human = (1e10 + 1e9 + 1e6) * 1e-10 = 1.1001
      const price: string = '10000'; // quote currency / base currency = 1e8 * 1e-8 * 1e-6 / 1e-10 = 1e4
      const totalMakerOrderFilled: string = '0.1001';
      const totalTakerOrderFilled: string = '0.1002';
      await expectOrderInDatabase({
        subaccountId: testConstants.defaultSubaccountId,
        clientId: '0',
        size: makerOrderSize,
        totalFilled: totalMakerOrderFilled,
        price,
        status: isOrderCanceled
          ? OrderStatus.CANCELED
          : OrderStatus.OPEN, // orderSize > totalFilled so status is open
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
        orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
        timeInForce: TimeInForce.GTT,
        reduceOnly: true,
        goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
        updatedAt: defaultDateTime.toISO(),
        updatedAtHeight: defaultHeight.toString(),
      });

      const takerOrderSize: string = '0.1002'; // quantums in human = (1e9 + 2e6) * 1e-10 = 0.1002
      await expectOrderInDatabase({
        subaccountId: testConstants.defaultSubaccountId2,
        clientId: '0',
        size: takerOrderSize,
        totalFilled: totalTakerOrderFilled,
        price,
        status: OrderStatus.FILLED, // orderSize == totalFilled so status is filled
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
        orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
        timeInForce: TimeInForce.POST_ONLY,
        reduceOnly: false,
        goodTilBlock: protocolTranslations.getGoodTilBlock(takerOrderProto)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(takerOrderProto),
        clientMetadata: takerOrderProto.clientMetadata.toString(),
        updatedAt: defaultDateTime.toISO(),
        updatedAtHeight: defaultHeight.toString(),
      });

      const eventId: Buffer = TendermintEventTable.createEventId(
        defaultHeight,
        transactionIndex,
        eventIndex,
      );
      const quoteAmount: string = '1'; // quote amount is price * fillAmount = 1e4 * 1e-4 = 2
      const fillAmountInHuman: string = '0.0001'; // fillAmount in human = 1e6 * 1e-10 = 1e-4
      await expectFillInDatabase({
        subaccountId: testConstants.defaultSubaccountId,
        clientId: '0',
        liquidity: Liquidity.MAKER,
        size: fillAmountInHuman,
        price,
        quoteAmount,
        eventId,
        transactionHash: defaultTxHash,
        createdAt: defaultDateTime.toISO(),
        createdAtHeight: defaultHeight,
        type: FillType.LIMIT,
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
        orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
        fee: defaultMakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
        positionSideBefore: 'LONG',
        entryPriceBefore: '15000',
        positionSizeBefore: '10',
      });
      await expectFillInDatabase({
        subaccountId: testConstants.defaultSubaccountId2,
        clientId: '0',
        liquidity: Liquidity.TAKER,
        size: fillAmountInHuman,
        price,
        quoteAmount,
        eventId,
        transactionHash: defaultTxHash,
        createdAt: defaultDateTime.toISO(),
        createdAtHeight: defaultHeight,
        type: FillType.LIMIT,
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
        orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
        clientMetadata: takerOrderProto.clientMetadata.toString(),
        fee: defaultTakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
        positionSideBefore: 'LONG',
        entryPriceBefore: '15000',
        positionSizeBefore: '10',
      });

      const expectedMakerUpdateOffchainUpdate: OffChainUpdateV1 = {
        orderUpdate: {
          orderId: makerOrderProto.orderId,
          totalFilledQuantums: orderFillEvent.totalFilledMaker,
        },
      };
      const expectedTakerUpdateOffchainUpdate: OffChainUpdateV1 = {
        orderUpdate: {
          orderId: takerOrderProto.orderId,
          totalFilledQuantums: orderFillEvent.totalFilledTaker,
        },
      };
      const expectedTakerRemoveOffchainUpdate: OffChainUpdateV1 = {
        orderRemove: {
          removedOrderId: takerOrderProto.orderId,
          reason: OrderRemovalReason.ORDER_REMOVAL_REASON_FULLY_FILLED,
          removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_FILLED,
        },
      };

      await Promise.all([
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: makerOrderProto.orderId!,
          offchainUpdate: expectedMakerUpdateOffchainUpdate,
        }),
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: takerOrderProto.orderId!,
          offchainUpdate: expectedTakerUpdateOffchainUpdate,
        }),
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: takerOrderProto.orderId!,
          offchainUpdate: expectedTakerRemoveOffchainUpdate,
        }),
        expectDefaultOrderAndFillSubaccountKafkaMessages(
          producerSendMock,
          eventId,
          ORDER_FLAG_SHORT_TERM,
          ORDER_FLAG_LONG_TERM,
        ),
        expectDefaultTradeKafkaMessageFromTakerFillId(
          producerSendMock,
          eventId,
        ),
        expectCandlesUpdated(),
        expectStateFilledQuantums(
          OrderTable.orderIdToUuid(makerOrderProto.orderId!),
          orderFillEvent.totalFilledMaker.toString(),
        ),
        expectStateFilledQuantums(
          OrderTable.orderIdToUuid(takerOrderProto.orderId!),
          orderFillEvent.totalFilledTaker.toString(),
        ),
      ]);
    });

  it.each([
    [
      'goodTilBlock',
      {
        goodTilBlock: 10,
      },
      {
        goodTilBlock: 15,
      },
      '5',
      undefined,
    ],
    [
      'goodTilBlockTime',
      {
        goodTilBlockTime: 1_000_000_000,
      },
      {
        goodTilBlockTime: 1_000_005_000,
      },
      undefined,
      '1970-01-11T13:46:40.000Z',
    ],
  ])(
    'replaces existing orders (with %s), upserting a new order with the same order id',
    async (
      _name: string,
      makerGoodTilOneof: Partial<IndexerOrder>,
      takerGoodTilOneof: Partial<IndexerOrder>,
      existingGoodTilBlock?: string,
      existingGoodTilBlockTime?: string,
    ) => {
      // create initial orders
      await Promise.all([
        // maker order
        OrderTable.create({
          subaccountId: testConstants.defaultSubaccountId,
          clientId: '0',
          clobPairId: defaultClobPairId,
          side: OrderSide.BUY,
          size: '1',
          totalFilled: '0.1',
          price: '10000',
          type: OrderType.LIMIT,
          status: OrderStatus.OPEN,
          timeInForce: TimeInForce.GTT,
          reduceOnly: false,
          goodTilBlock: existingGoodTilBlock,
          goodTilBlockTime: existingGoodTilBlockTime,
          orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
          clientMetadata: '0',
          updatedAt: DateTime.fromMillis(0).toISO(),
          updatedAtHeight: '0',
        }),
        // taker order
        OrderTable.create({
          subaccountId: testConstants.defaultSubaccountId2,
          clientId: '0',
          clobPairId: defaultClobPairId,
          side: OrderSide.SELL,
          size: '1',
          totalFilled: '0.1',
          price: '10000',
          type: OrderType.LIMIT,
          status: OrderStatus.OPEN,
          timeInForce: TimeInForce.GTT,
          reduceOnly: false,
          goodTilBlock: existingGoodTilBlock,
          goodTilBlockTime: existingGoodTilBlockTime,
          orderFlags: ORDER_FLAG_LONG_TERM.toString(),
          clientMetadata: '0',
          updatedAt: DateTime.fromMillis(0).toISO(),
          updatedAtHeight: '0',
        }),
      ]);

      // create initial PerpetualPositions
      await Promise.all([
        PerpetualPositionTable.create(defaultPerpetualPosition),
        PerpetualPositionTable.create({
          ...defaultPerpetualPosition,
          subaccountId: testConstants.defaultSubaccountId2,
        }),
      ]);

      const transactionIndex: number = 0;
      const eventIndex: number = 0;
      const newMakerQuantums: number = 21_001_000_000;
      const newSubticks: number = 200_000_000;
      const newTakerQuantums: number = 2_002_000_000;

      const makerOrderProto: IndexerOrder = createOrder({
        subaccountId: defaultSubaccountId,
        clientId: 0,
        side: IndexerOrder_Side.SIDE_SELL,
        quantums: newMakerQuantums,
        subticks: newSubticks,
        goodTilOneof: makerGoodTilOneof,
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY,
        reduceOnly: false,
        clientMetadata: 0,
      });
      const takerOrderProto: IndexerOrder = createOrder({
        subaccountId: defaultSubaccountId2,
        clientId: 0,
        side: IndexerOrder_Side.SIDE_BUY,
        quantums: newTakerQuantums,
        subticks: newSubticks,
        goodTilOneof: takerGoodTilOneof,
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
        reduceOnly: true,
        clientMetadata: 0,
      });

      const fillAmount: number = 1_000_000;
      const totalMakerFilled: number = 1_000_000;
      const totalTakerFilled: number = 1_000_000;
      const orderFillEvent: OrderFillEventV1 = createOrderFillEvent(
        makerOrderProto,
        takerOrderProto,
        fillAmount,
        totalMakerFilled,
        totalTakerFilled,
      );
      const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
        orderFillEvent,
        transactionIndex,
        eventIndex,
        height: parseInt(defaultHeight, 10),
        time: defaultTime,
        txHash: defaultTxHash,
      });
      const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
      await onMessage(kafkaMessage);

      const makerOrderSize: string = '2.1001'; // quantums in human = (2e10 + 1e9 + 1e6) * 1e-10 = 2.1001
      const price: string = '20000'; // quote currency / base currency = 2e8 * 1e-8 * 1e-6 / 1e-10 = 2e4
      const totalMakerOrderFilled: string = '0.0001';
      const totalTakerOrderFilled: string = '0.0001';
      await expectOrderInDatabase({
        subaccountId: testConstants.defaultSubaccountId,
        clientId: '0',
        size: makerOrderSize,
        totalFilled: totalMakerOrderFilled,
        price,
        status: OrderStatus.OPEN, // orderSize > totalFilled so status is OPEN
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
        orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
        timeInForce: TimeInForce.POST_ONLY,
        reduceOnly: false,
        goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
        updatedAt: defaultDateTime.toISO(),
        updatedAtHeight: defaultHeight.toString(),
      });

      const takerOrderSize: string = '0.2002'; // quantums in human = (2e9 + 2e6) * 1e-10 = 0.1002
      await expectOrderInDatabase({
        subaccountId: testConstants.defaultSubaccountId2,
        clientId: '0',
        size: takerOrderSize,
        totalFilled: totalTakerOrderFilled,
        price,
        status: OrderStatus.OPEN, // orderSize > totalFilled so status is OPEN
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
        orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
        timeInForce: TimeInForce.GTT,
        reduceOnly: true,
        goodTilBlock: protocolTranslations.getGoodTilBlock(takerOrderProto)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(takerOrderProto),
        clientMetadata: takerOrderProto.clientMetadata.toString(),
        updatedAt: defaultDateTime.toISO(),
        updatedAtHeight: defaultHeight.toString(),
      });

      const eventId: Buffer = TendermintEventTable.createEventId(
        defaultHeight,
        transactionIndex,
        eventIndex,
      );
      const quoteAmount: string = '2'; // quote amount is price * fillAmount = 2e4 * 1e-4 = 1
      const fillAmountInHuman: string = '0.0001'; // fillAmount in human = 1e6 * 1e-10 = 1e-4
      await expectFillInDatabase({
        subaccountId: testConstants.defaultSubaccountId,
        clientId: '0',
        liquidity: Liquidity.MAKER,
        size: fillAmountInHuman,
        price,
        quoteAmount,
        eventId,
        transactionHash: defaultTxHash,
        createdAt: defaultDateTime.toISO(),
        createdAtHeight: defaultHeight,
        type: FillType.LIMIT,
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
        orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
        fee: defaultMakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
        positionSideBefore: 'LONG',
        entryPriceBefore: '15000',
        positionSizeBefore: '10',
      });
      await expectFillInDatabase({
        subaccountId: testConstants.defaultSubaccountId2,
        clientId: '0',
        liquidity: Liquidity.TAKER,
        size: fillAmountInHuman,
        price,
        quoteAmount,
        eventId,
        transactionHash: defaultTxHash,
        createdAt: defaultDateTime.toISO(),
        createdAtHeight: defaultHeight,
        type: FillType.LIMIT,
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
        orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
        clientMetadata: takerOrderProto.clientMetadata.toString(),
        fee: defaultTakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
        positionSideBefore: 'LONG',
        entryPriceBefore: '15000',
        positionSizeBefore: '10',
      });

      const expectedMakerUpdateOffchainUpdate: OffChainUpdateV1 = {
        orderUpdate: {
          orderId: makerOrderProto.orderId,
          totalFilledQuantums: orderFillEvent.totalFilledMaker,
        },
      };
      const expectedTakerUpdateOffchainUpdate: OffChainUpdateV1 = {
        orderUpdate: {
          orderId: takerOrderProto.orderId,
          totalFilledQuantums: orderFillEvent.totalFilledTaker,
        },
      };

      await Promise.all([
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: makerOrderProto.orderId!,
          offchainUpdate: expectedMakerUpdateOffchainUpdate,
        }),
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: takerOrderProto.orderId!,
          offchainUpdate: expectedTakerUpdateOffchainUpdate,
        }),
        expectDefaultOrderAndFillSubaccountKafkaMessages(
          producerSendMock,
          eventId,
          ORDER_FLAG_SHORT_TERM,
          ORDER_FLAG_LONG_TERM,
        ),
        expectDefaultTradeKafkaMessageFromTakerFillId(
          producerSendMock,
          eventId,
        ),
        expectCandlesUpdated(),
        expectStateFilledQuantums(
          OrderTable.orderIdToUuid(makerOrderProto.orderId!),
          orderFillEvent.totalFilledMaker.toString(),
        ),
        expectStateFilledQuantums(
          OrderTable.orderIdToUuid(takerOrderProto.orderId!),
          orderFillEvent.totalFilledTaker.toString(),
        ),
      ]);
    });

  it('creates fills and orders with fixed-point notation quoteAmount', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const makerQuantums: number = 100;
    const makerSubticks: number = 1_000_000;

    const makerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: makerQuantums,
      subticks: makerSubticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: defaultClobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
      reduceOnly: false,
      clientMetadata: 0,
    });

    const takerSubticks: number = 150_000;
    const takerQuantums: number = 10;
    const takerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: takerQuantums,
      subticks: takerSubticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: defaultClobPairId,
      orderFlags: ORDER_FLAG_LONG_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: true,
      clientMetadata: 0,
    });

    // create initial PerpetualPositions with closed previous positions
    await Promise.all([
      // previous position for subaccount 1
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId,
      }),
      // previous position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        subaccountId: testConstants.defaultSubaccountId2,
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId,
      }),
      // initial position for subaccount 2
      PerpetualPositionTable.create(defaultPerpetualPosition),
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        subaccountId: testConstants.defaultSubaccountId2,
      }),
    ]);

    const fillAmount: number = 10;
    const orderFillEvent: OrderFillEventV1 = createOrderFillEvent(
      makerOrderProto,
      takerOrderProto,
      fillAmount,
      fillAmount,
      fillAmount,
    );
    const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
      orderFillEvent,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10),
      time: defaultTime,
      txHash: defaultTxHash,
    });

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    // This size should be in fixed-point notation rather than exponential notation (1e-8)
    const makerOrderSize: string = '0.00000001'; // quantums in human = 1e2 * 1e-10 = 1e-8
    const makerPrice: string = '100'; // quote currency / base currency = 1e6 * 1e-8 * 1e-6 / 1e-10 = 1e2
    const takerPrice: string = '15'; // quote currency / base currency = 1.5e5 * 1e-8 * 1e-6 / 1e-10 = 1.5e1
    const totalFilled: string = '0.000000001'; // fillAmount in human = 1e1 * 1e-10 = 1e-9
    await expectOrderInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '0',
      size: makerOrderSize,
      totalFilled,
      price: makerPrice,
      status: OrderStatus.OPEN, // orderSize > totalFilled so status is open
      clobPairId: defaultClobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
      orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
      timeInForce: TimeInForce.GTT,
      reduceOnly: false,
      goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    });

    // This size should be in fixed-point notation rather than exponential notation (1e-9)
    const takerOrderSize: string = '0.000000001'; // quantums in human = 1e1 * 1e-10 = 1e-9
    await expectOrderInDatabase({
      subaccountId: testConstants.defaultSubaccountId2,
      clientId: '0',
      size: takerOrderSize,
      totalFilled,
      price: takerPrice,
      status: OrderStatus.FILLED, // orderSize == totalFilled so status is filled
      clobPairId: defaultClobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
      orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
      timeInForce: TimeInForce.IOC,
      reduceOnly: true,
      goodTilBlock: protocolTranslations.getGoodTilBlock(takerOrderProto)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(takerOrderProto),
      clientMetadata: takerOrderProto.clientMetadata.toString(),
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    });

    const eventId: Buffer = TendermintEventTable.createEventId(
      defaultHeight,
      transactionIndex,
      eventIndex,
    );

    // This size should be in fixed-point notation rather than exponential notation (1e-5)
    const quoteAmount: string = '0.0000001'; // quote amount is price * fillAmount = 1e2 * 1e-9 = 1e-7
    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '0',
      liquidity: Liquidity.MAKER,
      size: totalFilled,
      price: makerPrice,
      quoteAmount,
      eventId,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight,
      type: FillType.LIMIT,
      clobPairId: defaultClobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
      orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
      fee: defaultMakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      positionSideBefore: 'LONG',
      entryPriceBefore: '15000',
      positionSizeBefore: '10',
    });
    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId2,
      clientId: '0',
      liquidity: Liquidity.TAKER,
      size: totalFilled,
      price: makerPrice,
      quoteAmount,
      eventId,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight,
      type: FillType.LIMIT,
      clobPairId: defaultClobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
      orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
      clientMetadata: takerOrderProto.clientMetadata.toString(),
      fee: defaultTakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      positionSideBefore: 'LONG',
      entryPriceBefore: '15000',
      positionSizeBefore: '10',
    });

    await Promise.all([
      expectDefaultOrderAndFillSubaccountKafkaMessages(
        producerSendMock,
        eventId,
        ORDER_FLAG_SHORT_TERM,
        ORDER_FLAG_LONG_TERM,
      ),
      expectDefaultTradeKafkaMessageFromTakerFillId(
        producerSendMock,
        eventId,
      ),
      expectCandlesUpdated(),
      expectStateFilledQuantums(
        OrderTable.orderIdToUuid(makerOrderProto.orderId!),
        orderFillEvent.totalFilledMaker.toString(),
      ),
      expectStateFilledQuantums(
        OrderTable.orderIdToUuid(takerOrderProto.orderId!),
        orderFillEvent.totalFilledTaker.toString(),
      ),
    ]);
  });

  it('creates fills and orders with fixed-point notation price', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const makerQuantums: number = 100;
    const makerSubticks: number = 1_000_000;

    const makerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: makerQuantums,
      subticks: makerSubticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
      reduceOnly: false,
      clientMetadata: 0,
    });

    const takerSubticks: number = 150_000;
    const takerQuantums: number = 10;
    const takerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: takerQuantums,
      subticks: takerSubticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_LONG_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: true,
      clientMetadata: 0,
    });

    // create initial PerpetualPositions with closed previous positions
    await Promise.all([
      // previous position for subaccount 1
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId,
      }),
      // previous position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
        subaccountId: testConstants.defaultSubaccountId2,
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId,
      }),
      // initial position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
      }),
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
        subaccountId: testConstants.defaultSubaccountId2,
      }),
    ]);

    const fillAmount: number = 10;
    const orderFillEvent: OrderFillEventV1 = createOrderFillEvent(
      makerOrderProto,
      takerOrderProto,
      fillAmount,
      fillAmount,
      fillAmount,
    );
    const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
      orderFillEvent,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10),
      time: defaultTime,
      txHash: defaultTxHash,
    });

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    // This price should be in fixed-point notation rather than exponential notation (1e-8)
    const makerOrderSize: string = '1'; // quantums in human = 1e2 * 1e-2 = 1
    const makerPrice: string = '0.00000000000001'; // quote currency / base currency = 1e6 * 1e-16 * 1e-6 / 1e-2 = 1e-14
    const takerPrice: string = '0.0000000000000015'; // quote currency / base currency = 1.5e5 * 1e-16 * 1e-6 / 1e-2 = 1.5e-15
    const totalFilled: string = '0.1'; // fillAmount in human = 1e1 * 1e-2 = 1e-1
    await expectOrderInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '0',
      size: makerOrderSize,
      totalFilled,
      price: makerPrice,
      status: OrderStatus.OPEN, // orderSize > totalFilled so status is open
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
      orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
      timeInForce: TimeInForce.GTT,
      reduceOnly: false,
      goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    });

    const takerOrderSize: string = '0.1'; // quantums in human = 1e1 * 1e-2 = 1e-1
    await expectOrderInDatabase({
      subaccountId: testConstants.defaultSubaccountId2,
      clientId: '0',
      size: takerOrderSize,
      totalFilled,
      price: takerPrice,
      status: OrderStatus.FILLED, // orderSize == totalFilled so status is filled
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
      orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
      timeInForce: TimeInForce.IOC,
      reduceOnly: true,
      goodTilBlock: protocolTranslations.getGoodTilBlock(takerOrderProto)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(takerOrderProto),
      clientMetadata: takerOrderProto.clientMetadata.toString(),
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    });

    const eventId: Buffer = TendermintEventTable.createEventId(
      defaultHeight,
      transactionIndex,
      eventIndex,
    );

    const quoteAmount: string = '0.000000000000001'; // quote amount is price * fillAmount = 1e-14 * 1e-1 = 1e-15
    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '0',
      liquidity: Liquidity.MAKER,
      size: totalFilled,
      price: makerPrice,
      quoteAmount,
      eventId,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight,
      type: FillType.LIMIT,
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
      orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
      fee: defaultMakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      positionSideBefore: 'LONG',
      entryPriceBefore: '15000',
      positionSizeBefore: '10',
    });
    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId2,
      clientId: '0',
      liquidity: Liquidity.TAKER,
      size: totalFilled,
      price: makerPrice,
      quoteAmount,
      eventId,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight,
      type: FillType.LIMIT,
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
      orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
      clientMetadata: takerOrderProto.clientMetadata.toString(),
      fee: defaultTakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      positionSideBefore: 'LONG',
      entryPriceBefore: '15000',
      positionSizeBefore: '10',
    });

    await Promise.all([
      expectDefaultOrderAndFillSubaccountKafkaMessages(
        producerSendMock,
        eventId,
        ORDER_FLAG_SHORT_TERM,
        ORDER_FLAG_LONG_TERM,
        testConstants.defaultPerpetualMarket3.id,
        testConstants.defaultPerpetualMarket3.clobPairId,
      ),
      expectDefaultTradeKafkaMessageFromTakerFillId(
        producerSendMock,
        eventId,
      ),
      expectCandlesUpdated(),
      expectStateFilledQuantums(
        OrderTable.orderIdToUuid(makerOrderProto.orderId!),
        orderFillEvent.totalFilledMaker.toString(),
      ),
      expectStateFilledQuantums(
        OrderTable.orderIdToUuid(takerOrderProto.orderId!),
        orderFillEvent.totalFilledTaker.toString(),
      ),
    ]);
  });

  it('creates twap fills and orders', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const makerQuantums: number = 100;
    const makerSubticks: number = 1_000_000;

    const makerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: makerQuantums,
      subticks: makerSubticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
      reduceOnly: false,
      clientMetadata: 0,
    });

    const takerSubticks: number = 150_000;
    const takerQuantums: number = 10;
    const takerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: takerQuantums,
      subticks: takerSubticks,
      goodTilOneof: {
        goodTilBlockTime: 10,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_TWAP.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: true,
      clientMetadata: 0,
      duration: 300,
      interval: 30,
      priceTolerance: 10000,
    });

    const takerSuborderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: 1, // takerQuantums / (duration / interval)
      subticks: takerSubticks,
      goodTilOneof: {
        goodTilBlockTime: 10,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_TWAP_SUBORDER.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: true,
      clientMetadata: 0,
    });

    // create initial PerpetualPositions with closed previous positions
    await Promise.all([
      // previous position for subaccount 1
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId,
      }),
      // previous position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
        subaccountId: testConstants.defaultSubaccountId2,
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId,
      }),
      // initial position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
      }),
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
        subaccountId: testConstants.defaultSubaccountId2,
      }),
    ]);

    const fillAmount: number = 1;
    const orderFillEvent: OrderFillEventV1 = createOrderFillEvent(
      makerOrderProto,
      takerSuborderProto,
      fillAmount,
      fillAmount,
      fillAmount,
    );
    const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
      orderFillEvent,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10),
      time: defaultTime,
      txHash: defaultTxHash,
    });

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    // This price should be in fixed-point notation rather than exponential notation (1e-8)
    const makerOrderSize: string = '1'; // quantums in human = 1e2 * 1e-2 = 1
    const makerPrice: string = '0.00000000000001'; // quote currency / base currency = 1e6 * 1e-16 * 1e-6 / 1e-2 = 1e-14
    const totalFilled: string = '0.01'; // fillAmount in human = 1e1 * 1e-2 = 1e-1
    await expectOrderInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '0',
      size: makerOrderSize,
      totalFilled,
      price: makerPrice,
      status: OrderStatus.OPEN, // orderSize > totalFilled so status is open
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
      orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
      timeInForce: TimeInForce.GTT,
      reduceOnly: false,
      goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
    });

    const takerOrderSize: string = '0.01'; // quantums in human = 1e1 * 1e-2 = 1e-1
    await expectOrderInDatabase({
      subaccountId: testConstants.defaultSubaccountId2,
      clientId: '0',
      size: takerOrderSize,
      totalFilled,
      price: makerPrice,
      status: OrderStatus.FILLED, // orderSize == totalFilled so status is filled
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
      orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
      timeInForce: TimeInForce.IOC,
      reduceOnly: true,
      goodTilBlock: protocolTranslations.getGoodTilBlock(takerOrderProto)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(takerOrderProto),
      clientMetadata: takerOrderProto.clientMetadata.toString(),
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
      orderType: OrderType.TWAP,
    });

    const eventId: Buffer = TendermintEventTable.createEventId(
      defaultHeight,
      transactionIndex,
      eventIndex,
    );

    const quoteAmount: string = '0.0000000000000001'; // quote amount is price * fillAmount = 1e-14 * 1e-1 = 1e-15
    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '0',
      liquidity: Liquidity.MAKER,
      size: totalFilled,
      price: makerPrice,
      quoteAmount,
      eventId,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight,
      type: FillType.LIMIT,
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
      orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
      fee: defaultMakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      positionSideBefore: 'LONG',
      entryPriceBefore: '15000',
      positionSizeBefore: '10',
    });
    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId2,
      clientId: '0',
      liquidity: Liquidity.TAKER,
      size: totalFilled,
      price: makerPrice,
      quoteAmount,
      eventId,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight,
      type: FillType.TWAP_SUBORDER,
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
      orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
      clientMetadata: takerOrderProto.clientMetadata.toString(),
      fee: defaultTakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      positionSideBefore: 'LONG',
      entryPriceBefore: '15000',
      positionSizeBefore: '10',
    });

    await Promise.all([
      expectDefaultOrderAndFillSubaccountKafkaMessages(
        producerSendMock,
        eventId,
        ORDER_FLAG_SHORT_TERM,
        ORDER_FLAG_TWAP,
        testConstants.defaultPerpetualMarket3.id,
        testConstants.defaultPerpetualMarket3.clobPairId,
      ),
      expectDefaultTradeKafkaMessageFromTakerFillId(
        producerSendMock,
        eventId,
      ),
      expectCandlesUpdated(),
      expectStateFilledQuantums(
        OrderTable.orderIdToUuid(makerOrderProto.orderId!),
        orderFillEvent.totalFilledMaker.toString(),
      ),
      expectStateFilledQuantums(
        OrderTable.orderIdToUuid(takerOrderProto.orderId!),
        orderFillEvent.totalFilledTaker.toString(),
      ),
    ]);
  });

  it('creates and updates twaps through suborder fills', async () => {
    // Create a parent TWAP order with duration 300 and interval 30
    const parentTwapOrder = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 123,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: 100_000_000_000, // 100 units
      subticks: 1_000_000, // price
      goodTilOneof: {
        goodTilBlock: 100,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_TWAP.toString(), // 128
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: false,
      clientMetadata: 2,
      duration: 300,
      interval: 30,
      priceTolerance: 0.01,
    });

    // The parent order is inserted by the handler, so we don't need to insert it manually.

    // Create two suborders (TWAP suborders) that will be filled
    const suborder1 = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 123,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: 10_000_000_000, // 10 units
      subticks: 1_000_000,
      goodTilOneof: {
        goodTilBlock: 100,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_TWAP_SUBORDER.toString(), // 256
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: false,
      clientMetadata: 2,
    });

    const suborder2 = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 123,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: 10_000_000_000, // 10 units
      subticks: 2_000_000,
      goodTilOneof: {
        goodTilBlock: 100,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_TWAP_SUBORDER.toString(), // 256
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: false,
      clientMetadata: 2,
    });

    // Create a maker order that will match the suborders
    const makerOrder1 = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 200,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: 10_000_000_000,
      subticks: 1_000_000,
      goodTilOneof: {
        goodTilBlock: 100,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: false,
      clientMetadata: 55,
    });

    const makerOrder2 = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 201,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: 10_000_000_000,
      subticks: 2_000_000,
      goodTilOneof: {
        goodTilBlock: 100,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: false,
      clientMetadata: 56,
    });

    // create initial PerpetualPositions with closed previous positions
    await Promise.all([
      // previous position for subaccount 1
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId,
      }),
      // previous position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
        subaccountId: testConstants.defaultSubaccountId2,
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId,
      }),
      // initial position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
      }),
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
        subaccountId: testConstants.defaultSubaccountId2,
      }),
    ]);

    // First suborder fill event
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const orderFillEvent: OrderFillEventV1 = createOrderFillEvent(
      makerOrder1,
      suborder1,
      10_000_000_000,
      10_000_000_000,
      10_000_000_000,
    );
    const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
      orderFillEvent,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10),
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await onMessage(kafkaMessage);

    // After first fill, parent TWAP totalFilled should be 1000000
    await expectOrderInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '123',
      size: '100000000', // because the size is being set by the fill message
      totalFilled: '100000000', // 10
      price: '0.00000000000001',
      status: OrderStatus.FILLED,
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(suborder1.side),
      orderFlags: parentTwapOrder.orderId!.orderFlags.toString(),
      timeInForce: TimeInForce.IOC,
      reduceOnly: false,
      goodTilBlock: protocolTranslations.getGoodTilBlock(parentTwapOrder)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(parentTwapOrder),
      clientMetadata: parentTwapOrder.clientMetadata.toString(),
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
      orderType: OrderType.TWAP,
    });

    // // Second suborder fill event
    const orderFillEvent2: OrderFillEventV1 = createOrderFillEvent(
      makerOrder2,
      suborder2,
      10_000_000_000,
      10_000_000_000,
      10_000_000_000,
    );

    const kafkaMessage2: KafkaMessage = createKafkaMessageFromOrderFillEvent({
      orderFillEvent: orderFillEvent2,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10) + 1,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await onMessage(kafkaMessage2);

    await expectOrderInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '123',
      size: '100000000', // because the size is being set by the fill message
      totalFilled: '200000000', // 20
      price: '0.000000000000015',
      status: OrderStatus.FILLED, // orderSize > totalFilled so status is open
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(suborder1.side),
      orderFlags: parentTwapOrder.orderId!.orderFlags.toString(),
      timeInForce: TimeInForce.IOC,
      reduceOnly: false,
      goodTilBlock: protocolTranslations.getGoodTilBlock(parentTwapOrder)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(parentTwapOrder),
      clientMetadata: parentTwapOrder.clientMetadata.toString(),
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: '4',
      orderType: OrderType.TWAP,
    });

    const eventId: Buffer = TendermintEventTable.createEventId(
      defaultHeight,
      transactionIndex,
      eventIndex,
    );

    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '123',
      liquidity: Liquidity.TAKER,
      size: '100000000',
      price: '0.00000000000001',
      quoteAmount: '0.000001',
      eventId,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight,
      type: FillType.TWAP_SUBORDER,
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(suborder2.side),
      orderFlags: parentTwapOrder.orderId!.orderFlags.toString(),
      clientMetadata: suborder1.clientMetadata.toString(),
      fee: defaultTakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      positionSideBefore: 'LONG',
      entryPriceBefore: '15000',
      positionSizeBefore: '10',
    });

    const eventId2: Buffer = TendermintEventTable.createEventId(
      '4',
      transactionIndex,
      eventIndex,
    );

    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '123',
      liquidity: Liquidity.TAKER,
      size: '100000000',
      price: '0.00000000000002',
      quoteAmount: '0.000002',
      eventId: eventId2,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: '4',
      type: FillType.TWAP_SUBORDER,
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(suborder2.side),
      orderFlags: parentTwapOrder.orderId!.orderFlags.toString(),
      clientMetadata: suborder2.clientMetadata.toString(),
      fee: defaultTakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      positionSideBefore: 'LONG',
      entryPriceBefore: '0.001499999850010015',
      positionSizeBefore: '100000010',
    });
  });

  it.each([
    [
      undefined, // no maker order
    ],
    [
      IndexerOrder.fromPartial({ // no orderId
        orderId: undefined,
        side: IndexerOrder_Side.SIDE_BUY,
        quantums: 1,
        subticks: 1,
        goodTilBlock: 10,
      }),
    ],
    [
      IndexerOrder.fromPartial({ // no subaccountId
        orderId: {
          clientId: 0,
          clobPairId: Number(defaultClobPairId),
        },
        side: IndexerOrder_Side.SIDE_BUY,
        quantums: 1,
        subticks: 1,
        goodTilBlock: 10,
      }),
    ],
    [
      createOrder({ // Unspecified Order_Side
        subaccountId: defaultSubaccountId,
        clientId: 0,
        side: IndexerOrder_Side.SIDE_UNSPECIFIED,
        quantums: 10_000_000_000,
        subticks: 1,
        goodTilOneof: {
          goodTilBlock: 10,
        },
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
        reduceOnly: true,
        clientMetadata: 0,
      }),
    ],
    [
      createOrder({ // Undefined goodTilOneof oneofKind
        subaccountId: defaultSubaccountId,
        clientId: 0,
        side: IndexerOrder_Side.SIDE_UNSPECIFIED,
        quantums: 10_000_000_000,
        subticks: 1,
        goodTilOneof: {},
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
        reduceOnly: true,
        clientMetadata: 0,
      }),
    ],
  ])('fillOrderEvent fails validation', async (makerOrderProto: IndexerOrder | undefined) => {
    const subticks: number = 1_000_000_000_000_000_000;
    const takerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: 10_000_000_000,
      subticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: defaultClobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: true,
      clientMetadata: 0,
    });

    const fillAmount: number = 1_000_000;
    const orderFillEvent: OrderFillEventV1 = createOrderFillEvent(
      makerOrderProto,
      takerOrderProto,
      fillAmount,
      fillAmount,
      fillAmount,
    );

    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
      orderFillEvent,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10),
      time: defaultTime,
      txHash: defaultTxHash,
    });
    const loggerCrit = jest.spyOn(logger, 'crit');
    await expect(onMessage(kafkaMessage)).rejects.toThrowError();

    expect(loggerCrit).toHaveBeenCalledWith(expect.objectContaining({
      at: 'onMessage#onMessage',
      message: 'Error: Unable to parse message, this must be due to a bug in V4 node',
    }));
    await expectNoCandles();
  });

  it('correctly sets status for short term IOC orders', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const makerQuantums: number = 100;
    const makerSubticks: number = 1_000_000;

    const makerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: makerQuantums,
      subticks: makerSubticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: false,
      clientMetadata: 0,
    });

    const takerSubticks: number = 150_000;
    const takerQuantums: number = 10;
    const takerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: takerQuantums,
      subticks: takerSubticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: true,
      clientMetadata: 0,
    });

    const fillAmount: number = takerQuantums;
    const orderFillEvent: OrderFillEventV1 = createOrderFillEvent(
      makerOrderProto,
      takerOrderProto,
      fillAmount,
      fillAmount,
      fillAmount,
    );
    const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
      orderFillEvent,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10),
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await Promise.all([
      // initial position for subaccount 1
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
      }),
      // initial position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        subaccountId: testConstants.defaultSubaccountId2,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
      }),
    ]);

    await onMessage(kafkaMessage);

    const makerOrderId: string = OrderTable.orderIdToUuid(makerOrderProto.orderId!);
    const takerOrderId: string = OrderTable.orderIdToUuid(takerOrderProto.orderId!);

    const [makerOrder, takerOrder]: [
      OrderFromDatabase | undefined,
      OrderFromDatabase | undefined
    ] = await Promise.all([
      OrderTable.findById(makerOrderId),
      OrderTable.findById(takerOrderId),
    ]);

    expect(makerOrder).toBeDefined();
    expect(takerOrder).toBeDefined();

    // maker order is partially filled
    expect(makerOrder!.status).toEqual(OrderStatus.CANCELED);
    // taker order is fully filled
    expect(takerOrder!.status).toEqual(OrderStatus.FILLED);
  });

  it('correctly sets the maker and taker order router fees and addresses', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const makerQuantums: number = 100;
    const makerSubticks: number = 1_000_000;

    const makerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: makerQuantums,
      subticks: makerSubticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: false,
      clientMetadata: 0,
      orderRouterAddress: testConstants.defaultAddress,
    });

    const takerSubticks: number = 150_000;
    const takerQuantums: number = 10;
    const takerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: takerQuantums,
      subticks: takerSubticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
      reduceOnly: true,
      clientMetadata: 0,
      orderRouterAddress: testConstants.defaultAddress2,
    });

    const fillAmount: number = takerQuantums;
    const orderFillEvent: OrderFillEventV1 = createOrderFillEvent(
      makerOrderProto,
      takerOrderProto,
      fillAmount,
      fillAmount,
      fillAmount,
      0,
      0,
      testConstants.noBuilderAddress,
      testConstants.noBuilderAddress,
      10_000,
      20_000,
      testConstants.defaultAddress,
      testConstants.defaultAddress2,
    );
    const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
      orderFillEvent,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10),
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await Promise.all([
      // initial position for subaccount 1
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
      }),
      // initial position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        subaccountId: testConstants.defaultSubaccountId2,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
      }),
    ]);

    await onMessage(kafkaMessage);

    const makerOrderId: string = OrderTable.orderIdToUuid(makerOrderProto.orderId!);
    const takerOrderId: string = OrderTable.orderIdToUuid(takerOrderProto.orderId!);

    const [makerOrder, takerOrder]: [
      OrderFromDatabase | undefined,
      OrderFromDatabase | undefined
    ] = await Promise.all([
      OrderTable.findById(makerOrderId),
      OrderTable.findById(takerOrderId),
    ]);

    expect(makerOrder).toBeDefined();
    expect(takerOrder).toBeDefined();

    // maker order is partially filled
    expect(makerOrder!.status).toEqual(OrderStatus.CANCELED);

    // taker order is fully filled
    expect(takerOrder!.status).toEqual(OrderStatus.FILLED);

    expect(makerOrder!.orderRouterAddress).toEqual(testConstants.defaultAddress);

    expect(takerOrder!.orderRouterAddress).toEqual(testConstants.defaultAddress2);

    const eventId: Buffer = TendermintEventTable.createEventId(
      defaultHeight,
      transactionIndex,
      eventIndex,
    );

    const quoteAmount: string = '0.000000000000001'; // quote amount is price * fillAmount = 1e-14 * 1e-1 = 1e-15
    const price: string = '0.00000000000001';
    const totalFilledSize: string = (takerQuantums / makerQuantums).toString();

    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId2,
      clientId: '0',
      liquidity: Liquidity.TAKER,
      size: totalFilledSize,
      price,
      quoteAmount,
      eventId,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight,
      type: FillType.LIMIT,
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
      orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
      clientMetadata: takerOrderProto.clientMetadata.toString(),
      fee: defaultTakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      orderRouterAddress: testConstants.defaultAddress2,
      orderRouterFee: '0.02',
    });

    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '0',
      liquidity: Liquidity.MAKER,
      size: totalFilledSize,
      price,
      quoteAmount,
      eventId,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight,
      type: FillType.LIMIT,
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
      orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
      fee: defaultMakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      orderRouterAddress: testConstants.defaultAddress,
      orderRouterFee: '0.01',
    });
  });

  it.each([
    [
      'limit',
      IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
    ],
    [
      'post-only best effort canceled',
      IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY,
    ],
    [
      'post-only canceled',
      IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY,
      OrderStatus.CANCELED,
    ],
  ])('correctly sets status for short term %s orders', async (
    _orderType: string,
    timeInForce: IndexerOrder_TimeInForce,
    // either BEST_EFFORT_CANCELED or CANCELED
    status: OrderStatus = OrderStatus.BEST_EFFORT_CANCELED,
  ) => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const makerQuantums: number = 100;
    const makerSubticks: number = 1_000_000;

    const makerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: makerQuantums,
      subticks: makerSubticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce,
      reduceOnly: false,
      clientMetadata: 0,
    });

    const takerSubticks: number = 150_000;
    const takerQuantums: number = 100;
    const takerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: takerQuantums,
      subticks: takerSubticks,
      goodTilOneof: {
        goodTilBlock: 10,
      },
      clobPairId: testConstants.defaultPerpetualMarket3.clobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce,
      reduceOnly: true,
      clientMetadata: 0,
    });

    const makerOrderId: string = OrderTable.orderIdToUuid(makerOrderProto.orderId!);
    if (status === OrderStatus.BEST_EFFORT_CANCELED) {
      await CanceledOrdersCache.addBestEffortCanceledOrderId(makerOrderId, Date.now(), redisClient);
    } else { // Status is only over CANCELED or BEST_EFFORT_CANCELED
      await CanceledOrdersCache.addCanceledOrderId(makerOrderId, Date.now(), redisClient);
    }

    const fillAmount: number = 10;
    const orderFillEvent: OrderFillEventV1 = createOrderFillEvent(
      makerOrderProto,
      takerOrderProto,
      fillAmount,
      fillAmount,
      fillAmount,
    );
    const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
      orderFillEvent,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10),
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await Promise.all([
      // initial position for subaccount 1
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
      }),
      // initial position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        subaccountId: testConstants.defaultSubaccountId2,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
      }),
    ]);

    await onMessage(kafkaMessage);

    const takerOrderId: string = OrderTable.orderIdToUuid(takerOrderProto.orderId!);

    const [makerOrder, takerOrder]: [
      OrderFromDatabase | undefined,
      OrderFromDatabase | undefined
    ] = await Promise.all([
      OrderTable.findById(makerOrderId),
      OrderTable.findById(takerOrderId),
    ]);

    expect(makerOrder).toBeDefined();
    expect(takerOrder).toBeDefined();

    // maker order is partially filled, and in CanceledOrdersCache
    expect(makerOrder!.status).toEqual(status);
    // taker order is partially filled, and not in CanceledOrdersCache
    expect(takerOrder!.status).toEqual(OrderStatus.OPEN);
  });

  it('populates before position fields on fills for both maker and taker', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;

    const makerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: 1_000_000,
      subticks: 100_000_000,
      goodTilOneof: { goodTilBlock: 10 },
      clobPairId: defaultClobPairId,
      orderFlags: ORDER_FLAG_LONG_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL,
      reduceOnly: false,
      clientMetadata: 0,
    });

    const takerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: 10_000_000,
      subticks: 15_000_000,
      goodTilOneof: { goodTilBlock: 15 },
      clobPairId: defaultClobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
      reduceOnly: true,
      clientMetadata: 0,
    });

    await Promise.all([
      PerpetualPositionTable.create(defaultPerpetualPosition),
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        subaccountId: testConstants.defaultSubaccountId2,
      }),
    ]);

    const fillAmount = 1_000_000;
    const orderFillEvent = createOrderFillEvent(
      makerOrderProto,
      takerOrderProto,
      fillAmount,
      fillAmount,
      fillAmount,
    );
    const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
      orderFillEvent,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10),
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await onMessage(kafkaMessage);

    const eventId: Buffer = TendermintEventTable.createEventId(
      defaultHeight,
      transactionIndex,
      eventIndex,
    );

    const positionSizeBefore = '10';
    const entryPriceBefore = '15000';
    const positionSideBefore = 'LONG';

    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '0',
      liquidity: Liquidity.MAKER,
      size: '0.0001',
      price: '10000',
      quoteAmount: '1',
      eventId,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight,
      type: FillType.LIMIT,
      clobPairId: makerOrderProto.orderId!.clobPairId.toString(),
      side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
      orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
      fee: defaultMakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      positionSizeBefore,
      entryPriceBefore,
      positionSideBefore,
    });

    await expectFillInDatabase({
      subaccountId: testConstants.defaultSubaccountId2,
      clientId: '0',
      liquidity: Liquidity.TAKER,
      size: '0.0001',
      price: '10000',
      quoteAmount: '1',
      eventId,
      transactionHash: defaultTxHash,
      createdAt: defaultDateTime.toISO(),
      createdAtHeight: defaultHeight,
      type: FillType.LIMIT,
      clobPairId: defaultClobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(takerOrderProto.side),
      orderFlags: takerOrderProto.orderId!.orderFlags.toString(),
      clientMetadata: takerOrderProto.clientMetadata.toString(),
      fee: defaultTakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      positionSizeBefore,
      entryPriceBefore,
      positionSideBefore,
    });
  });

  it('updates totalRealizedPnl when a reduce-only order closes and realizes PnL', async () => {
    const transactionIndex = 0;
    const eventIndex = 0;

    await Promise.all([
      PerpetualPositionTable.create(defaultPerpetualPosition),
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        subaccountId: testConstants.defaultSubaccountId2,
      }),
    ]);

    const makerOrderProto = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: 1_000_000,
      subticks: 100_000_000,
      goodTilOneof: { goodTilBlock: 10 },
      clobPairId: defaultClobPairId,
      orderFlags: ORDER_FLAG_LONG_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL,
      reduceOnly: false,
      clientMetadata: 0,
    });

    const takerOrderProto = createOrder({
      subaccountId: defaultSubaccountId2,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_SELL,
      quantums: 1_000_000,
      subticks: 15_000_000,
      goodTilOneof: { goodTilBlock: 15 },
      clobPairId: defaultClobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
      reduceOnly: true,
      clientMetadata: 0,
    });

    const fillAmount = 1_000_000;
    const orderFillEvent = createOrderFillEvent(
      makerOrderProto,
      takerOrderProto,
      fillAmount,
      fillAmount,
      fillAmount,
    );

    const kafkaMessage: KafkaMessage = createKafkaMessageFromOrderFillEvent({
      orderFillEvent,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10),
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await onMessage(kafkaMessage);

    const expectedRealizedPnl = '-2.5000';
    const fillInHuman = '0.0001';
    const makerPriceHuman = '10000';

    await expectPerpetualPosition(
      PerpetualPositionTable.uuid(
        testConstants.defaultSubaccountId2,
        defaultPerpetualPosition.openEventId,
      ),
      {
        sumClose: fillInHuman,
        exitPrice: makerPriceHuman,
        totalRealizedPnl: expectedRealizedPnl,
      },
    );
  });

  async function expectDefaultOrderAndFillSubaccountKafkaMessages(
    producerSendMock: jest.SpyInstance,
    eventId: Buffer,
    makerOrderFlag: number,
    takerOrderFlag: number,
    perpetualId: string = testConstants.defaultPerpetualMarket.id,
    clobPairId: string = defaultClobPairId,
  ) {
    const positionId: string = (
      await PerpetualPositionTable.findOpenPositionForSubaccountPerpetual(
        testConstants.defaultSubaccountId,
        perpetualId,
      )
    )!.id;
    const positionId2: string = (
      await PerpetualPositionTable.findOpenPositionForSubaccountPerpetual(
        testConstants.defaultSubaccountId2,
        perpetualId,
      )
    )!.id;

    await Promise.all([
      expectOrderFillAndPositionSubaccountKafkaMessageFromIds(
        producerSendMock,
        defaultSubaccountId,
        OrderTable.uuid(
          testConstants.defaultSubaccountId,
          '0',
          clobPairId,
          makerOrderFlag.toString(),
        ),
        FillTable.uuid(eventId, Liquidity.MAKER),
        positionId,
      ),
      expectOrderFillAndPositionSubaccountKafkaMessageFromIds(
        producerSendMock,
        defaultSubaccountId2,
        OrderTable.uuid(
          testConstants.defaultSubaccountId2,
          '0',
          clobPairId,
          takerOrderFlag.toString(),
        ),
        FillTable.uuid(eventId, Liquidity.TAKER),
        positionId2,
      ),
    ]);
  }
});

function createOrderFillEvent(
  makerOrderProto: IndexerOrder | undefined,
  takerOrderProto: IndexerOrder,
  fillAmount: number,
  totalFilledMaker: number,
  totalFilledTaker: number,
  makerBuilderFee: number = 0,
  takerBuilderFee: number = 0,
  makerBuilderAddress: string = testConstants.noBuilderAddress,
  takerBuilderAddress: string = testConstants.noBuilderAddress,
  makerOrderRouterFee: number = 0,
  takerOrderRouterFee: number = 0,
  makerOrderRouterAddress: string = testConstants.noOrderRouterAddress,
  takerOrderRouterAddress: string = testConstants.noOrderRouterAddress,
): OrderFillEventV1 {
  return {
    makerOrder: makerOrderProto,
    order: takerOrderProto,
    fillAmount: Long.fromValue(fillAmount),
    makerFee: Long.fromValue(defaultMakerFeeQuantum, false),
    takerFee: Long.fromValue(defaultTakerFeeQuantum, false),
    totalFilledMaker: Long.fromValue(totalFilledMaker, true),
    totalFilledTaker: Long.fromValue(totalFilledTaker, true),
    affiliateRevShare: Long.fromValue(defaultAffiliateRevShareQuantum, false),
    makerBuilderFee: Long.fromValue(makerBuilderFee, false),
    takerBuilderFee: Long.fromValue(takerBuilderFee, false),
    makerBuilderAddress,
    takerBuilderAddress,
    makerOrderRouterFee: Long.fromValue(makerOrderRouterFee, true),
    takerOrderRouterFee: Long.fromValue(takerOrderRouterFee, true),
    makerOrderRouterAddress,
    takerOrderRouterAddress,
  } as OrderFillEventV1;
}

async function expectCandlesUpdated() {
  const candles: CandleFromDatabase[] = await CandleTable.findAll({}, []);
  expect(candles.length).toBeGreaterThan(0);
}

async function expectNoCandles() {
  const candles: CandleFromDatabase[] = await CandleTable.findAll({}, []);
  expect(candles.length).toEqual(0);
}
