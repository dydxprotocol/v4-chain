import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryPerpetualRequest, QueryPerpetualResponseSDKType, QueryAllPerpetualsRequest, QueryAllPerpetualsResponseSDKType, QueryPremiumVotesRequest, QueryPremiumVotesResponseSDKType, QueryPremiumSamplesRequest, QueryPremiumSamplesResponseSDKType, QueryParamsRequest, QueryParamsResponseSDKType } from "./query";
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
    this.premiumVotes = this.premiumVotes.bind(this);
    this.premiumSamples = this.premiumSamples.bind(this);
    this.params = this.params.bind(this);
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
  /* Queries a list of premium votes. */


  async premiumVotes(_params: QueryPremiumVotesRequest = {}): Promise<QueryPremiumVotesResponseSDKType> {
    const endpoint = `dydxprotocol/perpetuals/premium_votes`;
    return await this.req.get<QueryPremiumVotesResponseSDKType>(endpoint);
  }
  /* Queries a list of premium samples. */


  async premiumSamples(_params: QueryPremiumSamplesRequest = {}): Promise<QueryPremiumSamplesResponseSDKType> {
    const endpoint = `dydxprotocol/perpetuals/premium_samples`;
    return await this.req.get<QueryPremiumSamplesResponseSDKType>(endpoint);
  }
  /* Queries the perpetual params. */


  async params(_params: QueryParamsRequest = {}): Promise<QueryParamsResponseSDKType> {
    const endpoint = `dydxprotocol/perpetuals/params`;
    return await this.req.get<QueryParamsResponseSDKType>(endpoint);
  }

}