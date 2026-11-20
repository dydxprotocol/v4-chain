import { stats } from '@dydxprotocol-indexer/base';
import { redis } from '@dydxprotocol-indexer/redis';
import { IncomingHttpHeaders } from 'http';
import { createClient } from 'redis';
import config from '../../src/config';
import { getPointCost, rateLimiterMiddleware } from '../../src/lib/rate-limit';
import {
  ratelimitRedis,
  defaultRateLimiter,
  ordersRateLimiter,
  fillsRateLimiter,
  candlesRateLimiter,
  sparklinesRateLimiter,
  historicalPnlRateLimiter,
  pnlRateLimiter,
  fundingRateLimiter,
  getDefaultRateLimiter,
  getOrdersRateLimiter,
  getFillsRateLimiter,
  getCandlesRateLimiter,
  getSparklinesRateLimiter,
  getHistoricalPnlRateLimiter,
  getPnlRateLimiter,
  getFundingRateLimiter,
  turnkeyEmailRateLimiter,
  turnkeyIpRateLimiter,
  turnkeyNewSuborgGlobalLimiter,
} from '../../src/caches/rate-limiters';
import { TurnkeyError } from '../../src/lib/errors';
import { reserveNewSuborgSlot } from '../../src/lib/turnkey-helpers';
import * as utils from '../../src/lib/utils';
import { ipRateLimitKey } from '../../src/lib/utils';

describe('rateLimit', () => {
  const defaultHeaders: IncomingHttpHeaders = { 'cf-connecting-ip': '0.0.0.0' };
  let isIndexerIpSpy: jest.SpyInstance;
  let req: any;
  let res: any;
  let next: any;

  beforeAll(async () => {
    config.RATE_LIMIT_ENABLED = true;
    // Limiters skip the store until the client is ready, so wait for the connection.
    await ratelimitRedis.connect();
  });

  beforeEach(() => {
    isIndexerIpSpy = jest.spyOn(utils, 'isIndexerIp');
    req = {
      get: jest.fn().mockReturnThis(),
    };
    res = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn().mockReturnThis(),
      set: jest.fn().mockReturnThis(),
    };
    next = jest.fn();
  });

  afterEach(async () => {
    await redis.deleteAllAsync(ratelimitRedis.client);
    jest.restoreAllMocks();
  });

  afterAll(() => {
    config.RATE_LIMIT_ENABLED = false;
  });

  describe('getReqRateLimiter', () => {
    it('consumes points for external request', async () => {
      req.headers = defaultHeaders;

      await rateLimiterMiddleware(defaultRateLimiter)(req, res, next);

      expect(res.status).not.toHaveBeenCalled();
      expect(res.set).toHaveBeenCalledTimes(1);
      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': config.RATE_LIMIT_GET_POINTS - 1,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_GET_POINTS,
      });
      expect(next).toHaveBeenCalledTimes(1);
    });

    it('consumes no points for internal request', async () => {
      isIndexerIpSpy.mockReturnValueOnce(true);
      req.headers = defaultHeaders;

      await rateLimiterMiddleware(defaultRateLimiter)(req, res, next);
      expect(res.status).not.toHaveBeenCalled();
      expect(res.set).toHaveBeenCalledTimes(1);
      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': config.RATE_LIMIT_GET_POINTS,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_GET_POINTS,
      });
      expect(next).toHaveBeenCalledTimes(1);
    });

    it('sets response code to 429 if rate limit exceeded', async () => {
      req.headers = defaultHeaders;
      for (let i = 0; i < config.RATE_LIMIT_GET_POINTS + 1; i++) {
        await rateLimiterMiddleware(defaultRateLimiter)(req, res, next);
      }

      expect(res.status).toHaveBeenCalledTimes(1);
      expect(res.status).toHaveBeenNthCalledWith(1, 429);
      expect(res.set).toHaveBeenCalledTimes(config.RATE_LIMIT_GET_POINTS + 1);
      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': 0,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_GET_POINTS,
      });
      expect(next).toHaveBeenCalledTimes(config.RATE_LIMIT_GET_POINTS);
    });

    it('respects config changes to rate limit points', async () => {
      const originalPoints = config.RATE_LIMIT_GET_POINTS;
      config.RATE_LIMIT_GET_POINTS = 5;

      req.headers = defaultHeaders;

      await rateLimiterMiddleware(getDefaultRateLimiter())(req, res, next);

      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': expect.any(Number),
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_GET_POINTS,
      });

      config.RATE_LIMIT_GET_POINTS = originalPoints;
    });
  });

  describe('ordersRateLimiter', () => {
    it('consumes points for external request', async () => {
      req.headers = defaultHeaders;

      await rateLimiterMiddleware(ordersRateLimiter)(req, res, next);

      expect(res.status).not.toHaveBeenCalled();
      expect(res.set).toHaveBeenCalledTimes(1);
      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': config.RATE_LIMIT_ORDERS_POINTS - 1,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_ORDERS_POINTS,
      });
      expect(next).toHaveBeenCalledTimes(1);
    });

    it('sets response code to 429 if rate limit exceeded', async () => {
      req.headers = defaultHeaders;
      for (let i = 0; i < config.RATE_LIMIT_ORDERS_POINTS + 1; i++) {
        await rateLimiterMiddleware(ordersRateLimiter)(req, res, next);
      }

      expect(res.status).toHaveBeenCalledTimes(1);
      expect(res.status).toHaveBeenNthCalledWith(1, 429);
      expect(next).toHaveBeenCalledTimes(config.RATE_LIMIT_ORDERS_POINTS);
    });

    it('respects config changes to orders rate limit', async () => {
      const originalPoints = config.RATE_LIMIT_ORDERS_POINTS;
      config.RATE_LIMIT_ORDERS_POINTS = 10;

      req.headers = defaultHeaders;

      await rateLimiterMiddleware(getOrdersRateLimiter())(req, res, next);

      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': expect.any(Number),
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_ORDERS_POINTS,
      });

      config.RATE_LIMIT_ORDERS_POINTS = originalPoints;
    });
  });

  describe('fillsRateLimiter', () => {
    it('consumes points for external request', async () => {
      req.headers = defaultHeaders;

      await rateLimiterMiddleware(fillsRateLimiter)(req, res, next);

      expect(res.status).not.toHaveBeenCalled();
      expect(res.set).toHaveBeenCalledTimes(1);
      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': config.RATE_LIMIT_FILLS_POINTS - 1,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_FILLS_POINTS,
      });
      expect(next).toHaveBeenCalledTimes(1);
    });

    it('sets response code to 429 if rate limit exceeded', async () => {
      req.headers = defaultHeaders;
      for (let i = 0; i < config.RATE_LIMIT_FILLS_POINTS + 1; i++) {
        await rateLimiterMiddleware(fillsRateLimiter)(req, res, next);
      }

      expect(res.status).toHaveBeenCalledTimes(1);
      expect(res.status).toHaveBeenNthCalledWith(1, 429);
      expect(next).toHaveBeenCalledTimes(config.RATE_LIMIT_FILLS_POINTS);
    });

    it('respects config changes to fills rate limit', async () => {
      const originalPoints = config.RATE_LIMIT_FILLS_POINTS;
      config.RATE_LIMIT_FILLS_POINTS = 15;

      req.headers = defaultHeaders;

      await rateLimiterMiddleware(getFillsRateLimiter())(req, res, next);

      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': expect.any(Number),
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_FILLS_POINTS,
      });

      config.RATE_LIMIT_FILLS_POINTS = originalPoints;
    });
  });

  describe('candlesRateLimiter', () => {
    it('consumes points for external request', async () => {
      req.headers = defaultHeaders;

      await rateLimiterMiddleware(candlesRateLimiter)(req, res, next);

      expect(res.status).not.toHaveBeenCalled();
      expect(res.set).toHaveBeenCalledTimes(1);
      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': config.RATE_LIMIT_CANDLES_POINTS - 1,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_CANDLES_POINTS,
      });
      expect(next).toHaveBeenCalledTimes(1);
    });

    it('sets response code to 429 if rate limit exceeded', async () => {
      req.headers = defaultHeaders;
      for (let i = 0; i < config.RATE_LIMIT_CANDLES_POINTS + 1; i++) {
        await rateLimiterMiddleware(candlesRateLimiter)(req, res, next);
      }

      expect(res.status).toHaveBeenCalledTimes(1);
      expect(res.status).toHaveBeenNthCalledWith(1, 429);
      expect(next).toHaveBeenCalledTimes(config.RATE_LIMIT_CANDLES_POINTS);
    });

    it('respects config changes to candles rate limit', async () => {
      const originalPoints = config.RATE_LIMIT_CANDLES_POINTS;
      config.RATE_LIMIT_CANDLES_POINTS = 500;

      req.headers = defaultHeaders;

      await rateLimiterMiddleware(getCandlesRateLimiter())(req, res, next);

      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': expect.any(Number),
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_CANDLES_POINTS,
      });

      config.RATE_LIMIT_CANDLES_POINTS = originalPoints;
    });
  });

  describe('sparklinesRateLimiter', () => {
    it('consumes points for external request', async () => {
      req.headers = defaultHeaders;

      await rateLimiterMiddleware(sparklinesRateLimiter)(req, res, next);

      expect(res.status).not.toHaveBeenCalled();
      expect(res.set).toHaveBeenCalledTimes(1);
      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': config.RATE_LIMIT_SPARKLINES_POINTS - 1,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_SPARKLINES_POINTS,
      });
      expect(next).toHaveBeenCalledTimes(1);
    });

    it('sets response code to 429 if rate limit exceeded', async () => {
      req.headers = defaultHeaders;
      for (let i = 0; i < config.RATE_LIMIT_SPARKLINES_POINTS + 1; i++) {
        await rateLimiterMiddleware(sparklinesRateLimiter)(req, res, next);
      }

      expect(res.status).toHaveBeenCalledTimes(1);
      expect(res.status).toHaveBeenNthCalledWith(1, 429);
      expect(next).toHaveBeenCalledTimes(config.RATE_LIMIT_SPARKLINES_POINTS);
    });

    it('respects config changes to sparklines rate limit', async () => {
      const originalPoints = config.RATE_LIMIT_SPARKLINES_POINTS;
      config.RATE_LIMIT_SPARKLINES_POINTS = 50;

      req.headers = defaultHeaders;

      await rateLimiterMiddleware(getSparklinesRateLimiter())(req, res, next);

      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': expect.any(Number),
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_SPARKLINES_POINTS,
      });

      config.RATE_LIMIT_SPARKLINES_POINTS = originalPoints;
    });
  });

  describe('historicalPnlRateLimiter', () => {
    it('consumes points for external request', async () => {
      req.headers = defaultHeaders;

      await rateLimiterMiddleware(historicalPnlRateLimiter)(req, res, next);

      expect(res.status).not.toHaveBeenCalled();
      expect(res.set).toHaveBeenCalledTimes(1);
      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': config.RATE_LIMIT_HISTORICAL_PNL_POINTS - 1,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_HISTORICAL_PNL_POINTS,
      });
      expect(next).toHaveBeenCalledTimes(1);
    });

    it('sets response code to 429 if rate limit exceeded', async () => {
      req.headers = defaultHeaders;
      for (let i = 0; i < config.RATE_LIMIT_HISTORICAL_PNL_POINTS + 1; i++) {
        await rateLimiterMiddleware(historicalPnlRateLimiter)(req, res, next);
      }

      expect(res.status).toHaveBeenCalledTimes(1);
      expect(res.status).toHaveBeenNthCalledWith(1, 429);
      expect(next).toHaveBeenCalledTimes(config.RATE_LIMIT_HISTORICAL_PNL_POINTS);
    });

    it('respects config changes to historical pnl rate limit', async () => {
      const originalPoints = config.RATE_LIMIT_HISTORICAL_PNL_POINTS;
      config.RATE_LIMIT_HISTORICAL_PNL_POINTS = 25;

      req.headers = defaultHeaders;

      await rateLimiterMiddleware(getHistoricalPnlRateLimiter())(req, res, next);

      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': expect.any(Number),
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_HISTORICAL_PNL_POINTS,
      });

      config.RATE_LIMIT_HISTORICAL_PNL_POINTS = originalPoints;
    });
  });

  describe('pnlRateLimiter', () => {
    it('consumes points for external request', async () => {
      req.headers = defaultHeaders;

      await rateLimiterMiddleware(pnlRateLimiter)(req, res, next);

      expect(res.status).not.toHaveBeenCalled();
      expect(res.set).toHaveBeenCalledTimes(1);
      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': config.RATE_LIMIT_PNL_POINTS - 1,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_PNL_POINTS,
      });
      expect(next).toHaveBeenCalledTimes(1);
    });

    it('sets response code to 429 if rate limit exceeded', async () => {
      req.headers = defaultHeaders;
      for (let i = 0; i < config.RATE_LIMIT_PNL_POINTS + 1; i++) {
        await rateLimiterMiddleware(pnlRateLimiter)(req, res, next);
      }

      expect(res.status).toHaveBeenCalledTimes(1);
      expect(res.status).toHaveBeenNthCalledWith(1, 429);
      expect(next).toHaveBeenCalledTimes(config.RATE_LIMIT_PNL_POINTS);
    });

    it('respects config changes to pnl rate limit', async () => {
      const originalPoints = config.RATE_LIMIT_PNL_POINTS;
      config.RATE_LIMIT_PNL_POINTS = 30;

      req.headers = defaultHeaders;

      await rateLimiterMiddleware(getPnlRateLimiter())(req, res, next);

      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': expect.any(Number),
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_PNL_POINTS,
      });

      config.RATE_LIMIT_PNL_POINTS = originalPoints;
    });
  });

  describe('fundingRateLimiter', () => {
    it('consumes points for external request', async () => {
      req.headers = defaultHeaders;

      await rateLimiterMiddleware(fundingRateLimiter)(req, res, next);

      expect(res.status).not.toHaveBeenCalled();
      expect(res.set).toHaveBeenCalledTimes(1);
      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': config.RATE_LIMIT_FUNDING_POINTS - 1,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_FUNDING_POINTS,
      });
      expect(next).toHaveBeenCalledTimes(1);
    });

    it('sets response code to 429 if rate limit exceeded', async () => {
      req.headers = defaultHeaders;
      for (let i = 0; i < config.RATE_LIMIT_FUNDING_POINTS + 1; i++) {
        await rateLimiterMiddleware(fundingRateLimiter)(req, res, next);
      }

      expect(res.status).toHaveBeenCalledTimes(1);
      expect(res.status).toHaveBeenNthCalledWith(1, 429);
      expect(next).toHaveBeenCalledTimes(config.RATE_LIMIT_FUNDING_POINTS);
    });

    it('respects config changes to funding rate limit', async () => {
      const originalPoints = config.RATE_LIMIT_FUNDING_POINTS;
      config.RATE_LIMIT_FUNDING_POINTS = 20;

      req.headers = defaultHeaders;

      await rateLimiterMiddleware(getFundingRateLimiter())(req, res, next);

      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': expect.any(Number),
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_FUNDING_POINTS,
      });

      config.RATE_LIMIT_FUNDING_POINTS = originalPoints;
    });
  });

  describe('getPointCost offset weighting', () => {
    const ip: string = '0.0.0.0';

    it('costs 1 point when the request has no page/limit in the query', () => {
      expect(getPointCost(ip, { query: undefined } as any)).toEqual(1);
      expect(getPointCost(ip, { query: {} } as any)).toEqual(1);
    });

    it('costs 1 point for page 1 (matches every pre-existing assertion in this file)', () => {
      expect(getPointCost(ip, { query: { page: '1', limit: '1000' } } as any)).toEqual(1);
    });

    it('costs 1 point for a shallow page regardless of limit', () => {
      expect(getPointCost(ip, { query: { page: '2', limit: '100' } } as any)).toEqual(1);
    });

    it('weights cost by offset depth: 1 + floor(offset / API_LIMIT_V4)', () => {
      // offset = (50 - 1) * 1000 = 49,000 -> 1 + floor(49000 / 1000) = 50
      expect(getPointCost(ip, { query: { page: '50', limit: '1000' } } as any)).toEqual(50);
      // offset = (2 - 1) * 1000 = 1,000 -> 1 + floor(1000 / 1000) = 2
      expect(getPointCost(ip, { query: { page: '2', limit: '1000' } } as any)).toEqual(2);
    });

    it('defaults limit to API_LIMIT_V4 when limit is omitted', () => {
      // offset = (10 - 1) * API_LIMIT_V4 -> 1 + floor(offset / API_LIMIT_V4) = 1 + 9 = 10
      expect(getPointCost(ip, { query: { page: '10' } } as any)).toEqual(10);
    });

    it('never throws and falls back to 1 point on adversarial/garbage page values', () => {
      expect(getPointCost(ip, { query: { page: 'not-a-number', limit: '1000' } } as any))
        .toEqual(1);
      expect(getPointCost(ip, { query: { page: '-5', limit: '1000' } } as any)).toEqual(1);
      expect(getPointCost(ip, { query: { page: 'Infinity', limit: '1000' } } as any)).toEqual(1);
      expect(getPointCost(ip, { query: { page: '10', limit: '0' } } as any)).toEqual(1);
      expect(getPointCost(ip, { query: { page: '10', limit: '-100' } } as any)).toEqual(1);
    });

    it('costs 0 for internal IPs regardless of page depth', () => {
      isIndexerIpSpy.mockReturnValueOnce(true);
      expect(getPointCost(ip, { query: { page: '10000', limit: '1000' } } as any)).toEqual(0);
    });

    it('drains the rate-limit budget faster for deep pagination end-to-end', async () => {
      req.headers = defaultHeaders;
      req.query = { page: '50', limit: '1000' };

      await rateLimiterMiddleware(fillsRateLimiter)(req, res, next);

      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': config.RATE_LIMIT_FILLS_POINTS - 50,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_FILLS_POINTS,
      });
      expect(next).toHaveBeenCalledTimes(1);
    });

    it('clamps an astronomically large offset-derived cost to the limiter capacity', async () => {
      // An unclamped cost this large (offset ~1e29, so 1 + floor(offset/API_LIMIT_V4) is on the
      // order of 1e26) serializes to exponential notation, which Redis's INCRBY rejects as "not
      // an integer" - rateLimiterMiddleware's catch block treats that rejection as a plain Error
      // and falls through to next(), silently disabling rate limiting entirely for the request.
      // Clamping to the limiter's own point capacity avoids ever sending such a value to Redis.
      req.headers = defaultHeaders;
      req.query = { page: '2', limit: '99999999999999999999999999999' };

      await rateLimiterMiddleware(fillsRateLimiter)(req, res, next);

      expect(res.set).toHaveBeenCalledWith({
        'RateLimit-Remaining': 0,
        'RateLimit-Reset': expect.any(Number),
        'RateLimit-Limit': config.RATE_LIMIT_FILLS_POINTS,
      });
      expect(res.status).not.toHaveBeenCalled();
      expect(next).toHaveBeenCalledTimes(1);
    });
  });

  describe('rate limiter duration configuration', () => {
    it('respects config changes to duration for get requests', async () => {
      const originalDuration = config.RATE_LIMIT_GET_DURATION_SECONDS;
      config.RATE_LIMIT_GET_DURATION_SECONDS = 5;

      req.headers = defaultHeaders;

      await rateLimiterMiddleware(getDefaultRateLimiter())(req, res, next);

      const resetTime = res.set.mock.calls[0][0]['RateLimit-Reset'];
      const currentTime = Date.now();
      const expectedResetWindow = config.RATE_LIMIT_GET_DURATION_SECONDS * 1000;

      expect(resetTime - currentTime).toBeLessThanOrEqual(expectedResetWindow);
      expect(resetTime - currentTime).toBeGreaterThan(0);

      config.RATE_LIMIT_GET_DURATION_SECONDS = originalDuration;
    });

    it('respects config changes to duration for orders', async () => {
      const originalDuration = config.RATE_LIMIT_ORDERS_DURATION_SECONDS;
      config.RATE_LIMIT_ORDERS_DURATION_SECONDS = 20;

      req.headers = defaultHeaders;

      await rateLimiterMiddleware(getOrdersRateLimiter())(req, res, next);

      const resetTime = res.set.mock.calls[0][0]['RateLimit-Reset'];
      const currentTime = Date.now();
      const expectedResetWindow = config.RATE_LIMIT_ORDERS_DURATION_SECONDS * 1000;

      expect(resetTime - currentTime).toBeLessThanOrEqual(expectedResetWindow);
      expect(resetTime - currentTime).toBeGreaterThan(0);

      config.RATE_LIMIT_ORDERS_DURATION_SECONDS = originalDuration;
    });
  });

  describe('ipRateLimitKey', () => {
    it.each([
      ['1.2.3.4', '1.2.3.4'],
      ['::ffff:1.2.3.4', '1.2.3.4'],
      ['::ffff:102:304', '1.2.3.4'],
      ['2a03:4000:13:1a88:4403:f7ff:fe05:9e08', '2a03:4000:13:1a88::/64'],
      ['2a03:4000:13:1a88::1', '2a03:4000:13:1a88::/64'],
      ['2001:db8::1', '2001:db8:0:0::/64'],
      ['2001:0DB8:0000:0001:0:0:0:1', '2001:db8:0:1::/64'],
    ])('maps %s to %s', (ip, key) => {
      expect(ipRateLimitKey(ip)).toBe(key);
    });

    it.each([
      ['2a03:4000:13:1a88:4403:f7ff:fe05:9e08', '2a03:4000:13::/48'],
      ['2a03:4000:13:ffff::1', '2a03:4000:13::/48'],
      ['1.2.3.4', '1.2.3.4'],
    ])('maps %s to %s at /48', (ip, key) => {
      expect(ipRateLimitKey(ip, 48)).toBe(key);
    });
  });

  describe('custom key', () => {
    it('shares one bucket across a /64 and blocks past the limit', async () => {
      const limit: number = config.RATE_LIMIT_TURNKEY_IP_POINTS;
      const getKey = (_r: any, ip: string) => ipRateLimitKey(ip);
      for (let i = 0; i < limit; i += 1) {
        req.headers = { 'cf-connecting-ip': `2a03:4000:13:1a88::${i + 1}` };
        await rateLimiterMiddleware(turnkeyIpRateLimiter, undefined, getKey)(req, res, next);
      }
      expect(next).toHaveBeenCalledTimes(limit);

      req.headers = { 'cf-connecting-ip': '2a03:4000:13:1a88:ffff::1' };
      await rateLimiterMiddleware(turnkeyIpRateLimiter, undefined, getKey)(req, res, next);
      expect(res.status).toHaveBeenCalledWith(429);
      expect(next).toHaveBeenCalledTimes(limit);
    });

    it('skips the limiter when the key function returns undefined', async () => {
      req.headers = defaultHeaders;
      await rateLimiterMiddleware(turnkeyEmailRateLimiter, undefined, () => undefined)(
        req, res, next,
      );
      expect(res.set).not.toHaveBeenCalled();
      expect(next).toHaveBeenCalledTimes(1);
    });
  });

  describe('reserveNewSuborgSlot', () => {
    it('refuses once the global cap is used up', async () => {
      const increment = jest.spyOn(stats, 'increment');
      for (let i = 0; i < config.RATE_LIMIT_TURNKEY_NEW_SUBORG_GLOBAL_POINTS; i += 1) {
        await reserveNewSuborgSlot();
      }
      await expect(reserveNewSuborgSlot()).rejects.toThrow(TurnkeyError);
      expect(increment).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.turnkey_new_suborg_rejected`, 1,
        { reason: 'cap_reached', limit: 'global' },
      );
    });

    it('limits new sub-orgs per caller before the global cap', async () => {
      for (let i = 0; i < config.RATE_LIMIT_TURNKEY_NEW_SUBORG_IP_POINTS; i += 1) {
        await reserveNewSuborgSlot(`2a03:4000:13:1a88::${i + 1}`);
      }
      const increment = jest.spyOn(stats, 'increment');
      await expect(reserveNewSuborgSlot('2a03:4000:13:1a88::ff')).rejects.toThrow(TurnkeyError);
      // Another /64 in the same /48 shares the budget.
      await expect(reserveNewSuborgSlot('2a03:4000:13:ffff::1')).rejects.toThrow(TurnkeyError);
      expect(increment).toHaveBeenCalledWith(
        `${config.SERVICE_NAME}.turnkey_new_suborg_rejected`, 1,
        { reason: 'cap_reached', limit: 'ip' },
      );
      await expect(reserveNewSuborgSlot('203.0.113.7')).resolves.toBeUndefined();
    });

    it('refuses at once, without touching the store, when Redis is not ready', async () => {
      const consume = jest.spyOn(turnkeyNewSuborgGlobalLimiter, 'consume');
      const client: any = ratelimitRedis.client;
      const ready: boolean = client.ready;
      client.ready = false;
      try {
        await expect(reserveNewSuborgSlot()).rejects.toThrow(TurnkeyError);
        expect(consume).not.toHaveBeenCalled();
      } finally {
        client.ready = ready;
      }
    });

    it('fails closed when the rate-limit store errors', async () => {
      jest.spyOn(turnkeyNewSuborgGlobalLimiter, 'consume').mockRejectedValueOnce(new Error('redis down'));
      await expect(reserveNewSuborgSlot()).rejects.toThrow(TurnkeyError);
    });
  });

  describe('rate-limit store not ready', () => {
    it('refuses with 503, without touching the store, when the route opts in', async () => {
      const client: any = ratelimitRedis.client;
      const ready: boolean = client.ready;
      const consume = jest.spyOn(defaultRateLimiter, 'consume');
      const increment = jest.spyOn(stats, 'increment');
      client.ready = false;
      req.headers = defaultHeaders;
      try {
        await rateLimiterMiddleware(defaultRateLimiter, undefined, undefined, true)(req, res, next);
        expect(res.status).toHaveBeenCalledWith(503);
        expect(increment).toHaveBeenCalledWith(
          `${config.SERVICE_NAME}.rate_limit_store_not_ready`, 1,
          { key_prefix: defaultRateLimiter.keyPrefix },
        );
        expect(next).not.toHaveBeenCalled();
        expect(consume).not.toHaveBeenCalled();
      } finally {
        client.ready = ready;
      }
    });

    it('still uses the store when the route does not opt in', async () => {
      const client: any = ratelimitRedis.client;
      const ready: boolean = client.ready;
      // Mocked: with ready=false a real command would wait in the client's offline queue.
      const consume = jest.spyOn(defaultRateLimiter, 'consume')
        .mockResolvedValue({ remainingPoints: 1, msBeforeNext: 1000 } as any);
      client.ready = false;
      req.headers = defaultHeaders;
      try {
        await rateLimiterMiddleware(defaultRateLimiter)(req, res, next);
        expect(consume).toHaveBeenCalledTimes(1);
      } finally {
        client.ready = ready;
      }
    });
  });

  describe('rate-limit store reconnects', () => {
    const waitFor = async (check: () => boolean, ms: number): Promise<boolean> => {
      const deadline: number = Date.now() + ms;
      while (Date.now() < deadline) {
        if (check()) {
          return true;
        }
        await new Promise((resolve) => setTimeout(resolve, 20));
      }
      return check();
    };
    const unreachableClient = (connectTimeout: number): any => {
      const client: any = createClient({ host: '127.0.0.1', port: 1, retry_strategy: () => 10 });
      client.on('error', () => {});
      client.connect_timeout = connectTimeout;
      return client;
    };

    it('never gives up on the rate-limit client', () => {
      expect((ratelimitRedis.client as any).connect_timeout).toBe(Infinity);
    });

    it('node-redis gives up once retries add up to connect_timeout, unless it is Infinity', async () => {
      const capped: any = unreachableClient(50);
      const uncapped: any = unreachableClient(Infinity);
      try {
        expect(await waitFor(() => capped.closing, 2000)).toBe(true);
        expect(await waitFor(() => uncapped.closing, 300)).toBe(false);
      } finally {
        capped.end(false);
        uncapped.end(false);
      }
    });
  });

  describe('internal vs external IP handling', () => {
    it('treats internal IPs differently across all rate limiters', async () => {
      isIndexerIpSpy.mockReturnValue(true);
      req.headers = defaultHeaders;

      const rateLimiters = [
        defaultRateLimiter,
        ordersRateLimiter,
        fillsRateLimiter,
        candlesRateLimiter,
        sparklinesRateLimiter,
        historicalPnlRateLimiter,
        pnlRateLimiter,
        fundingRateLimiter,
      ];

      for (const limiter of rateLimiters) {
        await rateLimiterMiddleware(limiter)(req, res, next);

        // Internal IPs should not consume points
        const remainingPoints = res.set.mock.calls[res.set.mock.calls.length - 1][0]['RateLimit-Remaining'];
        const limitPoints = res.set.mock.calls[res.set.mock.calls.length - 1][0]['RateLimit-Limit'];

        expect(remainingPoints).toBe(limitPoints);
      }
    });
  });
});
