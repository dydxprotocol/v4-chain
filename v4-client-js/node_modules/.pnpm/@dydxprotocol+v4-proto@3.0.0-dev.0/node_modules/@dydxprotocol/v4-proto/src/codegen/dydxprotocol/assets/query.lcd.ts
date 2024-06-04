import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryAssetRequest, QueryAssetResponseSDKType, QueryAllAssetsRequest, QueryAllAssetsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.asset = this.asset.bind(this);
    this.allAssets = this.allAssets.bind(this);
  }
  /* Queries a Asset by id. */


  async asset(params: QueryAssetRequest): Promise<QueryAssetResponseSDKType> {
    const endpoint = `dydxprotocol/assets/asset/${params.id}`;
    return await this.req.get<QueryAssetResponseSDKType>(endpoint);
  }
  /* Queries a list of Asset items. */


  async allAssets(params: QueryAllAssetsRequest = {
    pagination: undefined
  }): Promise<QueryAllAssetsResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/assets/asset`;
    return await this.req.get<QueryAllAssetsResponseSDKType>(endpoint, options);
  }

}