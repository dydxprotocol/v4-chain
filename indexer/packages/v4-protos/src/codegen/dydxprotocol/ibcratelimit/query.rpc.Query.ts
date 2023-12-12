import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { ListLimitParamsRequest, ListLimitParamsResponse, QueryCapacityByDenomRequest, QueryCapacityByDenomResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** List all limit params. */
  listLimitParams(request?: ListLimitParamsRequest): Promise<ListLimitParamsResponse>;
  /** Query capacity by denom. */

  capacityByDenom(request: QueryCapacityByDenomRequest): Promise<QueryCapacityByDenomResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.listLimitParams = this.listLimitParams.bind(this);
    this.capacityByDenom = this.capacityByDenom.bind(this);
  }

  listLimitParams(request: ListLimitParamsRequest = {}): Promise<ListLimitParamsResponse> {
    const data = ListLimitParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.ibcratelimit.Query", "ListLimitParams", data);
    return promise.then(data => ListLimitParamsResponse.decode(new _m0.Reader(data)));
  }

  capacityByDenom(request: QueryCapacityByDenomRequest): Promise<QueryCapacityByDenomResponse> {
    const data = QueryCapacityByDenomRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.ibcratelimit.Query", "CapacityByDenom", data);
    return promise.then(data => QueryCapacityByDenomResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    listLimitParams(request?: ListLimitParamsRequest): Promise<ListLimitParamsResponse> {
      return queryService.listLimitParams(request);
    },

    capacityByDenom(request: QueryCapacityByDenomRequest): Promise<QueryCapacityByDenomResponse> {
      return queryService.capacityByDenom(request);
    }

  };
};