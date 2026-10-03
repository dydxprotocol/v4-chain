# Indexer Packages

Internal packages for the dYdX Indexer monorepo. These libraries provide shared utilities, data-access layers, protocol-parsing helpers, and infrastructure integrations used by Indexer services. This is not a standalone service.

## Responsibilities & Scope

- Provide reusable building blocks for Indexer services:
  - Base utilities: configuration parsing, logging, metrics, error types, Axios helpers, timers.
  - Postgres data access: Objection.js models, store functions, types, and DB helpers.
  - Redis caches and Lua-scripted operations for order books and state.
  - Kafka client wrappers (producer/consumer/admin), batching helpers, and websocket message builders.
  - Notifications via Firebase Cloud Messaging.
  - Compliance provider integrations and geo-blocking helpers.
  - Protocol utilities (protobuf parsing and helpers).
- Offer typed public APIs per package for service consumption.
- Include configuration schemas for environment variables per package.

Non-goals:
- Running a network service by itself.
- Managing service lifecycles or top-level process orchestration.

## Public API

The main public surface is exposed via each package’s src/index.ts.

- Base utilities (@dydxprotocol-indexer/base)
  - Purpose: Logging, metrics, config parsing, Axios helpers, background task helpers, and misc utilities.
  - Key exports:
    - Config and parsing: baseConfigSchema, baseConfigSecrets, parseString/parseInteger/parseBoolean/parseBN/parseBig/parseBigInt/parseNumber, parseSchema. See ./base/src/config.ts and ./base/src/config-util.ts
    - Logging and transports: default logger, addTransportsToLogger side-effect on import, logger transports wired in ./base/src/add-transports-to-logger.ts
    - Metrics: default stats (hot-shots) ./base/src/stats.ts
    - Errors: CustomError, WrappedError, ConfigError, TooManyRequestsError, Axios* error types ./base/src/errors.ts and ./base/src/axios/errors.ts
    - Axios helpers: axiosRequest, safeAxiosRequest ./base/src/axios/axios-request.ts
    - Tasks/timers: delay, setIntervalNonOverlapping, wrapBackgroundTask ./base/src/tasks.ts
    - Utilities: sanitization (safeJsonStringify), date helpers, instance/availability zone IDs, cache-control middleware, constants. See ./base/src/*
- Postgres access (@dydxprotocol-indexer/postgres)
  - Purpose: Typed database models (Objection), store functions (query/build/update), typed data models, and DB helpers.
  - Key exports:
    - Types and constants: ./postgres/src/types/*, ./postgres/src/constants.ts
    - Models: default exports for key tables (e.g., OrderModel, MarketModel, …) ./postgres/src/models/*
    - Stores (query modules): e.g., OrderTable, BlockTable, FillTable, SubaccountTable, etc. ./postgres/src/stores/*
    - Loops (in-memory refreshers): perpetual-market, asset, block-height, liquidity-tier ./postgres/src/loops/*
    - Helpers: protocol translations, order translations, db helpers, store helpers, uuid. ./postgres/src/lib/* and ./postgres/src/helpers/*
    - Config schema: postgresConfigSchema ./postgres/src/config.ts
    - Public index: ./postgres/src/index.ts
- Redis caches (@dydxprotocol-indexer/redis)
  - Purpose: Typed Redis helpers, order book caches, order lifecycle scripts, stateful order updates, PnL/vault caches, and related utilities.
  - Key exports:
    - Redis helpers: createRedisClient, pttl/ttl/lockWithExpiry/… and async helpers for common commands ./redis/src/helpers/redis.ts
    - Order lifecycle: placeOrder, updateOrder, removeOrder using Lua SHA scripts ./redis/src/caches/*.ts and ./redis/src/caches/scripts.ts
    - Order book: updatePriceLevel/getOrderBookLevels/getOrderBookMidPrice + cleanup helpers ./redis/src/caches/orderbook-levels-cache.ts
    - Additional caches: open orders, expiries, mid prices, next funding, state-filled quantums, stateful order updates, leaderboard/aggregate reward processed flags, vault caches. See ./redis/src/caches/*
    - Types: cache and script types ./redis/src/types.ts
    - Config schema: redisConfigSchema ./redis/src/config.ts
    - Public index: ./redis/src/index.ts
- Kafka helpers (@dydxprotocol-indexer/kafka)
  - Purpose: Kafka client configuration and helpers for producers/consumers, batching, and websocket message generation.
  - Key exports:
    - Client singletons: kafka, producer, admin; consumer init/start/stop APIs ./kafka/src/*
    - Consumer control: initConsumer, startConsumer, stopConsumer, updateOnMessageFunction, updateOnBatchFunction ./kafka/src/consumer.ts
    - Batch producer utility: BatchKafkaProducer ./kafka/src/batch-kafka-producer.ts
    - Websocket message helpers (order/subaccount): ./kafka/src/websocket-helper.ts
    - Constants and topics/types: ./kafka/src/constants.ts, ./kafka/src/types.ts
    - Config schema: kafkaConfigSchema ./kafka/src/config.ts
- Notifications (@dydxprotocol-indexer/notifications)
  - Purpose: Firebase Cloud Messaging integration, localization, and typed notification payloads.
  - Key exports:
    - sendFirebaseMessage; initialization is internal with config ./notifications/src/message.ts and ./notifications/src/lib/firebase.ts
    - Localization helpers and messages: deriveLocalizedNotificationMessage ./notifications/src/localization.ts, ./notifications/src/localized-messages.ts
    - Typed notifications and factory: NotificationType, createNotification, language codes, topics ./notifications/src/types.ts
    - Config schema: notificationsConfigSchema ./notifications/src/config.ts
- Compliance (@dydxprotocol-indexer/compliance)
  - Purpose: Compliance data clients and geo restrictions.
  - Key exports:
    - getComplianceClient with support for PLACEHOLDER, BLOCKLIST, ELLIPTIC providers ./compliance/src/clients/clients.ts
    - Providers: EllipticProviderClient, BlocklistProviderClient, PlaceHolderProviderClient ./compliance/src/clients/*
    - Geo-blocking: isRestrictedCountryHeaders, isWhitelistedAddress ./compliance/src/geoblocking/restrict-countries.ts
    - Config schema & constants: ./compliance/src/config.ts, ./compliance/src/constants.ts
    - Types: ./compliance/src/types.ts
- Proto parser helpers (@dydxprotocol-indexer/v4-proto-parser)
  - Purpose: Utilities around v4 protos (bytes conversions, order ID hashing, flags, position sign).
  - Key exports: bytesToBigInt/bigIntToBytes/bytesToBase64/base64ToBytes, getOrderIdHash, isStatefulOrder, isLongTermOrder, requiresImmediateExecution, ORDER_FLAG_* constants. See ./v4-proto-parser/src/*

## Usage

- Base: configuration, logging, metrics, safe Axios

- Postgres: read latest block (read-replica-safe)

- Redis: connect and read mid-price levels

- Kafka: initialize consumer and handle messages

- Notifications: create and send a localized notification

- Compliance: resolve provider and check address

- Proto parser helpers: order utilities

## Dependencies

- Internal
  - @dydxprotocol-indexer/base is foundational for most packages.
  - @dydxprotocol-indexer/postgres types are shared with other packages.
  - @dydxprotocol-indexer/v4-protos are used across Kafka, Redis, parser helpers, and Postgres translations.

- External (not exhaustive; see package.jsons)
  - Base: axios, hot-shots, winston, bugsnag, lodash, traverse, uuid, @aws-sdk/client-ec2
  - Postgres: objection, knex, pg, luxon, big.js, lodash, long, uuid
  - Redis: redis (v2.x), bluebird, big.js, lodash, luxon, long
  - Kafka: kafkajs, lodash, uuid
  - Notifications: firebase-admin
  - Compliance: axios, lodash, crypto

Constraints and notes:
- Base logger transports are added on import (side effect in base/src/index.ts); configure env first.
- Stats use UDP via hot-shots; in tests, metrics are mocked.
- Postgres queries default to primary unless Options.readReplica is set.
- Redis scripts are lazily loaded on client “ready”; always await connect() before invoking script-backed APIs.
- Kafka producer/consumer require explicit connect/start via exported helpers.

## Performance & Constraints

- Redis order book operations and order lifecycle operations use Lua scripts (evalsha) for atomicity and speed; avoid calling before scripts are loaded (connect).
- Orderbook level updates are integer-quantum based to avoid float inaccuracy; use helpers to convert to human strings where needed.
- Kafka batch producer tracks and limits batch size (config.KAFKA_MAX_BATCH_WEBSOCKET_MESSAGE_SIZE_BYTES); large messages are split into multiple sends.
- Postgres store modules support pagination and orderBy; some heavy aggregations are pushed into SQL (materialized views, lateral joins) for performance.

## Directory Layout

- base/ — foundational utilities (config, logging, metrics, axios, tasks, helpers)
- postgres/ — DB types/models/stores, loops, helpers, and config
- redis/ — Redis helpers, caches, and Lua scripts
- kafka/ — Kafka client singletons, consumer/producer wrappers, batching, message helpers
- notifications/ — Firebase integration, localization, message sending
- compliance/ — provider clients, geo-blocking, config/types
- v4-proto-parser/ — proto-related helpers (bytes, flags, order hashing)
- example-package/ — template package

Each package follows:
- src/ — core implementation and public index.ts
- __tests__/ — unit/integration tests (where present)
- jest config and tsconfig files per package

## Integration Notes

- Configuration
  - Each package exposes a config schema (parseSchema(schema)) to validate and parse environment variables. Ensure envs are defined before importing packages that perform side-effectful initialization (e.g., base logger transports).

- Initialization
  - Base: importing @dydxprotocol-indexer/base attaches logger transports. No explicit init required beyond env.
  - Postgres: Objection is bound to the primary Knex connection on module load; use Options.readReplica for read paths. Transactions are managed via postgres/src/helpers/transaction.ts.
  - Redis: Always createRedisClient(url, reconnectTimeout) and await connect() before using caches; this loads and verifies Lua scripts.
  - Kafka: Call initConsumer() before startConsumer(); set handlers via updateOnMessageFunction/updateOnBatchFunction. Producer reconnects on disconnect; consumers auto-reconnect unless explicitly stopped.

- Error handling
  - Base Axios helpers throw typed errors (AxiosError/AxiosServerError or AxiosSafe* variants). Many store functions throw on invalid query shape (e.g., verifyAllRequiredFields) or validation errors.
  - Redis cache functions throw on validation errors (e.g., InvalidRedisOrderError, InvalidTotalFilledQuantumsError).
  - Postgres “upsert” builders rely on model idColumns; unknown enum/status inputs throw synchronously in translators.

- Concurrency
  - Redis Lua scripts ensure atomic updates; application code should avoid parallel conflicting updates to the same logical resource where possible.
  - Kafka consumer run is single process per consumer group member; batch vs per-message modes are selectable.

## Related Packages

- Protocol buffers: ../../v4-proto-js (publishing scripts) and repo-level proto/ definitions.
- Indexer services consuming these packages live under ../services (see service READMEs for integration specifics).
