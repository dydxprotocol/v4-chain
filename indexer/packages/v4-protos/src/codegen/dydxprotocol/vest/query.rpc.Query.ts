import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryVestEntryRequest, QueryVestEntryResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the VestEntry. */
  vestEntry(request: QueryVestEntryRequest): Promise<QueryVestEntryResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.vestEntry = this.vestEntry.bind(this);
  }

  vestEntry(request: QueryVestEntryRequest): Promise<QueryVestEntryResponse> {
    const data = QueryVestEntryRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vest.Query", "VestEntry", data);
    return promise.then(data => QueryVestEntryResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    vestEntry(request: QueryVestEntryRequest): Promise<QueryVestEntryResponse> {
      return queryService.vestEntry(request);
    }

  };
};