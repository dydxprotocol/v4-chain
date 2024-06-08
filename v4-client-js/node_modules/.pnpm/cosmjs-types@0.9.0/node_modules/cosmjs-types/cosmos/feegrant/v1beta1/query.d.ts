import { PageRequest, PageResponse } from "../../base/query/v1beta1/pagination";
import { Grant } from "./feegrant";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.feegrant.v1beta1";
/** QueryAllowanceRequest is the request type for the Query/Allowance RPC method. */
export interface QueryAllowanceRequest {
    /** granter is the address of the user granting an allowance of their funds. */
    granter: string;
    /** grantee is the address of the user being granted an allowance of another user's funds. */
    grantee: string;
}
/** QueryAllowanceResponse is the response type for the Query/Allowance RPC method. */
export interface QueryAllowanceResponse {
    /** allowance is a allowance granted for grantee by granter. */
    allowance?: Grant;
}
/** QueryAllowancesRequest is the request type for the Query/Allowances RPC method. */
export interface QueryAllowancesRequest {
    grantee: string;
    /** pagination defines an pagination for the request. */
    pagination?: PageRequest;
}
/** QueryAllowancesResponse is the response type for the Query/Allowances RPC method. */
export interface QueryAllowancesResponse {
    /** allowances are allowance's granted for grantee by granter. */
    allowances: Grant[];
    /** pagination defines an pagination for the response. */
    pagination?: PageResponse;
}
/**
 * QueryAllowancesByGranterRequest is the request type for the Query/AllowancesByGranter RPC method.
 *
 * Since: cosmos-sdk 0.46
 */
export interface QueryAllowancesByGranterRequest {
    granter: string;
    /** pagination defines an pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryAllowancesByGranterResponse is the response type for the Query/AllowancesByGranter RPC method.
 *
 * Since: cosmos-sdk 0.46
 */
export interface QueryAllowancesByGranterResponse {
    /** allowances that have been issued by the granter. */
    allowances: Grant[];
    /** pagination defines an pagination for the response. */
    pagination?: PageResponse;
}
export declare const QueryAllowanceRequest: {
    typeUrl: string;
    encode(message: QueryAllowanceRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAllowanceRequest;
    fromJSON(object: any): QueryAllowanceRequest;
    toJSON(message: QueryAllowanceRequest): unknown;
    fromPartial<I extends {
        granter?: string | undefined;
        grantee?: string | undefined;
    } & {
        granter?: string | undefined;
        grantee?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryAllowanceRequest>, never>>(object: I): QueryAllowanceRequest;
};
export declare const QueryAllowanceResponse: {
    typeUrl: string;
    encode(message: QueryAllowanceResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAllowanceResponse;
    fromJSON(object: any): QueryAllowanceResponse;
    toJSON(message: QueryAllowanceResponse): unknown;
    fromPartial<I extends {
        allowance?: {
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
    } & {
        allowance?: ({
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["allowance"]["allowance"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["allowance"], keyof Grant>, never>) | undefined;
    } & Record<Exclude<keyof I, "allowance">, never>>(object: I): QueryAllowanceResponse;
};
export declare const QueryAllowancesRequest: {
    typeUrl: string;
    encode(message: QueryAllowancesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAllowancesRequest;
    fromJSON(object: any): QueryAllowancesRequest;
    toJSON(message: QueryAllowancesRequest): unknown;
    fromPartial<I extends {
        grantee?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        grantee?: string | undefined;
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryAllowancesRequest>, never>>(object: I): QueryAllowancesRequest;
};
export declare const QueryAllowancesResponse: {
    typeUrl: string;
    encode(message: QueryAllowancesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAllowancesResponse;
    fromJSON(object: any): QueryAllowancesResponse;
    toJSON(message: QueryAllowancesResponse): unknown;
    fromPartial<I extends {
        allowances?: {
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        allowances?: ({
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] & ({
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["allowances"][number]["allowance"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["allowances"][number], keyof Grant>, never>)[] & Record<Exclude<keyof I["allowances"], keyof {
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryAllowancesResponse>, never>>(object: I): QueryAllowancesResponse;
};
export declare const QueryAllowancesByGranterRequest: {
    typeUrl: string;
    encode(message: QueryAllowancesByGranterRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAllowancesByGranterRequest;
    fromJSON(object: any): QueryAllowancesByGranterRequest;
    toJSON(message: QueryAllowancesByGranterRequest): unknown;
    fromPartial<I extends {
        granter?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        granter?: string | undefined;
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryAllowancesByGranterRequest>, never>>(object: I): QueryAllowancesByGranterRequest;
};
export declare const QueryAllowancesByGranterResponse: {
    typeUrl: string;
    encode(message: QueryAllowancesByGranterResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAllowancesByGranterResponse;
    fromJSON(object: any): QueryAllowancesByGranterResponse;
    toJSON(message: QueryAllowancesByGranterResponse): unknown;
    fromPartial<I extends {
        allowances?: {
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        allowances?: ({
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] & ({
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["allowances"][number]["allowance"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["allowances"][number], keyof Grant>, never>)[] & Record<Exclude<keyof I["allowances"], keyof {
            granter?: string | undefined;
            grantee?: string | undefined;
            allowance?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryAllowancesByGranterResponse>, never>>(object: I): QueryAllowancesByGranterResponse;
};
/** Query defines the gRPC querier service. */
export interface Query {
    /** Allowance returns fee granted to the grantee by the granter. */
    Allowance(request: QueryAllowanceRequest): Promise<QueryAllowanceResponse>;
    /** Allowances returns all the grants for address. */
    Allowances(request: QueryAllowancesRequest): Promise<QueryAllowancesResponse>;
    /**
     * AllowancesByGranter returns all the grants given by an address
     *
     * Since: cosmos-sdk 0.46
     */
    AllowancesByGranter(request: QueryAllowancesByGranterRequest): Promise<QueryAllowancesByGranterResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Allowance(request: QueryAllowanceRequest): Promise<QueryAllowanceResponse>;
    Allowances(request: QueryAllowancesRequest): Promise<QueryAllowancesResponse>;
    AllowancesByGranter(request: QueryAllowancesByGranterRequest): Promise<QueryAllowancesByGranterResponse>;
}
