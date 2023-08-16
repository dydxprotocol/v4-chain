import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { Perpetual, PerpetualSDKType } from "./perpetual";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Queries a Perpetual by id. */

export interface QueryPerpetualRequest {
  /** Queries a Perpetual by id. */
  id: number;
}
/** Queries a Perpetual by id. */

export interface QueryPerpetualRequestSDKType {
  /** Queries a Perpetual by id. */
  id: number;
}
/** QueryPerpetualResponse is response type for the Perpetual RPC method. */

export interface QueryPerpetualResponse {
  perpetual?: Perpetual;
}
/** QueryPerpetualResponse is response type for the Perpetual RPC method. */

export interface QueryPerpetualResponseSDKType {
  perpetual?: PerpetualSDKType;
}
/** Queries a list of Perpetual items. */

export interface QueryAllPerpetualsRequest {
  pagination?: PageRequest;
}
/** Queries a list of Perpetual items. */

export interface QueryAllPerpetualsRequestSDKType {
  pagination?: PageRequestSDKType;
}
/** QueryAllPerpetualsResponse is response type for the AllPerpetuals RPC method. */

export interface QueryAllPerpetualsResponse {
  perpetual: Perpetual[];
  pagination?: PageResponse;
}
/** QueryAllPerpetualsResponse is response type for the AllPerpetuals RPC method. */

export interface QueryAllPerpetualsResponseSDKType {
  perpetual: PerpetualSDKType[];
  pagination?: PageResponseSDKType;
}

function createBaseQueryPerpetualRequest(): QueryPerpetualRequest {
  return {
    id: 0
  };
}

export const QueryPerpetualRequest = {
  encode(message: QueryPerpetualRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPerpetualRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPerpetualRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryPerpetualRequest>): QueryPerpetualRequest {
    const message = createBaseQueryPerpetualRequest();
    message.id = object.id ?? 0;
    return message;
  }

};

function createBaseQueryPerpetualResponse(): QueryPerpetualResponse {
  return {
    perpetual: undefined
  };
}

export const QueryPerpetualResponse = {
  encode(message: QueryPerpetualResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.perpetual !== undefined) {
      Perpetual.encode(message.perpetual, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPerpetualResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPerpetualResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.perpetual = Perpetual.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryPerpetualResponse>): QueryPerpetualResponse {
    const message = createBaseQueryPerpetualResponse();
    message.perpetual = object.perpetual !== undefined && object.perpetual !== null ? Perpetual.fromPartial(object.perpetual) : undefined;
    return message;
  }

};

function createBaseQueryAllPerpetualsRequest(): QueryAllPerpetualsRequest {
  return {
    pagination: undefined
  };
}

export const QueryAllPerpetualsRequest = {
  encode(message: QueryAllPerpetualsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllPerpetualsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllPerpetualsRequest();

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

  fromPartial(object: DeepPartial<QueryAllPerpetualsRequest>): QueryAllPerpetualsRequest {
    const message = createBaseQueryAllPerpetualsRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryAllPerpetualsResponse(): QueryAllPerpetualsResponse {
  return {
    perpetual: [],
    pagination: undefined
  };
}

export const QueryAllPerpetualsResponse = {
  encode(message: QueryAllPerpetualsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.perpetual) {
      Perpetual.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllPerpetualsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllPerpetualsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.perpetual.push(Perpetual.decode(reader, reader.uint32()));
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

  fromPartial(object: DeepPartial<QueryAllPerpetualsResponse>): QueryAllPerpetualsResponse {
    const message = createBaseQueryAllPerpetualsResponse();
    message.perpetual = object.perpetual?.map(e => Perpetual.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};