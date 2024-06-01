import { Rpc } from "../../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryAccountRequest, AccountResponse, QueryAccountsRequest, AccountsResponse, QueryDisabledListRequest, DisabledListResponse } from "./query";
/** Query defines the circuit gRPC querier service. */
export interface Query {
    /** Account returns account permissions. */
    account(request: QueryAccountRequest): Promise<AccountResponse>;
    /** Account returns account permissions. */
    accounts(request?: QueryAccountsRequest): Promise<AccountsResponse>;
    /** DisabledList returns a list of disabled message urls */
    disabledList(request?: QueryDisabledListRequest): Promise<DisabledListResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    account(request: QueryAccountRequest): Promise<AccountResponse>;
    accounts(request?: QueryAccountsRequest): Promise<AccountsResponse>;
    disabledList(request?: QueryDisabledListRequest): Promise<DisabledListResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    account(request: QueryAccountRequest): Promise<AccountResponse>;
    accounts(request?: QueryAccountsRequest): Promise<AccountsResponse>;
    disabledList(request?: QueryDisabledListRequest): Promise<DisabledListResponse>;
};
