import { Params, ParamsSDKType } from "./params";
import { StatsMetadata, StatsMetadataSDKType, GlobalStats, GlobalStatsSDKType, UserStats, UserStatsSDKType } from "./stats";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QueryParamsRequest is a request type for the Params RPC method. */
export interface QueryParamsRequest {
}
/** QueryParamsRequest is a request type for the Params RPC method. */
export interface QueryParamsRequestSDKType {
}
/** QueryParamsResponse is a response type for the Params RPC method. */
export interface QueryParamsResponse {
    params?: Params;
}
/** QueryParamsResponse is a response type for the Params RPC method. */
export interface QueryParamsResponseSDKType {
    params?: ParamsSDKType;
}
/** QueryStatsMetadataRequest is a request type for the StatsMetadata RPC method. */
export interface QueryStatsMetadataRequest {
}
/** QueryStatsMetadataRequest is a request type for the StatsMetadata RPC method. */
export interface QueryStatsMetadataRequestSDKType {
}
/**
 * QueryStatsMetadataResponse is a response type for the StatsMetadata RPC
 * method.
 */
export interface QueryStatsMetadataResponse {
    /**
     * QueryStatsMetadataResponse is a response type for the StatsMetadata RPC
     * method.
     */
    metadata?: StatsMetadata;
}
/**
 * QueryStatsMetadataResponse is a response type for the StatsMetadata RPC
 * method.
 */
export interface QueryStatsMetadataResponseSDKType {
    metadata?: StatsMetadataSDKType;
}
/** QueryGlobalStatsRequest is a request type for the GlobalStats RPC method. */
export interface QueryGlobalStatsRequest {
}
/** QueryGlobalStatsRequest is a request type for the GlobalStats RPC method. */
export interface QueryGlobalStatsRequestSDKType {
}
/** QueryGlobalStatsResponse is a response type for the GlobalStats RPC method. */
export interface QueryGlobalStatsResponse {
    /** QueryGlobalStatsResponse is a response type for the GlobalStats RPC method. */
    stats?: GlobalStats;
}
/** QueryGlobalStatsResponse is a response type for the GlobalStats RPC method. */
export interface QueryGlobalStatsResponseSDKType {
    stats?: GlobalStatsSDKType;
}
/** QueryUserStatsRequest is a request type for the UserStats RPC method. */
export interface QueryUserStatsRequest {
    /** QueryUserStatsRequest is a request type for the UserStats RPC method. */
    user: string;
}
/** QueryUserStatsRequest is a request type for the UserStats RPC method. */
export interface QueryUserStatsRequestSDKType {
    user: string;
}
/** QueryUserStatsResponse is a request type for the UserStats RPC method. */
export interface QueryUserStatsResponse {
    /** QueryUserStatsResponse is a request type for the UserStats RPC method. */
    stats?: UserStats;
}
/** QueryUserStatsResponse is a request type for the UserStats RPC method. */
export interface QueryUserStatsResponseSDKType {
    stats?: UserStatsSDKType;
}
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
export declare const QueryStatsMetadataRequest: {
    encode(_: QueryStatsMetadataRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryStatsMetadataRequest;
    fromPartial(_: DeepPartial<QueryStatsMetadataRequest>): QueryStatsMetadataRequest;
};
export declare const QueryStatsMetadataResponse: {
    encode(message: QueryStatsMetadataResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryStatsMetadataResponse;
    fromPartial(object: DeepPartial<QueryStatsMetadataResponse>): QueryStatsMetadataResponse;
};
export declare const QueryGlobalStatsRequest: {
    encode(_: QueryGlobalStatsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryGlobalStatsRequest;
    fromPartial(_: DeepPartial<QueryGlobalStatsRequest>): QueryGlobalStatsRequest;
};
export declare const QueryGlobalStatsResponse: {
    encode(message: QueryGlobalStatsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryGlobalStatsResponse;
    fromPartial(object: DeepPartial<QueryGlobalStatsResponse>): QueryGlobalStatsResponse;
};
export declare const QueryUserStatsRequest: {
    encode(message: QueryUserStatsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryUserStatsRequest;
    fromPartial(object: DeepPartial<QueryUserStatsRequest>): QueryUserStatsRequest;
};
export declare const QueryUserStatsResponse: {
    encode(message: QueryUserStatsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryUserStatsResponse;
    fromPartial(object: DeepPartial<QueryUserStatsResponse>): QueryUserStatsResponse;
};
