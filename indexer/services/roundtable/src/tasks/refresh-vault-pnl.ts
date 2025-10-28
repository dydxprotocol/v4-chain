import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  VaultPnlTicksView,
  VaultPnlView,
} from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';

import config from '../config';

/**
 * Refresh the vault pnl ticks views.
 */
export default async function runTask(): Promise<void> {
  const taskStart: number = Date.now();
  try {
    const currentTime: DateTime = DateTime.utc();

    // Refresh hourly views
    if (currentTime.diff(
      currentTime.startOf('hour'),
    ).toMillis() < config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS) {
      logger.info({
        at: 'refresh-vault-pnl#runTask',
        message: 'Refreshing vault hourly pnl views (old and new)',
        currentTime,
      });

      const hourlyStart: number = Date.now();

      // Refresh both old and new views in parallel
      await Promise.all([
        VaultPnlTicksView.refreshHourlyView().then(() => {
          logger.info({
            at: 'refresh-vault-pnl#runTask',
            message: 'Successfully refreshed old hourly view (vaults_hourly_pnl)',
          });
        }),
        VaultPnlView.refreshHourlyView().then(() => {
          logger.info({
            at: 'refresh-vault-pnl#runTask',
            message: 'Successfully refreshed new hourly view (vaults_hourly_pnl_v2)',
          });
        }),
      ]);

      stats.timing(
        `${config.SERVICE_NAME}.refresh-vault-pnl.hourly-view.timing`,
        Date.now() - hourlyStart,
      );
    }

    // Refresh daily views
    if (currentTime.diff(
      currentTime.startOf('day'),
    ).toMillis() < config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS) {
      logger.info({
        at: 'refresh-vault-pnl#runTask',
        message: 'Refreshing vault daily pnl views (old and new)',
        currentTime,
      });

      const dailyStart: number = Date.now();

      // Refresh both old and new views in parallel
      await Promise.all([
        VaultPnlTicksView.refreshDailyView().then(() => {
          logger.info({
            at: 'refresh-vault-pnl#runTask',
            message: 'Successfully refreshed old daily view (vaults_daily_pnl)',
          });
        }),
        VaultPnlView.refreshDailyView().then(() => {
          logger.info({
            at: 'refresh-vault-pnl#runTask',
            message: 'Successfully refreshed new daily view (vaults_daily_pnl_v2)',
          });
        }),
      ]);

      stats.timing(
        `${config.SERVICE_NAME}.refresh-vault-pnl.daily-view.timing`,
        Date.now() - dailyStart,
      );
    }
    stats.timing(
      `${config.SERVICE_NAME}.refresh-vault-pnl.total.timing`,
      Date.now() - taskStart,
    );
  } catch (error) {
    logger.error({
      at: 'refresh-vault-pnl#runTask',
      message: 'Failed to refresh vault pnl views',
      error,
    });
  }
}
