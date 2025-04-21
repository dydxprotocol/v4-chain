import { PnlTicksCreateObject, IsoString } from '@dydxprotocol-indexer/postgres';
import { RedisOrder } from '@dydxprotocol-indexer/v4-protos';

// Type for the result of an order being placed
export interface PlaceOrderResult {
  // true if an order was placed
  placed: boolean,
  // true if an order was replaced
  replaced: boolean,
  // total filled of the old order in quantums, undefined if an order was not replaced
  oldTotalFilledQuantums?: number,
  // whether the old order was resting on the book, undefined if an order was not replaced
  restingOnBook?: boolean,
  // old order if an order was replaced, undefined if an order was not replaced
  oldOrder?: RedisOrder,
}

export interface RemoveOrderResult {
  // true if an order was removed
  removed: boolean,
  // total filled of the removed order in quantums, undefined if an order was not removed
  totalFilledQuantums?: number,
  // whether the removed order was resting on the book, undefined if an order was not removed
  restingOnBook?: boolean,
  // removed order if an order was removed, undefined if an order was not removed
  removedOrder?: RedisOrder,
}

export interface UpdateOrderResult {
  // true if an order was updated
  updated: boolean,
  // previous total filled of the order in quantums, undefined if an order was not updated
  oldTotalFilledQuantums?: number,
  // whether the updated order was resting on the book before the update, undefined if an order was
  // not updated
  oldRestingOnBook?: boolean,
  // order that was updated, undefined if an order was not updated
  order?: RedisOrder,
}

export interface OrderbookLevels {
  bids: PriceLevel[],
  asks: PriceLevel[],
}

export interface PriceLevel {
  // Total size of orders at a price in quantums
  quantums: string,
  // Human-readable price for the orderbook level
  humanPrice: string,
  // Timestamp of the most-recent edit to this value
  lastUpdated: string,
}

export interface OrderData {
  goodTilBlock: string,
  totalFilledQuantums: string,
  restingOnBook: boolean,
}

export type LuaScript = {
  // The name of the script
  readonly name: string,
  // The contents of the script
  readonly script: string,
  // The SHA1 hash of the contents of the script
  readonly hash: string,
};

export enum CanceledOrderStatus {
  CANCELED = 'CANCELED',
  BEST_EFFORT_CANCELED = 'BEST_EFFORT_CANCELED',
  NOT_CANCELED = 'NOT_CANCELED',
}

/* ------- PNL Creation TYPES ------- */
export type PnlTickForSubaccounts = {
  // Stores a PnlTicksCreateObject for the most recent pnl tick for each subaccount.
  // Opted for PnlTicksCreateObject instead ofPnlTicksFromDatabase as we don't need to store
  // the uuid.
  [subaccountId: string]: PnlTicksCreateObject,
};

/* -------- Stateful order update cache types -------- */
export interface StatefulOrderUpdateInfo {
  orderId: string,
  timestamp: number,
}

export interface CachedPnlTicks {
  equity: string,
  totalPnl: string,
  netTransfers: string,
  createdAt: string,
  blockHeight: string,
  blockTime: IsoString,
}

export interface CachedVaultHistoricalPnl {
  ticker: string,
  historicalPnl: CachedPnlTicks[],
}

/**
 * Redis internal, space-efficient representation of historical PNL data.
 * Format: [ticker, array of [equity, totalPnl, netTransfers, createdAtTimestamp,
 * blockHeight, blockTimeTimestamp]]
 */
export type RedisVaultsArray = [
  string,  // ticker
  [number, number, number, number, number, number][]  // [e, p, n, c, h, t][]
];

export interface CachedMegavaultPnl {
  pnlTicks: CachedPnlTicks[],
}

/**
 * Redis internal, space-efficient representation of megavault PNL data.
 * Format: array of [equity, totalPnl, netTransfers, createdAtTimestamp,
 * blockHeight, blockTimeTimestamp]
 */
export type RedisMegavaultPnl = [number, number, number, number, number, number][];
