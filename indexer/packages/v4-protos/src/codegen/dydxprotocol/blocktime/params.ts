import { Duration, DurationSDKType } from "../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** DowntimeParams defines the parameters for downtime. */

export interface DowntimeParams {
  /**
   * Durations tracked for downtime. The durations must be sorted from
   * shortest to longest and must all be positive.
   */
  durations: Duration[];
}
/** DowntimeParams defines the parameters for downtime. */

export interface DowntimeParamsSDKType {
  /**
   * Durations tracked for downtime. The durations must be sorted from
   * shortest to longest and must all be positive.
   */
  durations: DurationSDKType[];
}

function createBaseDowntimeParams(): DowntimeParams {
  return {
    durations: []
  };
}

export const DowntimeParams = {
  encode(message: DowntimeParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.durations) {
      Duration.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DowntimeParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDowntimeParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.durations.push(Duration.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<DowntimeParams>): DowntimeParams {
    const message = createBaseDowntimeParams();
    message.durations = object.durations?.map(e => Duration.fromPartial(e)) || [];
    return message;
  }

};