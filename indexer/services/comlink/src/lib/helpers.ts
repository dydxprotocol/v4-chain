import { logger } from '@dydxprotocol-indexer/base';
import {
  AssetPositionFromDatabase,
  BlockFromDatabase,
  CHILD_SUBACCOUNT_MULTIPLIER,
  FundingIndexMap,
  FundingIndexUpdatesTable,
  helpers,
  liquidityTierRefresher,
  LiquidityTiersFromDatabase,
  MarketFromDatabase,
  MarketsMap,
  MAX_PARENT_SUBACCOUNTS,
  PerpetualMarketFromDatabase,
  PerpetualMarketsMap,
  PerpetualMarketTable,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PnlTicksFromDatabase,
  PositionSide,
  SubaccountFromDatabase,
  SubaccountTable,
  USDC_SYMBOL,
  AssetFromDatabase,
  AssetColumns,
  MarketColumns,
  VaultFromDatabase, VaultTable, perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import express from 'express';
import _ from 'lodash';
import { DateTime } from 'luxon';

import config from '../config';
import {
  assetPositionToResponseObject,
  perpetualPositionToResponseObject,
  subaccountToResponseObject,
} from '../request-helpers/request-transformer';
import {
  AggregatedPnlTick,
  AssetById,
  AssetPositionResponseObject,
  AssetPositionsMap,
  MarketType,
  PerpetualPositionResponseObject,
  PerpetualPositionsMap,
  PerpetualPositionWithFunding,
  Risk,
  SubaccountResponseObject,
  VaultMapping,
} from '../types';
import { ZERO, ZERO_USDC_POSITION } from './constants';
import { InvalidParamError, NotFoundError, TurnkeyError } from './errors';

/* ------- GENERIC HELPERS ------- */

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function isDefined(val: any): boolean {
  return val !== undefined && val !== null;
}

/* ------- ERROR HELPERS ------- */

export function handleControllerError(
  at: string,
  message: string,
  error: Error,
  req: express.Request,
  res: express.Response,
): express.Response {
  if (error instanceof NotFoundError) {
    return handleNotFoundError(error.message, res);
  }
  if (error instanceof InvalidParamError) {
    return handleInvalidParamError(error.message, res);
  }
  if (error instanceof TurnkeyError) {
    return handleTurnkeyError(error, res);
  }
  return handleInternalServerError(
    at,
    message,
    error,
    req,
    res,
  );
}

export function handleInternalServerError(
  at: string,
  message: string,
  error: Error,
  req: express.Request,
  res: express.Response,
): express.Response {
  if (config.isDevelopment()) {
    // eslint-disable-next-line no-console
    console.error(error);
  }

  logger.error({
    at,
    message,
    error,
    stacktrace: error.stack,
    params: JSON.stringify(req.params),
    query: JSON.stringify(req.query),
  });
  return createInternalServerErrorResponse(res);
}

function handleInvalidParamError(
  message: string,
  res: express.Response,
): express.Response {
  return res.status(400).json({
    errors: [{
      msg: message,
    }],
  });
}

function handleNotFoundError(
  message: string,
  res: express.Response,
): express.Response {
  return res.status(404).json({
    errors: [{
      msg: message,
    }],
  });
}

function handleTurnkeyError(
  error: TurnkeyError,
  res: express.Response,
): express.Response {
  return res.status(400).json({
    errors: [{
      msg: error.message,
      type: 'TURNKEY_ERROR',
    }],
  });
}

export function createInternalServerErrorResponse(
  res: express.Response,
): express.Response {
  return res.status(500).json({
    errors: [{
      msg: 'Internal Server Error',
    }],
  });
}

export function create4xxResponse(
  res: express.Response,
  msg: string,
  status: number = 400,
  additionalParams: object = {},
): express.Response {
  return res.status(status).json({
    errors: [{
      ...additionalParams,
      msg,
    }],
  });
}

/* ------- MARKET HELPERS ------- */

export async function getClobPairId(
  market: string,
  marketType: MarketType,
): Promise<string | undefined> {
  if (marketType === MarketType.PERPETUAL) {
    const perpetualMarket: (
      PerpetualMarketFromDatabase | undefined
    ) = await PerpetualMarketTable.findByTicker(market);

    if (perpetualMarket !== undefined) {
      return perpetualMarket.clobPairId;
    }
  }
  // spot markets are not supported in V4 yet

  return undefined;
}

/* ------- ACCOUNT HELPERS ------- */

/**
 * Calculate the equity and free collateral for a subaccount given all the positions for it.
 * 1. Equity for a subaccount is sum of the notional value of all the positions held by the
 * subaccount and the USDC asset position of the subaccount.
 * 2. Free collateral for a subaccount is the sum of the initial margin required for all positions
 * held by the subaccount subtracted from the equity.
 * @param perpetualPositions List of positions in perpetual markets held by the subaccount
 * @param perpetualMarketsMap Map of perpetual ids to perpetual markets
 * @returns Equity and free collateral of the subaccount
 */
// TODO(DEC-1257): take into account all asset positions when calculating collateral
export function calculateEquityAndFreeCollateral(
  perpetualPositions: PerpetualPositionFromDatabase[],
  perpetualMarketsMap: PerpetualMarketsMap,
  marketsMap: MarketsMap,
  usdcPositionSize: string,
): {
  equity: string,
  freeCollateral: string,
} {
  const {
    signedPositionNotional,
    totalPositionRisk,
  }: {
    signedPositionNotional: Big,
    totalPositionRisk: Big,
  } = perpetualPositions.reduce((acc, position) => {
    // get the positionNotional for each position and the individualRisk of the position
    const {
      signedNotional,
      individualRisk,
    }: {
      signedNotional: Big,
      individualRisk: Risk,
    } = getSignedNotionalAndRisk({
      size: new Big(position.size),
      perpetualMarket: perpetualMarketsMap[position.perpetualId],
      market: marketsMap[perpetualMarketsMap[position.perpetualId].marketId],
    });

    // Add positionNotional and totalPositionRisk to the accumulator
    acc.signedPositionNotional = acc.signedPositionNotional.plus(signedNotional);
    acc.totalPositionRisk = acc.totalPositionRisk.plus(individualRisk.initial);
    return acc;
  },
  {
    signedPositionNotional: ZERO,
    totalPositionRisk: ZERO,
  },
  );

  // Derive equity and freeCollateral of the account from the PositionNotional
  // and totalPositionRisk of positions
  const equity: Big = signedPositionNotional.plus(usdcPositionSize);
  return {
    equity: equity.toFixed(),
    freeCollateral: equity.minus(totalPositionRisk).toFixed(),
  };
}

/* ------- POSITION HELPERS ------- */

/**
 * Calculates the notional value and risk of a perpetual position in a given market.
 * 1. Notional value of the position is the size of the position multiplied by the index price of
 * the market.
 * 2. Risk consists of the initial margin required for the position and the maintenance margin
 * required for the position. This is calculated using the notional value of the position multiplied
 * by the initial and maintenance margin fractions of the market the position is in.
 * @param param0 Object containing the size of the position and the perpetual market.
 * @returns Notional value of the position and the risk of the position.
 */
export function getSignedNotionalAndRisk({
  size,
  perpetualMarket,
  market,
}: {
  size: Big,
  perpetualMarket: PerpetualMarketFromDatabase,
  market: MarketFromDatabase,
}): {
  signedNotional: Big,
  individualRisk: Risk,
} {
  const signedNotional: Big = size.times(market.oraclePrice!);
  const liquidityTier:
  LiquidityTiersFromDatabase | undefined = liquidityTierRefresher.getLiquidityTierFromId(
    perpetualMarket.liquidityTierId,
  );
  if (liquidityTier === undefined) {
    throw new NotFoundError(`Liquidity tier with id ${perpetualMarket.liquidityTierId} not found for perpetual market ${perpetualMarket.ticker}`);
  }
  // Used to calculate risk / margin fracitons, as risk of a position should always be positive
  const positionNotional: Big = signedNotional.abs();
  const {
    initialMarginFraction,
    maintenanceMarginFraction,
  }: {
    initialMarginFraction: Big,
    maintenanceMarginFraction: Big,
  } = getMarginFractions(liquidityTier);
  return {
    signedNotional,
    individualRisk: {
      initial: positionNotional.times(initialMarginFraction),
      maintenance: positionNotional.times(maintenanceMarginFraction),
    },
  };
}

export function getMarginFractions(liquidityTier: LiquidityTiersFromDatabase): {
  initialMarginFraction: Big,
  maintenanceMarginFraction: Big,
} {
  const initialMarginFraction: Big = getMarginFraction({
    liquidityTier,
    initial: true,
  });
  const maintenanceMarginFraction: Big = getMarginFraction({
    liquidityTier,
    initial: false,
  });
  return {
    initialMarginFraction,
    maintenanceMarginFraction,
  };
}

/**
 * Get the margin fraction for a position in a given perpetual market.
 *
 * @param liquidityTier The liquidity tier of the position.
 * @param initial Whether to compute the initial margin fraction or the maintenance margin fraction.
 *
 * @returns The margin fraction for the position in human-readable form.
 */
export function getMarginFraction(
  {
    liquidityTier,
    initial,
  }: {
    liquidityTier: LiquidityTiersFromDatabase,
    initial: boolean,
  },
): Big {

  const margin: string = initial ? helpers.ppmToString(Number(liquidityTier.initialMarginPpm))
    : helpers.ppmToString(
      helpers.getMaintenanceMarginPpm(
        Number(liquidityTier.initialMarginPpm),
        Number(liquidityTier.maintenanceFractionPpm),
      ),
    );
  return Big(margin);
}

/**
 * Filter out asset positions of 0 size.
 *
 * @param assetPositions
 */
export function filterAssetPositions(assetPositions: AssetPositionFromDatabase[]):
  AssetPositionFromDatabase[] {
  return assetPositions.filter((elem: AssetPositionFromDatabase) => {
    return !Big(elem.size).eq(ZERO);
  });
}

/**
 * De-depulicate a list of perpetual positions by the perpetual id of the position, keeping the
 * position with the latest `lastEventId` for the perpetual id. Perpetual positions will be ordered
 * chronologically in descending order by the last event id in the returned list.
 * @param positions List of perpetual positions with funding.
 * @returns De-duplicated list of perpetual positions. Positions will be ordered in descending
 * chronological order by the last event id of the position.
 */
export function filterPositionsByLatestEventIdPerPerpetual(
  positions: PerpetualPositionWithFunding[],
): PerpetualPositionWithFunding[] {
  const sortedPositionsArray: PerpetualPositionWithFunding[] = positions.sort(
    (a: PerpetualPositionWithFunding, b: PerpetualPositionWithFunding): number => {
      // eventId is a 96 bit value, pad both hex-strings to (96/4) = 24 hex chars
      const eventAHex: string = a.lastEventId.toString('hex').padStart(24, '0');
      const eventBHex: string = b.lastEventId.toString('hex').padStart(24, '0');
      return eventBHex.localeCompare(eventAHex);
    },
  );

  // NOTE: A subaccount should only have one open position per perpetual market, this de-duplication
  // will ensure that this invariant is true.
  // TODO(DEC-698): Remove this if deemed unecessary after e2e testing of Indexer.
  return _.uniqBy(sortedPositionsArray, 'perpetualId');
}

/**
 * Get the last updated funding index map and the latest funding index map given a subaccount
 * and the latest block
 * @param subaccount
 * @param latestBlock
 * @returns
 */
export async function getFundingIndexMaps(
  subaccount: SubaccountFromDatabase,
  latestBlock: BlockFromDatabase,
): Promise<{
  lastUpdatedFundingIndexMap: FundingIndexMap,
  latestFundingIndexMap: FundingIndexMap,
}> {
  const [lastUpdatedFundingIndexMap, latestFundingIndexMap]:
  [FundingIndexMap, FundingIndexMap] = await Promise.all([
    FundingIndexUpdatesTable.findFundingIndexMap(
      subaccount.updatedAtHeight,
    ),
    FundingIndexUpdatesTable.findFundingIndexMap(
      latestBlock.blockHeight,
    ),
  ]);
  return {
    lastUpdatedFundingIndexMap,
    latestFundingIndexMap,
  };
}

/**
 * Compute the total unsettled funding across a set of perpetual positions given the last updated
 * funding indexes and the latest funding indexes
 * @param perpetualPositions
 * @param lastUpdatedFundingIndexes
 * @param latestFundingIndexes
 * @returns
 */
export function getTotalUnsettledFunding(
  perpetualPositions: PerpetualPositionFromDatabase[],
  latestFundingIndexes: FundingIndexMap,
  lastUpdatedFundingIndexes: FundingIndexMap,
): Big {
  return _.reduce(
    perpetualPositions,
    (acc: Big, perpetualPosition: PerpetualPositionFromDatabase): Big => {
      return acc.plus(
        helpers.getUnsettledFunding(
          perpetualPosition,
          latestFundingIndexes,
          lastUpdatedFundingIndexes,
        ),
      );
    },
    ZERO,
  );
}

/**
 * Gets and adjusts the USDC asset position within a map of AssetPositions given the unsettled
 * funding
 * @param assetPositionsMap
 * @param unsettledFunding
 * @returns
 */
export function adjustUSDCAssetPosition(
  assetPositionsMap: AssetPositionsMap,
  unsettledFunding: Big,
): {
  assetPositionsMap: AssetPositionsMap,
  adjustedUSDCAssetPositionSize: string,
} {
  let adjustedAssetPositionsMap: AssetPositionsMap = _.cloneDeep(assetPositionsMap);
  const usdcPosition: AssetPositionResponseObject = _.get(assetPositionsMap, USDC_SYMBOL);
  let signedUsdcPositionSize: Big;
  if (usdcPosition?.size !== undefined) {
    signedUsdcPositionSize = Big(
      usdcPosition.side === PositionSide.LONG
        ? usdcPosition.size
        : -usdcPosition.size,
    );
  } else {
    signedUsdcPositionSize = ZERO;
  }
  const adjustedSize: Big = signedUsdcPositionSize.plus(unsettledFunding);
  // Update the USDC position in the map if the adjusted size is non-zero
  if (!adjustedSize.eq(ZERO)) {
    _.set(
      adjustedAssetPositionsMap,
      USDC_SYMBOL,
      getUSDCAssetPosition(adjustedSize,
        adjustedAssetPositionsMap[USDC_SYMBOL]?.subaccountNumber ?? 0),
    );
    // Remove the USDC position in the map if the adjusted size is zero
  } else {
    adjustedAssetPositionsMap = _.omit(adjustedAssetPositionsMap, USDC_SYMBOL);
  }

  return {
    assetPositionsMap: adjustedAssetPositionsMap,
    adjustedUSDCAssetPositionSize: adjustedSize.toFixed(),
  };
}

function getUSDCAssetPosition(signedSize: Big, subaccountNumber: number):
    AssetPositionResponseObject {
  const side: PositionSide = signedSize.gt(ZERO) ? PositionSide.LONG : PositionSide.SHORT;
  return {
    ...ZERO_USDC_POSITION,
    side,
    size: signedSize.abs().toFixed(),
    subaccountNumber,
  };
}

export function getPerpetualPositionsWithUpdatedFunding(
  positions: PerpetualPositionWithFunding[],
  latestFundingIndexMap: FundingIndexMap,
  lastUpdatedFundingIndexMap: FundingIndexMap,
): PerpetualPositionWithFunding[] {
  return _.map(
    positions,
    (position: PerpetualPositionWithFunding): PerpetualPositionWithFunding => {
      const clonedPosition: PerpetualPositionWithFunding = _.cloneDeep(position);
      // Positions that are not open have 0 unsettled funding
      if (position.status !== PerpetualPositionStatus.OPEN) {
        clonedPosition.unsettledFunding = '0';
        return clonedPosition;
      }

      const unsettledFunding: Big = helpers.getUnsettledFunding(
        position,
        latestFundingIndexMap,
        lastUpdatedFundingIndexMap,
      );
      clonedPosition.unsettledFunding = unsettledFunding.toFixed();
      return clonedPosition;
    },
  );
}

export function initializePerpetualPositionsWithFunding(
  perpetualPositions: PerpetualPositionFromDatabase[],
): PerpetualPositionWithFunding[] {
  return perpetualPositions.map((pos: PerpetualPositionFromDatabase) => {
    return {
      ...pos,
      unsettledFunding: '0',
    };
  });
}

/* ------- PARENT/CHILD SUBACCOUNT HELPERS ------- */

/**
 * Gets a list of all possible child subaccount numbers for a parent subaccount number
 * Child subaccounts = [128*0+parentSubaccount, 128*1+parentSubaccount ... 128*999+parentSubaccount]
 * @param parentSubaccount
 * @returns
 */
export function getChildSubaccountNums(parentSubaccountNum: number): number[] {
  if (parentSubaccountNum >= MAX_PARENT_SUBACCOUNTS) {
    throw new NotFoundError(`Parent subaccount number must be less than ${MAX_PARENT_SUBACCOUNTS}`);
  }
  return Array.from({ length: CHILD_SUBACCOUNT_MULTIPLIER },
    // eslint-disable-next-line @typescript-eslint/no-shadow
    (_, i) => MAX_PARENT_SUBACCOUNTS * i + parentSubaccountNum);
}

/**
 * Gets the subaccount uuids of all the child subaccounts given a parent subaccount number
 * @param address
 * @param parentSubaccountNum
 * @returns
 */
export function getChildSubaccountIds(address: string, parentSubaccountNum: number): string[] {
  return getChildSubaccountNums(parentSubaccountNum).map(
    (subaccountNumber: number): string => SubaccountTable.uuid(address, subaccountNumber),
  );
}

export function checkIfValidDydxAddress(address: string): boolean {
  const pattern: RegExp = /^dydx[0-9a-z]{39}$/;
  return pattern.test(address);
}

/**
 * Gets subaccount response objects given the subaccount, perpetual positions and perpetual markets
 * @param subaccount Subaccount to get response for, from the database
 * @param positions List of perpetual positions held by the subaccount, from the database
 * @param markets List of perpetual markets, from the database
 * @param assetPositions List of asset positions held by the subaccount, from the database
 * @param assets List of assets from the database
 * @param perpetualMarketsMap Mapping of perpetual markets to clob pairs, perpetual ids,
 *                            tickers from the database.
 * @param latestBlockHeight Latest block height from the database
 * @param latestFundingIndexMap Latest funding indices per perpetual from the database.
 * @param lastUpdatedFundingIndexMap Funding indices per perpetual for the last updated block of
 *                                   the subaccount.
 *
 * @returns Response object for the subaccount
 */
export function getSubaccountResponse(
  subaccount: SubaccountFromDatabase,
  perpetualPositions: PerpetualPositionFromDatabase[],
  assetPositions: AssetPositionFromDatabase[],
  assets: AssetFromDatabase[],
  markets: MarketFromDatabase[],
  perpetualMarketsMap: PerpetualMarketsMap,
  latestBlockHeight: string,
  latestFundingIndexMap: FundingIndexMap,
  lastUpdatedFundingIndexMap: FundingIndexMap,
): SubaccountResponseObject {
  const marketIdToMarket: MarketsMap = _.keyBy(
    markets,
    MarketColumns.id,
  );

  const unsettledFunding: Big = getTotalUnsettledFunding(
    perpetualPositions,
    latestFundingIndexMap,
    lastUpdatedFundingIndexMap,
  );

  const updatedPerpetualPositions:
  PerpetualPositionWithFunding[] = getPerpetualPositionsWithUpdatedFunding(
    initializePerpetualPositionsWithFunding(perpetualPositions),
    latestFundingIndexMap,
    lastUpdatedFundingIndexMap,
  );

  const filteredPerpetualPositions: PerpetualPositionWithFunding[
  ] = filterPositionsByLatestEventIdPerPerpetual(updatedPerpetualPositions);

  const perpetualPositionResponses:
  PerpetualPositionResponseObject[] = filteredPerpetualPositions.map(
    (perpetualPosition: PerpetualPositionWithFunding): PerpetualPositionResponseObject => {
      return perpetualPositionToResponseObject(
        perpetualPosition,
        perpetualMarketsMap,
        marketIdToMarket,
        subaccount.subaccountNumber,
      );
    },
  );

  const perpetualPositionsMap: PerpetualPositionsMap = _.keyBy(
    perpetualPositionResponses,
    'market',
  );

  const assetIdToAsset: AssetById = _.keyBy(
    assets,
    AssetColumns.id,
  );

  const sortedAssetPositions:
  AssetPositionFromDatabase[] = filterAssetPositions(assetPositions);

  const assetPositionResponses: AssetPositionResponseObject[] = sortedAssetPositions.map(
    (assetPosition: AssetPositionFromDatabase): AssetPositionResponseObject => {
      return assetPositionToResponseObject(
        assetPosition,
        assetIdToAsset,
        subaccount.subaccountNumber,
      );
    },
  );

  const assetPositionsMap: AssetPositionsMap = _.keyBy(
    assetPositionResponses,
    'symbol',
  );

  const {
    assetPositionsMap: adjustedAssetPositionsMap,
    adjustedUSDCAssetPositionSize,
  }: {
    assetPositionsMap: AssetPositionsMap,
    adjustedUSDCAssetPositionSize: string,
  } = adjustUSDCAssetPosition(assetPositionsMap, unsettledFunding);

  const {
    equity,
    freeCollateral,
  }: {
    equity: string,
    freeCollateral: string,
  } = calculateEquityAndFreeCollateral(
    filteredPerpetualPositions,
    perpetualMarketsMap,
    marketIdToMarket,
    adjustedUSDCAssetPositionSize,
  );

  return subaccountToResponseObject({
    subaccount,
    equity,
    freeCollateral,
    latestBlockHeight,
    openPerpetualPositions: perpetualPositionsMap,
    assetPositions: adjustedAssetPositionsMap,
  });
}

/* ------- PNL HELPERS ------- */

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
          equity: Big(aggregatedTick.equity).plus(pnlTick.equity).toString(),
          totalPnl: Big(aggregatedTick.totalPnl).plus(pnlTick.totalPnl).toString(),
          netTransfers: Big(aggregatedTick.netTransfers).plus(pnlTick.netTransfers).toString(),
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

/* ------- VAULT HELPERS ------- */

export async function getVaultMapping(): Promise<VaultMapping> {
  const vaults: VaultFromDatabase[] = await VaultTable.findAll(
    {},
    [],
    {},
  );
  const vaultMapping: VaultMapping = _.zipObject(
    vaults.map((vault: VaultFromDatabase): string => {
      return SubaccountTable.uuid(vault.address, 0);
    }),
    vaults,
  );
  const validVaultMapping: VaultMapping = {};
  for (const subaccountId of _.keys(vaultMapping)) {
    const perpetual: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher
      .getPerpetualMarketFromClobPairId(
        vaultMapping[subaccountId].clobPairId,
      );
    if (perpetual === undefined) {
      logger.warning({
        at: 'get-vault-mapping',
        message: `Vault clob pair id ${vaultMapping[subaccountId]} does not correspond to a ` +
          'perpetual market.',
        subaccountId,
      });
      continue;
    }
    validVaultMapping[subaccountId] = vaultMapping[subaccountId];
  }
  return validVaultMapping;
}

export function getVaultPnlStartDate(): DateTime {
  const startDate: DateTime = DateTime.fromISO(config.VAULT_PNL_START_DATE).toUTC();
  return startDate;
}
