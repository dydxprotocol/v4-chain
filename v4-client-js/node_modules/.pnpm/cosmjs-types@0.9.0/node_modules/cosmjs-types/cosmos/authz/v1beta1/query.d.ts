import { PageRequest, PageResponse } from "../../base/query/v1beta1/pagination";
import { Grant, GrantAuthorization } from "./authz";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.authz.v1beta1";
/** QueryGrantsRequest is the request type for the Query/Grants RPC method. */
export interface QueryGrantsRequest {
    granter: string;
    grantee: string;
    /** Optional, msg_type_url, when set, will query only grants matching given msg type. */
    msgTypeUrl: string;
    /** pagination defines an pagination for the request. */
    pagination?: PageRequest;
}
/** QueryGrantsResponse is the response type for the Query/Authorizations RPC method. */
export interface QueryGrantsResponse {
    /** authorizations is a list of grants granted for grantee by granter. */
    grants: Grant[];
    /** pagination defines an pagination for the response. */
    pagination?: PageResponse;
}
/** QueryGranterGrantsRequest is the request type for the Query/GranterGrants RPC method. */
export interface QueryGranterGrantsRequest {
    granter: string;
    /** pagination defines an pagination for the request. */
    pagination?: PageRequest;
}
/** QueryGranterGrantsResponse is the response type for the Query/GranterGrants RPC method. */
export interface QueryGranterGrantsResponse {
    /** grants is a list of grants granted by the granter. */
    grants: GrantAuthorization[];
    /** pagination defines an pagination for the response. */
    pagination?: PageResponse;
}
/** QueryGranteeGrantsRequest is the request type for the Query/IssuedGrants RPC method. */
export interface QueryGranteeGrantsRequest {
    grantee: string;
    /** pagination defines an pagination for the request. */
    pagination?: PageRequest;
}
/** QueryGranteeGrantsResponse is the response type for the Query/GranteeGrants RPC method. */
export interface QueryGranteeGrantsResponse {
    /** grants is a list of grants granted to the grantee. */
    grants: GrantAuthorization[];
    /** pagination defines an pagination for the response. */
    pagination?: PageResponse;
}
export declare const QueryGrantsRequest: {
    typeUrl: string;
    encode(message: QueryGrantsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGrantsRequest;
    fromJSON(object: any): QueryGrantsRequest;
    toJSON(message: QueryGrantsRequest): unknown;
    fromPartial<I extends {
        granter?: string | undefined;
        grantee?: string | undefined;
        msgTypeUrl?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        granter?: string | undefined;
        grantee?: string | undefined;
        msgTypeUrl?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryGrantsRequest>, never>>(object: I): QueryGrantsRequest;
};
export declare const QueryGrantsResponse: {
    typeUrl: string;
    encode(message: QueryGrantsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGrantsResponse;
    fromJSON(object: any): QueryGrantsResponse;
    toJSON(message: QueryGrantsResponse): unknown;
    fromPartial<I extends {
        grants?: {
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        grants?: ({
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] & ({
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            authorization?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["grants"][number]["authorization"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
            expiration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["grants"][number]["expiration"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["grants"][number], keyof Grant>, never>)[] & Record<Exclude<keyof I["grants"], keyof {
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryGrantsResponse>, never>>(object: I): QueryGrantsResponse;
};
export declare const QueryGranterGrantsRequest: {
    typeUrl: string;
    encode(message: QueryGranterGrantsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGranterGrantsRequest;
    fromJSON(object: any): QueryGranterGrantsRequest;
    toJSON(message: QueryGranterGrantsRequest): unknown;
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
    } & Record<Exclude<keyof I, keyof QueryGranterGrantsRequest>, never>>(object: I): QueryGranterGrantsRequest;
};
export declare const QueryGranterGrantsResponse: {
    typeUrl: string;
    encode(message: QueryGranterGrantsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGranterGrantsResponse;
    fromJSON(object: any): QueryGranterGrantsResponse;
    toJSON(message: QueryGranterGrantsResponse): unknown;
    fromPartial<I extends {
        grants?: {
            granter?: string | undefined;
            grantee?: string | undefined;
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        grants?: ({
            granter?: string | undefined;
            grantee?: string | undefined;
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] & ({
            granter?: string | undefined;
            grantee?: string | undefined;
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            granter?: string | undefined;
            grantee?: string | undefined;
            authorization?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["grants"][number]["authorization"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
            expiration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["grants"][number]["expiration"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["grants"][number], keyof GrantAuthorization>, never>)[] & Record<Exclude<keyof I["grants"], keyof {
            granter?: string | undefined;
            grantee?: string | undefined;
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryGranterGrantsResponse>, never>>(object: I): QueryGranterGrantsResponse;
};
export declare const QueryGranteeGrantsRequest: {
    typeUrl: string;
    encode(message: QueryGranteeGrantsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGranteeGrantsRequest;
    fromJSON(object: any): QueryGranteeGrantsRequest;
    toJSON(message: QueryGranteeGrantsRequest): unknown;
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
    } & Record<Exclude<keyof I, keyof QueryGranteeGrantsRequest>, never>>(object: I): QueryGranteeGrantsRequest;
};
export declare const QueryGranteeGrantsResponse: {
    typeUrl: string;
    encode(message: QueryGranteeGrantsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryGranteeGrantsResponse;
    fromJSON(object: any): QueryGranteeGrantsResponse;
    toJSON(message: QueryGranteeGrantsResponse): unknown;
    fromPartial<I extends {
        grants?: {
            granter?: string | undefined;
            grantee?: string | undefined;
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        grants?: ({
            granter?: string | undefined;
            grantee?: string | undefined;
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[] & ({
            granter?: string | undefined;
            grantee?: string | undefined;
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        } & {
            granter?: string | undefined;
            grantee?: string | undefined;
            authorization?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["grants"][number]["authorization"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
            expiration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["grants"][number]["expiration"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
        } & Record<Exclude<keyof I["grants"][number], keyof GrantAuthorization>, never>)[] & Record<Exclude<keyof I["grants"], keyof {
            granter?: string | undefined;
            grantee?: string | undefined;
            authorization?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            expiration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryGranteeGrantsResponse>, never>>(object: I): QueryGranteeGrantsResponse;
};
/** Query defines the gRPC querier service. */
export interface Query {
    /** Returns list of `Authorization`, granted to the grantee by the granter. */
    Grants(request: QueryGrantsRequest): Promise<QueryGrantsResponse>;
    /**
     * GranterGrants returns list of `GrantAuthorization`, granted by granter.
     *
     * Since: cosmos-sdk 0.46
     */
    GranterGrants(request: QueryGranterGrantsRequest): Promise<QueryGranterGrantsResponse>;
    /**
     * GranteeGrants returns a list of `GrantAuthorization` by grantee.
     *
     * Since: cosmos-sdk 0.46
     */
    GranteeGrants(request: QueryGranteeGrantsRequest): Promise<QueryGranteeGrantsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Grants(request: QueryGrantsRequest): Promise<QueryGrantsResponse>;
    GranterGrants(request: QueryGranterGrantsRequest): Promise<QueryGranterGrantsResponse>;
    GranteeGrants(request: QueryGranteeGrantsRequest): Promise<QueryGranteeGrantsResponse>;
}
