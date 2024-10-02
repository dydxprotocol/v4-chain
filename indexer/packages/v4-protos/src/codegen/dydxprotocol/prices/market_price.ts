import * as _m0 from "protobufjs/minimal";
import { Long, DeepPartial } from "../../helpers";
/** MarketPrice is used by the application to store/retrieve oracle price. */

export interface MarketPrice {
  /** Unique, sequentially-generated value that matches `MarketParam`. */
  id: number;
  /**
   * Static value. The exponent of the price. See the comment on the duplicate
   * MarketParam field for more information.
   * 
   * As of v7.1.x, this value is determined from the marketmap instead of
   * needing to match the MarketParam field.
   */

  exponent: number;
  /**
   * The variable value that is updated by oracle price updates. `0` if it has
   * never been updated, `>0` otherwise.
   */

  price: Long;
}
/** MarketPrice is used by the application to store/retrieve oracle price. */

export interface MarketPriceSDKType {
  /** Unique, sequentially-generated value that matches `MarketParam`. */
  id: number;
  /**
   * Static value. The exponent of the price. See the comment on the duplicate
   * MarketParam field for more information.
   * 
   * As of v7.1.x, this value is determined from the marketmap instead of
   * needing to match the MarketParam field.
   */

  exponent: number;
  /**
   * The variable value that is updated by oracle price updates. `0` if it has
   * never been updated, `>0` otherwise.
   */

  price: Long;
}

function createBaseMarketPrice(): MarketPrice {
  return {
    id: 0,
    exponent: 0,
    price: Long.UZERO
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

    if (!message.price.isZero()) {
      writer.uint32(24).uint64(message.price);
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
          message.price = (reader.uint64() as Long);
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
    message.price = object.price !== undefined && object.price !== null ? Long.fromValue(object.price) : Long.UZERO;
    return message;
  }

};