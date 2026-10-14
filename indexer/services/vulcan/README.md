# Vulcan

Vulcan is the off-chain update processing service for the dYdX v4 Indexer. It consumes off-chain order updates (placements, updates, and removals) from Kafka, maintains real-time order state in Redis, updates orderbook price levels, and forwards relevant updates to WebSocket clients via Kafka topics.

## Responsibilities & Scope

- Consume off-chain order updates from the `to-vulcan` Kafka topic.
- Maintain real-time order state in Redis caches (OrdersCache, OrdersDataCache, SubaccountOrderIdsCache).
- Update orderbook price levels in Redis (OrderbookLevelsCache) based on order placements, updates, and removals.
- Track canceled orders in Redis (CanceledOrdersCache) for proper status handling.
- Cache stateful order updates (StatefulOrderUpdatesCache) for orders placed before their on-chain events are processed.
- Forward order updates to WebSocket clients by publishing messages to Kafka topics (`to-websockets-subaccounts`, `to-websockets-orderbooks`).
- Handle order replacements and expiry verification for indexer-expired orders.
- Support both individual message processing and batch processing modes.

**Out of scope:**
- Does not persist orders to PostgreSQL (handled by Ender service).
- Does not manage WebSocket connections directly (handled by Socks service).
- Does not process on-chain events or blockchain data.
- Does not perform order matching or settlement.

## Architecture & Dependencies

### Internal Structure

- **Handlers** (`src/handlers/`): Process specific off-chain update types.
  - `OrderPlaceHandler`: Handles order placements and replacements.
  - `OrderUpdateHandler`: Handles order fill updates and orderbook level adjustments.
  - `OrderRemoveHandler`: Handles order cancellations and removals.
- **Message Processing** (`src/lib/`): Kafka message consumption and routing.
  - `on-message.ts`: Routes individual messages to appropriate handlers.
  - `on-batch.ts`: Processes batches of messages with heartbeat management.
  - `send-message-helper.ts`: Batches and sends WebSocket messages to Kafka.
- **Helpers** (`src/helpers/`): Utility functions for Redis, Kafka, and order operations.

### Internal Dependencies

- `@dydxprotocol-indexer/base`: Shared utilities, logging, stats, and error handling.
- `@dydxprotocol-indexer/kafka`: Kafka consumer/producer and message utilities.
- `@dydxprotocol-indexer/postgres`: Database models, perpetual market cache, and block height cache.
- `@dydxprotocol-indexer/redis`: Redis caches for orders, orderbook levels, and canceled orders.
- `@dydxprotocol-indexer/v4-protos`: Protocol buffer definitions for off-chain updates.
- `@dydxprotocol-indexer/v4-proto-parser`: Order flag parsing and validation utilities.

### External Dependencies

- **Kafka**: Message bus for consuming off-chain updates and publishing WebSocket messages.
- **Redis**: In-memory cache for order state, orderbook levels, and canceled orders.
- **PostgreSQL**: Indirectly via perpetual market and block height caches (read-only).

### Processing Flow

1. Kafka consumer receives off-chain update message from `to-vulcan` topic.
2. Message is decoded into `OffChainUpdateV1` protobuf.
3. Message is routed to appropriate handler based on update type (orderPlace, orderUpdate, orderRemove).
4. Handler validates the message and updates Redis caches:
   - **OrderPlace**: Adds/replaces order in caches, removes from canceled cache, sends cached stateful order updates.
   - **OrderUpdate**: Updates total filled quantums, adjusts orderbook price levels.
   - **OrderRemove**: Removes order from caches, updates orderbook levels, adds to canceled cache.
5. Handler creates WebSocket messages (subaccount and/or orderbook updates).
6. Messages are batched and sent to Kafka topics (`to-websockets-subaccounts`, `to-websockets-orderbooks`).
7. Socks service consumes these messages and forwards to WebSocket clients.

## Public Interface

### Kafka Topics

#### Consumed Topics

- **`to-vulcan`**: Consumes off-chain order updates (placements, updates, removals).
  - Message format: `OffChainUpdateV1` protobuf.
  - Handlers: [`src/handlers/order-place-handler.ts`](src/handlers/order-place-handler.ts), [`src/handlers/order-update-handler.ts`](src/handlers/order-update-handler.ts), [`src/handlers/order-remove-handler.ts`](src/handlers/order-remove-handler.ts).

#### Published Topics

- **`to-websockets-subaccounts`**: Publishes subaccount order updates for WebSocket clients.
  - Message format: `SubaccountMessage` protobuf.
  - Contains order status changes, placements, and cancellations.
- **`to-websockets-orderbooks`**: Publishes orderbook level updates for WebSocket clients.
  - Message format: `OrderbookMessage` protobuf.
  - Contains price level changes (bids/asks).

### Message Handlers

- **OrderPlaceHandler** ([`src/handlers/order-place-handler.ts`](src/handlers/order-place-handler.ts)): Processes order placements and replacements.
- **OrderUpdateHandler** ([`src/handlers/order-update-handler.ts`](src/handlers/order-update-handler.ts)): Processes order fill updates.
- **OrderRemoveHandler** ([`src/handlers/order-remove-handler.ts`](src/handlers/order-remove-handler.ts)): Processes order cancellations and removals.

## Configuration

Configuration is loaded from environment variables and parsed in [`src/config.ts`](src/config.ts). Environment-specific defaults are in `.env`, `.env.test`.

### Core Settings

- **`SERVICE_NAME`** (string, default: `vulcan`): Service identifier for logging and metrics.
- **`NODE_ENV`** (string, default: `development`): Environment mode (`development`, `test`, `production`).

### Kafka Configuration

- **`KAFKA_BROKER_URLS`** (string, required): Comma-separated list of Kafka broker URLs.
- **`KAFKA_ENABLE_UNIQUE_CONSUMER_GROUP_IDS`** (boolean, default: `false`): Enable unique consumer group IDs per instance.
- **`BATCH_PROCESSING_ENABLED`** (boolean, default: `true`): Enable batch processing of Kafka messages.
- **`PROCESS_FROM_BEGINNING`** (boolean, default: `false`): Start consuming from the beginning of the topic (only matters if offset is lost).
- **`KAFKA_BATCH_PROCESSING_COMMIT_FREQUENCY_MS`** (number, default: `3000`): Frequency of offset commits during batch processing (3 seconds).

### Redis Configuration

- **`REDIS_URL`** (string, default: `redis://localhost:6382`): Primary Redis URL for caching.
- **`REDIS_RECONNECT_TIMEOUT_MS`** (integer, default: `5000`): Redis reconnection timeout.

### Message Batching Configuration

- **`FLUSH_KAFKA_MESSAGES_INTERVAL_MS`** (number, default: `10`): Interval for flushing batched WebSocket messages to Kafka (10ms).
- **`MAX_WEBSOCKET_MESSAGES_TO_QUEUE_PER_TOPIC`** (number, default: `20`): Maximum messages to queue per topic before flushing.

### Feature Flags

- **`SEND_WEBSOCKET_MESSAGES`** (boolean, default: `true`): Enable sending WebSocket messages to Kafka. Set to `false` during fast sync.
- **`SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_STATEFUL_ORDERS`** (boolean, default: `true`): Send subaccount messages for stateful order placements.
- **`SEND_SUBACCOUNT_WEBSOCKET_MESSAGE_FOR_CANCELS_MISSING_ORDERS`** (boolean, default: `true`): Send subaccount messages for cancellations of orders not found in Redis.

## Running Locally

### Prerequisites

- Node.js 18+ and `pnpm` installed.
- Kafka broker running and accessible.
- Redis instance running and accessible.
- PostgreSQL database running (for perpetual market and block height caches).
- Environment variables configured (see `.env` for examples).

### Steps

1. **Install dependencies** (from monorepo root):
   ```bash
   pnpm install
   ```

2. **Build the service**:
   ```bash
   cd indexer/services/vulcan
   pnpm run build
   ```

3. **Set up environment variables**:
   - Copy `.env` or create a `.env` file with required variables.
   - Ensure `KAFKA_BROKER_URLS`, `REDIS_URL`, `DB_HOSTNAME`, `DB_NAME`, `DB_USERNAME`, and `DB_PASSWORD` are set.

4. **Run the service**:
   ```bash
   pnpm run start
   ```
   The service will start consuming messages from the `to-vulcan` Kafka topic.

### Development Commands

- **`pnpm run build`**: Compile TypeScript to JavaScript in `build/` directory.
- **`pnpm run build:watch`**: Compile in watch mode.
- **`pnpm run start`**: Run the service in production mode with DataDog tracing.
- **`pnpm run lint`**: Run ESLint on the codebase.
- **`pnpm run lint:fix`**: Run ESLint with auto-fix.
