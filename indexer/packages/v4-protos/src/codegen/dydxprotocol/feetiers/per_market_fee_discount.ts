import { Timestamp } from "../../google/protobuf/timestamp";
import * as _m0 from "protobufjs/minimal";
import { toTimestamp, fromTimestamp, DeepPartial } from "../../helpers";
/**
 * PerMarketFeeDiscountParams defines a fee discount period for a specific CLOB
 * pair
 */

export interface PerMarketFeeDiscountParams {
  /** The CLOB Pair ID this fee holiday applies to */
  clobPairId: number;
  /** Start time */

  startTime?: Date;
  /** End time */

  endTime?: Date;
  /**
   * Percentage of normal fee to charge during the period (in parts per
   * million) 0 = completely free (100% discount) 500000 = charge 50% of normal
   * fee (50% discount) 1000000 = charge 100% of normal fee (no discount)
   */

  chargePpm: number;
}
/**
 * PerMarketFeeDiscountParams defines a fee discount period for a specific CLOB
 * pair
 */

export interface PerMarketFeeDiscountParamsSDKType {
  /** The CLOB Pair ID this fee holiday applies to */
  clob_pair_id: number;
  /** Start time */

  start_time?: Date;
  /** End time */

  end_time?: Date;
  /**
   * Percentage of normal fee to charge during the period (in parts per
   * million) 0 = completely free (100% discount) 500000 = charge 50% of normal
   * fee (50% discount) 1000000 = charge 100% of normal fee (no discount)
   */

  charge_ppm: number;
}

function createBasePerMarketFeeDiscountParams(): PerMarketFeeDiscountParams {
  return {
    clobPairId: 0,
    startTime: undefined,
    endTime: undefined,
    chargePpm: 0
  };
}

export const PerMarketFeeDiscountParams = {
  encode(message: PerMarketFeeDiscountParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.clobPairId !== 0) {
      writer.uint32(8).uint32(message.clobPairId);
    }

    if (message.startTime !== undefined) {
      Timestamp.encode(toTimestamp(message.startTime), writer.uint32(18).fork()).ldelim();
    }

    if (message.endTime !== undefined) {
      Timestamp.encode(toTimestamp(message.endTime), writer.uint32(26).fork()).ldelim();
    }

    if (message.chargePpm !== 0) {
      writer.uint32(32).uint32(message.chargePpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PerMarketFeeDiscountParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePerMarketFeeDiscountParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.clobPairId = reader.uint32();
          break;

        case 2:
          message.startTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;

        case 3:
          message.endTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
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

  fromPartial(object: DeepPartial<PerMarketFeeDiscountParams>): PerMarketFeeDiscountParams {
    const message = createBasePerMarketFeeDiscountParams();
    message.clobPairId = object.clobPairId ?? 0;
    message.startTime = object.startTime ?? undefined;
    message.endTime = object.endTime ?? undefined;
    message.chargePpm = object.chargePpm ?? 0;
    return message;
  }

};