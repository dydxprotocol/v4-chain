import { logger, wrapBackgroundTask } from '@dydxprotocol-indexer/base';
import {
  updateOnMessageFunction,
  consumer,
  producer,
  startConsumer,
  stopConsumer,
  TO_ENDER_TOPIC,
} from '@dydxprotocol-indexer/kafka';
import { IndexerTendermintBlock, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import _ from 'lodash';
import yargs from 'yargs';

import config from './config';
import { annotateIndexerTendermintEvent } from './helpers/block-helpers';
import { AnnotatedIndexerTendermintBlock, AnnotatedIndexerTendermintEvent } from './helpers/types';

/**
 * Creates an IndexerTendermintBlock from a KafkaMessage.
 * Throws an error if there's an issue.
 */
function getIndexerTendermintBlock(
  message: KafkaMessage,
): IndexerTendermintBlock | undefined {
  if (!message || !message.value || !message.timestamp) {
    throw Error('Empty message');
  }
  const messageValueBinary: Uint8Array = new Uint8Array(message.value);

  const block: IndexerTendermintBlock = IndexerTendermintBlock.decode(
    messageValueBinary,
  );

  return block;
}

export function seek(offset: bigint): void {
  logger.info({
    at: 'consumer#seek',
    message: 'Seeking...',
    offset: offset.toString(),
  });

  consumer!.seek({
    topic: TO_ENDER_TOPIC,
    partition: 0,
    offset: offset.toString(),
  });

  logger.info({
    at: 'consumer#seek',
    message: 'Seeked.',
    offset: offset.toString(),
  });
}

export async function connect(height: number): Promise<void> {
  await Promise.all([
    consumer!.connect(),
    producer.connect(),
  ]);

  await consumer!.subscribe({
    topic: TO_ENDER_TOPIC,
    fromBeginning: true,
  });

  updateOnMessageFunction((_topic: string, message: KafkaMessage): Promise<void> => {
    return printMessageAtHeight(message, height);
  });

  logger.info({
    at: 'consumers#connect',
    message: 'Connected to Kafka',
  });
}

export async function printMessageAtHeight(
  currentMessage: KafkaMessage,
  targetHeight: number,
): Promise<void> {
  const indexerTendermintBlock: IndexerTendermintBlock | undefined = getIndexerTendermintBlock(
    currentMessage,
  );
  if (indexerTendermintBlock === undefined) {
    return;
  }

  const currentBlockHeight: number = parseInt(indexerTendermintBlock.height.toString(), 10);
  if (currentBlockHeight < targetHeight) {
    const offsetToSeek: number = targetHeight - currentBlockHeight + Number(currentMessage.offset);
    await seek(BigInt(offsetToSeek));
  } else if (currentBlockHeight === targetHeight) {
    const annotatedEvents: AnnotatedIndexerTendermintEvent[] = [];
    _.forEach(indexerTendermintBlock.events, (event: IndexerTendermintEvent) => {
      const annotatedEvent:
      AnnotatedIndexerTendermintEvent | undefined = annotateIndexerTendermintEvent(
        event,
      );
      if (annotatedEvent === undefined) {
        logger.error({
          at: 'printMessageAtHeight',
          message: 'Failed to parse event',
          event,
        });
        throw Error('Failed to parse event');
      }
      annotatedEvents.push(annotatedEvent);
    });
    const annotatedBlock: AnnotatedIndexerTendermintBlock = {
      ...indexerTendermintBlock,
      events: [],
      annotatedEvents,
    };
    logger.info({
      at: 'printMessageAtHeight',
      message: 'Printing block',
      block: annotatedBlock,
    });
    await stopConsumer();
  }
}

async function startKafka(height: number): Promise<void> {
  logger.info({
    at: 'index#start',
    message: `Starting in env ${config.NODE_ENV}`,
  });

  await connect(height);
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

async function start(height: number): Promise<void> {
  logger.info({
    at: 'index#start',
    message: `Connecting to kafka brokers: ${config.KAFKA_BROKER_URLS}`,
  });
  await startKafka(height);
  logger.info({
    at: 'index#start',
    message: `Successfully connected to kafka brokers: ${config.KAFKA_BROKER_URLS}`,
  });
}

const args = yargs.options({
  height: {
    type: 'number',
    alias: 'h',
    description: 'Height to print block at',
    required: true,
  },
}).argv;

wrapBackgroundTask(start(args.height), false, 'main');
