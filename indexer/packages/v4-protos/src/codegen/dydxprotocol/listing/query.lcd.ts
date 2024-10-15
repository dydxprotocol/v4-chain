import { LCDClient } from "@osmonauts/lcd";
import { QueryMarketsHardCap, QueryMarketsHardCapResponseSDKType, QueryListingVaultDepositParams, QueryListingVaultDepositParamsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.marketsHardCap = this.marketsHardCap.bind(this);
    this.listingVaultDepositParams = this.listingVaultDepositParams.bind(this);
  }
  /* Queries for the hard cap number of listed markets */


  async marketsHardCap(_params: QueryMarketsHardCap = {}): Promise<QueryMarketsHardCapResponseSDKType> {
    const endpoint = `dydxprotocol/listing/markets_hard_cap`;
    return await this.req.get<QueryMarketsHardCapResponseSDKType>(endpoint);
  }
  /* Queries the listing vault deposit params */


  async listingVaultDepositParams(_params: QueryListingVaultDepositParams = {}): Promise<QueryListingVaultDepositParamsResponseSDKType> {
    const endpoint = `dydxprotocol/listing/vault_deposit_params`;
    return await this.req.get<QueryListingVaultDepositParamsResponseSDKType>(endpoint);
  }

}