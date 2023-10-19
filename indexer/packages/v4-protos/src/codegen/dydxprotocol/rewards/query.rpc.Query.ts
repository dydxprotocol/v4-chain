import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryRewardShareRequest, QueryRewardShareResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the Params. */
  params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a reward share by address. */

  rewardShare(request: QueryRewardShareRequest): Promise<QueryRewardShareResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.params = this.params.bind(this);
    this.rewardShare = this.rewardShare.bind(this);
  }

  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.rewards.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  rewardShare(request: QueryRewardShareRequest): Promise<QueryRewardShareResponse> {
    const data = QueryRewardShareRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.rewards.Query", "RewardShare", data);
    return promise.then(data => QueryRewardShareResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse> {
      return queryService.params(request);
    },

    rewardShare(request: QueryRewardShareRequest): Promise<QueryRewardShareResponse> {
      return queryService.rewardShare(request);
    }

  };
};