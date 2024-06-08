import { Class, NFT } from "./nft";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.nft.v1beta1";
/** GenesisState defines the nft module's genesis state. */
export interface GenesisState {
    /** class defines the class of the nft type. */
    classes: Class[];
    /** entry defines all nft owned by a person. */
    entries: Entry[];
}
/** Entry Defines all nft owned by a person */
export interface Entry {
    /** owner is the owner address of the following nft */
    owner: string;
    /** nfts is a group of nfts of the same owner */
    nfts: NFT[];
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
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
        entries?: {
            owner?: string | undefined;
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
        }[] | undefined;
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
        entries?: ({
            owner?: string | undefined;
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
        }[] & ({
            owner?: string | undefined;
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
        } & {
            owner?: string | undefined;
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
                } & Record<Exclude<keyof I["entries"][number]["nfts"][number]["data"], keyof import("../../../google/protobuf/any").Any>, never>) | undefined;
            } & Record<Exclude<keyof I["entries"][number]["nfts"][number], keyof NFT>, never>)[] & Record<Exclude<keyof I["entries"][number]["nfts"], keyof {
                classId?: string | undefined;
                id?: string | undefined;
                uri?: string | undefined;
                uriHash?: string | undefined;
                data?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["entries"][number], keyof Entry>, never>)[] & Record<Exclude<keyof I["entries"], keyof {
            owner?: string | undefined;
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
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
export declare const Entry: {
    typeUrl: string;
    encode(message: Entry, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Entry;
    fromJSON(object: any): Entry;
    toJSON(message: Entry): unknown;
    fromPartial<I extends {
        owner?: string | undefined;
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
    } & {
        owner?: string | undefined;
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
    } & Record<Exclude<keyof I, keyof Entry>, never>>(object: I): Entry;
};
