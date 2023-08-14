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

export const getReqRateLimiter: RateLimiterRedis = new RateLimiterRedis({
  storeClient: ratelimitRedis.client,
  points: config.RATE_LIMIT_GET_POINTS,
  duration: config.RATE_LIMIT_GET_DURATION_SECONDS,
  keyPrefix: `${config.SERVICE_NAME}/get`,
});
