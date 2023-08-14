import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryGetEpochInfoRequest, QueryEpochInfoResponseSDKType, QueryAllEpochInfoRequest, QueryEpochInfoAllResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.epochInfo = this.epochInfo.bind(this);
    this.epochInfoAll = this.epochInfoAll.bind(this);
  }
  /* Queries a EpochInfo by name. */


  async epochInfo(params: QueryGetEpochInfoRequest): Promise<QueryEpochInfoResponseSDKType> {
    const endpoint = `dydxprotocol/v4/epochs/epoch_info/${params.name}`;
    return await this.req.get<QueryEpochInfoResponseSDKType>(endpoint);
  }
  /* Queries a list of EpochInfo items. */


  async epochInfoAll(params: QueryAllEpochInfoRequest = {
    pagination: undefined
  }): Promise<QueryEpochInfoAllResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/v4/epochs/epoch_info`;
    return await this.req.get<QueryEpochInfoAllResponseSDKType>(endpoint, options);
  }

}