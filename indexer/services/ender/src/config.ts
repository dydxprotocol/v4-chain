/**
 * Environment variables required by Ender.
 */

import {
  parseSchema,
  baseConfigSchema,
  parseBoolean,
} from '@klyraprotocol-indexer/base';
import {
  kafkaConfigSchema,
} from '@klyraprotocol-indexer/kafka';
import {
  postgresConfigSchema,
} from '@klyraprotocol-indexer/postgres';
import { redisConfigSchema } from '@klyraprotocol-indexer/redis';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...redisConfigSchema,
  ...kafkaConfigSchema,
  SEND_WEBSOCKET_MESSAGES: parseBoolean({
    default: true,
  }),
  SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_STATEFUL_ORDERS: parseBoolean({
    default: false,
  }),
};

export default parseSchema(configSchema);
