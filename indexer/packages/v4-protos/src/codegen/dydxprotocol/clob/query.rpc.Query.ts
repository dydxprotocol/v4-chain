import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryGetClobPairRequest, QueryClobPairResponse, QueryAllClobPairRequest, QueryClobPairAllResponse, MevNodeToNodeCalculationRequest, MevNodeToNodeCalculationResponse, QueryEquityTierLimitConfigurationRequest, QueryEquityTierLimitConfigurationResponse, QueryBlockRateLimitConfigurationRequest, QueryBlockRateLimitConfigurationResponse, QueryLiquidationsConfigurationRequest, QueryLiquidationsConfigurationResponse, QueryStatefulOrderRequest, QueryStatefulOrderResponse, QueryNextClobPairIdRequest, QueryNextClobPairIdResponse, QueryLeverageRequest, QueryLeverageResponse, StreamOrderbookUpdatesRequest, StreamOrderbookUpdatesResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries a ClobPair by id. */
  clobPair(request: QueryGetClobPairRequest): Promise<QueryClobPairResponse>;
  /** Queries a list of ClobPair items. */

  clobPairAll(request?: QueryAllClobPairRequest): Promise<QueryClobPairAllResponse>;
  /** Runs the MEV node <> node calculation with the provided parameters. */

  mevNodeToNodeCalculation(request: MevNodeToNodeCalculationRequest): Promise<MevNodeToNodeCalculationResponse>;
  /** Queries EquityTierLimitConfiguration. */

  equityTierLimitConfiguration(request?: QueryEquityTierLimitConfigurationRequest): Promise<QueryEquityTierLimitConfigurationResponse>;
  /** Queries BlockRateLimitConfiguration. */

  blockRateLimitConfiguration(request?: QueryBlockRateLimitConfigurationRequest): Promise<QueryBlockRateLimitConfigurationResponse>;
  /** Queries LiquidationsConfiguration. */

  liquidationsConfiguration(request?: QueryLiquidationsConfigurationRequest): Promise<QueryLiquidationsConfigurationResponse>;
  /** Queries the stateful order for a given order id. */

  statefulOrder(request: QueryStatefulOrderRequest): Promise<QueryStatefulOrderResponse>;
  /** Queries the next clob pair id. */

  nextClobPairId(request?: QueryNextClobPairIdRequest): Promise<QueryNextClobPairIdResponse>;
  /** Queries leverage for a subaccount. */

  leverage(request: QueryLeverageRequest): Promise<QueryLeverageResponse>;
  /**
   * Streams orderbook updates. Updates contain orderbook data
   * such as order placements, updates, and fills.
   */

  streamOrderbookUpdates(request: StreamOrderbookUpdatesRequest): Promise<StreamOrderbookUpdatesResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.clobPair = this.clobPair.bind(this);
    this.clobPairAll = this.clobPairAll.bind(this);
    this.mevNodeToNodeCalculation = this.mevNodeToNodeCalculation.bind(this);
    this.equityTierLimitConfiguration = this.equityTierLimitConfiguration.bind(this);
    this.blockRateLimitConfiguration = this.blockRateLimitConfiguration.bind(this);
    this.liquidationsConfiguration = this.liquidationsConfiguration.bind(this);
    this.statefulOrder = this.statefulOrder.bind(this);
    this.nextClobPairId = this.nextClobPairId.bind(this);
    this.leverage = this.leverage.bind(this);
    this.streamOrderbookUpdates = this.streamOrderbookUpdates.bind(this);
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

  blockRateLimitConfiguration(request: QueryBlockRateLimitConfigurationRequest = {}): Promise<QueryBlockRateLimitConfigurationResponse> {
    const data = QueryBlockRateLimitConfigurationRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "BlockRateLimitConfiguration", data);
    return promise.then(data => QueryBlockRateLimitConfigurationResponse.decode(new _m0.Reader(data)));
  }

  liquidationsConfiguration(request: QueryLiquidationsConfigurationRequest = {}): Promise<QueryLiquidationsConfigurationResponse> {
    const data = QueryLiquidationsConfigurationRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "LiquidationsConfiguration", data);
    return promise.then(data => QueryLiquidationsConfigurationResponse.decode(new _m0.Reader(data)));
  }

  statefulOrder(request: QueryStatefulOrderRequest): Promise<QueryStatefulOrderResponse> {
    const data = QueryStatefulOrderRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "StatefulOrder", data);
    return promise.then(data => QueryStatefulOrderResponse.decode(new _m0.Reader(data)));
  }

  nextClobPairId(request: QueryNextClobPairIdRequest = {}): Promise<QueryNextClobPairIdResponse> {
    const data = QueryNextClobPairIdRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "NextClobPairId", data);
    return promise.then(data => QueryNextClobPairIdResponse.decode(new _m0.Reader(data)));
  }

  leverage(request: QueryLeverageRequest): Promise<QueryLeverageResponse> {
    const data = QueryLeverageRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "Leverage", data);
    return promise.then(data => QueryLeverageResponse.decode(new _m0.Reader(data)));
  }

  streamOrderbookUpdates(request: StreamOrderbookUpdatesRequest): Promise<StreamOrderbookUpdatesResponse> {
    const data = StreamOrderbookUpdatesRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.clob.Query", "StreamOrderbookUpdates", data);
    return promise.then(data => StreamOrderbookUpdatesResponse.decode(new _m0.Reader(data)));
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

    mevNodeToNodeCalculation(request: MevNodeToNodeCalculationRequest): Promise<MevNodeToNodeCalculationResponse> {
      return queryService.mevNodeToNodeCalculation(request);
    },

    equityTierLimitConfiguration(request?: QueryEquityTierLimitConfigurationRequest): Promise<QueryEquityTierLimitConfigurationResponse> {
      return queryService.equityTierLimitConfiguration(request);
    },

    blockRateLimitConfiguration(request?: QueryBlockRateLimitConfigurationRequest): Promise<QueryBlockRateLimitConfigurationResponse> {
      return queryService.blockRateLimitConfiguration(request);
    },

    liquidationsConfiguration(request?: QueryLiquidationsConfigurationRequest): Promise<QueryLiquidationsConfigurationResponse> {
      return queryService.liquidationsConfiguration(request);
    },

    statefulOrder(request: QueryStatefulOrderRequest): Promise<QueryStatefulOrderResponse> {
      return queryService.statefulOrder(request);
    },

    nextClobPairId(request?: QueryNextClobPairIdRequest): Promise<QueryNextClobPairIdResponse> {
      return queryService.nextClobPairId(request);
    },

    leverage(request: QueryLeverageRequest): Promise<QueryLeverageResponse> {
      return queryService.leverage(request);
    },

    streamOrderbookUpdates(request: StreamOrderbookUpdatesRequest): Promise<StreamOrderbookUpdatesResponse> {
      return queryService.streamOrderbookUpdates(request);
    }

  };
};