import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * RewardShare stores the relative weight of rewards that each address is
 * entitled to.
 */

export interface RewardShare {
  address: string;
  weight: Uint8Array;
}
/**
 * RewardShare stores the relative weight of rewards that each address is
 * entitled to.
 */

export interface RewardShareSDKType {
  address: string;
  weight: Uint8Array;
}

function createBaseRewardShare(): RewardShare {
  return {
    address: "",
    weight: new Uint8Array()
  };
}

export const RewardShare = {
  encode(message: RewardShare, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    if (message.weight.length !== 0) {
      writer.uint32(18).bytes(message.weight);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RewardShare {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRewardShare();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        case 2:
          message.weight = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<RewardShare>): RewardShare {
    const message = createBaseRewardShare();
    message.address = object.address ?? "";
    message.weight = object.weight ?? new Uint8Array();
    return message;
  }

};