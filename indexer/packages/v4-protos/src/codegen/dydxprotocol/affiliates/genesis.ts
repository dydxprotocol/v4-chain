import { AffiliateTiers, AffiliateTiersSDKType } from "./affiliates";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines generis state of `x/affiliates` */

export interface GenesisState {
  /** The list of affiliate tiers */
  affiliateTiers?: AffiliateTiers;
}
/** GenesisState defines generis state of `x/affiliates` */

export interface GenesisStateSDKType {
  /** The list of affiliate tiers */
  affiliate_tiers?: AffiliateTiersSDKType;
}

function createBaseGenesisState(): GenesisState {
  return {
    affiliateTiers: undefined
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.affiliateTiers !== undefined) {
      AffiliateTiers.encode(message.affiliateTiers, writer.uint32(10).fork()).ldelim();
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
          message.affiliateTiers = AffiliateTiers.decode(reader, reader.uint32());
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
    message.affiliateTiers = object.affiliateTiers !== undefined && object.affiliateTiers !== null ? AffiliateTiers.fromPartial(object.affiliateTiers) : undefined;
    return message;
  }

};