import * as _m0 from "protobufjs/minimal";
import { Long, DeepPartial } from "../../helpers";
/** MarketPrice is used by the application to store/retrieve oracle prices. */

export interface MarketPrice {
  /** Unique, sequentially-generated value that matches `MarketParam`. */
  id: number;
  /**
   * Static value. The exponent of the price. See the comment on the duplicate
   * MarketParam field for more information.
   */

  exponent: number;
  /**
   * The spot price value that is updated by oracle price updates. `0` if it has
   * never been updated, `>0` otherwise.
   */

  spotPrice: Long;
  /**
   * The pnl price value that is updated by oracle price updates. `0` if it has
   * never been updated, `>0` otherwise.
   */

  pnlPrice: Long;
}
/** MarketPrice is used by the application to store/retrieve oracle prices. */

export interface MarketPriceSDKType {
  /** Unique, sequentially-generated value that matches `MarketParam`. */
  id: number;
  /**
   * Static value. The exponent of the price. See the comment on the duplicate
   * MarketParam field for more information.
   */

  exponent: number;
  /**
   * The spot price value that is updated by oracle price updates. `0` if it has
   * never been updated, `>0` otherwise.
   */

  spot_price: Long;
  /**
   * The pnl price value that is updated by oracle price updates. `0` if it has
   * never been updated, `>0` otherwise.
   */

  pnl_price: Long;
}
/** MarketSpotPrice is used by the application to store/retrieve spot prices. */

export interface MarketSpotPrice {
  /** Unique, sequentially-generated value that matches `MarketParam`. */
  id: number;
  /**
   * Static value. The exponent of the price. See the comment on the duplicate
   * MarketParam field for more information.
   */

  exponent: number;
  /**
   * The spot price value that is updated by oracle price updates. `0` if it has
   * never been updated, `>0` otherwise.
   */

  spotPrice: Long;
}
/** MarketSpotPrice is used by the application to store/retrieve spot prices. */

export interface MarketSpotPriceSDKType {
  /** Unique, sequentially-generated value that matches `MarketParam`. */
  id: number;
  /**
   * Static value. The exponent of the price. See the comment on the duplicate
   * MarketParam field for more information.
   */

  exponent: number;
  /**
   * The spot price value that is updated by oracle price updates. `0` if it has
   * never been updated, `>0` otherwise.
   */

  spot_price: Long;
}
/** MarketPriceUpdate is used to update the price of a single market. */

export interface MarketPriceUpdate {
  /** The id of market to update */
  marketId: number;
  /** The updated spot price */

  spotPrice: Long;
  /** The updated pnl price */

  pnlPrice: Long;
}
/** MarketPriceUpdate is used to update the price of a single market. */

export interface MarketPriceUpdateSDKType {
  /** The id of market to update */
  market_id: number;
  /** The updated spot price */

  spot_price: Long;
  /** The updated pnl price */

  pnl_price: Long;
}
/** MarketSpotPriceUpdate is used to update the spot price of a single market. */

export interface MarketSpotPriceUpdate {
  /** The id of market to update */
  marketId: number;
  /** The updated spot price */

  spotPrice: Long;
}
/** MarketSpotPriceUpdate is used to update the spot price of a single market. */

export interface MarketSpotPriceUpdateSDKType {
  /** The id of market to update */
  market_id: number;
  /** The updated spot price */

  spot_price: Long;
}
/** MarketPnlPriceUpdate is used to update the pnl price of a single market. */

export interface MarketPnlPriceUpdate {
  /** The id of market to update */
  marketId: number;
  /** The updated pnl price */

  pnlPrice: Long;
}
/** MarketPnlPriceUpdate is used to update the pnl price of a single market. */

export interface MarketPnlPriceUpdateSDKType {
  /** The id of market to update */
  market_id: number;
  /** The updated pnl price */

  pnl_price: Long;
}
/** MarketPriceUpdates is a collection of MarketPriceUpdate messages. */

export interface MarketPriceUpdates {
  marketPriceUpdates: MarketPriceUpdate[];
}
/** MarketPriceUpdates is a collection of MarketPriceUpdate messages. */

export interface MarketPriceUpdatesSDKType {
  market_price_updates: MarketPriceUpdateSDKType[];
}

function createBaseMarketPrice(): MarketPrice {
  return {
    id: 0,
    exponent: 0,
    spotPrice: Long.UZERO,
    pnlPrice: Long.UZERO
  };
}

export const MarketPrice = {
  encode(message: MarketPrice, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    if (message.exponent !== 0) {
      writer.uint32(16).sint32(message.exponent);
    }

    if (!message.spotPrice.isZero()) {
      writer.uint32(24).uint64(message.spotPrice);
    }

    if (!message.pnlPrice.isZero()) {
      writer.uint32(32).uint64(message.pnlPrice);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketPrice {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketPrice();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        case 2:
          message.exponent = reader.sint32();
          break;

        case 3:
          message.spotPrice = (reader.uint64() as Long);
          break;

        case 4:
          message.pnlPrice = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketPrice>): MarketPrice {
    const message = createBaseMarketPrice();
    message.id = object.id ?? 0;
    message.exponent = object.exponent ?? 0;
    message.spotPrice = object.spotPrice !== undefined && object.spotPrice !== null ? Long.fromValue(object.spotPrice) : Long.UZERO;
    message.pnlPrice = object.pnlPrice !== undefined && object.pnlPrice !== null ? Long.fromValue(object.pnlPrice) : Long.UZERO;
    return message;
  }

};

function createBaseMarketSpotPrice(): MarketSpotPrice {
  return {
    id: 0,
    exponent: 0,
    spotPrice: Long.UZERO
  };
}

export const MarketSpotPrice = {
  encode(message: MarketSpotPrice, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    if (message.exponent !== 0) {
      writer.uint32(16).sint32(message.exponent);
    }

    if (!message.spotPrice.isZero()) {
      writer.uint32(24).uint64(message.spotPrice);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketSpotPrice {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketSpotPrice();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        case 2:
          message.exponent = reader.sint32();
          break;

        case 3:
          message.spotPrice = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketSpotPrice>): MarketSpotPrice {
    const message = createBaseMarketSpotPrice();
    message.id = object.id ?? 0;
    message.exponent = object.exponent ?? 0;
    message.spotPrice = object.spotPrice !== undefined && object.spotPrice !== null ? Long.fromValue(object.spotPrice) : Long.UZERO;
    return message;
  }

};

function createBaseMarketPriceUpdate(): MarketPriceUpdate {
  return {
    marketId: 0,
    spotPrice: Long.UZERO,
    pnlPrice: Long.UZERO
  };
}

export const MarketPriceUpdate = {
  encode(message: MarketPriceUpdate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.marketId !== 0) {
      writer.uint32(8).uint32(message.marketId);
    }

    if (!message.spotPrice.isZero()) {
      writer.uint32(16).uint64(message.spotPrice);
    }

    if (!message.pnlPrice.isZero()) {
      writer.uint32(24).uint64(message.pnlPrice);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketPriceUpdate {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketPriceUpdate();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketId = reader.uint32();
          break;

        case 2:
          message.spotPrice = (reader.uint64() as Long);
          break;

        case 3:
          message.pnlPrice = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketPriceUpdate>): MarketPriceUpdate {
    const message = createBaseMarketPriceUpdate();
    message.marketId = object.marketId ?? 0;
    message.spotPrice = object.spotPrice !== undefined && object.spotPrice !== null ? Long.fromValue(object.spotPrice) : Long.UZERO;
    message.pnlPrice = object.pnlPrice !== undefined && object.pnlPrice !== null ? Long.fromValue(object.pnlPrice) : Long.UZERO;
    return message;
  }

};

function createBaseMarketSpotPriceUpdate(): MarketSpotPriceUpdate {
  return {
    marketId: 0,
    spotPrice: Long.UZERO
  };
}

export const MarketSpotPriceUpdate = {
  encode(message: MarketSpotPriceUpdate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.marketId !== 0) {
      writer.uint32(8).uint32(message.marketId);
    }

    if (!message.spotPrice.isZero()) {
      writer.uint32(16).uint64(message.spotPrice);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketSpotPriceUpdate {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketSpotPriceUpdate();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketId = reader.uint32();
          break;

        case 2:
          message.spotPrice = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketSpotPriceUpdate>): MarketSpotPriceUpdate {
    const message = createBaseMarketSpotPriceUpdate();
    message.marketId = object.marketId ?? 0;
    message.spotPrice = object.spotPrice !== undefined && object.spotPrice !== null ? Long.fromValue(object.spotPrice) : Long.UZERO;
    return message;
  }

};

function createBaseMarketPnlPriceUpdate(): MarketPnlPriceUpdate {
  return {
    marketId: 0,
    pnlPrice: Long.UZERO
  };
}

export const MarketPnlPriceUpdate = {
  encode(message: MarketPnlPriceUpdate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.marketId !== 0) {
      writer.uint32(8).uint32(message.marketId);
    }

    if (!message.pnlPrice.isZero()) {
      writer.uint32(16).uint64(message.pnlPrice);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketPnlPriceUpdate {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketPnlPriceUpdate();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketId = reader.uint32();
          break;

        case 2:
          message.pnlPrice = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketPnlPriceUpdate>): MarketPnlPriceUpdate {
    const message = createBaseMarketPnlPriceUpdate();
    message.marketId = object.marketId ?? 0;
    message.pnlPrice = object.pnlPrice !== undefined && object.pnlPrice !== null ? Long.fromValue(object.pnlPrice) : Long.UZERO;
    return message;
  }

};

function createBaseMarketPriceUpdates(): MarketPriceUpdates {
  return {
    marketPriceUpdates: []
  };
}

export const MarketPriceUpdates = {
  encode(message: MarketPriceUpdates, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.marketPriceUpdates) {
      MarketPriceUpdate.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketPriceUpdates {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketPriceUpdates();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.marketPriceUpdates.push(MarketPriceUpdate.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketPriceUpdates>): MarketPriceUpdates {
    const message = createBaseMarketPriceUpdates();
    message.marketPriceUpdates = object.marketPriceUpdates?.map(e => MarketPriceUpdate.fromPartial(e)) || [];
    return message;
  }

};