import { LCDClient } from "@osmonauts/lcd";
import { QueryDowntimeParamsRequest, QueryDowntimeParamsResponseSDKType, QueryAllDowntimeInfoRequest, QueryAllDowntimeInfoResponseSDKType, QuerySynchronyParamsRequest, QuerySynchronyParamsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.downtimeParams = this.downtimeParams.bind(this);
    this.allDowntimeInfo = this.allDowntimeInfo.bind(this);
    this.synchronyParams = this.synchronyParams.bind(this);
  }
  /* Queries the DowntimeParams. */


  async downtimeParams(_params: QueryDowntimeParamsRequest = {}): Promise<QueryDowntimeParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/blocktime/downtime_params`;
    return await this.req.get<QueryDowntimeParamsResponseSDKType>(endpoint);
  }
  /* Queries all recorded downtime info. */


  async allDowntimeInfo(_params: QueryAllDowntimeInfoRequest = {}): Promise<QueryAllDowntimeInfoResponseSDKType> {
    const endpoint = `dydxprotocol/v4/blocktime/all_downtime_info`;
    return await this.req.get<QueryAllDowntimeInfoResponseSDKType>(endpoint);
  }
  /* Queries the SynchronyParams. */


  async synchronyParams(_params: QuerySynchronyParamsRequest = {}): Promise<QuerySynchronyParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/blocktime/synchrony_params`;
    return await this.req.get<QuerySynchronyParamsResponseSDKType>(endpoint);
  }

}