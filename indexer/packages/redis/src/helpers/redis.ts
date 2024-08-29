import { promisify } from 'util';

import {
  InfoObject,
  logger,
  stats,
} from '@dydxprotocol-indexer/base';
import bluebird from 'bluebird';
import _ from 'lodash';
import redis, {
  RedisClient,
  RetryStrategyOptions,
} from 'redis';

import { allLuaScripts, loadScript } from '../caches/scripts';
import config from '../config';

bluebird.promisifyAll(redis.RedisClient.prototype);
bluebird.promisifyAll(redis.Multi.prototype);

/**
 * Create a function for redis reconnection strategy. Exported for tests.
 * @param url URL for redis, used to tag stats on which redis cluster had reconnections
 * @param reconnectTimeoutMs Time to return for reconnecting
 * @returns The value passed in by `reconnectTimeoutMs`, this indicates the number of millisoneconds
 *          to wait before the next reconnection attempt
 */
export function createRetryConnectionStrategy(
  url: string,
  reconnectTimeoutMs: number,
): (options: RetryStrategyOptions) => number {
  return (options: RetryStrategyOptions) => {
    const errorLog: InfoObject = {
      at: 'redis#connect',
      message: `Failed to connect to redis: ${options.error?.code}`,
      error: options.error,
      attempt: options.attempt,
    };
    if (options.attempt <= config.REDIS_RECONNECT_ATTEMPT_ERROR_THRESHOLD &&
        (options.error === undefined || options.error === null)
    ) {
      // Don't log errors for transient connection failures with no error event emitted
      logger.info(errorLog);
    } else {
      logger.error(errorLog);
    }
    stats.distribution('redis.connection_retries', options.attempt, { url });
    return reconnectTimeoutMs;
  };
}

export function createRedisClient(
  url: string,
  reconnectTimeoutMs: number,
): {
  client: RedisClient,
  connect: () => Promise<void>,
} {
  const redisClient: RedisClient = redis.createClient({
    url,
    retry_strategy: createRetryConnectionStrategy(url, reconnectTimeoutMs),
  });

  redisClient.on('error', (error: Error) => {
    logger.error({
      at: 'redis#onError',
      message: error.message,
      error,
    });
  });

  const connectPromise: Promise<void> = new Promise((resolve) => redisClient.on('ready', async () => {
    logger.info({
      at: 'redis#ready',
      message: 'Connected to redis. Started to load scripts.',
    });

    await Promise.all(allLuaScripts.map((script) => loadScript(script, redisClient)));
    const scriptNames: string[] = allLuaScripts.map((script) => script.name);
    logger.info({
      at: 'redis#ready',
      message: `Scripts [${scriptNames}] loaded.`,
    });

    resolve();
  }));

  const redisConnect = async () => connectPromise;

  return { client: redisClient, connect: redisConnect };
}

/**
 * Time to live of key in milliseconds.
 * Returns -2 if not set, -1 if the key exists but does not expire.
 */
export async function pttl(
  redisClient: RedisClient,
  redisKey: string,
): Promise<number> {
  const clientPTTL = promisify(redisClient.pttl).bind(redisClient);
  return clientPTTL(redisKey);
}

/**
 * Time to live of a key in seconds.
 * Returns -2 if not set, -1 if the key exists but does not expire.
 */
export async function ttl(
  redisClient: RedisClient,
  redisKey: string,
): Promise<number> {
  const clientTTL = promisify(redisClient.ttl).bind(redisClient);
  return clientTTL(redisKey);
}

/**
 * Sets the timeout of a key in seconds.
 * Returns 1 if set, 0 if not set.
 */
export async function setExpiry(
  redisClient: RedisClient,
  redisKey: string,
  lockExpirySec: number,
): Promise<number> {
  const clientExpire = promisify(redisClient.expire).bind(redisClient);
  return clientExpire(redisKey, lockExpirySec);
}

/**
 * Sets a redis key if it's not already set. Returns true if it was set by this function. Returns
 * false if it was already set.
 */
export async function lockWithExpiry(
  redisClient: RedisClient,
  redisKey: string,
  redisValue: string,
  lockExpiryMs: number,
): Promise<boolean> {
  const clientSet = promisify<string, string, string, number, string, string | null>(
    redisClient.set as unknown as (
      arg0: string,
      arg1: string,
      arg2: string,
      arg3: number,
      arg4: string,
    ) => Promise<string | null>,
  ).bind(redisClient);
  const lockResult: string | null = await clientSet(
    redisKey,
    redisValue,
    'px',
    lockExpiryMs,
    'nx',
  );
  return !!lockResult;
}

/**
 * Sets a redis key if it's not already set. Returns true if it was set by this function. Returns
 * false if it was already set.
 */
export async function lockIndefinitely(
  redisClient: RedisClient,
  redisKey: string,
  redisValue: string,
): Promise<boolean> {
  const clientSet = promisify<string, string, string, string | null>(
    redisClient.set as unknown as (
      arg0: string,
      arg1: string,
      arg2: string,
    ) => Promise<string | null>,
  ).bind(redisClient);
  const lockResult: string | null = await clientSet(
    redisKey,
    redisValue,
    'nx',
  );
  return !!lockResult;
}

/**
 * Deletes a redis key. If redisValue is supplied, then the value at the key must match it.
 * Returns true if the key was deleted. Returns false if the key does not match redisValue or
 * the key was already deleted.
 */
export async function unlock(
  redisClient: RedisClient,
  redisKey: string,
  redisValue?: string,
): Promise<boolean> {
  // If a value is provided, check that redis holds the expected value.
  if (redisValue) {
    const clientGet = promisify<string, string | null>(redisClient.get).bind(redisClient);
    const currentValue = await clientGet(redisKey);
    if (currentValue !== redisValue) {
      return false;
    }
  }

  // Release the value.
  const clientDel = promisify<string, number>(redisClient.del).bind(redisClient);
  const lockResult: number = await clientDel(redisKey);

  // Returns true if a value was deleted.
  return !!lockResult;
}

/**
 * Atomic increment. If the key does not exist, it is treated as zero and incremented to one.
 *
 * Returns the new value.
 */
export async function increment(
  redisClient: RedisClient,
  redisKey: string,
  expiryMs: number | null = null,
): Promise<number> {
  const clientIncr = promisify(redisClient.incr).bind(redisClient);
  const clientPexpire = promisify(redisClient.pexpire).bind(redisClient);
  const newValue = await clientIncr(redisKey);
  if (expiryMs !== null) {
    await clientPexpire(redisKey, expiryMs);
  }
  return newValue;
}

/**
 * Atomic deccrement. If the key does not exist, it is treated as zero and decremented to -1.
 *
 * Returns the new value.
 */
export async function decrement(
  redisClient: RedisClient,
  redisKey: string,
  expiryMs: number | null = null,
): Promise<number> {
  const clientDecr = promisify(redisClient.decr).bind(redisClient);
  const clientPexpire = promisify(redisClient.pexpire).bind(redisClient);
  const newValue = await clientDecr(redisKey);
  if (expiryMs !== null) {
    await clientPexpire(redisKey, expiryMs);
  }
  return newValue;
}

export async function setAsync(
  {
    key,
    value,
  }: {
    key: string,
    value: string,
  },
  client: RedisClient,
): Promise<void> {
  const setAsyncFunc = promisify(client.set).bind(client);
  await setAsyncFunc(key, value);
}

export async function setexAsync(
  {
    key,
    value,
    timeToLiveSeconds,
  }: {
    key: string,
    value: string,
    timeToLiveSeconds: number,
  },
  client: RedisClient,
): Promise<string> {
  const setAsyncFunc = promisify(client.setex).bind(client);
  return setAsyncFunc(key, timeToLiveSeconds, value);
}

export async function getAsync(key: string, client: RedisClient): Promise<string | null> {
  const getAsyncFunc = promisify(client.get).bind(client);
  return getAsyncFunc(key);
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function deleteAsync(key: string, client: any): Promise<number> {
  const deleteAsyncFunc: (key: string) => Promise<number> = promisify(client.del).bind(client);
  return deleteAsyncFunc(key);
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function deleteAllAsync(client: any): Promise<void> {
  const deleteAllAsyncFunc: () => Promise<void> = promisify(client.flushall).bind(client);
  return deleteAllAsyncFunc();
}

export async function mgetAsync(keys: string[], client: RedisClient) {
  const mgetAsyncFunc = promisify<string[], string[]>(client.mget).bind(client);
  return mgetAsyncFunc(keys);
}

export async function hSetAsync(
  {
    hash,
    pairs,
  }: {
    hash: string,
    pairs: { [key: string]: string },
  },
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  client: any,
): Promise<number> {
  if (_.isEmpty(pairs)) {
    return 0;
  }

  const hSetAsyncFunc: (args: string[]) => Promise<number> = promisify(client.hset).bind(client);
  return hSetAsyncFunc([hash, ..._.flatten(Object.entries(pairs))]);
}

export async function hGetAsync(
  {
    hash,
    key,
  }: {
    hash: string,
    key: string,
  },
  client: RedisClient,
): Promise<string | null> {
  const hGetAsyncFunc = promisify(client.hget).bind(client);
  return hGetAsyncFunc(hash, key);
}

export async function rPushAsync(
  {
    key,
    value,
  }: {
    key: string,
    value: string,
  },
  client: RedisClient,
): Promise<number> {
  const rPushAsyncFunc: (key: string, arg1: string) => Promise<number> = promisify(
    client.rpush,
  ).bind(client);
  return rPushAsyncFunc(key, value);
}

export async function lRangeAsync(
  key: string,
  client: RedisClient,
): Promise<string[]> {
  const lRangeAsyncFunc = promisify(client.lrange).bind(client);
  return lRangeAsyncFunc(key, 0, -1);
}

export async function hMGetAsync(
  {
    hash,
    fields,
  }: {
    hash: string,
    fields: string[],
  },
  client: RedisClient,
): Promise<(string | null)[]> {
  const hMGetAsyncFunc: (args: string[]) => Promise<(string | null)[]> = promisify(
    client.hmget as unknown as (...args: string[]) => Promise<(string | null)[]>,
  ).bind(client);
  return hMGetAsyncFunc([hash, ...fields]);
}

export async function hSetnxAsync(
  {
    hash,
    key,
    value,
  }: {
    hash: string,
    key: string,
    value: string,
  },
  client: RedisClient,
): Promise<number> {
  const hSetnxAsyncFunction:
  (hash: string, key: string, value: string) => Promise<number> = promisify(
    client.hsetnx,
  ).bind(client);
  return hSetnxAsyncFunction(hash, key, value);
}

export async function hGetAllAsync(
  hash: string,
  client: RedisClient,
): Promise<{[field: string]: string}> {
  const hGetAllAsyncFunc = promisify(client.hgetall).bind(client);
  const result = await hGetAllAsyncFunc(hash);
  return result || {};
}

export async function hDelAsync(
  {
    hash,
    keys,
  }: {
    hash: string,
    keys: string[],
  },
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  client: any,
): Promise<number> {
  if (keys.length === 0) {
    return 0;
  }

  const hDelAsyncFunc: (args: string[]) => Promise<number> = promisify(client.hdel).bind(client);
  return hDelAsyncFunc([hash, ...keys]);
}

export async function hincrbyAsync(
  {
    hash,
    key,
    changeBy,
  }: {
    hash: string,
    key: string,
    changeBy: string,
  },
  client: RedisClient,
): Promise<number> {
  const hincrbyAsyncFunc = promisify(
    client.hincrby as unknown as (
      arg0: string,
      arg1: string,
      arg2: string,
    ) => Promise<number>,
  ).bind(client);
  const result = await hincrbyAsyncFunc(hash, key, changeBy);
  return result as number;
}

export async function hincrbyFloatAsync(
  {
    hash,
    key,
    changeBy,
  }: {
    hash: string,
    key: string,
    changeBy: string,
  },
  client: RedisClient,
): Promise<string> {
  const hincrbyFloatAsyncFunc = promisify(
    client.hincrbyfloat as unknown as (
      arg0: string,
      arg1: string,
      arg2: string,
    ) => Promise<number>,
  ).bind(client);
  const result = await hincrbyFloatAsyncFunc(hash, key, changeBy);
  return result as string;
}

export async function incrbyAsync(
  {
    key,
    changeBy,
  }: {
    key: string,
    changeBy: string,
  },
  client: RedisClient,
): Promise<number> {
  const incrbyAsyncFunc = promisify(
    client.incrby as unknown as (
      arg0: string,
      arg1: string,
    ) => Promise<number>,
  ).bind(client);
  const result = await incrbyAsyncFunc(key, changeBy);
  return result as number;
}

export async function incrbyFloatAsync(
  {
    key,
    changeBy,
  }: {
    key: string,
    changeBy: string,
  },
  client: RedisClient,
): Promise<string> {
  const incrbyFloatAsyncFunc = promisify(
    client.incrbyfloat as unknown as (
      arg0: string,
      arg1: string,
    ) => Promise<number>,
  ).bind(client);
  const result = await incrbyFloatAsyncFunc(key, changeBy);
  return result as string;
}

export async function zRevRangeAsync(
  {
    key,
    start,
    end,
  }: {
    key: string,
    start: number,
    end: number,
  },
  client: RedisClient,
): Promise<string[]> {
  const zRevRangeAsyncFunc = promisify(client.zrevrange).bind(client);
  return zRevRangeAsyncFunc(key, start, end) as Promise<string[]>;
}

export async function zRevRangeWithScoreAsync(
  {
    key,
    start,
    end,
  }: {
    key: string,
    start: number,
    end: number,
  },
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  client: any,
): Promise<string[]> {
  const zRevRangeAsyncFunc = promisify(client.zrevrange).bind(client);
  return zRevRangeAsyncFunc(key, start, end, 'withscores') as Promise<string[]>;
}

export async function zRemAsync(
  {
    hash,
    key,
  }: {
    hash: string,
    key: string,
  },
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  client: any,
): Promise<number> {
  const zRemAsyncFunc = promisify(client.zrem).bind(client);
  return zRemAsyncFunc(hash, key) as Promise<number>;
}

/**
 * Normalizes input values for zrangebyscore and zremrangebyscore.
 * Converts infinite values to Redis's "inf" string.
 * If isInclusive is false, adds a '(' character to the immediate
 * left of a value to indicate exclusivity.
 * @param value
 * @param isInclusive
 */
function normalizeInput(value: number, isInclusive: boolean): string {
  // Redis uses "inf" to denote Infinity. If the `value` is infinite, `isInclusive` is ignored.
  if (value === Infinity) {
    return '+inf';
  } else if (value === -Infinity) {
    return '-inf';
  }
  // The default is for each end of the range to be inclusive. A '(' character to the immediate
  // left of a value indicates exclusivity, i.e. the value should not be included.
  return `${isInclusive ? '' : '('}${value.toString()}`;
}

/** https://redis.io/commands/zrangebyscore/
 * @param {Object} args
 * * **key** `string` - The key to the ZSET.
 * * **start** `number` - The minimum score of the lookup.
 * * **startIsInclusive** `boolean[=true]` - Whether the score defined by `start` should be
 * included in the score range. If `start` is `+/-Infinity`, this option is ignored.
 * * **end** `number` - The maximum score of the lookup.
 * * **endIsInclusive** `boolean[=true]` - Whether the score defined by `end`
 * should be included in the score range. If `end` is `+/-Infinity`, this option is
 * ignored.
 * * **withScores** `boolean[=true]` - Whether to include the scores in the return value.
 * If this value is `true`, the return value will be altered, but the type will remain the same.
 * See @returns for more information.
 * @param {RedisClient} client - The Redis client.
 * @returns {Promise<string[]>} A list of the values whose scores were contained within the range
 * provided by `start` & `end`. If `withScores=true`, the return value will alternate
 * `[<value1>, <value1score>, <value2>, <value2score>, ...]`.
 */
export async function zRangeByScoreAsync(
  {
    key,
    start,
    startIsInclusive = true,
    end,
    endIsInclusive = true,
    withScores = true,
  }: {
    key: string,
    start: number,
    startIsInclusive?: boolean,
    end: number,
    endIsInclusive?: boolean,
    withScores?: boolean,
  },
  client: RedisClient,
): Promise<string[]> {

  const inputs: [string, string, string] = [
    key,
    normalizeInput(start, startIsInclusive),
    normalizeInput(end, endIsInclusive),
  ];
  if (withScores) {
    const zRangeByScoreAsyncFunc = promisify<
      string,
      string | number,
      string | number,
      string,
      string[]
    >(client.zrangebyscore).bind(client);
    return zRangeByScoreAsyncFunc(...inputs, 'withscores');
  } else {
    const zRangeByScoreAsyncFunc = promisify<
      string,
      string | number,
      string | number,
      string[]
    >(client.zrangebyscore).bind(client);
    return zRangeByScoreAsyncFunc(...inputs);
  }
}

/**
 * Removes all elements in a sorted set stored at key with scores between the specified range.
 * @param {Object} args
 * * **key** `string` - The key to the sorted set.
 * * **start** `number` - The minimum score of the range.
 * * **startIsInclusive** `boolean[=true]` - Whether the score defined by `start` should
 * be included in the range.
 *    If `start` is `+/-Infinity`, this option is ignored.
 * * **end** `number` - The maximum score of the range.
 * * **endIsInclusive** `boolean[=true]` - Whether the score defined by `end` should be
 * included in the range.
 *    If `end` is `+/-Infinity`, this option is ignored.
 * @param {RedisClient} client - The Redis client.
 * @returns {Promise<number>} The number of elements removed from the sorted set.
 */
export async function zRemRangeByScoreAsync(
  {
    key,
    start,
    startIsInclusive = true,
    end,
    endIsInclusive = true,
  }: {
    key: string,
    start: number,
    startIsInclusive?: boolean,
    end: number,
    endIsInclusive?: boolean,
  },
  client: RedisClient,
): Promise<number> {
  const inputs: [string, string, string] = [
    key,
    normalizeInput(start, startIsInclusive),
    normalizeInput(end, endIsInclusive),
  ];

  const zRemRangeByScoreAsyncFunc = promisify<
    string,
    string | number,
    string | number,
    number
    >(client.zremrangebyscore).bind(client);

  return zRemRangeByScoreAsyncFunc(...inputs);
}

/** https://redis.io/commands/zscore/
 * @param {Object} args
 * * **hash** `string` - The key to the ZSET.
 * * **key** `string` - The member in the ZSET to lookup.
 * @param {RedisClient} client - The Redis client.
 * @returns {Promise<string|null>} The score of the member in the ZSET. `null` if not found.
 */
export async function zScoreAsync(
  {
    hash,
    key,
  }: {
    hash: string,
    key: string,
  },
  client: RedisClient,
): Promise<string|null> {
  const zScoreAsyncFunc = promisify<string, string, string|null>(client.zscore).bind(client);
  return zScoreAsyncFunc(hash, key);
}

export async function zAddAsync(
  {
    key,
    value,
    score,
  }: {
    key: string,
    value: string,
    score: number,
  },
  client: RedisClient,
): Promise<number> {
  const zAddAsyncFunc = promisify<string, number, string, number>(client.zadd).bind(client);
  return zAddAsyncFunc(key, score, value);
}
