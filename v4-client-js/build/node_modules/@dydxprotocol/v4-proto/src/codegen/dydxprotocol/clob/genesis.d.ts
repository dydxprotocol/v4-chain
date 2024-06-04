import { ClobPair, ClobPairSDKType } from "./clob_pair";
import { LiquidationsConfig, LiquidationsConfigSDKType } from "./liquidations_config";
import { BlockRateLimitConfiguration, BlockRateLimitConfigurationSDKType } from "./block_rate_limit_config";
import { EquityTierLimitConfiguration, EquityTierLimitConfigurationSDKType } from "./equity_tier_limit_config";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the clob module's genesis state. */
export interface GenesisState {
    clobPairs: ClobPair[];
    liquidationsConfig?: LiquidationsConfig;
    blockRateLimitConfig?: BlockRateLimitConfiguration;
    equityTierLimitConfig?: EquityTierLimitConfiguration;
}
/** GenesisState defines the clob module's genesis state. */
export interface GenesisStateSDKType {
    clob_pairs: ClobPairSDKType[];
    liquidations_config?: LiquidationsConfigSDKType;
    block_rate_limit_config?: BlockRateLimitConfigurationSDKType;
    equity_tier_limit_config?: EquityTierLimitConfigurationSDKType;
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
