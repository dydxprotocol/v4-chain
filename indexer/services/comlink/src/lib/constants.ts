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
};

export const SPARKLINE_TIME_PERIOD_TO_LIMIT_MAP:
Record<SparklineTimePeriod, number> = {
  [SparklineTimePeriod.ONE_DAY]: 24, // 24 hours in a day
  [SparklineTimePeriod.SEVEN_DAYS]: 7 * 6, // 7 days times (6 * 4 hr candles per day)
};

export const SPARKLINE_TIME_PERIOD_TO_RESOLUTION_MAP:
Record<SparklineTimePeriod, CandleResolution> = {
  [SparklineTimePeriod.ONE_DAY]: CandleResolution.ONE_HOUR,
  [SparklineTimePeriod.SEVEN_DAYS]: CandleResolution.FOUR_HOURS,
};
