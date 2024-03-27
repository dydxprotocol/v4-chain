import {
  MILLIS_IN_NANOS,
  SECONDS_IN_MILLIS,
  Timestamp,
} from '@dydxprotocol-indexer/v4-protos';
import Long from 'long';
import { DateTime } from 'luxon';

const defaultDateTime: DateTime = DateTime.utc(2022, 6, 1, 12, 1, 1, 2);
export const defaultTime: Timestamp = {
  seconds: Long.fromValue(Math.floor(defaultDateTime.toSeconds()), true),
  nanos: (defaultDateTime.toMillis() % SECONDS_IN_MILLIS) * MILLIS_IN_NANOS,
};
