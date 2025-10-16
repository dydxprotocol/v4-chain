import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryPerpetualFeeParamsRequest, QueryPerpetualFeeParamsResponse, QueryUserFeeTierRequest, QueryUserFeeTierResponse, QueryStakingTiersRequest, QueryStakingTiersResponse, QueryUserStakingTierRequest, QueryUserStakingTierResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the PerpetualFeeParams. */
  perpetualFeeParams(request?: QueryPerpetualFeeParamsRequest): Promise<QueryPerpetualFeeParamsResponse>;
  /** Queries a user's fee tier */

  userFeeTier(request: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponse>;
  /** Get all staking tiers */

  stakingTiers(request?: QueryStakingTiersRequest): Promise<QueryStakingTiersResponse>;
  /** Get user's current staked amount and staking tier */

  userStakingTier(request: QueryUserStakingTierRequest): Promise<QueryUserStakingTierResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.perpetualFeeParams = this.perpetualFeeParams.bind(this);
    this.userFeeTier = this.userFeeTier.bind(this);
    this.stakingTiers = this.stakingTiers.bind(this);
    this.userStakingTier = this.userStakingTier.bind(this);
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

  stakingTiers(request: QueryStakingTiersRequest = {}): Promise<QueryStakingTiersResponse> {
    const data = QueryStakingTiersRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "StakingTiers", data);
    return promise.then(data => QueryStakingTiersResponse.decode(new _m0.Reader(data)));
  }

  userStakingTier(request: QueryUserStakingTierRequest): Promise<QueryUserStakingTierResponse> {
    const data = QueryUserStakingTierRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "UserStakingTier", data);
    return promise.then(data => QueryUserStakingTierResponse.decode(new _m0.Reader(data)));
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

    stakingTiers(request?: QueryStakingTiersRequest): Promise<QueryStakingTiersResponse> {
      return queryService.stakingTiers(request);
    },

    userStakingTier(request: QueryUserStakingTierRequest): Promise<QueryUserStakingTierResponse> {
      return queryService.userStakingTier(request);
    }

  };
};