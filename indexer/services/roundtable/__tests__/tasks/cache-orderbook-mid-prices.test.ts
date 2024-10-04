import {
  dbHelpers,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import {
  OrderbookLevelsCache,
  OrderbookMidPricesCache,
  redis,
} from '@dydxprotocol-indexer/redis';
import { redisClient } from '../../src/helpers/redis';
import runTask from '../../src/tasks/cache-orderbook-mid-prices';

describe('cache-orderbook-mid-prices', () => {
  beforeEach(async () => {
    await redis.deleteAllAsync(redisClient);
    await testMocks.seedData();
  });

  afterAll(() => {
    jest.restoreAllMocks();
  });

  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  it('caches mid prices for all markets', async () => {
    const market1 = await PerpetualMarketTable
      .findByMarketId(
        testConstants.defaultMarket.id,
      );
    const market2 = await PerpetualMarketTable
      .findByMarketId(
        testConstants.defaultMarket2.id,
      );
    if (!market1) {
      throw new Error('Market 1 not found');
    }
    if (!market2) {
      throw new Error('Market 2 not found');
    }

    jest.spyOn(PerpetualMarketTable, 'findAll')
      .mockReturnValueOnce(Promise.resolve([
        market1,
        // Passing market2 twice so that it will call getOrderbookMidPrice twice and
        // cache the last two prices from the mock below
        market2,
        market2,
      ] as PerpetualMarketFromDatabase[]));

    jest.spyOn(OrderbookLevelsCache, 'getOrderBookMidPrice')
      .mockReturnValueOnce(Promise.resolve('200'))
      .mockReturnValueOnce(Promise.resolve('300'))
      .mockReturnValueOnce(Promise.resolve('400'));

    await runTask();

    const prices = await OrderbookMidPricesCache.getMedianPrices(
      redisClient,
      [market1.ticker, market2.ticker],
    );

    expect(prices[market1.ticker]).toBe('200');
    expect(prices[market2.ticker]).toBe('350');
  });

  it('handles undefined prices', async () => {
    const market1 = await PerpetualMarketTable
      .findByMarketId(
        testConstants.defaultMarket.id,
      );

    if (!market1) {
      throw new Error('Market 1 not found');
    }

    jest.spyOn(PerpetualMarketTable, 'findAll')
      .mockReturnValueOnce(Promise.resolve([market1] as PerpetualMarketFromDatabase[]));

    jest.spyOn(OrderbookLevelsCache, 'getOrderBookMidPrice')
      .mockReturnValueOnce(Promise.resolve(undefined));

    await runTask();

    const price = await OrderbookMidPricesCache.getMedianPrices(redisClient, [market1.ticker]);
    expect(price).toEqual({ 'BTC-USD': undefined });
  });
});
