import { ClobMatch, ClobMatchAmino, ClobMatchSDKType } from "./matches";
import { MsgPlaceOrder, MsgPlaceOrderAmino, MsgPlaceOrderSDKType, MsgCancelOrder, MsgCancelOrderAmino, MsgCancelOrderSDKType } from "./tx";
import { OrderId, OrderIdAmino, OrderIdSDKType } from "./order";
import { OrderRemoval, OrderRemovalAmino, OrderRemovalSDKType } from "./order_removals";
import { BinaryReader, BinaryWriter } from "../../binary";
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
export interface OperationProtoMsg {
  typeUrl: "/dydxprotocol.clob.Operation";
  value: Uint8Array;
}
/**
 * Operation represents an operation in the proposed operations. Operation is
 * used internally within the memclob only.
 */
export interface OperationAmino {
  match?: ClobMatchAmino;
  short_term_order_placement?: MsgPlaceOrderAmino;
  short_term_order_cancellation?: MsgCancelOrderAmino;
  preexisting_stateful_order?: OrderIdAmino;
}
export interface OperationAminoMsg {
  type: "/dydxprotocol.clob.Operation";
  value: OperationAmino;
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
export interface InternalOperationProtoMsg {
  typeUrl: "/dydxprotocol.clob.InternalOperation";
  value: Uint8Array;
}
/**
 * InternalOperation represents an internal operation in the operations to
 * propose. InternalOperation is used internally within the memclob only.
 */
export interface InternalOperationAmino {
  match?: ClobMatchAmino;
  short_term_order_placement?: MsgPlaceOrderAmino;
  preexisting_stateful_order?: OrderIdAmino;
  order_removal?: OrderRemovalAmino;
}
export interface InternalOperationAminoMsg {
  type: "/dydxprotocol.clob.InternalOperation";
  value: InternalOperationAmino;
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
  typeUrl: "/dydxprotocol.clob.Operation",
  encode(message: Operation, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): Operation {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<Operation>): Operation {
    const message = createBaseOperation();
    message.match = object.match !== undefined && object.match !== null ? ClobMatch.fromPartial(object.match) : undefined;
    message.shortTermOrderPlacement = object.shortTermOrderPlacement !== undefined && object.shortTermOrderPlacement !== null ? MsgPlaceOrder.fromPartial(object.shortTermOrderPlacement) : undefined;
    message.shortTermOrderCancellation = object.shortTermOrderCancellation !== undefined && object.shortTermOrderCancellation !== null ? MsgCancelOrder.fromPartial(object.shortTermOrderCancellation) : undefined;
    message.preexistingStatefulOrder = object.preexistingStatefulOrder !== undefined && object.preexistingStatefulOrder !== null ? OrderId.fromPartial(object.preexistingStatefulOrder) : undefined;
    return message;
  },
  fromAmino(object: OperationAmino): Operation {
    const message = createBaseOperation();
    if (object.match !== undefined && object.match !== null) {
      message.match = ClobMatch.fromAmino(object.match);
    }
    if (object.short_term_order_placement !== undefined && object.short_term_order_placement !== null) {
      message.shortTermOrderPlacement = MsgPlaceOrder.fromAmino(object.short_term_order_placement);
    }
    if (object.short_term_order_cancellation !== undefined && object.short_term_order_cancellation !== null) {
      message.shortTermOrderCancellation = MsgCancelOrder.fromAmino(object.short_term_order_cancellation);
    }
    if (object.preexisting_stateful_order !== undefined && object.preexisting_stateful_order !== null) {
      message.preexistingStatefulOrder = OrderId.fromAmino(object.preexisting_stateful_order);
    }
    return message;
  },
  toAmino(message: Operation): OperationAmino {
    const obj: any = {};
    obj.match = message.match ? ClobMatch.toAmino(message.match) : undefined;
    obj.short_term_order_placement = message.shortTermOrderPlacement ? MsgPlaceOrder.toAmino(message.shortTermOrderPlacement) : undefined;
    obj.short_term_order_cancellation = message.shortTermOrderCancellation ? MsgCancelOrder.toAmino(message.shortTermOrderCancellation) : undefined;
    obj.preexisting_stateful_order = message.preexistingStatefulOrder ? OrderId.toAmino(message.preexistingStatefulOrder) : undefined;
    return obj;
  },
  fromAminoMsg(object: OperationAminoMsg): Operation {
    return Operation.fromAmino(object.value);
  },
  fromProtoMsg(message: OperationProtoMsg): Operation {
    return Operation.decode(message.value);
  },
  toProto(message: Operation): Uint8Array {
    return Operation.encode(message).finish();
  },
  toProtoMsg(message: Operation): OperationProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.Operation",
      value: Operation.encode(message).finish()
    };
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
  typeUrl: "/dydxprotocol.clob.InternalOperation",
  encode(message: InternalOperation, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
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
  decode(input: BinaryReader | Uint8Array, length?: number): InternalOperation {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<InternalOperation>): InternalOperation {
    const message = createBaseInternalOperation();
    message.match = object.match !== undefined && object.match !== null ? ClobMatch.fromPartial(object.match) : undefined;
    message.shortTermOrderPlacement = object.shortTermOrderPlacement !== undefined && object.shortTermOrderPlacement !== null ? MsgPlaceOrder.fromPartial(object.shortTermOrderPlacement) : undefined;
    message.preexistingStatefulOrder = object.preexistingStatefulOrder !== undefined && object.preexistingStatefulOrder !== null ? OrderId.fromPartial(object.preexistingStatefulOrder) : undefined;
    message.orderRemoval = object.orderRemoval !== undefined && object.orderRemoval !== null ? OrderRemoval.fromPartial(object.orderRemoval) : undefined;
    return message;
  },
  fromAmino(object: InternalOperationAmino): InternalOperation {
    const message = createBaseInternalOperation();
    if (object.match !== undefined && object.match !== null) {
      message.match = ClobMatch.fromAmino(object.match);
    }
    if (object.short_term_order_placement !== undefined && object.short_term_order_placement !== null) {
      message.shortTermOrderPlacement = MsgPlaceOrder.fromAmino(object.short_term_order_placement);
    }
    if (object.preexisting_stateful_order !== undefined && object.preexisting_stateful_order !== null) {
      message.preexistingStatefulOrder = OrderId.fromAmino(object.preexisting_stateful_order);
    }
    if (object.order_removal !== undefined && object.order_removal !== null) {
      message.orderRemoval = OrderRemoval.fromAmino(object.order_removal);
    }
    return message;
  },
  toAmino(message: InternalOperation): InternalOperationAmino {
    const obj: any = {};
    obj.match = message.match ? ClobMatch.toAmino(message.match) : undefined;
    obj.short_term_order_placement = message.shortTermOrderPlacement ? MsgPlaceOrder.toAmino(message.shortTermOrderPlacement) : undefined;
    obj.preexisting_stateful_order = message.preexistingStatefulOrder ? OrderId.toAmino(message.preexistingStatefulOrder) : undefined;
    obj.order_removal = message.orderRemoval ? OrderRemoval.toAmino(message.orderRemoval) : undefined;
    return obj;
  },
  fromAminoMsg(object: InternalOperationAminoMsg): InternalOperation {
    return InternalOperation.fromAmino(object.value);
  },
  fromProtoMsg(message: InternalOperationProtoMsg): InternalOperation {
    return InternalOperation.decode(message.value);
  },
  toProto(message: InternalOperation): Uint8Array {
    return InternalOperation.encode(message).finish();
  },
  toProtoMsg(message: InternalOperation): InternalOperationProtoMsg {
    return {
      typeUrl: "/dydxprotocol.clob.InternalOperation",
      value: InternalOperation.encode(message).finish()
    };
  }
};