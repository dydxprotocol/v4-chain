import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  VaultPnlTicksView,
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
    if (currentTime.diff(
      currentTime.startOf('hour'),
    ).toMillis() < config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS) {
      logger.info({
        at: 'refresh-vault-pnl#runTask',
        message: 'Refreshing vault hourly pnl view',
        currentTime,
      });
      await VaultPnlTicksView.refreshHourlyView();
      stats.timing(
        `${config.SERVICE_NAME}.refresh-vault-pnl.hourly-view.timing`,
        Date.now() - taskStart,
      );
    }

    if (currentTime.diff(
      currentTime.startOf('day'),
    ).toMillis() < config.TIME_WINDOW_FOR_REFRESH_VAULT_PNL_MS) {
      const refreshDailyStart: number = Date.now();
      logger.info({
        at: 'refresh-vault-pnl#runTask',
        message: 'Refreshing vault daily pnl view',
        currentTime,
      });
      await VaultPnlTicksView.refreshDailyView();
      stats.timing(
        `${config.SERVICE_NAME}.refresh-vault-pnl.daily-view.timing`,
        Date.now() - refreshDailyStart,
      );
    }
  } catch (error) {
    logger.error({
      at: 'refresh-vault-pnl#runTask',
      message: 'Failed to refresh vault pnl views',
      error,
    });
  }
}
