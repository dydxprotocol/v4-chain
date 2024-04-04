/* eslint-disable @typescript-eslint/require-await */
import { logger, safeJsonStringify, wrapBackgroundTask } from '@dydxprotocol-indexer/base';
import {
  consumer,
  producer,
  startConsumer,
  stopConsumer,
  updateOnMessageFunction,
} from '@dydxprotocol-indexer/kafka';
import { KafkaMessage } from 'kafkajs';

import config from './config';
import {
  Channel, getChannel, getMessageToForward, WebsocketTopics,
} from './helpers/kafka-helpers';

export async function connect(): Promise<void> {
  await Promise.all([
    consumer.connect(),
    producer.connect(),
  ]);

  await consumer.subscribe({
    topic: WebsocketTopics.TO_WEBSOCKETS_SUBACCOUNTS,
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

function replacer(key: string, value: any): any {
  if (value instanceof Buffer) {
    return value.toString('base64'); // Convert Buffer to base64 string
  }
  return value;
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
    logger.info({
      at: 'printMessageWithTimestampHeader',
      message: 'Printing message & headers',
      headers: JSON.stringify(message.headers, replacer),
      // update,
    });
    logger.info({
      at: 'printTimestamp',
      message: 'Printing timestamp',
      timestamp: Number(message.timestamp),
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
