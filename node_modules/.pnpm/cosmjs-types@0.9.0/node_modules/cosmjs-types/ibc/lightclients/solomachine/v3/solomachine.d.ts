import { Any } from "../../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.lightclients.solomachine.v3";
/**
 * ClientState defines a solo machine client that tracks the current consensus
 * state and if the client is frozen.
 */
export interface ClientState {
    /** latest sequence of the client state */
    sequence: bigint;
    /** frozen sequence of the solo machine */
    isFrozen: boolean;
    consensusState?: ConsensusState;
}
/**
 * ConsensusState defines a solo machine consensus state. The sequence of a
 * consensus state is contained in the "height" key used in storing the
 * consensus state.
 */
export interface ConsensusState {
    /** public key of the solo machine */
    publicKey?: Any;
    /**
     * diversifier allows the same public key to be re-used across different solo
     * machine clients (potentially on different chains) without being considered
     * misbehaviour.
     */
    diversifier: string;
    timestamp: bigint;
}
/** Header defines a solo machine consensus header */
export interface Header {
    timestamp: bigint;
    signature: Uint8Array;
    newPublicKey?: Any;
    newDiversifier: string;
}
/**
 * Misbehaviour defines misbehaviour for a solo machine which consists
 * of a sequence and two signatures over different messages at that sequence.
 */
export interface Misbehaviour {
    sequence: bigint;
    signatureOne?: SignatureAndData;
    signatureTwo?: SignatureAndData;
}
/**
 * SignatureAndData contains a signature and the data signed over to create that
 * signature.
 */
export interface SignatureAndData {
    signature: Uint8Array;
    path: Uint8Array;
    data: Uint8Array;
    timestamp: bigint;
}
/**
 * TimestampedSignatureData contains the signature data and the timestamp of the
 * signature.
 */
export interface TimestampedSignatureData {
    signatureData: Uint8Array;
    timestamp: bigint;
}
/** SignBytes defines the signed bytes used for signature verification. */
export interface SignBytes {
    /** the sequence number */
    sequence: bigint;
    /** the proof timestamp */
    timestamp: bigint;
    /** the public key diversifier */
    diversifier: string;
    /** the standardised path bytes */
    path: Uint8Array;
    /** the marshaled data bytes */
    data: Uint8Array;
}
/** HeaderData returns the SignBytes data for update verification. */
export interface HeaderData {
    /** header public key */
    newPubKey?: Any;
    /** header diversifier */
    newDiversifier: string;
}
export declare const ClientState: {
    typeUrl: string;
    encode(message: ClientState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ClientState;
    fromJSON(object: any): ClientState;
    toJSON(message: ClientState): unknown;
    fromPartial<I extends {
        sequence?: bigint | undefined;
        isFrozen?: boolean | undefined;
        consensusState?: {
            publicKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            diversifier?: string | undefined;
            timestamp?: bigint | undefined;
        } | undefined;
    } & {
        sequence?: bigint | undefined;
        isFrozen?: boolean | undefined;
        consensusState?: ({
            publicKey?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
            diversifier?: string | undefined;
            timestamp?: bigint | undefined;
        } & {
            publicKey?: ({
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["consensusState"]["publicKey"], keyof Any>, never>) | undefined;
            diversifier?: string | undefined;
            timestamp?: bigint | undefined;
        } & Record<Exclude<keyof I["consensusState"], keyof ConsensusState>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ClientState>, never>>(object: I): ClientState;
};
export declare const ConsensusState: {
    typeUrl: string;
    encode(message: ConsensusState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ConsensusState;
    fromJSON(object: any): ConsensusState;
    toJSON(message: ConsensusState): unknown;
    fromPartial<I extends {
        publicKey?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        diversifier?: string | undefined;
        timestamp?: bigint | undefined;
    } & {
        publicKey?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["publicKey"], keyof Any>, never>) | undefined;
        diversifier?: string | undefined;
        timestamp?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof ConsensusState>, never>>(object: I): ConsensusState;
};
export declare const Header: {
    typeUrl: string;
    encode(message: Header, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Header;
    fromJSON(object: any): Header;
    toJSON(message: Header): unknown;
    fromPartial<I extends {
        timestamp?: bigint | undefined;
        signature?: Uint8Array | undefined;
        newPublicKey?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        newDiversifier?: string | undefined;
    } & {
        timestamp?: bigint | undefined;
        signature?: Uint8Array | undefined;
        newPublicKey?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["newPublicKey"], keyof Any>, never>) | undefined;
        newDiversifier?: string | undefined;
    } & Record<Exclude<keyof I, keyof Header>, never>>(object: I): Header;
};
export declare const Misbehaviour: {
    typeUrl: string;
    encode(message: Misbehaviour, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Misbehaviour;
    fromJSON(object: any): Misbehaviour;
    toJSON(message: Misbehaviour): unknown;
    fromPartial<I extends {
        sequence?: bigint | undefined;
        signatureOne?: {
            signature?: Uint8Array | undefined;
            path?: Uint8Array | undefined;
            data?: Uint8Array | undefined;
            timestamp?: bigint | undefined;
        } | undefined;
        signatureTwo?: {
            signature?: Uint8Array | undefined;
            path?: Uint8Array | undefined;
            data?: Uint8Array | undefined;
            timestamp?: bigint | undefined;
        } | undefined;
    } & {
        sequence?: bigint | undefined;
        signatureOne?: ({
            signature?: Uint8Array | undefined;
            path?: Uint8Array | undefined;
            data?: Uint8Array | undefined;
            timestamp?: bigint | undefined;
        } & {
            signature?: Uint8Array | undefined;
            path?: Uint8Array | undefined;
            data?: Uint8Array | undefined;
            timestamp?: bigint | undefined;
        } & Record<Exclude<keyof I["signatureOne"], keyof SignatureAndData>, never>) | undefined;
        signatureTwo?: ({
            signature?: Uint8Array | undefined;
            path?: Uint8Array | undefined;
            data?: Uint8Array | undefined;
            timestamp?: bigint | undefined;
        } & {
            signature?: Uint8Array | undefined;
            path?: Uint8Array | undefined;
            data?: Uint8Array | undefined;
            timestamp?: bigint | undefined;
        } & Record<Exclude<keyof I["signatureTwo"], keyof SignatureAndData>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Misbehaviour>, never>>(object: I): Misbehaviour;
};
export declare const SignatureAndData: {
    typeUrl: string;
    encode(message: SignatureAndData, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SignatureAndData;
    fromJSON(object: any): SignatureAndData;
    toJSON(message: SignatureAndData): unknown;
    fromPartial<I extends {
        signature?: Uint8Array | undefined;
        path?: Uint8Array | undefined;
        data?: Uint8Array | undefined;
        timestamp?: bigint | undefined;
    } & {
        signature?: Uint8Array | undefined;
        path?: Uint8Array | undefined;
        data?: Uint8Array | undefined;
        timestamp?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof SignatureAndData>, never>>(object: I): SignatureAndData;
};
export declare const TimestampedSignatureData: {
    typeUrl: string;
    encode(message: TimestampedSignatureData, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): TimestampedSignatureData;
    fromJSON(object: any): TimestampedSignatureData;
    toJSON(message: TimestampedSignatureData): unknown;
    fromPartial<I extends {
        signatureData?: Uint8Array | undefined;
        timestamp?: bigint | undefined;
    } & {
        signatureData?: Uint8Array | undefined;
        timestamp?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof TimestampedSignatureData>, never>>(object: I): TimestampedSignatureData;
};
export declare const SignBytes: {
    typeUrl: string;
    encode(message: SignBytes, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): SignBytes;
    fromJSON(object: any): SignBytes;
    toJSON(message: SignBytes): unknown;
    fromPartial<I extends {
        sequence?: bigint | undefined;
        timestamp?: bigint | undefined;
        diversifier?: string | undefined;
        path?: Uint8Array | undefined;
        data?: Uint8Array | undefined;
    } & {
        sequence?: bigint | undefined;
        timestamp?: bigint | undefined;
        diversifier?: string | undefined;
        path?: Uint8Array | undefined;
        data?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof SignBytes>, never>>(object: I): SignBytes;
};
export declare const HeaderData: {
    typeUrl: string;
    encode(message: HeaderData, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): HeaderData;
    fromJSON(object: any): HeaderData;
    toJSON(message: HeaderData): unknown;
    fromPartial<I extends {
        newPubKey?: {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } | undefined;
        newDiversifier?: string | undefined;
    } & {
        newPubKey?: ({
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & {
            typeUrl?: string | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["newPubKey"], keyof Any>, never>) | undefined;
        newDiversifier?: string | undefined;
    } & Record<Exclude<keyof I, keyof HeaderData>, never>>(object: I): HeaderData;
};
