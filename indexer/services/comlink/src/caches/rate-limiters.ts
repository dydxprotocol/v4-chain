import { redis as redisLib } from '@dydxprotocol-indexer/redis';
import { RateLimiterRedis } from 'rate-limiter-flexible';
import { RedisClient } from 'redis';

import config from '../config';

export const ratelimitRedis: {
  client: RedisClient,
  connect: () => Promise<void>
} = redisLib.createRedisClient(
  config.RATE_LIMIT_REDIS_URL, config.REDIS_RECONNECT_TIMEOUT_MS,
);

// Generic rate limiter for all GET requests, limits per IP
export const getReqRateLimiter: RateLimiterRedis = new RateLimiterRedis({
  storeClient: ratelimitRedis.client,
  points: config.RATE_LIMIT_GET_POINTS,
  duration: config.RATE_LIMIT_GET_DURATION_SECONDS,
  keyPrefix: `${config.SERVICE_NAME}/get`,
});

// Rate-limiter for /screen endpoint querying a compliance provider, limits per IP
export const screenProviderLimiter: RateLimiterRedis = new RateLimiterRedis({
  storeClient: ratelimitRedis.client,
  points: config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_POINTS,
  duration: config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_DURATION_SECONDS,
  keyPrefix: `${config.SERVICE_NAME}/screen_providers`,
});

// Rate-limiter for /screen endpoint querying a compliance provider, limits the total calls made
// across all IPs
export const screenProviderGlobalLimiter: RateLimiterRedis = new RateLimiterRedis({
  storeClient: ratelimitRedis.client,
  points: config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_POINTS,
  duration: config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_DURATION_SECONDS,
  keyPrefix: `${config.SERVICE_NAME}/screen_providers_global`,
});
