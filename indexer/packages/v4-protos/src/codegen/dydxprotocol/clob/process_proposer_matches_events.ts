import { OrderId, OrderIdSDKType } from "./order";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
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
  /** @deprecated */
  placedLongTermOrderIds: OrderId[];
  expiredStatefulOrderIds: OrderId[];
  orderIdsFilledInLastBlock: OrderId[];
  /** @deprecated */

  placedStatefulCancellationOrderIds: OrderId[];
  removedStatefulOrderIds: OrderId[];
  conditionalOrderIdsTriggeredInLastBlock: OrderId[];
  /** @deprecated */

  placedConditionalOrderIds: OrderId[];
  blockHeight: number;
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
  /** @deprecated */
  placed_long_term_order_ids: OrderIdSDKType[];
  expired_stateful_order_ids: OrderIdSDKType[];
  order_ids_filled_in_last_block: OrderIdSDKType[];
  /** @deprecated */

  placed_stateful_cancellation_order_ids: OrderIdSDKType[];
  removed_stateful_order_ids: OrderIdSDKType[];
  conditional_order_ids_triggered_in_last_block: OrderIdSDKType[];
  /** @deprecated */

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
  encode(message: ProcessProposerMatchesEvents, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): ProcessProposerMatchesEvents {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<ProcessProposerMatchesEvents>): ProcessProposerMatchesEvents {
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
  }

};