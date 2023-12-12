import { LimitParams, LimitParamsAmino, LimitParamsSDKType } from "./limit_params";
import { BinaryReader, BinaryWriter } from "../../binary";
/** MsgSetLimitParams is the Msg/SetLimitParams request type. */
export interface MsgSetLimitParams {
  authority: string;
  /** Defines the parameters to set. All parameters must be supplied. */
  limitParams?: LimitParams;
}
export interface MsgSetLimitParamsProtoMsg {
  typeUrl: "/dydxprotocol.ibcratelimit.MsgSetLimitParams";
  value: Uint8Array;
}
/** MsgSetLimitParams is the Msg/SetLimitParams request type. */
export interface MsgSetLimitParamsAmino {
  authority?: string;
  /** Defines the parameters to set. All parameters must be supplied. */
  limit_params?: LimitParamsAmino;
}
export interface MsgSetLimitParamsAminoMsg {
  type: "/dydxprotocol.ibcratelimit.MsgSetLimitParams";
  value: MsgSetLimitParamsAmino;
}
/** MsgSetLimitParams is the Msg/SetLimitParams request type. */
export interface MsgSetLimitParamsSDKType {
  authority: string;
  limit_params?: LimitParamsSDKType;
}
/** MsgSetLimitParamsResponse is the Msg/SetLimitParams response type. */
export interface MsgSetLimitParamsResponse {}
export interface MsgSetLimitParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.ibcratelimit.MsgSetLimitParamsResponse";
  value: Uint8Array;
}
/** MsgSetLimitParamsResponse is the Msg/SetLimitParams response type. */
export interface MsgSetLimitParamsResponseAmino {}
export interface MsgSetLimitParamsResponseAminoMsg {
  type: "/dydxprotocol.ibcratelimit.MsgSetLimitParamsResponse";
  value: MsgSetLimitParamsResponseAmino;
}
/** MsgSetLimitParamsResponse is the Msg/SetLimitParams response type. */
export interface MsgSetLimitParamsResponseSDKType {}
/** MsgDeleteLimitParams is the Msg/SetLimitParams request type. */
export interface MsgDeleteLimitParams {
  authority: string;
  /** The denom for which the `LimitParams` should be deleted. */
  denom: string;
}
export interface MsgDeleteLimitParamsProtoMsg {
  typeUrl: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParams";
  value: Uint8Array;
}
/** MsgDeleteLimitParams is the Msg/SetLimitParams request type. */
export interface MsgDeleteLimitParamsAmino {
  authority?: string;
  /** The denom for which the `LimitParams` should be deleted. */
  denom?: string;
}
export interface MsgDeleteLimitParamsAminoMsg {
  type: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParams";
  value: MsgDeleteLimitParamsAmino;
}
/** MsgDeleteLimitParams is the Msg/SetLimitParams request type. */
export interface MsgDeleteLimitParamsSDKType {
  authority: string;
  denom: string;
}
/** MsgDeleteLimitParamsResponse is the Msg/DeleteLimitParams response type. */
export interface MsgDeleteLimitParamsResponse {}
export interface MsgDeleteLimitParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParamsResponse";
  value: Uint8Array;
}
/** MsgDeleteLimitParamsResponse is the Msg/DeleteLimitParams response type. */
export interface MsgDeleteLimitParamsResponseAmino {}
export interface MsgDeleteLimitParamsResponseAminoMsg {
  type: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParamsResponse";
  value: MsgDeleteLimitParamsResponseAmino;
}
/** MsgDeleteLimitParamsResponse is the Msg/DeleteLimitParams response type. */
export interface MsgDeleteLimitParamsResponseSDKType {}
function createBaseMsgSetLimitParams(): MsgSetLimitParams {
  return {
    authority: "",
    limitParams: undefined
  };
}
export const MsgSetLimitParams = {
  typeUrl: "/dydxprotocol.ibcratelimit.MsgSetLimitParams",
  encode(message: MsgSetLimitParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.limitParams !== undefined) {
      LimitParams.encode(message.limitParams, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgSetLimitParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetLimitParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.limitParams = LimitParams.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgSetLimitParams>): MsgSetLimitParams {
    const message = createBaseMsgSetLimitParams();
    message.authority = object.authority ?? "";
    message.limitParams = object.limitParams !== undefined && object.limitParams !== null ? LimitParams.fromPartial(object.limitParams) : undefined;
    return message;
  },
  fromAmino(object: MsgSetLimitParamsAmino): MsgSetLimitParams {
    const message = createBaseMsgSetLimitParams();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.limit_params !== undefined && object.limit_params !== null) {
      message.limitParams = LimitParams.fromAmino(object.limit_params);
    }
    return message;
  },
  toAmino(message: MsgSetLimitParams): MsgSetLimitParamsAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.limit_params = message.limitParams ? LimitParams.toAmino(message.limitParams) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgSetLimitParamsAminoMsg): MsgSetLimitParams {
    return MsgSetLimitParams.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSetLimitParamsProtoMsg): MsgSetLimitParams {
    return MsgSetLimitParams.decode(message.value);
  },
  toProto(message: MsgSetLimitParams): Uint8Array {
    return MsgSetLimitParams.encode(message).finish();
  },
  toProtoMsg(message: MsgSetLimitParams): MsgSetLimitParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.ibcratelimit.MsgSetLimitParams",
      value: MsgSetLimitParams.encode(message).finish()
    };
  }
};
function createBaseMsgSetLimitParamsResponse(): MsgSetLimitParamsResponse {
  return {};
}
export const MsgSetLimitParamsResponse = {
  typeUrl: "/dydxprotocol.ibcratelimit.MsgSetLimitParamsResponse",
  encode(_: MsgSetLimitParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgSetLimitParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetLimitParamsResponse();
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
  fromPartial(_: Partial<MsgSetLimitParamsResponse>): MsgSetLimitParamsResponse {
    const message = createBaseMsgSetLimitParamsResponse();
    return message;
  },
  fromAmino(_: MsgSetLimitParamsResponseAmino): MsgSetLimitParamsResponse {
    const message = createBaseMsgSetLimitParamsResponse();
    return message;
  },
  toAmino(_: MsgSetLimitParamsResponse): MsgSetLimitParamsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgSetLimitParamsResponseAminoMsg): MsgSetLimitParamsResponse {
    return MsgSetLimitParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSetLimitParamsResponseProtoMsg): MsgSetLimitParamsResponse {
    return MsgSetLimitParamsResponse.decode(message.value);
  },
  toProto(message: MsgSetLimitParamsResponse): Uint8Array {
    return MsgSetLimitParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgSetLimitParamsResponse): MsgSetLimitParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.ibcratelimit.MsgSetLimitParamsResponse",
      value: MsgSetLimitParamsResponse.encode(message).finish()
    };
  }
};
function createBaseMsgDeleteLimitParams(): MsgDeleteLimitParams {
  return {
    authority: "",
    denom: ""
  };
}
export const MsgDeleteLimitParams = {
  typeUrl: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParams",
  encode(message: MsgDeleteLimitParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.denom !== "") {
      writer.uint32(18).string(message.denom);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgDeleteLimitParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteLimitParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.denom = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgDeleteLimitParams>): MsgDeleteLimitParams {
    const message = createBaseMsgDeleteLimitParams();
    message.authority = object.authority ?? "";
    message.denom = object.denom ?? "";
    return message;
  },
  fromAmino(object: MsgDeleteLimitParamsAmino): MsgDeleteLimitParams {
    const message = createBaseMsgDeleteLimitParams();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.denom !== undefined && object.denom !== null) {
      message.denom = object.denom;
    }
    return message;
  },
  toAmino(message: MsgDeleteLimitParams): MsgDeleteLimitParamsAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.denom = message.denom;
    return obj;
  },
  fromAminoMsg(object: MsgDeleteLimitParamsAminoMsg): MsgDeleteLimitParams {
    return MsgDeleteLimitParams.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgDeleteLimitParamsProtoMsg): MsgDeleteLimitParams {
    return MsgDeleteLimitParams.decode(message.value);
  },
  toProto(message: MsgDeleteLimitParams): Uint8Array {
    return MsgDeleteLimitParams.encode(message).finish();
  },
  toProtoMsg(message: MsgDeleteLimitParams): MsgDeleteLimitParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParams",
      value: MsgDeleteLimitParams.encode(message).finish()
    };
  }
};
function createBaseMsgDeleteLimitParamsResponse(): MsgDeleteLimitParamsResponse {
  return {};
}
export const MsgDeleteLimitParamsResponse = {
  typeUrl: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParamsResponse",
  encode(_: MsgDeleteLimitParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgDeleteLimitParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteLimitParamsResponse();
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
  fromPartial(_: Partial<MsgDeleteLimitParamsResponse>): MsgDeleteLimitParamsResponse {
    const message = createBaseMsgDeleteLimitParamsResponse();
    return message;
  },
  fromAmino(_: MsgDeleteLimitParamsResponseAmino): MsgDeleteLimitParamsResponse {
    const message = createBaseMsgDeleteLimitParamsResponse();
    return message;
  },
  toAmino(_: MsgDeleteLimitParamsResponse): MsgDeleteLimitParamsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgDeleteLimitParamsResponseAminoMsg): MsgDeleteLimitParamsResponse {
    return MsgDeleteLimitParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgDeleteLimitParamsResponseProtoMsg): MsgDeleteLimitParamsResponse {
    return MsgDeleteLimitParamsResponse.decode(message.value);
  },
  toProto(message: MsgDeleteLimitParamsResponse): Uint8Array {
    return MsgDeleteLimitParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgDeleteLimitParamsResponse): MsgDeleteLimitParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.ibcratelimit.MsgDeleteLimitParamsResponse",
      value: MsgDeleteLimitParamsResponse.encode(message).finish()
    };
  }
};