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
export declare const Operation: {
    encode(message: Operation, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Operation;
    fromPartial(object: DeepPartial<Operation>): Operation;
};
export declare const InternalOperation: {
    encode(message: InternalOperation, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): InternalOperation;
    fromPartial(object: DeepPartial<InternalOperation>): InternalOperation;
};
