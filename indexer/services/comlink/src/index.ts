import {
  logger,
  wrapBackgroundTask,
  startBugsnag,
} from '@dydxprotocol-indexer/base';
import { perpetualMarketRefresher, liquidityTierRefresher } from '@dydxprotocol-indexer/postgres';

import { startVaultStartPnlCache } from './caches/vault-start-pnl';
import config from './config';
import IndexV4 from './controllers/api/index-v4';
import { connect as connectToRedis, connectReadOnly as connectToRedisReadOnly } from './helpers/redis/redis-controller';
import { initializeAmplitude } from './lib/amplitude-helpers';
import Server from './request-helpers/server';

process.on('SIGTERM', () => {
  logger.info({
    at: 'index#SIGTERM',
    message: 'Received SIGTERM, shutting down',
  });
  process.exit(0);
});

function startServer() {
  const app = Server(IndexV4);
  const port: number = config.PORT;
  const server = app.listen(port, () => {
    logger.info({
      at: 'index#startServer',
      message: `Api server is listening on ${port}`,
    });
  });

  server.keepAliveTimeout = config.KEEP_ALIVE_MS;
  server.headersTimeout = config.HEADERS_TIMEOUT_MS;
}

async function start() {
  startBugsnag();

  // Initialize Amplitude
  initializeAmplitude();

  // Initialize PerpetualMarkets cache
  await Promise.all([
    perpetualMarketRefresher.updatePerpetualMarkets(),
    liquidityTierRefresher.updateLiquidityTiers(),
  ]);
  wrapBackgroundTask(perpetualMarketRefresher.start(), true, 'startUpdatePerpetualMarkets');
  wrapBackgroundTask(liquidityTierRefresher.start(), true, 'startUpdateLiquidityTiers');
  // Initialize cache for vault start PnL
  await startVaultStartPnlCache();

  await connectToRedis();
  logger.info({
    at: 'index#start',
    message: `Connected to redis at ${config.REDIS_URL}`,
  });

  await connectToRedisReadOnly();
  logger.info({
    at: 'index#start',
    message: `Connected to read-only redis at ${config.REDIS_READONLY_URL}`,
  });

  startServer();
}

wrapBackgroundTask(start(), true, 'main');
