import { SynchronyParams, SynchronyParamsSDKType, DowntimeParams, DowntimeParamsSDKType } from "./params";
import { BlockInfo, BlockInfoSDKType, AllDowntimeInfo, AllDowntimeInfoSDKType } from "./blocktime";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QuerySynchronyParamsRequest is a request type for the SynchronyParams */

export interface QuerySynchronyParamsRequest {}
/** QuerySynchronyParamsRequest is a request type for the SynchronyParams */

export interface QuerySynchronyParamsRequestSDKType {}
/** QuerySynchronyParamsResponse is a response type for the SynchronyParams */

export interface QuerySynchronyParamsResponse {
  params?: SynchronyParams;
}
/** QuerySynchronyParamsResponse is a response type for the SynchronyParams */

export interface QuerySynchronyParamsResponseSDKType {
  params?: SynchronyParamsSDKType;
}
/**
 * QueryDowntimeParamsRequest is a request type for the DowntimeParams
 * RPC method.
 */

export interface QueryDowntimeParamsRequest {}
/**
 * QueryDowntimeParamsRequest is a request type for the DowntimeParams
 * RPC method.
 */

export interface QueryDowntimeParamsRequestSDKType {}
/**
 * QueryDowntimeParamsResponse is a response type for the DowntimeParams
 * RPC method.
 */

export interface QueryDowntimeParamsResponse {
  params?: DowntimeParams;
}
/**
 * QueryDowntimeParamsResponse is a response type for the DowntimeParams
 * RPC method.
 */

export interface QueryDowntimeParamsResponseSDKType {
  params?: DowntimeParamsSDKType;
}
/**
 * QueryPreviousBlockInfoRequest is a request type for the PreviousBlockInfo
 * RPC method.
 */

export interface QueryPreviousBlockInfoRequest {}
/**
 * QueryPreviousBlockInfoRequest is a request type for the PreviousBlockInfo
 * RPC method.
 */

export interface QueryPreviousBlockInfoRequestSDKType {}
/**
 * QueryPreviousBlockInfoResponse is a request type for the PreviousBlockInfo
 * RPC method.
 */

export interface QueryPreviousBlockInfoResponse {
  /**
   * QueryPreviousBlockInfoResponse is a request type for the PreviousBlockInfo
   * RPC method.
   */
  info?: BlockInfo;
}
/**
 * QueryPreviousBlockInfoResponse is a request type for the PreviousBlockInfo
 * RPC method.
 */

export interface QueryPreviousBlockInfoResponseSDKType {
  /**
   * QueryPreviousBlockInfoResponse is a request type for the PreviousBlockInfo
   * RPC method.
   */
  info?: BlockInfoSDKType;
}
/**
 * QueryAllDowntimeInfoRequest is a request type for the AllDowntimeInfo
 * RPC method.
 */

export interface QueryAllDowntimeInfoRequest {}
/**
 * QueryAllDowntimeInfoRequest is a request type for the AllDowntimeInfo
 * RPC method.
 */

export interface QueryAllDowntimeInfoRequestSDKType {}
/**
 * QueryAllDowntimeInfoResponse is a request type for the AllDowntimeInfo
 * RPC method.
 */

export interface QueryAllDowntimeInfoResponse {
  /**
   * QueryAllDowntimeInfoResponse is a request type for the AllDowntimeInfo
   * RPC method.
   */
  info?: AllDowntimeInfo;
}
/**
 * QueryAllDowntimeInfoResponse is a request type for the AllDowntimeInfo
 * RPC method.
 */

export interface QueryAllDowntimeInfoResponseSDKType {
  /**
   * QueryAllDowntimeInfoResponse is a request type for the AllDowntimeInfo
   * RPC method.
   */
  info?: AllDowntimeInfoSDKType;
}

function createBaseQuerySynchronyParamsRequest(): QuerySynchronyParamsRequest {
  return {};
}

export const QuerySynchronyParamsRequest = {
  encode(_: QuerySynchronyParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySynchronyParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySynchronyParamsRequest();

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

  fromPartial(_: DeepPartial<QuerySynchronyParamsRequest>): QuerySynchronyParamsRequest {
    const message = createBaseQuerySynchronyParamsRequest();
    return message;
  }

};

function createBaseQuerySynchronyParamsResponse(): QuerySynchronyParamsResponse {
  return {
    params: undefined
  };
}

export const QuerySynchronyParamsResponse = {
  encode(message: QuerySynchronyParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      SynchronyParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySynchronyParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySynchronyParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = SynchronyParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QuerySynchronyParamsResponse>): QuerySynchronyParamsResponse {
    const message = createBaseQuerySynchronyParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? SynchronyParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryDowntimeParamsRequest(): QueryDowntimeParamsRequest {
  return {};
}

export const QueryDowntimeParamsRequest = {
  encode(_: QueryDowntimeParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDowntimeParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDowntimeParamsRequest();

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

  fromPartial(_: DeepPartial<QueryDowntimeParamsRequest>): QueryDowntimeParamsRequest {
    const message = createBaseQueryDowntimeParamsRequest();
    return message;
  }

};

function createBaseQueryDowntimeParamsResponse(): QueryDowntimeParamsResponse {
  return {
    params: undefined
  };
}

export const QueryDowntimeParamsResponse = {
  encode(message: QueryDowntimeParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      DowntimeParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryDowntimeParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDowntimeParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = DowntimeParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryDowntimeParamsResponse>): QueryDowntimeParamsResponse {
    const message = createBaseQueryDowntimeParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? DowntimeParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryPreviousBlockInfoRequest(): QueryPreviousBlockInfoRequest {
  return {};
}

export const QueryPreviousBlockInfoRequest = {
  encode(_: QueryPreviousBlockInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPreviousBlockInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPreviousBlockInfoRequest();

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

  fromPartial(_: DeepPartial<QueryPreviousBlockInfoRequest>): QueryPreviousBlockInfoRequest {
    const message = createBaseQueryPreviousBlockInfoRequest();
    return message;
  }

};

function createBaseQueryPreviousBlockInfoResponse(): QueryPreviousBlockInfoResponse {
  return {
    info: undefined
  };
}

export const QueryPreviousBlockInfoResponse = {
  encode(message: QueryPreviousBlockInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.info !== undefined) {
      BlockInfo.encode(message.info, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPreviousBlockInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPreviousBlockInfoResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.info = BlockInfo.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryPreviousBlockInfoResponse>): QueryPreviousBlockInfoResponse {
    const message = createBaseQueryPreviousBlockInfoResponse();
    message.info = object.info !== undefined && object.info !== null ? BlockInfo.fromPartial(object.info) : undefined;
    return message;
  }

};

function createBaseQueryAllDowntimeInfoRequest(): QueryAllDowntimeInfoRequest {
  return {};
}

export const QueryAllDowntimeInfoRequest = {
  encode(_: QueryAllDowntimeInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllDowntimeInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllDowntimeInfoRequest();

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

  fromPartial(_: DeepPartial<QueryAllDowntimeInfoRequest>): QueryAllDowntimeInfoRequest {
    const message = createBaseQueryAllDowntimeInfoRequest();
    return message;
  }

};

function createBaseQueryAllDowntimeInfoResponse(): QueryAllDowntimeInfoResponse {
  return {
    info: undefined
  };
}

export const QueryAllDowntimeInfoResponse = {
  encode(message: QueryAllDowntimeInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.info !== undefined) {
      AllDowntimeInfo.encode(message.info, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllDowntimeInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllDowntimeInfoResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.info = AllDowntimeInfo.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllDowntimeInfoResponse>): QueryAllDowntimeInfoResponse {
    const message = createBaseQueryAllDowntimeInfoResponse();
    message.info = object.info !== undefined && object.info !== null ? AllDowntimeInfo.fromPartial(object.info) : undefined;
    return message;
  }

};