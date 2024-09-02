import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** LiquidationsConfig stores all configurable fields related to liquidations. */

export interface LiquidationsConfig {
  /**
   * The maximum liquidation fee (in parts-per-million). This fee goes
   * 100% to the insurance fund.
   */
  maxLiquidationFeePpm: number;
  /**
   * Limits around how many quote quantums from a single subaccount can
   * be liquidated within a single block.
   */

  subaccountBlockLimits?: SubaccountBlockLimits;
}
/** LiquidationsConfig stores all configurable fields related to liquidations. */

export interface LiquidationsConfigSDKType {
  /**
   * The maximum liquidation fee (in parts-per-million). This fee goes
   * 100% to the insurance fund.
   */
  max_liquidation_fee_ppm: number;
  /**
   * Limits around how many quote quantums from a single subaccount can
   * be liquidated within a single block.
   */

  subaccount_block_limits?: SubaccountBlockLimitsSDKType;
}
/**
 * SubaccountBlockLimits stores all configurable fields related to limits
 * around how many quote quantums from a single subaccount can
 * be liquidated within a single block.
 */

export interface SubaccountBlockLimits {
  /**
   * The maximum insurance-fund payout amount for a given subaccount
   * per block. I.e. how much it can cover for that subaccount.
   */
  maxQuantumsInsuranceLost: Long;
}
/**
 * SubaccountBlockLimits stores all configurable fields related to limits
 * around how many quote quantums from a single subaccount can
 * be liquidated within a single block.
 */

export interface SubaccountBlockLimitsSDKType {
  /**
   * The maximum insurance-fund payout amount for a given subaccount
   * per block. I.e. how much it can cover for that subaccount.
   */
  max_quantums_insurance_lost: Long;
}

function createBaseLiquidationsConfig(): LiquidationsConfig {
  return {
    maxLiquidationFeePpm: 0,
    subaccountBlockLimits: undefined
  };
}

export const LiquidationsConfig = {
  encode(message: LiquidationsConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.maxLiquidationFeePpm !== 0) {
      writer.uint32(8).uint32(message.maxLiquidationFeePpm);
    }

    if (message.subaccountBlockLimits !== undefined) {
      SubaccountBlockLimits.encode(message.subaccountBlockLimits, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LiquidationsConfig {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLiquidationsConfig();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.maxLiquidationFeePpm = reader.uint32();
          break;

        case 2:
          message.subaccountBlockLimits = SubaccountBlockLimits.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<LiquidationsConfig>): LiquidationsConfig {
    const message = createBaseLiquidationsConfig();
    message.maxLiquidationFeePpm = object.maxLiquidationFeePpm ?? 0;
    message.subaccountBlockLimits = object.subaccountBlockLimits !== undefined && object.subaccountBlockLimits !== null ? SubaccountBlockLimits.fromPartial(object.subaccountBlockLimits) : undefined;
    return message;
  }

};

function createBaseSubaccountBlockLimits(): SubaccountBlockLimits {
  return {
    maxQuantumsInsuranceLost: Long.UZERO
  };
}

export const SubaccountBlockLimits = {
  encode(message: SubaccountBlockLimits, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.maxQuantumsInsuranceLost.isZero()) {
      writer.uint32(8).uint64(message.maxQuantumsInsuranceLost);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubaccountBlockLimits {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubaccountBlockLimits();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.maxQuantumsInsuranceLost = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<SubaccountBlockLimits>): SubaccountBlockLimits {
    const message = createBaseSubaccountBlockLimits();
    message.maxQuantumsInsuranceLost = object.maxQuantumsInsuranceLost !== undefined && object.maxQuantumsInsuranceLost !== null ? Long.fromValue(object.maxQuantumsInsuranceLost) : Long.UZERO;
    return message;
  }

};