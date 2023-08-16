import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryPerpetualRequest, QueryPerpetualResponseSDKType, QueryAllPerpetualsRequest, QueryAllPerpetualsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.perpetual = this.perpetual.bind(this);
    this.allPerpetuals = this.allPerpetuals.bind(this);
  }
  /* Queries a Perpetual by id. */


  async perpetual(params: QueryPerpetualRequest): Promise<QueryPerpetualResponseSDKType> {
    const endpoint = `dydxprotocol/perpetuals/perpetual/${params.id}`;
    return await this.req.get<QueryPerpetualResponseSDKType>(endpoint);
  }
  /* Queries a list of Perpetual items. */


  async allPerpetuals(params: QueryAllPerpetualsRequest = {
    pagination: undefined
  }): Promise<QueryAllPerpetualsResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/perpetuals/perpetual`;
    return await this.req.get<QueryAllPerpetualsResponseSDKType>(endpoint, options);
  }

}