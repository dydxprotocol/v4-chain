import { logger, startBugsnag, wrapBackgroundTask } from '@dydxprotocol-indexer/base';
import { stopConsumer, startConsumer } from '@dydxprotocol-indexer/kafka';
import {
  assetRefresher, perpetualMarketRefresher, liquidityTierRefresher,
} from '@dydxprotocol-indexer/postgres';

import { initializeAllCaches } from './caches/block-cache';
import * as OrderbookMidPriceMemoryCache from './caches/orderbook-mid-price-memory-cache';
import config from './config';
import { connect } from './helpers/kafka/kafka-controller';
import { createPostgresFunctions } from './helpers/postgres/postgres-functions';
import {
  connect as connectToRedis,
  redisClient,
} from './helpers/redis/redis-controller';

async function startKafka(): Promise<void> {
  logger.info({
    at: 'index#start',
    message: `Starting in env ${config.NODE_ENV}`,
  });

  // TODO(DEC-1655): When PerpetualMarkets can be updated with events, create a custom cache for
  // Ender. Initialize PerpetualMarkets cache
  await Promise.all([
    perpetualMarketRefresher.updatePerpetualMarkets(),
    assetRefresher.updateAssets(),
    liquidityTierRefresher.updateLiquidityTiers(),
  ]);
  // Ender does not need to refresh its caches in a loop because Ender is the only service that
  // writes to the key attributes of perpetual_markets, asset_refresher, and market_refresher
  // The two exceptions are the aggregated properties of perpetual_markets and the
  // OrderbookMidPriceMemoryCache
  await initializeAllCaches();
  wrapBackgroundTask(OrderbookMidPriceMemoryCache.start(), true, 'startUpdateOrderbookMidPrices');

  await connect();
  await startConsumer();

  logger.info({
    at: 'index#start',
    message: 'Successfully started',
  });
}

process.on('SIGTERM', async () => {
  logger.info({
    at: 'index#SIGTERM',
    message: 'Received SIGTERM, shutting down',
  });
  await stopConsumer();
  redisClient.quit();
});

async function start(): Promise<void> {
  startBugsnag();
  logger.info({
    at: 'index#start',
    message: `Connecting to redis: ${config.REDIS_URL}`,
  });
  logger.info({
    at: 'index#start',
    message: `Connecting to kafka brokers: ${config.KAFKA_BROKER_URLS}`,
  });
  await Promise.all([
    connectToRedis(),
    startKafka(),
    createPostgresFunctions(),
  ]);
  logger.info({
    at: 'index#start',
    message: `Successfully connected to redis: ${config.REDIS_URL}`,
  });
  logger.info({
    at: 'index#start',
    message: `Successfully connected to kafka brokers: ${config.KAFKA_BROKER_URLS}`,
  });
}

wrapBackgroundTask(start(), true, 'main');
