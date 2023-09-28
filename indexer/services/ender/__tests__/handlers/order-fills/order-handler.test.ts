import { logger, stats, STATS_FUNCTION_NAME } from '@dydxprotocol-indexer/base';
import {
  IndexerOrder,
  IndexerOrder_Side,
  IndexerOrder_TimeInForce,
  IndexerOrderId,
  IndexerSubaccountId,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  OrderFillEventV1,
  Timestamp,
  OffChainUpdateV1,
  OrderRemovalReason,
  OrderRemoveV1_OrderRemovalStatus,
} from '@dydxprotocol-indexer/v4-protos';
import { redis, CanceledOrdersCache } from '@dydxprotocol-indexer/redis';
import {
  assetRefresher,
  CandleFromDatabase,
  CandleTable,
  dbHelpers,
  FillTable,
  FillType,
  Liquidity,
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
import { KafkaMessage } from 'kafkajs';
import { DateTime } from 'luxon';
import {
  MILLIS_IN_NANOS,
  SECONDS_IN_MILLIS,
  STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE,
  SUBACCOUNT_ORDER_FILL_EVENT_TYPE,
} from '../../../src/constants';
import { producer } from '@dydxprotocol-indexer/kafka';
import { onMessage } from '../../../src/lib/on-message';
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
import Big from 'big.js';
import { getWeightedAverage } from '../../../src/lib/helper';
import { ORDER_FLAG_LONG_TERM, ORDER_FLAG_SHORT_TERM } from '@dydxprotocol-indexer/v4-proto-parser';
import { updateBlockCache } from '../../../src/caches/block-cache';
import {
  defaultOrder, defaultOrderEvent, defaultPreviousHeight, defaultTakerOrder,
} from '../../helpers/constants';
import { DydxIndexerSubtypes } from '../../../src/lib/types';
import { OrderHandler } from '../../../src/handlers/order-fills/order-handler';
import { clearCandlesMap } from '../../../src/caches/candle-cache';
import Long from 'long';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';
import config from '../../../src/config';
import { redisClient } from '../../../src/helpers/redis/redis-controller';

const defaultClobPairId: string = testConstants.defaultPerpetualMarket.clobPairId;
const defaultMakerFeeQuantum: number = 1_000_000;
const defaultTakerFeeQuantum: number = 2_000_000;
const defaultMakerFee: string = protocolTranslations.quantumsToHumanFixedString(
  defaultMakerFeeQuantum.toString(),
  testConstants.defaultAsset.atomicResolution,
);
const defaultTakerFee: string = protocolTranslations.quantumsToHumanFixedString(
  defaultTakerFeeQuantum.toString(),
  testConstants.defaultAsset.atomicResolution,
);

describe('OrderHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
    jest.spyOn(stats, 'gauge');
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
    openEventId: testConstants.defaultTendermintEventId,
    lastEventId: testConstants.defaultTendermintEventId,
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
      'goodTilBlock via knex',
      {
        goodTilBlock: 10,
      },
      {
        goodTilBlock: 15,
      },
      false,
    ],
    [
      'goodTilBlock via SQL function',
      {
        goodTilBlock: 10,
      },
      {
        goodTilBlock: 15,
      },
      true,
    ],
    [
      'goodTilBlockTime via knex',
      {
        goodTilBlockTime: 1_000_000_000,
      },
      {
        goodTilBlockTime: 1_000_005_000,
      },
      false,
    ],
    [
      'goodTilBlockTime via SQL function',
      {
        goodTilBlockTime: 1_000_000_000,
      },
      {
        goodTilBlockTime: 1_000_005_000,
      },
      false,
    ],
  ])(
    'creates fills and orders (with %s), sends vulcan messages for order updates and order ' +
    'removal for maker order fully filled, and updates perpetualPosition',
    async (
      _name: string,
      makerGoodTilOneof: Partial<IndexerOrder>,
      takerGoodTilOneof: Partial<IndexerOrder>,
      useSqlFunction: boolean,
    ) => {
      config.USE_ORDER_HANDLER_SQL_FUNCTION = useSqlFunction;
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
          createdAtHeight: '0',
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
        timeInForce: TimeInForce.IOC,
        reduceOnly: true,
        goodTilBlock: protocolTranslations.getGoodTilBlock(takerOrderProto)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(takerOrderProto),
        clientMetadata: takerOrderProto.clientMetadata.toString(),
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
          },
        ),
        expectCandlesUpdated(),
      ]);

      if (!useSqlFunction) {
        expectTimingStats();
      }
    });

  it.each([
    [
      'goodTilBlock via knex',
      {
        goodTilBlock: 10,
      },
      {
        goodTilBlock: 15,
      },
      false,
      false,
      '5',
      undefined,
    ],
    [
      'goodTilBlock via SQL function',
      {
        goodTilBlock: 10,
      },
      {
        goodTilBlock: 15,
      },
      true,
      false,
      '5',
      undefined,
    ],
    [
      'goodTilBlockTime via knex',
      {
        goodTilBlockTime: 1_000_000_000,
      },
      {
        goodTilBlockTime: 1_000_005_000,
      },
      false,
      false,
      undefined,
      '1970-01-11T13:46:40.000Z',
    ],
    [
      'goodTilBlockTime via SQL function',
      {
        goodTilBlockTime: 1_000_000_000,
      },
      {
        goodTilBlockTime: 1_000_005_000,
      },
      true,
      false,
      undefined,
      '1970-01-11T13:46:40.000Z',
    ],
    [
      'goodTilBlockTime w/ cancelled maker order via knex',
      {
        goodTilBlockTime: 1_000_000_000,
      },
      {
        goodTilBlockTime: 1_000_005_000,
      },
      false,
      true,
      undefined,
      '1970-01-11T13:46:40.000Z',
    ],
    [
      'goodTilBlockTime w/ cancelled maker order via SQL function',
      {
        goodTilBlockTime: 1_000_000_000,
      },
      {
        goodTilBlockTime: 1_000_005_000,
      },
      true,
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
      useSqlFunction: boolean,
      isOrderCanceled: boolean,
      existingGoodTilBlock?: string,
      existingGoodTilBlockTime?: string,
    ) => {
      config.USE_ORDER_HANDLER_SQL_FUNCTION = useSqlFunction;
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
          ? OrderStatus.BEST_EFFORT_CANCELED
          : OrderStatus.OPEN, // orderSize > totalFilled so status is open
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
        orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
        timeInForce: TimeInForce.GTT,
        reduceOnly: true,
        goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
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
      });

      const eventId: Buffer = TendermintEventTable.createEventId(
        defaultHeight,
        transactionIndex,
        eventIndex,
      );
      const quoteAmount: string = '1'; // quote amount is price * fillAmount = 1e4 * 1e-4 = 1
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
      ]);

      if (!useSqlFunction) {
        expectTimingStats();
      }
    });

  it.each([
    [
      'via knex',
      false,
    ],
    [
      'via SQL function',
      true,
    ],
  ])('creates fills and orders with fixed-point notation quoteAmount (%s)', async (
    _name: string,
    useSqlFunction: boolean,
  ) => {
    config.USE_ORDER_HANDLER_SQL_FUNCTION = useSqlFunction;
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
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL,
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
        createdAtHeight: '1',
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId2,
      }),
      // previous position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        subaccountId: testConstants.defaultSubaccountId2,
        createdAtHeight: '1',
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId2,
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
      timeInForce: TimeInForce.FOK,
      reduceOnly: false,
      goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
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
    ]);
  });

  it.each([
    [
      'via knex',
      false,
    ],
    [
      'via SQL function',
      true,
    ],
  ])('creates fills and orders with fixed-point notation price (%s)', async (
    _name: string,
    useSqlFunction: boolean,
  ) => {
    config.USE_ORDER_HANDLER_SQL_FUNCTION = useSqlFunction;
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
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL,
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
        createdAtHeight: '1',
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId2,
      }),
      // previous position for subaccount 2
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: testConstants.defaultPerpetualMarket3.id,
        subaccountId: testConstants.defaultSubaccountId2,
        createdAtHeight: '1',
        size: '0',
        status: PerpetualPositionStatus.CLOSED,
        openEventId: testConstants.defaultTendermintEventId2,
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
      timeInForce: TimeInForce.FOK,
      reduceOnly: false,
      goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
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
    ]);
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
): OrderFillEventV1 {
  return {
    makerOrder: makerOrderProto,
    order: takerOrderProto,
    fillAmount: Long.fromValue(fillAmount),
    makerFee: Long.fromValue(defaultMakerFeeQuantum, false),
    takerFee: Long.fromValue(defaultTakerFeeQuantum, false),
    totalFilledMaker: Long.fromValue(totalFilledMaker, true),
    totalFilledTaker: Long.fromValue(totalFilledTaker, true),
  } as OrderFillEventV1;
}

function expectTimingStats() {
  expectTimingStat('upsert_orders');
  expectTimingStat('create_fill');
  expectTimingStat('update_perpetual_position');
}

function expectTimingStat(fnName: string) {
  expect(stats.timing).toHaveBeenCalledWith(
    `ender.${STATS_FUNCTION_NAME}.timing`,
    expect.any(Number),
    { className: 'OrderHandler', eventType: 'OrderFillEvent', fnName },
  );

}

async function expectCandlesUpdated() {
  const candles: CandleFromDatabase[] = await CandleTable.findAll({}, []);
  expect(candles.length).toBeGreaterThan(0);
}

async function expectNoCandles() {
  const candles: CandleFromDatabase[] = await CandleTable.findAll({}, []);
  expect(candles.length).toEqual(0);
}
