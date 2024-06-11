import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  BlockFromDatabase,
  BlockTable,
  PnlTicksTable,
  SubaccountFromDatabase,
  SubaccountTable,
  TransferTable,
} from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';

/**
 * Instrument data on PNL to be used for analytics.
 */
export default async function runTask(): Promise<void> {
  const startTaskTime = DateTime.utc();

  const block: BlockFromDatabase = await
  BlockTable.getLatest({ readReplica: true });
  const latestBlockHeight: string = block.blockHeight;

  // Get all subaccounts with transfers
  const subaccountsWithTransfers: SubaccountFromDatabase[] = await
  SubaccountTable.getSubaccountsWithTransfers(latestBlockHeight, { readReplica: true });

  const subaccountIds: string[] = subaccountsWithTransfers.map(
    (subaccount: SubaccountFromDatabase) => subaccount.id,
  );

  // Get the most recent PNL ticks for each subaccount in DB
  const mostRecentPnlTicks: {
    [subaccountId: string]: string
  } = await PnlTicksTable.findMostRecentPnlTickTimeForEachAccount(
    '1',
  );

  // Check last PNL computation for each subaccount
  const stalePnlSubaccounts: string[] = [];
  const subaccountsWithPnl: string[] = Object.keys(mostRecentPnlTicks);
  subaccountIds.forEach((id: string) => {
    const lastPnlTick: string = mostRecentPnlTicks[id];
    if (lastPnlTick) {
      const lastPnlTime: DateTime = DateTime.fromISO(lastPnlTick);
      const hoursSinceLastPnl = startTaskTime.diff(lastPnlTime, 'hours').hours;

      if (hoursSinceLastPnl >= 2) {
        stalePnlSubaccounts.push(id);
      }
    }
  });

  // Get the subaccounts without PNL data
  const subaccountsWithoutPnl: string[] = subaccountIds.filter(
    (id: string) => !subaccountsWithPnl.includes(id),
  );

  // Get the last transfer time for each subaccount without PNL data
  const transferTimes: { [subaccountId: string]: string } = await
  TransferTable.getLastTransferTimeForSubaccounts(
    subaccountsWithoutPnl,
  );

  const staleTransferSubaccounts: string[] = [];
  // Check last transfer time for each subaccount without PNL data
  // If the last transfer time is more than 2 hours ago, add to stale subaccounts
  Object.entries(transferTimes).forEach(([subaccountId, time]) => {
    const lastTransferTime: DateTime = DateTime.fromISO(time);
    const hoursSinceLastTransfer = startTaskTime.diff(lastTransferTime, 'hours').hours;

    if (hoursSinceLastTransfer >= 2) {
      staleTransferSubaccounts.push(subaccountId);
    }
  });

  stats.gauge('pnl_stale_subaccounts', stalePnlSubaccounts.length + staleTransferSubaccounts.length);
  stats.gauge('pnl_stale_subaccounts_with_prior_pnl', stalePnlSubaccounts.length);
  stats.gauge('pnl_stale_subaccounts_without_prior_pnl', staleTransferSubaccounts.length);
  if (stalePnlSubaccounts.length > 0 || staleTransferSubaccounts.length > 0) {
    logger.error({
      at: 'pnl-instrumentation#statPnl',
      message: 'Subaccount ids with stale PNL data',
      stalePnlSubaccounts,
      staleTransferSubaccounts,
    });
  }
}
