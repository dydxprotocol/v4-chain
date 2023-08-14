import {
  logger,
} from '@dydxprotocol-indexer/base';
import { Consumer, ConsumerRunConfig, KafkaMessage } from 'kafkajs';
import { v4 as uuidv4 } from 'uuid';

import config from './config';
import { kafka } from './kafka';

const groupIdPrefix: string = config.SERVICE_NAME;
const groupIdSuffix: string = config.KAFKA_ENABLE_UNIQUE_CONSUMER_GROUP_IDS ? `_${uuidv4()}` : '';
const groupId: string = `${groupIdPrefix}${groupIdSuffix}`;

export const consumer: Consumer = kafka.consumer({
  groupId,
  sessionTimeout: config.KAFKA_SESSION_TIMEOUT_MS,
  rebalanceTimeout: config.KAFKA_REBALANCE_TIMEOUT_MS,
  heartbeatInterval: config.KAFKA_HEARTBEAT_INTERVAL_MS,
  readUncommitted: false,
  maxBytes: 4194304, // 4MB
});

// List of functions to run per message consumed.
const onMessageFunctions: ((topic: string, message: KafkaMessage) => Promise<void>)[] = [];

export function addOnMessageFunction(
  onMessage: (topic: string, message: KafkaMessage) => Promise<void>,
): void {
  onMessageFunctions.push(onMessage);
}

// Whether the consumer is stopped.
let stopped: boolean = false;

consumer.on('consumer.disconnect', async () => {
  logger.info({
    at: 'consumers#disconnect',
    message: 'Kafka consumer disconnected',
    groupId,
  });

  if (!stopped) {
    await consumer.connect();
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

export async function stopConsumer(): Promise<void> {
  logger.info({
    at: 'kafka-consumer#stop',
    message: 'Stopping kafka consumer',
    groupId,
  });

  stopped = true;
  await consumer.disconnect();
}

export async function startConsumer(): Promise<void> {
  const consumerRunConfig: ConsumerRunConfig = {
    partitionsConsumedConcurrently: config.KAFKA_CONCURRENT_PARTITIONS,
    autoCommit: true,
  };

  consumerRunConfig.eachMessage = async ({ topic, message }) => {
    await Promise.all(
      onMessageFunctions.map(
        async (onMessage: (topic: string, message: KafkaMessage) => Promise<void>) => {
          await onMessage(topic, message);
        }),
    );
  };

  await consumer.run(consumerRunConfig);

  logger.info({
    at: 'consumers#connect',
    message: 'Started kafka consumer',
    groupId,
  });
}
