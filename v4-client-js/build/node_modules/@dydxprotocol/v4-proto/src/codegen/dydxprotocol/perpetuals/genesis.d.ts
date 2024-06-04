import { Perpetual, PerpetualSDKType, LiquidityTier, LiquidityTierSDKType } from "./perpetual";
import { Params, ParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the perpetuals module's genesis state. */
export interface GenesisState {
    perpetuals: Perpetual[];
    liquidityTiers: LiquidityTier[];
    params?: Params;
}
/** GenesisState defines the perpetuals module's genesis state. */
export interface GenesisStateSDKType {
    perpetuals: PerpetualSDKType[];
    liquidity_tiers: LiquidityTierSDKType[];
    params?: ParamsSDKType;
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
