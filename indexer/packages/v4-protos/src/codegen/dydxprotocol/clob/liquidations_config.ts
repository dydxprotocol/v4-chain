import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** LiquidationsConfig stores all configurable fields related to liquidations. */

export interface LiquidationsConfig {
  /**
   * The maximum liquidation fee (in parts-per-million). This fee goes
   * 100% to the insurance fund.
   */
  insuranceFundFeePpm: number;
  /** The fraction of the remaining collateral taken as a validator fee. */

  validatorFeePpm: number;
  /** The fraction of the remaining collateral taken as a liquidity fee. */

  liquidityFeePpm: number;
  /**
   * Limits around how many quote quantums from a single subaccount can
   * be liquidated within a single block.
   */

  subaccountBlockLimits?: SubaccountBlockLimits;
  /**
   * Config about how the fillable-price spread from the oracle price
   * increases based on the adjusted bankruptcy rating of the subaccount.
   */

  fillablePriceConfig?: FillablePriceConfig;
}
/** LiquidationsConfig stores all configurable fields related to liquidations. */

export interface LiquidationsConfigSDKType {
  /**
   * The maximum liquidation fee (in parts-per-million). This fee goes
   * 100% to the insurance fund.
   */
  insurance_fund_fee_ppm: number;
  /** The fraction of the remaining collateral taken as a validator fee. */

  validator_fee_ppm: number;
  /** The fraction of the remaining collateral taken as a liquidity fee. */

  liquidity_fee_ppm: number;
  /**
   * Limits around how many quote quantums from a single subaccount can
   * be liquidated within a single block.
   */

  subaccount_block_limits?: SubaccountBlockLimitsSDKType;
  /**
   * Config about how the fillable-price spread from the oracle price
   * increases based on the adjusted bankruptcy rating of the subaccount.
   */

  fillable_price_config?: FillablePriceConfigSDKType;
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
/**
 * FillablePriceConfig stores all configurable fields related to calculating
 * the fillable price for liquidating a position.
 */

export interface FillablePriceConfig {
  /** The rate at which the Adjusted Bankruptcy Rating increases. */
  bankruptcyAdjustmentPpm: number;
  /**
   * The maximum value that the liquidation spread can take, as
   * a ratio against the position's maintenance margin.
   */

  spreadToMaintenanceMarginRatioPpm: number;
}
/**
 * FillablePriceConfig stores all configurable fields related to calculating
 * the fillable price for liquidating a position.
 */

export interface FillablePriceConfigSDKType {
  /** The rate at which the Adjusted Bankruptcy Rating increases. */
  bankruptcy_adjustment_ppm: number;
  /**
   * The maximum value that the liquidation spread can take, as
   * a ratio against the position's maintenance margin.
   */

  spread_to_maintenance_margin_ratio_ppm: number;
}

function createBaseLiquidationsConfig(): LiquidationsConfig {
  return {
    insuranceFundFeePpm: 0,
    validatorFeePpm: 0,
    liquidityFeePpm: 0,
    subaccountBlockLimits: undefined,
    fillablePriceConfig: undefined
  };
}

export const LiquidationsConfig = {
  encode(message: LiquidationsConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.insuranceFundFeePpm !== 0) {
      writer.uint32(8).uint32(message.insuranceFundFeePpm);
    }

    if (message.validatorFeePpm !== 0) {
      writer.uint32(16).uint32(message.validatorFeePpm);
    }

    if (message.liquidityFeePpm !== 0) {
      writer.uint32(24).uint32(message.liquidityFeePpm);
    }

    if (message.subaccountBlockLimits !== undefined) {
      SubaccountBlockLimits.encode(message.subaccountBlockLimits, writer.uint32(34).fork()).ldelim();
    }

    if (message.fillablePriceConfig !== undefined) {
      FillablePriceConfig.encode(message.fillablePriceConfig, writer.uint32(42).fork()).ldelim();
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
          message.insuranceFundFeePpm = reader.uint32();
          break;

        case 2:
          message.validatorFeePpm = reader.uint32();
          break;

        case 3:
          message.liquidityFeePpm = reader.uint32();
          break;

        case 4:
          message.subaccountBlockLimits = SubaccountBlockLimits.decode(reader, reader.uint32());
          break;

        case 5:
          message.fillablePriceConfig = FillablePriceConfig.decode(reader, reader.uint32());
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
    message.insuranceFundFeePpm = object.insuranceFundFeePpm ?? 0;
    message.validatorFeePpm = object.validatorFeePpm ?? 0;
    message.liquidityFeePpm = object.liquidityFeePpm ?? 0;
    message.subaccountBlockLimits = object.subaccountBlockLimits !== undefined && object.subaccountBlockLimits !== null ? SubaccountBlockLimits.fromPartial(object.subaccountBlockLimits) : undefined;
    message.fillablePriceConfig = object.fillablePriceConfig !== undefined && object.fillablePriceConfig !== null ? FillablePriceConfig.fromPartial(object.fillablePriceConfig) : undefined;
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

function createBaseFillablePriceConfig(): FillablePriceConfig {
  return {
    bankruptcyAdjustmentPpm: 0,
    spreadToMaintenanceMarginRatioPpm: 0
  };
}

export const FillablePriceConfig = {
  encode(message: FillablePriceConfig, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.bankruptcyAdjustmentPpm !== 0) {
      writer.uint32(8).uint32(message.bankruptcyAdjustmentPpm);
    }

    if (message.spreadToMaintenanceMarginRatioPpm !== 0) {
      writer.uint32(16).uint32(message.spreadToMaintenanceMarginRatioPpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FillablePriceConfig {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFillablePriceConfig();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.bankruptcyAdjustmentPpm = reader.uint32();
          break;

        case 2:
          message.spreadToMaintenanceMarginRatioPpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<FillablePriceConfig>): FillablePriceConfig {
    const message = createBaseFillablePriceConfig();
    message.bankruptcyAdjustmentPpm = object.bankruptcyAdjustmentPpm ?? 0;
    message.spreadToMaintenanceMarginRatioPpm = object.spreadToMaintenanceMarginRatioPpm ?? 0;
    return message;
  }

};