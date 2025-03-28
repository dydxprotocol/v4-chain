import {
  logger, getInstanceId, startBugsnag, setInstanceId, wrapBackgroundTask,
} from '@dydxprotocol-indexer/base';
import { stopConsumer, startConsumer } from '@dydxprotocol-indexer/kafka';
import { blockHeightRefresher, perpetualMarketRefresher } from '@dydxprotocol-indexer/postgres';

import config from './config';
import { connect as connectToKafka } from './helpers/kafka/kafka-controller';
import {
  connect as connectToRedis,
  redisClient,
} from './helpers/redis/redis-controller';
import { flushAllQueues } from './lib/send-message-helper';

async function startService(): Promise<void> {
  logger.info({
    at: 'index#start',
    message: `Starting in env ${config.NODE_ENV}`,
  });

  startBugsnag();

  logger.info({
    at: 'index#start',
    message: 'Getting instance id...',
  });

  await setInstanceId();

  logger.info({
    at: 'index#start',
    message: `Got instance id ${getInstanceId()}.`,
  });

  // Initialize caches.
  await Promise.all([
    blockHeightRefresher.updateBlockHeight(),
    perpetualMarketRefresher.updatePerpetualMarkets(),
  ]);
  wrapBackgroundTask(blockHeightRefresher.start(), true, 'startUpdateBlockHeight');
  wrapBackgroundTask(perpetualMarketRefresher.start(), true, 'startUpdatePerpetualMarkets');

  logger.info({
    at: 'index#start',
    message: `Connecting to kafka brokers: ${config.KAFKA_BROKER_URLS}`,
  });

  logger.info({
    at: 'index#start',
    message: `Connecting to redis: ${config.REDIS_URL}`,
  });

  await Promise.all([
    connectToKafka(),
    connectToRedis(),
  ]);

  await startConsumer(config.BATCH_PROCESSING_ENABLED);

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
  await flushAllQueues();
  redisClient.quit();

  process.exit(0);
});

async function start(): Promise<void> {
  await startService();
}

wrapBackgroundTask(start(), true, 'main');
