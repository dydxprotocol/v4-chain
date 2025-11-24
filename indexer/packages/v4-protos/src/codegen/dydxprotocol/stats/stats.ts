import { Timestamp } from "../../google/protobuf/timestamp";
import * as _m0 from "protobufjs/minimal";
import { Long, DeepPartial, toTimestamp, fromTimestamp } from "../../helpers";
/** Role indicates whether this attribution is for the taker or maker */

export enum AffiliateAttribution_Role {
  ROLE_UNSPECIFIED = 0,
  ROLE_TAKER = 1,
  ROLE_MAKER = 2,
  UNRECOGNIZED = -1,
}
/** Role indicates whether this attribution is for the taker or maker */

export enum AffiliateAttribution_RoleSDKType {
  ROLE_UNSPECIFIED = 0,
  ROLE_TAKER = 1,
  ROLE_MAKER = 2,
  UNRECOGNIZED = -1,
}
export function affiliateAttribution_RoleFromJSON(object: any): AffiliateAttribution_Role {
  switch (object) {
    case 0:
    case "ROLE_UNSPECIFIED":
      return AffiliateAttribution_Role.ROLE_UNSPECIFIED;

    case 1:
    case "ROLE_TAKER":
      return AffiliateAttribution_Role.ROLE_TAKER;

    case 2:
    case "ROLE_MAKER":
      return AffiliateAttribution_Role.ROLE_MAKER;

    case -1:
    case "UNRECOGNIZED":
    default:
      return AffiliateAttribution_Role.UNRECOGNIZED;
  }
}
export function affiliateAttribution_RoleToJSON(object: AffiliateAttribution_Role): string {
  switch (object) {
    case AffiliateAttribution_Role.ROLE_UNSPECIFIED:
      return "ROLE_UNSPECIFIED";

    case AffiliateAttribution_Role.ROLE_TAKER:
      return "ROLE_TAKER";

    case AffiliateAttribution_Role.ROLE_MAKER:
      return "ROLE_MAKER";

    case AffiliateAttribution_Role.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/** AffiliateAttribution represents the affiliate attribution for a fill. */

export interface AffiliateAttribution {
  /** Role of the trader (taker or maker) whose affiliate is being attributed */
  role: AffiliateAttribution_Role;
  /** Referrer address (the affiliate receiving the fee) */

  referrerAddress: string;
  /** Referred volume in quote quantums (capped based on 30-day volume limits) */

  referredVolumeQuoteQuantums: Long;
}
/** AffiliateAttribution represents the affiliate attribution for a fill. */

export interface AffiliateAttributionSDKType {
  /** Role of the trader (taker or maker) whose affiliate is being attributed */
  role: AffiliateAttribution_RoleSDKType;
  /** Referrer address (the affiliate receiving the fee) */

  referrer_address: string;
  /** Referred volume in quote quantums (capped based on 30-day volume limits) */

  referred_volume_quote_quantums: Long;
}
/** BlockStats is used to store stats transiently within the scope of a block. */

export interface BlockStats {
  /** The fills that occured on this block. */
  fills: BlockStats_Fill[];
}
/** BlockStats is used to store stats transiently within the scope of a block. */

export interface BlockStatsSDKType {
  /** The fills that occured on this block. */
  fills: BlockStats_FillSDKType[];
}
/** Fill records data about a fill on this block. */

export interface BlockStats_Fill {
  /** Taker wallet address */
  taker: string;
  /** Maker wallet address */

  maker: string;
  /**
   * Notional USDC filled in quantums
   * Used to calculate fee tier, and affiliate revenue attributed for taker
   */

  notional: Long;
  /**
   * Affiliate fee generated in quantums of the taker fee for the affiliate
   * Used to calculate affiliate revenue attributed for taker. This is dynamic
   * per affiliate tier
   */

  affiliateFeeGeneratedQuantums: Long;
  /**
   * Affiliate revenue attributions for this fill (can include both taker and
   * maker)
   */

  affiliateAttributions: AffiliateAttribution[];
}
/** Fill records data about a fill on this block. */

export interface BlockStats_FillSDKType {
  /** Taker wallet address */
  taker: string;
  /** Maker wallet address */

  maker: string;
  /**
   * Notional USDC filled in quantums
   * Used to calculate fee tier, and affiliate revenue attributed for taker
   */

  notional: Long;
  /**
   * Affiliate fee generated in quantums of the taker fee for the affiliate
   * Used to calculate affiliate revenue attributed for taker. This is dynamic
   * per affiliate tier
   */

  affiliate_fee_generated_quantums: Long;
  /**
   * Affiliate revenue attributions for this fill (can include both taker and
   * maker)
   */

  affiliate_attributions: AffiliateAttributionSDKType[];
}
/** StatsMetadata stores metadata for the x/stats module */

export interface StatsMetadata {
  /**
   * The oldest epoch that is included in the stats. The next epoch to be
   * removed from the window.
   */
  trailingEpoch: number;
}
/** StatsMetadata stores metadata for the x/stats module */

export interface StatsMetadataSDKType {
  /**
   * The oldest epoch that is included in the stats. The next epoch to be
   * removed from the window.
   */
  trailing_epoch: number;
}
/** EpochStats stores stats for a particular epoch */

export interface EpochStats {
  /** Epoch end time */
  epochEndTime?: Date;
  /** Stats for each user in this epoch. Sorted by user. */

  stats: EpochStats_UserWithStats[];
}
/** EpochStats stores stats for a particular epoch */

export interface EpochStatsSDKType {
  /** Epoch end time */
  epoch_end_time?: Date;
  /** Stats for each user in this epoch. Sorted by user. */

  stats: EpochStats_UserWithStatsSDKType[];
}
/** A user and its associated stats */

export interface EpochStats_UserWithStats {
  user: string;
  stats?: UserStats;
}
/** A user and its associated stats */

export interface EpochStats_UserWithStatsSDKType {
  user: string;
  stats?: UserStatsSDKType;
}
/** GlobalStats stores global stats for the rolling window (default 30d). */

export interface GlobalStats {
  /** Notional USDC traded in quantums */
  notionalTraded: Long;
}
/** GlobalStats stores global stats for the rolling window (default 30d). */

export interface GlobalStatsSDKType {
  /** Notional USDC traded in quantums */
  notional_traded: Long;
}
/**
 * UserStats stores stats for a User. This is the sum of all stats for a user in
 * the rolling window (default 30d).
 */

export interface UserStats {
  /** Taker USDC in quantums */
  takerNotional: Long;
  /** Maker USDC in quantums */

  makerNotional: Long;
  /** Affiliate revenue generated in quantums with this user being a referee */

  affiliate_30dRevenueGeneratedQuantums: Long;
  /** Referred volume in quote quantums with this user being an affiliate */

  affiliate_30dReferredVolumeQuoteQuantums: Long;
  /**
   * Attributed volume in quote quantums - volume from this user (as referee)
   * that has been attributed to their affiliate in the last 30 days
   */

  affiliate_30dAttributedVolumeQuoteQuantums: Long;
}
/**
 * UserStats stores stats for a User. This is the sum of all stats for a user in
 * the rolling window (default 30d).
 */

export interface UserStatsSDKType {
  /** Taker USDC in quantums */
  taker_notional: Long;
  /** Maker USDC in quantums */

  maker_notional: Long;
  /** Affiliate revenue generated in quantums with this user being a referee */

  affiliate_30d_revenue_generated_quantums: Long;
  /** Referred volume in quote quantums with this user being an affiliate */

  affiliate_30d_referred_volume_quote_quantums: Long;
  /**
   * Attributed volume in quote quantums - volume from this user (as referee)
   * that has been attributed to their affiliate in the last 30 days
   */

  affiliate_30d_attributed_volume_quote_quantums: Long;
}
/** CachedStakedBaseTokens stores the last calculated total staked base tokens */

export interface CachedStakedBaseTokens {
  /** Last calculated total staked base tokens by the delegator. */
  stakedBaseTokens: Uint8Array;
  /**
   * Block time at which the calculation is cached (in Unix Epoch seconds)
   * Rounded down to nearest second.
   */

  cachedAt: Long;
}
/** CachedStakedBaseTokens stores the last calculated total staked base tokens */

export interface CachedStakedBaseTokensSDKType {
  /** Last calculated total staked base tokens by the delegator. */
  staked_base_tokens: Uint8Array;
  /**
   * Block time at which the calculation is cached (in Unix Epoch seconds)
   * Rounded down to nearest second.
   */

  cached_at: Long;
}

function createBaseAffiliateAttribution(): AffiliateAttribution {
  return {
    role: 0,
    referrerAddress: "",
    referredVolumeQuoteQuantums: Long.UZERO
  };
}

export const AffiliateAttribution = {
  encode(message: AffiliateAttribution, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.role !== 0) {
      writer.uint32(8).int32(message.role);
    }

    if (message.referrerAddress !== "") {
      writer.uint32(18).string(message.referrerAddress);
    }

    if (!message.referredVolumeQuoteQuantums.isZero()) {
      writer.uint32(24).uint64(message.referredVolumeQuoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AffiliateAttribution {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAffiliateAttribution();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.role = (reader.int32() as any);
          break;

        case 2:
          message.referrerAddress = reader.string();
          break;

        case 3:
          message.referredVolumeQuoteQuantums = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AffiliateAttribution>): AffiliateAttribution {
    const message = createBaseAffiliateAttribution();
    message.role = object.role ?? 0;
    message.referrerAddress = object.referrerAddress ?? "";
    message.referredVolumeQuoteQuantums = object.referredVolumeQuoteQuantums !== undefined && object.referredVolumeQuoteQuantums !== null ? Long.fromValue(object.referredVolumeQuoteQuantums) : Long.UZERO;
    return message;
  }

};

function createBaseBlockStats(): BlockStats {
  return {
    fills: []
  };
}

export const BlockStats = {
  encode(message: BlockStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.fills) {
      BlockStats_Fill.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BlockStats {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockStats();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.fills.push(BlockStats_Fill.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<BlockStats>): BlockStats {
    const message = createBaseBlockStats();
    message.fills = object.fills?.map(e => BlockStats_Fill.fromPartial(e)) || [];
    return message;
  }

};

function createBaseBlockStats_Fill(): BlockStats_Fill {
  return {
    taker: "",
    maker: "",
    notional: Long.UZERO,
    affiliateFeeGeneratedQuantums: Long.UZERO,
    affiliateAttributions: []
  };
}

export const BlockStats_Fill = {
  encode(message: BlockStats_Fill, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.taker !== "") {
      writer.uint32(10).string(message.taker);
    }

    if (message.maker !== "") {
      writer.uint32(18).string(message.maker);
    }

    if (!message.notional.isZero()) {
      writer.uint32(24).uint64(message.notional);
    }

    if (!message.affiliateFeeGeneratedQuantums.isZero()) {
      writer.uint32(32).uint64(message.affiliateFeeGeneratedQuantums);
    }

    for (const v of message.affiliateAttributions) {
      AffiliateAttribution.encode(v!, writer.uint32(42).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BlockStats_Fill {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockStats_Fill();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.taker = reader.string();
          break;

        case 2:
          message.maker = reader.string();
          break;

        case 3:
          message.notional = (reader.uint64() as Long);
          break;

        case 4:
          message.affiliateFeeGeneratedQuantums = (reader.uint64() as Long);
          break;

        case 5:
          message.affiliateAttributions.push(AffiliateAttribution.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<BlockStats_Fill>): BlockStats_Fill {
    const message = createBaseBlockStats_Fill();
    message.taker = object.taker ?? "";
    message.maker = object.maker ?? "";
    message.notional = object.notional !== undefined && object.notional !== null ? Long.fromValue(object.notional) : Long.UZERO;
    message.affiliateFeeGeneratedQuantums = object.affiliateFeeGeneratedQuantums !== undefined && object.affiliateFeeGeneratedQuantums !== null ? Long.fromValue(object.affiliateFeeGeneratedQuantums) : Long.UZERO;
    message.affiliateAttributions = object.affiliateAttributions?.map(e => AffiliateAttribution.fromPartial(e)) || [];
    return message;
  }

};

function createBaseStatsMetadata(): StatsMetadata {
  return {
    trailingEpoch: 0
  };
}

export const StatsMetadata = {
  encode(message: StatsMetadata, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.trailingEpoch !== 0) {
      writer.uint32(8).uint32(message.trailingEpoch);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StatsMetadata {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatsMetadata();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.trailingEpoch = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StatsMetadata>): StatsMetadata {
    const message = createBaseStatsMetadata();
    message.trailingEpoch = object.trailingEpoch ?? 0;
    return message;
  }

};

function createBaseEpochStats(): EpochStats {
  return {
    epochEndTime: undefined,
    stats: []
  };
}

export const EpochStats = {
  encode(message: EpochStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.epochEndTime !== undefined) {
      Timestamp.encode(toTimestamp(message.epochEndTime), writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.stats) {
      EpochStats_UserWithStats.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EpochStats {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEpochStats();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.epochEndTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;

        case 2:
          message.stats.push(EpochStats_UserWithStats.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<EpochStats>): EpochStats {
    const message = createBaseEpochStats();
    message.epochEndTime = object.epochEndTime ?? undefined;
    message.stats = object.stats?.map(e => EpochStats_UserWithStats.fromPartial(e)) || [];
    return message;
  }

};

function createBaseEpochStats_UserWithStats(): EpochStats_UserWithStats {
  return {
    user: "",
    stats: undefined
  };
}

export const EpochStats_UserWithStats = {
  encode(message: EpochStats_UserWithStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== "") {
      writer.uint32(10).string(message.user);
    }

    if (message.stats !== undefined) {
      UserStats.encode(message.stats, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EpochStats_UserWithStats {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEpochStats_UserWithStats();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.user = reader.string();
          break;

        case 2:
          message.stats = UserStats.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<EpochStats_UserWithStats>): EpochStats_UserWithStats {
    const message = createBaseEpochStats_UserWithStats();
    message.user = object.user ?? "";
    message.stats = object.stats !== undefined && object.stats !== null ? UserStats.fromPartial(object.stats) : undefined;
    return message;
  }

};

function createBaseGlobalStats(): GlobalStats {
  return {
    notionalTraded: Long.UZERO
  };
}

export const GlobalStats = {
  encode(message: GlobalStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.notionalTraded.isZero()) {
      writer.uint32(8).uint64(message.notionalTraded);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GlobalStats {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGlobalStats();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.notionalTraded = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<GlobalStats>): GlobalStats {
    const message = createBaseGlobalStats();
    message.notionalTraded = object.notionalTraded !== undefined && object.notionalTraded !== null ? Long.fromValue(object.notionalTraded) : Long.UZERO;
    return message;
  }

};

function createBaseUserStats(): UserStats {
  return {
    takerNotional: Long.UZERO,
    makerNotional: Long.UZERO,
    affiliate_30dRevenueGeneratedQuantums: Long.UZERO,
    affiliate_30dReferredVolumeQuoteQuantums: Long.UZERO,
    affiliate_30dAttributedVolumeQuoteQuantums: Long.UZERO
  };
}

export const UserStats = {
  encode(message: UserStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.takerNotional.isZero()) {
      writer.uint32(8).uint64(message.takerNotional);
    }

    if (!message.makerNotional.isZero()) {
      writer.uint32(16).uint64(message.makerNotional);
    }

    if (!message.affiliate_30dRevenueGeneratedQuantums.isZero()) {
      writer.uint32(24).uint64(message.affiliate_30dRevenueGeneratedQuantums);
    }

    if (!message.affiliate_30dReferredVolumeQuoteQuantums.isZero()) {
      writer.uint32(32).uint64(message.affiliate_30dReferredVolumeQuoteQuantums);
    }

    if (!message.affiliate_30dAttributedVolumeQuoteQuantums.isZero()) {
      writer.uint32(40).uint64(message.affiliate_30dAttributedVolumeQuoteQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserStats {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserStats();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.takerNotional = (reader.uint64() as Long);
          break;

        case 2:
          message.makerNotional = (reader.uint64() as Long);
          break;

        case 3:
          message.affiliate_30dRevenueGeneratedQuantums = (reader.uint64() as Long);
          break;

        case 4:
          message.affiliate_30dReferredVolumeQuoteQuantums = (reader.uint64() as Long);
          break;

        case 5:
          message.affiliate_30dAttributedVolumeQuoteQuantums = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<UserStats>): UserStats {
    const message = createBaseUserStats();
    message.takerNotional = object.takerNotional !== undefined && object.takerNotional !== null ? Long.fromValue(object.takerNotional) : Long.UZERO;
    message.makerNotional = object.makerNotional !== undefined && object.makerNotional !== null ? Long.fromValue(object.makerNotional) : Long.UZERO;
    message.affiliate_30dRevenueGeneratedQuantums = object.affiliate_30dRevenueGeneratedQuantums !== undefined && object.affiliate_30dRevenueGeneratedQuantums !== null ? Long.fromValue(object.affiliate_30dRevenueGeneratedQuantums) : Long.UZERO;
    message.affiliate_30dReferredVolumeQuoteQuantums = object.affiliate_30dReferredVolumeQuoteQuantums !== undefined && object.affiliate_30dReferredVolumeQuoteQuantums !== null ? Long.fromValue(object.affiliate_30dReferredVolumeQuoteQuantums) : Long.UZERO;
    message.affiliate_30dAttributedVolumeQuoteQuantums = object.affiliate_30dAttributedVolumeQuoteQuantums !== undefined && object.affiliate_30dAttributedVolumeQuoteQuantums !== null ? Long.fromValue(object.affiliate_30dAttributedVolumeQuoteQuantums) : Long.UZERO;
    return message;
  }

};

function createBaseCachedStakedBaseTokens(): CachedStakedBaseTokens {
  return {
    stakedBaseTokens: new Uint8Array(),
    cachedAt: Long.ZERO
  };
}

export const CachedStakedBaseTokens = {
  encode(message: CachedStakedBaseTokens, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stakedBaseTokens.length !== 0) {
      writer.uint32(10).bytes(message.stakedBaseTokens);
    }

    if (!message.cachedAt.isZero()) {
      writer.uint32(16).int64(message.cachedAt);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CachedStakedBaseTokens {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCachedStakedBaseTokens();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.stakedBaseTokens = reader.bytes();
          break;

        case 2:
          message.cachedAt = (reader.int64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<CachedStakedBaseTokens>): CachedStakedBaseTokens {
    const message = createBaseCachedStakedBaseTokens();
    message.stakedBaseTokens = object.stakedBaseTokens ?? new Uint8Array();
    message.cachedAt = object.cachedAt !== undefined && object.cachedAt !== null ? Long.fromValue(object.cachedAt) : Long.ZERO;
    return message;
  }

};