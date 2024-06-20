import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryMarketMapperRevenueShareParams, QueryMarketMapperRevenueShareParamsResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** MarketMapperRevenueShareParams queries the revenue share params for the market mapper */
  marketMapperRevenueShareParams(request?: QueryMarketMapperRevenueShareParams): Promise<QueryMarketMapperRevenueShareParamsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.marketMapperRevenueShareParams = this.marketMapperRevenueShareParams.bind(this);
  }

  marketMapperRevenueShareParams(request: QueryMarketMapperRevenueShareParams = {}): Promise<QueryMarketMapperRevenueShareParamsResponse> {
    const data = QueryMarketMapperRevenueShareParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.revshare.Query", "MarketMapperRevenueShareParams", data);
    return promise.then(data => QueryMarketMapperRevenueShareParamsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    marketMapperRevenueShareParams(request?: QueryMarketMapperRevenueShareParams): Promise<QueryMarketMapperRevenueShareParamsResponse> {
      return queryService.marketMapperRevenueShareParams(request);
    }

  };
};