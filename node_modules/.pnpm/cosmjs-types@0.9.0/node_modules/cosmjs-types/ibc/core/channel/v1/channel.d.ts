import { Height } from "../../client/v1/client";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.core.channel.v1";
/**
 * State defines if a channel is in one of the following states:
 * CLOSED, INIT, TRYOPEN, OPEN or UNINITIALIZED.
 */
export declare enum State {
    /** STATE_UNINITIALIZED_UNSPECIFIED - Default State */
    STATE_UNINITIALIZED_UNSPECIFIED = 0,
    /** STATE_INIT - A channel has just started the opening handshake. */
    STATE_INIT = 1,
    /** STATE_TRYOPEN - A channel has acknowledged the handshake step on the counterparty chain. */
    STATE_TRYOPEN = 2,
    /**
     * STATE_OPEN - A channel has completed the handshake. Open channels are
     * ready to send and receive packets.
     */
    STATE_OPEN = 3,
    /**
     * STATE_CLOSED - A channel has been closed and can no longer be used to send or receive
     * packets.
     */
    STATE_CLOSED = 4,
    UNRECOGNIZED = -1
}
export declare function stateFromJSON(object: any): State;
export declare function stateToJSON(object: State): string;
/** Order defines if a channel is ORDERED or UNORDERED */
export declare enum Order {
    /** ORDER_NONE_UNSPECIFIED - zero-value for channel ordering */
    ORDER_NONE_UNSPECIFIED = 0,
    /**
     * ORDER_UNORDERED - packets can be delivered in any order, which may differ from the order in
     * which they were sent.
     */
    ORDER_UNORDERED = 1,
    /** ORDER_ORDERED - packets are delivered exactly in the order which they were sent */
    ORDER_ORDERED = 2,
    UNRECOGNIZED = -1
}
export declare function orderFromJSON(object: any): Order;
export declare function orderToJSON(object: Order): string;
/**
 * Channel defines pipeline for exactly-once packet delivery between specific
 * modules on separate blockchains, which has at least one end capable of
 * sending packets and one end capable of receiving packets.
 */
export interface Channel {
    /** current state of the channel end */
    state: State;
    /** whether the channel is ordered or unordered */
    ordering: Order;
    /** counterparty channel end */
    counterparty: Counterparty;
    /**
     * list of connection identifiers, in order, along which packets sent on
     * this channel will travel
     */
    connectionHops: string[];
    /** opaque channel version, which is agreed upon during the handshake */
    version: string;
}
/**
 * IdentifiedChannel defines a channel with additional port and channel
 * identifier fields.
 */
export interface IdentifiedChannel {
    /** current state of the channel end */
    state: State;
    /** whether the channel is ordered or unordered */
    ordering: Order;
    /** counterparty channel end */
    counterparty: Counterparty;
    /**
     * list of connection identifiers, in order, along which packets sent on
     * this channel will travel
     */
    connectionHops: string[];
    /** opaque channel version, which is agreed upon during the handshake */
    version: string;
    /** port identifier */
    portId: string;
    /** channel identifier */
    channelId: string;
}
/** Counterparty defines a channel end counterparty */
export interface Counterparty {
    /** port on the counterparty chain which owns the other end of the channel. */
    portId: string;
    /** channel end on the counterparty chain */
    channelId: string;
}
/** Packet defines a type that carries data across different chains through IBC */
export interface Packet {
    /**
     * number corresponds to the order of sends and receives, where a Packet
     * with an earlier sequence number must be sent and received before a Packet
     * with a later sequence number.
     */
    sequence: bigint;
    /** identifies the port on the sending chain. */
    sourcePort: string;
    /** identifies the channel end on the sending chain. */
    sourceChannel: string;
    /** identifies the port on the receiving chain. */
    destinationPort: string;
    /** identifies the channel end on the receiving chain. */
    destinationChannel: string;
    /** actual opaque bytes transferred directly to the application module */
    data: Uint8Array;
    /** block height after which the packet times out */
    timeoutHeight: Height;
    /** block timestamp (in nanoseconds) after which the packet times out */
    timeoutTimestamp: bigint;
}
/**
 * PacketState defines the generic type necessary to retrieve and store
 * packet commitments, acknowledgements, and receipts.
 * Caller is responsible for knowing the context necessary to interpret this
 * state as a commitment, acknowledgement, or a receipt.
 */
export interface PacketState {
    /** channel port identifier. */
    portId: string;
    /** channel unique identifier. */
    channelId: string;
    /** packet sequence. */
    sequence: bigint;
    /** embedded data that represents packet state. */
    data: Uint8Array;
}
/**
 * PacketId is an identifer for a unique Packet
 * Source chains refer to packets by source port/channel
 * Destination chains refer to packets by destination port/channel
 */
export interface PacketId {
    /** channel port identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
    /** packet sequence */
    sequence: bigint;
}
/**
 * Acknowledgement is the recommended acknowledgement format to be used by
 * app-specific protocols.
 * NOTE: The field numbers 21 and 22 were explicitly chosen to avoid accidental
 * conflicts with other protobuf message formats used for acknowledgements.
 * The first byte of any message with this format will be the non-ASCII values
 * `0xaa` (result) or `0xb2` (error). Implemented as defined by ICS:
 * https://github.com/cosmos/ibc/tree/master/spec/core/ics-004-channel-and-packet-semantics#acknowledgement-envelope
 */
export interface Acknowledgement {
    result?: Uint8Array;
    error?: string;
}
export declare const Channel: {
    typeUrl: string;
    encode(message: Channel, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Channel;
    fromJSON(object: any): Channel;
    toJSON(message: Channel): unknown;
    fromPartial<I extends {
        state?: State | undefined;
        ordering?: Order | undefined;
        counterparty?: {
            portId?: string | undefined;
            channelId?: string | undefined;
        } | undefined;
        connectionHops?: string[] | undefined;
        version?: string | undefined;
    } & {
        state?: State | undefined;
        ordering?: Order | undefined;
        counterparty?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
        } & Record<Exclude<keyof I["counterparty"], keyof Counterparty>, never>) | undefined;
        connectionHops?: (string[] & string[] & Record<Exclude<keyof I["connectionHops"], keyof string[]>, never>) | undefined;
        version?: string | undefined;
    } & Record<Exclude<keyof I, keyof Channel>, never>>(object: I): Channel;
};
export declare const IdentifiedChannel: {
    typeUrl: string;
    encode(message: IdentifiedChannel, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): IdentifiedChannel;
    fromJSON(object: any): IdentifiedChannel;
    toJSON(message: IdentifiedChannel): unknown;
    fromPartial<I extends {
        state?: State | undefined;
        ordering?: Order | undefined;
        counterparty?: {
            portId?: string | undefined;
            channelId?: string | undefined;
        } | undefined;
        connectionHops?: string[] | undefined;
        version?: string | undefined;
        portId?: string | undefined;
        channelId?: string | undefined;
    } & {
        state?: State | undefined;
        ordering?: Order | undefined;
        counterparty?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
        } & Record<Exclude<keyof I["counterparty"], keyof Counterparty>, never>) | undefined;
        connectionHops?: (string[] & string[] & Record<Exclude<keyof I["connectionHops"], keyof string[]>, never>) | undefined;
        version?: string | undefined;
        portId?: string | undefined;
        channelId?: string | undefined;
    } & Record<Exclude<keyof I, keyof IdentifiedChannel>, never>>(object: I): IdentifiedChannel;
};
export declare const Counterparty: {
    typeUrl: string;
    encode(message: Counterparty, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Counterparty;
    fromJSON(object: any): Counterparty;
    toJSON(message: Counterparty): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & Record<Exclude<keyof I, keyof Counterparty>, never>>(object: I): Counterparty;
};
export declare const Packet: {
    typeUrl: string;
    encode(message: Packet, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Packet;
    fromJSON(object: any): Packet;
    toJSON(message: Packet): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["timeoutHeight"], keyof Height>, never>) | undefined;
        timeoutTimestamp?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof Packet>, never>>(object: I): Packet;
};
export declare const PacketState: {
    typeUrl: string;
    encode(message: PacketState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PacketState;
    fromJSON(object: any): PacketState;
    toJSON(message: PacketState): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
        data?: Uint8Array | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof PacketState>, never>>(object: I): PacketState;
};
export declare const PacketId: {
    typeUrl: string;
    encode(message: PacketId, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PacketId;
    fromJSON(object: any): PacketId;
    toJSON(message: PacketId): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof PacketId>, never>>(object: I): PacketId;
};
export declare const Acknowledgement: {
    typeUrl: string;
    encode(message: Acknowledgement, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Acknowledgement;
    fromJSON(object: any): Acknowledgement;
    toJSON(message: Acknowledgement): unknown;
    fromPartial<I extends {
        result?: Uint8Array | undefined;
        error?: string | undefined;
    } & {
        result?: Uint8Array | undefined;
        error?: string | undefined;
    } & Record<Exclude<keyof I, keyof Acknowledgement>, never>>(object: I): Acknowledgement;
};
