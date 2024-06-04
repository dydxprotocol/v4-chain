import { ParamChange } from "./params";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.params.v1beta1";
/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
    /** subspace defines the module to query the parameter for. */
    subspace: string;
    /** key defines the key of the parameter in the subspace. */
    key: string;
}
/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
    /** param defines the queried parameter. */
    param: ParamChange;
}
/**
 * QuerySubspacesRequest defines a request type for querying for all registered
 * subspaces and all keys for a subspace.
 *
 * Since: cosmos-sdk 0.46
 */
export interface QuerySubspacesRequest {
}
/**
 * QuerySubspacesResponse defines the response types for querying for all
 * registered subspaces and all keys for a subspace.
 *
 * Since: cosmos-sdk 0.46
 */
export interface QuerySubspacesResponse {
    subspaces: Subspace[];
}
/**
 * Subspace defines a parameter subspace name and all the keys that exist for
 * the subspace.
 *
 * Since: cosmos-sdk 0.46
 */
export interface Subspace {
    subspace: string;
    keys: string[];
}
export declare const QueryParamsRequest: {
    typeUrl: string;
    encode(message: QueryParamsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsRequest;
    fromJSON(object: any): QueryParamsRequest;
    toJSON(message: QueryParamsRequest): unknown;
    fromPartial<I extends {
        subspace?: string | undefined;
        key?: string | undefined;
    } & {
        subspace?: string | undefined;
        key?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryParamsRequest>, never>>(object: I): QueryParamsRequest;
};
export declare const QueryParamsResponse: {
    typeUrl: string;
    encode(message: QueryParamsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryParamsResponse;
    fromJSON(object: any): QueryParamsResponse;
    toJSON(message: QueryParamsResponse): unknown;
    fromPartial<I extends {
        param?: {
            subspace?: string | undefined;
            key?: string | undefined;
            value?: string | undefined;
        } | undefined;
    } & {
        param?: ({
            subspace?: string | undefined;
            key?: string | undefined;
            value?: string | undefined;
        } & {
            subspace?: string | undefined;
            key?: string | undefined;
            value?: string | undefined;
        } & Record<Exclude<keyof I["param"], keyof ParamChange>, never>) | undefined;
    } & Record<Exclude<keyof I, "param">, never>>(object: I): QueryParamsResponse;
};
export declare const QuerySubspacesRequest: {
    typeUrl: string;
    encode(_: QuerySubspacesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySubspacesRequest;
    fromJSON(_: any): QuerySubspacesRequest;
    toJSON(_: QuerySubspacesRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): QuerySubspacesRequest;
};
export declare const QuerySubspacesResponse: {
    typeUrl: string;
    encode(message: QuerySubspacesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySubspacesResponse;
    fromJSON(object: any): QuerySubspacesResponse;
    toJSON(message: QuerySubspacesResponse): unknown;
    fromPartial<I extends {
        subspaces?: {
            subspace?: string | undefined;
            keys?: string[] | undefined;
        }[] | undefined;
    } & {
        subspaces?: ({
            subspace?: string | undefined;
            keys?: string[] | undefined;
        }[] & ({
            subspace?: string | undefined;
            keys?: string[] | undefined;
        } & {
            subspace?: string | undefined;
            keys?: (string[] & string[] & Record<Exclude<keyof I["subspaces"][number]["keys"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["subspaces"][number], keyof Subspace>, never>)[] & Record<Exclude<keyof I["subspaces"], keyof {
            subspace?: string | undefined;
            keys?: string[] | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "subspaces">, never>>(object: I): QuerySubspacesResponse;
};
export declare const Subspace: {
    typeUrl: string;
    encode(message: Subspace, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Subspace;
    fromJSON(object: any): Subspace;
    toJSON(message: Subspace): unknown;
    fromPartial<I extends {
        subspace?: string | undefined;
        keys?: string[] | undefined;
    } & {
        subspace?: string | undefined;
        keys?: (string[] & string[] & Record<Exclude<keyof I["keys"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Subspace>, never>>(object: I): Subspace;
};
/** Query defines the gRPC querier service. */
export interface Query {
    /**
     * Params queries a specific parameter of a module, given its subspace and
     * key.
     */
    Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
    /**
     * Subspaces queries for all registered subspaces and all keys for a subspace.
     *
     * Since: cosmos-sdk 0.46
     */
    Subspaces(request?: QuerySubspacesRequest): Promise<QuerySubspacesResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
    Subspaces(request?: QuerySubspacesRequest): Promise<QuerySubspacesResponse>;
}
