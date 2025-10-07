import { PerpetualFeeParams, PerpetualFeeParamsSDKType, PerpetualFeeTier, PerpetualFeeTierSDKType } from "./params";
import { FeeDiscountCampaignParams, FeeDiscountCampaignParamsSDKType } from "./fee_discount_campaign";
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
 * QueryFeeDiscountCampaignParamsRequest is the request type for the
 * Query/FeeDiscountCampaignParams RPC method.
 */

export interface QueryFeeDiscountCampaignParamsRequest {
  /**
   * QueryFeeDiscountCampaignParamsRequest is the request type for the
   * Query/FeeDiscountCampaignParams RPC method.
   */
  clobPairId: number;
}
/**
 * QueryFeeDiscountCampaignParamsRequest is the request type for the
 * Query/FeeDiscountCampaignParams RPC method.
 */

export interface QueryFeeDiscountCampaignParamsRequestSDKType {
  /**
   * QueryFeeDiscountCampaignParamsRequest is the request type for the
   * Query/FeeDiscountCampaignParams RPC method.
   */
  clob_pair_id: number;
}
/**
 * QueryFeeDiscountCampaignParamsResponse is the response type for the
 * Query/FeeDiscountCampaignParams RPC method.
 */

export interface QueryFeeDiscountCampaignParamsResponse {
  params?: FeeDiscountCampaignParams;
}
/**
 * QueryFeeDiscountCampaignParamsResponse is the response type for the
 * Query/FeeDiscountCampaignParams RPC method.
 */

export interface QueryFeeDiscountCampaignParamsResponseSDKType {
  params?: FeeDiscountCampaignParamsSDKType;
}
/**
 * QueryAllFeeDiscountCampaignParamsRequest is the request type for the
 * Query/AllFeeDiscountCampaignParams RPC method.
 */

export interface QueryAllFeeDiscountCampaignParamsRequest {}
/**
 * QueryAllFeeDiscountCampaignParamsRequest is the request type for the
 * Query/AllFeeDiscountCampaignParams RPC method.
 */

export interface QueryAllFeeDiscountCampaignParamsRequestSDKType {}
/**
 * QueryAllFeeDiscountCampaignParamsResponse is the response type for the
 * Query/AllFeeDiscountCampaignParams RPC method.
 */

export interface QueryAllFeeDiscountCampaignParamsResponse {
  params: FeeDiscountCampaignParams[];
}
/**
 * QueryAllFeeDiscountCampaignParamsResponse is the response type for the
 * Query/AllFeeDiscountCampaignParams RPC method.
 */

export interface QueryAllFeeDiscountCampaignParamsResponseSDKType {
  params: FeeDiscountCampaignParamsSDKType[];
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

function createBaseQueryFeeDiscountCampaignParamsRequest(): QueryFeeDiscountCampaignParamsRequest {
  return {
    clobPairId: 0
  };
}

export const QueryFeeDiscountCampaignParamsRequest = {
  encode(message: QueryFeeDiscountCampaignParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobPairId !== 0) {
      writer.uint32(8).uint32(message.clobPairId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFeeDiscountCampaignParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFeeDiscountCampaignParamsRequest();

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

  fromPartial(object: DeepPartial<QueryFeeDiscountCampaignParamsRequest>): QueryFeeDiscountCampaignParamsRequest {
    const message = createBaseQueryFeeDiscountCampaignParamsRequest();
    message.clobPairId = object.clobPairId ?? 0;
    return message;
  }

};

function createBaseQueryFeeDiscountCampaignParamsResponse(): QueryFeeDiscountCampaignParamsResponse {
  return {
    params: undefined
  };
}

export const QueryFeeDiscountCampaignParamsResponse = {
  encode(message: QueryFeeDiscountCampaignParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      FeeDiscountCampaignParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFeeDiscountCampaignParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFeeDiscountCampaignParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = FeeDiscountCampaignParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryFeeDiscountCampaignParamsResponse>): QueryFeeDiscountCampaignParamsResponse {
    const message = createBaseQueryFeeDiscountCampaignParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? FeeDiscountCampaignParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryAllFeeDiscountCampaignParamsRequest(): QueryAllFeeDiscountCampaignParamsRequest {
  return {};
}

export const QueryAllFeeDiscountCampaignParamsRequest = {
  encode(_: QueryAllFeeDiscountCampaignParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllFeeDiscountCampaignParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllFeeDiscountCampaignParamsRequest();

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

  fromPartial(_: DeepPartial<QueryAllFeeDiscountCampaignParamsRequest>): QueryAllFeeDiscountCampaignParamsRequest {
    const message = createBaseQueryAllFeeDiscountCampaignParamsRequest();
    return message;
  }

};

function createBaseQueryAllFeeDiscountCampaignParamsResponse(): QueryAllFeeDiscountCampaignParamsResponse {
  return {
    params: []
  };
}

export const QueryAllFeeDiscountCampaignParamsResponse = {
  encode(message: QueryAllFeeDiscountCampaignParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.params) {
      FeeDiscountCampaignParams.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllFeeDiscountCampaignParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllFeeDiscountCampaignParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params.push(FeeDiscountCampaignParams.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllFeeDiscountCampaignParamsResponse>): QueryAllFeeDiscountCampaignParamsResponse {
    const message = createBaseQueryAllFeeDiscountCampaignParamsResponse();
    message.params = object.params?.map(e => FeeDiscountCampaignParams.fromPartial(e)) || [];
    return message;
  }

};