import { Order, OrderSDKType, OrderId, OrderIdSDKType } from "./order";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * ProcessProposerMatchesEvents is used for communicating which events occurred
 * in the last block that require updating the state of the memclob in the
 * Commit blocker. It contains information about the following state updates:
 * - Stateful orders that were placed in the last block.
 * - Stateful order IDs that were expired in the last block.
 * - Order IDs that were filled in the last block.
 * - Stateful cancellations that were placed in the last block.
 * - Stateful order removals in the last block.
 * - Conditional order ids triggered in the last block.
 * - The height of the block in which the events occurred.
 */

export interface ProcessProposerMatchesEvents {
  placedStatefulOrders: Order[];
  expiredStatefulOrderIds: OrderId[];
  ordersIdsFilledInLastBlock: OrderId[];
  placedStatefulCancellations: OrderId[];
  removedStatefulOrderIds: OrderId[];
  conditionalOrderIdsTriggeredInLastBlock: OrderId[];
  blockHeight: number;
}
/**
 * ProcessProposerMatchesEvents is used for communicating which events occurred
 * in the last block that require updating the state of the memclob in the
 * Commit blocker. It contains information about the following state updates:
 * - Stateful orders that were placed in the last block.
 * - Stateful order IDs that were expired in the last block.
 * - Order IDs that were filled in the last block.
 * - Stateful cancellations that were placed in the last block.
 * - Stateful order removals in the last block.
 * - Conditional order ids triggered in the last block.
 * - The height of the block in which the events occurred.
 */

export interface ProcessProposerMatchesEventsSDKType {
  placed_stateful_orders: OrderSDKType[];
  expired_stateful_order_ids: OrderIdSDKType[];
  orders_ids_filled_in_last_block: OrderIdSDKType[];
  placed_stateful_cancellations: OrderIdSDKType[];
  removed_stateful_order_ids: OrderIdSDKType[];
  conditional_order_ids_triggered_in_last_block: OrderIdSDKType[];
  block_height: number;
}

function createBaseProcessProposerMatchesEvents(): ProcessProposerMatchesEvents {
  return {
    placedStatefulOrders: [],
    expiredStatefulOrderIds: [],
    ordersIdsFilledInLastBlock: [],
    placedStatefulCancellations: [],
    removedStatefulOrderIds: [],
    conditionalOrderIdsTriggeredInLastBlock: [],
    blockHeight: 0
  };
}

export const ProcessProposerMatchesEvents = {
  encode(message: ProcessProposerMatchesEvents, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.placedStatefulOrders) {
      Order.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    for (const v of message.expiredStatefulOrderIds) {
      OrderId.encode(v!, writer.uint32(18).fork()).ldelim();
    }

    for (const v of message.ordersIdsFilledInLastBlock) {
      OrderId.encode(v!, writer.uint32(26).fork()).ldelim();
    }

    for (const v of message.placedStatefulCancellations) {
      OrderId.encode(v!, writer.uint32(34).fork()).ldelim();
    }

    for (const v of message.removedStatefulOrderIds) {
      OrderId.encode(v!, writer.uint32(42).fork()).ldelim();
    }

    for (const v of message.conditionalOrderIdsTriggeredInLastBlock) {
      OrderId.encode(v!, writer.uint32(50).fork()).ldelim();
    }

    if (message.blockHeight !== 0) {
      writer.uint32(56).uint32(message.blockHeight);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProcessProposerMatchesEvents {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProcessProposerMatchesEvents();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.placedStatefulOrders.push(Order.decode(reader, reader.uint32()));
          break;

        case 2:
          message.expiredStatefulOrderIds.push(OrderId.decode(reader, reader.uint32()));
          break;

        case 3:
          message.ordersIdsFilledInLastBlock.push(OrderId.decode(reader, reader.uint32()));
          break;

        case 4:
          message.placedStatefulCancellations.push(OrderId.decode(reader, reader.uint32()));
          break;

        case 5:
          message.removedStatefulOrderIds.push(OrderId.decode(reader, reader.uint32()));
          break;

        case 6:
          message.conditionalOrderIdsTriggeredInLastBlock.push(OrderId.decode(reader, reader.uint32()));
          break;

        case 7:
          message.blockHeight = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ProcessProposerMatchesEvents>): ProcessProposerMatchesEvents {
    const message = createBaseProcessProposerMatchesEvents();
    message.placedStatefulOrders = object.placedStatefulOrders?.map(e => Order.fromPartial(e)) || [];
    message.expiredStatefulOrderIds = object.expiredStatefulOrderIds?.map(e => OrderId.fromPartial(e)) || [];
    message.ordersIdsFilledInLastBlock = object.ordersIdsFilledInLastBlock?.map(e => OrderId.fromPartial(e)) || [];
    message.placedStatefulCancellations = object.placedStatefulCancellations?.map(e => OrderId.fromPartial(e)) || [];
    message.removedStatefulOrderIds = object.removedStatefulOrderIds?.map(e => OrderId.fromPartial(e)) || [];
    message.conditionalOrderIdsTriggeredInLastBlock = object.conditionalOrderIdsTriggeredInLastBlock?.map(e => OrderId.fromPartial(e)) || [];
    message.blockHeight = object.blockHeight ?? 0;
    return message;
  }

};