/**
 * Environment variables required by Ender.
 */

import {
  parseSchema,
  baseConfigSchema,
  parseBoolean,
  parseString,
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
  // Config var to skip processing stateful order events with specific uuids.
  // Order UUIDs should be in a string delimited by commas.
  // Only set if invalid order events are being included in a block and preventing ender from
  // progressing.
  SKIP_STATEFUL_ORDER_UUIDS: parseString({
    default: '',
  }),
  ORDERBOOK_MID_PRICE_REFRESH_INTERVAL_MS: parseInteger({ default: 10_000 }), // 10 seconds
};

export default parseSchema(configSchema);
