/* ------- DATABASE MODEL TYPES ------- */
import Big from 'big.js';

import { CandleResolution } from './candle-types';
import { FillType, Liquidity } from './fill-types';
import {
  OrderSide, OrderStatus, OrderType, TimeInForce,
} from './order-types';
import { PerpetualMarketStatus, PerpetualMarketType } from './perpetual-market-types';
import { PerpetualPositionStatus } from './perpetual-position-types';
import { PositionSide } from './position-types';

type IsoString = string;

export interface IdBasedModelFromDatabase {
  id: string;
}

export interface SubaccountFromDatabase extends IdBasedModelFromDatabase {
  address: string,
  assetYieldIndex: string,
  subaccountNumber: number,
  updatedAt: IsoString,
  updatedAtHeight: string,
}

export interface PerpetualPositionFromDatabase extends IdBasedModelFromDatabase {
  id: string;
  subaccountId: string;
  perpetualId: string;
  side: PositionSide;
  status: PerpetualPositionStatus;
  size: string;  // The size of the position. Positive for long, negative for short.
  maxSize: string;
  entryPrice: string;
  exitPrice?: string;
  sumOpen: string;
  sumClose: string;
  createdAt: IsoString;
  closedAt?: IsoString;
  createdAtHeight: string;
  closedAtHeight?: string;
  openEventId: Buffer;
  closeEventId?: Buffer;
  lastEventId: Buffer;
  settledFunding: string;
  perpYieldIndex: string;
}

export interface OrderFromDatabase extends IdBasedModelFromDatabase {
  subaccountId: string;
  clientId: string;
  clobPairId: string;
  side: OrderSide;
  size: string;
  totalFilled: string;
  price: string;
  type: OrderType;
  status: OrderStatus;
  timeInForce: TimeInForce;
  reduceOnly: boolean;
  orderFlags: string;
  updatedAt: IsoString;
  updatedAtHeight: string;
  goodTilBlock?: string;
  goodTilBlockTime?: string;
  // createdAtHeight is optional because short term orders do not have a createdAtHeight.
  createdAtHeight?: string;
  clientMetadata: string;
  triggerPrice?: string;
}

export interface PerpetualMarketFromDatabase {
  id: string;
  clobPairId: string;
  ticker: string;
  marketId: number;
  status: PerpetualMarketStatus;
  priceChange24H: string;
  volume24H: string;
  trades24H: number;
  nextFundingRate: string;
  openInterest: string;
  quantumConversionExponent: number;
  atomicResolution: number;
  dangerIndexPpm: number;
  isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock: string;
  subticksPerTick: number;
  stepBaseQuantums: number;
  liquidityTierId: number;
  marketType: PerpetualMarketType;
  baseOpenInterest: string;
  perpYieldIndex: string;
}

export interface FillFromDatabase {
  id: string;
  subaccountId: string;
  side: OrderSide;
  liquidity: Liquidity;
  type: FillType;
  clobPairId: string;
  size: string;
  price: string;
  quoteAmount: string;
  eventId: Buffer;
  transactionHash: string;
  createdAt: IsoString;
  createdAtHeight: string;
  orderId?: string;
  clientMetadata?: string;
  fee: string;
}

export interface BlockFromDatabase {
  blockHeight: string;
  time: IsoString;
}

export interface TendermintEventFromDatabase {
  id: Buffer;
  blockHeight: string;
  transactionIndex: number;
  eventIndex: number;
}

export interface TransactionFromDatabase extends IdBasedModelFromDatabase {
  id: string,
  blockHeight: string,
  transactionIndex: number,
  transactionHash: string,
}

export interface AssetFromDatabase {
  id: string;
  symbol: string;
  atomicResolution: number;
  hasMarket: boolean;
  marketId?: number;
}

export interface AssetPositionFromDatabase {
  id: string;
  assetId: string;
  subaccountId: string;
  size: string;
  isLong: boolean;
}

export interface TransferFromDatabase extends IdBasedModelFromDatabase {
  senderSubaccountId?: string;
  recipientSubaccountId?: string;
  senderWalletAddress?: string;
  recipientWalletAddress?: string;
  assetId: string;
  size: string;
  eventId: Buffer;
  transactionHash: string;
  createdAt: IsoString;
  createdAtHeight: string;
}

export interface MarketFromDatabase {
  id: number;
  pair: string;
  exponent: number;
  minPriceChangePpm: number;
  spotPrice?: string;
  pnlPrice?: string;
}

export interface OraclePriceFromDatabase extends IdBasedModelFromDatabase {
  marketId: number;
  spotPrice: string;
  pnlPrice: string;
  effectiveAt: IsoString;
  effectiveAtHeight: string;
}

export interface LiquidityTiersFromDatabase {
  id: number;
  name: string;
  initialMarginPpm: string;
  maintenanceFractionPpm: string;
  openInterestLowerCap?: string;
  openInterestUpperCap?: string;
}

export interface CandleFromDatabase extends IdBasedModelFromDatabase {
  startedAt: IsoString;
  ticker: string;
  resolution: CandleResolution;
  low: string;
  high: string;
  open: string;
  close: string;
  baseTokenVolume: string;
  usdVolume: string;
  trades: number;
  startingOpenInterest: string;
}

export interface PnlTicksFromDatabase extends IdBasedModelFromDatabase {
  subaccountId: string;
  equity: string;
  totalPnl: string;
  netTransfers: string;
  createdAt: IsoString;
  blockHeight: string;
  blockTime: IsoString;
}

export interface FundingIndexUpdatesFromDatabase extends IdBasedModelFromDatabase {
  perpetualId: string;
  eventId: Buffer;
  rate: string;
  oraclePrice: string;
  fundingIndex: string;
  effectiveAt: string;
  effectiveAtHeight: string;
}

export interface ComplianceDataFromDatabase {
  address: string;
  chain?: string;
  blocked: boolean;
  riskScore?: string;
  updatedAt: string;
}

export interface YieldParamsFromDatabase extends IdBasedModelFromDatabase {
  sDAIPrice: string;
  assetYieldIndex: string;
  createdAt: IsoString;
  createdAtHeight: string;
}

export type SubaccountAssetNetTransferMap = { [subaccountId: string]:
{ [assetId: string]: string } };
export type SubaccountToPerpetualPositionsMap = { [subaccountId: string]:
{ [perpetualId: string]: PerpetualPositionFromDatabase } };
export type PerpetualPositionsMap = { [perpetualMarketId: string]: PerpetualPositionFromDatabase };
export type PerpetualMarketsMap = { [perpetualMarketId: string]: PerpetualMarketFromDatabase };
export type AssetsMap = { [assetId: string]: AssetFromDatabase };
export type LiquidityTiersMap = { [liquidityTierId: number]: LiquidityTiersFromDatabase };
export type SubaccountTDaiMap = { [subaccountId: string]: Big };
export type AssetPositionsMap = { [subaccountId: string]: AssetPositionFromDatabase[] };
export type MarketsMap = { [marketId: number]: MarketFromDatabase };
export type OraclePricesMap = { [marketId: number]: OraclePriceFromDatabase[] };
export type PriceMap = { [marketId: number]: { spotPrice: string; pnlPrice: string } };
export type FundingIndexMap = { [perpetualId: string]: Big };
export type CandlesResolutionMap = { [resolution: string]: CandleFromDatabase };
export type CandlesMap = { [ticker: string]: CandlesResolutionMap };
