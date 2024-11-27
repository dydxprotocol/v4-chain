import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryPerpetualRequest, QueryPerpetualResponse, QueryAllPerpetualsRequest, QueryAllPerpetualsResponse, QueryAllLiquidityTiersRequest, QueryAllLiquidityTiersResponse, QueryPremiumVotesRequest, QueryPremiumVotesResponse, QueryPremiumSamplesRequest, QueryPremiumSamplesResponse, QueryParamsRequest, QueryParamsResponse, QueryNextPerpetualIdRequest, QueryNextPerpetualIdResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries a Perpetual by id. */
  perpetual(request: QueryPerpetualRequest): Promise<QueryPerpetualResponse>;
  /** Queries a list of Perpetual items. */

  allPerpetuals(request?: QueryAllPerpetualsRequest): Promise<QueryAllPerpetualsResponse>;
  /** Queries a list of LiquidityTiers. */

  allLiquidityTiers(request?: QueryAllLiquidityTiersRequest): Promise<QueryAllLiquidityTiersResponse>;
  /** Queries a list of premium votes. */

  premiumVotes(request?: QueryPremiumVotesRequest): Promise<QueryPremiumVotesResponse>;
  /** Queries a list of premium samples. */

  premiumSamples(request?: QueryPremiumSamplesRequest): Promise<QueryPremiumSamplesResponse>;
  /** Queries the perpetual params. */

  params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries the next perpetual id. */

  nextPerpetualId(request?: QueryNextPerpetualIdRequest): Promise<QueryNextPerpetualIdResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.perpetual = this.perpetual.bind(this);
    this.allPerpetuals = this.allPerpetuals.bind(this);
    this.allLiquidityTiers = this.allLiquidityTiers.bind(this);
    this.premiumVotes = this.premiumVotes.bind(this);
    this.premiumSamples = this.premiumSamples.bind(this);
    this.params = this.params.bind(this);
    this.nextPerpetualId = this.nextPerpetualId.bind(this);
  }

  perpetual(request: QueryPerpetualRequest): Promise<QueryPerpetualResponse> {
    const data = QueryPerpetualRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "Perpetual", data);
    return promise.then(data => QueryPerpetualResponse.decode(new _m0.Reader(data)));
  }

  allPerpetuals(request: QueryAllPerpetualsRequest = {
    pagination: undefined
  }): Promise<QueryAllPerpetualsResponse> {
    const data = QueryAllPerpetualsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "AllPerpetuals", data);
    return promise.then(data => QueryAllPerpetualsResponse.decode(new _m0.Reader(data)));
  }

  allLiquidityTiers(request: QueryAllLiquidityTiersRequest = {
    pagination: undefined
  }): Promise<QueryAllLiquidityTiersResponse> {
    const data = QueryAllLiquidityTiersRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "AllLiquidityTiers", data);
    return promise.then(data => QueryAllLiquidityTiersResponse.decode(new _m0.Reader(data)));
  }

  premiumVotes(request: QueryPremiumVotesRequest = {}): Promise<QueryPremiumVotesResponse> {
    const data = QueryPremiumVotesRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "PremiumVotes", data);
    return promise.then(data => QueryPremiumVotesResponse.decode(new _m0.Reader(data)));
  }

  premiumSamples(request: QueryPremiumSamplesRequest = {}): Promise<QueryPremiumSamplesResponse> {
    const data = QueryPremiumSamplesRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "PremiumSamples", data);
    return promise.then(data => QueryPremiumSamplesResponse.decode(new _m0.Reader(data)));
  }

  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  nextPerpetualId(request: QueryNextPerpetualIdRequest = {}): Promise<QueryNextPerpetualIdResponse> {
    const data = QueryNextPerpetualIdRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "NextPerpetualId", data);
    return promise.then(data => QueryNextPerpetualIdResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    perpetual(request: QueryPerpetualRequest): Promise<QueryPerpetualResponse> {
      return queryService.perpetual(request);
    },

    allPerpetuals(request?: QueryAllPerpetualsRequest): Promise<QueryAllPerpetualsResponse> {
      return queryService.allPerpetuals(request);
    },

    allLiquidityTiers(request?: QueryAllLiquidityTiersRequest): Promise<QueryAllLiquidityTiersResponse> {
      return queryService.allLiquidityTiers(request);
    },

    premiumVotes(request?: QueryPremiumVotesRequest): Promise<QueryPremiumVotesResponse> {
      return queryService.premiumVotes(request);
    },

    premiumSamples(request?: QueryPremiumSamplesRequest): Promise<QueryPremiumSamplesResponse> {
      return queryService.premiumSamples(request);
    },

    params(request?: QueryParamsRequest): Promise<QueryParamsResponse> {
      return queryService.params(request);
    },

    nextPerpetualId(request?: QueryNextPerpetualIdRequest): Promise<QueryNextPerpetualIdResponse> {
      return queryService.nextPerpetualId(request);
    }

  };
};