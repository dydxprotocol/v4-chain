import { LCDClient } from "@osmonauts/lcd";
import { QueryGetClobPairRequest, QueryClobPairResponseSDKType, QueryAllClobPairRequest, QueryClobPairAllResponseSDKType, QueryEquityTierLimitConfigurationRequest, QueryEquityTierLimitConfigurationResponseSDKType, QueryBlockRateLimitConfigurationRequest, QueryBlockRateLimitConfigurationResponseSDKType, QueryLiquidationsConfigurationRequest, QueryLiquidationsConfigurationResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    clobPair(params: QueryGetClobPairRequest): Promise<QueryClobPairResponseSDKType>;
    clobPairAll(params?: QueryAllClobPairRequest): Promise<QueryClobPairAllResponseSDKType>;
    equityTierLimitConfiguration(_params?: QueryEquityTierLimitConfigurationRequest): Promise<QueryEquityTierLimitConfigurationResponseSDKType>;
    blockRateLimitConfiguration(_params?: QueryBlockRateLimitConfigurationRequest): Promise<QueryBlockRateLimitConfigurationResponseSDKType>;
    liquidationsConfiguration(_params?: QueryLiquidationsConfigurationRequest): Promise<QueryLiquidationsConfigurationResponseSDKType>;
}
