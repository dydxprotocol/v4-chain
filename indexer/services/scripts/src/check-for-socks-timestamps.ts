/* eslint-disable @typescript-eslint/require-await */
import { logger, safeJsonStringify, wrapBackgroundTask } from '@dydxprotocol-indexer/base';
import {
  consumer,
  producer,
  startConsumer,
  stopConsumer,
  TO_VULCAN_TOPIC,
  updateOnMessageFunction,
} from '@dydxprotocol-indexer/kafka';
import { KafkaMessage } from 'kafkajs';

import config from './config';
import { Channel, getChannel, getMessageToForward } from './helpers/kafka-helpers';

export async function connect(): Promise<void> {
  await Promise.all([
    consumer.connect(),
    producer.connect(),
  ]);

  await consumer.subscribe({
    topic: TO_VULCAN_TOPIC,
    fromBeginning: true,
  });

  updateOnMessageFunction((_topic: string, message: KafkaMessage): Promise<void> => {
    return onMessage(_topic, message);
  });

  logger.info({
    at: 'consumers#connect',
    message: 'Connected to Kafka',
  });
}

async function onMessage(topic: string, message: KafkaMessage): Promise<void> {
  const channel: Channel | undefined = getChannel(topic);
  if (channel !== Channel.V4_ACCOUNTS) {
    return;
  }
  try {
    getMessageToForward(channel, message);
    logger.info({
      at: 'onMessage',
      message: 'Forwarded message',
      kafkaMessage: safeJsonStringify(message),
    });
  } catch (error) {}
}

async function startKafka(): Promise<void> {
  logger.info({
    at: 'index#start',
    message: `Starting in env ${config.NODE_ENV}`,
  });

  await connect();
  await startConsumer();

  logger.info({
    at: 'index#start',
    message: 'Successfully started',
  });
}

process.on('SIGTERM', async () => {
  logger.info({
    at: 'index#SIGTERM',
    message: 'Received SIGTERM, shutting down',
  });
  await stopConsumer();
});

async function start(): Promise<void> {
  logger.info({
    at: 'index#start',
    message: `Connecting to kafka brokers: ${config.KAFKA_BROKER_URLS}`,
  });
  await startKafka();
  logger.info({
    at: 'index#start',
    message: `Successfully connected to kafka brokers: ${config.KAFKA_BROKER_URLS}`,
  });
}

wrapBackgroundTask(start(), false, 'main');
