import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.applications.transfer.v1";
/**
 * DenomTrace contains the base denomination for ICS20 fungible tokens and the
 * source tracing information path.
 */
export interface DenomTrace {
    /**
     * path defines the chain of port/channel identifiers used for tracing the
     * source of the fungible token.
     */
    path: string;
    /** base denomination of the relayed fungible token. */
    baseDenom: string;
}
/**
 * Params defines the set of IBC transfer parameters.
 * NOTE: To prevent a single token from being transferred, set the
 * TransfersEnabled parameter to true and then set the bank module's SendEnabled
 * parameter for the denomination to false.
 */
export interface Params {
    /**
     * send_enabled enables or disables all cross-chain token transfers from this
     * chain.
     */
    sendEnabled: boolean;
    /**
     * receive_enabled enables or disables all cross-chain token transfers to this
     * chain.
     */
    receiveEnabled: boolean;
}
export declare const DenomTrace: {
    typeUrl: string;
    encode(message: DenomTrace, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): DenomTrace;
    fromJSON(object: any): DenomTrace;
    toJSON(message: DenomTrace): unknown;
    fromPartial<I extends {
        path?: string | undefined;
        baseDenom?: string | undefined;
    } & {
        path?: string | undefined;
        baseDenom?: string | undefined;
    } & Record<Exclude<keyof I, keyof DenomTrace>, never>>(object: I): DenomTrace;
};
export declare const Params: {
    typeUrl: string;
    encode(message: Params, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial<I extends {
        sendEnabled?: boolean | undefined;
        receiveEnabled?: boolean | undefined;
    } & {
        sendEnabled?: boolean | undefined;
        receiveEnabled?: boolean | undefined;
    } & Record<Exclude<keyof I, keyof Params>, never>>(object: I): Params;
};
