import { deleteAllAsync } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';
import {
  fetchAndCacheOrderbookMidPrices,
  getMedianPrice,
  ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX,
} from '../../src/caches/orderbook-mid-prices-cache';
import * as OrderbookLevelsCache from '../../src/caches/orderbook-levels-cache';

// Mock the OrderbookLevelsCache module
jest.mock('../../src/caches/orderbook-levels-cache', () => ({
  getOrderBookMidPrice: jest.fn(),
}));

describe('orderbook-mid-prices-cache', () => {
  const ticker: string = 'BTC-USD';

  // Helper function to set a price for a given market ticker
  const setPrice = (marketTicker: string, price: string) => {
    const now = Date.now();
    client.zadd(`${ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX}${marketTicker}`, now, price);
  };

  afterAll(async () => {
    await deleteAllAsync(client);
  });

  beforeEach(async () => {
    await deleteAllAsync(client);
    jest.resetAllMocks();
    (OrderbookLevelsCache.getOrderBookMidPrice as jest.Mock).mockReset();
  });

  describe('fetchAndCacheOrderbookMidPrices', () => {
    it('sets a price for a ticker', async () => {
      (OrderbookLevelsCache.getOrderBookMidPrice as jest.Mock).mockResolvedValue('50000');

      await fetchAndCacheOrderbookMidPrices(client, [ticker]);

      expect(OrderbookLevelsCache.getOrderBookMidPrice).toHaveBeenCalledTimes(1);
      expect(OrderbookLevelsCache.getOrderBookMidPrice).toHaveBeenCalledWith(
        `${ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX}${ticker}`,
        client,
      );

      client.zrange(
        `${ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX}${ticker}`,
        0,
        -1,
        (_: any, response: string[]) => {
          expect(response[0]).toBe('50000');
        },
      );
    });

    it('sets multiple prices for a ticker', async () => {
      const mockPrices = ['49000', '50000', '51000'];
      for (const price of mockPrices) {
        (OrderbookLevelsCache.getOrderBookMidPrice as jest.Mock).mockResolvedValue(price);
        await fetchAndCacheOrderbookMidPrices(client, [ticker]);
      }

      client.zrange(
        `${ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX}${ticker}`,
        0,
        -1,
        (_: any, response: string[]) => {
          expect(response).toEqual(mockPrices);
        },
      );
      expect(OrderbookLevelsCache.getOrderBookMidPrice).toHaveBeenCalledTimes(3);
    });

    it('sets prices for multiple tickers', async () => {
      const ticker2 = 'SHIB-USD';
      const ticker3 = 'SOL-USD';
      const mockPrices = {
        [ticker]: '49000',
        [ticker2]: '50000',
        [ticker3]: '51000',
      };

      // Mock the getOrderBookMidPrice function for each ticker
      (OrderbookLevelsCache.getOrderBookMidPrice as jest.Mock)
        .mockResolvedValueOnce(mockPrices[ticker])
        .mockResolvedValueOnce(mockPrices[ticker2])
        .mockResolvedValueOnce(mockPrices[ticker3]);

      await fetchAndCacheOrderbookMidPrices(client, [ticker, ticker2, ticker3]);
      expect(OrderbookLevelsCache.getOrderBookMidPrice).toHaveBeenCalledTimes(3);

      for (const [key, price] of Object.entries(mockPrices)) {
        client.zrange(`${ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX}${key}`,
          0,
          -1,
          (err: Error, res: string[]) => {
            expect(res).toHaveLength(1);
            expect(res[0]).toEqual(price);
          });
      }
    });
  });

  describe('getMedianPrice', () => {
    it('returns null when no prices are set', async () => {
      const result = await getMedianPrice(client, ticker);
      expect(result).toBeNull();
    });

    it('returns the median price for odd number of prices', async () => {
      setPrice(ticker, '50000');
      setPrice(ticker, '51000');
      setPrice(ticker, '49000');

      const result = await getMedianPrice(client, ticker);
      expect(result).toBe('50000');
    });

    it('returns the median price for even number of prices', async () => {
      setPrice(ticker, '50000');
      setPrice(ticker, '51000');
      setPrice(ticker, '49000');
      setPrice(ticker, '52000');

      const result = await getMedianPrice(client, ticker);
      expect(result).toBe('50500');
    });

    it('returns the correct median price after 5 seconds', async () => {
      jest.useFakeTimers();
      // Mock the getOrderBookMidPrice function for the ticker
      const mockPrices = ['50000', '51000', '49000', '48000', '52000', '53000'];

      (OrderbookLevelsCache.getOrderBookMidPrice as jest.Mock)
        .mockResolvedValueOnce(mockPrices[0])
        .mockResolvedValueOnce(mockPrices[1])
        .mockResolvedValueOnce(mockPrices[2])
        .mockResolvedValueOnce(mockPrices[3])
        .mockResolvedValueOnce(mockPrices[4])
        .mockResolvedValueOnce(mockPrices[5]);

      // Fetch and cache initial prices
      await fetchAndCacheOrderbookMidPrices(client, [ticker, ticker]);
      expect(OrderbookLevelsCache.getOrderBookMidPrice).toHaveBeenCalledTimes(2);

      // Advance time and fetch more prices
      jest.advanceTimersByTime(6000); // Advance time by 6 seconds
      await fetchAndCacheOrderbookMidPrices(client, [ticker, ticker, ticker, ticker]);
      expect(OrderbookLevelsCache.getOrderBookMidPrice).toHaveBeenCalledTimes(6);

      // Check the median price
      const result = await getMedianPrice(client, ticker);
      // Median of last 4 prices, as first two should have expired after moving clock forward
      expect(result).toBe('50500');

      jest.useRealTimers();
    });

    it('returns the correct median price for small numbers with even number of prices', async () => {
      setPrice(ticker, '0.00000000002345');
      setPrice(ticker, '0.00000000002346');

      const midPrice1 = await getMedianPrice(client, ticker);
      expect(midPrice1).toEqual('0.000000000023455');
    });

    it('returns the correct median price for small numbers with odd number of prices', async () => {
      setPrice(ticker, '0.00000000001');
      setPrice(ticker, '0.00000000002');
      setPrice(ticker, '0.00000000003');
      setPrice(ticker, '0.00000000004');
      setPrice(ticker, '0.00000000005');

      const midPrice1 = await getMedianPrice(client, ticker);
      expect(midPrice1).toEqual('0.00000000003');

      await deleteAllAsync(client);

      setPrice(ticker, '0.00000847007');
      setPrice(ticker, '0.00000847006');
      setPrice(ticker, '0.00000847008');

      const midPrice2 = await getMedianPrice(client, ticker);
      expect(midPrice2).toEqual('0.00000847007');
    });
  });
});
