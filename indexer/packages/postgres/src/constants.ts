import { CandleMessage_Resolution, ClobPairStatus } from '@dydxprotocol-indexer/v4-protos';

import config from './config';
import AffiliateReferredUsersModel from './models/affiliate-referred-users-model';
import AssetModel from './models/asset-model';
import AssetPositionModel from './models/asset-position-model';
import FillModel from './models/fill-model';
import FundingIndexUpdatesModel from './models/funding-index-updates-model';
import LiquidityTiersModel from './models/liquidity-tiers-model';
import MarketModel from './models/market-model';
import OraclePriceModel from './models/oracle-price-model';
import OrderModel from './models/order-model';
import PerpetualMarketModel from './models/perpetual-market-model';
import PerpetualPositionModel from './models/perpetual-position-model';
import SubaccountModel from './models/subaccount-model';
import TradingRewardModel from './models/trading-reward-model';
import TransferModel from './models/transfer-model';
import VaultModel from './models/vault-model';
import {
  APITimeInForce,
  CandleResolution,
  Options,
  PerpetualMarketStatus,
  TimeInForce,
} from './types';

export const BUFFER_ENCODING_UTF_8: BufferEncoding = 'utf-8';

// Sourced from protocol https://github.com/dydxprotocol/v4/blob/main/lib/constants.go#L6
export const QUOTE_CURRENCY_ATOMIC_RESOLUTION: number = -6;

export const USDC_SYMBOL: string = 'USDC';

export const USDC_DENOM: string = 'ibc/xxx';

export const ZERO_TIME_ISO_8601: string = '1970-01-01T00:00:00.000Z';

export const ONE_MILLION: number = 1_000_000;

export const NUM_SECONDS_IN_CANDLE_RESOLUTIONS: Record<CandleResolution, number> = {
  [CandleResolution.ONE_DAY]: 60 * 60 * 24,
  [CandleResolution.FOUR_HOURS]: 60 * 60 * 4,
  [CandleResolution.ONE_HOUR]: 60 * 60,
  [CandleResolution.THIRTY_MINUTES]: 60 * 30,
  [CandleResolution.FIFTEEN_MINUTES]: 60 * 15,
  [CandleResolution.FIVE_MINUTES]: 60 * 5,
  [CandleResolution.ONE_MINUTE]: 60,
};

export const CANDLE_RESOLUTION_TO_PROTO: Record<CandleResolution, CandleMessage_Resolution> = {
  [CandleResolution.ONE_DAY]: CandleMessage_Resolution.ONE_DAY,
  [CandleResolution.FOUR_HOURS]: CandleMessage_Resolution.FOUR_HOURS,
  [CandleResolution.ONE_HOUR]: CandleMessage_Resolution.ONE_HOUR,
  [CandleResolution.THIRTY_MINUTES]: CandleMessage_Resolution.THIRTY_MINUTES,
  [CandleResolution.FIFTEEN_MINUTES]: CandleMessage_Resolution.FIFTEEN_MINUTES,
  [CandleResolution.FIVE_MINUTES]: CandleMessage_Resolution.FIVE_MINUTES,
  [CandleResolution.ONE_MINUTE]: CandleMessage_Resolution.ONE_MINUTE,
};

export type SpecifiedCandleResolution = Exclude<
  CandleMessage_Resolution,
  CandleMessage_Resolution.UNRECOGNIZED
  >;

export const PROTO_TO_CANDLE_RESOLUTION: Record<SpecifiedCandleResolution, CandleResolution> = {
  [CandleMessage_Resolution.ONE_DAY]: CandleResolution.ONE_DAY,
  [CandleMessage_Resolution.FOUR_HOURS]: CandleResolution.FOUR_HOURS,
  [CandleMessage_Resolution.ONE_HOUR]: CandleResolution.ONE_HOUR,
  [CandleMessage_Resolution.THIRTY_MINUTES]: CandleResolution.THIRTY_MINUTES,
  [CandleMessage_Resolution.FIFTEEN_MINUTES]: CandleResolution.FIFTEEN_MINUTES,
  [CandleMessage_Resolution.FIVE_MINUTES]: CandleResolution.FIVE_MINUTES,
  [CandleMessage_Resolution.ONE_MINUTE]: CandleResolution.ONE_MINUTE,
};

export const USDC_ASSET_ID: string = '0';

// Parts-per-million exponent, 1 million = 10 ^ 6, exponent to convert from ppm units is 10 ^ -6
export const PPM_EXPONENT: number = -6;

// On the protocol, funding is returned in 8 hour increments but we want to store funding in 1 hour
// increments
export const FUNDING_RATE_FROM_PROTOCOL_IN_HOURS: number = 8;

export const TIME_IN_FORCE_TO_API_TIME_IN_FORCE: Record<TimeInForce, APITimeInForce> = {
  [TimeInForce.GTT]: APITimeInForce.GTT,
  [TimeInForce.IOC]: APITimeInForce.IOC,
  [TimeInForce.FOK]: APITimeInForce.FOK,
  [TimeInForce.POST_ONLY]: APITimeInForce.GTT,
};

// A list of models that have sqlToJsonConversions defined.
export const SQL_TO_JSON_DEFINED_MODELS = [
  AffiliateReferredUsersModel,
  AssetModel,
  AssetPositionModel,
  FillModel,
  FundingIndexUpdatesModel,
  LiquidityTiersModel,
  MarketModel,
  OraclePriceModel,
  OrderModel,
  PerpetualMarketModel,
  PerpetualPositionModel,
  SubaccountModel,
  TransferModel,
  TradingRewardModel,
  VaultModel,
];

export type SpecifiedClobPairStatus =
  Exclude<ClobPairStatus, ClobPairStatus.CLOB_PAIR_STATUS_UNSPECIFIED> &
  Exclude<ClobPairStatus, ClobPairStatus.UNRECOGNIZED>;

export const CLOB_STATUS_TO_MARKET_STATUS:
Record<SpecifiedClobPairStatus, PerpetualMarketStatus> = {
  [ClobPairStatus.CLOB_PAIR_STATUS_ACTIVE]: PerpetualMarketStatus.ACTIVE,
  [ClobPairStatus.CLOB_PAIR_STATUS_CANCEL_ONLY]: PerpetualMarketStatus.CANCEL_ONLY,
  [ClobPairStatus.CLOB_PAIR_STATUS_PAUSED]: PerpetualMarketStatus.PAUSED,
  [ClobPairStatus.CLOB_PAIR_STATUS_POST_ONLY]: PerpetualMarketStatus.POST_ONLY,
  [ClobPairStatus.CLOB_PAIR_STATUS_INITIALIZING]: PerpetualMarketStatus.INITIALIZING,
  [ClobPairStatus.CLOB_PAIR_STATUS_FINAL_SETTLEMENT]: PerpetualMarketStatus.FINAL_SETTLEMENT,
};

export const DEFAULT_POSTGRES_OPTIONS : Options = config.USE_READ_REPLICA
  ? {
    readReplica: true,
  } : {};

// The maximum number of parent subaccounts per address.
export const MAX_PARENT_SUBACCOUNTS: number = 128;
// The maximum number of child subaccounts per parent subaccount.
export const CHILD_SUBACCOUNT_MULTIPLIER: number = 1000;

// From https://github.com/dydxprotocol/v4-chain/blob/protocol/v7.0.0-dev0/protocol/app/module_accounts_test.go#L41
export const MEGAVAULT_MODULE_ADDRESS: string = 'dydx18tkxrnrkqc2t0lr3zxr5g6a4hdvqksylxqje4r';
// Generated from the module address + subaccount number 0.
export const MEGAVAULT_SUBACCOUNT_ID: string = 'c7169f81-0c80-54c5-a41f-9cbb6a538fdf';
