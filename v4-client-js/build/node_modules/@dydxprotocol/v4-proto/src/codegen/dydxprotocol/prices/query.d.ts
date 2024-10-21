import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { MarketPrice, MarketPriceSDKType } from "./market_price";
import { MarketParam, MarketParamSDKType } from "./market_param";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * QueryMarketPriceRequest is request type for the Query/Params `MarketPrice`
 * RPC method.
 */
export interface QueryMarketPriceRequest {
    /**
     * QueryMarketPriceRequest is request type for the Query/Params `MarketPrice`
     * RPC method.
     */
    id: number;
}
/**
 * QueryMarketPriceRequest is request type for the Query/Params `MarketPrice`
 * RPC method.
 */
export interface QueryMarketPriceRequestSDKType {
    id: number;
}
/**
 * QueryMarketPriceResponse is response type for the Query/Params `MarketPrice`
 * RPC method.
 */
export interface QueryMarketPriceResponse {
    marketPrice?: MarketPrice;
}
/**
 * QueryMarketPriceResponse is response type for the Query/Params `MarketPrice`
 * RPC method.
 */
export interface QueryMarketPriceResponseSDKType {
    market_price?: MarketPriceSDKType;
}
/**
 * QueryAllMarketPricesRequest is request type for the Query/Params
 * `AllMarketPrices` RPC method.
 */
export interface QueryAllMarketPricesRequest {
    pagination?: PageRequest;
}
/**
 * QueryAllMarketPricesRequest is request type for the Query/Params
 * `AllMarketPrices` RPC method.
 */
export interface QueryAllMarketPricesRequestSDKType {
    pagination?: PageRequestSDKType;
}
/**
 * QueryAllMarketPricesResponse is response type for the Query/Params
 * `AllMarketPrices` RPC method.
 */
export interface QueryAllMarketPricesResponse {
    marketPrices: MarketPrice[];
    pagination?: PageResponse;
}
/**
 * QueryAllMarketPricesResponse is response type for the Query/Params
 * `AllMarketPrices` RPC method.
 */
export interface QueryAllMarketPricesResponseSDKType {
    market_prices: MarketPriceSDKType[];
    pagination?: PageResponseSDKType;
}
/**
 * QueryMarketParamsRequest is request type for the Query/Params `MarketParams`
 * RPC method.
 */
export interface QueryMarketParamRequest {
    /**
     * QueryMarketParamsRequest is request type for the Query/Params `MarketParams`
     * RPC method.
     */
    id: number;
}
/**
 * QueryMarketParamsRequest is request type for the Query/Params `MarketParams`
 * RPC method.
 */
export interface QueryMarketParamRequestSDKType {
    id: number;
}
/**
 * QueryMarketParamResponse is response type for the Query/Params `MarketParams`
 * RPC method.
 */
export interface QueryMarketParamResponse {
    marketParam?: MarketParam;
}
/**
 * QueryMarketParamResponse is response type for the Query/Params `MarketParams`
 * RPC method.
 */
export interface QueryMarketParamResponseSDKType {
    market_param?: MarketParamSDKType;
}
/**
 * QueryAllMarketParamsRequest is request type for the Query/Params
 * `AllMarketParams` RPC method.
 */
export interface QueryAllMarketParamsRequest {
    pagination?: PageRequest;
}
/**
 * QueryAllMarketParamsRequest is request type for the Query/Params
 * `AllMarketParams` RPC method.
 */
export interface QueryAllMarketParamsRequestSDKType {
    pagination?: PageRequestSDKType;
}
/**
 * QueryAllMarketParamsResponse is response type for the Query/Params
 * `AllMarketParams` RPC method.
 */
export interface QueryAllMarketParamsResponse {
    marketParams: MarketParam[];
    pagination?: PageResponse;
}
/**
 * QueryAllMarketParamsResponse is response type for the Query/Params
 * `AllMarketParams` RPC method.
 */
export interface QueryAllMarketParamsResponseSDKType {
    market_params: MarketParamSDKType[];
    pagination?: PageResponseSDKType;
}
export declare const QueryMarketPriceRequest: {
    encode(message: QueryMarketPriceRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketPriceRequest;
    fromPartial(object: DeepPartial<QueryMarketPriceRequest>): QueryMarketPriceRequest;
};
export declare const QueryMarketPriceResponse: {
    encode(message: QueryMarketPriceResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketPriceResponse;
    fromPartial(object: DeepPartial<QueryMarketPriceResponse>): QueryMarketPriceResponse;
};
export declare const QueryAllMarketPricesRequest: {
    encode(message: QueryAllMarketPricesRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMarketPricesRequest;
    fromPartial(object: DeepPartial<QueryAllMarketPricesRequest>): QueryAllMarketPricesRequest;
};
export declare const QueryAllMarketPricesResponse: {
    encode(message: QueryAllMarketPricesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMarketPricesResponse;
    fromPartial(object: DeepPartial<QueryAllMarketPricesResponse>): QueryAllMarketPricesResponse;
};
export declare const QueryMarketParamRequest: {
    encode(message: QueryMarketParamRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketParamRequest;
    fromPartial(object: DeepPartial<QueryMarketParamRequest>): QueryMarketParamRequest;
};
export declare const QueryMarketParamResponse: {
    encode(message: QueryMarketParamResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketParamResponse;
    fromPartial(object: DeepPartial<QueryMarketParamResponse>): QueryMarketParamResponse;
};
export declare const QueryAllMarketParamsRequest: {
    encode(message: QueryAllMarketParamsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMarketParamsRequest;
    fromPartial(object: DeepPartial<QueryAllMarketParamsRequest>): QueryAllMarketParamsRequest;
};
export declare const QueryAllMarketParamsResponse: {
    encode(message: QueryAllMarketParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMarketParamsResponse;
    fromPartial(object: DeepPartial<QueryAllMarketParamsResponse>): QueryAllMarketParamsResponse;
};
