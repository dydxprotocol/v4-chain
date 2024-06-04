import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmwasm.wasm.v1";
/** MsgIBCSend */
export interface MsgIBCSend {
    /** the channel by which the packet will be sent */
    channel: string;
    /**
     * Timeout height relative to the current block height.
     * The timeout is disabled when set to 0.
     */
    timeoutHeight: bigint;
    /**
     * Timeout timestamp (in nanoseconds) relative to the current block timestamp.
     * The timeout is disabled when set to 0.
     */
    timeoutTimestamp: bigint;
    /**
     * Data is the payload to transfer. We must not make assumption what format or
     * content is in here.
     */
    data: Uint8Array;
}
/** MsgIBCSendResponse */
export interface MsgIBCSendResponse {
    /** Sequence number of the IBC packet sent */
    sequence: bigint;
}
/** MsgIBCCloseChannel port and channel need to be owned by the contract */
export interface MsgIBCCloseChannel {
    channel: string;
}
export declare const MsgIBCSend: {
    typeUrl: string;
    encode(message: MsgIBCSend, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgIBCSend;
    fromJSON(object: any): MsgIBCSend;
    toJSON(message: MsgIBCSend): unknown;
    fromPartial<I extends {
        channel?: string | undefined;
        timeoutHeight?: bigint | undefined;
        timeoutTimestamp?: bigint | undefined;
        data?: Uint8Array | undefined;
    } & {
        channel?: string | undefined;
        timeoutHeight?: bigint | undefined;
        timeoutTimestamp?: bigint | undefined;
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof MsgIBCSend>, never>>(object: I): MsgIBCSend;
};
export declare const MsgIBCSendResponse: {
    typeUrl: string;
    encode(message: MsgIBCSendResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgIBCSendResponse;
    fromJSON(object: any): MsgIBCSendResponse;
    toJSON(message: MsgIBCSendResponse): unknown;
    fromPartial<I extends {
        sequence?: bigint | undefined;
    } & {
        sequence?: bigint | undefined;
    } & Record<Exclude<keyof I, "sequence">, never>>(object: I): MsgIBCSendResponse;
};
export declare const MsgIBCCloseChannel: {
    typeUrl: string;
    encode(message: MsgIBCCloseChannel, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgIBCCloseChannel;
    fromJSON(object: any): MsgIBCCloseChannel;
    toJSON(message: MsgIBCCloseChannel): unknown;
    fromPartial<I extends {
        channel?: string | undefined;
    } & {
        channel?: string | undefined;
    } & Record<Exclude<keyof I, "channel">, never>>(object: I): MsgIBCCloseChannel;
};
