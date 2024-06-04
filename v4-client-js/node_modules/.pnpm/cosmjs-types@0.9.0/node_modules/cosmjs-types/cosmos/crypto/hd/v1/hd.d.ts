import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.crypto.hd.v1";
/** BIP44Params is used as path field in ledger item in Record. */
export interface BIP44Params {
    /** purpose is a constant set to 44' (or 0x8000002C) following the BIP43 recommendation */
    purpose: number;
    /** coin_type is a constant that improves privacy */
    coinType: number;
    /** account splits the key space into independent user identities */
    account: number;
    /**
     * change is a constant used for public derivation. Constant 0 is used for external chain and constant 1 for internal
     * chain.
     */
    change: boolean;
    /** address_index is used as child index in BIP32 derivation */
    addressIndex: number;
}
export declare const BIP44Params: {
    typeUrl: string;
    encode(message: BIP44Params, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): BIP44Params;
    fromJSON(object: any): BIP44Params;
    toJSON(message: BIP44Params): unknown;
    fromPartial<I extends {
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
    } & Record<Exclude<keyof I, keyof BIP44Params>, never>>(object: I): BIP44Params;
};
