import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** PerpetualFeeParams defines the parameters for perpetual fees. */

export interface PerpetualFeeParams {
  /** Sorted fee tiers (lowest requirements first). */
  tiers: PerpetualFeeTier[];
}
/** PerpetualFeeParams defines the parameters for perpetual fees. */

export interface PerpetualFeeParamsSDKType {
  /** Sorted fee tiers (lowest requirements first). */
  tiers: PerpetualFeeTierSDKType[];
}
/** A fee tier for perpetuals */

export interface PerpetualFeeTier {
  /** Human-readable name of the tier, e.g. "Gold". */
  name: string;
  /** The trader's absolute volume requirement in quote quantums. */

  absoluteVolumeRequirement: Long;
  /** The total volume share requirement. */

  totalVolumeShareRequirementPpm: number;
  /** The maker volume share requirement. */

  makerVolumeShareRequirementPpm: number;
  /** The maker fee once this tier is reached. */

  makerFeePpm: number;
  /** The taker fee once this tier is reached. */

  takerFeePpm: number;
}
/** A fee tier for perpetuals */

export interface PerpetualFeeTierSDKType {
  /** Human-readable name of the tier, e.g. "Gold". */
  name: string;
  /** The trader's absolute volume requirement in quote quantums. */

  absolute_volume_requirement: Long;
  /** The total volume share requirement. */

  total_volume_share_requirement_ppm: number;
  /** The maker volume share requirement. */

  maker_volume_share_requirement_ppm: number;
  /** The maker fee once this tier is reached. */

  maker_fee_ppm: number;
  /** The taker fee once this tier is reached. */

  taker_fee_ppm: number;
}

function createBasePerpetualFeeParams(): PerpetualFeeParams {
  return {
    tiers: []
  };
}

export const PerpetualFeeParams = {
  encode(message: PerpetualFeeParams, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.tiers) {
      PerpetualFeeTier.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PerpetualFeeParams {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePerpetualFeeParams();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.tiers.push(PerpetualFeeTier.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<PerpetualFeeParams>): PerpetualFeeParams {
    const message = createBasePerpetualFeeParams();
    message.tiers = object.tiers?.map(e => PerpetualFeeTier.fromPartial(e)) || [];
    return message;
  }

};

function createBasePerpetualFeeTier(): PerpetualFeeTier {
  return {
    name: "",
    absoluteVolumeRequirement: Long.UZERO,
    totalVolumeShareRequirementPpm: 0,
    makerVolumeShareRequirementPpm: 0,
    makerFeePpm: 0,
    takerFeePpm: 0
  };
}

export const PerpetualFeeTier = {
  encode(message: PerpetualFeeTier, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }

    if (!message.absoluteVolumeRequirement.isZero()) {
      writer.uint32(16).uint64(message.absoluteVolumeRequirement);
    }

    if (message.totalVolumeShareRequirementPpm !== 0) {
      writer.uint32(24).uint32(message.totalVolumeShareRequirementPpm);
    }

    if (message.makerVolumeShareRequirementPpm !== 0) {
      writer.uint32(32).uint32(message.makerVolumeShareRequirementPpm);
    }

    if (message.makerFeePpm !== 0) {
      writer.uint32(40).sint32(message.makerFeePpm);
    }

    if (message.takerFeePpm !== 0) {
      writer.uint32(48).sint32(message.takerFeePpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PerpetualFeeTier {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePerpetualFeeTier();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;

        case 2:
          message.absoluteVolumeRequirement = (reader.uint64() as Long);
          break;

        case 3:
          message.totalVolumeShareRequirementPpm = reader.uint32();
          break;

        case 4:
          message.makerVolumeShareRequirementPpm = reader.uint32();
          break;

        case 5:
          message.makerFeePpm = reader.sint32();
          break;

        case 6:
          message.takerFeePpm = reader.sint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<PerpetualFeeTier>): PerpetualFeeTier {
    const message = createBasePerpetualFeeTier();
    message.name = object.name ?? "";
    message.absoluteVolumeRequirement = object.absoluteVolumeRequirement !== undefined && object.absoluteVolumeRequirement !== null ? Long.fromValue(object.absoluteVolumeRequirement) : Long.UZERO;
    message.totalVolumeShareRequirementPpm = object.totalVolumeShareRequirementPpm ?? 0;
    message.makerVolumeShareRequirementPpm = object.makerVolumeShareRequirementPpm ?? 0;
    message.makerFeePpm = object.makerFeePpm ?? 0;
    message.takerFeePpm = object.takerFeePpm ?? 0;
    return message;
  }

};