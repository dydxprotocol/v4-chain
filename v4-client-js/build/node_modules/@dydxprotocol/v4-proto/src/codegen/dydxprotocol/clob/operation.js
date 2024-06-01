"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.InternalOperation = exports.Operation = void 0;
const matches_1 = require("./matches");
const tx_1 = require("./tx");
const order_1 = require("./order");
const order_removals_1 = require("./order_removals");
const _m0 = __importStar(require("protobufjs/minimal"));
function createBaseOperation() {
    return {
        match: undefined,
        shortTermOrderPlacement: undefined,
        shortTermOrderCancellation: undefined,
        preexistingStatefulOrder: undefined
    };
}
exports.Operation = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.match !== undefined) {
            matches_1.ClobMatch.encode(message.match, writer.uint32(10).fork()).ldelim();
        }
        if (message.shortTermOrderPlacement !== undefined) {
            tx_1.MsgPlaceOrder.encode(message.shortTermOrderPlacement, writer.uint32(18).fork()).ldelim();
        }
        if (message.shortTermOrderCancellation !== undefined) {
            tx_1.MsgCancelOrder.encode(message.shortTermOrderCancellation, writer.uint32(26).fork()).ldelim();
        }
        if (message.preexistingStatefulOrder !== undefined) {
            order_1.OrderId.encode(message.preexistingStatefulOrder, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseOperation();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.match = matches_1.ClobMatch.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.shortTermOrderPlacement = tx_1.MsgPlaceOrder.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.shortTermOrderCancellation = tx_1.MsgCancelOrder.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.preexistingStatefulOrder = order_1.OrderId.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        const message = createBaseOperation();
        message.match = object.match !== undefined && object.match !== null ? matches_1.ClobMatch.fromPartial(object.match) : undefined;
        message.shortTermOrderPlacement = object.shortTermOrderPlacement !== undefined && object.shortTermOrderPlacement !== null ? tx_1.MsgPlaceOrder.fromPartial(object.shortTermOrderPlacement) : undefined;
        message.shortTermOrderCancellation = object.shortTermOrderCancellation !== undefined && object.shortTermOrderCancellation !== null ? tx_1.MsgCancelOrder.fromPartial(object.shortTermOrderCancellation) : undefined;
        message.preexistingStatefulOrder = object.preexistingStatefulOrder !== undefined && object.preexistingStatefulOrder !== null ? order_1.OrderId.fromPartial(object.preexistingStatefulOrder) : undefined;
        return message;
    }
};
function createBaseInternalOperation() {
    return {
        match: undefined,
        shortTermOrderPlacement: undefined,
        preexistingStatefulOrder: undefined,
        orderRemoval: undefined
    };
}
exports.InternalOperation = {
    encode(message, writer = _m0.Writer.create()) {
        if (message.match !== undefined) {
            matches_1.ClobMatch.encode(message.match, writer.uint32(10).fork()).ldelim();
        }
        if (message.shortTermOrderPlacement !== undefined) {
            tx_1.MsgPlaceOrder.encode(message.shortTermOrderPlacement, writer.uint32(18).fork()).ldelim();
        }
        if (message.preexistingStatefulOrder !== undefined) {
            order_1.OrderId.encode(message.preexistingStatefulOrder, writer.uint32(26).fork()).ldelim();
        }
        if (message.orderRemoval !== undefined) {
            order_removals_1.OrderRemoval.encode(message.orderRemoval, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = createBaseInternalOperation();
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.match = matches_1.ClobMatch.decode(reader, reader.uint32());
                    break;
                case 2:
                    message.shortTermOrderPlacement = tx_1.MsgPlaceOrder.decode(reader, reader.uint32());
                    break;
                case 3:
                    message.preexistingStatefulOrder = order_1.OrderId.decode(reader, reader.uint32());
                    break;
                case 4:
                    message.orderRemoval = order_removals_1.OrderRemoval.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromPartial(object) {
        const message = createBaseInternalOperation();
        message.match = object.match !== undefined && object.match !== null ? matches_1.ClobMatch.fromPartial(object.match) : undefined;
        message.shortTermOrderPlacement = object.shortTermOrderPlacement !== undefined && object.shortTermOrderPlacement !== null ? tx_1.MsgPlaceOrder.fromPartial(object.shortTermOrderPlacement) : undefined;
        message.preexistingStatefulOrder = object.preexistingStatefulOrder !== undefined && object.preexistingStatefulOrder !== null ? order_1.OrderId.fromPartial(object.preexistingStatefulOrder) : undefined;
        message.orderRemoval = object.orderRemoval !== undefined && object.orderRemoval !== null ? order_removals_1.OrderRemoval.fromPartial(object.orderRemoval) : undefined;
        return message;
    }
};
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoib3BlcmF0aW9uLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vLi4vLi4vLi4vLi4vbm9kZV9tb2R1bGVzL0BkeWR4cHJvdG9jb2wvdjQtcHJvdG8vc3JjL2NvZGVnZW4vZHlkeHByb3RvY29sL2Nsb2Ivb3BlcmF0aW9uLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7O0FBQUEsdUNBQXdEO0FBQ3hELDZCQUFrRztBQUNsRyxtQ0FBa0Q7QUFDbEQscURBQXFFO0FBQ3JFLHdEQUEwQztBQStDMUMsU0FBUyxtQkFBbUI7SUFDMUIsT0FBTztRQUNMLEtBQUssRUFBRSxTQUFTO1FBQ2hCLHVCQUF1QixFQUFFLFNBQVM7UUFDbEMsMEJBQTBCLEVBQUUsU0FBUztRQUNyQyx3QkFBd0IsRUFBRSxTQUFTO0tBQ3BDLENBQUM7QUFDSixDQUFDO0FBRVksUUFBQSxTQUFTLEdBQUc7SUFDdkIsTUFBTSxDQUFDLE9BQWtCLEVBQUUsU0FBcUIsR0FBRyxDQUFDLE1BQU0sQ0FBQyxNQUFNLEVBQUU7UUFDakUsSUFBSSxPQUFPLENBQUMsS0FBSyxLQUFLLFNBQVMsRUFBRTtZQUMvQixtQkFBUyxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMsS0FBSyxFQUFFLE1BQU0sQ0FBQyxNQUFNLENBQUMsRUFBRSxDQUFDLENBQUMsSUFBSSxFQUFFLENBQUMsQ0FBQyxNQUFNLEVBQUUsQ0FBQztTQUNwRTtRQUVELElBQUksT0FBTyxDQUFDLHVCQUF1QixLQUFLLFNBQVMsRUFBRTtZQUNqRCxrQkFBYSxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMsdUJBQXVCLEVBQUUsTUFBTSxDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUMsQ0FBQyxJQUFJLEVBQUUsQ0FBQyxDQUFDLE1BQU0sRUFBRSxDQUFDO1NBQzFGO1FBRUQsSUFBSSxPQUFPLENBQUMsMEJBQTBCLEtBQUssU0FBUyxFQUFFO1lBQ3BELG1CQUFjLENBQUMsTUFBTSxD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