import {
  stats,
  cacheControlMiddleware,
} from '@dydxprotocol-indexer/base';
import {
  PnlTicksFromDatabase,
  perpetualMarketRefresher,
  PerpetualMarketFromDatabase,
  USDC_ASSET_ID,
  FundingIndexMap,
  AssetPositionFromDatabase,
  PerpetualPositionFromDatabase,
  SubaccountFromDatabase,
  AssetColumns,
  BlockTable,
  MarketTable,
  AssetPositionTable,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  AssetTable,
  SubaccountTable,
  AssetFromDatabase,
  MarketFromDatabase,
  BlockFromDatabase,
  FundingIndexUpdatesTable,
  PnlTickInterval,
  VaultFromDatabase,
  MEGAVAULT_SUBACCOUNT_ID,
  TransferFromDatabase,
  TransferTable,
  TransferColumns,
  Ordering,
  VaultPnlTicksView,
} from '@dydxprotocol-indexer/postgres';
import { VaultCache, CachedMegavaultPnl, CachedVaultHistoricalPnl } from '@dydxprotocol-indexer/redis';
import Big from 'big.js';
import bounds from 'binary-searching';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import _, { Dictionary } from 'lodash';
import { DateTime } from 'luxon';
import {
  Controller, Get, Query, Route,
} from 'tsoa';

import { defaultRateLimiter } from '../../../caches/rate-limiters';
import { getVaultStartPnl } from '../../../caches/vault-start-pnl';
import config from '../../../config';
import { redisClient, redisReadOnlyClient } from '../../../helpers/redis/redis-controller';
import {
  aggregateHourlyPnlTicks,
  getSubaccountResponse,
  getVaultMapping,
  getVaultPnlStartDate,
  handleControllerError,
} from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { pnlTicksToResponseObject } from '../../../request-helpers/request-transformer';
import {
  MegavaultHistoricalPnlResponse,
  VaultsHistoricalPnlResponse,
  VaultHistoricalPnl,
  VaultPosition,
  AssetById,
  MegavaultPositionResponse,
  SubaccountResponseObject,
  MegavaultHistoricalPnlRequest,
  VaultsHistoricalPnlRequest,
  AggregatedPnlTick,
  PnlTicksResponseObject,
} from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'vault-controller';
const vaultCacheControlMiddleware = cacheControlMiddleware(config.CACHE_CONTROL_DIRECTIVE_VAULTS);

interface VaultMapping {
  [subaccountId: string]: VaultFromDatabase,
}

@Route('vault/v1')
class VaultController extends Controller {
  @Get('/megavault/historicalPnl')
  async getMegavaultHistoricalPnl(
    @Query() resolution?: PnlTickInterval,
  ): Promise<MegavaultHistoricalPnlResponse> {
    const start: number = Date.now();

    // Get cache timestamp from read-only Redis
    const cacheTimestamp: Date | null = await VaultCache.getMegavaultPnlCacheTimestamp(
      getResolution(resolution),
      redisReadOnlyClient,
    );

    // Check (from read-only Redis) if the last cached result was less than the cache TTL
    if (config.VAULT_CACHE_TTL_MS > 0 &&
      cacheTimestamp !== null &&
      Date.now() - cacheTimestamp.getTime() < config.VAULT_CACHE_TTL_MS) {
      const cached: CachedMegavaultPnl | null = await VaultCache.getMegavaultPnl(
        getResolution(resolution),
        redisReadOnlyClient,
      );
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.megavault_historical_pnl_cache_hit.timing`,
        Date.now() - start,
        {
          resolution: getResolution(resolution),
        },
      );

      if (cached !== null) {
        stats.increment(
          `${config.SERVICE_NAME}.${controllerName}.megavault_historical_pnl_cache_hit`,
          {
            resolution: getResolution(resolution),
          },
        );

        return {
          megavaultPnl: _.sortBy(cached.pnlTicks, 'blockTime'),
        };
      }
    }

    stats.increment(
      `${config.SERVICE_NAME}.${controllerName}.megavault_historical_pnl_cache_miss`,
      {
        resolution: getResolution(resolution),
      },
    );

    const vaultSubaccounts: VaultMapping = await getVaultMapping();
    stats.timing(
      `${config.SERVICE_NAME}.${controllerName}.fetch_vaults.timing`,
      Date.now() - start,
    );

    const startTicksPositions: number = Date.now();
    const vaultSubaccountIdsWithMainSubaccount: string[] = _
      .keys(vaultSubaccounts)
      .concat([MEGAVAULT_SUBACCOUNT_ID]);
    const [
      vaultPnlTicks,
      vaultPositions,
      latestBlock,
      mainSubaccountEquity,
      latestPnlTick,
      firstMainVaultTransferTimestamp,
    ] : [
      PnlTicksFromDatabase[],
      Map<string, VaultPosition>,
      BlockFromDatabase,
      string,
      PnlTicksFromDatabase | undefined,
      DateTime | undefined,
    ] = await Promise.all([
      getVaultSubaccountPnlTicks(vaultSubaccountIdsWithMainSubaccount, getResolution(resolution)),
      getVaultPositions(vaultSubaccounts),
      BlockTable.getLatest(),
      getMainSubaccountEquity(),
      getLatestPnlTick(_.values(vaultSubaccounts)),
      getFirstMainVaultTransferDateTime(),
    ]);
    stats.timing(
      `${config.SERVICE_NAME}.${controllerName}.fetch_ticks_positions_equity.timing`,
      Date.now() - startTicksPositions,
    );
    // aggregate pnlTicks for all vault subaccounts grouped by blockHeight
    const aggregatedPnlTicks: PnlTicksFromDatabase[] = aggregateVaultPnlTicks(
      vaultPnlTicks,
      _.values(vaultSubaccounts),
      firstMainVaultTransferTimestamp,
    );

    const currentEquity: string = Array.from(vaultPositions.values())
      .map((position: VaultPosition): string => {
        return position.equity;
      }).reduce((acc: string, curr: string): string => {
        return (Big(acc).add(Big(curr))).toFixed();
      }, mainSubaccountEquity);
    const pnlTicksWithCurrentTick: PnlTicksFromDatabase[] = getPnlTicksWithCurrentTick(
      currentEquity,
      filterOutIntervalTicks(aggregatedPnlTicks, getResolution(resolution)),
      latestBlock,
      latestPnlTick,
    );

    const sortedPnlTicks: PnlTicksResponseObject[] = _.sortBy(pnlTicksWithCurrentTick, 'blockTime')
      .map(
        (pnlTick: PnlTicksFromDatabase) => {
          return pnlTicksToResponseObject(pnlTick);
        });

    // Cache in primary Redis instance
    await VaultCache.setMegavaultPnl(
      getResolution(resolution),
      sortedPnlTicks,
      redisClient,
    );

    stats.timing(
      `${config.SERVICE_NAME}.${controllerName}.megavault_historical_pnl.cache_miss.timing`,
      Date.now() - start,
      {
        resolution: getResolution(resolution),
      },
    );

    return {
      megavaultPnl: sortedPnlTicks,
    };
  }

  @Get('/vaults/historicalPnl')
  async getVaultsHistoricalPnl(
    @Query() resolution?: PnlTickInterval,
  ): Promise<VaultsHistoricalPnlResponse> {
    const start: number = Date.now();

    // Get cache timestamp from read-only Redis
    const cacheTimestamp: Date | null = await VaultCache.getVaultsHistoricalPnlCacheTimestamp(
      getResolution(resolution),
      redisReadOnlyClient,
    );
    // Check if the last cached result was less than the cache TTL
    if (config.VAULT_CACHE_TTL_MS > 0 &&
      cacheTimestamp !== null &&
      Date.now() - cacheTimestamp.getTime() < config.VAULT_CACHE_TTL_MS) {
      const cached: CachedVaultHistoricalPnl[] | null = await VaultCache.getVaultsHistoricalPnl(
        getResolution(resolution),
        redisReadOnlyClient,
      );

      if (cached !== null) {
        stats.timing(
          `${config.SERVICE_NAME}.${controllerName}.vaults_historical_pnl_cache_hit.timing`,
          Date.now() - start,
          {
            resolution: getResolution(resolution),
          },
        );

        stats.increment(
          `${config.SERVICE_NAME}.${controllerName}.vaults_historical_pnl_cache_hit`,
          {
            resolution: getResolution(resolution),
          },
        );

        return {
          vaultsPnl: cached,
        };
      }
    }

    stats.increment(
      `${config.SERVICE_NAME}.${controllerName}.vaults_historical_pnl_cache_miss`,
      {
        resolution: getResolution(resolution),
      },
    );

    const vaultSubaccounts: VaultMapping = await getVaultMapping();
    const [
      vaultPnlTicks,
      vaultPositions,
      latestBlock,
      latestTicks,
    ] : [
      PnlTicksFromDatabase[],
      Map<string, VaultPosition>,
      BlockFromDatabase,
      PnlTicksFromDatabase[],
    ] = await Promise.all([
      getVaultSubaccountPnlTicks(_.keys(vaultSubaccounts), getResolution(resolution)),
      getVaultPositions(vaultSubaccounts),
      BlockTable.getLatest(),
      getLatestPnlTicks(),
    ]);
    const latestTicksBySubaccountId: Dictionary<PnlTicksFromDatabase> = _.keyBy(
      latestTicks,
      'subaccountId',
    );

    const groupedVaultPnlTicks: VaultHistoricalPnl[] = _(vaultPnlTicks)
      .filter((pnlTickFromDatabsae: PnlTicksFromDatabase): boolean => {
        return vaultSubaccounts[pnlTickFromDatabsae.subaccountId] !== undefined;
      })
      .groupBy('subaccountId')
      .mapValues((pnlTicks: PnlTicksFromDatabase[], subaccountId: string): VaultHistoricalPnl => {
        const market: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
          .getPerpetualMarketFromClobPairId(
            vaultSubaccounts[subaccountId].clobPairId,
          );

        if (market === undefined) {
          throw new Error(
            `Vault clob pair id ${vaultSubaccounts[subaccountId]} does not correspond to ` +
            'a perpetual market.');
        }

        const vaultPosition: VaultPosition | undefined = vaultPositions.get(subaccountId);
        const currentEquity: string = vaultPosition === undefined ? '0' : vaultPosition.equity;
        const pnlTicksWithCurrentTick: PnlTicksFromDatabase[] = getPnlTicksWithCurrentTick(
          currentEquity,
          pnlTicks,
          latestBlock,
          latestTicksBySubaccountId[subaccountId],
        );

        // Only retain fields we need and excludes fields like `id`, `subaccountId`.
        const pnlTicksResponseObjects
        : PnlTicksResponseObject[] = pnlTicksWithCurrentTick.map(pnlTicksToResponseObject);
        return {
          ticker: market.ticker,
          historicalPnl: pnlTicksResponseObjects,
        };
      })
      .values()
      .value();

    const sortedVaultPnlTicks: VaultHistoricalPnl[] = _.sortBy(groupedVaultPnlTicks, 'ticker');

    stats.timing(
      `${config.SERVICE_NAME}.${controllerName}.vaults_historical_pnl_cache_miss.timing`,
      Date.now() - start,
      {
        resolution: getResolution(resolution),
      },
    );

    // Cache in primary Redis instance
    await VaultCache.setVaultsHistoricalPnl(
      getResolution(resolution),
      sortedVaultPnlTicks,
      redisClient,
    );

    return {
      vaultsPnl: sortedVaultPnlTicks,
    };
  }

  @Get('/megavault/positions')
  async getMegavaultPositions(): Promise<MegavaultPositionResponse> {
    const vaultSubaccounts: VaultMapping = await getVaultMapping();

    const vaultPositions: Map<string, VaultPosition> = await getVaultPositions(vaultSubaccounts);

    return {
      positions: _.sortBy(Array.from(vaultPositions.values()), 'ticker'),
    };
  }
}

router.get(
  '/v1/megavault/historicalPnl',
  rateLimiterMiddleware(defaultRateLimiter),
  vaultCacheControlMiddleware,
  ...checkSchema({
    resolution: {
      in: 'query',
      isIn: {
        options: [Object.values(PnlTickInterval)],
        errorMessage: `type must be one of ${Object.values(PnlTickInterval)}`,
      },
      optional: true,
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      resolution,
    }: MegavaultHistoricalPnlRequest = matchedData(req) as MegavaultHistoricalPnlRequest;

    try {
      const controllers: VaultController = new VaultController();
      const response: MegavaultHistoricalPnlResponse = await controllers
        .getMegavaultHistoricalPnl(resolution);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultController GET /megavault/historicalPnl',
        'Megavault Historical Pnl error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_megavault_historical_pnl.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/v1/vaults/historicalPnl',
  rateLimiterMiddleware(defaultRateLimiter),
  vaultCacheControlMiddleware,
  ...checkSchema({
    resolution: {
      in: 'query',
      isIn: {
        options: [Object.values(PnlTickInterval)],
        errorMessage: `type must be one of ${Object.values(PnlTickInterval)}`,
      },
      optional: true,
    },
  }),
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const {
      resolution,
    }: VaultsHistoricalPnlRequest = matchedData(req) as VaultsHistoricalPnlRequest;

    try {
      const controllers: VaultController = new VaultController();
      const response: VaultsHistoricalPnlResponse = await controllers
        .getVaultsHistoricalPnl(resolution);
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultHistoricalPnlController GET /vaults/historicalPnl',
        'Vaults Historical Pnl error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_vaults_historical_pnl.timing`,
        Date.now() - start,
      );
    }
  },
);

router.get(
  '/v1/megavault/positions',
  rateLimiterMiddleware(defaultRateLimiter),
  vaultCacheControlMiddleware,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    try {
      const controllers: VaultController = new VaultController();
      const response: MegavaultPositionResponse = await controllers.getMegavaultPositions();
      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'VaultController GET /megavault/positions',
        'Megavault Positions error',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_megavault_positions.timing`,
        Date.now() - start,
      );
    }
  },
);

async function getVaultSubaccountPnlTicks(
  vaultSubaccountIds: string[],
  resolution: PnlTickInterval,
): Promise<PnlTicksFromDatabase[]> {
  if (vaultSubaccountIds.length === 0) {
    return [];
  }

  let windowSeconds: number;
  if (resolution === PnlTickInterval.day) {
    windowSeconds = config.VAULT_PNL_HISTORY_DAYS * 24 * 60 * 60; // days to seconds
  } else {
    windowSeconds = config.VAULT_PNL_HISTORY_HOURS * 60 * 60; // hours to seconds
  }

  const pnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
    resolution,
    windowSeconds,
    getVaultPnlStartDate(),
  );

  return adjustVaultPnlTicks(pnlTicks, getVaultStartPnl());
}

async function getVaultPositions(
  vaultSubaccounts: VaultMapping,
): Promise<Map<string, VaultPosition>> {
  const start: number = Date.now();
  const vaultSubaccountIds: string[] = _.keys(vaultSubaccounts);
  if (vaultSubaccountIds.length === 0) {
    return new Map();
  }

  const [
    subaccounts,
    assets,
    openPerpetualPositions,
    assetPositions,
    markets,
    latestBlock,
  ]: [
    SubaccountFromDatabase[],
    AssetFromDatabase[],
    PerpetualPositionFromDatabase[],
    AssetPositionFromDatabase[],
    MarketFromDatabase[],
    BlockFromDatabase | undefined,
  ] = await Promise.all([
    SubaccountTable.findAll(
      {
        id: vaultSubaccountIds,
      },
      [],
    ),
    AssetTable.findAll(
      {},
      [],
    ),
    PerpetualPositionTable.findAll(
      {
        subaccountId: vaultSubaccountIds,
        status: [PerpetualPositionStatus.OPEN],
      },
      [],
    ),
    AssetPositionTable.findAll(
      {
        subaccountId: vaultSubaccountIds,
        assetId: [USDC_ASSET_ID],
      },
      [],
    ),
    MarketTable.findAll(
      {},
      [],
    ),
    BlockTable.getLatest(),
  ]);
  stats.timing(
    `${config.SERVICE_NAME}.${controllerName}.positions.fetch_subaccounts_positions.timing`,
    Date.now() - start,
  );

  const startFunding: number = Date.now();
  const updatedAtHeights: string[] = _(subaccounts).map('updatedAtHeight').uniq().value();
  const [
    latestFundingIndexMap,
    fundingIndexMaps,
  ]: [
    FundingIndexMap,
    {[blockHeight: string]: FundingIndexMap}
  ] = await Promise.all([
    FundingIndexUpdatesTable
      .findFundingIndexMap(
        latestBlock.blockHeight,
      ),
    getFundingIndexMapsChunked(updatedAtHeights),
  ]);
  stats.timing(
    `${config.SERVICE_NAME}.${controllerName}.positions.fetch_funding.timing`,
    Date.now() - startFunding,
  );

  const assetPositionsBySubaccount:
  { [subaccountId: string]: AssetPositionFromDatabase[] } = _.groupBy(
    assetPositions,
    'subaccountId',
  );
  const openPerpetualPositionsBySubaccount:
  { [subaccountId: string]: PerpetualPositionFromDatabase[] } = _.groupBy(
    openPerpetualPositions,
    'subaccountId',
  );
  const assetIdToAsset: AssetById = _.keyBy(
    assets,
    AssetColumns.id,
  );

  const vaultPositionsAndSubaccountId: {
    position: VaultPosition,
    subaccountId: string,
  }[] = subaccounts.map((subaccount: SubaccountFromDatabase) => {
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(vaultSubaccounts[subaccount.id].clobPairId);
    if (perpetualMarket === undefined) {
      throw new Error(
        `Vault clob pair id ${vaultSubaccounts[subaccount.id]} does not correspond to a ` +
          'perpetual market.');
    }
    const lastUpdatedFundingIndexMap: FundingIndexMap = fundingIndexMaps[
      subaccount.updatedAtHeight
    ];
    if (lastUpdatedFundingIndexMap === undefined) {
      throw new Error(
        `No funding indices could be found for vault with subaccount ${subaccount.id}`,
      );
    }

    const subaccountResponse: SubaccountResponseObject = getSubaccountResponse(
      subaccount,
      openPerpetualPositionsBySubaccount[subaccount.id] || [],
      assetPositionsBySubaccount[subaccount.id] || [],
      assets,
      markets,
      perpetualMarketRefresher.getPerpetualMarketsMap(),
      latestBlock.blockHeight,
      latestFundingIndexMap,
      lastUpdatedFundingIndexMap,
    );

    return {
      position: {
        ticker: perpetualMarket.ticker,
        assetPosition: subaccountResponse.assetPositions[
          assetIdToAsset[USDC_ASSET_ID].symbol
        ],
        perpetualPosition: subaccountResponse.openPerpetualPositions[
          perpetualMarket.ticker
        ] || undefined,
        equity: subaccountResponse.equity,
      },
      subaccountId: subaccount.id,
    };
  });

  return new Map(vaultPositionsAndSubaccountId.map(
    (obj: { position: VaultPosition, subaccountId: string }) : [string, VaultPosition] => {
      return [
        obj.subaccountId,
        obj.position,
      ];
    },
  ));
}

async function getMainSubaccountEquity(): Promise<string> {
  // Main vault subaccount should only ever hold a USDC and never any perpetuals.
  const usdcBalance: {[subaccountId: string]: Big} = await AssetPositionTable
    .findUsdcPositionForSubaccounts(
      [MEGAVAULT_SUBACCOUNT_ID],
    );
  return usdcBalance[MEGAVAULT_SUBACCOUNT_ID]?.toFixed() || '0';
}

function getPnlTicksWithCurrentTick(
  equity: string,
  pnlTicks: PnlTicksFromDatabase[],
  latestBlock: BlockFromDatabase,
  latestTick: PnlTicksFromDatabase | undefined = undefined,
): PnlTicksFromDatabase[] {
  if (latestTick !== undefined) {
    return pnlTicks.concat({
      ...latestTick,
      equity,
      blockHeight: latestBlock.blockHeight,
      blockTime: latestBlock.time,
      createdAt: latestBlock.time,
    });
  }
  if (pnlTicks.length === 0) {
    return [];
  }
  const currentTick: PnlTicksFromDatabase = {
    ...(_.maxBy(pnlTicks, 'blockTime')!),
    equity,
    blockHeight: latestBlock.blockHeight,
    blockTime: latestBlock.time,
    createdAt: latestBlock.time,
  };
  return pnlTicks.concat([currentTick]);
}

export async function getLatestPnlTicks(): Promise<PnlTicksFromDatabase[]> {
  const latestPnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getLatestVaultPnl();
  const adjustedPnlTicks: PnlTicksFromDatabase[] = adjustVaultPnlTicks(
    latestPnlTicks,
    getVaultStartPnl(),
  );
  return adjustedPnlTicks;
}

export async function getLatestPnlTick(
  vaults: VaultFromDatabase[],
): Promise<PnlTicksFromDatabase | undefined> {
  const pnlTicks: PnlTicksFromDatabase[] = await VaultPnlTicksView.getVaultsPnl(
    PnlTickInterval.hour,
    config.VAULT_LATEST_PNL_TICK_WINDOW_HOURS * 60 * 60,
    getVaultPnlStartDate(),
  );
  const adjustedPnlTicks: PnlTicksFromDatabase[] = adjustVaultPnlTicks(
    pnlTicks,
    getVaultStartPnl(),
  );
  // Aggregate and get pnl tick closest to the hour
  const aggregatedTicks: PnlTicksFromDatabase[] = aggregateVaultPnlTicks(
    adjustedPnlTicks,
    vaults,
  );
  const filteredTicks: PnlTicksFromDatabase[] = filterOutIntervalTicks(
    aggregatedTicks,
    PnlTickInterval.hour,
  );
  return _.maxBy(filteredTicks, 'blockTime');
}

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

function getResolution(resolution: PnlTickInterval = PnlTickInterval.day): PnlTickInterval {
  return resolution;
}

/**
 * Gets funding index maps in a chunked fashion to reduce database load and aggregates into a
 * a map of funding index maps.
 * @param updatedAtHeights
 * @returns
 */
async function getFundingIndexMapsChunked(
  updatedAtHeights: string[],
): Promise<{[blockHeight: string]: FundingIndexMap}> {
  const updatedAtHeightsNum: number[] = updatedAtHeights.map((height: string): number => {
    return parseInt(height, 10);
  }).sort();
  const aggregateFundingIndexMaps: {[blockHeight: string]: FundingIndexMap} = {};
  await Promise.all(getHeightWindows(updatedAtHeightsNum).map(
    async (heightWindow: number[]): Promise<void> => {
      const fundingIndexMaps: {[blockHeight: string]: FundingIndexMap} = await
      FundingIndexUpdatesTable
        .findFundingIndexMaps(
          heightWindow.map((heightNum: number): string => { return heightNum.toString(); }),
        );
      for (const height of _.keys(fundingIndexMaps)) {
        aggregateFundingIndexMaps[height] = fundingIndexMaps[height];
      }
    }));
  return aggregateFundingIndexMaps;
}

/**
 * Separates an array of heights into a chunks based on a window size. Each chunk should only
 * contain heights within a certain number of blocks of each other.
 * @param heights
 * @returns
 */
function getHeightWindows(
  heights: number[],
): number[][] {
  if (heights.length === 0) {
    return [];
  }
  const windows: number[][] = [];
  let windowStart: number = heights[0];
  let currentWindow: number[] = [];
  for (const height of heights) {
    if (height - windowStart < config.VAULT_FETCH_FUNDING_INDEX_BLOCK_WINDOWS) {
      currentWindow.push(height);
    } else {
      windows.push(currentWindow);
      currentWindow = [height];
      windowStart = height;
    }
  }
  windows.push(currentWindow);
  return windows;
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
    (createdAt: string) => { return DateTime.fromISO(createdAt); },
  ).concat(
    mainVaultCreatedAt === undefined ? [] : [mainVaultCreatedAt],
  ).sort(
    (a: DateTime, b: DateTime) => {
      return a.diff(b).milliseconds;
    },
  );
  return aggregatedPnlTicks.filter((aggregatedTick: AggregatedPnlTick) => {
    // Get number of vaults created before the pnl tick was created by binary-searching for the
    // index of the pnl ticks createdAt in a sorted array of vault createdAt times.
    const numVaultsCreated: number = bounds.le(
      vaultCreationTimes,
      DateTime.fromISO(aggregatedTick.pnlTick.createdAt),
      (a: DateTime, b: DateTime) => { return a.diff(b).milliseconds; },
    );
    // Number of ticks should be greater than number of vaults created before it as there should be
    // a tick for the main vault subaccount.
    return aggregatedTick.numTicks >= numVaultsCreated;
  }).map((aggregatedPnlTick: AggregatedPnlTick) => { return aggregatedPnlTick.pnlTick; });
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

export default router;
