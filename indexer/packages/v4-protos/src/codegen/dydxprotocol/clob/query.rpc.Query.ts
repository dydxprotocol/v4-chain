import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryGetClobPairRequest, QueryClobPairResponse, QueryAllClobPairRequest, QueryClobPairAllResponse, AreSubaccountsLiquidatableRequest, AreSubaccountsLiquidatableResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries a ClobPair by id. */
  clobPair(request: QueryGetClobPairRequest): Promise<QueryClobPairResponse>;
  /** Queries a list of ClobPair items. */

  clobPairAll(request?: QueryAllClobPairRequest): Promise<QueryClobPairAllResponse>;
  /** Returns whether a subaccount is liquidatable. */

  areSubaccountsLiquidatable(request: AreSubaccountsLiquidatableRequest): Promise<AreSubaccountsLiquidatableResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.clobPair = this.clobPair.bind(this);
    this.clobPairAll = this.clobPairAll.bind(this);
    this.areSubaccountsLiquidatable = this.areSubaccountsLiquidatable.bind(this);
  }

  clobPair(request: QueryGetClobPairRequest): Promise<QueryClobPairResponse> {
    const data = QueryGetClobPairRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "ClobPair", data);
    return promise.then(data => QueryClobPairResponse.decode(new _m0.Reader(data)));
  }

  clobPairAll(request: QueryAllClobPairRequest = {
    pagination: undefined
  }): Promise<QueryClobPairAllResponse> {
    const data = QueryAllClobPairRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "ClobPairAll", data);
    return promise.then(data => QueryClobPairAllResponse.decode(new _m0.Reader(data)));
  }

  areSubaccountsLiquidatable(request: AreSubaccountsLiquidatableRequest): Promise<AreSubaccountsLiquidatableResponse> {
    const data = AreSubaccountsLiquidatableRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "AreSubaccountsLiquidatable", data);
    return promise.then(data => AreSubaccountsLiquidatableResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    clobPair(request: QueryGetClobPairRequest): Promise<QueryClobPairResponse> {
      return queryService.clobPair(request);
    },

    clobPairAll(request?: QueryAllClobPairRequest): Promise<QueryClobPairAllResponse> {
      return queryService.clobPairAll(request);
    },

    areSubaccountsLiquidatable(request: AreSubaccountsLiquidatableRequest): Promise<AreSubaccountsLiquidatableResponse> {
      return queryService.areSubaccountsLiquidatable(request);
    }

  };
};