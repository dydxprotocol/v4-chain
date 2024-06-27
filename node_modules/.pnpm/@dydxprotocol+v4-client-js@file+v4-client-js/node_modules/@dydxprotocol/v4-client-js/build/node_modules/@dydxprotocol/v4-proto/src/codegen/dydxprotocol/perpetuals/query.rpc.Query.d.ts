import { Rpc } from "../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryPerpetualRequest, QueryPerpetualResponse, QueryAllPerpetualsRequest, QueryAllPerpetualsResponse, QueryAllLiquidityTiersRequest, QueryAllLiquidityTiersResponse, QueryPremiumVotesRequest, QueryPremiumVotesResponse, QueryPremiumSamplesRequest, QueryPremiumSamplesResponse, QueryParamsRequest, QueryParamsResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /** Queries a Perpetual by id. */
    perpetual(request: QueryPerpetualRequest): Promise<QueryPerpetualResponse>;
    /** Queries a list of Perpetual items. */
    allPerpetuals(request?: QueryAllPerpetualsRequest): Promise<QueryAllPerpetualsResponse>;
    /** Queries a list of LiquidityTiers. */
    allLiquidityTiers(request?: QueryAllLiquidityTiersRequest): Promise<QueryAllLiquidityTiersResponse>;
    /** Queries a list of premium votes. */
    premiumVotes(request?: QueryPremiumVotesRequest): Promise<QueryPremiumVotesResponse>;
    /** Queries a list of premium samples. */
    premiumSamples(request?: QueryPremiumSamplesRequest): Promise<QueryPremiumSamplesResponse>;
    /** Queries the perpetual params. */
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    perpetual(request: QueryPerpetualRequest): Promise<QueryPerpetualResponse>;
    allPerpetuals(request?: QueryAllPerpetualsRequest): Promise<QueryAllPerpetualsResponse>;
    allLiquidityTiers(request?: QueryAllLiquidityTiersRequest): Promise<QueryAllLiquidityTiersResponse>;
    premiumVotes(request?: QueryPremiumVotesRequest): Promise<QueryPremiumVotesResponse>;
    premiumSamples(request?: QueryPremiumSamplesRequest): Promise<QueryPremiumSamplesResponse>;
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    perpetual(request: QueryPerpetualRequest): Promise<QueryPerpetualResponse>;
    allPerpetuals(request?: QueryAllPerpetualsRequest): Promise<QueryAllPerpetualsResponse>;
    allLiquidityTiers(request?: QueryAllLiquidityTiersRequest): Promise<QueryAllLiquidityTiersResponse>;
    premiumVotes(request?: QueryPremiumVotesRequest): Promise<QueryPremiumVotesResponse>;
    premiumSamples(request?: QueryPremiumSamplesRequest): Promise<QueryPremiumSamplesResponse>;
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
};
