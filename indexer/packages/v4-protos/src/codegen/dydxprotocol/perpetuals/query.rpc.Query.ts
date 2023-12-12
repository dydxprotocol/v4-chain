import { Rpc } from "../../helpers";
import { BinaryReader } from "../../binary";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryPerpetualRequest, QueryPerpetualResponse, QueryAllPerpetualsRequest, QueryAllPerpetualsResponse, QueryAllLiquidityTiersRequest, QueryAllLiquidityTiersResponse, QueryPremiumVotesRequest, QueryPremiumVotesResponse, QueryPremiumSamplesRequest, QueryPremiumSamplesResponse, QueryParamsRequest, QueryParamsResponse } from "./query";
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
  }
  perpetual(request: QueryPerpetualRequest): Promise<QueryPerpetualResponse> {
    const data = QueryPerpetualRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "Perpetual", data);
    return promise.then(data => QueryPerpetualResponse.decode(new BinaryReader(data)));
  }
  allPerpetuals(request: QueryAllPerpetualsRequest = {
    pagination: undefined
  }): Promise<QueryAllPerpetualsResponse> {
    const data = QueryAllPerpetualsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "AllPerpetuals", data);
    return promise.then(data => QueryAllPerpetualsResponse.decode(new BinaryReader(data)));
  }
  allLiquidityTiers(request: QueryAllLiquidityTiersRequest = {
    pagination: undefined
  }): Promise<QueryAllLiquidityTiersResponse> {
    const data = QueryAllLiquidityTiersRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "AllLiquidityTiers", data);
    return promise.then(data => QueryAllLiquidityTiersResponse.decode(new BinaryReader(data)));
  }
  premiumVotes(request: QueryPremiumVotesRequest = {}): Promise<QueryPremiumVotesResponse> {
    const data = QueryPremiumVotesRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "PremiumVotes", data);
    return promise.then(data => QueryPremiumVotesResponse.decode(new BinaryReader(data)));
  }
  premiumSamples(request: QueryPremiumSamplesRequest = {}): Promise<QueryPremiumSamplesResponse> {
    const data = QueryPremiumSamplesRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "PremiumSamples", data);
    return promise.then(data => QueryPremiumSamplesResponse.decode(new BinaryReader(data)));
  }
  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.perpetuals.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new BinaryReader(data)));
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
    }
  };
};