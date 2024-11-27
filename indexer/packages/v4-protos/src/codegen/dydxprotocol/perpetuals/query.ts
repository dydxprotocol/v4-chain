import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { Perpetual, PerpetualSDKType, LiquidityTier, LiquidityTierSDKType, PremiumStore, PremiumStoreSDKType } from "./perpetual";
import { Params, ParamsSDKType } from "./params";
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
/** Queries a list of LiquidityTier items. */

export interface QueryAllLiquidityTiersRequest {
  pagination?: PageRequest;
}
/** Queries a list of LiquidityTier items. */

export interface QueryAllLiquidityTiersRequestSDKType {
  pagination?: PageRequestSDKType;
}
/**
 * QueryAllLiquidityTiersResponse is response type for the AllLiquidityTiers RPC
 * method.
 */

export interface QueryAllLiquidityTiersResponse {
  liquidityTiers: LiquidityTier[];
  pagination?: PageResponse;
}
/**
 * QueryAllLiquidityTiersResponse is response type for the AllLiquidityTiers RPC
 * method.
 */

export interface QueryAllLiquidityTiersResponseSDKType {
  liquidity_tiers: LiquidityTierSDKType[];
  pagination?: PageResponseSDKType;
}
/** QueryPremiumVotesRequest is the request type for the PremiumVotes RPC method. */

export interface QueryPremiumVotesRequest {}
/** QueryPremiumVotesRequest is the request type for the PremiumVotes RPC method. */

export interface QueryPremiumVotesRequestSDKType {}
/**
 * QueryPremiumVotesResponse is the response type for the PremiumVotes RPC
 * method.
 */

export interface QueryPremiumVotesResponse {
  premiumVotes?: PremiumStore;
}
/**
 * QueryPremiumVotesResponse is the response type for the PremiumVotes RPC
 * method.
 */

export interface QueryPremiumVotesResponseSDKType {
  premium_votes?: PremiumStoreSDKType;
}
/**
 * QueryPremiumSamplesRequest is the request type for the PremiumSamples RPC
 * method.
 */

export interface QueryPremiumSamplesRequest {}
/**
 * QueryPremiumSamplesRequest is the request type for the PremiumSamples RPC
 * method.
 */

export interface QueryPremiumSamplesRequestSDKType {}
/**
 * QueryPremiumSamplesResponse is the response type for the PremiumSamples RPC
 * method.
 */

export interface QueryPremiumSamplesResponse {
  premiumSamples?: PremiumStore;
}
/**
 * QueryPremiumSamplesResponse is the response type for the PremiumSamples RPC
 * method.
 */

export interface QueryPremiumSamplesResponseSDKType {
  premium_samples?: PremiumStoreSDKType;
}
/** QueryParamsResponse is the response type for the Params RPC method. */

export interface QueryParamsRequest {}
/** QueryParamsResponse is the response type for the Params RPC method. */

export interface QueryParamsRequestSDKType {}
/** QueryParamsResponse is the response type for the Params RPC method. */

export interface QueryParamsResponse {
  params?: Params;
}
/** QueryParamsResponse is the response type for the Params RPC method. */

export interface QueryParamsResponseSDKType {
  params?: ParamsSDKType;
}
/** QueryNextPerpetualIdRequest is the request type for the NextPerpetualId RPC */

export interface QueryNextPerpetualIdRequest {}
/** QueryNextPerpetualIdRequest is the request type for the NextPerpetualId RPC */

export interface QueryNextPerpetualIdRequestSDKType {}
/** QueryNextPerpetualIdResponse is the response type for the NextPerpetualId RPC */

export interface QueryNextPerpetualIdResponse {
  /** QueryNextPerpetualIdResponse is the response type for the NextPerpetualId RPC */
  nextPerpetualId: number;
}
/** QueryNextPerpetualIdResponse is the response type for the NextPerpetualId RPC */

export interface QueryNextPerpetualIdResponseSDKType {
  /** QueryNextPerpetualIdResponse is the response type for the NextPerpetualId RPC */
  next_perpetual_id: number;
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

function createBaseQueryAllLiquidityTiersRequest(): QueryAllLiquidityTiersRequest {
  return {
    pagination: undefined
  };
}

export const QueryAllLiquidityTiersRequest = {
  encode(message: QueryAllLiquidityTiersRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllLiquidityTiersRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllLiquidityTiersRequest();

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

  fromPartial(object: DeepPartial<QueryAllLiquidityTiersRequest>): QueryAllLiquidityTiersRequest {
    const message = createBaseQueryAllLiquidityTiersRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryAllLiquidityTiersResponse(): QueryAllLiquidityTiersResponse {
  return {
    liquidityTiers: [],
    pagination: undefined
  };
}

export const QueryAllLiquidityTiersResponse = {
  encode(message: QueryAllLiquidityTiersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.liquidityTiers) {
      LiquidityTier.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllLiquidityTiersResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllLiquidityTiersResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.liquidityTiers.push(LiquidityTier.decode(reader, reader.uint32()));
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

  fromPartial(object: DeepPartial<QueryAllLiquidityTiersResponse>): QueryAllLiquidityTiersResponse {
    const message = createBaseQueryAllLiquidityTiersResponse();
    message.liquidityTiers = object.liquidityTiers?.map(e => LiquidityTier.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryPremiumVotesRequest(): QueryPremiumVotesRequest {
  return {};
}

export const QueryPremiumVotesRequest = {
  encode(_: QueryPremiumVotesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPremiumVotesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPremiumVotesRequest();

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

  fromPartial(_: DeepPartial<QueryPremiumVotesRequest>): QueryPremiumVotesRequest {
    const message = createBaseQueryPremiumVotesRequest();
    return message;
  }

};

function createBaseQueryPremiumVotesResponse(): QueryPremiumVotesResponse {
  return {
    premiumVotes: undefined
  };
}

export const QueryPremiumVotesResponse = {
  encode(message: QueryPremiumVotesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.premiumVotes !== undefined) {
      PremiumStore.encode(message.premiumVotes, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPremiumVotesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPremiumVotesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.premiumVotes = PremiumStore.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryPremiumVotesResponse>): QueryPremiumVotesResponse {
    const message = createBaseQueryPremiumVotesResponse();
    message.premiumVotes = object.premiumVotes !== undefined && object.premiumVotes !== null ? PremiumStore.fromPartial(object.premiumVotes) : undefined;
    return message;
  }

};

function createBaseQueryPremiumSamplesRequest(): QueryPremiumSamplesRequest {
  return {};
}

export const QueryPremiumSamplesRequest = {
  encode(_: QueryPremiumSamplesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPremiumSamplesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPremiumSamplesRequest();

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

  fromPartial(_: DeepPartial<QueryPremiumSamplesRequest>): QueryPremiumSamplesRequest {
    const message = createBaseQueryPremiumSamplesRequest();
    return message;
  }

};

function createBaseQueryPremiumSamplesResponse(): QueryPremiumSamplesResponse {
  return {
    premiumSamples: undefined
  };
}

export const QueryPremiumSamplesResponse = {
  encode(message: QueryPremiumSamplesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.premiumSamples !== undefined) {
      PremiumStore.encode(message.premiumSamples, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPremiumSamplesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPremiumSamplesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.premiumSamples = PremiumStore.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryPremiumSamplesResponse>): QueryPremiumSamplesResponse {
    const message = createBaseQueryPremiumSamplesResponse();
    message.premiumSamples = object.premiumSamples !== undefined && object.premiumSamples !== null ? PremiumStore.fromPartial(object.premiumSamples) : undefined;
    return message;
  }

};

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();

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

  fromPartial(_: DeepPartial<QueryParamsRequest>): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  }

};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return {
    params: undefined
  };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryParamsResponse>): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryNextPerpetualIdRequest(): QueryNextPerpetualIdRequest {
  return {};
}

export const QueryNextPerpetualIdRequest = {
  encode(_: QueryNextPerpetualIdRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNextPerpetualIdRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNextPerpetualIdRequest();

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

  fromPartial(_: DeepPartial<QueryNextPerpetualIdRequest>): QueryNextPerpetualIdRequest {
    const message = createBaseQueryNextPerpetualIdRequest();
    return message;
  }

};

function createBaseQueryNextPerpetualIdResponse(): QueryNextPerpetualIdResponse {
  return {
    nextPerpetualId: 0
  };
}

export const QueryNextPerpetualIdResponse = {
  encode(message: QueryNextPerpetualIdResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nextPerpetualId !== 0) {
      writer.uint32(8).uint32(message.nextPerpetualId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNextPerpetualIdResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNextPerpetualIdResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.nextPerpetualId = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryNextPerpetualIdResponse>): QueryNextPerpetualIdResponse {
    const message = createBaseQueryNextPerpetualIdResponse();
    message.nextPerpetualId = object.nextPerpetualId ?? 0;
    return message;
  }

};