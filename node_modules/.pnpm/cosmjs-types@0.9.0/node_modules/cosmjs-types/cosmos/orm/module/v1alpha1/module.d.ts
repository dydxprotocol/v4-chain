import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.orm.module.v1alpha1";
/**
 * Module defines the ORM module which adds providers to the app container for
 * module-scoped DB's. In the future it may provide gRPC services for interacting
 * with ORM data.
 */
export interface Module {
}
export declare const Module: {
    typeUrl: string;
    encode(_: Module, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Module;
    fromJSON(_: any): Module;
    toJSON(_: Module): unknown;
    fromPartial<I extends {} & {} & Record<Exclude<keyof I, never>, never>>(_: I): Module;
};
