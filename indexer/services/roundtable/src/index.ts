import { logger, startBugsnag, wrapBackgroundTask } from '@dydxprotocol-indexer/base';
import { producer } from '@dydxprotocol-indexer/kafka';
import { TradingRewardAggregationPeriod } from '@dydxprotocol-indexer/postgres';

import config from './config';
import { complianceProvider } from './helpers/compliance-clients';
import { startLoop } from './helpers/loops-helper';
import {
  redisClient,
  connect as connectToRedis,
} from './helpers/redis';
import aggregateTradingRewardsTasks from './tasks/aggregate-trading-rewards';
import cancelStaleOrdersTask from './tasks/cancel-stale-orders';
import createPnlTicksTask from './tasks/create-pnl-ticks';
import deleteOldFastSyncSnapshots from './tasks/delete-old-fast-sync-snapshots';
import deleteZeroPriceLevelsTask from './tasks/delete-zero-price-levels';
import marketUpdaterTask from './tasks/market-updater';
import orderbookInstrumentationTask from './tasks/orderbook-instrumentation';
import removeExpiredOrdersTask from './tasks/remove-expired-orders';
import removeOldOrderUpdatesTask from './tasks/remove-old-order-updates';
import takeFastSyncSnapshotTask from './tasks/take-fast-sync-snapshot';
import trackLag from './tasks/track-lag';
import updateComplianceDataTask from './tasks/update-compliance-data';
import updateResearchEnvironmentTask from './tasks/update-research-environment';

process.on('SIGTERM', () => {
  logger.info({
    at: 'index#SIGTERM',
    message: 'Received SIGTERM, shutting down',
  });
  redisClient.quit();

  process.exit(0);
});

async function start(): Promise<void> {
  logger.info({
    at: 'index#start',
    message: `Starting in env ${config.NODE_ENV}`,
  });
  startBugsnag();

  await Promise.all([
    producer.connect(),
    connectToRedis(),
  ]);

  if (config.LOOPS_ENABLED_MARKET_UPDATER) {
    startLoop(
      marketUpdaterTask,
      'market_updater',
      config.LOOPS_INTERVAL_MS_MARKET_UPDATER,
      config.MARKET_UPDATER_LOCK_MULTIPLIER,
    );
  }

  if (config.LOOPS_ENABLED_DELETE_ZERO_PRICE_LEVELS) {
    startLoop(
      deleteZeroPriceLevelsTask,
      'delete_zero_price_levels',
      config.LOOPS_INTERVAL_MS_DELETE_ZERO_PRICE_LEVELS,
      config.DELETE_ZERO_PRICE_LEVELS_LOCK_MULTIPLIER,
    );
  }

  if (config.LOOPS_ENABLED_PNL_TICKS) {
    startLoop(
      createPnlTicksTask,
      'create_pnl_ticks',
      config.LOOPS_INTERVAL_MS_PNL_TICKS,
    );
  }

  if (config.LOOPS_ENABLED_REMOVE_EXPIRED_ORDERS) {
    startLoop(
      removeExpiredOrdersTask,
      'remove_expired_orders',
      config.LOOPS_INTERVAL_MS_REMOVE_EXPIRED_ORDERS,
    );
  }

  if (config.LOOPS_ORDERBOOK_INSTRUMENTATION) {
    startLoop(
      orderbookInstrumentationTask,
      'orderbook_instrumentation',
      config.LOOPS_INTERVAL_MS_ORDERBOOK_INSTRUMENTATION,
    );
  }

  if (config.LOOPS_CANCEL_STALE_ORDERS) {
    startLoop(
      cancelStaleOrdersTask,
      'cancel_stale_orders',
      config.LOOPS_INTERVAL_MS_CANCEL_STALE_ORDERS,
    );
  }

  if (config.LOOPS_ENABLED_UPDATE_RESEARCH_ENVIRONMENT) {
    startLoop(
      updateResearchEnvironmentTask,
      'update_research_environment',
      config.LOOPS_INTERVAL_MS_UPDATE_RESEARCH_ENVIRONMENT,
    );
  }

  if (config.LOOPS_ENABLED_TAKE_FAST_SYNC_SNAPSHOTS) {
    startLoop(
      takeFastSyncSnapshotTask,
      'take_fast_sync_snapshot',
      config.LOOPS_INTERVAL_MS_TAKE_FAST_SYNC_SNAPSHOTS,
    );
  }

  if (config.LOOPS_ENABLED_DELETE_OLD_FAST_SYNC_SNAPSHOTS) {
    startLoop(
      deleteOldFastSyncSnapshots,
      'delete_old_fast_sync_snapshots',
      config.LOOPS_INTERVAL_MS_DELETE_OLD_FAST_SYNC_SNAPSHOTS,
    );
  }

  startLoop(
    () => updateComplianceDataTask(complianceProvider),
    'update_compliance_data',
    config.LOOPS_INTERVAL_MS_UPDATE_COMPLIANCE_DATA,
  );

  if (config.LOOPS_ENABLED_TRACK_LAG) {
    startLoop(
      trackLag,
      'track_lag',
      config.LOOPS_INTERVAL_MS_TRACK_LAG,
    );
  }

  if (config.LOOPS_ENABLED_REMOVE_OLD_ORDER_UPDATES) {
    startLoop(
      removeOldOrderUpdatesTask,
      'remove_old_order_updates',
      config.LOOPS_INTERVAL_MS_REMOVE_OLD_ORDER_UPDATES,
    );
  }

  if (config.LOOPS_ENABLED_AGGREGATE_TRADING_REWARDS_DAILY) {
    startLoop(
      aggregateTradingRewardsTasks(TradingRewardAggregationPeriod.DAILY),
      'aggregate_trading_rewards_daily',
      config.LOOPS_INTERVAL_MS_AGGREGATE_TRADING_REWARDS,
    );
  }

  if (config.LOOPS_ENABLED_AGGREGATE_TRADING_REWARDS_WEEKLY) {
    startLoop(
      aggregateTradingRewardsTasks(TradingRewardAggregationPeriod.WEEKLY),
      'aggregate_trading_rewards_weekly',
      config.LOOPS_INTERVAL_MS_AGGREGATE_TRADING_REWARDS,
    );
  }

  if (config.LOOPS_ENABLED_AGGREGATE_TRADING_REWARDS_MONTHLY) {
    startLoop(
      aggregateTradingRewardsTasks(TradingRewardAggregationPeriod.MONTHLY),
      'aggregate_trading_rewards_monthly',
      config.LOOPS_INTERVAL_MS_AGGREGATE_TRADING_REWARDS,
    );
  }

  logger.info({
    at: 'index',
    message: 'Successfully started',
  });
}

wrapBackgroundTask(start(), true, 'main');
