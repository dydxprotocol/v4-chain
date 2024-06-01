import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.crypto.multisig.v1beta1";
/**
 * MultiSignature wraps the signatures from a multisig.LegacyAminoPubKey.
 * See cosmos.tx.v1betata1.ModeInfo.Multi for how to specify which signers
 * signed and with which modes.
 */
export interface MultiSignature {
    signatures: Uint8Array[];
}
/**
 * CompactBitArray is an implementation of a space efficient bit array.
 * This is used to ensure that the encoded data takes up a minimal amount of
 * space after proto encoding.
 * This is not thread safe, and is not intended for concurrent usage.
 */
export interface CompactBitArray {
    extraBitsStored: number;
    elems: Uint8Array;
}
export declare const MultiSignature: {
    typeUrl: string;
    encode(message: MultiSignature, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MultiSignature;
    fromJSON(object: any): MultiSignature;
    toJSON(message: MultiSignature): unknown;
    fromPartial<I extends {
        signatures?: Uint8Array[] | undefined;
    } & {
        signatures?: (Uint8Array[] & Uint8Array[] & Record<Exclude<keyof I["signatures"], keyof Uint8Array[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "signatures">, never>>(object: I): MultiSignature;
};
export declare const CompactBitArray: {
    typeUrl: string;
    encode(message: CompactBitArray, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CompactBitArray;
    fromJSON(object: any): CompactBitArray;
    toJSON(message: CompactBitArray): unknown;
    fromPartial<I extends {
        extraBitsStored?: number | undefined;
        elems?: Uint8Array | undefined;
    } & {
        extraBitsStored?: number | undefined;
        elems?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof CompactBitArray>, never>>(object: I): CompactBitArray;
};
