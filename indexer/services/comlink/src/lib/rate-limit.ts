import { logger, stats } from '@dydxprotocol-indexer/base';
import express from 'express';
import _ from 'lodash';
import { RateLimiterRedis, RateLimiterRes } from 'rate-limiter-flexible';

import { isRateLimitStoreReady } from '../caches/rate-limiters';
import config from '../config';
import { create4xxResponse } from './helpers';
import { getIpAddr, isIndexerIp } from './utils';

const INTERNAL_REQUEST_POINTS: number = 0;
const EXTERNAL_REQUEST_POINTS: number = 1;

// During a Redis outage every limited request errors at once, so log one line per limiter per
// interval; the store_error metric still counts every error.
const STORE_ERROR_LOG_INTERVAL_MS: number = 10_000;
const lastStoreErrorLogAt: Map<string, number> = new Map();

function logStoreError(keyPrefix: string, reject: Error): void {
  const now: number = Date.now();
  const last: number | undefined = lastStoreErrorLogAt.get(keyPrefix);
  if (last !== undefined && now - last < STORE_ERROR_LOG_INTERVAL_MS) {
    return;
  }
  lastStoreErrorLogAt.set(keyPrefix, now);
  logger.error({
    at: 'rate-limit',
    message: 'redis error when checking rate limit',
    keyPrefix,
    reject,
  });
}

// Note, the return-type of this method and the return type of the middle-ware function have are
// not defined in the Express types package, and so are omitted for readability as the inferred
// types are too verbose.
export function rateLimiterMiddleware(
  rateLimiter: RateLimiterRedis,
  postfixKey?: string,
  // Overrides the IP-based key; returning undefined skips this limiter for the request.
  getKey?: (req: express.Request, ipAddr: string) => string | undefined,
  // Refuse with 503 while the store is not ready or errors (READONLY, OOM). Opt-in, for routes
  // that must not run unlimited during a Redis outage; every other route fails open.
  refuseWhenStoreNotReady: boolean = false,
) {
  const refuseStoreUnavailable = (res: express.Response, reason: 'not_ready' | 'error') => {
    stats.increment(`${config.SERVICE_NAME}.rate_limit_store_not_ready`, 1, {
      key_prefix: rateLimiter.keyPrefix,
      reason,
    });
    return create4xxResponse(res, 'Service temporarily unavailable', 503);
  };

  return async (
    req: express.Request,
    res: express.Response,
    next: express.NextFunction,
  ) => {
    if (!config.RATE_LIMIT_ENABLED) {
      return next();
    }

    const ipAddr: string | undefined = getIpAddr(req);

    if (refuseWhenStoreNotReady && !isRateLimitStoreReady()) {
      return refuseStoreUnavailable(res, 'not_ready');
    }

    if (ipAddr === undefined) {
      return next();
    }

    // Clamp to the limiter's own capacity: a cost beyond that exhausts the bucket either way,
    // and an unclamped cost derived from a crafted page/limit can be large enough that its
    // decimal string trips Redis's integer parser (e.g. exponential notation), which throws a
    // plain Error and - via the catch branch below - lets the request through with no rate
    // limiting applied at all.
    const pointCost: number = Math.min(getPointCost(ipAddr, req), rateLimiter.points);

    // generate redis key
    let redisKey: string;
    if (getKey) {
      const key: string | undefined = getKey(req, ipAddr);
      if (key === undefined) {
        return next();
      }
      redisKey = key;
    } else {
      const postfix: string | undefined = postfixKey ? _.get(req, postfixKey) : undefined;
      redisKey = postfix ? ipAddr.concat(postfix) : ipAddr;
    }

    try {
      const limitRes: RateLimiterRes = await rateLimiter.consume(redisKey, pointCost);
      res.set({
        'RateLimit-Remaining': limitRes.remainingPoints,
        'RateLimit-Reset': Date.now() + limitRes.msBeforeNext,
        'RateLimit-Limit': rateLimiter.points,
      });
    } catch (reject) {
      if (reject instanceof Error) {
        stats.increment(`${config.SERVICE_NAME}.rate_limit_store_error`, 1, {
          key_prefix: rateLimiter.keyPrefix,
        });
        logStoreError(rateLimiter.keyPrefix, reject);
        if (refuseWhenStoreNotReady) {
          return refuseStoreUnavailable(res, 'error');
        }
      } else {
        // Response-code stats are registered after the limiters, so this is the only count of 429s.
        stats.increment(`${config.SERVICE_NAME}.rate_limited`, 1, {
          key_prefix: rateLimiter.keyPrefix,
        });
        const rejectRes: RateLimiterRes = reject as RateLimiterRes;
        res.set({
          'RateLimit-Remaining': rejectRes.remainingPoints,
          'RateLimit-Reset': Date.now() + rejectRes.msBeforeNext,
          'Retry-After': rejectRes.msBeforeNext,
          'RateLimit-Limit': rateLimiter.points,
        });
        return create4xxResponse(res, 'Too many requests', 429);
      }
    }

    return next();
  };
}

export function getPointCost(
  ipAddress: string,
  req: express.Request,
): number {
  if (isIndexerIp(ipAddress)) {
    return INTERNAL_REQUEST_POINTS;
  }

  // Weight cost by OFFSET depth: deep pagination (large page * limit) is far more expensive to the
  // DB than a shallow page, so it should drain the rate-limit budget faster. Read directly from
  // req.query since this middleware runs before checkSchema/handleValidationErrors.
  const rawPage: unknown = req.query?.page;
  const rawLimit: unknown = req.query?.limit;
  const page: number | undefined = rawPage !== undefined ? Number(rawPage) : undefined;
  const limit: number = rawLimit !== undefined ? Number(rawLimit) : config.API_LIMIT_V4;

  if (
    page === undefined || !Number.isFinite(page) || page <= 1 ||
    !Number.isFinite(limit) || limit <= 0
  ) {
    return EXTERNAL_REQUEST_POINTS;
  }

  const offset: number = (Math.max(1, Math.floor(page)) - 1) * Math.floor(limit);
  return EXTERNAL_REQUEST_POINTS + Math.floor(offset / config.API_LIMIT_V4);
}
