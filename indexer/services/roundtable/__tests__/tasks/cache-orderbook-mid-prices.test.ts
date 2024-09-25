import {
  dbHelpers,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import {
  OrderbookMidPricesCache,
  OrderbookLevelsCache,
  redis,
} from '@dydxprotocol-indexer/redis';
import { redisClient } from '../../src/helpers/redis';
import runTask from '../../src/tasks/cache-orderbook-mid-prices';

jest.mock('@dydxprotocol-indexer/base', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/base'),
  logger: {
    info: jest.fn(),
    error: jest.fn(),
  },
}));

jest.mock('@dydxprotocol-indexer/redis', () => ({
  ...jest.requireActual('@dydxprotocol-indexer/redis'),
  OrderbookLevelsCache: {
    getOrderBookMidPrice: jest.fn(),
  },
}));

describe('cache-orderbook-mid-prices', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await dbHelpers.clearData();
    await redis.deleteAllAsync(redisClient);
    await testMocks.seedData();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  it('caches mid prices for all markets', async () => {
    const market1 = testConstants.defaultPerpetualMarket;
    const market2 = testConstants.defaultPerpetualMarket2;

    const mockGetOrderBookMidPrice = jest.spyOn(OrderbookLevelsCache, 'getOrderBookMidPrice');
    mockGetOrderBookMidPrice.mockResolvedValueOnce('100.5'); // For market1
    mockGetOrderBookMidPrice.mockResolvedValueOnce('200.75'); // For market2

    await runTask();

    // Check if the mock was called with the correct arguments
    expect(mockGetOrderBookMidPrice).toHaveBeenCalledWith(market1.ticker, redisClient);
    expect(mockGetOrderBookMidPrice).toHaveBeenCalledWith(market2.ticker, redisClient);

    // Check if prices were cached correctly
    const price1 = await OrderbookMidPricesCache.getMedianPrice(redisClient, market1.ticker);
    const price2 = await OrderbookMidPricesCache.getMedianPrice(redisClient, market2.ticker);

    expect(price1).toBe('100.5');
    expect(price2).toBe('200.75');
  });

  it('handles undefined prices', async () => {
    const market = testConstants.defaultPerpetualMarket;

    const mockGetOrderBookMidPrice = jest.spyOn(OrderbookLevelsCache, 'getOrderBookMidPrice');
    mockGetOrderBookMidPrice.mockResolvedValueOnce(undefined);

    await runTask();

    const price = await OrderbookMidPricesCache.getMedianPrice(redisClient, market.ticker);
    expect(price).toBeNull();

    // Check that a log message was created
    expect(jest.requireMock('@dydxprotocol-indexer/base').logger.info).toHaveBeenCalledWith({
      at: 'cache-orderbook-mid-prices#runTask',
      message: `undefined price for ${market.ticker}`,
    });
  });

  it('handles errors', async () => {
    // Mock OrderbookLevelsCache.getOrderBookMidPrice to throw an error
    const mockGetOrderBookMidPrice = jest.spyOn(OrderbookLevelsCache, 'getOrderBookMidPrice');
    mockGetOrderBookMidPrice.mockRejectedValueOnce(new Error('Test error'));

    await runTask();

    expect(jest.requireMock('@dydxprotocol-indexer/base').logger.error).toHaveBeenCalledWith({
      at: 'cache-orderbook-mid-prices#runTask',
      message: 'Test error',
      error: expect.any(Error),
    });
  });
});
