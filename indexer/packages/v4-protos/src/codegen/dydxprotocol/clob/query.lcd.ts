import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryGetClobPairRequest, QueryClobPairResponseSDKType, QueryAllClobPairRequest, QueryClobPairAllResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.clobPair = this.clobPair.bind(this);
    this.clobPairAll = this.clobPairAll.bind(this);
  }
  /* Queries a ClobPair by id. */


  async clobPair(params: QueryGetClobPairRequest): Promise<QueryClobPairResponseSDKType> {
    const endpoint = `dydxprotocol/clob/clob_pair/${params.id}`;
    return await this.req.get<QueryClobPairResponseSDKType>(endpoint);
  }
  /* Queries a list of ClobPair items. */


  async clobPairAll(params: QueryAllClobPairRequest = {
    pagination: undefined
  }): Promise<QueryClobPairAllResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/clob/clob_pair`;
    return await this.req.get<QueryClobPairAllResponseSDKType>(endpoint, options);
  }

}