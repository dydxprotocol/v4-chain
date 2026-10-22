# Bazooka

Bazooka is a deployment and migration service for the dYdX Indexer infrastructure. It runs as an AWS Lambda function to orchestrate database migrations, Kafka topic management, Redis cache clearing, and other deployment-related operations. The service is designed to be invoked during deployment workflows to prepare or reset the indexer environment.

## Responsibilities & Scope

- Execute Knex database migrations and rollbacks
- Create and partition Kafka topics with appropriate replication factors
- Clear data from Kafka topics without deleting the topics themselves
- Clear all data from Redis cache
- Reset or clear database tables while preserving schema
- Send stateful orders to Vulcan during indexer fast sync operations

**Out of scope:**
- Does not handle application runtime logic or request processing
- Does not perform continuous monitoring or health checks
- Does not manage infrastructure provisioning (only operates on existing resources)

## Architecture & Dependencies

### Internal Structure

- `src/index.ts` – Main Lambda handler and orchestration logic
- `src/config.ts` – Environment variable parsing and configuration schema
- `src/vulcan-helpers.ts` – Logic for sending stateful order messages to Vulcan
- `src/redis.ts` – Redis client initialization
- `src/constants.ts` – Shared constants (e.g., `ZERO` for Big.js)

### Internal Dependencies

- `@dydxprotocol-indexer/base` – Logging, utilities, and base configuration
- `@dydxprotocol-indexer/kafka` – Kafka admin client and producer
- `@dydxprotocol-indexer/postgres` – Database helpers, migrations, and ORM models
- `@dydxprotocol-indexer/redis` – Redis client utilities
- `@dydxprotocol-indexer/v4-proto-parser` – Protocol buffer parsing utilities
- `@dydxprotocol-indexer/v4-protos` – Generated protobuf types

### External Dependencies

- **PostgreSQL** – Primary data store for indexer state
- **Kafka** – Message bus for inter-service communication
- **Redis** – Cache layer for indexer data
- **AWS Lambda** – Execution environment

### Processing Flow

1. Lambda receives event with operation flags (migrate, clear_db, etc.)
2. Validates force flag if attempting breaking operations in production
3. Executes operations in sequence:
   - Database reset/rollback/migration
   - Database data clearing
   - Kafka topic creation and partitioning
   - Kafka topic data clearing
   - Redis cache clearing
   - Stateful order synchronization to Vulcan
4. Returns success/failure status

## Public Interface

### Lambda Handler Event Schema

`BazookaEventJson`
  - `migrate: boolean` – Run Knex migrations
  - `rollback: boolean` – Rollback latest migration batch
  - `clear_db: boolean` – Clear database data without dropping tables
  - `reset_db: boolean` – Drop and recreate all database schemas
  - `create_kafka_topics: boolean` – Create Kafka topics with partitions
  - `clear_kafka_topics: boolean` – Delete all records from Kafka topics
  - `clear_redis: boolean` – Flush all Redis data
  - `send_stateful_orders_to_vulcan: boolean` – Sync stateful orders to Vulcan
  - `force: boolean` – Required for breaking operations in protected environments

### Kafka Topics Managed

- `TO_ENDER` (1 partition)
- `TO_VULCAN` (210 partitions)
- `TO_WEBSOCKETS_ORDERBOOKS` (1 partition)
- `TO_WEBSOCKETS_SUBACCOUNTS` (3 partitions)
- `TO_WEBSOCKETS_TRADES` (1 partition)
- `TO_WEBSOCKETS_MARKETS` (1 partition)
- `TO_WEBSOCKETS_CANDLES` (1 partition)
- `TO_WEBSOCKETS_BLOCK_HEIGHT` (1 partition)

All topics are created with 3 replicas by default.

## Configuration

### Environment Variables

- `PREVENT_BREAKING_CHANGES_WITHOUT_FORCE` (boolean, default: `true`) – Requires `force: true` flag for destructive operations
- `CLEAR_KAFKA_TOPIC_RETRY_MS` (integer, default: `3000`) – Base retry delay in milliseconds for Kafka topic clearing
- `CLEAR_KAFKA_TOPIC_MAX_RETRIES` (integer, default: `3`) – Maximum retry attempts for Kafka operations
- `PG_POOL_MAX` (integer, default: `10` in production) – PostgreSQL connection pool size
- `DB_NAME`, `DB_USERNAME`, `DB_PASSWORD`, `DB_PORT` – PostgreSQL connection parameters
- `REDIS_URL`, `REDIS_RECONNECT_TIMEOUT_MS` – Redis connection parameters
- Kafka configuration inherited from `@dydxprotocol-indexer/kafka` package

### Configuration Files

- `.env` – Base service-level variables
- `.env.production` – Production-specific overrides
- `.env.test` – Test environment configuration
- `src/config.ts` – Configuration schema and parsing logic

## Running Locally

Bazooka is designed to run as a Lambda function, but can be tested locally:

1. **Install dependencies:**
   ```bash
   pnpm install
   ```

2. **Set up local environment:**
   - Ensure PostgreSQL, Kafka, and Redis are running locally or via Docker
   - Copy `.env.test` or create a `.env.local` with appropriate connection strings

3. **Build the service:**
   ```bash
   pnpm run build
   ```

4. **Run locally (requires Lambda runtime simulation):**
   ```bash
   pnpm start
   ```

5. **Invoke with test event:**
   Create a test event JSON matching `BazookaEventJson` schema and invoke the handler programmatically or via AWS SAM CLI.

## Logging

- All operations log to stdout using structured logging from `@dydxprotocol-indexer/base`

- Key log fields:

  - `at` – Source location (e.g., `index#handler`)

  - `message` – Human-readable description

  - `error` – Error details for failures

  - `topic`, `attempt`, `topicMetadata` – Kafka operation context

## Health & Status

- Lambda execution status returned in response body

- Bugsnag integration for error tracking (initialized via `startBugsnag()`)

## Directory Layout

- `src/` – Source code

  - `index.ts` – Main Lambda handler and orchestration

  - `config.ts` – Configuration schema and parsing

  - `vulcan-helpers.ts` – Stateful order synchronization logic

  - `redis.ts` – Redis client setup

  - `constants.ts` – Shared constants

- `__tests__/` – Test files

  - `index.test.ts` – Handler and Kafka tests

  - `vulcan-helpers.test.ts` – Vulcan message generation tests

- `patches/` – npm package patches (currently empty, `.gitkeep` placeholder)

- `jest.config.js` – Jest test configuration

- `jest.setup.js` – Test setup (runs before each test file)

- `jest.globalSetup.js` – Global test setup (runs once before all tests)

- `tsconfig.json` – TypeScript compiler configuration

- `tsconfig.eslint.json` – TypeScript configuration for ESLint

- `.eslintrc.js` – ESLint configuration

- `package.json` – Dependencies and scripts

- `.env`, `.env.production`, `.env.test` – Environment-specific configuration

## Related Documents

- [Indexer README](../../README.md) – Overview of the indexer services

- [Root README](../../../README.md) – dYdX Chain repository overview
