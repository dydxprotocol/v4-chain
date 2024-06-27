import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.base.v1beta1";
/**
 * Coin defines a token with a denomination and an amount.
 *
 * NOTE: The amount field is an Int which implements the custom method
 * signatures required by gogoproto.
 */
export interface Coin {
    denom: string;
    amount: string;
}
/**
 * DecCoin defines a token with a denomination and a decimal amount.
 *
 * NOTE: The amount field is an Dec which implements the custom method
 * signatures required by gogoproto.
 */
export interface DecCoin {
    denom: string;
    amount: string;
}
/** IntProto defines a Protobuf wrapper around an Int object. */
export interface IntProto {
    int: string;
}
/** DecProto defines a Protobuf wrapper around a Dec object. */
export interface DecProto {
    dec: string;
}
export declare const Coin: {
    typeUrl: string;
    encode(message: Coin, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Coin;
    fromJSON(object: any): Coin;
    toJSON(message: Coin): unknown;
    fromPartial<I extends {
        denom?: string | undefined;
        amount?: string | undefined;
    } & {
        denom?: string | undefined;
        amount?: string | undefined;
    } & Record<Exclude<keyof I, keyof Coin>, never>>(object: I): Coin;
};
export declare const DecCoin: {
    typeUrl: string;
    encode(message: DecCoin, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DecCoin;
    fromJSON(object: any): DecCoin;
    toJSON(message: DecCoin): unknown;
    fromPartial<I extends {
        denom?: string | undefined;
        amount?: string | undefined;
    } & {
        denom?: string | undefined;
        amount?: string | undefined;
    } & Record<Exclude<keyof I, keyof DecCoin>, never>>(object: I): DecCoin;
};
export declare const IntProto: {
    typeUrl: string;
    encode(message: IntProto, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): IntProto;
    fromJSON(object: any): IntProto;
    toJSON(message: IntProto): unknown;
    fromPartial<I extends {
        int?: string | undefined;
    } & {
        int?: string | undefined;
    } & Record<Exclude<keyof I, "int">, never>>(object: I): IntProto;
};
export declare const DecProto: {
    typeUrl: string;
    encode(message: DecProto, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DecProto;
    fromJSON(object: any): DecProto;
    toJSON(message: DecProto): unknown;
    fromPartial<I extends {
        dec?: string | undefined;
    } & {
        dec?: string | undefined;
    } & Record<Exclude<keyof I, "dec">, never>>(object: I): DecProto;
};
