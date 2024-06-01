import { Rpc } from "../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryStatsMetadataRequest, QueryStatsMetadataResponse, QueryGlobalStatsRequest, QueryGlobalStatsResponse, QueryUserStatsRequest, QueryUserStatsResponse } from "./query";
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
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    statsMetadata(request?: QueryStatsMetadataRequest): Promise<QueryStatsMetadataResponse>;
    globalStats(request?: QueryGlobalStatsRequest): Promise<QueryGlobalStatsResponse>;
    userStats(request: QueryUserStatsRequest): Promise<QueryUserStatsResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
    statsMetadata(request?: QueryStatsMetadataRequest): Promise<QueryStatsMetadataResponse>;
    globalStats(request?: QueryGlobalStatsRequest): Promise<QueryGlobalStatsResponse>;
    userStats(request: QueryUserStatsRequest): Promise<QueryUserStatsResponse>;
};
