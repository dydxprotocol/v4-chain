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
<<<<<<< HEAD
    // Need to set fromBeginning to true, so when vulcan restarts, it will consume all messages
    // rather than ignoring the messages in queue that were produced before ender was started.
    fromBeginning: true,
=======
    // fromBeginning is by default set to false, so vulcan will only consume messages produced
    // after vulcan was started. This config should almost never matter, because by Vulcan should
    // read from the last read offset. fromBeginning will only matter if the offset is lost.
    // In the case where the offset is lost, Vulcan should read from head because in 60 seconds all
    // short term messages will expire and we can resend stateful orders through bazooka to Vulcan.
    fromBeginning: config.PROCESS_FROM_BEGINNING,
>>>>>>> 5644e389 (Set fromBeginning to false by default (#2101))
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
