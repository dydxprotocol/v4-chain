import { logger, wrapBackgroundTask } from '@dydxprotocol-indexer/base';
import {
  consumer,
  producer,
  startConsumer,
  stopConsumer,
  TO_VULCAN_TOPIC,
  updateOnMessageFunction,
} from '@dydxprotocol-indexer/kafka';
import { isStatefulOrder } from '@dydxprotocol-indexer/v4-proto-parser';
import { OffChainUpdateV1 } from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';

import config from './config';

/**
 * Creates an OffChainUpdateV1 from a KafkaMessage.
 * Throws an error if there's an issue.
 */
function getOffChainUpdate(
  message: KafkaMessage,
): OffChainUpdateV1 {
  if (!message || !message.value || !message.timestamp) {
    throw Error('Empty message');
  }
  const messageValue: Buffer = message.value!;
  const messageValueBinary: Uint8Array = new Uint8Array(messageValue);
  return OffChainUpdateV1.decode(messageValueBinary);
}

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
    return printMessageWithTimestampHeader(message);
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

export async function printMessageWithTimestampHeader(
  currentMessage: KafkaMessage,
): Promise<void> {
  const update: OffChainUpdateV1 | undefined = getOffChainUpdate(
    currentMessage,
  );
  if (update === undefined) {
    return;
  }
  logger.info({
    at: 'printMessageWithTimestampHeader',
    message: 'Printing update',
    update,
  });

  if (update.orderPlace === undefined) {
    return;
  }
  if (isStatefulOrder(update.orderPlace!.order!.orderId!.orderFlags)) {
    logger.info({
      at: 'printMessageWithTimestampHeader',
      message: 'Printing stateful order',
      update,
    });
    logger.info({
      at: 'printMessageWithTimestampHeader',
      message: 'Printing message & headers',
      headers: JSON.stringify(currentMessage.headers, replacer),
      // update,
    });
    await stopConsumer();
  }
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
