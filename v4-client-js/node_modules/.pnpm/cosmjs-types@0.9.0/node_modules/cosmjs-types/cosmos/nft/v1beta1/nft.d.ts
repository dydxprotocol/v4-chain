import { Any } from "../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.nft.v1beta1";
/** Class defines the class of the nft type. */
export interface Class {
    /** id defines the unique identifier of the NFT classification, similar to the contract address of ERC721 */
    id: string;
    /** name defines the human-readable name of the NFT classification. Optional */
    name: string;
    /** symbol is an abbreviated name for nft classification. Optional */
    symbol: string;
    /** description is a brief description of nft classification. Optional */
    description: string;
    /** uri for the class metadata stored off chain. It can define schema for Class and NFT `Data` attributes. Optional */
    uri: string;
    /** uri_hash is a hash of the document pointed by uri. Optional */
    uriHash: string;
    /** data is the app specific metadata of the NFT class. Optional */
    data?: Any;
}
/** NFT defines the NFT. */
export interface NFT {
    /** class_id associated with the NFT, similar to the contract address of ERC721 */
    classId: string;
    /** id is a unique identifier of the NFT */
    id: string;
    /** uri for the NFT metadata stored off chain */
    uri: string;
    /** uri_hash is a hash of the document pointed by uri */
    uriHash: string;
    /** data is an app specific data of the NFT. Optional */
    data?: Any;
}
export declare const Class: {
    typeUrl: string;
    encode(message: Class, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Class;
    fromJSON(object: any): Class;
    toJSON(message: Class): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["data"], keyof Any>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Class>, never>>(object: I): Class;
};
export declare const NFT: {
    typeUrl: string;
    encode(message: NFT, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): NFT;
    fromJSON(object: any): NFT;
    toJSON(message: NFT): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["data"], keyof Any>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof NFT>, never>>(object: I): NFT;
};
