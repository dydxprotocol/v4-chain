/// <reference types="long" />
import { MarketParam, MarketParamSDKType } from "./market_param";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/**
 * MsgCreateOracleMarket is a message used by x/gov for creating a new oracle
 * market.
 */
export interface MsgCreateOracleMarket {
    /** The address that controls the module. */
    authority: string;
    /** `params` defines parameters for the new oracle market. */
    params?: MarketParam;
}
/**
 * MsgCreateOracleMarket is a message used by x/gov for creating a new oracle
 * market.
 */
export interface MsgCreateOracleMarketSDKType {
    authority: string;
    params?: MarketParamSDKType;
}
/** MsgCreateOracleMarketResponse defines the CreateOracleMarket response type. */
export interface MsgCreateOracleMarketResponse {
}
/** MsgCreateOracleMarketResponse defines the CreateOracleMarket response type. */
export interface MsgCreateOracleMarketResponseSDKType {
}
/** MsgUpdateMarketPrices is a request type for the UpdateMarketPrices method. */
export interface MsgUpdateMarketPrices {
    marketPriceUpdates: MsgUpdateMarketPrices_MarketPrice[];
}
/** MsgUpdateMarketPrices is a request type for the UpdateMarketPrices method. */
export interface MsgUpdateMarketPricesSDKType {
    market_price_updates: MsgUpdateMarketPrices_MarketPriceSDKType[];
}
/** MarketPrice represents a price update for a single market */
export interface MsgUpdateMarketPrices_MarketPrice {
    /** The id of market to update */
    marketId: number;
    /** The updated price */
    price: Long;
}
/** MarketPrice represents a price update for a single market */
export interface MsgUpdateMarketPrices_MarketPriceSDKType {
    market_id: number;
    price: Long;
}
/**
 * MsgUpdateMarketPricesResponse defines the MsgUpdateMarketPrices response
 * type.
 */
export interface MsgUpdateMarketPricesResponse {
}
/**
 * MsgUpdateMarketPricesResponse defines the MsgUpdateMarketPrices response
 * type.
 */
export interface MsgUpdateMarketPricesResponseSDKType {
}
/**
 * MsgUpdateMarketParam is a message used by x/gov for updating the parameters
 * of an oracle market.
 */
export interface MsgUpdateMarketParam {
    authority: string;
    /** The market param to update. Each field must be set. */
    marketParam?: MarketParam;
}
/**
 * MsgUpdateMarketParam is a message used by x/gov for updating the parameters
 * of an oracle market.
 */
export interface MsgUpdateMarketParamSDKType {
    authority: string;
    market_param?: MarketParamSDKType;
}
/** MsgUpdateMarketParamResponse defines the UpdateMarketParam response type. */
export interface MsgUpdateMarketParamResponse {
}
/** MsgUpdateMarketParamResponse defines the UpdateMarketParam response type. */
export interface MsgUpdateMarketParamResponseSDKType {
}
export declare const MsgCreateOracleMarket: {
    encode(message: MsgCreateOracleMarket, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateOracleMarket;
    fromPartial(object: DeepPartial<MsgCreateOracleMarket>): MsgCreateOracleMarket;
};
export declare const MsgCreateOracleMarketResponse: {
    encode(_: MsgCreateOracleMarketResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateOracleMarketResponse;
    fromPartial(_: DeepPartial<MsgCreateOracleMarketResponse>): MsgCreateOracleMarketResponse;
};
export declare const MsgUpdateMarketPrices: {
    encode(message: MsgUpdateMarketPrices, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMarketPrices;
    fromPartial(object: DeepPartial<MsgUpdateMarketPrices>): MsgUpdateMarketPrices;
};
export declare const MsgUpdateMarketPrices_MarketPrice: {
    encode(message: MsgUpdateMarketPrices_MarketPrice, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMarketPrices_MarketPrice;
    fromPartial(object: DeepPartial<MsgUpdateMarketPrices_MarketPrice>): MsgUpdateMarketPrices_MarketPrice;
};
export declare const MsgUpdateMarketPricesResponse: {
    encode(_: MsgUpdateMarketPricesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMarketPricesResponse;
    fromPartial(_: DeepPartial<MsgUpdateMarketPricesResponse>): MsgUpdateMarketPricesResponse;
};
export declare const MsgUpdateMarketParam: {
    encode(message: MsgUpdateMarketParam, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMarketParam;
    fromPartial(object: DeepPartial<MsgUpdateMarketParam>): MsgUpdateMarketParam;
};
export declare const MsgUpdateMarketParamResponse: {
    encode(_: MsgUpdateMarketParamResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMarketParamResponse;
    fromPartial(_: DeepPartial<MsgUpdateMarketParamResponse>): MsgUpdateMarketParamResponse;
};
