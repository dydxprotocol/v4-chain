import { LCDClient } from "@osmonauts/lcd";
import { QueryDowntimeParamsRequest, QueryDowntimeParamsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.downtimeParams = this.downtimeParams.bind(this);
  }
  /* Queries the DowntimeParams. */


  async downtimeParams(_params: QueryDowntimeParamsRequest = {}): Promise<QueryDowntimeParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/blocktime/downtime_params`;
    return await this.req.get<QueryDowntimeParamsResponseSDKType>(endpoint);
  }

}