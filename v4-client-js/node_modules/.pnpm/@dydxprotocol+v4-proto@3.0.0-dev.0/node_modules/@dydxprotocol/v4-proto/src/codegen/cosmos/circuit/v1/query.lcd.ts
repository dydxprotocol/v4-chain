import { setPaginationParams } from "../../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryAccountRequest, AccountResponseSDKType, QueryAccountsRequest, AccountsResponseSDKType, QueryDisabledListRequest, DisabledListResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.account = this.account.bind(this);
    this.accounts = this.accounts.bind(this);
    this.disabledList = this.disabledList.bind(this);
  }
  /* Account returns account permissions. */


  async account(params: QueryAccountRequest): Promise<AccountResponseSDKType> {
    const endpoint = `cosmos/circuit/v1/accounts/${params.address}`;
    return await this.req.get<AccountResponseSDKType>(endpoint);
  }
  /* Account returns account permissions. */


  async accounts(params: QueryAccountsRequest = {
    pagination: undefined
  }): Promise<AccountsResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `cosmos/circuit/v1/accounts`;
    return await this.req.get<AccountsResponseSDKType>(endpoint, options);
  }
  /* DisabledList returns a list of disabled message urls */


  async disabledList(_params: QueryDisabledListRequest = {}): Promise<DisabledListResponseSDKType> {
    const endpoint = `cosmos/circuit/v1/disable_list`;
    return await this.req.get<DisabledListResponseSDKType>(endpoint);
  }

}