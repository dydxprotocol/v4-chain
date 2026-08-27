import { logger } from '@dydxprotocol-indexer/base';
import express from 'express';
import _ from 'lodash';
import { RateLimiterRedis, RateLimiterRes } from 'rate-limiter-flexible';

import config from '../config';
import { create4xxResponse } from './helpers';
import { getIpAddr, isIndexerIp } from './utils';

const INTERNAL_REQUEST_POINTS: number = 0;
const EXTERNAL_REQUEST_POINTS: number = 1;

// Note, the return-type of this method and the return type of the middle-ware function have are
// not defined in the Express types package, and so are omitted for readability as the inferred
// types are too verbose.
export function rateLimiterMiddleware(
  rateLimiter: RateLimiterRedis,
  postfixKey?: string,
) {
  return async (
    req: express.Request,
    res: express.Response,
    next: express.NextFunction,
  ) => {
    if (!config.RATE_LIMIT_ENABLED) {
      return next();
    }

    const ipAddr: string | undefined = getIpAddr(req);

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
    const postfix: string | undefined = postfixKey ? _.get(req, postfixKey) : undefined;
    const redisKey: string = postfix ? ipAddr.concat(postfix) : ipAddr;

    try {
      const limitRes: RateLimiterRes = await rateLimiter.consume(redisKey, pointCost);
      res.set({
        'RateLimit-Remaining': limitRes.remainingPoints,
        'RateLimit-Reset': Date.now() + limitRes.msBeforeNext,
        'RateLimit-Limit': rateLimiter.points,
      });
    } catch (reject) {
      if (reject instanceof Error) {
        logger.error({
          at: 'rate-limit',
          message: 'redis error when checking rate limit',
          reject,
        });
      } else {
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
