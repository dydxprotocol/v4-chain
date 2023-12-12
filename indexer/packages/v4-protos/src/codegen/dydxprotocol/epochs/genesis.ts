import { EpochInfo, EpochInfoAmino, EpochInfoSDKType } from "./epoch_info";
import { BinaryReader, BinaryWriter } from "../../binary";
/** GenesisState defines the epochs module's genesis state. */
export interface GenesisState {
  epochInfoList: EpochInfo[];
}
export interface GenesisStateProtoMsg {
  typeUrl: "/dydxprotocol.epochs.GenesisState";
  value: Uint8Array;
}
/** GenesisState defines the epochs module's genesis state. */
export interface GenesisStateAmino {
  epoch_info_list?: EpochInfoAmino[];
}
export interface GenesisStateAminoMsg {
  type: "/dydxprotocol.epochs.GenesisState";
  value: GenesisStateAmino;
}
/** GenesisState defines the epochs module's genesis state. */
export interface GenesisStateSDKType {
  epoch_info_list: EpochInfoSDKType[];
}
function createBaseGenesisState(): GenesisState {
  return {
    epochInfoList: []
  };
}
export const GenesisState = {
  typeUrl: "/dydxprotocol.epochs.GenesisState",
  encode(message: GenesisState, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.epochInfoList) {
      EpochInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.epochInfoList.push(EpochInfo.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<GenesisState>): GenesisState {
    const message = createBaseGenesisState();
    message.epochInfoList = object.epochInfoList?.map(e => EpochInfo.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: GenesisStateAmino): GenesisState {
    const message = createBaseGenesisState();
    message.epochInfoList = object.epoch_info_list?.map(e => EpochInfo.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: GenesisState): GenesisStateAmino {
    const obj: any = {};
    if (message.epochInfoList) {
      obj.epoch_info_list = message.epochInfoList.map(e => e ? EpochInfo.toAmino(e) : undefined);
    } else {
      obj.epoch_info_list = [];
    }
    return obj;
  },
  fromAminoMsg(object: GenesisStateAminoMsg): GenesisState {
    return GenesisState.fromAmino(object.value);
  },
  fromProtoMsg(message: GenesisStateProtoMsg): GenesisState {
    return GenesisState.decode(message.value);
  },
  toProto(message: GenesisState): Uint8Array {
    return GenesisState.encode(message).finish();
  },
  toProtoMsg(message: GenesisState): GenesisStateProtoMsg {
    return {
      typeUrl: "/dydxprotocol.epochs.GenesisState",
      value: GenesisState.encode(message).finish()
    };
  }
};