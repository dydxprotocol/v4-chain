import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * MarketParam represents the x/prices configuration for markets, including
 * representing price values, resolving markets on individual exchanges, and
 * generating price updates. This configuration is specific to the quote
 * currency.
 */

export interface MarketParam {
  /** Unique, sequentially-generated value. */
  id: number;
  /** The human-readable name of the market pair (e.g. `BTC-USD`). */

  pair: string;
  /**
   * Static value. The exponent of the price.
   * For example if `Exponent == -5` then a `Value` of `1,000,000,000`
   * represents ``$10,000`. Therefore `10 ^ Exponent` represents the smallest
   * price step (in dollars) that can be recorded.
   * 
   * Deprecated since v8.x. This value is now determined from the marketmap.
   */

  /** @deprecated */

  exponent: number;
  /**
   * The minimum number of exchanges that should be reporting a live price for
   * a price update to be considered valid.
   * 
   * Deprecated since v8.x. This value is now determined from the marketmap.
   */

  minExchanges: number;
  /**
   * The minimum allowable change in `price` value that would cause a price
   * update on the network. Measured as `1e-6` (parts per million).
   */

  minPriceChangePpm: number;
  /**
   * A string of json that encodes the configuration for resolving the price
   * of this market on various exchanges.
   * 
   * Deprecated since v8.x. This is now determined from the marketmap.
   */

  exchangeConfigJson: string;
}
/**
 * MarketParam represents the x/prices configuration for markets, including
 * representing price values, resolving markets on individual exchanges, and
 * generating price updates. This configuration is specific to the quote
 * currency.
 */

export interface MarketParamSDKType {
  /** Unique, sequentially-generated value. */
  id: number;
  /** The human-readable name of the market pair (e.g. `BTC-USD`). */

  pair: string;
  /**
   * Static value. The exponent of the price.
   * For example if `Exponent == -5` then a `Value` of `1,000,000,000`
   * represents ``$10,000`. Therefore `10 ^ Exponent` represents the smallest
   * price step (in dollars) that can be recorded.
   * 
   * Deprecated since v8.x. This value is now determined from the marketmap.
   */

  /** @deprecated */

  exponent: number;
  /**
   * The minimum number of exchanges that should be reporting a live price for
   * a price update to be considered valid.
   * 
   * Deprecated since v8.x. This value is now determined from the marketmap.
   */

  min_exchanges: number;
  /**
   * The minimum allowable change in `price` value that would cause a price
   * update on the network. Measured as `1e-6` (parts per million).
   */

  min_price_change_ppm: number;
  /**
   * A string of json that encodes the configuration for resolving the price
   * of this market on various exchanges.
   * 
   * Deprecated since v8.x. This is now determined from the marketmap.
   */

  exchange_config_json: string;
}

function createBaseMarketParam(): MarketParam {
  return {
    id: 0,
    pair: "",
    exponent: 0,
    minExchanges: 0,
    minPriceChangePpm: 0,
    exchangeConfigJson: ""
  };
}

export const MarketParam = {
  encode(message: MarketParam, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    if (message.pair !== "") {
      writer.uint32(18).string(message.pair);
    }

    if (message.exponent !== 0) {
      writer.uint32(24).sint32(message.exponent);
    }

    if (message.minExchanges !== 0) {
      writer.uint32(32).uint32(message.minExchanges);
    }

    if (message.minPriceChangePpm !== 0) {
      writer.uint32(40).uint32(message.minPriceChangePpm);
    }

    if (message.exchangeConfigJson !== "") {
      writer.uint32(50).string(message.exchangeConfigJson);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketParam {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketParam();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        case 2:
          message.pair = reader.string();
          break;

        case 3:
          message.exponent = reader.sint32();
          break;

        case 4:
          message.minExchanges = reader.uint32();
          break;

        case 5:
          message.minPriceChangePpm = reader.uint32();
          break;

        case 6:
          message.exchangeConfigJson = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketParam>): MarketParam {
    const message = createBaseMarketParam();
    message.id = object.id ?? 0;
    message.pair = object.pair ?? "";
    message.exponent = object.exponent ?? 0;
    message.minExchanges = object.minExchanges ?? 0;
    message.minPriceChangePpm = object.minPriceChangePpm ?? 0;
    message.exchangeConfigJson = object.exchangeConfigJson ?? "";
    return message;
  }

};