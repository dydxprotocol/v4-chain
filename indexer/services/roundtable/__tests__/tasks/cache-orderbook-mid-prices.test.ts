import {
  dbHelpers,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  testConstants,
  testMocks,
  perpetualMarketRefresher,
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
    await perpetualMarketRefresher.updatePerpetualMarkets();
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
    await perpetualMarketRefresher.clear();
  });

  it('caches mid prices for all markets', async () => {
    const market1: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .findByMarketId(
        testConstants.defaultMarket.id,
      );
    const market2: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .findByMarketId(
        testConstants.defaultMarket2.id,
      );
    const market3: PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable
      .findByMarketId(
        testConstants.defaultMarket3.id,
      );
    if (!market1 || !market2 || !market3) {
      throw new Error('Test market not found');
    }

    jest.spyOn(OrderbookLevelsCache, 'getOrderBookMidPrice')
      .mockReturnValueOnce(Promise.resolve('200'))
      .mockReturnValueOnce(Promise.resolve('300'))
      .mockReturnValueOnce(Promise.resolve('400'));

    await runTask();
    expect(OrderbookLevelsCache.getOrderBookMidPrice).toHaveBeenCalledTimes(5);

    const prices: {[ticker: string]: string | undefined} = await
    OrderbookMidPricesCache.getMedianPrices(
      redisClient,
      [market1.ticker, market2.ticker, market3.ticker],
    );

    expect(prices[market1.ticker]).toBe('200');
    expect(prices[market2.ticker]).toBe('300');
    expect(prices[market3.ticker]).toBe('400');
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
