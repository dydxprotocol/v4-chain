/**
 * Environment variables required by Ender.
 */

import {
  parseSchema,
  baseConfigSchema,
  parseBoolean,
  parseInteger,
} from '@dydxprotocol-indexer/base';
import {
  kafkaConfigSchema,
} from '@dydxprotocol-indexer/kafka';
import {
  postgresConfigSchema,
} from '@dydxprotocol-indexer/postgres';
import { redisConfigSchema } from '@dydxprotocol-indexer/redis';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...redisConfigSchema,
  ...kafkaConfigSchema,
  SEND_WEBSOCKET_MESSAGES: parseBoolean({
    default: true,
  }),
  ORDERBOOK_MID_PRICE_REFRESH_INTERVAL_MS: parseInteger({ default: 10_000 }), // 10 seconds
};

export default parseSchema(configSchema);
