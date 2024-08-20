import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryParamsRequest, QueryParamsResponseSDKType, QueryVaultRequest, QueryVaultResponseSDKType, QueryAllVaultsRequest, QueryAllVaultsResponseSDKType, QueryMegavaultTotalSharesRequest, QueryMegavaultTotalSharesResponseSDKType, QueryMegavaultOwnerSharesRequest, QueryMegavaultOwnerSharesResponseSDKType } from "./query";
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
    this.megavaultTotalShares = this.megavaultTotalShares.bind(this);
    this.megavaultOwnerShares = this.megavaultOwnerShares.bind(this);
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
  /* Queries total shares of megavault. */


  async megavaultTotalShares(_params: QueryMegavaultTotalSharesRequest = {}): Promise<QueryMegavaultTotalSharesResponseSDKType> {
    const endpoint = `dydxprotocol/vault/megavault/total_shares`;
    return await this.req.get<QueryMegavaultTotalSharesResponseSDKType>(endpoint);
  }
  /* Queries owner shares of megavault. */


  async megavaultOwnerShares(params: QueryMegavaultOwnerSharesRequest = {
    pagination: undefined
  }): Promise<QueryMegavaultOwnerSharesResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/vault/megavault/owner_shares`;
    return await this.req.get<QueryMegavaultOwnerSharesResponseSDKType>(endpoint, options);
  }

}