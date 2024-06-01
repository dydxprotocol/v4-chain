import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** EpochInfo stores metadata of an epoch timer. */
export interface EpochInfo {
    /** name is the unique identifier. */
    name: string;
    /**
     * next_tick indicates when the next epoch starts (in Unix Epoch seconds),
     * if `EpochInfo` has been initialized.
     * If `EpochInfo` is not initialized yet, `next_tick` indicates the earliest
     * initialization time (see `is_initialized` below).
     */
    nextTick: number;
    /** duration of the epoch in seconds. */
    duration: number;
    /**
     * current epoch is the number of the current epoch.
     * 0 if `next_tick` has never been reached, positive otherwise.
     */
    currentEpoch: number;
    /**
     * current_epoch_start_block indicates the block height when the current
     * epoch started. 0 if `current_epoch` is 0.
     */
    currentEpochStartBlock: number;
    /**
     * is_initialized indicates whether the `EpochInfo` has been initialized
     * and started ticking.
     * An `EpochInfo` is initialized when all below conditions are true:
     * - Not yet initialized
     * - `BlockHeight` >= 2
     * - `BlockTime` >= `next_tick`
     */
    isInitialized: boolean;
    /**
     * fast_forward_next_tick specifies whether during initialization, `next_tick`
     * should be fast-forwarded to be greater than the current block time.
     * If `false`, the original `next_tick` value is
     * unchanged during initialization.
     * If `true`, `next_tick` will be set to the smallest value `x` greater than
     * the current block time such that `(x - next_tick) % duration = 0`.
     */
    fastForwardNextTick: boolean;
}
/** EpochInfo stores metadata of an epoch timer. */
export interface EpochInfoSDKType {
    name: string;
    next_tick: number;
    duration: number;
    current_epoch: number;
    current_epoch_start_block: number;
    is_initialized: boolean;
    fast_forward_next_tick: boolean;
}
export declare const EpochInfo: {
    encode(message: EpochInfo, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): EpochInfo;
    fromPartial(object: DeepPartial<EpochInfo>): EpochInfo;
};
