import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { EpochInfo, EpochInfoSDKType } from "./epoch_info";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QueryGetEpochInfoRequest is request type for the GetEpochInfo RPC method. */

export interface QueryGetEpochInfoRequest {
  /** QueryGetEpochInfoRequest is request type for the GetEpochInfo RPC method. */
  name: string;
}
/** QueryGetEpochInfoRequest is request type for the GetEpochInfo RPC method. */

export interface QueryGetEpochInfoRequestSDKType {
  /** QueryGetEpochInfoRequest is request type for the GetEpochInfo RPC method. */
  name: string;
}
/** QueryEpochInfoResponse is response type for the GetEpochInfo RPC method. */

export interface QueryEpochInfoResponse {
  epochInfo?: EpochInfo;
}
/** QueryEpochInfoResponse is response type for the GetEpochInfo RPC method. */

export interface QueryEpochInfoResponseSDKType {
  epoch_info?: EpochInfoSDKType;
}
/** QueryAllEpochInfoRequest is request type for the AllEpochInfo RPC method. */

export interface QueryAllEpochInfoRequest {
  pagination?: PageRequest;
}
/** QueryAllEpochInfoRequest is request type for the AllEpochInfo RPC method. */

export interface QueryAllEpochInfoRequestSDKType {
  pagination?: PageRequestSDKType;
}
/** QueryEpochInfoAllResponse is response type for the AllEpochInfo RPC method. */

export interface QueryEpochInfoAllResponse {
  epochInfo: EpochInfo[];
  pagination?: PageResponse;
}
/** QueryEpochInfoAllResponse is response type for the AllEpochInfo RPC method. */

export interface QueryEpochInfoAllResponseSDKType {
  epoch_info: EpochInfoSDKType[];
  pagination?: PageResponseSDKType;
}

function createBaseQueryGetEpochInfoRequest(): QueryGetEpochInfoRequest {
  return {
    name: ""
  };
}

export const QueryGetEpochInfoRequest = {
  encode(message: QueryGetEpochInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetEpochInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetEpochInfoRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryGetEpochInfoRequest>): QueryGetEpochInfoRequest {
    const message = createBaseQueryGetEpochInfoRequest();
    message.name = object.name ?? "";
    return message;
  }

};

function createBaseQueryEpochInfoResponse(): QueryEpochInfoResponse {
  return {
    epochInfo: undefined
  };
}

export const QueryEpochInfoResponse = {
  encode(message: QueryEpochInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.epochInfo !== undefined) {
      EpochInfo.encode(message.epochInfo, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryEpochInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryEpochInfoResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.epochInfo = EpochInfo.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryEpochInfoResponse>): QueryEpochInfoResponse {
    const message = createBaseQueryEpochInfoResponse();
    message.epochInfo = object.epochInfo !== undefined && object.epochInfo !== null ? EpochInfo.fromPartial(object.epochInfo) : undefined;
    return message;
  }

};

function createBaseQueryAllEpochInfoRequest(): QueryAllEpochInfoRequest {
  return {
    pagination: undefined
  };
}

export const QueryAllEpochInfoRequest = {
  encode(message: QueryAllEpochInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllEpochInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllEpochInfoRequest();

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

  fromPartial(object: DeepPartial<QueryAllEpochInfoRequest>): QueryAllEpochInfoRequest {
    const message = createBaseQueryAllEpochInfoRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryEpochInfoAllResponse(): QueryEpochInfoAllResponse {
  return {
    epochInfo: [],
    pagination: undefined
  };
}

export const QueryEpochInfoAllResponse = {
  encode(message: QueryEpochInfoAllResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.epochInfo) {
      EpochInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryEpochInfoAllResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryEpochInfoAllResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.epochInfo.push(EpochInfo.decode(reader, reader.uint32()));
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

  fromPartial(object: DeepPartial<QueryEpochInfoAllResponse>): QueryEpochInfoAllResponse {
    const message = createBaseQueryEpochInfoAllResponse();
    message.epochInfo = object.epochInfo?.map(e => EpochInfo.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};