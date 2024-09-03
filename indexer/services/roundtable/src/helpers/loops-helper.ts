import {
  logger,
  stats,
  wrapBackgroundTask,
} from '@dydxprotocol-indexer/base';
import { redis } from '@dydxprotocol-indexer/redis';
import { v4 as uuidv4 } from 'uuid';

import config from '../config';
import {
  REDIS_VALUE,
  STATS_NO_SAMPLING,
} from '../lib/constants';
import { generateRandomStartDelayMs } from './helpers';
import { redisClient } from './redis';

let numRunningTasks = 0;

export const delay = (ms: number) => new Promise((r) => setTimeout(r, ms));

export function startLoop(
  loopTask: () => Promise<unknown>,
  taskName: string,
  loopIntervalMs: number,
  extendedLoopLockMultiplier?: number,
): void {
  wrapBackgroundTask(
    startLoopAsync(
      loopTask,
      taskName,
      loopIntervalMs,
      extendedLoopLockMultiplier,
    ),
    true,
    'taskName',
  );
}

export function loopMetricNames(taskName: string): {
  startedStat: string,
  completedStat: string,
  durationRatioStat: string,
  timingStat: string,
  exceededMaxConcurrentTasksStat: string,
  couldNotAcquireExtendedLockStat: string,
  redisKey: string,
  extendedLockKey: string,
} {
  return {
    startedStat: `${config.SERVICE_NAME}.loops.${taskName}.started`,
    completedStat: `${config.SERVICE_NAME}.loops.${taskName}.completed`,
    durationRatioStat: `${config.SERVICE_NAME}.loops.duration_ratio`,
    timingStat: `${config.SERVICE_NAME}.loops.${taskName}.timing`,
    exceededMaxConcurrentTasksStat: `${config.SERVICE_NAME}.loops.exceeded_max_concurrent_tasks`,
    couldNotAcquireExtendedLockStat: `${config.SERVICE_NAME}.loops.could_not_acquire_extended_lock`,
    redisKey: `${config.SERVICE_NAME}.loops.${taskName}.task_timeouts`,
    extendedLockKey: `${config.SERVICE_NAME}.loops.${taskName}.task_extended_timeouts`,
  };
}

async function startLoopAsync(
  loopTask: () => Promise<unknown>,
  taskName: string,
  loopIntervalMs: number,
  extendedLoopLockMultiplier?: number,
): Promise<void> {

  const {
    startedStat,
    completedStat,
    durationRatioStat,
    timingStat,
    exceededMaxConcurrentTasksStat,
    couldNotAcquireExtendedLockStat,
    redisKey,
    extendedLockKey,
  } = loopMetricNames(taskName);

  if (config.START_DELAY_ENABLED) {
    await delay(generateRandomStartDelayMs(loopIntervalMs));
  }
  for (;;) {
    if (numRunningTasks > config.MAX_CONCURRENT_RUNNING_TASKS) {
      stats.increment(
        exceededMaxConcurrentTasksStat,
        { taskName },
      );
      await delay(config.EXCEEDED_MAX_CONCURRENT_RUNNING_TASKS_DELAY_MS);
      continue;
    }

    // Create exclusive redis lock with timeout.
    const lockResult: boolean = await redis.lockWithExpiry(
      redisClient,
      redisKey,
      REDIS_VALUE,
      loopIntervalMs,
    );

    // If lock was not created, wait and try again
    if (!lockResult) {
      // Wait until next task.
      await waitUntilKeyExpiresPlusJitter(redisKey);
      continue;
    }

    // If lock was created, run the task.
    let extendedLockResult: boolean = true;

    // Generate a random string to save as the value for the extended lock key in Redis
    // This way, only this Roundtable instance can unlock the extended lock in the finally clause
    // Otherwise, other instances could unlock the lock while this instance is still working
    const extendedLockRedisValue: string = `${REDIS_VALUE} ${uuidv4()}`;
    if (extendedLoopLockMultiplier !== undefined) {
      extendedLockResult = await redis.lockWithExpiry(
        redisClient,
        extendedLockKey,
        extendedLockRedisValue,
        loopIntervalMs * extendedLoopLockMultiplier,
      );
    }

    // If lock was not created, try again
    if (!extendedLockResult) {
      stats.increment(
        couldNotAcquireExtendedLockStat,
        { taskName },
      );
      logger.error({
        at: 'loop-helpers#startLoopAsync/extended-lock',
        message: 'could not acquire extended lock to run task',
        taskName,
      });
      // Unlock the regular lock and wait until next task.
      await redis.unlock(
        redisClient,
        redisKey,
        REDIS_VALUE,
      );
      await waitUntilKeyExpiresPlusJitter(extendedLockKey);
      continue;
    }

    // Log start of task.
    const start: number = Date.now();
    stats.gauge(startedStat, 1);
    numRunningTasks += 1;

    try {
      await loopTask();
      stats.gauge(completedStat, 1);
    } catch (error) {
      stats.gauge(completedStat, 0);
      logger.error({
        at: `loops-helpers/${taskName}`,
        message: 'uncaught error in an individual loop',
        error,
        taskName,
        disableGroupingHash: true,
      });
    } finally {
      numRunningTasks -= 1;
      if (extendedLoopLockMultiplier !== undefined && extendedLockResult) {
        // Only unlock the extended lock key if the value matches
        // Otherwise, the extended lock was set by another Roundtable instance
        await redis.unlock(
          redisClient,
          extendedLockKey,
          extendedLockRedisValue,
        );
      }

      // Log timing of task.
      const end: number = Date.now();
      const loopDuration: number = end - start;
      stats.timing(timingStat, loopDuration);
      stats.gauge(
        durationRatioStat,
        loopDuration / loopIntervalMs,
        STATS_NO_SAMPLING,
        { taskName },
      );
    }

    // Wait until next task.
    await waitUntilKeyExpiresPlusJitter(redisKey);
  }
}

async function waitUntilKeyExpiresPlusJitter(redisKey: string): Promise<void> {
  const pttl: number = await redis.pttl(redisClient, redisKey);
  const jitter: number = Math.ceil(Math.random() * (pttl * config.JITTER_FRACTION_OF_DELAY) + 1);
  if (pttl > 0) {
    await delay(pttl + jitter);
  }
}
