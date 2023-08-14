import { LCDClient } from "@osmonauts/lcd";
import { QueryPerpetualFeeParamsRequest, QueryPerpetualFeeParamsResponseSDKType, QueryUserFeeTierRequest, QueryUserFeeTierResponseSDKType } from "./query";
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

}