import { PageRequest, PageRequestAmino, PageRequestSDKType, PageResponse, PageResponseAmino, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { Subaccount, SubaccountAmino, SubaccountSDKType } from "./subaccount";
import { BinaryReader, BinaryWriter } from "../../binary";
/** QueryGetSubaccountRequest is request type for the Query RPC method. */
export interface QueryGetSubaccountRequest {
  owner: string;
  number: number;
}
export interface QueryGetSubaccountRequestProtoMsg {
  typeUrl: "/dydxprotocol.subaccounts.QueryGetSubaccountRequest";
  value: Uint8Array;
}
/** QueryGetSubaccountRequest is request type for the Query RPC method. */
export interface QueryGetSubaccountRequestAmino {
  owner?: string;
  number?: number;
}
export interface QueryGetSubaccountRequestAminoMsg {
  type: "/dydxprotocol.subaccounts.QueryGetSubaccountRequest";
  value: QueryGetSubaccountRequestAmino;
}
/** QueryGetSubaccountRequest is request type for the Query RPC method. */
export interface QueryGetSubaccountRequestSDKType {
  owner: string;
  number: number;
}
/** QuerySubaccountResponse is response type for the Query RPC method. */
export interface QuerySubaccountResponse {
  subaccount: Subaccount;
}
export interface QuerySubaccountResponseProtoMsg {
  typeUrl: "/dydxprotocol.subaccounts.QuerySubaccountResponse";
  value: Uint8Array;
}
/** QuerySubaccountResponse is response type for the Query RPC method. */
export interface QuerySubaccountResponseAmino {
  subaccount?: SubaccountAmino;
}
export interface QuerySubaccountResponseAminoMsg {
  type: "/dydxprotocol.subaccounts.QuerySubaccountResponse";
  value: QuerySubaccountResponseAmino;
}
/** QuerySubaccountResponse is response type for the Query RPC method. */
export interface QuerySubaccountResponseSDKType {
  subaccount: SubaccountSDKType;
}
/** QueryAllSubaccountRequest is request type for the Query RPC method. */
export interface QueryAllSubaccountRequest {
  pagination?: PageRequest;
}
export interface QueryAllSubaccountRequestProtoMsg {
  typeUrl: "/dydxprotocol.subaccounts.QueryAllSubaccountRequest";
  value: Uint8Array;
}
/** QueryAllSubaccountRequest is request type for the Query RPC method. */
export interface QueryAllSubaccountRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllSubaccountRequestAminoMsg {
  type: "/dydxprotocol.subaccounts.QueryAllSubaccountRequest";
  value: QueryAllSubaccountRequestAmino;
}
/** QueryAllSubaccountRequest is request type for the Query RPC method. */
export interface QueryAllSubaccountRequestSDKType {
  pagination?: PageRequestSDKType;
}
/** QuerySubaccountAllResponse is response type for the Query RPC method. */
export interface QuerySubaccountAllResponse {
  subaccount: Subaccount[];
  pagination?: PageResponse;
}
export interface QuerySubaccountAllResponseProtoMsg {
  typeUrl: "/dydxprotocol.subaccounts.QuerySubaccountAllResponse";
  value: Uint8Array;
}
/** QuerySubaccountAllResponse is response type for the Query RPC method. */
export interface QuerySubaccountAllResponseAmino {
  subaccount?: SubaccountAmino[];
  pagination?: PageResponseAmino;
}
export interface QuerySubaccountAllResponseAminoMsg {
  type: "/dydxprotocol.subaccounts.QuerySubaccountAllResponse";
  value: QuerySubaccountAllResponseAmino;
}
/** QuerySubaccountAllResponse is response type for the Query RPC method. */
export interface QuerySubaccountAllResponseSDKType {
  subaccount: SubaccountSDKType[];
  pagination?: PageResponseSDKType;
}
function createBaseQueryGetSubaccountRequest(): QueryGetSubaccountRequest {
  return {
    owner: "",
    number: 0
  };
}
export const QueryGetSubaccountRequest = {
  typeUrl: "/dydxprotocol.subaccounts.QueryGetSubaccountRequest",
  encode(message: QueryGetSubaccountRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.owner !== "") {
      writer.uint32(10).string(message.owner);
    }
    if (message.number !== 0) {
      writer.uint32(16).uint32(message.number);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryGetSubaccountRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetSubaccountRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.owner = reader.string();
          break;
        case 2:
          message.number = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryGetSubaccountRequest>): QueryGetSubaccountRequest {
    const message = createBaseQueryGetSubaccountRequest();
    message.owner = object.owner ?? "";
    message.number = object.number ?? 0;
    return message;
  },
  fromAmino(object: QueryGetSubaccountRequestAmino): QueryGetSubaccountRequest {
    const message = createBaseQueryGetSubaccountRequest();
    if (object.owner !== undefined && object.owner !== null) {
      message.owner = object.owner;
    }
    if (object.number !== undefined && object.number !== null) {
      message.number = object.number;
    }
    return message;
  },
  toAmino(message: QueryGetSubaccountRequest): QueryGetSubaccountRequestAmino {
    const obj: any = {};
    obj.owner = message.owner;
    obj.number = message.number;
    return obj;
  },
  fromAminoMsg(object: QueryGetSubaccountRequestAminoMsg): QueryGetSubaccountRequest {
    return QueryGetSubaccountRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryGetSubaccountRequestProtoMsg): QueryGetSubaccountRequest {
    return QueryGetSubaccountRequest.decode(message.value);
  },
  toProto(message: QueryGetSubaccountRequest): Uint8Array {
    return QueryGetSubaccountRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryGetSubaccountRequest): QueryGetSubaccountRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.subaccounts.QueryGetSubaccountRequest",
      value: QueryGetSubaccountRequest.encode(message).finish()
    };
  }
};
function createBaseQuerySubaccountResponse(): QuerySubaccountResponse {
  return {
    subaccount: Subaccount.fromPartial({})
  };
}
export const QuerySubaccountResponse = {
  typeUrl: "/dydxprotocol.subaccounts.QuerySubaccountResponse",
  encode(message: QuerySubaccountResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.subaccount !== undefined) {
      Subaccount.encode(message.subaccount, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QuerySubaccountResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySubaccountResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.subaccount = Subaccount.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QuerySubaccountResponse>): QuerySubaccountResponse {
    const message = createBaseQuerySubaccountResponse();
    message.subaccount = object.subaccount !== undefined && object.subaccount !== null ? Subaccount.fromPartial(object.subaccount) : undefined;
    return message;
  },
  fromAmino(object: QuerySubaccountResponseAmino): QuerySubaccountResponse {
    const message = createBaseQuerySubaccountResponse();
    if (object.subaccount !== undefined && object.subaccount !== null) {
      message.subaccount = Subaccount.fromAmino(object.subaccount);
    }
    return message;
  },
  toAmino(message: QuerySubaccountResponse): QuerySubaccountResponseAmino {
    const obj: any = {};
    obj.subaccount = message.subaccount ? Subaccount.toAmino(message.subaccount) : undefined;
    return obj;
  },
  fromAminoMsg(object: QuerySubaccountResponseAminoMsg): QuerySubaccountResponse {
    return QuerySubaccountResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QuerySubaccountResponseProtoMsg): QuerySubaccountResponse {
    return QuerySubaccountResponse.decode(message.value);
  },
  toProto(message: QuerySubaccountResponse): Uint8Array {
    return QuerySubaccountResponse.encode(message).finish();
  },
  toProtoMsg(message: QuerySubaccountResponse): QuerySubaccountResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.subaccounts.QuerySubaccountResponse",
      value: QuerySubaccountResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllSubaccountRequest(): QueryAllSubaccountRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllSubaccountRequest = {
  typeUrl: "/dydxprotocol.subaccounts.QueryAllSubaccountRequest",
  encode(message: QueryAllSubaccountRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllSubaccountRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllSubaccountRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryAllSubaccountRequest>): QueryAllSubaccountRequest {
    const message = createBaseQueryAllSubaccountRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllSubaccountRequestAmino): QueryAllSubaccountRequest {
    const message = createBaseQueryAllSubaccountRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllSubaccountRequest): QueryAllSubaccountRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllSubaccountRequestAminoMsg): QueryAllSubaccountRequest {
    return QueryAllSubaccountRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllSubaccountRequestProtoMsg): QueryAllSubaccountRequest {
    return QueryAllSubaccountRequest.decode(message.value);
  },
  toProto(message: QueryAllSubaccountRequest): Uint8Array {
    return QueryAllSubaccountRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllSubaccountRequest): QueryAllSubaccountRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.subaccounts.QueryAllSubaccountRequest",
      value: QueryAllSubaccountRequest.encode(message).finish()
    };
  }
};
function createBaseQuerySubaccountAllResponse(): QuerySubaccountAllResponse {
  return {
    subaccount: [],
    pagination: undefined
  };
}
export const QuerySubaccountAllResponse = {
  typeUrl: "/dydxprotocol.subaccounts.QuerySubaccountAllResponse",
  encode(message: QuerySubaccountAllResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.subaccount) {
      Subaccount.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QuerySubaccountAllResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySubaccountAllResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.subaccount.push(Subaccount.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QuerySubaccountAllResponse>): QuerySubaccountAllResponse {
    const message = createBaseQuerySubaccountAllResponse();
    message.subaccount = object.subaccount?.map(e => Subaccount.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QuerySubaccountAllResponseAmino): QuerySubaccountAllResponse {
    const message = createBaseQuerySubaccountAllResponse();
    message.subaccount = object.subaccount?.map(e => Subaccount.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QuerySubaccountAllResponse): QuerySubaccountAllResponseAmino {
    const obj: any = {};
    if (message.subaccount) {
      obj.subaccount = message.subaccount.map(e => e ? Subaccount.toAmino(e) : undefined);
    } else {
      obj.subaccount = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QuerySubaccountAllResponseAminoMsg): QuerySubaccountAllResponse {
    return QuerySubaccountAllResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QuerySubaccountAllResponseProtoMsg): QuerySubaccountAllResponse {
    return QuerySubaccountAllResponse.decode(message.value);
  },
  toProto(message: QuerySubaccountAllResponse): Uint8Array {
    return QuerySubaccountAllResponse.encode(message).finish();
  },
  toProtoMsg(message: QuerySubaccountAllResponse): QuerySubaccountAllResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.subaccounts.QuerySubaccountAllResponse",
      value: QuerySubaccountAllResponse.encode(message).finish()
    };
  }
};