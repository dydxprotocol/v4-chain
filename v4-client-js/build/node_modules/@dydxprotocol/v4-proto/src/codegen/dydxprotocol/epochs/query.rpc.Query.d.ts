import { Rpc } from "../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryGetEpochInfoRequest, QueryEpochInfoResponse, QueryAllEpochInfoRequest, QueryEpochInfoAllResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /** Queries a EpochInfo by name. */
    epochInfo(request: QueryGetEpochInfoRequest): Promise<QueryEpochInfoResponse>;
    /** Queries a list of EpochInfo items. */
    epochInfoAll(request?: QueryAllEpochInfoRequest): Promise<QueryEpochInfoAllResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    epochInfo(request: QueryGetEpochInfoRequest): Promise<QueryEpochInfoResponse>;
    epochInfoAll(request?: QueryAllEpochInfoRequest): Promise<QueryEpochInfoAllResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    epochInfo(request: QueryGetEpochInfoRequest): Promise<QueryEpochInfoResponse>;
    epochInfoAll(request?: QueryAllEpochInfoRequest): Promise<QueryEpochInfoAllResponse>;
};
