/// <reference types="long" />
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../../../helpers";
/** ConfigRequest defines the request structure for the Config gRPC query. */
export interface ConfigRequest {
}
/** ConfigRequest defines the request structure for the Config gRPC query. */
export interface ConfigRequestSDKType {
}
/** ConfigResponse defines the response structure for the Config gRPC query. */
export interface ConfigResponse {
    minimumGasPrice: string;
    /** pruning settings */
    pruningKeepRecent: string;
    pruningInterval: string;
}
/** ConfigResponse defines the response structure for the Config gRPC query. */
export interface ConfigResponseSDKType {
    minimum_gas_price: string;
    pruning_keep_recent: string;
    pruning_interval: string;
}
/** StateRequest defines the request structure for the status of a node. */
export interface StatusRequest {
}
/** StateRequest defines the request structure for the status of a node. */
export interface StatusRequestSDKType {
}
/** StateResponse defines the response structure for the status of a node. */
export interface StatusResponse {
    /** earliest block height available in the store */
    earliestStoreHeight: Long;
    /** current block height */
    height: Long;
    /** block height timestamp */
    timestamp?: Date;
    /** app hash of the current block */
    appHash: Uint8Array;
    /** validator hash provided by the consensus header */
    validatorHash: Uint8Array;
}
/** StateResponse defines the response structure for the status of a node. */
export interface StatusResponseSDKType {
    earliest_store_height: Long;
    height: Long;
    timestamp?: Date;
    app_hash: Uint8Array;
    validator_hash: Uint8Array;
}
export declare const ConfigRequest: {
    encode(_: ConfigRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ConfigRequest;
    fromPartial(_: DeepPartial<ConfigRequest>): ConfigRequest;
};
export declare const ConfigResponse: {
    encode(message: ConfigResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ConfigResponse;
    fromPartial(object: DeepPartial<ConfigResponse>): ConfigResponse;
};
export declare const StatusRequest: {
    encode(_: StatusRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): StatusRequest;
    fromPartial(_: DeepPartial<StatusRequest>): StatusRequest;
};
export declare const StatusResponse: {
    encode(message: StatusResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): StatusResponse;
    fromPartial(object: DeepPartial<StatusResponse>): StatusResponse;
};
