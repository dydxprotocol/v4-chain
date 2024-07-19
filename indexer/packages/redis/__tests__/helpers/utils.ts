import { DateTime } from 'luxon';
import { RetryStrategyOptions } from 'redis';

import config from '../../src/config';
import { createRedisClient } from '../../src/helpers/redis';

const res = createRedisClient(config.REDIS_URL, config.REDIS_RECONNECT_TIMEOUT_MS);

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const redis: any = res?.client;
export const connect = res?.connect;

// eslint-disable-next-line
export function expectEqual(expectedObject: any, retrievedObject: any) {
  Object.keys(expectedObject).forEach(
    (key) => {

      if (!['createdAt', 'updatedAt'].includes(key)) {
        if (key === 'userData') {
          expect(expectedObject[key]).toEqual(JSON.stringify(retrievedObject[key]));
        } if (key === 'unfillableAt') {
          expect(
            DateTime.fromISO(expectedObject[key]).toISO(),
          ).toEqual(
            DateTime.fromISO(retrievedObject[key]).toISO(),
          );
        } else {
          expect(expectedObject[key]).toEqual(retrievedObject[key]);
        }
      }
    },
  );
}

export function callRetryStrategy(
  retryStrategy: (options: RetryStrategyOptions) => number,
  error: object | null | undefined,
  attempt: number = config.REDIS_RECONNECT_ATTEMPT_ERROR_THRESHOLD,
): number {
  // RetryStrategyOptions has non-null/undefined error object, however error can be null/undefined
  // Have to type cast to not get syntax errors
  return retryStrategy({
    total_retry_time: 1,
    times_connected: 1,
    error,
    attempt,
  } as RetryStrategyOptions);
}
