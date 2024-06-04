import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.capability.module.v1";
/** Module is the config object of the capability module. */
export interface Module {
    /**
     * seal_keeper defines if keeper.Seal() will run on BeginBlock() to prevent further modules from creating a scoped
     * keeper. For more details check x/capability/keeper.go.
     */
    sealKeeper: boolean;
}
export declare const Module: {
    typeUrl: string;
    encode(message: Module, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Module;
    fromJSON(object: any): Module;
    toJSON(message: Module): unknown;
    fromPartial<I extends {
        sealKeeper?: boolean | undefined;
    } & {
        sealKeeper?: boolean | undefined;
    } & Record<Exclude<keyof I, "sealKeeper">, never>>(object: I): Module;
};
