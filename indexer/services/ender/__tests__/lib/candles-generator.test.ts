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
  helpers,
  OrderSide,
} from '@dydxprotocol-indexer/postgres';
import { CandleMessage, CandleMessage_Resolution } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import {
  clearCandlesMap, getCandlesMap, startCandleCache,
} from '../../src/caches/candle-cache';
import config from '../../src/config';
import { CandlesGenerator, getOrderbookMidPriceMap } from '../../src/lib/candles-generator';
import { KafkaPublisher } from '../../src/lib/kafka-publisher';
import { ConsolidatedKafkaEvent } from '../../src/lib/types';
import { defaultTradeContent, defaultTradeKafkaEvent } from '../helpers/constants';
import { contentToSingleTradeMessage, createConsolidatedKafkaEventFromTrade } from '../helpers/kafka-publisher-helpers';
import { updatePriceLevel } from '../helpers/redis-helpers';
import { redisClient } from '../../src/helpers/redis/redis-controller';
import { redis } from '@dydxprotocol-indexer/redis';

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
    await redis.deleteAllAsync(redisClient);
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
    orderbookMidPriceClose: undefined,
    orderbookMidPriceOpen: undefined,
  };
  const startedAt: IsoString = helpers.calculateNormalizedCandleStartTime(
    testConstants.createdDateTime,
    CandleResolution.ONE_MINUTE,
  ).toISO();
  const previousStartedAt: IsoString = helpers.calculateNormalizedCandleStartTime(
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

    // Create Orderbook levels to set orderbookMidPrice open & close
    await updatePriceLevel('BTC-USD', '100000', OrderSide.BUY);
    await updatePriceLevel('BTC-USD', '110000', OrderSide.SELL);

    await runUpdateCandles(publisher);

    // Verify postgres is updated
    const expectedCandles: CandleFromDatabase[] = _.map(
      Object.values(CandleResolution),
      (resolution: CandleResolution) => {
        const currentStartedAt: IsoString = helpers.calculateNormalizedCandleStartTime(
          testConstants.createdDateTime,
          resolution,
        ).toISO();

        return {
          ...defaultCandle,
          id: CandleTable.uuid(currentStartedAt, defaultCandle.ticker, resolution),
          startedAt: currentStartedAt,
          resolution,
          orderbookMidPriceClose: '105000',
          orderbookMidPriceOpen: '105000',
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

    await updatePriceLevel('BTC-USD', '80000', OrderSide.BUY);
    await updatePriceLevel('BTC-USD', '81000', OrderSide.SELL);

    // Create Perpetual Position to set open position
    const openInterest: string = '100';
    await createOpenPosition(openInterest);

    await runUpdateCandles(publisher);

    // Verify postgres is updated
    const expectedCandles: CandleFromDatabase[] = _.map(
      Object.values(CandleResolution),
      (resolution: CandleResolution) => {
        const currentStartedAt: IsoString = helpers.calculateNormalizedCandleStartTime(
          testConstants.createdDateTime,
          resolution,
        ).toISO();

        return {
          ...defaultCandle,
          id: CandleTable.uuid(currentStartedAt, defaultCandle.ticker, resolution),
          startedAt: currentStartedAt,
          resolution,
          startingOpenInterest: openInterest,
          orderbookMidPriceClose: '80500',
          orderbookMidPriceOpen: '80500',
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
    const orderbookMidPriceClose = '7500';
    const orderbookMidPriceOpen = '8000';
    await Promise.all(
      _.map(Object.values(CandleResolution), (resolution: CandleResolution) => {
        const currentStartedAt: IsoString = helpers.calculateNormalizedCandleStartTime(
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
          orderbookMidPriceClose,
          orderbookMidPriceOpen,
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
        const currentStartedAt: IsoString = helpers.calculateNormalizedCandleStartTime(
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
          orderbookMidPriceClose,
          orderbookMidPriceOpen,
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
      'creates empty candle', // description
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
        orderbookMidPriceClose: undefined,
        orderbookMidPriceOpen: undefined,
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
        orderbookMidPriceClose: '1000',
        orderbookMidPriceOpen: '1000',
      },
      true,
      1000,
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
        orderbookMidPriceClose: '3000',
        orderbookMidPriceOpen: '3500',
      },
      '100', // open interest
      true, // block contains trades
      { // expected candle
        ...defaultCandle,
        id: CandleTable.uuid(startedAt, defaultCandle.ticker, CandleResolution.ONE_MINUTE),
        startedAt,
        resolution: CandleResolution.ONE_MINUTE,
        startingOpenInterest: '100',
        orderbookMidPriceClose: '1000',
        orderbookMidPriceOpen: '1000',
      },
      true, // contains kafka messages
      1000, // orderbook mid price
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
        orderbookMidPriceClose: undefined,
        orderbookMidPriceOpen: undefined,
      },
      '100', // open interest
      true, // block contains trades
      { // expected candle
        ...defaultCandle,
        id: CandleTable.uuid(startedAt, defaultCandle.ticker, CandleResolution.ONE_MINUTE),
        startedAt,
        resolution: CandleResolution.ONE_MINUTE,
        startingOpenInterest: existingStartingOpenInterest,
        orderbookMidPriceClose: null,
        orderbookMidPriceOpen: null,
      },
      true, // contains kafka messages
      1000, // orderbook mid price
    ],
    [
      'does nothing when there are no trades and no existing candle', // description
      undefined, // initial candle
      '100', // open interest
      false, // block contains trades
      undefined, // expected candle
      true, // contains kafka messages
      1000, // orderbook mid price
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
        orderbookMidPriceClose: '5000',
        orderbookMidPriceOpen: '6000',
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
        orderbookMidPriceClose: '5000',
        orderbookMidPriceOpen: '6000',
      },
      false, // contains kafka messages
      1000,
    ],
  ])('Successfully %s', async (
    _description: string,
    initialCandle: CandleCreateObject | undefined,
    openInterest: string,
    blockContainsTrades: boolean,
    expectedCandle: CandleFromDatabase | undefined,
    containsKafkaMessages: boolean = true,
    orderbookMidPrice: number,
  ) => {
    const midPriceSpread = 10;
    await updatePriceLevel('BTC-USD', String(orderbookMidPrice + midPriceSpread), OrderSide.SELL);
    await updatePriceLevel('BTC-USD', String(orderbookMidPrice - midPriceSpread), OrderSide.BUY);

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

  it('Updates previous candle orderBookMidPriceClose if startTime is past candle resolution', async () => {
    // Create existing candles
    const existingPrice: string = '7000';
    const startingOpenInterest: string = '200';
    const baseTokenVolume: string = '10';
    const usdVolume: string = Big(existingPrice).times(baseTokenVolume).toString();
    const orderbookMidPriceClose = '7500';
    const orderbookMidPriceOpen = '8000';
    await Promise.all(
      _.map(Object.values(CandleResolution), (resolution: CandleResolution) => {
        return CandleTable.create({
          startedAt: previousStartedAt,
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
          orderbookMidPriceClose,
          orderbookMidPriceOpen,
        });
      }),
    );
    await startCandleCache();

    // Update Orderbook levels
    await updatePriceLevel('BTC-USD', '10010', OrderSide.SELL);
    await updatePriceLevel('BTC-USD', '10000', OrderSide.BUY);

    const publisher: KafkaPublisher = new KafkaPublisher();
    publisher.addEvents([
      defaultTradeKafkaEvent,
      defaultTradeKafkaEvent2,
    ]);

    // Create new candles, with trades
    await runUpdateCandles(publisher);

    // Verify previous candles have orderbookMidPriceClose updated
    const previousExpectedCandles: CandleFromDatabase[] = _.map(
      Object.values(CandleResolution),
      (resolution: CandleResolution) => {
        return {
          id: CandleTable.uuid(previousStartedAt, defaultCandle.ticker, resolution),
          startedAt: previousStartedAt,
          ticker: defaultCandle.ticker,
          resolution,
          low: existingPrice,
          high: existingPrice,
          open: existingPrice,
          close: existingPrice,
          baseTokenVolume,
          usdVolume,
          trades: existingTrades,
          startingOpenInterest,
          orderbookMidPriceClose: '10005',
          orderbookMidPriceOpen,
        };
      },
    );
    await verifyCandlesInPostgres(previousExpectedCandles);

    // Verify new candles were created
    const expectedCandles: CandleFromDatabase[] = _.map(
      Object.values(CandleResolution),
      (resolution: CandleResolution) => {
        const currentStartedAt: IsoString = helpers.calculateNormalizedCandleStartTime(
          testConstants.createdDateTime,
          resolution,
        ).toISO();

        return {
          id: CandleTable.uuid(currentStartedAt, defaultCandle.ticker, resolution),
          startedAt: currentStartedAt,
          ticker: defaultCandle.ticker,
          resolution,
          low: '10000',
          high: defaultPrice2,
          open: '10000',
          close: defaultPrice2,
          baseTokenVolume: '20',
          usdVolume: '250000',
          trades: 2,
          startingOpenInterest: '0',
          orderbookMidPriceClose: '10005',
          orderbookMidPriceOpen: '10005',
        };
      },
    );
    await verifyCandlesInPostgres(expectedCandles);

    await validateCandlesCache();
    expectTimingStats();
  });

  it('creates an empty candle and updates the previous candle orderBookMidPriceClose if startTime is past candle resolution', async () => {
    // Create existing candles
    const existingPrice: string = '7000';
    const startingOpenInterest: string = '200';
    const baseTokenVolume: string = '10';
    const usdVolume: string = Big(existingPrice).times(baseTokenVolume).toString();
    const orderbookMidPriceClose = '7500';
    const orderbookMidPriceOpen = '8000';

    await Promise.all(
      _.map(Object.values(CandleResolution), (resolution: CandleResolution) => {
        return CandleTable.create({
          startedAt: previousStartedAt,
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
          orderbookMidPriceClose,
          orderbookMidPriceOpen,
        });
      }),
    );
    await startCandleCache();

    // Update Orderbook levels
    await updatePriceLevel('BTC-USD', '10010', OrderSide.SELL);
    await updatePriceLevel('BTC-USD', '10000', OrderSide.BUY);

    const publisher: KafkaPublisher = new KafkaPublisher();
    publisher.addEvents([]);

    // Create new candles, without trades
    await runUpdateCandles(publisher);

    // Verify previous candles have orderbookMidPriceClose updated
    const previousExpectedCandles: CandleFromDatabase[] = _.map(
      Object.values(CandleResolution),
      (resolution: CandleResolution) => {
        return {
          id: CandleTable.uuid(previousStartedAt, defaultCandle.ticker, resolution),
          startedAt: previousStartedAt,
          ticker: defaultCandle.ticker,
          resolution,
          low: existingPrice,
          high: existingPrice,
          open: existingPrice,
          close: existingPrice,
          baseTokenVolume,
          usdVolume,
          trades: existingTrades,
          startingOpenInterest,
          orderbookMidPriceClose: '10005',
          orderbookMidPriceOpen,
        };
      },
    );
    await verifyCandlesInPostgres(previousExpectedCandles);

    // Verify new empty candle was created
    const expectedCandles: CandleFromDatabase[] = _.map(
      Object.values(CandleResolution),
      (resolution: CandleResolution) => {
        const currentStartedAt: IsoString = helpers.calculateNormalizedCandleStartTime(
          testConstants.createdDateTime,
          resolution,
        ).toISO();

        return {
          id: CandleTable.uuid(currentStartedAt, defaultCandle.ticker, resolution),
          startedAt: currentStartedAt,
          ticker: defaultCandle.ticker,
          resolution,
          low: existingPrice,
          high: existingPrice,
          open: existingPrice,
          close: existingPrice,
          baseTokenVolume: '0',
          usdVolume: '0',
          trades: 0,
          startingOpenInterest: '0',
          orderbookMidPriceClose: '10005',
          orderbookMidPriceOpen: '10005',
        };
      },
    );
    await verifyCandlesInPostgres(expectedCandles);

  });

  it('successfully creates an orderbook price map for each market', async () => {
    await updatePriceLevel('BTC-USD', '100000', OrderSide.BUY);
    await updatePriceLevel('BTC-USD', '110000', OrderSide.SELL);

    await updatePriceLevel('ISO-USD', '110000', OrderSide.BUY);
    await updatePriceLevel('ISO-USD', '120000', OrderSide.SELL);

    await updatePriceLevel('ETH-USD', '100000', OrderSide.BUY);
    await updatePriceLevel('ETH-USD', '200000', OrderSide.SELL);

    const map = await getOrderbookMidPriceMap();
    expect(map).toEqual({
      'BTC-USD': '105000',
      'ETH-USD': '150000',
      'ISO-USD': '115000',
      'ISO2-USD': undefined,
      'SHIB-USD': undefined,
    });
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
