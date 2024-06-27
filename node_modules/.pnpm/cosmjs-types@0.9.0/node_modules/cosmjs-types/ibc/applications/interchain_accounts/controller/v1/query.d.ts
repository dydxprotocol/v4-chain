import { Params } from "./controller";
import { BinaryReader, BinaryWriter } from "../../../../../binary";
import { Rpc } from "../../../../../helpers";
export declare const protobufPackage = "ibc.applications.interchain_accounts.controller.v1";
/** QueryInterchainAccountRequest is the request type for the Query/InterchainAccount RPC method. */
export interface QueryInterchainAccountRequest {
    owner: string;
    connectionId: string;
}
/** QueryInterchainAccountResponse the response type for the Query/InterchainAccount RPC method. */
export interface QueryInterchainAccountResponse {
    address: string;
}
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
    /** params defines the parameters of the module. */
    params?: Params;
}
export declare const QueryInterchainAccountRequest: {
    typeUrl: string;
    encode(message: QueryInterchainAccountRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryInterchainAccountRequest;
    fromJSON(object: any): QueryInterchainAccountRequest;
    toJSON(message: QueryInterchainAccountRequest): unknown;
    fromPartial<I extends {
        owner?: string | undefined;
        connectionId?: string | undefined;
    } & {
        owner?: string | undefined;
        connectionId?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryInterchainAccountRequest>, never>>(object: I): QueryInterchainAccountRequest;
};
export declare const QueryInterchainAccountResponse: {
    typeUrl: string;
    encode(message: QueryInterchainAccountResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryInterchainAccountResponse;
    fromJSON(object: any): QueryInterchainAccountResponse;
    toJSON(message: QueryInterchainAccountResponse): unknown;
    fromPartial<I extends {
        address?: string | undefined;
    } & {
        address?: string | undefined;
    } & Record<Exclude<keyof I, "address">, never>>(object: I): QueryInterchainAccountResponse;
};
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
            controllerEnabled?: boolean | undefined;
        } | undefined;
    } & {
        params?: ({
            controllerEnabled?: boolean | undefined;
        } & {
            controllerEnabled?: boolean | undefined;
        } & Record<Exclude<keyof I["params"], "controllerEnabled">, never>) | undefined;
    } & Record<Exclude<keyof I, "params">, never>>(object: I): QueryParamsResponse;
};
/** Query provides defines the gRPC querier service. */
export interface Query {
    /** InterchainAccount returns the interchain account address for a given owner address on a given connection */
    InterchainAccount(request: QueryInterchainAccountRequest): Promise<QueryInterchainAccountResponse>;
    /** Params queries all parameters of the ICA controller submodule. */
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    InterchainAccount(request: QueryInterchainAccountRequest): Promise<QueryInterchainAccountResponse>;
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
}
