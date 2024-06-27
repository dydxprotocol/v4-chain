/// <reference types="long" />
import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../query/v1beta1/pagination";
import { Any, AnySDKType } from "../../../../google/protobuf/any";
import { BlockID, BlockIDSDKType } from "../../../../tendermint/types/types";
import { Block as Block1 } from "../../../../tendermint/types/block";
import { BlockSDKType as Block1SDKType } from "../../../../tendermint/types/block";
import { Block as Block2 } from "./types";
import { BlockSDKType as Block2SDKType } from "./types";
import { DefaultNodeInfo, DefaultNodeInfoSDKType } from "../../../../tendermint/p2p/types";
import { Long, DeepPartial } from "../../../../helpers";
import * as _m0 from "protobufjs/minimal";
/** GetValidatorSetByHeightRequest is the request type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetValidatorSetByHeightRequest {
    height: Long;
    /** pagination defines an pagination for the request. */
    pagination?: PageRequest;
}
/** GetValidatorSetByHeightRequest is the request type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetValidatorSetByHeightRequestSDKType {
    height: Long;
    pagination?: PageRequestSDKType;
}
/** GetValidatorSetByHeightResponse is the response type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetValidatorSetByHeightResponse {
    blockHeight: Long;
    validators: Validator[];
    /** pagination defines an pagination for the response. */
    pagination?: PageResponse;
}
/** GetValidatorSetByHeightResponse is the response type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetValidatorSetByHeightResponseSDKType {
    block_height: Long;
    validators: ValidatorSDKType[];
    pagination?: PageResponseSDKType;
}
/** GetLatestValidatorSetRequest is the request type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetLatestValidatorSetRequest {
    /** pagination defines an pagination for the request. */
    pagination?: PageRequest;
}
/** GetLatestValidatorSetRequest is the request type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetLatestValidatorSetRequestSDKType {
    pagination?: PageRequestSDKType;
}
/** GetLatestValidatorSetResponse is the response type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetLatestValidatorSetResponse {
    blockHeight: Long;
    validators: Validator[];
    /** pagination defines an pagination for the response. */
    pagination?: PageResponse;
}
/** GetLatestValidatorSetResponse is the response type for the Query/GetValidatorSetByHeight RPC method. */
export interface GetLatestValidatorSetResponseSDKType {
    block_height: Long;
    validators: ValidatorSDKType[];
    pagination?: PageResponseSDKType;
}
/** Validator is the type for the validator-set. */
export interface Validator {
    address: string;
    pubKey?: Any;
    votingPower: Long;
    proposerPriority: Long;
}
/** Validator is the type for the validator-set. */
export interface ValidatorSDKType {
    address: string;
    pub_key?: AnySDKType;
    voting_power: Long;
    proposer_priority: Long;
}
/** GetBlockByHeightRequest is the request type for the Query/GetBlockByHeight RPC method. */
export interface GetBlockByHeightRequest {
    height: Long;
}
/** GetBlockByHeightRequest is the request type for the Query/GetBlockByHeight RPC method. */
export interface GetBlockByHeightRequestSDKType {
    height: Long;
}
/** GetBlockByHeightResponse is the response type for the Query/GetBlockByHeight RPC method. */
export interface GetBlockByHeightResponse {
    blockId?: BlockID;
    /** Deprecated: please use `sdk_block` instead */
    block?: Block1;
    /** Since: cosmos-sdk 0.47 */
    sdkBlock?: Block2;
}
/** GetBlockByHeightResponse is the response type for the Query/GetBlockByHeight RPC method. */
export interface GetBlockByHeightResponseSDKType {
    block_id?: BlockIDSDKType;
    block?: Block1SDKType;
    sdk_block?: Block2SDKType;
}
/** GetLatestBlockRequest is the request type for the Query/GetLatestBlock RPC method. */
export interface GetLatestBlockRequest {
}
/** GetLatestBlockRequest is the request type for the Query/GetLatestBlock RPC method. */
export interface GetLatestBlockRequestSDKType {
}
/** GetLatestBlockResponse is the response type for the Query/GetLatestBlock RPC method. */
export interface GetLatestBlockResponse {
    blockId?: BlockID;
    /** Deprecated: please use `sdk_block` instead */
    block?: Block1;
    /** Since: cosmos-sdk 0.47 */
    sdkBlock?: Block2;
}
/** GetLatestBlockResponse is the response type for the Query/GetLatestBlock RPC method. */
export interface GetLatestBlockResponseSDKType {
    block_id?: BlockIDSDKType;
    block?: Block1SDKType;
    sdk_block?: Block2SDKType;
}
/** GetSyncingRequest is the request type for the Query/GetSyncing RPC method. */
export interface GetSyncingRequest {
}
/** GetSyncingRequest is the request type for the Query/GetSyncing RPC method. */
export interface GetSyncingRequestSDKType {
}
/** GetSyncingResponse is the response type for the Query/GetSyncing RPC method. */
export interface GetSyncingResponse {
    syncing: boolean;
}
/** GetSyncingResponse is the response type for the Query/GetSyncing RPC method. */
export interface GetSyncingResponseSDKType {
    syncing: boolean;
}
/** GetNodeInfoRequest is the request type for the Query/GetNodeInfo RPC method. */
export interface GetNodeInfoRequest {
}
/** GetNodeInfoRequest is the request type for the Query/GetNodeInfo RPC method. */
export interface GetNodeInfoRequestSDKType {
}
/** GetNodeInfoResponse is the response type for the Query/GetNodeInfo RPC method. */
export interface GetNodeInfoResponse {
    defaultNodeInfo?: DefaultNodeInfo;
    applicationVersion?: VersionInfo;
}
/** GetNodeInfoResponse is the response type for the Query/GetNodeInfo RPC method. */
export interface GetNodeInfoResponseSDKType {
    default_node_info?: DefaultNodeInfoSDKType;
    application_version?: VersionInfoSDKType;
}
/** VersionInfo is the type for the GetNodeInfoResponse message. */
export interface VersionInfo {
    name: string;
    appName: string;
    version: string;
    gitCommit: string;
    buildTags: string;
    goVersion: string;
    buildDeps: Module[];
    /** Since: cosmos-sdk 0.43 */
    cosmosSdkVersion: string;
}
/** VersionInfo is the type for the GetNodeInfoResponse message. */
export interface VersionInfoSDKType {
    name: string;
    app_name: string;
    version: string;
    git_commit: string;
    build_tags: string;
    go_version: string;
    build_deps: ModuleSDKType[];
    cosmos_sdk_version: string;
}
/** Module is the type for VersionInfo */
export interface Module {
    /** module path */
    path: string;
    /** module version */
    version: string;
    /** checksum */
    sum: string;
}
/** Module is the type for VersionInfo */
export interface ModuleSDKType {
    path: string;
    version: string;
    sum: string;
}
/** ABCIQueryRequest defines the request structure for the ABCIQuery gRPC query. */
export interface ABCIQueryRequest {
    data: Uint8Array;
    path: string;
    height: Long;
    prove: boolean;
}
/** ABCIQueryRequest defines the request structure for the ABCIQuery gRPC query. */
export interface ABCIQueryRequestSDKType {
    data: Uint8Array;
    path: string;
    height: Long;
    prove: boolean;
}
/**
 * ABCIQueryResponse defines the response structure for the ABCIQuery gRPC query.
 *
 * Note: This type is a duplicate of the ResponseQuery proto type defined in
 * Tendermint.
 */
export interface ABCIQueryResponse {
    code: number;
    /** nondeterministic */
    log: string;
    /** nondeterministic */
    info: string;
    index: Long;
    key: Uint8Array;
    value: Uint8Array;
    proofOps?: ProofOps;
    height: Long;
    codespace: string;
}
/**
 * ABCIQueryResponse defines the response structure for the ABCIQuery gRPC query.
 *
 * Note: This type is a duplicate of the ResponseQuery proto type defined in
 * Tendermint.
 */
export interface ABCIQueryResponseSDKType {
    code: number;
    log: string;
    info: string;
    index: Long;
    key: Uint8Array;
    value: Uint8Array;
    proof_ops?: ProofOpsSDKType;
    height: Long;
    codespace: string;
}
/**
 * ProofOp defines an operation used for calculating Merkle root. The data could
 * be arbitrary format, providing necessary data for example neighbouring node
 * hash.
 *
 * Note: This type is a duplicate of the ProofOp proto type defined in Tendermint.
 */
export interface ProofOp {
    type: string;
    key: Uint8Array;
    data: Uint8Array;
}
/**
 * ProofOp defines an operation used for calculating Merkle root. The data could
 * be arbitrary format, providing necessary data for example neighbouring node
 * hash.
 *
 * Note: This type is a duplicate of the ProofOp proto type defined in Tendermint.
 */
export interface ProofOpSDKType {
    type: string;
    key: Uint8Array;
    data: Uint8Array;
}
/**
 * ProofOps is Merkle proof defined by the list of ProofOps.
 *
 * Note: This type is a duplicate of the ProofOps proto type defined in Tendermint.
 */
export interface ProofOps {
    ops: ProofOp[];
}
/**
 * ProofOps is Merkle proof defined by the list of ProofOps.
 *
 * Note: This type is a duplicate of the ProofOps proto type defined in Tendermint.
 */
export interface ProofOpsSDKType {
    ops: ProofOpSDKType[];
}
export declare const GetValidatorSetByHeightRequest: {
    encode(message: GetValidatorSetByHeightRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetValidatorSetByHeightRequest;
    fromPartial(object: DeepPartial<GetValidatorSetByHeightRequest>): GetValidatorSetByHeightRequest;
};
export declare const GetValidatorSetByHeightResponse: {
    encode(message: GetValidatorSetByHeightResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetValidatorSetByHeightResponse;
    fromPartial(object: DeepPartial<GetValidatorSetByHeightResponse>): GetValidatorSetByHeightResponse;
};
export declare const GetLatestValidatorSetRequest: {
    encode(message: GetLatestValidatorSetRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetLatestValidatorSetRequest;
    fromPartial(object: DeepPartial<GetLatestValidatorSetRequest>): GetLatestValidatorSetRequest;
};
export declare const GetLatestValidatorSetResponse: {
    encode(message: GetLatestValidatorSetResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetLatestValidatorSetResponse;
    fromPartial(object: DeepPartial<GetLatestValidatorSetResponse>): GetLatestValidatorSetResponse;
};
export declare const Validator: {
    encode(message: Validator, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Validator;
    fromPartial(object: DeepPartial<Validator>): Validator;
};
export declare const GetBlockByHeightRequest: {
    encode(message: GetBlockByHeightRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetBlockByHeightRequest;
    fromPartial(object: DeepPartial<GetBlockByHeightRequest>): GetBlockByHeightRequest;
};
export declare const GetBlockByHeightResponse: {
    encode(message: GetBlockByHeightResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetBlockByHeightResponse;
    fromPartial(object: DeepPartial<GetBlockByHeightResponse>): GetBlockByHeightResponse;
};
export declare const GetLatestBlockRequest: {
    encode(_: GetLatestBlockRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetLatestBlockRequest;
    fromPartial(_: DeepPartial<GetLatestBlockRequest>): GetLatestBlockRequest;
};
export declare const GetLatestBlockResponse: {
    encode(message: GetLatestBlockResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetLatestBlockResponse;
    fromPartial(object: DeepPartial<GetLatestBlockResponse>): GetLatestBlockResponse;
};
export declare const GetSyncingRequest: {
    encode(_: GetSyncingRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetSyncingRequest;
    fromPartial(_: DeepPartial<GetSyncingRequest>): GetSyncingRequest;
};
export declare const GetSyncingResponse: {
    encode(message: GetSyncingResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetSyncingResponse;
    fromPartial(object: DeepPartial<GetSyncingResponse>): GetSyncingResponse;
};
export declare const GetNodeInfoRequest: {
    encode(_: GetNodeInfoRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetNodeInfoRequest;
    fromPartial(_: DeepPartial<GetNodeInfoRequest>): GetNodeInfoRequest;
};
export declare const GetNodeInfoResponse: {
    encode(message: GetNodeInfoResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetNodeInfoResponse;
    fromPartial(object: DeepPartial<GetNodeInfoResponse>): GetNodeInfoResponse;
};
export declare const VersionInfo: {
    encode(message: VersionInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): VersionInfo;
    fromPartial(object: DeepPartial<VersionInfo>): VersionInfo;
};
export declare const Module: {
    encode(message: Module, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Module;
    fromPartial(object: DeepPartial<Module>): Module;
};
export declare const ABCIQueryRequest: {
    encode(message: ABCIQueryRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ABCIQueryRequest;
    fromPartial(object: DeepPartial<ABCIQueryRequest>): ABCIQueryRequest;
};
export declare const ABCIQueryResponse: {
    encode(message: ABCIQueryResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ABCIQueryResponse;
    fromPartial(object: DeepPartial<ABCIQueryResponse>): ABCIQueryResponse;
};
export declare const ProofOp: {
    encode(message: ProofOp, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ProofOp;
    fromPartial(object: DeepPartial<ProofOp>): ProofOp;
};
export declare const ProofOps: {
    encode(message: ProofOps, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ProofOps;
    fromPartial(object: DeepPartial<ProofOps>): ProofOps;
};
