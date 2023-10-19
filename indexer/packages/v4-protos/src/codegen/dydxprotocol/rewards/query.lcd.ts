import { LCDClient } from "@osmonauts/lcd";
import { QueryParamsRequest, QueryParamsResponseSDKType, QueryRewardShareRequest, QueryRewardShareResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.params = this.params.bind(this);
    this.rewardShare = this.rewardShare.bind(this);
  }
  /* Queries the Params. */


  async params(_params: QueryParamsRequest = {}): Promise<QueryParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/rewards/params`;
    return await this.req.get<QueryParamsResponseSDKType>(endpoint);
  }
  /* Queries a reward share by address. */


  async rewardShare(params: QueryRewardShareRequest): Promise<QueryRewardShareResponseSDKType> {
    const endpoint = `dydxprotocol/v4/rewards/shares/${params.address}`;
    return await this.req.get<QueryRewardShareResponseSDKType>(endpoint);
  }

}