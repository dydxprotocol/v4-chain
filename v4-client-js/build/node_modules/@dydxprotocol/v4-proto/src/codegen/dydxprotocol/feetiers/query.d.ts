import { PerpetualFeeParams, PerpetualFeeParamsSDKType, PerpetualFeeTier, PerpetualFeeTierSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * QueryPerpetualFeeParamsRequest is a request type for the PerpetualFeeParams
 * RPC method.
 */
export interface QueryPerpetualFeeParamsRequest {
}
/**
 * QueryPerpetualFeeParamsRequest is a request type for the PerpetualFeeParams
 * RPC method.
 */
export interface QueryPerpetualFeeParamsRequestSDKType {
}
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
    index: number;
    tier?: PerpetualFeeTierSDKType;
}
export declare const QueryPerpetualFeeParamsRequest: {
    encode(_: QueryPerpetualFeeParamsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryPerpetualFeeParamsRequest;
    fromPartial(_: DeepPartial<QueryPerpetualFeeParamsRequest>): QueryPerpetualFeeParamsRequest;
};
export declare const QueryPerpetualFeeParamsResponse: {
    encode(message: QueryPerpetualFeeParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryPerpetualFeeParamsResponse;
    fromPartial(object: DeepPartial<QueryPerpetualFeeParamsResponse>): QueryPerpetualFeeParamsResponse;
};
export declare const QueryUserFeeTierRequest: {
    encode(message: QueryUserFeeTierRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryUserFeeTierRequest;
    fromPartial(object: DeepPartial<QueryUserFeeTierRequest>): QueryUserFeeTierRequest;
};
export declare const QueryUserFeeTierResponse: {
    encode(message: QueryUserFeeTierResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryUserFeeTierResponse;
    fromPartial(object: DeepPartial<QueryUserFeeTierResponse>): QueryUserFeeTierResponse;
};
