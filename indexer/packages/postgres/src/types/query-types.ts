/* ------- QUERY TYPES ------- */

import { CandleResolution } from './candle-types';
import { Liquidity } from './fill-types';
import { OrderSide, OrderStatus, OrderType } from './order-types';
import { PerpetualPositionStatus } from './perpetual-position-types';
import { PositionSide } from './position-types';
import { IsoString } from './utility-types';

export enum QueryableField {
  LIMIT = 'limit',
  ID = 'id',
  ADDRESS = 'address',
  ASSET_YIELD_INDEX = 'assetYieldIndex',
  PERP_YIELD_INDEX = 'perpYieldIndex',
  SUBACCOUNT_NUMBER = 'subaccountNumber',
  SUBACCOUNT_ID = 'subaccountId',
  SENDER_SUBACCOUNT_ID = 'senderSubaccountId',
  RECIPIENT_SUBACCOUNT_ID = 'recipientSubaccountId',
  SENDER_WALLET_ADDRESS = 'senderWalletAddress',
  RECIPIENT_WALLET_ADDRESS = 'recipientWalletAddress',
  CLIENT_ID = 'clientId',
  CLOB_PAIR_ID = 'clobPairId',
  SIDE = 'side',
  SIZE = 'size',
  TOTAL_FILLED = 'totalFilled',
  PRICE = 'price',
  SPOT_PRICE = 'spotPrice',
  PNL_PRICE = 'pnlPrice',
  TYPE = 'type',
  STATUS = 'status',
  STATUSES = 'statuses',
  POST_ONLY = 'postOnly',
  REDUCE_ONLY = 'reduceOnly',
  PERPETUAL_ID = 'perpetualId',
  LIQUIDITY = 'liquidity',
  CREATED_BEFORE_OR_AT = 'createdBeforeOrAt',
  CREATED_BEFORE_OR_AT_HEIGHT = 'createdBeforeOrAtHeight',
  CREATED_ON_OR_AFTER = 'createdOnOrAfter',
  CREATED_ON_OR_AFTER_HEIGHT = 'createdOnOrAfterHeight',
  CREATED_AFTER = 'createdAfter',
  CREATED_AFTER_HEIGHT = 'createdAfterHeight',
  CREATED_AT = 'createdAt',
  CREATED_AT_HEIGHT = 'createdAtHeight',
  EVENT_ID = 'eventId',
  TRANSACTION_HASH = 'transactionHash',
  BLOCK_HEIGHT = 'blockHeight',
  BLOCK_TIME = 'blockTime',
  TRANSACTION_INDEX = 'transactionIndex',
  EVENT_INDEX = 'eventIndex',
  SYMBOL = 'symbol',
  ATOMIC_RESOLUTION = 'atomicResolution',
  HAS_MARKET = 'hasMarket',
  MARKET_ID = 'marketId',
  ASSET_ID = 'assetId',
  IS_LONG = 'isLong',
  PAIR = 'pair',
  EFFECTIVE_AT = 'effectiveAt',
  EFFECTIVE_AT_HEIGHT = 'effectiveAtHeight',
  EFFECTIVE_BEFORE_OR_AT = 'effectiveBeforeOrAt',
  EFFECTIVE_BEFORE_OR_AT_HEIGHT = 'effectiveBeforeOrAtHeight',
  GOOD_TIL_BLOCK_BEFORE_OR_AT = 'goodTilBlockBeforeOrAt',
  GOOD_TIL_BLOCK_TIME_BEFORE_OR_AT = 'goodTilBlockTimeBeforeOrAt',
  TICKER = 'ticker',
  RESOLUTION = 'resolution',
  FROM_ISO = 'fromISO',
  TO_ISO = 'toISO',
  CREATED_BEFORE_OR_AT_BLOCK_HEIGHT = 'createdBeforeOrAtBlockHeight',
  CREATED_ON_OR_AFTER_BLOCK_HEIGHT = 'createdOnOrAfterBlockHeight',
  ORDER_FLAGS = 'orderFlags',
  CLIENT_METADATA = 'clientMetadata',
  LIQUIDITY_TIER_ID = 'liquidityTierId',
  FEE = 'fee',
  TRIGGER_PRICE = 'triggerPrice',
  UPDATED_BEFORE_OR_AT = 'updatedBeforeOrAt',
  UPDATED_ON_OR_AFTER = 'updatedOnOrAfter',
  PROVIDER = 'provider',
  BLOCKED = 'blocked',
  BLOCK_TIME_BEFORE_OR_AT = 'blockTimeBeforeOrAt',
  STARTED_AT_HEIGHT = 'startedAtHeight',
  PERIOD = 'period',
  STARTED_AT_HEIGHT_OR_AFTER = 'startedAtHeightOrAfter',
  BLOCK_TIME_AFTER_OR_AT = 'blockTimeAfterOrAt',
  BLOCK_TIME_BEFORE = 'blockTimeBefore',
  ADDRESSES = 'addresses',
  BLOCK_HEIGHT_BEFORE_OR_AT = 'blockHeightBeforeOrAt',
  STARTED_AT_BEFORE_OR_AT = 'startedAtBeforeOrAt',
  STARTED_AT_HEIGHT_BEFORE_OR_AT = 'startedAtHeightBeforeOrAt',
  S_DAI_PRICE = 'sDAIPrice',
}

export interface QueryConfig {
  [QueryableField.LIMIT]?: number;
}

export interface SubaccountQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.ADDRESS]?: string;
  [QueryableField.SUBACCOUNT_NUMBER]?: number;
  [QueryableField.UPDATED_BEFORE_OR_AT]?: string;
  [QueryableField.UPDATED_ON_OR_AFTER]?: string;
  [QueryableField.ASSET_YIELD_INDEX]?: string;
}

export interface YieldParamsQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.CREATED_AT_HEIGHT]?: string[];
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string;
  [QueryableField.CREATED_AFTER_HEIGHT]?: string;
  [QueryableField.CREATED_AT]?: string;
  [QueryableField.CREATED_BEFORE_OR_AT]?: string;
  [QueryableField.CREATED_AFTER]?: string;
  [QueryableField.ASSET_YIELD_INDEX]?: string;
  [QueryableField.S_DAI_PRICE]?: string;
}

export interface WalletQueryConfig extends QueryConfig {
  [QueryableField.ADDRESS]?: string;
}

export interface PerpetualPositionQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.SUBACCOUNT_ID]?: string[];
  [QueryableField.PERPETUAL_ID]?: string[];
  [QueryableField.SIDE]?: PositionSide;
  [QueryableField.STATUS]?: PerpetualPositionStatus[];
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string;
  [QueryableField.CREATED_BEFORE_OR_AT]?: string;
  [QueryableField.PERP_YIELD_INDEX]?: string;
}

export interface OrderQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.SUBACCOUNT_ID]?: string[];
  [QueryableField.CLIENT_ID]?: string;
  [QueryableField.CLOB_PAIR_ID]?: string;
  [QueryableField.SIDE]?: OrderSide;
  [QueryableField.SIZE]?: string;
  [QueryableField.TOTAL_FILLED]?: string;
  [QueryableField.PRICE]?: string;
  [QueryableField.TYPE]?: OrderType;
  [QueryableField.STATUSES]?: OrderStatus[];
  [QueryableField.POST_ONLY]?: boolean;
  [QueryableField.REDUCE_ONLY]?: boolean;
  [QueryableField.GOOD_TIL_BLOCK_BEFORE_OR_AT]?: string;
  [QueryableField.GOOD_TIL_BLOCK_TIME_BEFORE_OR_AT]?: string;
  [QueryableField.ORDER_FLAGS]?: string;
  [QueryableField.CLIENT_METADATA]?: string;
  [QueryableField.TRIGGER_PRICE]?: string;
}

export interface PerpetualMarketQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.MARKET_ID]?: number[];
  [QueryableField.LIQUIDITY_TIER_ID]?: number[];
}

export interface FillQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.SUBACCOUNT_ID]?: string[];
  [QueryableField.SIDE]?: OrderSide;
  [QueryableField.LIQUIDITY]?: Liquidity;
  [QueryableField.TYPE]?: OrderType;
  [QueryableField.CLOB_PAIR_ID]?: string;
  [QueryableField.EVENT_ID]?: Buffer;
  [QueryableField.TRANSACTION_HASH]?: string;
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string;
  [QueryableField.CREATED_BEFORE_OR_AT]?: string;
  [QueryableField.CREATED_ON_OR_AFTER_HEIGHT]?: string;
  [QueryableField.CREATED_ON_OR_AFTER]?: string;
  [QueryableField.CLIENT_METADATA]?: string;
  [QueryableField.FEE]?: string;
}

export interface BlockQueryConfig extends QueryConfig {
  [QueryableField.BLOCK_HEIGHT]?: string[];
  [QueryableField.CREATED_ON_OR_AFTER]?: string;
}

export interface TendermintEventQueryConfig extends QueryConfig {
  [QueryableField.ID]?: Buffer[];
  [QueryableField.BLOCK_HEIGHT]?: string[];
  [QueryableField.TRANSACTION_INDEX]?: number[];
  [QueryableField.EVENT_INDEX]?: number[];
}

export interface TransactionQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.BLOCK_HEIGHT]?: string[];
  [QueryableField.TRANSACTION_INDEX]?: number[];
  [QueryableField.TRANSACTION_HASH]?: string[];
}

export interface AssetQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.SYMBOL]?: string;
  [QueryableField.ATOMIC_RESOLUTION]?: number;
  [QueryableField.HAS_MARKET]?: boolean;
  [QueryableField.MARKET_ID]?: number;
}

export interface AssetPositionQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.ASSET_ID]?: string[];
  [QueryableField.SUBACCOUNT_ID]?: string[];
  [QueryableField.SIZE]?: string;
  [QueryableField.IS_LONG]?: boolean;
}

export interface TransferQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.SENDER_SUBACCOUNT_ID]?: string[];
  [QueryableField.RECIPIENT_SUBACCOUNT_ID]?: string[];
  [QueryableField.SENDER_WALLET_ADDRESS]?: string[];
  [QueryableField.RECIPIENT_WALLET_ADDRESS]?: string[];
  [QueryableField.ASSET_ID]?: string[];
  [QueryableField.SIZE]?: string;
  [QueryableField.EVENT_ID]?: Buffer[];
  [QueryableField.TRANSACTION_HASH]?: string[];
  [QueryableField.CREATED_AT]?: string;
  [QueryableField.CREATED_AT_HEIGHT]?: string[];
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string;
  [QueryableField.CREATED_BEFORE_OR_AT]?: string;
  [QueryableField.CREATED_AFTER]?: string;
  [QueryableField.CREATED_AFTER_HEIGHT]?: string;
}

export interface ToAndFromSubaccountTransferQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.SUBACCOUNT_ID]?: string[];
  [QueryableField.ASSET_ID]?: string[];
  [QueryableField.SIZE]?: string;
  [QueryableField.EVENT_ID]?: Buffer[];
  [QueryableField.TRANSACTION_HASH]?: string[];
  [QueryableField.CREATED_AT]?: string;
  [QueryableField.CREATED_AT_HEIGHT]?: string[];
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string | undefined;
  [QueryableField.CREATED_BEFORE_OR_AT]?: string | undefined;
  [QueryableField.CREATED_AFTER_HEIGHT]?: string | undefined;
  [QueryableField.CREATED_AFTER]?: string | undefined;
}

export interface OraclePriceQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.MARKET_ID]?: number[];
  [QueryableField.SPOT_PRICE]?: string[];
  [QueryableField.PNL_PRICE]?: string[];
  [QueryableField.EFFECTIVE_AT]?: string;
  [QueryableField.EFFECTIVE_AT_HEIGHT]?: string;
  [QueryableField.EFFECTIVE_BEFORE_OR_AT]?: string;
  [QueryableField.EFFECTIVE_BEFORE_OR_AT_HEIGHT]?: string;
}

export interface MarketQueryConfig extends QueryConfig {
  [QueryableField.ID]?: number[];
  [QueryableField.PAIR]?: string[];
}

export interface CandleQueryConfig extends QueryConfig {
  [QueryableField.ID]?: number[];
  [QueryableField.TICKER]?: string[];
  [QueryableField.RESOLUTION]?: CandleResolution;
  [QueryableField.FROM_ISO]?: IsoString;
  [QueryableField.TO_ISO]?: IsoString;
}

export interface PnlTicksQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.SUBACCOUNT_ID]?: string[];
  [QueryableField.CREATED_AT]?: string;
  [QueryableField.BLOCK_HEIGHT]?: string;
  [QueryableField.BLOCK_TIME]?: string;
  [QueryableField.CREATED_BEFORE_OR_AT]?: string;
  [QueryableField.CREATED_BEFORE_OR_AT_BLOCK_HEIGHT]?: string;
  [QueryableField.CREATED_ON_OR_AFTER]?: string;
  [QueryableField.CREATED_ON_OR_AFTER_BLOCK_HEIGHT]?: string;
}

export interface FundingIndexUpdatesQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
  [QueryableField.PERPETUAL_ID]?: string[];
  [QueryableField.EVENT_ID]?: Buffer;
  [QueryableField.EFFECTIVE_AT]?: string;
  [QueryableField.EFFECTIVE_AT_HEIGHT]?: string;
  [QueryableField.EFFECTIVE_BEFORE_OR_AT]?: string;
  [QueryableField.EFFECTIVE_BEFORE_OR_AT_HEIGHT]?: string;
}

export interface LiquidityTiersQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[];
}

export interface ComplianceDataQueryConfig extends QueryConfig {
  [QueryableField.ADDRESS]?: string[];
  [QueryableField.UPDATED_BEFORE_OR_AT]?: string;
  [QueryableField.PROVIDER]?: string;
  [QueryableField.BLOCKED]?: boolean;
}
