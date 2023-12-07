import { logger } from '@dydxprotocol-indexer/base';
import {
  AssetPositionFromDatabase,
  BlockFromDatabase,
  FundingIndexMap,
  FundingIndexUpdatesTable,
  helpers,
  liquidityTierRefresher,
  LiquidityTiersFromDatabase,
  MarketFromDatabase,
  MarketsMap,
  PerpetualMarketFromDatabase,
  PerpetualMarketsMap,
  PerpetualMarketTable,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PositionSide,
  SubaccountFromDatabase,
  TendermintEventFromDatabase,
  TendermintEventTable,
  USDC_SYMBOL,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import express from 'express';
import _ from 'lodash';

import config from '../config';
import {
  AssetPositionResponseObject,
  AssetPositionsMap,
  MarketType,
  PerpetualPositionWithFunding,
  Risk,
} from '../types';
import { ONE, ZERO, ZERO_USDC_POSITION } from './constants';
import { NotFoundError } from './errors';

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
  return handleInternalServerError(
    at,
    message,
    error,
    req,
    res,
  );
}

function handleInternalServerError(
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
    params: JSON.stringify(req.params),
    query: JSON.stringify(req.query),
  });
  return createInternalServerErrorResponse(res);
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
    totalPositionRisk: Big
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
    adjustedInitialMarginFraction,
    adjustedMaintenanceMarginFraction,
  }: {
    adjustedInitialMarginFraction: Big,
    adjustedMaintenanceMarginFraction: Big,
  } = getAdjustedMarginFractions({
    liquidityTier,
    positionNotional,
  });
  return {
    signedNotional,
    individualRisk: {
      initial: positionNotional.times(adjustedInitialMarginFraction),
      maintenance: positionNotional.times(adjustedMaintenanceMarginFraction),
    },
  };
}

export function getAdjustedMarginFractions(
  {
    liquidityTier,
    positionNotional,
  }: {
    liquidityTier: LiquidityTiersFromDatabase,
    positionNotional: Big,
  },
): {
  adjustedInitialMarginFraction: Big,
  adjustedMaintenanceMarginFraction: Big,
} {
  const adjustedInitialMarginFraction: Big = getAdjustedMarginFraction({
    liquidityTier,
    positionNotional,
    initial: true,
  });
  const adjustedMaintenanceMarginFraction: Big = getAdjustedMarginFraction({
    liquidityTier,
    positionNotional,
    initial: false,
  });
  return {
    adjustedInitialMarginFraction,
    adjustedMaintenanceMarginFraction,
  };
}

/**
 * Get the adjusted margin fraction for a position in a given perpetual market.
 * Uses the `positionNotional`, `initialMarginFraction`,  and `basePositionNotional`
 * of the associated liquidity tier to calculate the adjusted initial margin fraction.
 *
 * @param liquidityTier The liquidity tier of the position.
 * @param positionNotional The notional value of the position.
 * @param initial Whether to compute the initial margin fraction or the maintenance margin fraction.
 *
 * @returns The adjusted margin fraction for the position in human-readable form.
 */
export function getAdjustedMarginFraction(
  {
    liquidityTier,
    positionNotional,
    initial,
  }: {
    liquidityTier: LiquidityTiersFromDatabase,
    positionNotional: Big,
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

  if (positionNotional.lte(liquidityTier.basePositionNotional)) {
    return Big(margin);
  }
  const adjustedImf: Big = Big(
    positionNotional.div(liquidityTier.basePositionNotional),
  ).sqrt().times(margin);
  if (adjustedImf.gte(ONE)) {
    return ONE;
  }
  return adjustedImf;
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
export async function filterPositionsByLatestEventIdPerPerpetual(
  positions: PerpetualPositionWithFunding[],
): Promise<PerpetualPositionWithFunding[]> {
  const events: TendermintEventFromDatabase[] = await TendermintEventTable.findAll(
    {
      id: positions.map((position: PerpetualPositionWithFunding) => position.lastEventId),
    },
    [],
  );
  const eventByIdHex: { [eventId: string]: TendermintEventFromDatabase } = _.keyBy(
    events,
    (event) => event.id.toString('hex'),
  );
  const sortedPositionsArray: PerpetualPositionWithFunding[] = positions.sort(
    (a: PerpetualPositionWithFunding, b: PerpetualPositionWithFunding): number => {
      const eventA: TendermintEventFromDatabase = eventByIdHex[a.lastEventId.toString('hex')];
      const eventB: TendermintEventFromDatabase = eventByIdHex[b.lastEventId.toString('hex')];
      return -1 * TendermintEventTable.compare(eventA, eventB);
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
  adjustedUSDCAssetPositionSize: string
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
      getUSDCAssetPosition(adjustedSize),
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

function getUSDCAssetPosition(signedSize: Big): AssetPositionResponseObject {
  const side: PositionSide = signedSize.gt(ZERO) ? PositionSide.LONG : PositionSide.SHORT;
  return {
    ...ZERO_USDC_POSITION,
    side,
    size: signedSize.abs().toFixed(),
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
