/// <reference types="long" />
import { Duration, DurationSDKType } from "../../../google/protobuf/duration";
import { Long, DeepPartial } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
/**
 * ValidatorSigningInfo defines a validator's signing info for monitoring their
 * liveness activity.
 */
export interface ValidatorSigningInfo {
    address: string;
    /** Height at which validator was first a candidate OR was un-jailed */
    startHeight: Long;
    /**
     * Index which is incremented every time a validator is bonded in a block and
     * _may_ have signed a pre-commit or not. This in conjunction with the
     * signed_blocks_window param determines the index in the missed block bitmap.
     */
    indexOffset: Long;
    /** Timestamp until which the validator is jailed due to liveness downtime. */
    jailedUntil?: Date;
    /**
     * Whether or not a validator has been tombstoned (killed out of validator
     * set). It is set once the validator commits an equivocation or for any other
     * configured misbehavior.
     */
    tombstoned: boolean;
    /**
     * A counter of missed (unsigned) blocks. It is used to avoid unnecessary
     * reads in the missed block bitmap.
     */
    missedBlocksCounter: Long;
}
/**
 * ValidatorSigningInfo defines a validator's signing info for monitoring their
 * liveness activity.
 */
export interface ValidatorSigningInfoSDKType {
    address: string;
    start_height: Long;
    index_offset: Long;
    jailed_until?: Date;
    tombstoned: boolean;
    missed_blocks_counter: Long;
}
/** Params represents the parameters used for by the slashing module. */
export interface Params {
    signedBlocksWindow: Long;
    minSignedPerWindow: Uint8Array;
    downtimeJailDuration?: Duration;
    slashFractionDoubleSign: Uint8Array;
    slashFractionDowntime: Uint8Array;
}
/** Params represents the parameters used for by the slashing module. */
export interface ParamsSDKType {
    signed_blocks_window: Long;
    min_signed_per_window: Uint8Array;
    downtime_jail_duration?: DurationSDKType;
    slash_fraction_double_sign: Uint8Array;
    slash_fraction_downtime: Uint8Array;
}
export declare const ValidatorSigningInfo: {
    encode(message: ValidatorSigningInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ValidatorSigningInfo;
    fromPartial(object: DeepPartial<ValidatorSigningInfo>): ValidatorSigningInfo;
};
export declare const Params: {
    encode(message: Params, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Params;
    fromPartial(object: DeepPartial<Params>): Params;
};
