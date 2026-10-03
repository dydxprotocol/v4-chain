import {
  getAvailabilityZoneId,
  logger,
} from '@dydxprotocol-indexer/base';
import {
  Consumer, ConsumerRunConfig, EachBatchPayload, KafkaMessage,
} from 'kafkajs';
import { v4 as uuidv4 } from 'uuid';

import config from './config';
import { kafka } from './kafka';

const groupIdPrefix: string = config.SERVICE_NAME;
const groupIdSuffix: string = config.KAFKA_ENABLE_UNIQUE_CONSUMER_GROUP_IDS ? `_${uuidv4()}` : '';
const groupId: string = `${groupIdPrefix}${groupIdSuffix}`;

// As a hack, we made this mutable since CommonJS doesn't support top level await.
// Top level await would needed to fetch the az id (used as rack id).
// eslint-disable-next-line import/no-mutable-exports
export let consumer: Consumer | undefined;

// List of functions to run per message consumed.
let onMessageFunction: (topic: string, message: KafkaMessage) => Promise<void>;

// List of function to be run per batch consumed.
let onBatchFunction: (payload: EachBatchPayload) => Promise<void>;

/**
 * Overwrite function to be run on each kafka message
 * @param onMessage
 */
export function updateOnMessageFunction(
  onMessage: (topic: string, message: KafkaMessage) => Promise<void>,
): void {
  onMessageFunction = onMessage;
}

/**
 * Overwrite function to be run on each kafka batch
 */
export function updateOnBatchFunction(
  onBatch: (payload: EachBatchPayload) => Promise<void>,
): void {
  onBatchFunction = onBatch;
}

// Whether the consumer is stopped.
let stopped: boolean = false;

// Number of consecutive consumer crashes (each fired after KafkaJS exhausts its per-batch
// retries). Reset to 0 whenever a message is processed successfully. Used to distinguish a
// transient blip from a sustained downstream outage that should recycle the task.
let consecutiveCrashes: number = 0;

/**
 * @returns the current consecutive-crash count. Exposed for testing and observability.
 */
export function getConsecutiveCrashes(): number {
  return consecutiveCrashes;
}

/**
 * Reset the consecutive-crash count. Called after any successfully processed message/batch so a
 * consumer that is making progress never accrues toward the force-exit threshold.
 */
export function resetConsecutiveCrashes(): void {
  consecutiveCrashes = 0;
}

/**
 * Record a consumer crash and decide whether KafkaJS should restart the consumer internally.
 * Returns false once the consecutive-crash count reaches KAFKA_CONSUMER_MAX_CONSECUTIVE_CRASHES,
 * which makes KafkaJS emit a fatal `consumer.crash` (restart: false) instead of looping forever.
 */
export function recordCrashAndDecideRestart(error: Error): boolean {
  consecutiveCrashes += 1;
  const shouldRestart: boolean = consecutiveCrashes < config.KAFKA_CONSUMER_MAX_CONSECUTIVE_CRASHES;
  logger.error({
    at: 'kafka-consumer#recordCrashAndDecideRestart',
    message: shouldRestart
      ? 'Kafka consumer crashed; restarting internally'
      : 'Kafka consumer exceeded max consecutive crashes; not restarting internally',
    groupId,
    consecutiveCrashes,
    maxConsecutiveCrashes: config.KAFKA_CONSUMER_MAX_CONSECUTIVE_CRASHES,
    error,
  });
  return shouldRestart;
}

/**
 * Force the process to exit so the orchestrator (ECS) restarts the task with a fresh container.
 * Used when the consumer can no longer make progress on its own; exiting non-zero is strictly
 * safer than hanging, since Kafka offsets are committed and consumption resumes where it left off.
 */
function exitForRestart(reason: string, error?: Error): void {
  logger.crit({
    at: 'kafka-consumer#exitForRestart',
    message: `Kafka consumer cannot make progress (${reason}); exiting process for orchestrator restart`,
    groupId,
    consecutiveCrashes,
    error,
  });
  // Flush is best-effort; a hung process is worse than losing the tail of a log line.
  process.exit(1);
}

/**
 * Handle a KafkaJS `consumer.crash` event. When the crash is not restartable — either because the
 * error is non-retriable or because recordCrashAndDecideRestart has hit the crash budget — the
 * consumer would otherwise sit idle forever, so exit the process for an orchestrator restart.
 */
export function handleConsumerCrash(payload: { error: Error, restart: boolean }): void {
  if (!payload.restart && config.KAFKA_EXIT_PROCESS_ON_CONSUMER_CRASH && !stopped) {
    exitForRestart('fatal consumer.crash', payload.error);
  }
}

export async function stopConsumer(): Promise<void> {
  logger.info({
    at: 'kafka-consumer#stop',
    message: 'Stopping kafka consumer',
    groupId,
  });

  stopped = true;
  await consumer!.disconnect();
}

export async function initConsumer(): Promise<void> {
  // A fresh consumer is not shutting down; clear state left over from any prior lifecycle.
  stopped = false;
  resetConsecutiveCrashes();
  consumer = kafka.consumer({
    groupId,
    sessionTimeout: config.KAFKA_SESSION_TIMEOUT_MS,
    rebalanceTimeout: config.KAFKA_REBALANCE_TIMEOUT_MS,
    heartbeatInterval: config.KAFKA_HEARTBEAT_INTERVAL_MS,
    maxWaitTimeInMs: config.KAFKA_WAIT_MAX_TIME_MS,
    readUncommitted: false,
    maxBytes: 4194304, // 4MB
    rackId: await getAvailabilityZoneId(),
    // Bound KafkaJS's internal restart loop. Without this, a retriable error (such as an
    // unreachable Postgres during a disk-full event) causes the consumer to restart itself
    // indefinitely while the process still looks healthy to ECS. Returning false once the crash
    // budget is exhausted makes KafkaJS emit a fatal `consumer.crash` (restart: false), which
    // handleConsumerCrash turns into a process exit so the task is recycled.
    retry: {
      retries: config.KAFKA_CONSUMER_MAX_RETRIES,
      restartOnFailure: (error: Error): Promise<boolean> => Promise.resolve(
        recordCrashAndDecideRestart(error),
      ),
    },
  });

  consumer!.on('consumer.crash', (event) => handleConsumerCrash(event.payload));

  consumer!.on('consumer.disconnect', async () => {
    logger.info({
      at: 'consumers#disconnect',
      message: 'Kafka consumer disconnected',
      groupId,
    });

    if (!stopped) {
      await consumer!.connect();
      logger.info({
        at: 'kafka-consumer#disconnect',
        message: 'Kafka consumer reconnected',
        groupId,
      });
    } else {
      logger.info({
        at: 'kafka-consumer#disconnect',
        message: 'Not reconnecting since task is shutting down',
        groupId,
      });
    }
  });
}

export async function startConsumer(batchProcessing: boolean = false): Promise<void> {
  const consumerRunConfig: ConsumerRunConfig = {
    // The last offset of each batch will be committed if processing does not error.
    // The commit will still happen if the number of messages in the batch < autoCommitThreshold.
    eachBatchAutoResolve: true,
    partitionsConsumedConcurrently: config.KAFKA_CONCURRENT_PARTITIONS,
    autoCommit: true,
    autoCommitThreshold: config.KAFKA_CONSUMER_AUTO_COMMIT_THRESHOLD,
    autoCommitInterval: config.KAFKA_CONSUMER_AUTO_COMMIT_INTERVAL_MS,
  };

  if (batchProcessing) {
    consumerRunConfig.eachBatch = async (payload: EachBatchPayload) => {
      await onBatchFunction(payload);
      // A successful batch means the consumer is making progress; clear the crash counter so a
      // later isolated failure gets the full crash budget rather than an inherited count.
      resetConsecutiveCrashes();
    };
  } else {
    consumerRunConfig.eachMessage = async ({ topic, message }) => {
      await onMessageFunction(topic, message);
      resetConsecutiveCrashes();
    };
  }

  await consumer!.run(consumerRunConfig);

  logger.info({
    at: 'consumers#connect',
    message: 'Started kafka consumer',
    groupId,
  });
}
