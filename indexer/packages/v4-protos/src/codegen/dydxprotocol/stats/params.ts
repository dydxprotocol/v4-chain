import { Duration, DurationSDKType } from "../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Params defines the parameters for x/stats module. */

export interface Params {
  /** The desired number of seconds in the look-back window. */
  windowDuration?: Duration;
}
/** Params defines the parameters for x/stats module. */

export interface ParamsSDKType {
  /** The desired number of seconds in the look-back window. */
  window_duration?: DurationSDKType;
}

function createBaseParams(): Params {
  return {
    windowDuration: undefined
  };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.windowDuration !== undefined) {
      Duration.encode(message.windowDuration, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.windowDuration = Duration.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Params>): Params {
    const message = createBaseParams();
    message.windowDuration = object.windowDuration !== undefined && object.windowDuration !== null ? Duration.fromPartial(object.windowDuration) : undefined;
    return message;
  }

};