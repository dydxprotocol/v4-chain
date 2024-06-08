import { PageRequest, PageResponse } from "../../base/query/v1beta1/pagination";
import { NFT, Class } from "./nft";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.nft.v1beta1";
/** QueryBalanceRequest is the request type for the Query/Balance RPC method */
export interface QueryBalanceRequest {
    /** class_id associated with the nft */
    classId: string;
    /** owner is the owner address of the nft */
    owner: string;
}
/** QueryBalanceResponse is the response type for the Query/Balance RPC method */
export interface QueryBalanceResponse {
    /** amount is the number of all NFTs of a given class owned by the owner */
    amount: bigint;
}
/** QueryOwnerRequest is the request type for the Query/Owner RPC method */
export interface QueryOwnerRequest {
    /** class_id associated with the nft */
    classId: string;
    /** id is a unique identifier of the NFT */
    id: string;
}
/** QueryOwnerResponse is the response type for the Query/Owner RPC method */
export interface QueryOwnerResponse {
    /** owner is the owner address of the nft */
    owner: string;
}
/** QuerySupplyRequest is the request type for the Query/Supply RPC method */
export interface QuerySupplyRequest {
    /** class_id associated with the nft */
    classId: string;
}
/** QuerySupplyResponse is the response type for the Query/Supply RPC method */
export interface QuerySupplyResponse {
    /** amount is the number of all NFTs from the given class */
    amount: bigint;
}
/** QueryNFTstRequest is the request type for the Query/NFTs RPC method */
export interface QueryNFTsRequest {
    /** class_id associated with the nft */
    classId: string;
    /** owner is the owner address of the nft */
    owner: string;
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryNFTsResponse is the response type for the Query/NFTs RPC methods */
export interface QueryNFTsResponse {
    /** NFT defines the NFT */
    nfts: NFT[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
/** QueryNFTRequest is the request type for the Query/NFT RPC method */
export interface QueryNFTRequest {
    /** class_id associated with the nft */
    classId: string;
    /** id is a unique identifier of the NFT */
    id: string;
}
/** QueryNFTResponse is the response type for the Query/NFT RPC method */
export interface QueryNFTResponse {
    /** owner is the owner address of the nft */
    nft?: NFT;
}
/** QueryClassRequest is the request type for the Query/Class RPC method */
export interface QueryClassRequest {
    /** class_id associated with the nft */
    classId: string;
}
/** QueryClassResponse is the response type for the Query/Class RPC method */
export interface QueryClassResponse {
    /** class defines the class of the nft type. */
    class?: Class;
}
/** QueryClassesRequest is the request type for the Query/Classes RPC method */
export interface QueryClassesRequest {
    /** pagination defines an optional pagination for the request. */
    pagination?: PageRequest;
}
/** QueryClassesResponse is the response type for the Query/Classes RPC method */
export interface QueryClassesResponse {
    /** class defines the class of the nft type. */
    classes: Class[];
    /** pagination defines the pagination in the response. */
    pagination?: PageResponse;
}
export declare const QueryBalanceRequest: {
    typeUrl: string;
    encode(message: QueryBalanceRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryBalanceRequest;
    fromJSON(object: any): QueryBalanceRequest;
    toJSON(message: QueryBalanceRequest): unknown;
    fromPartial<I extends {
        classId?: string | undefined;
        owner?: string | undefined;
    } & {
        classId?: string | undefined;
        owner?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryBalanceRequest>, never>>(object: I): QueryBalanceRequest;
};
export declare const QueryBalanceResponse: {
    typeUrl: string;
    encode(message: QueryBalanceResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryBalanceResponse;
    fromJSON(object: any): QueryBalanceResponse;
    toJSON(message: QueryBalanceResponse): unknown;
    fromPartial<I extends {
        amount?: bigint | undefined;
    } & {
        amount?: bigint | undefined;
    } & Record<Exclude<keyof I, "amount">, never>>(object: I): QueryBalanceResponse;
};
export declare const QueryOwnerRequest: {
    typeUrl: string;
    encode(message: QueryOwnerRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryOwnerRequest;
    fromJSON(object: any): QueryOwnerRequest;
    toJSON(message: QueryOwnerRequest): unknown;
    fromPartial<I extends {
        classId?: string | undefined;
        id?: string | undefined;
    } & {
        classId?: string | undefined;
        id?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryOwnerRequest>, never>>(object: I): QueryOwnerRequest;
};
export declare const QueryOwnerResponse: {
    typeUrl: string;
    encode(message: QueryOwnerResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryOwnerResponse;
    fromJSON(object: any): QueryOwnerResponse;
    toJSON(message: QueryOwnerResponse): unknown;
    fromPartial<I extends {
        owner?: string | undefined;
    } & {
        owner?: string | undefined;
    } & Record<Exclude<keyof I, "owner">, never>>(object: I): QueryOwnerResponse;
};
export declare const QuerySupplyRequest: {
    typeUrl: string;
    encode(message: QuerySupplyRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySupplyRequest;
    fromJSON(object: any): QuerySupplyRequest;
    toJSON(message: QuerySupplyRequest): unknown;
    fromPartial<I extends {
        classId?: string | undefined;
    } & {
        classId?: string | undefined;
    } & Record<Exclude<keyof I, "classId">, never>>(object: I): QuerySupplyRequest;
};
export declare const QuerySupplyResponse: {
    typeUrl: string;
    encode(message: QuerySupplyResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QuerySupplyResponse;
    fromJSON(object: any): QuerySupplyResponse;
    toJSON(message: QuerySupplyResponse): unknown;
    fromPartial<I extends {
        amount?: bigint | undefined;
    } & {
        amount?: bigint | undefined;
    } & Record<Exclude<keyof I, "amount">, never>>(object: I): QuerySupplyResponse;
};
export declare const QueryNFTsRequest: {
    typeUrl: string;
    encode(message: QueryNFTsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryNFTsRequest;
    fromJSON(object: any): QueryNFTsRequest;
    toJSON(message: QueryNFTsRequest): unknown;
    fromPartial<I extends {
        classId?: string | undefined;
        owner?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        classId?: string | undefined;
        owner?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryNFTsRequest>, never>>(object: I): QueryNFTsRequest;
};
export declare const QueryNFTsResponse: {
    typeUrl: string;
    encode(message: QueryNFTsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryNFTsResponse;
    fromJSON(object: any): QueryNFTsResponse;
    toJSON(message: QueryNFTsResponse): unknown;
    fromPartial<I extends {
        nfts?: {
            classId?: string | undefined;
            id?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        nfts?: ({
            classId?: string | undefined;
            id?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] & ({
            classId?: string | undefined;
            id?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            classId?: string | undefined;
            id?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["nfts"][number]["data"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["nfts"][number], keyof NFT>, never>)[] & Record<Exclude<keyof I["nfts"], keyof {
            classId?: string | undefined;
            id?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
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
    } & Record<Exclude<keyof I, keyof QueryNFTsResponse>, never>>(object: I): QueryNFTsResponse;
};
export declare const QueryNFTRequest: {
    typeUrl: string;
    encode(message: QueryNFTRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryNFTRequest;
    fromJSON(object: any): QueryNFTRequest;
    toJSON(message: QueryNFTRequest): unknown;
    fromPartial<I extends {
        classId?: string | undefined;
        id?: string | undefined;
    } & {
        classId?: string | undefined;
        id?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryNFTRequest>, never>>(object: I): QueryNFTRequest;
};
export declare const QueryNFTResponse: {
    typeUrl: string;
    encode(message: QueryNFTResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryNFTResponse;
    fromJSON(object: any): QueryNFTResponse;
    toJSON(message: QueryNFTResponse): unknown;
    fromPartial<I extends {
        nft?: {
            classId?: string | undefined;
            id?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
    } & {
        nft?: ({
            classId?: string | undefined;
            id?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            classId?: string | undefined;
            id?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["nft"]["data"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["nft"], keyof NFT>, never>) | undefined;
    } & Record<Exclude<keyof I, "nft">, never>>(object: I): QueryNFTResponse;
};
export declare const QueryClassRequest: {
    typeUrl: string;
    encode(message: QueryClassRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryClassRequest;
    fromJSON(object: any): QueryClassRequest;
    toJSON(message: QueryClassRequest): unknown;
    fromPartial<I extends {
        classId?: string | undefined;
    } & {
        classId?: string | undefined;
    } & Record<Exclude<keyof I, "classId">, never>>(object: I): QueryClassRequest;
};
export declare const QueryClassResponse: {
    typeUrl: string;
    encode(message: QueryClassResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryClassResponse;
    fromJSON(object: any): QueryClassResponse;
    toJSON(message: QueryClassResponse): unknown;
    fromPartial<I extends {
        class?: {
            id?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            description?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
    } & {
        class?: ({
            id?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            description?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            id?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            description?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["class"]["data"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["class"], keyof Class>, never>) | undefined;
    } & Record<Exclude<keyof I, "class">, never>>(object: I): QueryClassResponse;
};
export declare const QueryClassesRequest: {
    typeUrl: string;
    encode(message: QueryClassesRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryClassesRequest;
    fromJSON(object: any): QueryClassesRequest;
    toJSON(message: QueryClassesRequest): unknown;
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
    } & Record<Exclude<keyof I, "pagination">, never>>(object: I): QueryClassesRequest;
};
export declare const QueryClassesResponse: {
    typeUrl: string;
    encode(message: QueryClassesResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryClassesResponse;
    fromJSON(object: any): QueryClassesResponse;
    toJSON(message: QueryClassesResponse): unknown;
    fromPartial<I extends {
        classes?: {
            id?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            description?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
    } & {
        classes?: ({
            id?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            description?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] & ({
            id?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            description?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            id?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            description?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["classes"][number]["data"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["classes"][number], keyof Class>, never>)[] & Record<Exclude<keyof I["classes"], keyof {
            id?: string | undefined;
            name?: string | undefined;
            symbol?: string | undefined;
            description?: string | undefined;
            uri?: string | undefined;
            uriHash?: string | undefined;
            data?: {
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
    } & Record<Exclude<keyof I, keyof QueryClassesResponse>, never>>(object: I): QueryClassesResponse;
};
/** Query defines the gRPC querier service. */
export interface Query {
    /** Balance queries the number of NFTs of a given class owned by the owner, same as balanceOf in ERC721 */
    Balance(request: QueryBalanceRequest): Promise<QueryBalanceResponse>;
    /** Owner queries the owner of the NFT based on its class and id, same as ownerOf in ERC721 */
    Owner(request: QueryOwnerRequest): Promise<QueryOwnerResponse>;
    /** Supply queries the number of NFTs from the given class, same as totalSupply of ERC721. */
    Supply(request: QuerySupplyRequest): Promise<QuerySupplyResponse>;
    /**
     * NFTs queries all NFTs of a given class or owner,choose at least one of the two, similar to tokenByIndex in
     * ERC721Enumerable
     */
    NFTs(request: QueryNFTsRequest): Promise<QueryNFTsResponse>;
    /** NFT queries an NFT based on its class and id. */
    NFT(request: QueryNFTRequest): Promise<QueryNFTResponse>;
    /** Class queries an NFT class based on its id */
    Class(request: QueryClassRequest): Promise<QueryClassResponse>;
    /** Classes queries all NFT classes */
    Classes(request?: QueryClassesRequest): Promise<QueryClassesResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Balance(request: QueryBalanceRequest): Promise<QueryBalanceResponse>;
    Owner(request: QueryOwnerRequest): Promise<QueryOwnerResponse>;
    Supply(request: QuerySupplyRequest): Promise<QuerySupplyResponse>;
    NFTs(request: QueryNFTsRequest): Promise<QueryNFTsResponse>;
    NFT(request: QueryNFTRequest): Promise<QueryNFTResponse>;
    Class(request: QueryClassRequest): Promise<QueryClassResponse>;
    Classes(request?: QueryClassesRequest): Promise<QueryClassesResponse>;
}
