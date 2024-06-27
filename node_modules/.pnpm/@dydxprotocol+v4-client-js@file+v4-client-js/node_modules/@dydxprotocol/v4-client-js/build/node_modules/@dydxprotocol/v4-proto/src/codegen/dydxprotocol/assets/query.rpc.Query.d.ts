import { Rpc } from "../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryAssetRequest, QueryAssetResponse, QueryAllAssetsRequest, QueryAllAssetsResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /** Queries a Asset by id. */
    asset(request: QueryAssetRequest): Promise<QueryAssetResponse>;
    /** Queries a list of Asset items. */
    allAssets(request?: QueryAllAssetsRequest): Promise<QueryAllAssetsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    asset(request: QueryAssetRequest): Promise<QueryAssetResponse>;
    allAssets(request?: QueryAllAssetsRequest): Promise<QueryAllAssetsResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    asset(request: QueryAssetRequest): Promise<QueryAssetResponse>;
    allAssets(request?: QueryAllAssetsRequest): Promise<QueryAllAssetsResponse>;
};
