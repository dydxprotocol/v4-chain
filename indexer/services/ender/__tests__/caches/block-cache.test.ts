import {
  BlockTable,
  CandleTable,
  OraclePriceTable,
  assetRefresher,
  dbHelpers,
  marketRefresher,
  perpetualMarketRefresher,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import {
  getCurrentBlockHeight,
  initializeAllCaches,
  refreshBlockCache,
  resetBlockCache,
  shouldSkipBlock,
} from '../../src/caches/block-cache';
import { clearCandlesMap, getCandlesMap } from '../../src/caches/candle-cache';
import { clearPriceMap, getPriceMap } from '../../src/caches/price-cache';

describe('block-cache', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await CandleTable.create(testConstants.defaultCandle);
    await OraclePriceTable.create(testConstants.defaultOraclePrice);
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    resetBlockCache();
    clearCandlesMap();
    clearPriceMap();
    perpetualMarketRefresher.clear();
    assetRefresher.clear();
    marketRefresher.clear();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  it('block cache initial height should be -1', () => {
    expect(getCurrentBlockHeight()).toEqual('-1');
  });

  it('successfully starts block cache', async () => {
    await refreshBlockCache();

    expect(getCurrentBlockHeight()).toEqual('2');
  });

  describe('shouldSkipBlock', () => {
    it.each([
      [true, 'block.height < currentBlockHeight', '0', false],
      [true, 'block.height == currentBlockHeight', '2', false],
      [false, 'block.height == currentBlockHeight + 1', '3', false],
      [false, 'block.height == currentBlockHeight + 1 with refresh', '4', true],
      [false, 'block.height >= currentBlockHeight + 1 with refresh', '5', true],
    ])('returns %s when %s', async (
      skip: boolean,
      _condition: string,
      lastBlockHeight: string,
      createNextBlock: boolean,
    ) => {
      await refreshBlockCache();

      if (createNextBlock) {
        await BlockTable.create({
          ...testConstants.defaultBlock2,
          blockHeight: '3',
        });
      }
      expect(await shouldSkipBlock(lastBlockHeight)).toEqual(skip);
      if (createNextBlock) {
        // validate that block, candles, and price cache are updated
        expect(getCurrentBlockHeight()).toEqual('3');
        expect(getCandlesMap()).not.toEqual({});
        expect(getPriceMap()).not.toEqual({});
        expect(perpetualMarketRefresher.getPerpetualMarketsMap()).not.toEqual({});
        expect(assetRefresher.getAssetsMap()).not.toEqual({});
        expect(marketRefresher.getMarketsMap()).not.toEqual({});
      }
    });
  });

  describe('initializeAllCaches', () => {
    it('successfully initializes all caches', async () => {
      // Validate that caches are empty
      expect(getCurrentBlockHeight()).toEqual('-1');
      expect(getCandlesMap()).toEqual({});
      expect(getPriceMap()).toEqual({});
      expect(perpetualMarketRefresher.getPerpetualMarketsMap()).toEqual({});
      expect(assetRefresher.getAssetsMap()).toEqual({});
      expect(marketRefresher.getMarketsMap()).toEqual({});

      await initializeAllCaches();

      // Validate that caches are populated
      expect(getCurrentBlockHeight()).toEqual('2');
      expect(getCandlesMap()).not.toEqual({});
      expect(getPriceMap()).not.toEqual({});
      expect(perpetualMarketRefresher.getPerpetualMarketsMap()).not.toEqual({});
      expect(assetRefresher.getAssetsMap()).not.toEqual({});
      expect(marketRefresher.getMarketsMap()).not.toEqual({});
    });
  });
});
