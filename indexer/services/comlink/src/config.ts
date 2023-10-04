import {
  baseConfigSchema,
  parseBoolean,
  parseInteger,
  parseSchema,
  parseString,
} from '@dydxprotocol-indexer/base';
import { complianceConfigSchema } from '@dydxprotocol-indexer/compliance';
import {
  postgresConfigSchema,
} from '@dydxprotocol-indexer/postgres';
import { redisConfigSchema } from '@dydxprotocol-indexer/redis';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...redisConfigSchema,
  ...complianceConfigSchema,

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

  // Rate limit for screening new / refreshed addresses
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_POINTS: parseInteger({ default: 2 }),
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_DURATION_SECONDS: parseInteger({ default: 60 }), // 2 reqs / min
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_POINTS: parseInteger({ default: 100 }),
  // 100 req / min
  RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_DURATION_SECONDS: parseInteger({ default: 60 }),
  // Threshold for refreshing compliance data for an address when screened
  MAX_AGE_SCREENED_ADDRESS_COMPLIANCE_DATA_SECONDS: parseInteger({ default: 86_400 }), //  1 day
};

////////////////////////////////////////////////////////////////////////////////
//                             CONFIG PROCESSING                              //
////////////////////////////////////////////////////////////////////////////////

// Process the top-level configuration.
const config = parseSchema(configSchema);

export default config;
