import { Timestamp } from './codegen/google/protobuf/timestamp';

export const MILLIS_IN_NANOS: number = 1_000_000;
export const SECONDS_IN_MILLIS: number = 1_000;
export function protoTimestampToDate(
  protoTime: Timestamp,
): Date {
  const timeInMillis: number = Number(protoTime.seconds) * SECONDS_IN_MILLIS +
    Math.floor(protoTime.nanos / MILLIS_IN_NANOS);

  return new Date(timeInMillis);
}
