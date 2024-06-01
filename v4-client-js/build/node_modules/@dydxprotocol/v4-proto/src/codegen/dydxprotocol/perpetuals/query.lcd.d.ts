import { LCDClient } from "@osmonauts/lcd";
import { QueryPerpetualRequest, QueryPerpetualResponseSDKType, QueryAllPerpetualsRequest, QueryAllPerpetualsResponseSDKType, QueryAllLiquidityTiersRequest, QueryAllLiquidityTiersResponseSDKType, QueryPremiumVotesRequest, QueryPremiumVotesResponseSDKType, QueryPremiumSamplesRequest, QueryPremiumSamplesResponseSDKType, QueryParamsRequest, QueryParamsResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    perpetual(params: QueryPerpetualRequest): Promise<QueryPerpetualResponseSDKType>;
    allPerpetuals(params?: QueryAllPerpetualsRequest): Promise<QueryAllPerpetualsResponseSDKType>;
    allLiquidityTiers(params?: QueryAllLiquidityTiersRequest): Promise<QueryAllLiquidityTiersResponseSDKType>;
    premiumVotes(_params?: QueryPremiumVotesRequest): Promise<QueryPremiumVotesResponseSDKType>;
    premiumSamples(_params?: QueryPremiumSamplesRequest): Promise<QueryPremiumSamplesResponseSDKType>;
    params(_params?: QueryParamsRequest): Promise<QueryParamsResponseSDKType>;
}
