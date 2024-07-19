/**
 * Floor of the date to the timeInMilliseconds.
 * For example with ONE_MINUTE_IN_MILLISECONDS, 12:02:33 will be rounded to 12:02:00
 */
export function floorDate(date: Date, timeInMilliseconds: number): Date {
  return new Date(
    Math.floor(
      date.getTime() / timeInMilliseconds,
    ) * timeInMilliseconds,
  );
}

/**
 * Ceiling of the date to the timeInMilliseconds.
 * For example with ONE_MINUTE_IN_MILLISECONDS, 12:02:33 will be rounded to 12:03:00
 */
export function ceilingDate(date: Date, timeInMilliseconds: number): Date {
  return new Date(
    Math.ceil(
      date.getTime() / timeInMilliseconds,
    ) * timeInMilliseconds,
  );
}
