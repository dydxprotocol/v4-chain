import { Duration, DurationSDKType } from "../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** DowntimeParams defines the parameters for downtime. */
export interface DowntimeParams {
    /**
     * Durations tracked for downtime. The durations must be sorted from
     * shortest to longest and must all be positive.
     */
    durations: Duration[];
}
/** DowntimeParams defines the parameters for downtime. */
export interface DowntimeParamsSDKType {
    durations: DurationSDKType[];
}
export declare const DowntimeParams: {
    encode(message: DowntimeParams, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DowntimeParams;
    fromPartial(object: DeepPartial<DowntimeParams>): DowntimeParams;
};
