import { PerpetualFeeParams, PerpetualFeeParamsAmino, PerpetualFeeParamsSDKType, PerpetualFeeTier, PerpetualFeeTierAmino, PerpetualFeeTierSDKType } from "./params";
import { BinaryReader, BinaryWriter } from "../../binary";
/**
 * QueryPerpetualFeeParamsRequest is a request type for the PerpetualFeeParams
 * RPC method.
 */
export interface QueryPerpetualFeeParamsRequest {}
export interface QueryPerpetualFeeParamsRequestProtoMsg {
  typeUrl: "/dydxprotocol.feetiers.QueryPerpetualFeeParamsRequest";
  value: Uint8Array;
}
/**
 * QueryPerpetualFeeParamsRequest is a request type for the PerpetualFeeParams
 * RPC method.
 */
export interface QueryPerpetualFeeParamsRequestAmino {}
export interface QueryPerpetualFeeParamsRequestAminoMsg {
  type: "/dydxprotocol.feetiers.QueryPerpetualFeeParamsRequest";
  value: QueryPerpetualFeeParamsRequestAmino;
}
/**
 * QueryPerpetualFeeParamsRequest is a request type for the PerpetualFeeParams
 * RPC method.
 */
export interface QueryPerpetualFeeParamsRequestSDKType {}
/**
 * QueryPerpetualFeeParamsResponse is a response type for the PerpetualFeeParams
 * RPC method.
 */
export interface QueryPerpetualFeeParamsResponse {
  params: PerpetualFeeParams;
}
export interface QueryPerpetualFeeParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.feetiers.QueryPerpetualFeeParamsResponse";
  value: Uint8Array;
}
/**
 * QueryPerpetualFeeParamsResponse is a response type for the PerpetualFeeParams
 * RPC method.
 */
export interface QueryPerpetualFeeParamsResponseAmino {
  params?: PerpetualFeeParamsAmino;
}
export interface QueryPerpetualFeeParamsResponseAminoMsg {
  type: "/dydxprotocol.feetiers.QueryPerpetualFeeParamsResponse";
  value: QueryPerpetualFeeParamsResponseAmino;
}
/**
 * QueryPerpetualFeeParamsResponse is a response type for the PerpetualFeeParams
 * RPC method.
 */
export interface QueryPerpetualFeeParamsResponseSDKType {
  params: PerpetualFeeParamsSDKType;
}
/** QueryUserFeeTierRequest is a request type for the UserFeeTier RPC method. */
export interface QueryUserFeeTierRequest {
  user: string;
}
export interface QueryUserFeeTierRequestProtoMsg {
  typeUrl: "/dydxprotocol.feetiers.QueryUserFeeTierRequest";
  value: Uint8Array;
}
/** QueryUserFeeTierRequest is a request type for the UserFeeTier RPC method. */
export interface QueryUserFeeTierRequestAmino {
  user?: string;
}
export interface QueryUserFeeTierRequestAminoMsg {
  type: "/dydxprotocol.feetiers.QueryUserFeeTierRequest";
  value: QueryUserFeeTierRequestAmino;
}
/** QueryUserFeeTierRequest is a request type for the UserFeeTier RPC method. */
export interface QueryUserFeeTierRequestSDKType {
  user: string;
}
/** QueryUserFeeTierResponse is a request type for the UserFeeTier RPC method. */
export interface QueryUserFeeTierResponse {
  /** Index of the fee tier in the list queried from PerpetualFeeParams. */
  index: number;
  tier?: PerpetualFeeTier;
}
export interface QueryUserFeeTierResponseProtoMsg {
  typeUrl: "/dydxprotocol.feetiers.QueryUserFeeTierResponse";
  value: Uint8Array;
}
/** QueryUserFeeTierResponse is a request type for the UserFeeTier RPC method. */
export interface QueryUserFeeTierResponseAmino {
  /** Index of the fee tier in the list queried from PerpetualFeeParams. */
  index?: number;
  tier?: PerpetualFeeTierAmino;
}
export interface QueryUserFeeTierResponseAminoMsg {
  type: "/dydxprotocol.feetiers.QueryUserFeeTierResponse";
  value: QueryUserFeeTierResponseAmino;
}
/** QueryUserFeeTierResponse is a request type for the UserFeeTier RPC method. */
export interface QueryUserFeeTierResponseSDKType {
  index: number;
  tier?: PerpetualFeeTierSDKType;
}
function createBaseQueryPerpetualFeeParamsRequest(): QueryPerpetualFeeParamsRequest {
  return {};
}
export const QueryPerpetualFeeParamsRequest = {
  typeUrl: "/dydxprotocol.feetiers.QueryPerpetualFeeParamsRequest",
  encode(_: QueryPerpetualFeeParamsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryPerpetualFeeParamsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPerpetualFeeParamsRequest();
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
  fromPartial(_: Partial<QueryPerpetualFeeParamsRequest>): QueryPerpetualFeeParamsRequest {
    const message = createBaseQueryPerpetualFeeParamsRequest();
    return message;
  },
  fromAmino(_: QueryPerpetualFeeParamsRequestAmino): QueryPerpetualFeeParamsRequest {
    const message = createBaseQueryPerpetualFeeParamsRequest();
    return message;
  },
  toAmino(_: QueryPerpetualFeeParamsRequest): QueryPerpetualFeeParamsRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryPerpetualFeeParamsRequestAminoMsg): QueryPerpetualFeeParamsRequest {
    return QueryPerpetualFeeParamsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryPerpetualFeeParamsRequestProtoMsg): QueryPerpetualFeeParamsRequest {
    return QueryPerpetualFeeParamsRequest.decode(message.value);
  },
  toProto(message: QueryPerpetualFeeParamsRequest): Uint8Array {
    return QueryPerpetualFeeParamsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryPerpetualFeeParamsRequest): QueryPerpetualFeeParamsRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.feetiers.QueryPerpetualFeeParamsRequest",
      value: QueryPerpetualFeeParamsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryPerpetualFeeParamsResponse(): QueryPerpetualFeeParamsResponse {
  return {
    params: PerpetualFeeParams.fromPartial({})
  };
}
export const QueryPerpetualFeeParamsResponse = {
  typeUrl: "/dydxprotocol.feetiers.QueryPerpetualFeeParamsResponse",
  encode(message: QueryPerpetualFeeParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.params !== undefined) {
      PerpetualFeeParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryPerpetualFeeParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPerpetualFeeParamsResponse();
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
  fromPartial(object: Partial<QueryPerpetualFeeParamsResponse>): QueryPerpetualFeeParamsResponse {
    const message = createBaseQueryPerpetualFeeParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? PerpetualFeeParams.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: QueryPerpetualFeeParamsResponseAmino): QueryPerpetualFeeParamsResponse {
    const message = createBaseQueryPerpetualFeeParamsResponse();
    if (object.params !== undefined && object.params !== null) {
      message.params = PerpetualFeeParams.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: QueryPerpetualFeeParamsResponse): QueryPerpetualFeeParamsResponseAmino {
    const obj: any = {};
    obj.params = message.params ? PerpetualFeeParams.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryPerpetualFeeParamsResponseAminoMsg): QueryPerpetualFeeParamsResponse {
    return QueryPerpetualFeeParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryPerpetualFeeParamsResponseProtoMsg): QueryPerpetualFeeParamsResponse {
    return QueryPerpetualFeeParamsResponse.decode(message.value);
  },
  toProto(message: QueryPerpetualFeeParamsResponse): Uint8Array {
    return QueryPerpetualFeeParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryPerpetualFeeParamsResponse): QueryPerpetualFeeParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.feetiers.QueryPerpetualFeeParamsResponse",
      value: QueryPerpetualFeeParamsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryUserFeeTierRequest(): QueryUserFeeTierRequest {
  return {
    user: ""
  };
}
export const QueryUserFeeTierRequest = {
  typeUrl: "/dydxprotocol.feetiers.QueryUserFeeTierRequest",
  encode(message: QueryUserFeeTierRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.user !== "") {
      writer.uint32(10).string(message.user);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryUserFeeTierRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryUserFeeTierRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.user = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryUserFeeTierRequest>): QueryUserFeeTierRequest {
    const message = createBaseQueryUserFeeTierRequest();
    message.user = object.user ?? "";
    return message;
  },
  fromAmino(object: QueryUserFeeTierRequestAmino): QueryUserFeeTierRequest {
    const message = createBaseQueryUserFeeTierRequest();
    if (object.user !== undefined && object.user !== null) {
      message.user = object.user;
    }
    return message;
  },
  toAmino(message: QueryUserFeeTierRequest): QueryUserFeeTierRequestAmino {
    const obj: any = {};
    obj.user = message.user;
    return obj;
  },
  fromAminoMsg(object: QueryUserFeeTierRequestAminoMsg): QueryUserFeeTierRequest {
    return QueryUserFeeTierRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryUserFeeTierRequestProtoMsg): QueryUserFeeTierRequest {
    return QueryUserFeeTierRequest.decode(message.value);
  },
  toProto(message: QueryUserFeeTierRequest): Uint8Array {
    return QueryUserFeeTierRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryUserFeeTierRequest): QueryUserFeeTierRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.feetiers.QueryUserFeeTierRequest",
      value: QueryUserFeeTierRequest.encode(message).finish()
    };
  }
};
function createBaseQueryUserFeeTierResponse(): QueryUserFeeTierResponse {
  return {
    index: 0,
    tier: undefined
  };
}
export const QueryUserFeeTierResponse = {
  typeUrl: "/dydxprotocol.feetiers.QueryUserFeeTierResponse",
  encode(message: QueryUserFeeTierResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.index !== 0) {
      writer.uint32(8).uint32(message.index);
    }
    if (message.tier !== undefined) {
      PerpetualFeeTier.encode(message.tier, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryUserFeeTierResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryUserFeeTierResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.index = reader.uint32();
          break;
        case 2:
          message.tier = PerpetualFeeTier.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryUserFeeTierResponse>): QueryUserFeeTierResponse {
    const message = createBaseQueryUserFeeTierResponse();
    message.index = object.index ?? 0;
    message.tier = object.tier !== undefined && object.tier !== null ? PerpetualFeeTier.fromPartial(object.tier) : undefined;
    return message;
  },
  fromAmino(object: QueryUserFeeTierResponseAmino): QueryUserFeeTierResponse {
    const message = createBaseQueryUserFeeTierResponse();
    if (object.index !== undefined && object.index !== null) {
      message.index = object.index;
    }
    if (object.tier !== undefined && object.tier !== null) {
      message.tier = PerpetualFeeTier.fromAmino(object.tier);
    }
    return message;
  },
  toAmino(message: QueryUserFeeTierResponse): QueryUserFeeTierResponseAmino {
    const obj: any = {};
    obj.index = message.index;
    obj.tier = message.tier ? PerpetualFeeTier.toAmino(message.tier) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryUserFeeTierResponseAminoMsg): QueryUserFeeTierResponse {
    return QueryUserFeeTierResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryUserFeeTierResponseProtoMsg): QueryUserFeeTierResponse {
    return QueryUserFeeTierResponse.decode(message.value);
  },
  toProto(message: QueryUserFeeTierResponse): Uint8Array {
    return QueryUserFeeTierResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryUserFeeTierResponse): QueryUserFeeTierResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.feetiers.QueryUserFeeTierResponse",
      value: QueryUserFeeTierResponse.encode(message).finish()
    };
  }
};