import { BinaryReader, BinaryWriter } from "../../binary";
/** Defines the block rate limits for CLOB specific operations. */
export interface BlockRateLimitConfiguration {
  /**
   * How many short term order attempts (successful and failed) are allowed for
   * an account per N blocks. Note that the rate limits are applied
   * in an AND fashion such that an order placement must pass all rate limit
   * configurations.
   * 
   * Specifying 0 values disables this rate limit.
   */
  maxShortTermOrdersPerNBlocks: MaxPerNBlocksRateLimit[];
  /**
   * How many stateful order attempts (successful and failed) are allowed for
   * an account per N blocks. Note that the rate limits are applied
   * in an AND fashion such that an order placement must pass all rate limit
   * configurations.
   * 
   * Specifying 0 values disables this rate limit.
   */
  maxStatefulOrdersPerNBlocks: MaxPerNBlocksRateLimit[];
  maxShortTermOrderCancellationsPerNBlocks: MaxPerNBlocksRateLimit[];
}
export interface BlockRateLimitConfigurationProtoMsg {
  typeUrl: "/dydxprotocol.clob.BlockRateLimitConfiguration";
  value: Uint8Array;
}
/** Defines the block rate limits for CLOB specific operations. */
export interface BlockRateLimitConfigurationAmino {
  /**
   * How many short term order attempts (successful and failed) are allowed for
   * an account per N blocks. Note that the rate limits are applied
   * in an AND fashion such that an order placement must pass all rate limit
   * configurations.
   * 
   * Specifying 0 values disables this rate limit.
   */
  max_short_term_orders_per_n_blocks?: MaxPerNBlocksRateLimitAmino[];
  /**
   * How many stateful order attempts (successful and failed) are allowed for
   * an account per N blocks. Note that the rate limits are applied
   * in an AND fashion such that an order placement must pass all rate limit
   * configurations.
   * 
   * Specifying 0 values disables this rate limit.
   */
  max_stateful_orders_per_n_blocks?: MaxPerNBlocksRateLimitAmino[];
  max_short_term_order_cancellations_per_n_blocks?: MaxPerNBlocksRateLimitAmino[];
}
export interface BlockRateLimitConfigurationAminoMsg {
  type: "/dydxprotocol.clob.BlockRateLimitConfiguration";
  value: BlockRateLimitConfigurationAmino;
}
/** Defines the block rate limits for CLOB specific operations. */
export interface BlockRateLimitConfigurationSDKType {
  max_short_term_orders_per_n_blocks: MaxPerNBlocksRateLimitSDKType[];
  max_stateful_orders_per_n_blocks: MaxPerNBlocksRateLimitSDKType[];
  max_short_term_order_cancellations_per_n_blocks: MaxPerNBlocksRateLimitSDKType[];
}
/** Defines a rate limit over a specific number of blocks. */
export interface MaxPerNBlocksRateLimit {
  /**
   * How many blocks the rate limit is over.
   * Specifying 0 is invalid.
   */
  numBlocks: number;
  /**
   * What the limit is for `num_blocks`.
   * Specifying 0 is invalid.
   */
  limit: number;
}
export interface MaxPerNBlocksRateLimitProtoMsg {
  typeUrl: "/dydxprotocol.clob.MaxPerNBlocksRateLimit";
  value: Uint8Array;
}
/** Defines a rate limit over a specific number of blocks. */
export interface MaxPerNBlocksRateLimitAmino {
  /**
   * How many blocks the rate limit is over.
   * Specifying 0 is invalid.
   */
  num_blocks?: number;
  /**
   * What the limit is for `num_blocks`.
   * Specifying 0 is invalid.
   */
  limit?: number;
}
export interface MaxPerNBlocksRateLimitAminoMsg {
  type: "/dydxprotocol.clob.MaxPerNBlocksRateLimit";
  value: MaxPerNBlocksRateLimitAmino;
}
/** Defines a rate limit over a specific number of blocks. */
export interface MaxPerNBlocksRateLimitSDKType {
  num_blocks: number;
  limit: number;
}
function createBaseBlockRateLimitConfiguration(): BlockRateLimitConfiguration {
  return {
    maxShortTermOrdersPerNBlocks: [],
    maxStatefulOrdersPerNBlocks: [],
    maxShortTermOrderCancellationsPerNBlocks: []
  };
}
export const BlockRateLimitConfiguration = {
  typeUrl: "/dydxprotocol.clob.BlockRateLimitConfiguration",
  encode(message: BlockRateLimitConfiguration, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.maxShortTermOrdersPerNBlocks) {
      MaxPerNBlocksRateLimit.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.maxStatefulOrdersPerNBlocks) {
      MaxPerNBlocksRateLimit.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.maxShortTermOrderCancellationsPerNBlocks) {
      MaxPerNBlocksRateLimit.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): BlockRateLimitConfiguration {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockRateLimitConfiguration();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.maxShortTermOrdersPerNBlocks.push(MaxPerNBlocksRateLimit.decode(reader, reader.uint32()));
          break;
        case 2:
          message.maxStatefulOrdersPerNBlocks.push(MaxPerNBlocksRateLimit.decode(reader, reader.uint32()));
          break;
        case 3:
          message.maxShortTermOrderCancellationsPerNBlocks.push(MaxPerNBlocksRateLimit.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<BlockRateLimitConfiguration>): BlockRateLimitConfiguration {
    const message = createBaseBlockRateLimitConfiguration();
    message.maxShortTermOrdersPerNBlocks = object.maxShortTermOrdersPerNBlocks?.map(e => MaxPerNBlocksRateLimit.fromPartial(e)) || [];
    message.maxStatefulOrdersPerNBlocks = object.maxStatefulOrdersPerNBlocks?.map(e => MaxPerNBlocksRateLimit.fromPartial(e)) || [];
    message.maxShortTermOrderCancellationsPerNBlocks = object.maxShortTermOrderCancellationsPerNBlocks?.map(e => MaxPerNBlocksRateLimit.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: BlockRateLimitConfigurationAmino): BlockRateLimitConfiguration {
    const message = createBaseBlockRateLimitConfiguration();
    message.maxShortTermOrdersPerNBlocks = object.max_short_term_orders_per_n_blocks?.map(e => MaxPerNBlocksRateLimit.fromAmino(e)) || [];
    message.maxStatefulOrdersPerNBlocks = object.max_stateful_orders_per_n_blocks?.map(e => MaxPerNBlocksRateLimit.fromAmino(e)) || [];
    message.maxShortTermOrderCancellationsPerNBlocks = object.max_short_term_order_cancellations_per_n_blocks?.map(e => MaxPerNBlocksRateLimit.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: BlockRateLimitConfiguration): BlockRateLimitConfigurationAmino {
    const obj: any = {};
    if (message.maxShortTermOrdersPerNBlocks) {
      obj.max_short_term_orders_per_n_blocks = message.maxShortTermOrdersPerNBlocks.map(e => e ? MaxPerNBlocksRateLimit.toAmino(e) : undefined);
    } else {
      obj.max_short_term_orders_per_n_blocks = [];
    }
    if (message.maxStatefulOrdersPerNBlocks) {
      obj.max_stateful_orders_per_n_blocks = message.maxStatefulOrdersPerNBlocks.map(e => e ? MaxPerNBlocksRateLimit.toAmino(e) : undefined);
    } else {
      obj.max_stateful_orders_per_n_blocks = [];
    }
    if (message.maxShortTermOrderCancellationsPerNBlocks) {
      obj.max_short_term_order_cancellations_per_n_blocks = message.maxShortTermOrderCancellationsPerNBlocks.map(e => e ? MaxPerNBlocksRateLimit.toAmino(e) : undefined);
    } else {
      obj.max_short_term_order_cancellations_per_n_blocks = [];
    }
    return obj;
  },
  fromAminoMsg(object: BlockRateLimitConfigurationAminoMsg): BlockRateLimitConfiguration {
    return BlockRateLimitConfiguration.fromAmino(object.value);
  },
  fromProtoMsg(message: BlockRateLimitConfigurationProtoMsg): BlockRateLimitConfiguration {
    return BlockRateLimitConfiguration.decode(message.value);
  },
  toProto(message: BlockRateLimitConfiguration): Uint8Array {
    return BlockRateLimitConfiguration.encode(message).finish();
  },
  toProtoMsg(message: BlockRateLimitConfiguration): BlockRateLimitConfigurationProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.BlockRateLimitConfiguration",
      value: BlockRateLimitConfiguration.encode(message).finish()
    };
  }
};
function createBaseMaxPerNBlocksRateLimit(): MaxPerNBlocksRateLimit {
  return {
    numBlocks: 0,
    limit: 0
  };
}
export const MaxPerNBlocksRateLimit = {
  typeUrl: "/dydxprotocol.clob.MaxPerNBlocksRateLimit",
  encode(message: MaxPerNBlocksRateLimit, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.numBlocks !== 0) {
      writer.uint32(8).uint32(message.numBlocks);
    }
    if (message.limit !== 0) {
      writer.uint32(16).uint32(message.limit);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MaxPerNBlocksRateLimit {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMaxPerNBlocksRateLimit();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.numBlocks = reader.uint32();
          break;
        case 2:
          message.limit = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MaxPerNBlocksRateLimit>): MaxPerNBlocksRateLimit {
    const message = createBaseMaxPerNBlocksRateLimit();
    message.numBlocks = object.numBlocks ?? 0;
    message.limit = object.limit ?? 0;
    return message;
  },
  fromAmino(object: MaxPerNBlocksRateLimitAmino): MaxPerNBlocksRateLimit {
    const message = createBaseMaxPerNBlocksRateLimit();
    if (object.num_blocks !== undefined && object.num_blocks !== null) {
      message.numBlocks = object.num_blocks;
    }
    if (object.limit !== undefined && object.limit !== null) {
      message.limit = object.limit;
    }
    return message;
  },
  toAmino(message: MaxPerNBlocksRateLimit): MaxPerNBlocksRateLimitAmino {
    const obj: any = {};
    obj.num_blocks = message.numBlocks;
    obj.limit = message.limit;
    return obj;
  },
  fromAminoMsg(object: MaxPerNBlocksRateLimitAminoMsg): MaxPerNBlocksRateLimit {
    return MaxPerNBlocksRateLimit.fromAmino(object.value);
  },
  fromProtoMsg(message: MaxPerNBlocksRateLimitProtoMsg): MaxPerNBlocksRateLimit {
    return MaxPerNBlocksRateLimit.decode(message.value);
  },
  toProto(message: MaxPerNBlocksRateLimit): Uint8Array {
    return MaxPerNBlocksRateLimit.encode(message).finish();
  },
  toProtoMsg(message: MaxPerNBlocksRateLimit): MaxPerNBlocksRateLimitProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.MaxPerNBlocksRateLimit",
      value: MaxPerNBlocksRateLimit.encode(message).finish()
    };
  }
};