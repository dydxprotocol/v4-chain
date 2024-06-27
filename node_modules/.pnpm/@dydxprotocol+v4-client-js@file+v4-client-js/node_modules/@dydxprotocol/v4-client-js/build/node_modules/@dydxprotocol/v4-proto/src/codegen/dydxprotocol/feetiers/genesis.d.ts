import { PerpetualFeeParams, PerpetualFeeParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the feetiers module's genesis state. */
export interface GenesisState {
    /** The parameters for perpetual fees. */
    params?: PerpetualFeeParams;
}
/** GenesisState defines the feetiers module's genesis state. */
export interface GenesisStateSDKType {
    params?: PerpetualFeeParamsSDKType;
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
