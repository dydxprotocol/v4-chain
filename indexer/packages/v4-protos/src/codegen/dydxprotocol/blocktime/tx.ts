import { DowntimeParams, DowntimeParamsAmino, DowntimeParamsSDKType } from "./params";
import { BinaryReader, BinaryWriter } from "../../binary";
/** MsgUpdateDowntimeParams is the Msg/UpdateDowntimeParams request type. */
export interface MsgUpdateDowntimeParams {
  authority: string;
  /** Defines the parameters to update. All parameters must be supplied. */
  params: DowntimeParams;
}
export interface MsgUpdateDowntimeParamsProtoMsg {
  typeUrl: "/dydxprotocol.blocktime.MsgUpdateDowntimeParams";
  value: Uint8Array;
}
/** MsgUpdateDowntimeParams is the Msg/UpdateDowntimeParams request type. */
export interface MsgUpdateDowntimeParamsAmino {
  authority?: string;
  /** Defines the parameters to update. All parameters must be supplied. */
  params?: DowntimeParamsAmino;
}
export interface MsgUpdateDowntimeParamsAminoMsg {
  type: "/dydxprotocol.blocktime.MsgUpdateDowntimeParams";
  value: MsgUpdateDowntimeParamsAmino;
}
/** MsgUpdateDowntimeParams is the Msg/UpdateDowntimeParams request type. */
export interface MsgUpdateDowntimeParamsSDKType {
  authority: string;
  params: DowntimeParamsSDKType;
}
/**
 * MsgUpdateDowntimeParamsResponse is the Msg/UpdateDowntimeParams response
 * type.
 */
export interface MsgUpdateDowntimeParamsResponse {}
export interface MsgUpdateDowntimeParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse";
  value: Uint8Array;
}
/**
 * MsgUpdateDowntimeParamsResponse is the Msg/UpdateDowntimeParams response
 * type.
 */
export interface MsgUpdateDowntimeParamsResponseAmino {}
export interface MsgUpdateDowntimeParamsResponseAminoMsg {
  type: "/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse";
  value: MsgUpdateDowntimeParamsResponseAmino;
}
/**
 * MsgUpdateDowntimeParamsResponse is the Msg/UpdateDowntimeParams response
 * type.
 */
export interface MsgUpdateDowntimeParamsResponseSDKType {}
function createBaseMsgUpdateDowntimeParams(): MsgUpdateDowntimeParams {
  return {
    authority: "",
    params: DowntimeParams.fromPartial({})
  };
}
export const MsgUpdateDowntimeParams = {
  typeUrl: "/dydxprotocol.blocktime.MsgUpdateDowntimeParams",
  encode(message: MsgUpdateDowntimeParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.params !== undefined) {
      DowntimeParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateDowntimeParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDowntimeParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.params = DowntimeParams.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateDowntimeParams>): MsgUpdateDowntimeParams {
    const message = createBaseMsgUpdateDowntimeParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? DowntimeParams.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateDowntimeParamsAmino): MsgUpdateDowntimeParams {
    const message = createBaseMsgUpdateDowntimeParams();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.params !== undefined && object.params !== null) {
      message.params = DowntimeParams.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: MsgUpdateDowntimeParams): MsgUpdateDowntimeParamsAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.params = message.params ? DowntimeParams.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateDowntimeParamsAminoMsg): MsgUpdateDowntimeParams {
    return MsgUpdateDowntimeParams.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateDowntimeParamsProtoMsg): MsgUpdateDowntimeParams {
    return MsgUpdateDowntimeParams.decode(message.value);
  },
  toProto(message: MsgUpdateDowntimeParams): Uint8Array {
    return MsgUpdateDowntimeParams.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateDowntimeParams): MsgUpdateDowntimeParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.blocktime.MsgUpdateDowntimeParams",
      value: MsgUpdateDowntimeParams.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateDowntimeParamsResponse(): MsgUpdateDowntimeParamsResponse {
  return {};
}
export const MsgUpdateDowntimeParamsResponse = {
  typeUrl: "/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse",
  encode(_: MsgUpdateDowntimeParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateDowntimeParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateDowntimeParamsResponse();
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
  fromPartial(_: Partial<MsgUpdateDowntimeParamsResponse>): MsgUpdateDowntimeParamsResponse {
    const message = createBaseMsgUpdateDowntimeParamsResponse();
    return message;
  },
  fromAmino(_: MsgUpdateDowntimeParamsResponseAmino): MsgUpdateDowntimeParamsResponse {
    const message = createBaseMsgUpdateDowntimeParamsResponse();
    return message;
  },
  toAmino(_: MsgUpdateDowntimeParamsResponse): MsgUpdateDowntimeParamsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateDowntimeParamsResponseAminoMsg): MsgUpdateDowntimeParamsResponse {
    return MsgUpdateDowntimeParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateDowntimeParamsResponseProtoMsg): MsgUpdateDowntimeParamsResponse {
    return MsgUpdateDowntimeParamsResponse.decode(message.value);
  },
  toProto(message: MsgUpdateDowntimeParamsResponse): Uint8Array {
    return MsgUpdateDowntimeParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateDowntimeParamsResponse): MsgUpdateDowntimeParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse",
      value: MsgUpdateDowntimeParamsResponse.encode(message).finish()
    };
  }
};