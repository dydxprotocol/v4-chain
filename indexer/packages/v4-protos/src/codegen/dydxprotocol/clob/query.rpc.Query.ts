import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryGetClobPairRequest, QueryClobPairResponse, QueryAllClobPairRequest, QueryClobPairAllResponse, AreSubaccountsLiquidatableRequest, AreSubaccountsLiquidatableResponse, MevNodeToNodeCalculationRequest, MevNodeToNodeCalculationResponse, QueryEquityTierLimitConfigurationRequest, QueryEquityTierLimitConfigurationResponse, QueryAllStatefulOrdersRequest, QueryAllStatefulOrdersResponse, QueryStatefulOrderCountRequest, QueryStatefulOrderCountResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries a ClobPair by id. */
  clobPair(request: QueryGetClobPairRequest): Promise<QueryClobPairResponse>;
  /** Queries a list of ClobPair items. */

  clobPairAll(request?: QueryAllClobPairRequest): Promise<QueryClobPairAllResponse>;
  /** Returns whether a subaccount is liquidatable. */

  areSubaccountsLiquidatable(request: AreSubaccountsLiquidatableRequest): Promise<AreSubaccountsLiquidatableResponse>;
  /** Runs the MEV node <> node calculation with the provided parameters. */

  mevNodeToNodeCalculation(request: MevNodeToNodeCalculationRequest): Promise<MevNodeToNodeCalculationResponse>;
  /** Queries EquityTierLimitConfiguration. */

  equityTierLimitConfiguration(request?: QueryEquityTierLimitConfigurationRequest): Promise<QueryEquityTierLimitConfigurationResponse>;
  /** Queries all stateful orders. */

  allStatefulOrders(request?: QueryAllStatefulOrdersRequest): Promise<QueryAllStatefulOrdersResponse>;
  /** Queries the count of all stateful orders. */

  statefulOrderCount(request: QueryStatefulOrderCountRequest): Promise<QueryStatefulOrderCountResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.clobPair = this.clobPair.bind(this);
    this.clobPairAll = this.clobPairAll.bind(this);
    this.areSubaccountsLiquidatable = this.areSubaccountsLiquidatable.bind(this);
    this.mevNodeToNodeCalculation = this.mevNodeToNodeCalculation.bind(this);
    this.equityTierLimitConfiguration = this.equityTierLimitConfiguration.bind(this);
    this.allStatefulOrders = this.allStatefulOrders.bind(this);
    this.statefulOrderCount = this.statefulOrderCount.bind(this);
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

  mevNodeToNodeCalculation(request: MevNodeToNodeCalculationRequest): Promise<MevNodeToNodeCalculationResponse> {
    const data = MevNodeToNodeCalculationRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "MevNodeToNodeCalculation", data);
    return promise.then(data => MevNodeToNodeCalculationResponse.decode(new _m0.Reader(data)));
  }

  equityTierLimitConfiguration(request: QueryEquityTierLimitConfigurationRequest = {}): Promise<QueryEquityTierLimitConfigurationResponse> {
    const data = QueryEquityTierLimitConfigurationRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "EquityTierLimitConfiguration", data);
    return promise.then(data => QueryEquityTierLimitConfigurationResponse.decode(new _m0.Reader(data)));
  }

  allStatefulOrders(request: QueryAllStatefulOrdersRequest = {
    pagination: undefined
  }): Promise<QueryAllStatefulOrdersResponse> {
    const data = QueryAllStatefulOrdersRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "AllStatefulOrders", data);
    return promise.then(data => QueryAllStatefulOrdersResponse.decode(new _m0.Reader(data)));
  }

  statefulOrderCount(request: QueryStatefulOrderCountRequest): Promise<QueryStatefulOrderCountResponse> {
    const data = QueryStatefulOrderCountRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "StatefulOrderCount", data);
    return promise.then(data => QueryStatefulOrderCountResponse.decode(new _m0.Reader(data)));
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
    },

    mevNodeToNodeCalculation(request: MevNodeToNodeCalculationRequest): Promise<MevNodeToNodeCalculationResponse> {
      return queryService.mevNodeToNodeCalculation(request);
    },

    equityTierLimitConfiguration(request?: QueryEquityTierLimitConfigurationRequest): Promise<QueryEquityTierLimitConfigurationResponse> {
      return queryService.equityTierLimitConfiguration(request);
    },

    allStatefulOrders(request?: QueryAllStatefulOrdersRequest): Promise<QueryAllStatefulOrdersResponse> {
      return queryService.allStatefulOrders(request);
    },

    statefulOrderCount(request: QueryStatefulOrderCountRequest): Promise<QueryStatefulOrderCountResponse> {
      return queryService.statefulOrderCount(request);
    }

  };
};