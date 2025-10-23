import { redis } from '@dydxprotocol-indexer/redis';
import { IncomingHttpHeaders } from 'http';
import config from '../../src/config';
import { rateLimiterMiddleware } from '../../src/lib/rate-limit';
import { ratelimitRedis, getReqRateLimiter } from '../../src/caches/rate-limiters';
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

  it('consumes points for external request', async () => {
    req.headers = defaultHeaders;

    await rateLimiterMiddleware(getReqRateLimiter)(req, res, next);

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

    await rateLimiterMiddleware(getReqRateLimiter)(req, res, next);
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
      await rateLimiterMiddleware(getReqRateLimiter)(req, res, next);
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
});
