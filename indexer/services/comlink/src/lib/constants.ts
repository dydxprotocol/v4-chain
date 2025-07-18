import {
  CandleResolution, PositionSide, USDC_ASSET_ID, USDC_SYMBOL,
} from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';

import { AssetPositionResponseObject, SparklineTimePeriod } from '../types';

export const ZERO: Big = new Big(0);
export const ONE: Big = new Big(1.0);

export const ZERO_USDC_POSITION: AssetPositionResponseObject = {
  size: '0',
  symbol: USDC_SYMBOL,
  side: PositionSide.LONG,
  assetId: USDC_ASSET_ID,
  subaccountNumber: 0,
};

export const SPARKLINE_TIME_PERIOD_TO_RESOLUTION_MAP:
Record<SparklineTimePeriod, CandleResolution> = {
  [SparklineTimePeriod.ONE_DAY]: CandleResolution.ONE_HOUR,
  [SparklineTimePeriod.SEVEN_DAYS]: CandleResolution.FOUR_HOURS,
};

export const ONE_DAY_MS = 24 * 60 * 60 * 1000;
export const SEVEN_DAYS_MS = 7 * ONE_DAY_MS;

export const SPARKLINE_TIME_PERIOD_TO_LOOKBACK_MAP
: Record<SparklineTimePeriod, number> = {
  [SparklineTimePeriod.ONE_DAY]: ONE_DAY_MS,
  [SparklineTimePeriod.SEVEN_DAYS]: SEVEN_DAYS_MS,
};

export const DYDX_ADDRESS_PREFIX: string = 'dydx';

export const WRITE_REQUEST_TTL_SECONDS: number = 30;
