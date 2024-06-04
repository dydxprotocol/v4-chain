import { MarketParam, MarketParamSDKType } from "./market_param";
import { MarketPrice, MarketPriceSDKType } from "./market_price";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** GenesisState defines the prices module's genesis state. */
export interface GenesisState {
    marketParams: MarketParam[];
    marketPrices: MarketPrice[];
}
/** GenesisState defines the prices module's genesis state. */
export interface GenesisStateSDKType {
    market_params: MarketParamSDKType[];
    market_prices: MarketPriceSDKType[];
}
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
