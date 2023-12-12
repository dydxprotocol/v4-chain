import { LimitParams, LimitParamsAmino, LimitParamsSDKType } from "./limit_params";
import { BinaryReader, BinaryWriter } from "../../binary";
/** GenesisState defines the ibcratelimit module's genesis state. */
export interface GenesisState {
  /** limit_params_list defines the list of `LimitParams` at genesis. */
  limitParamsList: LimitParams[];
}
export interface GenesisStateProtoMsg {
  typeUrl: "/dydxprotocol.ibcratelimit.GenesisState";
  value: Uint8Array;
}
/** GenesisState defines the ibcratelimit module's genesis state. */
export interface GenesisStateAmino {
  /** limit_params_list defines the list of `LimitParams` at genesis. */
  limit_params_list?: LimitParamsAmino[];
}
export interface GenesisStateAminoMsg {
  type: "/dydxprotocol.ibcratelimit.GenesisState";
  value: GenesisStateAmino;
}
/** GenesisState defines the ibcratelimit module's genesis state. */
export interface GenesisStateSDKType {
  limit_params_list: LimitParamsSDKType[];
}
function createBaseGenesisState(): GenesisState {
  return {
    limitParamsList: []
  };
}
export const GenesisState = {
  typeUrl: "/dydxprotocol.ibcratelimit.GenesisState",
  encode(message: GenesisState, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.limitParamsList) {
      LimitParams.encode(v!, writer.uint32(10).fork()).ldelim();
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
          message.limitParamsList.push(LimitParams.decode(reader, reader.uint32()));
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
    message.limitParamsList = object.limitParamsList?.map(e => LimitParams.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: GenesisStateAmino): GenesisState {
    const message = createBaseGenesisState();
    message.limitParamsList = object.limit_params_list?.map(e => LimitParams.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: GenesisState): GenesisStateAmino {
    const obj: any = {};
    if (message.limitParamsList) {
      obj.limit_params_list = message.limitParamsList.map(e => e ? LimitParams.toAmino(e) : undefined);
    } else {
      obj.limit_params_list = [];
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
      typeUrl: "/dydxprotocol.ibcratelimit.GenesisState",
      value: GenesisState.encode(message).finish()
    };
  }
};