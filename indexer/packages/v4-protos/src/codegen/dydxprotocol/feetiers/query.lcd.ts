import { LCDClient } from "@osmonauts/lcd";
import { QueryPerpetualFeeParamsRequest, QueryPerpetualFeeParamsResponseSDKType, QueryUserFeeTierRequest, QueryUserFeeTierResponseSDKType, QueryPerMarketFeeDiscountParamsRequest, QueryPerMarketFeeDiscountParamsResponseSDKType, QueryAllMarketFeeDiscountParamsRequest, QueryAllMarketFeeDiscountParamsResponseSDKType, QueryStakingTiersRequest, QueryStakingTiersResponseSDKType, QueryUserStakingTierRequest, QueryUserStakingTierResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.perpetualFeeParams = this.perpetualFeeParams.bind(this);
    this.userFeeTier = this.userFeeTier.bind(this);
    this.perMarketFeeDiscountParams = this.perMarketFeeDiscountParams.bind(this);
    this.allMarketFeeDiscountParams = this.allMarketFeeDiscountParams.bind(this);
    this.stakingTiers = this.stakingTiers.bind(this);
    this.userStakingTier = this.userStakingTier.bind(this);
  }
  /* Queries the PerpetualFeeParams. */


  async perpetualFeeParams(_params: QueryPerpetualFeeParamsRequest = {}): Promise<QueryPerpetualFeeParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/feetiers/perpetual_fee_params`;
    return await this.req.get<QueryPerpetualFeeParamsResponseSDKType>(endpoint);
  }
  /* Queries a user's fee tier */


  async userFeeTier(params: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.user !== "undefined") {
      options.params.user = params.user;
    }

    const endpoint = `dydxprotocol/v4/feetiers/user_fee_tier`;
    return await this.req.get<QueryUserFeeTierResponseSDKType>(endpoint, options);
  }
  /* PerMarketFeeDiscountParams queries fee discount parameters for a
   specific market/CLOB pair. */


  async perMarketFeeDiscountParams(params: QueryPerMarketFeeDiscountParamsRequest): Promise<QueryPerMarketFeeDiscountParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/feetiers/market_fee_discount_params/${params.clobPairId}`;
    return await this.req.get<QueryPerMarketFeeDiscountParamsResponseSDKType>(endpoint);
  }
  /* AllMarketFeeDiscountParams queries all per-market fee discount parameters. */


  async allMarketFeeDiscountParams(_params: QueryAllMarketFeeDiscountParamsRequest = {}): Promise<QueryAllMarketFeeDiscountParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/feetiers/market_fee_discount_params`;
    return await this.req.get<QueryAllMarketFeeDiscountParamsResponseSDKType>(endpoint);
  }
  /* Get all staking tiers */


  async stakingTiers(_params: QueryStakingTiersRequest = {}): Promise<QueryStakingTiersResponseSDKType> {
    const endpoint = `dydxprotocol/v4/feetiers/staking_tiers`;
    return await this.req.get<QueryStakingTiersResponseSDKType>(endpoint);
  }
  /* Get user's current staked amount and staking tier */


  async userStakingTier(params: QueryUserStakingTierRequest): Promise<QueryUserStakingTierResponseSDKType> {
    const endpoint = `dydxprotocol/v4/feetiers/user_staking_tier/${params.address}`;
    return await this.req.get<QueryUserStakingTierResponseSDKType>(endpoint);
  }

}