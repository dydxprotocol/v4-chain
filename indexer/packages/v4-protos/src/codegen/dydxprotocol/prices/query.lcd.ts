import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryMarketPriceRequest, QueryMarketPriceResponseSDKType, QueryAllMarketPricesRequest, QueryAllMarketPricesResponseSDKType, QueryMarketParamRequest, QueryMarketParamResponseSDKType, QueryAllMarketParamsRequest, QueryAllMarketParamsResponseSDKType, QueryNextMarketIdRequest, QueryNextMarketIdResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.marketPrice = this.marketPrice.bind(this);
    this.allMarketPrices = this.allMarketPrices.bind(this);
    this.marketParam = this.marketParam.bind(this);
    this.allMarketParams = this.allMarketParams.bind(this);
    this.nextMarketId = this.nextMarketId.bind(this);
  }
  /* Queries a MarketPrice by id. */


  async marketPrice(params: QueryMarketPriceRequest): Promise<QueryMarketPriceResponseSDKType> {
    const endpoint = `dydxprotocol/prices/market/${params.id}`;
    return await this.req.get<QueryMarketPriceResponseSDKType>(endpoint);
  }
  /* Queries a list of MarketPrice items. */


  async allMarketPrices(params: QueryAllMarketPricesRequest = {
    pagination: undefined
  }): Promise<QueryAllMarketPricesResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/prices/market`;
    return await this.req.get<QueryAllMarketPricesResponseSDKType>(endpoint, options);
  }
  /* Queries a MarketParam by id. */


  async marketParam(params: QueryMarketParamRequest): Promise<QueryMarketParamResponseSDKType> {
    const endpoint = `dydxprotocol/prices/params/market/${params.id}`;
    return await this.req.get<QueryMarketParamResponseSDKType>(endpoint);
  }
  /* Queries a list of MarketParam items. */


  async allMarketParams(params: QueryAllMarketParamsRequest = {
    pagination: undefined
  }): Promise<QueryAllMarketParamsResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/prices/params/market`;
    return await this.req.get<QueryAllMarketParamsResponseSDKType>(endpoint, options);
  }
  /* Queries the next market id. */


  async nextMarketId(_params: QueryNextMarketIdRequest = {}): Promise<QueryNextMarketIdResponseSDKType> {
    const endpoint = `dydxprotocol/prices/next_market_id`;
    return await this.req.get<QueryNextMarketIdResponseSDKType>(endpoint);
  }

}