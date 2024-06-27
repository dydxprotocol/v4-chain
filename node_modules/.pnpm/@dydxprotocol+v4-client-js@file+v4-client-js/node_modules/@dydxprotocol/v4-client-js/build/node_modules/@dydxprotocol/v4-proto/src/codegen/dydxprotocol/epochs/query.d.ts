import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { EpochInfo, EpochInfoSDKType } from "./epoch_info";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QueryGetEpochInfoRequest is request type for the GetEpochInfo RPC method. */
export interface QueryGetEpochInfoRequest {
    /** QueryGetEpochInfoRequest is request type for the GetEpochInfo RPC method. */
    name: string;
}
/** QueryGetEpochInfoRequest is request type for the GetEpochInfo RPC method. */
export interface QueryGetEpochInfoRequestSDKType {
    name: string;
}
/** QueryEpochInfoResponse is response type for the GetEpochInfo RPC method. */
export interface QueryEpochInfoResponse {
    epochInfo?: EpochInfo;
}
/** QueryEpochInfoResponse is response type for the GetEpochInfo RPC method. */
export interface QueryEpochInfoResponseSDKType {
    epoch_info?: EpochInfoSDKType;
}
/** QueryAllEpochInfoRequest is request type for the AllEpochInfo RPC method. */
export interface QueryAllEpochInfoRequest {
    pagination?: PageRequest;
}
/** QueryAllEpochInfoRequest is request type for the AllEpochInfo RPC method. */
export interface QueryAllEpochInfoRequestSDKType {
    pagination?: PageRequestSDKType;
}
/** QueryEpochInfoAllResponse is response type for the AllEpochInfo RPC method. */
export interface QueryEpochInfoAllResponse {
    epochInfo: EpochInfo[];
    pagination?: PageResponse;
}
/** QueryEpochInfoAllResponse is response type for the AllEpochInfo RPC method. */
export interface QueryEpochInfoAllResponseSDKType {
    epoch_info: EpochInfoSDKType[];
    pagination?: PageResponseSDKType;
}
export declare const QueryGetEpochInfoRequest: {
    encode(message: QueryGetEpochInfoRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryGetEpochInfoRequest;
    fromPartial(object: DeepPartial<QueryGetEpochInfoRequest>): QueryGetEpochInfoRequest;
};
export declare const QueryEpochInfoResponse: {
    encode(message: QueryEpochInfoResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryEpochInfoResponse;
    fromPartial(object: DeepPartial<QueryEpochInfoResponse>): QueryEpochInfoResponse;
};
export declare const QueryAllEpochInfoRequest: {
    encode(message: QueryAllEpochInfoRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllEpochInfoRequest;
    fromPartial(object: DeepPartial<QueryAllEpochInfoRequest>): QueryAllEpochInfoRequest;
};
export declare const QueryEpochInfoAllResponse: {
    encode(message: QueryEpochInfoAllResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryEpochInfoAllResponse;
    fromPartial(object: DeepPartial<QueryEpochInfoAllResponse>): QueryEpochInfoAllResponse;
};
