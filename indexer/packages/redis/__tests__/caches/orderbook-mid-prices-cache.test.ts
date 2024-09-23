import { deleteAllAsync } from '../../src/helpers/redis';
import { redis as client } from '../helpers/utils';
import {
  setPrice,
  getMedianPrice,
  ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX,
} from '../../src/caches/orderbook-mid-prices-cache';

describe('orderbook-mid-prices-cache', () => {
  const ticker: string = 'BTC-USD';

  beforeEach(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('setPrice', () => {
    it('sets a price for a ticker', async () => {
      await setPrice(client, ticker, '50000');

      await client.zrange(
        `${ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX}${ticker}`,
        0,
        -1,
        (_: any, response: string[]) => {
          expect(response[0]).toBe('50000');
        },
      );
    });

    it('sets multiple prices for a ticker', async () => {
      await Promise.all([
        setPrice(client, ticker, '50000'),
        setPrice(client, ticker, '51000'),
        setPrice(client, ticker, '49000'),
      ]);

      await client.zrange(
        `${ORDERBOOK_MID_PRICES_CACHE_KEY_PREFIX}${ticker}`,
        0,
        -1,
        (_: any, response: string[]) => {
          expect(response).toEqual(['49000', '50000', '51000']);
        },
      );
    });
  });

  describe('getMedianPrice', () => {
    it('returns null when no prices are set', async () => {
      const result = await getMedianPrice(client, ticker);
      expect(result).toBeNull();
    });

    it('returns the median price for odd number of prices', async () => {
      await Promise.all([
        setPrice(client, ticker, '50000'),
        setPrice(client, ticker, '51000'),
        setPrice(client, ticker, '49000'),
      ]);

      const result = await getMedianPrice(client, ticker);
      expect(result).toBe('50000');
    });

    it('returns the median price for even number of prices', async () => {
      await Promise.all([
        setPrice(client, ticker, '50000'),
        setPrice(client, ticker, '51000'),
        setPrice(client, ticker, '49000'),
        setPrice(client, ticker, '52000'),
      ]);

      const result = await getMedianPrice(client, ticker);
      expect(result).toBe('50500');
    });

    it('returns the correct median price after 5 seconds', async () => {
      jest.useFakeTimers();

      const nowSeconds = Math.floor(Date.now() / 1000);
      jest.setSystemTime(nowSeconds * 1000);

      await Promise.all([
        setPrice(client, ticker, '50000'),
        setPrice(client, ticker, '51000'),
      ]);

      jest.advanceTimersByTime(6000); // Advance time by 6 seconds
      await Promise.all([
        setPrice(client, ticker, '49000'),
        setPrice(client, ticker, '48000'),
        setPrice(client, ticker, '52000'),
        setPrice(client, ticker, '53000'),
      ]);

      const result = await getMedianPrice(client, ticker);
      expect(result).toBe('50500');

      jest.useRealTimers();
    });
  });
});
