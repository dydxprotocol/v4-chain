import { PerpetualFeeParams, PerpetualFeeParamsSDKType, PerpetualFeeTier, PerpetualFeeTierSDKType } from "./params";
import { PerMarketFeeDiscountParams, PerMarketFeeDiscountParamsSDKType } from "./per_market_fee_discount";
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
/**
 * QueryPerMarketFeeDiscountParamsRequest is the request type for the
 * Query/PerMarketFeeDiscountParams RPC method.
 */

export interface QueryPerMarketFeeDiscountParamsRequest {
  /**
   * QueryPerMarketFeeDiscountParamsRequest is the request type for the
   * Query/PerMarketFeeDiscountParams RPC method.
   */
  clobPairId: number;
}
/**
 * QueryPerMarketFeeDiscountParamsRequest is the request type for the
 * Query/PerMarketFeeDiscountParams RPC method.
 */

export interface QueryPerMarketFeeDiscountParamsRequestSDKType {
  /**
   * QueryPerMarketFeeDiscountParamsRequest is the request type for the
   * Query/PerMarketFeeDiscountParams RPC method.
   */
  clob_pair_id: number;
}
/**
 * QueryPerMarketFeeDiscountParamsResponse is the response type for the
 * Query/PerMarketFeeDiscountParams RPC method.
 */

export interface QueryPerMarketFeeDiscountParamsResponse {
  params?: PerMarketFeeDiscountParams;
}
/**
 * QueryPerMarketFeeDiscountParamsResponse is the response type for the
 * Query/PerMarketFeeDiscountParams RPC method.
 */

export interface QueryPerMarketFeeDiscountParamsResponseSDKType {
  params?: PerMarketFeeDiscountParamsSDKType;
}
/**
 * QueryAllMarketFeeDiscountParamsRequest is the request type for the
 * Query/AllMarketFeeDiscountParams RPC method.
 */

export interface QueryAllMarketFeeDiscountParamsRequest {}
/**
 * QueryAllMarketFeeDiscountParamsRequest is the request type for the
 * Query/AllMarketFeeDiscountParams RPC method.
 */

export interface QueryAllMarketFeeDiscountParamsRequestSDKType {}
/**
 * QueryAllMarketFeeDiscountParamsResponse is the response type for the
 * Query/AllMarketFeeDiscountParams RPC method.
 */

export interface QueryAllMarketFeeDiscountParamsResponse {
  params: PerMarketFeeDiscountParams[];
}
/**
 * QueryAllMarketFeeDiscountParamsResponse is the response type for the
 * Query/AllMarketFeeDiscountParams RPC method.
 */

export interface QueryAllMarketFeeDiscountParamsResponseSDKType {
  params: PerMarketFeeDiscountParamsSDKType[];
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

function createBaseQueryPerMarketFeeDiscountParamsRequest(): QueryPerMarketFeeDiscountParamsRequest {
  return {
    clobPairId: 0
  };
}

export const QueryPerMarketFeeDiscountParamsRequest = {
  encode(message: QueryPerMarketFeeDiscountParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobPairId !== 0) {
      writer.uint32(8).uint32(message.clobPairId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPerMarketFeeDiscountParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPerMarketFeeDiscountParamsRequest();

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

  fromPartial(object: DeepPartial<QueryPerMarketFeeDiscountParamsRequest>): QueryPerMarketFeeDiscountParamsRequest {
    const message = createBaseQueryPerMarketFeeDiscountParamsRequest();
    message.clobPairId = object.clobPairId ?? 0;
    return message;
  }

};

function createBaseQueryPerMarketFeeDiscountParamsResponse(): QueryPerMarketFeeDiscountParamsResponse {
  return {
    params: undefined
  };
}

export const QueryPerMarketFeeDiscountParamsResponse = {
  encode(message: QueryPerMarketFeeDiscountParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      PerMarketFeeDiscountParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPerMarketFeeDiscountParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPerMarketFeeDiscountParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = PerMarketFeeDiscountParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryPerMarketFeeDiscountParamsResponse>): QueryPerMarketFeeDiscountParamsResponse {
    const message = createBaseQueryPerMarketFeeDiscountParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? PerMarketFeeDiscountParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryAllMarketFeeDiscountParamsRequest(): QueryAllMarketFeeDiscountParamsRequest {
  return {};
}

export const QueryAllMarketFeeDiscountParamsRequest = {
  encode(_: QueryAllMarketFeeDiscountParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMarketFeeDiscountParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMarketFeeDiscountParamsRequest();

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

  fromPartial(_: DeepPartial<QueryAllMarketFeeDiscountParamsRequest>): QueryAllMarketFeeDiscountParamsRequest {
    const message = createBaseQueryAllMarketFeeDiscountParamsRequest();
    return message;
  }

};

function createBaseQueryAllMarketFeeDiscountParamsResponse(): QueryAllMarketFeeDiscountParamsResponse {
  return {
    params: []
  };
}

export const QueryAllMarketFeeDiscountParamsResponse = {
  encode(message: QueryAllMarketFeeDiscountParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.params) {
      PerMarketFeeDiscountParams.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMarketFeeDiscountParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMarketFeeDiscountParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params.push(PerMarketFeeDiscountParams.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllMarketFeeDiscountParamsResponse>): QueryAllMarketFeeDiscountParamsResponse {
    const message = createBaseQueryAllMarketFeeDiscountParamsResponse();
    message.params = object.params?.map(e => PerMarketFeeDiscountParams.fromPartial(e)) || [];
    return message;
  }

};