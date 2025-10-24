/* ------- DATABASE MODEL TYPES ------- */
import Big from 'big.js';

import { CandleResolution } from './candle-types';
import { ComplianceReason, ComplianceStatus } from './compliance-status-types';
import { FillType, Liquidity } from './fill-types';
import {
  OrderSide, OrderStatus, OrderType, TimeInForce,
} from './order-types';
import { PerpetualMarketStatus, PerpetualMarketType } from './perpetual-market-types';
import { PerpetualPositionStatus } from './perpetual-position-types';
import { PositionSide } from './position-types';
import { TradingRewardAggregationPeriod } from './trading-reward-aggregation-types';
import { VaultStatus } from './vault-types';

type IsoString = string;

export interface IdBasedModelFromDatabase {
  id: string,
}

export interface SubaccountFromDatabase extends IdBasedModelFromDatabase {
  address: string,
  subaccountNumber: number,
  updatedAt: IsoString,
  updatedAtHeight: string,
}

export interface WalletFromDatabase {
  address: string,
  totalTradingRewards: string,
  totalVolume: string,
}

export interface PerpetualPositionFromDatabase extends IdBasedModelFromDatabase {
  id: string,
  subaccountId: string,
  perpetualId: string,
  side: PositionSide,
  status: PerpetualPositionStatus,
  size: string,  // The size of the position. Positive for long, negative for short.
  maxSize: string,
  entryPrice: string,
  exitPrice?: string,
  sumOpen: string,
  sumClose: string,
  createdAt: IsoString,
  closedAt?: IsoString,
  createdAtHeight: string,
  closedAtHeight?: string,
  openEventId: Buffer,
  closeEventId?: Buffer,
  lastEventId: Buffer,
  settledFunding: string,
  totalRealizedPnl?: string,
}

export interface OrderFromDatabase extends IdBasedModelFromDatabase {
  subaccountId: string,
  clientId: string,
  clobPairId: string,
  side: OrderSide,
  size: string,
  totalFilled: string,
  price: string,
  type: OrderType,
  status: OrderStatus,
  timeInForce: TimeInForce,
  reduceOnly: boolean,
  orderFlags: string,
  updatedAt: IsoString,
  updatedAtHeight: string,
  goodTilBlock?: string,
  goodTilBlockTime?: string,
  // createdAtHeight is optional because short term orders do not have a createdAtHeight.
  createdAtHeight?: string,
  clientMetadata: string,
  triggerPrice?: string,
  builderAddress?: string,
  feePpm?: string,
  orderRouterAddress?: string,
  // these fields only exist for twap orders
  duration?: string,
  interval?: string,
  priceTolerance?: string,
}

export interface PerpetualMarketFromDatabase {
  id: string,
  clobPairId: string,
  ticker: string,
  marketId: number,
  status: PerpetualMarketStatus,
  priceChange24H: string,
  volume24H: string,
  trades24H: number,
  nextFundingRate: string,
  openInterest: string,
  quantumConversionExponent: number,
  atomicResolution: number,
  subticksPerTick: number,
  stepBaseQuantums: number,
  liquidityTierId: number,
  marketType: PerpetualMarketType,
  baseOpenInterest: string,
  defaultFundingRate1H?: string,
}

export interface FillFromDatabase {
  id: string,
  subaccountId: string,
  side: OrderSide,
  liquidity: Liquidity,
  type: FillType,
  clobPairId: string,
  size: string,
  price: string,
  quoteAmount: string,
  eventId: Buffer,
  transactionHash: string,
  createdAt: IsoString,
  createdAtHeight: string,
  orderId?: string,
  clientMetadata?: string,
  fee: string,
  affiliateRevShare: string,
  builderAddress?: string,
  builderFee?: string,
  orderRouterAddress?: string,
  orderRouterFee?: string,
  positionSizeBefore?: string,
  entryPriceBefore?: string,
  positionSideBefore?: PositionSide,
}

export interface BlockFromDatabase {
  blockHeight: string,
  time: IsoString,
}

export interface TendermintEventFromDatabase {
  id: Buffer,
  blockHeight: string,
  transactionIndex: number,
  eventIndex: number,
}

export interface TransactionFromDatabase extends IdBasedModelFromDatabase {
  id: string,
  blockHeight: string,
  transactionIndex: number,
  transactionHash: string,
}

export interface AssetFromDatabase {
  id: string,
  symbol: string,
  atomicResolution: number,
  hasMarket: boolean,
  marketId?: number,
}

export interface AssetPositionFromDatabase {
  id: string,
  assetId: string,
  subaccountId: string,
  size: string,
  isLong: boolean,
}

export interface TransferFromDatabase extends IdBasedModelFromDatabase {
  senderSubaccountId?: string,
  recipientSubaccountId?: string,
  senderWalletAddress?: string,
  recipientWalletAddress?: string,
  assetId: string,
  size: string,
  eventId: Buffer,
  transactionHash: string,
  createdAt: IsoString,
  createdAtHeight: string,
}

export interface MarketFromDatabase {
  id: number,
  pair: string,
  exponent: number,
  minPriceChangePpm: number,
  oraclePrice?: string,
}

export interface OraclePriceFromDatabase extends IdBasedModelFromDatabase {
  marketId: number,
  price: string,
  effectiveAt: IsoString,
  effectiveAtHeight: string,
}

export interface LiquidityTiersFromDatabase {
  id: number,
  name: string,
  initialMarginPpm: string,
  maintenanceFractionPpm: string,
  openInterestLowerCap?: string,
  openInterestUpperCap?: string,
}

export interface CandleFromDatabase extends IdBasedModelFromDatabase {
  startedAt: IsoString,
  ticker: string,
  resolution: CandleResolution,
  low: string,
  high: string,
  open: string,
  close: string,
  baseTokenVolume: string,
  usdVolume: string,
  trades: number,
  startingOpenInterest: string,
  orderbookMidPriceOpen?: string | null,
  orderbookMidPriceClose?: string | null,
}

export interface PnlTicksFromDatabase extends IdBasedModelFromDatabase {
  subaccountId: string,
  equity: string,
  totalPnl: string,
  netTransfers: string,
  createdAt: IsoString,
  blockHeight: string,
  blockTime: IsoString,
}

export interface FundingIndexUpdatesFromDatabase extends IdBasedModelFromDatabase {
  perpetualId: string,
  eventId: Buffer,
  rate: string,
  oraclePrice: string,
  fundingIndex: string,
  effectiveAt: string,
  effectiveAtHeight: string,
}

export interface ComplianceDataFromDatabase {
  address: string,
  chain?: string,
  blocked: boolean,
  riskScore?: string,
  updatedAt: string,
}

export interface ComplianceStatusFromDatabase {
  address: string,
  status: ComplianceStatus,
  reason?: ComplianceReason,
  createdAt: IsoString,
  updatedAt: IsoString,
}

export interface TradingRewardFromDatabase {
  id: string,
  address: string,
  blockTime: IsoString,
  blockHeight: string,
  amount: string,
}

export interface TradingRewardAggregationFromDatabase {
  id: string,
  address: string,
  startedAt: IsoString,
  startedAtHeight: string,
  endedAt?: IsoString,
  endedAtHeight?: string,
  period: TradingRewardAggregationPeriod,
  amount: string,
}

export interface SubaccountUsernamesFromDatabase {
  username: string,
  subaccountId: string,
}

export interface AddressUsername {
  address: string,
  username: string,
}

export interface LeaderboardPnlFromDatabase {
  address: string,
  timeSpan: string,
  pnl: string,
  currentEquity: string,
  rank: number,
}

export interface PersistentCacheFromDatabase {
  key: string,
  value: string,
}

export interface AffiliateInfoFromDatabase {
  address: string,
  affiliateEarnings: string,
  referredMakerTrades: number,
  referredTakerTrades: number,
  totalReferredMakerFees: string,
  totalReferredTakerFees: string,
  totalReferredMakerRebates: string,
  totalReferredUsers: number,
  firstReferralBlockHeight: string,
  referredTotalVolume: string,
}

export interface AffiliateReferredUserFromDatabase {
  affiliateAddress: string,
  refereeAddress: string,
  referredAtBlock: string,
}

export interface FirebaseNotificationTokenFromDatabase {
  address: string,
  token: string,
  updatedAt: IsoString,
  language: string,
}

export interface VaultFromDatabase {
  address: string,
  clobPairId: string,
  status: VaultStatus,
  createdAt: IsoString,
  updatedAt: IsoString,
}

export interface FundingPaymentsFromDatabase {
  subaccountId: string,
  perpetualId: string,
  ticker: string,
  createdAt: string,
  createdAtHeight: string,
  oraclePrice: string,
  size: string,
  side: PositionSide,
  rate: string,
  payment: string,
  fundingIndex: string,
}

export interface TurnkeyUserFromDatabase {
  suborg_id: string,
  username?: string,
  email?: string,
  svm_address: string,
  evm_address: string,
  smart_account_address?: string,
  salt: string,
  dydx_address?: string,
  created_at: string,
}

export interface PnlFromDatabase {
  subaccountId: string,
  createdAt: IsoString,
  createdAtHeight: string,
  equity: string,
  netTransfers: string,
  totalPnl: string,
}

export interface BridgeInformationFromDatabase {
  id: string,
  from_address: string,
  chain_id: string,
  amount: string,
  transaction_hash?: string,
  created_at: IsoString,
}

export type SubaccountAssetNetTransferMap = { [subaccountId: string]:
{ [assetId: string]: string }, };
export type SubaccountToPerpetualPositionsMap = { [subaccountId: string]:
{ [perpetualId: string]: PerpetualPositionFromDatabase }, };
export type PerpetualPositionsMap = { [perpetualMarketId: string]: PerpetualPositionFromDatabase };
export type PerpetualMarketsMap = { [perpetualMarketId: string]: PerpetualMarketFromDatabase };
export type AssetsMap = { [assetId: string]: AssetFromDatabase };
export type LiquidityTiersMap = { [liquidityTierId: number]: LiquidityTiersFromDatabase };
export type SubaccountUsdcMap = { [subaccountId: string]: Big };
export type AssetPositionsMap = { [subaccountId: string]: AssetPositionFromDatabase[] };
export type MarketsMap = { [marketId: number]: MarketFromDatabase };
export type OraclePricesMap = { [marketId: number]: OraclePriceFromDatabase[] };
export type PriceMap = { [marketId: number]: string };
export type FundingIndexMap = { [perpetualId: string]: Big };
export type CandlesResolutionMap = { [resolution: string]: CandleFromDatabase };
export type CandlesMap = { [ticker: string]: CandlesResolutionMap };
