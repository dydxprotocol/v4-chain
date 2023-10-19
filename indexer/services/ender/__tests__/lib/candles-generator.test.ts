import { stats } from '@dydxprotocol-indexer/base';
import { CANDLES_WEBSOCKET_MESSAGE_VERSION } from '@dydxprotocol-indexer/kafka';
import {
  CandlesMap, CandlesResolutionMap,
  CandleColumns,
  CandleCreateObject,
  CandleFromDatabase,
  CandleMessageContents,
  CandleResolution,
  CandleTable,
  dbHelpers,
  IsolationLevel,
  IsoString,
  perpetualMarketRefresher,
  PerpetualPositionTable,
  PROTO_TO_CANDLE_RESOLUTION,
  testConstants,
  testMocks,
  Transaction,
} from '@dydxprotocol-indexer/postgres';
import { CandleMessage, CandleMessage_Resolution } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import {
  clearCandlesMap, getCandlesMap, startCandleCache,
} from '../../src/caches/candle-cache';
import config from '../../src/config';
import { CandlesGenerator } from '../../src/lib/candles-generator';
import { KafkaPublisher } from '../../src/lib/kafka-publisher';
import { ConsolidatedKafkaEvent } from '../../src/lib/types';
import { defaultTradeContent, defaultTradeKafkaEvent } from '../helpers/constants';
import { contentToSingleTradeMessage, createConsolidatedKafkaEventFromTrade } from '../helpers/kafka-publisher-helpers';

describe('candleHelper', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
    jest.spyOn(stats, 'timing');
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    clearCandlesMap();
    jest.clearAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  const defaultPrice: string = defaultTradeContent.price;
  const defaultPrice2: string = '15000';
  const defaultCandle: CandleCreateObject = {
    startedAt: '',
    ticker: testConstants.defaultPerpetualMarket.ticker,
    resolution: CandleResolution.ONE_MINUTE,
    low: defaultPrice,
    high: defaultPrice2,
    open: defaultPrice,
    close: defaultPrice2,
    baseTokenVolume: Big(defaultTradeContent.size).times(2).toString(),
    usdVolume: Big(defaultTradeContent.size).times(defaultPrice).plus(
      Big(defaultTradeContent.size).times(defaultPrice2),
    ).toString(),
    trades: 2,
    startingOpenInterest: '0',
  };
  const startedAt: IsoString = CandlesGenerator.calculateNormalizedCandleStartTime(
    testConstants.createdDateTime,
    CandleResolution.ONE_MINUTE,
  ).toISO();
  const previousStartedAt: IsoString = CandlesGenerator.calculateNormalizedCandleStartTime(
    testConstants.createdDateTime.minus({ minutes: 1 }),
    CandleResolution.ONE_MINUTE,
  ).toISO();
  const lowPrice: string = '7500';
  const openPrice: string = '7000';
  const closePrice: string = '8000';
  const highPrice: string = '8500';
  const existingStartingOpenInterest: string = '200';
  const existingTrades: number = 4;

  const defaultTradeKafkaEvent2:
  ConsolidatedKafkaEvent = createConsolidatedKafkaEventFromTrade(
    contentToSingleTradeMessage(
      {
        ...defaultTradeContent,
        price: defaultPrice2,
      },
      testConstants.defaultPerpetualMarket.clobPairId,
    ),
  );

  it('successfully creates candles with no open positions', async () => {
    // Create publisher and add events
    const publisher: KafkaPublisher = new KafkaPublisher();
    publisher.addEvents([
      defaultTradeKafkaEvent,
      defaultTradeKafkaEvent2,
    ]);

    await runUpdateCandles(publisher);

    // Verify postgres is updated
    const expectedCandles: CandleFromDatabase[] = _.map(
      Object.values(CandleResolution),
      (resolution: CandleResolution) => {
        const currentStartedAt: IsoString = CandlesGenerator.calculateNormalizedCandleStartTime(
          testConstants.createdDateTime,
          resolution,
        ).toISO();

        return {
          ...defaultCandle,
          id: CandleTable.uuid(currentStartedAt, defaultCandle.ticker, resolution),
          startedAt: currentStartedAt,
          resolution,
        };
      },
    );
    await verifyCandlesInPostgres(expectedCandles);

    // Verify publisher contains candles
    verifyAllCandlesEqualsKafkaMessages(publisher, expectedCandles);

    await validateCandlesCache();
    expectTimingStats();
  });

  it('successfully creates first candles with open interest', async () => {
    // Create publisher and add events
    const publisher: KafkaPublisher = new KafkaPublisher();
    publisher.addEvents([
      defaultTradeKafkaEvent,
      defaultTradeKafkaEvent2,
    ]);

    // Create Perpetual Position to set open position
    const openInterest: string = '100';
    await createOpenPosition(openInterest);

    await runUpdateCandles(publisher);

    // Verify postgres is updated
    const expectedCandles: CandleFromDatabase[] = _.map(
      Object.values(CandleResolution),
      (resolution: CandleResolution) => {
        const currentStartedAt: IsoString = CandlesGenerator.calculateNormalizedCandleStartTime(
          testConstants.createdDateTime,
          resolution,
        ).toISO();

        return {
          ...defaultCandle,
          id: CandleTable.uuid(currentStartedAt, defaultCandle.ticker, resolution),
          startedAt: currentStartedAt,
          resolution,
          startingOpenInterest: openInterest,
        };
      },
    );
    await verifyCandlesInPostgres(expectedCandles);

    // Verify publisher contains candles
    verifyAllCandlesEqualsKafkaMessages(publisher, expectedCandles);

    await validateCandlesCache();
    expectTimingStats();
  });

  it('successfully updates existing candles', async () => {
    const existingPrice: string = '7000';
    const startingOpenInterest: string = '200';
    const baseTokenVolume: string = '10';
    const usdVolume: string = Big(existingPrice).times(baseTokenVolume).toString();
    await Promise.all(
      _.map(Object.values(CandleResolution), (resolution: CandleResolution) => {
        const currentStartedAt: IsoString = CandlesGenerator.calculateNormalizedCandleStartTime(
          testConstants.createdDateTime,
          resolution,
        ).toISO();

        return CandleTable.create({
          startedAt: currentStartedAt,
          ticker: testConstants.defaultPerpetualMarket.ticker,
          resolution,
          low: existingPrice,
          high: existingPrice,
          open: existingPrice,
          close: existingPrice,
          baseTokenVolume,
          usdVolume,
          trades: existingTrades,
          startingOpenInterest,
        });
      }),
    );
    await startCandleCache();

    const publisher: KafkaPublisher = new KafkaPublisher();
    publisher.addEvents([
      defaultTradeKafkaEvent,
      defaultTradeKafkaEvent2,
    ]);

    await runUpdateCandles(publisher);

    // Verify postgres is updated
    const expectedCandles: CandleFromDatabase[] = _.map(
      Object.values(CandleResolution),
      (resolution: CandleResolution) => {
        const currentStartedAt: IsoString = CandlesGenerator.calculateNormalizedCandleStartTime(
          testConstants.createdDateTime,
          resolution,
        ).toISO();

        return {
          id: CandleTable.uuid(currentStartedAt, defaultCandle.ticker, resolution),
          startedAt: currentStartedAt,
          ticker: defaultCandle.ticker,
          resolution,
          low: existingPrice,
          high: defaultCandle.high,
          open: existingPrice,
          close: defaultCandle.close,
          baseTokenVolume: Big(defaultCandle.baseTokenVolume).plus(baseTokenVolume).toString(),
          usdVolume: Big(defaultCandle.usdVolume).plus(usdVolume).toString(),
          trades: existingTrades + 2,
          startingOpenInterest,
        };
      },
    );
    await verifyCandlesInPostgres(expectedCandles);

    // Verify publisher contains candles
    verifyAllCandlesEqualsKafkaMessages(publisher, expectedCandles);

    await validateCandlesCache();
    expectTimingStats();
  });

  it.each([
    [
      'creates empty candles', // description
      { // initial candle
        startedAt: previousStartedAt,
        ticker: testConstants.defaultPerpetualMarket.ticker,
        resolution: CandleResolution.ONE_MINUTE,
        low: lowPrice,
        high: highPrice,
        open: openPrice,
        close: closePrice,
        baseTokenVolume: '10',
        usdVolume: '10000',
        trades: existingTrades,
        startingOpenInterest: existingStartingOpenInterest,
      },
      '100', // open interest
      false, // block contains trades
      { // expected candle
        id: CandleTable.uuid(startedAt, defaultCandle.ticker, CandleResolution.ONE_MINUTE),
        startedAt,
        ticker: testConstants.defaultPerpetualMarket.ticker,
        resolution: CandleResolution.ONE_MINUTE,
        low: closePrice,
        high: closePrice,
        open: closePrice,
        close: closePrice,
        baseTokenVolume: '0',
        usdVolume: '0',
        trades: 0,
        startingOpenInterest: '100',
      },
    ],
    [
      'creates new candle if existing candle is from a past normalized candle start time', // description
      { // initial candle
        startedAt: previousStartedAt,
        ticker: testConstants.defaultPerpetualMarket.ticker,
        resolution: CandleResolution.ONE_MINUTE,
        low: lowPrice,
        high: highPrice,
        open: openPrice,
        close: closePrice,
        baseTokenVolume: '10',
        usdVolume: '10000',
        trades: existingTrades,
        startingOpenInterest: existingStartingOpenInterest,
      },
      '100', // open interest
      true, // block contains trades
      { // expected candle
        ...defaultCandle,
        id: CandleTable.uuid(startedAt, defaultCandle.ticker, CandleResolution.ONE_MINUTE),
        startedAt,
        resolution: CandleResolution.ONE_MINUTE,
        startingOpenInterest: '100',
      },
    ],
    [
      'updates empty candle', // description
      { // initial candle
        startedAt,
        ticker: testConstants.defaultPerpetualMarket.ticker,
        resolution: CandleResolution.ONE_MINUTE,
        low: closePrice,
        high: closePrice,
        open: closePrice,
        close: closePrice,
        baseTokenVolume: '0',
        usdVolume: '0',
        trades: 0,
        startingOpenInterest: existingStartingOpenInterest,
      },
      '100', // open interest
      true, // block contains trades
      { // expected candle
        ...defaultCandle,
        id: CandleTable.uuid(startedAt, defaultCandle.ticker, CandleResolution.ONE_MINUTE),
        startedAt,
        resolution: CandleResolution.ONE_MINUTE,
        startingOpenInterest: existingStartingOpenInterest,
      },
    ],
    [
      'does nothing when there are no trades and no existing candle', // description
      undefined,
      '100', // open interest
      false, // block contains trades
      undefined,
    ],
    [
      'does not update candle when there are no trades and an existing candle', // description
      { // initial candle
        startedAt,
        ticker: testConstants.defaultPerpetualMarket.ticker,
        resolution: CandleResolution.ONE_MINUTE,
        low: lowPrice,
        high: highPrice,
        open: openPrice,
        close: closePrice,
        baseTokenVolume: '10',
        usdVolume: '10000',
        trades: existingTrades,
        startingOpenInterest: existingStartingOpenInterest,
      },
      '100', // open interest
      false, // block contains trades
      { // expected candle
        id: CandleTable.uuid(startedAt, defaultCandle.ticker, CandleResolution.ONE_MINUTE),
        startedAt,
        ticker: testConstants.defaultPerpetualMarket.ticker,
        resolution: CandleResolution.ONE_MINUTE,
        low: lowPrice,
        high: highPrice,
        open: openPrice,
        close: closePrice,
        baseTokenVolume: '10',
        usdVolume: '10000',
        trades: existingTrades,
        startingOpenInterest: existingStartingOpenInterest,
      },
      false, // contains kafka messages
    ],
  ])('Successfully %s', async (
    _description: string,
    initialCandle: CandleCreateObject | undefined,
    openInterest: string,
    blockContainsTrades: boolean,
    expectedCandle: CandleFromDatabase | undefined,
    containsKafkaMessages: boolean = true,
  ) => {
    if (initialCandle !== undefined) {
      await CandleTable.create(initialCandle);
    }
    await startCandleCache();

    if (openInterest !== '0') {
      await createOpenPosition(openInterest);
    }

    const publisher: KafkaPublisher = new KafkaPublisher();
    if (blockContainsTrades) {
      publisher.addEvents([
        defaultTradeKafkaEvent,
        defaultTradeKafkaEvent2,
      ]);
    }
    await runUpdateCandles(publisher);

    if (expectedCandle === undefined) {
      // Verify no candles in postgres and no kafka messages
      await verifyNoCandleInPostgres(CandleResolution.ONE_MINUTE, startedAt);
      verifyNoCandlesKafkaMessages(publisher, CandleResolution.ONE_MINUTE);
    } else {
      const expectedCandles: CandleFromDatabase[] = [expectedCandle];
      await verifyCandlesInPostgres(expectedCandles);
      if (containsKafkaMessages) {
        verifyContainsCandlesKafkaMessages(publisher, expectedCandles);
      }
    }

    await validateCandlesCache();
    expectTimingStats();
  });
});

async function createOpenPosition(
  size: string,
): Promise<void> {
  // Create Perpetual Position to set open position
  await PerpetualPositionTable.create({
    ...testConstants.defaultPerpetualPosition,
    size,
  });
}

async function runUpdateCandles(
  publisher: KafkaPublisher,
) {
  const txId: number = await Transaction.start();
  await Transaction.setIsolationLevel(txId, IsolationLevel.READ_UNCOMMITTED);
  const candlesGenerator: CandlesGenerator = new CandlesGenerator(
    publisher,
    testConstants.createdDateTime,
    txId,
  );
  await candlesGenerator.updateCandles();
  await Transaction.commit(txId);
}

async function verifyCandlesInPostgres(
  expectedCandles: CandleFromDatabase[],
): Promise<void> {
  const candles: CandleFromDatabase[] = await CandleTable.findAll({}, []);
  _.forEach(expectedCandles, (expectedCandle: CandleFromDatabase) => {
    expect(candles).toContainEqual(expectedCandle);
  });
}

async function verifyNoCandleInPostgres(
  resolution: CandleResolution,
  startedAt: IsoString,
): Promise<void> {
  const candle: CandleFromDatabase | undefined = await CandleTable.findById(
    CandleTable.uuid(
      startedAt,
      testConstants.defaultPerpetualMarket.ticker,
      resolution,
    ),
  );
  expect(candle).toBeUndefined();
}

function verifyNoCandlesKafkaMessages(
  publisher: KafkaPublisher,
  resolution: CandleResolution,
) {
  _.forEach(publisher.candleMessages, (candleMessage: CandleMessage) => {
    expect(candleMessage.clobPairId).toEqual(testConstants.defaultPerpetualMarket.clobPairId);
    expect(candleMessage.version).toEqual(CANDLES_WEBSOCKET_MESSAGE_VERSION);

    if (candleMessage.resolution !== CandleMessage_Resolution.UNRECOGNIZED) {
      expect(PROTO_TO_CANDLE_RESOLUTION[candleMessage.resolution]).not.toEqual(resolution);
    }
  });
}

function verifyAllCandlesEqualsKafkaMessages(
  publisher: KafkaPublisher,
  expectedCandles: CandleFromDatabase[],
) {
  const resolutionToExpectedContent:
  Partial<Record<CandleResolution, CandleMessageContents>> = _.chain(expectedCandles)
    .keyBy(CandleColumns.resolution)
    .mapValues((candle: CandleFromDatabase) => {
      return _.omit(candle, [CandleColumns.id]);
    })
    .value();

  _.forEach(publisher.candleMessages, (candleMessage: CandleMessage) => {
    expect(candleMessage.clobPairId).toEqual(testConstants.defaultPerpetualMarket.clobPairId);
    expect(candleMessage.version).toEqual(CANDLES_WEBSOCKET_MESSAGE_VERSION);

    if (candleMessage.resolution !== CandleMessage_Resolution.UNRECOGNIZED) {
      const resolution: CandleResolution = PROTO_TO_CANDLE_RESOLUTION[candleMessage.resolution];
      const expectedContent: CandleMessageContents = resolutionToExpectedContent[resolution]!;
      expect(expectedContent).toEqual(JSON.parse(candleMessage.contents));
    }
  });
}

/**
 * Verifies that candles kafka messages contain the expected candles
 */
function verifyContainsCandlesKafkaMessages(
  publisher: KafkaPublisher,
  expectedCandles: CandleFromDatabase[],
) {
  const resolutionToContent: Partial<Record<CandleResolution, CandleMessageContents>> = {};
  _.forEach(publisher.candleMessages, (candleMessage: CandleMessage) => {
    expect(candleMessage.clobPairId).toEqual(testConstants.defaultPerpetualMarket.clobPairId);
    expect(candleMessage.version).toEqual(CANDLES_WEBSOCKET_MESSAGE_VERSION);

    if (candleMessage.resolution !== CandleMessage_Resolution.UNRECOGNIZED) {
      const resolution: CandleResolution = PROTO_TO_CANDLE_RESOLUTION[candleMessage.resolution];
      resolutionToContent[resolution] = JSON.parse(candleMessage.contents);
    }
  });

  _.forEach(expectedCandles, (expectedCandle: CandleFromDatabase) => {
    expect(
      _.omit(expectedCandle, [CandleColumns.id]),
    ).toEqual(resolutionToContent[expectedCandle.resolution]);
  });
}

async function validateCandlesCache() {
  const candlesMap: CandlesMap = getCandlesMap();
  const promises: Promise<CandleFromDatabase | undefined >[] = [];
  _.forEach(candlesMap, (candlesResolutionMap: CandlesResolutionMap, _ticker: string) => {
    _.forEach(candlesResolutionMap, (candle: CandleFromDatabase, _resolution: string) => {
      promises.push(CandleTable.findById(candle.id));
    });
  });

  const candlesInPostgres: (CandleFromDatabase | undefined)[] = await Promise.all(promises);

  _.forEach(candlesInPostgres, (candle: CandleFromDatabase | undefined) => {
    expect(candle).toBeDefined();
  });
}

function expectTimingStats() {
  expectTimingStat('update_candles');
  expectTimingStat('update_postgres_candles');
}

function expectTimingStat(statName: string) {
  expect(stats.timing).toHaveBeenCalledWith(
    `${config.SERVICE_NAME}.${statName}.timing`,
    expect.any(Number),
  );
}
