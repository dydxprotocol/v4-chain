import {
  Kafka,
} from 'kafkajs';

import config from './config';

export const KAFKA_BROKERS: string[] = config.KAFKA_BROKER_URLS.split(',');

export const kafka: Kafka = new Kafka({
  clientId: config.SERVICE_NAME,
  brokers: KAFKA_BROKERS,
  connectionTimeout: config.KAFKA_CONNECTION_TIMEOUT_MS,
});
