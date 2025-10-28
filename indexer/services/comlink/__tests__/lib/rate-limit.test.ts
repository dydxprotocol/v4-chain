import { redis } from '@dydxprotocol-indexer/redis';
import { IncomingHttpHeaders } from 'http';
import config from '../../src/config';
import { rateLimiterMiddleware } from '../../src/lib/rate-limit';
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
} from '../../src/caches/rate-limiters';
import * as utils from '../../src/lib/utils';

describe('rateLimit', () => {
  const defaultHeaders: IncomingHttpHeaders = { 'cf-connecting-ip': '0.0.0.0' };
  let isIndexerIpSpy: jest.SpyInstance;
  let req: any;
  let res: any;
  let next: any;

  beforeAll(() => {
    config.RATE_LIMIT_ENABLED = true;
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
