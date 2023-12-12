import { BinaryReader, BinaryWriter } from "../../binary";
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
   */
  exponent: number;
  /**
   * The minimum number of exchanges that should be reporting a live price for
   * a price update to be considered valid.
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
   */
  exchangeConfigJson: string;
}
export interface MarketParamProtoMsg {
  typeUrl: "/dydxprotocol.prices.MarketParam";
  value: Uint8Array;
}
/**
 * MarketParam represents the x/prices configuration for markets, including
 * representing price values, resolving markets on individual exchanges, and
 * generating price updates. This configuration is specific to the quote
 * currency.
 */
export interface MarketParamAmino {
  /** Unique, sequentially-generated value. */
  id?: number;
  /** The human-readable name of the market pair (e.g. `BTC-USD`). */
  pair?: string;
  /**
   * Static value. The exponent of the price.
   * For example if `Exponent == -5` then a `Value` of `1,000,000,000`
   * represents ``$10,000`. Therefore `10 ^ Exponent` represents the smallest
   * price step (in dollars) that can be recorded.
   */
  exponent?: number;
  /**
   * The minimum number of exchanges that should be reporting a live price for
   * a price update to be considered valid.
   */
  min_exchanges?: number;
  /**
   * The minimum allowable change in `price` value that would cause a price
   * update on the network. Measured as `1e-6` (parts per million).
   */
  min_price_change_ppm?: number;
  /**
   * A string of json that encodes the configuration for resolving the price
   * of this market on various exchanges.
   */
  exchange_config_json?: string;
}
export interface MarketParamAminoMsg {
  type: "/dydxprotocol.prices.MarketParam";
  value: MarketParamAmino;
}
/**
 * MarketParam represents the x/prices configuration for markets, including
 * representing price values, resolving markets on individual exchanges, and
 * generating price updates. This configuration is specific to the quote
 * currency.
 */
export interface MarketParamSDKType {
  id: number;
  pair: string;
  exponent: number;
  min_exchanges: number;
  min_price_change_ppm: number;
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
  typeUrl: "/dydxprotocol.prices.MarketParam",
  encode(message: MarketParam, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): MarketParam {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<MarketParam>): MarketParam {
    const message = createBaseMarketParam();
    message.id = object.id ?? 0;
    message.pair = object.pair ?? "";
    message.exponent = object.exponent ?? 0;
    message.minExchanges = object.minExchanges ?? 0;
    message.minPriceChangePpm = object.minPriceChangePpm ?? 0;
    message.exchangeConfigJson = object.exchangeConfigJson ?? "";
    return message;
  },
  fromAmino(object: MarketParamAmino): MarketParam {
    const message = createBaseMarketParam();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    if (object.pair !== undefined && object.pair !== null) {
      message.pair = object.pair;
    }
    if (object.exponent !== undefined && object.exponent !== null) {
      message.exponent = object.exponent;
    }
    if (object.min_exchanges !== undefined && object.min_exchanges !== null) {
      message.minExchanges = object.min_exchanges;
    }
    if (object.min_price_change_ppm !== undefined && object.min_price_change_ppm !== null) {
      message.minPriceChangePpm = object.min_price_change_ppm;
    }
    if (object.exchange_config_json !== undefined && object.exchange_config_json !== null) {
      message.exchangeConfigJson = object.exchange_config_json;
    }
    return message;
  },
  toAmino(message: MarketParam): MarketParamAmino {
    const obj: any = {};
    obj.id = message.id;
    obj.pair = message.pair;
    obj.exponent = message.exponent;
    obj.min_exchanges = message.minExchanges;
    obj.min_price_change_ppm = message.minPriceChangePpm;
    obj.exchange_config_json = message.exchangeConfigJson;
    return obj;
  },
  fromAminoMsg(object: MarketParamAminoMsg): MarketParam {
    return MarketParam.fromAmino(object.value);
  },
  fromProtoMsg(message: MarketParamProtoMsg): MarketParam {
    return MarketParam.decode(message.value);
  },
  toProto(message: MarketParam): Uint8Array {
    return MarketParam.encode(message).finish();
  },
  toProtoMsg(message: MarketParam): MarketParamProtoMsg {
    return {
      typeUrl: "/dydxprotocol.prices.MarketParam",
      value: MarketParam.encode(message).finish()
    };
  }
};