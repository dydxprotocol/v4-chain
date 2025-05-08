/* eslint-disable no-console */
import { CandleResolution, NUM_SECONDS_IN_CANDLE_RESOLUTIONS } from '@dydxprotocol-indexer/postgres';
import { DateTime } from 'luxon';
import yargs from 'yargs';

export function calculateNormalizedCandleStartTime(
  time: DateTime,
  resolution: CandleResolution,
): DateTime {
  const epochSeconds: number = Math.floor(time.toUTC().toSeconds());
  const normalizedTimeSeconds: number = epochSeconds - (
    epochSeconds % NUM_SECONDS_IN_CANDLE_RESOLUTIONS[resolution]
  );

  return DateTime.fromSeconds(normalizedTimeSeconds).toUTC();
}

// Get normalized boundaries for a given time and resolutions
function getNormalizedBoundaries(time: string, resolutions: CandleResolution[]): void {
  const date = DateTime.fromISO(time);

  resolutions.forEach((resolution: CandleResolution) => {
    const startTime = calculateNormalizedCandleStartTime(date, resolution);
    const endTime = startTime.plus({ seconds: NUM_SECONDS_IN_CANDLE_RESOLUTIONS[resolution] });

    console.log(`Resolution: ${resolution}, Start Time: ${startTime.toISO()}, End Time: ${endTime.toISO()}`);
  });
}

const resolutions: CandleResolution[] = [
  CandleResolution.ONE_DAY,
  CandleResolution.FOUR_HOURS,
  CandleResolution.ONE_HOUR,
  CandleResolution.THIRTY_MINUTES,
  CandleResolution.FIFTEEN_MINUTES,
  CandleResolution.FIVE_MINUTES,
  CandleResolution.ONE_MINUTE,
];

const args = yargs.options({
  time: {
    type: 'string',
    alias: 't',
    description: 'Time to compute normalized boundaries for, e.g. 2024-02-28T10:01:36.17+00:00',
    required: true,
  },
}).argv;

getNormalizedBoundaries(args.time, resolutions);
