import * as _m0 from "protobufjs/minimal";
import { Long, DeepPartial } from "../../helpers";
/**
 * MarketMapperRevShareDetails specifies any details associated with the market
 * mapper revenue share
 */

export interface MarketMapperRevShareDetails {
  /** Unix timestamp recorded when the market revenue share expires */
  expirationTs: Long;
}
/**
 * MarketMapperRevShareDetails specifies any details associated with the market
 * mapper revenue share
 */

export interface MarketMapperRevShareDetailsSDKType {
  /** Unix timestamp recorded when the market revenue share expires */
  expiration_ts: Long;
}

function createBaseMarketMapperRevShareDetails(): MarketMapperRevShareDetails {
  return {
    expirationTs: Long.UZERO
  };
}

export const MarketMapperRevShareDetails = {
  encode(message: MarketMapperRevShareDetails, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.expirationTs.isZero()) {
      writer.uint32(8).uint64(message.expirationTs);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketMapperRevShareDetails {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketMapperRevShareDetails();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.expirationTs = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketMapperRevShareDetails>): MarketMapperRevShareDetails {
    const message = createBaseMarketMapperRevShareDetails();
    message.expirationTs = object.expirationTs !== undefined && object.expirationTs !== null ? Long.fromValue(object.expirationTs) : Long.UZERO;
    return message;
  }

};