import { redis as redisLib } from '@dydxprotocol-indexer/redis';
import {
  RedisClient,
} from 'redis';

import config from '../../config';

// Primary read-write Redis client
const res: {
  client: RedisClient,
  connect: () => Promise<void>,
} = redisLib.createRedisClient(config.REDIS_URL, config.REDIS_RECONNECT_TIMEOUT_MS);

export const redisClient: RedisClient = res.client;
export const connect = res.connect;

// Read-only Redis client
const resReadOnly: {
  client: RedisClient,
  connect: () => Promise<void>,
} = redisLib.createRedisClient(config.REDIS_READONLY_URL, config.REDIS_RECONNECT_TIMEOUT_MS);

export const redisReadOnlyClient: RedisClient = resReadOnly.client;
export const connectReadOnly = resReadOnly.connect;
