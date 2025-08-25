import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock, IndexerTendermintEvent, Timestamp,
  LiquidationOrderV1,
  IndexerOrder,
  IndexerOrder_Side,
  OrderFillEventV1,
  IndexerSubaccountId,
  IndexerOrder_TimeInForce,
  OffChainUpdateV1,
  OrderRemovalReason, OrderRemoveV1_OrderRemovalStatus,
} from '@dydxprotocol-indexer/v4-protos';
import {
  assetRefresher,
  CandleFromDatabase,
  CandleTable,
  dbHelpers,
  FillTable,
  FillType,
  Liquidity,
  OrderCreateObject,
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
  createKafkaMessageFromOrderFillEvent,
  createLiquidationOrder,
  createOrder,
  expectFillInDatabase,
  expectFillSubaccountKafkaMessageFromLiquidationEvent,
  expectNoOrdersExistForSubaccountClobPairId,
  expectOrderFillAndPositionSubaccountKafkaMessageFromIds,
  expectOrderInDatabase,
  expectPerpetualPosition,
  expectDefaultTradeKafkaMessageFromTakerFillId,
  liquidationOrderToOrderSide,
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectVulcanKafkaMessage,
} from '../../helpers/indexer-proto-helpers';
import Big from 'big.js';
import { getWeightedAverage } from '../../../src/lib/helper';
import {
  ORDER_FLAG_SHORT_TERM,
  ORDER_FLAG_LONG_TERM,
  ORDER_FLAG_TWAP_SUBORDER,
  ORDER_FLAG_TWAP,
} from '@dydxprotocol-indexer/v4-proto-parser';
import { updateBlockCache } from '../../../src/caches/block-cache';
import { defaultLiquidation, defaultLiquidationEvent, defaultPreviousHeight } from '../../helpers/constants';
import { DydxIndexerSubtypes } from '../../../src/lib/types';
import { LiquidationHandler } from '../../../src/handlers/order-fills/liquidation-handler';
import { clearCandlesMap } from '../../../src/caches/candle-cache';
import Long from 'long';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';
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

describe('LiquidationHandler', () => {
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
    createdAtHeight: '1',
    openEventId: testConstants.defaultTendermintEventId4,
    lastEventId: testConstants.defaultTendermintEventId4,
    settledFunding: '200000',
  };

  describe('getParallelizationIds', () => {
    it.each([
      [
        'maker',
        Liquidity.MAKER,
        OrderTable.orderIdToUuid(defaultLiquidationEvent.makerOrder!.orderId!),
        defaultLiquidationEvent.makerOrder!.orderId!.subaccountId!,
      ],
      [
        'taker',
        Liquidity.TAKER,
        undefined,
        defaultSubaccountId,
      ],
    ])('returns the correct %s parallelization ids', (
      _name: string,
      liquidity: Liquidity,
      orderId: string | undefined,
      subaccountId: IndexerSubaccountId,
    ) => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.ORDER_FILL,
        Uint8Array.from(OrderFillEventV1.encode(defaultLiquidationEvent).finish()),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: LiquidationHandler = new LiquidationHandler(
        block,
        0,
        indexerTendermintEvent,
        0,
        {
          ...defaultLiquidation,
          liquidity,
        },
      );

      const parallelizationIds: string[] = [
        `${handler.eventType}_${SubaccountTable.subaccountIdToUuid(subaccountId)}_${defaultLiquidationEvent.makerOrder!.orderId!.clobPairId}`,
        `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${SubaccountTable.subaccountIdToUuid(subaccountId)}`,
      ];
      if (orderId !== undefined) {
        parallelizationIds.push(`${STATEFUL_ORDER_ORDER_FILL_EVENT_TYPE}_${orderId}`);
      }
      expect(handler.getParallelizationIds()).toEqual(parallelizationIds);
    });
  });

  it.each([
    [
      'goodTilBlock',
      {
        goodTilBlock: 10,
      },
    ],
    [
      'goodTilBlockTime',
      {
        goodTilBlockTime: 1_000_000_000,
      },
    ],
  ])(
    'creates fills and orders (with %s), sends vulcan message for maker order update and updates ' +
    'perpetualPosition',
    async (
      _name: string,
      goodTilOneof: Partial<IndexerOrder>,
    ) => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;
      const makerQuantums: number = 10_000_000;
      const makerSubticks: number = 100_000_000;

      const makerOrderProto: IndexerOrder = createOrder({
        subaccountId: defaultSubaccountId,
        clientId: 0,
        side: IndexerOrder_Side.SIDE_BUY,
        quantums: makerQuantums,
        subticks: makerSubticks,
        goodTilOneof,
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
        reduceOnly: true,
        clientMetadata: 0,
      });

      const takerSubticks: number = 15_000_000;
      const takerQuantums: number = 1_000_000;
      const liquidationOrder: LiquidationOrderV1 = createLiquidationOrder({
        subaccountId: defaultSubaccountId2,
        clobPairId: defaultClobPairId,
        perpetualId: defaultPerpetualPosition.perpetualId,
        quantums: takerQuantums,
        isBuy: false,
        subticks: takerSubticks,
      });

      const fillAmount: number = 1_000_000;
      const orderFillEvent: OrderFillEventV1 = createLiquidationOrderFillEvent(
        makerOrderProto,
        liquidationOrder,
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
          openEventId: testConstants.defaultTendermintEventId,
        }),
      ]);

      const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
      await onMessage(kafkaMessage);

      const makerOrderSize: string = '0.001'; // quantums in human = 1e7 * 1e-10 = 1e-3
      const makerPrice: string = '10000'; // quote currency / base currency = 1e8 * 1e-8 * 1e-6 / 1e-10 = 1e4
      const totalFilled: string = '0.0001'; // fillAmount in human = 1e6 * 1e-10 = 1e-4
      await expectOrderInDatabase({
        subaccountId: testConstants.defaultSubaccountId,
        clientId: '0',
        size: makerOrderSize,
        totalFilled,
        price: makerPrice,
        status: OrderStatus.OPEN, // orderSize > totalFilled so status is open
        clobPairId: defaultClobPairId,
        side: makerOrderProto.side === IndexerOrder_Side.SIDE_BUY ? OrderSide.BUY : OrderSide.SELL,
        orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
        timeInForce: TimeInForce.GTT,
        reduceOnly: true,
        goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
        updatedAt: defaultDateTime.toISO(),
        updatedAtHeight: defaultHeight.toString(),
      });

      // No orders should exist for the liquidated account since none are created, and there
      // are no orders to begin with.
      await expectNoOrdersExistForSubaccountClobPairId({
        subaccountId: testConstants.defaultSubaccountId2,
        clobPairId: defaultClobPairId,
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
        type: FillType.LIQUIDATION,
        clobPairId: makerOrderProto.orderId!.clobPairId.toString(),
        side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
        orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
        fee: defaultMakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
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
        type: FillType.LIQUIDATED,
        clobPairId: liquidationOrder.clobPairId.toString(),
        side: liquidationOrderToOrderSide(liquidationOrder),
        orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
        clientMetadata: null,
        fee: defaultTakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
        hasOrderId: false,
      });

      const expectedMakerOffchainUpdate: OffChainUpdateV1 = {
        orderUpdate: {
          orderId: makerOrderProto.orderId,
          totalFilledQuantums: orderFillEvent.totalFilledMaker,
        },
      };

      await Promise.all([
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: makerOrderProto.orderId!,
          offchainUpdate: expectedMakerOffchainUpdate,
        }),
        expectDefaultOrderFillAndPositionSubaccountKafkaMessages(
          producerSendMock,
          eventId,
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
        expectStateFilledQuantums(
          OrderTable.orderIdToUuid(makerOrderProto.orderId!),
          orderFillEvent.totalFilledMaker.toString(),
        ),
        expectCandlesUpdated(),
      ]);
    });

  it.each([
    [
      'goodTilBlock',
      {
        goodTilBlock: 10,
      },
      '5',
      undefined,
    ],
    [
      'goodTilBlockTime',
      {
        goodTilBlockTime: 1_000_000,
      },
      undefined,
      '1970-01-11T13:46:40.000Z',
    ],
  ])(
    'updates existing maker order (with %s), and sends vulcan message for maker order update and ' +
    'order removal for maker order fully-filled',
    async (
      _name: string,
      goodTilOneof: Partial<IndexerOrder>,
      existingGoodTilBlock?: string,
      existingGoodTilBlockTime?: string,
    ) => {

      // create initial orders
      const existingMakerOrder: OrderCreateObject = {
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
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
        goodTilBlock: existingGoodTilBlock,
        goodTilBlockTime: existingGoodTilBlockTime,
        clientMetadata: '0',
        updatedAt: DateTime.fromMillis(0).toISO(),
        updatedAtHeight: '0',
      };

      await Promise.all([
      // maker order
        OrderTable.create(existingMakerOrder),
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
      const makerQuantums: number = 1_001_000_000;
      const subticks: number = 100_000_000;
      const takerQuantums: number = 10_000_000_000;

      const makerOrderProto: IndexerOrder = createOrder({
        subaccountId: defaultSubaccountId,
        clientId: 0,
        side: IndexerOrder_Side.SIDE_BUY,
        quantums: makerQuantums,
        subticks,
        goodTilOneof,
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY,
        reduceOnly: true,
        clientMetadata: 0,
      });
      const liquidationOrder: LiquidationOrderV1 = createLiquidationOrder({
        subaccountId: defaultSubaccountId2,
        clobPairId: defaultClobPairId,
        perpetualId: defaultPerpetualPosition.perpetualId,
        quantums: takerQuantums,
        isBuy: false,
        subticks,
      });

      const fillAmount: number = 1_000_000;
      const makerTotalFilled: number = 1_001_000_000;
      const orderFillEvent: OrderFillEventV1 = createLiquidationOrderFillEvent(
        makerOrderProto,
        liquidationOrder,
        fillAmount,
        makerTotalFilled,
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

      const makerOrderSize: string = '0.1001'; // quantums in human = (1e9 + 1e6) * 1e-10 = 0.1001
      const price: string = '10000'; // quote currency / base currency = 1e8 * 1e-8 * 1e-6 / 1e-10 = 1e4
      const totalMakerOrderFilled: string = '0.1001';
      await expectOrderInDatabase({
        subaccountId: testConstants.defaultSubaccountId,
        clientId: '0',
        size: makerOrderSize,
        totalFilled: totalMakerOrderFilled,
        price,
        status: OrderStatus.FILLED, // orderSize == totalFilled so status is open
        clobPairId: defaultClobPairId,
        side: makerOrderProto.side === IndexerOrder_Side.SIDE_BUY ? OrderSide.BUY : OrderSide.SELL,
        orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
        timeInForce: TimeInForce.POST_ONLY,
        reduceOnly: true,
        goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
        updatedAt: defaultDateTime.toISO(),
        updatedAtHeight: defaultHeight.toString(),
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
        type: FillType.LIQUIDATION,
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
        fee: defaultMakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
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
        type: FillType.LIQUIDATED,
        clobPairId: liquidationOrder.clobPairId.toString(),
        side: liquidationOrderToOrderSide(liquidationOrder),
        orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
        clientMetadata: null,
        fee: defaultTakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
        hasOrderId: false,
      });

      const expectedMakerUpdateOffchainUpdate: OffChainUpdateV1 = {
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

      await Promise.all([
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: makerOrderProto.orderId!,
          offchainUpdate: expectedMakerUpdateOffchainUpdate,
        }),
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: makerOrderProto.orderId!,
          offchainUpdate: expectedMakerRemoveOffchainUpdate,
        }),
        expectDefaultOrderFillAndPositionSubaccountKafkaMessages(
          producerSendMock,
          eventId,
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
      ]);
    });

  it.each([
    [
      'goodTilBlock',
      {
        goodTilBlock: 10,
      },
      '5',
      undefined,
    ],
    [
      'goodTilBlockTime',
      {
        goodTilBlockTime: 1_000_000,
      },
      undefined,
      '1970-01-11T13:46:40.000Z',
    ],
  ])(
    'replaces existing maker order (with %s), upserting a new order with the same order id',
    async (
      _name: string,
      goodTilOneof: Partial<IndexerOrder>,
      existingGoodTilBlock?: string,
      existingGoodTilBlockTime?: string,
    ) => {

      // create initial orders
      const existingMakerOrder: OrderCreateObject = {
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
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
        goodTilBlock: existingGoodTilBlock,
        goodTilBlockTime: existingGoodTilBlockTime,
        clientMetadata: '0',
        updatedAt: DateTime.fromMillis(0).toISO(),
        updatedAtHeight: '0',
      };

      await Promise.all([
        // maker order
        OrderTable.create(existingMakerOrder),
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
      const newMakerQuantums: number = 2_001_000_000;
      const newSubticks: number = 200_000_000;
      const takerQuantums: number = 10_000_000_000;

      const makerOrderProto: IndexerOrder = createOrder({
        subaccountId: defaultSubaccountId,
        clientId: 0,
        side: IndexerOrder_Side.SIDE_SELL,
        quantums: newMakerQuantums,
        subticks: newSubticks,
        goodTilOneof,
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY,
        reduceOnly: true,
        clientMetadata: 0,
      });
      const liquidationOrder: LiquidationOrderV1 = createLiquidationOrder({
        subaccountId: defaultSubaccountId2,
        clobPairId: defaultClobPairId,
        perpetualId: defaultPerpetualPosition.perpetualId,
        quantums: takerQuantums,
        isBuy: true,
        subticks: newSubticks,
      });

      const fillAmount: number = 1_000_000;
      const makerTotalFilled: number = 1_000_000;
      const orderFillEvent: OrderFillEventV1 = createLiquidationOrderFillEvent(
        makerOrderProto,
        liquidationOrder,
        fillAmount,
        makerTotalFilled,
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

      const makerOrderSize: string = '0.2001'; // quantums in human = (2e9 + 1e6) * 1e-10 = 0.1001
      const price: string = '20000'; // quote currency / base currency = 2e8 * 1e-8 * 1e-6 / 1e-10 = 1e4
      const totalMakerOrderFilled: string = '0.0001';
      await expectOrderInDatabase({
        subaccountId: testConstants.defaultSubaccountId,
        clientId: '0',
        size: makerOrderSize,
        totalFilled: totalMakerOrderFilled,
        price,
        status: OrderStatus.OPEN, // orderSize > totalFilled so status is open
        clobPairId: defaultClobPairId,
        side: makerOrderProto.side === IndexerOrder_Side.SIDE_BUY ? OrderSide.BUY : OrderSide.SELL,
        orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
        timeInForce: TimeInForce.POST_ONLY,
        reduceOnly: true,
        goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
        goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
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
        type: FillType.LIQUIDATION,
        clobPairId: defaultClobPairId,
        side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
        orderFlags: ORDER_FLAG_LONG_TERM.toString(),
        clientMetadata: makerOrderProto.clientMetadata.toString(),
        fee: defaultMakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
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
        type: FillType.LIQUIDATED,
        clobPairId: liquidationOrder.clobPairId.toString(),
        side: liquidationOrderToOrderSide(liquidationOrder),
        orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
        clientMetadata: null,
        fee: defaultTakerFee,
        affiliateRevShare: defaultAffiliateRevShare,
        hasOrderId: false,
      });

      const expectedMakerUpdateOffchainUpdate: OffChainUpdateV1 = {
        orderUpdate: {
          orderId: makerOrderProto.orderId,
          totalFilledQuantums: orderFillEvent.totalFilledMaker,
        },
      };

      await Promise.all([
        expectVulcanKafkaMessage({
          producerSendMock,
          orderId: makerOrderProto.orderId!,
          offchainUpdate: expectedMakerUpdateOffchainUpdate,
        }),
        expectDefaultOrderFillAndPositionSubaccountKafkaMessages(
          producerSendMock,
          eventId,
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
      goodTilOneof: { goodTilBlock: 10 },
      clobPairId: defaultClobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
      reduceOnly: false,
      clientMetadata: 0,
    });

    const takerSubticks: number = 150_000;
    const takerQuantums: number = 10;
    const liquidationOrder: LiquidationOrderV1 = createLiquidationOrder({
      subaccountId: defaultSubaccountId2,
      clobPairId: defaultClobPairId,
      perpetualId: defaultPerpetualPosition.perpetualId,
      quantums: takerQuantums,
      isBuy: false,
      subticks: takerSubticks,
    });

    const fillAmount: number = 10;
    const orderFillEvent: OrderFillEventV1 = createLiquidationOrderFillEvent(
      makerOrderProto,
      liquidationOrder,
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

    // create initial PerpetualPositions
    await Promise.all([
      PerpetualPositionTable.create(defaultPerpetualPosition),
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        subaccountId: testConstants.defaultSubaccountId2,
      }),
    ]);

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    // This size should be in fixed-point notation rather than exponential notation (1e-8)
    const makerOrderSize: string = '0.00000001'; // quantums in human = 1e2 * 1e-10 = 1e-8
    const makerPrice: string = '100'; // quote currency / base currency = 1e6 * 1e-8 * 1e-6 / 1e-10 = 1e2
    const totalFilled: string = '0.000000001'; // fillAmount in human = 1e1 * 1e-10 = 1e-9
    await expectOrderInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '0',
      size: makerOrderSize,
      totalFilled,
      price: makerPrice,
      status: OrderStatus.OPEN, // orderSize > totalFilled so status is open
      clobPairId: defaultClobPairId,
      side: makerOrderProto.side === IndexerOrder_Side.SIDE_BUY ? OrderSide.BUY : OrderSide.SELL,
      orderFlags: makerOrderProto.orderId!.orderFlags.toString(),
      timeInForce: TimeInForce.GTT,
      reduceOnly: false,
      goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
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
      type: FillType.LIQUIDATION,
      clobPairId: defaultClobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
      fee: defaultMakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
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
      type: FillType.LIQUIDATED,
      clobPairId: defaultClobPairId,
      side: liquidationOrderToOrderSide(liquidationOrder),
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      clientMetadata: null,
      fee: defaultTakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      hasOrderId: false,
    });

    await Promise.all([
      expectDefaultOrderFillAndPositionSubaccountKafkaMessages(
        producerSendMock,
        eventId,
        ORDER_FLAG_SHORT_TERM,
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
      expectTimingStats(),
    ]);
  });

  it('LiquidationOrderFillEvent for TWAP suborder handles fills correctly', async () => {
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
      goodTilOneof: { goodTilBlock: 10 },
      clobPairId: defaultClobPairId,
      orderFlags: ORDER_FLAG_TWAP_SUBORDER.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED,
      reduceOnly: false,
      clientMetadata: 0,
    });

    const takerSubticks: number = 150_000;
    const takerQuantums: number = 10;
    const liquidationOrder: LiquidationOrderV1 = createLiquidationOrder({
      subaccountId: defaultSubaccountId2,
      clobPairId: defaultClobPairId,
      perpetualId: defaultPerpetualPosition.perpetualId,
      quantums: takerQuantums,
      isBuy: false,
      subticks: takerSubticks,
    });

    const fillAmount: number = 10;
    const orderFillEvent: OrderFillEventV1 = createLiquidationOrderFillEvent(
      makerOrderProto,
      liquidationOrder,
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

    // create initial PerpetualPositions
    await Promise.all([
      PerpetualPositionTable.create(defaultPerpetualPosition),
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        subaccountId: testConstants.defaultSubaccountId2,
      }),
    ]);

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const makerOrderSize: string = '0.00000001'; // quantums in human = 1e2 * 1e-10 = 1e-8
    const makerPrice: string = '100'; // quote currency / base currency = 1e6 * 1e-8 * 1e-6 / 1e-10 = 1e2
    const totalFilled: string = '0.000000001'; // fillAmount in human = 1e1 * 1e-10 = 1e-9
    await expectOrderInDatabase({
      subaccountId: testConstants.defaultSubaccountId,
      clientId: '0',
      size: makerOrderSize,
      totalFilled,
      price: makerPrice,
      status: OrderStatus.OPEN, // orderSize > totalFilled so status is open
      clobPairId: defaultClobPairId,
      side: makerOrderProto.side === IndexerOrder_Side.SIDE_BUY ? OrderSide.BUY : OrderSide.SELL,
      orderFlags: ORDER_FLAG_TWAP.toString(),
      timeInForce: TimeInForce.GTT,
      reduceOnly: false,
      goodTilBlock: protocolTranslations.getGoodTilBlock(makerOrderProto)?.toString(),
      goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(makerOrderProto),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: defaultHeight.toString(),
      orderType: OrderType.TWAP,
    });

    const eventId: Buffer = TendermintEventTable.createEventId(
      defaultHeight,
      transactionIndex,
      eventIndex,
    );

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
      type: FillType.LIQUIDATION,
      clobPairId: defaultClobPairId,
      side: protocolTranslations.protocolOrderSideToOrderSide(makerOrderProto.side),
      orderFlags: ORDER_FLAG_TWAP.toString(),
      clientMetadata: makerOrderProto.clientMetadata.toString(),
      fee: defaultMakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
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
      type: FillType.LIQUIDATED,
      clobPairId: defaultClobPairId,
      side: liquidationOrderToOrderSide(liquidationOrder),
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      clientMetadata: null,
      fee: defaultTakerFee,
      affiliateRevShare: defaultAffiliateRevShare,
      hasOrderId: false,
    });

    await Promise.all([
      expectDefaultOrderFillAndPositionSubaccountKafkaMessages(
        producerSendMock,
        eventId,
        ORDER_FLAG_TWAP,
      ),
      expectDefaultTradeKafkaMessageFromTakerFillId(
        producerSendMock,
        eventId,
      ),
      expectCandlesUpdated(),
      expectStateFilledQuantums(
        OrderTable.orderIdToUuid({
          ...makerOrderProto.orderId!,
          orderFlags: ORDER_FLAG_TWAP,
        }),
        orderFillEvent.totalFilledMaker.toString(),
      ),
      expectTimingStats(),
    ]);
  });

  it('LiquidationOrderFillEvent fails liquidationOrder validation', async () => {
    const makerQuantums: number = 10_000_000;
    const makerSubticks: number = 100_000_000;

    const makerOrderProto: IndexerOrder = createOrder({
      subaccountId: defaultSubaccountId,
      clientId: 0,
      side: IndexerOrder_Side.SIDE_BUY,
      quantums: makerQuantums,
      subticks: makerSubticks,
      goodTilOneof: { goodTilBlock: 10 },
      clobPairId: defaultClobPairId,
      orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
      timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL,
      reduceOnly: true,
      clientMetadata: 0,
    });

    const liquidationOrder: LiquidationOrderV1 = LiquidationOrderV1
      .fromPartial({ // no liquidated subaccount
        liquidated: undefined,
        clobPairId: 0,
        perpetualId: 0,
        totalSize: 1,
        subticks: 1,
        isBuy: false,
      });

    const fillAmount: number = 1_000_000;
    const orderFillEvent: OrderFillEventV1 = createLiquidationOrderFillEvent(
      makerOrderProto,
      liquidationOrder,
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

  it.each([
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
        goodTilOneof: { goodTilBlock: 10 },
        clobPairId: defaultClobPairId,
        orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
        timeInForce: IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC,
        reduceOnly: false,
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
        reduceOnly: false,
        clientMetadata: 0,
      }),
    ],
  ])('LiquidationOrderFillEvent fails makerOrder validation', async (
    makerOrderProto: IndexerOrder,
  ) => {
    const subticks: number = 1_000_000_000_000_000_000;
    const quantums: number = 1_000;
    const liquidationOrder: LiquidationOrderV1 = createLiquidationOrder({
      subaccountId: defaultSubaccountId2,
      clobPairId: defaultClobPairId,
      perpetualId: defaultPerpetualPosition.perpetualId,
      quantums,
      isBuy: false,
      subticks,
    });

    const fillAmount: number = 1_000_000;
    const orderFillEvent: OrderFillEventV1 = createLiquidationOrderFillEvent(
      makerOrderProto,
      liquidationOrder,
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

  async function expectDefaultOrderFillAndPositionSubaccountKafkaMessages(
    producerSendMock: jest.SpyInstance,
    eventId: Buffer,
    makerOrderFlag: number,
  ) {
    const positionId: string = (
      await PerpetualPositionTable.findOpenPositionForSubaccountPerpetual(
        testConstants.defaultSubaccountId,
        testConstants.defaultPerpetualMarket.id,
      )
    )!.id;
    const positionId2: string = (
      await PerpetualPositionTable.findOpenPositionForSubaccountPerpetual(
        testConstants.defaultSubaccountId2,
        testConstants.defaultPerpetualMarket.id,
      )
    )!.id;

    await Promise.all([
      expectOrderFillAndPositionSubaccountKafkaMessageFromIds(
        producerSendMock,
        defaultSubaccountId,
        OrderTable.uuid(
          testConstants.defaultSubaccountId,
          '0',
          defaultClobPairId,
          makerOrderFlag.toString(),
        ),
        FillTable.uuid(eventId, Liquidity.MAKER),
        positionId,
      ),
      expectFillSubaccountKafkaMessageFromLiquidationEvent(
        producerSendMock,
        defaultSubaccountId2,
        FillTable.uuid(eventId, Liquidity.TAKER),
        positionId2,
      ),
    ]);
  }
});

function createLiquidationOrderFillEvent(
  makerOrderProto: IndexerOrder,
  liquidationOrderProto: LiquidationOrderV1,
  fillAmount: number,
  totalFilledMaker: number,
): OrderFillEventV1 {
  return {
    makerOrder: makerOrderProto,
    liquidationOrder: liquidationOrderProto,
    fillAmount: Long.fromValue(fillAmount, true),
    makerFee: Long.fromValue(defaultMakerFeeQuantum, false),
    takerFee: Long.fromValue(defaultTakerFeeQuantum, false),
    totalFilledMaker: Long.fromValue(totalFilledMaker, true),
    totalFilledTaker: Long.fromValue(fillAmount, true),
    affiliateRevShare: Long.fromValue(defaultAffiliateRevShareQuantum, false),
    makerBuilderFee: Long.fromValue(0, false),
    takerBuilderFee: Long.fromValue(0, false),
    makerBuilderAddress: testConstants.noBuilderAddress,
    takerBuilderAddress: testConstants.noBuilderAddress,
    makerOrderRouterFee: Long.fromValue(0, true),
    takerOrderRouterFee: Long.fromValue(0, true),
    makerOrderRouterAddress: testConstants.noOrderRouterAddress,
    takerOrderRouterAddress: testConstants.noOrderRouterAddress,
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

function expectTimingStats() {
  expect(stats.timing).toHaveBeenCalledWith(
    'ender.handle_event.timing',
    expect.any(Number),
    {
      className: 'LiquidationHandler',
      eventType: 'LiquidationEvent',
    },
  );
}
