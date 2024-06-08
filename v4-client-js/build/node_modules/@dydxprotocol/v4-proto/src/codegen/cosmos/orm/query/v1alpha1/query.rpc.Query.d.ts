import { Rpc } from "../../../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { GetRequest, GetResponse, ListRequest, ListResponse } from "./query";
/** Query is a generic gRPC service for querying ORM data. */
export interface Query {
    /** Get queries an ORM table against an unique index. */
    get(request: GetRequest): Promise<GetResponse>;
    /** List queries an ORM table against an index. */
    list(request: ListRequest): Promise<ListResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    get(request: GetRequest): Promise<GetResponse>;
    list(request: ListRequest): Promise<ListResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    get(request: GetRequest): Promise<GetResponse>;
    list(request: ListRequest): Promise<ListResponse>;
};
