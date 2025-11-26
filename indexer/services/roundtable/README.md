# Roundtable

Roundtable is a task scheduling service for the dYdX v4 Indexer that runs recurring background jobs (cron-like tasks) to maintain data consistency, update aggregations, and perform periodic maintenance operations. It orchestrates dozens of independent loops that process data from PostgreSQL, Redis, and Kafka to keep the indexer state up-to-date.

## Responsibilities & Scope

- Execute recurring background tasks on configurable intervals (market updates, PnL calculations, orderbook maintenance, etc.).
- Maintain data consistency across PostgreSQL and Redis caches.
- Aggregate trading rewards and compute leaderboard rankings.
- Update perpetual market statistics (24h volume, open interest, funding rates, price changes).
- Generate and maintain PnL ticks for all subaccounts.
- Calculate funding payments for perpetual positions.
- Clean up stale data (expired orders, old cached updates, zero price levels).
- Export database snapshots to S3 for research and fast-sync purposes.
- Update compliance data from external providers (Elliptic).
- Manage affiliate program data and wallet statistics.
- Refresh materialized views for vault PnL calculations.
- Generate usernames for new subaccounts.
- Publish market updates to Kafka for WebSocket distribution.

**Out of scope:**
- Does not process real-time blockchain events (handled by Ender).
- Does not serve HTTP requests or expose public APIs.
- Does not manage WebSocket connections (handled by Socks).
- Does not write orders or execute trades.

## Architecture & Dependencies

### Internal Structure

- **Task Loops** (`src/tasks/`): Independent background jobs that run on fixed intervals with distributed locking via Redis.
- **Loop Management** (`src/helpers/loops-helper.ts`): Orchestrates task execution with timing, locking, jitter, and concurrency control.
- **AWS Helpers** (`src/helpers/aws.ts`): Utilities for RDS snapshots, S3 exports, and Athena table management.
- **PnL Helpers** (`src/helpers/pnl-ticks-helper.ts`): Complex calculations for equity, total PnL, net transfers, and funding payments.
- **Compliance Clients** (`src/helpers/compliance-clients.ts`): Integration with Elliptic for address screening.
- **Configuration** (`src/config.ts`): Centralized environment variable parsing with validation and defaults.

### Internal Dependencies

- `@dydxprotocol-indexer/base`: Shared utilities, logging, stats, and error handling.
- `@dydxprotocol-indexer/postgres`: Database models, queries, and transaction management.
- `@dydxprotocol-indexer/redis`: Redis caches for orders, orderbooks, PnL ticks, and distributed locking.
- `@dydxprotocol-indexer/kafka`: Kafka producer for sending market updates to WebSocket service.
- `@dydxprotocol-indexer/compliance`: Compliance provider integrations (Elliptic).
- `@dydxprotocol-indexer/v4-protos`: Protocol buffer definitions for messages.

### External Dependencies

- **PostgreSQL**: Primary datastore for all indexer data (orders, fills, positions, markets, transfers, etc.).
- **Redis**: Distributed locking, caching, and real-time data (orderbooks, order state, PnL ticks).
- **Kafka**: Message bus for publishing market updates to WebSocket service.
- **AWS RDS**: Database snapshot management for fast-sync and research exports.
- **AWS S3**: Storage for database exports.
- **AWS Athena**: Query engine for research data analysis.
- **Elliptic API**: Compliance screening for wallet addresses.

### Processing Flow

1. Service starts and initializes all enabled task loops based on configuration.
2. Each loop runs independently on its configured interval (e.g., every 30 seconds, 1 hour, 1 day).
3. Before executing, each loop attempts to acquire a distributed lock in Redis with TTL equal to the interval.
4. If lock is acquired, the task executes (e.g., update market stats, calculate PnL, clean up data).
5. Task reads from PostgreSQL/Redis, performs calculations, and writes results back.
6. Some tasks publish messages to Kafka (e.g., market updates for WebSocket clients).
7. Task completion is logged with timing metrics, and the lock expires naturally.
8. Random jitter is added to prevent thundering herd when multiple instances restart.

## Public Interface

Roundtable does not expose HTTP endpoints or public APIs. It operates entirely through scheduled background tasks.

### Scheduled Tasks

All tasks are defined in [`src/tasks/`](src/tasks/) and enabled/disabled via environment variables.

#### Market & Orderbook Tasks

- **Market Updater** ([`market-updater.ts`](src/tasks/market-updater.ts)): Updates 24h volume, trades, open interest, price changes, and funding rates for all perpetual markets. Publishes updates to Kafka for WebSocket clients.
- **Delete Zero Price Levels** ([`delete-zero-price-levels.ts`](src/tasks/delete-zero-price-levels.ts)): Removes orderbook price levels with zero size from Redis.
- **Uncross Orderbook** ([`uncross-orderbook.ts`](src/tasks/uncross-orderbook.ts)): Detects and removes crossed orderbook levels (bid >= ask) by deleting stale entries.
- **Orderbook Instrumentation** ([`orderbook-instrumentation.ts`](src/tasks/orderbook-instrumentation.ts)): Emits metrics on best bid/ask prices and orderbook depth for monitoring.
- **Cache Orderbook Mid Prices** ([`cache-orderbook-mid-prices.ts`](src/tasks/cache-orderbook-mid-prices.ts)): Caches mid prices for all markets in Redis for fast access.

#### PnL & Position Tasks

- **Create PnL Ticks** ([`create-pnl-ticks.ts`](src/tasks/create-pnl-ticks.ts)): Calculates hourly PnL ticks for all subaccounts based on positions, funding payments, and transfers.
- **Update PnL** ([`update-pnl.ts`](src/tasks/update-pnl.ts)): Computes PnL changes between funding periods using SQL script ([`update_pnl.sql`](src/scripts/update_pnl.sql)) that accounts for funding payments, position value changes, and trade cash flows.
- **PnL Instrumentation** ([`pnl-instrumentation.ts`](src/tasks/pnl-instrumentation.ts)): Monitors for stale PnL data and emits alerts.
- **Refresh Vault PnL** ([`refresh-vault-pnl.ts`](src/tasks/refresh-vault-pnl.ts)): Refreshes materialized views for vault PnL calculations (hourly and daily).

#### Order Management Tasks

- **Remove Expired Orders** ([`remove-expired-orders.ts`](src/tasks/remove-expired-orders.ts)): Sends expiry messages to Vulcan for short-term orders past their goodTilBlock.
- **Cancel Stale Orders** ([`cancel-stale-orders.ts`](src/tasks/cancel-stale-orders.ts)): Updates order status to CANCELED for orders that should have expired.
- **Remove Old Order Updates** ([`remove-old-order-updates.ts`](src/tasks/remove-old-order-updates.ts)): Cleans up cached stateful order updates from Redis after 30 seconds.

#### Funding & Rewards Tasks

- **Update Funding Payments** ([`update-funding-payments.ts`](src/tasks/update-funding-payments.ts)): Calculates funding payments for all positions at each funding index update using SQL script ([`update_funding_payments.sql`](src/scripts/update_funding_payments.sql)).
- **Aggregate Trading Rewards** ([`aggregate-trading-rewards.ts`](src/tasks/aggregate-trading-rewards.ts)): Aggregates trading rewards into daily, weekly, and monthly periods.
- **Create Leaderboard** ([`create-leaderboard.ts`](src/tasks/create-leaderboard.ts)): Generates leaderboard rankings based on PnL for various time periods (daily, weekly, monthly, yearly, all-time).

#### Compliance & User Management Tasks

- **Update Compliance Data** ([`update-compliance-data.ts`](src/tasks/update-compliance-data.ts)): Screens addresses using Elliptic API and updates compliance status.
- **Perform Compliance Status Transitions** ([`perform-compliance-status-transitions.ts`](src/tasks/perform-compliance-status-transitions.ts)): Transitions CLOSE_ONLY addresses to BLOCKED after configured period (default 7 days).
- **Subaccount Username Generator** ([`subaccount-username-generator.ts`](src/tasks/subaccount-username-generator.ts)): Generates unique usernames for new subaccounts using deterministic algorithm with adjectives and nouns.
- **Update Wallet Total Volume** ([`update-wallet-total-volume.ts`](src/tasks/update-wallet-total-volume.ts)): Updates total trading volume for each wallet address.
- **Update Affiliate Info** ([`update-affiliate-info.ts`](src/tasks/update-affiliate-info.ts)): Updates affiliate earnings, referral counts, and trading statistics.
- **Delete Old Firebase Notification Tokens** ([`delete-old-firebase-notification-tokens.ts`](src/tasks/delete-old-firebase-notification-tokens.ts)): Removes notification tokens older than 30 days.

#### Infrastructure & Maintenance Tasks

- **Update Research Environment** ([`update-research-environment.ts`](src/tasks/update-research-environment.ts)): Exports RDS snapshots to S3 and creates Athena tables for research queries.
- **Take Fast Sync Snapshot** ([`take-fast-sync-snapshot.ts`](src/tasks/take-fast-sync-snapshot.ts)): Creates RDS snapshots for fast-sync purposes.
- **Delete Old Fast Sync Snapshots** ([`delete-old-fast-sync-snapshots.ts`](src/tasks/delete-old-fast-sync-snapshots.ts)): Removes fast-sync snapshots older than 7 days.
- **Track Lag** ([`track-lag.ts`](src/tasks/track-lag.ts)): Monitors block height and time lag between validator, full nodes, and indexer.

## Configuration

Configuration is loaded from environment variables and parsed in [`src/config.ts`](src/config.ts). Environment-specific defaults are in `.env`, `.env.development`, `.env.test`.

### Core Settings

- **`SERVICE_NAME`** (string, default: `roundtable`): Service identifier for logging and metrics.
- **`NODE_ENV`** (string, default: `development`): Environment mode (`development`, `test`, `production`).
- **`START_DELAY_ENABLED`** (boolean, default: `true`): Enable random start delay to prevent thundering herd.
- **`MAX_START_DELAY_MS`** (integer, default: `180000`): Maximum start delay in milliseconds (3 minutes).
- **`MAX_CONCURRENT_RUNNING_TASKS`** (integer, default: `2`): Maximum number of tasks that can run concurrently. Should match `PG_POOL_MIN`.
- **`EXCEEDED_MAX_CONCURRENT_RUNNING_TASKS_DELAY_MS`** (integer, default: `1000`): Delay when max concurrent tasks exceeded.
- **`JITTER_FRACTION_OF_DELAY`** (number, default: `0.01`): Fraction of delay to add as random jitter.

### Database Configuration

- **`DB_HOSTNAME`** (string, required): PostgreSQL primary host.
- **`DB_READONLY_HOSTNAME`** (string, optional): PostgreSQL read replica host.
- **`DB_PORT`** (integer, default: `5432`): PostgreSQL port.
- **`DB_NAME`** (string, required): Database name.
- **`DB_USERNAME`** (string, required): Database username.
- **`DB_PASSWORD`** (string, required): Database password.
- **`PG_POOL_MAX`** (integer, default: `20`): Maximum PostgreSQL connection pool size.
- **`PG_POOL_MIN`** (integer, default: `2`): Minimum PostgreSQL connection pool size.

### Redis Configuration

- **`REDIS_URL`** (string, default: `redis://localhost:6382`): Primary Redis URL.
- **`REDIS_RECONNECT_TIMEOUT_MS`** (integer, default: `5000`): Redis reconnection timeout.

### Kafka Configuration

- **`KAFKA_BROKER_URLS`** (string, required): Comma-separated list of Kafka broker URLs.

### Loop Enablement Flags

Each task can be enabled/disabled independently (all boolean):
- `LOOPS_ENABLED_MARKET_UPDATER` (default: `true`)
- `LOOPS_ENABLED_DELETE_ZERO_PRICE_LEVELS` (default: `true`)
- `LOOPS_ENABLED_UNCROSS_ORDERBOOK` (default: `true`)
- `LOOPS_ENABLED_PNL_TICKS` (default: `true`)
- `LOOPS_ENABLED_REMOVE_EXPIRED_ORDERS` (default: `true`)
- `LOOPS_ORDERBOOK_INSTRUMENTATION` (default: `true`)
- `LOOPS_PNL_INSTRUMENTATION` (default: `true`)
- `LOOPS_CANCEL_STALE_ORDERS` (default: `true`)
- `LOOPS_ENABLED_UPDATE_RESEARCH_ENVIRONMENT` (default: `false`)
- `LOOPS_ENABLED_TAKE_FAST_SYNC_SNAPSHOTS` (default: `true`)
- `LOOPS_ENABLED_DELETE_OLD_FAST_SYNC_SNAPSHOTS` (default: `true`)
- `LOOPS_ENABLED_TRACK_LAG` (default: `false`)
- `LOOPS_ENABLED_REMOVE_OLD_ORDER_UPDATES` (default: `true`)
- `LOOPS_ENABLED_AGGREGATE_TRADING_REWARDS_DAILY` (default: `true`)
- `LOOPS_ENABLED_AGGREGATE_TRADING_REWARDS_WEEKLY` (default: `true`)
- `LOOPS_ENABLED_AGGREGATE_TRADING_REWARDS_MONTHLY` (default: `true`)
- `LOOPS_ENABLED_SUBACCOUNT_USERNAME_GENERATOR` (default: `true`)
- `LOOPS_ENABLED_LEADERBOARD_PNL_ALL_TIME` (default: `false`)
- `LOOPS_ENABLED_LEADERBOARD_PNL_DAILY` (default: `false`)
- `LOOPS_ENABLED_LEADERBOARD_PNL_WEEKLY` (default: `false`)
- `LOOPS_ENABLED_LEADERBOARD_PNL_MONTHLY` (default: `false`)
- `LOOPS_ENABLED_LEADERBOARD_PNL_YEARLY` (default: `false`)
- `LOOPS_ENABLED_UPDATE_WALLET_TOTAL_VOLUME` (default: `true`)
- `LOOPS_ENABLED_UPDATE_AFFILIATE_INFO` (default: `true`)
- `LOOPS_ENABLED_DELETE_OLD_FIREBASE_NOTIFICATION_TOKENS` (default: `true`)
- `LOOPS_ENABLED_REFRESH_VAULT_PNL` (default: `true`)
- `LOOPS_ENABLED_CACHE_ORDERBOOK_MID_PRICES` (default: `true`)
- `LOOPS_ENABLED_UPDATE_FUNDING_PAYMENTS` (default: `true`)
- `LOOPS_ENABLED_UPDATE_PNL` (default: `true`)

### Loop Interval Configuration

Each task has a configurable interval in milliseconds (all integer):
- `LOOPS_INTERVAL_MS_MARKET_UPDATER` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_DELETE_ZERO_PRICE_LEVELS` (default: `120000` - 2 minutes)
- `LOOPS_INTERVAL_MS_UNCROSS_ORDERBOOK` (default: `15000` - 15 seconds)
- `LOOPS_INTERVAL_MS_PNL_TICKS` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_REMOVE_EXPIRED_ORDERS` (default: `120000` - 2 minutes)
- `LOOPS_INTERVAL_MS_ORDERBOOK_INSTRUMENTATION` (default: `60000` - 1 minute)
- `LOOPS_INTERVAL_MS_PNL_INSTRUMENTATION` (default: `3600000` - 1 hour)
- `LOOPS_INTERVAL_MS_CANCEL_STALE_ORDERS` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_UPDATE_RESEARCH_ENVIRONMENT` (default: `3600000` - 1 hour)
- `LOOPS_INTERVAL_MS_TAKE_FAST_SYNC_SNAPSHOTS` (default: `14400000` - 4 hours)
- `LOOPS_INTERVAL_MS_DELETE_OLD_FAST_SYNC_SNAPSHOTS` (default: `86400000` - 1 day)
- `LOOPS_INTERVAL_MS_UPDATE_COMPLIANCE_DATA` (default: `300000` - 5 minutes)
- `LOOPS_INTERVAL_MS_PERFORM_COMPLIANCE_STATUS_TRANSITIONS` (default: `3600000` - 1 hour)
- `LOOPS_INTERVAL_MS_TRACK_LAG` (default: `10000` - 10 seconds)
- `LOOPS_INTERVAL_MS_REMOVE_OLD_ORDER_UPDATES` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_AGGREGATE_TRADING_REWARDS` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_SUBACCOUNT_USERNAME_GENERATOR` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_LEADERBOARD_PNL_ALL_TIME` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_LEADERBOARD_PNL_DAILY` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_LEADERBOARD_PNL_WEEKLY` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_LEADERBOARD_PNL_MONTHLY` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_LEADERBOARD_PNL_YEARLY` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_UPDATE_WALLET_TOTAL_VOLUME` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_UPDATE_AFFILIATE_INFO` (default: `30000` - 30 seconds)
- `LOOPS_INTERVAL_MS_DELETE_FIREBASE_NOTIFICATION_TOKENS_MONTHLY` (default: `2592000000` - 30 days)
- `LOOPS_INTERVAL_MS_REFRESH_VAULT_PNL` (default: `300000` - 5 minutes)
- `LOOPS_INTERVAL_MS_CACHE_ORDERBOOK_MID_PRICES` (default: `10000` - 10 seconds)
- `LOOPS_INTERVAL_MS_UPDATE_FUNDING_PAYMENTS` (default: `60000` - 1 minute)
- `LOOPS_INTERVAL_MS_UPDATE_PNL` (default: `60000` - 1 minute)

### Lock Multipliers

Extended lock multipliers for long-running tasks (all integer):
- `MARKET_UPDATER_LOCK_MULTIPLIER` (default: `10`)
- `DELETE_ZERO_PRICE_LEVELS_LOCK_MULTIPLIER` (default: `1`)
- `UNCROSS_ORDERBOOK_LOCK_MULTIPLIER` (default: `1`)
- `PNL_TICK_UPDATE_LOCK_MULTIPLIER` (default: `20`)
- `SUBACCOUNT_USERNAME_GENERATOR_LOCK_MULTIPLIER` (default: `5`)
- `UPDATE_FUNDING_PAYMENTS_LOCK_MULTIPLIER` (default: `3` - increase for initial backfill)
- `UPDATE_PNL_LOCK_MULTIPLIER` (default: `3` - increase for initial backfill)

### Task-Specific Configuration

- **`PNL_TICK_UPDATE_INTERVAL_MS`** (integer, default: `3600000` - 1 hour): Interval for PnL tick updates.
- **`PNL_TICK_MAX_ROWS_PER_UPSERT`** (integer, default: `1000`): Maximum rows per PnL tick upsert.
- **`PNL_TICK_MAX_ACCOUNTS_PER_RUN`** (integer, default: `65000`): Maximum accounts to process per PnL tick run.
- **`LEADERBOARD_PNL_MAX_ROWS_PER_UPSERT`** (integer, default: `1000`): Maximum rows per leaderboard upsert.
- **`BLOCKS_TO_DELAY_EXPIRY_BEFORE_SENDING_REMOVES`** (integer, default: `20`): Block delay before sending order expiry messages.
- **`CANCEL_STALE_ORDERS_QUERY_BATCH_SIZE`** (integer, default: `10000`): Batch size for stale order queries.
- **`OLD_CACHED_ORDER_UPDATES_WINDOW_MS`** (integer, default: `30000` - 30 seconds): Age threshold for removing cached order updates.
- **`AGGREGATE_TRADING_REWARDS_MAX_INTERVAL_SIZE_MS`** (integer, default: `3600000` - 1 hour): Maximum interval size for trading reward aggregation.
- **`AGGREGATE_TRADING_REWARDS_CHUNK_SIZE`** (integer, default: `50`): Chunk size for trading reward operations.
- **`STALE_ORDERBOOK_LEVEL_THRESHOLD_SECONDS`** (integer, default: `10`): Threshold for considering orderbook levels stale.
- **`SUBACCOUNT_USERNAME_SUFFIX_RANDOM_DIGITS`** (integer, default: `3`): Number of random digits in username suffix.
- **`SUBACCOUNT_USERNAME_BATCH_SIZE`** (integer, default: `2000`): Batch size for username generation.
- **`ATTEMPT_PER_SUBACCOUNT`** (integer, default: `3`): Number of attempts to generate username per subaccount.
- **`TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS`** (integer, default: `900000` - 15 minutes): Time window for refreshing vault PnL views.
- **`VOLUME_ELIGIBILITY_THRESHOLD`** (integer, default: `10000`): Minimum trading volume (USDC) for affiliate eligibility.

### Compliance Configuration

- **`ACTIVE_ADDRESS_THRESHOLD_SECONDS`** (integer, default: `86400` - 1 day): Threshold for considering an address active.
- **`MAX_COMPLIANCE_DATA_AGE_SECONDS`** (integer, default: `2630000` - ~1 month): Maximum age for compliance data.
- **`MAX_ACTIVE_COMPLIANCE_DATA_AGE_SECONDS`** (integer, default: `86400` - 1 day): Maximum age for active address compliance data.
- **`MAX_COMPLIANCE_DATA_QUERY_PER_LOOP`** (integer, default: `100`): Maximum compliance queries per loop.
- **`COMPLIANCE_PROVIDER_QUERY_BATCH_SIZE`** (integer, default: `15`): Batch size for compliance provider queries (Elliptic rate limit: 15/s).
- **`COMPLIANCE_PROVIDER_QUERY_DELAY_MS`** (integer, default: `1000` - 1 second): Delay between compliance query batches.
- **`CLOSE_ONLY_TO_BLOCKED_DAYS`** (integer, default: `7`): Days before transitioning CLOSE_ONLY to BLOCKED status.

### AWS Configuration (for research environment)

- **`AWS_ACCOUNT_ID`** (string, required): AWS account ID.
- **`AWS_REGION`** (string, required): AWS region.
- **`S3_BUCKET_ARN`** (string, required): S3 bucket ARN for exports.
- **`FAST_SYNC_SNAPSHOT_IDENTIFIER_PREFIX`** (string, default: `fast-sync`): Prefix for fast-sync snapshots.
- **`ECS_TASK_ROLE_ARN`** (string, required): ECS task role ARN.
- **`KMS_KEY_ARN`** (string, required): KMS key ARN for encryption.
- **`RDS_INSTANCE_NAME`** (string, required): RDS instance name.
- **`ATHENA_CATALOG_NAME`** (string, default: `AwsDataCatalog`): Athena catalog name.
- **`ATHENA_DATABASE_NAME`** (string, default: `default`): Athena database name.
- **`ATHENA_WORKING_GROUP`** (string, default: `primary`): Athena working group.
- **`SKIP_TO_ATHENA_TABLE_WRITING`** (boolean, default: `false`): Skip to Athena table writing (for testing).

### Lag Tracking Configuration

- **`TRACK_LAG_INDEXER_FULL_NODE_URL`** (string, default: `""`): Indexer full node URL (e.g., `http://11.11.11.11:26657`).
- **`TRACK_LAG_VALIDATOR_URL`** (string, default: `""`): Validator URL.
- **`TRACK_LAG_OTHER_FULL_NODE_URL`** (string, default: `""`): Other full node URL for comparison.

## Running Locally

### Prerequisites

- Node.js 18+ and `pnpm` installed.
- PostgreSQL database running and accessible.
- Redis instance running and accessible.
- Kafka broker running and accessible (for market update tasks).
- Environment variables configured (see `.env.development` for examples).

### Steps

1. **Install dependencies** (from monorepo root):
   ```bash
   pnpm install
   ```

2. **Build the service**:
   ```bash
   cd indexer/services/roundtable
   pnpm run build
   ```

3. **Set up environment variables**:
   - Copy `.env.development` or create a `.env` file with required variables.
   - Ensure `DB_HOSTNAME`, `DB_NAME`, `DB_USERNAME`, `DB_PASSWORD`, `REDIS_URL`, and `KAFKA_BROKER_URLS` are set.
   - For AWS tasks, set `AWS_ACCOUNT_ID`, `AWS_REGION`, `S3_BUCKET_ARN`, etc.

4. **Run the service**:
   ```bash
   pnpm run start
   ```
   The service will start all enabled task loops based on configuration.

5. **Run in development mode** (with dotenv-flow):
   ```bash
   pnpm run build
   node -r dotenv-flow/config build/src/index.js
   ```

### Development Commands

- **`pnpm run build`**: Compile TypeScript to JavaScript in `build/` directory (includes copying SQL scripts).
- **`pnpm run build:watch`**: Compile in watch mode.
- **`pnpm run start`**: Run the service in production mode with DataDog tracing.
- **`pnpm run lint`**: Run ESLint on the codebase.
- **`pnpm run lint:fix`**: Run ESLint with auto-fix.

## Testing

### Run All Tests

From the `indexer/services/roundtable` directory:

pnpm test
### Run Tests with Coverage
pnpm run coverage
### Test Environment
- Tests use the `NODE_ENV=test` environment.

- Test database configuration is loaded from `.env.test`.

- Default test database settings:

  - `DB_NAME=dydx_test`

  - `DB_USERNAME=dydx_test`

  - `DB_PASSWORD=dydxserver123`

  - `DB_PORT=5436`
### Test Setup
- Global setup: [`jest.globalSetup.js`](jest.globalSetup.js) loads environment variables via dotenv-flow.

- Per-test setup: [`jest.setup.js`](jest.setup.js) runs before each test file (currently empty).

- Configuration: [`jest.config.js`](jest.config.js) uses the base configuration from `@dydxprotocol-indexer/dev`.
## Observability & Operations
### Logging
- Logs are structured JSON written to stdout via `@dydxprotocol-indexer/base` logger.

- Key log fields:

  - `at`: Source location (e.g., `market-updater#runTask`).

  - `message`: Human-readable message.

  - `error`: Error object with message and stack trace (if applicable).

  - `taskName`: Name of the task being executed.

  - Task-specific fields (e.g., `blockHeight`, `subaccountId`, `ticker`).
### Metrics
- Metrics are emitted via `@dydxprotocol-indexer/base` stats module (StatsD format).

- Key metrics:

  - `roundtable.loops.<taskName>.started`: Task start counter.

  - `roundtable.loops.<taskName>.completed`: Task completion counter (1 = success, 0 = failure).

  - `roundtable.loops.<taskName>.timing`: Task execution duration in milliseconds.

  - `roundtable.loops.duration_ratio`: Ratio of task duration to interval (tagged by `taskName`).

  - `roundtable.loops.exceeded_max_concurrent_tasks`: Counter for when max concurrent tasks exceeded.

  - `roundtable.loops.could_not_acquire_extended_lock`: Counter for extended lock acquisition failures.

  - `roundtable.num_connections`: Number of active connections (for various resources).

  - `roundtable.pnl_stale_subaccounts`: Number of subaccounts with stale PnL data.

  - `roundtable.num_stale_orders`: Number of stale orders found.

  - `roundtable.crossed_orderbook`: Counter for crossed orderbooks (tagged by `ticker`).

  - `roundtable.uncross_orderbook_succeed`: Counter for successful orderbook uncrossing.

  - `roundtable.uncross_orderbook_failed`: Counter for failed orderbook uncrossing.

  - Task-specific metrics (see individual task files for details).
### Common Failure Modes
- **Database connection failures**: Service will fail to start or tasks will fail. Check `DB_HOSTNAME`, `DB_PORT`, and credentials. Verify PostgreSQL is running and accessible.

- **Redis connection failures**: Distributed locking will fail, causing tasks to skip execution. Check `REDIS_URL` and verify Redis is running.

- **Kafka connection failures**: Market update tasks will fail to publish messages. Check `KAFKA_BROKER_URLS` and Kafka broker status.

- **Lock contention**: Multiple instances may compete for locks. This is expected behavior; only one instance will execute each task.

- **Task timeout**: Long-running tasks may exceed their interval. Increase lock multipliers or optimize task logic.

- **Compliance provider rate limits**: Elliptic API has rate limits (15 requests/second). Tasks will delay between batches to avoid hitting limits.

- **AWS API failures**: Snapshot and export tasks may fail due to AWS API errors. Check IAM permissions and service quotas.

- **Stale PnL data**: If PnL calculation tasks fall behind, alerts will be emitted. Check task execution frequency and database performance.
### Known Limitations
- Tasks are not guaranteed to execute exactly on schedule due to lock contention and processing time.

- Extended locks prevent multiple instances from running the same task, but don't guarantee exactly-once execution.

- PnL calculations are eventually consistent and may lag behind real-time data.

- Compliance data is cached and may be stale for up to 1 day (active addresses) or 1 month (inactive addresses).

- Username generation uses a deterministic algorithm, so collisions are possible (handled with retries).

- Orderbook uncrossing relies on timestamp-based staleness detection, which may not catch all crossed levels.
## Deployment & Runtime
### Deployment
- Roundtable is deployed as a containerized service (Docker).

- Dockerfile: Located at monorepo root (`../../Dockerfile`).

- Build command: `docker build -t roundtable -f indexer/Dockerfile --target roundtable .`

- The service uses multi-stage builds with the `roundtable` target.
### Runtime Characteristics
- **Stateless**: No local state; all coordination is via Redis locks.

- **Horizontal scaling**: Multiple replicas can run with distributed locking ensuring only one instance executes each task.

- **Resource requirements**: CPU and memory usage depend on enabled tasks and data volume. Typical memory usage is 512MB-2GB per replica.

- **Expected replicas**: Typically 2-3 replicas for high availability (only one will execute each task at a time).

- **Graceful shutdown**: Service handles `SIGTERM` signals, disconnects from Redis, and exits cleanly.
### Kubernetes Manifests
Deployment manifests are not present in this directory. Check the monorepo root or a separate `infra/` or `k8s/` directory for Helm charts or Kubernetes YAML files.
## Directory Layout
- **`src/`**: Source code.

  - **`tasks/`**: Individual task implementations (one file per task).

  - **`helpers/`**: Utility modules.

    - **`loops-helper.ts`**: Loop orchestration with locking and timing.

    - **`aws.ts`**: AWS SDK utilities for RDS, S3, and Athena.

    - **`pnl-ticks-helper.ts`**: PnL calculation logic.

    - **`compliance-clients.ts`**: Compliance provider integration.

    - **`redis.ts`**: Redis client initialization.

    - **`websocket.ts`**: WebSocket message creation and publishing.

    - **`sql.ts`**: SQL query generation for Athena.

    - **`helpers.ts`**: Miscellaneous utilities.

    - **`types.ts`**: TypeScript type definitions.

    - **`constants.ts`**: Shared constants.

    - **`adjectives.json`**, **`nouns.json`**: Word lists for username generation.

  - **`lib/`**: Shared library code.

    - **`athena-ddl-tables/`**: Athena table DDL generation for each database table.

    - **`constants.ts`**: Service-level constants.

  - **`scripts/`**: SQL scripts for complex operations.

    - **`update_funding_payments.sql`**: SQL for calculating funding payments.

    - **`update_pnl.sql`**: SQL for calculating PnL changes.

  - **`config.ts`**: Configuration schema and environment variable parsing.

  - **`index.ts`**: Service entrypoint.

- **`__tests__/`**: Test files (not documented here).

- **`build/`**: Compiled JavaScript output (generated by `pnpm run build`).

- **`patches/`**: Patch files for npm dependencies (currently empty, `.gitkeep` only).

- **`.env`**: Environment variables for local development (not committed, create from `.env.development`).

- **`.env.development`**: Default environment variables for development.

- **`.env.test`**: Environment variables for testing.

- **`package.json`**: Node.js package manifest and scripts.

- **`tsconfig.json`**: TypeScript compiler configuration.

- **`tsconfig.eslint.json`**: TypeScript configuration for ESLint.

- **`jest.config.js`**: Jest test configuration.

- **`jest.setup.js`**: Jest setup script.

- **`jest.globalSetup.js`**: Jest global setup script.

- **`.eslintrc.js`**: ESLint configuration.
## Related Documents
- **Indexer Overview**: [`../../README.md`](../../README.md) - Overview of the entire indexer system.

- **Monorepo Root README**: [`../../../README.md`](../../../README.md) - dYdX v4 Chain overview and getting started guide.

- **Comlink Service**: [`../comlink/README.md`](../comlink/README.md) - REST API service that serves data maintained by Roundtable.

- **Socks Service**: [`../socks/README.md`](../socks/README.md) - WebSocket service that receives market updates from Roundtable.
