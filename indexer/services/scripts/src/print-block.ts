import { logger, wrapBackgroundTask } from '@dydxprotocol-indexer/base';
import {
  addOnMessageFunction,
  consumer,
  producer,
  startConsumer,
  stopConsumer,
  TO_ENDER_TOPIC,
} from '@dydxprotocol-indexer/kafka';
import { IndexerTendermintBlock } from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import yargs from 'yargs';

import config from './config';
import { runAsyncScript } from './helpers/util';

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
  logger.info({
    at: 'onMessage#getIndexerTendermintBlock',
    message: 'Received message',
    offset: message.offset,
  });

  const block: IndexerTendermintBlock = IndexerTendermintBlock.decode(
    messageValueBinary,
  );
  logger.info({
    at: 'onMessage#getIndexerTendermintBlock',
    message: 'Parsed message',
    offset: message.offset,
    height: block.height,
    block,
  });

  return block;
}

export function seek(offset: bigint): void {
  logger.info({
    at: 'consumer#seek',
    message: 'Seeking...',
    offset: offset.toString(),
  });

  consumer.seek({
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
    consumer.connect(),
    producer.connect(),
  ]);

  await consumer.subscribe({
    topic: TO_ENDER_TOPIC,
    // Need to set fromBeginning to true, so when ender restarts, it will consume all messages
    // rather than ignoring the messages in queue that were produced before ender was started.
    fromBeginning: true,
  });

  addOnMessageFunction((_topic: string, message: KafkaMessage): Promise<void> => {
    return printMessageAtHeight(message, height);
  });

  logger.info({
    at: 'consumer#connect',
    message: 'Added onMessage function',
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
  logger.info({
    at: 'consumer#printMessageAtHeight',
    message: 'Received message',
  });
  const indexerTendermintBlock: IndexerTendermintBlock | undefined = getIndexerTendermintBlock(
    currentMessage,
  );
  if (indexerTendermintBlock === undefined) {
    return;
  }

  const currentBlockHeight: number = parseInt(indexerTendermintBlock.height.toString(), 10);
  console.log(`Current block height: ${currentBlockHeight}`);
  if (currentBlockHeight < targetHeight) {
    const offsetToSeek = BigInt(targetHeight - currentBlockHeight) + currentMessage.offset;
    console.log(`Seeking to offset: ${offsetToSeek}`);
    const desiredMessage = await seek(BigInt(offsetToSeek));
    console.log(JSON.stringify(desiredMessage));
  } else if (currentBlockHeight === targetHeight) {
    console.log(JSON.stringify(currentMessage));
  } else {
    throw Error(`Current block height ${currentBlockHeight} is greater than target height ${targetHeight}`);
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
    description: 'Height to print message for',
    required: true,
  },
}).argv;

wrapBackgroundTask(start(args.height), true, 'main');
