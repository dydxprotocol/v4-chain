import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryPerpetualFeeParamsRequest, QueryPerpetualFeeParamsResponse, QueryUserFeeTierRequest, QueryUserFeeTierResponse, QueryPerMarketFeeDiscountParamsRequest, QueryPerMarketFeeDiscountParamsResponse, QueryAllMarketFeeDiscountParamsRequest, QueryAllMarketFeeDiscountParamsResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the PerpetualFeeParams. */
  perpetualFeeParams(request?: QueryPerpetualFeeParamsRequest): Promise<QueryPerpetualFeeParamsResponse>;
  /** Queries a user's fee tier */

  userFeeTier(request: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponse>;
  /**
   * PerMarketFeeDiscountParams queries fee discount parameters for a
   * specific market/CLOB pair.
   */

  perMarketFeeDiscountParams(request: QueryPerMarketFeeDiscountParamsRequest): Promise<QueryPerMarketFeeDiscountParamsResponse>;
  /** AllMarketFeeDiscountParams queries all per-market fee discount parameters. */

  allMarketFeeDiscountParams(request?: QueryAllMarketFeeDiscountParamsRequest): Promise<QueryAllMarketFeeDiscountParamsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.perpetualFeeParams = this.perpetualFeeParams.bind(this);
    this.userFeeTier = this.userFeeTier.bind(this);
    this.perMarketFeeDiscountParams = this.perMarketFeeDiscountParams.bind(this);
    this.allMarketFeeDiscountParams = this.allMarketFeeDiscountParams.bind(this);
  }

  perpetualFeeParams(request: QueryPerpetualFeeParamsRequest = {}): Promise<QueryPerpetualFeeParamsResponse> {
    const data = QueryPerpetualFeeParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "PerpetualFeeParams", data);
    return promise.then(data => QueryPerpetualFeeParamsResponse.decode(new _m0.Reader(data)));
  }

  userFeeTier(request: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponse> {
    const data = QueryUserFeeTierRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "UserFeeTier", data);
    return promise.then(data => QueryUserFeeTierResponse.decode(new _m0.Reader(data)));
  }

  perMarketFeeDiscountParams(request: QueryPerMarketFeeDiscountParamsRequest): Promise<QueryPerMarketFeeDiscountParamsResponse> {
    const data = QueryPerMarketFeeDiscountParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "PerMarketFeeDiscountParams", data);
    return promise.then(data => QueryPerMarketFeeDiscountParamsResponse.decode(new _m0.Reader(data)));
  }

  allMarketFeeDiscountParams(request: QueryAllMarketFeeDiscountParamsRequest = {}): Promise<QueryAllMarketFeeDiscountParamsResponse> {
    const data = QueryAllMarketFeeDiscountParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "AllMarketFeeDiscountParams", data);
    return promise.then(data => QueryAllMarketFeeDiscountParamsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    perpetualFeeParams(request?: QueryPerpetualFeeParamsRequest): Promise<QueryPerpetualFeeParamsResponse> {
      return queryService.perpetualFeeParams(request);
    },

    userFeeTier(request: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponse> {
      return queryService.userFeeTier(request);
    },

    perMarketFeeDiscountParams(request: QueryPerMarketFeeDiscountParamsRequest): Promise<QueryPerMarketFeeDiscountParamsResponse> {
      return queryService.perMarketFeeDiscountParams(request);
    },

    allMarketFeeDiscountParams(request?: QueryAllMarketFeeDiscountParamsRequest): Promise<QueryAllMarketFeeDiscountParamsResponse> {
      return queryService.allMarketFeeDiscountParams(request);
    }

  };
};