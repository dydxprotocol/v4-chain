import { LCDClient } from "@osmonauts/lcd";
import { QueryAccountRequest, AccountResponseSDKType, QueryAccountsRequest, AccountsResponseSDKType, QueryDisabledListRequest, DisabledListResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    account(params: QueryAccountRequest): Promise<AccountResponseSDKType>;
    accounts(params?: QueryAccountsRequest): Promise<AccountsResponseSDKType>;
    disabledList(_params?: QueryDisabledListRequest): Promise<DisabledListResponseSDKType>;
}
