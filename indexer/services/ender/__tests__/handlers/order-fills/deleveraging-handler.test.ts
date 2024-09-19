import { logger, stats } from '@dydxprotocol-indexer/base';
import { redis } from '@dydxprotocol-indexer/redis';
import {
  assetRefresher,
  dbHelpers,
  FillTable,
  FillType,
  Liquidity,
  OrderSide,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualPositionCreateObject,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  PositionSide,
  SubaccountCreateObject,
  SubaccountTable,
  TendermintEventTable,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { updateBlockCache } from '../../../src/caches/block-cache';
import { defaultDeleveragingEvent, defaultPreviousHeight } from '../../helpers/constants';
import { clearCandlesMap } from '../../../src/caches/candle-cache';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';
import { redisClient } from '../../../src/helpers/redis/redis-controller';
import {
  DeleveragingEventV1,
  IndexerSubaccountId,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  Timestamp,
} from '@dydxprotocol-indexer/v4-protos';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  createKafkaMessageFromDeleveragingEvent,
  expectDefaultTradeKafkaMessageFromTakerFillId,
  expectFillInDatabase,
  expectFillSubaccountKafkaMessageFromLiquidationEvent, expectPerpetualPosition,
} from '../../helpers/indexer-proto-helpers';
import { DydxIndexerSubtypes } from '../../../src/lib/types';
import {
  MILLIS_IN_NANOS,
  SECONDS_IN_MILLIS,
  SUBACCOUNT_ORDER_FILL_EVENT_TYPE,
} from '../../../src/constants';
import { DateTime } from 'luxon';
import Long from 'long';
import { DeleveragingHandler } from '../../../src/handlers/order-fills/deleveraging-handler';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../../src/lib/on-message';
import { producer } from '@dydxprotocol-indexer/kafka';
import { createdDateTime, createdHeight } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import Big from 'big.js';

describe('DeleveragingHandler', () => {
  const offsettingSubaccount: SubaccountCreateObject = {
    address: defaultDeleveragingEvent.offsetting!.owner,
    subaccountNumber: defaultDeleveragingEvent.offsetting!.number,
    updatedAt: createdDateTime.toISO(),
    updatedAtHeight: createdHeight,
  };

  const deleveragedSubaccount: SubaccountCreateObject = {
    address: defaultDeleveragingEvent.liquidated!.owner,
    subaccountNumber: defaultDeleveragingEvent.liquidated!.number,
    updatedAt: createdDateTime.toISO(),
    updatedAtHeight: createdHeight,
  };

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
  const transactionIndex: number = 0;
  const eventIndex: number = 0;

  const offsettingPerpetualPosition: PerpetualPositionCreateObject = {
    subaccountId: SubaccountTable.subaccountIdToUuid(defaultDeleveragingEvent.offsetting!),
    perpetualId: testConstants.defaultPerpetualMarket2.id,
    side: PositionSide.LONG,
    status: PerpetualPositionStatus.OPEN,
    size: '10',
    maxSize: '25',
    sumOpen: '10',
    entryPrice: '15000',
    createdAt: DateTime.utc().toISO(),
    createdAtHeight: '1',
    openEventId: testConstants.defaultTendermintEventId,
    lastEventId: testConstants.defaultTendermintEventId,
    settledFunding: '200000',
  };
  const deleveragedPerpetualPosition: PerpetualPositionCreateObject = {
    ...offsettingPerpetualPosition,
    subaccountId: SubaccountTable.subaccountIdToUuid(defaultDeleveragingEvent.liquidated!),
    side: PositionSide.SHORT,
  };

  it('getParallelizationIds', () => {
    const offsettingSubaccountId: IndexerSubaccountId = defaultDeleveragingEvent.offsetting!;
    const deleveragedSubaccountId: IndexerSubaccountId = defaultDeleveragingEvent.liquidated!;

    const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
      DydxIndexerSubtypes.DELEVERAGING,
      DeleveragingEventV1.encode(defaultDeleveragingEvent).finish(),
      transactionIndex,
      eventIndex,
    );
    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      0,
      defaultTime,
      [indexerTendermintEvent],
      [defaultTxHash],
    );

    const handler: DeleveragingHandler = new DeleveragingHandler(
      block,
      0,
      indexerTendermintEvent,
      0,
      defaultDeleveragingEvent,
    );

    const offsettingSubaccountUuid: string = SubaccountTable.subaccountIdToUuid(
      offsettingSubaccountId,
    );
    const deleveragedSubaccountUuid: string = SubaccountTable.subaccountIdToUuid(
      deleveragedSubaccountId,
    );

    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromId(
        defaultDeleveragingEvent.perpetualId.toString(),
      );
    expect(perpetualMarket).toBeDefined();

    expect(handler.getParallelizationIds()).toEqual([
      `${handler.eventType}_${offsettingSubaccountUuid}_${perpetualMarket!.clobPairId}`,
      `${handler.eventType}_${deleveragedSubaccountUuid}_${perpetualMarket!.clobPairId}`,
      // To ensure that SubaccountUpdateEvents, OrderFillEvents, and DeleveragingEvents for
      // the same subaccount are not processed in parallel
      `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${offsettingSubaccountUuid}`,
      `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${deleveragedSubaccountUuid}`,
    ]);
  });

  it('DeleveragingEvent fails validation', async () => {
    const deleveragingEvent: DeleveragingEventV1 = DeleveragingEventV1
      .fromPartial({ // no liquidated subaccount
        ...defaultDeleveragingEvent,
        liquidated: undefined,
      });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromDeleveragingEvent({
      deleveragingEvent,
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
  });

  it('creates fills and updates perpetual positions', async () => {
    const kafkaMessage: KafkaMessage = createKafkaMessageFromDeleveragingEvent({
      deleveragingEvent: defaultDeleveragingEvent,
      transactionIndex,
      eventIndex,
      height: parseInt(defaultHeight, 10),
      time: defaultTime,
      txHash: defaultTxHash,
    });

    // create initial Subaccounts
    await Promise.all([
      SubaccountTable.create(offsettingSubaccount),
      SubaccountTable.create(deleveragedSubaccount),
    ]);
    // create initial PerpetualPositions
    await Promise.all([
      PerpetualPositionTable.create(offsettingPerpetualPosition),
      PerpetualPositionTable.create(deleveragedPerpetualPosition),
    ]);

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const eventId: Buffer = TendermintEventTable.createEventId(
      defaultHeight,
      transactionIndex,
      eventIndex,
    );

    // This size should be in fixed-point notation rather than exponential notation.
    const quoteAmount: string = '1000'; // quote amount is event->price * QUOTE_CURRENCY_ATOMIC_RESOLUTION = 1e9*1e-6=1e3
    const totalFilled: string = '0.00000000000001'; // fillAmount in human = 10^4*10^-18=10^-14
    const price: string = '100000000000000000'; // 1e3/1e-14=1e17
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromId(
        defaultDeleveragingEvent.perpetualId.toString(),
      );

    await Promise.all([
      expectFillInDatabase({
        subaccountId: SubaccountTable.subaccountIdToUuid(defaultDeleveragingEvent.offsetting!),
        clientId: '0',
        liquidity: Liquidity.MAKER,
        size: totalFilled,
        price,
        quoteAmount,
        eventId,
        transactionHash: defaultTxHash,
        createdAt: defaultDateTime.toISO(),
        createdAtHeight: defaultHeight,
        type: FillType.OFFSETTING,
        clobPairId: perpetualMarket!.clobPairId,
        side: OrderSide.SELL,
        orderFlags: '0',
        clientMetadata: null,
        hasOrderId: false,
        fee: '0',
        affiliateRevShare: '0',
      }),
      expectFillInDatabase({
        subaccountId: SubaccountTable.subaccountIdToUuid(defaultDeleveragingEvent.liquidated!),
        clientId: '0',
        liquidity: Liquidity.TAKER,
        size: totalFilled,
        price,
        quoteAmount,
        eventId,
        transactionHash: defaultTxHash,
        createdAt: defaultDateTime.toISO(),
        createdAtHeight: defaultHeight,
        type: FillType.DELEVERAGED,
        clobPairId: perpetualMarket!.clobPairId,
        side: OrderSide.BUY,
        orderFlags: '0',
        clientMetadata: null,
        hasOrderId: false,
        fee: '0',
        affiliateRevShare: '0',
      }),
      expectPerpetualPosition(
        PerpetualPositionTable.uuid(
          offsettingPerpetualPosition.subaccountId,
          offsettingPerpetualPosition.openEventId,
        ),
        {
          sumClose: Big(totalFilled).toFixed(),
          exitPrice: price,
        },
      ),
      expectPerpetualPosition(
        PerpetualPositionTable.uuid(
          deleveragedPerpetualPosition.subaccountId,
          deleveragedPerpetualPosition.openEventId,
        ),
        {
          sumClose: Big(totalFilled).toFixed(),
          exitPrice: price,
        },
      ),
    ]);

    await Promise.all([
      expectFillsAndPositionsSubaccountKafkaMessages(
        producerSendMock,
        eventId,
        true,
      ),
      expectFillsAndPositionsSubaccountKafkaMessages(
        producerSendMock,
        eventId,
        false,
      ),
      expectDefaultTradeKafkaMessageFromTakerFillId(
        producerSendMock,
        eventId,
      ),
      expectTimingStats(),
    ]);
  });

  async function expectFillsAndPositionsSubaccountKafkaMessages(
    producerSendMock: jest.SpyInstance,
    eventId: Buffer,
    deleveraged: boolean,
  ) {
    const subaccountId: IndexerSubaccountId = deleveraged
      ? defaultDeleveragingEvent.liquidated! : defaultDeleveragingEvent.offsetting!;
    const liquidity: Liquidity = deleveraged ? Liquidity.TAKER : Liquidity.MAKER;
    const positionId: string = (
      await PerpetualPositionTable.findOpenPositionForSubaccountPerpetual(
        SubaccountTable.subaccountIdToUuid(subaccountId),
        testConstants.defaultPerpetualMarket2.id,
      )
    )!.id;

    await expectFillSubaccountKafkaMessageFromLiquidationEvent(
      producerSendMock,
      subaccountId,
      FillTable.uuid(eventId, liquidity),
      positionId,
      defaultHeight,
      transactionIndex,
      eventIndex,
      testConstants.defaultPerpetualMarket2.ticker,
    );
  }
});

function expectTimingStats() {
  expect(stats.timing).toHaveBeenCalledWith(
    'ender.handle_event.timing',
    expect.any(Number),
    {
      className: 'DeleveragingHandler',
      eventType: 'DeleveragingEvent',
    },
  );
}
