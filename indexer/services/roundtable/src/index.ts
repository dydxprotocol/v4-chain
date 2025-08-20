import { logger, startBugsnag, wrapBackgroundTask } from '@dydxprotocol-indexer/base';
import { producer } from '@dydxprotocol-indexer/kafka';
import {
  LeaderboardPnlTimeSpan,
  TradingRewardAggregationPeriod,
} from '@dydxprotocol-indexer/postgres';

import config from './config';
import { complianceProvider } from './helpers/compliance-clients';
import { startLoop } from './helpers/loops-helper';
import { redisClient, connect as connectToRedis } from './helpers/redis';
import aggregateTradingRewardsTasks from './tasks/aggregate-trading-rewards';
import cacheOrderbookMidPrices from './tasks/cache-orderbook-mid-prices';
import cancelStaleOrdersTask from './tasks/cancel-stale-orders';
import createLeaderboardTask from './tasks/create-leaderboard';
import createPnlTicksTask from './tasks/create-pnl-ticks';
import deleteOldFastSyncSnapshots from './tasks/delete-old-fast-sync-snapshots';
import deleteOldFirebaseNotificationTokensTask from './tasks/delete-old-firebase-notification-tokens';
import deleteZeroPriceLevelsTask from './tasks/delete-zero-price-levels';
import marketUpdaterTask from './tasks/market-updater';
import orderbookInstrumentationTask from './tasks/orderbook-instrumentation';
import performComplianceStatusTransitionsTask from './tasks/perform-compliance-status-transitions';
import pnlInstrumentationTask from './tasks/pnl-instrumentation';
import refreshVaultPnlTask from './tasks/refresh-vault-pnl';
import removeExpiredOrdersTask from './tasks/remove-expired-orders';
import removeOldOrderUpdatesTask from './tasks/remove-old-order-updates';
import subaccountUsernameGeneratorTask from './tasks/subaccount-username-generator';
import takeFastSyncSnapshotTask from './tasks/take-fast-sync-snapshot';
import trackLag from './tasks/track-lag';
import uncrossOrderbookTask from './tasks/uncross-orderbook';
import updateAffiliateInfoTask from './tasks/update-affiliate-info';
import updateComplianceDataTask from './tasks/update-compliance-data';
import updateFundingPaymentsTask from './tasks/update-funding-payments';
import updatePnlTask from './tasks/update-pnl';
import updateResearchEnvironmentTask from './tasks/update-research-environment';
import updateWalletTotalVolumeTask from './tasks/update-wallet-total-volume';

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

  await Promise.all([producer.connect(), connectToRedis()]);

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

  if (config.LOOPS_ENABLED_UNCROSS_ORDERBOOK) {
    startLoop(
      uncrossOrderbookTask,
      'uncross_orderbook',
      config.LOOPS_INTERVAL_MS_UNCROSS_ORDERBOOK,
      config.UNCROSS_ORDERBOOK_LOCK_MULTIPLIER,
    );
  }

  if (config.LOOPS_ENABLED_PNL_TICKS) {
    startLoop(
      createPnlTicksTask,
      'create_pnl_ticks',
      config.LOOPS_INTERVAL_MS_PNL_TICKS,
      config.PNL_TICK_UPDATE_LOCK_MULTIPLIER,
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

  if (config.LOOPS_PNL_INSTRUMENTATION) {
    startLoop(
      pnlInstrumentationTask,
      'pnl_instrumentation',
      config.LOOPS_INTERVAL_MS_PNL_INSTRUMENTATION,
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

  startLoop(
    () => performComplianceStatusTransitionsTask(),
    'update_compliance_status',
    config.LOOPS_INTERVAL_MS_PERFORM_COMPLIANCE_STATUS_TRANSITIONS,
  );

  if (config.LOOPS_ENABLED_TRACK_LAG) {
    startLoop(trackLag, 'track_lag', config.LOOPS_INTERVAL_MS_TRACK_LAG);
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

  if (config.LOOPS_ENABLED_SUBACCOUNT_USERNAME_GENERATOR) {
    startLoop(
      subaccountUsernameGeneratorTask,
      'subaccount_username_generator',
      config.LOOPS_INTERVAL_MS_SUBACCOUNT_USERNAME_GENERATOR,
      config.SUBACCOUNT_USERNAME_GENERATOR_LOCK_MULTIPLIER,
    );
  }

  if (config.LOOPS_ENABLED_LEADERBOARD_PNL_ALL_TIME) {
    const allTimeLeaderboardTask: () => Promise<void> = createLeaderboardTask(
      LeaderboardPnlTimeSpan.ALL_TIME,
    );
    startLoop(
      allTimeLeaderboardTask,
      'create_leaderboard_pnl_all_time',
      config.LOOPS_INTERVAL_MS_LEADERBOARD_PNL_ALL_TIME,
    );
  }
  if (config.LOOPS_ENABLED_LEADERBOARD_PNL_DAILY) {
    const dailyLeaderboardTask: () => Promise<void> = createLeaderboardTask(
      LeaderboardPnlTimeSpan.ONE_DAY,
    );
    startLoop(
      dailyLeaderboardTask,
      'create_leaderboard_pnl_daily',
      config.LOOPS_INTERVAL_MS_LEADERBOARD_PNL_DAILY,
    );
  }
  if (config.LOOPS_ENABLED_LEADERBOARD_PNL_WEEKLY) {
    const weeklyLeaderboardTask: () => Promise<void> = createLeaderboardTask(
      LeaderboardPnlTimeSpan.SEVEN_DAYS,
    );
    startLoop(
      weeklyLeaderboardTask,
      'create_leaderboard_pnl_weekly',
      config.LOOPS_INTERVAL_MS_LEADERBOARD_PNL_WEEKLY,
    );
  }
  if (config.LOOPS_ENABLED_LEADERBOARD_PNL_MONTHLY) {
    const monthlyLeaderboardTask: () => Promise<void> = createLeaderboardTask(
      LeaderboardPnlTimeSpan.THIRTY_DAYS,
    );
    startLoop(
      monthlyLeaderboardTask,
      'create_leaderboard_pnl_monthly',
      config.LOOPS_INTERVAL_MS_LEADERBOARD_PNL_MONTHLY,
    );
  }
  if (config.LOOPS_ENABLED_LEADERBOARD_PNL_YEARLY) {
    const yearlyLeaderboardTask: () => Promise<void> = createLeaderboardTask(
      LeaderboardPnlTimeSpan.ONE_YEAR,
    );
    startLoop(
      yearlyLeaderboardTask,
      'create_leaderboard_pnl_yearly',
      config.LOOPS_INTERVAL_MS_LEADERBOARD_PNL_YEARLY,
    );
  }
  if (config.LOOPS_ENABLED_UPDATE_WALLET_TOTAL_VOLUME) {
    startLoop(
      updateWalletTotalVolumeTask,
      'update_wallet_total_volume',
      config.LOOPS_INTERVAL_MS_UPDATE_WALLET_TOTAL_VOLUME,
    );
  }
  if (config.LOOPS_ENABLED_UPDATE_AFFILIATE_INFO) {
    startLoop(
      updateAffiliateInfoTask,
      'update_affiliate_info',
      config.LOOPS_INTERVAL_MS_UPDATE_AFFILIATE_INFO,
    );
  }
  if (config.LOOPS_ENABLED_DELETE_OLD_FIREBASE_NOTIFICATION_TOKENS) {
    startLoop(
      deleteOldFirebaseNotificationTokensTask,
      'delete-old-firebase-notification-tokens',
      config.LOOPS_INTERVAL_MS_DELETE_FIREBASE_NOTIFICATION_TOKENS_MONTHLY,
    );
  }
  if (config.LOOPS_ENABLED_REFRESH_VAULT_PNL) {
    startLoop(
      refreshVaultPnlTask,
      'refresh-vault-pnl',
      config.LOOPS_INTERVAL_MS_REFRESH_VAULT_PNL,
    );
  }
  if (config.LOOPS_ENABLED_CACHE_ORDERBOOK_MID_PRICES) {
    startLoop(
      cacheOrderbookMidPrices,
      'cache-orderbook-mid-prices',
      config.LOOPS_INTERVAL_MS_CACHE_ORDERBOOK_MID_PRICES,
    );
  }
  if (config.LOOPS_ENABLED_UPDATE_FUNDING_PAYMENTS) {
    startLoop(
      updateFundingPaymentsTask,
      'update-funding-payments',
      config.LOOPS_INTERVAL_MS_UPDATE_FUNDING_PAYMENTS,
      // extended lock multiplier for 12 hours since on the first run,
      // the task takes a while to complete.
      config.UPDATE_FUNDING_PAYMENTS_LOCK_MULTIPLIER,
    );
  }

  if (config.LOOPS_ENABLED_UPDATE_PNL) {
    startLoop(
      updatePnlTask,
      'update-pnl',
      config.LOOPS_INTERVAL_MS_UPDATE_PNL,
      // extended lock multiplier for 12 hours since on the first run,
      // the task takes a while to complete.
      config.UPDATE_PNL_LOCK_MULTIPLIER,
    );
  }

  logger.info({
    at: 'index',
    message: 'Successfully started',
  });
}

wrapBackgroundTask(start(), true, 'main');
