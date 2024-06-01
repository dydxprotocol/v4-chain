import { Channel, Packet } from "./channel";
import { Height } from "../../client/v1/client";
import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "ibc.core.channel.v1";
/** ResponseResultType defines the possible outcomes of the execution of a message */
export declare enum ResponseResultType {
    /** RESPONSE_RESULT_TYPE_UNSPECIFIED - Default zero value enumeration */
    RESPONSE_RESULT_TYPE_UNSPECIFIED = 0,
    /** RESPONSE_RESULT_TYPE_NOOP - The message did not call the IBC application callbacks (because, for example, the packet had already been relayed) */
    RESPONSE_RESULT_TYPE_NOOP = 1,
    /** RESPONSE_RESULT_TYPE_SUCCESS - The message was executed successfully */
    RESPONSE_RESULT_TYPE_SUCCESS = 2,
    UNRECOGNIZED = -1
}
export declare function responseResultTypeFromJSON(object: any): ResponseResultType;
export declare function responseResultTypeToJSON(object: ResponseResultType): string;
/**
 * MsgChannelOpenInit defines an sdk.Msg to initialize a channel handshake. It
 * is called by a relayer on Chain A.
 */
export interface MsgChannelOpenInit {
    portId: string;
    channel: Channel;
    signer: string;
}
/** MsgChannelOpenInitResponse defines the Msg/ChannelOpenInit response type. */
export interface MsgChannelOpenInitResponse {
    channelId: string;
    version: string;
}
/**
 * MsgChannelOpenInit defines a msg sent by a Relayer to try to open a channel
 * on Chain B. The version field within the Channel field has been deprecated. Its
 * value will be ignored by core IBC.
 */
export interface MsgChannelOpenTry {
    portId: string;
    /** Deprecated: this field is unused. Crossing hello's are no longer supported in core IBC. */
    /** @deprecated */
    previousChannelId: string;
    /** NOTE: the version field within the channel has been deprecated. Its value will be ignored by core IBC. */
    channel: Channel;
    counterpartyVersion: string;
    proofInit: Uint8Array;
    proofHeight: Height;
    signer: string;
}
/** MsgChannelOpenTryResponse defines the Msg/ChannelOpenTry response type. */
export interface MsgChannelOpenTryResponse {
    version: string;
}
/**
 * MsgChannelOpenAck defines a msg sent by a Relayer to Chain A to acknowledge
 * the change of channel state to TRYOPEN on Chain B.
 */
export interface MsgChannelOpenAck {
    portId: string;
    channelId: string;
    counterpartyChannelId: string;
    counterpartyVersion: string;
    proofTry: Uint8Array;
    proofHeight: Height;
    signer: string;
}
/** MsgChannelOpenAckResponse defines the Msg/ChannelOpenAck response type. */
export interface MsgChannelOpenAckResponse {
}
/**
 * MsgChannelOpenConfirm defines a msg sent by a Relayer to Chain B to
 * acknowledge the change of channel state to OPEN on Chain A.
 */
export interface MsgChannelOpenConfirm {
    portId: string;
    channelId: string;
    proofAck: Uint8Array;
    proofHeight: Height;
    signer: string;
}
/**
 * MsgChannelOpenConfirmResponse defines the Msg/ChannelOpenConfirm response
 * type.
 */
export interface MsgChannelOpenConfirmResponse {
}
/**
 * MsgChannelCloseInit defines a msg sent by a Relayer to Chain A
 * to close a channel with Chain B.
 */
export interface MsgChannelCloseInit {
    portId: string;
    channelId: string;
    signer: string;
}
/** MsgChannelCloseInitResponse defines the Msg/ChannelCloseInit response type. */
export interface MsgChannelCloseInitResponse {
}
/**
 * MsgChannelCloseConfirm defines a msg sent by a Relayer to Chain B
 * to acknowledge the change of channel state to CLOSED on Chain A.
 */
export interface MsgChannelCloseConfirm {
    portId: string;
    channelId: string;
    proofInit: Uint8Array;
    proofHeight: Height;
    signer: string;
}
/**
 * MsgChannelCloseConfirmResponse defines the Msg/ChannelCloseConfirm response
 * type.
 */
export interface MsgChannelCloseConfirmResponse {
}
/** MsgRecvPacket receives incoming IBC packet */
export interface MsgRecvPacket {
    packet: Packet;
    proofCommitment: Uint8Array;
    proofHeight: Height;
    signer: string;
}
/** MsgRecvPacketResponse defines the Msg/RecvPacket response type. */
export interface MsgRecvPacketResponse {
    result: ResponseResultType;
}
/** MsgTimeout receives timed-out packet */
export interface MsgTimeout {
    packet: Packet;
    proofUnreceived: Uint8Array;
    proofHeight: Height;
    nextSequenceRecv: bigint;
    signer: string;
}
/** MsgTimeoutResponse defines the Msg/Timeout response type. */
export interface MsgTimeoutResponse {
    result: ResponseResultType;
}
/** MsgTimeoutOnClose timed-out packet upon counterparty channel closure. */
export interface MsgTimeoutOnClose {
    packet: Packet;
    proofUnreceived: Uint8Array;
    proofClose: Uint8Array;
    proofHeight: Height;
    nextSequenceRecv: bigint;
    signer: string;
}
/** MsgTimeoutOnCloseResponse defines the Msg/TimeoutOnClose response type. */
export interface MsgTimeoutOnCloseResponse {
    result: ResponseResultType;
}
/** MsgAcknowledgement receives incoming IBC acknowledgement */
export interface MsgAcknowledgement {
    packet: Packet;
    acknowledgement: Uint8Array;
    proofAcked: Uint8Array;
    proofHeight: Height;
    signer: string;
}
/** MsgAcknowledgementResponse defines the Msg/Acknowledgement response type. */
export interface MsgAcknowledgementResponse {
    result: ResponseResultType;
}
export declare const MsgChannelOpenInit: {
    typeUrl: string;
    encode(message: MsgChannelOpenInit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelOpenInit;
    fromJSON(object: any): MsgChannelOpenInit;
    toJSON(message: MsgChannelOpenInit): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channel?: {
            state?: import("./channel").State | undefined;
            ordering?: import("./channel").Order | undefined;
            counterparty?: {
                portId?: string | undefined;
                channelId?: string | undefined;
            } | undefined;
            connectionHops?: string[] | undefined;
            version?: string | undefined;
        } | undefined;
        signer?: string | undefined;
    } & {
        portId?: string | undefined;
        channel?: ({
            state?: import("./channel").State | undefined;
            ordering?: import("./channel").Order | undefined;
            counterparty?: {
                portId?: string | undefined;
                channelId?: string | undefined;
            } | undefined;
            connectionHops?: string[] | undefined;
            version?: string | undefined;
        } & {
            state?: import("./channel").State | undefined;
            ordering?: import("./channel").Order | undefined;
            counterparty?: ({
                portId?: string | undefined;
                channelId?: string | undefined;
            } & {
                portId?: string | undefined;
                channelId?: string | undefined;
            } & Record<Exclude<keyof I["channel"]["counterparty"], keyof import("./channel").Counterparty>, never>) | undefined;
            connectionHops?: (string[] & string[] & Record<Exclude<keyof I["channel"]["connectionHops"], keyof string[]>, never>) | undefined;
            version?: string | undefined;
        } & Record<Exclude<keyof I["channel"], keyof Channel>, never>) | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgChannelOpenInit>, never>>(object: I): MsgChannelOpenInit;
};
export declare const MsgChannelOpenInitResponse: {
    typeUrl: string;
    encode(message: MsgChannelOpenInitResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelOpenInitResponse;
    fromJSON(object: any): MsgChannelOpenInitResponse;
    toJSON(message: MsgChannelOpenInitResponse): unknown;
    fromPartial<I extends {
        channelId?: string | undefined;
        version?: string | undefined;
    } & {
        channelId?: string | undefined;
        version?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgChannelOpenInitResponse>, never>>(object: I): MsgChannelOpenInitResponse;
};
export declare const MsgChannelOpenTry: {
    typeUrl: string;
    encode(message: MsgChannelOpenTry, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelOpenTry;
    fromJSON(object: any): MsgChannelOpenTry;
    toJSON(message: MsgChannelOpenTry): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        previousChannelId?: string | undefined;
        channel?: {
            state?: import("./channel").State | undefined;
            ordering?: import("./channel").Order | undefined;
            counterparty?: {
                portId?: string | undefined;
                channelId?: string | undefined;
            } | undefined;
            connectionHops?: string[] | undefined;
            version?: string | undefined;
        } | undefined;
        counterpartyVersion?: string | undefined;
        proofInit?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        signer?: string | undefined;
    } & {
        portId?: string | undefined;
        previousChannelId?: string | undefined;
        channel?: ({
            state?: import("./channel").State | undefined;
            ordering?: import("./channel").Order | undefined;
            counterparty?: {
                portId?: string | undefined;
                channelId?: string | undefined;
            } | undefined;
            connectionHops?: string[] | undefined;
            version?: string | undefined;
        } & {
            state?: import("./channel").State | undefined;
            ordering?: import("./channel").Order | undefined;
            counterparty?: ({
                portId?: string | undefined;
                channelId?: string | undefined;
            } & {
                portId?: string | undefined;
                channelId?: string | undefined;
            } & Record<Exclude<keyof I["channel"]["counterparty"], keyof import("./channel").Counterparty>, never>) | undefined;
            connectionHops?: (string[] & string[] & Record<Exclude<keyof I["channel"]["connectionHops"], keyof string[]>, never>) | undefined;
            version?: string | undefined;
        } & Record<Exclude<keyof I["channel"], keyof Channel>, never>) | undefined;
        counterpartyVersion?: string | undefined;
        proofInit?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgChannelOpenTry>, never>>(object: I): MsgChannelOpenTry;
};
export declare const MsgChannelOpenTryResponse: {
    typeUrl: string;
    encode(message: MsgChannelOpenTryResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelOpenTryResponse;
    fromJSON(object: any): MsgChannelOpenTryResponse;
    toJSON(message: MsgChannelOpenTryResponse): unknown;
    fromPartial<I extends {
        version?: string | undefined;
    } & {
        version?: string | undefined;
    } & Record<Exclude<keyof I, "version">, never>>(object: I): MsgChannelOpenTryResponse;
};
export declare const MsgChannelOpenAck: {
    typeUrl: string;
    encode(message: MsgChannelOpenAck, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelOpenAck;
    fromJSON(object: any): MsgChannelOpenAck;
    toJSON(message: MsgChannelOpenAck): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        counterpartyChannelId?: string | undefined;
        counterpartyVersion?: string | undefined;
        proofTry?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        signer?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        counterpartyChannelId?: string | undefined;
        counterpartyVersion?: string | undefined;
        proofTry?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgChannelOpenAck>, never>>(object: I): MsgChannelOpenAck;
};
export declare const MsgChannelOpenAckResponse: {
    typeUrl: string;
    encode(_: MsgChannelOpenAckResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelOpenAckResponse;
    fromJSON(_: any): MsgChannelOpenAckResponse;
    toJSON(_: MsgChannelOpenAckResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgChannelOpenAckResponse;
};
export declare const MsgChannelOpenConfirm: {
    typeUrl: string;
    encode(message: MsgChannelOpenConfirm, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelOpenConfirm;
    fromJSON(object: any): MsgChannelOpenConfirm;
    toJSON(message: MsgChannelOpenConfirm): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        proofAck?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        signer?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        proofAck?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgChannelOpenConfirm>, never>>(object: I): MsgChannelOpenConfirm;
};
export declare const MsgChannelOpenConfirmResponse: {
    typeUrl: string;
    encode(_: MsgChannelOpenConfirmResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelOpenConfirmResponse;
    fromJSON(_: any): MsgChannelOpenConfirmResponse;
    toJSON(_: MsgChannelOpenConfirmResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgChannelOpenConfirmResponse;
};
export declare const MsgChannelCloseInit: {
    typeUrl: string;
    encode(message: MsgChannelCloseInit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelCloseInit;
    fromJSON(object: any): MsgChannelCloseInit;
    toJSON(message: MsgChannelCloseInit): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        signer?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgChannelCloseInit>, never>>(object: I): MsgChannelCloseInit;
};
export declare const MsgChannelCloseInitResponse: {
    typeUrl: string;
    encode(_: MsgChannelCloseInitResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelCloseInitResponse;
    fromJSON(_: any): MsgChannelCloseInitResponse;
    toJSON(_: MsgChannelCloseInitResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgChannelCloseInitResponse;
};
export declare const MsgChannelCloseConfirm: {
    typeUrl: string;
    encode(message: MsgChannelCloseConfirm, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelCloseConfirm;
    fromJSON(object: any): MsgChannelCloseConfirm;
    toJSON(message: MsgChannelCloseConfirm): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        proofInit?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        signer?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        proofInit?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgChannelCloseConfirm>, never>>(object: I): MsgChannelCloseConfirm;
};
export declare const MsgChannelCloseConfirmResponse: {
    typeUrl: string;
    encode(_: MsgChannelCloseConfirmResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgChannelCloseConfirmResponse;
    fromJSON(_: any): MsgChannelCloseConfirmResponse;
    toJSON(_: MsgChannelCloseConfirmResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgChannelCloseConfirmResponse;
};
export declare const MsgRecvPacket: {
    typeUrl: string;
    encode(message: MsgRecvPacket, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRecvPacket;
    fromJSON(object: any): MsgRecvPacket;
    toJSON(message: MsgRecvPacket): unknown;
    fromPartial<I extends {
        packet?: {
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            timeoutTimestamp?: bigint | undefined;
        } | undefined;
        proofCommitment?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        signer?: string | undefined;
    } & {
        packet?: ({
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            timeoutTimestamp?: bigint | undefined;
        } & {
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: ({
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & Record<Exclude<keyof I["packet"]["timeoutHeight"], keyof Height>, never>) | undefined;
            timeoutTimestamp?: bigint | undefined;
        } & Record<Exclude<keyof I["packet"], keyof Packet>, never>) | undefined;
        proofCommitment?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgRecvPacket>, never>>(object: I): MsgRecvPacket;
};
export declare const MsgRecvPacketResponse: {
    typeUrl: string;
    encode(message: MsgRecvPacketResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgRecvPacketResponse;
    fromJSON(object: any): MsgRecvPacketResponse;
    toJSON(message: MsgRecvPacketResponse): unknown;
    fromPartial<I extends {
        result?: ResponseResultType | undefined;
    } & {
        result?: ResponseResultType | undefined;
    } & Record<Exclude<keyof I, "result">, never>>(object: I): MsgRecvPacketResponse;
};
export declare const MsgTimeout: {
    typeUrl: string;
    encode(message: MsgTimeout, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgTimeout;
    fromJSON(object: any): MsgTimeout;
    toJSON(message: MsgTimeout): unknown;
    fromPartial<I extends {
        packet?: {
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            timeoutTimestamp?: bigint | undefined;
        } | undefined;
        proofUnreceived?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        nextSequenceRecv?: bigint | undefined;
        signer?: string | undefined;
    } & {
        packet?: ({
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            timeoutTimestamp?: bigint | undefined;
        } & {
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: ({
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & Record<Exclude<keyof I["packet"]["timeoutHeight"], keyof Height>, never>) | undefined;
            timeoutTimestamp?: bigint | undefined;
        } & Record<Exclude<keyof I["packet"], keyof Packet>, never>) | undefined;
        proofUnreceived?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
        nextSequenceRecv?: bigint | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgTimeout>, never>>(object: I): MsgTimeout;
};
export declare const MsgTimeoutResponse: {
    typeUrl: string;
    encode(message: MsgTimeoutResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgTimeoutResponse;
    fromJSON(object: any): MsgTimeoutResponse;
    toJSON(message: MsgTimeoutResponse): unknown;
    fromPartial<I extends {
        result?: ResponseResultType | undefined;
    } & {
        result?: ResponseResultType | undefined;
    } & Record<Exclude<keyof I, "result">, never>>(object: I): MsgTimeoutResponse;
};
export declare const MsgTimeoutOnClose: {
    typeUrl: string;
    encode(message: MsgTimeoutOnClose, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgTimeoutOnClose;
    fromJSON(object: any): MsgTimeoutOnClose;
    toJSON(message: MsgTimeoutOnClose): unknown;
    fromPartial<I extends {
        packet?: {
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            timeoutTimestamp?: bigint | undefined;
        } | undefined;
        proofUnreceived?: Uint8Array | undefined;
        proofClose?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        nextSequenceRecv?: bigint | undefined;
        signer?: string | undefined;
    } & {
        packet?: ({
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            timeoutTimestamp?: bigint | undefined;
        } & {
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: ({
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & Record<Exclude<keyof I["packet"]["timeoutHeight"], keyof Height>, never>) | undefined;
            timeoutTimestamp?: bigint | undefined;
        } & Record<Exclude<keyof I["packet"], keyof Packet>, never>) | undefined;
        proofUnreceived?: Uint8Array | undefined;
        proofClose?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
        nextSequenceRecv?: bigint | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgTimeoutOnClose>, never>>(object: I): MsgTimeoutOnClose;
};
export declare const MsgTimeoutOnCloseResponse: {
    typeUrl: string;
    encode(message: MsgTimeoutOnCloseResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgTimeoutOnCloseResponse;
    fromJSON(object: any): MsgTimeoutOnCloseResponse;
    toJSON(message: MsgTimeoutOnCloseResponse): unknown;
    fromPartial<I extends {
        result?: ResponseResultType | undefined;
    } & {
        result?: ResponseResultType | undefined;
    } & Record<Exclude<keyof I, "result">, never>>(object: I): MsgTimeoutOnCloseResponse;
};
export declare const MsgAcknowledgement: {
    typeUrl: string;
    encode(message: MsgAcknowledgement, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgAcknowledgement;
    fromJSON(object: any): MsgAcknowledgement;
    toJSON(message: MsgAcknowledgement): unknown;
    fromPartial<I extends {
        packet?: {
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            timeoutTimestamp?: bigint | undefined;
        } | undefined;
        acknowledgement?: Uint8Array | undefined;
        proofAcked?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        signer?: string | undefined;
    } & {
        packet?: ({
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } | undefined;
            timeoutTimestamp?: bigint | undefined;
        } & {
            sequence?: bigint | undefined;
            sourcePort?: string | undefined;
            sourceChannel?: string | undefined;
            destinationPort?: string | undefined;
            destinationChannel?: string | undefined;
            data?: Uint8Array | undefined;
            timeoutHeight?: ({
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & {
                revisionNumber?: bigint | undefined;
                revisionHeight?: bigint | undefined;
            } & Record<Exclude<keyof I["packet"]["timeoutHeight"], keyof Height>, never>) | undefined;
            timeoutTimestamp?: bigint | undefined;
        } & Record<Exclude<keyof I["packet"], keyof Packet>, never>) | undefined;
        acknowledgement?: Uint8Array | undefined;
        proofAcked?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgAcknowledgement>, never>>(object: I): MsgAcknowledgement;
};
export declare const MsgAcknowledgementResponse: {
    typeUrl: string;
    encode(message: MsgAcknowledgementResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgAcknowledgementResponse;
    fromJSON(object: any): MsgAcknowledgementResponse;
    toJSON(message: MsgAcknowledgementResponse): unknown;
    fromPartial<I extends {
        result?: ResponseResultType | undefined;
    } & {
        result?: ResponseResultType | undefined;
    } & Record<Exclude<keyof I, "result">, never>>(object: I): MsgAcknowledgementResponse;
};
/** Msg defines the ibc/channel Msg service. */
export interface Msg {
    /** ChannelOpenInit defines a rpc handler method for MsgChannelOpenInit. */
    ChannelOpenInit(request: MsgChannelOpenInit): Promise<MsgChannelOpenInitResponse>;
    /** ChannelOpenTry defines a rpc handler method for MsgChannelOpenTry. */
    ChannelOpenTry(request: MsgChannelOpenTry): Promise<MsgChannelOpenTryResponse>;
    /** ChannelOpenAck defines a rpc handler method for MsgChannelOpenAck. */
    ChannelOpenAck(request: MsgChannelOpenAck): Promise<MsgChannelOpenAckResponse>;
    /** ChannelOpenConfirm defines a rpc handler method for MsgChannelOpenConfirm. */
    ChannelOpenConfirm(request: MsgChannelOpenConfirm): Promise<MsgChannelOpenConfirmResponse>;
    /** ChannelCloseInit defines a rpc handler method for MsgChannelCloseInit. */
    ChannelCloseInit(request: MsgChannelCloseInit): Promise<MsgChannelCloseInitResponse>;
    /**
     * ChannelCloseConfirm defines a rpc handler method for
     * MsgChannelCloseConfirm.
     */
    ChannelCloseConfirm(request: MsgChannelCloseConfirm): Promise<MsgChannelCloseConfirmResponse>;
    /** RecvPacket defines a rpc handler method for MsgRecvPacket. */
    RecvPacket(request: MsgRecvPacket): Promise<MsgRecvPacketResponse>;
    /** Timeout defines a rpc handler method for MsgTimeout. */
    Timeout(request: MsgTimeout): Promise<MsgTimeoutResponse>;
    /** TimeoutOnClose defines a rpc handler method for MsgTimeoutOnClose. */
    TimeoutOnClose(request: MsgTimeoutOnClose): Promise<MsgTimeoutOnCloseResponse>;
    /** Acknowledgement defines a rpc handler method for MsgAcknowledgement. */
    Acknowledgement(request: MsgAcknowledgement): Promise<MsgAcknowledgementResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    ChannelOpenInit(request: MsgChannelOpenInit): Promise<MsgChannelOpenInitResponse>;
    ChannelOpenTry(request: MsgChannelOpenTry): Promise<MsgChannelOpenTryResponse>;
    ChannelOpenAck(request: MsgChannelOpenAck): Promise<MsgChannelOpenAckResponse>;
    ChannelOpenConfirm(request: MsgChannelOpenConfirm): Promise<MsgChannelOpenConfirmResponse>;
    ChannelCloseInit(request: MsgChannelCloseInit): Promise<MsgChannelCloseInitResponse>;
    ChannelCloseConfirm(request: MsgChannelCloseConfirm): Promise<MsgChannelCloseConfirmResponse>;
    RecvPacket(request: MsgRecvPacket): Promise<MsgRecvPacketResponse>;
    Timeout(request: MsgTimeout): Promise<MsgTimeoutResponse>;
    TimeoutOnClose(request: MsgTimeoutOnClose): Promise<MsgTimeoutOnCloseResponse>;
    Acknowledgement(request: MsgAcknowledgement): Promise<MsgAcknowledgementResponse>;
}
