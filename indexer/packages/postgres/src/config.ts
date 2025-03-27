/**
 * Environment variables required by postgres module.
 */

import {
  parseBoolean,
  parseInteger,
  parseString,
  parseSchema,
  baseConfigSchema,
} from '@dydxprotocol-indexer/base';

export const configSecrets: (keyof typeof postgresConfigSchema)[] = [
  'DB_PASSWORD',
];

export const DEFAULT_KMS_KEY = 'default kms key';

export const postgresConfigSchema = {
  ...baseConfigSchema,

  // Required environment variables.
  DB_HOSTNAME: parseString({ default: 'localhost' }),
  DB_READONLY_HOSTNAME: parseString({ default: 'localhost' }),
  IS_USING_DB_READONLY: parseBoolean({ default: true }),

  DB_PORT: parseInteger({ default: 5435 }),
  DB_NAME: parseString({ default: 'dydx_dev' }),
  DB_USERNAME: parseString({ default: 'dydx_dev' }),
  DB_PASSWORD: parseString({ default: 'dydxserver123' }),
  PG_POOL_MIN: parseInteger({ default: 1 }),
  PG_POOL_MAX: parseInteger({ default: 2 }),
  PG_ACQUIRE_CONNECTION_TIMEOUT_MS: parseInteger({ default: 10_000 }),
  PERPETUAL_MARKETS_REFRESHER_INTERVAL_MS: parseInteger({ default: 30_000 }), // 30 seconds
  ASSET_REFRESHER_INTERVAL_MS: parseInteger({ default: 30_000 }), // 30 seconds
  MARKET_REFRESHER_INTERVAL_MS: parseInteger({ default: 30_000 }), // 30 seconds
  LIQUIDITY_TIER_REFRESHER_INTERVAL_MS: parseInteger({ default: 30_000 }), // 30 seconds
  BLOCK_HEIGHT_REFRESHER_INTERVAL_MS: parseInteger({ default: 1_000 }), // 1 second
  USE_READ_REPLICA: parseBoolean({ default: false }),

  // Optional environment variables.
  NODE_ENV: parseString({ default: null }),
};

export default parseSchema(postgresConfigSchema);
