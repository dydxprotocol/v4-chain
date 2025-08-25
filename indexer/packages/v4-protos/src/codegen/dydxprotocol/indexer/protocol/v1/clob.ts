import { IndexerSubaccountId, IndexerSubaccountIdSDKType } from "./subaccount";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../../../helpers";
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
/**
 * Represents the side of the orderbook the order will be placed on.
 * Note that Side.SIDE_UNSPECIFIED is an invalid order and cannot be
 * placed on the orderbook.
 */

export enum IndexerOrder_SideSDKType {
  /** SIDE_UNSPECIFIED - Default value. This value is invalid and unused. */
  SIDE_UNSPECIFIED = 0,

  /** SIDE_BUY - SIDE_BUY is used to represent a BUY order. */
  SIDE_BUY = 1,

  /** SIDE_SELL - SIDE_SELL is used to represent a SELL order. */
  SIDE_SELL = 2,
  UNRECOGNIZED = -1,
}
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
/**
 * TimeInForce indicates how long an order will remain active before it
 * is executed or expires.
 */

export enum IndexerOrder_TimeInForceSDKType {
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
export enum IndexerOrder_ConditionTypeSDKType {
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
   * is blocked. All open positions are closed and open stateful orders canceled
   * by the protocol when the clob pair transitions to this status. All
   * short-term orders are left to expire.
   */
  CLOB_PAIR_STATUS_FINAL_SETTLEMENT = 6,
  UNRECOGNIZED = -1,
}
/**
 * Status of the CLOB.
 * Defined in clob.clob_pair
 */

export enum ClobPairStatusSDKType {
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
   * is blocked. All open positions are closed and open stateful orders canceled
   * by the protocol when the clob pair transitions to this status. All
   * short-term orders are left to expire.
   */
  CLOB_PAIR_STATUS_FINAL_SETTLEMENT = 6,
  UNRECOGNIZED = -1,
}
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
  subaccountId?: IndexerSubaccountId;
  /**
   * The client ID of this order, unique with respect to the specific
   * sub account (I.E., the same subaccount can't have two orders with
   * the same ClientId).
   */

  clientId: number;
  /**
   * order_flags represent order flags for the order. This field is invalid if
   * it's greater than 257. Each bit represents a different flag.
   * 
   * The following are the valid orderId flags:
   * ShortTerm    = uint32(0)
   * Conditional  = uint32(32)
   * LongTerm     = uint32(64)
   * Twap         = uint32(128)
   * TwapSuborder = uint32(256) (for internal use only)
   * 
   * If both bits are set or bits other than the 2nd and 3rd are set, the order
   * ID is invalid.
   */

  orderFlags: number;
  /** ID of the CLOB the order is created for. */

  clobPairId: number;
}
/** IndexerOrderId refers to a single order belonging to a Subaccount. */

export interface IndexerOrderIdSDKType {
  /**
   * The subaccount ID that opened this order.
   * Note that this field has `gogoproto.nullable = false` so that it is
   * generated as a value instead of a pointer. This is because the `OrderId`
   * proto is used as a key within maps, and map comparisons will compare
   * pointers for equality (when the desired behavior is to compare the values).
   */
  subaccount_id?: IndexerSubaccountIdSDKType;
  /**
   * The client ID of this order, unique with respect to the specific
   * sub account (I.E., the same subaccount can't have two orders with
   * the same ClientId).
   */

  client_id: number;
  /**
   * order_flags represent order flags for the order. This field is invalid if
   * it's greater than 257. Each bit represents a different flag.
   * 
   * The following are the valid orderId flags:
   * ShortTerm    = uint32(0)
   * Conditional  = uint32(32)
   * LongTerm     = uint32(64)
   * Twap         = uint32(128)
   * TwapSuborder = uint32(256) (for internal use only)
   * 
   * If both bits are set or bits other than the 2nd and 3rd are set, the order
   * ID is invalid.
   */

  order_flags: number;
  /** ID of the CLOB the order is created for. */

  clob_pair_id: number;
}
/**
 * IndexerOrderV1 represents a single order belonging to a `Subaccount`
 * for a particular `ClobPair`.
 */

export interface IndexerOrder {
  /** The unique ID of this order. Meant to be unique across all orders. */
  orderId?: IndexerOrderId;
  side: IndexerOrder_Side;
  /**
   * The size of this order in base quantums. Must be a multiple of
   * `ClobPair.StepBaseQuantums` (where `ClobPair.Id = orderId.ClobPairId`).
   */

  quantums: Long;
  /**
   * The price level that this order will be placed at on the orderbook,
   * in subticks. Must be a multiple of ClobPair.SubticksPerTick
   * (where `ClobPair.Id = orderId.ClobPairId`).
   */

  subticks: Long;
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

  conditionalOrderTriggerSubticks: Long;
  /** builder_code_params is the metadata for the partner or builder of an order. */

  builderCodeParams?: BuilderCodeParameters;
  /**
   * order_router_address is the address of the order router that forwarded this
   * order.
   */

  orderRouterAddress: string;
  /**
   * twap_parameters represent the configuration for a TWAP order. This must be
   * set for twap orders and will be ignored for all other order types.
   */

  twapParameters?: TwapParameters;
}
/**
 * IndexerOrderV1 represents a single order belonging to a `Subaccount`
 * for a particular `ClobPair`.
 */

export interface IndexerOrderSDKType {
  /** The unique ID of this order. Meant to be unique across all orders. */
  order_id?: IndexerOrderIdSDKType;
  side: IndexerOrder_SideSDKType;
  /**
   * The size of this order in base quantums. Must be a multiple of
   * `ClobPair.StepBaseQuantums` (where `ClobPair.Id = orderId.ClobPairId`).
   */

  quantums: Long;
  /**
   * The price level that this order will be placed at on the orderbook,
   * in subticks. Must be a multiple of ClobPair.SubticksPerTick
   * (where `ClobPair.Id = orderId.ClobPairId`).
   */

  subticks: Long;
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

  time_in_force: IndexerOrder_TimeInForceSDKType;
  /**
   * Enforces that the order can only reduce the size of an existing position.
   * If a ReduceOnly order would change the side of the existing position,
   * its size is reduced to that of the remaining size of the position.
   * If existing orders on the book with ReduceOnly
   * would already close the position, the least aggressive (out-of-the-money)
   * ReduceOnly orders are resized and canceled first.
   */

  reduce_only: boolean;
  /**
   * Set of bit flags set arbitrarily by clients and ignored by the protocol.
   * Used by indexer to infer information about a placed order.
   */

  client_metadata: number;
  condition_type: IndexerOrder_ConditionTypeSDKType;
  /**
   * conditional_order_trigger_subticks represents the price at which this order
   * will be triggered. If the condition_type is CONDITION_TYPE_UNSPECIFIED,
   * this value is enforced to be 0. If this value is nonzero, condition_type
   * cannot be CONDITION_TYPE_UNSPECIFIED. Value is in subticks.
   * Must be a multiple of ClobPair.SubticksPerTick (where `ClobPair.Id =
   * orderId.ClobPairId`).
   */

  conditional_order_trigger_subticks: Long;
  /** builder_code_params is the metadata for the partner or builder of an order. */

  builder_code_params?: BuilderCodeParametersSDKType;
  /**
   * order_router_address is the address of the order router that forwarded this
   * order.
   */

  order_router_address: string;
  /**
   * twap_parameters represent the configuration for a TWAP order. This must be
   * set for twap orders and will be ignored for all other order types.
   */

  twap_parameters?: TwapParametersSDKType;
}
/** TwapParameters represents the necessary configuration for a TWAP order. */

export interface TwapParameters {
  /**
   * Duration of the TWAP order execution in seconds. Must be in the range
   * [300 (5 minutes), 86400 (24 hours)].
   */
  duration: number;
  /**
   * Interval in seconds for each suborder to execute. Must be a
   * whole number, a factor of the duration, and in the range
   * [30 (30 seconds), 3600 (1 hour)].
   */

  interval: number;
  /**
   * Price tolerance in ppm for each suborder. This will be applied to
   * the oracle price each time a suborder is triggered. Must be
   * be in the range [0, 1_000_000).
   */

  priceTolerance: number;
}
/** TwapParameters represents the necessary configuration for a TWAP order. */

export interface TwapParametersSDKType {
  /**
   * Duration of the TWAP order execution in seconds. Must be in the range
   * [300 (5 minutes), 86400 (24 hours)].
   */
  duration: number;
  /**
   * Interval in seconds for each suborder to execute. Must be a
   * whole number, a factor of the duration, and in the range
   * [30 (30 seconds), 3600 (1 hour)].
   */

  interval: number;
  /**
   * Price tolerance in ppm for each suborder. This will be applied to
   * the oracle price each time a suborder is triggered. Must be
   * be in the range [0, 1_000_000).
   */

  price_tolerance: number;
}
/**
 * BuilderCodeParameters represents the metadata for the partner or builder of
 * an order. This allows them to specify a fee for providing there service which
 * will be paid out in the event of an order fill.
 */

export interface BuilderCodeParameters {
  /** The address of the builder to which the fee will be paid. */
  builderAddress: string;
  /** The fee enforced on the order in ppm. */

  feePpm: number;
}
/**
 * BuilderCodeParameters represents the metadata for the partner or builder of
 * an order. This allows them to specify a fee for providing there service which
 * will be paid out in the event of an order fill.
 */

export interface BuilderCodeParametersSDKType {
  /** The address of the builder to which the fee will be paid. */
  builder_address: string;
  /** The fee enforced on the order in ppm. */

  fee_ppm: number;
}

function createBaseIndexerOrderId(): IndexerOrderId {
  return {
    subaccountId: undefined,
    clientId: 0,
    orderFlags: 0,
    clobPairId: 0
  };
}

export const IndexerOrderId = {
  encode(message: IndexerOrderId, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): IndexerOrderId {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<IndexerOrderId>): IndexerOrderId {
    const message = createBaseIndexerOrderId();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? IndexerSubaccountId.fromPartial(object.subaccountId) : undefined;
    message.clientId = object.clientId ?? 0;
    message.orderFlags = object.orderFlags ?? 0;
    message.clobPairId = object.clobPairId ?? 0;
    return message;
  }

};

function createBaseIndexerOrder(): IndexerOrder {
  return {
    orderId: undefined,
    side: 0,
    quantums: Long.UZERO,
    subticks: Long.UZERO,
    goodTilBlock: undefined,
    goodTilBlockTime: undefined,
    timeInForce: 0,
    reduceOnly: false,
    clientMetadata: 0,
    conditionType: 0,
    conditionalOrderTriggerSubticks: Long.UZERO,
    builderCodeParams: undefined,
    orderRouterAddress: "",
    twapParameters: undefined
  };
}

export const IndexerOrder = {
  encode(message: IndexerOrder, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orderId !== undefined) {
      IndexerOrderId.encode(message.orderId, writer.uint32(10).fork()).ldelim();
    }

    if (message.side !== 0) {
      writer.uint32(16).int32(message.side);
    }

    if (!message.quantums.isZero()) {
      writer.uint32(24).uint64(message.quantums);
    }

    if (!message.subticks.isZero()) {
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

    if (!message.conditionalOrderTriggerSubticks.isZero()) {
      writer.uint32(88).uint64(message.conditionalOrderTriggerSubticks);
    }

    if (message.builderCodeParams !== undefined) {
      BuilderCodeParameters.encode(message.builderCodeParams, writer.uint32(98).fork()).ldelim();
    }

    if (message.orderRouterAddress !== "") {
      writer.uint32(106).string(message.orderRouterAddress);
    }

    if (message.twapParameters !== undefined) {
      TwapParameters.encode(message.twapParameters, writer.uint32(114).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IndexerOrder {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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
          message.quantums = (reader.uint64() as Long);
          break;

        case 4:
          message.subticks = (reader.uint64() as Long);
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
          message.conditionalOrderTriggerSubticks = (reader.uint64() as Long);
          break;

        case 12:
          message.builderCodeParams = BuilderCodeParameters.decode(reader, reader.uint32());
          break;

        case 13:
          message.orderRouterAddress = reader.string();
          break;

        case 14:
          message.twapParameters = TwapParameters.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<IndexerOrder>): IndexerOrder {
    const message = createBaseIndexerOrder();
    message.orderId = object.orderId !== undefined && object.orderId !== null ? IndexerOrderId.fromPartial(object.orderId) : undefined;
    message.side = object.side ?? 0;
    message.quantums = object.quantums !== undefined && object.quantums !== null ? Long.fromValue(object.quantums) : Long.UZERO;
    message.subticks = object.subticks !== undefined && object.subticks !== null ? Long.fromValue(object.subticks) : Long.UZERO;
    message.goodTilBlock = object.goodTilBlock ?? undefined;
    message.goodTilBlockTime = object.goodTilBlockTime ?? undefined;
    message.timeInForce = object.timeInForce ?? 0;
    message.reduceOnly = object.reduceOnly ?? false;
    message.clientMetadata = object.clientMetadata ?? 0;
    message.conditionType = object.conditionType ?? 0;
    message.conditionalOrderTriggerSubticks = object.conditionalOrderTriggerSubticks !== undefined && object.conditionalOrderTriggerSubticks !== null ? Long.fromValue(object.conditionalOrderTriggerSubticks) : Long.UZERO;
    message.builderCodeParams = object.builderCodeParams !== undefined && object.builderCodeParams !== null ? BuilderCodeParameters.fromPartial(object.builderCodeParams) : undefined;
    message.orderRouterAddress = object.orderRouterAddress ?? "";
    message.twapParameters = object.twapParameters !== undefined && object.twapParameters !== null ? TwapParameters.fromPartial(object.twapParameters) : undefined;
    return message;
  }

};

function createBaseTwapParameters(): TwapParameters {
  return {
    duration: 0,
    interval: 0,
    priceTolerance: 0
  };
}

export const TwapParameters = {
  encode(message: TwapParameters, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.duration !== 0) {
      writer.uint32(8).uint32(message.duration);
    }

    if (message.interval !== 0) {
      writer.uint32(16).uint32(message.interval);
    }

    if (message.priceTolerance !== 0) {
      writer.uint32(24).uint32(message.priceTolerance);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TwapParameters {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTwapParameters();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.duration = reader.uint32();
          break;

        case 2:
          message.interval = reader.uint32();
          break;

        case 3:
          message.priceTolerance = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<TwapParameters>): TwapParameters {
    const message = createBaseTwapParameters();
    message.duration = object.duration ?? 0;
    message.interval = object.interval ?? 0;
    message.priceTolerance = object.priceTolerance ?? 0;
    return message;
  }

};

function createBaseBuilderCodeParameters(): BuilderCodeParameters {
  return {
    builderAddress: "",
    feePpm: 0
  };
}

export const BuilderCodeParameters = {
  encode(message: BuilderCodeParameters, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.builderAddress !== "") {
      writer.uint32(10).string(message.builderAddress);
    }

    if (message.feePpm !== 0) {
      writer.uint32(16).uint32(message.feePpm);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BuilderCodeParameters {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBuilderCodeParameters();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.builderAddress = reader.string();
          break;

        case 2:
          message.feePpm = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<BuilderCodeParameters>): BuilderCodeParameters {
    const message = createBaseBuilderCodeParameters();
    message.builderAddress = object.builderAddress ?? "";
    message.feePpm = object.feePpm ?? 0;
    return message;
  }

};