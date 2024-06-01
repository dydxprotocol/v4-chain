import { Counterparty, Version } from "./connection";
import { Any } from "../../../../google/protobuf/any";
import { Height } from "../../client/v1/client";
import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "ibc.core.connection.v1";
/**
 * MsgConnectionOpenInit defines the msg sent by an account on Chain A to
 * initialize a connection with Chain B.
 */
export interface MsgConnectionOpenInit {
    clientId: string;
    counterparty: Counterparty;
    version?: Version;
    delayPeriod: bigint;
    signer: string;
}
/**
 * MsgConnectionOpenInitResponse defines the Msg/ConnectionOpenInit response
 * type.
 */
export interface MsgConnectionOpenInitResponse {
}
/**
 * MsgConnectionOpenTry defines a msg sent by a Relayer to try to open a
 * connection on Chain B.
 */
export interface MsgConnectionOpenTry {
    clientId: string;
    /** Deprecated: this field is unused. Crossing hellos are no longer supported in core IBC. */
    /** @deprecated */
    previousConnectionId: string;
    clientState?: Any;
    counterparty: Counterparty;
    delayPeriod: bigint;
    counterpartyVersions: Version[];
    proofHeight: Height;
    /**
     * proof of the initialization the connection on Chain A: `UNITIALIZED ->
     * INIT`
     */
    proofInit: Uint8Array;
    /** proof of client state included in message */
    proofClient: Uint8Array;
    /** proof of client consensus state */
    proofConsensus: Uint8Array;
    consensusHeight: Height;
    signer: string;
    /** optional proof data for host state machines that are unable to introspect their own consensus state */
    hostConsensusStateProof: Uint8Array;
}
/** MsgConnectionOpenTryResponse defines the Msg/ConnectionOpenTry response type. */
export interface MsgConnectionOpenTryResponse {
}
/**
 * MsgConnectionOpenAck defines a msg sent by a Relayer to Chain A to
 * acknowledge the change of connection state to TRYOPEN on Chain B.
 */
export interface MsgConnectionOpenAck {
    connectionId: string;
    counterpartyConnectionId: string;
    version?: Version;
    clientState?: Any;
    proofHeight: Height;
    /**
     * proof of the initialization the connection on Chain B: `UNITIALIZED ->
     * TRYOPEN`
     */
    proofTry: Uint8Array;
    /** proof of client state included in message */
    proofClient: Uint8Array;
    /** proof of client consensus state */
    proofConsensus: Uint8Array;
    consensusHeight: Height;
    signer: string;
    /** optional proof data for host state machines that are unable to introspect their own consensus state */
    hostConsensusStateProof: Uint8Array;
}
/** MsgConnectionOpenAckResponse defines the Msg/ConnectionOpenAck response type. */
export interface MsgConnectionOpenAckResponse {
}
/**
 * MsgConnectionOpenConfirm defines a msg sent by a Relayer to Chain B to
 * acknowledge the change of connection state to OPEN on Chain A.
 */
export interface MsgConnectionOpenConfirm {
    connectionId: string;
    /** proof for the change of the connection state on Chain A: `INIT -> OPEN` */
    proofAck: Uint8Array;
    proofHeight: Height;
    signer: string;
}
/**
 * MsgConnectionOpenConfirmResponse defines the Msg/ConnectionOpenConfirm
 * response type.
 */
export interface MsgConnectionOpenConfirmResponse {
}
export declare const MsgConnectionOpenInit: {
    typeUrl: string;
    encode(message: MsgConnectionOpenInit, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgConnectionOpenInit;
    fromJSON(object: any): MsgConnectionOpenInit;
    toJSON(message: MsgConnectionOpenInit): unknown;
    fromPartial<I extends {
        clientId?: string | undefined;
        counterparty?: {
            clientId?: string | undefined;
            connectionId?: string | undefined;
            prefix?: {
                keyPrefix?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        version?: {
            identifier?: string | undefined;
            features?: string[] | undefined;
        } | undefined;
        delayPeriod?: bigint | undefined;
        signer?: string | undefined;
    } & {
        clientId?: string | undefined;
        counterparty?: ({
            clientId?: string | undefined;
            connectionId?: string | undefined;
            prefix?: {
                keyPrefix?: Uint8Array | undefined;
            } | undefined;
        } & {
            clientId?: string | undefined;
            connectionId?: string | undefined;
            prefix?: ({
                keyPrefix?: Uint8Array | undefined;
            } & {
                keyPrefix?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["counterparty"]["prefix"], "keyPrefix">, never>) | undefined;
        } & Record<Exclude<keyof I["counterparty"], keyof Counterparty>, never>) | undefined;
        version?: ({
            identifier?: string | undefined;
            features?: string[] | undefined;
        } & {
            identifier?: string | undefined;
            features?: (string[] & string[] & Record<Exclude<keyof I["version"]["features"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["version"], keyof Version>, never>) | undefined;
        delayPeriod?: bigint | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgConnectionOpenInit>, never>>(object: I): MsgConnectionOpenInit;
};
export declare const MsgConnectionOpenInitResponse: {
    typeUrl: string;
    encode(_: MsgConnectionOpenInitResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgConnectionOpenInitResponse;
    fromJSON(_: any): MsgConnectionOpenInitResponse;
    toJSON(_: MsgConnectionOpenInitResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgConnectionOpenInitResponse;
};
export declare const MsgConnectionOpenTry: {
    typeUrl: string;
    encode(message: MsgConnectionOpenTry, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgConnectionOpenTry;
    fromJSON(object: any): MsgConnectionOpenTry;
    toJSON(message: MsgConnectionOpenTry): unknown;
    fromPartial<I extends {
        clientId?: string | undefined;
        previousConnectionId?: string | undefined;
        clientState?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        counterparty?: {
            clientId?: string | undefined;
            connectionId?: string | undefined;
            prefix?: {
                keyPrefix?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        delayPeriod?: bigint | undefined;
        counterpartyVersions?: {
            identifier?: string | undefined;
            features?: string[] | undefined;
        }[] | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        proofInit?: Uint8Array | undefined;
        proofClient?: Uint8Array | undefined;
        proofConsensus?: Uint8Array | undefined;
        consensusHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        signer?: string | undefined;
        hostConsensusStateProof?: Uint8Array | undefined;
    } & {
        clientId?: string | undefined;
        previousConnectionId?: string | undefined;
        clientState?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["clientState"], keyof Any>, never>) | undefined;
        counterparty?: ({
            clientId?: string | undefined;
            connectionId?: string | undefined;
            prefix?: {
                keyPrefix?: Uint8Array | undefined;
            } | undefined;
        } & {
            clientId?: string | undefined;
            connectionId?: string | undefined;
            prefix?: ({
                keyPrefix?: Uint8Array | undefined;
            } & {
                keyPrefix?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["counterparty"]["prefix"], "keyPrefix">, never>) | undefined;
        } & Record<Exclude<keyof I["counterparty"], keyof Counterparty>, never>) | undefined;
        delayPeriod?: bigint | undefined;
        counterpartyVersions?: ({
            identifier?: string | undefined;
            features?: string[] | undefined;
        }[] & ({
            identifier?: string | undefined;
            features?: string[] | undefined;
        } & {
            identifier?: string | undefined;
            features?: (string[] & string[] & Record<Exclude<keyof I["counterpartyVersions"][number]["features"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["counterpartyVersions"][number], keyof Version>, never>)[] & Record<Exclude<keyof I["counterpartyVersions"], keyof {
            identifier?: string | undefined;
            features?: string[] | undefined;
        }[]>, never>) | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
        proofInit?: Uint8Array | undefined;
        proofClient?: Uint8Array | undefined;
        proofConsensus?: Uint8Array | undefined;
        consensusHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["consensusHeight"], keyof Height>, never>) | undefined;
        signer?: string | undefined;
        hostConsensusStateProof?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof MsgConnectionOpenTry>, never>>(object: I): MsgConnectionOpenTry;
};
export declare const MsgConnectionOpenTryResponse: {
    typeUrl: string;
    encode(_: MsgConnectionOpenTryResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgConnectionOpenTryResponse;
    fromJSON(_: any): MsgConnectionOpenTryResponse;
    toJSON(_: MsgConnectionOpenTryResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgConnectionOpenTryResponse;
};
export declare const MsgConnectionOpenAck: {
    typeUrl: string;
    encode(message: MsgConnectionOpenAck, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgConnectionOpenAck;
    fromJSON(object: any): MsgConnectionOpenAck;
    toJSON(message: MsgConnectionOpenAck): unknown;
    fromPartial<I extends {
        connectionId?: string | undefined;
        counterpartyConnectionId?: string | undefined;
        version?: {
            identifier?: string | undefined;
            features?: string[] | undefined;
        } | undefined;
        clientState?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        proofTry?: Uint8Array | undefined;
        proofClient?: Uint8Array | undefined;
        proofConsensus?: Uint8Array | undefined;
        consensusHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        signer?: string | undefined;
        hostConsensusStateProof?: Uint8Array | undefined;
    } & {
        connectionId?: string | undefined;
        counterpartyConnectionId?: string | undefined;
        version?: ({
            identifier?: string | undefined;
            features?: string[] | undefined;
        } & {
            identifier?: string | undefined;
            features?: (string[] & string[] & Record<Exclude<keyof I["version"]["features"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["version"], keyof Version>, never>) | undefined;
        clientState?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["clientState"], keyof Any>, never>) | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
        proofTry?: Uint8Array | undefined;
        proofClient?: Uint8Array | undefined;
        proofConsensus?: Uint8Array | undefined;
        consensusHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["consensusHeight"], keyof Height>, never>) | undefined;
        signer?: string | undefined;
        hostConsensusStateProof?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof MsgConnectionOpenAck>, never>>(object: I): MsgConnectionOpenAck;
};
export declare const MsgConnectionOpenAckResponse: {
    typeUrl: string;
    encode(_: MsgConnectionOpenAckResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgConnectionOpenAckResponse;
    fromJSON(_: any): MsgConnectionOpenAckResponse;
    toJSON(_: MsgConnectionOpenAckResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgConnectionOpenAckResponse;
};
export declare const MsgConnectionOpenConfirm: {
    typeUrl: string;
    encode(message: MsgConnectionOpenConfirm, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgConnectionOpenConfirm;
    fromJSON(object: any): MsgConnectionOpenConfirm;
    toJSON(message: MsgConnectionOpenConfirm): unknown;
    fromPartial<I extends {
        connectionId?: string | undefined;
        proofAck?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
        signer?: string | undefined;
    } & {
        connectionId?: string | undefined;
        proofAck?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
        signer?: string | undefined;
    } & Record<Exclude<keyof I, keyof MsgConnectionOpenConfirm>, never>>(object: I): MsgConnectionOpenConfirm;
};
export declare const MsgConnectionOpenConfirmResponse: {
    typeUrl: string;
    encode(_: MsgConnectionOpenConfirmResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): MsgConnectionOpenConfirmResponse;
    fromJSON(_: any): MsgConnectionOpenConfirmResponse;
    toJSON(_: MsgConnectionOpenConfirmResponse): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): MsgConnectionOpenConfirmResponse;
};
/** Msg defines the ibc/connection Msg service. */
export interface Msg {
    /** ConnectionOpenInit defines a rpc handler method for MsgConnectionOpenInit. */
    ConnectionOpenInit(request: MsgConnectionOpenInit): Promise<MsgConnectionOpenInitResponse>;
    /** ConnectionOpenTry defines a rpc handler method for MsgConnectionOpenTry. */
    ConnectionOpenTry(request: MsgConnectionOpenTry): Promise<MsgConnectionOpenTryResponse>;
    /** ConnectionOpenAck defines a rpc handler method for MsgConnectionOpenAck. */
    ConnectionOpenAck(request: MsgConnectionOpenAck): Promise<MsgConnectionOpenAckResponse>;
    /**
     * ConnectionOpenConfirm defines a rpc handler method for
     * MsgConnectionOpenConfirm.
     */
    ConnectionOpenConfirm(request: MsgConnectionOpenConfirm): Promise<MsgConnectionOpenConfirmResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    ConnectionOpenInit(request: MsgConnectionOpenInit): Promise<MsgConnectionOpenInitResponse>;
    ConnectionOpenTry(request: MsgConnectionOpenTry): Promise<MsgConnectionOpenTryResponse>;
    ConnectionOpenAck(request: MsgConnectionOpenAck): Promise<MsgConnectionOpenAckResponse>;
    ConnectionOpenConfirm(request: MsgConnectionOpenConfirm): Promise<MsgConnectionOpenConfirmResponse>;
}
