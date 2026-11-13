# Comlink

Comlink is the public-facing REST API service for the dYdX v4 Indexer. It serves all read-only queries for market data, account information, orders, fills, transfers, and compliance checks. The service aggregates data from PostgreSQL (historical state) and Redis (real-time orderbook and order state) to provide a unified view of the exchange.

## Responsibilities & Scope

- Serve REST API endpoints for querying market data, perpetual positions, orders, fills, transfers, PnL, and funding payments.
- Aggregate data from PostgreSQL (historical) and Redis (real-time) to provide consistent responses.
- Enforce rate limiting per IP address and per endpoint.
- Perform compliance and geo-blocking checks for restricted addresses and regions.
- Support Turnkey-based wallet operations including account creation, authentication, and bridging.
- Facilitate cross-chain bridging via Skip Protocol and Alchemy webhooks.
- Track affiliate relationships and trading rewards.
- Expose vault positions and historical PnL for megavault and individual vaults.
- Send push notifications via Firebase for order fills and other events.

**Out of scope:**
- Does not write to the blockchain or submit transactions (read-only API).
- Does not run the WebSocket server (WebSocket functionality is provided by a separate service).
- Does not perform order matching or settlement.

## Architecture & Dependencies

### Internal Structure

- **Controllers** (`src/controllers/api/v4/`): Express route handlers organized by domain (addresses, orders, fills, markets, compliance, turnkey, bridging, vaults, etc.).
- **Request Helpers** (`src/request-helpers/`): Middleware for logging, validation, error handling, response transformation, and cache control.
- **Helpers** (`src/helpers/`): Business logic for compliance checks, Turnkey integration, Alchemy webhooks, Skip Protocol bridging, policy enforcement, and email validation.
- **Caches** (`src/caches/`): Rate limiter configurations and vault start PnL cache.
- **Types** (`src/types.ts`): TypeScript interfaces for requests and responses.
- **Config** (`src/config.ts`): Centralized configuration schema with environment variable parsing.

### Internal Dependencies

- `@dydxprotocol-indexer/base`: Shared utilities, logging, stats, and error handling.
- `@dydxprotocol-indexer/postgres`: Database models and query helpers.
- `@dydxprotocol-indexer/redis`: Redis clients and caching utilities.
- `@dydxprotocol-indexer/compliance`: Compliance provider integrations (Elliptic).
- `@dydxprotocol-indexer/notifications`: Firebase notification support.
- `@dydxprotocol-indexer/v4-protos`: Protocol buffer definitions for orders and events.

### External Dependencies

- **PostgreSQL**: Primary datastore for historical orders, fills, positions, transfers, and market data. Supports read replicas.
- **Redis**: Real-time orderbook levels, open orders, and caching layer. Supports read-only replicas.
- **Elliptic API**: Compliance screening for wallet addresses.
- **Turnkey API**: Wallet creation, authentication, and transaction signing for custodial wallet features.
- **Alchemy API**: Webhook management for EVM and Solana chain events to trigger bridging.
- **Skip Protocol**: Cross-chain bridging and token swaps.
- **ZeroDev**: EIP-7702 and smart account sponsorship for gasless transactions.
- **Firebase**: Push notifications for order fills and other events.
- **Amplitude**: Analytics tracking for user events (optional).

### Request Flow

1. Client sends HTTP request to Express server.
2. Middleware applies rate limiting (via Redis), geo-blocking (via headers), and compliance checks (via Elliptic API or cached data).
3. Controller validates request parameters using express-validator schemas.
4. Controller queries PostgreSQL and/or Redis for data.
5. Response transformers convert database models to API response objects.
6. Cache-Control headers are set based on endpoint configuration.
7. Response is returned to client with appropriate status code and body.
8. Request/response details are logged and metrics are emitted.

## Public Interface

### REST API Endpoints

All endpoints are prefixed with `/v4/`. The service exposes approximately 50+ endpoints organized by domain.

#### Core Account & Position Endpoints

- **`GET /v4/addresses/:address`**: Retrieve all subaccounts and total trading rewards for an address.
  - Handler: [`src/controllers/api/v4/addresses-controller.ts`](src/controllers/api/v4/addresses-controller.ts)
- **`GET /v4/addresses/:address/subaccountNumber/:subaccountNumber`**: Retrieve a specific subaccount with positions and balances.
- **`GET /v4/addresses/:address/parentSubaccountNumber/:parentSubaccountNumber`**: Retrieve aggregated data for a parent subaccount and its children.
- **`POST /v4/addresses/:address/registerToken`**: Register a Firebase notification token for an address (requires signature verification).
- **`GET /v4/perpetualPositions`**: List perpetual positions for a subaccount with optional filters.
  - Handler: [`src/controllers/api/v4/perpetual-positions-controller.ts`](src/controllers/api/v4/perpetual-positions-controller.ts)
- **`GET /v4/perpetualPositions/parentSubaccountNumber`**: List perpetual positions for a parent subaccount.
- **`GET /v4/assetPositions`**: List asset positions (USDC) for a subaccount.
  - Handler: [`src/controllers/api/v4/asset-positions-controller.ts`](src/controllers/api/v4/asset-positions-controller.ts)
- **`GET /v4/assetPositions/parentSubaccountNumber`**: List asset positions for a parent subaccount.

#### Order & Fill Endpoints

- **`GET /v4/orders`**: List orders for a subaccount with optional filters (status, side, type, ticker, goodTilBlock, etc.).
  - Handler: [`src/controllers/api/v4/orders-controller.ts`](src/controllers/api/v4/orders-controller.ts)
- **`GET /v4/orders/parentSubaccountNumber`**: List orders for a parent subaccount.
- **`GET /v4/orders/:orderId`**: Retrieve a specific order by ID.
- **`GET /v4/fills`**: List fills for a subaccount with optional filters (market, side, liquidity, type).
  - Handler: [`src/controllers/api/v4/fills-controller.ts`](src/controllers/api/v4/fills-controller.ts)
- **`GET /v4/fills/parentSubaccountNumber`**: List fills for a parent subaccount.

#### Transfer & PnL Endpoints

- **`GET /v4/transfers`**: List transfers for a subaccount.
  - Handler: [`src/controllers/api/v4/transfers-controller.ts`](src/controllers/api/v4/transfers-controller.ts)
- **`GET /v4/transfers/parentSubaccountNumber`**: List transfers for a parent subaccount.
- **`GET /v4/transfers/between`**: List transfers between two specific subaccounts.
- **`GET /v4/historical-pnl`**: Retrieve historical PnL ticks for a subaccount.
  - Handler: [`src/controllers/api/v4/historical-pnl-controller.ts`](src/controllers/api/v4/historical-pnl-controller.ts)
- **`GET /v4/historical-pnl/parentSubaccountNumber`**: Retrieve historical PnL for a parent subaccount.
- **`GET /v4/pnl`**: Retrieve aggregated PnL (hourly or daily) for a subaccount.
  - Handler: [`src/controllers/api/v4/pnl-controller.ts`](src/controllers/api/v4/pnl-controller.ts)
- **`GET /v4/pnl/parentSubaccountNumber`**: Retrieve aggregated PnL for a parent subaccount.

#### Market Data Endpoints

- **`GET /v4/perpetualMarkets`**: List all perpetual markets with current prices and funding rates.
  - Handler: [`src/controllers/api/v4/perpetual-markets-controller.ts`](src/controllers/api/v4/perpetual-markets-controller.ts)
- **`GET /v4/orderbooks/perpetualMarket/:ticker`**: Retrieve orderbook levels for a specific market.
  - Handler: [`src/controllers/api/v4/orderbook-controller.ts`](src/controllers/api/v4/orderbook-controller.ts)
- **`GET /v4/trades/perpetualMarket/:ticker`**: List recent trades for a specific market.
  - Handler: [`src/controllers/api/v4/trades-controller.ts`](src/controllers/api/v4/trades-controller.ts)
- **`GET /v4/candles/perpetualMarkets/:ticker`**: Retrieve OHLCV candles for a market.
  - Handler: [`src/controllers/api/v4/candles-controller.ts`](src/controllers/api/v4/candles-controller.ts)
- **`GET /v4/sparklines`**: Retrieve sparkline data (recent price points) for all markets.
  - Handler: [`src/controllers/api/v4/sparklines-controller.ts`](src/controllers/api/v4/sparklines-controller.ts)
- **`GET /v4/historicalFunding/:ticker`**: Retrieve historical funding rates for a market.
  - Handler: [`src/controllers/api/v4/historical-funding-controller.ts`](src/controllers/api/v4/historical-funding-controller.ts)
- **`GET /v4/fundingPayments`**: List funding payments for a subaccount.
  - Handler: [`src/controllers/api/v4/funding-payments-controller.ts`](src/controllers/api/v4/funding-payments-controller.ts)
- **`GET /v4/fundingPayments/parentSubaccount`**: List funding payments for a parent subaccount.

#### Compliance Endpoints

- **`GET /v4/screen/:address`**: Screen an address for compliance (deprecated, use `/v4/compliance/screen/:address`).
  - Handler: [`src/controllers/api/v4/compliance-controller.ts`](src/controllers/api/v4/compliance-controller.ts)
- **`GET /v4/compliance/screen/:address`**: Screen an address for compliance and geo-blocking.
  - Handler: [`src/controllers/api/v4/compliance-v2-controller.ts`](src/controllers/api/v4/compliance-v2-controller.ts)
- **`POST /v4/compliance/geoblock`**: Update compliance status based on geo-blocking survey (requires signature).
- **`POST /v4/compliance/geoblock-keplr`**: Update compliance status using Keplr wallet signature.
- **`POST /v4/compliance/setStatus`**: Manually set compliance status (dev/staging only, requires `EXPOSE_SET_COMPLIANCE_ENDPOINT=true`).

#### Affiliate Endpoints

- **`GET /v4/affiliates/metadata`**: Retrieve affiliate metadata (referral code, eligibility).
  - Handler: [`src/controllers/api/v4/affiliates-controller.ts`](src/controllers/api/v4/affiliates-controller.ts)
- **`GET /v4/affiliates/address`**: Retrieve address for a given referral code.
- **`POST /v4/affiliates/referralCode`**: Create or update a referral code for an address (requires signature).
- **`POST /v4/affiliates/referralCode-keplr`**: Create or update a referral code using Keplr wallet signature.
- **`GET /v4/affiliates/snapshot`**: Retrieve affiliate earnings and referral statistics.
- **`GET /v4/affiliates/total_volume`**: Retrieve total trading volume for an address.

#### Vault Endpoints

- **`GET /v4/vault/v1/megavault/historicalPnl`**: Retrieve historical PnL for the megavault (hourly or daily).
  - Handler: [`src/controllers/api/v4/vault-controller.ts`](src/controllers/api/v4/vault-controller.ts)
- **`GET /v4/vault/v1/vaults/historicalPnl`**: Retrieve historical PnL for all individual vaults.
- **`GET /v4/vault/v1/megavault/positions`**: Retrieve current positions held by the megavault.

#### Turnkey Wallet Endpoints

- **`POST /v4/turnkey/signin`**: Authenticate or create a Turnkey wallet (email, social OAuth, or passkey).
  - Handler: [`src/controllers/api/v4/turnkey-controller.ts`](src/controllers/api/v4/turnkey-controller.ts)
- **`POST /v4/turnkey/uploadAddress`**: Upload dYdX address and configure bridging policies (requires signature).
- **`GET /v4/turnkey/appleLoginRedirect`**: Handle Apple OAuth redirect and return session data.

#### Bridging Endpoints

- **`POST /v4/bridging/startBridge`**: Webhook endpoint to initiate cross-chain bridging via Alchemy events (requires Alchemy auth token).
  - Handler: [`src/controllers/api/v4/skip-bridge-controller.ts`](src/controllers/api/v4/skip-bridge-controller.ts)
- **`GET /v4/bridging/getDepositAddress/:dydxAddress`**: Retrieve deposit addresses (EVM, SVM, Avalanche) for a dYdX address.
- **`GET /v4/bridging/getDeposits/:dydxAddress`**: Retrieve deposit history for a dYdX address.

#### Utility Endpoints

- **`GET /v4/time`**: Retrieve current server time in ISO and epoch formats.
  - Handler: [`src/controllers/api/v4/time-controller.ts`](src/controllers/api/v4/time-controller.ts)
- **`GET /v4/height`**: Retrieve latest indexed block height and time.
  - Handler: [`src/controllers/api/v4/height-controller.ts`](src/controllers/api/v4/height-controller.ts)
- **`GET /v4/trader/search`**: Search for a trader by address or username.
  - Handler: [`src/controllers/api/v4/social-trading-controller.ts`](src/controllers/api/v4/social-trading-controller.ts)
- **`GET /v4/historicalTradingRewardAggregations/:address`**: Retrieve historical trading reward aggregations.
  - Handler: [`src/controllers/api/v4/historical-trading-reward-aggregations-controller.ts`](src/controllers/api/v4/historical-trading-reward-aggregations-controller.ts)
- **`GET /v4/historicalBlockTradingRewards/:address`**: Retrieve block-level trading rewards.
  - Handler: [`src/controllers/api/v4/historical-block-trading-rewards-controller.ts`](src/controllers/api/v4/historical-block-trading-rewards-controller.ts)

### Documentation Endpoints

- **`GET /docs`**: Swagger UI for interactive API documentation.
- **`GET /health`**: Health check endpoint returning `{"ok": true}`.

### API Documentation

- OpenAPI specification: [`public/swagger.json`](public/swagger.json)
- Markdown documentation: [`public/api-documentation.md`](public/api-documentation.md)
- WebSocket documentation: [`public/websocket-documentation.md`](public/websocket-documentation.md)

## Configuration

Configuration is loaded from environment variables and parsed in [`src/config.ts`](src/config.ts). Environment-specific defaults are in `.env.development`, `.env.test`, and `.env.production`.

### Core Settings

- **`SERVICE_NAME`** (string, default: `comlink`): Service identifier for logging and metrics.
- **`NODE_ENV`** (string, default: `development`): Environment mode (`development`, `test`, `production`).
- **`PORT`** (integer, default: `8080`): HTTP server port.
- **`CORS_ORIGIN`** (string, default: `*`): CORS allowed origins.
- **`LOG_GETS`** (boolean, default: `false`): Whether to log GET requests.
- **`KEEP_ALIVE_MS`** (integer, default: `61000`): Server keep-alive timeout.
- **`HEADERS_TIMEOUT_MS`** (integer, default: `65000`): Server headers timeout.

### Database Configuration

- **`DB_HOSTNAME`** (string, required): PostgreSQL primary host.
- **`DB_READONLY_HOSTNAME`** (string, optional): PostgreSQL read replica host.
- **`DB_PORT`** (integer, default: `5432`): PostgreSQL port.
- **`DB_NAME`** (string, required): Database name.
- **`DB_USERNAME`** (string, required): Database username.
- **`DB_PASSWORD`** (string, required): Database password.

### Redis Configuration

- **`REDIS_URL`** (string, default: `redis://localhost:6382`): Primary Redis URL for caching and orderbook data.
- **`REDIS_READONLY_URL`** (string, optional): Read-only Redis URL.
- **`REDIS_RECONNECT_TIMEOUT_MS`** (integer, default: `5000`): Redis reconnection timeout.
- **`RATE_LIMIT_REDIS_URL`** (string, default: `redis://localhost:6382`): Redis URL for rate limiting.

### Rate Limiting Configuration

- **`RATE_LIMIT_ENABLED`** (boolean, default: `true`): Enable rate limiting.
- **`RATE_LIMIT_GET_POINTS`** (integer, default: `100`): Points per request for general GET endpoints.
- **`RATE_LIMIT_GET_DURATION_SECONDS`** (integer, default: `10`): Rate limit window duration.
- **`INDEXER_INTERNAL_IPS`** (string, default: `""`): Comma-separated list of internal IPs exempt from rate limiting.

Per-endpoint rate limits (all follow the pattern `RATE_LIMIT_<ENDPOINT>_POINTS` and `RATE_LIMIT_<ENDPOINT>_DURATION_SECONDS`):
- `RATE_LIMIT_ORDERS_*`: Orders endpoint limits.
- `RATE_LIMIT_FILLS_*`: Fills endpoint limits.
- `RATE_LIMIT_CANDLES_*`: Candles endpoint limits.
- `RATE_LIMIT_SPARKLINES_*`: Sparklines endpoint limits.
- `RATE_LIMIT_HISTORICAL_PNL_*`: Historical PnL endpoint limits.
- `RATE_LIMIT_PNL_*`: PnL endpoint limits.
- `RATE_LIMIT_FUNDING_*`: Funding payments endpoint limits.
- `RATE_LIMIT_SCREEN_QUERY_PROVIDER_*`: Compliance screening limits per IP.
- `RATE_LIMIT_SCREEN_QUERY_PROVIDER_GLOBAL_*`: Global compliance screening limits.

### Compliance Configuration

- **`INDEXER_LEVEL_GEOBLOCKING_ENABLED`** (boolean, default: `false` in dev, `true` in production): Enable geo-blocking checks.
- **`ELLIPTIC_API_KEY`** (string, optional): Elliptic API key for compliance screening.
- **`MAX_AGE_SCREENED_ADDRESS_COMPLIANCE_DATA_SECONDS`** (integer, default: `86400`): Cache TTL for compliance data (1 day).
- **`EXPOSE_SET_COMPLIANCE_ENDPOINT`** (boolean, default: `false`): Expose internal compliance status endpoint (dev/staging only).

### Turnkey & Bridging Configuration

- **`TURNKEY_API_BASE_URL`** (string, default: `https://api.turnkey.com`): Turnkey API base URL.
- **`TURNKEY_API_PRIVATE_KEY`** (string, required): Turnkey API private key for root user.
- **`TURNKEY_API_PUBLIC_KEY`** (string, required): Turnkey API public key for root user.
- **`TURNKEY_API_SENDER_PRIVATE_KEY`** (string, required): Turnkey API private key for bridge sender.
- **`TURNKEY_API_SENDER_PUBLIC_KEY`** (string, required): Turnkey API public key for bridge sender.
- **`TURNKEY_ORGANIZATION_ID`** (string, required): Turnkey parent organization ID.
- **`TURNKEY_MAGIC_LINK_TEMPLATE`** (string, optional): Custom magic link template for email authentication.
- **`TURNKEY_EMAIL_SENDER_ADDRESS`** (string, default: `notifications@mail.dydx.trade`): Email sender address.
- **`TURNKEY_EMAIL_SENDER_NAME`** (string, default: `dYdX Notifications`): Email sender name.
- **`SOLANA_SPONSOR_PUBLIC_KEY`** (string, required): Solana fee payer public key for sponsored transactions.
- **`ALCHEMY_API_KEY`** (string, required): Alchemy API key for webhook management and RPC calls.
- **`ALCHEMY_AUTH_TOKEN`** (string, required): Alchemy auth token for webhook updates.
- **`ALCHEMY_WEBHOOK_UPDATE_URL`** (string, default: `https://dashboard.alchemy.com/api/update-webhook-addresses`): Alchemy webhook update endpoint.
- **`ETHEREUM_WEBHOOK_ID`** (string, required): Alchemy webhook ID for Ethereum mainnet.
- **`ARBITRUM_WEBHOOK_ID`** (string, required): Alchemy webhook ID for Arbitrum.
- **`AVALANCHE_WEBHOOK_ID`** (string, required): Alchemy webhook ID for Avalanche.
- **`BASE_WEBHOOK_ID`** (string, required): Alchemy webhook ID for Base.
- **`OPTIMISM_WEBHOOK_ID`** (string, required): Alchemy webhook ID for Optimism.
- **`SOLANA_WEBHOOK_ID`** (string, required): Alchemy webhook ID for Solana.
- **`ZERODEV_API_KEY`** (string, required): ZeroDev API key for smart account sponsorship.
- **`ZERODEV_API_BASE_URL`** (string, default: `https://rpc.zerodev.app/api/v3`): ZeroDev RPC base URL.
- **`SKIP_SLIPPAGE_TOLERANCE_PERCENTAGE`** (string, default: `"0"`): Maximum slippage percentage for Skip Protocol swaps.
- **`SKIP_SLIPPAGE_TOLERANCE_USDC`** (integer, default: `100`): Maximum slippage in USDC for Skip Protocol swaps.
- **`BRIDGE_THRESHOLD_USDC`** (integer, default: `20`): Minimum USDC amount to trigger automatic bridging.
- **`MAXIMUM_BRIDGE_AMOUNT_USDC`** (integer, default: `999000`): Maximum USDC amount per bridge transaction.
- **`APPROVAL_SIGNER_PUBLIC_ADDRESS`** (string, required): Public address of the approval signer for policy enforcement.
- **`APPROVAL_ENABLED`** (boolean, default: `true`): Enable policy approvals for bridging.
- **`CALL_POLICY_VALUE_LIMIT`** (bigint, default: `100000000000000000000`): Maximum value limit for call policies.

### Apple Sign-In Configuration

- **`APPLE_TEAM_ID`** (string, required): Apple Developer Team ID.
- **`APPLE_SERVICE_ID`** (string, required): Apple Service ID (client_id).
- **`APPLE_KEY_ID`** (string, required): Apple Key ID.
- **`APPLE_PRIVATE_KEY`** (string, required): Apple private key (base64-encoded).
- **`APPLE_APP_SCHEME`** (string, default: `dydxV4`): Deep link scheme for Apple OAuth redirect.

### Firebase Configuration

- **`FIREBASE_PROJECT_ID`** (string, required): Firebase project ID.
- **`FIREBASE_PRIVATE_KEY`** (string, required): Firebase private key.
- **`FIREBASE_CLIENT_EMAIL`** (string, required): Firebase client email.

### Amplitude Analytics Configuration

- **`AMPLITUDE_API_KEY`** (string, optional): Amplitude API key for event tracking.
- **`AMPLITUDE_SERVER_URL`** (string, default: `https://api.eu.amplitude.com/2/httpapi`): Amplitude server URL.

### Cache-Control Directives

Each endpoint has a configurable `Cache-Control` header. All follow the pattern `CACHE_CONTROL_DIRECTIVE_<ENDPOINT>` (string, default varies by endpoint):
- `CACHE_CONTROL_DIRECTIVE_ADDRESSES` (default: `public, max-age=1`)
- `CACHE_CONTROL_DIRECTIVE_AFFILIATES` (default: `public, max-age=10`)
- `CACHE_CONTROL_DIRECTIVE_AFFILIATES_METADATA` (default: `no-cache`)
- `CACHE_CONTROL_DIRECTIVE_ASSET_POSITIONS` (default: `public, max-age=1`)
- `CACHE_CONTROL_DIRECTIVE_CANDLES` (default: `public, max-age=1`)
- `CACHE_CONTROL_DIRECTIVE_FILLS` (default: `public, max-age=1`)
- `CACHE_CONTROL_DIRECTIVE_FUNDING` (default: `public, max-age=10`)
- `CACHE_CONTROL_DIRECTIVE_HISTORICAL_BLOCK_TRADING_REWARDS` (default: `public, max-age=10`)
- `CACHE_CONTROL_DIRECTIVE_HISTORICAL_FUNDING` (default: `public, max-age=10`)
- `CACHE_CONTROL_DIRECTIVE_HISTORICAL_PNL` (default: `public, max-age=10`)
- `CACHE_CONTROL_DIRECTIVE_HISTORICAL_TRADING_REWARDS` (default: `public, max-age=10`)
- `CACHE_CONTROL_DIRECTIVE_ORDERBOOK` (default: `public, max-age=1`)
- `CACHE_CONTROL_DIRECTIVE_ORDERS` (default: `public, max-age=1`)
- `CACHE_CONTROL_DIRECTIVE_PERPETUAL_MARKETS` (default: `public, max-age=1`)
- `CACHE_CONTROL_DIRECTIVE_PERPETUAL_POSITIONS` (default: `public, max-age=1`)
- `CACHE_CONTROL_DIRECTIVE_SOCIAL_TRADING` (default: `public, max-age=10`)
- `CACHE_CONTROL_DIRECTIVE_SPARKLINES` (default: `public, max-age=10`)
- `CACHE_CONTROL_DIRECTIVE_TIME` (default: `no-cache, no-store, no-transform`)
- `CACHE_CONTROL_DIRECTIVE_TRADES` (default: `public, max-age=1`)
- `CACHE_CONTROL_DIRECTIVE_TRANSFERS` (default: `public, max-age=1`)
- `CACHE_CONTROL_DIRECTIVE_VAULTS` (default: `public, max-age=10`)
- `CACHE_CONTROL_DIRECTIVE_PNL` (default: `public, max-age=10`)

### Vault Configuration

- **`VAULT_PNL_HISTORY_DAYS`** (integer, default: `90`): Number of days of daily PnL history to return.
- **`VAULT_PNL_HISTORY_HOURS`** (integer, default: `72`): Number of hours of hourly PnL history to return.
- **`VAULT_PNL_START_DATE`** (string, default: `2024-01-01T00:00:00Z`): Start date for vault PnL calculations.
- **`VAULT_LATEST_PNL_TICK_WINDOW_HOURS`** (integer, default: `1`): Window for fetching latest PnL tick.
- **`VAULT_FETCH_FUNDING_INDEX_BLOCK_WINDOWS`** (integer, default: `250000`): Block window size for fetching funding indices.
- **`VAULT_CACHE_TTL_MS`** (integer, default: `120000`): Cache TTL for vault PnL data (2 minutes).

### Affiliate Configuration

- **`VOLUME_ELIGIBILITY_THRESHOLD`** (integer, default: `10000`): Minimum trading volume (in USDC) required for affiliate eligibility.

### API Limits

- **`API_LIMIT_V4`** (integer, default: `1000`): Maximum number of results per API request.
- **`API_ORDERBOOK_LEVELS_PER_SIDE_LIMIT`** (integer, default: `100`): Maximum orderbook levels per side.

## Running Locally

### Prerequisites

- Node.js 18+ and `pnpm` installed.
- PostgreSQL database running and accessible.
- Redis instance running and accessible.
- Environment variables configured (see `.env.development` or `.env.test` for examples).

### Steps

1. **Install dependencies** (from monorepo root):
   ```bash
   pnpm install
   ```

2. **Build the service**:
   ```bash
   cd indexer/services/comlink
   pnpm run build
   ```

3. **Set up environment variables**:
   - Copy `.env.development` or create a `.env` file with required variables.
   - Ensure `DB_HOSTNAME`, `DB_NAME`, `DB_USERNAME`, `DB_PASSWORD`, `REDIS_URL`, and `RATE_LIMIT_REDIS_URL` are set.
   - For Turnkey and bridging features, set the corresponding `TURNKEY_*`, `ALCHEMY_*`, and `ZERODEV_*` variables.

4. **Run the service**:
   ```bash
   pnpm run dev
   ```
   The service will start on `http://localhost:8080` (or the port specified in `PORT`).

5. **Access Swagger UI**:
   - Navigate to `http://localhost:8080/docs` to view interactive API documentation.

### Docker Compose (if available)

If a `docker-compose.yml` exists at the monorepo root, you can start all dependencies (PostgreSQL, Redis) with:
