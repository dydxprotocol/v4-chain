import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  AssetPositionTable,
  FundingIndexMap,
  FundingIndexUpdatesTable,
  IsoString,
  OraclePriceTable,
  PerpetualPositionFromDatabase,
  PerpetualPositionTable,
  PnlTicksCreateObject,
  PnlTicksTable,
  PriceMap,
  SubaccountAssetNetTransferMap,
  SubaccountFromDatabase,
  SubaccountTable,
  SubaccountToPerpetualPositionsMap,
  TransferTable,
  helpers,
} from '@dydxprotocol-indexer/postgres';
import { LatestAccountPnlTicksCache, PnlTickForSubaccounts } from '@dydxprotocol-indexer/redis';
import Big from 'big.js';
import _ from 'lodash';
import { DateTime } from 'luxon';

import config from '../config';
import { USDC_ASSET_ID, ZERO } from '../lib/constants';
import { redisClient } from './redis';
import { SubaccountUsdcTransferMap } from './types';

/**
 * Gets a batch of new pnl ticks to write to the database and set in the cache.
 * @param blockHeight: consider transfers up until this block height.
 */
export async function getPnlTicksCreateObjects(
  blockHeight: string,
  blockTime: IsoString,
  txId: number,
): Promise<PnlTicksCreateObject[]> {
  const startGetPnlTicksCreateObjects: number = Date.now();
  const pnlTicksToBeCreatedAt: DateTime = DateTime.utc();
  const [
    mostRecentPnlTicks,
    subaccountsWithTransfers,
  ]: [
    PnlTickForSubaccounts,
    SubaccountFromDatabase[],
  ] = await Promise.all([
    getMostRecentPnlTicksForEachAccount(),
    SubaccountTable.getSubaccountsWithTransfers(blockHeight, { readReplica: true, txId }),
  ]);
  stats.timing(
    `${config.SERVICE_NAME}_get_ticks_relevant_accounts`,
    new Date().getTime() - startGetPnlTicksCreateObjects,
  );
  const accountToLastUpdatedBlockTime: _.Dictionary<IsoString> = _.mapValues(
    mostRecentPnlTicks,
    (pnlTick: PnlTicksCreateObject) => pnlTick.blockTime,
  );
  const subaccountIdsWithTranfers: string[] = _.map(subaccountsWithTransfers, 'id');
  const newSubaccountIds: string[] = _.difference(
    subaccountIdsWithTranfers, _.keys(accountToLastUpdatedBlockTime),
  );
  // get accounts to update based on last updated block height
  const accountsToUpdate: string[] = [
    ..._.keys(accountToLastUpdatedBlockTime).filter(
      (accountId) => {
        const lastUpdatedBlockTime: string = accountToLastUpdatedBlockTime[accountId];
        return new Date(blockTime).getTime() - new Date(lastUpdatedBlockTime).getTime() >=
          config.PNL_TICK_UPDATE_INTERVAL_MS;
      },
    ),
    ...newSubaccountIds,
  ];
  stats.gauge(
    `${config.SERVICE_NAME}_get_ticks_accounts_to_update`,
    accountsToUpdate.length,
  );
  const idToSubaccount: _.Dictionary<SubaccountFromDatabase> = _.keyBy(
    subaccountsWithTransfers,
    'id',
  );
  const getFundingIndexStart: number = Date.now();
  const blockHeightToFundingIndexMap: _.Dictionary<FundingIndexMap> = await
  getBlockHeightToFundingIndexMap(
    subaccountsWithTransfers,
    accountsToUpdate,
    txId,
  );
  stats.timing(
    `${config.SERVICE_NAME}_get_ticks_funding_indices`,
    new Date().getTime() - getFundingIndexStart,
  );

  const getAccountInfoStart: number = Date.now();
  const [
    subaccountTotalTransfersMap,
    openPerpetualPositions,
    usdcAssetPositions,
    netUsdcTransfers,
    markets,
    currentFundingIndexMap,
  ]: [
    SubaccountAssetNetTransferMap,
    SubaccountToPerpetualPositionsMap,
    { [subaccountId: string]: Big },
    SubaccountUsdcTransferMap,
    PriceMap,
    FundingIndexMap,
  ] = await Promise.all([
    TransferTable.getNetTransfersPerSubaccount(
      blockHeight,
      {
        readReplica: true,
        txId,
      },
    ),
    PerpetualPositionTable.findOpenPositionsForSubaccounts(
      accountsToUpdate,
      {
        readReplica: true,
        txId,
      },
    ),
    AssetPositionTable.findUsdcPositionForSubaccounts(
      accountsToUpdate,
      {
        readReplica: true,
        txId,
      },
    ),
    getUsdcTransfersSinceLastPnlTick(
      accountsToUpdate,
      mostRecentPnlTicks,
      blockHeight,
      txId,
    ),
    OraclePriceTable.findLatestPrices(blockHeight),
    FundingIndexUpdatesTable.findFundingIndexMap(blockHeight),
  ]);
  stats.timing(
    `${config.SERVICE_NAME}_get_ticks_account_info`,
    new Date().getTime() - getAccountInfoStart,
  );

  const computePnlStart: number = Date.now();
  const newTicksToCreate: PnlTicksCreateObject[] = accountsToUpdate.map(
    (account: string) => getNewPnlTick(
      account,
      subaccountTotalTransfersMap,
      markets,
      Object.values(openPerpetualPositions[account] || {}),
      usdcAssetPositions[account] || ZERO,
      netUsdcTransfers[account] || ZERO,
      pnlTicksToBeCreatedAt,
      blockHeight,
      blockTime,
      mostRecentPnlTicks,
      blockHeightToFundingIndexMap[idToSubaccount[account].updatedAtHeight],
      currentFundingIndexMap,
    ),
  );
  stats.timing(
    `${config.SERVICE_NAME}_get_ticks_compute_pnl`,
    new Date().getTime() - computePnlStart,
  );
  return newTicksToCreate;
}

/**
 * Get a map of block height to funding index state.
 * Funding index state represents the most recent funding index value for every perpetual market.
 *
 * @param subaccountsWithTransfers
 * @param accountsToUpdate
 */
export async function getBlockHeightToFundingIndexMap(
  subaccountsWithTransfers: SubaccountFromDatabase[],
  accountsToUpdate: string[],
  txId: number | undefined = undefined,
): Promise<_.Dictionary<FundingIndexMap>> {
  const idToSubaccount: _.Dictionary<SubaccountFromDatabase> = _.keyBy(
    subaccountsWithTransfers,
    'id',
  );
  // get the subaccount id to last updated block height
  const blockHeights: Set<string> = _.reduce(
    accountsToUpdate,
    (acc: Set<string>, accountId: string) => {
      acc.add(idToSubaccount[accountId].updatedAtHeight);
      return acc;
    },
    new Set<string>(),
  );

  const fundingIndexMaps: FundingIndexMap[] = await Promise.all(
    [...blockHeights].map(
      (blockHeight: string) => FundingIndexUpdatesTable.findFundingIndexMap(
        blockHeight,
        {
          readReplica: true,
          txId,
        },
      ),
    ),
  );

  const blockHeightToFundingIndexMap: _.Dictionary<FundingIndexMap> = _.zipObject(
    Object.values([...blockHeights]), fundingIndexMaps,
  );
  return blockHeightToFundingIndexMap;
}

/**
 * Get the most recent pnl tick for a given subaccount
 * @param subaccountId: subaccountId to compute pnl tick for
 * @param subaccountTotalTransfersMap: total historical transfers across all subaccounts
 * @param markets: latest market prices effectiveBeforeOrAt latestBlockHeight
 * @param openPerpetualPositionsForSubaccount: list of open perpetual positions for given subaccount
 * @param usdcPositionSize: USDC asset position of subaccount
 * @param usdcNetTransfersSinceLastPnlTick: net USDC transfers since last pnl tick
 * @param pnlTicksToBeCreatedAt: time at which new pnl tick will be created
 * @param latestBlockHeight: block height at which new pnl tick will be created
 * @param latestBlockTime: block time for above block height
 * @param mostRecentPnlTicks: most recent pnl ticks for all subaccounts
 */
// TODO(IND-126): Add support for multiple assets.
export function getNewPnlTick(
  subaccountId: string,
  subaccountTotalTransfersMap: SubaccountAssetNetTransferMap,
  marketPrices: PriceMap,
  openPerpetualPositionsForSubaccount: PerpetualPositionFromDatabase[],
  usdcPositionSize: Big,
  usdcNetTransfersSinceLastPnlTick: Big,
  pnlTicksToBeCreatedAt: DateTime,
  latestBlockHeight: string,
  latestBlockTime: IsoString,
  mostRecentPnlTicks: PnlTickForSubaccounts,
  lastUpdatedFundingIndexMap: FundingIndexMap,
  currentFundingIndexMap: FundingIndexMap,
): PnlTicksCreateObject {
  const currentEquity: Big = calculateEquity(
    usdcPositionSize,
    openPerpetualPositionsForSubaccount,
    marketPrices,
    lastUpdatedFundingIndexMap,
    currentFundingIndexMap,
  );

  const totalPnl: Big = calculateTotalPnl(
    currentEquity,
    subaccountTotalTransfersMap[subaccountId][USDC_ASSET_ID],
  );

  const mostRecentPnlTick: PnlTicksCreateObject | undefined = mostRecentPnlTicks[subaccountId];

  // if there has been a significant chagne in equity or totalPnl, log it for debugging purposes.
  if (
    mostRecentPnlTick &&
    Big(mostRecentPnlTick.equity).gt(0) &&
    currentEquity.div(mostRecentPnlTick.equity).gt(2) &&
    totalPnl.gte(10000) &&
    Big(mostRecentPnlTick.totalPnl).lt(-1000) &&
    usdcNetTransfersSinceLastPnlTick === ZERO
  ) {
    logger.info({
      at: 'createPnlTicks#getNewPnlTick',
      message: 'large change of equity and totalPnl',
      subaccountId,
      previousEquity: mostRecentPnlTick.equity,
      previousTotalPnl: mostRecentPnlTick.totalPnl,
      currentEquity: currentEquity.toFixed(),
      currentTotalPnl: totalPnl.toFixed(),
      usdcPositionSize,
      openPerpetualPositionsForSubaccount: JSON.stringify(openPerpetualPositionsForSubaccount),
      currentMarkets: marketPrices,
    });
  }

  return {
    totalPnl: totalPnl.toFixed(6),
    netTransfers: usdcNetTransfersSinceLastPnlTick.toFixed(6),
    subaccountId,
    createdAt: pnlTicksToBeCreatedAt.toISO(),
    equity: currentEquity.toFixed(6),
    blockHeight: latestBlockHeight,
    blockTime: latestBlockTime,
  };
}

/**
 * Gets a map of subaccount id to net USDC transfers between lastUpdatedHeight and blockHeight
 * @param subaccountIds: list of subaccount ids to get net USDC transfers for.
 * @param mostRecentPnlTicks: most recent pnl tick for each subaccount.
 * @param blockHeight: block height to get net USDC transfers up to.
 * @param txId: optional transaction id to use for query.
 */
export async function getUsdcTransfersSinceLastPnlTick(
  subaccountIds: string[],
  mostRecentPnlTicks: PnlTickForSubaccounts,
  blockHeight: string,
  txId?: number,
): Promise<SubaccountUsdcTransferMap> {
  const netTransfers: SubaccountUsdcTransferMap = {};
  const promises = [];
  for (const subaccountId of subaccountIds) {
    const mostRecentPnlTick: PnlTicksCreateObject | undefined = mostRecentPnlTicks[subaccountId];
    const lastUpdatedHeight: string = mostRecentPnlTick === undefined ? '0'
      : mostRecentPnlTick!.blockHeight;
    const transferPromise = TransferTable.getNetTransfersBetweenBlockHeightsForSubaccount(
      subaccountId,
      lastUpdatedHeight,
      blockHeight,
      { readReplica: true, txId },
    ).then((transfers: TransferTable.AssetTransferMap) => {
      if (USDC_ASSET_ID in transfers) {
        netTransfers[subaccountId] = transfers[USDC_ASSET_ID];
      }
    });
    promises.push(transferPromise);
  }
  await Promise.all(promises);
  return netTransfers;
}

/**
 * Calculate the current equity of a subaccount based on USDC position size and open
 * perpetual positions and any unsettled funding payments.
 * @param usdcPositionSize
 * @param positions
 * @param marketPrices
 */
// TODO(IND-226): De-duplicate this with the same function in `comlink`
export function calculateEquity(
  usdcPositionSize: Big,
  positions: PerpetualPositionFromDatabase[],
  marketPrices: PriceMap,
  lastUpdatedFundingIndexMap: FundingIndexMap,
  currentFundingIndexMap: FundingIndexMap,
): Big {
  const totalUnsettledFundingPayment: Big = positions.reduce(
    (acc: Big, position: PerpetualPositionFromDatabase) => {
      return acc.plus(
        helpers.getUnsettledFunding(
          position,
          currentFundingIndexMap,
          lastUpdatedFundingIndexMap,
        ),
      );
    },
    ZERO,
  );

  const signedPositionNotional: Big = positions.reduce(
    (acc: Big, position: PerpetualPositionFromDatabase) => {
      const positionNotional: Big = Big(position.size).times(
        marketPrices[Number(position.perpetualId)],
      );
      // Add positionNotional to the accumulator
      return acc.plus(positionNotional);
    },
    ZERO,
  );

  return signedPositionNotional.plus(usdcPositionSize).plus(totalUnsettledFundingPayment);
}

/**
 * Calculate the total pnl of a subaccount based on current equity and total historical
 * transfers.
 * @param currentEquity
 * @param totalTransfers
 */
export function calculateTotalPnl(
  currentEquity: Big,
  totalTransfers: string,
): Big {
  return currentEquity.minus(totalTransfers);
}

/**
 * Fetches the most recent pnl tick for each account from Redis. If Redis is empty,
 * fetches from the db. Redis will be empty if the service has just been started.
 */
export async function getMostRecentPnlTicksForEachAccount():
  Promise<PnlTickForSubaccounts> {
  const mostRecentCachedPnlTicks: PnlTickForSubaccounts = await
  LatestAccountPnlTicksCache.getAll(redisClient);

  return !_.isEmpty(mostRecentCachedPnlTicks)
    ? mostRecentCachedPnlTicks
    // If Redis is empty, fetch the most recent pnl ticks from the db created on or after
    // block height '1'.
    : PnlTicksTable.findMostRecentPnlTickForEachAccount(
      '1',
    );
}
