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
export interface QueryPremiumVotesRequest {
}
/** QueryPremiumVotesRequest is the request type for the PremiumVotes RPC method. */
export interface QueryPremiumVotesRequestSDKType {
}
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
export interface QueryPremiumSamplesRequest {
}
/**
 * QueryPremiumSamplesRequest is the request type for the PremiumSamples RPC
 * method.
 */
export interface QueryPremiumSamplesRequestSDKType {
}
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
export interface QueryParamsRequest {
}
/** QueryParamsResponse is the response type for the Params RPC method. */
export interface QueryParamsRequestSDKType {
}
/** QueryParamsResponse is the response type for the Params RPC method. */
export interface QueryParamsResponse {
    params?: Params;
}
/** QueryParamsResponse is the response type for the Params RPC method. */
export interface QueryParamsResponseSDKType {
    params?: ParamsSDKType;
}
export declare const QueryPerpetualRequest: {
    encode(message: QueryPerpetualRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryPerpetualRequest;
    fromPartial(object: DeepPartial<QueryPerpetualRequest>): QueryPerpetualRequest;
};
export declare const QueryPerpetualResponse: {
    encode(message: QueryPerpetualResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryPerpetualResponse;
    fromPartial(object: DeepPartial<QueryPerpetualResponse>): QueryPerpetualResponse;
};
export declare const QueryAllPerpetualsRequest: {
    encode(message: QueryAllPerpetualsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllPerpetualsRequest;
    fromPartial(object: DeepPartial<QueryAllPerpetualsRequest>): QueryAllPerpetualsRequest;
};
export declare const QueryAllPerpetualsResponse: {
    encode(message: QueryAllPerpetualsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllPerpetualsResponse;
    fromPartial(object: DeepPartial<QueryAllPerpetualsResponse>): QueryAllPerpetualsResponse;
};
export declare const QueryAllLiquidityTiersRequest: {
    encode(message: QueryAllLiquidityTiersRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllLiquidityTiersRequest;
    fromPartial(object: DeepPartial<QueryAllLiquidityTiersRequest>): QueryAllLiquidityTiersRequest;
};
export declare const QueryAllLiquidityTiersResponse: {
    encode(message: QueryAllLiquidityTiersResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllLiquidityTiersResponse;
    fromPartial(object: DeepPartial<QueryAllLiquidityTiersResponse>): QueryAllLiquidityTiersResponse;
};
export declare const QueryPremiumVotesRequest: {
    encode(_: QueryPremiumVotesRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryPremiumVotesRequest;
    fromPartial(_: DeepPartial<QueryPremiumVotesRequest>): QueryPremiumVotesRequest;
};
export declare const QueryPremiumVotesResponse: {
    encode(message: QueryPremiumVotesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryPremiumVotesResponse;
    fromPartial(object: DeepPartial<QueryPremiumVotesResponse>): QueryPremiumVotesResponse;
};
export declare const QueryPremiumSamplesRequest: {
    encode(_: QueryPremiumSamplesRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryPremiumSamplesRequest;
    fromPartial(_: DeepPartial<QueryPremiumSamplesRequest>): QueryPremiumSamplesRequest;
};
export declare const QueryPremiumSamplesResponse: {
    encode(message: QueryPremiumSamplesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryPremiumSamplesResponse;
    fromPartial(object: DeepPartial<QueryPremiumSamplesResponse>): QueryPremiumSamplesResponse;
};
export declare const QueryParamsRequest: {
    encode(_: QueryParamsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest;
    fromPartial(_: DeepPartial<QueryParamsRequest>): QueryParamsRequest;
};
export declare const QueryParamsResponse: {
    encode(message: QueryParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse;
    fromPartial(object: DeepPartial<QueryParamsResponse>): QueryParamsResponse;
};
