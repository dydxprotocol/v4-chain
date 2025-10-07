import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryPerpetualFeeParamsRequest, QueryPerpetualFeeParamsResponse, QueryUserFeeTierRequest, QueryUserFeeTierResponse, QueryFeeDiscountCampaignParamsRequest, QueryFeeDiscountCampaignParamsResponse, QueryAllFeeDiscountCampaignParamsRequest, QueryAllFeeDiscountCampaignParamsResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the PerpetualFeeParams. */
  perpetualFeeParams(request?: QueryPerpetualFeeParamsRequest): Promise<QueryPerpetualFeeParamsResponse>;
  /** Queries a user's fee tier */

  userFeeTier(request: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponse>;
  /**
   * FeeDiscountCampaignParams queries fee discount campaign parameters for a
   * specific CLOB pair.
   */

  feeDiscountCampaignParams(request: QueryFeeDiscountCampaignParamsRequest): Promise<QueryFeeDiscountCampaignParamsResponse>;
  /** AllFeeDiscountCampaignParams queries all fee discount campaign parameters. */

  allFeeDiscountCampaignParams(request?: QueryAllFeeDiscountCampaignParamsRequest): Promise<QueryAllFeeDiscountCampaignParamsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.perpetualFeeParams = this.perpetualFeeParams.bind(this);
    this.userFeeTier = this.userFeeTier.bind(this);
    this.feeDiscountCampaignParams = this.feeDiscountCampaignParams.bind(this);
    this.allFeeDiscountCampaignParams = this.allFeeDiscountCampaignParams.bind(this);
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

  feeDiscountCampaignParams(request: QueryFeeDiscountCampaignParamsRequest): Promise<QueryFeeDiscountCampaignParamsResponse> {
    const data = QueryFeeDiscountCampaignParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "FeeDiscountCampaignParams", data);
    return promise.then(data => QueryFeeDiscountCampaignParamsResponse.decode(new _m0.Reader(data)));
  }

  allFeeDiscountCampaignParams(request: QueryAllFeeDiscountCampaignParamsRequest = {}): Promise<QueryAllFeeDiscountCampaignParamsResponse> {
    const data = QueryAllFeeDiscountCampaignParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "AllFeeDiscountCampaignParams", data);
    return promise.then(data => QueryAllFeeDiscountCampaignParamsResponse.decode(new _m0.Reader(data)));
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

    feeDiscountCampaignParams(request: QueryFeeDiscountCampaignParamsRequest): Promise<QueryFeeDiscountCampaignParamsResponse> {
      return queryService.feeDiscountCampaignParams(request);
    },

    allFeeDiscountCampaignParams(request?: QueryAllFeeDiscountCampaignParamsRequest): Promise<QueryAllFeeDiscountCampaignParamsResponse> {
      return queryService.allFeeDiscountCampaignParams(request);
    }

  };
};