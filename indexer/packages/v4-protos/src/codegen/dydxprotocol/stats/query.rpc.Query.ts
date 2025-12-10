import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryStatsMetadataRequest, QueryStatsMetadataResponse, QueryGlobalStatsRequest, QueryGlobalStatsResponse, QueryUserStatsRequest, QueryUserStatsResponse, QueryEpochStatsRequest, QueryEpochStatsResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the Params. */
  params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries StatsMetadata. */

  statsMetadata(request?: QueryStatsMetadataRequest): Promise<QueryStatsMetadataResponse>;
  /** Queries GlobalStats. */

  globalStats(request?: QueryGlobalStatsRequest): Promise<QueryGlobalStatsResponse>;
  /** Queries UserStats. */

  userStats(request: QueryUserStatsRequest): Promise<QueryUserStatsResponse>;
  /** Queries EpochStats for a specific epoch. */

  epochStats(request: QueryEpochStatsRequest): Promise<QueryEpochStatsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.params = this.params.bind(this);
    this.statsMetadata = this.statsMetadata.bind(this);
    this.globalStats = this.globalStats.bind(this);
    this.userStats = this.userStats.bind(this);
    this.epochStats = this.epochStats.bind(this);
  }

  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.stats.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  statsMetadata(request: QueryStatsMetadataRequest = {}): Promise<QueryStatsMetadataResponse> {
    const data = QueryStatsMetadataRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.stats.Query", "StatsMetadata", data);
    return promise.then(data => QueryStatsMetadataResponse.decode(new _m0.Reader(data)));
  }

  globalStats(request: QueryGlobalStatsRequest = {}): Promise<QueryGlobalStatsResponse> {
    const data = QueryGlobalStatsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.stats.Query", "GlobalStats", data);
    return promise.then(data => QueryGlobalStatsResponse.decode(new _m0.Reader(data)));
  }

  userStats(request: QueryUserStatsRequest): Promise<QueryUserStatsResponse> {
    const data = QueryUserStatsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.stats.Query", "UserStats", data);
    return promise.then(data => QueryUserStatsResponse.decode(new _m0.Reader(data)));
  }

  epochStats(request: QueryEpochStatsRequest): Promise<QueryEpochStatsResponse> {
    const data = QueryEpochStatsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.stats.Query", "EpochStats", data);
    return promise.then(data => QueryEpochStatsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse> {
      return queryService.params(request);
    },

    statsMetadata(request?: QueryStatsMetadataRequest): Promise<QueryStatsMetadataResponse> {
      return queryService.statsMetadata(request);
    },

    globalStats(request?: QueryGlobalStatsRequest): Promise<QueryGlobalStatsResponse> {
      return queryService.globalStats(request);
    },

    userStats(request: QueryUserStatsRequest): Promise<QueryUserStatsResponse> {
      return queryService.userStats(request);
    },

    epochStats(request: QueryEpochStatsRequest): Promise<QueryEpochStatsResponse> {
      return queryService.epochStats(request);
    }

  };
};