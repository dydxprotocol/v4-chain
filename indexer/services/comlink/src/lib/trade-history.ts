import {
  FillFromDatabase,
  OrderSide,
  OrderType,
  PerpetualMarketType,
  PositionSide,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import _ from 'lodash';

import {
  MarketAndTypeByClobPairId,
  TradeHistoryResponseObject,
  TradeHistoryType,
} from '../types';

// ---------------------------------------------------------------------------
// Internal types
// ---------------------------------------------------------------------------

/** A group of fills belonging to the same order (or a single liquidation fill). */
interface FillGroup {
  orderId: string | null;
  fills: FillFromDatabase[];
  side: OrderSide;
  totalSize: Big;
  weightedPriceSum: Big; // sum(price * size) for weighted avg
  totalFee: Big;
  isLiquidation: boolean;
  clobPairId: string;
  lastCreatedAt: string;
  lastCreatedAtHeight: string;
}

/** Running state for a single market's position lifecycle. */
interface MarketState {
  positionSize: Big; // signed: positive = LONG, negative = SHORT
  entryPrice: Big;
  cumulativePnl: Big;
  cumulativeFee: Big;
}

// ---------------------------------------------------------------------------
// Exported functions
// ---------------------------------------------------------------------------

/**
 * Computes trade history rows from a chronologically-ordered list of fills.
 *
 * @param fills       All fills for the subaccount(s), ordered by createdAt ASC
 * @param orderTypeMap  orderId -> OrderType lookup
 * @param clobPairIdToMarket  clobPairId -> { market, marketType } lookup
 * @returns Trade history rows sorted by time DESC (most recent first)
 */
export function computeTradeHistory(
  fills: FillFromDatabase[],
  orderTypeMap: Record<string, OrderType>,
  clobPairIdToMarket: MarketAndTypeByClobPairId,
): TradeHistoryResponseObject[] {
  // Group fills by market
  const fillsByMarket: Record<string, FillFromDatabase[]> = _.groupBy(fills, 'clobPairId');

  const allRows: TradeHistoryResponseObject[] = [];

  for (const [clobPairId, marketFills] of Object.entries(fillsByMarket)) {
    const marketInfo = clobPairIdToMarket[clobPairId];
    if (!marketInfo?.market) continue;

    const rows = processMarketFills(
      marketFills, marketInfo.market, marketInfo.perpetualMarketType!, orderTypeMap,
    );
    allRows.push(...rows);
  }

  // Sort by time DESC (most recent first)
  allRows.sort((a, b) => (a.time > b.time ? -1 : a.time < b.time ? 1 : 0));

  return allRows;
}

/**
 * Applies in-memory pagination to sorted trade history rows.
 */
export function paginateTradeHistory(
  rows: TradeHistoryResponseObject[],
  limit: number,
  page?: number,
): {
  tradeHistory: TradeHistoryResponseObject[];
  pageSize: number;
  totalResults: number;
  offset: number;
} {
  const total = rows.length;
  if (page !== undefined) {
    const currentPage = Math.max(1, page);
    const offset = (currentPage - 1) * limit;
    return {
      tradeHistory: rows.slice(offset, offset + limit),
      pageSize: limit,
      totalResults: total,
      offset,
    };
  }
  return {
    tradeHistory: rows.slice(0, limit),
    pageSize: limit,
    totalResults: total,
    offset: 0,
  };
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

function processMarketFills(
  fills: FillFromDatabase[],
  marketId: string,
  marginMode: PerpetualMarketType,
  orderTypeMap: Record<string, OrderType>,
): TradeHistoryResponseObject[] {
  const fillGroups: FillGroup[] = groupFillsByOrder(fills);

  const state: MarketState = {
    positionSize: new Big(0),
    entryPrice: new Big(0),
    cumulativePnl: new Big(0),
    cumulativeFee: new Big(0),
  };

  const rows: TradeHistoryResponseObject[] = [];

  for (const group of fillGroups) {
    const newRows = processOrderGroup(group, state, marketId, marginMode, orderTypeMap);
    rows.push(...newRows);
  }

  return rows;
}

/**
 * Groups fills by orderId. Liquidation fills (orderId is null/undefined)
 * each become their own group. Regular fills with the same orderId are combined.
 * Order of groups follows the chronological order of their first fill.
 */
function groupFillsByOrder(fills: FillFromDatabase[]): FillGroup[] {
  const groups: FillGroup[] = [];
  const orderMap = new Map<string, FillGroup>();

  for (const fill of fills) {
    const isLiquidation = fill.orderId === null || fill.orderId === undefined;

    if (isLiquidation) {
      groups.push(createFillGroup(fill, true));
    } else {
      let group = orderMap.get(fill.orderId!);
      if (!group) {
        group = createFillGroup(fill, false);
        orderMap.set(fill.orderId!, group);
        groups.push(group);
      } else {
        addFillToGroup(group, fill);
      }
    }
  }

  return groups;
}

function createFillGroup(fill: FillFromDatabase, isLiquidation: boolean): FillGroup {
  return {
    orderId: fill.orderId ?? null,
    fills: [fill],
    side: fill.side,
    totalSize: new Big(fill.size),
    weightedPriceSum: new Big(fill.price).times(fill.size),
    totalFee: new Big(fill.fee),
    isLiquidation,
    clobPairId: fill.clobPairId,
    lastCreatedAt: fill.createdAt,
    lastCreatedAtHeight: fill.createdAtHeight,
  };
}

function addFillToGroup(group: FillGroup, fill: FillFromDatabase): void {
  group.fills.push(fill);
  group.totalSize = group.totalSize.plus(fill.size);
  group.weightedPriceSum = group.weightedPriceSum.plus(
    new Big(fill.price).times(fill.size),
  );
  group.totalFee = group.totalFee.plus(fill.fee);
  if (fill.createdAt > group.lastCreatedAt) {
    group.lastCreatedAt = fill.createdAt;
    group.lastCreatedAtHeight = fill.createdAtHeight;
  }
}

/**
 * Processes a single order group against the running market state.
 * Returns 1 row normally, or 2 rows if the order crosses zero.
 * Mutates `state` in place.
 */
function processOrderGroup(
  group: FillGroup,
  state: MarketState,
  marketId: string,
  marginMode: PerpetualMarketType,
  orderTypeMap: Record<string, OrderType>,
): TradeHistoryResponseObject[] {
  const avgPrice = group.weightedPriceSum.div(group.totalSize);

  // Signed delta: BUY is positive, SELL is negative
  const signedDelta = group.side === OrderSide.BUY
    ? group.totalSize
    : group.totalSize.times(-1);

  const positionBefore = state.positionSize;
  const positionAfter = positionBefore.plus(signedDelta);

  // Detect cross-zero: sign changed AND neither before nor after is zero
  const crossesZero = !positionBefore.eq(0)
    && !positionAfter.eq(0)
    && ((positionBefore.gt(0) && positionAfter.lt(0))
      || (positionBefore.lt(0) && positionAfter.gt(0)));

  if (crossesZero) {
    return handleCrossZero(group, state, marketId, marginMode, avgPrice, positionBefore,
      positionAfter, orderTypeMap);
  }

  return [computeSingleRow(group, state, marketId, marginMode, avgPrice, positionBefore,
    positionAfter, orderTypeMap)];
}

function handleCrossZero(
  group: FillGroup,
  state: MarketState,
  marketId: string,
  marginMode: PerpetualMarketType,
  avgPrice: Big,
  positionBefore: Big,
  positionAfter: Big,
  orderTypeMap: Record<string, OrderType>,
): TradeHistoryResponseObject[] {
  const closingSize = positionBefore.abs();
  const openingSize = positionAfter.abs();
  const orderType = group.orderId ? (orderTypeMap[group.orderId] ?? null) : null;

  // --- Row 1: CLOSE ---
  const closePnl = computeClosingPnl(positionBefore.gt(0), closingSize, avgPrice,
    state.entryPrice);
  const closeFee = group.totalFee.times(closingSize).div(group.totalSize);

  state.cumulativePnl = state.cumulativePnl.plus(closePnl);
  state.cumulativeFee = state.cumulativeFee.plus(closeFee);

  const closeType = group.isLiquidation
    ? TradeHistoryType.LIQUIDATION_CLOSE
    : TradeHistoryType.CLOSE;

  const closeRow: TradeHistoryResponseObject = {
    id: makeRowId(group, 'close'),
    action: closeType,
    executionPrice: avgPrice.toFixed(),
    side: group.side,
    positionSide: null, // position fully closed
    prevSize: positionBefore.abs().toFixed(),
    additionalSize: closingSize.times(-1).toFixed(), // negative (reducing)
    value: closingSize.times(avgPrice).toFixed(),
    orderType,
    netFee: state.cumulativeFee.toFixed(),
    netRealizedPnl: state.cumulativePnl.toFixed(),
    time: group.lastCreatedAt,
    orderId: group.orderId,
    marketId,
    marginMode,
  };

  // --- Reset for new lifecycle ---
  state.cumulativePnl = new Big(0);
  state.cumulativeFee = new Big(0);

  // --- Row 2: OPEN ---
  const openFee = group.totalFee.minus(closeFee);
  state.cumulativeFee = state.cumulativeFee.plus(openFee);
  state.positionSize = positionAfter;
  state.entryPrice = avgPrice;

  const openRow: TradeHistoryResponseObject = {
    id: makeRowId(group, 'open'),
    action: TradeHistoryType.OPEN,
    executionPrice: avgPrice.toFixed(),
    side: group.side,
    positionSide: positionAfter.gt(0) ? PositionSide.LONG : PositionSide.SHORT,
    prevSize: '0',
    additionalSize: openingSize.toFixed(), // positive (opening)
    value: openingSize.times(avgPrice).toFixed(),
    orderType,
    netFee: state.cumulativeFee.toFixed(),
    netRealizedPnl: state.cumulativePnl.toFixed(),
    time: group.lastCreatedAt,
    orderId: group.orderId,
    marketId,
    marginMode,
  };

  return [closeRow, openRow];
}

function computeSingleRow(
  group: FillGroup,
  state: MarketState,
  marketId: string,
  marginMode: PerpetualMarketType,
  avgPrice: Big,
  positionBefore: Big,
  positionAfter: Big,
  orderTypeMap: Record<string, OrderType>,
): TradeHistoryResponseObject {
  const isFlat = positionBefore.eq(0);
  const becomesFlat = positionAfter.eq(0);
  const isReducing = !isFlat && positionAfter.abs().lt(positionBefore.abs());

  // Determine action type
  let action: TradeHistoryType;
  if (isFlat) {
    action = TradeHistoryType.OPEN;
  } else if (becomesFlat) {
    action = group.isLiquidation
      ? TradeHistoryType.LIQUIDATION_CLOSE
      : TradeHistoryType.CLOSE;
  } else if (isReducing) {
    action = group.isLiquidation
      ? TradeHistoryType.LIQUIDATION_PARTIAL_CLOSE
      : TradeHistoryType.PARTIAL_CLOSE;
  } else {
    action = TradeHistoryType.EXTEND;
  }

  // Compute per-trade realized PnL for closing trades
  let tradePnl = new Big(0);
  if (isReducing || becomesFlat) {
    const closingAmount = becomesFlat
      ? positionBefore.abs()
      : positionBefore.abs().minus(positionAfter.abs());
    tradePnl = computeClosingPnl(positionBefore.gt(0), closingAmount, avgPrice,
      state.entryPrice);
  }

  // Update cumulative state
  state.cumulativePnl = state.cumulativePnl.plus(tradePnl);
  state.cumulativeFee = state.cumulativeFee.plus(group.totalFee);

  // Update entry price
  if (isFlat) {
    // Opening fresh position
    state.entryPrice = avgPrice;
  } else if (!isReducing && !becomesFlat) {
    // Extending: weighted average
    const existingValue = state.entryPrice.times(positionBefore.abs());
    const newValue = avgPrice.times(group.totalSize);
    state.entryPrice = existingValue.plus(newValue).div(positionAfter.abs());
  }
  // Partial close / full close: entry price stays the same

  // Compute additionalSize (signed: positive = increasing, negative = reducing)
  const signedDelta = group.side === OrderSide.BUY
    ? group.totalSize
    : group.totalSize.times(-1);

  // Determine position side after this trade
  let positionSide: PositionSide | null = null;
  if (positionAfter.gt(0)) {
    positionSide = PositionSide.LONG;
  } else if (positionAfter.lt(0)) {
    positionSide = PositionSide.SHORT;
  }

  // Build the row
  const row: TradeHistoryResponseObject = {
    id: makeRowId(group),
    action,
    executionPrice: avgPrice.toFixed(),
    side: group.side,
    positionSide,
    prevSize: positionBefore.abs().toFixed(),
    additionalSize: signedDelta.toFixed(),
    value: group.totalSize.times(avgPrice).toFixed(),
    orderType: group.orderId ? (orderTypeMap[group.orderId] ?? null) : null,
    netFee: state.cumulativeFee.toFixed(),
    netRealizedPnl: state.cumulativePnl.toFixed(),
    time: group.lastCreatedAt,
    orderId: group.orderId,
    marketId,
    marginMode,
  };

  // Update position
  state.positionSize = positionAfter;

  // If fully closed, reset for next lifecycle
  if (becomesFlat) {
    state.cumulativePnl = new Big(0);
    state.cumulativeFee = new Big(0);
    state.entryPrice = new Big(0);
  }

  return row;
}

/**
 * Generates a unique ID for a trade history row.
 *   Normal row:    orderId (or fillId for liquidations)
 *   Cross-zero:    orderId:close / orderId:open
 */
function makeRowId(group: FillGroup, suffix?: 'close' | 'open'): string {
  const base = group.orderId ?? group.fills[0].id;
  return suffix ? `${base}:${suffix}` : base;
}

/**
 * PnL formula (mirrors dydx_apply_fill_realized_effects.sql):
 *   LONG closing:  (fillPrice - entryPrice) * closingAmount
 *   SHORT closing: (entryPrice - fillPrice) * closingAmount
 */
function computeClosingPnl(
  isLong: boolean,
  closingAmount: Big,
  fillPrice: Big,
  entryPrice: Big,
): Big {
  if (isLong) {
    return fillPrice.minus(entryPrice).times(closingAmount);
  }
  return entryPrice.minus(fillPrice).times(closingAmount);
}
