import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  PnlTickInterval,
  PnlTicksFromDatabase,
  VaultFromDatabase,
  MEGAVAULT_SUBACCOUNT_ID,
  TransferFromDatabase,
  TransferTable,
  Ordering,
  TransferColumns,
  VaultPnlTicksView,
} from '@dydxprotocol-indexer/postgres';
import * as MegavaultHistoricalPnlCache from '@dydxprotocol-indexer/redis/src/caches/megavault-historical-pnl-cache';
import Big from 'big.js';
import bounds from 'binary-searching';
import _ from 'lodash';
import { DateTime } from 'luxon';

import {
  getVaultStartPnl,
} from '../caches/vault-start-pnl';
import config from '../config';
import { redisClient } from '../helpers/redis';
import {
  getVaultMapping,
  getVaultPnlStartDate,
} from '../lib/helpers';

interface AggregatedPnlTick {
  pnlTick: PnlTicksFromDatabase,
  numTicks: number,
}

/**
 * Cache megavault PNL data for both hourly and daily resolutions.
 */
export default async function runTask(): Promise<void> {
  const taskStart: number = Date.now();
  try {
    const vaultSubaccounts: { [subaccountId: string]: VaultFromDatabase } = await getVaultMapping();

    // Cache both hourly and daily PNL data
    await Promise.all([
      cacheMegavaultPnl(PnlTickInterval.hour, vaultSubaccounts),
      cacheMegavaultPnl(PnlTickInterval.day, vaultSubaccounts),
    ]);

    stats.timing(
      `${config.SERVICE_NAME}.cache-megavault-pnl.timing`,
      Date.now() - taskStart,
    );
  } catch (error) {
    logger.error({
      at: 'cache-megavault-pnl#runTask',
      message: 'Failed to cache megavault PNL data',
      error,
    });
  }
}

async function getVaultSubaccountPnlTicks(
  vaultSubaccountIds: string[],
  resolution: PnlTickInterval,
): Promise<PnlTicksFromDatabase[]> {
  if (vaultSubaccountIds.length === 0) {
    return [];
  }

  let windowSeconds: number;
  if (resolution === PnlTickInterval.day) {
    // windowSeconds = config.VAULT_PNL_HISTORY_DAYS * 24 * 60 * 60; // days to seconds
    windowSeconds = 90 * 24 * 60 * 60; // NEXT: use config
  } else {
    // windowSeconds = config.VAULT_PNL_HISTORY_HOURS * 60 * 60; // hours to seconds
    windowSeconds = 72 * 60 * 60; // hours to seconds // NEXT use config
  }

  const pnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
    resolution,
    windowSeconds,
    getVaultPnlStartDate(),
  );

  return adjustVaultPnlTicks(pnlTicks, getVaultStartPnl());
}

async function cacheMegavaultPnl(
  resolution: PnlTickInterval,
  vaultSubaccounts: { [subaccountId: string]: VaultFromDatabase },
): Promise<void> {
  const vaultSubaccountIdsWithMainSubaccount: string[] = _
    .keys(vaultSubaccounts)
    .concat([MEGAVAULT_SUBACCOUNT_ID]);

  const [
    vaultPnlTicks,
    firstMainVaultTransferTimestamp,
  ] = await Promise.all([
    getVaultSubaccountPnlTicks(vaultSubaccountIdsWithMainSubaccount, resolution),
    getFirstMainVaultTransferDateTime(),
  ]);

  // Aggregate pnlTicks for all vault subaccounts grouped by blockHeight
  const aggregatedPnlTicks: PnlTicksFromDatabase[] = aggregateVaultPnlTicks(
    vaultPnlTicks,
    _.values(vaultSubaccounts),
    firstMainVaultTransferTimestamp,
  );

  const filteredPnlTicks: PnlTicksFromDatabase[] = filterOutIntervalTicks(
    aggregatedPnlTicks,
    resolution,
  );

  // Cache the computed PNL ticks
  await MegavaultHistoricalPnlCache.set(
    resolution,
    filteredPnlTicks,
    redisClient,
  );

  logger.info({
    at: 'cache-megavault-pnl#cacheMegavaultPnl',
    message: 'Successfully cached megavault PNL data',
    resolution,
    numTicks: filteredPnlTicks.length,
  });
}

// NEXT: collapse with comlink vault-controller.

/**
 * Takes in an array of PnlTicks and filters out the closest pnl tick per interval.
 * @param pnlTicks Array of pnl ticks.
 * @param resolution Resolution of interval.
 * @returns Array of PnlTicksFromDatabase, one per interval.
 */
function filterOutIntervalTicks(
  pnlTicks: PnlTicksFromDatabase[],
  resolution: PnlTickInterval,
): PnlTicksFromDatabase[] {
  // Track start of intervals to closest Pnl tick.
  const ticksPerInterval: Map<string, PnlTicksFromDatabase> = new Map();
  pnlTicks.forEach((pnlTick: PnlTicksFromDatabase): void => {
    const blockTime: DateTime = DateTime.fromISO(pnlTick.blockTime).toUTC();

    const startOfInterval: DateTime = blockTime.toUTC().startOf(resolution);
    const startOfIntervalStr: string = startOfInterval.toISO();
    const tickForInterval: PnlTicksFromDatabase | undefined = ticksPerInterval.get(
      startOfIntervalStr,
    );
    // No tick for the start of interval, set this tick as the block for the interval.
    if (tickForInterval === undefined) {
      ticksPerInterval.set(startOfIntervalStr, pnlTick);
      return;
    }
    const tickPerIntervalBlockTime: DateTime = DateTime.fromISO(tickForInterval.blockTime);

    // This tick is closer to the start of the interval, set it as the tick for the interval.
    if (blockTime.diff(startOfInterval) < tickPerIntervalBlockTime.diff(startOfInterval)) {
      ticksPerInterval.set(startOfIntervalStr, pnlTick);
    }
  });
  return Array.from(ticksPerInterval.values());
}

/**
 * Aggregates vault pnl ticks per hour, filtering out pnl ticks made up of less ticks than expected.
 * Expected number of pnl ticks is calculated from the number of vaults that were created before
 * the pnl tick was created.
 * @param vaultPnlTicks Pnl ticks to aggregate.
 * @param vaults List of all valid vaults.
 * @param mainVaultCreatedAt Date time when the main vault was created or undefined if it does not
 * exist yet.
 * @returns
 */
function aggregateVaultPnlTicks(
  vaultPnlTicks: PnlTicksFromDatabase[],
  vaults: VaultFromDatabase[],
  mainVaultCreatedAt?: DateTime,
): PnlTicksFromDatabase[] {
  // aggregate pnlTicks for all vault subaccounts grouped by blockHeight
  const aggregatedPnlTicks: AggregatedPnlTick[] = aggregateHourlyPnlTicks(vaultPnlTicks);
  const vaultCreationTimes: DateTime[] = _.map(vaults, 'createdAt').map(
    (createdAt: string) => DateTime.fromISO(createdAt),
  ).concat(
    mainVaultCreatedAt === undefined ? [] : [mainVaultCreatedAt],
  ).sort(
    (a: DateTime, b: DateTime) => a.diff(b).milliseconds,
  );
  return aggregatedPnlTicks.filter((aggregatedTick: AggregatedPnlTick) => {
    // Get number of vaults created before the pnl tick was created by binary-searching for the
    // index of the pnl ticks createdAt in a sorted array of vault createdAt times.
    const numVaultsCreated: number = bounds.le(
      vaultCreationTimes,
      DateTime.fromISO(aggregatedTick.pnlTick.createdAt),
      (a: DateTime, b: DateTime) => a.diff(b).milliseconds,
    );
    // Number of ticks should be greater than number of vaults created before it
    return aggregatedTick.numTicks >= numVaultsCreated;
  }).map((aggregatedPnlTick: AggregatedPnlTick) => aggregatedPnlTick.pnlTick);
}

function adjustVaultPnlTicks(
  pnlTicks: PnlTicksFromDatabase[],
  pnlTicksToAdjustBy: PnlTicksFromDatabase[],
): PnlTicksFromDatabase[] {
  const subaccountToPnlTick: {[subaccountId: string]: PnlTicksFromDatabase} = {};
  for (const pnlTickToAdjustBy of pnlTicksToAdjustBy) {
    subaccountToPnlTick[pnlTickToAdjustBy.subaccountId] = pnlTickToAdjustBy;
  }

  return pnlTicks.map((pnlTick: PnlTicksFromDatabase): PnlTicksFromDatabase => {
    const adjustByPnlTick: PnlTicksFromDatabase | undefined = subaccountToPnlTick[
      pnlTick.subaccountId
    ];
    if (adjustByPnlTick === undefined) {
      return pnlTick;
    }
    return {
      ...pnlTick,
      totalPnl: Big(pnlTick.totalPnl).sub(Big(adjustByPnlTick.totalPnl)).toFixed(),
    };
  });
}

/**
 * Aggregates a list of PnL ticks, combining any PnL ticks for the same hour by summing
 * the equity, totalPnl, and net transfers.
 * Returns a map of aggregated pnl ticks and the number of ticks the aggreated tick is made up of.
 * @param pnlTicks
 * @returns
 */
export function aggregateHourlyPnlTicks(
  pnlTicks: PnlTicksFromDatabase[],
): AggregatedPnlTick[] {
  const hourlyPnlTicks: Map<string, PnlTicksFromDatabase> = new Map();
  const hourlySubaccountIds: Map<string, Set<string>> = new Map();
  for (const pnlTick of pnlTicks) {
    const truncatedTime: string = DateTime.fromISO(pnlTick.createdAt).startOf('hour').toISO();
    if (hourlyPnlTicks.has(truncatedTime)) {
      const subaccountIds: Set<string> = hourlySubaccountIds.get(truncatedTime) as Set<string>;
      if (subaccountIds.has(pnlTick.subaccountId)) {
        continue;
      }
      subaccountIds.add(pnlTick.subaccountId);
      const aggregatedTick: PnlTicksFromDatabase = hourlyPnlTicks.get(
        truncatedTime,
      ) as PnlTicksFromDatabase;
      hourlyPnlTicks.set(
        truncatedTime,
        {
          ...aggregatedTick,
          equity: (parseFloat(aggregatedTick.equity) + parseFloat(pnlTick.equity)).toString(),
          totalPnl: (parseFloat(aggregatedTick.totalPnl) + parseFloat(pnlTick.totalPnl)).toString(),
          netTransfers: (
            parseFloat(aggregatedTick.netTransfers) + parseFloat(pnlTick.netTransfers)
          ).toString(),
        },
      );
      hourlySubaccountIds.set(truncatedTime, subaccountIds);
    } else {
      hourlyPnlTicks.set(truncatedTime, pnlTick);
      hourlySubaccountIds.set(truncatedTime, new Set([pnlTick.subaccountId]));
    }
  }
  return Array.from(hourlyPnlTicks.keys()).map((hour: string): AggregatedPnlTick => {
    return {
      pnlTick: hourlyPnlTicks.get(hour) as PnlTicksFromDatabase,
      numTicks: (hourlySubaccountIds.get(hour) as Set<string>).size,
    };
  });
}

async function getFirstMainVaultTransferDateTime(): Promise<DateTime | undefined> {
  const { results }: {
    results: TransferFromDatabase[],
  } = await TransferTable.findAllToOrFromSubaccountId(
    {
      subaccountId: [MEGAVAULT_SUBACCOUNT_ID],
      limit: 1,
    },
    [],
    {
      orderBy: [[TransferColumns.createdAt, Ordering.ASC]],
    },
  );
  if (results.length === 0) {
    return undefined;
  }
  return DateTime.fromISO(results[0].createdAt);
}
