import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.mint.module.v1";
/** Module is the config object of the mint module. */
export interface Module {
    feeCollectorName: string;
    /** authority defines the custom module authority. If not set, defaults to the governance module. */
    authority: string;
}
export declare const Module: {
    typeUrl: string;
    encode(message: Module, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Module;
    fromJSON(object: any): Module;
    toJSON(message: Module): unknown;
    fromPartial<I extends {
        feeCollectorName?: string | undefined;
        authority?: string | undefined;
    } & {
        feeCollectorName?: string | undefined;
        authority?: string | undefined;
    } & Record<Exclude<keyof I, keyof Module>, never>>(object: I): Module;
};
