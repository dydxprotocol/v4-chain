/**
 * Utility functions for promises, timers, and intervals.
 */

import util from 'util';

import logger from './logger';
import stats from './stats';

export const delay = util.promisify(setTimeout);

/**
 * Run a periodic task.
 *
 * Like built-in setInterval(), but won't start a run if the previous run has not finished.
 * Also, unlike setInterval(), there is no delay before the initial run.
 *
 * Errors thrown by the task will be caught and logged.
 */
export function setIntervalNonOverlapping<T>(
  handler: () => T | Promise<T>,
  intervalMs: number,
  intervalTimingMetricName?: string,
): void {
  wrapBackgroundTask(
    (async () => {
      for (; ;) {
        const start: number = Date.now();

        try {
          await handler();
        } catch (error) {
          logger.error({
            at: 'tasks#setIntervalNonOverlapping',
            message: 'Uncaught error in periodic task',
            intervalMs,
            intervalTimingMetricName,
            error,
          });
        }

        const elapsedMillis: number = Date.now() - start;
        if (intervalTimingMetricName) {
          stats.timing(intervalTimingMetricName, elapsedMillis);
        }
        if (elapsedMillis < intervalMs) {
          await delay(intervalMs - elapsedMillis);
        }
      }
    })(),
    true, // Exit the process on uncaught error. Don't expect this to ever occur.
    `setIntervalNonOverlapping:${intervalTimingMetricName}`,
  );
}

/**
 * Wrapper that should be used around all hanging promises--any promise which is not awaited.
 *
 * The purpose of the wrapper is to help enforce explicit handling of all promises in a standardized
 * way. The wrapper will act as a last fallback for logging uncaught errors, and requires the caller
 * to specify whether the task should re-throw, potentially killing the Node process.
 *
 * Context:
 *   https://blog.heroku.com/best-practices-nodejs-errors
 *   https://github.com/palantir/tslint/issues/4653
 */
export function wrapBackgroundTask<T>(
  promise: Promise<T>,
  shouldRethrowErrors: boolean,
  taskName: string | null = null,
): void {
  promise.catch((error) => {
    // Paranoid: don't want the catch body to throw unexpectedly.
    try {
      logger.error({
        at: 'tasks#wrapBackgroundTask',
        message: 'Background task had an uncaught error',
        error,
        shouldRethrowErrors,
        taskName,
      });
    } catch (innerError) {
      /* eslint-disable-next-line no-console */
      console.error(
        `wrapBackgroundTask catch block failed (taskName=${taskName}):`,
        innerError?.stack,
      );
    }
    if (shouldRethrowErrors) {
      throw error;
    }
  });
}
