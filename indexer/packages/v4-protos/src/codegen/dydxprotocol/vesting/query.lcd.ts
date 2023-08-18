import { LCDClient } from "@osmonauts/lcd";
import { QueryVestingEntryRequest, QueryVestingEntryResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.vestingEntry = this.vestingEntry.bind(this);
  }
  /* Queries the VestingEntry. */


  async vestingEntry(params: QueryVestingEntryRequest): Promise<QueryVestingEntryResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.vesterAccount !== "undefined") {
      options.params.vester_account = params.vesterAccount;
    }

    const endpoint = `dydxprotocol/v4/vesting/vesting_entry`;
    return await this.req.get<QueryVestingEntryResponseSDKType>(endpoint, options);
  }

}