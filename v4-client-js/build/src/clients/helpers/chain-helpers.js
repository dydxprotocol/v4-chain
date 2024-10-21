"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.round = round;
exports.calculateQuantums = calculateQuantums;
exports.calculateSubticks = calculateSubticks;
exports.calculateSide = calculateSide;
exports.calculateTimeInForce = calculateTimeInForce;
exports.calculateOrderFlags = calculateOrderFlags;
exports.calculateClientMetadata = calculateClientMetadata;
exports.calculateConditionType = calculateConditionType;
exports.calculateConditionalOrderTriggerSubticks = calculateConditionalOrderTriggerSubticks;
const long_1 = __importDefault(require("long"));
const constants_1 = require("../constants");
const proto_includes_1 = require("../modules/proto-includes");
const types_1 = require("../types");
function round(input, base) {
    return Math.floor(input / base) * base;
}
function calculateQuantums(size, atomicResolution, stepBaseQuantums) {
    const rawQuantums = size * 10 ** (-1 * atomicResolution);
    const quantums = round(rawQuantums, stepBaseQuantums);
    // stepBaseQuantums functions as minimum order size
    const result = Math.max(quantums, stepBaseQuantums);
    return long_1.default.fromNumber(result);
}
function calculateSubticks(price, atomicResolution, quantumConversionExponent, subticksPerTick) {
    const QUOTE_QUANTUMS_ATOMIC_RESOLUTION = -6;
    const exponent = atomicResolution - quantumConversionExponent - QUOTE_QUANTUMS_ATOMIC_RESOLUTION;
    const rawSubticks = price * 10 ** exponent;
    const subticks = round(rawSubticks, subticksPerTick);
    const result = Math.max(subticks, subticksPerTick);
    return long_1.default.fromNumber(result);
}
function calculateSide(side) {
    return side === constants_1.OrderSide.BUY ? proto_includes_1.Order_Side.SIDE_BUY : proto_includes_1.Order_Side.SIDE_SELL;
}
function calculateTimeInForce(type, timeInForce, execution, postOnly) {
    switch (type) {
        case constants_1.OrderType.MARKET:
            switch (timeInForce) {
                case constants_1.OrderTimeInForce.IOC:
                    return proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_IOC;
                default:
                    return proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL;
            }
        case constants_1.OrderType.LIMIT:
            switch (timeInForce) {
                case constants_1.OrderTimeInForce.GTT:
                    if (postOnly == null) {
                        throw new Error('postOnly must be set if order type is LIMIT and timeInForce is GTT');
                    }
                    return postOnly
                        ? proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_POST_ONLY
                        : proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED;
                case constants_1.OrderTimeInForce.FOK:
                    return proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL;
                case constants_1.OrderTimeInForce.IOC:
                    return proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_IOC;
                default:
                    throw new Error('Unexpected code path: timeInForce');
            }
        case constants_1.OrderType.STOP_LIMIT:
        case constants_1.OrderType.TAKE_PROFIT_LIMIT:
            if (execution == null) {
                throw new Error('execution must be set if order type is STOP_LIMIT or TAKE_PROFIT_LIMIT');
            }
            switch (execution) {
                case constants_1.OrderExecution.DEFAULT:
                    return proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_UNSPECIFIED;
                case constants_1.OrderExecution.POST_ONLY:
                    return proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_POST_ONLY;
                case constants_1.OrderExecution.FOK:
                    return proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL;
                case constants_1.OrderExecution.IOC:
                    return proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_IOC;
                default:
                    throw new Error('Unexpected code path: timeInForce');
            }
        case constants_1.OrderType.STOP_MARKET:
        case constants_1.OrderType.TAKE_PROFIT_MARKET:
            if (execution == null) {
                throw new Error('execution must be set if order type is STOP_MARKET or TAKE_PROFIT_MARKET');
            }
            switch (execution) {
                case constants_1.OrderExecution.DEFAULT:
                    throw new Error('Execution value DEFAULT not supported for STOP_MARKET or TAKE_PROFIT_MARKET');
                case constants_1.OrderExecution.POST_ONLY:
                    throw new Error('Execution value POST_ONLY not supported for STOP_MARKET or TAKE_PROFIT_MARKET');
                case constants_1.OrderExecution.FOK:
                    return proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_FILL_OR_KILL;
                case constants_1.OrderExecution.IOC:
                    return proto_includes_1.Order_TimeInForce.TIME_IN_FORCE_IOC;
                default:
                    throw new Error('Unexpected code path: timeInForce');
            }
        default:
            throw new Error('Unexpected code path: timeInForce');
    }
}
function calculateOrderFlags(type, timeInForce) {
    switch (type) {
        case constants_1.OrderType.MARKET:
            return types_1.OrderFlags.SHORT_TERM;
        case constants_1.OrderType.LIMIT:
            if (timeInForce === undefined) {
                throw new Error('timeInForce must be set if orderType is LIMIT');
            }
            if (timeInForce === constants_1.OrderTimeInForce.GTT) {
                return types_1.OrderFlags.LONG_TERM;
            }
            else {
                return types_1.OrderFlags.SHORT_TERM;
            }
        case constants_1.OrderType.STOP_LIMIT:
        case constants_1.OrderType.TAKE_PROFIT_LIMIT:
        case constants_1.OrderType.STOP_MARKET:
        case constants_1.OrderType.TAKE_PROFIT_MARKET:
            return types_1.OrderFlags.CONDITIONAL;
        default:
            throw new Error('Unexpected code path: orderFlags');
    }
}
function calculateClientMetadata(orderType) {
    switch (orderType) {
        case constants_1.OrderType.MARKET:
        case constants_1.OrderType.STOP_MARKET:
        case constants_1.OrderType.TAKE_PROFIT_MARKET:
            return 1;
        default:
            return 0;
    }
}
function calculateConditionType(orderType) {
    switch (orderType) {
        case constants_1.OrderType.LIMIT:
            return proto_includes_1.Order_ConditionType.CONDITION_TYPE_UNSPECIFIED;
        case constants_1.OrderType.MARKET:
            return proto_includes_1.Order_ConditionType.CONDITION_TYPE_UNSPECIFIED;
        case constants_1.OrderType.STOP_LIMIT:
        case constants_1.OrderType.STOP_MARKET:
            return proto_includes_1.Order_ConditionType.CONDITION_TYPE_STOP_LOSS;
        case constants_1.OrderType.TAKE_PROFIT_LIMIT:
        case constants_1.OrderType.TAKE_PROFIT_MARKET:
            return proto_includes_1.Order_ConditionType.CONDITION_TYPE_TAKE_PROFIT;
        default:
            return proto_includes_1.Order_ConditionType.CONDITION_TYPE_UNSPECIFIED;
    }
}
function calculateConditionalOrderTriggerSubticks(orderType, atomicResolution, quantumConversionExponent, subticksPerTick, triggerPrice) {
    switch (orderType) {
        case constants_1.OrderType.LIMIT:
        case constants_1.OrderType.MARKET:
            return long_1.default.fromNumber(0);
        case constants_1.OrderType.STOP_LIMIT:
        case constants_1.OrderType.STOP_MARKET:
        case constants_1.OrderType.TAKE_PROFIT_LIMIT:
        case constants_1.OrderType.TAKE_PROFIT_MARKET:
            if (triggerPrice === undefined) {
                throw new Error('triggerPrice must be set if orderType is STOP_LIMIT, STOP_MARKET, TAKE_PROFIT_LIMIT, or TAKE_PROFIT_MARKET');
            }
            return calculateSubticks(triggerPrice, atomicResolution, quantumConversionExponent, subticksPerTick);
        default:
            return long_1.default.fromNumber(0);
    }
}
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY2hhaW4taGVscGVycy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9jbGllbnRzL2hlbHBlcnMvY2hhaW4taGVscGVycy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7OztBQVFBLHNCQUVDO0FBRUQsOENBVUM7QUFFRCw4Q0FZQztBQUVELHNDQUVDO0FBRUQsb0RBbUZDO0FBRUQsa0RBd0JDO0FBRUQsMERBVUM7QUFFRCx3REFtQkM7QUFFRCw0RkE2QkM7QUF2TkQsZ0RBQXdCO0FBRXhCLDRDQUVzQjtBQUN0Qiw4REFBK0Y7QUFDL0Ysb0NBQXNDO0FBRXRDLFNBQWdCLEtBQUssQ0FBQyxLQUFhLEVBQUUsSUFBWTtJQUMvQyxPQUFPLElBQUksQ0FBQyxLQUFLLENBQUMsS0FBSyxHQUFHLElBQUksQ0FBQyxHQUFHLElBQUksQ0FBQztBQUN6QyxDQUFDO0FBRUQsU0FBZ0IsaUJBQWlCLENBQy9CLElBQVksRUFDWixnQkFBd0IsRUFDeEIsZ0JBQXdCO0lBRXhCLE1BQU0sV0FBVyxHQUFHLElBQUksR0FBRyxFQUFFLElBQUksQ0FBQyxDQUFDLENBQUMsR0FBRyxnQkFBZ0IsQ0FBQyxDQUFDO0lBQ3pELE1BQU0sUUFBUSxHQUFHLEtBQUssQ0FBQyxXQUFXLEVBQUUsZ0JBQWdCLENBQUMsQ0FBQztJQUN0RCxtREFBbUQ7SUFDbkQsTUFBTSxNQUFNLEdBQUcsSUFBSSxDQUFDLEdBQUcsQ0FBQyxRQUFRLEVBQUUsZ0JBQWdCLENBQUMsQ0FBQztJQUNwRCxPQUFPLGNBQUksQ0FBQyxVQUFVLENBQUMsTUFBTSxDQUFDLENBQUM7QUFDakMsQ0FBQztBQUVELFNBQWdCLGlCQUFpQixDQUMvQixLQUFhLEVBQ2IsZ0JBQXdCLEVBQ3hCLHlCQUFpQyxFQUNqQyxlQUF1QjtJQUV2QixNQUFNLGdDQUFnQyxHQUFHLENBQUMsQ0FBQyxDQUFDO0lBQzVDLE1BQU0sUUFBUSxHQUFHLGdCQUFnQixHQUFHLHlCQUF5QixHQUFHLGdDQUFnQyxD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