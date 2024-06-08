import { PerpetualParams, PerpetualParamsSDKType, LiquidityTier, LiquidityTierSDKType } from "./perpetual";
import { Params, ParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgCreatePerpetual is a message used by x/gov to create a new perpetual. */
export interface MsgCreatePerpetual {
    /** The address that controls the module. */
    authority: string;
    /** `params` defines parameters for the new perpetual market. */
    params?: PerpetualParams;
}
/** MsgCreatePerpetual is a message used by x/gov to create a new perpetual. */
export interface MsgCreatePerpetualSDKType {
    authority: string;
    params?: PerpetualParamsSDKType;
}
/**
 * MsgCreatePerpetualResponse defines the CreatePerpetual
 * response type.
 */
export interface MsgCreatePerpetualResponse {
}
/**
 * MsgCreatePerpetualResponse defines the CreatePerpetual
 * response type.
 */
export interface MsgCreatePerpetualResponseSDKType {
}
/**
 * MsgSetLiquidityTier is a message used by x/gov to create or update a
 * liquidity tier.
 */
export interface MsgSetLiquidityTier {
    /** The address that controls the module. */
    authority: string;
    /** The liquidity tier to create or update. */
    liquidityTier?: LiquidityTier;
}
/**
 * MsgSetLiquidityTier is a message used by x/gov to create or update a
 * liquidity tier.
 */
export interface MsgSetLiquidityTierSDKType {
    authority: string;
    liquidity_tier?: LiquidityTierSDKType;
}
/** MsgSetLiquidityTierResponse defines the SetLiquidityTier response type. */
export interface MsgSetLiquidityTierResponse {
}
/** MsgSetLiquidityTierResponse defines the SetLiquidityTier response type. */
export interface MsgSetLiquidityTierResponseSDKType {
}
/**
 * MsgUpdatePerpetualParams is a message used by x/gov to update the parameters
 * of a perpetual.
 */
export interface MsgUpdatePerpetualParams {
    authority: string;
    /** The perpetual to update. Each field must be set. */
    perpetualParams?: PerpetualParams;
}
/**
 * MsgUpdatePerpetualParams is a message used by x/gov to update the parameters
 * of a perpetual.
 */
export interface MsgUpdatePerpetualParamsSDKType {
    authority: string;
    perpetual_params?: PerpetualParamsSDKType;
}
/**
 * MsgUpdatePerpetualParamsResponse defines the UpdatePerpetualParams
 * response type.
 */
export interface MsgUpdatePerpetualParamsResponse {
}
/**
 * MsgUpdatePerpetualParamsResponse defines the UpdatePerpetualParams
 * response type.
 */
export interface MsgUpdatePerpetualParamsResponseSDKType {
}
/**
 * FundingPremium represents a funding premium value for a perpetual
 * market. Can be used to represent a premium vote or a premium sample.
 */
export interface FundingPremium {
    /** The id of the perpetual market. */
    perpetualId: number;
    /** The sampled premium rate. In parts-per-million. */
    premiumPpm: number;
}
/**
 * FundingPremium represents a funding premium value for a perpetual
 * market. Can be used to represent a premium vote or a premium sample.
 */
export interface FundingPremiumSDKType {
    perpetual_id: number;
    premium_ppm: number;
}
/** MsgAddPremiumVotes is a request type for the AddPremiumVotes method. */
export interface MsgAddPremiumVotes {
    votes: FundingPremium[];
}
/** MsgAddPremiumVotes is a request type for the AddPremiumVotes method. */
export interface MsgAddPremiumVotesSDKType {
    votes: FundingPremiumSDKType[];
}
/**
 * MsgAddPremiumVotesResponse defines the AddPremiumVotes
 * response type.
 */
export interface MsgAddPremiumVotesResponse {
}
/**
 * MsgAddPremiumVotesResponse defines the AddPremiumVotes
 * response type.
 */
export interface MsgAddPremiumVotesResponseSDKType {
}
/**
 * MsgUpdateParams is a message used by x/gov to update the parameters of the
 * perpetuals module.
 */
export interface MsgUpdateParams {
    authority: string;
    /** The parameters to update. Each field must be set. */
    params?: Params;
}
/**
 * MsgUpdateParams is a message used by x/gov to update the parameters of the
 * perpetuals module.
 */
export interface MsgUpdateParamsSDKType {
    authority: string;
    params?: ParamsSDKType;
}
/** MsgUpdateParamsResponse defines the UpdateParams response type. */
export interface MsgUpdateParamsResponse {
}
/** MsgUpdateParamsResponse defines the UpdateParams response type. */
export interface MsgUpdateParamsResponseSDKType {
}
export declare const MsgCreatePerpetual: {
    encode(message: MsgCreatePerpetual, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreatePerpetual;
    fromPartial(object: DeepPartial<MsgCreatePerpetual>): MsgCreatePerpetual;
};
export declare const MsgCreatePerpetualResponse: {
    encode(_: MsgCreatePerpetualResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreatePerpetualResponse;
    fromPartial(_: DeepPartial<MsgCreatePerpetualResponse>): MsgCreatePerpetualResponse;
};
export declare const MsgSetLiquidityTier: {
    encode(message: MsgSetLiquidityTier, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetLiquidityTier;
    fromPartial(object: DeepPartial<MsgSetLiquidityTier>): MsgSetLiquidityTier;
};
export declare const MsgSetLiquidityTierResponse: {
    encode(_: MsgSetLiquidityTierResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetLiquidityTierResponse;
    fromPartial(_: DeepPartial<MsgSetLiquidityTierResponse>): MsgSetLiquidityTierResponse;
};
export declare const MsgUpdatePerpetualParams: {
    encode(message: MsgUpdatePerpetualParams, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePerpetualParams;
    fromPartial(object: DeepPartial<MsgUpdatePerpetualParams>): MsgUpdatePerpetualParams;
};
export declare const MsgUpdatePerpetualParamsResponse: {
    encode(_: MsgUpdatePerpetualParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePerpetualParamsResponse;
    fromPartial(_: DeepPartial<MsgUpdatePerpetualParamsResponse>): MsgUpdatePerpetualParamsResponse;
};
export declare const FundingPremium: {
    encode(message: FundingPremium, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): FundingPremium;
    fromPartial(object: DeepPartial<FundingPremium>): FundingPremium;
};
export declare const MsgAddPremiumVotes: {
    encode(message: MsgAddPremiumVotes, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAddPremiumVotes;
    fromPartial(object: DeepPartial<MsgAddPremiumVotes>): MsgAddPremiumVotes;
};
export declare const MsgAddPremiumVotesResponse: {
    encode(_: MsgAddPremiumVotesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAddPremiumVotesResponse;
    fromPartial(_: DeepPartial<MsgAddPremiumVotesResponse>): MsgAddPremiumVotesResponse;
};
export declare const MsgUpdateParams: {
    encode(message: MsgUpdateParams, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateParams;
    fromPartial(object: DeepPartial<MsgUpdateParams>): MsgUpdateParams;
};
export declare const MsgUpdateParamsResponse: {
    encode(_: MsgUpdateParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateParamsResponse;
    fromPartial(_: DeepPartial<MsgUpdateParamsResponse>): MsgUpdateParamsResponse;
};
