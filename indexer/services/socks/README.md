# Socks

Socks is the WebSocket server for the dYdX v4 Indexer. It provides real-time streaming of market data, account updates, trades, orderbook changes, and block height information to connected clients. The service consumes messages from Kafka topics and forwards them to subscribed WebSocket clients based on their channel subscriptions.

## Responsibilities & Scope

- Accept and manage WebSocket connections from clients.
- Handle subscription requests for various data channels (markets, orderbooks, trades, accounts, candles, block height).
- Consume messages from Kafka topics and forward them to subscribed clients.
- Support both individual message delivery and batched message delivery per subscription.
- Enforce per-connection subscription limits for each channel.
- Rate limit subscription requests and invalid messages to prevent abuse.
- Maintain heartbeat/ping-pong mechanism to detect and close stale connections.
- Provide initial snapshot data when clients subscribe to a channel.

**Out of scope:**
- Does not write to Kafka or modify any data.
- Does not perform authentication or authorization (relies on geo-blocking headers from upstream).
- Does not persist WebSocket connection state (stateless, connections are ephemeral).

## Architecture & Dependencies

### Internal Structure

- **WebSocket Server** (`src/helpers/wss.ts`, `src/websocket/index.ts`): Manages WebSocket connections, handles connection lifecycle, and implements heartbeat mechanism.
- **Subscriptions** (`src/lib/subscription.ts`): Tracks client subscriptions by channel and ID, validates subscription requests, fetches initial data from Comlink API.
- **Message Forwarder** (`src/lib/message-forwarder.ts`): Consumes Kafka messages, routes them to subscribed clients, and handles message batching.
- **Rate Limiters** (`src/lib/rate-limit.ts`, `src/lib/invalid-message.ts`): Enforces rate limits on subscriptions and invalid messages per connection.
- **Kafka Helpers** (`src/helpers/from-kafka-helpers.ts`, `src/helpers/kafka/kafka-controller.ts`): Decodes Kafka messages and maps topics to channels.

### Internal Dependencies

- `@dydxprotocol-indexer/base`: Shared utilities, logging, stats, and error handling.
- `@dydxprotocol-indexer/postgres`: Database models, perpetual market cache, and block height cache.
- `@dydxprotocol-indexer/kafka`: Kafka consumer and producer utilities.
- `@dydxprotocol-indexer/v4-protos`: Protocol buffer definitions for messages.
- `@dydxprotocol-indexer/compliance`: Geo-blocking header types.

### External Dependencies

- **Kafka**: Message bus for receiving real-time updates from the indexer pipeline.
- **Comlink API**: REST API for fetching initial snapshot data when clients subscribe.
- **PostgreSQL**: Indirectly via perpetual market and block height caches.
- **Redis**: Not directly used by this service.

### Message Flow

1. Client connects via WebSocket and receives a `connected` message with a unique connection ID.
2. Client sends a `subscribe` message specifying a channel and optional ID (e.g., `v4_trades` for `BTC-USD`).
3. Socks validates the subscription, fetches initial data from Comlink, and sends a `subscribed` message with the snapshot.
4. Socks consumes messages from Kafka topics (e.g., `to-websockets-trades`).
5. Message Forwarder decodes Kafka messages, maps them to channels, and forwards to subscribed clients.
6. For batched subscriptions, messages are buffered and sent in batches at regular intervals.
7. Client can send `unsubscribe` messages to stop receiving updates for a channel.
8. Heartbeat mechanism sends periodic pings; connections are closed if pongs are not received.

## Public Interface

### WebSocket Endpoints

- **`ws://host:port/`**: Main WebSocket endpoint for all client connections.

### WebSocket Message Types

#### Incoming Messages (Client → Server)

- **`subscribe`**: Subscribe to a channel.
  - Fields: `type`, `channel`, `id` (optional for markets/block height), `batched` (optional).
  - Handler: [`src/websocket/index.ts`](src/websocket/index.ts) → [`src/lib/subscription.ts`](src/lib/subscription.ts)
- **`unsubscribe`**: Unsubscribe from a channel.
  - Fields: `type`, `channel`, `id` (optional for markets/block height).
  - Handler: [`src/websocket/index.ts`](src/websocket/index.ts) → [`src/lib/subscription.ts`](src/lib/subscription.ts)
- **`ping`**: Client-initiated ping (server automatically responds to WebSocket pings).
  - Fields: `type`, `id` (optional).

#### Outgoing Messages (Server → Client)

- **`connected`**: Sent immediately after connection is established.
  - Fields: `type`, `connection_id`, `message_id`.
- **`subscribed`**: Sent after successful subscription with initial data snapshot.
  - Fields: `type`, `connection_id`, `message_id`, `channel`, `id`, `contents`.
- **`unsubscribed`**: Sent after successful unsubscription.
  - Fields: `type`, `connection_id`, `message_id`, `channel`, `id`.
- **`channel_data`**: Real-time update for a subscribed channel (non-batched).
  - Fields: `type`, `connection_id`, `message_id`, `channel`, `id`, `version`, `contents`, `subaccountNumber` (for parent accounts).
- **`channel_batch_data`**: Batched real-time updates for a subscribed channel.
  - Fields: `type`, `connection_id`, `message_id`, `channel`, `id`, `version`, `contents` (array), `subaccountNumber` (for parent accounts).
- **`error`**: Error message for invalid requests or internal errors.
  - Fields: `type`, `connection_id`, `message_id`, `message`, `channel`, `id`.
- **`pong`**: Response to client ping (automatic WebSocket pongs are also sent).
  - Fields: `type`, `connection_id`, `message_id`, `id`.

### Supported Channels

- **`v4_markets`**: All perpetual markets data (no ID required).
- **`v4_orderbook`**: Orderbook levels for a specific market (ID: ticker, e.g., `BTC-USD`).
- **`v4_trades`**: Recent trades for a specific market (ID: ticker).
- **`v4_subaccounts`**: Account data for a specific subaccount (ID: `address/subaccountNumber`).
- **`v4_parent_subaccounts`**: Aggregated account data for a parent subaccount and its children (ID: `address/parentSubaccountNumber`).
- **`v4_candles`**: OHLCV candles for a specific market and resolution (ID: `ticker/resolution`, e.g., `BTC-USD/1MIN`).
- **`v4_block_height`**: Latest block height and timestamp (no ID required).

### HTTP Endpoints

- **`GET /health`**: Health check endpoint returning `{"ok": true}`.
  - Handler: [`src/server.ts`](src/server.ts)

## Configuration

Configuration is loaded from environment variables and parsed in [`src/config.ts`](src/config.ts). Environment-specific defaults are in `.env.development`, `.env.test`, and `.env.production`.

### Core Settings

- **`SERVICE_NAME`** (string, default: `socks`): Service identifier for logging and metrics.
- **`NODE_ENV`** (string, default: `development`): Environment mode (`development`, `test`, `production`).
- **`PORT`** (integer, default: `8000`): HTTP server port for health check endpoint.
- **`WS_PORT`** (integer, default: `8080`): WebSocket server port.
- **`CORS_ORIGIN`** (string, optional): CORS allowed origins for HTTP endpoints.

### WebSocket Configuration

- **`WS_HEARTBEAT_INTERVAL_MS`** (integer, default: `30000`): Interval between heartbeat pings (30 seconds).
- **`WS_HEARTBEAT_TIMEOUT_MS`** (integer, default: `10000`): Timeout for receiving pong after ping (10 seconds).
- **`BATCH_SEND_INTERVAL_MS`** (integer, default: `250`): Interval for sending batched messages (250ms).

### Rate Limiting

- **`RATE_LIMIT_ENABLED`** (boolean, default: `true`): Enable rate limiting.
- **`RATE_LIMIT_SUBSCRIBE_POINTS`** (number, default: `2`): Points consumed per subscription request.
- **`RATE_LIMIT_SUBSCRIBE_DURATION_MS`** (integer, default: `1000`): Rate limit window duration for subscriptions.
- **`RATE_LIMIT_PING_POINTS`** (number, default: `5`): Points consumed per ping request.
- **`RATE_LIMIT_PING_DURATION_MS`** (integer, default: `1000`): Rate limit window duration for pings.
- **`RATE_LIMIT_INVALID_MESSAGE_POINTS`** (number, default: `2`): Points consumed per invalid message.
- **`RATE_LIMIT_INVALID_MESSAGE_DURATION_MS`** (integer, default: `1000`): Rate limit window duration for invalid messages.

### Per-Channel Subscription Limits

- **`V4_ACCOUNTS_CHANNEL_LIMIT`** (integer, default: `256`): Maximum subscriptions per connection for `v4_subaccounts`.
- **`V4_CANDLES_CHANNEL_LIMIT`** (integer, default: `32`): Maximum subscriptions per connection for `v4_candles`.
- **`V4_MARKETS_CHANNEL_LIMIT`** (integer, default: `32`): Maximum subscriptions per connection for `v4_markets`.
- **`V4_ORDERBOOK_CHANNEL_LIMIT`** (integer, default: `32`): Maximum subscriptions per connection for `v4_orderbook`.
- **`V4_PARENT_ACCOUNTS_CHANNEL_LIMIT`** (integer, default: `256`): Maximum subscriptions per connection for `v4_parent_subaccounts`.
- **`V4_TRADES_CHANNEL_LIMIT`** (integer, default: `32`): Maximum subscriptions per connection for `v4_trades`.

### Database Configuration

- **`DB_HOSTNAME`** (string, required): PostgreSQL primary host.
- **`DB_READONLY_HOSTNAME`** (string, optional): PostgreSQL read replica host.
- **`DB_PORT`** (integer, default: `5432`): PostgreSQL port.
- **`DB_NAME`** (string, required): Database name.
- **`DB_USERNAME`** (string, required): Database username.
- **`DB_PASSWORD`** (string, required): Database password.

### Kafka Configuration

- **`KAFKA_BROKER_URLS`** (string, required): Comma-separated list of Kafka broker URLs.
- **`KAFKA_ENABLE_UNIQUE_CONSUMER_GROUP_IDS`** (boolean, default: `false`): Enable unique consumer group IDs per instance.
- **`BATCH_PROCESSING_ENABLED`** (boolean, default: `true`): Enable batch processing of Kafka messages.
- **`KAFKA_BATCH_PROCESSING_COMMIT_FREQUENCY_MS`** (number, default: `3000`): Frequency of offset commits during batch processing.

### Comlink API Configuration

- **`COMLINK_URL`** (string, required): Base URL for Comlink API (e.g., `localhost:3002`).
- **`AXIOS_TIMEOUT_MS`** (integer, default: `5000`): Timeout for Comlink API requests.
- **`INITIAL_GET_TIMEOUT_MS`** (integer, default: `20000`): Timeout for fetching initial subscription data.

### Compliance Configuration

- **`INDEXER_LEVEL_GEOBLOCKING_ENABLED`** (boolean, default: `false` in dev, `true` in production): Enable geo-blocking checks.

### Metrics Configuration

- **`MESSAGE_FORWARDER_STATSD_SAMPLE_RATE`** (number, default: `1.0`): Sample rate for message forwarding metrics.
- **`ENABLE_ORDERBOOK_LOGS`** (boolean, default: `true`): Enable detailed orderbook logging.
- **`SUBSCRIPTION_METRIC_INTERVAL_MS`** (integer, default: `60000`): Interval for emitting subscription metrics (1 minute).
- **`PERPETUAL_MARKETS_REFRESHER_INTERVAL_MS`** (integer, default: `300000`): Interval for refreshing perpetual markets cache (5 minutes).

## Running Locally

### Prerequisites

- Node.js 16+ and `pnpm` 8+ installed.
- PostgreSQL database running and accessible.
- Kafka broker running and accessible.
- Comlink service running (for initial subscription data).
- Environment variables configured (see `.env.development` for examples).

### Steps

1. **Install dependencies** (from monorepo root):
   ```bash
   pnpm install
   ```

2. **Build the service**:
   ```bash
   cd indexer/services/socks
   pnpm run build
   ```

3. **Set up environment variables**:
   - Copy `.env.development` or create a `.env` file with required variables.
   - Ensure `DB_HOSTNAME`, `DB_NAME`, `DB_USERNAME`, `DB_PASSWORD`, `KAFKA_BROKER_URLS`, and `COMLINK_URL` are set.

4. **Start dependencies** (if using Docker Compose):
   ```bash
   docker compose up -d postgres kafka
   ```

5. **Run the service**:
   ```bash
   pnpm run start
   ```
   The WebSocket server will start on `ws://localhost:8080` (or the port specified in `WS_PORT`).
   The HTTP health check server will start on `http://localhost:8000` (or the port specified in `PORT`).

6. **Test the connection**:
   - Use a WebSocket client (e.g., `wscat`, browser console) to connect to `ws://localhost:8080`.
   - Send a subscribe message: `{"type":"subscribe","channel":"v4_markets"}`.

### Development Commands

- **`pnpm run build`**: Compile TypeScript to JavaScript in `build/` directory.
- **`pnpm run build:watch`**: Compile in watch mode.
- **`pnpm run start`**: Run the service in production mode with DataDog tracing.
- **`pnpm run lint`**: Run ESLint on the codebase.
- **`pnpm run lint:fix`**: Run ESLint with auto-fix.

### Logging
- Logs are structured JSON written to stdout via `@dydxprotocol-indexer/base` logger.

- Key log fields:

  - `at`: Source location (e.g., `index#onConnection`).

  - `message`: Human-readable message.

  - `error`: Error object with message and stack trace (if applicable).

  - `connectionId`: Unique identifier for WebSocket connections.

  - `url`, `protocol`, `headers`: Connection metadata.

  - `numConcurrentConnections`: Current number of active connections.

### Metrics
- Metrics are emitted via `@dydxprotocol-indexer/base` stats module (StatsD format).

- Key metrics:

  - `socks.num_connections`: Counter for new connections.

  - `socks.num_concurrent_connections`: Gauge for active connections.

  - `socks.num_disconnects`: Counter for disconnections (tagged by close code and reason).

  - `socks.on_message`: Counter for messages received from clients.

  - `socks.message_received_<type>`: Counter for messages by type (subscribe, unsubscribe, ping).

  - `socks.message_to_forward`: Counter for messages forwarded to clients.

  - `socks.forward_message_with_subscribers`: Counter for messages with active subscribers.

  - `socks.forward_to_client_success`: Counter for successful message deliveries.

  - `socks.forward_to_client_error`: Counter for failed message deliveries.

  - `socks.forward_to_client_batch_success`: Counter for successful batched message deliveries.

  - `socks.ws_send.error`: Counter for WebSocket send errors (tagged by error code).

  - `socks.ws_send.stream_destroyed_errors`: Counter for stream destroyed errors.

  - `socks.ws_send.write_epipe_errors`: Counter for EPIPE write errors.

  - `socks.ws_message_not_sent`: Counter for messages not sent (tagged by reason and readyState).

  - `socks.message_time_in_queue`: Timing for Kafka message queue time.

  - `socks.message_time_since_received`: Timing from message receipt to forwarding.

  - `socks.forward_message`: Timing for message forwarding.

  - `socks.subscribe_send_message`: Timing for sending subscription messages.

  - `socks.subscriptions_limit_reached`: Counter for subscription limit violations (tagged by channel).

  - `socks.initial_response_error`: Counter for initial data fetch errors (tagged by channel).

  - `socks.initial_response_get`: Timing for initial data fetch (tagged by channel).

  - `socks.largest_subscriber`: Gauge for largest number of subscriptions per connection (tagged by channel).

  - `socks.subscriptions.channel_size`: Gauge for total subscriptions per channel.

  - `socks.batch_time_in_queue`: Timing for Kafka batch queue time (batch processing mode).

  - `socks.batch_processing_time`: Timing for Kafka batch processing (batch processing mode).

  - `socks.batch_size`: Gauge for Kafka batch size (batch processing mode).

## Directory Layout

- **`src/`**: Source code.

  - **`websocket/index.ts`**: Main WebSocket server logic, connection handling, and message routing.

  - **`lib/`**: Core business logic.

    - **`subscription.ts`**: Subscription management, validation, and initial data fetching.

    - **`message-forwarder.ts`**: Kafka message consumption and forwarding to WebSocket clients.

    - **`rate-limit.ts`**: General-purpose rate limiter.

    - **`invalid-message.ts`**: Invalid message handler with rate limiting.

    - **`constants.ts`**: WebSocket close codes, error messages, and topic-to-channel mappings.

    - **`errors.ts`**: Custom error classes.

    - **`axios.ts`**: Axios wrapper for Comlink API requests.

  - **`helpers/`**: Utility functions.

    - **`wss.ts`**: WebSocket server wrapper and message sending utilities.

    - **`message.ts`**: Message creation helpers for outgoing WebSocket messages.

    - **`from-kafka-helpers.ts`**: Kafka message decoding and channel mapping.

    - **`header-utils.ts`**: Geo-blocking header extraction.

    - **`kafka/kafka-controller.ts`**: Kafka consumer initialization and connection management.

  - **`middlewares/`**: Express middleware.

    - **`request-logger.ts`**: HTTP request/response logging.

    - **`res-body-capture.ts`**: Response body capture for logging.

  - **`types.ts`**: TypeScript type definitions for messages and internal structures.

  - **`config.ts`**: Configuration schema and environment variable parsing.

  - **`server.ts`**: Express HTTP server for health check endpoint.

  - **`index.ts`**: Service entrypoint.

- **`__tests__/`**: Test files (unit and integration tests).

  - **`scripts/`**: Test utility scripts.

    - **`dydx_ws_limit_test.py`**: Python script for testing subscription limits.

- **`build/`**: Compiled JavaScript output (generated by `pnpm run build`).

- **`patches/`**: Patch files for npm dependencies.

  - **`kafkajs+2.2.4.patch`**: Performance optimizations for KafkaJS (reduces memory allocations).

- **`.env`**: Environment variables for local development (not committed, create from `.env.development`).

- **`.env.development`**: Default environment variables for development.

- **`.env.test`**: Environment variables for testing.

- **`.env.production`**: Production-specific environment variables.

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

- **WebSocket API Documentation**: Refer to the Comlink service's `public/websocket-documentation.md` for client-facing WebSocket API documentation.
