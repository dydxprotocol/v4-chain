import { IdentifiedChannel, PacketState } from "./channel";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.core.channel.v1";
/** GenesisState defines the ibc channel submodule's genesis state. */
export interface GenesisState {
    channels: IdentifiedChannel[];
    acknowledgements: PacketState[];
    commitments: PacketState[];
    receipts: PacketState[];
    sendSequences: PacketSequence[];
    recvSequences: PacketSequence[];
    ackSequences: PacketSequence[];
    /** the sequence for the next generated channel identifier */
    nextChannelSequence: bigint;
}
/**
 * PacketSequence defines the genesis type necessary to retrieve and store
 * next send and receive sequences.
 */
export interface PacketSequence {
    portId: string;
    channelId: string;
    sequence: bigint;
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        channels?: {
            state?: import("./channel").State | undefined;
            ordering?: import("./channel").Order | undefined;
            counterparty?: {
                portId?: string | undefined;
                channelId?: string | undefined;
            } | undefined;
            connectionHops?: string[] | undefined;
            version?: string | undefined;
            portId?: string | undefined;
            channelId?: string | undefined;
        }[] | undefined;
        acknowledgements?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        }[] | undefined;
        commitments?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        }[] | undefined;
        receipts?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        }[] | undefined;
        sendSequences?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        }[] | undefined;
        recvSequences?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        }[] | undefined;
        ackSequences?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        }[] | undefined;
        nextChannelSequence?: bigint | undefined;
    } & {
        channels?: ({
            state?: import("./channel").State | undefined;
            ordering?: import("./channel").Order | undefined;
            counterparty?: {
                portId?: string | undefined;
                channelId?: string | undefined;
            } | undefined;
            connectionHops?: string[] | undefined;
            version?: string | undefined;
            portId?: string | undefined;
            channelId?: string | undefined;
        }[] & ({
            state?: import("./channel").State | undefined;
            ordering?: import("./channel").Order | undefined;
            counterparty?: {
                portId?: string | undefined;
                channelId?: string | undefined;
            } | undefined;
            connectionHops?: string[] | undefined;
            version?: string | undefined;
            portId?: string | undefined;
            channelId?: string | undefined;
        } & {
            state?: import("./channel").State | undefined;
            ordering?: import("./channel").Order | undefined;
            counterparty?: ({
                portId?: string | undefined;
                channelId?: string | undefined;
            } & {
                portId?: string | undefined;
                channelId?: string | undefined;
            } & Record<Exclude<keyof I["channels"][number]["counterparty"], keyof import("./channel").Counterparty>, never>) | undefined;
            connectionHops?: (string[] & string[] & Record<Exclude<keyof I["channels"][number]["connectionHops"], keyof string[]>, never>) | undefined;
            version?: string | undefined;
            portId?: string | undefined;
            channelId?: string | undefined;
        } & Record<Exclude<keyof I["channels"][number], keyof IdentifiedChannel>, never>)[] & Record<Exclude<keyof I["channels"], keyof {
            state?: import("./channel").State | undefined;
            ordering?: import("./channel").Order | undefined;
            counterparty?: {
                portId?: string | undefined;
                channelId?: string | undefined;
            } | undefined;
            connectionHops?: string[] | undefined;
            version?: string | undefined;
            portId?: string | undefined;
            channelId?: string | undefined;
        }[]>, never>) | undefined;
        acknowledgements?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        }[] & ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["acknowledgements"][number], keyof PacketState>, never>)[] & Record<Exclude<keyof I["acknowledgements"], keyof {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        }[]>, never>) | undefined;
        commitments?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        }[] & ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["commitments"][number], keyof PacketState>, never>)[] & Record<Exclude<keyof I["commitments"], keyof {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        }[]>, never>) | undefined;
        receipts?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        }[] & ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["receipts"][number], keyof PacketState>, never>)[] & Record<Exclude<keyof I["receipts"], keyof {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        }[]>, never>) | undefined;
        sendSequences?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        }[] & ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & Record<Exclude<keyof I["sendSequences"][number], keyof PacketSequence>, never>)[] & Record<Exclude<keyof I["sendSequences"], keyof {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        }[]>, never>) | undefined;
        recvSequences?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        }[] & ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & Record<Exclude<keyof I["recvSequences"][number], keyof PacketSequence>, never>)[] & Record<Exclude<keyof I["recvSequences"], keyof {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        }[]>, never>) | undefined;
        ackSequences?: ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        }[] & ({
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        } & Record<Exclude<keyof I["ackSequences"][number], keyof PacketSequence>, never>)[] & Record<Exclude<keyof I["ackSequences"], keyof {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
        }[]>, never>) | undefined;
        nextChannelSequence?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
export declare const PacketSequence: {
    typeUrl: string;
    encode(message: PacketSequence, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): PacketSequence;
    fromJSON(object: any): PacketSequence;
    toJSON(message: PacketSequence): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof PacketSequence>, never>>(object: I): PacketSequence;
};
