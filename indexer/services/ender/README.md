# Ender

Ender is the on-chain data archival service for the dYdX v4 indexer. It consumes blocks from the dYdX Chain via Kafka, processes all events within each block, persists data to PostgreSQL, and publishes updates to downstream services via Kafka topics for real-time websocket delivery.

Ender serves as the primary ingestion point for blockchain data, transforming protocol events into queryable database records and real-time message streams.

## Responsibilities & Scope

- Consume `IndexerTendermintBlock` messages from Kafka (`to-ender` topic).
- Validate and parse all on-chain events (order fills, subaccount updates, transfers, market events, stateful orders, funding updates, etc.).
- Execute SQL-based event handlers to persist data to PostgreSQL (blocks, transactions, orders, fills, positions, transfers, markets, assets, etc.).
- Generate and publish Kafka messages for downstream websocket services (subaccounts, trades, markets, candles, block height).
- Maintain in-memory caches for perpetual markets, assets, liquidity tiers, candles, and block height.
- Update Redis caches for order state and funding rates.
- Generate candles from trade data at multiple resolutions (1m, 5m, 15m, 30m, 1h, 4h, 1d).

**Out of scope:**
- Does not serve HTTP/gRPC endpoints for external clients.
- Does not perform user authentication or authorization.
- Does not generate orderbook data (handled by Vulcan).

## Architecture & Dependencies

### Internal Structure

- **Event Processing Pipeline:**
  - `BlockProcessor` orchestrates processing of all events in a block.
  - `Validator` classes validate event structure and data integrity.
  - `Handler` classes process individual events and generate Kafka messages.
  - SQL functions in `src/scripts/` perform database operations.
  - `BatchedHandlers` and `SyncHandlers` manage parallelization and ordering.

- **Caching Layer:**
  - In-memory caches for perpetual markets, assets, liquidity tiers, candles, and block height.
  - Redis integration for order state and funding rates.

- **Kafka Integration:**
  - `KafkaPublisher` aggregates and batches outbound messages.
  - Separate topics for subaccounts, trades, markets, candles, block height, and Vulcan updates.

### Internal Dependencies

- `@dydxprotocol-indexer/postgres` – database models, queries, and refreshers.
- `@dydxprotocol-indexer/kafka` – Kafka consumer/producer utilities.
- `@dydxprotocol-indexer/redis` – Redis client and cache utilities.
- `@dydxprotocol-indexer/v4-protos` – Protocol buffer definitions.
- `@dydxprotocol-indexer/v4-proto-parser` – Protocol parsing utilities.
- `@dydxprotocol-indexer/notifications` – Firebase notification integration.

### External Dependencies

- **PostgreSQL** – primary datastore for all indexed blockchain data.
- **Kafka** – message bus for consuming blocks and publishing updates.
- **Redis** – cache for order state, funding rates, and orderbook mid-prices.
- **Firebase** (optional) – push notifications for order fills and triggers.

### Processing Flow

1. Consume `IndexerTendermintBlock` from Kafka `to-ender` topic.
2. Validate block height against cache to prevent reprocessing.
3. Group events by transaction and block events.
4. Validate each event using appropriate `Validator` class.
5. Organize events into parallelizable batches based on parallelization IDs.
6. Execute SQL block processor function with all decoded events.
7. Process handlers in batches (parallel where safe, sequential where required).
8. Generate candles from trade data.
9. Publish consolidated Kafka messages to downstream topics.
10. Update in-memory caches and commit transaction.

## Public Interface

Ender does not expose HTTP/gRPC endpoints. It operates as a Kafka consumer and producer.

### Kafka Topics Consumed

- **`to-ender`** – Receives `IndexerTendermintBlock` messages from the protocol node.
  - Handler: [`src/lib/on-message.ts`](src/lib/on-message.ts)

### Kafka Topics Published

- **`to-websockets-subaccounts`** – Subaccount updates (positions, orders, fills, transfers).
  - Message type: `SubaccountMessage`
  - Key: `IndexerSubaccountId`

- **`to-websockets-trades`** – Trade events for market data.
  - Message type: `TradeMessage`
  - Grouped by `clobPairId`

- **`to-websockets-markets`** – Market and oracle price updates.
  - Message type: `MarketMessage`

- **`to-websockets-candles`** – Candle data at multiple resolutions.
  - Message type: `CandleMessage`

- **`to-websockets-block-height`** – Block height updates.
  - Message type: `BlockHeightMessage`

- **`to-vulcan`** – Off-chain order updates for orderbook management.
  - Message type: `OffChainUpdateV1`
  - Key: Order ID hash

## Configuration

Configuration is managed via environment variables and `.env` files.

### Environment Variables

- **`NODE_ENV`** (string) – Environment mode (`development`, `test`, `production`). Default: `development`.
- **`SERVICE_NAME`** (string) – Service identifier for logging and metrics. Default: `ender`.
- **`DB_NAME`** (string) – PostgreSQL database name. Default: `dydx_dev`.
- **`DB_USERNAME`** (string) – PostgreSQL username. Default: `dydx_dev`.
- **`DB_PASSWORD`** (string) – PostgreSQL password. Default: `dydxserver123`.
- **`DB_HOSTNAME`** (string) – PostgreSQL host. Default: `postgres`.
- **`DB_PORT`** (integer) – PostgreSQL port. Default: `5432`.
- **`DB_READONLY_HOSTNAME`** (string) – Read-only PostgreSQL replica host. Default: `postgres`.
- **`PG_POOL_MAX`** (integer) – Maximum PostgreSQL connection pool size. Default: `10` (production), `2` (development).
- **`PG_POOL_MIN`** (integer) – Minimum PostgreSQL connection pool size. Default: `1` (development).
- **`REDIS_URL`** (string) – Redis connection URL. Required.
- **`REDIS_RECONNECT_TIMEOUT_MS`** (integer) – Redis reconnection timeout. Default: `5000`.
- **`KAFKA_BROKER_URLS`** (string) – Comma-separated Kafka broker URLs. Default: `kafka:9092`.
- **`KAFKA_MAX_BATCH_WEBSOCKET_MESSAGE_SIZE_BYTES`** (integer) – Maximum Kafka batch size for websocket messages.
- **`SEND_WEBSOCKET_MESSAGES`** (boolean) – Enable/disable publishing to websocket topics. Default: `true`.
- **`SKIP_STATEFUL_ORDER_UUIDS`** (string) – Comma-separated list of order UUIDs to skip processing. Default: empty.
- **`ORDERBOOK_MID_PRICE_REFRESH_INTERVAL_MS`** (integer) – Interval for refreshing orderbook mid-prices from Redis. Default: `10000` (10 seconds).

### Configuration Files

- [`.env`](src/.env) – Base environment variables.
- [`.env.development`](src/.env.development) – Development overrides.
- [`.env.production`](src/.env.production) – Production overrides.
- [`.env.test`](src/.env.test) – Test environment overrides.
- [`tsconfig.json`](tsconfig.json) – TypeScript compiler configuration.
- [`package.json`](package.json) – Node.js dependencies and scripts.

## Running Locally

### Prerequisites

- Node.js 18+
- PostgreSQL 13+
- Kafka
- Redis
- Docker and Docker Compose (recommended)

### Steps

1. **Start dependencies:**
   ```bash
   # From repository root
   docker compose up postgres kafka redis
   ```

2. **Install dependencies:**
   ```bash
   # From indexer/services/ender
   pnpm install
   ```

3. **Run database migrations:**
   ```bash
   # From indexer/packages/postgres
   pnpm run migrate:up
   ```

4. **Build the service:**
   ```bash
   pnpm run build
   ```

5. **Start Ender:**
   ```bash
   pnpm start
   ```

   Or for development with auto-reload:
   ```bash
   pnpm run build:watch
   # In another terminal:
   pnpm start
   ```

### Environment Setup

Copy `.env.development` to `.env` and adjust values as needed:
