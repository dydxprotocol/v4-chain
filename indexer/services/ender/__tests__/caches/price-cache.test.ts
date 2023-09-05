import {
  dbHelpers,
  MarketTable,
  OraclePriceFromDatabase,
  OraclePriceTable,
  PriceMap,
  testConstants,
  testMocks,
  MarketCreateObject,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';
import {
  clearPriceMap,
  getPrice,
  getPriceMap,
  startPriceCache,
  updatePriceCacheWithPrice,
} from '../../src/caches/price-cache';

describe('priceCache', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    clearPriceMap();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  const priceFromDatabase: OraclePriceFromDatabase = {
    ...testConstants.defaultOraclePrice,
    id: testConstants.defaultOraclePriceId,
    price: '5000',
  };

  it('uses both Markets table and Oracle prices to populate price cache', async () => {
    await testMocks.seedData();
    const marketWithOraclePrice: MarketCreateObject = {
      id: 127,
      pair: 'NEAR-USD',
      exponent: -9,
      minPriceChangePpm: 1000,
      oraclePrice: '30',
    };
    const marketWithoutOraclePrice: MarketCreateObject = {
      id: 128,
      pair: 'UNI-USD',
      exponent: -9,
      minPriceChangePpm: 1000,
    };
    await Promise.all([
      MarketTable.create(marketWithOraclePrice),
      MarketTable.create(marketWithoutOraclePrice),
      OraclePriceTable.create(testConstants.defaultOraclePrice),
      OraclePriceTable.create(testConstants.defaultOraclePrice2),
    ]);
    await startPriceCache(testConstants.defaultBlock2.blockHeight);

    const map:PriceMap = getPriceMap();
    expect(_.size(map)).toEqual(4);
    expect(
      getPrice(0),
    ).toEqual(
      testConstants.defaultOraclePrice.price,
    );
    expect(
      getPrice(1),
    ).toEqual(
      testConstants.defaultOraclePrice2.price,
    );
    expect(
      getPrice(marketWithOraclePrice.id),
    ).toEqual(
      marketWithOraclePrice.oraclePrice,
    );
    await expect(() => getPrice(marketWithoutOraclePrice.id)).toThrow(
      new Error(`price not found for marketId ${marketWithoutOraclePrice.id} in price cache`),
    );
  });

  it('getPrice throws error on empty price cache', async () => {
    await startPriceCache(testConstants.defaultBlock2.blockHeight);

    const map: PriceMap = getPriceMap();
    expect(_.size(map)).toEqual(0);
    const marketId: number = 0;
    await expect(() => getPrice(marketId)).toThrow(
      new Error(`price not found for marketId ${marketId} in price cache`),
    );
  });

  it('successfully populates price cache', async () => {
    await testMocks.seedData();
    await Promise.all([
      OraclePriceTable.create(testConstants.defaultOraclePrice),
      OraclePriceTable.create(testConstants.defaultOraclePrice2),
    ]);
    await startPriceCache(testConstants.defaultBlock2.blockHeight);

    expect(
      getPrice(0),
    ).toEqual(
      testConstants.defaultOraclePrice.price,
    );
    expect(
      getPrice(1),
    ).toEqual(
      testConstants.defaultOraclePrice2.price,
    );
  });

  it('successfully updates price cache', async () => {
    await testMocks.seedData();
    await Promise.all([
      OraclePriceTable.create(testConstants.defaultOraclePrice),
      OraclePriceTable.create(testConstants.defaultOraclePrice2),
    ]);
    await startPriceCache(testConstants.defaultBlock2.blockHeight);
    updatePriceCacheWithPrice(priceFromDatabase);

    const map: PriceMap = getPriceMap();
    expect(_.size(map)).toEqual(3);
    expect(
      getPrice(0),
    ).toEqual(
      priceFromDatabase.price,
    );
  });
});
