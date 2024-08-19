import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryParamsRequest, QueryParamsResponseSDKType, QueryVaultRequest, QueryVaultResponseSDKType, QueryAllVaultsRequest, QueryAllVaultsResponseSDKType, QueryTotalSharesRequest, QueryTotalSharesResponseSDKType, QueryOwnerSharesRequest, QueryOwnerSharesResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.params = this.params.bind(this);
    this.vault = this.vault.bind(this);
    this.allVaults = this.allVaults.bind(this);
    this.totalShares = this.totalShares.bind(this);
    this.ownerShares = this.ownerShares.bind(this);
  }
  /* Queries the Params. */


  async params(_params: QueryParamsRequest = {}): Promise<QueryParamsResponseSDKType> {
    const endpoint = `dydxprotocol/vault/params`;
    return await this.req.get<QueryParamsResponseSDKType>(endpoint);
  }
  /* Queries a Vault by type and number. */


  async vault(params: QueryVaultRequest): Promise<QueryVaultResponseSDKType> {
    const endpoint = `dydxprotocol/vault/vault/${params.type}/${params.number}`;
    return await this.req.get<QueryVaultResponseSDKType>(endpoint);
  }
  /* Queries all vaults. */


  async allVaults(params: QueryAllVaultsRequest = {
    pagination: undefined
  }): Promise<QueryAllVaultsResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/vault/vault`;
    return await this.req.get<QueryAllVaultsResponseSDKType>(endpoint, options);
  }
  /* Queries total shares. */


  async totalShares(_params: QueryTotalSharesRequest = {}): Promise<QueryTotalSharesResponseSDKType> {
    const endpoint = `dydxprotocol/vault/total_shares`;
    return await this.req.get<QueryTotalSharesResponseSDKType>(endpoint);
  }
  /* Queries owner shares. */


  async ownerShares(params: QueryOwnerSharesRequest = {
    pagination: undefined
  }): Promise<QueryOwnerSharesResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/vault/owner_shares`;
    return await this.req.get<QueryOwnerSharesResponseSDKType>(endpoint, options);
  }

}