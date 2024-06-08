import { BinaryReader, BinaryWriter } from "../../../binary";
import { Rpc } from "../../../helpers";
export declare const protobufPackage = "cosmos.nft.v1beta1";
/** MsgSend represents a message to send a nft from one account to another account. */
export interface MsgSend {
    /** class_id defines the unique identifier of the nft classification, similar to the contract address of ERC721 */
    classId: string;
    /** id defines the unique identification of nft */
    id: string;
    /** sender is the address of the owner of nft */
    sender: string;
    /** receiver is the receiver address of nft */
    receiver: string;
}
/** MsgSendResponse defines the Msg/Send response type. */
export interface MsgSendResponse {
}
export declare const MsgSend: {
    typeUrl: string;
    encode(message: MsgSend, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgSend;
    fromJSON(object: any): MsgSend;
    toJSON(message: MsgSend): unknown;
    fromPartial<I extends {
        classId?: string | undefined;
        id?: string | undefined;
        sender?: string | undefined;
        receiver?: string | undefined;
    } & {
        classId?: string | undefined;
        id?: string | undefined;
        sender?: string | undefined;
        receiver?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgSend>, never>>(object: I): MsgSend;
};
export declare const MsgSendResponse: {
    typeUrl: string;
    encode(_: MsgSendResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgSendResponse;
    fromJSON(_: any): MsgSendResponse;
    toJSON(_: MsgSendResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgSendResponse;
};
/** Msg defines the nft Msg service. */
export interface Msg {
    /** Send defines a method to send a nft from one account to another account. */
    Send(request: MsgSend): Promise<MsgSendResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    Send(request: MsgSend): Promise<MsgSendResponse>;
}
