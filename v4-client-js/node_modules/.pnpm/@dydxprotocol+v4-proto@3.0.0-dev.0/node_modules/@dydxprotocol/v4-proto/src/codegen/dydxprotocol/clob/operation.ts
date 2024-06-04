import { ClobMatch, ClobMatchSDKType } from "./matches";
import { MsgPlaceOrder, MsgPlaceOrderSDKType, MsgCancelOrder, MsgCancelOrderSDKType } from "./tx";
import { OrderId, OrderIdSDKType } from "./order";
import { OrderRemoval, OrderRemovalSDKType } from "./order_removals";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * Operation represents an operation in the proposed operations. Operation is
 * used internally within the memclob only.
 */

export interface Operation {
  match?: ClobMatch;
  shortTermOrderPlacement?: MsgPlaceOrder;
  shortTermOrderCancellation?: MsgCancelOrder;
  preexistingStatefulOrder?: OrderId;
}
/**
 * Operation represents an operation in the proposed operations. Operation is
 * used internally within the memclob only.
 */

export interface OperationSDKType {
  match?: ClobMatchSDKType;
  short_term_order_placement?: MsgPlaceOrderSDKType;
  short_term_order_cancellation?: MsgCancelOrderSDKType;
  preexisting_stateful_order?: OrderIdSDKType;
}
/**
 * InternalOperation represents an internal operation in the operations to
 * propose. InternalOperation is used internally within the memclob only.
 */

export interface InternalOperation {
  match?: ClobMatch;
  shortTermOrderPlacement?: MsgPlaceOrder;
  preexistingStatefulOrder?: OrderId;
  orderRemoval?: OrderRemoval;
}
/**
 * InternalOperation represents an internal operation in the operations to
 * propose. InternalOperation is used internally within the memclob only.
 */

export interface InternalOperationSDKType {
  match?: ClobMatchSDKType;
  short_term_order_placement?: MsgPlaceOrderSDKType;
  preexisting_stateful_order?: OrderIdSDKType;
  order_removal?: OrderRemovalSDKType;
}

function createBaseOperation(): Operation {
  return {
    match: undefined,
    shortTermOrderPlacement: undefined,
    shortTermOrderCancellation: undefined,
    preexistingStatefulOrder: undefined
  };
}

export const Operation = {
  encode(message: Operation, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.match !== undefined) {
      ClobMatch.encode(message.match, writer.uint32(10).fork()).ldelim();
    }

    if (message.shortTermOrderPlacement !== undefined) {
      MsgPlaceOrder.encode(message.shortTermOrderPlacement, writer.uint32(18).fork()).ldelim();
    }

    if (message.shortTermOrderCancellation !== undefined) {
      MsgCancelOrder.encode(message.shortTermOrderCancellation, writer.uint32(26).fork()).ldelim();
    }

    if (message.preexistingStatefulOrder !== undefined) {
      OrderId.encode(message.preexistingStatefulOrder, writer.uint32(34).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Operation {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOperation();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.match = ClobMatch.decode(reader, reader.uint32());
          break;

        case 2:
          message.shortTermOrderPlacement = MsgPlaceOrder.decode(reader, reader.uint32());
          break;

        case 3:
          message.shortTermOrderCancellation = MsgCancelOrder.decode(reader, reader.uint32());
          break;

        case 4:
          message.preexistingStatefulOrder = OrderId.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<Operation>): Operation {
    const message = createBaseOperation();
    message.match = object.match !== undefined && object.match !== null ? ClobMatch.fromPartial(object.match) : undefined;
    message.shortTermOrderPlacement = object.shortTermOrderPlacement !== undefined && object.shortTermOrderPlacement !== null ? MsgPlaceOrder.fromPartial(object.shortTermOrderPlacement) : undefined;
    message.shortTermOrderCancellation = object.shortTermOrderCancellation !== undefined && object.shortTermOrderCancellation !== null ? MsgCancelOrder.fromPartial(object.shortTermOrderCancellation) : undefined;
    message.preexistingStatefulOrder = object.preexistingStatefulOrder !== undefined && object.preexistingStatefulOrder !== null ? OrderId.fromPartial(object.preexistingStatefulOrder) : undefined;
    return message;
  }

};

function createBaseInternalOperation(): InternalOperation {
  return {
    match: undefined,
    shortTermOrderPlacement: undefined,
    preexistingStatefulOrder: undefined,
    orderRemoval: undefined
  };
}

export const InternalOperation = {
  encode(message: InternalOperation, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.match !== undefined) {
      ClobMatch.encode(message.match, writer.uint32(10).fork()).ldelim();
    }

    if (message.shortTermOrderPlacement !== undefined) {
      MsgPlaceOrder.encode(message.shortTermOrderPlacement, writer.uint32(18).fork()).ldelim();
    }

    if (message.preexistingStatefulOrder !== undefined) {
      OrderId.encode(message.preexistingStatefulOrder, writer.uint32(26).fork()).ldelim();
    }

    if (message.orderRemoval !== undefined) {
      OrderRemoval.encode(message.orderRemoval, writer.uint32(34).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InternalOperation {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInternalOperation();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.match = ClobMatch.decode(reader, reader.uint32());
          break;

        case 2:
          message.shortTermOrderPlacement = MsgPlaceOrder.decode(reader, reader.uint32());
          break;

        case 3:
          message.preexistingStatefulOrder = OrderId.decode(reader, reader.uint32());
          break;

        case 4:
          message.orderRemoval = OrderRemoval.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<InternalOperation>): InternalOperation {
    const message = createBaseInternalOperation();
    message.match = object.match !== undefined && object.match !== null ? ClobMatch.fromPartial(object.match) : undefined;
    message.shortTermOrderPlacement = object.shortTermOrderPlacement !== undefined && object.shortTermOrderPlacement !== null ? MsgPlaceOrder.fromPartial(object.shortTermOrderPlacement) : undefined;
    message.preexistingStatefulOrder = object.preexistingStatefulOrder !== undefined && object.preexistingStatefulOrder !== null ? OrderId.fromPartial(object.preexistingStatefulOrder) : undefined;
    message.orderRemoval = object.orderRemoval !== undefined && object.orderRemoval !== null ? OrderRemoval.fromPartial(object.orderRemoval) : undefined;
    return message;
  }

};