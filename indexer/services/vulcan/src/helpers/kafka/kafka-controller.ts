import { logger } from '@dydxprotocol-indexer/base';
import {
  consumer, producer, KafkaTopics, updateOnMessageFunction, updateOnBatchFunction,
} from '@dydxprotocol-indexer/kafka';
import { KafkaMessage } from 'kafkajs';

import config from '../../config';
import { onBatch } from '../../lib/on-batch';
import { onMessage } from '../../lib/on-message';

export async function connect(): Promise<void> {
  await Promise.all([
    consumer.connect(),
    producer.connect(),
  ]);

  await consumer.subscribe({
    topic: KafkaTopics.TO_VULCAN,
    // https://kafka.js.org/docs/consuming#a-name-from-beginning-a-frombeginning
    // Need to set fromBeginning to true, so when vulcan restarts, it will consume all messages
    // rather than ignoring the messages in queue that were produced before ender was started.
    fromBeginning: true,
  });

  if (config.BATCH_PROCESSING_ENABLED) {
    logger.info({
      at: 'consumers#connect',
      message: 'Batch processing enabled',
    });
    updateOnBatchFunction(onBatch);
  } else {
    logger.info({
      at: 'consumers#connect',
      message: 'Batch processing disabled. Processing each message individually',
    });
    updateOnMessageFunction((_topic: string, message: KafkaMessage): Promise<void> => {
      return onMessage(message);
    });
  }

  logger.info({
    at: 'consumers#connect',
    message: 'Connected to Kafka',
  });
}
