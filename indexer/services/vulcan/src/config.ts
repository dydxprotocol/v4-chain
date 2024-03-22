/**
 * Environment variables required by Vulcan.
 */

import {
  parseNumber,
  parseSchema,
  baseConfigSchema,
  parseBoolean,
} from '@dydxprotocol-indexer/base';
import {
  kafkaConfigSchema,
} from '@dydxprotocol-indexer/kafka';
import {
  redisConfigSchema,
} from '@dydxprotocol-indexer/redis';

export const configSchema = {
  ...baseConfigSchema,
  ...kafkaConfigSchema,
  ...redisConfigSchema,

  BATCH_PROCESSING_ENABLED: parseBoolean({ default: true }),
  KAFKA_BATCH_PROCESSING_COMMIT_FREQUENCY_MS: parseNumber({
    default: 3_000,
  }),
  FLUSH_KAFKA_MESSAGES_INTERVAL_MS: parseNumber({
    default: 10,
  }),
  MAX_WEBSOCKET_MESSAGES_TO_QUEUE_PER_TOPIC: parseNumber({
    default: 20,
  }),
  // Set this flag to false during fast sync.
  SEND_WEBSOCKET_MESSAGES: parseBoolean({
    default: true,
  }),
  DONT_SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_STATEFUL_ORDERS: parseBoolean({
    default: false,
  }),
};

export default parseSchema(configSchema);
