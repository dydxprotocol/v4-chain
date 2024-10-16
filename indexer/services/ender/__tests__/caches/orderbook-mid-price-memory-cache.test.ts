import { OrderbookMidPricesCache } from '@dydxprotocol-indexer/redis';
import * as orderbookMidPriceMemoryCache from '../../src/caches/orderbook-mid-price-memory-cache';
import {
  dbHelpers,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import config from '../../src/config';
import { logger, stats } from '@dydxprotocol-indexer/base';

describe('orderbook-mid-price-memory-cache', () => {

  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
  });

  beforeEach(async () => {
    await testMocks.seedData();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  describe('getOrderbookMidPrice', () => {
    it('should return the mid price for a given ticker', async () => {
      jest.spyOn(OrderbookMidPricesCache, 'getMedianPrices')
        .mockReturnValue(Promise.resolve({ 'BTC-USD': '300', 'ETH-USD': '200' }));

      await orderbookMidPriceMemoryCache.updateOrderbookMidPrices();

      expect(orderbookMidPriceMemoryCache.getOrderbookMidPrice('BTC-USD')).toBe('300');
      expect(orderbookMidPriceMemoryCache.getOrderbookMidPrice('ETH-USD')).toBe('200');
    });
  });

  describe('updateOrderbookMidPrices', () => {
    it('should update the orderbook mid price cache', async () => {
      const mockMedianPrices = {
        'BTC-USD': '50000',
        'ETH-USD': '3000',
        'SOL-USD': '1000',
      };

      jest.spyOn(OrderbookMidPricesCache, 'getMedianPrices')
        .mockResolvedValue(mockMedianPrices);

      await orderbookMidPriceMemoryCache.updateOrderbookMidPrices();

      expect(orderbookMidPriceMemoryCache.getOrderbookMidPrice('BTC-USD')).toBe('50000');
      expect(orderbookMidPriceMemoryCache.getOrderbookMidPrice('ETH-USD')).toBe('3000');
      expect(orderbookMidPriceMemoryCache.getOrderbookMidPrice('SOL-USD')).toBe('1000');
    });

    it('should handle errors and log them', async () => {
      const mockError = new Error('Test error');
      jest.spyOn(OrderbookMidPricesCache, 'getMedianPrices').mockImplementation(() => {
        throw mockError;
      });

      jest.spyOn(logger, 'error');
      await orderbookMidPriceMemoryCache.updateOrderbookMidPrices();

      expect(logger.error).toHaveBeenCalledWith(
        expect.objectContaining({
          at: 'orderbook-mid-price-cache#updateOrderbookMidPrices',
          message: 'Failed to fetch OrderbookMidPrices',
          error: mockError,
        }),
      );
    });

    it('should record timing stats', async () => {
      jest.spyOn(stats, 'timing');
      await orderbookMidPriceMemoryCache.updateOrderbookMidPrices();

      expect(stats.timing).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.update_orderbook_mid_prices_cache.timing`,
        expect.any(Number),
      );
    });
  });
});
