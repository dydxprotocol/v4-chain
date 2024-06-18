import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines `x/listing`'s genesis state. */

export interface GenesisState {
  /**
   * hard_cap_for_markets is the hard cap for the number of markets that can be
   * listed
   */
  hardCapForMarkets: number;
}
/** GenesisState defines `x/listing`'s genesis state. */

export interface GenesisStateSDKType {
  /**
   * hard_cap_for_markets is the hard cap for the number of markets that can be
   * listed
   */
  hard_cap_for_markets: number;
}

function createBaseGenesisState(): GenesisState {
  return {
    hardCapForMarkets: 0
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.hardCapForMarkets !== 0) {
      writer.uint32(8).uint32(message.hardCapForMarkets);
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
          message.hardCapForMarkets = reader.uint32();
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
    message.hardCapForMarkets = object.hardCapForMarkets ?? 0;
    return message;
  }

};