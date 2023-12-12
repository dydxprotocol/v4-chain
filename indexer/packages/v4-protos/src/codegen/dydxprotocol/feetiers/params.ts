import { BinaryReader, BinaryWriter } from "../../binary";
/** PerpetualFeeParams defines the parameters for perpetual fees. */
export interface PerpetualFeeParams {
  /** Sorted fee tiers (lowest requirements first). */
  tiers: PerpetualFeeTier[];
}
export interface PerpetualFeeParamsProtoMsg {
  typeUrl: "/dydxprotocol.feetiers.PerpetualFeeParams";
  value: Uint8Array;
}
/** PerpetualFeeParams defines the parameters for perpetual fees. */
export interface PerpetualFeeParamsAmino {
  /** Sorted fee tiers (lowest requirements first). */
  tiers?: PerpetualFeeTierAmino[];
}
export interface PerpetualFeeParamsAminoMsg {
  type: "/dydxprotocol.feetiers.PerpetualFeeParams";
  value: PerpetualFeeParamsAmino;
}
/** PerpetualFeeParams defines the parameters for perpetual fees. */
export interface PerpetualFeeParamsSDKType {
  tiers: PerpetualFeeTierSDKType[];
}
/** A fee tier for perpetuals */
export interface PerpetualFeeTier {
  /** Human-readable name of the tier, e.g. "Gold". */
  name: string;
  /** The trader's absolute volume requirement in quote quantums. */
  absoluteVolumeRequirement: bigint;
  /** The total volume share requirement. */
  totalVolumeShareRequirementPpm: number;
  /** The maker volume share requirement. */
  makerVolumeShareRequirementPpm: number;
  /** The maker fee once this tier is reached. */
  makerFeePpm: number;
  /** The taker fee once this tier is reached. */
  takerFeePpm: number;
}
export interface PerpetualFeeTierProtoMsg {
  typeUrl: "/dydxprotocol.feetiers.PerpetualFeeTier";
  value: Uint8Array;
}
/** A fee tier for perpetuals */
export interface PerpetualFeeTierAmino {
  /** Human-readable name of the tier, e.g. "Gold". */
  name?: string;
  /** The trader's absolute volume requirement in quote quantums. */
  absolute_volume_requirement?: string;
  /** The total volume share requirement. */
  total_volume_share_requirement_ppm?: number;
  /** The maker volume share requirement. */
  maker_volume_share_requirement_ppm?: number;
  /** The maker fee once this tier is reached. */
  maker_fee_ppm?: number;
  /** The taker fee once this tier is reached. */
  taker_fee_ppm?: number;
}
export interface PerpetualFeeTierAminoMsg {
  type: "/dydxprotocol.feetiers.PerpetualFeeTier";
  value: PerpetualFeeTierAmino;
}
/** A fee tier for perpetuals */
export interface PerpetualFeeTierSDKType {
  name: string;
  absolute_volume_requirement: bigint;
  total_volume_share_requirement_ppm: number;
  maker_volume_share_requirement_ppm: number;
  maker_fee_ppm: number;
  taker_fee_ppm: number;
}
function createBasePerpetualFeeParams(): PerpetualFeeParams {
  return {
    tiers: []
  };
}
export const PerpetualFeeParams = {
  typeUrl: "/dydxprotocol.feetiers.PerpetualFeeParams",
  encode(message: PerpetualFeeParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.tiers) {
      PerpetualFeeTier.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): PerpetualFeeParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<PerpetualFeeParams>): PerpetualFeeParams {
    const message = createBasePerpetualFeeParams();
    message.tiers = object.tiers?.map(e => PerpetualFeeTier.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: PerpetualFeeParamsAmino): PerpetualFeeParams {
    const message = createBasePerpetualFeeParams();
    message.tiers = object.tiers?.map(e => PerpetualFeeTier.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: PerpetualFeeParams): PerpetualFeeParamsAmino {
    const obj: any = {};
    if (message.tiers) {
      obj.tiers = message.tiers.map(e => e ? PerpetualFeeTier.toAmino(e) : undefined);
    } else {
      obj.tiers = [];
    }
    return obj;
  },
  fromAminoMsg(object: PerpetualFeeParamsAminoMsg): PerpetualFeeParams {
    return PerpetualFeeParams.fromAmino(object.value);
  },
  fromProtoMsg(message: PerpetualFeeParamsProtoMsg): PerpetualFeeParams {
    return PerpetualFeeParams.decode(message.value);
  },
  toProto(message: PerpetualFeeParams): Uint8Array {
    return PerpetualFeeParams.encode(message).finish();
  },
  toProtoMsg(message: PerpetualFeeParams): PerpetualFeeParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.feetiers.PerpetualFeeParams",
      value: PerpetualFeeParams.encode(message).finish()
    };
  }
};
function createBasePerpetualFeeTier(): PerpetualFeeTier {
  return {
    name: "",
    absoluteVolumeRequirement: BigInt(0),
    totalVolumeShareRequirementPpm: 0,
    makerVolumeShareRequirementPpm: 0,
    makerFeePpm: 0,
    takerFeePpm: 0
  };
}
export const PerpetualFeeTier = {
  typeUrl: "/dydxprotocol.feetiers.PerpetualFeeTier",
  encode(message: PerpetualFeeTier, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.absoluteVolumeRequirement !== BigInt(0)) {
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
  decode(input: BinaryReader | Uint8Array, length?: number): PerpetualFeeTier {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePerpetualFeeTier();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.absoluteVolumeRequirement = reader.uint64();
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
  fromPartial(object: Partial<PerpetualFeeTier>): PerpetualFeeTier {
    const message = createBasePerpetualFeeTier();
    message.name = object.name ?? "";
    message.absoluteVolumeRequirement = object.absoluteVolumeRequirement !== undefined && object.absoluteVolumeRequirement !== null ? BigInt(object.absoluteVolumeRequirement.toString()) : BigInt(0);
    message.totalVolumeShareRequirementPpm = object.totalVolumeShareRequirementPpm ?? 0;
    message.makerVolumeShareRequirementPpm = object.makerVolumeShareRequirementPpm ?? 0;
    message.makerFeePpm = object.makerFeePpm ?? 0;
    message.takerFeePpm = object.takerFeePpm ?? 0;
    return message;
  },
  fromAmino(object: PerpetualFeeTierAmino): PerpetualFeeTier {
    const message = createBasePerpetualFeeTier();
    if (object.name !== undefined && object.name !== null) {
      message.name = object.name;
    }
    if (object.absolute_volume_requirement !== undefined && object.absolute_volume_requirement !== null) {
      message.absoluteVolumeRequirement = BigInt(object.absolute_volume_requirement);
    }
    if (object.total_volume_share_requirement_ppm !== undefined && object.total_volume_share_requirement_ppm !== null) {
      message.totalVolumeShareRequirementPpm = object.total_volume_share_requirement_ppm;
    }
    if (object.maker_volume_share_requirement_ppm !== undefined && object.maker_volume_share_requirement_ppm !== null) {
      message.makerVolumeShareRequirementPpm = object.maker_volume_share_requirement_ppm;
    }
    if (object.maker_fee_ppm !== undefined && object.maker_fee_ppm !== null) {
      message.makerFeePpm = object.maker_fee_ppm;
    }
    if (object.taker_fee_ppm !== undefined && object.taker_fee_ppm !== null) {
      message.takerFeePpm = object.taker_fee_ppm;
    }
    return message;
  },
  toAmino(message: PerpetualFeeTier): PerpetualFeeTierAmino {
    const obj: any = {};
    obj.name = message.name;
    obj.absolute_volume_requirement = message.absoluteVolumeRequirement ? message.absoluteVolumeRequirement.toString() : undefined;
    obj.total_volume_share_requirement_ppm = message.totalVolumeShareRequirementPpm;
    obj.maker_volume_share_requirement_ppm = message.makerVolumeShareRequirementPpm;
    obj.maker_fee_ppm = message.makerFeePpm;
    obj.taker_fee_ppm = message.takerFeePpm;
    return obj;
  },
  fromAminoMsg(object: PerpetualFeeTierAminoMsg): PerpetualFeeTier {
    return PerpetualFeeTier.fromAmino(object.value);
  },
  fromProtoMsg(message: PerpetualFeeTierProtoMsg): PerpetualFeeTier {
    return PerpetualFeeTier.decode(message.value);
  },
  toProto(message: PerpetualFeeTier): Uint8Array {
    return PerpetualFeeTier.encode(message).finish();
  },
  toProtoMsg(message: PerpetualFeeTier): PerpetualFeeTierProtoMsg {
    return {
      typeUrl: "/dydxprotocol.feetiers.PerpetualFeeTier",
      value: PerpetualFeeTier.encode(message).finish()
    };
  }
};