import { logger } from '@dydxprotocol-indexer/base';
import {
  consumer, initConsumer, producer, TO_ENDER_TOPIC, updateOnMessageFunction,
} from '@dydxprotocol-indexer/kafka';
import { KafkaMessage } from 'kafkajs';

import { onMessage } from '../../lib/on-message';

export async function connect(): Promise<void> {
  await Promise.all([
    initConsumer(),
    producer.connect(),
  ]);

  await consumer!.subscribe({
    topic: TO_ENDER_TOPIC,
    // https://kafka.js.org/docs/consuming#a-name-from-beginning-a-frombeginning
    // Need to set fromBeginning to true, so when ender restarts, it will consume all messages
    // rather than ignoring the messages in queue that were produced before ender was started.
    fromBeginning: true,
  });

  updateOnMessageFunction(
    (_topic: string, message: KafkaMessage): Promise<void> => onMessage(message),
  );

  logger.info({
    at: 'consumers#connect',
    message: 'Connected to Kafka',
  });
}
