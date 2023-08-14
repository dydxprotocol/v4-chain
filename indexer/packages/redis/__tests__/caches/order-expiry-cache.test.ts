import {
  ORDER_EXPIRY_CACHE_KEY,
  getOrderExpiries,
  getOrdersAndExpiries,
} from '../../src/caches/order-expiry-cache';
import { deleteAllAsync, zAddAsync } from '../../src/helpers/redis';
import {
  redis as client,
} from '../helpers/utils';
import { redisOrder, redisOrderGoodTilBlockTime, secondRedisOrder } from './constants';

describe('orderExpiryCache', () => {
  beforeAll(async () => {
    await deleteAllAsync(client);
  });

  afterEach(async () => {
    await deleteAllAsync(client);
  });

  describe('getOrdersAndExpiries', () => {
    beforeEach(async () => {
      await Promise.all([
        zAddAsync({ key: ORDER_EXPIRY_CACHE_KEY, value: redisOrder.id, score: 10 }, client),
        zAddAsync({ key: ORDER_EXPIRY_CACHE_KEY, value: secondRedisOrder.id, score: 20 }, client),
        zAddAsync(
          { key: ORDER_EXPIRY_CACHE_KEY, value: redisOrderGoodTilBlockTime.id, score: 30 },
          client,
        ),
      ]);
    });
    it('gets all scores up to given, inclusive', async () => {
      const actual20: Record<string, Number> = await getOrdersAndExpiries(
        {
          latestExpiry: 20,
        },
        client,
      );
      expect(actual20).toEqual({
        [redisOrder.id]: 10,
        [secondRedisOrder.id]: 20,
      });
      const actual25: Record<string, Number> = await getOrdersAndExpiries(
        {
          latestExpiry: 25,
        },
        client,
      );
      expect(actual25).toEqual({
        [redisOrder.id]: 10,
        [secondRedisOrder.id]: 20,
      });
      const actual30: Record<string, Number> = await getOrdersAndExpiries(
        {
          latestExpiry: 30,
        },
        client,
      );
      expect(actual30).toEqual({
        [redisOrder.id]: 10,
        [secondRedisOrder.id]: 20,
        [redisOrderGoodTilBlockTime.id]: 30,
      });
    });

    it('gets all scores up to given, exclusive', async () => {
      const actual20: Record<string, Number> = await getOrdersAndExpiries(
        {
          latestExpiry: 20,
          latestExpiryIsInclusive: false,
        },
        client,
      );
      expect(actual20).toEqual({
        [redisOrder.id]: 10,
      });
      const actual25: Record<string, Number> = await getOrdersAndExpiries(
        {
          latestExpiry: 25,
          latestExpiryIsInclusive: false,
        },
        client,
      );
      expect(actual25).toEqual({
        [redisOrder.id]: 10,
        [secondRedisOrder.id]: 20,
      });
      const actual30: Record<string, Number> = await getOrdersAndExpiries(
        {
          latestExpiry: 30,
          latestExpiryIsInclusive: false,
        },
        client,
      );
      expect(actual30).toEqual({
        [redisOrder.id]: 10,
        [secondRedisOrder.id]: 20,
      });
    });

    it('returns empty object if no expiries exist within range', async () => {
      const actual9: Record<string, Number> = await getOrdersAndExpiries(
        {
          latestExpiry: 9,
          latestExpiryIsInclusive: true,
        },
        client,
      );
      expect(actual9).toStrictEqual({});
      const actual10: Record<string, Number> = await getOrdersAndExpiries(
        {
          latestExpiry: 10,
          latestExpiryIsInclusive: false,
        },
        client,
      );
      expect(actual10).toStrictEqual({});
    });
  });

  describe('getOrderExpiries', () => {
    beforeEach(async () => {
      await Promise.all([
        zAddAsync({ key: ORDER_EXPIRY_CACHE_KEY, value: redisOrder.id, score: 10 }, client),
        zAddAsync({ key: ORDER_EXPIRY_CACHE_KEY, value: secondRedisOrder.id, score: 20 }, client),
        zAddAsync(
          { key: ORDER_EXPIRY_CACHE_KEY, value: redisOrderGoodTilBlockTime.id, score: 30 },
          client,
        ),
      ]);
    });
    it('gets all scores up to given, inclusive', async () => {
      const actual20: string[] = await getOrderExpiries({ latestExpiry: 20 }, client);
      expect(actual20).toEqual([redisOrder.id, secondRedisOrder.id]);
      const actual25: string[] = await getOrderExpiries({ latestExpiry: 25 }, client);
      expect(actual25).toEqual([redisOrder.id, secondRedisOrder.id]);
      const actual30: string[] = await getOrderExpiries({ latestExpiry: 30 }, client);
      expect(actual30).toEqual([redisOrder.id, secondRedisOrder.id, redisOrderGoodTilBlockTime.id]);
    });

    it('gets all scores up to given, exclusive', async () => {
      const actual20: string[] = await getOrderExpiries(
        {
          latestExpiry: 20,
          latestExpiryIsInclusive: false,
        },
        client,
      );
      expect(actual20).toEqual([redisOrder.id]);
      const actual25: string[] = await getOrderExpiries(
        {
          latestExpiry: 25,
          latestExpiryIsInclusive: false,
        },
        client,
      );
      expect(actual25).toEqual([redisOrder.id, secondRedisOrder.id]);
      const actual30: string[] = await getOrderExpiries(
        {
          latestExpiry: 30,
          latestExpiryIsInclusive: false,
        },
        client,
      );
      expect(actual30).toEqual([redisOrder.id, secondRedisOrder.id]);
    });

    it('returns empty object if no expiries exist within range', async () => {
      const actual9: string[] = await getOrderExpiries(
        {
          latestExpiry: 9,
          latestExpiryIsInclusive: true,
        },
        client,
      );
      expect(actual9).toStrictEqual([]);
      const actual10: string[] = await getOrderExpiries(
        {
          latestExpiry: 10,
          latestExpiryIsInclusive: false,
        },
        client,
      );
      expect(actual10).toStrictEqual([]);
    });
  });
});
