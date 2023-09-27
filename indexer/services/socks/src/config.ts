import {
  baseConfigSchema,
  parseBoolean,
  parseInteger,
  parseNumber,
  parseSchema,
  parseString,
} from '@dydxprotocol-indexer/base';
import {
  complianceConfigSchema,
} from '@dydxprotocol-indexer/compliance';
import { kafkaConfigSchema } from '@dydxprotocol-indexer/kafka';
import {
  postgresConfigSchema,
} from '@dydxprotocol-indexer/postgres';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...kafkaConfigSchema,
  ...complianceConfigSchema,

  PORT: parseString({
    default: '8000',
  }),
  WS_PORT: parseInteger({
    default: 8080,
  }),
  BATCH_SEND_INTERVAL_MS: parseInteger({
    default: 250,
  }),

  WS_HEARTBEAT_INTERVAL_MS: parseInteger(),
  WS_HEARTBEAT_TIMEOUT_MS: parseInteger(),

  RATE_LIMIT_ENABLED: parseBoolean({ default: false }),
  RATE_LIMIT_SUBSCRIBE_POINTS: parseNumber({ default: 2 }),
  RATE_LIMIT_SUBSCRIBE_DURATION_MS: parseInteger({ default: 1000 }),
  RATE_LIMIT_PING_POINTS: parseNumber({ default: 5 }),
  RATE_LIMIT_PING_DURATION_MS: parseInteger({ default: 1000 }),
  RATE_LIMIT_INVALID_MESSAGE_POINTS: parseNumber({ default: 2 }),
  RATE_LIMIT_INVALID_MESSAGE_DURATION_MS: parseInteger({ default: 1000 }),

  MESSAGE_FORWARDER_STATSD_SAMPLE_RATE: parseNumber({ default: 1.0 }),
  ENABLE_ORDERBOOK_LOGS: parseBoolean({ default: true }),
  PERPETUAL_MARKETS_REFRESHER_INTERVAL_MS: parseInteger({ default: 300_000 }), // 5 minutes

  // Config for getting initial response data
  COMLINK_URL: parseString(),
  AXIOS_TIMEOUT_MS: parseInteger({ default: 5000 }), // 5 seconds
  INITIAL_GET_TIMEOUT_MS: parseInteger({ default: 20_000 }), // 20 seconds
};

////////////////////////////////////////////////////////////////////////////////
//                             CONFIG PROCESSING                              //
////////////////////////////////////////////////////////////////////////////////

// Process the top-level configuration.
const config = parseSchema(configSchema);

export default config;
