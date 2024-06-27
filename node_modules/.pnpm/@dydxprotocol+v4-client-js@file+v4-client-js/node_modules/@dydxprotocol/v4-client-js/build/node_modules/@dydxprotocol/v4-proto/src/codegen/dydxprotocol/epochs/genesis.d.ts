import { EpochInfo, EpochInfoSDKType } from "./epoch_info";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the epochs module's genesis state. */
export interface GenesisState {
    epochInfoList: EpochInfo[];
}
/** GenesisState defines the epochs module's genesis state. */
export interface GenesisStateSDKType {
    epoch_info_list: EpochInfoSDKType[];
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
