import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryAssetRequest, QueryAssetResponse, QueryAllAssetsRequest, QueryAllAssetsResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries a Asset by id. */
  asset(request: QueryAssetRequest): Promise<QueryAssetResponse>;
  /** Queries a list of Asset items. */

  allAssets(request?: QueryAllAssetsRequest): Promise<QueryAllAssetsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.asset = this.asset.bind(this);
    this.allAssets = this.allAssets.bind(this);
  }

  asset(request: QueryAssetRequest): Promise<QueryAssetResponse> {
    const data = QueryAssetRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.assets.Query", "Asset", data);
    return promise.then(data => QueryAssetResponse.decode(new _m0.Reader(data)));
  }

  allAssets(request: QueryAllAssetsRequest = {
    pagination: undefined
  }): Promise<QueryAllAssetsResponse> {
    const data = QueryAllAssetsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.assets.Query", "AllAssets", data);
    return promise.then(data => QueryAllAssetsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    asset(request: QueryAssetRequest): Promise<QueryAssetResponse> {
      return queryService.asset(request);
    },

    allAssets(request?: QueryAllAssetsRequest): Promise<QueryAllAssetsResponse> {
      return queryService.allAssets(request);
    }

  };
};