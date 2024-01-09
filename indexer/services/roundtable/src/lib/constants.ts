import Big from 'big.js';
import { DateTimeOptions } from 'luxon';

export const REDIS_VALUE: string = 'TIMEOUT NOT EXPIRED';

// Per datadog: A sample rate of 1 sends metrics 100% of the time, while a sample rate of 0
// sends metrics 0% of the time.
export const STATS_NO_SAMPLING: number = 1;

export const ZERO: Big = new Big(0);

export const USDC_ASSET_ID: string = '0';

export const UTC_OPTIONS: DateTimeOptions = { zone: 'utc' };
