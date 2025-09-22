/**
 * Environment variables required by Roundtable.
 */

import {
  parseSchema,
  baseConfigSchema,
  parseBoolean,
  parseInteger,
  parseNumber,
  parseString,
  ONE_MINUTE_IN_MILLISECONDS,
  THIRTY_SECONDS_IN_MILLISECONDS,
  FIVE_MINUTES_IN_MILLISECONDS,
  ONE_HOUR_IN_MILLISECONDS,
  ONE_SECOND_IN_MILLISECONDS,
  TEN_SECONDS_IN_MILLISECONDS,
  FOUR_HOURS_IN_MILLISECONDS,
  ONE_DAY_IN_MILLISECONDS,
} from '@dydxprotocol-indexer/base';
import {
  kafkaConfigSchema,
} from '@dydxprotocol-indexer/kafka';
import {
  postgresConfigSchema,
} from '@dydxprotocol-indexer/postgres';
import {
  redisConfigSchema,
} from '@dydxprotocol-indexer/redis';

export const configSchema = {
  ...baseConfigSchema,
  ...postgresConfigSchema,
  ...kafkaConfigSchema,
  ...redisConfigSchema,

  // Loop Enablement
  LOOPS_ENABLED_MARKET_UPDATER: parseBoolean({ default: true }),
  LOOPS_ENABLED_DELETE_ZERO_PRICE_LEVELS: parseBoolean({ default: true }),
  LOOPS_ENABLED_UNCROSS_ORDERBOOK: parseBoolean({ default: true }),
  LOOPS_ENABLED_PNL_TICKS: parseBoolean({ default: true }),
  LOOPS_ENABLED_REMOVE_EXPIRED_ORDERS: parseBoolean({ default: true }),
  LOOPS_ORDERBOOK_INSTRUMENTATION: parseBoolean({ default: true }),
  LOOPS_PNL_INSTRUMENTATION: parseBoolean({ default: true }),
  LOOPS_CANCEL_STALE_ORDERS: parseBoolean({ default: true }),
  LOOPS_ENABLED_UPDATE_RESEARCH_ENVIRONMENT: parseBoolean({ default: false }),
  LOOPS_ENABLED_TAKE_FAST_SYNC_SNAPSHOTS: parseBoolean({ default: true }),
  LOOPS_ENABLED_DELETE_OLD_FAST_SYNC_SNAPSHOTS: parseBoolean({ default: true }),
  LOOPS_ENABLED_TRACK_LAG: parseBoolean({ default: false }),
  LOOPS_ENABLED_REMOVE_OLD_ORDER_UPDATES: parseBoolean({ default: true }),
  LOOPS_ENABLED_AGGREGATE_TRADING_REWARDS_DAILY: parseBoolean({ default: true }),
  LOOPS_ENABLED_AGGREGATE_TRADING_REWARDS_WEEKLY: parseBoolean({ default: true }),
  LOOPS_ENABLED_AGGREGATE_TRADING_REWARDS_MONTHLY: parseBoolean({ default: true }),
  LOOPS_ENABLED_SUBACCOUNT_USERNAME_GENERATOR: parseBoolean({ default: true }),
  LOOPS_ENABLED_LEADERBOARD_PNL_ALL_TIME: parseBoolean({ default: false }),
  LOOPS_ENABLED_LEADERBOARD_PNL_DAILY: parseBoolean({ default: false }),
  LOOPS_ENABLED_LEADERBOARD_PNL_WEEKLY: parseBoolean({ default: false }),
  LOOPS_ENABLED_LEADERBOARD_PNL_MONTHLY: parseBoolean({ default: false }),
  LOOPS_ENABLED_LEADERBOARD_PNL_YEARLY: parseBoolean({ default: false }),
  LOOPS_ENABLED_UPDATE_WALLET_TOTAL_VOLUME: parseBoolean({ default: true }),
  LOOPS_ENABLED_UPDATE_AFFILIATE_INFO: parseBoolean({ default: true }),
  LOOPS_ENABLED_DELETE_OLD_FIREBASE_NOTIFICATION_TOKENS: parseBoolean({ default: true }),
  LOOPS_ENABLED_REFRESH_VAULT_PNL: parseBoolean({ default: true }),
  LOOPS_ENABLED_CACHE_ORDERBOOK_MID_PRICES: parseBoolean({ default: true }),
  LOOPS_ENABLED_UPDATE_FUNDING_PAYMENTS: parseBoolean({ default: true }),
  LOOPS_ENABLED_UPDATE_PNL: parseBoolean({ default: true }),

  // Loop Timing
  LOOPS_INTERVAL_MS_MARKET_UPDATER: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_DELETE_ZERO_PRICE_LEVELS: parseInteger({
    default: 2 * ONE_MINUTE_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_UNCROSS_ORDERBOOK: parseInteger({
    default: 15 * ONE_SECOND_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_PNL_TICKS: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_REMOVE_EXPIRED_ORDERS: parseInteger({
    default: 2 * ONE_MINUTE_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_ORDERBOOK_INSTRUMENTATION: parseInteger({
    default: 1 * ONE_MINUTE_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_PNL_INSTRUMENTATION: parseInteger({
    default: ONE_HOUR_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_CANCEL_STALE_ORDERS: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_UPDATE_RESEARCH_ENVIRONMENT: parseInteger({
    default: ONE_HOUR_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_TAKE_FAST_SYNC_SNAPSHOTS: parseInteger({
    default: FOUR_HOURS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_DELETE_OLD_FAST_SYNC_SNAPSHOTS: parseInteger({
    default: ONE_DAY_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_UPDATE_COMPLIANCE_DATA: parseInteger({
    default: FIVE_MINUTES_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_TRACK_LAG: parseInteger({
    default: TEN_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_REMOVE_OLD_ORDER_UPDATES: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_AGGREGATE_TRADING_REWARDS: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_PERFORM_COMPLIANCE_STATUS_TRANSITIONS: parseInteger({
    default: ONE_HOUR_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_SUBACCOUNT_USERNAME_GENERATOR: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_LEADERBOARD_PNL_ALL_TIME: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_LEADERBOARD_PNL_DAILY: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_LEADERBOARD_PNL_WEEKLY: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_LEADERBOARD_PNL_MONTHLY: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_LEADERBOARD_PNL_YEARLY: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_UPDATE_WALLET_TOTAL_VOLUME: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_UPDATE_AFFILIATE_INFO: parseInteger({
    default: THIRTY_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_DELETE_FIREBASE_NOTIFICATION_TOKENS_MONTHLY: parseInteger({
    default: 30 * ONE_DAY_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_CACHE_ORDERBOOK_MID_PRICES: parseInteger({
    default: TEN_SECONDS_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_REFRESH_VAULT_PNL: parseInteger({
    default: 5 * ONE_MINUTE_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_UPDATE_FUNDING_PAYMENTS: parseInteger({
    default: ONE_MINUTE_IN_MILLISECONDS,
  }),
  LOOPS_INTERVAL_MS_UPDATE_PNL: parseInteger({
    default: ONE_MINUTE_IN_MILLISECONDS,
  }),
  // Start delay
  START_DELAY_ENABLED: parseBoolean({ default: true }),
  MAX_START_DELAY_MS: parseInteger({ default: 3 * ONE_MINUTE_IN_MILLISECONDS }),
  MAX_START_DELAY_FRACTION_OF_INTERVAL: parseNumber({ default: 0.1 }),
  JITTER_FRACTION_OF_DELAY: parseNumber({ default: 0.01 }),

  // Lock multipliers
  MARKET_UPDATER_LOCK_MULTIPLIER: parseInteger({ default: 10 }),
  DELETE_ZERO_PRICE_LEVELS_LOCK_MULTIPLIER: parseInteger({ default: 1 }),
  UNCROSS_ORDERBOOK_LOCK_MULTIPLIER: parseInteger({ default: 1 }),
  PNL_TICK_UPDATE_LOCK_MULTIPLIER: parseInteger({ default: 20 }),
  SUBACCOUNT_USERNAME_GENERATOR_LOCK_MULTIPLIER: parseInteger({ default: 5 }),
  // for initial backfill from genesis, can configure this to be a higher value
  UPDATE_FUNDING_PAYMENTS_LOCK_MULTIPLIER: parseInteger({ default: 3 }),
  // for initial backfill from genesis, can configure this to be a higher value
  UPDATE_PNL_LOCK_MULTIPLIER: parseInteger({ default: 3 }),

  // Maximum number of running tasks - set this equal to PG_POOL_MIN in .env, default is 2
  MAX_CONCURRENT_RUNNING_TASKS: parseInteger({ default: 2 }),
  EXCEEDED_MAX_CONCURRENT_RUNNING_TASKS_DELAY_MS: parseInteger({ default: 1000 }),

  // PNL ticks
  PNL_TICK_UPDATE_INTERVAL_MS: parseInteger({ default: ONE_HOUR_IN_MILLISECONDS }),
  PNL_TICK_MAX_ROWS_PER_UPSERT: parseInteger({ default: 1000 }),
  PNL_TICK_MAX_ACCOUNTS_PER_RUN: parseInteger({ default: 65000 }),
  LEADERBOARD_PNL_MAX_ROWS_PER_UPSERT: parseInteger({ default: 1000 }),

  // Remove expired orders
  BLOCKS_TO_DELAY_EXPIRY_BEFORE_SENDING_REMOVES: parseInteger({ default: 20 }),

  // Cancel stale orders
  CANCEL_STALE_ORDERS_QUERY_BATCH_SIZE: parseInteger({ default: 10000 }),

  // Tracking indexer lag
  TRACK_LAG_INDEXER_FULL_NODE_URL: parseString({ default: '' }), // i.e. http://11.11.11.11:26657
  TRACK_LAG_VALIDATOR_URL: parseString({ default: '' }), // i.e. http://11.11.11.11:26657
  TRACK_LAG_OTHER_FULL_NODE_URL: parseString({ default: '' }), // i.e. http://11.11.11.11:26657

  // Update research environment
  AWS_ACCOUNT_ID: parseString(),
  AWS_REGION: parseString(),
  S3_BUCKET_ARN: parseString(),
  FAST_SYNC_SNAPSHOT_IDENTIFIER_PREFIX: parseString({ default: 'fast-sync' }),
  ECS_TASK_ROLE_ARN: parseString(),
  KMS_KEY_ARN: parseString(),
  RDS_INSTANCE_NAME: parseString(),
  ATHENA_CATALOG_NAME: parseString({ default: 'AwsDataCatalog' }),
  ATHENA_DATABASE_NAME: parseString({ default: 'default' }),
  ATHENA_WORKING_GROUP: parseString({ default: 'primary' }),
  SKIP_TO_ATHENA_TABLE_WRITING: parseBoolean({ default: false }),

  // Update compliance data
  ACTIVE_ADDRESS_THRESHOLD_SECONDS: parseInteger({ default: 86_400 }),
  MAX_COMPLIANCE_DATA_AGE_SECONDS: parseInteger({ default: 2_630_000 }), // 1 month
  MAX_ACTIVE_COMPLIANCE_DATA_AGE_SECONDS: parseInteger({ default: 86_400 }), // 1 day
  MAX_COMPLIANCE_DATA_QUERY_PER_LOOP: parseInteger({ default: 100 }),
  // v2/wallet/synchronous rate limit is 15/s https://developers.elliptic.co/docs/configuration
  COMPLIANCE_PROVIDER_QUERY_BATCH_SIZE: parseInteger({ default: 15 }),
  COMPLIANCE_PROVIDER_QUERY_DELAY_MS: parseInteger({ default: ONE_SECOND_IN_MILLISECONDS }),
  CLOSE_ONLY_TO_BLOCKED_DAYS: parseInteger({ default: 7 }),

  // Remove old cached order updates
  OLD_CACHED_ORDER_UPDATES_WINDOW_MS: parseInteger({ default: 30 * ONE_SECOND_IN_MILLISECONDS }),

  // Aggregate Trading Rewards
  AGGREGATE_TRADING_REWARDS_MAX_INTERVAL_SIZE_MS: parseInteger({
    default: ONE_HOUR_IN_MILLISECONDS,
  }),
  AGGREGATE_TRADING_REWARDS_CHUNK_SIZE: parseInteger({ default: 50 }),

  // Uncross orderbook
  STALE_ORDERBOOK_LEVEL_THRESHOLD_SECONDS: parseInteger({ default: 10 }),

  // Subaccount username generator
  SUBACCOUNT_USERNAME_SUFFIX_RANDOM_DIGITS: parseInteger({ default: 3 }),
  SUBACCOUNT_USERNAME_BATCH_SIZE: parseInteger({ default: 2000 }),
  // number of attempts to generate username for a subaccount
  ATTEMPT_PER_SUBACCOUNT: parseInteger({ default: 3 }),

  // Refresh vault pnl view
  TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS: parseInteger({ default: 15 * ONE_MINUTE_IN_MILLISECONDS }),
};

export default parseSchema(configSchema);
