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

    const pointCost: number = getPointCost(ipAddr);

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
): number {
  if (isIndexerIp(ipAddress)) {
    return INTERNAL_REQUEST_POINTS;
  }

  return EXTERNAL_REQUEST_POINTS;
}
