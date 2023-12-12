import { BinaryReader, BinaryWriter } from "../../binary";
/** Params defines the parameters for x/perpetuals module. */
export interface Params {
  /**
   * Funding rate clamp factor in parts-per-million, used for clamping 8-hour
   * funding rates according to equation: |R| <= funding_rate_clamp_factor *
   * (initial margin - maintenance margin).
   */
  fundingRateClampFactorPpm: number;
  /**
   * Premium vote clamp factor in parts-per-million, used for clamping premium
   * votes according to equation: |V| <= premium_vote_clamp_factor *
   * (initial margin - maintenance margin).
   */
  premiumVoteClampFactorPpm: number;
  /**
   * Minimum number of premium votes per premium sample. If number of premium
   * votes is smaller than this number, pad with zeros up to this number.
   */
  minNumVotesPerSample: number;
}
export interface ParamsProtoMsg {
  typeUrl: "/dydxprotocol.perpetuals.Params";
  value: Uint8Array;
}
/** Params defines the parameters for x/perpetuals module. */
export interface ParamsAmino {
  /**
   * Funding rate clamp factor in parts-per-million, used for clamping 8-hour
   * funding rates according to equation: |R| <= funding_rate_clamp_factor *
   * (initial margin - maintenance margin).
   */
  funding_rate_clamp_factor_ppm?: number;
  /**
   * Premium vote clamp factor in parts-per-million, used for clamping premium
   * votes according to equation: |V| <= premium_vote_clamp_factor *
   * (initial margin - maintenance margin).
   */
  premium_vote_clamp_factor_ppm?: number;
  /**
   * Minimum number of premium votes per premium sample. If number of premium
   * votes is smaller than this number, pad with zeros up to this number.
   */
  min_num_votes_per_sample?: number;
}
export interface ParamsAminoMsg {
  type: "/dydxprotocol.perpetuals.Params";
  value: ParamsAmino;
}
/** Params defines the parameters for x/perpetuals module. */
export interface ParamsSDKType {
  funding_rate_clamp_factor_ppm: number;
  premium_vote_clamp_factor_ppm: number;
  min_num_votes_per_sample: number;
}
function createBaseParams(): Params {
  return {
    fundingRateClampFactorPpm: 0,
    premiumVoteClampFactorPpm: 0,
    minNumVotesPerSample: 0
  };
}
export const Params = {
  typeUrl: "/dydxprotocol.perpetuals.Params",
  encode(message: Params, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.fundingRateClampFactorPpm !== 0) {
      writer.uint32(8).uint32(message.fundingRateClampFactorPpm);
    }
    if (message.premiumVoteClampFactorPpm !== 0) {
      writer.uint32(16).uint32(message.premiumVoteClampFactorPpm);
    }
    if (message.minNumVotesPerSample !== 0) {
      writer.uint32(24).uint32(message.minNumVotesPerSample);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): Params {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fundingRateClampFactorPpm = reader.uint32();
          break;
        case 2:
          message.premiumVoteClampFactorPpm = reader.uint32();
          break;
        case 3:
          message.minNumVotesPerSample = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<Params>): Params {
    const message = createBaseParams();
    message.fundingRateClampFactorPpm = object.fundingRateClampFactorPpm ?? 0;
    message.premiumVoteClampFactorPpm = object.premiumVoteClampFactorPpm ?? 0;
    message.minNumVotesPerSample = object.minNumVotesPerSample ?? 0;
    return message;
  },
  fromAmino(object: ParamsAmino): Params {
    const message = createBaseParams();
    if (object.funding_rate_clamp_factor_ppm !== undefined && object.funding_rate_clamp_factor_ppm !== null) {
      message.fundingRateClampFactorPpm = object.funding_rate_clamp_factor_ppm;
    }
    if (object.premium_vote_clamp_factor_ppm !== undefined && object.premium_vote_clamp_factor_ppm !== null) {
      message.premiumVoteClampFactorPpm = object.premium_vote_clamp_factor_ppm;
    }
    if (object.min_num_votes_per_sample !== undefined && object.min_num_votes_per_sample !== null) {
      message.minNumVotesPerSample = object.min_num_votes_per_sample;
    }
    return message;
  },
  toAmino(message: Params): ParamsAmino {
    const obj: any = {};
    obj.funding_rate_clamp_factor_ppm = message.fundingRateClampFactorPpm;
    obj.premium_vote_clamp_factor_ppm = message.premiumVoteClampFactorPpm;
    obj.min_num_votes_per_sample = message.minNumVotesPerSample;
    return obj;
  },
  fromAminoMsg(object: ParamsAminoMsg): Params {
    return Params.fromAmino(object.value);
  },
  fromProtoMsg(message: ParamsProtoMsg): Params {
    return Params.decode(message.value);
  },
  toProto(message: Params): Uint8Array {
    return Params.encode(message).finish();
  },
  toProtoMsg(message: Params): ParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.perpetuals.Params",
      value: Params.encode(message).finish()
    };
  }
};