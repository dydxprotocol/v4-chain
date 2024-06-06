import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines `x/listing`'s genesis state. */

export interface GenesisState {
  /**
   * permissionless_listing_enabled defines whether permissionless listing is
   * enabled.
   */
  permissionlessListingEnabled: boolean;
}
/** GenesisState defines `x/listing`'s genesis state. */

export interface GenesisStateSDKType {
  /**
   * permissionless_listing_enabled defines whether permissionless listing is
   * enabled.
   */
  permissionless_listing_enabled: boolean;
}

function createBaseGenesisState(): GenesisState {
  return {
    permissionlessListingEnabled: false
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.permissionlessListingEnabled === true) {
      writer.uint32(8).bool(message.permissionlessListingEnabled);
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
          message.permissionlessListingEnabled = reader.bool();
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
    message.permissionlessListingEnabled = object.permissionlessListingEnabled ?? false;
    return message;
  }

};