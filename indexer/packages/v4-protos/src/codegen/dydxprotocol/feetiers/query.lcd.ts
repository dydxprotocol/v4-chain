import { LCDClient } from "@osmonauts/lcd";
import { QueryPerpetualFeeParamsRequest, QueryPerpetualFeeParamsResponseSDKType, QueryUserFeeTierRequest, QueryUserFeeTierResponseSDKType, QueryFeeHolidayParamsRequest, QueryFeeHolidayParamsResponseSDKType, QueryAllFeeHolidayParamsRequest, QueryAllFeeHolidayParamsResponseSDKType } from "./query";
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
    this.feeHolidayParams = this.feeHolidayParams.bind(this);
    this.allFeeHolidayParams = this.allFeeHolidayParams.bind(this);
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
  /* Queries the FeeHolidayParams */


  async feeHolidayParams(params: QueryFeeHolidayParamsRequest): Promise<QueryFeeHolidayParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/feetiers/fee_holiday_params/${params.clobPairId}`;
    return await this.req.get<QueryFeeHolidayParamsResponseSDKType>(endpoint);
  }
  /* Queries all fee holiday params */


  async allFeeHolidayParams(_params: QueryAllFeeHolidayParamsRequest = {}): Promise<QueryAllFeeHolidayParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/feetiers/fee_holiday_params`;
    return await this.req.get<QueryAllFeeHolidayParamsResponseSDKType>(endpoint);
  }

}