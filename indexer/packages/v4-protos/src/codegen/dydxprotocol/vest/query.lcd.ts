import { LCDClient } from "@osmonauts/lcd";
import { QueryVestEntryRequest, QueryVestEntryResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.vestEntry = this.vestEntry.bind(this);
  }
  /* Queries the VestEntry. */


  async vestEntry(params: QueryVestEntryRequest): Promise<QueryVestEntryResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.vesterAccount !== "undefined") {
      options.params.vester_account = params.vesterAccount;
    }

    const endpoint = `dydxprotocol/v4/vest/vest_entry`;
    return await this.req.get<QueryVestEntryResponseSDKType>(endpoint, options);
  }

}