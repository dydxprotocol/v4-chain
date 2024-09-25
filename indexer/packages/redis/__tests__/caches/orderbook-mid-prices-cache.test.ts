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

    it('returns the correct median price for small numbers with even number of prices', async () => {
      await Promise.all([
        setPrice(client, ticker, '0.00000000002345'),
        setPrice(client, ticker, '0.00000000002346'),
      ]);

      const midPrice1 = await getMedianPrice(client, ticker);
      expect(midPrice1).toEqual('0.000000000023455');
    });

    it('returns the correct median price for small numbers with odd number of prices', async () => {
      await Promise.all([
        setPrice(client, ticker, '0.00000000001'),
        setPrice(client, ticker, '0.00000000002'),
        setPrice(client, ticker, '0.00000000003'),
        setPrice(client, ticker, '0.00000000004'),
        setPrice(client, ticker, '0.00000000005'),
      ]);

      const midPrice1 = await getMedianPrice(client, ticker);
      expect(midPrice1).toEqual('0.00000000003');

      await deleteAllAsync(client);

      await Promise.all([
        setPrice(client, ticker, '0.00000847007'),
        setPrice(client, ticker, '0.00000847006'),
        setPrice(client, ticker, '0.00000847008'),
      ]);

      const midPrice2 = await getMedianPrice(client, ticker);
      expect(midPrice2).toEqual('0.00000847007');
    });
  });
});
