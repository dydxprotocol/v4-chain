/**
 * Environment variables required for Kafka.
 */

import {
  parseInteger,
  parseString,
  parseSchema,
  NodeEnv,
  parseBoolean,
} from '@dydxprotocol-indexer/base';

export const kafkaConfigSchema = {
  // Required to set groupId and clientId for kafka.
  SERVICE_NAME: parseString(),

  KAFKA_BROKER_URLS: parseString({
    default: 'localhost:9092',
    requireInEnv: [NodeEnv.PRODUCTION, NodeEnv.STAGING],
  }),
  KAFKA_CONNECTION_TIMEOUT_MS: parseInteger({ default: 5_000 }),
  KAFKA_SESSION_TIMEOUT_MS: parseInteger({ default: 60_000 }),
  KAFKA_REBALANCE_TIMEOUT_MS: parseInteger({ default: 50_000 }),
  KAFKA_HEARTBEAT_INTERVAL_MS: parseInteger({ default: 5_000 }),
  KAFKA_CONCURRENT_PARTITIONS: parseInteger({ default: 1 }),
  // The number of messages to process before committing the offset.
  KAFKA_CONSUMER_AUTO_COMMIT_THRESHOLD: parseInteger({ default: 100 }),
  // The interval at which the consumer will commit the offset.
  // Note that the consumer will respect both the commit threshold and the commit interval
  // config, whichever comes first.
  KAFKA_CONSUMER_AUTO_COMMIT_INTERVAL_MS: parseInteger({ default: 5_000 }),
  // If true, consumers will have unique group ids, and SERVICE_NAME will be a common prefix for
  // the consumer group ids.
  KAFKA_ENABLE_UNIQUE_CONSUMER_GROUP_IDS: parseBoolean({ default: false }),
  // Set to a number smaller than the max message size for the Kafka broker
  KAFKA_MAX_BATCH_WEBSOCKET_MESSAGE_SIZE_BYTES: parseInteger({
    default: 900000, // ~900 kB, 100 kB smaller than the 1 MB default max size of messages in Kafka
  }),
  KAFKA_WAIT_MAX_TIME_MS: parseInteger({ default: 5_000 }),
};

export default parseSchema(kafkaConfigSchema);
