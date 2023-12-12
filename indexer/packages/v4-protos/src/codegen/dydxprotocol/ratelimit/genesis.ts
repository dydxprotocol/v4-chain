import { LimitParams, LimitParamsSDKType } from "./limit_params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the ratelimit module's genesis state. */

export interface GenesisState {
  /** limit_params_list defines the list of `LimitParams` at genesis. */
  limitParamsList: LimitParams[];
}
/** GenesisState defines the ratelimit module's genesis state. */

export interface GenesisStateSDKType {
  /** limit_params_list defines the list of `LimitParams` at genesis. */
  limit_params_list: LimitParamsSDKType[];
}

function createBaseGenesisState(): GenesisState {
  return {
    limitParamsList: []
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.limitParamsList) {
      LimitParams.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.limitParamsList.push(LimitParams.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.limitParamsList = object.limitParamsList?.map(e => LimitParams.fromPartial(e)) || [];
    return message;
  }

};