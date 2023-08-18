import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { SubaccountId, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { ClobPair, ClobPairSDKType } from "./clob_pair";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QueryGetClobPairRequest is request type for the QueryRPC method. */

export interface QueryGetClobPairRequest {
  /** QueryGetClobPairRequest is request type for the QueryRPC method. */
  id: number;
}
/** QueryGetClobPairRequest is request type for the QueryRPC method. */

export interface QueryGetClobPairRequestSDKType {
  /** QueryGetClobPairRequest is request type for the QueryRPC method. */
  id: number;
}
/** QueryClobPairResponse is response type for the QueryRPC method. */

export interface QueryClobPairResponse {
  clobPair?: ClobPair;
}
/** QueryClobPairResponse is response type for the QueryRPC method. */

export interface QueryClobPairResponseSDKType {
  clob_pair?: ClobPairSDKType;
}
/** QueryAllClobPairRequest is request type for the QueryRPC method. */

export interface QueryAllClobPairRequest {
  pagination?: PageRequest;
}
/** QueryAllClobPairRequest is request type for the QueryRPC method. */

export interface QueryAllClobPairRequestSDKType {
  pagination?: PageRequestSDKType;
}
/** QueryClobPairAllResponse is request type for the QueryRPC method. */

export interface QueryClobPairAllResponse {
  clobPair: ClobPair[];
  pagination?: PageResponse;
}
/** QueryClobPairAllResponse is request type for the QueryRPC method. */

export interface QueryClobPairAllResponseSDKType {
  clob_pair: ClobPairSDKType[];
  pagination?: PageResponseSDKType;
}
/**
 * AreSubaccountsLiquidatableRequest is a request message used to check whether
 * the given subaccounts are liquidatable.
 * The subaccount ids should not contain duplicates.
 */

export interface AreSubaccountsLiquidatableRequest {
  subaccountIds: SubaccountId[];
}
/**
 * AreSubaccountsLiquidatableRequest is a request message used to check whether
 * the given subaccounts are liquidatable.
 * The subaccount ids should not contain duplicates.
 */

export interface AreSubaccountsLiquidatableRequestSDKType {
  subaccount_ids: SubaccountIdSDKType[];
}
/**
 * AreSubaccountsLiquidatableResponse is a response message that contains the
 * liquidation status for each subaccount.
 */

export interface AreSubaccountsLiquidatableResponse {
  results: AreSubaccountsLiquidatableResponse_Result[];
}
/**
 * AreSubaccountsLiquidatableResponse is a response message that contains the
 * liquidation status for each subaccount.
 */

export interface AreSubaccountsLiquidatableResponseSDKType {
  results: AreSubaccountsLiquidatableResponse_ResultSDKType[];
}
/** Result returns whether a subaccount should be liquidated. */

export interface AreSubaccountsLiquidatableResponse_Result {
  subaccountId?: SubaccountId;
  isLiquidatable: boolean;
}
/** Result returns whether a subaccount should be liquidated. */

export interface AreSubaccountsLiquidatableResponse_ResultSDKType {
  subaccount_id?: SubaccountIdSDKType;
  is_liquidatable: boolean;
}

function createBaseQueryGetClobPairRequest(): QueryGetClobPairRequest {
  return {
    id: 0
  };
}

export const QueryGetClobPairRequest = {
  encode(message: QueryGetClobPairRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetClobPairRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryGetClobPairRequest();

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

  fromPartial(object: DeepPartial<QueryGetClobPairRequest>): QueryGetClobPairRequest {
    const message = createBaseQueryGetClobPairRequest();
    message.id = object.id ?? 0;
    return message;
  }

};

function createBaseQueryClobPairResponse(): QueryClobPairResponse {
  return {
    clobPair: undefined
  };
}

export const QueryClobPairResponse = {
  encode(message: QueryClobPairResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobPair !== undefined) {
      ClobPair.encode(message.clobPair, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClobPairResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClobPairResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPair = ClobPair.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryClobPairResponse>): QueryClobPairResponse {
    const message = createBaseQueryClobPairResponse();
    message.clobPair = object.clobPair !== undefined && object.clobPair !== null ? ClobPair.fromPartial(object.clobPair) : undefined;
    return message;
  }

};

function createBaseQueryAllClobPairRequest(): QueryAllClobPairRequest {
  return {
    pagination: undefined
  };
}

export const QueryAllClobPairRequest = {
  encode(message: QueryAllClobPairRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllClobPairRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllClobPairRequest();

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

  fromPartial(object: DeepPartial<QueryAllClobPairRequest>): QueryAllClobPairRequest {
    const message = createBaseQueryAllClobPairRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryClobPairAllResponse(): QueryClobPairAllResponse {
  return {
    clobPair: [],
    pagination: undefined
  };
}

export const QueryClobPairAllResponse = {
  encode(message: QueryClobPairAllResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.clobPair) {
      ClobPair.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryClobPairAllResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryClobPairAllResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPair.push(ClobPair.decode(reader, reader.uint32()));
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

  fromPartial(object: DeepPartial<QueryClobPairAllResponse>): QueryClobPairAllResponse {
    const message = createBaseQueryClobPairAllResponse();
    message.clobPair = object.clobPair?.map(e => ClobPair.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseAreSubaccountsLiquidatableRequest(): AreSubaccountsLiquidatableRequest {
  return {
    subaccountIds: []
  };
}

export const AreSubaccountsLiquidatableRequest = {
  encode(message: AreSubaccountsLiquidatableRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.subaccountIds) {
      SubaccountId.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AreSubaccountsLiquidatableRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAreSubaccountsLiquidatableRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountIds.push(SubaccountId.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AreSubaccountsLiquidatableRequest>): AreSubaccountsLiquidatableRequest {
    const message = createBaseAreSubaccountsLiquidatableRequest();
    message.subaccountIds = object.subaccountIds?.map(e => SubaccountId.fromPartial(e)) || [];
    return message;
  }

};

function createBaseAreSubaccountsLiquidatableResponse(): AreSubaccountsLiquidatableResponse {
  return {
    results: []
  };
}

export const AreSubaccountsLiquidatableResponse = {
  encode(message: AreSubaccountsLiquidatableResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.results) {
      AreSubaccountsLiquidatableResponse_Result.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AreSubaccountsLiquidatableResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAreSubaccountsLiquidatableResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.results.push(AreSubaccountsLiquidatableResponse_Result.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AreSubaccountsLiquidatableResponse>): AreSubaccountsLiquidatableResponse {
    const message = createBaseAreSubaccountsLiquidatableResponse();
    message.results = object.results?.map(e => AreSubaccountsLiquidatableResponse_Result.fromPartial(e)) || [];
    return message;
  }

};

function createBaseAreSubaccountsLiquidatableResponse_Result(): AreSubaccountsLiquidatableResponse_Result {
  return {
    subaccountId: undefined,
    isLiquidatable: false
  };
}

export const AreSubaccountsLiquidatableResponse_Result = {
  encode(message: AreSubaccountsLiquidatableResponse_Result, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.subaccountId !== undefined) {
      SubaccountId.encode(message.subaccountId, writer.uint32(10).fork()).ldelim();
    }

    if (message.isLiquidatable === true) {
      writer.uint32(16).bool(message.isLiquidatable);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AreSubaccountsLiquidatableResponse_Result {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAreSubaccountsLiquidatableResponse_Result();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountId = SubaccountId.decode(reader, reader.uint32());
          break;

        case 2:
          message.isLiquidatable = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AreSubaccountsLiquidatableResponse_Result>): AreSubaccountsLiquidatableResponse_Result {
    const message = createBaseAreSubaccountsLiquidatableResponse_Result();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? SubaccountId.fromPartial(object.subaccountId) : undefined;
    message.isLiquidatable = object.isLiquidatable ?? false;
    return message;
  }

};