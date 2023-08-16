import {
  BlockTable,
  CandleTable,
  OraclePriceTable,
  dbHelpers,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { getCurrentBlockHeight, refreshBlockCache, shouldSkipBlock } from '../../src/caches/block-cache';
import { getCandlesMap } from '../../src/caches/candle-cache';
import { getPriceMap } from '../../src/caches/price-cache';

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
  });

  afterAll(async () => {
    await dbHelpers.teardown();
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
      }
    });
  });
});
