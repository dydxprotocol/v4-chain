import { LCDClient } from "@osmonauts/lcd";
import { QueryParamsRequest, QueryParamsResponseSDKType, QueryStatsMetadataRequest, QueryStatsMetadataResponseSDKType, QueryGlobalStatsRequest, QueryGlobalStatsResponseSDKType, QueryUserStatsRequest, QueryUserStatsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.params = this.params.bind(this);
    this.statsMetadata = this.statsMetadata.bind(this);
    this.globalStats = this.globalStats.bind(this);
    this.userStats = this.userStats.bind(this);
  }
  /* Queries the Params. */


  async params(_params: QueryParamsRequest = {}): Promise<QueryParamsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/stats/params`;
    return await this.req.get<QueryParamsResponseSDKType>(endpoint);
  }
  /* Queries StatsMetadata. */


  async statsMetadata(_params: QueryStatsMetadataRequest = {}): Promise<QueryStatsMetadataResponseSDKType> {
    const endpoint = `dydxprotocol/v4/stats/stats_metadata`;
    return await this.req.get<QueryStatsMetadataResponseSDKType>(endpoint);
  }
  /* Queries GlobalStats. */


  async globalStats(_params: QueryGlobalStatsRequest = {}): Promise<QueryGlobalStatsResponseSDKType> {
    const endpoint = `dydxprotocol/v4/stats/global_stats`;
    return await this.req.get<QueryGlobalStatsResponseSDKType>(endpoint);
  }
  /* Queries UserStats. */


  async userStats(params: QueryUserStatsRequest): Promise<QueryUserStatsResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.user !== "undefined") {
      options.params.user = params.user;
    }

    const endpoint = `dydxprotocol/v4/stats/user_stats`;
    return await this.req.get<QueryUserStatsResponseSDKType>(endpoint, options);
  }

}