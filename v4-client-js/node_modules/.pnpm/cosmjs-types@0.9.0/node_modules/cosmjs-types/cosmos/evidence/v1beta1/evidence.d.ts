import { Timestamp } from "../../../google/protobuf/timestamp";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.evidence.v1beta1";
/**
 * Equivocation implements the Evidence interface and defines evidence of double
 * signing misbehavior.
 */
export interface Equivocation {
    /** height is the equivocation height. */
    height: bigint;
    /** time is the equivocation time. */
    time: Timestamp;
    /** power is the equivocation validator power. */
    power: bigint;
    /** consensus_address is the equivocation validator consensus address. */
    consensusAddress: string;
}
export declare const Equivocation: {
    typeUrl: string;
    encode(message: Equivocation, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Equivocation;
    fromJSON(object: any): Equivocation;
    toJSON(message: Equivocation): unknown;
    fromPartial<I extends {
        height?: bigint | undefined;
        time?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        power?: bigint | undefined;
        consensusAddress?: string | undefined;
    } & {
        height?: bigint | undefined;
        time?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["time"], keyof Timestamp>, never>) | undefined;
        power?: bigint | undefined;
        consensusAddress?: string | undefined;
    } & Record<Exclude<keyof I, keyof Equivocation>, never>>(object: I): Equivocation;
};
