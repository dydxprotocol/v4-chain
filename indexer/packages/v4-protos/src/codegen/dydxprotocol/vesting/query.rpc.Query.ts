import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryVestingEntryRequest, QueryVestingEntryResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the VestingEntry. */
  vestingEntry(request: QueryVestingEntryRequest): Promise<QueryVestingEntryResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.vestingEntry = this.vestingEntry.bind(this);
  }

  vestingEntry(request: QueryVestingEntryRequest): Promise<QueryVestingEntryResponse> {
    const data = QueryVestingEntryRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vesting.Query", "VestingEntry", data);
    return promise.then(data => QueryVestingEntryResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    vestingEntry(request: QueryVestingEntryRequest): Promise<QueryVestingEntryResponse> {
      return queryService.vestingEntry(request);
    }

  };
};