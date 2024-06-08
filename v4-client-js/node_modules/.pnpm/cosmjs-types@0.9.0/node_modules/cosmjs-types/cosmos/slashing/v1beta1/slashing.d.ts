import { Timestamp } from "../../../google/protobuf/timestamp";
import { Duration } from "../../../google/protobuf/duration";
import { BinaryReader, BinaryWriter } from "../../../binary";
export declare const protobufPackage = "cosmos.slashing.v1beta1";
/**
 * ValidatorSigningInfo defines a validator's signing info for monitoring their
 * liveness activity.
 */
export interface ValidatorSigningInfo {
    address: string;
    /** Height at which validator was first a candidate OR was unjailed */
    startHeight: bigint;
    /**
     * Index which is incremented each time the validator was a bonded
     * in a block and may have signed a precommit or not. This in conjunction with the
     * `SignedBlocksWindow` param determines the index in the `MissedBlocksBitArray`.
     */
    indexOffset: bigint;
    /** Timestamp until which the validator is jailed due to liveness downtime. */
    jailedUntil: Timestamp;
    /**
     * Whether or not a validator has been tombstoned (killed out of validator set). It is set
     * once the validator commits an equivocation or for any other configured misbehiavor.
     */
    tombstoned: boolean;
    /**
     * A counter kept to avoid unnecessary array reads.
     * Note that `Sum(MissedBlocksBitArray)` always equals `MissedBlocksCounter`.
     */
    missedBlocksCounter: bigint;
}
/** Params represents the parameters used for by the slashing module. */
export interface Params {
    signedBlocksWindow: bigint;
    minSignedPerWindow: Uint8Array;
    downtimeJailDuration: Duration;
    slashFractionDoubleSign: Uint8Array;
    slashFractionDowntime: Uint8Array;
}
export declare const ValidatorSigningInfo: {
    typeUrl: string;
    encode(message: ValidatorSigningInfo, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): ValidatorSigningInfo;
    fromJSON(object: any): ValidatorSigningInfo;
    toJSON(message: ValidatorSigningInfo): unknown;
    fromPartial<I extends {
        address?: string | undefined;
        startHeight?: bigint | undefined;
        indexOffset?: bigint | undefined;
        jailedUntil?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        tombstoned?: boolean | undefined;
        missedBlocksCounter?: bigint | undefined;
    } & {
        address?: string | undefined;
        startHeight?: bigint | undefined;
        indexOffset?: bigint | undefined;
        jailedUntil?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["jailedUntil"], keyof Timestamp>, never>) | undefined;
        tombstoned?: boolean | undefined;
        missedBlocksCounter?: bigint | undefined;
    } & Record<Exclude<keyof I, keyof ValidatorSigningInfo>, never>>(object: I): ValidatorSigningInfo;
};
export declare const Params: {
    typeUrl: string;
    encode(message: Params, writer?: BinaryWriter): BinaryWriter;
    decode(input: BinaryReader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial<I extends {
        signedBlocksWindow?: bigint | undefined;
        minSignedPerWindow?: Uint8Array | undefined;
        downtimeJailDuration?: {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } | undefined;
        slashFractionDoubleSign?: Uint8Array | undefined;
        slashFractionDowntime?: Uint8Array | undefined;
    } & {
        signedBlocksWindow?: bigint | undefined;
        minSignedPerWindow?: Uint8Array | undefined;
        downtimeJailDuration?: ({
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & {
            seconds?: bigint | undefined;
            nanos?: number | undefined;
        } & Record<Exclude<keyof I["downtimeJailDuration"], keyof Duration>, never>) | undefined;
        slashFractionDoubleSign?: Uint8Array | undefined;
        slashFractionDowntime?: Uint8Array | undefined;
    } & Record<Exclude<keyof I, keyof Params>, never>>(object: I): Params;
};
