import { CandleResolution } from './candle-types';
import { FillType, Liquidity } from './fill-types';
import {
  OrderSide,
  OrderStatus,
  OrderType,
} from './order-types';
import { PerpetualMarketStatus, PerpetualMarketType } from './perpetual-market-types';
import { PerpetualPositionStatus } from './perpetual-position-types';
import { PositionSide } from './position-types';
import { TradeType } from './trade-types';
import { TransferType } from './transfer-types';
import { IsoString } from './utility-types';

/**
 * All types that will be stringified into the contents field of a websocket message.
 * Any changes to this file should also update the corresponding websocket message version in
 * packages/kafka/src/constants.ts.
 */

/* ------- OrderbookMessageContents ------- */

export interface OrderbookMessageContents {
  bids?: PriceLevel[],
  asks?: PriceLevel[],
}

// The first string indicates the price, the second string indicates the size
export type PriceLevel = [string, string];

/* ------- SubaccountMessageContents ------- */

export interface SubaccountMessageContents {
  perpetualPositions?: PerpetualPositionSubaccountMessageContents[],
  assetPositions?: AssetPositionSubaccountMessageContents[],
  orders?: OrderSubaccountMessageContents[],
  fills?: FillSubaccountMessageContents[],
  transfers?: TransferSubaccountMessageContents,
  tradingReward?: TradingRewardSubaccountMessageContents,
  blockHeight?: string,
}

export interface PerpetualPositionSubaccountMessageContents {
  address: string,
  subaccountNumber: number,
  positionId: string,
  market: string,
  side: PositionSide,
  status: PerpetualPositionStatus,
  size: string,
  maxSize: string,
  netFunding: string,
  entryPrice: string,
  exitPrice?: string,
  sumOpen: string,
  sumClose: string,
  realizedPnl?: string,
  unrealizedPnl?: string,
}

export interface AssetPositionSubaccountMessageContents {
  address: string,
  subaccountNumber: number,
  positionId: string,
  assetId: string,
  symbol: string,
  side: PositionSide,
  size: string,
}

// TODO(DEC-1659): Change this to match API specification devised by Product team.
// API does not have POST_ONLY as a valid TimeInForce value
// Chose to implement manually instead of with `Exclude<TimeInForce, TimeInForce.POST_ONLY>`
// because Exclude breaks conversion from swagger to markdown
export enum APITimeInForce {
  // GTT represents Good-Til-Time, where an order will first match with existing orders on the book
  // and any remaining size will be added to the book as a maker order, which will expire at a
  // given expiry time.
  GTT = 'GTT',
  // FOK represents Fill-Or-KILl where it's enforced that an order will either be filled
  // completely and immediately by maker orders on the book or canceled if the entire amount can't
  // be filled.
  FOK = 'FOK',
  // IOC represents Immediate-Or-Cancel, where it's enforced that an order only be matched with
  // maker orders on the book. If the order has remaining size after matching with existing orders
  // on the book, the remaining size is not placed on the book.
  IOC = 'IOC',
}

// Create a superset of the `OrderStatus` type including `BEST_EFFORT_OPENED` which should only
// exist in API responses.
export enum BestEffortOpenedStatus {
  BEST_EFFORT_OPENED = 'BEST_EFFORT_OPENED',
}
export type APIOrderStatus = OrderStatus | BestEffortOpenedStatus;
export const APIOrderStatusEnum = {
  ...OrderStatus,
  ...BestEffortOpenedStatus,
};

export interface OrderSubaccountMessageContents {
  id: string,
  subaccountId: string,
  clientId: string,
  clobPairId?: string,
  side?: OrderSide,
  size?: string,
  ticker?: string,
  price?: string,
  type?: OrderType,
  timeInForce?: APITimeInForce,
  postOnly?: boolean,
  reduceOnly?: boolean,
  status: APIOrderStatus,
  orderFlags: string,

  totalFilled?: string,
  totalOptimisticFilled?: string,
  goodTilBlock?: string,
  goodTilBlockTime?: string,
  triggerPrice?: string,
  updatedAt?: IsoString,
  updatedAtHeight?: string,

  // This will only be filled if the order was removed
  removalReason?: string,
  // This will only be set for stateful orders
  createdAtHeight?: string,
  clientMetadata?: string,
  // This will only be set for twap orders
  duration?: string,
  interval?: string,
  priceTolerance?: string,
}

export interface FillSubaccountMessageContents {
  id: string,
  subaccountId: string,
  side: OrderSide,
  liquidity: Liquidity,
  type: FillType,
  clobPairId: string,
  size: string,
  price: string,
  quoteAmount: string,
  eventId: string,
  transactionHash: string,
  createdAt: IsoString,
  createdAtHeight: string,
  orderId?: string,
  ticker: string,
  clientMetadata?: string,
}

export interface TransferSubaccountMessageContents {
  sender: {
    address: string,
    subaccountNumber?: number,
  },
  recipient: {
    address: string,
    subaccountNumber?: number,
  },
  symbol: string,
  size: string,
  type: TransferType,
  transactionHash: string,
  createdAt: IsoString,
  createdAtHeight: string,
}

export interface TradingRewardSubaccountMessageContents {
  tradingReward: string,
  createdAtHeight: string,
  createdAt: string,
}

/* ------- TradeMessageContents ------- */

export interface TradeMessageContents {
  trades: TradeContent[],
}

export interface TradeContent {
// Unique id of the trade, which is the taker fill id.
  id: string,
  size: string,
  price: string,
  side: string,
  createdAt: IsoString,
  type: TradeType,
}

/* ------- MarketMessageContents ------- */

export interface MarketMessageContents {
  trading?: TradingMarketMessageContents,
  oraclePrices?: OraclePriceMarketMessageContentsMapping,
}

export type TradingMarketMessageContents = {
  [ticker: string]: TradingPerpetualMarketMessage,
};

// All the fields in PerpetualMarketFromDatabase, but optional
export interface TradingPerpetualMarketMessage {
  // These fields are very unlikely to change
  id?: string,
  clobPairId?: string,
  ticker?: string,
  marketId?: number,
  status?: PerpetualMarketStatus,
  initialMarginFraction?: string,
  maintenanceMarginFraction?: string,
  openInterest?: string,
  quantumConversionExponent?: number,
  atomicResolution?: number,
  subticksPerTick?: number,
  stepBaseQuantums?: number,
  marketType?: PerpetualMarketType,
  openInterestLowerCap?: string,
  openInterestUpperCap?: string,
  baseOpenInterest?: string,
  defaultFundingRate1H?: string,

  // Fields that are likely to change
  priceChange24H?: string,
  volume24H?: string,
  trades24H?: number,
  nextFundingRate?: string,

  // Derived fields
  tickSize?: string,
  stepSize?: string,
}

export type OraclePriceMarketMessageContentsMapping = {
  [ticker: string]: OraclePriceMarket,
};

export interface OraclePriceMarket {
  oraclePrice: string,
  effectiveAt: IsoString,
  effectiveAtHeight: string,
  marketId: number,
}

/* ------- CandleMessageContents ------- */

export interface CandleMessageContents {
  resolution: CandleResolution,
  startedAt: IsoString,
  ticker: string,
  low: string,
  high: string,
  open: string,
  close: string,
  baseTokenVolume: string,
  trades: number,
  usdVolume: string,
  startingOpenInterest: string,
}
