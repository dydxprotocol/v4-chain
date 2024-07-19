// Per datadog: A sample rate of 1 sends metrics 100% of the time, while a sample rate of 0
// sends metrics 0% of the time.
export const STATS_NO_SAMPLING: number = 1;

export const STATS_FUNCTION_NAME: string = 'generic_function';

export const ONE_SECOND_IN_MILLISECONDS: number = 1000;
export const FIVE_SECONDS_IN_MILLISECONDS: number = 5 * ONE_SECOND_IN_MILLISECONDS;
export const TEN_SECONDS_IN_MILLISECONDS: number = 10 * ONE_SECOND_IN_MILLISECONDS;
export const THIRTY_SECONDS_IN_MILLISECONDS: number = 30 * ONE_SECOND_IN_MILLISECONDS;
export const ONE_MINUTE_IN_MILLISECONDS: number = 60 * ONE_SECOND_IN_MILLISECONDS;
export const FIVE_MINUTES_IN_MILLISECONDS: number = 5 * ONE_MINUTE_IN_MILLISECONDS;
export const TEN_MINUTES_IN_MILLISECONDS: number = 10 * ONE_MINUTE_IN_MILLISECONDS;
export const ONE_HOUR_IN_MILLISECONDS: number = 60 * ONE_MINUTE_IN_MILLISECONDS;
export const FOUR_HOURS_IN_MILLISECONDS: number = 4 * ONE_HOUR_IN_MILLISECONDS;
export const ONE_DAY_IN_MILLISECONDS: number = 24 * ONE_HOUR_IN_MILLISECONDS;
