import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.capability.v1beta1";
/**
 * Capability defines an implementation of an object capability. The index
 * provided to a Capability must be globally unique.
 */
export interface Capability {
    index: bigint;
}
/**
 * Owner defines a single capability owner. An owner is defined by the name of
 * capability and the module name.
 */
export interface Owner {
    module: string;
    name: string;
}
/**
 * CapabilityOwners defines a set of owners of a single Capability. The set of
 * owners must be unique.
 */
export interface CapabilityOwners {
    owners: Owner[];
}
export declare const Capability: {
    typeUrl: string;
    encode(message: Capability, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Capability;
    fromJSON(object: any): Capability;
    toJSON(message: Capability): unknown;
    fromPartial<I extends {
        index?: bigint | undefined;
    } & {
        index?: bigint | undefined;
    } & Record<Exclude<keyof I, "index">, never>>(object: I): Capability;
};
export declare const Owner: {
    typeUrl: string;
    encode(message: Owner, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Owner;
    fromJSON(object: any): Owner;
    toJSON(message: Owner): unknown;
    fromPartial<I extends {
        module?: string | undefined;
        name?: string | undefined;
    } & {
        module?: string | undefined;
        name?: string | undefined;
    } & Record<Exclude<keyof I, keyof Owner>, never>>(object: I): Owner;
};
export declare const CapabilityOwners: {
    typeUrl: string;
    encode(message: CapabilityOwners, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): CapabilityOwners;
    fromJSON(object: any): CapabilityOwners;
    toJSON(message: CapabilityOwners): unknown;
    fromPartial<I extends {
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
        } & Record<Exclude<keyof I["owners"][number], keyof Owner>, never>)[] & Record<Exclude<keyof I["owners"], keyof {
            module?: string | undefined;
            name?: string | undefined;
        }[]>, never>) | undefined;
    } & Record<Exclude<keyof I, "owners">, never>>(object: I): CapabilityOwners;
};
