import { LCDClient } from "@osmonauts/lcd";
import { ListLimitParamsRequest, ListLimitParamsResponseSDKType, QueryCapacityByDenomRequest, QueryCapacityByDenomResponseSDKType, QueryAllPendingSendPacketsRequest, QueryAllPendingSendPacketsResponseSDKType, GetSDAIPriceQueryRequest, GetSDAIPriceQueryResponseSDKType, GetAssetYieldIndexQueryRequest, GetAssetYieldIndexQueryResponseSDKType } from "./query";
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
    this.allPendingSendPackets = this.allPendingSendPackets.bind(this);
    this.getSDAIPriceQuery = this.getSDAIPriceQuery.bind(this);
    this.getAssetYieldIndexQuery = this.getAssetYieldIndexQuery.bind(this);
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
  /* Get all pending send packets */


  async allPendingSendPackets(_params: QueryAllPendingSendPacketsRequest = {}): Promise<QueryAllPendingSendPacketsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/ratelimit/get_all_pending_send_packet`;
    return await this.req.get<QueryAllPendingSendPacketsResponseSDKType>(endpoint);
  }
  /* Get the price of sDAI. */


  async getSDAIPriceQuery(_params: GetSDAIPriceQueryRequest = {}): Promise<GetSDAIPriceQueryResponseSDKType> {
    const endpoint = `dydxprotocol/v4/ratelimit/get_sdai_price`;
    return await this.req.get<GetSDAIPriceQueryResponseSDKType>(endpoint);
  }
  /* Get the price of sDAI. */


  async getAssetYieldIndexQuery(_params: GetAssetYieldIndexQueryRequest = {}): Promise<GetAssetYieldIndexQueryResponseSDKType> {
    const endpoint = `dydxprotocol/v4/ratelimit/get_asset_yield_index`;
    return await this.req.get<GetAssetYieldIndexQueryResponseSDKType>(endpoint);
  }

}