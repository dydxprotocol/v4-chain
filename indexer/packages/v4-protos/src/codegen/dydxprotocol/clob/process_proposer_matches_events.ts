import { OrderId, OrderIdAmino, OrderIdSDKType } from "./order";
import { BinaryReader, BinaryWriter } from "../../binary";
/**
 * ProcessProposerMatchesEvents is used for communicating which events occurred
 * in the last block that require updating the state of the memclob in the
 * Commit blocker. It contains information about the following state updates:
 * - Long term order IDs that were placed in the last block.
 * - Stateful order IDs that were expired in the last block.
 * - Order IDs that were filled in the last block.
 * - Stateful cancellations order IDs that were placed in the last block.
 * - Stateful order IDs forcefully removed in the last block.
 * - Conditional order IDs triggered in the last block.
 * - Conditional order IDs placed, but not triggered in the last block.
 * - The height of the block in which the events occurred.
 */
export interface ProcessProposerMatchesEvents {
  placedLongTermOrderIds: OrderId[];
  expiredStatefulOrderIds: OrderId[];
  orderIdsFilledInLastBlock: OrderId[];
  placedStatefulCancellationOrderIds: OrderId[];
  removedStatefulOrderIds: OrderId[];
  conditionalOrderIdsTriggeredInLastBlock: OrderId[];
  placedConditionalOrderIds: OrderId[];
  blockHeight: number;
}
export interface ProcessProposerMatchesEventsProtoMsg {
  typeUrl: "/dydxprotocol.clob.ProcessProposerMatchesEvents";
  value: Uint8Array;
}
/**
 * ProcessProposerMatchesEvents is used for communicating which events occurred
 * in the last block that require updating the state of the memclob in the
 * Commit blocker. It contains information about the following state updates:
 * - Long term order IDs that were placed in the last block.
 * - Stateful order IDs that were expired in the last block.
 * - Order IDs that were filled in the last block.
 * - Stateful cancellations order IDs that were placed in the last block.
 * - Stateful order IDs forcefully removed in the last block.
 * - Conditional order IDs triggered in the last block.
 * - Conditional order IDs placed, but not triggered in the last block.
 * - The height of the block in which the events occurred.
 */
export interface ProcessProposerMatchesEventsAmino {
  placed_long_term_order_ids?: OrderIdAmino[];
  expired_stateful_order_ids?: OrderIdAmino[];
  order_ids_filled_in_last_block?: OrderIdAmino[];
  placed_stateful_cancellation_order_ids?: OrderIdAmino[];
  removed_stateful_order_ids?: OrderIdAmino[];
  conditional_order_ids_triggered_in_last_block?: OrderIdAmino[];
  placed_conditional_order_ids?: OrderIdAmino[];
  block_height?: number;
}
export interface ProcessProposerMatchesEventsAminoMsg {
  type: "/dydxprotocol.clob.ProcessProposerMatchesEvents";
  value: ProcessProposerMatchesEventsAmino;
}
/**
 * ProcessProposerMatchesEvents is used for communicating which events occurred
 * in the last block that require updating the state of the memclob in the
 * Commit blocker. It contains information about the following state updates:
 * - Long term order IDs that were placed in the last block.
 * - Stateful order IDs that were expired in the last block.
 * - Order IDs that were filled in the last block.
 * - Stateful cancellations order IDs that were placed in the last block.
 * - Stateful order IDs forcefully removed in the last block.
 * - Conditional order IDs triggered in the last block.
 * - Conditional order IDs placed, but not triggered in the last block.
 * - The height of the block in which the events occurred.
 */
export interface ProcessProposerMatchesEventsSDKType {
  placed_long_term_order_ids: OrderIdSDKType[];
  expired_stateful_order_ids: OrderIdSDKType[];
  order_ids_filled_in_last_block: OrderIdSDKType[];
  placed_stateful_cancellation_order_ids: OrderIdSDKType[];
  removed_stateful_order_ids: OrderIdSDKType[];
  conditional_order_ids_triggered_in_last_block: OrderIdSDKType[];
  placed_conditional_order_ids: OrderIdSDKType[];
  block_height: number;
}
function createBaseProcessProposerMatchesEvents(): ProcessProposerMatchesEvents {
  return {
    placedLongTermOrderIds: [],
    expiredStatefulOrderIds: [],
    orderIdsFilledInLastBlock: [],
    placedStatefulCancellationOrderIds: [],
    removedStatefulOrderIds: [],
    conditionalOrderIdsTriggeredInLastBlock: [],
    placedConditionalOrderIds: [],
    blockHeight: 0
  };
}
export const ProcessProposerMatchesEvents = {
  typeUrl: "/dydxprotocol.clob.ProcessProposerMatchesEvents",
  encode(message: ProcessProposerMatchesEvents, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.placedLongTermOrderIds) {
      OrderId.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.expiredStatefulOrderIds) {
      OrderId.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.orderIdsFilledInLastBlock) {
      OrderId.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.placedStatefulCancellationOrderIds) {
      OrderId.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.removedStatefulOrderIds) {
      OrderId.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.conditionalOrderIdsTriggeredInLastBlock) {
      OrderId.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.placedConditionalOrderIds) {
      OrderId.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    if (message.blockHeight !== 0) {
      writer.uint32(64).uint32(message.blockHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): ProcessProposerMatchesEvents {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProcessProposerMatchesEvents();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.placedLongTermOrderIds.push(OrderId.decode(reader, reader.uint32()));
          break;
        case 2:
          message.expiredStatefulOrderIds.push(OrderId.decode(reader, reader.uint32()));
          break;
        case 3:
          message.orderIdsFilledInLastBlock.push(OrderId.decode(reader, reader.uint32()));
          break;
        case 4:
          message.placedStatefulCancellationOrderIds.push(OrderId.decode(reader, reader.uint32()));
          break;
        case 5:
          message.removedStatefulOrderIds.push(OrderId.decode(reader, reader.uint32()));
          break;
        case 6:
          message.conditionalOrderIdsTriggeredInLastBlock.push(OrderId.decode(reader, reader.uint32()));
          break;
        case 7:
          message.placedConditionalOrderIds.push(OrderId.decode(reader, reader.uint32()));
          break;
        case 8:
          message.blockHeight = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<ProcessProposerMatchesEvents>): ProcessProposerMatchesEvents {
    const message = createBaseProcessProposerMatchesEvents();
    message.placedLongTermOrderIds = object.placedLongTermOrderIds?.map(e => OrderId.fromPartial(e)) || [];
    message.expiredStatefulOrderIds = object.expiredStatefulOrderIds?.map(e => OrderId.fromPartial(e)) || [];
    message.orderIdsFilledInLastBlock = object.orderIdsFilledInLastBlock?.map(e => OrderId.fromPartial(e)) || [];
    message.placedStatefulCancellationOrderIds = object.placedStatefulCancellationOrderIds?.map(e => OrderId.fromPartial(e)) || [];
    message.removedStatefulOrderIds = object.removedStatefulOrderIds?.map(e => OrderId.fromPartial(e)) || [];
    message.conditionalOrderIdsTriggeredInLastBlock = object.conditionalOrderIdsTriggeredInLastBlock?.map(e => OrderId.fromPartial(e)) || [];
    message.placedConditionalOrderIds = object.placedConditionalOrderIds?.map(e => OrderId.fromPartial(e)) || [];
    message.blockHeight = object.blockHeight ?? 0;
    return message;
  },
  fromAmino(object: ProcessProposerMatchesEventsAmino): ProcessProposerMatchesEvents {
    const message = createBaseProcessProposerMatchesEvents();
    message.placedLongTermOrderIds = object.placed_long_term_order_ids?.map(e => OrderId.fromAmino(e)) || [];
    message.expiredStatefulOrderIds = object.expired_stateful_order_ids?.map(e => OrderId.fromAmino(e)) || [];
    message.orderIdsFilledInLastBlock = object.order_ids_filled_in_last_block?.map(e => OrderId.fromAmino(e)) || [];
    message.placedStatefulCancellationOrderIds = object.placed_stateful_cancellation_order_ids?.map(e => OrderId.fromAmino(e)) || [];
    message.removedStatefulOrderIds = object.removed_stateful_order_ids?.map(e => OrderId.fromAmino(e)) || [];
    message.conditionalOrderIdsTriggeredInLastBlock = object.conditional_order_ids_triggered_in_last_block?.map(e => OrderId.fromAmino(e)) || [];
    message.placedConditionalOrderIds = object.placed_conditional_order_ids?.map(e => OrderId.fromAmino(e)) || [];
    if (object.block_height !== undefined && object.block_height !== null) {
      message.blockHeight = object.block_height;
    }
    return message;
  },
  toAmino(message: ProcessProposerMatchesEvents): ProcessProposerMatchesEventsAmino {
    const obj: any = {};
    if (message.placedLongTermOrderIds) {
      obj.placed_long_term_order_ids = message.placedLongTermOrderIds.map(e => e ? OrderId.toAmino(e) : undefined);
    } else {
      obj.placed_long_term_order_ids = [];
    }
    if (message.expiredStatefulOrderIds) {
      obj.expired_stateful_order_ids = message.expiredStatefulOrderIds.map(e => e ? OrderId.toAmino(e) : undefined);
    } else {
      obj.expired_stateful_order_ids = [];
    }
    if (message.orderIdsFilledInLastBlock) {
      obj.order_ids_filled_in_last_block = message.orderIdsFilledInLastBlock.map(e => e ? OrderId.toAmino(e) : undefined);
    } else {
      obj.order_ids_filled_in_last_block = [];
    }
    if (message.placedStatefulCancellationOrderIds) {
      obj.placed_stateful_cancellation_order_ids = message.placedStatefulCancellationOrderIds.map(e => e ? OrderId.toAmino(e) : undefined);
    } else {
      obj.placed_stateful_cancellation_order_ids = [];
    }
    if (message.removedStatefulOrderIds) {
      obj.removed_stateful_order_ids = message.removedStatefulOrderIds.map(e => e ? OrderId.toAmino(e) : undefined);
    } else {
      obj.removed_stateful_order_ids = [];
    }
    if (message.conditionalOrderIdsTriggeredInLastBlock) {
      obj.conditional_order_ids_triggered_in_last_block = message.conditionalOrderIdsTriggeredInLastBlock.map(e => e ? OrderId.toAmino(e) : undefined);
    } else {
      obj.conditional_order_ids_triggered_in_last_block = [];
    }
    if (message.placedConditionalOrderIds) {
      obj.placed_conditional_order_ids = message.placedConditionalOrderIds.map(e => e ? OrderId.toAmino(e) : undefined);
    } else {
      obj.placed_conditional_order_ids = [];
    }
    obj.block_height = message.blockHeight;
    return obj;
  },
  fromAminoMsg(object: ProcessProposerMatchesEventsAminoMsg): ProcessProposerMatchesEvents {
    return ProcessProposerMatchesEvents.fromAmino(object.value);
  },
  fromProtoMsg(message: ProcessProposerMatchesEventsProtoMsg): ProcessProposerMatchesEvents {
    return ProcessProposerMatchesEvents.decode(message.value);
  },
  toProto(message: ProcessProposerMatchesEvents): Uint8Array {
    return ProcessProposerMatchesEvents.encode(message).finish();
  },
  toProtoMsg(message: ProcessProposerMatchesEvents): ProcessProposerMatchesEventsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.ProcessProposerMatchesEvents",
      value: ProcessProposerMatchesEvents.encode(message).finish()
    };
  }
};