import * as _m0 from 'protobufjs/minimal';

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

/**
 * Pass this to a generated `decode` instead of the raw buffer. Given a buffer, generated code
 * builds a base `Reader`, which decodes strings in a JS loop; for a Node `Buffer` this returns a
 * `BufferReader`, which uses the native UTF-8 decoder and is several times faster on large strings.
 */
export function createProtoReader(buffer: Uint8Array): _m0.Reader {
  return _m0.Reader.create(buffer);
}
