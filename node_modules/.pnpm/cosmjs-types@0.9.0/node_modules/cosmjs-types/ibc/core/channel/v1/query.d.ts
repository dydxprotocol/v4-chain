import { PageRequest, PageResponse } from "../../../../cosmos/base/query/v1beta1/pagination";
import { Channel, IdentifiedChannel, PacketState } from "./channel";
import { Height, IdentifiedClientState } from "../../client/v1/client";
import { Any } from "../../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "ibc.core.channel.v1";
/** QueryChannelRequest is the request type for the Query/Channel RPC method */
export interface QueryChannelRequest {
    /** port unique identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
}
/**
 * QueryChannelResponse is the response type for the Query/Channel RPC method.
 * Besides the Channel end, it includes a proof and the height from which the
 * proof was retrieved.
 */
export interface QueryChannelResponse {
    /** channel associated with the request identifiers */
    channel?: Channel;
    /** merkle proof of existence */
    proof: Uint8Array;
    /** height at which the proof was retrieved */
    proofHeight: Height;
}
/** QueryChannelsRequest is the request type for the Query/Channels RPC method */
export interface QueryChannelsRequest {
    /** pagination request */
    pagination?: PageRequest;
}
/** QueryChannelsResponse is the response type for the Query/Channels RPC method. */
export interface QueryChannelsResponse {
    /** list of stored channels of the chain. */
    channels: IdentifiedChannel[];
    /** pagination response */
    pagination?: PageResponse;
    /** query block height */
    height: Height;
}
/**
 * QueryConnectionChannelsRequest is the request type for the
 * Query/QueryConnectionChannels RPC method
 */
export interface QueryConnectionChannelsRequest {
    /** connection unique identifier */
    connection: string;
    /** pagination request */
    pagination?: PageRequest;
}
/**
 * QueryConnectionChannelsResponse is the Response type for the
 * Query/QueryConnectionChannels RPC method
 */
export interface QueryConnectionChannelsResponse {
    /** list of channels associated with a connection. */
    channels: IdentifiedChannel[];
    /** pagination response */
    pagination?: PageResponse;
    /** query block height */
    height: Height;
}
/**
 * QueryChannelClientStateRequest is the request type for the Query/ClientState
 * RPC method
 */
export interface QueryChannelClientStateRequest {
    /** port unique identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
}
/**
 * QueryChannelClientStateResponse is the Response type for the
 * Query/QueryChannelClientState RPC method
 */
export interface QueryChannelClientStateResponse {
    /** client state associated with the channel */
    identifiedClientState?: IdentifiedClientState;
    /** merkle proof of existence */
    proof: Uint8Array;
    /** height at which the proof was retrieved */
    proofHeight: Height;
}
/**
 * QueryChannelConsensusStateRequest is the request type for the
 * Query/ConsensusState RPC method
 */
export interface QueryChannelConsensusStateRequest {
    /** port unique identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
    /** revision number of the consensus state */
    revisionNumber: bigint;
    /** revision height of the consensus state */
    revisionHeight: bigint;
}
/**
 * QueryChannelClientStateResponse is the Response type for the
 * Query/QueryChannelClientState RPC method
 */
export interface QueryChannelConsensusStateResponse {
    /** consensus state associated with the channel */
    consensusState?: Any;
    /** client ID associated with the consensus state */
    clientId: string;
    /** merkle proof of existence */
    proof: Uint8Array;
    /** height at which the proof was retrieved */
    proofHeight: Height;
}
/**
 * QueryPacketCommitmentRequest is the request type for the
 * Query/PacketCommitment RPC method
 */
export interface QueryPacketCommitmentRequest {
    /** port unique identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
    /** packet sequence */
    sequence: bigint;
}
/**
 * QueryPacketCommitmentResponse defines the client query response for a packet
 * which also includes a proof and the height from which the proof was
 * retrieved
 */
export interface QueryPacketCommitmentResponse {
    /** packet associated with the request fields */
    commitment: Uint8Array;
    /** merkle proof of existence */
    proof: Uint8Array;
    /** height at which the proof was retrieved */
    proofHeight: Height;
}
/**
 * QueryPacketCommitmentsRequest is the request type for the
 * Query/QueryPacketCommitments RPC method
 */
export interface QueryPacketCommitmentsRequest {
    /** port unique identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
    /** pagination request */
    pagination?: PageRequest;
}
/**
 * QueryPacketCommitmentsResponse is the request type for the
 * Query/QueryPacketCommitments RPC method
 */
export interface QueryPacketCommitmentsResponse {
    commitments: PacketState[];
    /** pagination response */
    pagination?: PageResponse;
    /** query block height */
    height: Height;
}
/**
 * QueryPacketReceiptRequest is the request type for the
 * Query/PacketReceipt RPC method
 */
export interface QueryPacketReceiptRequest {
    /** port unique identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
    /** packet sequence */
    sequence: bigint;
}
/**
 * QueryPacketReceiptResponse defines the client query response for a packet
 * receipt which also includes a proof, and the height from which the proof was
 * retrieved
 */
export interface QueryPacketReceiptResponse {
    /** success flag for if receipt exists */
    received: boolean;
    /** merkle proof of existence */
    proof: Uint8Array;
    /** height at which the proof was retrieved */
    proofHeight: Height;
}
/**
 * QueryPacketAcknowledgementRequest is the request type for the
 * Query/PacketAcknowledgement RPC method
 */
export interface QueryPacketAcknowledgementRequest {
    /** port unique identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
    /** packet sequence */
    sequence: bigint;
}
/**
 * QueryPacketAcknowledgementResponse defines the client query response for a
 * packet which also includes a proof and the height from which the
 * proof was retrieved
 */
export interface QueryPacketAcknowledgementResponse {
    /** packet associated with the request fields */
    acknowledgement: Uint8Array;
    /** merkle proof of existence */
    proof: Uint8Array;
    /** height at which the proof was retrieved */
    proofHeight: Height;
}
/**
 * QueryPacketAcknowledgementsRequest is the request type for the
 * Query/QueryPacketCommitments RPC method
 */
export interface QueryPacketAcknowledgementsRequest {
    /** port unique identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
    /** pagination request */
    pagination?: PageRequest;
    /** list of packet sequences */
    packetCommitmentSequences: bigint[];
}
/**
 * QueryPacketAcknowledgemetsResponse is the request type for the
 * Query/QueryPacketAcknowledgements RPC method
 */
export interface QueryPacketAcknowledgementsResponse {
    acknowledgements: PacketState[];
    /** pagination response */
    pagination?: PageResponse;
    /** query block height */
    height: Height;
}
/**
 * QueryUnreceivedPacketsRequest is the request type for the
 * Query/UnreceivedPackets RPC method
 */
export interface QueryUnreceivedPacketsRequest {
    /** port unique identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
    /** list of packet sequences */
    packetCommitmentSequences: bigint[];
}
/**
 * QueryUnreceivedPacketsResponse is the response type for the
 * Query/UnreceivedPacketCommitments RPC method
 */
export interface QueryUnreceivedPacketsResponse {
    /** list of unreceived packet sequences */
    sequences: bigint[];
    /** query block height */
    height: Height;
}
/**
 * QueryUnreceivedAcks is the request type for the
 * Query/UnreceivedAcks RPC method
 */
export interface QueryUnreceivedAcksRequest {
    /** port unique identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
    /** list of acknowledgement sequences */
    packetAckSequences: bigint[];
}
/**
 * QueryUnreceivedAcksResponse is the response type for the
 * Query/UnreceivedAcks RPC method
 */
export interface QueryUnreceivedAcksResponse {
    /** list of unreceived acknowledgement sequences */
    sequences: bigint[];
    /** query block height */
    height: Height;
}
/**
 * QueryNextSequenceReceiveRequest is the request type for the
 * Query/QueryNextSequenceReceiveRequest RPC method
 */
export interface QueryNextSequenceReceiveRequest {
    /** port unique identifier */
    portId: string;
    /** channel unique identifier */
    channelId: string;
}
/**
 * QuerySequenceResponse is the request type for the
 * Query/QueryNextSequenceReceiveResponse RPC method
 */
export interface QueryNextSequenceReceiveResponse {
    /** next sequence receive number */
    nextSequenceReceive: bigint;
    /** merkle proof of existence */
    proof: Uint8Array;
    /** height at which the proof was retrieved */
    proofHeight: Height;
}
export declare const QueryChannelRequest: {
    typeUrl: string;
    encode(message: QueryChannelRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryChannelRequest;
    fromJSON(object: any): QueryChannelRequest;
    toJSON(message: QueryChannelRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryChannelRequest>, never>>(object: I): QueryChannelRequest;
};
export declare const QueryChannelResponse: {
    typeUrl: string;
    encode(message: QueryChannelResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryChannelResponse;
    fromJSON(object: any): QueryChannelResponse;
    toJSON(message: QueryChannelResponse): unknown;
    fromPartial<I extends {
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
        proof?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
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
        proof?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryChannelResponse>, never>>(object: I): QueryChannelResponse;
};
export declare const QueryChannelsRequest: {
    typeUrl: string;
    encode(message: QueryChannelsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryChannelsRequest;
    fromJSON(object: any): QueryChannelsRequest;
    toJSON(message: QueryChannelsRequest): unknown;
    fromPartial<I extends {
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
    } & Record<Exclude<keyof I, "pagination">, never>>(object: I): QueryChannelsRequest;
};
export declare const QueryChannelsResponse: {
    typeUrl: string;
    encode(message: QueryChannelsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryChannelsResponse;
    fromJSON(object: any): QueryChannelsResponse;
    toJSON(message: QueryChannelsResponse): unknown;
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
        height?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
        height?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["height"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryChannelsResponse>, never>>(object: I): QueryChannelsResponse;
};
export declare const QueryConnectionChannelsRequest: {
    typeUrl: string;
    encode(message: QueryConnectionChannelsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionChannelsRequest;
    fromJSON(object: any): QueryConnectionChannelsRequest;
    toJSON(message: QueryConnectionChannelsRequest): unknown;
    fromPartial<I extends {
        connection?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        connection?: string | undefined;
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryConnectionChannelsRequest>, never>>(object: I): QueryConnectionChannelsRequest;
};
export declare const QueryConnectionChannelsResponse: {
    typeUrl: string;
    encode(message: QueryConnectionChannelsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionChannelsResponse;
    fromJSON(object: any): QueryConnectionChannelsResponse;
    toJSON(message: QueryConnectionChannelsResponse): unknown;
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
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
        height?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
        height?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["height"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryConnectionChannelsResponse>, never>>(object: I): QueryConnectionChannelsResponse;
};
export declare const QueryChannelClientStateRequest: {
    typeUrl: string;
    encode(message: QueryChannelClientStateRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryChannelClientStateRequest;
    fromJSON(object: any): QueryChannelClientStateRequest;
    toJSON(message: QueryChannelClientStateRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryChannelClientStateRequest>, never>>(object: I): QueryChannelClientStateRequest;
};
export declare const QueryChannelClientStateResponse: {
    typeUrl: string;
    encode(message: QueryChannelClientStateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryChannelClientStateResponse;
    fromJSON(object: any): QueryChannelClientStateResponse;
    toJSON(message: QueryChannelClientStateResponse): unknown;
    fromPartial<I extends {
        identifiedClientState?: {
            clientId?: string | undefined;
            clientState?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
        identifiedClientState?: ({
            clientId?: string | undefined;
            clientState?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        } & {
            clientId?: string | undefined;
            clientState?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["identifiedClientState"]["clientState"], keyof Any>, never>) | undefined;
        } & Record<Exclude<keyof I["identifiedClientState"], keyof IdentifiedClientState>, never>) | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryChannelClientStateResponse>, never>>(object: I): QueryChannelClientStateResponse;
};
export declare const QueryChannelConsensusStateRequest: {
    typeUrl: string;
    encode(message: QueryChannelConsensusStateRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryChannelConsensusStateRequest;
    fromJSON(object: any): QueryChannelConsensusStateRequest;
    toJSON(message: QueryChannelConsensusStateRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        revisionNumber?: bigint | undefined;
        revisionHeight?: bigint | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        revisionNumber?: bigint | undefined;
        revisionHeight?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof QueryChannelConsensusStateRequest>, never>>(object: I): QueryChannelConsensusStateRequest;
};
export declare const QueryChannelConsensusStateResponse: {
    typeUrl: string;
    encode(message: QueryChannelConsensusStateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryChannelConsensusStateResponse;
    fromJSON(object: any): QueryChannelConsensusStateResponse;
    toJSON(message: QueryChannelConsensusStateResponse): unknown;
    fromPartial<I extends {
        consensusState?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        clientId?: string | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
        consensusState?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["consensusState"], keyof Any>, never>) | undefined;
        clientId?: string | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryChannelConsensusStateResponse>, never>>(object: I): QueryChannelConsensusStateResponse;
};
export declare const QueryPacketCommitmentRequest: {
    typeUrl: string;
    encode(message: QueryPacketCommitmentRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPacketCommitmentRequest;
    fromJSON(object: any): QueryPacketCommitmentRequest;
    toJSON(message: QueryPacketCommitmentRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof QueryPacketCommitmentRequest>, never>>(object: I): QueryPacketCommitmentRequest;
};
export declare const QueryPacketCommitmentResponse: {
    typeUrl: string;
    encode(message: QueryPacketCommitmentResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPacketCommitmentResponse;
    fromJSON(object: any): QueryPacketCommitmentResponse;
    toJSON(message: QueryPacketCommitmentResponse): unknown;
    fromPartial<I extends {
        commitment?: Uint8Array | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
        commitment?: Uint8Array | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryPacketCommitmentResponse>, never>>(object: I): QueryPacketCommitmentResponse;
};
export declare const QueryPacketCommitmentsRequest: {
    typeUrl: string;
    encode(message: QueryPacketCommitmentsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPacketCommitmentsRequest;
    fromJSON(object: any): QueryPacketCommitmentsRequest;
    toJSON(message: QueryPacketCommitmentsRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryPacketCommitmentsRequest>, never>>(object: I): QueryPacketCommitmentsRequest;
};
export declare const QueryPacketCommitmentsResponse: {
    typeUrl: string;
    encode(message: QueryPacketCommitmentsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPacketCommitmentsResponse;
    fromJSON(object: any): QueryPacketCommitmentsResponse;
    toJSON(message: QueryPacketCommitmentsResponse): unknown;
    fromPartial<I extends {
        commitments?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
        height?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
        height?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["height"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryPacketCommitmentsResponse>, never>>(object: I): QueryPacketCommitmentsResponse;
};
export declare const QueryPacketReceiptRequest: {
    typeUrl: string;
    encode(message: QueryPacketReceiptRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPacketReceiptRequest;
    fromJSON(object: any): QueryPacketReceiptRequest;
    toJSON(message: QueryPacketReceiptRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof QueryPacketReceiptRequest>, never>>(object: I): QueryPacketReceiptRequest;
};
export declare const QueryPacketReceiptResponse: {
    typeUrl: string;
    encode(message: QueryPacketReceiptResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPacketReceiptResponse;
    fromJSON(object: any): QueryPacketReceiptResponse;
    toJSON(message: QueryPacketReceiptResponse): unknown;
    fromPartial<I extends {
        received?: boolean | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
        received?: boolean | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryPacketReceiptResponse>, never>>(object: I): QueryPacketReceiptResponse;
};
export declare const QueryPacketAcknowledgementRequest: {
    typeUrl: string;
    encode(message: QueryPacketAcknowledgementRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPacketAcknowledgementRequest;
    fromJSON(object: any): QueryPacketAcknowledgementRequest;
    toJSON(message: QueryPacketAcknowledgementRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        sequence?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof QueryPacketAcknowledgementRequest>, never>>(object: I): QueryPacketAcknowledgementRequest;
};
export declare const QueryPacketAcknowledgementResponse: {
    typeUrl: string;
    encode(message: QueryPacketAcknowledgementResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPacketAcknowledgementResponse;
    fromJSON(object: any): QueryPacketAcknowledgementResponse;
    toJSON(message: QueryPacketAcknowledgementResponse): unknown;
    fromPartial<I extends {
        acknowledgement?: Uint8Array | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
        acknowledgement?: Uint8Array | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryPacketAcknowledgementResponse>, never>>(object: I): QueryPacketAcknowledgementResponse;
};
export declare const QueryPacketAcknowledgementsRequest: {
    typeUrl: string;
    encode(message: QueryPacketAcknowledgementsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPacketAcknowledgementsRequest;
    fromJSON(object: any): QueryPacketAcknowledgementsRequest;
    toJSON(message: QueryPacketAcknowledgementsRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        pagination?: {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } | undefined;
        packetCommitmentSequences?: bigint[] | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        pagination?: ({
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & {
            key?: Uint8Array | undefined;
            offset?: bigint | undefined;
            limit?: bigint | undefined;
            countTotal?: boolean | undefined;
            reverse?: boolean | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageRequest>, never>) | undefined;
        packetCommitmentSequences?: (bigint[] & bigint[] & Record<Exclude<keyof I["packetCommitmentSequences"], keyof bigint[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryPacketAcknowledgementsRequest>, never>>(object: I): QueryPacketAcknowledgementsRequest;
};
export declare const QueryPacketAcknowledgementsResponse: {
    typeUrl: string;
    encode(message: QueryPacketAcknowledgementsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryPacketAcknowledgementsResponse;
    fromJSON(object: any): QueryPacketAcknowledgementsResponse;
    toJSON(message: QueryPacketAcknowledgementsResponse): unknown;
    fromPartial<I extends {
        acknowledgements?: {
            portId?: string | undefined;
            channelId?: string | undefined;
            sequence?: bigint | undefined;
            data?: Uint8Array | undefined;
        }[] | undefined;
        pagination?: {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } | undefined;
        height?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
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
        pagination?: ({
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & {
            nextKey?: Uint8Array | undefined;
            total?: bigint | undefined;
        } & Record<Exclude<keyof I["pagination"], keyof PageResponse>, never>) | undefined;
        height?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["height"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryPacketAcknowledgementsResponse>, never>>(object: I): QueryPacketAcknowledgementsResponse;
};
export declare const QueryUnreceivedPacketsRequest: {
    typeUrl: string;
    encode(message: QueryUnreceivedPacketsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryUnreceivedPacketsRequest;
    fromJSON(object: any): QueryUnreceivedPacketsRequest;
    toJSON(message: QueryUnreceivedPacketsRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        packetCommitmentSequences?: bigint[] | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        packetCommitmentSequences?: (bigint[] & bigint[] & Record<Exclude<keyof I["packetCommitmentSequences"], keyof bigint[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryUnreceivedPacketsRequest>, never>>(object: I): QueryUnreceivedPacketsRequest;
};
export declare const QueryUnreceivedPacketsResponse: {
    typeUrl: string;
    encode(message: QueryUnreceivedPacketsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryUnreceivedPacketsResponse;
    fromJSON(object: any): QueryUnreceivedPacketsResponse;
    toJSON(message: QueryUnreceivedPacketsResponse): unknown;
    fromPartial<I extends {
        sequences?: bigint[] | undefined;
        height?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
        sequences?: (bigint[] & bigint[] & Record<Exclude<keyof I["sequences"], keyof bigint[]>, never>) | undefined;
        height?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["height"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryUnreceivedPacketsResponse>, never>>(object: I): QueryUnreceivedPacketsResponse;
};
export declare const QueryUnreceivedAcksRequest: {
    typeUrl: string;
    encode(message: QueryUnreceivedAcksRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryUnreceivedAcksRequest;
    fromJSON(object: any): QueryUnreceivedAcksRequest;
    toJSON(message: QueryUnreceivedAcksRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
        packetAckSequences?: bigint[] | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
        packetAckSequences?: (bigint[] & bigint[] & Record<Exclude<keyof I["packetAckSequences"], keyof bigint[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryUnreceivedAcksRequest>, never>>(object: I): QueryUnreceivedAcksRequest;
};
export declare const QueryUnreceivedAcksResponse: {
    typeUrl: string;
    encode(message: QueryUnreceivedAcksResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryUnreceivedAcksResponse;
    fromJSON(object: any): QueryUnreceivedAcksResponse;
    toJSON(message: QueryUnreceivedAcksResponse): unknown;
    fromPartial<I extends {
        sequences?: bigint[] | undefined;
        height?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
        sequences?: (bigint[] & bigint[] & Record<Exclude<keyof I["sequences"], keyof bigint[]>, never>) | undefined;
        height?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["height"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryUnreceivedAcksResponse>, never>>(object: I): QueryUnreceivedAcksResponse;
};
export declare const QueryNextSequenceReceiveRequest: {
    typeUrl: string;
    encode(message: QueryNextSequenceReceiveRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryNextSequenceReceiveRequest;
    fromJSON(object: any): QueryNextSequenceReceiveRequest;
    toJSON(message: QueryNextSequenceReceiveRequest): unknown;
    fromPartial<I extends {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & {
        portId?: string | undefined;
        channelId?: string | undefined;
    } & Record<Exclude<keyof I, keyof QueryNextSequenceReceiveRequest>, never>>(object: I): QueryNextSequenceReceiveRequest;
};
export declare const QueryNextSequenceReceiveResponse: {
    typeUrl: string;
    encode(message: QueryNextSequenceReceiveResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryNextSequenceReceiveResponse;
    fromJSON(object: any): QueryNextSequenceReceiveResponse;
    toJSON(message: QueryNextSequenceReceiveResponse): unknown;
    fromPartial<I extends {
        nextSequenceReceive?: bigint | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
        nextSequenceReceive?: bigint | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryNextSequenceReceiveResponse>, never>>(object: I): QueryNextSequenceReceiveResponse;
};
/** Query provides defines the gRPC querier service */
export interface Query {
    /** Channel queries an IBC Channel. */
    Channel(request: QueryChannelRequest): Promise<QueryChannelResponse>;
    /** Channels queries all the IBC channels of a chain. */
    Channels(request?: QueryChannelsRequest): Promise<QueryChannelsResponse>;
    /**
     * ConnectionChannels queries all the channels associated with a connection
     * end.
     */
    ConnectionChannels(request: QueryConnectionChannelsRequest): Promise<QueryConnectionChannelsResponse>;
    /**
     * ChannelClientState queries for the client state for the channel associated
     * with the provided channel identifiers.
     */
    ChannelClientState(request: QueryChannelClientStateRequest): Promise<QueryChannelClientStateResponse>;
    /**
     * ChannelConsensusState queries for the consensus state for the channel
     * associated with the provided channel identifiers.
     */
    ChannelConsensusState(request: QueryChannelConsensusStateRequest): Promise<QueryChannelConsensusStateResponse>;
    /** PacketCommitment queries a stored packet commitment hash. */
    PacketCommitment(request: QueryPacketCommitmentRequest): Promise<QueryPacketCommitmentResponse>;
    /**
     * PacketCommitments returns all the packet commitments hashes associated
     * with a channel.
     */
    PacketCommitments(request: QueryPacketCommitmentsRequest): Promise<QueryPacketCommitmentsResponse>;
    /**
     * PacketReceipt queries if a given packet sequence has been received on the
     * queried chain
     */
    PacketReceipt(request: QueryPacketReceiptRequest): Promise<QueryPacketReceiptResponse>;
    /** PacketAcknowledgement queries a stored packet acknowledgement hash. */
    PacketAcknowledgement(request: QueryPacketAcknowledgementRequest): Promise<QueryPacketAcknowledgementResponse>;
    /**
     * PacketAcknowledgements returns all the packet acknowledgements associated
     * with a channel.
     */
    PacketAcknowledgements(request: QueryPacketAcknowledgementsRequest): Promise<QueryPacketAcknowledgementsResponse>;
    /**
     * UnreceivedPackets returns all the unreceived IBC packets associated with a
     * channel and sequences.
     */
    UnreceivedPackets(request: QueryUnreceivedPacketsRequest): Promise<QueryUnreceivedPacketsResponse>;
    /**
     * UnreceivedAcks returns all the unreceived IBC acknowledgements associated
     * with a channel and sequences.
     */
    UnreceivedAcks(request: QueryUnreceivedAcksRequest): Promise<QueryUnreceivedAcksResponse>;
    /** NextSequenceReceive returns the next receive sequence for a given channel. */
    NextSequenceReceive(request: QueryNextSequenceReceiveRequest): Promise<QueryNextSequenceReceiveResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Channel(request: QueryChannelRequest): Promise<QueryChannelResponse>;
    Channels(request?: QueryChannelsRequest): Promise<QueryChannelsResponse>;
    ConnectionChannels(request: QueryConnectionChannelsRequest): Promise<QueryConnectionChannelsResponse>;
    ChannelClientState(request: QueryChannelClientStateRequest): Promise<QueryChannelClientStateResponse>;
    ChannelConsensusState(request: QueryChannelConsensusStateRequest): Promise<QueryChannelConsensusStateResponse>;
    PacketCommitment(request: QueryPacketCommitmentRequest): Promise<QueryPacketCommitmentResponse>;
    PacketCommitments(request: QueryPacketCommitmentsRequest): Promise<QueryPacketCommitmentsResponse>;
    PacketReceipt(request: QueryPacketReceiptRequest): Promise<QueryPacketReceiptResponse>;
    PacketAcknowledgement(request: QueryPacketAcknowledgementRequest): Promise<QueryPacketAcknowledgementResponse>;
    PacketAcknowledgements(request: QueryPacketAcknowledgementsRequest): Promise<QueryPacketAcknowledgementsResponse>;
    UnreceivedPackets(request: QueryUnreceivedPacketsRequest): Promise<QueryUnreceivedPacketsResponse>;
    UnreceivedAcks(request: QueryUnreceivedAcksRequest): Promise<QueryUnreceivedAcksResponse>;
    NextSequenceReceive(request: QueryNextSequenceReceiveRequest): Promise<QueryNextSequenceReceiveResponse>;
}
