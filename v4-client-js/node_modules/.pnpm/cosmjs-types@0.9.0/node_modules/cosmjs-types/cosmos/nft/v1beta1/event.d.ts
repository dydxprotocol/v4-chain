import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.nft.v1beta1";
/** EventSend is emitted on Msg/Send */
export interface EventSend {
    /** class_id associated with the nft */
    classId: string;
    /** id is a unique identifier of the nft */
    id: string;
    /** sender is the address of the owner of nft */
    sender: string;
    /** receiver is the receiver address of nft */
    receiver: string;
}
/** EventMint is emitted on Mint */
export interface EventMint {
    /** class_id associated with the nft */
    classId: string;
    /** id is a unique identifier of the nft */
    id: string;
    /** owner is the owner address of the nft */
    owner: string;
}
/** EventBurn is emitted on Burn */
export interface EventBurn {
    /** class_id associated with the nft */
    classId: string;
    /** id is a unique identifier of the nft */
    id: string;
    /** owner is the owner address of the nft */
    owner: string;
}
export declare const EventSend: {
    typeUrl: string;
    encode(message: EventSend, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventSend;
    fromJSON(object: any): EventSend;
    toJSON(message: EventSend): unknown;
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
    } & Record<Exclude<keyof I, keyof EventSend>, never>>(object: I): EventSend;
};
export declare const EventMint: {
    typeUrl: string;
    encode(message: EventMint, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventMint;
    fromJSON(object: any): EventMint;
    toJSON(message: EventMint): unknown;
    fromPartial<I extends {
        classId?: string | undefined;
        id?: string | undefined;
        owner?: string | undefined;
    } & {
        classId?: string | undefined;
        id?: string | undefined;
        owner?: string | undefined;
    } & Record<Exclude<keyof I, keyof EventMint>, never>>(object: I): EventMint;
};
export declare const EventBurn: {
    typeUrl: string;
    encode(message: EventBurn, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): EventBurn;
    fromJSON(object: any): EventBurn;
    toJSON(message: EventBurn): unknown;
    fromPartial<I extends {
        classId?: string | undefined;
        id?: string | undefined;
        owner?: string | undefined;
    } & {
        classId?: string | undefined;
        id?: string | undefined;
        owner?: string | undefined;
    } & Record<Exclude<keyof I, keyof EventBurn>, never>>(object: I): EventBurn;
};
