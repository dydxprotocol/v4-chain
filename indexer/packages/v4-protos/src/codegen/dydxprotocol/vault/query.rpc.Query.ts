import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryParamsRequest, QueryParamsResponse, QueryVaultRequest, QueryVaultResponse, QueryAllVaultsRequest, QueryAllVaultsResponse, QueryOwnerSharesRequest, QueryOwnerSharesResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the Params. */
  params(request?: QueryParamsRequest): Promise<QueryParamsResponse>;
  /** Queries a Vault by type and number. */

  vault(request: QueryVaultRequest): Promise<QueryVaultResponse>;
  /** Queries all vaults. */

  allVaults(request?: QueryAllVaultsRequest): Promise<QueryAllVaultsResponse>;
  /** Queries owner shares of a vault. */

  ownerShares(request: QueryOwnerSharesRequest): Promise<QueryOwnerSharesResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.params = this.params.bind(this);
    this.vault = this.vault.bind(this);
    this.allVaults = this.allVaults.bind(this);
    this.ownerShares = this.ownerShares.bind(this);
  }

  params(request: QueryParamsRequest = {}): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Query", "Params", data);
    return promise.then(data => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  vault(request: QueryVaultRequest): Promise<QueryVaultResponse> {
    const data = QueryVaultRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Query", "Vault", data);
    return promise.then(data => QueryVaultResponse.decode(new _m0.Reader(data)));
  }

  allVaults(request: QueryAllVaultsRequest = {
    pagination: undefined
  }): Promise<QueryAllVaultsResponse> {
    const data = QueryAllVaultsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Query", "AllVaults", data);
    return promise.then(data => QueryAllVaultsResponse.decode(new _m0.Reader(data)));
  }

  ownerShares(request: QueryOwnerSharesRequest): Promise<QueryOwnerSharesResponse> {
    const data = QueryOwnerSharesRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.vault.Query", "OwnerShares", data);
    return promise.then(data => QueryOwnerSharesResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    params(request?: QueryParamsRequest): Promise<QueryParamsResponse> {
      return queryService.params(request);
    },

    vault(request: QueryVaultRequest): Promise<QueryVaultResponse> {
      return queryService.vault(request);
    },

    allVaults(request?: QueryAllVaultsRequest): Promise<QueryAllVaultsResponse> {
      return queryService.allVaults(request);
    },

    ownerShares(request: QueryOwnerSharesRequest): Promise<QueryOwnerSharesResponse> {
      return queryService.ownerShares(request);
    }

  };
};