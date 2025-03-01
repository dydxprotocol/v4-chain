import { DateTime } from 'luxon';

import {
  CandleFromDatabase, CandleResolution, CandlesMap, CandleUpdateObject,
} from '../../src/types';
import * as BlockTable from '../../src/stores/block-table';
import * as CandleTable from '../../src/stores/candle-table';
import {
  clearData,
  teardown,
} from '../../src/helpers/db-helpers';
import {
  defaultCandle, defaultCandleId, defaultPerpetualMarket2,
} from '../helpers/constants';

describe('CandleTable', () => {
  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a Candle', async () => {
    await CandleTable.create(defaultCandle);
  });

  it('Successfully finds all Candles', async () => {
    await Promise.all([
      CandleTable.create(defaultCandle),
      CandleTable.create({
        ...defaultCandle,
        ticker: defaultPerpetualMarket2.ticker,
      }),
    ]);

    const candles: CandleFromDatabase[] = await CandleTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(candles.length).toEqual(2);
    expect(candles[0]).toEqual(expect.objectContaining({
      ...defaultCandle,
      ticker: defaultPerpetualMarket2.ticker,
    }));
    expect(candles[1]).toEqual(expect.objectContaining(defaultCandle));
  });

  it('Successfully finds a Candle', async () => {
    await CandleTable.create(defaultCandle);

    const candle: CandleFromDatabase | undefined = await CandleTable.findById(
      defaultCandleId,
    );

    expect(candle).toEqual(expect.objectContaining(defaultCandle));
  });

  it('Fails to finds a nonexistent Candle', async () => {
    const candle: CandleFromDatabase | undefined = await CandleTable.findById(
      defaultCandleId,
    );
    expect(candle).toEqual(undefined);
  });

  it('Successfully updates a candle', async () => {
    await CandleTable.create(defaultCandle);

    const updatedCandle: CandleUpdateObject = {
      id: defaultCandleId,
      open: '100',
      orderbookMidPriceClose: '200',
      orderbookMidPriceOpen: '300',
    };

    await CandleTable.update(updatedCandle);

    const candle: CandleFromDatabase | undefined = await CandleTable.findById(
      defaultCandleId,
    );

    expect(candle).toEqual(expect.objectContaining(updatedCandle));
  });

  it('Successfully finds candles map based on latest block', async () => {
    const latestBlockTime: DateTime = DateTime.utc(2022, 6, 1);
    const candleMinuteResThreeHoursAgo = {
      ...defaultCandle,
      resolution: CandleResolution.ONE_MINUTE,
      startedAt: latestBlockTime.minus({ hours: 3 }).toISO(),
    };
    const candleFiveMinuteResTwoHoursAgo = {
      ...defaultCandle,
      resolution: CandleResolution.FIVE_MINUTES,
      startedAt: latestBlockTime.minus({ hours: 2 }).toISO(),
    };
    const candleHourResThreeDaysAgo = {
      ...defaultCandle,
      resolution: CandleResolution.FOUR_HOURS,
      startedAt: latestBlockTime.minus({ days: 3 }).toISO(),
    };
    const candleDayResOneDayAgo = {
      ...defaultCandle,
      resolution: CandleResolution.ONE_DAY,
      startedAt: latestBlockTime.minus({ days: 1 }).toISO(),
    };
    await Promise.all([
      // Create two blocks with block time 1 second apart.
      BlockTable.create({
        blockHeight: '1',
        time: latestBlockTime.minus({ seconds: 1 }).toISO(),
      }),
      BlockTable.create({
        blockHeight: '2',
        time: latestBlockTime.toISO(),
      }),
      // Create candles.
      CandleTable.create(candleMinuteResThreeHoursAgo), // should not be part of candles map
      CandleTable.create(candleFiveMinuteResTwoHoursAgo), // should be part of candles map
      CandleTable.create(candleHourResThreeDaysAgo), // should not be part of candles map
      CandleTable.create(candleDayResOneDayAgo), // should be part of candles map
    ]);

    const candlesMap: CandlesMap = await CandleTable.findCandlesMap(
      [defaultCandle.ticker],
      latestBlockTime.toISO(),
    );

    expect(Object.keys(candlesMap)).toEqual([defaultCandle.ticker]);
    expect(
      Object.keys(candlesMap[defaultCandle.ticker]),
    ).toEqual([CandleResolution.FIVE_MINUTES, CandleResolution.ONE_DAY]);
    expect(
      candlesMap[defaultCandle.ticker][CandleResolution.FIVE_MINUTES],
    ).toEqual(expect.objectContaining(candleFiveMinuteResTwoHoursAgo));
    expect(
      candlesMap[defaultCandle.ticker][CandleResolution.ONE_DAY],
    ).toEqual(expect.objectContaining(candleDayResOneDayAgo));
  });

  it('Successfully finds candles within lookback period', async () => {
    const now: DateTime = DateTime.utc();
    const candle23HoursAgo = {
      ...defaultCandle,
      startedAt: now.minus({ hours: 23 }).toISO(),
    };
    const candle25HoursAgo = {
      ...defaultCandle,
      startedAt: now.minus({ hours: 25 }).toISO(),
    };
    const candle6DaysAgo = {
      ...defaultCandle,
      startedAt: now.minus({ days: 6 }).toISO(),
    };
    const candle8DaysAgo = {
      ...defaultCandle,
      startedAt: now.minus({ days: 8 }).toISO(),
    };
    const otherTickerCandle23HoursAgo = {
      ...defaultCandle,
      ticker: defaultPerpetualMarket2.ticker,
      startedAt: now.minus({ hours: 23 }).toISO(),
    };
    const otherTickerCandle6DaysAgo = {
      ...defaultCandle,
      ticker: defaultPerpetualMarket2.ticker,
      startedAt: now.minus({ days: 6 }).toISO(),
    };

    await Promise.all([
      CandleTable.create(candle23HoursAgo),
      CandleTable.create(candle25HoursAgo),
      CandleTable.create(candle6DaysAgo),
      CandleTable.create(candle8DaysAgo),
      CandleTable.create(otherTickerCandle23HoursAgo),
      CandleTable.create(otherTickerCandle6DaysAgo),
    ]);

    // Test ONE_DAY lookback
    const oneDayCandles = await CandleTable.findByResAndLookbackPeriod(
      defaultCandle.resolution,
      24 * 60 * 60 * 1000,
    );

    expect(oneDayCandles.length).toEqual(2);
    expect(oneDayCandles).toEqual(
      expect.arrayContaining([
        expect.objectContaining(candle23HoursAgo),
        expect.objectContaining(otherTickerCandle23HoursAgo),
      ]),
    );

    // Test SEVEN_DAYS lookback
    const sevenDayCandles = await CandleTable.findByResAndLookbackPeriod(
      defaultCandle.resolution,
      7 * 24 * 60 * 60 * 1000,
    );

    expect(sevenDayCandles.length).toEqual(5);
    expect(sevenDayCandles).toEqual(
      expect.arrayContaining([
        expect.objectContaining(candle23HoursAgo),
        expect.objectContaining(candle25HoursAgo),
        expect.objectContaining(candle6DaysAgo),
        expect.objectContaining(otherTickerCandle23HoursAgo),
        expect.objectContaining(otherTickerCandle6DaysAgo),
      ]),
    );
  });
});
