import { IndexerSubaccountId, IndexerSubaccountIdAmino, IndexerSubaccountIdSDKType } from "./subaccount";
import { BinaryReader, BinaryWriter } from "../../../../binary";
/**
 * Represents the side of the orderbook the order will be placed on.
 * Note that Side.SIDE_UNSPECIFIED is an invalid order and cannot be
 * placed on the orderbook.
 */
export enum IndexerOrder_Side {
  /** SIDE_UNSPECIFIED - Default value. This value is invalid and unused. */
  SIDE_UNSPECIFIED = 0,
  /** SIDE_BUY - SIDE_BUY is used to represent a BUY order. */
  SIDE_BUY = 1,
  /** SIDE_SELL - SIDE_SELL is used to represent a SELL order. */
  SIDE_SELL = 2,
  UNRECOGNIZED = -1,
}
export const IndexerOrder_SideSDKType = IndexerOrder_Side;
export const IndexerOrder_SideAmino = IndexerOrder_Side;
export function indexerOrder_SideFromJSON(object: any): IndexerOrder_Side {
  switch (object) {
    case 0:
    case "SIDE_UNSPECIFIED":
      return IndexerOrder_Side.SIDE_UNSPECIFIED;
    case 1:
    case "SIDE_BUY":
      return IndexerOrder_Side.SIDE_BUY;
    case 2:
    case "SIDE_SELL":
      return IndexerOrder_Side.SIDE_SELL;
    case -1:
    case "UNRECOGNIZED":
    default:
      return IndexerOrder_Side.UNRECOGNIZED;
  }
}
export function indexerOrder_SideToJSON(object: IndexerOrder_Side): string {
  switch (object) {
    case IndexerOrder_Side.SIDE_UNSPECIFIED:
      return "SIDE_UNSPECIFIED";
    case IndexerOrder_Side.SIDE_BUY:
      return "SIDE_BUY";
    case IndexerOrder_Side.SIDE_SELL:
      return "SIDE_SELL";
    case IndexerOrder_Side.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/**
 * TimeInForce indicates how long an order will remain active before it
 * is executed or expires.
 */
export enum IndexerOrder_TimeInForce {
  /**
   * TIME_IN_FORCE_UNSPECIFIED - TIME_IN_FORCE_UNSPECIFIED represents the default behavior where an
   * order will first match with existing orders on the book, and any
   * remaining size will be added to the book as a maker order.
   */
  TIME_IN_FORCE_UNSPECIFIED = 0,
  /**
   * TIME_IN_FORCE_IOC - TIME_IN_FORCE_IOC enforces that an order only be matched with
   * maker orders on the book. If the order has remaining size after
   * matching with existing orders on the book, the remaining size
   * is not placed on the book.
   */
  TIME_IN_FORCE_IOC = 1,
  /**
   * TIME_IN_FORCE_POST_ONLY - TIME_IN_FORCE_POST_ONLY enforces that an order only be placed
   * on the book as a maker order. Note this means that validators will cancel
   * any newly-placed post only orders that would cross with other maker
   * orders.
   */
  TIME_IN_FORCE_POST_ONLY = 2,
  /**
   * TIME_IN_FORCE_FILL_OR_KILL - TIME_IN_FORCE_FILL_OR_KILL enforces that an order will either be filled
   * completely and immediately by maker orders on the book or canceled if the
   * entire amount can‘t be matched.
   */
  TIME_IN_FORCE_FILL_OR_KILL = 3,
  UNRECOGNIZED = -1,
}
export const IndexerOrder_TimeInForceSDKType = IndexerOrder_TimeInForce;
export const IndexerOrder_TimeInForceAmino = IndexerOrder_TimeInForce;
export function indexerOrder_TimeInForceFromJSON(object: any): IndexerOrder_TimeInForce {
  switch (object) {
    case 0:
    case "TIME_IN_FORCE_UNSPECIFIED":
      return IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED;
    case 1:
    case "TIME_IN_FORCE_IOC":
      return IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC;
    case 2:
    case "TIME_IN_FORCE_POST_ONLY":
      return IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY;
    case 3:
    case "TIME_IN_FORCE_FILL_OR_KILL":
      return IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL;
    case -1:
    case "UNRECOGNIZED":
    default:
      return IndexerOrder_TimeInForce.UNRECOGNIZED;
  }
}
export function indexerOrder_TimeInForceToJSON(object: IndexerOrder_TimeInForce): string {
  switch (object) {
    case IndexerOrder_TimeInForce.TIME_IN_FORCE_UNSPECIFIED:
      return "TIME_IN_FORCE_UNSPECIFIED";
    case IndexerOrder_TimeInForce.TIME_IN_FORCE_IOC:
      return "TIME_IN_FORCE_IOC";
    case IndexerOrder_TimeInForce.TIME_IN_FORCE_POST_ONLY:
      return "TIME_IN_FORCE_POST_ONLY";
    case IndexerOrder_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL:
      return "TIME_IN_FORCE_FILL_OR_KILL";
    case IndexerOrder_TimeInForce.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export enum IndexerOrder_ConditionType {
  /**
   * CONDITION_TYPE_UNSPECIFIED - CONDITION_TYPE_UNSPECIFIED represents the default behavior where an
   * order will be placed immediately on the orderbook.
   */
  CONDITION_TYPE_UNSPECIFIED = 0,
  /**
   * CONDITION_TYPE_STOP_LOSS - CONDITION_TYPE_STOP_LOSS represents a stop order. A stop order will
   * trigger when the oracle price moves at or above the trigger price for
   * buys, and at or below the trigger price for sells.
   */
  CONDITION_TYPE_STOP_LOSS = 1,
  /**
   * CONDITION_TYPE_TAKE_PROFIT - CONDITION_TYPE_TAKE_PROFIT represents a take profit order. A take profit
   * order will trigger when the oracle price moves at or below the trigger
   * price for buys and at or above the trigger price for sells.
   */
  CONDITION_TYPE_TAKE_PROFIT = 2,
  UNRECOGNIZED = -1,
}
export const IndexerOrder_ConditionTypeSDKType = IndexerOrder_ConditionType;
export const IndexerOrder_ConditionTypeAmino = IndexerOrder_ConditionType;
export function indexerOrder_ConditionTypeFromJSON(object: any): IndexerOrder_ConditionType {
  switch (object) {
    case 0:
    case "CONDITION_TYPE_UNSPECIFIED":
      return IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED;
    case 1:
    case "CONDITION_TYPE_STOP_LOSS":
      return IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS;
    case 2:
    case "CONDITION_TYPE_TAKE_PROFIT":
      return IndexerOrder_ConditionType.CONDITION_TYPE_TAKE_PROFIT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return IndexerOrder_ConditionType.UNRECOGNIZED;
  }
}
export function indexerOrder_ConditionTypeToJSON(object: IndexerOrder_ConditionType): string {
  switch (object) {
    case IndexerOrder_ConditionType.CONDITION_TYPE_UNSPECIFIED:
      return "CONDITION_TYPE_UNSPECIFIED";
    case IndexerOrder_ConditionType.CONDITION_TYPE_STOP_LOSS:
      return "CONDITION_TYPE_STOP_LOSS";
    case IndexerOrder_ConditionType.CONDITION_TYPE_TAKE_PROFIT:
      return "CONDITION_TYPE_TAKE_PROFIT";
    case IndexerOrder_ConditionType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/**
 * Status of the CLOB.
 * Defined in clob.clob_pair
 */
export enum ClobPairStatus {
  /** CLOB_PAIR_STATUS_UNSPECIFIED - Default value. This value is invalid and unused. */
  CLOB_PAIR_STATUS_UNSPECIFIED = 0,
  /**
   * CLOB_PAIR_STATUS_ACTIVE - CLOB_PAIR_STATUS_ACTIVE behavior is unfinalized.
   * TODO(DEC-600): update this documentation.
   */
  CLOB_PAIR_STATUS_ACTIVE = 1,
  /**
   * CLOB_PAIR_STATUS_PAUSED - CLOB_PAIR_STATUS_PAUSED behavior is unfinalized.
   * TODO(DEC-600): update this documentation.
   */
  CLOB_PAIR_STATUS_PAUSED = 2,
  /**
   * CLOB_PAIR_STATUS_CANCEL_ONLY - CLOB_PAIR_STATUS_CANCEL_ONLY behavior is unfinalized.
   * TODO(DEC-600): update this documentation.
   */
  CLOB_PAIR_STATUS_CANCEL_ONLY = 3,
  /**
   * CLOB_PAIR_STATUS_POST_ONLY - CLOB_PAIR_STATUS_POST_ONLY behavior is unfinalized.
   * TODO(DEC-600): update this documentation.
   */
  CLOB_PAIR_STATUS_POST_ONLY = 4,
  /**
   * CLOB_PAIR_STATUS_INITIALIZING - CLOB_PAIR_STATUS_INITIALIZING represents a newly-added clob pair.
   * Clob pairs in this state only accept orders which are
   * both short-term and post-only.
   */
  CLOB_PAIR_STATUS_INITIALIZING = 5,
  /**
   * CLOB_PAIR_STATUS_FINAL_SETTLEMENT - CLOB_PAIR_STATUS_FINAL_SETTLEMENT represents a clob pair that has been
   * deactivated. Clob pairs in this state do not accept new orders and trading
   * is blocked. All open positions are closed by the protocol when the clob
   * pair gains this status.
   */
  CLOB_PAIR_STATUS_FINAL_SETTLEMENT = 6,
  UNRECOGNIZED = -1,
}
export const ClobPairStatusSDKType = ClobPairStatus;
export const ClobPairStatusAmino = ClobPairStatus;
export function clobPairStatusFromJSON(object: any): ClobPairStatus {
  switch (object) {
    case 0:
    case "CLOB_PAIR_STATUS_UNSPECIFIED":
      return ClobPairStatus.CLOB_PAIR_STATUS_UNSPECIFIED;
    case 1:
    case "CLOB_PAIR_STATUS_ACTIVE":
      return ClobPairStatus.CLOB_PAIR_STATUS_ACTIVE;
    case 2:
    case "CLOB_PAIR_STATUS_PAUSED":
      return ClobPairStatus.CLOB_PAIR_STATUS_PAUSED;
    case 3:
    case "CLOB_PAIR_STATUS_CANCEL_ONLY":
      return ClobPairStatus.CLOB_PAIR_STATUS_CANCEL_ONLY;
    case 4:
    case "CLOB_PAIR_STATUS_POST_ONLY":
      return ClobPairStatus.CLOB_PAIR_STATUS_POST_ONLY;
    case 5:
    case "CLOB_PAIR_STATUS_INITIALIZING":
      return ClobPairStatus.CLOB_PAIR_STATUS_INITIALIZING;
    case 6:
    case "CLOB_PAIR_STATUS_FINAL_SETTLEMENT":
      return ClobPairStatus.CLOB_PAIR_STATUS_FINAL_SETTLEMENT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ClobPairStatus.UNRECOGNIZED;
  }
}
export function clobPairStatusToJSON(object: ClobPairStatus): string {
  switch (object) {
    case ClobPairStatus.CLOB_PAIR_STATUS_UNSPECIFIED:
      return "CLOB_PAIR_STATUS_UNSPECIFIED";
    case ClobPairStatus.CLOB_PAIR_STATUS_ACTIVE:
      return "CLOB_PAIR_STATUS_ACTIVE";
    case ClobPairStatus.CLOB_PAIR_STATUS_PAUSED:
      return "CLOB_PAIR_STATUS_PAUSED";
    case ClobPairStatus.CLOB_PAIR_STATUS_CANCEL_ONLY:
      return "CLOB_PAIR_STATUS_CANCEL_ONLY";
    case ClobPairStatus.CLOB_PAIR_STATUS_POST_ONLY:
      return "CLOB_PAIR_STATUS_POST_ONLY";
    case ClobPairStatus.CLOB_PAIR_STATUS_INITIALIZING:
      return "CLOB_PAIR_STATUS_INITIALIZING";
    case ClobPairStatus.CLOB_PAIR_STATUS_FINAL_SETTLEMENT:
      return "CLOB_PAIR_STATUS_FINAL_SETTLEMENT";
    case ClobPairStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/** IndexerOrderId refers to a single order belonging to a Subaccount. */
export interface IndexerOrderId {
  /**
   * The subaccount ID that opened this order.
   * Note that this field has `gogoproto.nullable = false` so that it is
   * generated as a value instead of a pointer. This is because the `OrderId`
   * proto is used as a key within maps, and map comparisons will compare
   * pointers for equality (when the desired behavior is to compare the values).
   */
  subaccountId: IndexerSubaccountId;
  /**
   * The client ID of this order, unique with respect to the specific
   * sub account (I.E., the same subaccount can't have two orders with
   * the same ClientId).
   */
  clientId: number;
  /**
   * order_flags represent order flags for the order. This field is invalid if
   * it's greater than 127 (larger than one byte). Each bit in the first byte
   * represents a different flag. Currently only two flags are supported.
   * 
   * Starting from the bit after the most MSB (note that the MSB is used in
   * proto varint encoding, and therefore cannot be used): Bit 1 is set if this
   * order is a Long-Term order (0x40, or 64 as a uint8). Bit 2 is set if this
   * order is a Conditional order (0x20, or 32 as a uint8).
   * 
   * If neither bit is set, the order is assumed to be a Short-Term order.
   * 
   * If both bits are set or bits other than the 2nd and 3rd are set, the order
   * ID is invalid.
   */
  orderFlags: number;
  /** ID of the CLOB the order is created for. */
  clobPairId: number;
}
export interface IndexerOrderIdProtoMsg {
  typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerOrderId";
  value: Uint8Array;
}
/** IndexerOrderId refers to a single order belonging to a Subaccount. */
export interface IndexerOrderIdAmino {
  /**
   * The subaccount ID that opened this order.
   * Note that this field has `gogoproto.nullable = false` so that it is
   * generated as a value instead of a pointer. This is because the `OrderId`
   * proto is used as a key within maps, and map comparisons will compare
   * pointers for equality (when the desired behavior is to compare the values).
   */
  subaccount_id?: IndexerSubaccountIdAmino;
  /**
   * The client ID of this order, unique with respect to the specific
   * sub account (I.E., the same subaccount can't have two orders with
   * the same ClientId).
   */
  client_id?: number;
  /**
   * order_flags represent order flags for the order. This field is invalid if
   * it's greater than 127 (larger than one byte). Each bit in the first byte
   * represents a different flag. Currently only two flags are supported.
   * 
   * Starting from the bit after the most MSB (note that the MSB is used in
   * proto varint encoding, and therefore cannot be used): Bit 1 is set if this
   * order is a Long-Term order (0x40, or 64 as a uint8). Bit 2 is set if this
   * order is a Conditional order (0x20, or 32 as a uint8).
   * 
   * If neither bit is set, the order is assumed to be a Short-Term order.
   * 
   * If both bits are set or bits other than the 2nd and 3rd are set, the order
   * ID is invalid.
   */
  order_flags?: number;
  /** ID of the CLOB the order is created for. */
  clob_pair_id?: number;
}
export interface IndexerOrderIdAminoMsg {
  type: "/dydxprotocol.indexer.protocol.v1.IndexerOrderId";
  value: IndexerOrderIdAmino;
}
/** IndexerOrderId refers to a single order belonging to a Subaccount. */
export interface IndexerOrderIdSDKType {
  subaccount_id: IndexerSubaccountIdSDKType;
  client_id: number;
  order_flags: number;
  clob_pair_id: number;
}
/**
 * IndexerOrderV1 represents a single order belonging to a `Subaccount`
 * for a particular `ClobPair`.
 */
export interface IndexerOrder {
  /** The unique ID of this order. Meant to be unique across all orders. */
  orderId: IndexerOrderId;
  side: IndexerOrder_Side;
  /**
   * The size of this order in base quantums. Must be a multiple of
   * `ClobPair.StepBaseQuantums` (where `ClobPair.Id = orderId.ClobPairId`).
   */
  quantums: bigint;
  /**
   * The price level that this order will be placed at on the orderbook,
   * in subticks. Must be a multiple of ClobPair.SubticksPerTick
   * (where `ClobPair.Id = orderId.ClobPairId`).
   */
  subticks: bigint;
  /**
   * The last block this order can be executed at (after which it will be
   * unfillable). Used only for Short-Term orders. If this value is non-zero
   * then the order is assumed to be a Short-Term order.
   */
  goodTilBlock?: number;
  /**
   * good_til_block_time represents the unix timestamp (in seconds) at which a
   * stateful order will be considered expired. The
   * good_til_block_time is always evaluated against the previous block's
   * `BlockTime` instead of the block in which the order is committed. If this
   * value is non-zero then the order is assumed to be a stateful or
   * conditional order.
   */
  goodTilBlockTime?: number;
  /** The time in force of this order. */
  timeInForce: IndexerOrder_TimeInForce;
  /**
   * Enforces that the order can only reduce the size of an existing position.
   * If a ReduceOnly order would change the side of the existing position,
   * its size is reduced to that of the remaining size of the position.
   * If existing orders on the book with ReduceOnly
   * would already close the position, the least aggressive (out-of-the-money)
   * ReduceOnly orders are resized and canceled first.
   */
  reduceOnly: boolean;
  /**
   * Set of bit flags set arbitrarily by clients and ignored by the protocol.
   * Used by indexer to infer information about a placed order.
   */
  clientMetadata: number;
  conditionType: IndexerOrder_ConditionType;
  /**
   * conditional_order_trigger_subticks represents the price at which this order
   * will be triggered. If the condition_type is CONDITION_TYPE_UNSPECIFIED,
   * this value is enforced to be 0. If this value is nonzero, condition_type
   * cannot be CONDITION_TYPE_UNSPECIFIED. Value is in subticks.
   * Must be a multiple of ClobPair.SubticksPerTick (where `ClobPair.Id =
   * orderId.ClobPairId`).
   */
  conditionalOrderTriggerSubticks: bigint;
}
export interface IndexerOrderProtoMsg {
  typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerOrder";
  value: Uint8Array;
}
/**
 * IndexerOrderV1 represents a single order belonging to a `Subaccount`
 * for a particular `ClobPair`.
 */
export interface IndexerOrderAmino {
  /** The unique ID of this order. Meant to be unique across all orders. */
  order_id?: IndexerOrderIdAmino;
  side?: IndexerOrder_Side;
  /**
   * The size of this order in base quantums. Must be a multiple of
   * `ClobPair.StepBaseQuantums` (where `ClobPair.Id = orderId.ClobPairId`).
   */
  quantums?: string;
  /**
   * The price level that this order will be placed at on the orderbook,
   * in subticks. Must be a multiple of ClobPair.SubticksPerTick
   * (where `ClobPair.Id = orderId.ClobPairId`).
   */
  subticks?: string;
  /**
   * The last block this order can be executed at (after which it will be
   * unfillable). Used only for Short-Term orders. If this value is non-zero
   * then the order is assumed to be a Short-Term order.
   */
  good_til_block?: number;
  /**
   * good_til_block_time represents the unix timestamp (in seconds) at which a
   * stateful order will be considered expired. The
   * good_til_block_time is always evaluated against the previous block's
   * `BlockTime` instead of the block in which the order is committed. If this
   * value is non-zero then the order is assumed to be a stateful or
   * conditional order.
   */
  good_til_block_time?: number;
  /** The time in force of this order. */
  time_in_force?: IndexerOrder_TimeInForce;
  /**
   * Enforces that the order can only reduce the size of an existing position.
   * If a ReduceOnly order would change the side of the existing position,
   * its size is reduced to that of the remaining size of the position.
   * If existing orders on the book with ReduceOnly
   * would already close the position, the least aggressive (out-of-the-money)
   * ReduceOnly orders are resized and canceled first.
   */
  reduce_only?: boolean;
  /**
   * Set of bit flags set arbitrarily by clients and ignored by the protocol.
   * Used by indexer to infer information about a placed order.
   */
  client_metadata?: number;
  condition_type?: IndexerOrder_ConditionType;
  /**
   * conditional_order_trigger_subticks represents the price at which this order
   * will be triggered. If the condition_type is CONDITION_TYPE_UNSPECIFIED,
   * this value is enforced to be 0. If this value is nonzero, condition_type
   * cannot be CONDITION_TYPE_UNSPECIFIED. Value is in subticks.
   * Must be a multiple of ClobPair.SubticksPerTick (where `ClobPair.Id =
   * orderId.ClobPairId`).
   */
  conditional_order_trigger_subticks?: string;
}
export interface IndexerOrderAminoMsg {
  type: "/dydxprotocol.indexer.protocol.v1.IndexerOrder";
  value: IndexerOrderAmino;
}
/**
 * IndexerOrderV1 represents a single order belonging to a `Subaccount`
 * for a particular `ClobPair`.
 */
export interface IndexerOrderSDKType {
  order_id: IndexerOrderIdSDKType;
  side: IndexerOrder_Side;
  quantums: bigint;
  subticks: bigint;
  good_til_block?: number;
  good_til_block_time?: number;
  time_in_force: IndexerOrder_TimeInForce;
  reduce_only: boolean;
  client_metadata: number;
  condition_type: IndexerOrder_ConditionType;
  conditional_order_trigger_subticks: bigint;
}
function createBaseIndexerOrderId(): IndexerOrderId {
  return {
    subaccountId: IndexerSubaccountId.fromPartial({}),
    clientId: 0,
    orderFlags: 0,
    clobPairId: 0
  };
}
export const IndexerOrderId = {
  typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerOrderId",
  encode(message: IndexerOrderId, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.subaccountId !== undefined) {
      IndexerSubaccountId.encode(message.subaccountId, writer.uint32(10).fork()).ldelim();
    }
    if (message.clientId !== 0) {
      writer.uint32(21).fixed32(message.clientId);
    }
    if (message.orderFlags !== 0) {
      writer.uint32(24).uint32(message.orderFlags);
    }
    if (message.clobPairId !== 0) {
      writer.uint32(32).uint32(message.clobPairId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): IndexerOrderId {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIndexerOrderId();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.subaccountId = IndexerSubaccountId.decode(reader, reader.uint32());
          break;
        case 2:
          message.clientId = reader.fixed32();
          break;
        case 3:
          message.orderFlags = reader.uint32();
          break;
        case 4:
          message.clobPairId = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<IndexerOrderId>): IndexerOrderId {
    const message = createBaseIndexerOrderId();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? IndexerSubaccountId.fromPartial(object.subaccountId) : undefined;
    message.clientId = object.clientId ?? 0;
    message.orderFlags = object.orderFlags ?? 0;
    message.clobPairId = object.clobPairId ?? 0;
    return message;
  },
  fromAmino(object: IndexerOrderIdAmino): IndexerOrderId {
    const message = createBaseIndexerOrderId();
    if (object.subaccount_id !== undefined && object.subaccount_id !== null) {
      message.subaccountId = IndexerSubaccountId.fromAmino(object.subaccount_id);
    }
    if (object.client_id !== undefined && object.client_id !== null) {
      message.clientId = object.client_id;
    }
    if (object.order_flags !== undefined && object.order_flags !== null) {
      message.orderFlags = object.order_flags;
    }
    if (object.clob_pair_id !== undefined && object.clob_pair_id !== null) {
      message.clobPairId = object.clob_pair_id;
    }
    return message;
  },
  toAmino(message: IndexerOrderId): IndexerOrderIdAmino {
    const obj: any = {};
    obj.subaccount_id = message.subaccountId ? IndexerSubaccountId.toAmino(message.subaccountId) : undefined;
    obj.client_id = message.clientId;
    obj.order_flags = message.orderFlags;
    obj.clob_pair_id = message.clobPairId;
    return obj;
  },
  fromAminoMsg(object: IndexerOrderIdAminoMsg): IndexerOrderId {
    return IndexerOrderId.fromAmino(object.value);
  },
  fromProtoMsg(message: IndexerOrderIdProtoMsg): IndexerOrderId {
    return IndexerOrderId.decode(message.value);
  },
  toProto(message: IndexerOrderId): Uint8Array {
    return IndexerOrderId.encode(message).finish();
  },
  toProtoMsg(message: IndexerOrderId): IndexerOrderIdProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerOrderId",
      value: IndexerOrderId.encode(message).finish()
    };
  }
};
function createBaseIndexerOrder(): IndexerOrder {
  return {
    orderId: IndexerOrderId.fromPartial({}),
    side: 0,
    quantums: BigInt(0),
    subticks: BigInt(0),
    goodTilBlock: undefined,
    goodTilBlockTime: undefined,
    timeInForce: 0,
    reduceOnly: false,
    clientMetadata: 0,
    conditionType: 0,
    conditionalOrderTriggerSubticks: BigInt(0)
  };
}
export const IndexerOrder = {
  typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerOrder",
  encode(message: IndexerOrder, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.orderId !== undefined) {
      IndexerOrderId.encode(message.orderId, writer.uint32(10).fork()).ldelim();
    }
    if (message.side !== 0) {
      writer.uint32(16).int32(message.side);
    }
    if (message.quantums !== BigInt(0)) {
      writer.uint32(24).uint64(message.quantums);
    }
    if (message.subticks !== BigInt(0)) {
      writer.uint32(32).uint64(message.subticks);
    }
    if (message.goodTilBlock !== undefined) {
      writer.uint32(40).uint32(message.goodTilBlock);
    }
    if (message.goodTilBlockTime !== undefined) {
      writer.uint32(53).fixed32(message.goodTilBlockTime);
    }
    if (message.timeInForce !== 0) {
      writer.uint32(56).int32(message.timeInForce);
    }
    if (message.reduceOnly === true) {
      writer.uint32(64).bool(message.reduceOnly);
    }
    if (message.clientMetadata !== 0) {
      writer.uint32(72).uint32(message.clientMetadata);
    }
    if (message.conditionType !== 0) {
      writer.uint32(80).int32(message.conditionType);
    }
    if (message.conditionalOrderTriggerSubticks !== BigInt(0)) {
      writer.uint32(88).uint64(message.conditionalOrderTriggerSubticks);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): IndexerOrder {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIndexerOrder();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orderId = IndexerOrderId.decode(reader, reader.uint32());
          break;
        case 2:
          message.side = (reader.int32() as any);
          break;
        case 3:
          message.quantums = reader.uint64();
          break;
        case 4:
          message.subticks = reader.uint64();
          break;
        case 5:
          message.goodTilBlock = reader.uint32();
          break;
        case 6:
          message.goodTilBlockTime = reader.fixed32();
          break;
        case 7:
          message.timeInForce = (reader.int32() as any);
          break;
        case 8:
          message.reduceOnly = reader.bool();
          break;
        case 9:
          message.clientMetadata = reader.uint32();
          break;
        case 10:
          message.conditionType = (reader.int32() as any);
          break;
        case 11:
          message.conditionalOrderTriggerSubticks = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<IndexerOrder>): IndexerOrder {
    const message = createBaseIndexerOrder();
    message.orderId = object.orderId !== undefined && object.orderId !== null ? IndexerOrderId.fromPartial(object.orderId) : undefined;
    message.side = object.side ?? 0;
    message.quantums = object.quantums !== undefined && object.quantums !== null ? BigInt(object.quantums.toString()) : BigInt(0);
    message.subticks = object.subticks !== undefined && object.subticks !== null ? BigInt(object.subticks.toString()) : BigInt(0);
    message.goodTilBlock = object.goodTilBlock ?? undefined;
    message.goodTilBlockTime = object.goodTilBlockTime ?? undefined;
    message.timeInForce = object.timeInForce ?? 0;
    message.reduceOnly = object.reduceOnly ?? false;
    message.clientMetadata = object.clientMetadata ?? 0;
    message.conditionType = object.conditionType ?? 0;
    message.conditionalOrderTriggerSubticks = object.conditionalOrderTriggerSubticks !== undefined && object.conditionalOrderTriggerSubticks !== null ? BigInt(object.conditionalOrderTriggerSubticks.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: IndexerOrderAmino): IndexerOrder {
    const message = createBaseIndexerOrder();
    if (object.order_id !== undefined && object.order_id !== null) {
      message.orderId = IndexerOrderId.fromAmino(object.order_id);
    }
    if (object.side !== undefined && object.side !== null) {
      message.side = indexerOrder_SideFromJSON(object.side);
    }
    if (object.quantums !== undefined && object.quantums !== null) {
      message.quantums = BigInt(object.quantums);
    }
    if (object.subticks !== undefined && object.subticks !== null) {
      message.subticks = BigInt(object.subticks);
    }
    if (object.good_til_block !== undefined && object.good_til_block !== null) {
      message.goodTilBlock = object.good_til_block;
    }
    if (object.good_til_block_time !== undefined && object.good_til_block_time !== null) {
      message.goodTilBlockTime = object.good_til_block_time;
    }
    if (object.time_in_force !== undefined && object.time_in_force !== null) {
      message.timeInForce = indexerOrder_TimeInForceFromJSON(object.time_in_force);
    }
    if (object.reduce_only !== undefined && object.reduce_only !== null) {
      message.reduceOnly = object.reduce_only;
    }
    if (object.client_metadata !== undefined && object.client_metadata !== null) {
      message.clientMetadata = object.client_metadata;
    }
    if (object.condition_type !== undefined && object.condition_type !== null) {
      message.conditionType = indexerOrder_ConditionTypeFromJSON(object.condition_type);
    }
    if (object.conditional_order_trigger_subticks !== undefined && object.conditional_order_trigger_subticks !== null) {
      message.conditionalOrderTriggerSubticks = BigInt(object.conditional_order_trigger_subticks);
    }
    return message;
  },
  toAmino(message: IndexerOrder): IndexerOrderAmino {
    const obj: any = {};
    obj.order_id = message.orderId ? IndexerOrderId.toAmino(message.orderId) : undefined;
    obj.side = indexerOrder_SideToJSON(message.side);
    obj.quantums = message.quantums ? message.quantums.toString() : undefined;
    obj.subticks = message.subticks ? message.subticks.toString() : undefined;
    obj.good_til_block = message.goodTilBlock;
    obj.good_til_block_time = message.goodTilBlockTime;
    obj.time_in_force = indexerOrder_TimeInForceToJSON(message.timeInForce);
    obj.reduce_only = message.reduceOnly;
    obj.client_metadata = message.clientMetadata;
    obj.condition_type = indexerOrder_ConditionTypeToJSON(message.conditionType);
    obj.conditional_order_trigger_subticks = message.conditionalOrderTriggerSubticks ? message.conditionalOrderTriggerSubticks.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: IndexerOrderAminoMsg): IndexerOrder {
    return IndexerOrder.fromAmino(object.value);
  },
  fromProtoMsg(message: IndexerOrderProtoMsg): IndexerOrder {
    return IndexerOrder.decode(message.value);
  },
  toProto(message: IndexerOrder): Uint8Array {
    return IndexerOrder.encode(message).finish();
  },
  toProtoMsg(message: IndexerOrder): IndexerOrderProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.protocol.v1.IndexerOrder",
      value: IndexerOrder.encode(message).finish()
    };
  }
};