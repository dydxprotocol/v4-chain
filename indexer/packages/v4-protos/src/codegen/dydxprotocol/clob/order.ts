import { SubaccountId, SubaccountIdAmino, SubaccountIdSDKType } from "../subaccounts/subaccount";
import { BinaryReader, BinaryWriter } from "../../binary";
/**
 * Represents the side of the orderbook the order will be placed on.
 * Note that Side.SIDE_UNSPECIFIED is an invalid order and cannot be
 * placed on the orderbook.
 */
export enum Order_Side {
  /** SIDE_UNSPECIFIED - Default value. This value is invalid and unused. */
  SIDE_UNSPECIFIED = 0,
  /** SIDE_BUY - SIDE_BUY is used to represent a BUY order. */
  SIDE_BUY = 1,
  /** SIDE_SELL - SIDE_SELL is used to represent a SELL order. */
  SIDE_SELL = 2,
  UNRECOGNIZED = -1,
}
export const Order_SideSDKType = Order_Side;
export const Order_SideAmino = Order_Side;
export function order_SideFromJSON(object: any): Order_Side {
  switch (object) {
    case 0:
    case "SIDE_UNSPECIFIED":
      return Order_Side.SIDE_UNSPECIFIED;
    case 1:
    case "SIDE_BUY":
      return Order_Side.SIDE_BUY;
    case 2:
    case "SIDE_SELL":
      return Order_Side.SIDE_SELL;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Order_Side.UNRECOGNIZED;
  }
}
export function order_SideToJSON(object: Order_Side): string {
  switch (object) {
    case Order_Side.SIDE_UNSPECIFIED:
      return "SIDE_UNSPECIFIED";
    case Order_Side.SIDE_BUY:
      return "SIDE_BUY";
    case Order_Side.SIDE_SELL:
      return "SIDE_SELL";
    case Order_Side.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/**
 * TimeInForce indicates how long an order will remain active before it
 * is executed or expires.
 */
export enum Order_TimeInForce {
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
export const Order_TimeInForceSDKType = Order_TimeInForce;
export const Order_TimeInForceAmino = Order_TimeInForce;
export function order_TimeInForceFromJSON(object: any): Order_TimeInForce {
  switch (object) {
    case 0:
    case "TIME_IN_FORCE_UNSPECIFIED":
      return Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED;
    case 1:
    case "TIME_IN_FORCE_IOC":
      return Order_TimeInForce.TIME_IN_FORCE_IOC;
    case 2:
    case "TIME_IN_FORCE_POST_ONLY":
      return Order_TimeInForce.TIME_IN_FORCE_POST_ONLY;
    case 3:
    case "TIME_IN_FORCE_FILL_OR_KILL":
      return Order_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Order_TimeInForce.UNRECOGNIZED;
  }
}
export function order_TimeInForceToJSON(object: Order_TimeInForce): string {
  switch (object) {
    case Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED:
      return "TIME_IN_FORCE_UNSPECIFIED";
    case Order_TimeInForce.TIME_IN_FORCE_IOC:
      return "TIME_IN_FORCE_IOC";
    case Order_TimeInForce.TIME_IN_FORCE_POST_ONLY:
      return "TIME_IN_FORCE_POST_ONLY";
    case Order_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL:
      return "TIME_IN_FORCE_FILL_OR_KILL";
    case Order_TimeInForce.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
export enum Order_ConditionType {
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
export const Order_ConditionTypeSDKType = Order_ConditionType;
export const Order_ConditionTypeAmino = Order_ConditionType;
export function order_ConditionTypeFromJSON(object: any): Order_ConditionType {
  switch (object) {
    case 0:
    case "CONDITION_TYPE_UNSPECIFIED":
      return Order_ConditionType.CONDITION_TYPE_UNSPECIFIED;
    case 1:
    case "CONDITION_TYPE_STOP_LOSS":
      return Order_ConditionType.CONDITION_TYPE_STOP_LOSS;
    case 2:
    case "CONDITION_TYPE_TAKE_PROFIT":
      return Order_ConditionType.CONDITION_TYPE_TAKE_PROFIT;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Order_ConditionType.UNRECOGNIZED;
  }
}
export function order_ConditionTypeToJSON(object: Order_ConditionType): string {
  switch (object) {
    case Order_ConditionType.CONDITION_TYPE_UNSPECIFIED:
      return "CONDITION_TYPE_UNSPECIFIED";
    case Order_ConditionType.CONDITION_TYPE_STOP_LOSS:
      return "CONDITION_TYPE_STOP_LOSS";
    case Order_ConditionType.CONDITION_TYPE_TAKE_PROFIT:
      return "CONDITION_TYPE_TAKE_PROFIT";
    case Order_ConditionType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/** OrderId refers to a single order belonging to a Subaccount. */
export interface OrderId {
  /**
   * The subaccount ID that opened this order.
   * Note that this field has `gogoproto.nullable = false` so that it is
   * generated as a value instead of a pointer. This is because the `OrderId`
   * proto is used as a key within maps, and map comparisons will compare
   * pointers for equality (when the desired behavior is to compare the values).
   */
  subaccountId: SubaccountId;
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
export interface OrderIdProtoMsg {
  typeUrl: "/dydxprotocol.clob.OrderId";
  value: Uint8Array;
}
/** OrderId refers to a single order belonging to a Subaccount. */
export interface OrderIdAmino {
  /**
   * The subaccount ID that opened this order.
   * Note that this field has `gogoproto.nullable = false` so that it is
   * generated as a value instead of a pointer. This is because the `OrderId`
   * proto is used as a key within maps, and map comparisons will compare
   * pointers for equality (when the desired behavior is to compare the values).
   */
  subaccount_id?: SubaccountIdAmino;
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
export interface OrderIdAminoMsg {
  type: "/dydxprotocol.clob.OrderId";
  value: OrderIdAmino;
}
/** OrderId refers to a single order belonging to a Subaccount. */
export interface OrderIdSDKType {
  subaccount_id: SubaccountIdSDKType;
  client_id: number;
  order_flags: number;
  clob_pair_id: number;
}
/**
 * OrdersFilledDuringLatestBlock represents a list of `OrderIds` that were
 * filled by any non-zero amount in the latest block.
 */
export interface OrdersFilledDuringLatestBlock {
  /**
   * A list of unique order_ids that were filled by any non-zero amount in the
   * latest block.
   */
  orderIds: OrderId[];
}
export interface OrdersFilledDuringLatestBlockProtoMsg {
  typeUrl: "/dydxprotocol.clob.OrdersFilledDuringLatestBlock";
  value: Uint8Array;
}
/**
 * OrdersFilledDuringLatestBlock represents a list of `OrderIds` that were
 * filled by any non-zero amount in the latest block.
 */
export interface OrdersFilledDuringLatestBlockAmino {
  /**
   * A list of unique order_ids that were filled by any non-zero amount in the
   * latest block.
   */
  order_ids?: OrderIdAmino[];
}
export interface OrdersFilledDuringLatestBlockAminoMsg {
  type: "/dydxprotocol.clob.OrdersFilledDuringLatestBlock";
  value: OrdersFilledDuringLatestBlockAmino;
}
/**
 * OrdersFilledDuringLatestBlock represents a list of `OrderIds` that were
 * filled by any non-zero amount in the latest block.
 */
export interface OrdersFilledDuringLatestBlockSDKType {
  order_ids: OrderIdSDKType[];
}
/**
 * PotentiallyPrunableOrders represents a list of orders that may be prunable
 * from state at a future block height.
 */
export interface PotentiallyPrunableOrders {
  /**
   * A list of unique order_ids that may potentially be pruned from state at a
   * future block height.
   */
  orderIds: OrderId[];
}
export interface PotentiallyPrunableOrdersProtoMsg {
  typeUrl: "/dydxprotocol.clob.PotentiallyPrunableOrders";
  value: Uint8Array;
}
/**
 * PotentiallyPrunableOrders represents a list of orders that may be prunable
 * from state at a future block height.
 */
export interface PotentiallyPrunableOrdersAmino {
  /**
   * A list of unique order_ids that may potentially be pruned from state at a
   * future block height.
   */
  order_ids?: OrderIdAmino[];
}
export interface PotentiallyPrunableOrdersAminoMsg {
  type: "/dydxprotocol.clob.PotentiallyPrunableOrders";
  value: PotentiallyPrunableOrdersAmino;
}
/**
 * PotentiallyPrunableOrders represents a list of orders that may be prunable
 * from state at a future block height.
 */
export interface PotentiallyPrunableOrdersSDKType {
  order_ids: OrderIdSDKType[];
}
/**
 * OrderFillState represents the fill amount of an order according to on-chain
 * state. This proto includes both the current on-chain fill amount of the
 * order, as well as the block at which this information can be pruned from
 * state.
 */
export interface OrderFillState {
  /** The current fillAmount of the order according to on-chain state. */
  fillAmount: bigint;
  /**
   * The block height at which the fillAmount state for this order can be
   * pruned.
   */
  prunableBlockHeight: number;
}
export interface OrderFillStateProtoMsg {
  typeUrl: "/dydxprotocol.clob.OrderFillState";
  value: Uint8Array;
}
/**
 * OrderFillState represents the fill amount of an order according to on-chain
 * state. This proto includes both the current on-chain fill amount of the
 * order, as well as the block at which this information can be pruned from
 * state.
 */
export interface OrderFillStateAmino {
  /** The current fillAmount of the order according to on-chain state. */
  fill_amount?: string;
  /**
   * The block height at which the fillAmount state for this order can be
   * pruned.
   */
  prunable_block_height?: number;
}
export interface OrderFillStateAminoMsg {
  type: "/dydxprotocol.clob.OrderFillState";
  value: OrderFillStateAmino;
}
/**
 * OrderFillState represents the fill amount of an order according to on-chain
 * state. This proto includes both the current on-chain fill amount of the
 * order, as well as the block at which this information can be pruned from
 * state.
 */
export interface OrderFillStateSDKType {
  fill_amount: bigint;
  prunable_block_height: number;
}
/**
 * StatefulOrderTimeSliceValue represents the type of the value of the
 * `StatefulOrdersTimeSlice` in state. The `StatefulOrdersTimeSlice`
 * in state consists of key/value pairs where the keys are UTF-8-encoded
 * `RFC3339NANO` timestamp strings with right-padded zeroes and no
 * time zone info, and the values are of type `StatefulOrderTimeSliceValue`.
 * This `StatefulOrderTimeSliceValue` in state is used for managing stateful
 * order expiration. Stateful order expirations can be for either long term
 * or conditional orders.
 */
export interface StatefulOrderTimeSliceValue {
  /**
   * A unique list of order_ids that expire at this timestamp, sorted in
   * ascending order by block height and transaction index of each stateful
   * order.
   */
  orderIds: OrderId[];
}
export interface StatefulOrderTimeSliceValueProtoMsg {
  typeUrl: "/dydxprotocol.clob.StatefulOrderTimeSliceValue";
  value: Uint8Array;
}
/**
 * StatefulOrderTimeSliceValue represents the type of the value of the
 * `StatefulOrdersTimeSlice` in state. The `StatefulOrdersTimeSlice`
 * in state consists of key/value pairs where the keys are UTF-8-encoded
 * `RFC3339NANO` timestamp strings with right-padded zeroes and no
 * time zone info, and the values are of type `StatefulOrderTimeSliceValue`.
 * This `StatefulOrderTimeSliceValue` in state is used for managing stateful
 * order expiration. Stateful order expirations can be for either long term
 * or conditional orders.
 */
export interface StatefulOrderTimeSliceValueAmino {
  /**
   * A unique list of order_ids that expire at this timestamp, sorted in
   * ascending order by block height and transaction index of each stateful
   * order.
   */
  order_ids?: OrderIdAmino[];
}
export interface StatefulOrderTimeSliceValueAminoMsg {
  type: "/dydxprotocol.clob.StatefulOrderTimeSliceValue";
  value: StatefulOrderTimeSliceValueAmino;
}
/**
 * StatefulOrderTimeSliceValue represents the type of the value of the
 * `StatefulOrdersTimeSlice` in state. The `StatefulOrdersTimeSlice`
 * in state consists of key/value pairs where the keys are UTF-8-encoded
 * `RFC3339NANO` timestamp strings with right-padded zeroes and no
 * time zone info, and the values are of type `StatefulOrderTimeSliceValue`.
 * This `StatefulOrderTimeSliceValue` in state is used for managing stateful
 * order expiration. Stateful order expirations can be for either long term
 * or conditional orders.
 */
export interface StatefulOrderTimeSliceValueSDKType {
  order_ids: OrderIdSDKType[];
}
/**
 * LongTermOrderPlacement represents the placement of a stateful order in
 * state. It stores the stateful order itself and the `BlockHeight` and
 * `TransactionIndex` at which the order was placed.
 */
export interface LongTermOrderPlacement {
  order: Order;
  /**
   * The block height and transaction index at which the order was placed.
   * Used for ordering by time priority when the chain is restarted.
   */
  placementIndex: TransactionOrdering;
}
export interface LongTermOrderPlacementProtoMsg {
  typeUrl: "/dydxprotocol.clob.LongTermOrderPlacement";
  value: Uint8Array;
}
/**
 * LongTermOrderPlacement represents the placement of a stateful order in
 * state. It stores the stateful order itself and the `BlockHeight` and
 * `TransactionIndex` at which the order was placed.
 */
export interface LongTermOrderPlacementAmino {
  order?: OrderAmino;
  /**
   * The block height and transaction index at which the order was placed.
   * Used for ordering by time priority when the chain is restarted.
   */
  placement_index?: TransactionOrderingAmino;
}
export interface LongTermOrderPlacementAminoMsg {
  type: "/dydxprotocol.clob.LongTermOrderPlacement";
  value: LongTermOrderPlacementAmino;
}
/**
 * LongTermOrderPlacement represents the placement of a stateful order in
 * state. It stores the stateful order itself and the `BlockHeight` and
 * `TransactionIndex` at which the order was placed.
 */
export interface LongTermOrderPlacementSDKType {
  order: OrderSDKType;
  placement_index: TransactionOrderingSDKType;
}
/**
 * ConditionalOrderPlacement represents the placement of a conditional order in
 * state. It stores the stateful order itself, the `BlockHeight` and
 * `TransactionIndex` at which the order was placed and triggered.
 */
export interface ConditionalOrderPlacement {
  order: Order;
  /** The block height and transaction index at which the order was placed. */
  placementIndex: TransactionOrdering;
  /**
   * The block height and transaction index at which the order was triggered.
   * Set to be nil if the transaction has not been triggered.
   * Used for ordering by time priority when the chain is restarted.
   */
  triggerIndex?: TransactionOrdering;
}
export interface ConditionalOrderPlacementProtoMsg {
  typeUrl: "/dydxprotocol.clob.ConditionalOrderPlacement";
  value: Uint8Array;
}
/**
 * ConditionalOrderPlacement represents the placement of a conditional order in
 * state. It stores the stateful order itself, the `BlockHeight` and
 * `TransactionIndex` at which the order was placed and triggered.
 */
export interface ConditionalOrderPlacementAmino {
  order?: OrderAmino;
  /** The block height and transaction index at which the order was placed. */
  placement_index?: TransactionOrderingAmino;
  /**
   * The block height and transaction index at which the order was triggered.
   * Set to be nil if the transaction has not been triggered.
   * Used for ordering by time priority when the chain is restarted.
   */
  trigger_index?: TransactionOrderingAmino;
}
export interface ConditionalOrderPlacementAminoMsg {
  type: "/dydxprotocol.clob.ConditionalOrderPlacement";
  value: ConditionalOrderPlacementAmino;
}
/**
 * ConditionalOrderPlacement represents the placement of a conditional order in
 * state. It stores the stateful order itself, the `BlockHeight` and
 * `TransactionIndex` at which the order was placed and triggered.
 */
export interface ConditionalOrderPlacementSDKType {
  order: OrderSDKType;
  placement_index: TransactionOrderingSDKType;
  trigger_index?: TransactionOrderingSDKType;
}
/**
 * Order represents a single order belonging to a `Subaccount`
 * for a particular `ClobPair`.
 */
export interface Order {
  /** The unique ID of this order. Meant to be unique across all orders. */
  orderId: OrderId;
  side: Order_Side;
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
  timeInForce: Order_TimeInForce;
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
  conditionType: Order_ConditionType;
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
export interface OrderProtoMsg {
  typeUrl: "/dydxprotocol.clob.Order";
  value: Uint8Array;
}
/**
 * Order represents a single order belonging to a `Subaccount`
 * for a particular `ClobPair`.
 */
export interface OrderAmino {
  /** The unique ID of this order. Meant to be unique across all orders. */
  order_id?: OrderIdAmino;
  side?: Order_Side;
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
  time_in_force?: Order_TimeInForce;
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
  condition_type?: Order_ConditionType;
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
export interface OrderAminoMsg {
  type: "/dydxprotocol.clob.Order";
  value: OrderAmino;
}
/**
 * Order represents a single order belonging to a `Subaccount`
 * for a particular `ClobPair`.
 */
export interface OrderSDKType {
  order_id: OrderIdSDKType;
  side: Order_Side;
  quantums: bigint;
  subticks: bigint;
  good_til_block?: number;
  good_til_block_time?: number;
  time_in_force: Order_TimeInForce;
  reduce_only: boolean;
  client_metadata: number;
  condition_type: Order_ConditionType;
  conditional_order_trigger_subticks: bigint;
}
/**
 * TransactionOrdering represents a unique location in the block where a
 * transaction was placed. This proto includes both block height and the
 * transaction index that the specific transaction was placed. This information
 * is used for ordering by time priority when the chain is restarted.
 */
export interface TransactionOrdering {
  /** Block height in which the transaction was placed. */
  blockHeight: number;
  /** Within the block, the unique transaction index. */
  transactionIndex: number;
}
export interface TransactionOrderingProtoMsg {
  typeUrl: "/dydxprotocol.clob.TransactionOrdering";
  value: Uint8Array;
}
/**
 * TransactionOrdering represents a unique location in the block where a
 * transaction was placed. This proto includes both block height and the
 * transaction index that the specific transaction was placed. This information
 * is used for ordering by time priority when the chain is restarted.
 */
export interface TransactionOrderingAmino {
  /** Block height in which the transaction was placed. */
  block_height?: number;
  /** Within the block, the unique transaction index. */
  transaction_index?: number;
}
export interface TransactionOrderingAminoMsg {
  type: "/dydxprotocol.clob.TransactionOrdering";
  value: TransactionOrderingAmino;
}
/**
 * TransactionOrdering represents a unique location in the block where a
 * transaction was placed. This proto includes both block height and the
 * transaction index that the specific transaction was placed. This information
 * is used for ordering by time priority when the chain is restarted.
 */
export interface TransactionOrderingSDKType {
  block_height: number;
  transaction_index: number;
}
function createBaseOrderId(): OrderId {
  return {
    subaccountId: SubaccountId.fromPartial({}),
    clientId: 0,
    orderFlags: 0,
    clobPairId: 0
  };
}
export const OrderId = {
  typeUrl: "/dydxprotocol.clob.OrderId",
  encode(message: OrderId, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.subaccountId !== undefined) {
      SubaccountId.encode(message.subaccountId, writer.uint32(10).fork()).ldelim();
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
  decode(input: BinaryReader | Uint8Array, length?: number): OrderId {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrderId();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.subaccountId = SubaccountId.decode(reader, reader.uint32());
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
  fromPartial(object: Partial<OrderId>): OrderId {
    const message = createBaseOrderId();
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? SubaccountId.fromPartial(object.subaccountId) : undefined;
    message.clientId = object.clientId ?? 0;
    message.orderFlags = object.orderFlags ?? 0;
    message.clobPairId = object.clobPairId ?? 0;
    return message;
  },
  fromAmino(object: OrderIdAmino): OrderId {
    const message = createBaseOrderId();
    if (object.subaccount_id !== undefined && object.subaccount_id !== null) {
      message.subaccountId = SubaccountId.fromAmino(object.subaccount_id);
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
  toAmino(message: OrderId): OrderIdAmino {
    const obj: any = {};
    obj.subaccount_id = message.subaccountId ? SubaccountId.toAmino(message.subaccountId) : undefined;
    obj.client_id = message.clientId;
    obj.order_flags = message.orderFlags;
    obj.clob_pair_id = message.clobPairId;
    return obj;
  },
  fromAminoMsg(object: OrderIdAminoMsg): OrderId {
    return OrderId.fromAmino(object.value);
  },
  fromProtoMsg(message: OrderIdProtoMsg): OrderId {
    return OrderId.decode(message.value);
  },
  toProto(message: OrderId): Uint8Array {
    return OrderId.encode(message).finish();
  },
  toProtoMsg(message: OrderId): OrderIdProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.OrderId",
      value: OrderId.encode(message).finish()
    };
  }
};
function createBaseOrdersFilledDuringLatestBlock(): OrdersFilledDuringLatestBlock {
  return {
    orderIds: []
  };
}
export const OrdersFilledDuringLatestBlock = {
  typeUrl: "/dydxprotocol.clob.OrdersFilledDuringLatestBlock",
  encode(message: OrdersFilledDuringLatestBlock, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.orderIds) {
      OrderId.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): OrdersFilledDuringLatestBlock {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrdersFilledDuringLatestBlock();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orderIds.push(OrderId.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<OrdersFilledDuringLatestBlock>): OrdersFilledDuringLatestBlock {
    const message = createBaseOrdersFilledDuringLatestBlock();
    message.orderIds = object.orderIds?.map(e => OrderId.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: OrdersFilledDuringLatestBlockAmino): OrdersFilledDuringLatestBlock {
    const message = createBaseOrdersFilledDuringLatestBlock();
    message.orderIds = object.order_ids?.map(e => OrderId.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: OrdersFilledDuringLatestBlock): OrdersFilledDuringLatestBlockAmino {
    const obj: any = {};
    if (message.orderIds) {
      obj.order_ids = message.orderIds.map(e => e ? OrderId.toAmino(e) : undefined);
    } else {
      obj.order_ids = [];
    }
    return obj;
  },
  fromAminoMsg(object: OrdersFilledDuringLatestBlockAminoMsg): OrdersFilledDuringLatestBlock {
    return OrdersFilledDuringLatestBlock.fromAmino(object.value);
  },
  fromProtoMsg(message: OrdersFilledDuringLatestBlockProtoMsg): OrdersFilledDuringLatestBlock {
    return OrdersFilledDuringLatestBlock.decode(message.value);
  },
  toProto(message: OrdersFilledDuringLatestBlock): Uint8Array {
    return OrdersFilledDuringLatestBlock.encode(message).finish();
  },
  toProtoMsg(message: OrdersFilledDuringLatestBlock): OrdersFilledDuringLatestBlockProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.OrdersFilledDuringLatestBlock",
      value: OrdersFilledDuringLatestBlock.encode(message).finish()
    };
  }
};
function createBasePotentiallyPrunableOrders(): PotentiallyPrunableOrders {
  return {
    orderIds: []
  };
}
export const PotentiallyPrunableOrders = {
  typeUrl: "/dydxprotocol.clob.PotentiallyPrunableOrders",
  encode(message: PotentiallyPrunableOrders, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.orderIds) {
      OrderId.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): PotentiallyPrunableOrders {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePotentiallyPrunableOrders();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orderIds.push(OrderId.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<PotentiallyPrunableOrders>): PotentiallyPrunableOrders {
    const message = createBasePotentiallyPrunableOrders();
    message.orderIds = object.orderIds?.map(e => OrderId.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: PotentiallyPrunableOrdersAmino): PotentiallyPrunableOrders {
    const message = createBasePotentiallyPrunableOrders();
    message.orderIds = object.order_ids?.map(e => OrderId.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: PotentiallyPrunableOrders): PotentiallyPrunableOrdersAmino {
    const obj: any = {};
    if (message.orderIds) {
      obj.order_ids = message.orderIds.map(e => e ? OrderId.toAmino(e) : undefined);
    } else {
      obj.order_ids = [];
    }
    return obj;
  },
  fromAminoMsg(object: PotentiallyPrunableOrdersAminoMsg): PotentiallyPrunableOrders {
    return PotentiallyPrunableOrders.fromAmino(object.value);
  },
  fromProtoMsg(message: PotentiallyPrunableOrdersProtoMsg): PotentiallyPrunableOrders {
    return PotentiallyPrunableOrders.decode(message.value);
  },
  toProto(message: PotentiallyPrunableOrders): Uint8Array {
    return PotentiallyPrunableOrders.encode(message).finish();
  },
  toProtoMsg(message: PotentiallyPrunableOrders): PotentiallyPrunableOrdersProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.PotentiallyPrunableOrders",
      value: PotentiallyPrunableOrders.encode(message).finish()
    };
  }
};
function createBaseOrderFillState(): OrderFillState {
  return {
    fillAmount: BigInt(0),
    prunableBlockHeight: 0
  };
}
export const OrderFillState = {
  typeUrl: "/dydxprotocol.clob.OrderFillState",
  encode(message: OrderFillState, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.fillAmount !== BigInt(0)) {
      writer.uint32(8).uint64(message.fillAmount);
    }
    if (message.prunableBlockHeight !== 0) {
      writer.uint32(16).uint32(message.prunableBlockHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): OrderFillState {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrderFillState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.fillAmount = reader.uint64();
          break;
        case 2:
          message.prunableBlockHeight = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<OrderFillState>): OrderFillState {
    const message = createBaseOrderFillState();
    message.fillAmount = object.fillAmount !== undefined && object.fillAmount !== null ? BigInt(object.fillAmount.toString()) : BigInt(0);
    message.prunableBlockHeight = object.prunableBlockHeight ?? 0;
    return message;
  },
  fromAmino(object: OrderFillStateAmino): OrderFillState {
    const message = createBaseOrderFillState();
    if (object.fill_amount !== undefined && object.fill_amount !== null) {
      message.fillAmount = BigInt(object.fill_amount);
    }
    if (object.prunable_block_height !== undefined && object.prunable_block_height !== null) {
      message.prunableBlockHeight = object.prunable_block_height;
    }
    return message;
  },
  toAmino(message: OrderFillState): OrderFillStateAmino {
    const obj: any = {};
    obj.fill_amount = message.fillAmount ? message.fillAmount.toString() : undefined;
    obj.prunable_block_height = message.prunableBlockHeight;
    return obj;
  },
  fromAminoMsg(object: OrderFillStateAminoMsg): OrderFillState {
    return OrderFillState.fromAmino(object.value);
  },
  fromProtoMsg(message: OrderFillStateProtoMsg): OrderFillState {
    return OrderFillState.decode(message.value);
  },
  toProto(message: OrderFillState): Uint8Array {
    return OrderFillState.encode(message).finish();
  },
  toProtoMsg(message: OrderFillState): OrderFillStateProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.OrderFillState",
      value: OrderFillState.encode(message).finish()
    };
  }
};
function createBaseStatefulOrderTimeSliceValue(): StatefulOrderTimeSliceValue {
  return {
    orderIds: []
  };
}
export const StatefulOrderTimeSliceValue = {
  typeUrl: "/dydxprotocol.clob.StatefulOrderTimeSliceValue",
  encode(message: StatefulOrderTimeSliceValue, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.orderIds) {
      OrderId.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): StatefulOrderTimeSliceValue {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatefulOrderTimeSliceValue();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orderIds.push(OrderId.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<StatefulOrderTimeSliceValue>): StatefulOrderTimeSliceValue {
    const message = createBaseStatefulOrderTimeSliceValue();
    message.orderIds = object.orderIds?.map(e => OrderId.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: StatefulOrderTimeSliceValueAmino): StatefulOrderTimeSliceValue {
    const message = createBaseStatefulOrderTimeSliceValue();
    message.orderIds = object.order_ids?.map(e => OrderId.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: StatefulOrderTimeSliceValue): StatefulOrderTimeSliceValueAmino {
    const obj: any = {};
    if (message.orderIds) {
      obj.order_ids = message.orderIds.map(e => e ? OrderId.toAmino(e) : undefined);
    } else {
      obj.order_ids = [];
    }
    return obj;
  },
  fromAminoMsg(object: StatefulOrderTimeSliceValueAminoMsg): StatefulOrderTimeSliceValue {
    return StatefulOrderTimeSliceValue.fromAmino(object.value);
  },
  fromProtoMsg(message: StatefulOrderTimeSliceValueProtoMsg): StatefulOrderTimeSliceValue {
    return StatefulOrderTimeSliceValue.decode(message.value);
  },
  toProto(message: StatefulOrderTimeSliceValue): Uint8Array {
    return StatefulOrderTimeSliceValue.encode(message).finish();
  },
  toProtoMsg(message: StatefulOrderTimeSliceValue): StatefulOrderTimeSliceValueProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.StatefulOrderTimeSliceValue",
      value: StatefulOrderTimeSliceValue.encode(message).finish()
    };
  }
};
function createBaseLongTermOrderPlacement(): LongTermOrderPlacement {
  return {
    order: Order.fromPartial({}),
    placementIndex: TransactionOrdering.fromPartial({})
  };
}
export const LongTermOrderPlacement = {
  typeUrl: "/dydxprotocol.clob.LongTermOrderPlacement",
  encode(message: LongTermOrderPlacement, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.order !== undefined) {
      Order.encode(message.order, writer.uint32(10).fork()).ldelim();
    }
    if (message.placementIndex !== undefined) {
      TransactionOrdering.encode(message.placementIndex, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): LongTermOrderPlacement {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLongTermOrderPlacement();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.order = Order.decode(reader, reader.uint32());
          break;
        case 2:
          message.placementIndex = TransactionOrdering.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<LongTermOrderPlacement>): LongTermOrderPlacement {
    const message = createBaseLongTermOrderPlacement();
    message.order = object.order !== undefined && object.order !== null ? Order.fromPartial(object.order) : undefined;
    message.placementIndex = object.placementIndex !== undefined && object.placementIndex !== null ? TransactionOrdering.fromPartial(object.placementIndex) : undefined;
    return message;
  },
  fromAmino(object: LongTermOrderPlacementAmino): LongTermOrderPlacement {
    const message = createBaseLongTermOrderPlacement();
    if (object.order !== undefined && object.order !== null) {
      message.order = Order.fromAmino(object.order);
    }
    if (object.placement_index !== undefined && object.placement_index !== null) {
      message.placementIndex = TransactionOrdering.fromAmino(object.placement_index);
    }
    return message;
  },
  toAmino(message: LongTermOrderPlacement): LongTermOrderPlacementAmino {
    const obj: any = {};
    obj.order = message.order ? Order.toAmino(message.order) : undefined;
    obj.placement_index = message.placementIndex ? TransactionOrdering.toAmino(message.placementIndex) : undefined;
    return obj;
  },
  fromAminoMsg(object: LongTermOrderPlacementAminoMsg): LongTermOrderPlacement {
    return LongTermOrderPlacement.fromAmino(object.value);
  },
  fromProtoMsg(message: LongTermOrderPlacementProtoMsg): LongTermOrderPlacement {
    return LongTermOrderPlacement.decode(message.value);
  },
  toProto(message: LongTermOrderPlacement): Uint8Array {
    return LongTermOrderPlacement.encode(message).finish();
  },
  toProtoMsg(message: LongTermOrderPlacement): LongTermOrderPlacementProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.LongTermOrderPlacement",
      value: LongTermOrderPlacement.encode(message).finish()
    };
  }
};
function createBaseConditionalOrderPlacement(): ConditionalOrderPlacement {
  return {
    order: Order.fromPartial({}),
    placementIndex: TransactionOrdering.fromPartial({}),
    triggerIndex: undefined
  };
}
export const ConditionalOrderPlacement = {
  typeUrl: "/dydxprotocol.clob.ConditionalOrderPlacement",
  encode(message: ConditionalOrderPlacement, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.order !== undefined) {
      Order.encode(message.order, writer.uint32(10).fork()).ldelim();
    }
    if (message.placementIndex !== undefined) {
      TransactionOrdering.encode(message.placementIndex, writer.uint32(18).fork()).ldelim();
    }
    if (message.triggerIndex !== undefined) {
      TransactionOrdering.encode(message.triggerIndex, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ConditionalOrderPlacement {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConditionalOrderPlacement();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.order = Order.decode(reader, reader.uint32());
          break;
        case 2:
          message.placementIndex = TransactionOrdering.decode(reader, reader.uint32());
          break;
        case 3:
          message.triggerIndex = TransactionOrdering.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ConditionalOrderPlacement>): ConditionalOrderPlacement {
    const message = createBaseConditionalOrderPlacement();
    message.order = object.order !== undefined && object.order !== null ? Order.fromPartial(object.order) : undefined;
    message.placementIndex = object.placementIndex !== undefined && object.placementIndex !== null ? TransactionOrdering.fromPartial(object.placementIndex) : undefined;
    message.triggerIndex = object.triggerIndex !== undefined && object.triggerIndex !== null ? TransactionOrdering.fromPartial(object.triggerIndex) : undefined;
    return message;
  },
  fromAmino(object: ConditionalOrderPlacementAmino): ConditionalOrderPlacement {
    const message = createBaseConditionalOrderPlacement();
    if (object.order !== undefined && object.order !== null) {
      message.order = Order.fromAmino(object.order);
    }
    if (object.placement_index !== undefined && object.placement_index !== null) {
      message.placementIndex = TransactionOrdering.fromAmino(object.placement_index);
    }
    if (object.trigger_index !== undefined && object.trigger_index !== null) {
      message.triggerIndex = TransactionOrdering.fromAmino(object.trigger_index);
    }
    return message;
  },
  toAmino(message: ConditionalOrderPlacement): ConditionalOrderPlacementAmino {
    const obj: any = {};
    obj.order = message.order ? Order.toAmino(message.order) : undefined;
    obj.placement_index = message.placementIndex ? TransactionOrdering.toAmino(message.placementIndex) : undefined;
    obj.trigger_index = message.triggerIndex ? TransactionOrdering.toAmino(message.triggerIndex) : undefined;
    return obj;
  },
  fromAminoMsg(object: ConditionalOrderPlacementAminoMsg): ConditionalOrderPlacement {
    return ConditionalOrderPlacement.fromAmino(object.value);
  },
  fromProtoMsg(message: ConditionalOrderPlacementProtoMsg): ConditionalOrderPlacement {
    return ConditionalOrderPlacement.decode(message.value);
  },
  toProto(message: ConditionalOrderPlacement): Uint8Array {
    return ConditionalOrderPlacement.encode(message).finish();
  },
  toProtoMsg(message: ConditionalOrderPlacement): ConditionalOrderPlacementProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.ConditionalOrderPlacement",
      value: ConditionalOrderPlacement.encode(message).finish()
    };
  }
};
function createBaseOrder(): Order {
  return {
    orderId: OrderId.fromPartial({}),
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
export const Order = {
  typeUrl: "/dydxprotocol.clob.Order",
  encode(message: Order, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.orderId !== undefined) {
      OrderId.encode(message.orderId, writer.uint32(10).fork()).ldelim();
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
  decode(input: BinaryReader | Uint8Array, length?: number): Order {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrder();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orderId = OrderId.decode(reader, reader.uint32());
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
  fromPartial(object: Partial<Order>): Order {
    const message = createBaseOrder();
    message.orderId = object.orderId !== undefined && object.orderId !== null ? OrderId.fromPartial(object.orderId) : undefined;
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
  fromAmino(object: OrderAmino): Order {
    const message = createBaseOrder();
    if (object.order_id !== undefined && object.order_id !== null) {
      message.orderId = OrderId.fromAmino(object.order_id);
    }
    if (object.side !== undefined && object.side !== null) {
      message.side = order_SideFromJSON(object.side);
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
      message.timeInForce = order_TimeInForceFromJSON(object.time_in_force);
    }
    if (object.reduce_only !== undefined && object.reduce_only !== null) {
      message.reduceOnly = object.reduce_only;
    }
    if (object.client_metadata !== undefined && object.client_metadata !== null) {
      message.clientMetadata = object.client_metadata;
    }
    if (object.condition_type !== undefined && object.condition_type !== null) {
      message.conditionType = order_ConditionTypeFromJSON(object.condition_type);
    }
    if (object.conditional_order_trigger_subticks !== undefined && object.conditional_order_trigger_subticks !== null) {
      message.conditionalOrderTriggerSubticks = BigInt(object.conditional_order_trigger_subticks);
    }
    return message;
  },
  toAmino(message: Order): OrderAmino {
    const obj: any = {};
    obj.order_id = message.orderId ? OrderId.toAmino(message.orderId) : undefined;
    obj.side = order_SideToJSON(message.side);
    obj.quantums = message.quantums ? message.quantums.toString() : undefined;
    obj.subticks = message.subticks ? message.subticks.toString() : undefined;
    obj.good_til_block = message.goodTilBlock;
    obj.good_til_block_time = message.goodTilBlockTime;
    obj.time_in_force = order_TimeInForceToJSON(message.timeInForce);
    obj.reduce_only = message.reduceOnly;
    obj.client_metadata = message.clientMetadata;
    obj.condition_type = order_ConditionTypeToJSON(message.conditionType);
    obj.conditional_order_trigger_subticks = message.conditionalOrderTriggerSubticks ? message.conditionalOrderTriggerSubticks.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: OrderAminoMsg): Order {
    return Order.fromAmino(object.value);
  },
  fromProtoMsg(message: OrderProtoMsg): Order {
    return Order.decode(message.value);
  },
  toProto(message: Order): Uint8Array {
    return Order.encode(message).finish();
  },
  toProtoMsg(message: Order): OrderProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.Order",
      value: Order.encode(message).finish()
    };
  }
};
function createBaseTransactionOrdering(): TransactionOrdering {
  return {
    blockHeight: 0,
    transactionIndex: 0
  };
}
export const TransactionOrdering = {
  typeUrl: "/dydxprotocol.clob.TransactionOrdering",
  encode(message: TransactionOrdering, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.blockHeight !== 0) {
      writer.uint32(8).uint32(message.blockHeight);
    }
    if (message.transactionIndex !== 0) {
      writer.uint32(16).uint32(message.transactionIndex);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): TransactionOrdering {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTransactionOrdering();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.blockHeight = reader.uint32();
          break;
        case 2:
          message.transactionIndex = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<TransactionOrdering>): TransactionOrdering {
    const message = createBaseTransactionOrdering();
    message.blockHeight = object.blockHeight ?? 0;
    message.transactionIndex = object.transactionIndex ?? 0;
    return message;
  },
  fromAmino(object: TransactionOrderingAmino): TransactionOrdering {
    const message = createBaseTransactionOrdering();
    if (object.block_height !== undefined && object.block_height !== null) {
      message.blockHeight = object.block_height;
    }
    if (object.transaction_index !== undefined && object.transaction_index !== null) {
      message.transactionIndex = object.transaction_index;
    }
    return message;
  },
  toAmino(message: TransactionOrdering): TransactionOrderingAmino {
    const obj: any = {};
    obj.block_height = message.blockHeight;
    obj.transaction_index = message.transactionIndex;
    return obj;
  },
  fromAminoMsg(object: TransactionOrderingAminoMsg): TransactionOrdering {
    return TransactionOrdering.fromAmino(object.value);
  },
  fromProtoMsg(message: TransactionOrderingProtoMsg): TransactionOrdering {
    return TransactionOrdering.decode(message.value);
  },
  toProto(message: TransactionOrdering): Uint8Array {
    return TransactionOrdering.encode(message).finish();
  },
  toProtoMsg(message: TransactionOrdering): TransactionOrderingProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.TransactionOrdering",
      value: TransactionOrdering.encode(message).finish()
    };
  }
};