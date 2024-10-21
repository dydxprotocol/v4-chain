import { LCDClient } from "@osmonauts/lcd";
import { QueryMarketMapperRevenueShareParams, QueryMarketMapperRevenueShareParamsResponseSDKType, QueryMarketMapperRevShareDetails, QueryMarketMapperRevShareDetailsResponseSDKType, QueryUnconditionalRevShareConfig, QueryUnconditionalRevShareConfigResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.marketMapperRevenueShareParams = this.marketMapperRevenueShareParams.bind(this);
    this.marketMapperRevShareDetails = this.marketMapperRevShareDetails.bind(this);
    this.unconditionalRevShareConfig = this.unconditionalRevShareConfig.bind(this);
  }
  /* MarketMapperRevenueShareParams queries the revenue share params for the
   market mapper */


  async marketMapperRevenueShareParams(_params: QueryMarketMapperRevenueShareParams = {}): Promise<QueryMarketMapperRevenueShareParamsResponseSDKType> {
    const endpoint = `dydxprotocol/revshare/market_mapper_rev_share_params`;
    return await this.req.get<QueryMarketMapperRevenueShareParamsResponseSDKType>(endpoint);
  }
  /* Queries market mapper revenue share details for a specific market */


  async marketMapperRevShareDetails(params: QueryMarketMapperRevShareDetails): Promise<QueryMarketMapperRevShareDetailsResponseSDKType> {
    const endpoint = `dydxprotocol/revshare/market_mapper_rev_share_details/${params.marketId}`;
    return await this.req.get<QueryMarketMapperRevShareDetailsResponseSDKType>(endpoint);
  }
  /* Queries unconditional revenue share config */


  async unconditionalRevShareConfig(_params: QueryUnconditionalRevShareConfig = {}): Promise<QueryUnconditionalRevShareConfigResponseSDKType> {
    const endpoint = `dydxprotocol/revshare/unconditional_rev_share`;
    return await this.req.get<QueryUnconditionalRevShareConfigResponseSDKType>(endpoint);
  }

}