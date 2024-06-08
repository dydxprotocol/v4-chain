import { Plan, ModuleVersion } from "./upgrade";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.upgrade.v1beta1";
/**
 * QueryCurrentPlanRequest is the request type for the Query/CurrentPlan RPC
 * method.
 */
export interface QueryCurrentPlanRequest {
}
/**
 * QueryCurrentPlanResponse is the response type for the Query/CurrentPlan RPC
 * method.
 */
export interface QueryCurrentPlanResponse {
    /** plan is the current upgrade plan. */
    plan?: Plan;
}
/**
 * QueryCurrentPlanRequest is the request type for the Query/AppliedPlan RPC
 * method.
 */
export interface QueryAppliedPlanRequest {
    /** name is the name of the applied plan to query for. */
    name: string;
}
/**
 * QueryAppliedPlanResponse is the response type for the Query/AppliedPlan RPC
 * method.
 */
export interface QueryAppliedPlanResponse {
    /** height is the block height at which the plan was applied. */
    height: bigint;
}
/**
 * QueryUpgradedConsensusStateRequest is the request type for the Query/UpgradedConsensusState
 * RPC method.
 */
/** @deprecated */
export interface QueryUpgradedConsensusStateRequest {
    /**
     * last height of the current chain must be sent in request
     * as this is the height under which next consensus state is stored
     */
    lastHeight: bigint;
}
/**
 * QueryUpgradedConsensusStateResponse is the response type for the Query/UpgradedConsensusState
 * RPC method.
 */
/** @deprecated */
export interface QueryUpgradedConsensusStateResponse {
    /** Since: cosmos-sdk 0.43 */
    upgradedConsensusState: Uint8Array;
}
/**
 * QueryModuleVersionsRequest is the request type for the Query/ModuleVersions
 * RPC method.
 *
 * Since: cosmos-sdk 0.43
 */
export interface QueryModuleVersionsRequest {
    /**
     * module_name is a field to query a specific module
     * consensus version from state. Leaving this empty will
     * fetch the full list of module versions from state
     */
    moduleName: string;
}
/**
 * QueryModuleVersionsResponse is the response type for the Query/ModuleVersions
 * RPC method.
 *
 * Since: cosmos-sdk 0.43
 */
export interface QueryModuleVersionsResponse {
    /** module_versions is a list of module names with their consensus versions. */
    moduleVersions: ModuleVersion[];
}
/**
 * QueryAuthorityRequest is the request type for Query/Authority
 *
 * Since: cosmos-sdk 0.46
 */
export interface QueryAuthorityRequest {
}
/**
 * QueryAuthorityResponse is the response type for Query/Authority
 *
 * Since: cosmos-sdk 0.46
 */
export interface QueryAuthorityResponse {
    address: string;
}
export declare const QueryCurrentPlanRequest: {
    typeUrl: string;
    encode(_: QueryCurrentPlanRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryCurrentPlanRequest;
    fromJSON(_: any): QueryCurrentPlanRequest;
    toJSON(_: QueryCurrentPlanRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): QueryCurrentPlanRequest;
};
export declare const QueryCurrentPlanResponse: {
    typeUrl: string;
    encode(message: QueryCurrentPlanResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryCurrentPlanResponse;
    fromJSON(object: any): QueryCurrentPlanResponse;
    toJSON(message: QueryCurrentPlanResponse): unknown;
    fromPartial<I extends {
        plan?: {
            name?: string | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            height?: bigint | undefined;
            info?: string | undefined;
            upgradedClientState?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
    } & {
        plan?: ({
            name?: string | undefined;
            time?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            height?: bigint | undefined;
            info?: string | undefined;
            upgradedClientState?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            name?: string | undefined;
            time?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["plan"]["time"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            height?: bigint | undefined;
            info?: string | undefined;
            upgradedClientState?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["plan"]["upgradedClientState"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["plan"], keyof Plan>, never>) | undefined;
    } & Record<Exclude<keyof I, "plan">, never>>(object: I): QueryCurrentPlanResponse;
};
export declare const QueryAppliedPlanRequest: {
    typeUrl: string;
    encode(message: QueryAppliedPlanRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAppliedPlanRequest;
    fromJSON(object: any): QueryAppliedPlanRequest;
    toJSON(message: QueryAppliedPlanRequest): unknown;
    fromPartial<I extends {
        name?: string | undefined;
    } & {
        name?: string | undefined;
    } & Record<Exclude<keyof I, "name">, never>>(object: I): QueryAppliedPlanRequest;
};
export declare const QueryAppliedPlanResponse: {
    typeUrl: string;
    encode(message: QueryAppliedPlanResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAppliedPlanResponse;
    fromJSON(object: any): QueryAppliedPlanResponse;
    toJSON(message: QueryAppliedPlanResponse): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
    } & {
        height?: bigint | undefined;
    } & Record<Exclude<keyof I, "height">, never>>(object: I): QueryAppliedPlanResponse;
};
export declare const QueryUpgradedConsensusStateRequest: {
    typeUrl: string;
    encode(message: QueryUpgradedConsensusStateRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryUpgradedConsensusStateRequest;
    fromJSON(object: any): QueryUpgradedConsensusStateRequest;
    toJSON(message: QueryUpgradedConsensusStateRequest): unknown;
    fromPartial<I extends {
        lastHeight?: bigint | undefined;
    } & {
        lastHeight?: bigint | undefined;
    } & Record<Exclude<keyof I, "lastHeight">, never>>(object: I): QueryUpgradedConsensusStateRequest;
};
export declare const QueryUpgradedConsensusStateResponse: {
    typeUrl: string;
    encode(message: QueryUpgradedConsensusStateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryUpgradedConsensusStateResponse;
    fromJSON(object: any): QueryUpgradedConsensusStateResponse;
    toJSON(message: QueryUpgradedConsensusStateResponse): unknown;
    fromPartial<I extends {
        upgradedConsensusState?: Uint8Array | undefined;
    } & {
        upgradedConsensusState?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "upgradedConsensusState">, never>>(object: I): QueryUpgradedConsensusStateResponse;
};
export declare const QueryModuleVersionsRequest: {
    typeUrl: string;
    encode(message: QueryModuleVersionsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryModuleVersionsRequest;
    fromJSON(object: any): QueryModuleVersionsRequest;
    toJSON(message: QueryModuleVersionsRequest): unknown;
    fromPartial<I extends {
        moduleName?: string | undefined;
    } & {
        moduleName?: string | undefined;
    } & Record<Exclude<keyof I, "moduleName">, never>>(object: I): QueryModuleVersionsRequest;
};
export declare const QueryModuleVersionsResponse: {
    typeUrl: string;
    encode(message: QueryModuleVersionsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryModuleVersionsResponse;
    fromJSON(object: any): QueryModuleVersionsResponse;
    toJSON(message: QueryModuleVersionsResponse): unknown;
    fromPartial<I extends {
        moduleVersions?: {
            name?: string | undefined;
            version?: bigint | undefined;
        }[] | undefined;
    } & {
        moduleVersions?: ({
            name?: string | undefined;
            version?: bigint | undefined;
        }[] & ({
            name?: string | undefined;
            version?: bigint | undefined;
        } & {
            name?: string | undefined;
            version?: bigint | undefined;
        } & Record<Exclude<keyof I["moduleVersions"][number], keyof ModuleVersion>, never>)[] & Record<Exclude<keyof I["moduleVersions"], keyof {
            name?: string | undefined;
            version?: bigint | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "moduleVersions">, never>>(object: I): QueryModuleVersionsResponse;
};
export declare const QueryAuthorityRequest: {
    typeUrl: string;
    encode(_: QueryAuthorityRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAuthorityRequest;
    fromJSON(_: any): QueryAuthorityRequest;
    toJSON(_: QueryAuthorityRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): QueryAuthorityRequest;
};
export declare const QueryAuthorityResponse: {
    typeUrl: string;
    encode(message: QueryAuthorityResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAuthorityResponse;
    fromJSON(object: any): QueryAuthorityResponse;
    toJSON(message: QueryAuthorityResponse): unknown;
    fromPartial<I extends {
        address?: string | undefined;
    } & {
        address?: string | undefined;
    } & Record<Exclude<keyof I, "address">, never>>(object: I): QueryAuthorityResponse;
};
/** Query defines the gRPC upgrade querier service. */
export interface Query {
    /** CurrentPlan queries the current upgrade plan. */
    CurrentPlan(request?: QueryCurrentPlanRequest): Promise<QueryCurrentPlanResponse>;
    /** AppliedPlan queries a previously applied upgrade plan by its name. */
    AppliedPlan(request: QueryAppliedPlanRequest): Promise<QueryAppliedPlanResponse>;
    /**
     * UpgradedConsensusState queries the consensus state that will serve
     * as a trusted kernel for the next version of this chain. It will only be
     * stored at the last height of this chain.
     * UpgradedConsensusState RPC not supported with legacy querier
     * This rpc is deprecated now that IBC has its own replacement
     * (https://github.com/cosmos/ibc-go/blob/2c880a22e9f9cc75f62b527ca94aa75ce1106001/proto/ibc/core/client/v1/query.proto#L54)
     */
    UpgradedConsensusState(request: QueryUpgradedConsensusStateRequest): Promise<QueryUpgradedConsensusStateResponse>;
    /**
     * ModuleVersions queries the list of module versions from state.
     *
     * Since: cosmos-sdk 0.43
     */
    ModuleVersions(request: QueryModuleVersionsRequest): Promise<QueryModuleVersionsResponse>;
    /**
     * Returns the account with authority to conduct upgrades
     *
     * Since: cosmos-sdk 0.46
     */
    Authority(request?: QueryAuthorityRequest): Promise<QueryAuthorityResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    CurrentPlan(request?: QueryCurrentPlanRequest): Promise<QueryCurrentPlanResponse>;
    AppliedPlan(request: QueryAppliedPlanRequest): Promise<QueryAppliedPlanResponse>;
    UpgradedConsensusState(request: QueryUpgradedConsensusStateRequest): Promise<QueryUpgradedConsensusStateResponse>;
    ModuleVersions(request: QueryModuleVersionsRequest): Promise<QueryModuleVersionsResponse>;
    Authority(request?: QueryAuthorityRequest): Promise<QueryAuthorityResponse>;
}
