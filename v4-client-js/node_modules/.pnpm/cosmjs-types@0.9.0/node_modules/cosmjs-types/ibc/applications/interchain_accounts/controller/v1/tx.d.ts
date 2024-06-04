import { InterchainAccountPacketData } from "../../v1/packet";
import { BinaryReader, BinaryWriter } from "../../../../../binary";
import { Rpc } from "../../../../../helpers";
export declare const protobufPackage = "ibc.applications.interchain_accounts.controller.v1";
/** MsgRegisterInterchainAccount defines the payload for Msg/RegisterAccount */
export interface MsgRegisterInterchainAccount {
    owner: string;
    connectionId: string;
    version: string;
}
/** MsgRegisterInterchainAccountResponse defines the response for Msg/RegisterAccount */
export interface MsgRegisterInterchainAccountResponse {
    channelId: string;
}
/** MsgSendTx defines the payload for Msg/SendTx */
export interface MsgSendTx {
    owner: string;
    connectionId: string;
    packetData: InterchainAccountPacketData;
    /**
     * Relative timeout timestamp provided will be added to the current block time during transaction execution.
     * The timeout timestamp must be non-zero.
     */
    relativeTimeout: bigint;
}
/** MsgSendTxResponse defines the response for MsgSendTx */
export interface MsgSendTxResponse {
    sequence: bigint;
}
export declare const MsgRegisterInterchainAccount: {
    typeUrl: string;
    encode(message: MsgRegisterInterchainAccount, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRegisterInterchainAccount;
    fromJSON(object: any): MsgRegisterInterchainAccount;
    toJSON(message: MsgRegisterInterchainAccount): unknown;
    fromPartial<I extends {
        owner?: string | undefined;
        connectionId?: string | undefined;
        version?: string | undefined;
    } & {
        owner?: string | undefined;
        connectionId?: string | undefined;
        version?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgRegisterInterchainAccount>, never>>(object: I): MsgRegisterInterchainAccount;
};
export declare const MsgRegisterInterchainAccountResponse: {
    typeUrl: string;
    encode(message: MsgRegisterInterchainAccountResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRegisterInterchainAccountResponse;
    fromJSON(object: any): MsgRegisterInterchainAccountResponse;
    toJSON(message: MsgRegisterInterchainAccountResponse): unknown;
    fromPartial<I extends {
        channelId?: string | undefined;
    } & {
        channelId?: string | undefined;
    } & Record<Exclude<keyof I, "channelId">, never>>(object: I): MsgRegisterInterchainAccountResponse;
};
export declare const MsgSendTx: {
    typeUrl: string;
    encode(message: MsgSendTx, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgSendTx;
    fromJSON(object: any): MsgSendTx;
    toJSON(message: MsgSendTx): unknown;
    fromPartial<I extends {
        owner?: string | undefined;
        connectionId?: string | undefined;
        packetData?: {
            type?: import("../../v1/packet").Type | undefined;
            data?: Uint8Array | undefined;
            memo?: string | undefined;
        } | undefined;
        relativeTimeout?: bigint | undefined;
    } & {
        owner?: string | undefined;
        connectionId?: string | undefined;
        packetData?: ({
            type?: import("../../v1/packet").Type | undefined;
            data?: Uint8Array | undefined;
            memo?: string | undefined;
        } & {
            type?: import("../../v1/packet").Type | undefined;
            data?: Uint8Array | undefined;
            memo?: string | undefined;
        } & Record<Exclude<keyof I["packetData"], keyof InterchainAccountPacketData>, never>) | undefined;
        relativeTimeout?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof MsgSendTx>, never>>(object: I): MsgSendTx;
};
export declare const MsgSendTxResponse: {
    typeUrl: string;
    encode(message: MsgSendTxResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgSendTxResponse;
    fromJSON(object: any): MsgSendTxResponse;
    toJSON(message: MsgSendTxResponse): unknown;
    fromPartial<I extends {
        sequence?: bigint | undefined;
    } & {
        sequence?: bigint | undefined;
    } & Record<Exclude<keyof I, "sequence">, never>>(object: I): MsgSendTxResponse;
};
/** Msg defines the 27-interchain-accounts/controller Msg service. */
export interface Msg {
    /** RegisterInterchainAccount defines a rpc handler for MsgRegisterInterchainAccount. */
    RegisterInterchainAccount(request: MsgRegisterInterchainAccount): Promise<MsgRegisterInterchainAccountResponse>;
    /** SendTx defines a rpc handler for MsgSendTx. */
    SendTx(request: MsgSendTx): Promise<MsgSendTxResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    RegisterInterchainAccount(request: MsgRegisterInterchainAccount): Promise<MsgRegisterInterchainAccountResponse>;
    SendTx(request: MsgSendTx): Promise<MsgSendTxResponse>;
}
