import { MarketParam, MarketParamSDKType } from "./market_param";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/**
 * MsgCreateOracleMarket is a message used by x/gov for creating a new oracle
 * market.
 */

export interface MsgCreateOracleMarket {
  /** The address that controls the module. */
  auhority: string;
  /** `params` defines parameters for the new oracle market. */

  params?: MarketParam;
}
/**
 * MsgCreateOracleMarket is a message used by x/gov for creating a new oracle
 * market.
 */

export interface MsgCreateOracleMarketSDKType {
  /** The address that controls the module. */
  auhority: string;
  /** `params` defines parameters for the new oracle market. */

  params?: MarketParamSDKType;
}
/** MsgCreateOracleMarketResponse defines the CreateOracleMarket response type. */

export interface MsgCreateOracleMarketResponse {}
/** MsgCreateOracleMarketResponse defines the CreateOracleMarket response type. */

export interface MsgCreateOracleMarketResponseSDKType {}
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
  /** The id of market to update */
  market_id: number;
  /** The updated price */

  price: Long;
}
/**
 * MsgUpdateMarketPricesResponse defines the MsgUpdateMarketPrices response
 * type.
 */

export interface MsgUpdateMarketPricesResponse {}
/**
 * MsgUpdateMarketPricesResponse defines the MsgUpdateMarketPrices response
 * type.
 */

export interface MsgUpdateMarketPricesResponseSDKType {}

function createBaseMsgCreateOracleMarket(): MsgCreateOracleMarket {
  return {
    auhority: "",
    params: undefined
  };
}

export const MsgCreateOracleMarket = {
  encode(message: MsgCreateOracleMarket, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.auhority !== "") {
      writer.uint32(10).string(message.auhority);
    }

    if (message.params !== undefined) {
      MarketParam.encode(message.params, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateOracleMarket {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateOracleMarket();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.auhority = reader.string();
          break;

        case 2:
          message.params = MarketParam.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgCreateOracleMarket>): MsgCreateOracleMarket {
    const message = createBaseMsgCreateOracleMarket();
    message.auhority = object.auhority ?? "";
    message.params = object.params !== undefined && object.params !== null ? MarketParam.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseMsgCreateOracleMarketResponse(): MsgCreateOracleMarketResponse {
  return {};
}

export const MsgCreateOracleMarketResponse = {
  encode(_: MsgCreateOracleMarketResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCreateOracleMarketResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCreateOracleMarketResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<MsgCreateOracleMarketResponse>): MsgCreateOracleMarketResponse {
    const message = createBaseMsgCreateOracleMarketResponse();
    return message;
  }

};

function createBaseMsgUpdateMarketPrices(): MsgUpdateMarketPrices {
  return {
    marketPriceUpdates: []
  };
}

export const MsgUpdateMarketPrices = {
  encode(message: MsgUpdateMarketPrices, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.marketPriceUpdates) {
      MsgUpdateMarketPrices_MarketPrice.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMarketPrices {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateMarketPrices();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketPriceUpdates.push(MsgUpdateMarketPrices_MarketPrice.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdateMarketPrices>): MsgUpdateMarketPrices {
    const message = createBaseMsgUpdateMarketPrices();
    message.marketPriceUpdates = object.marketPriceUpdates?.map(e => MsgUpdateMarketPrices_MarketPrice.fromPartial(e)) || [];
    return message;
  }

};

function createBaseMsgUpdateMarketPrices_MarketPrice(): MsgUpdateMarketPrices_MarketPrice {
  return {
    marketId: 0,
    price: Long.UZERO
  };
}

export const MsgUpdateMarketPrices_MarketPrice = {
  encode(message: MsgUpdateMarketPrices_MarketPrice, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.marketId !== 0) {
      writer.uint32(8).uint32(message.marketId);
    }

    if (!message.price.isZero()) {
      writer.uint32(16).uint64(message.price);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMarketPrices_MarketPrice {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateMarketPrices_MarketPrice();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketId = reader.uint32();
          break;

        case 2:
          message.price = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgUpdateMarketPrices_MarketPrice>): MsgUpdateMarketPrices_MarketPrice {
    const message = createBaseMsgUpdateMarketPrices_MarketPrice();
    message.marketId = object.marketId ?? 0;
    message.price = object.price !== undefined && object.price !== null ? Long.fromValue(object.price) : Long.UZERO;
    return message;
  }

};

function createBaseMsgUpdateMarketPricesResponse(): MsgUpdateMarketPricesResponse {
  return {};
}

export const MsgUpdateMarketPricesResponse = {
  encode(_: MsgUpdateMarketPricesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateMarketPricesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateMarketPricesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<MsgUpdateMarketPricesResponse>): MsgUpdateMarketPricesResponse {
    const message = createBaseMsgUpdateMarketPricesResponse();
    return message;
  }

};