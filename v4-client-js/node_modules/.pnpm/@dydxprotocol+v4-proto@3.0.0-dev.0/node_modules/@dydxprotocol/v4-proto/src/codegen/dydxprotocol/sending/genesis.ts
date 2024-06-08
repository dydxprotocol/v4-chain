import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the sending module's genesis state. */

export interface GenesisState {}
/** GenesisState defines the sending module's genesis state. */

export interface GenesisStateSDKType {}

function createBaseGenesisState(): GenesisState {
  return {};
}

export const GenesisState = {
  encode(_: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    return message;
  }

};