import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  VaultPnlTicksView
} from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';
import config from '../config';

/**
 * Update the affiliate info for all affiliate addresses.
 */
export default async function runTask(): Promise<void> {
  const taskStart: number = Date.now();
  
  const currentTime: DateTime = DateTime.utc();
  if (currentTime.diff(
      currentTime.startOf('hour')
    ).toMillis() < config.TIME_WINDOW_FOR_REFRESH_MS) {
    logger.info({
      at: 'refresh-vault-pnl#runTask',
      message: 'Refreshing vault hourly pnl view',
      currentTime,
    });
    await VaultPnlTicksView.refreshHourlyView();
    stats.timing(
      `${config.SERVICE_NAME}.refresh-vault-pnl.hourly-view.timing`,
      Date.now() - taskStart,
    )
  }

  const refreshDailyStart: number = Date.now();
  if (currentTime.diff(
      currentTime.startOf('day')
    ).toMillis() < config.TIME_WINDOW_FOR_REFRESH_MS) {
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
}
