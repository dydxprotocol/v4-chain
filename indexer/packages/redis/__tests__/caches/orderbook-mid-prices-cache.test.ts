import { deleteAllAsync } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';
import {
  fetchAndCacheOrderbookMidPrices,
  getMedianPrices,
  ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX,
} from '../../src/caches/orderbook-mid-prices-cache';
import * as OrderbookLevelsCache from '../../src/caches/orderbook-levels-cache';

// Mock the OrderbookLevelsCache module
jest.mock('../../src/caches/orderbook-levels-cache', () => ({
  getOrderBookMidPrice: jest.fn(),
}));

describe('orderbook-mid-prices-cache', () => {
  const defaultTicker: string = 'BTC-USD';

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

      await fetchAndCacheOrderbookMidPrices(client, [defaultTicker]);

      expect(OrderbookLevelsCache.getOrderBookMidPrice).toHaveBeenCalledTimes(1);
      expect(OrderbookLevelsCache.getOrderBookMidPrice).toHaveBeenCalledWith(
        defaultTicker,
        client,
      );

      client.zrange(
        `${ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX}${defaultTicker}`,
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
        await fetchAndCacheOrderbookMidPrices(client, [defaultTicker]);
      }

      client.zrange(
        `${ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX}${defaultTicker}`,
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
        [defaultTicker]: '49000',
        [ticker2]: '50000',
        [ticker3]: '51000',
      };

      // Mock the getOrderBookMidPrice function for each ticker
      (OrderbookLevelsCache.getOrderBookMidPrice as jest.Mock)
        .mockResolvedValueOnce(mockPrices[defaultTicker])
        .mockResolvedValueOnce(mockPrices[ticker2])
        .mockResolvedValueOnce(mockPrices[ticker3]);

      await fetchAndCacheOrderbookMidPrices(client, [defaultTicker, ticker2, ticker3]);
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
      const result: {[ticker: string]: string | undefined} = await getMedianPrices(
        client,
        [defaultTicker],
      );
      expect(result).toEqual({ 'BTC-USD': undefined });
    });

    it('returns the median price for odd number of prices', async () => {
      setPrice(defaultTicker, '51000');
      setPrice(defaultTicker, '50000');
      setPrice(defaultTicker, '49000');

      const result: {[ticker: string]: string | undefined} = await getMedianPrices(
        client,
        [defaultTicker],
      );
      expect(result).toEqual({ 'BTC-USD': '50000' });
    });

    it('returns the median price for even number of prices', async () => {
      setPrice(defaultTicker, '50000');
      setPrice(defaultTicker, '51000');
      setPrice(defaultTicker, '49000');
      setPrice(defaultTicker, '52000');

      const result: {[ticker: string]: string | undefined} = await getMedianPrices(
        client,
        [defaultTicker],
      );
      expect(result).toEqual({ 'BTC-USD': '50500' });
    });

    it('returns the correct median price after 60 seconds', async () => {
      jest.useFakeTimers();
      // Mock the getOrderBookMidPrice function for the ticker
      const mockPrices: string[] = ['50000', '51000', '49000', '48000', '52000', '53000'];

      (OrderbookLevelsCache.getOrderBookMidPrice as jest.Mock)
        .mockResolvedValueOnce(mockPrices[0])
        .mockResolvedValueOnce(mockPrices[1])
        .mockResolvedValueOnce(mockPrices[2])
        .mockResolvedValueOnce(mockPrices[3])
        .mockResolvedValueOnce(mockPrices[4])
        .mockResolvedValueOnce(mockPrices[5]);

      // Fetch and cache initial prices
      await fetchAndCacheOrderbookMidPrices(client, [defaultTicker, defaultTicker]);
      expect(OrderbookLevelsCache.getOrderBookMidPrice).toHaveBeenCalledTimes(2);

      // Advance time and fetch more prices
      jest.advanceTimersByTime(61000); // Advance time by 61 seconds
      await fetchAndCacheOrderbookMidPrices(
        client,
        [defaultTicker, defaultTicker, defaultTicker, defaultTicker],
      );

      client.zrange(`${ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX}${defaultTicker}`,
        0,
        -1,
        (err: Error, res: string[]) => {
          expect(res).toHaveLength(4);
        });
      expect(OrderbookLevelsCache.getOrderBookMidPrice).toHaveBeenCalledTimes(6);

      // Check the median price
      const result:{[ticker: string]: string | undefined} = await getMedianPrices(
        client,
        [defaultTicker],
      );
      // Median of last 4 prices, as first two should have expired after moving clock forward
      expect(result).toEqual({ 'BTC-USD': '50500' });

      jest.useRealTimers();
    });

    it('returns the correct median price for small numbers with even number of prices', async () => {
      setPrice(defaultTicker, '0.00000000002345');
      setPrice(defaultTicker, '0.00000000002346');

      const midPrice1: { [ticker: string]: string | undefined } = await getMedianPrices(
        client,
        [defaultTicker],
      );
      expect(midPrice1).toEqual({ 'BTC-USD': '0.000000000023455' });
    });

    it('returns the correct median price for small numbers with odd number of prices', async () => {
      setPrice(defaultTicker, '0.00000000001');
      setPrice(defaultTicker, '0.00000000002');
      setPrice(defaultTicker, '0.00000000003');
      setPrice(defaultTicker, '0.00000000004');
      setPrice(defaultTicker, '0.00000000005');

      const midPrice1: { [ticker: string]: string | undefined } = await getMedianPrices(
        client,
        [defaultTicker],
      );
      expect(midPrice1).toEqual({ 'BTC-USD': '0.00000000003' });

      await deleteAllAsync(client);

      setPrice(defaultTicker, '0.00000847007');
      setPrice(defaultTicker, '0.00000847006');
      setPrice(defaultTicker, '0.00000847008');

      const midPrice2: { [ticker: string]: string | undefined } = await getMedianPrices(
        client,
        [defaultTicker],
      );
      expect(midPrice2).toEqual({ 'BTC-USD': '0.00000847007' });
    });
  });

  describe('getMedianPrices for multiple markets', () => {
    const btcUsdTicker = 'BTC-USD';
    const ethUsdTicker = 'ETH-USD';
    const solUsdTicker = 'SOL-USD';

    beforeEach(async () => {
      await deleteAllAsync(client);
    });

    it('returns correct median prices for multiple markets with odd number of prices', async () => {
      // Set prices for BTC-USD
      setPrice(btcUsdTicker, '50000');
      setPrice(btcUsdTicker, '51000');
      setPrice(btcUsdTicker, '49000');

      // Set prices for ETH-USD
      setPrice(ethUsdTicker, '3000');
      setPrice(ethUsdTicker, '3100');
      setPrice(ethUsdTicker, '2900');

      // Set prices for SOL-USD
      setPrice(solUsdTicker, '100');
      setPrice(solUsdTicker, '102');
      setPrice(solUsdTicker, '98');

      const result: { [ticker: string]: string | undefined } = await getMedianPrices(
        client,
        [btcUsdTicker, ethUsdTicker, solUsdTicker],
      );
      expect(result).toEqual({
        'BTC-USD': '50000',
        'ETH-USD': '3000',
        'SOL-USD': '100',
      });
    });

    it('returns correct median prices for multiple markets with even number of prices', async () => {
      // Set prices for BTC-USD
      setPrice(btcUsdTicker, '50000');
      setPrice(btcUsdTicker, '51000');
      setPrice(btcUsdTicker, '49000');
      setPrice(btcUsdTicker, '52000');

      // Set prices for ETH-USD
      setPrice(ethUsdTicker, '3000');
      setPrice(ethUsdTicker, '3100');
      setPrice(ethUsdTicker, '2900');
      setPrice(ethUsdTicker, '3200');

      const result: { [ticker: string]: string | undefined } = await getMedianPrices(
        client,
        [btcUsdTicker, ethUsdTicker],
      );
      expect(result).toEqual({
        'BTC-USD': '50500',
        'ETH-USD': '3050',
      });
    });

    it('handles markets with different numbers of prices', async () => {
      // Set prices for BTC-USD (odd number)
      setPrice(btcUsdTicker, '50000');
      setPrice(btcUsdTicker, '51000');
      setPrice(btcUsdTicker, '49000');

      // Set prices for ETH-USD (even number)
      setPrice(ethUsdTicker, '3000');
      setPrice(ethUsdTicker, '3100');
      setPrice(ethUsdTicker, '2900');
      setPrice(ethUsdTicker, '3200');

      // Set no prices for SOL-USD

      const result: { [ticker: string]: string | undefined } = await getMedianPrices(
        client,
        [btcUsdTicker, ethUsdTicker, solUsdTicker],
      );
      expect(result).toEqual({
        'BTC-USD': '50000',
        'ETH-USD': '3050',
        'SOL-USD': undefined,
      });
    });

    it('calculates correct median prices for markets with small and large numbers', async () => {
      // Set prices for BTC-USD (large numbers)
      setPrice(btcUsdTicker, '50000.12345');
      setPrice(btcUsdTicker, '50000.12346');

      // Set prices for ETH-USD (medium numbers)
      setPrice(ethUsdTicker, '3000.5');
      setPrice(ethUsdTicker, '3000.6');
      setPrice(ethUsdTicker, '3000.7');

      // Set prices for SOL-USD (small numbers)
      setPrice(solUsdTicker, '0.00000123');
      setPrice(solUsdTicker, '0.00000124');
      setPrice(solUsdTicker, '0.00000125');
      setPrice(solUsdTicker, '0.00000126');

      const result: { [ticker: string]: string | undefined } = await getMedianPrices(
        client,
        [btcUsdTicker, ethUsdTicker, solUsdTicker],
      );
      expect(result).toEqual({
        'BTC-USD': '50000.123455',
        'ETH-USD': '3000.6',
        'SOL-USD': '0.000001245',
      });
    });
  });
});
