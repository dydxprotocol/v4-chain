import { PageRequest, PageResponse } from "../../../../cosmos/base/query/v1beta1/pagination";
import { ConnectionEnd, IdentifiedConnection } from "./connection";
import { Height, IdentifiedClientState, Params } from "../../client/v1/client";
import { Any } from "../../../../google/protobuf/any";
import { BinaryReader, BinaryWriter } from "../../../../binary";
import { Rpc } from "../../../../helpers";
export declare const protobufPackage = "ibc.core.connection.v1";
/**
 * QueryConnectionRequest is the request type for the Query/Connection RPC
 * method
 */
export interface QueryConnectionRequest {
    /** connection unique identifier */
    connectionId: string;
}
/**
 * QueryConnectionResponse is the response type for the Query/Connection RPC
 * method. Besides the connection end, it includes a proof and the height from
 * which the proof was retrieved.
 */
export interface QueryConnectionResponse {
    /** connection associated with the request identifier */
    connection?: ConnectionEnd;
    /** merkle proof of existence */
    proof: Uint8Array;
    /** height at which the proof was retrieved */
    proofHeight: Height;
}
/**
 * QueryConnectionsRequest is the request type for the Query/Connections RPC
 * method
 */
export interface QueryConnectionsRequest {
    pagination?: PageRequest;
}
/**
 * QueryConnectionsResponse is the response type for the Query/Connections RPC
 * method.
 */
export interface QueryConnectionsResponse {
    /** list of stored connections of the chain. */
    connections: IdentifiedConnection[];
    /** pagination response */
    pagination?: PageResponse;
    /** query block height */
    height: Height;
}
/**
 * QueryClientConnectionsRequest is the request type for the
 * Query/ClientConnections RPC method
 */
export interface QueryClientConnectionsRequest {
    /** client identifier associated with a connection */
    clientId: string;
}
/**
 * QueryClientConnectionsResponse is the response type for the
 * Query/ClientConnections RPC method
 */
export interface QueryClientConnectionsResponse {
    /** slice of all the connection paths associated with a client. */
    connectionPaths: string[];
    /** merkle proof of existence */
    proof: Uint8Array;
    /** height at which the proof was generated */
    proofHeight: Height;
}
/**
 * QueryConnectionClientStateRequest is the request type for the
 * Query/ConnectionClientState RPC method
 */
export interface QueryConnectionClientStateRequest {
    /** connection identifier */
    connectionId: string;
}
/**
 * QueryConnectionClientStateResponse is the response type for the
 * Query/ConnectionClientState RPC method
 */
export interface QueryConnectionClientStateResponse {
    /** client state associated with the channel */
    identifiedClientState?: IdentifiedClientState;
    /** merkle proof of existence */
    proof: Uint8Array;
    /** height at which the proof was retrieved */
    proofHeight: Height;
}
/**
 * QueryConnectionConsensusStateRequest is the request type for the
 * Query/ConnectionConsensusState RPC method
 */
export interface QueryConnectionConsensusStateRequest {
    /** connection identifier */
    connectionId: string;
    revisionNumber: bigint;
    revisionHeight: bigint;
}
/**
 * QueryConnectionConsensusStateResponse is the response type for the
 * Query/ConnectionConsensusState RPC method
 */
export interface QueryConnectionConsensusStateResponse {
    /** consensus state associated with the channel */
    consensusState?: Any;
    /** client ID associated with the consensus state */
    clientId: string;
    /** merkle proof of existence */
    proof: Uint8Array;
    /** height at which the proof was retrieved */
    proofHeight: Height;
}
/** QueryConnectionParamsRequest is the request type for the Query/ConnectionParams RPC method. */
export interface QueryConnectionParamsRequest {
}
/** QueryConnectionParamsResponse is the response type for the Query/ConnectionParams RPC method. */
export interface QueryConnectionParamsResponse {
    /** params defines the parameters of the module. */
    params?: Params;
}
export declare const QueryConnectionRequest: {
    typeUrl: string;
    encode(message: QueryConnectionRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionRequest;
    fromJSON(object: any): QueryConnectionRequest;
    toJSON(message: QueryConnectionRequest): unknown;
    fromPartial<I extends {
        connectionId?: string | undefined;
    } & {
        connectionId?: string | undefined;
    } & Record<Exclude<keyof I, "connectionId">, never>>(object: I): QueryConnectionRequest;
};
export declare const QueryConnectionResponse: {
    typeUrl: string;
    encode(message: QueryConnectionResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionResponse;
    fromJSON(object: any): QueryConnectionResponse;
    toJSON(message: QueryConnectionResponse): unknown;
    fromPartial<I extends {
        connection?: {
            clientId?: string | undefined;
            versions?: {
                identifier?: string | undefined;
                features?: string[] | undefined;
            }[] | undefined;
            state?: import("./connection").State | undefined;
            counterparty?: {
                clientId?: string | undefined;
                connectionId?: string | undefined;
                prefix?: {
                    keyPrefix?: Uint8Array | undefined;
                } | undefined;
            } | undefined;
            delayPeriod?: bigint | undefined;
        } | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
        connection?: ({
            clientId?: string | undefined;
            versions?: {
                identifier?: string | undefined;
                features?: string[] | undefined;
            }[] | undefined;
            state?: import("./connection").State | undefined;
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
                features?: (string[] & string[] & Record<Exclude<keyof I["connection"]["versions"][number]["features"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["connection"]["versions"][number], keyof import("./connection").Version>, never>)[] & Record<Exclude<keyof I["connection"]["versions"], keyof {
                identifier?: string | undefined;
                features?: string[] | undefined;
            }[]>, never>) | undefined;
            state?: import("./connection").State | undefined;
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
                } & Record<Exclude<keyof I["connection"]["counterparty"]["prefix"], "keyPrefix">, never>) | undefined;
            } & Record<Exclude<keyof I["connection"]["counterparty"], keyof import("./connection").Counterparty>, never>) | undefined;
            delayPeriod?: bigint | undefined;
        } & Record<Exclude<keyof I["connection"], keyof ConnectionEnd>, never>) | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryConnectionResponse>, never>>(object: I): QueryConnectionResponse;
};
export declare const QueryConnectionsRequest: {
    typeUrl: string;
    encode(message: QueryConnectionsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionsRequest;
    fromJSON(object: any): QueryConnectionsRequest;
    toJSON(message: QueryConnectionsRequest): unknown;
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
    } & Record<Exclude<keyof I, "pagination">, never>>(object: I): QueryConnectionsRequest;
};
export declare const QueryConnectionsResponse: {
    typeUrl: string;
    encode(message: QueryConnectionsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionsResponse;
    fromJSON(object: any): QueryConnectionsResponse;
    toJSON(message: QueryConnectionsResponse): unknown;
    fromPartial<I extends {
        connections?: {
            id?: string | undefined;
            clientId?: string | undefined;
            versions?: {
                identifier?: string | undefined;
                features?: string[] | undefined;
            }[] | undefined;
            state?: import("./connection").State | undefined;
            counterparty?: {
                clientId?: string | undefined;
                connectionId?: string | undefined;
                prefix?: {
                    keyPrefix?: Uint8Array | undefined;
                } | undefined;
            } | undefined;
            delayPeriod?: bigint | undefined;
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
        connections?: ({
            id?: string | undefined;
            clientId?: string | undefined;
            versions?: {
                identifier?: string | undefined;
                features?: string[] | undefined;
            }[] | undefined;
            state?: import("./connection").State | undefined;
            counterparty?: {
                clientId?: string | undefined;
                connectionId?: string | undefined;
                prefix?: {
                    keyPrefix?: Uint8Array | undefined;
                } | undefined;
            } | undefined;
            delayPeriod?: bigint | undefined;
        }[] & ({
            id?: string | undefined;
            clientId?: string | undefined;
            versions?: {
                identifier?: string | undefined;
                features?: string[] | undefined;
            }[] | undefined;
            state?: import("./connection").State | undefined;
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
                features?: (string[] & string[] & Record<Exclude<keyof I["connections"][number]["versions"][number]["features"], keyof string[]>, never>) | undefined;
            } & Record<Exclude<keyof I["connections"][number]["versions"][number], keyof import("./connection").Version>, never>)[] & Record<Exclude<keyof I["connections"][number]["versions"], keyof {
                identifier?: string | undefined;
                features?: string[] | undefined;
            }[]>, never>) | undefined;
            state?: import("./connection").State | undefined;
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
                } & Record<Exclude<keyof I["connections"][number]["counterparty"]["prefix"], "keyPrefix">, never>) | undefined;
            } & Record<Exclude<keyof I["connections"][number]["counterparty"], keyof import("./connection").Counterparty>, never>) | undefined;
            delayPeriod?: bigint | undefined;
        } & Record<Exclude<keyof I["connections"][number], keyof IdentifiedConnection>, never>)[] & Record<Exclude<keyof I["connections"], keyof {
            id?: string | undefined;
            clientId?: string | undefined;
            versions?: {
                identifier?: string | undefined;
                features?: string[] | undefined;
            }[] | undefined;
            state?: import("./connection").State | undefined;
            counterparty?: {
                clientId?: string | undefined;
                connectionId?: string | undefined;
                prefix?: {
                    keyPrefix?: Uint8Array | undefined;
                } | undefined;
            } | undefined;
            delayPeriod?: bigint | undefined;
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
    } & Record<Exclude<keyof I, keyof QueryConnectionsResponse>, never>>(object: I): QueryConnectionsResponse;
};
export declare const QueryClientConnectionsRequest: {
    typeUrl: string;
    encode(message: QueryClientConnectionsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryClientConnectionsRequest;
    fromJSON(object: any): QueryClientConnectionsRequest;
    toJSON(message: QueryClientConnectionsRequest): unknown;
    fromPartial<I extends {
        clientId?: string | undefined;
    } & {
        clientId?: string | undefined;
    } & Record<Exclude<keyof I, "clientId">, never>>(object: I): QueryClientConnectionsRequest;
};
export declare const QueryClientConnectionsResponse: {
    typeUrl: string;
    encode(message: QueryClientConnectionsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryClientConnectionsResponse;
    fromJSON(object: any): QueryClientConnectionsResponse;
    toJSON(message: QueryClientConnectionsResponse): unknown;
    fromPartial<I extends {
        connectionPaths?: string[] | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } | undefined;
    } & {
        connectionPaths?: (string[] & string[] & Record<Exclude<keyof I["connectionPaths"], keyof string[]>, never>) | undefined;
        proof?: Uint8Array | undefined;
        proofHeight?: ({
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & {
            revisionNumber?: bigint | undefined;
            revisionHeight?: bigint | undefined;
        } & Record<Exclude<keyof I["proofHeight"], keyof Height>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof QueryClientConnectionsResponse>, never>>(object: I): QueryClientConnectionsResponse;
};
export declare const QueryConnectionClientStateRequest: {
    typeUrl: string;
    encode(message: QueryConnectionClientStateRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionClientStateRequest;
    fromJSON(object: any): QueryConnectionClientStateRequest;
    toJSON(message: QueryConnectionClientStateRequest): unknown;
    fromPartial<I extends {
        connectionId?: string | undefined;
    } & {
        connectionId?: string | undefined;
    } & Record<Exclude<keyof I, "connectionId">, never>>(object: I): QueryConnectionClientStateRequest;
};
export declare const QueryConnectionClientStateResponse: {
    typeUrl: string;
    encode(message: QueryConnectionClientStateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionClientStateResponse;
    fromJSON(object: any): QueryConnectionClientStateResponse;
    toJSON(message: QueryConnectionClientStateResponse): unknown;
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
    } & Record<Exclude<keyof I, keyof QueryConnectionClientStateResponse>, never>>(object: I): QueryConnectionClientStateResponse;
};
export declare const QueryConnectionConsensusStateRequest: {
    typeUrl: string;
    encode(message: QueryConnectionConsensusStateRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionConsensusStateRequest;
    fromJSON(object: any): QueryConnectionConsensusStateRequest;
    toJSON(message: QueryConnectionConsensusStateRequest): unknown;
    fromPartial<I extends {
        connectionId?: string | undefined;
        revisionNumber?: bigint | undefined;
        revisionHeight?: bigint | undefined;
    } & {
        connectionId?: string | undefined;
        revisionNumber?: bigint | undefined;
        revisionHeight?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof QueryConnectionConsensusStateRequest>, never>>(object: I): QueryConnectionConsensusStateRequest;
};
export declare const QueryConnectionConsensusStateResponse: {
    typeUrl: string;
    encode(message: QueryConnectionConsensusStateResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionConsensusStateResponse;
    fromJSON(object: any): QueryConnectionConsensusStateResponse;
    toJSON(message: QueryConnectionConsensusStateResponse): unknown;
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
    } & Record<Exclude<keyof I, keyof QueryConnectionConsensusStateResponse>, never>>(object: I): QueryConnectionConsensusStateResponse;
};
export declare const QueryConnectionParamsRequest: {
    typeUrl: string;
    encode(_: QueryConnectionParamsRequest, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionParamsRequest;
    fromJSON(_: any): QueryConnectionParamsRequest;
    toJSON(_: QueryConnectionParamsRequest): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): QueryConnectionParamsRequest;
};
export declare const QueryConnectionParamsResponse: {
    typeUrl: string;
    encode(message: QueryConnectionParamsResponse, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): QueryConnectionParamsResponse;
    fromJSON(object: any): QueryConnectionParamsResponse;
    toJSON(message: QueryConnectionParamsResponse): unknown;
    fromPartial<I extends {
        params?: {
            allowedClients?: string[] | undefined;
        } | undefined;
    } & {
        params?: ({
            allowedClients?: string[] | undefined;
        } & {
            allowedClients?: (string[] & string[] & Record<Exclude<keyof I["params"]["allowedClients"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["params"], "allowedClients">, never>) | undefined;
    } & Record<Exclude<keyof I, "params">, never>>(object: I): QueryConnectionParamsResponse;
};
/** Query provides defines the gRPC querier service */
export interface Query {
    /** Connection queries an IBC connection end. */
    Connection(request: QueryConnectionRequest): Promise<QueryConnectionResponse>;
    /** Connections queries all the IBC connections of a chain. */
    Connections(request?: QueryConnectionsRequest): Promise<QueryConnectionsResponse>;
    /**
     * ClientConnections queries the connection paths associated with a client
     * state.
     */
    ClientConnections(request: QueryClientConnectionsRequest): Promise<QueryClientConnectionsResponse>;
    /**
     * ConnectionClientState queries the client state associated with the
     * connection.
     */
    ConnectionClientState(request: QueryConnectionClientStateRequest): Promise<QueryConnectionClientStateResponse>;
    /**
     * ConnectionConsensusState queries the consensus state associated with the
     * connection.
     */
    ConnectionConsensusState(request: QueryConnectionConsensusStateRequest): Promise<QueryConnectionConsensusStateResponse>;
    /** ConnectionParams queries all parameters of the ibc connection submodule. */
    ConnectionParams(request?: QueryConnectionParamsRequest): Promise<QueryConnectionParamsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    Connection(request: QueryConnectionRequest): Promise<QueryConnectionResponse>;
    Connections(request?: QueryConnectionsRequest): Promise<QueryConnectionsResponse>;
    ClientConnections(request: QueryClientConnectionsRequest): Promise<QueryClientConnectionsResponse>;
    ConnectionClientState(request: QueryConnectionClientStateRequest): Promise<QueryConnectionClientStateResponse>;
    ConnectionConsensusState(request: QueryConnectionConsensusStateRequest): Promise<QueryConnectionConsensusStateResponse>;
    ConnectionParams(request?: QueryConnectionParamsRequest): Promise<QueryConnectionParamsResponse>;
}
