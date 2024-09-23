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
  PositionSide,
  SubaccountFromDatabase,
  SubaccountTable,
  TendermintEventFromDatabase,
  TendermintEventTable,
  TDAI_SYMBOL,
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
import { ZERO, ZERO_TDAI_POSITION } from './constants';
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
 * subaccount and the TDAI asset position of the subaccount.
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
  tdaiPositionSize: string,
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
  const equity: Big = signedPositionNotional.plus(tdaiPositionSize);
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
  const signedNotional: Big = size.times(market.pnlPrice!);
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
 * Gets and adjusts the TDAI asset position within a map of AssetPositions given the unsettled
 * funding
 * @param assetPositionsMap
 * @param unsettledFunding
 * @returns
 */
export function adjustTDAIAssetPosition(
  assetPositionsMap: AssetPositionsMap,
  unsettledFunding: Big,
): {
  assetPositionsMap: AssetPositionsMap,
  adjustedTDAIAssetPositionSize: string
} {
  let adjustedAssetPositionsMap: AssetPositionsMap = _.cloneDeep(assetPositionsMap);
  const tdaiPosition: AssetPositionResponseObject = _.get(assetPositionsMap, TDAI_SYMBOL);
  let signedTDaiPositionSize: Big;
  if (tdaiPosition?.size !== undefined) {
    signedTDaiPositionSize = Big(
      tdaiPosition.side === PositionSide.LONG
        ? tdaiPosition.size
        : -tdaiPosition.size,
    );
  } else {
    signedTDaiPositionSize = ZERO;
  }
  const adjustedSize: Big = signedTDaiPositionSize.plus(unsettledFunding);
  // Update the TDAI position in the map if the adjusted size is non-zero
  if (!adjustedSize.eq(ZERO)) {
    _.set(
      adjustedAssetPositionsMap,
      TDAI_SYMBOL,
      getTDAIAssetPosition(adjustedSize,
        adjustedAssetPositionsMap[TDAI_SYMBOL]?.subaccountNumber ?? 0),
    );
    // Remove the TDAI position in the map if the adjusted size is zero
  } else {
    adjustedAssetPositionsMap = _.omit(adjustedAssetPositionsMap, TDAI_SYMBOL);
  }

  return {
    assetPositionsMap: adjustedAssetPositionsMap,
    adjustedTDAIAssetPositionSize: adjustedSize.toFixed(),
  };
}

function getTDAIAssetPosition(signedSize: Big, subaccountNumber: number):
    AssetPositionResponseObject {
  const side: PositionSide = signedSize.gt(ZERO) ? PositionSide.LONG : PositionSide.SHORT;
  return {
    ...ZERO_TDAI_POSITION,
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
