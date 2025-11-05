import { redis as redisLib } from '@dydxprotocol-indexer/redis';
import { RateLimiterRedis } from 'rate-limiter-flexible';
import { RedisClient } from 'redis';

import config from '../config';

export const ratelimitRedis: {
  client: RedisClient,
  connect: () => Promise<void>,
} = redisLib.createRedisClient(
  config.RATE_LIMIT_REDIS_URL, config.REDIS_RECONNECT_TIMEOUT_MS,
);

export function getDefaultRateLimiter(): RateLimiterRedis {
  return new RateLimiterRedis({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_GET_POINTS,
    duration: config.RATE_LIMIT_GET_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/get`,
  });
}

// Generic rate limiter for all GET requests, limits per IP
export const defaultRateLimiter: RateLimiterRedis = getDefaultRateLimiter();

// Rate-limiter for /screen endpoint querying a compliance provider, limits per IP
export function getScreenProviderLimiter(): RateLimiterRedis {
  return new RateLimiterRedis({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_POINTS,
    duration: config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/screen_providers`,
  });
}
export const screenProviderLimiter: RateLimiterRedis = getScreenProviderLimiter();

// Rate-limiter for /screen endpoint querying a compliance provider, limits the total calls made
// across all IPs
export function getScreenProviderGlobalLimiter(): RateLimiterRedis {
  return new RateLimiterRedis({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_POINTS,
    duration: config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/screen_providers_global`,
  });
}
export const screenProviderGlobalLimiter: RateLimiterRedis = getScreenProviderGlobalLimiter();

// Rate-limiter for /orders endpoint, limits per IP
export function getOrdersRateLimiter(): RateLimiterRedis {
  return new RateLimiterRedis({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_ORDERS_POINTS,
    duration: config.RATE_LIMIT_ORDERS_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/orders`,
  });
}
export const ordersRateLimiter: RateLimiterRedis = getOrdersRateLimiter();

// Rate-limiter for /fills endpoint, limits per IP
export function getFillsRateLimiter(): RateLimiterRedis {
  return new RateLimiterRedis({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_FILLS_POINTS,
    duration: config.RATE_LIMIT_FILLS_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/fills`,
  });
}
export const fillsRateLimiter: RateLimiterRedis = getFillsRateLimiter();

// Rate-limiter for /candles endpoint, limits per IP
export function getCandlesRateLimiter(): RateLimiterRedis {
  return new RateLimiterRedis({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_CANDLES_POINTS,
    duration: config.RATE_LIMIT_CANDLES_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/candles`,
  });
}
export const candlesRateLimiter: RateLimiterRedis = getCandlesRateLimiter();

// Rate-limiter for /sparklines endpoint, limits per IP
export function getSparklinesRateLimiter(): RateLimiterRedis {
  return new RateLimiterRedis({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_SPARKLINES_POINTS,
    duration: config.RATE_LIMIT_SPARKLINES_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/sparklines`,
  });
}
export const sparklinesRateLimiter: RateLimiterRedis = getSparklinesRateLimiter();

// Rate-limiter for /historical-pnl endpoint, limits per IP
export function getHistoricalPnlRateLimiter(): RateLimiterRedis {
  return new RateLimiterRedis({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_HISTORICAL_PNL_POINTS,
    duration: config.RATE_LIMIT_HISTORICAL_PNL_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/historical_pnl`,
  });
}
export const historicalPnlRateLimiter: RateLimiterRedis = getHistoricalPnlRateLimiter();

// Rate-limiter for /pnl endpoint, limits per IP
export function getPnlRateLimiter(): RateLimiterRedis {
  return new RateLimiterRedis({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_PNL_POINTS,
    duration: config.RATE_LIMIT_PNL_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/pnl`,
  });
}
export const pnlRateLimiter: RateLimiterRedis = getPnlRateLimiter();

// Rate-limiter for /funding endpoint, limits per IP
export function getFundingRateLimiter(): RateLimiterRedis {
  return new RateLimiterRedis({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_FUNDING_POINTS,
    duration: config.RATE_LIMIT_FUNDING_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/funding`,
  });
}
export const fundingRateLimiter: RateLimiterRedis = getFundingRateLimiter();
