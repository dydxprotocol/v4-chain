import { LimitParams, LimitParamsSDKType } from "./limit_params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the ratelimit module's genesis state. */
export interface GenesisState {
    /** limit_params_list defines the list of `LimitParams` at genesis. */
    limitParamsList: LimitParams[];
}
/** GenesisState defines the ratelimit module's genesis state. */
export interface GenesisStateSDKType {
    limit_params_list: LimitParamsSDKType[];
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
