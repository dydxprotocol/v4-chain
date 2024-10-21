import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Params defines the parameters for x/rewards module. */
export interface Params {
    /** The module account to distribute rewards from. */
    treasuryAccount: string;
    /** The denom of the rewards token. */
    denom: string;
    /**
     * The exponent of converting one unit of `denom` to a full coin.
     * For example, `denom=uatom, denom_exponent=-6` defines that
     * `1 uatom = 10^(-6) ATOM`. This conversion is needed since the
     * `market_id` retrieves the price of a full coin of the reward token.
     */
    denomExponent: number;
    /** The id of the market that has the price of the rewards token. */
    marketId: number;
    /**
     * The amount (in ppm) that fees are multiplied by to get
     * the maximum rewards amount.
     */
    feeMultiplierPpm: number;
}
/** Params defines the parameters for x/rewards module. */
export interface ParamsSDKType {
    treasury_account: string;
    denom: string;
    denom_exponent: number;
    market_id: number;
    fee_multiplier_ppm: number;
}
export declare const Params: {
    encode(message: Params, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Params;
    fromPartial(object: DeepPartial<Params>): Params;
};
