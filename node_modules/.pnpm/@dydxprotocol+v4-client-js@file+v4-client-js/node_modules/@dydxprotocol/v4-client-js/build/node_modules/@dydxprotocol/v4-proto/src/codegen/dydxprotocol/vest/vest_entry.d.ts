import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * VestEntry specifies a Vester Account and the rate at which tokens are
 * dripped into the corresponding Treasury Account.
 */
export interface VestEntry {
    /**
     * The module account to vest tokens from.
     * This is also the key to this `VestEntry` in state.
     */
    vesterAccount: string;
    /** The module account to vest tokens to. */
    treasuryAccount: string;
    /** The denom of the token to vest. */
    denom: string;
    /** The start time of vest. Before this time, no vest will occur. */
    startTime?: Date;
    /**
     * The end time of vest. At this target date, all funds should be in the
     * Treasury Account and none left in the Vester Account.
     */
    endTime?: Date;
}
/**
 * VestEntry specifies a Vester Account and the rate at which tokens are
 * dripped into the corresponding Treasury Account.
 */
export interface VestEntrySDKType {
    vester_account: string;
    treasury_account: string;
    denom: string;
    start_time?: Date;
    end_time?: Date;
}
export declare const VestEntry: {
    encode(message: VestEntry, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): VestEntry;
    fromPartial(object: DeepPartial<VestEntry>): VestEntry;
};
