import { logger, stats } from '@dydxprotocol-indexer/base';
import { redis as redisLib } from '@dydxprotocol-indexer/redis';
import { RateLimiterMemory, RateLimiterRedis, RateLimiterRes } from 'rate-limiter-flexible';
import { RedisClient, RetryStrategyOptions } from 'redis';

import config from '../config';

export const ratelimitRedis: {
  client: RedisClient,
  connect: () => Promise<void>,
} = redisLib.createRedisClient(
  config.RATE_LIMIT_REDIS_URL, config.REDIS_RECONNECT_TIMEOUT_MS,
);
const rateLimitClient = ratelimitRedis.client as RedisClient & {
  enable_offline_queue: boolean,
  connect_timeout: number,
  ready: boolean,
  stream: { destroy: () => void },
  options: { retry_strategy: (options: RetryStrategyOptions) => number },
};
// node-redis 2.8 reads these properties on each command or retry, so setting them here leaves the
// shared redis package unchanged. Without the offline queue a limiter call fails at once while
// Redis is unreachable, instead of holding the request for the whole outage and replaying it on
// recovery. With nothing queued it is safe to keep reconnecting for good: the default gives up
// after 1 hour of retries, leaving every limiter failing until comlink restarts.
rateLimitClient.enable_offline_queue = false;
rateLimitClient.connect_timeout = Infinity;

// The shared retry strategy logs an error on every attempt (every 500 ms), which with no
// connect_timeout lasts the whole outage; run it once per interval and just retry in between.
const RECONNECT_LOG_INTERVAL_MS: number = 10_000;
const sharedRetryStrategy = rateLimitClient.options.retry_strategy;
let lastRetryLogAt: number = -Infinity;
rateLimitClient.options.retry_strategy = (options: RetryStrategyOptions): number => {
  const now: number = Date.now();
  if (now - lastRetryLogAt < RECONNECT_LOG_INTERVAL_MS) {
    stats.distribution('redis.connection_retries', options.attempt, {
      url: config.RATE_LIMIT_REDIS_URL,
    });
    return config.REDIS_RECONNECT_TIMEOUT_MS;
  }
  lastRetryLogAt = now;
  return sharedRetryStrategy(options);
};
// So the first failed attempt of the next outage is logged.
rateLimitClient.on('ready', () => {
  lastRetryLogAt = -Infinity;
});

// node-redis 2.8 never leaves two states by itself: a server demoted to replica keeps answering
// READONLY on a live connection, and a failed ready check (e.g. NOAUTH after a credential change)
// leaves the client connected but never ready. Dropping the socket sends both through the normal
// reconnect, which also resolves the host again.
function reconnectRateLimitStore(): void {
  rateLimitClient.stream.destroy();
}

// Limiter calls drop the socket on a timeout (the connection may be half-open, and dropping it
// fails the calls queued on it instead of charging them late) or on READONLY. At most once per
// floor, since a reconnect to a healthy replica becomes ready and answers READONLY again. A
// timeout while not ready skips the floor: node-redis queues even its fast-fail errors behind the
// pending ready check, so that check got no answer and every call would wait out the timeout.
const STORE_FAILURE_RECONNECT_FLOOR_MS: number = 5_000;
let lastStoreFailureReconnectAt: number = -Infinity;

function reconnectAfterStoreFailure(readyCheckHangs: boolean): void {
  const now: number = Date.now();
  if (!readyCheckHangs && now - lastStoreFailureReconnectAt < STORE_FAILURE_RECONNECT_FLOOR_MS) {
    return;
  }
  lastStoreFailureReconnectAt = now;
  logger.info({ at: 'rate-limiters#storeFailure', message: 'Reconnecting rate-limit Redis' });
  reconnectRateLimitStore();
}

export function resetRateLimitStoreThrottles(): void {
  lastRetryLogAt = -Infinity;
  lastStoreFailureReconnectAt = -Infinity;
}
rateLimitClient.on('error', (error: Error) => {
  if (error.message.startsWith('Ready check failed')) {
    // Delayed, so a lasting cause is retried (and logged) every few seconds, not every 500 ms.
    setTimeout(() => {
      if (!rateLimitClient.ready) {
        logger.info({ at: 'rate-limiters#readyCheck', message: 'Reconnecting rate-limit Redis' });
        reconnectRateLimitStore();
      }
    }, config.RATE_LIMIT_STORE_READY_RETRY_MS).unref();
  }
});

// rate-limiter-flexible's rejectIfRedisNotReady cannot see node-redis 2.8's state, so callers
// that must refuse while Redis is unreachable check this instead.
export function isRateLimitStoreReady(): boolean {
  return rateLimitClient.ready;
}

class StoreTimeoutError extends Error {}

function withStoreTimeout<T>(call: Promise<T>): Promise<T> {
  let timer: NodeJS.Timeout | undefined;
  const timeout: Promise<never> = new Promise((_resolve, reject) => {
    // After the event loop itself was blocked (GC, a long synchronous handler), Node runs expired
    // timers before it reads the socket. setImmediate runs after that read, so a reply already
    // received wins and a slow comlink is not taken for a failed store (which drops the socket).
    timer = setTimeout(
      () => setImmediate(() => reject(new StoreTimeoutError(
        `Rate-limit store did not answer within ${config.RATE_LIMIT_STORE_TIMEOUT_MS} ms`,
      ))),
      config.RATE_LIMIT_STORE_TIMEOUT_MS,
    );
  });
  return Promise.race([call, timeout]).finally(() => clearTimeout(timer));
}

// A stalled or half-open connection stays `ready`, so without a timeout a limiter call waits until
// TCP gives up and then runs late. Timeouts and store errors reject with an Error, as
// rate-limiter-flexible does for store errors, so each caller applies its own policy; a limiter
// built with a fallback answers from memory instead, with the same limits, per comlink task.
class StoreRateLimiter extends RateLimiterRedis {
  private readonly fallback?: RateLimiterMemory;

  constructor(opts: ConstructorParameters<typeof RateLimiterRedis>[0], withFallback = false) {
    super(opts);
    if (withFallback) {
      this.fallback = new RateLimiterMemory({
        points: opts.points,
        duration: opts.duration,
        keyPrefix: opts.keyPrefix,
      });
    }
  }

  async consume(
    key: string | number,
    pointsToConsume?: number,
    options?: { [key: string]: unknown },
  ): Promise<RateLimiterRes> {
    try {
      return await withStoreTimeout(super.consume(key, pointsToConsume, options));
    } catch (error) {
      if (!(error instanceof Error)) {
        throw error;
      }
      // A write inside the limiter's Lua script fails with code ERR, so match the message.
      if (error instanceof StoreTimeoutError || error.message.includes('READONLY')) {
        reconnectAfterStoreFailure(error instanceof StoreTimeoutError && !rateLimitClient.ready);
      }
      if (this.fallback === undefined) {
        throw error;
      }
      stats.increment(`${config.SERVICE_NAME}.rate_limit_fallback_used`, 1, {
        key_prefix: this.keyPrefix,
      });
      return this.fallback.consume(key, pointsToConsume, options);
    }
  }
}

export function getDefaultRateLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_GET_POINTS,
    duration: config.RATE_LIMIT_GET_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/get`,
  });
}

// Generic rate limiter for all GET requests, limits per IP
export const defaultRateLimiter: RateLimiterRedis = getDefaultRateLimiter();

// Rate-limiter for Turnkey /signin and /uploadAddress together, per IP (IPv6 per /64)
export function getTurnkeyIpRateLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_TURNKEY_IP_POINTS,
    duration: config.RATE_LIMIT_TURNKEY_IP_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/turnkey_ip`,
  });
}
export const turnkeyIpRateLimiter: RateLimiterRedis = getTurnkeyIpRateLimiter();

// New Turnkey sub-orgs per IPv4 address or IPv6 /48, which bounds how much of the global cap one
// network can take
export function getTurnkeyNewSuborgIpLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_TURNKEY_NEW_SUBORG_IP_POINTS,
    duration: config.RATE_LIMIT_TURNKEY_NEW_SUBORG_IP_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/turnkey_new_suborg_ip`,
  });
}
export const turnkeyNewSuborgIpLimiter: RateLimiterRedis = getTurnkeyNewSuborgIpLimiter();

// Cap on new Turnkey sub-orgs across all callers
export function getTurnkeyNewSuborgGlobalLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_TURNKEY_NEW_SUBORG_GLOBAL_POINTS,
    duration: config.RATE_LIMIT_TURNKEY_NEW_SUBORG_GLOBAL_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/turnkey_new_suborg_global`,
  });
}
export const turnkeyNewSuborgGlobalLimiter: RateLimiterRedis = getTurnkeyNewSuborgGlobalLimiter();

// The screen limiters guard a paid provider, so on store errors and timeouts they fall back to
// memory: screening keeps working and stays capped.
function getScreenProviderStoreLimiter(
  points: number,
  duration: number,
  keyPrefix: string,
): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points,
    duration,
    keyPrefix,
  }, true);
}

// Rate-limiter for /screen endpoint querying a compliance provider, limits per IP
export function getScreenProviderLimiter(): RateLimiterRedis {
  return getScreenProviderStoreLimiter(
    config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_POINTS,
    config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_DURATION_SECONDS,
    `${config.SERVICE_NAME}/screen_providers`,
  );
}
export const screenProviderLimiter: RateLimiterRedis = getScreenProviderLimiter();

// Rate-limiter for /screen endpoint querying a compliance provider, limits the total calls made
// across all IPs
export function getScreenProviderGlobalLimiter(): RateLimiterRedis {
  return getScreenProviderStoreLimiter(
    config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_POINTS,
    config.RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_DURATION_SECONDS,
    `${config.SERVICE_NAME}/screen_providers_global`,
  );
}
export const screenProviderGlobalLimiter: RateLimiterRedis = getScreenProviderGlobalLimiter();

// Rate-limiter for /orders endpoint, limits per IP
export function getOrdersRateLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_ORDERS_POINTS,
    duration: config.RATE_LIMIT_ORDERS_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/orders`,
  });
}
export const ordersRateLimiter: RateLimiterRedis = getOrdersRateLimiter();

// Rate-limiter for /fills endpoint, limits per IP
export function getFillsRateLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_FILLS_POINTS,
    duration: config.RATE_LIMIT_FILLS_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/fills`,
  });
}
export const fillsRateLimiter: RateLimiterRedis = getFillsRateLimiter();

// Rate-limiter for /candles endpoint, limits per IP
export function getCandlesRateLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_CANDLES_POINTS,
    duration: config.RATE_LIMIT_CANDLES_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/candles`,
  });
}
export const candlesRateLimiter: RateLimiterRedis = getCandlesRateLimiter();

// Rate-limiter for /sparklines endpoint, limits per IP
export function getSparklinesRateLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_SPARKLINES_POINTS,
    duration: config.RATE_LIMIT_SPARKLINES_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/sparklines`,
  });
}
export const sparklinesRateLimiter: RateLimiterRedis = getSparklinesRateLimiter();

// Rate-limiter for /historical-pnl endpoint, limits per IP
export function getHistoricalPnlRateLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_HISTORICAL_PNL_POINTS,
    duration: config.RATE_LIMIT_HISTORICAL_PNL_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/historical_pnl`,
  });
}
export const historicalPnlRateLimiter: RateLimiterRedis = getHistoricalPnlRateLimiter();

// Rate-limiter for /pnl endpoint, limits per IP
export function getPnlRateLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_PNL_POINTS,
    duration: config.RATE_LIMIT_PNL_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/pnl`,
  });
}
export const pnlRateLimiter: RateLimiterRedis = getPnlRateLimiter();

// Rate-limiter for /funding endpoint, limits per IP
export function getFundingRateLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_FUNDING_POINTS,
    duration: config.RATE_LIMIT_FUNDING_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/funding`,
  });
}
export const fundingRateLimiter: RateLimiterRedis = getFundingRateLimiter();

// Rate-limiter for /gasSponsorship endpoint, limits per IP
export function getGasSponsorshipRateLimiter(): RateLimiterRedis {
  return new StoreRateLimiter({
    storeClient: ratelimitRedis.client,
    points: config.RATE_LIMIT_GAS_SPONSORSHIP_POINTS,
    duration: config.RATE_LIMIT_GAS_SPONSORSHIP_DURATION_SECONDS,
    keyPrefix: `${config.SERVICE_NAME}/gasSponsorship`,
  });
}
export const gasSponsorshipRateLimiter: RateLimiterRedis = getGasSponsorshipRateLimiter();
