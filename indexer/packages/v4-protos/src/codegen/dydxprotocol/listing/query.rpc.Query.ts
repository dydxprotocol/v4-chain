import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryPermissionlessMarketListingStatus, QueryPermissionlessMarketListingStatusResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries if permissionless listings are enabled */
  permissionlessMarketListingStatus(request?: QueryPermissionlessMarketListingStatus): Promise<QueryPermissionlessMarketListingStatusResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.permissionlessMarketListingStatus = this.permissionlessMarketListingStatus.bind(this);
  }

  permissionlessMarketListingStatus(request: QueryPermissionlessMarketListingStatus = {}): Promise<QueryPermissionlessMarketListingStatusResponse> {
    const data = QueryPermissionlessMarketListingStatus.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.listing.Query", "PermissionlessMarketListingStatus", data);
    return promise.then(data => QueryPermissionlessMarketListingStatusResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    permissionlessMarketListingStatus(request?: QueryPermissionlessMarketListingStatus): Promise<QueryPermissionlessMarketListingStatusResponse> {
      return queryService.permissionlessMarketListingStatus(request);
    }

  };
};