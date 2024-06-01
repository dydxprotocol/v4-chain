import { Any } from "../../../../google/protobuf/any";
import { BIP44Params } from "../../hd/v1/hd";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.crypto.keyring.v1";
/** Record is used for representing a key in the keyring. */
export interface Record {
    /** name represents a name of Record */
    name: string;
    /** pub_key represents a public key in any format */
    pubKey?: Any;
    /** local stores the private key locally. */
    local?: Record_Local;
    /** ledger stores the information about a Ledger key. */
    ledger?: Record_Ledger;
    /** Multi does not store any other information. */
    multi?: Record_Multi;
    /** Offline does not store any other information. */
    offline?: Record_Offline;
}
/**
 * Item is a keyring item stored in a keyring backend.
 * Local item
 */
export interface Record_Local {
    privKey?: Any;
}
/** Ledger item */
export interface Record_Ledger {
    path?: BIP44Params;
}
/** Multi item */
export interface Record_Multi {
}
/** Offline item */
export interface Record_Offline {
}
export declare const Record: {
    typeUrl: string;
    encode(message: Record, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Record;
    fromJSON(object: any): Record;
    toJSON(message: Record): unknown;
    fromPartial<I extends {
        name?: string | undefined;
        pubKey?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        local?: {
            privKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        ledger?: {
            path?: {
                purpose?: number | undefined;
                coinType?: number | undefined;
                account?: number | undefined;
                change?: boolean | undefined;
                addressIndex?: number | undefined;
            } | undefined;
        } | undefined;
        multi?: {} | undefined;
        offline?: {} | undefined;
    } & {
        name?: string | undefined;
        pubKey?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & globalThis.Record<Exclude<keyof I["pubKey"], keyof Any>, never>) | undefined;
        local?: ({
            privKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            privKey?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & globalThis.Record<Exclude<keyof I["local"]["privKey"], keyof Any>, never>) | undefined;
        } & globalThis.Record<Exclude<keyof I["local"], "privKey">, never>) | undefined;
        ledger?: ({
            path?: {
                purpose?: number | undefined;
                coinType?: number | undefined;
                account?: number | undefined;
                change?: boolean | undefined;
                addressIndex?: number | undefined;
            } | undefined;
        } & {
            path?: ({
                purpose?: number | undefined;
                coinType?: number | undefined;
                account?: number | undefined;
                change?: boolean | undefined;
                addressIndex?: number | undefined;
            } & {
                purpose?: number | undefined;
                coinType?: number | undefined;
                account?: number | undefined;
                change?: boolean | undefined;
                addressIndex?: number | undefined;
            } & globalThis.Record<Exclude<keyof I["ledger"]["path"], keyof BIP44Params>, never>) | undefined;
        } & globalThis.Record<Exclude<keyof I["ledger"], "path">, never>) | undefined;
        multi?: ({} & {} & globalThis.Record<Exclude<keyof I["multi"], never>, never>) | undefined;
        offline?: ({} & {} & globalThis.Record<Exclude<keyof I["offline"], never>, never>) | undefined;
    } & globalThis.Record<Exclude<keyof I, keyof Record>, never>>(object: I): Record;
};
export declare const Record_Local: {
    typeUrl: string;
    encode(message: Record_Local, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Record_Local;
    fromJSON(object: any): Record_Local;
    toJSON(message: Record_Local): unknown;
    fromPartial<I extends {
        privKey?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
    } & {
        privKey?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & globalThis.Record<Exclude<keyof I["privKey"], keyof Any>, never>) | undefined;
    } & globalThis.Record<Exclude<keyof I, "privKey">, never>>(object: I): Record_Local;
};
export declare const Record_Ledger: {
    typeUrl: string;
    encode(message: Record_Ledger, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Record_Ledger;
    fromJSON(object: any): Record_Ledger;
    toJSON(message: Record_Ledger): unknown;
    fromPartial<I extends {
        path?: {
            purpose?: number | undefined;
            coinType?: number | undefined;
            account?: number | undefined;
            change?: boolean | undefined;
            addressIndex?: number | undefined;
        } | undefined;
    } & {
        path?: ({
            purpose?: number | undefined;
            coinType?: number | undefined;
            account?: number | undefined;
            change?: boolean | undefined;
            addressIndex?: number | undefined;
        } & {
            purpose?: number | undefined;
            coinType?: number | undefined;
            account?: number | undefined;
            change?: boolean | undefined;
            addressIndex?: number | undefined;
        } & globalThis.Record<Exclude<keyof I["path"], keyof BIP44Params>, never>) | undefined;
    } & globalThis.Record<Exclude<keyof I, "path">, never>>(object: I): Record_Ledger;
};
export declare const Record_Multi: {
    typeUrl: string;
    encode(_: Record_Multi, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Record_Multi;
    fromJSON(_: any): Record_Multi;
    toJSON(_: Record_Multi): unknown;
    fromPartial<I extends {} & {} & globalThis.Record<Exclude<keyof I, never>, never>>(_: I): Record_Multi;
};
export declare const Record_Offline: {
    typeUrl: string;
    encode(_: Record_Offline, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Record_Offline;
    fromJSON(_: any): Record_Offline;
    toJSON(_: Record_Offline): unknown;
    fromPartial<I extends {} & {} & globalThis.Record<Exclude<keyof I, never>, never>>(_: I): Record_Offline;
};
