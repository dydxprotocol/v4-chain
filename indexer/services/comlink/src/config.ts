import {
  baseConfigSchema,
  parseBoolean,
  parseInteger,
  parseSchema,
  parseString,
} from '@dydxprotocol-indexer/base';
import {
  postgresConfigSchema,
} from '@dydxprotocol-indexer/postgres';
import { redisConfigSchema } from '@dydxprotocol-indexer/redis';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...redisConfigSchema,

  CHAIN_ID: parseString({ default: 'dydxprotocol' }),
  API_LIMIT_V4: parseInteger({
    default: 100,
  }),
  API_ORDERBOOK_LEVELS_PER_SIDE_LIMIT: parseInteger({ default: 100 }),

  // Logging config
  LOG_GETS: parseBoolean({ default: false }),

  // Express server config
  PORT: parseInteger({ default: 8080 }),
  CORS_ORIGIN: parseString({ default: '*' }),
  KEEP_ALIVE_MS: parseInteger({ default: 61_000 }),
  HEADERS_TIMEOUT_MS: parseInteger({ default: 65_000 }),

  // Rate limit Redis URL
  RATE_LIMIT_REDIS_URL: parseString({
    default: 'redis://localhost:6382',
  }),
  // Rate limits
  RATE_LIMIT_ENABLED: parseBoolean({ default: true }),
  // IP addresses internal to the Indexer have no rate-limit
  INDEXER_INTERNAL_IPS: parseString({ default: '' }),
  // Points / duration determines the maximum rate of requests given that each requests costs 1
  // point
  RATE_LIMIT_GET_POINTS: parseInteger({ default: 100 }),
  RATE_LIMIT_GET_DURATION_SECONDS: parseInteger({ default: 10 }), // 100 requests / 10 seconds

  // Geo-blocking
  RESTRICTED_COUNTRIES: parseString({
    default: '',
  }),
};

////////////////////////////////////////////////////////////////////////////////
//                             CONFIG PROCESSING                              //
////////////////////////////////////////////////////////////////////////////////

// Process the top-level configuration.
const config = parseSchema(configSchema);

export default config;
