import { Duration, DurationSDKType } from "../../google/protobuf/duration";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Params defines the parameters for x/stats module. */
export interface Params {
    /** The desired number of seconds in the look-back window. */
    windowDuration?: Duration;
}
/** Params defines the parameters for x/stats module. */
export interface ParamsSDKType {
    window_duration?: DurationSDKType;
}
export declare const Params: {
    encode(message: Params, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Params;
    fromPartial(object: DeepPartial<Params>): Params;
};
