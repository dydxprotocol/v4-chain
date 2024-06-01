import { Any } from "../../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.applications.interchain_accounts.v1";
/**
 * Type defines a classification of message issued from a controller chain to its associated interchain accounts
 * host
 */
export declare enum Type {
    /** TYPE_UNSPECIFIED - Default zero value enumeration */
    TYPE_UNSPECIFIED = 0,
    /** TYPE_EXECUTE_TX - Execute a transaction on an interchain accounts host chain */
    TYPE_EXECUTE_TX = 1,
    UNRECOGNIZED = -1
}
export declare function typeFromJSON(object: any): Type;
export declare function typeToJSON(object: Type): string;
/** InterchainAccountPacketData is comprised of a raw transaction, type of transaction and optional memo field. */
export interface InterchainAccountPacketData {
    type: Type;
    data: Uint8Array;
    memo: string;
}
/** CosmosTx contains a list of sdk.Msg's. It should be used when sending transactions to an SDK host chain. */
export interface CosmosTx {
    messages: Any[];
}
export declare const InterchainAccountPacketData: {
    typeUrl: string;
    encode(message: InterchainAccountPacketData, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): InterchainAccountPacketData;
    fromJSON(object: any): InterchainAccountPacketData;
    toJSON(message: InterchainAccountPacketData): unknown;
    fromPartial<I extends {
        type?: Type | undefined;
        data?: Uint8Array | undefined;
        memo?: string | undefined;
    } & {
        type?: Type | undefined;
        data?: Uint8Array | undefined;
        memo?: string | undefined;
    } & Record<Exclude<keyof I, keyof InterchainAccountPacketData>, never>>(object: I): InterchainAccountPacketData;
};
export declare const CosmosTx: {
    typeUrl: string;
    encode(message: CosmosTx, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CosmosTx;
    fromJSON(object: any): CosmosTx;
    toJSON(message: CosmosTx): unknown;
    fromPartial<I extends {
        messages?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        messages?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[] & ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["messages"][number], keyof Any>, never>)[] & Record<Exclude<keyof I["messages"], keyof {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "messages">, never>>(object: I): CosmosTx;
};
