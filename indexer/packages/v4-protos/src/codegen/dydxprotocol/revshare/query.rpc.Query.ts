import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryMarketMapperRevenueShareParams, QueryMarketMapperRevenueShareParamsResponse, QueryMarketMapperRevShareDetails, QueryMarketMapperRevShareDetailsResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /**
   * MarketMapperRevenueShareParams queries the revenue share params for the
   * market mapper
   */
  marketMapperRevenueShareParams(request?: QueryMarketMapperRevenueShareParams): Promise<QueryMarketMapperRevenueShareParamsResponse>;
  /** Queries market mapper revenue share details for a specific market */

  marketMapperRevShareDetails(request: QueryMarketMapperRevShareDetails): Promise<QueryMarketMapperRevShareDetailsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.marketMapperRevenueShareParams = this.marketMapperRevenueShareParams.bind(this);
    this.marketMapperRevShareDetails = this.marketMapperRevShareDetails.bind(this);
  }

  marketMapperRevenueShareParams(request: QueryMarketMapperRevenueShareParams = {}): Promise<QueryMarketMapperRevenueShareParamsResponse> {
    const data = QueryMarketMapperRevenueShareParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.revshare.Query", "MarketMapperRevenueShareParams", data);
    return promise.then(data => QueryMarketMapperRevenueShareParamsResponse.decode(new _m0.Reader(data)));
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
    marketMapperRevenueShareParams(request?: QueryMarketMapperRevenueShareParams): Promise<QueryMarketMapperRevenueShareParamsResponse> {
      return queryService.marketMapperRevenueShareParams(request);
    },

    marketMapperRevShareDetails(request: QueryMarketMapperRevShareDetails): Promise<QueryMarketMapperRevShareDetailsResponse> {
      return queryService.marketMapperRevShareDetails(request);
    }

  };
};