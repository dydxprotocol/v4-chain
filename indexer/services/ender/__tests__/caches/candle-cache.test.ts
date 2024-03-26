import {
  CandleFromDatabase,
  CandlesMap,
  CandleTable,
  dbHelpers,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualMarketTable,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';
import {
  getCandle, getCandlesMap, updateCandleCacheWithCandle, startCandleCache, clearCandlesMap,
} from '../../src/caches/candle-cache';

describe('candleCache', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    clearCandlesMap();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  const candleFromDatabase: CandleFromDatabase = {
    ...testConstants.defaultCandle,
    id: testConstants.defaultCandleId,
  };
  it('successfully populates candle cache with no perpetual markets', async () => {
    await startCandleCache();

    const map: CandlesMap = getCandlesMap();
    expect(_.size(map)).toEqual(0);
    expect(
      getCandle(
        testConstants.defaultPerpetualMarket.ticker,
        testConstants.defaultCandle.resolution,
      ),
    ).toBeUndefined();
  });

  it('successfully populates candle cache with no candles', async () => {
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
    await startCandleCache();

    const map: CandlesMap = getCandlesMap();

    const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [],
    );
    expect(_.size(map)).toEqual(perpetualMarkets.length);
    _.forEach(perpetualMarkets, (perpetualMarket: PerpetualMarketFromDatabase) => {
      expect(_.size(map[perpetualMarket.ticker])).toEqual(0);
    });
    expect(
      getCandle(
        testConstants.defaultPerpetualMarket.ticker,
        testConstants.defaultCandle.resolution,
      ),
    ).toBeUndefined();
  });

  it('successfully populates candle cache', async () => {
    await Promise.all([
      testMocks.seedData(),
      CandleTable.create(testConstants.defaultCandle),
    ]);
    await startCandleCache();

    const map: CandlesMap = getCandlesMap();
    expect(_.size(map[testConstants.defaultPerpetualMarket.ticker])).toEqual(1);
    expect(_.size(map[testConstants.defaultPerpetualMarket2.ticker])).toEqual(0);
    expect(_.size(map[testConstants.defaultPerpetualMarket3.ticker])).toEqual(0);
    expect(
      map[testConstants.defaultPerpetualMarket.ticker][testConstants.defaultCandle.resolution],
    ).toEqual(
      candleFromDatabase,
    );

    expect(
      getCandle(
        testConstants.defaultPerpetualMarket.ticker,
        testConstants.defaultCandle.resolution,
      ),
    ).toEqual(
      candleFromDatabase,
    );
  });

  it('successfully updates candle cache', () => {
    updateCandleCacheWithCandle(candleFromDatabase);

    const map: CandlesMap = getCandlesMap();
    expect(_.size(map)).toEqual(1);
    expect(_.size(map[testConstants.defaultPerpetualMarket.ticker])).toEqual(1);
    expect(
      map[testConstants.defaultPerpetualMarket.ticker][testConstants.defaultCandle.resolution],
    ).toEqual(
      candleFromDatabase,
    );

    expect(
      getCandle(
        testConstants.defaultPerpetualMarket.ticker,
        testConstants.defaultCandle.resolution,
      ),
    ).toEqual(
      candleFromDatabase,
    );
  });
});
