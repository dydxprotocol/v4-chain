import { IdentifiedClientState, ClientConsensusStates, Params } from "./client";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.core.client.v1";
/** GenesisState defines the ibc client submodule's genesis state. */
export interface GenesisState {
    /** client states with their corresponding identifiers */
    clients: IdentifiedClientState[];
    /** consensus states from each client */
    clientsConsensus: ClientConsensusStates[];
    /** metadata from each client */
    clientsMetadata: IdentifiedGenesisMetadata[];
    params: Params;
    /** create localhost on initialization */
    createLocalhost: boolean;
    /** the sequence for the next generated client identifier */
    nextClientSequence: bigint;
}
/**
 * GenesisMetadata defines the genesis type for metadata that clients may return
 * with ExportMetadata
 */
export interface GenesisMetadata {
    /** store key of metadata without clientID-prefix */
    key: Uint8Array;
    /** metadata value */
    value: Uint8Array;
}
/**
 * IdentifiedGenesisMetadata has the client metadata with the corresponding
 * client id.
 */
export interface IdentifiedGenesisMetadata {
    clientId: string;
    clientMetadata: GenesisMetadata[];
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        clients?: {
            clientId?: string | undefined;
            clientState?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] | undefined;
        clientsConsensus?: {
            clientId?: string | undefined;
            consensusStates?: {
                height?: {
                    revisionNumber?: bigint | undefined;
                    revisionHeight?: bigint | undefined;
                } | undefined;
                consensusState?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
            }[] | undefined;
        }[] | undefined;
        clientsMetadata?: {
            clientId?: string | undefined;
            clientMetadata?: {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
        }[] | undefined;
        params?: {
            allowedClients?: string[] | undefined;
        } | undefined;
        createLocalhost?: boolean | undefined;
        nextClientSequence?: bigint | undefined;
    } & {
        clients?: ({
            clientId?: string | undefined;
            clientState?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[] & ({
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
            } & Record<Exclude<keyof I["clients"][number]["clientState"], keyof import("../../../../google/protobuf/any").Any>, never>) | undefined;
        } & Record<Exclude<keyof I["clients"][number], keyof IdentifiedClientState>, never>)[] & Record<Exclude<keyof I["clients"], keyof {
            clientId?: string | undefined;
            clientState?: {
                typeUrl?: string | undefined;
                value?: Uint8Array | undefined;
            } | undefined;
        }[]>, never>) | undefined;
        clientsConsensus?: ({
            clientId?: string | undefined;
            consensusStates?: {
                height?: {
                    revisionNumber?: bigint | undefined;
                    revisionHeight?: bigint | undefined;
                } | undefined;
                consensusState?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
            }[] | undefined;
        }[] & ({
            clientId?: string | undefined;
            consensusStates?: {
                height?: {
                    revisionNumber?: bigint | undefined;
                    revisionHeight?: bigint | undefined;
                } | undefined;
                consensusState?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
            }[] | undefined;
        } & {
            clientId?: string | undefined;
            consensusStates?: ({
                height?: {
                    revisionNumber?: bigint | undefined;
                    revisionHeight?: bigint | undefined;
                } | undefined;
                consensusState?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
            }[] & ({
                height?: {
                    revisionNumber?: bigint | undefined;
                    revisionHeight?: bigint | undefined;
                } | undefined;
                consensusState?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
            } & {
                height?: ({
                    revisionNumber?: bigint | undefined;
                    revisionHeight?: bigint | undefined;
                } & {
                    revisionNumber?: bigint | undefined;
                    revisionHeight?: bigint | undefined;
                } & Record<Exclude<keyof I["clientsConsensus"][number]["consensusStates"][number]["height"], keyof import("./client").Height>, never>) | undefined;
                consensusState?: ({
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } & Record<Exclude<keyof I["clientsConsensus"][number]["consensusStates"][number]["consensusState"], keyof import("../../../../google/protobuf/any").Any>, never>) | undefined;
            } & Record<Exclude<keyof I["clientsConsensus"][number]["consensusStates"][number], keyof import("./client").ConsensusStateWithHeight>, never>)[] & Record<Exclude<keyof I["clientsConsensus"][number]["consensusStates"], keyof {
                height?: {
                    revisionNumber?: bigint | undefined;
                    revisionHeight?: bigint | undefined;
                } | undefined;
                consensusState?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["clientsConsensus"][number], keyof ClientConsensusStates>, never>)[] & Record<Exclude<keyof I["clientsConsensus"], keyof {
            clientId?: string | undefined;
            consensusStates?: {
                height?: {
                    revisionNumber?: bigint | undefined;
                    revisionHeight?: bigint | undefined;
                } | undefined;
                consensusState?: {
                    typeUrl?: string | undefined;
                    value?: Uint8Array | undefined;
                } | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        clientsMetadata?: ({
            clientId?: string | undefined;
            clientMetadata?: {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
        }[] & ({
            clientId?: string | undefined;
            clientMetadata?: {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
        } & {
            clientId?: string | undefined;
            clientMetadata?: ({
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[] & ({
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            } & {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            } & Record<Exclude<keyof I["clientsMetadata"][number]["clientMetadata"][number], keyof GenesisMetadata>, never>)[] & Record<Exclude<keyof I["clientsMetadata"][number]["clientMetadata"], keyof {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["clientsMetadata"][number], keyof IdentifiedGenesisMetadata>, never>)[] & Record<Exclude<keyof I["clientsMetadata"], keyof {
            clientId?: string | undefined;
            clientMetadata?: {
                key?: Uint8Array | undefined;
                value?: Uint8Array | undefined;
            }[] | undefined;
        }[]>, never>) | undefined;
        params?: ({
            allowedClients?: string[] | undefined;
        } & {
            allowedClients?: (string[] & string[] & Record<Exclude<keyof I["params"]["allowedClients"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["params"], "allowedClients">, never>) | undefined;
        createLocalhost?: boolean | undefined;
        nextClientSequence?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
export declare const GenesisMetadata: {
    typeUrl: string;
    encode(message: GenesisMetadata, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisMetadata;
    fromJSON(object: any): GenesisMetadata;
    toJSON(message: GenesisMetadata): unknown;
    fromPartial<I extends {
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
    } & {
        key?: Uint8Array | undefined;
        value?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof GenesisMetadata>, never>>(object: I): GenesisMetadata;
};
export declare const IdentifiedGenesisMetadata: {
    typeUrl: string;
    encode(message: IdentifiedGenesisMetadata, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): IdentifiedGenesisMetadata;
    fromJSON(object: any): IdentifiedGenesisMetadata;
    toJSON(message: IdentifiedGenesisMetadata): unknown;
    fromPartial<I extends {
        clientId?: string | undefined;
        clientMetadata?: {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        }[] | undefined;
    } & {
        clientId?: string | undefined;
        clientMetadata?: ({
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        }[] & ({
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        } & {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        } & Record<Exclude<keyof I["clientMetadata"][number], keyof GenesisMetadata>, never>)[] & Record<Exclude<keyof I["clientMetadata"], keyof {
            key?: Uint8Array | undefined;
            value?: Uint8Array | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof IdentifiedGenesisMetadata>, never>>(object: I): IdentifiedGenesisMetadata;
};
