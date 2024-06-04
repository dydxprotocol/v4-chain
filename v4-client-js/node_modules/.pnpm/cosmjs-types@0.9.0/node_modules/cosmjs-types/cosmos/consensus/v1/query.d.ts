import { ConsensusParams } from "../../../tendermint/types/params";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.consensus.v1";
/** QueryParamsRequest defines the request type for querying x/consensus parameters. */
export interface QueryParamsRequest {
}
/** QueryParamsResponse defines the response type for querying x/consensus parameters. */
export interface QueryParamsResponse {
    /**
     * params are the tendermint consensus params stored in the consensus module.
     * Please note that `params.version` is not populated in this response, it is
     * tracked separately in the x/upgrade module.
     */
    params?: ConsensusParams;
}
export declare const QueryParamsRequest: {
    typeUrl: string;
    encode(_: QueryParamsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsRequest;
    fromJSON(_: any): QueryParamsRequest;
    toJSON(_: QueryParamsRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): QueryParamsRequest;
};
export declare const QueryParamsResponse: {
    typeUrl: string;
    encode(message: QueryParamsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsResponse;
    fromJSON(object: any): QueryParamsResponse;
    toJSON(message: QueryParamsResponse): unknown;
    fromPartial<I extends {
        params?: {
            block?: {
                maxBytes?: bigint | undefined;
                maxGas?: bigint | undefined;
            } | undefined;
            evidence?: {
                maxAgeNumBlocks?: bigint | undefined;
                maxAgeDuration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                maxBytes?: bigint | undefined;
            } | undefined;
            validator?: {
                pubKeyTypes?: string[] | undefined;
            } | undefined;
            version?: {
                app?: bigint | undefined;
            } | undefined;
        } | undefined;
    } & {
        params?: ({
            block?: {
                maxBytes?: bigint | undefined;
                maxGas?: bigint | undefined;
            } | undefined;
            evidence?: {
                maxAgeNumBlocks?: bigint | undefined;
                maxAgeDuration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                maxBytes?: bigint | undefined;
            } | undefined;
            validator?: {
                pubKeyTypes?: string[] | undefined;
            } | undefined;
            version?: {
                app?: bigint | undefined;
            } | undefined;
        } & {
            block?: ({
                maxBytes?: bigint | undefined;
                maxGas?: bigint | undefined;
            } & {
                maxBytes?: bigint | undefined;
                maxGas?: bigint | undefined;
            } & Record<Exclude<keyof I["params"]["block"], keyof import("../../../tendermint/types/params").BlockParams>, never>) | undefined;
            evidence?: ({
                maxAgeNumBlocks?: bigint | undefined;
                maxAgeDuration?: {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } | undefined;
                maxBytes?: bigint | undefined;
            } & {
                maxAgeNumBlocks?: bigint | undefined;
                maxAgeDuration?: ({
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & {
                    seconds?: bigint | undefined;
                    nanos?: number | undefined;
                } & Record<Exclude<keyof I["params"]["evidence"]["maxAgeDuration"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
                maxBytes?: bigint | undefined;
            } & Record<Exclude<keyof I["params"]["evidence"], keyof import("../../../tendermint/types/params").EvidenceParams>, never>) | undefined;
            validator?: ({
                pubKeyTypes?: string[] | undefined;
            } & {
                pubKeyTypes?: (string[] & string[] & Record<Exclude<keyof I["params"]["validator"]["pubKeyTypes"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["params"]["validator"], "pubKeyTypes">, never>) | undefined;
            version?: ({
                app?: bigint | undefined;
            } & {
                app?: bigint | undefined;
            } & Record<Exclude<keyof I["params"]["version"], "app">, never>) | undefined;
        } & Record<Exclude<keyof I["params"], keyof ConsensusParams>, never>) | undefined;
    } & Record<Exclude<keyof I, "params">, never>>(object: I): QueryParamsResponse;
};
/** Query defines the gRPC querier service. */
export interface Query {
    /** Params queries the parameters of x/consensus_param module. */
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
}
