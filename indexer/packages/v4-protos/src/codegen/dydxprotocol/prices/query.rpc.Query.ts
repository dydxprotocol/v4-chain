import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryMarketPriceRequest, QueryMarketPriceResponse, QueryAllMarketPricesRequest, QueryAllMarketPricesResponse, QueryMarketParamRequest, QueryMarketParamResponse, QueryAllMarketParamsRequest, QueryAllMarketParamsResponse, QueryNextMarketIdRequest, QueryNextMarketIdResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries a MarketPrice by id. */
  marketPrice(request: QueryMarketPriceRequest): Promise<QueryMarketPriceResponse>;
  /** Queries a list of MarketPrice items. */

  allMarketPrices(request?: QueryAllMarketPricesRequest): Promise<QueryAllMarketPricesResponse>;
  /** Queries a MarketParam by id. */

  marketParam(request: QueryMarketParamRequest): Promise<QueryMarketParamResponse>;
  /** Queries a list of MarketParam items. */

  allMarketParams(request?: QueryAllMarketParamsRequest): Promise<QueryAllMarketParamsResponse>;
  /** Queries the next market id. */

  nextMarketId(request?: QueryNextMarketIdRequest): Promise<QueryNextMarketIdResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.marketPrice = this.marketPrice.bind(this);
    this.allMarketPrices = this.allMarketPrices.bind(this);
    this.marketParam = this.marketParam.bind(this);
    this.allMarketParams = this.allMarketParams.bind(this);
    this.nextMarketId = this.nextMarketId.bind(this);
  }

  marketPrice(request: QueryMarketPriceRequest): Promise<QueryMarketPriceResponse> {
    const data = QueryMarketPriceRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.prices.Query", "MarketPrice", data);
    return promise.then(data => QueryMarketPriceResponse.decode(new _m0.Reader(data)));
  }

  allMarketPrices(request: QueryAllMarketPricesRequest = {
    pagination: undefined
  }): Promise<QueryAllMarketPricesResponse> {
    const data = QueryAllMarketPricesRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.prices.Query", "AllMarketPrices", data);
    return promise.then(data => QueryAllMarketPricesResponse.decode(new _m0.Reader(data)));
  }

  marketParam(request: QueryMarketParamRequest): Promise<QueryMarketParamResponse> {
    const data = QueryMarketParamRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.prices.Query", "MarketParam", data);
    return promise.then(data => QueryMarketParamResponse.decode(new _m0.Reader(data)));
  }

  allMarketParams(request: QueryAllMarketParamsRequest = {
    pagination: undefined
  }): Promise<QueryAllMarketParamsResponse> {
    const data = QueryAllMarketParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.prices.Query", "AllMarketParams", data);
    return promise.then(data => QueryAllMarketParamsResponse.decode(new _m0.Reader(data)));
  }

  nextMarketId(request: QueryNextMarketIdRequest = {}): Promise<QueryNextMarketIdResponse> {
    const data = QueryNextMarketIdRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.prices.Query", "NextMarketId", data);
    return promise.then(data => QueryNextMarketIdResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    marketPrice(request: QueryMarketPriceRequest): Promise<QueryMarketPriceResponse> {
      return queryService.marketPrice(request);
    },

    allMarketPrices(request?: QueryAllMarketPricesRequest): Promise<QueryAllMarketPricesResponse> {
      return queryService.allMarketPrices(request);
    },

    marketParam(request: QueryMarketParamRequest): Promise<QueryMarketParamResponse> {
      return queryService.marketParam(request);
    },

    allMarketParams(request?: QueryAllMarketParamsRequest): Promise<QueryAllMarketParamsResponse> {
      return queryService.allMarketParams(request);
    },

    nextMarketId(request?: QueryNextMarketIdRequest): Promise<QueryNextMarketIdResponse> {
      return queryService.nextMarketId(request);
    }

  };
};