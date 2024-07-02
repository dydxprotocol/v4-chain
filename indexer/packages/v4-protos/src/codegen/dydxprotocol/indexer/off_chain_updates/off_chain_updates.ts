import { IndexerOrder, IndexerOrderSDKType, IndexerOrderId, IndexerOrderIdSDKType } from "../protocol/v1/clob";
import { Timestamp } from "../../../google/protobuf/timestamp";
import { OrderRemovalReason, OrderRemovalReasonSDKType } from "../shared/removal_reason";
import * as _m0 from "protobufjs/minimal";
import { toTimestamp, fromTimestamp, DeepPartial, Long } from "../../../helpers";
/**
 * OrderPlacementStatus is an enum for the resulting status after an order is
 * placed.
 */

export enum OrderPlaceV1_OrderPlacementStatus {
  /** ORDER_PLACEMENT_STATUS_UNSPECIFIED - Default value, this is invalid and unused. */
  ORDER_PLACEMENT_STATUS_UNSPECIFIED = 0,

  /**
   * ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED - A best effort opened order is one that has only been confirmed to be
   * placed on the dYdX node sending the off-chain update message.
   * The cases where this happens includes:
   * - The dYdX node places an order in it's in-memory orderbook during the
   *   CheckTx flow.
   * A best effort placed order may not have been placed on other dYdX
   * nodes including other dYdX validator nodes and may still be excluded in
   * future order matches.
   */
  ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED = 1,

  /**
   * ORDER_PLACEMENT_STATUS_OPENED - An opened order is one that is confirmed to be placed on all dYdX nodes
   * (discounting dishonest dYdX nodes) and will be included in any future
   * order matches.
   * This status is used internally by the indexer and will not be sent
   * out by protocol.
   */
  ORDER_PLACEMENT_STATUS_OPENED = 2,
  UNRECOGNIZED = -1,
}
/**
 * OrderPlacementStatus is an enum for the resulting status after an order is
 * placed.
 */

export enum OrderPlaceV1_OrderPlacementStatusSDKType {
  /** ORDER_PLACEMENT_STATUS_UNSPECIFIED - Default value, this is invalid and unused. */
  ORDER_PLACEMENT_STATUS_UNSPECIFIED = 0,

  /**
   * ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED - A best effort opened order is one that has only been confirmed to be
   * placed on the dYdX node sending the off-chain update message.
   * The cases where this happens includes:
   * - The dYdX node places an order in it's in-memory orderbook during the
   *   CheckTx flow.
   * A best effort placed order may not have been placed on other dYdX
   * nodes including other dYdX validator nodes and may still be excluded in
   * future order matches.
   */
  ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED = 1,

  /**
   * ORDER_PLACEMENT_STATUS_OPENED - An opened order is one that is confirmed to be placed on all dYdX nodes
   * (discounting dishonest dYdX nodes) and will be included in any future
   * order matches.
   * This status is used internally by the indexer and will not be sent
   * out by protocol.
   */
  ORDER_PLACEMENT_STATUS_OPENED = 2,
  UNRECOGNIZED = -1,
}
export function orderPlaceV1_OrderPlacementStatusFromJSON(object: any): OrderPlaceV1_OrderPlacementStatus {
  switch (object) {
    case 0:
    case "ORDER_PLACEMENT_STATUS_UNSPECIFIED":
      return OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_UNSPECIFIED;

    case 1:
    case "ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED":
      return OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED;

    case 2:
    case "ORDER_PLACEMENT_STATUS_OPENED":
      return OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED;

    case -1:
    case "UNRECOGNIZED":
    default:
      return OrderPlaceV1_OrderPlacementStatus.UNRECOGNIZED;
  }
}
export function orderPlaceV1_OrderPlacementStatusToJSON(object: OrderPlaceV1_OrderPlacementStatus): string {
  switch (object) {
    case OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_UNSPECIFIED:
      return "ORDER_PLACEMENT_STATUS_UNSPECIFIED";

    case OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED:
      return "ORDER_PLACEMENT_STATUS_BEST_EFFORT_OPENED";

    case OrderPlaceV1_OrderPlacementStatus.ORDER_PLACEMENT_STATUS_OPENED:
      return "ORDER_PLACEMENT_STATUS_OPENED";

    case OrderPlaceV1_OrderPlacementStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/**
 * OrderRemovalStatus is an enum for the resulting status after an order is
 * removed.
 */

export enum OrderRemoveV1_OrderRemovalStatus {
  /** ORDER_REMOVAL_STATUS_UNSPECIFIED - Default value, this is invalid and unused. */
  ORDER_REMOVAL_STATUS_UNSPECIFIED = 0,

  /**
   * ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED - A best effort canceled order is one that has only been confirmed to be
   * removed on the dYdX node sending the off-chain update message.
   * The cases where this happens includes:
   * - the order was removed due to the dYdX node receiving a CancelOrder
   *   transaction for the order.
   * - the order was removed due to being undercollateralized during
   *   optimistic matching.
   * A best effort canceled order may not have been removed on other dYdX
   * nodes including other dYdX validator nodes and may still be included in
   * future order matches.
   */
  ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED = 1,

  /**
   * ORDER_REMOVAL_STATUS_CANCELED - A canceled order is one that is confirmed to be removed on all dYdX nodes
   * (discounting dishonest dYdX nodes) and will not be included in any future
   * order matches.
   * The cases where this happens includes:
   * - the order is expired.
   */
  ORDER_REMOVAL_STATUS_CANCELED = 2,

  /** ORDER_REMOVAL_STATUS_FILLED - An order was fully-filled. Only sent by the Indexer for stateful orders. */
  ORDER_REMOVAL_STATUS_FILLED = 3,
  UNRECOGNIZED = -1,
}
/**
 * OrderRemovalStatus is an enum for the resulting status after an order is
 * removed.
 */

export enum OrderRemoveV1_OrderRemovalStatusSDKType {
  /** ORDER_REMOVAL_STATUS_UNSPECIFIED - Default value, this is invalid and unused. */
  ORDER_REMOVAL_STATUS_UNSPECIFIED = 0,

  /**
   * ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED - A best effort canceled order is one that has only been confirmed to be
   * removed on the dYdX node sending the off-chain update message.
   * The cases where this happens includes:
   * - the order was removed due to the dYdX node receiving a CancelOrder
   *   transaction for the order.
   * - the order was removed due to being undercollateralized during
   *   optimistic matching.
   * A best effort canceled order may not have been removed on other dYdX
   * nodes including other dYdX validator nodes and may still be included in
   * future order matches.
   */
  ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED = 1,

  /**
   * ORDER_REMOVAL_STATUS_CANCELED - A canceled order is one that is confirmed to be removed on all dYdX nodes
   * (discounting dishonest dYdX nodes) and will not be included in any future
   * order matches.
   * The cases where this happens includes:
   * - the order is expired.
   */
  ORDER_REMOVAL_STATUS_CANCELED = 2,

  /** ORDER_REMOVAL_STATUS_FILLED - An order was fully-filled. Only sent by the Indexer for stateful orders. */
  ORDER_REMOVAL_STATUS_FILLED = 3,
  UNRECOGNIZED = -1,
}
export function orderRemoveV1_OrderRemovalStatusFromJSON(object: any): OrderRemoveV1_OrderRemovalStatus {
  switch (object) {
    case 0:
    case "ORDER_REMOVAL_STATUS_UNSPECIFIED":
      return OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_UNSPECIFIED;

    case 1:
    case "ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED":
      return OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED;

    case 2:
    case "ORDER_REMOVAL_STATUS_CANCELED":
      return OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED;

    case 3:
    case "ORDER_REMOVAL_STATUS_FILLED":
      return OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_FILLED;

    case -1:
    case "UNRECOGNIZED":
    default:
      return OrderRemoveV1_OrderRemovalStatus.UNRECOGNIZED;
  }
}
export function orderRemoveV1_OrderRemovalStatusToJSON(object: OrderRemoveV1_OrderRemovalStatus): string {
  switch (object) {
    case OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_UNSPECIFIED:
      return "ORDER_REMOVAL_STATUS_UNSPECIFIED";

    case OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED:
      return "ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED";

    case OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED:
      return "ORDER_REMOVAL_STATUS_CANCELED";

    case OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_FILLED:
      return "ORDER_REMOVAL_STATUS_FILLED";

    case OrderRemoveV1_OrderRemovalStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/** OrderPlace messages contain the order placed/replaced. */

export interface OrderPlaceV1 {
  order?: IndexerOrder;
  placementStatus: OrderPlaceV1_OrderPlacementStatus;
  /** The timestamp of the order placement. */

  timeStamp?: Date;
}
/** OrderPlace messages contain the order placed/replaced. */

export interface OrderPlaceV1SDKType {
  order?: IndexerOrderSDKType;
  placement_status: OrderPlaceV1_OrderPlacementStatusSDKType;
  /** The timestamp of the order placement. */

  time_stamp?: Date;
}
/**
 * OrderRemove messages contain the id of the order removed, the reason for the
 * removal and the resulting status from the removal.
 */

export interface OrderRemoveV1 {
  removedOrderId?: IndexerOrderId;
  reason: OrderRemovalReason;
  removalStatus: OrderRemoveV1_OrderRemovalStatus;
  /** The timestamp of the order removal. */

  timeStamp?: Date;
}
/**
 * OrderRemove messages contain the id of the order removed, the reason for the
 * removal and the resulting status from the removal.
 */

export interface OrderRemoveV1SDKType {
  removed_order_id?: IndexerOrderIdSDKType;
  reason: OrderRemovalReasonSDKType;
  removal_status: OrderRemoveV1_OrderRemovalStatusSDKType;
  /** The timestamp of the order removal. */

  time_stamp?: Date;
}
/**
 * OrderUpdate messages contain the id of the order being updated, and the
 * updated total filled quantums of the order.
 */

export interface OrderUpdateV1 {
  orderId?: IndexerOrderId;
  totalFilledQuantums: Long;
}
/**
 * OrderUpdate messages contain the id of the order being updated, and the
 * updated total filled quantums of the order.
 */

export interface OrderUpdateV1SDKType {
  order_id?: IndexerOrderIdSDKType;
  total_filled_quantums: Long;
}
/** OrderReplace messages contain the old order ID and the replacement order. */

export interface OrderReplaceV1 {
  /** vault replaces orders with a different order ID */
  oldOrderId?: IndexerOrderId;
  order?: IndexerOrder;
  placementStatus: OrderPlaceV1_OrderPlacementStatus;
  timeStamp?: Date;
}
/** OrderReplace messages contain the old order ID and the replacement order. */

export interface OrderReplaceV1SDKType {
  /** vault replaces orders with a different order ID */
  old_order_id?: IndexerOrderIdSDKType;
  order?: IndexerOrderSDKType;
  placement_status: OrderPlaceV1_OrderPlacementStatusSDKType;
  time_stamp?: Date;
}
/**
 * An OffChainUpdate message is the message type which will be sent on Kafka to
 * the Indexer.
 */

export interface OffChainUpdateV1 {
  orderPlace?: OrderPlaceV1;
  orderRemove?: OrderRemoveV1;
  orderUpdate?: OrderUpdateV1;
  orderReplace?: OrderReplaceV1;
}
/**
 * An OffChainUpdate message is the message type which will be sent on Kafka to
 * the Indexer.
 */

export interface OffChainUpdateV1SDKType {
  order_place?: OrderPlaceV1SDKType;
  order_remove?: OrderRemoveV1SDKType;
  order_update?: OrderUpdateV1SDKType;
  order_replace?: OrderReplaceV1SDKType;
}

function createBaseOrderPlaceV1(): OrderPlaceV1 {
  return {
    order: undefined,
    placementStatus: 0,
    timeStamp: undefined
  };
}

export const OrderPlaceV1 = {
  encode(message: OrderPlaceV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.order !== undefined) {
      IndexerOrder.encode(message.order, writer.uint32(10).fork()).ldelim();
    }

    if (message.placementStatus !== 0) {
      writer.uint32(16).int32(message.placementStatus);
    }

    if (message.timeStamp !== undefined) {
      Timestamp.encode(toTimestamp(message.timeStamp), writer.uint32(26).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OrderPlaceV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrderPlaceV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.order = IndexerOrder.decode(reader, reader.uint32());
          break;

        case 2:
          message.placementStatus = (reader.int32() as any);
          break;

        case 3:
          message.timeStamp = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OrderPlaceV1>): OrderPlaceV1 {
    const message = createBaseOrderPlaceV1();
    message.order = object.order !== undefined && object.order !== null ? IndexerOrder.fromPartial(object.order) : undefined;
    message.placementStatus = object.placementStatus ?? 0;
    message.timeStamp = object.timeStamp ?? undefined;
    return message;
  }

};

function createBaseOrderRemoveV1(): OrderRemoveV1 {
  return {
    removedOrderId: undefined,
    reason: 0,
    removalStatus: 0,
    timeStamp: undefined
  };
}

export const OrderRemoveV1 = {
  encode(message: OrderRemoveV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.removedOrderId !== undefined) {
      IndexerOrderId.encode(message.removedOrderId, writer.uint32(10).fork()).ldelim();
    }

    if (message.reason !== 0) {
      writer.uint32(16).int32(message.reason);
    }

    if (message.removalStatus !== 0) {
      writer.uint32(24).int32(message.removalStatus);
    }

    if (message.timeStamp !== undefined) {
      Timestamp.encode(toTimestamp(message.timeStamp), writer.uint32(34).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OrderRemoveV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrderRemoveV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.removedOrderId = IndexerOrderId.decode(reader, reader.uint32());
          break;

        case 2:
          message.reason = (reader.int32() as any);
          break;

        case 3:
          message.removalStatus = (reader.int32() as any);
          break;

        case 4:
          message.timeStamp = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OrderRemoveV1>): OrderRemoveV1 {
    const message = createBaseOrderRemoveV1();
    message.removedOrderId = object.removedOrderId !== undefined && object.removedOrderId !== null ? IndexerOrderId.fromPartial(object.removedOrderId) : undefined;
    message.reason = object.reason ?? 0;
    message.removalStatus = object.removalStatus ?? 0;
    message.timeStamp = object.timeStamp ?? undefined;
    return message;
  }

};

function createBaseOrderUpdateV1(): OrderUpdateV1 {
  return {
    orderId: undefined,
    totalFilledQuantums: Long.UZERO
  };
}

export const OrderUpdateV1 = {
  encode(message: OrderUpdateV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orderId !== undefined) {
      IndexerOrderId.encode(message.orderId, writer.uint32(10).fork()).ldelim();
    }

    if (!message.totalFilledQuantums.isZero()) {
      writer.uint32(16).uint64(message.totalFilledQuantums);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OrderUpdateV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrderUpdateV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.orderId = IndexerOrderId.decode(reader, reader.uint32());
          break;

        case 2:
          message.totalFilledQuantums = (reader.uint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OrderUpdateV1>): OrderUpdateV1 {
    const message = createBaseOrderUpdateV1();
    message.orderId = object.orderId !== undefined && object.orderId !== null ? IndexerOrderId.fromPartial(object.orderId) : undefined;
    message.totalFilledQuantums = object.totalFilledQuantums !== undefined && object.totalFilledQuantums !== null ? Long.fromValue(object.totalFilledQuantums) : Long.UZERO;
    return message;
  }

};

function createBaseOrderReplaceV1(): OrderReplaceV1 {
  return {
    oldOrderId: undefined,
    order: undefined,
    placementStatus: 0,
    timeStamp: undefined
  };
}

export const OrderReplaceV1 = {
  encode(message: OrderReplaceV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.oldOrderId !== undefined) {
      IndexerOrderId.encode(message.oldOrderId, writer.uint32(10).fork()).ldelim();
    }

    if (message.order !== undefined) {
      IndexerOrder.encode(message.order, writer.uint32(18).fork()).ldelim();
    }

    if (message.placementStatus !== 0) {
      writer.uint32(24).int32(message.placementStatus);
    }

    if (message.timeStamp !== undefined) {
      Timestamp.encode(toTimestamp(message.timeStamp), writer.uint32(34).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OrderReplaceV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrderReplaceV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.oldOrderId = IndexerOrderId.decode(reader, reader.uint32());
          break;

        case 2:
          message.order = IndexerOrder.decode(reader, reader.uint32());
          break;

        case 3:
          message.placementStatus = (reader.int32() as any);
          break;

        case 4:
          message.timeStamp = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OrderReplaceV1>): OrderReplaceV1 {
    const message = createBaseOrderReplaceV1();
    message.oldOrderId = object.oldOrderId !== undefined && object.oldOrderId !== null ? IndexerOrderId.fromPartial(object.oldOrderId) : undefined;
    message.order = object.order !== undefined && object.order !== null ? IndexerOrder.fromPartial(object.order) : undefined;
    message.placementStatus = object.placementStatus ?? 0;
    message.timeStamp = object.timeStamp ?? undefined;
    return message;
  }

};

function createBaseOffChainUpdateV1(): OffChainUpdateV1 {
  return {
    orderPlace: undefined,
    orderRemove: undefined,
    orderUpdate: undefined,
    orderReplace: undefined
  };
}

export const OffChainUpdateV1 = {
  encode(message: OffChainUpdateV1, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orderPlace !== undefined) {
      OrderPlaceV1.encode(message.orderPlace, writer.uint32(10).fork()).ldelim();
    }

    if (message.orderRemove !== undefined) {
      OrderRemoveV1.encode(message.orderRemove, writer.uint32(18).fork()).ldelim();
    }

    if (message.orderUpdate !== undefined) {
      OrderUpdateV1.encode(message.orderUpdate, writer.uint32(26).fork()).ldelim();
    }

    if (message.orderReplace !== undefined) {
      OrderReplaceV1.encode(message.orderReplace, writer.uint32(34).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OffChainUpdateV1 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOffChainUpdateV1();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.orderPlace = OrderPlaceV1.decode(reader, reader.uint32());
          break;

        case 2:
          message.orderRemove = OrderRemoveV1.decode(reader, reader.uint32());
          break;

        case 3:
          message.orderUpdate = OrderUpdateV1.decode(reader, reader.uint32());
          break;

        case 4:
          message.orderReplace = OrderReplaceV1.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OffChainUpdateV1>): OffChainUpdateV1 {
    const message = createBaseOffChainUpdateV1();
    message.orderPlace = object.orderPlace !== undefined && object.orderPlace !== null ? OrderPlaceV1.fromPartial(object.orderPlace) : undefined;
    message.orderRemove = object.orderRemove !== undefined && object.orderRemove !== null ? OrderRemoveV1.fromPartial(object.orderRemove) : undefined;
    message.orderUpdate = object.orderUpdate !== undefined && object.orderUpdate !== null ? OrderUpdateV1.fromPartial(object.orderUpdate) : undefined;
    message.orderReplace = object.orderReplace !== undefined && object.orderReplace !== null ? OrderReplaceV1.fromPartial(object.orderReplace) : undefined;
    return message;
  }

};