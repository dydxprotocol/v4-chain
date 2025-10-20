import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** StakingTier defines all staking levels of a fee tier */

export interface StakingTier {
  /** Which fee tier this corresponds to */
  feeTierName: string;
  /**
   * Staking levels for this tier (in strictly increasing order of staking
   * requirement)
   */

  levels: StakingLevel[];
}
/** StakingTier defines all staking levels of a fee tier */

export interface StakingTierSDKType {
  /** Which fee tier this corresponds to */
  fee_tier_name: string;
  /**
   * Staking levels for this tier (in strictly increasing order of staking
   * requirement)
   */

  levels: StakingLevelSDKType[];
}
/**
 * StakingLevel defines how many staked native tokens qualifies for how much
 * discount
 */

export interface StakingLevel {
  /** Minimum native tokens to stake (in base unit) */
  minStakedBaseTokens: Uint8Array;
  /** Fee discount in ppm (e.g. 1_000_000 is 100% discount) */

  feeDiscountPpm: number;
}
/**
 * StakingLevel defines how many staked native tokens qualifies for how much
 * discount
 */

export interface StakingLevelSDKType {
  /** Minimum native tokens to stake (in base unit) */
  min_staked_base_tokens: Uint8Array;
  /** Fee discount in ppm (e.g. 1_000_000 is 100% discount) */

  fee_discount_ppm: number;
}

function createBaseStakingTier(): StakingTier {
  return {
    feeTierName: "",
    levels: []
  };
}

export const StakingTier = {
  encode(message: StakingTier, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.feeTierName !== "") {
      writer.uint32(10).string(message.feeTierName);
    }

    for (const v of message.levels) {
      StakingLevel.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StakingTier {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStakingTier();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.feeTierName = reader.string();
          break;

        case 2:
          message.levels.push(StakingLevel.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StakingTier>): StakingTier {
    const message = createBaseStakingTier();
    message.feeTierName = object.feeTierName ?? "";
    message.levels = object.levels?.map(e => StakingLevel.fromPartial(e)) || [];
    return message;
  }

};

function createBaseStakingLevel(): StakingLevel {
  return {
    minStakedBaseTokens: new Uint8Array(),
    feeDiscountPpm: 0
  };
}

export const StakingLevel = {
  encode(message: StakingLevel, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.minStakedBaseTokens.length !== 0) {
      writer.uint32(10).bytes(message.minStakedBaseTokens);
    }

    if (message.feeDiscountPpm !== 0) {
      writer.uint32(16).uint32(message.feeDiscountPpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StakingLevel {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStakingLevel();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.minStakedBaseTokens = reader.bytes();
          break;

        case 2:
          message.feeDiscountPpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StakingLevel>): StakingLevel {
    const message = createBaseStakingLevel();
    message.minStakedBaseTokens = object.minStakedBaseTokens ?? new Uint8Array();
    message.feeDiscountPpm = object.feeDiscountPpm ?? 0;
    return message;
  }

};