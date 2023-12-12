import { PerpetualFeeParams, PerpetualFeeParamsAmino, PerpetualFeeParamsSDKType } from "./params";
import { BinaryReader, BinaryWriter } from "../../binary";
/** GenesisState defines the feetiers module's genesis state. */
export interface GenesisState {
  /** The parameters for perpetual fees. */
  params: PerpetualFeeParams;
}
export interface GenesisStateProtoMsg {
  typeUrl: "/dydxprotocol.feetiers.GenesisState";
  value: Uint8Array;
}
/** GenesisState defines the feetiers module's genesis state. */
export interface GenesisStateAmino {
  /** The parameters for perpetual fees. */
  params?: PerpetualFeeParamsAmino;
}
export interface GenesisStateAminoMsg {
  type: "/dydxprotocol.feetiers.GenesisState";
  value: GenesisStateAmino;
}
/** GenesisState defines the feetiers module's genesis state. */
export interface GenesisStateSDKType {
  params: PerpetualFeeParamsSDKType;
}
function createBaseGenesisState(): GenesisState {
  return {
    params: PerpetualFeeParams.fromPartial({})
  };
}
export const GenesisState = {
  typeUrl: "/dydxprotocol.feetiers.GenesisState",
  encode(message: GenesisState, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.params !== undefined) {
      PerpetualFeeParams.encode(message.params, writer.uint32(10).fork()).ldelim();
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
          message.params = PerpetualFeeParams.decode(reader, reader.uint32());
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
    message.params = object.params !== undefined && object.params !== null ? PerpetualFeeParams.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: GenesisStateAmino): GenesisState {
    const message = createBaseGenesisState();
    if (object.params !== undefined && object.params !== null) {
      message.params = PerpetualFeeParams.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: GenesisState): GenesisStateAmino {
    const obj: any = {};
    obj.params = message.params ? PerpetualFeeParams.toAmino(message.params) : undefined;
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
      typeUrl: "/dydxprotocol.feetiers.GenesisState",
      value: GenesisState.encode(message).finish()
    };
  }
};