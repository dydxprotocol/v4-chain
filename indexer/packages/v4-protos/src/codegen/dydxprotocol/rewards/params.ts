import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Params defines the parameters for x/rewards module. */

export interface Params {
  /** The module account to distribute rewards from. */
  treasuryAccount: string;
  /** The denom of the rewards token. */

  denom: string;
  /**
   * The exponent of converting one unit of `denom` to a full coin.
   * For example, `denom=uatom, denom_exponent=-6` defines that
   * `1 uatom = 10^(-6) ATOM`. This conversion is needed since the
   * `market_id` retrieves the price of a full coin of the reward token.
   */

  denomExponent: number;
  /** The id of the market that has the price of the rewards token. */

  marketId: number;
  /**
   * The amount (in ppm) that fees are multiplied by to get
   * the maximum rewards amount.
   */

  feeMultiplierPpm: number;
}
/** Params defines the parameters for x/rewards module. */

export interface ParamsSDKType {
  /** The module account to distribute rewards from. */
  treasury_account: string;
  /** The denom of the rewards token. */

  denom: string;
  /**
   * The exponent of converting one unit of `denom` to a full coin.
   * For example, `denom=uatom, denom_exponent=-6` defines that
   * `1 uatom = 10^(-6) ATOM`. This conversion is needed since the
   * `market_id` retrieves the price of a full coin of the reward token.
   */

  denom_exponent: number;
  /** The id of the market that has the price of the rewards token. */

  market_id: number;
  /**
   * The amount (in ppm) that fees are multiplied by to get
   * the maximum rewards amount.
   */

  fee_multiplier_ppm: number;
}

function createBaseParams(): Params {
  return {
    treasuryAccount: "",
    denom: "",
    denomExponent: 0,
    marketId: 0,
    feeMultiplierPpm: 0
  };
}

export const Params = {
  encode(message: Params, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.treasuryAccount !== "") {
      writer.uint32(10).string(message.treasuryAccount);
    }

    if (message.denom !== "") {
      writer.uint32(18).string(message.denom);
    }

    if (message.denomExponent !== 0) {
      writer.uint32(24).sint32(message.denomExponent);
    }

    if (message.marketId !== 0) {
      writer.uint32(32).uint32(message.marketId);
    }

    if (message.feeMultiplierPpm !== 0) {
      writer.uint32(40).uint32(message.feeMultiplierPpm);
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
          message.treasuryAccount = reader.string();
          break;

        case 2:
          message.denom = reader.string();
          break;

        case 3:
          message.denomExponent = reader.sint32();
          break;

        case 4:
          message.marketId = reader.uint32();
          break;

        case 5:
          message.feeMultiplierPpm = reader.uint32();
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
    message.treasuryAccount = object.treasuryAccount ?? "";
    message.denom = object.denom ?? "";
    message.denomExponent = object.denomExponent ?? 0;
    message.marketId = object.marketId ?? 0;
    message.feeMultiplierPpm = object.feeMultiplierPpm ?? 0;
    return message;
  }

};