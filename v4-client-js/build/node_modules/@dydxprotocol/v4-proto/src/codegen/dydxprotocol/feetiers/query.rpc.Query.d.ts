import { Rpc } from "../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryPerpetualFeeParamsRequest, QueryPerpetualFeeParamsResponse, QueryUserFeeTierRequest, QueryUserFeeTierResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /** Queries the PerpetualFeeParams. */
    perpetualFeeParams(request?: QueryPerpetualFeeParamsRequest): Promise<QueryPerpetualFeeParamsResponse>;
    /** Queries a user's fee tier */
    userFeeTier(request: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    perpetualFeeParams(request?: QueryPerpetualFeeParamsRequest): Promise<QueryPerpetualFeeParamsResponse>;
    userFeeTier(request: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    perpetualFeeParams(request?: QueryPerpetualFeeParamsRequest): Promise<QueryPerpetualFeeParamsResponse>;
    userFeeTier(request: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponse>;
};
