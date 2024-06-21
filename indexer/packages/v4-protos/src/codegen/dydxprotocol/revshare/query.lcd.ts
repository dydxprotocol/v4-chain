import { LCDClient } from "@osmonauts/lcd";
import { QueryMarketMapperRevShareDetails, QueryMarketMapperRevShareDetailsResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.marketMapperRevShareDetails = this.marketMapperRevShareDetails.bind(this);
  }
  /* Queries market mapper revenue share details for a specific market */


  async marketMapperRevShareDetails(params: QueryMarketMapperRevShareDetails): Promise<QueryMarketMapperRevShareDetailsResponseSDKType> {
    const endpoint = `dydxprotocol/revshare/market_mapper_rev_share_details/${params.marketId}`;
    return await this.req.get<QueryMarketMapperRevShareDetailsResponseSDKType>(endpoint);
  }

}