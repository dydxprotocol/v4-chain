import { PageRequest, PageResponse } from "../../../cosmos/base/query/v1beta1/pagination";
import { ContractInfo, ContractCodeHistoryEntry, Model, AccessConfig, Params } from "./types";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmwasm.wasm.v1";
/**
 * QueryContractInfoRequest is the request type for the Query/ContractInfo RPC
 * method
 */
export interface QueryContractInfoRequest {
    /** address is the address of the contract to query */
    address: string;
}
/**
 * QueryContractInfoResponse is the response type for the Query/ContractInfo RPC
 * method
 */
export interface QueryContractInfoResponse {
    /** address is the address of the contract */
    address: string;
    contractInfo: ContractInfo;
}
/**
 * QueryContractHistoryRequest is the request type for the Query/ContractHistory
 * RPC method
 */
export interface QueryContractHistoryRequest {
    /** address is the address of the contract to query */
    address: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryContractHistoryResponse is the response type for the
 * Query/ContractHistory RPC method
 */
export interface QueryContractHistoryResponse {
    entries: ContractCodeHistoryEntry[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QueryContractsByCodeRequest is the request type for the Query/ContractsByCode
 * RPC method
 */
export interface QueryContractsByCodeRequest {
    /**
     * grpc-gateway_out does not support Go style CodID
     * pagination defines an optional pagination for the request.
     */
    codeId: bigint;
    pagination?: PageRequest;
}
/**
 * QueryContractsByCodeResponse is the response type for the
 * Query/ContractsByCode RPC method
 */
export interface QueryContractsByCodeResponse {
    /** contracts are a set of contract addresses */
    contracts: string[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QueryAllContractStateRequest is the request type for the
 * Query/AllContractState RPC method
 */
export interface QueryAllContractStateRequest {
    /** address is the address of the contract */
    address: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryAllContractStateResponse is the response type for the
 * Query/AllContractState RPC method
 */
export interface QueryAllContractStateResponse {
    models: Model[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QueryRawContractStateRequest is the request type for the
 * Query/RawContractState RPC method
 */
export interface QueryRawContractStateRequest {
    /** address is the address of the contract */
    address: string;
    queryData: Uint8Array;
}
/**
 * QueryRawContractStateResponse is the response type for the
 * Query/RawContractState RPC method
 */
export interface QueryRawContractStateResponse {
    /** Data contains the raw store data */
    data: Uint8Array;
}
/**
 * QuerySmartContractStateRequest is the request type for the
 * Query/SmartContractState RPC method
 */
export interface QuerySmartContractStateRequest {
    /** address is the address of the contract */
    address: string;
    /** QueryData contains the query data passed to the contract */
    queryData: Uint8Array;
}
/**
 * QuerySmartContractStateResponse is the response type for the
 * Query/SmartContractState RPC method
 */
export interface QuerySmartContractStateResponse {
    /** Data contains the json data returned from the smart contract */
    data: Uint8Array;
}
/** QueryCodeRequest is the request type for the Query/Code RPC method */
export interface QueryCodeRequest {
    /** grpc-gateway_out does not support Go style CodID */
    codeId: bigint;
}
/** CodeInfoResponse contains code meta data from CodeInfo */
export interface CodeInfoResponse {
    codeId: bigint;
    creator: string;
    dataHash: Uint8Array;
    instantiatePermission: AccessConfig;
}
/** QueryCodeResponse is the response type for the Query/Code RPC method */
export interface QueryCodeResponse {
    codeInfo?: CodeInfoResponse;
    data: Uint8Array;
}
/** QueryCodesRequest is the request type for the Query/Codes RPC method */
export interface QueryCodesRequest {
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryCodesResponse is the response type for the Query/Codes RPC method */
export interface QueryCodesResponse {
    codeInfos: CodeInfoResponse[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/**
 * QueryPinnedCodesRequest is the request type for the Query/PinnedCodes
 * RPC method
 */
export interface QueryPinnedCodesRequest {
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryPinnedCodesResponse is the response type for the
 * Query/PinnedCodes RPC method
 */
export interface QueryPinnedCodesResponse {
    codeIds: bigint[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryParamsRequest is the request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}
/** QueryParamsResponse is the response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
    /** params defines the parameters of the module. */
    params: Params;
}
/**
 * QueryContractsByCreatorRequest is the request type for the
 * Query/ContractsByCreator RPC method.
 */
export interface QueryContractsByCreatorRequest {
    /** CreatorAddress is the address of contract creator */
    creatorAddress: string;
    /** Pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/**
 * QueryContractsByCreatorResponse is the response type for the
 * Query/ContractsByCreator RPC method.
 */
export interface QueryContractsByCreatorResponse {
    /** ContractAddresses result set */
    contractAddresses: string[];
    /** Pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
export declare const QueryContractInfoRequest: {
    typeUrl: string;
    encode(message: QueryContractInfoRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryContractInfoRequest;
    fromJSON(object: any): QueryContractInfoRequest;
    toJSON(message: QueryContractInfoRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
    } & {
        address?: string | undefined;
    } & Record<Exclude<keyof I, "address">, never>>(object: I): QueryContractInfoRequest;
};
export declare const QueryContractInfoResponse: {
    typeUrl: string;
    encode(message: QueryContractInfoResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryContractInfoResponse;
    fromJSON(object: any): QueryContractInfoResponse;
    toJSON(message: QueryContractInfoResponse): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        contractInfo?: {
            codeId?: bigint | undefined;
            creator?: string | undefined;
            admin?: string | undefined;
            label?: string | undefined;
            created?: {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } | undefined;
            ibcPortId?: string | undefined;
            extension?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
        contractInfo?: ({
            codeId?: bigint | undefined;
            creator?: string | undefined;
            admin?: string | undefined;
            label?: string | undefined;
            created?: {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } | undefined;
            ibcPortId?: string | undefined;
            extension?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            codeId?: bigint | undefined;
            creator?: string | undefined;
            admin?: string | undefined;
            label?: string | undefined;
            created?: ({
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } & {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } & Record<Exclude<keyof I["contractInfo"]["created"], keyof import("./types").AbsoluteTxPosition>, never>) | undefined;
            ibcPortId?: string | undefined;
            extension?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["contractInfo"]["extension"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["contractInfo"], keyof ContractInfo>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryContractInfoResponse>, never>>(object: I): QueryContractInfoResponse;
};
export declare const QueryContractHistoryRequest: {
    typeUrl: string;
    encode(message: QueryContractHistoryRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryContractHistoryRequest;
    fromJSON(object: any): QueryContractHistoryRequest;
    toJSON(message: QueryContractHistoryRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryContractHistoryRequest>, never>>(object: I): QueryContractHistoryRequest;
};
export declare const QueryContractHistoryResponse: {
    typeUrl: string;
    encode(message: QueryContractHistoryResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryContractHistoryResponse;
    fromJSON(object: any): QueryContractHistoryResponse;
    toJSON(message: QueryContractHistoryResponse): unknown;
    fromPartial<I extends {
        entries?: {
            operation?: import("./types").ContractCodeHistoryOperationType | undefined;
            codeId?: bigint | undefined;
            updated?: {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } | undefined;
            msg?: Uint8Array | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        entries?: ({
            operation?: import("./types").ContractCodeHistoryOperationType | undefined;
            codeId?: bigint | undefined;
            updated?: {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } | undefined;
            msg?: Uint8Array | undefined;
        }[] & ({
            operation?: import("./types").ContractCodeHistoryOperationType | undefined;
            codeId?: bigint | undefined;
            updated?: {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } | undefined;
            msg?: Uint8Array | undefined;
        } & {
            operation?: import("./types").ContractCodeHistoryOperationType | undefined;
            codeId?: bigint | undefined;
            updated?: ({
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } & {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } & Record<Exclude<keyof I["entries"][number]["updated"], keyof import("./types").AbsoluteTxPosition>, never>) | undefined;
            msg?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["entries"][number], keyof ContractCodeHistoryEntry>, never>)[] & Record<Exclude<keyof I["entries"], keyof {
            operation?: import("./types").ContractCodeHistoryOperationType | undefined;
            codeId?: bigint | undefined;
            updated?: {
                blockHeight?: bigint | undefined;
                txIndex?: bigint | undefined;
            } | undefined;
            msg?: Uint8Array | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryContractHistoryResponse>, never>>(object: I): QueryContractHistoryResponse;
};
export declare const QueryContractsByCodeRequest: {
    typeUrl: string;
    encode(message: QueryContractsByCodeRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryContractsByCodeRequest;
    fromJSON(object: any): QueryContractsByCodeRequest;
    toJSON(message: QueryContractsByCodeRequest): unknown;
    fromPartial<I extends {
        codeId?: bigint | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        codeId?: bigint | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryContractsByCodeRequest>, never>>(object: I): QueryContractsByCodeRequest;
};
export declare const QueryContractsByCodeResponse: {
    typeUrl: string;
    encode(message: QueryContractsByCodeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryContractsByCodeResponse;
    fromJSON(object: any): QueryContractsByCodeResponse;
    toJSON(message: QueryContractsByCodeResponse): unknown;
    fromPartial<I extends {
        contracts?: string[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        contracts?: (string[] & string[] & Record<Exclude<keyof I["contracts"], keyof string[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryContractsByCodeResponse>, never>>(object: I): QueryContractsByCodeResponse;
};
export declare const QueryAllContractStateRequest: {
    typeUrl: string;
    encode(message: QueryAllContractStateRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAllContractStateRequest;
    fromJSON(object: any): QueryAllContractStateRequest;
    toJSON(message: QueryAllContractStateRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        address?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryAllContractStateRequest>, never>>(object: I): QueryAllContractStateRequest;
};
export declare const QueryAllContractStateResponse: {
    typeUrl: string;
    encode(message: QueryAllContractStateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryAllContractStateResponse;
    fromJSON(object: any): QueryAllContractStateResponse;
    toJSON(message: QueryAllContractStateResponse): unknown;
    fromPartial<I extends {
        models?: {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        models?: ({
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        }[] & ({
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        } & {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["models"][number], keyof Model>, never>)[] & Record<Exclude<keyof I["models"], keyof {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryAllContractStateResponse>, never>>(object: I): QueryAllContractStateResponse;
};
export declare const QueryRawContractStateRequest: {
    typeUrl: string;
    encode(message: QueryRawContractStateRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryRawContractStateRequest;
    fromJSON(object: any): QueryRawContractStateRequest;
    toJSON(message: QueryRawContractStateRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        queryData?: Uint8Array | undefined;
    } & {
        address?: string | undefined;
        queryData?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof QueryRawContractStateRequest>, never>>(object: I): QueryRawContractStateRequest;
};
export declare const QueryRawContractStateResponse: {
    typeUrl: string;
    encode(message: QueryRawContractStateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryRawContractStateResponse;
    fromJSON(object: any): QueryRawContractStateResponse;
    toJSON(message: QueryRawContractStateResponse): unknown;
    fromPartial<I extends {
        data?: Uint8Array | undefined;
    } & {
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "data">, never>>(object: I): QueryRawContractStateResponse;
};
export declare const QuerySmartContractStateRequest: {
    typeUrl: string;
    encode(message: QuerySmartContractStateRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySmartContractStateRequest;
    fromJSON(object: any): QuerySmartContractStateRequest;
    toJSON(message: QuerySmartContractStateRequest): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        queryData?: Uint8Array | undefined;
    } & {
        address?: string | undefined;
        queryData?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof QuerySmartContractStateRequest>, never>>(object: I): QuerySmartContractStateRequest;
};
export declare const QuerySmartContractStateResponse: {
    typeUrl: string;
    encode(message: QuerySmartContractStateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySmartContractStateResponse;
    fromJSON(object: any): QuerySmartContractStateResponse;
    toJSON(message: QuerySmartContractStateResponse): unknown;
    fromPartial<I extends {
        data?: Uint8Array | undefined;
    } & {
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, "data">, never>>(object: I): QuerySmartContractStateResponse;
};
export declare const QueryCodeRequest: {
    typeUrl: string;
    encode(message: QueryCodeRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryCodeRequest;
    fromJSON(object: any): QueryCodeRequest;
    toJSON(message: QueryCodeRequest): unknown;
    fromPartial<I extends {
        codeId?: bigint | undefined;
    } & {
        codeId?: bigint | undefined;
    } & Record<Exclude<keyof I, "codeId">, never>>(object: I): QueryCodeRequest;
};
export declare const CodeInfoResponse: {
    typeUrl: string;
    encode(message: CodeInfoResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CodeInfoResponse;
    fromJSON(object: any): CodeInfoResponse;
    toJSON(message: CodeInfoResponse): unknown;
    fromPartial<I extends {
        codeId?: bigint | undefined;
        creator?: string | undefined;
        dataHash?: Uint8Array | undefined;
        instantiatePermission?: {
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } | undefined;
    } & {
        codeId?: bigint | undefined;
        creator?: string | undefined;
        dataHash?: Uint8Array | undefined;
        instantiatePermission?: ({
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: string[] | undefined;
        } & {
            permission?: import("./types").AccessType | undefined;
            address?: string | undefined;
            addresses?: (string[] & string[] & Record<Exclude<keyof I["instantiatePermission"]["addresses"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["instantiatePermission"], keyof AccessConfig>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof CodeInfoResponse>, never>>(object: I): CodeInfoResponse;
};
export declare const QueryCodeResponse: {
    typeUrl: string;
    encode(message: QueryCodeResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryCodeResponse;
    fromJSON(object: any): QueryCodeResponse;
    toJSON(message: QueryCodeResponse): unknown;
    fromPartial<I extends {
        codeInfo?: {
            codeId?: bigint | undefined;
            creator?: string | undefined;
            dataHash?: Uint8Array | undefined;
            instantiatePermission?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        } | undefined;
        data?: Uint8Array | undefined;
    } & {
        codeInfo?: ({
            codeId?: bigint | undefined;
            creator?: string | undefined;
            dataHash?: Uint8Array | undefined;
            instantiatePermission?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        } & {
            codeId?: bigint | undefined;
            creator?: string | undefined;
            dataHash?: Uint8Array | undefined;
            instantiatePermission?: ({
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } & {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: (string[] & string[] & Record<Exclude<keyof I["codeInfo"]["instantiatePermission"]["addresses"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["codeInfo"]["instantiatePermission"], keyof AccessConfig>, never>) | undefined;
        } & Record<Exclude<keyof I["codeInfo"], keyof CodeInfoResponse>, never>) | undefined;
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof QueryCodeResponse>, never>>(object: I): QueryCodeResponse;
};
export declare const QueryCodesRequest: {
    typeUrl: string;
    encode(message: QueryCodesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryCodesRequest;
    fromJSON(object: any): QueryCodesRequest;
    toJSON(message: QueryCodesRequest): unknown;
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
    } & Record<Exclude<keyof I, "pagination">, never>>(object: I): QueryCodesRequest;
};
export declare const QueryCodesResponse: {
    typeUrl: string;
    encode(message: QueryCodesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryCodesResponse;
    fromJSON(object: any): QueryCodesResponse;
    toJSON(message: QueryCodesResponse): unknown;
    fromPartial<I extends {
        codeInfos?: {
            codeId?: bigint | undefined;
            creator?: string | undefined;
            dataHash?: Uint8Array | undefined;
            instantiatePermission?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        codeInfos?: ({
            codeId?: bigint | undefined;
            creator?: string | undefined;
            dataHash?: Uint8Array | undefined;
            instantiatePermission?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        }[] & ({
            codeId?: bigint | undefined;
            creator?: string | undefined;
            dataHash?: Uint8Array | undefined;
            instantiatePermission?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        } & {
            codeId?: bigint | undefined;
            creator?: string | undefined;
            dataHash?: Uint8Array | undefined;
            instantiatePermission?: ({
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } & {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: (string[] & string[] & Record<Exclude<keyof I["codeInfos"][number]["instantiatePermission"]["addresses"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["codeInfos"][number]["instantiatePermission"], keyof AccessConfig>, never>) | undefined;
        } & Record<Exclude<keyof I["codeInfos"][number], keyof CodeInfoResponse>, never>)[] & Record<Exclude<keyof I["codeInfos"], keyof {
            codeId?: bigint | undefined;
            creator?: string | undefined;
            dataHash?: Uint8Array | undefined;
            instantiatePermission?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryCodesResponse>, never>>(object: I): QueryCodesResponse;
};
export declare const QueryPinnedCodesRequest: {
    typeUrl: string;
    encode(message: QueryPinnedCodesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPinnedCodesRequest;
    fromJSON(object: any): QueryPinnedCodesRequest;
    toJSON(message: QueryPinnedCodesRequest): unknown;
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
    } & Record<Exclude<keyof I, "pagination">, never>>(object: I): QueryPinnedCodesRequest;
};
export declare const QueryPinnedCodesResponse: {
    typeUrl: string;
    encode(message: QueryPinnedCodesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPinnedCodesResponse;
    fromJSON(object: any): QueryPinnedCodesResponse;
    toJSON(message: QueryPinnedCodesResponse): unknown;
    fromPartial<I extends {
        codeIds?: bigint[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        codeIds?: (bigint[] & bigint[] & Record<Exclude<keyof I["codeIds"], keyof bigint[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryPinnedCodesResponse>, never>>(object: I): QueryPinnedCodesResponse;
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
            codeUploadAccess?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
            instantiateDefaultPermission?: import("./types").AccessType | undefined;
        } | undefined;
    } & {
        params?: ({
            codeUploadAccess?: {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } | undefined;
            instantiateDefaultPermission?: import("./types").AccessType | undefined;
        } & {
            codeUploadAccess?: ({
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: string[] | undefined;
            } & {
                permission?: import("./types").AccessType | undefined;
                address?: string | undefined;
                addresses?: (string[] & string[] & Record<Exclude<keyof I["params"]["codeUploadAccess"]["addresses"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["params"]["codeUploadAccess"], keyof AccessConfig>, never>) | undefined;
            instantiateDefaultPermission?: import("./types").AccessType | undefined;
        } & Record<Exclude<keyof I["params"], keyof Params>, never>) | undefined;
    } & Record<Exclude<keyof I, "params">, never>>(object: I): QueryParamsResponse;
};
export declare const QueryContractsByCreatorRequest: {
    typeUrl: string;
    encode(message: QueryContractsByCreatorRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryContractsByCreatorRequest;
    fromJSON(object: any): QueryContractsByCreatorRequest;
    toJSON(message: QueryContractsByCreatorRequest): unknown;
    fromPartial<I extends {
        creatorAddress?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        creatorAddress?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryContractsByCreatorRequest>, never>>(object: I): QueryContractsByCreatorRequest;
};
export declare const QueryContractsByCreatorResponse: {
    typeUrl: string;
    encode(message: QueryContractsByCreatorResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryContractsByCreatorResponse;
    fromJSON(object: any): QueryContractsByCreatorResponse;
    toJSON(message: QueryContractsByCreatorResponse): unknown;
    fromPartial<I extends {
        contractAddresses?: string[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        contractAddresses?: (string[] & string[] & Record<Exclude<keyof I["contractAddresses"], keyof string[]>, never>) | undefined;
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryContractsByCreatorResponse>, never>>(object: I): QueryContractsByCreatorResponse;
};
/** Query provides defines the gRPC querier service */
export interface Query {
    /** ContractInfo gets the contract meta data */
    ContractInfo(request: QueryContractInfoRequest): Promise<QueryContractInfoResponse>;
    /** ContractHistory gets the contract code history */
    ContractHistory(request: QueryContractHistoryRequest): Promise<QueryContractHistoryResponse>;
    /** ContractsByCode lists all smart contracts for a code id */
    ContractsByCode(request: QueryContractsByCodeRequest): Promise<QueryContractsByCodeResponse>;
    /** AllContractState gets all raw store data for a single contract */
    AllContractState(request: QueryAllContractStateRequest): Promise<QueryAllContractStateResponse>;
    /** RawContractState gets single key from the raw store data of a contract */
    RawContractState(request: QueryRawContractStateRequest): Promise<QueryRawContractStateResponse>;
    /** SmartContractState get smart query result from the contract */
    SmartContractState(request: QuerySmartContractStateRequest): Promise<QuerySmartContractStateResponse>;
    /** Code gets the binary code and metadata for a singe wasm code */
    Code(request: QueryCodeRequest): Promise<QueryCodeResponse>;
    /** Codes gets the metadata for all stored wasm codes */
    Codes(request?: QueryCodesRequest): Promise<QueryCodesResponse>;
    /** PinnedCodes gets the pinned code ids */
    PinnedCodes(request?: QueryPinnedCodesRequest): Promise<QueryPinnedCodesResponse>;
    /** Params gets the module params */
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    /** ContractsByCreator gets the contracts by creator */
    ContractsByCreator(request: QueryContractsByCreatorRequest): Promise<QueryContractsByCreatorResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    ContractInfo(request: QueryContractInfoRequest): Promise<QueryContractInfoResponse>;
    ContractHistory(request: QueryContractHistoryRequest): Promise<QueryContractHistoryResponse>;
    ContractsByCode(request: QueryContractsByCodeRequest): Promise<QueryContractsByCodeResponse>;
    AllContractState(request: QueryAllContractStateRequest): Promise<QueryAllContractStateResponse>;
    RawContractState(request: QueryRawContractStateRequest): Promise<QueryRawContractStateResponse>;
    SmartContractState(request: QuerySmartContractStateRequest): Promise<QuerySmartContractStateResponse>;
    Code(request: QueryCodeRequest): Promise<QueryCodeResponse>;
    Codes(request?: QueryCodesRequest): Promise<QueryCodesResponse>;
    PinnedCodes(request?: QueryPinnedCodesRequest): Promise<QueryPinnedCodesResponse>;
    Params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    ContractsByCreator(request: QueryContractsByCreatorRequest): Promise<QueryContractsByCreatorResponse>;
}
