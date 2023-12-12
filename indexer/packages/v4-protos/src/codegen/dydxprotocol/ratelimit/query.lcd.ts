import { LCDClient } from "@osmonauts/lcd";
import { ListLimitParamsRequest, ListLimitParamsResponseSDKType, QueryCapacityByDenomRequest, QueryCapacityByDenomResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.listLimitParams = this.listLimitParams.bind(this);
    this.capacityByDenom = this.capacityByDenom.bind(this);
  }
  /* List all limit params. */


  async listLimitParams(_params: ListLimitParamsRequest = {}): Promise<ListLimitParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/ratelimit/list_limit_params`;
    return await this.req.get<ListLimitParamsResponseSDKType>(endpoint);
  }
  /* Query capacity by denom. */


  async capacityByDenom(params: QueryCapacityByDenomRequest): Promise<QueryCapacityByDenomResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.denom !== "undefined") {
      options.params.denom = params.denom;
    }

    const endpoint = `dydxprotocol/v4/ratelimit/capacity_by_denom`;
    return await this.req.get<QueryCapacityByDenomResponseSDKType>(endpoint, options);
  }

}