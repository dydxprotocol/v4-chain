import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryMarketsHardCap, QueryMarketsHardCapResponse, QueryListingVaultDepositParams, QueryListingVaultDepositParamsResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries for the hard cap number of listed markets */
  marketsHardCap(request?: QueryMarketsHardCap): Promise<QueryMarketsHardCapResponse>;
  /** Queries the listing vault deposit params */

  listingVaultDepositParams(request?: QueryListingVaultDepositParams): Promise<QueryListingVaultDepositParamsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.marketsHardCap = this.marketsHardCap.bind(this);
    this.listingVaultDepositParams = this.listingVaultDepositParams.bind(this);
  }

  marketsHardCap(request: QueryMarketsHardCap = {}): Promise<QueryMarketsHardCapResponse> {
    const data = QueryMarketsHardCap.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.listing.Query", "MarketsHardCap", data);
    return promise.then(data => QueryMarketsHardCapResponse.decode(new _m0.Reader(data)));
  }

  listingVaultDepositParams(request: QueryListingVaultDepositParams = {}): Promise<QueryListingVaultDepositParamsResponse> {
    const data = QueryListingVaultDepositParams.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.listing.Query", "ListingVaultDepositParams", data);
    return promise.then(data => QueryListingVaultDepositParamsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    marketsHardCap(request?: QueryMarketsHardCap): Promise<QueryMarketsHardCapResponse> {
      return queryService.marketsHardCap(request);
    },

    listingVaultDepositParams(request?: QueryListingVaultDepositParams): Promise<QueryListingVaultDepositParamsResponse> {
      return queryService.listingVaultDepositParams(request);
    }

  };
};