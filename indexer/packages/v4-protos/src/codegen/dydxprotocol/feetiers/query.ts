import { PerpetualFeeParams, PerpetualFeeParamsSDKType, PerpetualFeeTier, PerpetualFeeTierSDKType } from "./params";
import { FeeHolidayParams, FeeHolidayParamsSDKType } from "./fee_holiday";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * QueryPerpetualFeeParamsRequest is a request type for the PerpetualFeeParams
 * RPC method.
 */

export interface QueryPerpetualFeeParamsRequest {}
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
  params?: PerpetualFeeParams;
}
/**
 * QueryPerpetualFeeParamsResponse is a response type for the PerpetualFeeParams
 * RPC method.
 */

export interface QueryPerpetualFeeParamsResponseSDKType {
  params?: PerpetualFeeParamsSDKType;
}
/** QueryUserFeeTierRequest is a request type for the UserFeeTier RPC method. */

export interface QueryUserFeeTierRequest {
  user: string;
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
/** QueryUserFeeTierResponse is a request type for the UserFeeTier RPC method. */

export interface QueryUserFeeTierResponseSDKType {
  /** Index of the fee tier in the list queried from PerpetualFeeParams. */
  index: number;
  tier?: PerpetualFeeTierSDKType;
}
/** QueryFeeHolidayParamsRequest is a request type for the FeeHolidayParams RPC method. */

export interface QueryFeeHolidayParamsRequest {
  clobPairId: number;
}
/** QueryFeeHolidayParamsRequest is a request type for the FeeHolidayParams RPC method. */

export interface QueryFeeHolidayParamsRequestSDKType {
  clob_pair_id: number;
}
/** QueryFeeHolidayParamsResponse is a response type for the FeeHolidayParams RPC method. */

export interface QueryFeeHolidayParamsResponse {
  params?: FeeHolidayParams;
}
/** QueryFeeHolidayParamsResponse is a response type for the FeeHolidayParams RPC method. */

export interface QueryFeeHolidayParamsResponseSDKType {
  params?: FeeHolidayParamsSDKType;
}
/** QueryAllFeeHolidayParamsRequest is a request type for the AllFeeHolidayParams RPC method. */

export interface QueryAllFeeHolidayParamsRequest {}
/** QueryAllFeeHolidayParamsRequest is a request type for the AllFeeHolidayParams RPC method. */

export interface QueryAllFeeHolidayParamsRequestSDKType {}
/** QueryAllFeeHolidayParamsResponse is a response type for the AllFeeHolidayParams RPC method. */

export interface QueryAllFeeHolidayParamsResponse {
  params: FeeHolidayParams[];
}
/** QueryAllFeeHolidayParamsResponse is a response type for the AllFeeHolidayParams RPC method. */

export interface QueryAllFeeHolidayParamsResponseSDKType {
  params: FeeHolidayParamsSDKType[];
}

function createBaseQueryPerpetualFeeParamsRequest(): QueryPerpetualFeeParamsRequest {
  return {};
}

export const QueryPerpetualFeeParamsRequest = {
  encode(_: QueryPerpetualFeeParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPerpetualFeeParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(_: DeepPartial<QueryPerpetualFeeParamsRequest>): QueryPerpetualFeeParamsRequest {
    const message = createBaseQueryPerpetualFeeParamsRequest();
    return message;
  }

};

function createBaseQueryPerpetualFeeParamsResponse(): QueryPerpetualFeeParamsResponse {
  return {
    params: undefined
  };
}

export const QueryPerpetualFeeParamsResponse = {
  encode(message: QueryPerpetualFeeParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      PerpetualFeeParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPerpetualFeeParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<QueryPerpetualFeeParamsResponse>): QueryPerpetualFeeParamsResponse {
    const message = createBaseQueryPerpetualFeeParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? PerpetualFeeParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryUserFeeTierRequest(): QueryUserFeeTierRequest {
  return {
    user: ""
  };
}

export const QueryUserFeeTierRequest = {
  encode(message: QueryUserFeeTierRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== "") {
      writer.uint32(10).string(message.user);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryUserFeeTierRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<QueryUserFeeTierRequest>): QueryUserFeeTierRequest {
    const message = createBaseQueryUserFeeTierRequest();
    message.user = object.user ?? "";
    return message;
  }

};

function createBaseQueryUserFeeTierResponse(): QueryUserFeeTierResponse {
  return {
    index: 0,
    tier: undefined
  };
}

export const QueryUserFeeTierResponse = {
  encode(message: QueryUserFeeTierResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.index !== 0) {
      writer.uint32(8).uint32(message.index);
    }

    if (message.tier !== undefined) {
      PerpetualFeeTier.encode(message.tier, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryUserFeeTierResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<QueryUserFeeTierResponse>): QueryUserFeeTierResponse {
    const message = createBaseQueryUserFeeTierResponse();
    message.index = object.index ?? 0;
    message.tier = object.tier !== undefined && object.tier !== null ? PerpetualFeeTier.fromPartial(object.tier) : undefined;
    return message;
  }

};

function createBaseQueryFeeHolidayParamsRequest(): QueryFeeHolidayParamsRequest {
  return {
    clobPairId: 0
  };
}

export const QueryFeeHolidayParamsRequest = {
  encode(message: QueryFeeHolidayParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobPairId !== 0) {
      writer.uint32(8).uint32(message.clobPairId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFeeHolidayParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFeeHolidayParamsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPairId = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryFeeHolidayParamsRequest>): QueryFeeHolidayParamsRequest {
    const message = createBaseQueryFeeHolidayParamsRequest();
    message.clobPairId = object.clobPairId ?? 0;
    return message;
  }

};

function createBaseQueryFeeHolidayParamsResponse(): QueryFeeHolidayParamsResponse {
  return {
    params: undefined
  };
}

export const QueryFeeHolidayParamsResponse = {
  encode(message: QueryFeeHolidayParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      FeeHolidayParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFeeHolidayParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFeeHolidayParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = FeeHolidayParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryFeeHolidayParamsResponse>): QueryFeeHolidayParamsResponse {
    const message = createBaseQueryFeeHolidayParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? FeeHolidayParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryAllFeeHolidayParamsRequest(): QueryAllFeeHolidayParamsRequest {
  return {};
}

export const QueryAllFeeHolidayParamsRequest = {
  encode(_: QueryAllFeeHolidayParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllFeeHolidayParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllFeeHolidayParamsRequest();

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

  fromPartial(_: DeepPartial<QueryAllFeeHolidayParamsRequest>): QueryAllFeeHolidayParamsRequest {
    const message = createBaseQueryAllFeeHolidayParamsRequest();
    return message;
  }

};

function createBaseQueryAllFeeHolidayParamsResponse(): QueryAllFeeHolidayParamsResponse {
  return {
    params: []
  };
}

export const QueryAllFeeHolidayParamsResponse = {
  encode(message: QueryAllFeeHolidayParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.params) {
      FeeHolidayParams.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllFeeHolidayParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllFeeHolidayParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params.push(FeeHolidayParams.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllFeeHolidayParamsResponse>): QueryAllFeeHolidayParamsResponse {
    const message = createBaseQueryAllFeeHolidayParamsResponse();
    message.params = object.params?.map(e => FeeHolidayParams.fromPartial(e)) || [];
    return message;
  }

};