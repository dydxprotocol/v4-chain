import { Asset, AssetSDKType } from "./asset";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the assets module's genesis state. */
export interface GenesisState {
    assets: Asset[];
}
/** GenesisState defines the assets module's genesis state. */
export interface GenesisStateSDKType {
    assets: AssetSDKType[];
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
