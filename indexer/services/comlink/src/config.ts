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
    default: 1000,
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
  // Expose setting compliance status, only set to true in dev/staging.
  EXPOSE_SET_COMPLIANCE_ENDPOINT: parseBoolean({ default: false }),

  // Affiliates config
  VOLUME_ELIGIBILITY_THRESHOLD: parseInteger({ default: 10_000 }),

  // Vaults config
  VAULT_PNL_HISTORY_DAYS: parseInteger({ default: 90 }),
  VAULT_PNL_HISTORY_HOURS: parseInteger({ default: 72 }),
  VAULT_PNL_START_DATE: parseString({ default: '2024-01-01T00:00:00Z' }),
  VAULT_LATEST_PNL_TICK_WINDOW_HOURS: parseInteger({ default: 1 }),
  VAULT_FETCH_FUNDING_INDEX_BLOCK_WINDOWS: parseInteger({ default: 250_000 }),
  VAULT_CACHE_TTL_MS: parseInteger({ default: 120_000 }), // 2 minutes
  // Alchemy webhook config
  ALCHEMY_AUTH_TOKEN: parseString({ default: '' }),
  // Cache-Control directives
  CACHE_CONTROL_DIRECTIVE_ADDRESSES: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_AFFILIATES: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_ASSET_POSITIONS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_CANDLES: parseString({ default: 'public, max-age=1' }),
  // omit compliance
  CACHE_CONTROL_DIRECTIVE_FILLS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_FUNDING: parseString({ default: 'public, max-age=10' }),
  // omit height
  CACHE_CONTROL_DIRECTIVE_HISTORICAL_BLOCK_TRADING_REWARDS: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_HISTORICAL_FUNDING: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_HISTORICAL_PNL: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_HISTORICAL_TRADING_REWARDS: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_ORDERBOOK: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_ORDERS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_PERPETUAL_MARKETS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_PERPETUAL_POSITIONS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_SOCIAL_TRADING: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_SPARKLINES: parseString({ default: 'public, max-age=10' }),
  CACHE_CONTROL_DIRECTIVE_TIME: parseString({ default: 'no-cache, no-store, no-transform' }),
  CACHE_CONTROL_DIRECTIVE_TRADES: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_TRANSFERS: parseString({ default: 'public, max-age=1' }),
  CACHE_CONTROL_DIRECTIVE_VAULTS: parseString({ default: 'public, max-age=10' }),
};

////////////////////////////////////////////////////////////////////////////////
//                             CONFIG PROCESSING                              //
////////////////////////////////////////////////////////////////////////////////

// Process the top-level configuration.
const config = parseSchema(configSchema);

export default config;
