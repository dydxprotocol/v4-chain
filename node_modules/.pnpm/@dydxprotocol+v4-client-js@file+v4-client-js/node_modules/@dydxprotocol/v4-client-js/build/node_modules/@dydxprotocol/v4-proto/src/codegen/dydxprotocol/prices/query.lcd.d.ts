import { LCDClient } from "@osmonauts/lcd";
import { QueryMarketPriceRequest, QueryMarketPriceResponseSDKType, QueryAllMarketPricesRequest, QueryAllMarketPricesResponseSDKType, QueryMarketParamRequest, QueryMarketParamResponseSDKType, QueryAllMarketParamsRequest, QueryAllMarketParamsResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    marketPrice(params: QueryMarketPriceRequest): Promise<QueryMarketPriceResponseSDKType>;
    allMarketPrices(params?: QueryAllMarketPricesRequest): Promise<QueryAllMarketPricesResponseSDKType>;
    marketParam(params: QueryMarketParamRequest): Promise<QueryMarketParamResponseSDKType>;
    allMarketParams(params?: QueryAllMarketParamsRequest): Promise<QueryAllMarketParamsResponseSDKType>;
}
