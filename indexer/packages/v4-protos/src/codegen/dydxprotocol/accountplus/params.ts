import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Params defines the parameters for the module. */

export interface Params {
  /**
   * IsSmartAccountActive defines the state of the authenticator.
   * If set to false, the authenticator module will not be used
   * and the classic cosmos sdk authentication will be used instead.
   */
  isSmartAccountActive: boolean;
}
/** Params defines the parameters for the module. */

export interface ParamsSDKType {
  /**
   * IsSmartAccountActive defines the state of the authenticator.
   * If set to false, the authenticator module will not be used
   * and the classic cosmos sdk authentication will be used instead.
   */
  is_smart_account_active: boolean;
}

function createBaseParams(): Params {
  return {
    isSmartAccountActive: false
  };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.isSmartAccountActive === true) {
      writer.uint32(8).bool(message.isSmartAccountActive);
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
          message.isSmartAccountActive = reader.bool();
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
    message.isSmartAccountActive = object.isSmartAccountActive ?? false;
    return message;
  }

};