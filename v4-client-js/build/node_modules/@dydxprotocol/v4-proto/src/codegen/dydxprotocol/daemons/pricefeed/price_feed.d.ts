/// <reference types="long" />
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../../helpers";
/** UpdateMarketPriceRequest is a request message updating market prices. */
export interface UpdateMarketPricesRequest {
    marketPriceUpdates: MarketPriceUpdate[];
}
/** UpdateMarketPriceRequest is a request message updating market prices. */
export interface UpdateMarketPricesRequestSDKType {
    market_price_updates: MarketPriceUpdateSDKType[];
}
/** UpdateMarketPricesResponse is a response message for updating market prices. */
export interface UpdateMarketPricesResponse {
}
/** UpdateMarketPricesResponse is a response message for updating market prices. */
export interface UpdateMarketPricesResponseSDKType {
}
/** ExchangePrice represents a specific exchange's market price */
export interface ExchangePrice {
    exchangeId: string;
    price: Long;
    lastUpdateTime?: Date;
}
/** ExchangePrice represents a specific exchange's market price */
export interface ExchangePriceSDKType {
    exchange_id: string;
    price: Long;
    last_update_time?: Date;
}
/** MarketPriceUpdate represents an update to a single market */
export interface MarketPriceUpdate {
    marketId: number;
    exchangePrices: ExchangePrice[];
}
/** MarketPriceUpdate represents an update to a single market */
export interface MarketPriceUpdateSDKType {
    market_id: number;
    exchange_prices: ExchangePriceSDKType[];
}
export declare const UpdateMarketPricesRequest: {
    encode(message: UpdateMarketPricesRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): UpdateMarketPricesRequest;
    fromPartial(object: DeepPartial<UpdateMarketPricesRequest>): UpdateMarketPricesRequest;
};
export declare const UpdateMarketPricesResponse: {
    encode(_: UpdateMarketPricesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): UpdateMarketPricesResponse;
    fromPartial(_: DeepPartial<UpdateMarketPricesResponse>): UpdateMarketPricesResponse;
};
export declare const ExchangePrice: {
    encode(message: ExchangePrice, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ExchangePrice;
    fromPartial(object: DeepPartial<ExchangePrice>): ExchangePrice;
};
export declare const MarketPriceUpdate: {
    encode(message: MarketPriceUpdate, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MarketPriceUpdate;
    fromPartial(object: DeepPartial<MarketPriceUpdate>): MarketPriceUpdate;
};
