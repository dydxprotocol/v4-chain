import { LCDClient } from "@osmonauts/lcd";
import { QueryParamsRequest, QueryParamsResponseSDKType, QueryStatsMetadataRequest, QueryStatsMetadataResponseSDKType, QueryGlobalStatsRequest, QueryGlobalStatsResponseSDKType, QueryUserStatsRequest, QueryUserStatsResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    params(_params?: QueryParamsRequest): Promise<QueryParamsResponseSDKType>;
    statsMetadata(_params?: QueryStatsMetadataRequest): Promise<QueryStatsMetadataResponseSDKType>;
    globalStats(_params?: QueryGlobalStatsRequest): Promise<QueryGlobalStatsResponseSDKType>;
    userStats(params: QueryUserStatsRequest): Promise<QueryUserStatsResponseSDKType>;
}
