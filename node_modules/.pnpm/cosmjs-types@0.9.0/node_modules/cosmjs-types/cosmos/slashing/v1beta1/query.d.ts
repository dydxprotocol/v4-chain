import { PageRequest, PageResponse } from "../../base/query/v1beta1/pagination";
import { Params, ValidatorSigningInfo } from "./slashing";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.slashing.v1beta1";
/** QueryParamsRequest is the request type for the Query/Params RPC method */
export interface QueryParamsRequest {
}
/** QueryParamsResponse is the response type for the Query/Params RPC method */
export interface QueryParamsResponse {
    params: Params;
}
/**
 * QuerySigningInfoRequest is the request type for the Query/SigningInfo RPC
 * method
 */
export interface QuerySigningInfoRequest {
    /** cons_address is the address to query signing info of */
    consAddress: string;
}
/**
 * QuerySigningInfoResponse is the response type for the Query/SigningInfo RPC
 * method
 */
export interface QuerySigningInfoResponse {
    /** val_signing_info is the signing info of requested val cons address */
    valSigningInfo: ValidatorSigningInfo;
}
/**
 * QuerySigningInfosRequest is the request type for the Query/SigningInfos RPC
 * method
 */
export interface QuerySigningInfosRequest {
    pagination?: PageRequest;
}
/**
 * QuerySigningInfosResponse is the response type for the Query/SigningInfos RPC
 * method
 */
export interface QuerySigningInfosResponse {
    /** info is the signing info of all validators */
    info: ValidatorSigningInfo[];
    pagination?: PageResponse;
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
            signedBlocksWindow?: bigint | undefined;
            minSignedPerWindow?: Uint8Array | undefined;
            downtimeJailDuration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            slashFractionDoubleSign?: Uint8Array | undefined;
            slashFractionDowntime?: Uint8Array | undefined;
        } | undefined;
    } & {
        params?: ({
            signedBlocksWindow?: bigint | undefined;
            minSignedPerWindow?: Uint8Array | undefined;
            downtimeJailDuration?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            slashFractionDoubleSign?: Uint8Array | undefined;
            slashFractionDowntime?: Uint8Array | undefined;
        } & {
            signedBlocksWindow?: bigint | undefined;
            minSignedPerWindow?: Uint8Array | undefined;
            downtimeJailDuration?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["params"]["downtimeJailDuration"], keyof import("../../../google/protobuf/duration").Duration>, never>) | undefined;
            slashFractionDoubleSign?: Uint8Array | undefined;
            slashFractionDowntime?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["params"], keyof Params>, never>) | undefined;
    } & Record<Exclude<keyof I, "params">, never>>(object: I): QueryParamsResponse;
};
export declare const QuerySigningInfoRequest: {
    typeUrl: string;
    encode(message: QuerySigningInfoRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySigningInfoRequest;
    fromJSON(object: any): QuerySigningInfoRequest;
    toJSON(message: QuerySigningInfoRequest): unknown;
    fromPartial<I extends {
        consAddress?: string | undefined;
    } & {
        consAddress?: string | undefined;
    } & Record<Exclude<keyof I, "consAddress">, never>>(object: I): QuerySigningInfoRequest;
};
export declare const QuerySigningInfoResponse: {
    typeUrl: string;
    encode(message: QuerySigningInfoResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySigningInfoResponse;
    fromJSON(object: any): QuerySigningInfoResponse;
    toJSON(message: QuerySigningInfoResponse): unknown;
    fromPartial<I extends {
        valSigningInfo?: {
            address?: string | undefined;
            startHeight?: bigint | undefined;
            indexOffset?: bigint | undefined;
            jailedUntil?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            tombstoned?: boolean | undefined;
            missedBlocksCounter?: bigint | undefined;
        } | undefined;
    } & {
        valSigningInfo?: ({
            address?: string | undefined;
            startHeight?: bigint | undefined;
            indexOffset?: bigint | undefined;
            jailedUntil?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            tombstoned?: boolean | undefined;
            missedBlocksCounter?: bigint | undefined;
        } & {
            address?: string | undefined;
            startHeight?: bigint | undefined;
            indexOffset?: bigint | undefined;
            jailedUntil?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["valSigningInfo"]["jailedUntil"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            tombstoned?: boolean | undefined;
            missedBlocksCounter?: bigint | undefined;
        } & Record<Exclude<keyof I["valSigningInfo"], keyof ValidatorSigningInfo>, never>) | undefined;
    } & Record<Exclude<keyof I, "valSigningInfo">, never>>(object: I): QuerySigningInfoResponse;
};
export declare const QuerySigningInfosRequest: {
    typeUrl: string;
    encode(message: QuerySigningInfosRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySigningInfosRequest;
    fromJSON(object: any): QuerySigningInfosRequest;
    toJSON(message: QuerySigningInfosRequest): unknown;
    fromPartial<I extends {
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
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
    } & Record<Exclude<keyof I, "pagination">, never>>(object: I): QuerySigningInfosRequest;
};
export declare const QuerySigningInfosResponse: {
    typeUrl: string;
    encode(message: QuerySigningInfosResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySigningInfosResponse;
    fromJSON(object: any): QuerySigningInfosResponse;
    toJSON(message: QuerySigningInfosResponse): unknown;
    fromPartial<I extends {
        info?: {
            address?: string | undefined;
            startHeight?: bigint | undefined;
            indexOffset?: bigint | undefined;
            jailedUntil?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            tombstoned?: boolean | undefined;
            missedBlocksCounter?: bigint | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        info?: ({
            address?: string | undefined;
            startHeight?: bigint | undefined;
            indexOffset?: bigint | undefined;
            jailedUntil?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            tombstoned?: boolean | undefined;
            missedBlocksCounter?: bigint | undefined;
        }[] & ({
            address?: string | undefined;
            startHeight?: bigint | undefined;
            indexOffset?: bigint | undefined;
            jailedUntil?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            tombstoned?: boolean | undefined;
            missedBlocksCounter?: bigint | undefined;
        } & {
            address?: string | undefined;
            startHeight?: bigint | undefined;
            indexOffset?: bigint | undefined;
            jailedUntil?: ({
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } & Record<Exclude<keyof I["info"][number]["jailedUntil"], keyof import("../../../google/protobuf/timestamp").Timestamp>, never>) | undefined;
            tombstoned?: boolean | undefined;
            missedBlocksCounter?: bigint | undefined;
        } & Record<Exclude<keyof I["info"][number], keyof ValidatorSigningInfo>, never>)[] & Record<Exclude<keyof I["info"], keyof {
            address?: string | undefined;
            startHeight?: bigint | undefined;
            indexOffset?: bigint | undefined;
            jailedUntil?: {
                seconds?: bigint | undefined;
                nanos?: number | undefined;
            } | undefined;
            tombstoned?: boolean | undefined;
            missedBlocksCounter?: bigint | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QuerySigningInfosResponse>, never>>(object: I): QuerySigningInfosResponse;
};
/** Query provides defines the gRPC querier service */
export interface Query {
    /** Params queries the parameters of slashing module */
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    /** SigningInfo queries the signing info of given cons address */
    SigningInfo(request: QuerySigningInfoRequest): Promise<QuerySigningInfoResponse>;
    /** SigningInfos queries signing info of all validators */
    SigningInfos(request?: QuerySigningInfosRequest): Promise<QuerySigningInfosResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    SigningInfo(request: QuerySigningInfoRequest): Promise<QuerySigningInfoResponse>;
    SigningInfos(request?: QuerySigningInfosRequest): Promise<QuerySigningInfosResponse>;
}
