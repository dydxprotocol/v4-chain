import { CapabilityOwners } from "./capability";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.capability.v1beta1";
/** GenesisOwners defines the capability owners with their corresponding index. */
export interface GenesisOwners {
    /** index is the index of the capability owner. */
    index: bigint;
    /** index_owners are the owners at the given index. */
    indexOwners: CapabilityOwners;
}
/** GenesisState defines the capability module's genesis state. */
export interface GenesisState {
    /** index is the capability global index. */
    index: bigint;
    /**
     * owners represents a map from index to owners of the capability index
     * index key is string to allow amino marshalling.
     */
    owners: GenesisOwners[];
}
export declare const GenesisOwners: {
    typeUrl: string;
    encode(message: GenesisOwners, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisOwners;
    fromJSON(object: any): GenesisOwners;
    toJSON(message: GenesisOwners): unknown;
    fromPartial<I extends {
        index?: bigint | undefined;
        indexOwners?: {
            owners?: {
                module?: string | undefined;
                name?: string | undefined;
            }[] | undefined;
        } | undefined;
    } & {
        index?: bigint | undefined;
        indexOwners?: ({
            owners?: {
                module?: string | undefined;
                name?: string | undefined;
            }[] | undefined;
        } & {
            owners?: ({
                module?: string | undefined;
                name?: string | undefined;
            }[] & ({
                module?: string | undefined;
                name?: string | undefined;
            } & {
                module?: string | undefined;
                name?: string | undefined;
            } & Record<Exclude<keyof I["indexOwners"]["owners"][number], keyof import("./capability").Owner>, never>)[] & Record<Exclude<keyof I["indexOwners"]["owners"], keyof {
                module?: string | undefined;
                name?: string | undefined;
            }[]>, never>) | undefined;
        } & Record<Exclude<keyof I["indexOwners"], "owners">, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisOwners>, never>>(object: I): GenesisOwners;
};
export declare const GenesisState: {
    typeUrl: string;
    encode(message: GenesisState, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial<I extends {
        index?: bigint | undefined;
        owners?: {
            index?: bigint | undefined;
            indexOwners?: {
                owners?: {
                    module?: string | undefined;
                    name?: string | undefined;
                }[] | undefined;
            } | undefined;
        }[] | undefined;
    } & {
        index?: bigint | undefined;
        owners?: ({
            index?: bigint | undefined;
            indexOwners?: {
                owners?: {
                    module?: string | undefined;
                    name?: string | undefined;
                }[] | undefined;
            } | undefined;
        }[] & ({
            index?: bigint | undefined;
            indexOwners?: {
                owners?: {
                    module?: string | undefined;
                    name?: string | undefined;
                }[] | undefined;
            } | undefined;
        } & {
            index?: bigint | undefined;
            indexOwners?: ({
                owners?: {
                    module?: string | undefined;
                    name?: string | undefined;
                }[] | undefined;
            } & {
                owners?: ({
                    module?: string | undefined;
                    name?: string | undefined;
                }[] & ({
                    module?: string | undefined;
                    name?: string | undefined;
                } & {
                    module?: string | undefined;
                    name?: string | undefined;
                } & Record<Exclude<keyof I["owners"][number]["indexOwners"]["owners"][number], keyof import("./capability").Owner>, never>)[] & Record<Exclude<keyof I["owners"][number]["indexOwners"]["owners"], keyof {
                    module?: string | undefined;
                    name?: string | undefined;
                }[]>, never>) | undefined;
            } & Record<Exclude<keyof I["owners"][number]["indexOwners"], "owners">, never>) | undefined;
        } & Record<Exclude<keyof I["owners"][number], keyof GenesisOwners>, never>)[] & Record<Exclude<keyof I["owners"], keyof {
            index?: bigint | undefined;
            indexOwners?: {
                owners?: {
                    module?: string | undefined;
                    name?: string | undefined;
                }[] | undefined;
            } | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, keyof GenesisState>, never>>(object: I): GenesisState;
};
