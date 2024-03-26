import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Params stores `x/vault` parameters. */

export interface Params {
  /**
   * The number of layers of orders a vault places. For example if
   * `layers=2`, a vault places 2 asks and 2 bids.
   */
  layers: number;
  /** The minimum base spread when a vault quotes around reservation price. */

  spreadMinPpm: number;
  /**
   * The buffer amount to add to min_price_change_ppm to arrive at `spread`
   * according to formula:
   * `spread = max(spread_min_ppm, min_price_change_ppm + spread_buffer_ppm)`.
   */

  spreadBufferPpm: number;
  /** The factor that determines how aggressive a vault skews its orders. */

  skewFactorPpm: number;
  /** The percentage of vault equity that each order is sized at. */

  orderSizePpm: number;
  /** The duration that a vault's orders are valid for. */

  orderExpirationSeconds: number;
}
/** Params stores `x/vault` parameters. */

export interface ParamsSDKType {
  /**
   * The number of layers of orders a vault places. For example if
   * `layers=2`, a vault places 2 asks and 2 bids.
   */
  layers: number;
  /** The minimum base spread when a vault quotes around reservation price. */

  spread_min_ppm: number;
  /**
   * The buffer amount to add to min_price_change_ppm to arrive at `spread`
   * according to formula:
   * `spread = max(spread_min_ppm, min_price_change_ppm + spread_buffer_ppm)`.
   */

  spread_buffer_ppm: number;
  /** The factor that determines how aggressive a vault skews its orders. */

  skew_factor_ppm: number;
  /** The percentage of vault equity that each order is sized at. */

  order_size_ppm: number;
  /** The duration that a vault's orders are valid for. */

  order_expiration_seconds: number;
}

function createBaseParams(): Params {
  return {
    layers: 0,
    spreadMinPpm: 0,
    spreadBufferPpm: 0,
    skewFactorPpm: 0,
    orderSizePpm: 0,
    orderExpirationSeconds: 0
  };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.layers !== 0) {
      writer.uint32(8).uint32(message.layers);
    }

    if (message.spreadMinPpm !== 0) {
      writer.uint32(16).uint32(message.spreadMinPpm);
    }

    if (message.spreadBufferPpm !== 0) {
      writer.uint32(24).uint32(message.spreadBufferPpm);
    }

    if (message.skewFactorPpm !== 0) {
      writer.uint32(32).uint32(message.skewFactorPpm);
    }

    if (message.orderSizePpm !== 0) {
      writer.uint32(40).uint32(message.orderSizePpm);
    }

    if (message.orderExpirationSeconds !== 0) {
      writer.uint32(48).uint32(message.orderExpirationSeconds);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Params {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.layers = reader.uint32();
          break;

        case 2:
          message.spreadMinPpm = reader.uint32();
          break;

        case 3:
          message.spreadBufferPpm = reader.uint32();
          break;

        case 4:
          message.skewFactorPpm = reader.uint32();
          break;

        case 5:
          message.orderSizePpm = reader.uint32();
          break;

        case 6:
          message.orderExpirationSeconds = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Params>): Params {
    const message = createBaseParams();
    message.layers = object.layers ?? 0;
    message.spreadMinPpm = object.spreadMinPpm ?? 0;
    message.spreadBufferPpm = object.spreadBufferPpm ?? 0;
    message.skewFactorPpm = object.skewFactorPpm ?? 0;
    message.orderSizePpm = object.orderSizePpm ?? 0;
    message.orderExpirationSeconds = object.orderExpirationSeconds ?? 0;
    return message;
  }

};