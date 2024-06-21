import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryMarketMapperRevShareDetails, QueryMarketMapperRevShareDetailsResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries market mapper revenue share details for a specific market */
  marketMapperRevShareDetails(request: QueryMarketMapperRevShareDetails): Promise<QueryMarketMapperRevShareDetailsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.marketMapperRevShareDetails = this.marketMapperRevShareDetails.bind(this);
  }

  marketMapperRevShareDetails(request: QueryMarketMapperRevShareDetails): Promise<QueryMarketMapperRevShareDetailsResponse> {
    const data = QueryMarketMapperRevShareDetails.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.revshare.Query", "MarketMapperRevShareDetails", data);
    return promise.then(data => QueryMarketMapperRevShareDetailsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    marketMapperRevShareDetails(request: QueryMarketMapperRevShareDetails): Promise<QueryMarketMapperRevShareDetailsResponse> {
      return queryService.marketMapperRevShareDetails(request);
    }

  };
};