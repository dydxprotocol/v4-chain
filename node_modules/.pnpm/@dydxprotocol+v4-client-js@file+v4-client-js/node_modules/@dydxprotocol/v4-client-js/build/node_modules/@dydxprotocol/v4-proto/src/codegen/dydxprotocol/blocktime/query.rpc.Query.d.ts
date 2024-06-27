import { Rpc } from "../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryDowntimeParamsRequest, QueryDowntimeParamsResponse, QueryPreviousBlockInfoRequest, QueryPreviousBlockInfoResponse, QueryAllDowntimeInfoRequest, QueryAllDowntimeInfoResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /** Queries the DowntimeParams. */
    downtimeParams(request?: QueryDowntimeParamsRequest): Promise<QueryDowntimeParamsResponse>;
    /** Queries the information of the previous block */
    previousBlockInfo(request?: QueryPreviousBlockInfoRequest): Promise<QueryPreviousBlockInfoResponse>;
    /** Queries all recorded downtime info. */
    allDowntimeInfo(request?: QueryAllDowntimeInfoRequest): Promise<QueryAllDowntimeInfoResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    downtimeParams(request?: QueryDowntimeParamsRequest): Promise<QueryDowntimeParamsResponse>;
    previousBlockInfo(request?: QueryPreviousBlockInfoRequest): Promise<QueryPreviousBlockInfoResponse>;
    allDowntimeInfo(request?: QueryAllDowntimeInfoRequest): Promise<QueryAllDowntimeInfoResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    downtimeParams(request?: QueryDowntimeParamsRequest): Promise<QueryDowntimeParamsResponse>;
    previousBlockInfo(request?: QueryPreviousBlockInfoRequest): Promise<QueryPreviousBlockInfoResponse>;
    allDowntimeInfo(request?: QueryAllDowntimeInfoRequest): Promise<QueryAllDowntimeInfoResponse>;
};
