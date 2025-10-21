import { PerpetualFeeParams, PerpetualFeeParamsSDKType, PerpetualFeeTier, PerpetualFeeTierSDKType } from "./params";
import { PerMarketFeeDiscountParams, PerMarketFeeDiscountParamsSDKType } from "./per_market_fee_discount";
import { StakingTier, StakingTierSDKType } from "./staking_tier";
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
/** QueryStakingTiersRequest is a request type for the StakingTiers RPC method. */

export interface QueryStakingTiersRequest {}
/** QueryStakingTiersRequest is a request type for the StakingTiers RPC method. */

export interface QueryStakingTiersRequestSDKType {}
/** QueryStakingTiersResponse is a response type for the StakingTiers RPC method. */

export interface QueryStakingTiersResponse {
  /** QueryStakingTiersResponse is a response type for the StakingTiers RPC method. */
  stakingTiers: StakingTier[];
}
/** QueryStakingTiersResponse is a response type for the StakingTiers RPC method. */

export interface QueryStakingTiersResponseSDKType {
  /** QueryStakingTiersResponse is a response type for the StakingTiers RPC method. */
  staking_tiers: StakingTierSDKType[];
}
/**
 * QueryUserStakingTierRequest is a request type for the UserStakingTier RPC
 * method.
 */

export interface QueryUserStakingTierRequest {
  address: string;
}
/**
 * QueryUserStakingTierRequest is a request type for the UserStakingTier RPC
 * method.
 */

export interface QueryUserStakingTierRequestSDKType {
  address: string;
}
/**
 * QueryUserStakingTierResponse is a response type for the UserStakingTier RPC
 * method.
 */

export interface QueryUserStakingTierResponse {
  /** The user's current fee tier name */
  feeTierName: string;
  /** Amount of tokens staked by the user (in base units) */

  stakedBaseTokens: Uint8Array;
  /** The discount percentage in ppm that user qualifies for */

  discountPpm: number;
}
/**
 * QueryUserStakingTierResponse is a response type for the UserStakingTier RPC
 * method.
 */

export interface QueryUserStakingTierResponseSDKType {
  /** The user's current fee tier name */
  fee_tier_name: string;
  /** Amount of tokens staked by the user (in base units) */

  staked_base_tokens: Uint8Array;
  /** The discount percentage in ppm that user qualifies for */

  discount_ppm: number;
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

function createBaseQueryStakingTiersRequest(): QueryStakingTiersRequest {
  return {};
}

export const QueryStakingTiersRequest = {
  encode(_: QueryStakingTiersRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryStakingTiersRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryStakingTiersRequest();

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

  fromPartial(_: DeepPartial<QueryStakingTiersRequest>): QueryStakingTiersRequest {
    const message = createBaseQueryStakingTiersRequest();
    return message;
  }

};

function createBaseQueryStakingTiersResponse(): QueryStakingTiersResponse {
  return {
    stakingTiers: []
  };
}

export const QueryStakingTiersResponse = {
  encode(message: QueryStakingTiersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.stakingTiers) {
      StakingTier.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryStakingTiersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryStakingTiersResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.stakingTiers.push(StakingTier.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryStakingTiersResponse>): QueryStakingTiersResponse {
    const message = createBaseQueryStakingTiersResponse();
    message.stakingTiers = object.stakingTiers?.map(e => StakingTier.fromPartial(e)) || [];
    return message;
  }

};

function createBaseQueryUserStakingTierRequest(): QueryUserStakingTierRequest {
  return {
    address: ""
  };
}

export const QueryUserStakingTierRequest = {
  encode(message: QueryUserStakingTierRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryUserStakingTierRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryUserStakingTierRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryUserStakingTierRequest>): QueryUserStakingTierRequest {
    const message = createBaseQueryUserStakingTierRequest();
    message.address = object.address ?? "";
    return message;
  }

};

function createBaseQueryUserStakingTierResponse(): QueryUserStakingTierResponse {
  return {
    feeTierName: "",
    stakedBaseTokens: new Uint8Array(),
    discountPpm: 0
  };
}

export const QueryUserStakingTierResponse = {
  encode(message: QueryUserStakingTierResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.feeTierName !== "") {
      writer.uint32(10).string(message.feeTierName);
    }

    if (message.stakedBaseTokens.length !== 0) {
      writer.uint32(18).bytes(message.stakedBaseTokens);
    }

    if (message.discountPpm !== 0) {
      writer.uint32(24).uint32(message.discountPpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryUserStakingTierResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryUserStakingTierResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.feeTierName = reader.string();
          break;

        case 2:
          message.stakedBaseTokens = reader.bytes();
          break;

        case 3:
          message.discountPpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryUserStakingTierResponse>): QueryUserStakingTierResponse {
    const message = createBaseQueryUserStakingTierResponse();
    message.feeTierName = object.feeTierName ?? "";
    message.stakedBaseTokens = object.stakedBaseTokens ?? new Uint8Array();
    message.discountPpm = object.discountPpm ?? 0;
    return message;
  }

};