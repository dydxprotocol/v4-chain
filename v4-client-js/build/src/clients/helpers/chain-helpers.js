"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.calculateConditionalOrderTriggerSubticks = exports.calculateConditionType = exports.calculateClientMetadata = exports.calculateOrderFlags = exports.calculateTimeInForce = exports.calculateSide = exports.calculateSubticks = exports.calculateQuantums = exports.round = void 0;
const long_1 = __importDefault(require("long"));
const constants_1 = require("../constants");
const proto_includes_1 = require("../modules/proto-includes");
const types_1 = require("../types");
function round(input, base) {
    return Math.floor(input / base) * base;
}
exports.round = round;
function calculateQuantums(size, atomicResolution, stepBaseQuantums) {
    const rawQuantums = size * 10 ** (-1 * atomicResolution);
    const quantums = round(rawQuantums, stepBaseQuantums);
    // stepBaseQuantums functions as minimum order size
    const result = Math.max(quantums, stepBaseQuantums);
    return long_1.default.fromNumber(result);
}
exports.calculateQuantums = calculateQuantums;
function calculateSubticks(price, atomicResolution, quantumConversionExponent, subticksPerTick) {
    const QUOTE_QUANTUMS_ATOMIC_RESOLUTION = -6;
    const exponent = atomicResolution - quantumConversionExponent - QUOTE_QUANTUMS_ATOMIC_RESOLUTION;
    const rawSubticks = price * 10 ** exponent;
    const subticks = round(rawSubticks, subticksPerTick);
    const result = Math.max(subticks, subticksPerTick);
    return long_1.default.fromNumber(result);
}
exports.calculateSubticks = calculateSubticks;
function calculateSide(side) {
    return side === constants_1.OrderSide.BUY ? proto_includes_1.Order_Side.SIDE_BUY : proto_includes_1.Order_Side.SIDE_SELL;
}
exports.calculateSide = calculateSide;
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
exports.calculateTimeInForce = calculateTimeInForce;
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
exports.calculateOrderFlags = calculateOrderFlags;
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
exports.calculateClientMetadata = calculateClientMetadata;
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
exports.calculateConditionType = calculateConditionType;
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
exports.calculateConditionalOrderTriggerSubticks = calculateConditionalOrderTriggerSubticks;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY2hhaW4taGVscGVycy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3NyYy9jbGllbnRzL2hlbHBlcnMvY2hhaW4taGVscGVycy50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFBQSxnREFBd0I7QUFFeEIsNENBRXNCO0FBQ3RCLDhEQUErRjtBQUMvRixvQ0FBc0M7QUFFdEMsU0FBZ0IsS0FBSyxDQUFDLEtBQWEsRUFBRSxJQUFZO0lBQy9DLE9BQU8sSUFBSSxDQUFDLEtBQUssQ0FBQyxLQUFLLEdBQUcsSUFBSSxDQUFDLEdBQUcsSUFBSSxDQUFDO0FBQ3pDLENBQUM7QUFGRCxzQkFFQztBQUVELFNBQWdCLGlCQUFpQixDQUMvQixJQUFZLEVBQ1osZ0JBQXdCLEVBQ3hCLGdCQUF3QjtJQUV4QixNQUFNLFdBQVcsR0FBRyxJQUFJLEdBQUcsRUFBRSxJQUFJLENBQUMsQ0FBQyxDQUFDLEdBQUcsZ0JBQWdCLENBQUMsQ0FBQztJQUN6RCxNQUFNLFFBQVEsR0FBRyxLQUFLLENBQUMsV0FBVyxFQUFFLGdCQUFnQixDQUFDLENBQUM7SUFDdEQsbURBQW1EO0lBQ25ELE1BQU0sTUFBTSxHQUFHLElBQUksQ0FBQyxHQUFHLENBQUMsUUFBUSxFQUFFLGdCQUFnQixDQUFDLENBQUM7SUFDcEQsT0FBTyxjQUFJLENBQUMsVUFBVSxDQUFDLE1BQU0sQ0FBQyxDQUFDO0FBQ2pDLENBQUM7QUFWRCw4Q0FVQztBQUVELFNBQWdCLGlCQUFpQixDQUMvQixLQUFhLEVBQ2IsZ0JBQXdCLEVBQ3hCLHlCQUFpQyxFQUNqQyxlQUF1QjtJQUV2QixNQUFNLGdDQUFnQyxHQUFHLENBQUMsQ0FBQyxDQUFDO0lBQzVDLE1BQU0sUUFBUSxHQUFHLGdCQUFnQixHQUFHLHlCQUF5QixHQUFHLGdDQUFnQyxDQUFDO0lBQ2pHLE1BQU0sV0FBVyxHQUFHLEtBQUssR0FBRyxFQUFFLElBQUksUUFBUSxDQUFDO0lBQzNDLE1BQU0sUUFBUSxHQUFHLEtBQUss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