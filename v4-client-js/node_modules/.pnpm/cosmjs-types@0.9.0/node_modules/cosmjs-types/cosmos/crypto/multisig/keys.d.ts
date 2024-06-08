import { Any } from "../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.crypto.multisig";
/**
 * LegacyAminoPubKey specifies a public key type
 * which nests multiple public keys and a threshold,
 * it uses legacy amino address rules.
 */
export interface LegacyAminoPubKey {
    threshold: number;
    publicKeys: Any[];
}
export declare const LegacyAminoPubKey: {
    typeUrl: string;
    encode(message: LegacyAminoPubKey, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): LegacyAminoPubKey;
    fromJSON(object: any): LegacyAminoPubKey;
    toJSON(message: LegacyAminoPubKey): unknown;
    fromPartial<I extends {
        threshold?: number | undefined;
        publicKeys?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        threshold?: number | undefined;
        publicKeys?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] & ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["publicKeys"][number], keyof Any>, never>)[] & Record<Exclude<keyof I["publicKeys"], keyof {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof LegacyAminoPubKey>, never>>(object: I): LegacyAminoPubKey;
};
