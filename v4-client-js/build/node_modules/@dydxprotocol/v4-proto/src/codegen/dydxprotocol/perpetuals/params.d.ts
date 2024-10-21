import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Params defines the parameters for x/perpetuals module. */
export interface Params {
    /**
     * Funding rate clamp factor in parts-per-million, used for clamping 8-hour
     * funding rates according to equation: |R| <= funding_rate_clamp_factor *
     * (initial margin - maintenance margin).
     */
    fundingRateClampFactorPpm: number;
    /**
     * Premium vote clamp factor in parts-per-million, used for clamping premium
     * votes according to equation: |V| <= premium_vote_clamp_factor *
     * (initial margin - maintenance margin).
     */
    premiumVoteClampFactorPpm: number;
    /**
     * Minimum number of premium votes per premium sample. If number of premium
     * votes is smaller than this number, pad with zeros up to this number.
     */
    minNumVotesPerSample: number;
}
/** Params defines the parameters for x/perpetuals module. */
export interface ParamsSDKType {
    funding_rate_clamp_factor_ppm: number;
    premium_vote_clamp_factor_ppm: number;
    min_num_votes_per_sample: number;
}
export declare const Params: {
    encode(message: Params, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Params;
    fromPartial(object: DeepPartial<Params>): Params;
};
