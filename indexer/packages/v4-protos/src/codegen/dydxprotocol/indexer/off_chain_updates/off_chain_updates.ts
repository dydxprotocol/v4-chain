import { IndexerOrder, IndexerOrderAmino, IndexerOrderSDKType, IndexerOrderId, IndexerOrderIdAmino, IndexerOrderIdSDKType } from "../protocol/v1/clob";
import { OrderRemovalReason, orderRemovalReasonFromJSON, orderRemovalReasonToJSON } from "../shared/removal_reason";
import { BinaryReader, BinaryWriter } from "../../../binary";
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
export const OrderPlaceV1_OrderPlacementStatusSDKType = OrderPlaceV1_OrderPlacementStatus;
export const OrderPlaceV1_OrderPlacementStatusAmino = OrderPlaceV1_OrderPlacementStatus;
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
export const OrderRemoveV1_OrderRemovalStatusSDKType = OrderRemoveV1_OrderRemovalStatus;
export const OrderRemoveV1_OrderRemovalStatusAmino = OrderRemoveV1_OrderRemovalStatus;
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
}
export interface OrderPlaceV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.off_chain_updates.OrderPlaceV1";
  value: Uint8Array;
}
/** OrderPlace messages contain the order placed/replaced. */
export interface OrderPlaceV1Amino {
  order?: IndexerOrderAmino;
  placement_status?: OrderPlaceV1_OrderPlacementStatus;
}
export interface OrderPlaceV1AminoMsg {
  type: "/dydxprotocol.indexer.off_chain_updates.OrderPlaceV1";
  value: OrderPlaceV1Amino;
}
/** OrderPlace messages contain the order placed/replaced. */
export interface OrderPlaceV1SDKType {
  order?: IndexerOrderSDKType;
  placement_status: OrderPlaceV1_OrderPlacementStatus;
}
/**
 * OrderRemove messages contain the id of the order removed, the reason for the
 * removal and the resulting status from the removal.
 */
export interface OrderRemoveV1 {
  removedOrderId?: IndexerOrderId;
  reason: OrderRemovalReason;
  removalStatus: OrderRemoveV1_OrderRemovalStatus;
}
export interface OrderRemoveV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.off_chain_updates.OrderRemoveV1";
  value: Uint8Array;
}
/**
 * OrderRemove messages contain the id of the order removed, the reason for the
 * removal and the resulting status from the removal.
 */
export interface OrderRemoveV1Amino {
  removed_order_id?: IndexerOrderIdAmino;
  reason?: OrderRemovalReason;
  removal_status?: OrderRemoveV1_OrderRemovalStatus;
}
export interface OrderRemoveV1AminoMsg {
  type: "/dydxprotocol.indexer.off_chain_updates.OrderRemoveV1";
  value: OrderRemoveV1Amino;
}
/**
 * OrderRemove messages contain the id of the order removed, the reason for the
 * removal and the resulting status from the removal.
 */
export interface OrderRemoveV1SDKType {
  removed_order_id?: IndexerOrderIdSDKType;
  reason: OrderRemovalReason;
  removal_status: OrderRemoveV1_OrderRemovalStatus;
}
/**
 * OrderUpdate messages contain the id of the order being updated, and the
 * updated total filled quantums of the order.
 */
export interface OrderUpdateV1 {
  orderId?: IndexerOrderId;
  totalFilledQuantums: bigint;
}
export interface OrderUpdateV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.off_chain_updates.OrderUpdateV1";
  value: Uint8Array;
}
/**
 * OrderUpdate messages contain the id of the order being updated, and the
 * updated total filled quantums of the order.
 */
export interface OrderUpdateV1Amino {
  order_id?: IndexerOrderIdAmino;
  total_filled_quantums?: string;
}
export interface OrderUpdateV1AminoMsg {
  type: "/dydxprotocol.indexer.off_chain_updates.OrderUpdateV1";
  value: OrderUpdateV1Amino;
}
/**
 * OrderUpdate messages contain the id of the order being updated, and the
 * updated total filled quantums of the order.
 */
export interface OrderUpdateV1SDKType {
  order_id?: IndexerOrderIdSDKType;
  total_filled_quantums: bigint;
}
/**
 * An OffChainUpdate message is the message type which will be sent on Kafka to
 * the Indexer.
 */
export interface OffChainUpdateV1 {
  orderPlace?: OrderPlaceV1;
  orderRemove?: OrderRemoveV1;
  orderUpdate?: OrderUpdateV1;
}
export interface OffChainUpdateV1ProtoMsg {
  typeUrl: "/dydxprotocol.indexer.off_chain_updates.OffChainUpdateV1";
  value: Uint8Array;
}
/**
 * An OffChainUpdate message is the message type which will be sent on Kafka to
 * the Indexer.
 */
export interface OffChainUpdateV1Amino {
  order_place?: OrderPlaceV1Amino;
  order_remove?: OrderRemoveV1Amino;
  order_update?: OrderUpdateV1Amino;
}
export interface OffChainUpdateV1AminoMsg {
  type: "/dydxprotocol.indexer.off_chain_updates.OffChainUpdateV1";
  value: OffChainUpdateV1Amino;
}
/**
 * An OffChainUpdate message is the message type which will be sent on Kafka to
 * the Indexer.
 */
export interface OffChainUpdateV1SDKType {
  order_place?: OrderPlaceV1SDKType;
  order_remove?: OrderRemoveV1SDKType;
  order_update?: OrderUpdateV1SDKType;
}
function createBaseOrderPlaceV1(): OrderPlaceV1 {
  return {
    order: undefined,
    placementStatus: 0
  };
}
export const OrderPlaceV1 = {
  typeUrl: "/dydxprotocol.indexer.off_chain_updates.OrderPlaceV1",
  encode(message: OrderPlaceV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.order !== undefined) {
      IndexerOrder.encode(message.order, writer.uint32(10).fork()).ldelim();
    }
    if (message.placementStatus !== 0) {
      writer.uint32(16).int32(message.placementStatus);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): OrderPlaceV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<OrderPlaceV1>): OrderPlaceV1 {
    const message = createBaseOrderPlaceV1();
    message.order = object.order !== undefined && object.order !== null ? IndexerOrder.fromPartial(object.order) : undefined;
    message.placementStatus = object.placementStatus ?? 0;
    return message;
  },
  fromAmino(object: OrderPlaceV1Amino): OrderPlaceV1 {
    const message = createBaseOrderPlaceV1();
    if (object.order !== undefined && object.order !== null) {
      message.order = IndexerOrder.fromAmino(object.order);
    }
    if (object.placement_status !== undefined && object.placement_status !== null) {
      message.placementStatus = orderPlaceV1_OrderPlacementStatusFromJSON(object.placement_status);
    }
    return message;
  },
  toAmino(message: OrderPlaceV1): OrderPlaceV1Amino {
    const obj: any = {};
    obj.order = message.order ? IndexerOrder.toAmino(message.order) : undefined;
    obj.placement_status = orderPlaceV1_OrderPlacementStatusToJSON(message.placementStatus);
    return obj;
  },
  fromAminoMsg(object: OrderPlaceV1AminoMsg): OrderPlaceV1 {
    return OrderPlaceV1.fromAmino(object.value);
  },
  fromProtoMsg(message: OrderPlaceV1ProtoMsg): OrderPlaceV1 {
    return OrderPlaceV1.decode(message.value);
  },
  toProto(message: OrderPlaceV1): Uint8Array {
    return OrderPlaceV1.encode(message).finish();
  },
  toProtoMsg(message: OrderPlaceV1): OrderPlaceV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.off_chain_updates.OrderPlaceV1",
      value: OrderPlaceV1.encode(message).finish()
    };
  }
};
function createBaseOrderRemoveV1(): OrderRemoveV1 {
  return {
    removedOrderId: undefined,
    reason: 0,
    removalStatus: 0
  };
}
export const OrderRemoveV1 = {
  typeUrl: "/dydxprotocol.indexer.off_chain_updates.OrderRemoveV1",
  encode(message: OrderRemoveV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.removedOrderId !== undefined) {
      IndexerOrderId.encode(message.removedOrderId, writer.uint32(10).fork()).ldelim();
    }
    if (message.reason !== 0) {
      writer.uint32(16).int32(message.reason);
    }
    if (message.removalStatus !== 0) {
      writer.uint32(24).int32(message.removalStatus);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): OrderRemoveV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<OrderRemoveV1>): OrderRemoveV1 {
    const message = createBaseOrderRemoveV1();
    message.removedOrderId = object.removedOrderId !== undefined && object.removedOrderId !== null ? IndexerOrderId.fromPartial(object.removedOrderId) : undefined;
    message.reason = object.reason ?? 0;
    message.removalStatus = object.removalStatus ?? 0;
    return message;
  },
  fromAmino(object: OrderRemoveV1Amino): OrderRemoveV1 {
    const message = createBaseOrderRemoveV1();
    if (object.removed_order_id !== undefined && object.removed_order_id !== null) {
      message.removedOrderId = IndexerOrderId.fromAmino(object.removed_order_id);
    }
    if (object.reason !== undefined && object.reason !== null) {
      message.reason = orderRemovalReasonFromJSON(object.reason);
    }
    if (object.removal_status !== undefined && object.removal_status !== null) {
      message.removalStatus = orderRemoveV1_OrderRemovalStatusFromJSON(object.removal_status);
    }
    return message;
  },
  toAmino(message: OrderRemoveV1): OrderRemoveV1Amino {
    const obj: any = {};
    obj.removed_order_id = message.removedOrderId ? IndexerOrderId.toAmino(message.removedOrderId) : undefined;
    obj.reason = orderRemovalReasonToJSON(message.reason);
    obj.removal_status = orderRemoveV1_OrderRemovalStatusToJSON(message.removalStatus);
    return obj;
  },
  fromAminoMsg(object: OrderRemoveV1AminoMsg): OrderRemoveV1 {
    return OrderRemoveV1.fromAmino(object.value);
  },
  fromProtoMsg(message: OrderRemoveV1ProtoMsg): OrderRemoveV1 {
    return OrderRemoveV1.decode(message.value);
  },
  toProto(message: OrderRemoveV1): Uint8Array {
    return OrderRemoveV1.encode(message).finish();
  },
  toProtoMsg(message: OrderRemoveV1): OrderRemoveV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.off_chain_updates.OrderRemoveV1",
      value: OrderRemoveV1.encode(message).finish()
    };
  }
};
function createBaseOrderUpdateV1(): OrderUpdateV1 {
  return {
    orderId: undefined,
    totalFilledQuantums: BigInt(0)
  };
}
export const OrderUpdateV1 = {
  typeUrl: "/dydxprotocol.indexer.off_chain_updates.OrderUpdateV1",
  encode(message: OrderUpdateV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.orderId !== undefined) {
      IndexerOrderId.encode(message.orderId, writer.uint32(10).fork()).ldelim();
    }
    if (message.totalFilledQuantums !== BigInt(0)) {
      writer.uint32(16).uint64(message.totalFilledQuantums);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): OrderUpdateV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrderUpdateV1();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.orderId = IndexerOrderId.decode(reader, reader.uint32());
          break;
        case 2:
          message.totalFilledQuantums = reader.uint64();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<OrderUpdateV1>): OrderUpdateV1 {
    const message = createBaseOrderUpdateV1();
    message.orderId = object.orderId !== undefined && object.orderId !== null ? IndexerOrderId.fromPartial(object.orderId) : undefined;
    message.totalFilledQuantums = object.totalFilledQuantums !== undefined && object.totalFilledQuantums !== null ? BigInt(object.totalFilledQuantums.toString()) : BigInt(0);
    return message;
  },
  fromAmino(object: OrderUpdateV1Amino): OrderUpdateV1 {
    const message = createBaseOrderUpdateV1();
    if (object.order_id !== undefined && object.order_id !== null) {
      message.orderId = IndexerOrderId.fromAmino(object.order_id);
    }
    if (object.total_filled_quantums !== undefined && object.total_filled_quantums !== null) {
      message.totalFilledQuantums = BigInt(object.total_filled_quantums);
    }
    return message;
  },
  toAmino(message: OrderUpdateV1): OrderUpdateV1Amino {
    const obj: any = {};
    obj.order_id = message.orderId ? IndexerOrderId.toAmino(message.orderId) : undefined;
    obj.total_filled_quantums = message.totalFilledQuantums ? message.totalFilledQuantums.toString() : undefined;
    return obj;
  },
  fromAminoMsg(object: OrderUpdateV1AminoMsg): OrderUpdateV1 {
    return OrderUpdateV1.fromAmino(object.value);
  },
  fromProtoMsg(message: OrderUpdateV1ProtoMsg): OrderUpdateV1 {
    return OrderUpdateV1.decode(message.value);
  },
  toProto(message: OrderUpdateV1): Uint8Array {
    return OrderUpdateV1.encode(message).finish();
  },
  toProtoMsg(message: OrderUpdateV1): OrderUpdateV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.off_chain_updates.OrderUpdateV1",
      value: OrderUpdateV1.encode(message).finish()
    };
  }
};
function createBaseOffChainUpdateV1(): OffChainUpdateV1 {
  return {
    orderPlace: undefined,
    orderRemove: undefined,
    orderUpdate: undefined
  };
}
export const OffChainUpdateV1 = {
  typeUrl: "/dydxprotocol.indexer.off_chain_updates.OffChainUpdateV1",
  encode(message: OffChainUpdateV1, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.orderPlace !== undefined) {
      OrderPlaceV1.encode(message.orderPlace, writer.uint32(10).fork()).ldelim();
    }
    if (message.orderRemove !== undefined) {
      OrderRemoveV1.encode(message.orderRemove, writer.uint32(18).fork()).ldelim();
    }
    if (message.orderUpdate !== undefined) {
      OrderUpdateV1.encode(message.orderUpdate, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): OffChainUpdateV1 {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<OffChainUpdateV1>): OffChainUpdateV1 {
    const message = createBaseOffChainUpdateV1();
    message.orderPlace = object.orderPlace !== undefined && object.orderPlace !== null ? OrderPlaceV1.fromPartial(object.orderPlace) : undefined;
    message.orderRemove = object.orderRemove !== undefined && object.orderRemove !== null ? OrderRemoveV1.fromPartial(object.orderRemove) : undefined;
    message.orderUpdate = object.orderUpdate !== undefined && object.orderUpdate !== null ? OrderUpdateV1.fromPartial(object.orderUpdate) : undefined;
    return message;
  },
  fromAmino(object: OffChainUpdateV1Amino): OffChainUpdateV1 {
    const message = createBaseOffChainUpdateV1();
    if (object.order_place !== undefined && object.order_place !== null) {
      message.orderPlace = OrderPlaceV1.fromAmino(object.order_place);
    }
    if (object.order_remove !== undefined && object.order_remove !== null) {
      message.orderRemove = OrderRemoveV1.fromAmino(object.order_remove);
    }
    if (object.order_update !== undefined && object.order_update !== null) {
      message.orderUpdate = OrderUpdateV1.fromAmino(object.order_update);
    }
    return message;
  },
  toAmino(message: OffChainUpdateV1): OffChainUpdateV1Amino {
    const obj: any = {};
    obj.order_place = message.orderPlace ? OrderPlaceV1.toAmino(message.orderPlace) : undefined;
    obj.order_remove = message.orderRemove ? OrderRemoveV1.toAmino(message.orderRemove) : undefined;
    obj.order_update = message.orderUpdate ? OrderUpdateV1.toAmino(message.orderUpdate) : undefined;
    return obj;
  },
  fromAminoMsg(object: OffChainUpdateV1AminoMsg): OffChainUpdateV1 {
    return OffChainUpdateV1.fromAmino(object.value);
  },
  fromProtoMsg(message: OffChainUpdateV1ProtoMsg): OffChainUpdateV1 {
    return OffChainUpdateV1.decode(message.value);
  },
  toProto(message: OffChainUpdateV1): Uint8Array {
    return OffChainUpdateV1.encode(message).finish();
  },
  toProtoMsg(message: OffChainUpdateV1): OffChainUpdateV1ProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.off_chain_updates.OffChainUpdateV1",
      value: OffChainUpdateV1.encode(message).finish()
    };
  }
};