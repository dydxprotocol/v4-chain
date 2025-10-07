import * as _m0 from "protobufjs/minimal";
import { Long, DeepPartial } from "../../helpers";
/**
 * FeeDiscountCampaignParams defines a fee discount period for a specific CLOB
 * pair
 */

export interface FeeDiscountCampaignParams {
  /** The CLOB Pair ID this fee holiday applies to */
  clobPairId: number;
  /** Start time (Unix timestamp in seconds) */

  startTimeUnix: Long;
  /** End time (Unix timestamp in seconds) */

  endTimeUnix: Long;
  /**
   * Percentage of normal fee to charge during the campaign (in parts per
   * million) 0 = completely free (100% discount) 500000 = charge 50% of normal
   * fee (50% discount) 1000000 = charge 100% of normal fee (no discount)
   */

  chargePpm: number;
}
/**
 * FeeDiscountCampaignParams defines a fee discount period for a specific CLOB
 * pair
 */

export interface FeeDiscountCampaignParamsSDKType {
  /** The CLOB Pair ID this fee holiday applies to */
  clob_pair_id: number;
  /** Start time (Unix timestamp in seconds) */

  start_time_unix: Long;
  /** End time (Unix timestamp in seconds) */

  end_time_unix: Long;
  /**
   * Percentage of normal fee to charge during the campaign (in parts per
   * million) 0 = completely free (100% discount) 500000 = charge 50% of normal
   * fee (50% discount) 1000000 = charge 100% of normal fee (no discount)
   */

  charge_ppm: number;
}

function createBaseFeeDiscountCampaignParams(): FeeDiscountCampaignParams {
  return {
    clobPairId: 0,
    startTimeUnix: Long.ZERO,
    endTimeUnix: Long.ZERO,
    chargePpm: 0
  };
}

export const FeeDiscountCampaignParams = {
  encode(message: FeeDiscountCampaignParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobPairId !== 0) {
      writer.uint32(8).uint32(message.clobPairId);
    }

    if (!message.startTimeUnix.isZero()) {
      writer.uint32(16).int64(message.startTimeUnix);
    }

    if (!message.endTimeUnix.isZero()) {
      writer.uint32(24).int64(message.endTimeUnix);
    }

    if (message.chargePpm !== 0) {
      writer.uint32(32).uint32(message.chargePpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FeeDiscountCampaignParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFeeDiscountCampaignParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPairId = reader.uint32();
          break;

        case 2:
          message.startTimeUnix = (reader.int64() as Long);
          break;

        case 3:
          message.endTimeUnix = (reader.int64() as Long);
          break;

        case 4:
          message.chargePpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<FeeDiscountCampaignParams>): FeeDiscountCampaignParams {
    const message = createBaseFeeDiscountCampaignParams();
    message.clobPairId = object.clobPairId ?? 0;
    message.startTimeUnix = object.startTimeUnix !== undefined && object.startTimeUnix !== null ? Long.fromValue(object.startTimeUnix) : Long.ZERO;
    message.endTimeUnix = object.endTimeUnix !== undefined && object.endTimeUnix !== null ? Long.fromValue(object.endTimeUnix) : Long.ZERO;
    message.chargePpm = object.chargePpm ?? 0;
    return message;
  }

};