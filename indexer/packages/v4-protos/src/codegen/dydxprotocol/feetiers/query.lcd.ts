import { LCDClient } from "@osmonauts/lcd";
import { QueryPerpetualFeeParamsRequest, QueryPerpetualFeeParamsResponseSDKType, QueryUserFeeTierRequest, QueryUserFeeTierResponseSDKType, QueryFeeDiscountCampaignParamsRequest, QueryFeeDiscountCampaignParamsResponseSDKType, QueryAllFeeDiscountCampaignParamsRequest, QueryAllFeeDiscountCampaignParamsResponseSDKType } from "./query";
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
    this.feeDiscountCampaignParams = this.feeDiscountCampaignParams.bind(this);
    this.allFeeDiscountCampaignParams = this.allFeeDiscountCampaignParams.bind(this);
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
  /* FeeDiscountCampaignParams queries fee discount campaign parameters for a
   specific CLOB pair. */


  async feeDiscountCampaignParams(params: QueryFeeDiscountCampaignParamsRequest): Promise<QueryFeeDiscountCampaignParamsResponseSDKType> {
    const endpoint = `dydxprotocol/feetiers/fee_discount_campaign_params/${params.clobPairId}`;
    return await this.req.get<QueryFeeDiscountCampaignParamsResponseSDKType>(endpoint);
  }
  /* AllFeeDiscountCampaignParams queries all fee discount campaign parameters. */


  async allFeeDiscountCampaignParams(_params: QueryAllFeeDiscountCampaignParamsRequest = {}): Promise<QueryAllFeeDiscountCampaignParamsResponseSDKType> {
    const endpoint = `dydxprotocol/feetiers/fee_discount_campaign_params`;
    return await this.req.get<QueryAllFeeDiscountCampaignParamsResponseSDKType>(endpoint);
  }

}