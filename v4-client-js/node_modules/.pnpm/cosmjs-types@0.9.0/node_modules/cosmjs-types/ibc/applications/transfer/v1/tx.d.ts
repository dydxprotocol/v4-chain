import { Coin } from "../../../../cosmos/base/v1beta1/coin";
import { Height } from "../../../core/client/v1/client";
import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "ibc.applications.transfer.v1";
/**
 * MsgTransfer defines a msg to transfer fungible tokens (i.e Coins) between
 * ICS20 enabled chains. See ICS Spec here:
 * https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#data-structures
 */
export interface MsgTransfer {
    /** the port on which the packet will be sent */
    sourcePort: string;
    /** the channel by which the packet will be sent */
    sourceChannel: string;
    /** the tokens to be transferred */
    token: Coin;
    /** the sender address */
    sender: string;
    /** the recipient address on the destination chain */
    receiver: string;
    /**
     * Timeout height relative to the current block height.
     * The timeout is disabled when set to 0.
     */
    timeoutHeight: Height;
    /**
     * Timeout timestamp in absolute nanoseconds since unix epoch.
     * The timeout is disabled when set to 0.
     */
    timeoutTimestamp: bigint;
    /** optional memo */
    memo: string;
}
/** MsgTransferResponse defines the Msg/Transfer response type. */
export interface MsgTransferResponse {
    /** sequence number of the transfer packet sent */
    sequence: bigint;
}
export declare const MsgTransfer: {
    typeUrl: string;
    encode(message: MsgTransfer, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgTransfer;
    fromJSON(object: any): MsgTransfer;
    toJSON(message: MsgTransfer): unknown;
    fromPartial<I extends {
        sourcePort?: string | undefined;
        sourceChannel?: string | undefined;
        token?: {
            denom?: string | undefined;
            amount?: string | undefined;
        } | undefined;
        sender?: string | undefined;
        receiver?: string | undefined;
        timeoutHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        timeoutTimestamp?: bigint | undefined;
        memo?: string | undefined;
    } & {
        sourcePort?: string | undefined;
        sourceChannel?: string | undefined;
        token?: ({
            denom?: string | undefined;
            amount?: string | undefined;
        } & {
            denom?: string | undefined;
            amount?: string | undefined;
        } & Record<Exclude<keyof I["token"], keyof Coin>, never>) | undefined;
        sender?: string | undefined;
        receiver?: string | undefined;
        timeoutHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["timeoutHeight"], keyof Height>, never>) | undefined;
        timeoutTimestamp?: bigint | undefined;
        memo?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgTransfer>, never>>(object: I): MsgTransfer;
};
export declare const MsgTransferResponse: {
    typeUrl: string;
    encode(message: MsgTransferResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgTransferResponse;
    fromJSON(object: any): MsgTransferResponse;
    toJSON(message: MsgTransferResponse): unknown;
    fromPartial<I extends {
        sequence?: bigint | undefined;
    } & {
        sequence?: bigint | undefined;
    } & Record<Exclude<keyof I, "sequence">, never>>(object: I): MsgTransferResponse;
};
/** Msg defines the ibc/transfer Msg service. */
export interface Msg {
    /** Transfer defines a rpc handler method for MsgTransfer. */
    Transfer(request: MsgTransfer): Promise<MsgTransferResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    Transfer(request: MsgTransfer): Promise<MsgTransferResponse>;
}
