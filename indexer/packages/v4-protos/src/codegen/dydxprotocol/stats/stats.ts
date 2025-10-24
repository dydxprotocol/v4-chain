import { Timestamp } from "../../google/protobuf/timestamp";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long, toTimestamp, fromTimestamp } from "../../helpers";
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
  /** Notional USDC filled in quantums */

  notional: Long;
}
/** Fill records data about a fill on this block. */

export interface BlockStats_FillSDKType {
  /** Taker wallet address */
  taker: string;
  /** Maker wallet address */

  maker: string;
  /** Notional USDC filled in quantums */

  notional: Long;
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
/** GlobalStats stores global stats */

export interface GlobalStats {
  /** Notional USDC traded in quantums */
  notionalTraded: Long;
}
/** GlobalStats stores global stats */

export interface GlobalStatsSDKType {
  /** Notional USDC traded in quantums */
  notional_traded: Long;
}
/** UserStats stores stats for a User */

export interface UserStats {
  /** Taker USDC in quantums */
  takerNotional: Long;
  /** Maker USDC in quantums */

  makerNotional: Long;
}
/** UserStats stores stats for a User */

export interface UserStatsSDKType {
  /** Taker USDC in quantums */
  taker_notional: Long;
  /** Maker USDC in quantums */

  maker_notional: Long;
}
/** CachedStakeAmount stores the last calculated total staked amount for address */

export interface CachedStakeAmount {
  /** Last calculated total staked amount by the delegator (in coin amount). */
  stakedAmount: Uint8Array;
  /**
   * Block time at which the calculation is cached (in Unix Epoch seconds)
   * Rounded down to nearest second.
   */

  cachedAt: Long;
}
/** CachedStakeAmount stores the last calculated total staked amount for address */

export interface CachedStakeAmountSDKType {
  /** Last calculated total staked amount by the delegator (in coin amount). */
  staked_amount: Uint8Array;
  /**
   * Block time at which the calculation is cached (in Unix Epoch seconds)
   * Rounded down to nearest second.
   */

  cached_at: Long;
}

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
    notional: Long.UZERO
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
    makerNotional: Long.UZERO
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
    return message;
  }

};

function createBaseCachedStakeAmount(): CachedStakeAmount {
  return {
    stakedAmount: new Uint8Array(),
    cachedAt: Long.ZERO
  };
}

export const CachedStakeAmount = {
  encode(message: CachedStakeAmount, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.stakedAmount.length !== 0) {
      writer.uint32(10).bytes(message.stakedAmount);
    }

    if (!message.cachedAt.isZero()) {
      writer.uint32(16).int64(message.cachedAt);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CachedStakeAmount {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCachedStakeAmount();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.stakedAmount = reader.bytes();
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

  fromPartial(object: DeepPartial<CachedStakeAmount>): CachedStakeAmount {
    const message = createBaseCachedStakeAmount();
    message.stakedAmount = object.stakedAmount ?? new Uint8Array();
    message.cachedAt = object.cachedAt !== undefined && object.cachedAt !== null ? Long.fromValue(object.cachedAt) : Long.ZERO;
    return message;
  }

};