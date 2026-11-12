import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryStateRequest, QueryStateResponse, QueryPendingSettlementsRequest, QueryPendingSettlementsResponse, QueryPendingSettlementRequest, QueryPendingSettlementResponse } from "./query";
/** Query defines the canonical USDC query service. */

export interface Query {
  /** State returns controls and aggregate backing state. */
  state(request?: QueryStateRequest): Promise<QueryStateResponse>;
  /** PendingSettlements returns a page of in-flight canonical settlements. */

  pendingSettlements(request?: QueryPendingSettlementsRequest): Promise<QueryPendingSettlementsResponse>;
  /**
   * PendingSettlement returns the in-flight canonical settlement for an IBC
   * packet identified by source channel and sequence.
   */

  pendingSettlement(request: QueryPendingSettlementRequest): Promise<QueryPendingSettlementResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.state = this.state.bind(this);
    this.pendingSettlements = this.pendingSettlements.bind(this);
    this.pendingSettlement = this.pendingSettlement.bind(this);
  }

  state(request: QueryStateRequest = {}): Promise<QueryStateResponse> {
    const data = QueryStateRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.canonicalusdc.Query", "State", data);
    return promise.then(data => QueryStateResponse.decode(new _m0.Reader(data)));
  }

  pendingSettlements(request: QueryPendingSettlementsRequest = {
    pagination: undefined
  }): Promise<QueryPendingSettlementsResponse> {
    const data = QueryPendingSettlementsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.canonicalusdc.Query", "PendingSettlements", data);
    return promise.then(data => QueryPendingSettlementsResponse.decode(new _m0.Reader(data)));
  }

  pendingSettlement(request: QueryPendingSettlementRequest): Promise<QueryPendingSettlementResponse> {
    const data = QueryPendingSettlementRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.canonicalusdc.Query", "PendingSettlement", data);
    return promise.then(data => QueryPendingSettlementResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    state(request?: QueryStateRequest): Promise<QueryStateResponse> {
      return queryService.state(request);
    },

    pendingSettlements(request?: QueryPendingSettlementsRequest): Promise<QueryPendingSettlementsResponse> {
      return queryService.pendingSettlements(request);
    },

    pendingSettlement(request: QueryPendingSettlementRequest): Promise<QueryPendingSettlementResponse> {
      return queryService.pendingSettlement(request);
    }

  };
};