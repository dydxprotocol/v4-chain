import { DowntimeParams, DowntimeParamsSDKType } from "./params";
import { BlockInfo, BlockInfoSDKType, AllDowntimeInfo, AllDowntimeInfoSDKType } from "./blocktime";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * QueryDowntimeParamsRequest is a request type for the DowntimeParams
 * RPC method.
 */
export interface QueryDowntimeParamsRequest {
}
/**
 * QueryDowntimeParamsRequest is a request type for the DowntimeParams
 * RPC method.
 */
export interface QueryDowntimeParamsRequestSDKType {
}
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
export interface QueryPreviousBlockInfoRequest {
}
/**
 * QueryPreviousBlockInfoRequest is a request type for the PreviousBlockInfo
 * RPC method.
 */
export interface QueryPreviousBlockInfoRequestSDKType {
}
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
    info?: BlockInfoSDKType;
}
/**
 * QueryAllDowntimeInfoRequest is a request type for the AllDowntimeInfo
 * RPC method.
 */
export interface QueryAllDowntimeInfoRequest {
}
/**
 * QueryAllDowntimeInfoRequest is a request type for the AllDowntimeInfo
 * RPC method.
 */
export interface QueryAllDowntimeInfoRequestSDKType {
}
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
    info?: AllDowntimeInfoSDKType;
}
export declare const QueryDowntimeParamsRequest: {
    encode(_: QueryDowntimeParamsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDowntimeParamsRequest;
    fromPartial(_: DeepPartial<QueryDowntimeParamsRequest>): QueryDowntimeParamsRequest;
};
export declare const QueryDowntimeParamsResponse: {
    encode(message: QueryDowntimeParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDowntimeParamsResponse;
    fromPartial(object: DeepPartial<QueryDowntimeParamsResponse>): QueryDowntimeParamsResponse;
};
export declare const QueryPreviousBlockInfoRequest: {
    encode(_: QueryPreviousBlockInfoRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryPreviousBlockInfoRequest;
    fromPartial(_: DeepPartial<QueryPreviousBlockInfoRequest>): QueryPreviousBlockInfoRequest;
};
export declare const QueryPreviousBlockInfoResponse: {
    encode(message: QueryPreviousBlockInfoResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryPreviousBlockInfoResponse;
    fromPartial(object: DeepPartial<QueryPreviousBlockInfoResponse>): QueryPreviousBlockInfoResponse;
};
export declare const QueryAllDowntimeInfoRequest: {
    encode(_: QueryAllDowntimeInfoRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllDowntimeInfoRequest;
    fromPartial(_: DeepPartial<QueryAllDowntimeInfoRequest>): QueryAllDowntimeInfoRequest;
};
export declare const QueryAllDowntimeInfoResponse: {
    encode(message: QueryAllDowntimeInfoResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllDowntimeInfoResponse;
    fromPartial(object: DeepPartial<QueryAllDowntimeInfoResponse>): QueryAllDowntimeInfoResponse;
};
