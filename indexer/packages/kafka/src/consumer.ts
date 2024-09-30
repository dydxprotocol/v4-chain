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
  consumer = kafka.consumer({
    groupId,
    sessionTimeout: config.KAFKA_SESSION_TIMEOUT_MS,
    rebalanceTimeout: config.KAFKA_REBALANCE_TIMEOUT_MS,
    heartbeatInterval: config.KAFKA_HEARTBEAT_INTERVAL_MS,
    maxWaitTimeInMs: config.KAFKA_WAIT_MAX_TIME_MS,
    readUncommitted: false,
    maxBytes: 4194304, // 4MB
    rackId: await getAvailabilityZoneId(),
  });

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
    consumerRunConfig.eachBatch = onBatchFunction;
  } else {
    consumerRunConfig.eachMessage = async ({ topic, message }) => {
      await onMessageFunction(topic, message);
    };
  }

  await consumer!.run(consumerRunConfig);

  logger.info({
    at: 'consumers#connect',
    message: 'Started kafka consumer',
    groupId,
  });
}
