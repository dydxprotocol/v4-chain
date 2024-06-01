import { Rpc } from "../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryVestEntryRequest, QueryVestEntryResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /** Queries the VestEntry. */
    vestEntry(request: QueryVestEntryRequest): Promise<QueryVestEntryResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    vestEntry(request: QueryVestEntryRequest): Promise<QueryVestEntryResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    vestEntry(request: QueryVestEntryRequest): Promise<QueryVestEntryResponse>;
};
