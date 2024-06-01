import { Rpc } from "../../helpers";
import { QueryClient } from "@cosmjs/stargate";
/** Query defines the gRPC querier service. */
export interface Query {
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
}
export declare const createRpcQueryExtension: (base: QueryClient) => {};
