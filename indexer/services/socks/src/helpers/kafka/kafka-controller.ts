import { logger } from '@dydxprotocol-indexer/base';
import {
  WebsocketTopics,
  consumer,
  stopConsumer,
} from '@dydxprotocol-indexer/kafka';

export async function connect(): Promise<void> {
  await consumer.connect();

  logger.info({
    at: 'kafka-controller#connect',
    message: 'Connected to Kafka',
  });

  await consumer.subscribe({ topics: Object.values(WebsocketTopics) });
}

export async function disconnect(): Promise<void> {
  logger.info({
    at: 'kafka-controller#stop',
    message: 'Stopping kafka consumer',
  });

  await stopConsumer();
}
