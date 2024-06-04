import { BinaryReader, BinaryWriter } from "../../binary";
export declare const protobufPackage = "tendermint.crypto";
/** PublicKey defines the keys available for use with Validators */
export interface PublicKey {
    ed25519?: Uint8Array;
    secp256k1?: Uint8Array;
}
export declare const PublicKey: {
    typeUrl: string;
    encode(message: PublicKey, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PublicKey;
    fromJSON(object: any): PublicKey;
    toJSON(message: PublicKey): unknown;
    fromPartial<I extends {
        ed25519?: Uint8Array | undefined;
        secp256k1?: Uint8Array | undefined;
    } & {
        ed25519?: Uint8Array | undefined;
        secp256k1?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof PublicKey>, never>>(object: I): PublicKey;
};
