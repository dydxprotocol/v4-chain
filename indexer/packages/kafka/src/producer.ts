import { logger } from '@dydxprotocol-indexer/base';
import { Partitioners, Producer } from 'kafkajs';

import { kafka } from './kafka';

export const producer: Producer = kafka.producer({
  createPartitioner: Partitioners.DefaultPartitioner,
});

let stopped: boolean = false;

producer.on('producer.disconnect', async () => {
  logger.info({
    at: 'kafka-producer#disconnect',
    message: 'Kafka producer disconnected',
  });

  if (!stopped) {
    await producer.connect();

    logger.info({
      at: 'kafka-producer#disconnect',
      message: 'Kafka producer reconnected',
    });
  }
});

export async function disconnect(): Promise<void> {
  stopped = true;
  await producer.disconnect();
}
