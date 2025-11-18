/* ------- QUERY TYPES ------- */

import { CandleResolution } from './candle-types';
import { FillType, Liquidity } from './fill-types';
import { OrderSide, OrderStatus, OrderType } from './order-types';
import { PerpetualPositionStatus } from './perpetual-position-types';
import { PositionSide } from './position-types';
import { ParentSubaccount } from './subaccount-types';
import { TradingRewardAggregationPeriod } from './trading-reward-aggregation-types';
import { IsoString } from './utility-types';

export enum QueryableField {
  LIMIT = 'limit',
  PAGE = 'page',
  ID = 'id',
  ADDRESS = 'address',
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
  TYPE = 'type',
  STATUS = 'status',
  STATUSES = 'statuses',
  INCLUDE_TYPES = 'includeTypes',
  EXCLUDE_TYPES = 'excludeTypes',
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
  EFFECTIVE_AT_OR_AFTER_HEIGHT = 'effectiveAtOrAfterHeight',
  GOOD_TIL_BLOCK_BEFORE_OR_AT = 'goodTilBlockBeforeOrAt',
  GOOD_TIL_BLOCK_AFTER = 'goodTilBlockAfter',
  GOOD_TIL_BLOCK_TIME_BEFORE_OR_AT = 'goodTilBlockTimeBeforeOrAt',
  GOOD_TIL_BLOCK_TIME_AFTER = 'goodTilBlockTimeAfter',
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
  REASON = 'reason',
  USERNAME = 'username',
  EMAIL = 'email',
  TIMESPAN = 'timeSpan',
  RANK = 'rank',
  AFFILIATE_ADDRESS = 'affiliateAddress',
  REFEREE_ADDRESS = 'refereeAddress',
  KEY = 'key',
  TOKEN = 'token',
  ADDRESS_IN_WALLETS_TABLE = 'addressInWalletsTable',
  PARENT_SUBACCOUNT = 'parentSubaccount',
  DISTINCT_FIELDS = 'distinctFields',
  ZERO_PAYMENTS = 'zeroPayments',
  SUBORG_ID = 'suborg_id',
  SVM_ADDRESS = 'svm_address',
  EVM_ADDRESS = 'evm_address',
  DYDX_ADDRESS = 'dydx_address',
}

export interface QueryConfig {
  [QueryableField.LIMIT]?: number,
  [QueryableField.PAGE]?: number,
}

export interface SubaccountQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.ADDRESS]?: string,
  [QueryableField.SUBACCOUNT_NUMBER]?: number,
  [QueryableField.UPDATED_BEFORE_OR_AT]?: string,
  [QueryableField.UPDATED_ON_OR_AFTER]?: string,
}

export interface SubaccountUsernamesQueryConfig extends QueryConfig {
  [QueryableField.USERNAME]?: string[],
  [QueryableField.SUBACCOUNT_ID]?: string[],
}

export interface WalletQueryConfig extends QueryConfig {
  [QueryableField.ADDRESS]?: string,
}

export interface PerpetualPositionQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.SUBACCOUNT_ID]?: string[],
  [QueryableField.PERPETUAL_ID]?: string[],
  [QueryableField.SIDE]?: PositionSide,
  [QueryableField.STATUS]?: PerpetualPositionStatus[],
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string,
  [QueryableField.CREATED_BEFORE_OR_AT]?: string,
}

export interface OrderQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.SUBACCOUNT_ID]?: string[],
  [QueryableField.CLIENT_ID]?: string,
  [QueryableField.CLOB_PAIR_ID]?: string,
  [QueryableField.SIDE]?: OrderSide,
  [QueryableField.SIZE]?: string,
  [QueryableField.TOTAL_FILLED]?: string,
  [QueryableField.PRICE]?: string,
  [QueryableField.TYPE]?: OrderType,
  [QueryableField.INCLUDE_TYPES]?: OrderType[],
  [QueryableField.EXCLUDE_TYPES]?: OrderType[],
  [QueryableField.STATUSES]?: OrderStatus[],
  [QueryableField.POST_ONLY]?: boolean,
  [QueryableField.REDUCE_ONLY]?: boolean,
  [QueryableField.GOOD_TIL_BLOCK_BEFORE_OR_AT]?: string,
  [QueryableField.GOOD_TIL_BLOCK_AFTER]?: string,
  [QueryableField.GOOD_TIL_BLOCK_TIME_BEFORE_OR_AT]?: string,
  [QueryableField.GOOD_TIL_BLOCK_TIME_AFTER]?: string,
  [QueryableField.ORDER_FLAGS]?: string,
  [QueryableField.CLIENT_METADATA]?: string,
  [QueryableField.TRIGGER_PRICE]?: string,
  [QueryableField.PARENT_SUBACCOUNT]?: ParentSubaccount,
}

export interface PerpetualMarketQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.MARKET_ID]?: number[],
  [QueryableField.LIQUIDITY_TIER_ID]?: number[],
}

export interface FillQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.SUBACCOUNT_ID]?: string[],
  [QueryableField.SIDE]?: OrderSide,
  [QueryableField.LIQUIDITY]?: Liquidity,
  [QueryableField.TYPE]?: FillType,
  [QueryableField.INCLUDE_TYPES]?: FillType[],
  [QueryableField.EXCLUDE_TYPES]?: FillType[],
  [QueryableField.CLOB_PAIR_ID]?: string,
  [QueryableField.EVENT_ID]?: Buffer,
  [QueryableField.TRANSACTION_HASH]?: string,
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string,
  [QueryableField.CREATED_BEFORE_OR_AT]?: string,
  [QueryableField.CREATED_ON_OR_AFTER_HEIGHT]?: string,
  [QueryableField.CREATED_ON_OR_AFTER]?: string,
  [QueryableField.CLIENT_METADATA]?: string,
  [QueryableField.FEE]?: string,
  [QueryableField.PARENT_SUBACCOUNT]?: ParentSubaccount,
}

export interface BlockQueryConfig extends QueryConfig {
  [QueryableField.BLOCK_HEIGHT]?: string[],
  [QueryableField.CREATED_ON_OR_AFTER]?: string,
}

export interface TendermintEventQueryConfig extends QueryConfig {
  [QueryableField.ID]?: Buffer[],
  [QueryableField.BLOCK_HEIGHT]?: string[],
  [QueryableField.TRANSACTION_INDEX]?: number[],
  [QueryableField.EVENT_INDEX]?: number[],
}

export interface TransactionQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.BLOCK_HEIGHT]?: string[],
  [QueryableField.TRANSACTION_INDEX]?: number[],
  [QueryableField.TRANSACTION_HASH]?: string[],
}

export interface AssetQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.SYMBOL]?: string,
  [QueryableField.ATOMIC_RESOLUTION]?: number,
  [QueryableField.HAS_MARKET]?: boolean,
  [QueryableField.MARKET_ID]?: number,
}

export interface AssetPositionQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.ASSET_ID]?: string[],
  [QueryableField.SUBACCOUNT_ID]?: string[],
  [QueryableField.SIZE]?: string,
  [QueryableField.IS_LONG]?: boolean,
}

export interface TransferQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.SENDER_SUBACCOUNT_ID]?: string[],
  [QueryableField.RECIPIENT_SUBACCOUNT_ID]?: string[],
  [QueryableField.SENDER_WALLET_ADDRESS]?: string[],
  [QueryableField.RECIPIENT_WALLET_ADDRESS]?: string[],
  [QueryableField.ASSET_ID]?: string[],
  [QueryableField.SIZE]?: string,
  [QueryableField.EVENT_ID]?: Buffer[],
  [QueryableField.TRANSACTION_HASH]?: string[],
  [QueryableField.CREATED_AT]?: string,
  [QueryableField.CREATED_AT_HEIGHT]?: string[],
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string,
  [QueryableField.CREATED_BEFORE_OR_AT]?: string,
  [QueryableField.CREATED_AFTER]?: string,
  [QueryableField.CREATED_AFTER_HEIGHT]?: string,
}

export interface ToAndFromSubaccountTransferQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.SUBACCOUNT_ID]?: string[],
  [QueryableField.ASSET_ID]?: string[],
  [QueryableField.SIZE]?: string,
  [QueryableField.EVENT_ID]?: Buffer[],
  [QueryableField.TRANSACTION_HASH]?: string[],
  [QueryableField.CREATED_AT]?: string,
  [QueryableField.CREATED_AT_HEIGHT]?: string[],
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string | undefined,
  [QueryableField.CREATED_BEFORE_OR_AT]?: string | undefined,
  [QueryableField.CREATED_AFTER_HEIGHT]?: string | undefined,
  [QueryableField.CREATED_AFTER]?: string | undefined,
}

export interface ParentSubaccountTransferQueryConfig extends QueryConfig {
  [QueryableField.SUBACCOUNT_ID]: string[],
  [QueryableField.LIMIT]?: number,
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string,
  [QueryableField.CREATED_BEFORE_OR_AT]?: string,
  [QueryableField.PAGE]?: number,
}

export interface OraclePriceQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.MARKET_ID]?: number[],
  [QueryableField.PRICE]?: string[],
  [QueryableField.EFFECTIVE_AT]?: string,
  [QueryableField.EFFECTIVE_AT_HEIGHT]?: string,
  [QueryableField.EFFECTIVE_BEFORE_OR_AT]?: string,
  [QueryableField.EFFECTIVE_BEFORE_OR_AT_HEIGHT]?: string,
}

export interface MarketQueryConfig extends QueryConfig {
  [QueryableField.ID]?: number[],
  [QueryableField.PAIR]?: string[],
}

export interface CandleQueryConfig extends QueryConfig {
  [QueryableField.ID]?: number[],
  [QueryableField.TICKER]?: string[],
  [QueryableField.RESOLUTION]?: CandleResolution,
  [QueryableField.FROM_ISO]?: IsoString,
  [QueryableField.TO_ISO]?: IsoString,
}

export interface PnlTicksQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.SUBACCOUNT_ID]?: string[],
  [QueryableField.CREATED_AT]?: string,
  [QueryableField.BLOCK_HEIGHT]?: string,
  [QueryableField.BLOCK_TIME]?: string,
  [QueryableField.CREATED_BEFORE_OR_AT]?: string,
  [QueryableField.CREATED_BEFORE_OR_AT_BLOCK_HEIGHT]?: string,
  [QueryableField.CREATED_ON_OR_AFTER]?: string,
  [QueryableField.CREATED_ON_OR_AFTER_BLOCK_HEIGHT]?: string,
  [QueryableField.PARENT_SUBACCOUNT]?: ParentSubaccount,
}

export interface FundingIndexUpdatesQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
  [QueryableField.PERPETUAL_ID]?: string[],
  [QueryableField.EVENT_ID]?: Buffer,
  [QueryableField.EFFECTIVE_AT]?: string,
  [QueryableField.EFFECTIVE_AT_HEIGHT]?: string,
  [QueryableField.EFFECTIVE_BEFORE_OR_AT]?: string,
  [QueryableField.EFFECTIVE_BEFORE_OR_AT_HEIGHT]?: string,
  [QueryableField.EFFECTIVE_AT_OR_AFTER_HEIGHT]?: string,
  [QueryableField.DISTINCT_FIELDS]?: string[],
}

export interface LiquidityTiersQueryConfig extends QueryConfig {
  [QueryableField.ID]?: string[],
}

export interface ComplianceDataQueryConfig extends QueryConfig {
  [QueryableField.ADDRESS]?: string[],
  [QueryableField.UPDATED_BEFORE_OR_AT]?: string,
  [QueryableField.PROVIDER]?: string,
  [QueryableField.BLOCKED]?: boolean,
  [QueryableField.ADDRESS_IN_WALLETS_TABLE]?: boolean,
}

export interface ComplianceStatusQueryConfig extends QueryConfig {
  [QueryableField.ADDRESS]?: string[],
  [QueryableField.STATUS]?: string[],
  [QueryableField.CREATED_BEFORE_OR_AT]?: string,
  [QueryableField.UPDATED_BEFORE_OR_AT]?: string,
  [QueryableField.REASON]?: string,
}

export interface TradingRewardQueryConfig extends QueryConfig {
  [QueryableField.ADDRESS]?: string,
  [QueryableField.BLOCK_HEIGHT]?: string,
  [QueryableField.BLOCK_TIME_BEFORE_OR_AT]?: IsoString,
  [QueryableField.BLOCK_TIME_AFTER_OR_AT]?: IsoString,
  [QueryableField.BLOCK_TIME_BEFORE]?: IsoString,
  [QueryableField.BLOCK_HEIGHT_BEFORE_OR_AT]?: IsoString,
}

export interface TradingRewardAggregationQueryConfig extends QueryConfig {
  [QueryableField.ADDRESS]?: string,
  [QueryableField.ADDRESSES]?: string[],
  [QueryableField.STARTED_AT_HEIGHT]?: string,
  [QueryableField.STARTED_AT_HEIGHT_OR_AFTER]?: string,
  [QueryableField.PERIOD]?: TradingRewardAggregationPeriod,
  [QueryableField.STARTED_AT_BEFORE_OR_AT]?: IsoString,
  [QueryableField.STARTED_AT_HEIGHT_BEFORE_OR_AT]?: string,
}

export interface AffiliateReferredUsersQueryConfig extends QueryConfig {
  [QueryableField.AFFILIATE_ADDRESS]?: string[],
  [QueryableField.REFEREE_ADDRESS]?: string[],
}

export interface LeaderboardPnlQueryConfig extends QueryConfig {
  [QueryableField.ADDRESS]?: string[],
  [QueryableField.TIMESPAN]?: string[],
  [QueryableField.RANK]?: number[],
}

export interface PersistentCacheQueryConfig extends QueryConfig {
  [QueryableField.KEY]?: string,
}

export interface AffiliateInfoQueryConfig extends QueryConfig {
  [QueryableField.ADDRESS]?: string,
}

export interface FirebaseNotificationTokenQueryConfig extends QueryConfig {
  [QueryableField.ADDRESS]?: string,
  [QueryableField.TOKEN]?: string,
  [QueryableField.UPDATED_BEFORE_OR_AT]?: IsoString,
}

export interface VaultQueryConfig extends QueryConfig {
  [QueryableField.ADDRESS]?: string[],
  [QueryableField.CLOB_PAIR_ID]?: string[],
  [QueryableField.STATUS]?: string[],
}

export interface FundingPaymentsQueryConfig extends QueryConfig {
  [QueryableField.SUBACCOUNT_ID]?: string[],
  [QueryableField.PERPETUAL_ID]?: string[],
  [QueryableField.TICKER]?: string,
  [QueryableField.CREATED_AT_HEIGHT]?: string,
  [QueryableField.CREATED_AT]?: string,
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string,
  [QueryableField.CREATED_BEFORE_OR_AT]?: string,
  [QueryableField.CREATED_ON_OR_AFTER_HEIGHT]?: string,
  [QueryableField.CREATED_ON_OR_AFTER]?: string,
  [QueryableField.PARENT_SUBACCOUNT]?: ParentSubaccount,
  [QueryableField.ZERO_PAYMENTS]?: boolean,
  [QueryableField.DISTINCT_FIELDS]?: string[],
}

export interface PnlQueryConfig extends QueryConfig {
  [QueryableField.SUBACCOUNT_ID]?: string[],
  [QueryableField.CREATED_AT_HEIGHT]?: string,
  [QueryableField.CREATED_AT]?: string,
  [QueryableField.CREATED_BEFORE_OR_AT_HEIGHT]?: string,
  [QueryableField.CREATED_BEFORE_OR_AT]?: string,
  [QueryableField.CREATED_ON_OR_AFTER_HEIGHT]?: string,
  [QueryableField.CREATED_ON_OR_AFTER]?: string,
  [QueryableField.PARENT_SUBACCOUNT]?: ParentSubaccount,
}
