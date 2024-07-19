import config from '../config';

/**
 * Generate a delay duration for the start of a roundtable loop.
 * @param {number} intervalMs: time in milliseconds of roundtable loop
 * @return {number} start delay in milliseconds
 */
export function generateRandomStartDelayMs(intervalMs: number): number {
  const maxDelayMs: number = Math.max(
    config.MAX_START_DELAY_MS,
    config.MAX_START_DELAY_FRACTION_OF_INTERVAL * intervalMs,
  );
  return Math.floor(Math.random() * maxDelayMs);
}
