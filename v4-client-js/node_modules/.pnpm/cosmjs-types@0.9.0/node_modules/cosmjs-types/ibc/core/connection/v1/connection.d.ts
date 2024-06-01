import { MerklePrefix } from "../../commitment/v1/commitment";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.core.connection.v1";
/**
 * State defines if a connection is in one of the following states:
 * INIT, TRYOPEN, OPEN or UNINITIALIZED.
 */
export declare enum State {
    /** STATE_UNINITIALIZED_UNSPECIFIED - Default State */
    STATE_UNINITIALIZED_UNSPECIFIED = 0,
    /** STATE_INIT - A connection end has just started the opening handshake. */
    STATE_INIT = 1,
    /**
     * STATE_TRYOPEN - A connection end has acknowledged the handshake step on the counterparty
     * chain.
     */
    STATE_TRYOPEN = 2,
    /** STATE_OPEN - A connection end has completed the handshake. */
    STATE_OPEN = 3,
    UNRECOGNIZED = -1
}
export declare function stateFromJSON(object: any): State;
export declare function stateToJSON(object: State): string;
/**
 * ConnectionEnd defines a stateful object on a chain connected to another
 * separate one.
 * NOTE: there must only be 2 defined ConnectionEnds to establish
 * a connection between two chains.
 */
export interface ConnectionEnd {
    /** client associated with this connection. */
    clientId: string;
    /**
     * IBC version which can be utilised to determine encodings or protocols for
     * channels or packets utilising this connection.
     */
    versions: Version[];
    /** current state of the connection end. */
    state: State;
    /** counterparty chain associated with this connection. */
    counterparty: Counterparty;
    /**
     * delay period that must pass before a consensus state can be used for
     * packet-verification NOTE: delay period logic is only implemented by some
     * clients.
     */
    delayPeriod: bigint;
}
/**
 * IdentifiedConnection defines a connection with additional connection
 * identifier field.
 */
export interface IdentifiedConnection {
    /** connection identifier. */
    id: string;
    /** client associated with this connection. */
    clientId: string;
    /**
     * IBC version which can be utilised to determine encodings or protocols for
     * channels or packets utilising this connection
     */
    versions: Version[];
    /** current state of the connection end. */
    state: State;
    /** counterparty chain associated with this connection. */
    counterparty: Counterparty;
    /** delay period associated with this connection. */
    delayPeriod: bigint;
}
/** Counterparty defines the counterparty chain associated with a connection end. */
export interface Counterparty {
    /**
     * identifies the client on the counterparty chain associated with a given
     * connection.
     */
    clientId: string;
    /**
     * identifies the connection end on the counterparty chain associated with a
     * given connection.
     */
    connectionId: string;
    /** commitment merkle prefix of the counterparty chain. */
    prefix: MerklePrefix;
}
/** ClientPaths define all the connection paths for a client state. */
export interface ClientPaths {
    /** list of connection paths */
    paths: string[];
}
/** ConnectionPaths define all the connection paths for a given client state. */
export interface ConnectionPaths {
    /** client state unique identifier */
    clientId: string;
    /** list of connection paths */
    paths: string[];
}
/**
 * Version defines the versioning scheme used to negotiate the IBC verison in
 * the connection handshake.
 */
export interface Version {
    /** unique version identifier */
    identifier: string;
    /** list of features compatible with the specified identifier */
    features: string[];
}
/** Params defines the set of Connection parameters. */
export interface Params {
    /**
     * maximum expected time per block (in nanoseconds), used to enforce block delay. This parameter should reflect the
     * largest amount of time that the chain might reasonably take to produce the next block under normal operating
     * conditions. A safe choice is 3-5x the expected time per block.
     */
    maxExpectedTimePerBlock: bigint;
}
export declare const ConnectionEnd: {
    typeUrl: string;
    encode(message: ConnectionEnd, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ConnectionEnd;
    fromJSON(object: any): ConnectionEnd;
    toJSON(message: ConnectionEnd): unknown;
    fromPartial<I extends {
        clientId?: string | undefined;
        versions?: {
            identifier?: string | undefined;
            features?: string[] | undefined;
        }[] | undefined;
        state?: State | undefined;
        counterparty?: {
            clientId?: string | undefined;
            connectionId?: string | undefined;
            prefix?: {
                keyPrefix?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        delayPeriod?: bigint | undefined;
    } & {
        clientId?: string | undefined;
        versions?: ({
            identifier?: string | undefined;
            features?: string[] | undefined;
        }[] & ({
            identifier?: string | undefined;
            features?: string[] | undefined;
        } & {
            identifier?: string | undefined;
            features?: (string[] & string[] & Record<Exclude<keyof I["versions"][number]["features"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["versions"][number], keyof Version>, never>)[] & Record<Exclude<keyof I["versions"], keyof {
            identifier?: string | undefined;
            features?: string[] | undefined;
        }[]>, never>) | undefined;
        state?: State | undefined;
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
    } & Record<Exclude<keyof I, keyof ConnectionEnd>, never>>(object: I): ConnectionEnd;
};
export declare const IdentifiedConnection: {
    typeUrl: string;
    encode(message: IdentifiedConnection, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): IdentifiedConnection;
    fromJSON(object: any): IdentifiedConnection;
    toJSON(message: IdentifiedConnection): unknown;
    fromPartial<I extends {
        id?: string | undefined;
        clientId?: string | undefined;
        versions?: {
            identifier?: string | undefined;
            features?: string[] | undefined;
        }[] | undefined;
        state?: State | undefined;
        counterparty?: {
            clientId?: string | undefined;
            connectionId?: string | undefined;
            prefix?: {
                keyPrefix?: Uint8Array | undefined;
            } | undefined;
        } | undefined;
        delayPeriod?: bigint | undefined;
    } & {
        id?: string | undefined;
        clientId?: string | undefined;
        versions?: ({
            identifier?: string | undefined;
            features?: string[] | undefined;
        }[] & ({
            identifier?: string | undefined;
            features?: string[] | undefined;
        } & {
            identifier?: string | undefined;
            features?: (string[] & string[] & Record<Exclude<keyof I["versions"][number]["features"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["versions"][number], keyof Version>, never>)[] & Record<Exclude<keyof I["versions"], keyof {
            identifier?: string | undefined;
            features?: string[] | undefined;
        }[]>, never>) | undefined;
        state?: State | undefined;
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
    } & Record<Exclude<keyof I, keyof IdentifiedConnection>, never>>(object: I): IdentifiedConnection;
};
export declare const Counterparty: {
    typeUrl: string;
    encode(message: Counterparty, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Counterparty;
    fromJSON(object: any): Counterparty;
    toJSON(message: Counterparty): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["prefix"], "keyPrefix">, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Counterparty>, never>>(object: I): Counterparty;
};
export declare const ClientPaths: {
    typeUrl: string;
    encode(message: ClientPaths, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ClientPaths;
    fromJSON(object: any): ClientPaths;
    toJSON(message: ClientPaths): unknown;
    fromPartial<I extends {
        paths?: string[] | undefined;
    } & {
        paths?: (string[] & string[] & Record<Exclude<keyof I["paths"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "paths">, never>>(object: I): ClientPaths;
};
export declare const ConnectionPaths: {
    typeUrl: string;
    encode(message: ConnectionPaths, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ConnectionPaths;
    fromJSON(object: any): ConnectionPaths;
    toJSON(message: ConnectionPaths): unknown;
    fromPartial<I extends {
        clientId?: string | undefined;
        paths?: string[] | undefined;
    } & {
        clientId?: string | undefined;
        paths?: (string[] & string[] & Record<Exclude<keyof I["paths"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof ConnectionPaths>, never>>(object: I): ConnectionPaths;
};
export declare const Version: {
    typeUrl: string;
    encode(message: Version, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Version;
    fromJSON(object: any): Version;
    toJSON(message: Version): unknown;
    fromPartial<I extends {
        identifier?: string | undefined;
        features?: string[] | undefined;
    } & {
        identifier?: string | undefined;
        features?: (string[] & string[] & Record<Exclude<keyof I["features"], keyof string[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof Version>, never>>(object: I): Version;
};
export declare const Params: {
    typeUrl: string;
    encode(message: Params, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial<I extends {
        maxExpectedTimePerBlock?: bigint | undefined;
    } & {
        maxExpectedTimePerBlock?: bigint | undefined;
    } & Record<Exclude<keyof I, "maxExpectedTimePerBlock">, never>>(object: I): Params;
};
