import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.staking.module.v1";
/** Module is the config object of the staking module. */
export interface Module {
    /**
     * hooks_order specifies the order of staking hooks and should be a list
     * of module names which provide a staking hooks instance. If no order is
     * provided, then hooks will be applied in alphabetical order of module names.
     */
    hooksOrder: string[];
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
        hooksOrder?: string[] | undefined;
        authority?: string | undefined;
    } & {
        hooksOrder?: (string[] & string[] & Record<Exclude<keyof I["hooksOrder"], keyof string[]>, never>) | undefined;
        authority?: string | undefined;
    } & Record<Exclude<keyof I, keyof Module>, never>>(object: I): Module;
};
