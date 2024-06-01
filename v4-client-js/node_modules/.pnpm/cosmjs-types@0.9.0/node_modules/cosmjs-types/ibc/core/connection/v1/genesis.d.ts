import { IdentifiedConnection, ConnectionPaths, Params } from "./connection";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "ibc.core.connection.v1";
/** GenesisState defines the ibc connection submodule's genesis state. */
export interface GenesisState {
    connections: IdentifiedConnection[];
    clientConnectionPaths: ConnectionPaths[];
    /** the sequence for the next generated connection identifier */
    nextConnectionSequence: bigint;
    params: Params;
}
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
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
        clientConnectionPaths?: {
            clientId?: string | undefined;
            paths?: string[] | undefined;
        }[] | undefined;
        nextConnectionSequence?: bigint | undefined;
        params?: {
            maxExpectedTimePerBlock?: bigint | undefined;
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
        clientConnectionPaths?: ({
            clientId?: string | undefined;
            paths?: string[] | undefined;
        }[] & ({
            clientId?: string | undefined;
            paths?: string[] | undefined;
        } & {
            clientId?: string | undefined;
            paths?: (string[] & string[] & Record<Exclude<keyof I["clientConnectionPaths"][number]["paths"], keyof string[]>, never>) | undefined;
        } & Record<Exclude<keyof I["clientConnectionPaths"][number], keyof ConnectionPaths>, never>)[] & Record<Exclude<keyof I["clientConnectionPaths"], keyof {
            clientId?: string | undefined;
            paths?: string[] | undefined;
        }[]>, never>) | undefined;
        nextConnectionSequence?: bigint | undefined;
        params?: ({
            maxExpectedTimePerBlock?: bigint | undefined;
        } & {
            maxExpectedTimePerBlock?: bigint | undefined;
        } & Record<Exclude<keyof I["params"], "maxExpectedTimePerBlock">, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
