import { BinaryReader, BinaryWriter } from "../../binary";
import { bytesFromBase64, base64FromBytes } from "../../helpers";
/**
 * RewardShare stores the relative weight of rewards that each address is
 * entitled to.
 */
export interface RewardShare {
  address: string;
  weight: Uint8Array;
}
export interface RewardShareProtoMsg {
  typeUrl: "/dydxprotocol.rewards.RewardShare";
  value: Uint8Array;
}
/**
 * RewardShare stores the relative weight of rewards that each address is
 * entitled to.
 */
export interface RewardShareAmino {
  address?: string;
  weight?: string;
}
export interface RewardShareAminoMsg {
  type: "/dydxprotocol.rewards.RewardShare";
  value: RewardShareAmino;
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
  typeUrl: "/dydxprotocol.rewards.RewardShare",
  encode(message: RewardShare, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.weight.length !== 0) {
      writer.uint32(18).bytes(message.weight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): RewardShare {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<RewardShare>): RewardShare {
    const message = createBaseRewardShare();
    message.address = object.address ?? "";
    message.weight = object.weight ?? new Uint8Array();
    return message;
  },
  fromAmino(object: RewardShareAmino): RewardShare {
    const message = createBaseRewardShare();
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    }
    if (object.weight !== undefined && object.weight !== null) {
      message.weight = bytesFromBase64(object.weight);
    }
    return message;
  },
  toAmino(message: RewardShare): RewardShareAmino {
    const obj: any = {};
    obj.address = message.address;
    obj.weight = message.weight ? base64FromBytes(message.weight) : undefined;
    return obj;
  },
  fromAminoMsg(object: RewardShareAminoMsg): RewardShare {
    return RewardShare.fromAmino(object.value);
  },
  fromProtoMsg(message: RewardShareProtoMsg): RewardShare {
    return RewardShare.decode(message.value);
  },
  toProto(message: RewardShare): Uint8Array {
    return RewardShare.encode(message).finish();
  },
  toProtoMsg(message: RewardShare): RewardShareProtoMsg {
    return {
      typeUrl: "/dydxprotocol.rewards.RewardShare",
      value: RewardShare.encode(message).finish()
    };
  }
};