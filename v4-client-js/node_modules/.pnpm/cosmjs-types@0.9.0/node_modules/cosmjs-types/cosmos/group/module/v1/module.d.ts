import { Duration } from "../../../../google/protobuf/duration";
import { BinaryReader, BinaryWriter } from "../../../../binary";
export declare const protobufPackage = "cosmos.group.module.v1";
/** Module is the config object of the group module. */
export interface Module {
    /**
     * max_execution_period defines the max duration after a proposal's voting period ends that members can send a MsgExec
     * to execute the proposal.
     */
    maxExecutionPeriod: Duration;
    /**
     * max_metadata_len defines the max length of the metadata bytes field for various entities within the group module.
     * Defaults to 255 if not explicitly set.
     */
    maxMetadataLen: bigint;
}
export declare const Module: {
    typeUrl: string;
    encode(message: Module, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Module;
    fromJSON(object: any): Module;
    toJSON(message: Module): unknown;
    fromPartial<I extends {
        maxExecutionPeriod?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        maxMetadataLen?: bigint | undefined;
    } & {
        maxExecutionPeriod?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["maxExecutionPeriod"], keyof Duration>, never>) | undefined;
        maxMetadataLen?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof Module>, never>>(object: I): Module;
};
