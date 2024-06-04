"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.isValidAddress = exports.isStatefulOrder = exports.verifyOrderFlags = exports.validateTransferMessage = exports.validateCancelOrderMessage = exports.validatePlaceOrderMessage = void 0;
const bech32_1 = require("bech32");
const long_1 = __importDefault(require("long"));
const constants_1 = require("../clients/constants");
const types_1 = require("../clients/types");
const errors_1 = require("./errors");
/**
 * @describe validatePlaceOrderMessage validates that an order to place has fields that would be
 *  valid on-chain.
 */
function validatePlaceOrderMessage(subaccountNumber, order) {
    if (!verifyNumberIsUint32(order.clientId)) {
        return new errors_1.UserError(`clientId: ${order.clientId} is not a valid uint32`);
    }
    if (order.quantums.lessThanOrEqual(long_1.default.ZERO)) {
        return new errors_1.UserError(`quantums: ${order.quantums} cannot be <= 0`);
    }
    if (order.subticks.lessThanOrEqual(long_1.default.ZERO)) {
        return new errors_1.UserError(`subticks: ${order.subticks} cannot be <= 0`);
    }
    if (!verifySubaccountNumber(subaccountNumber)) {
        return new errors_1.UserError(`subaccountNumber: ${subaccountNumber} cannot be < 0 or > ${constants_1.MAX_SUBACCOUNT_NUMBER}`);
    }
    if (!isStatefulOrder(order.orderFlags) && !verifyGoodTilBlock(order.goodTilBlock)) {
        return new errors_1.UserError(`goodTilBlock: ${order.goodTilBlock} is not a valid uint32 or is 0`);
    }
    if (isStatefulOrder(order.orderFlags) && !verifyGoodTilBlockTime(order.goodTilBlockTime)) {
        return new errors_1.UserError(`goodTilBlockTime: ${order.goodTilBlockTime} is not a valid uint32 or is 0`);
    }
    return undefined;
}
exports.validatePlaceOrderMessage = validatePlaceOrderMessage;
/**
 * @describe validateCancelOrderMessage validates that an order to cancel has fields that would be
 *  valid on-chain.
 */
function validateCancelOrderMessage(subaccountNumber, order) {
    if (!verifyNumberIsUint32(order.clientId)) {
        return new errors_1.UserError(`clientId: ${order.clientId} is not a valid uint32`);
    }
    if (!isStatefulOrder(order.orderFlags) && !verifyGoodTilBlock(order.goodTilBlock)) {
        return new errors_1.UserError(`goodTilBlock: ${order.goodTilBlock} is not a valid uint32 or is 0`);
    }
    if (!isStatefulOrder(order.orderFlags) && order.goodTilBlockTime !== undefined) {
        return new errors_1.UserError(`goodTilBlockTime is ${order.goodTilBlockTime}, but should not be set for non-stateful orders`);
    }
    if (isStatefulOrder(order.orderFlags) && !verifyGoodTilBlockTime(order.goodTilBlockTime)) {
        return new errors_1.UserError(`goodTilBlockTime: ${order.goodTilBlockTime} is not a valid uint32 or is 0`);
    }
    if (isStatefulOrder(order.orderFlags) && order.goodTilBlock !== undefined) {
        return new errors_1.UserError(`goodTilBlock is ${order.goodTilBlock}, but should not be set for stateful orders`);
    }
    if (!verifySubaccountNumber(subaccountNumber)) {
        return new errors_1.UserError(`subaccountNumber: ${subaccountNumber} cannot be < 0 or > ${constants_1.MAX_SUBACCOUNT_NUMBER}`);
    }
    return undefined;
}
exports.validateCancelOrderMessage = validateCancelOrderMessage;
/**
 * @describe validateTransferMessage validates that a transfer to place has fields that would be
 *  valid on-chain.
 */
function validateTransferMessage(transfer) {
    if (!verifySubaccountNumber(transfer.sender.number || 0)) {
        return new errors_1.UserError(`senderSubaccountNumber: ${transfer.sender.number || 0} cannot be < 0 or > ${constants_1.MAX_SUBACCOUNT_NUMBER}`);
    }
    if (!verifySubaccountNumber(transfer.recipient.number || 0)) {
        return new errors_1.UserError(`recipientSubaccountNumber: ${transfer.recipient.number || 0} cannot be < 0 or > ${constants_1.MAX_SUBACCOUNT_NUMBER}`);
    }
    if (transfer.assetId !== 0) {
        return new errors_1.UserError(`asset id: ${transfer.assetId} not supported`);
    }
    if (transfer.amount.lessThanOrEqual(long_1.default.ZERO)) {
        return new errors_1.UserError(`amount: ${transfer.amount} cannot be <= 0`);
    }
    const addressError = verifyIsBech32(transfer.recipient.owner);
    if (addressError !== undefined) {
        return new errors_1.UserError(addressError.toString());
    }
    return undefined;
}
exports.validateTransferMessage = validateTransferMessage;
function verifyGoodTilBlock(goodTilBlock) {
    if (goodTilBlock === undefined) {
        return false;
    }
    return verifyNumberIsUint32(goodTilBlock) && goodTilBlock > 0;
}
function verifyGoodTilBlockTime(goodTilBlockTime) {
    if (goodTilBlockTime === undefined) {
        return false;
    }
    return verifyNumberIsUint32(goodTilBlockTime) && goodTilBlockTime > 0;
}
function verifySubaccountNumber(subaccountNumber) {
    return subaccountNumber >= 0 && subaccountNumber <= constants_1.MAX_SUBACCOUNT_NUMBER;
}
function verifyNumberIsUint32(num) {
    return num >= 0 && num <= constants_1.MAX_UINT_32;
}
function verifyOrderFlags(orderFlags) {
    return orderFlags === types_1.OrderFlags.SHORT_TERM ||
        orderFlags === types_1.OrderFlags.LONG_TERM || orderFlags === types_1.OrderFlags.CONDITIONAL;
}
exports.verifyOrderFlags = verifyOrderFlags;
function isStatefulOrder(orderFlags) {
    return orderFlags === types_1.OrderFlags.LONG_TERM || orderFlags === types_1.OrderFlags.CONDITIONAL;
}
exports.isStatefulOrder = isStatefulOrder;
function verifyIsBech32(address) {
    try {
        (0, bech32_1.decode)(address);
    }
    catch (error) {
        return error;
    }
    return undefined;
}
function isValidAddress(address) {
    // An address is valid if it starts with `dydx1` and is Bech32 format.
    return address.startsWith('dydx1') && (verifyIsBech32(address) === undefined);
}
exports.isValidAddress = isValidAddress;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidmFsaWRhdGlvbi5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9saWIvdmFsaWRhdGlvbi50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFBQSxtQ0FBZ0M7QUFDaEMsZ0RBQXdCO0FBRXhCLG9EQUEwRTtBQUMxRSw0Q0FLMEI7QUFDMUIscUNBQXFDO0FBRXJDOzs7R0FHRztBQUNILFNBQWdCLHlCQUF5QixDQUN2QyxnQkFBd0IsRUFDeEIsS0FBa0I7SUFFbEIsSUFBSSxDQUFDLG9CQUFvQixDQUFDLEtBQUssQ0FBQyxRQUFRLENBQUMsRUFBRTtRQUN6QyxPQUFPLElBQUksa0JBQVMsQ0FBQyxhQUFhLEtBQUssQ0FBQyxRQUFRLHdCQUF3QixDQUFDLENBQUM7S0FDM0U7SUFDRCxJQUFJLEtBQUssQ0FBQyxRQUFRLENBQUMsZUFBZSxDQUFDLGNBQUksQ0FBQyxJQUFJLENBQUMsRUFBRTtRQUM3QyxPQUFPLElBQUksa0JBQVMsQ0FBQyxhQUFhLEtBQUssQ0FBQyxRQUFRLGlCQUFpQixDQUFDLENBQUM7S0FDcEU7SUFDRCxJQUFJLEtBQUssQ0FBQyxRQUFRLENBQUMsZUFBZSxDQUFDLGNBQUksQ0FBQyxJQUFJLENBQUMsRUFBRTtRQUM3QyxPQUFPLElBQUksa0JBQVMsQ0FBQyxhQUFhLEtBQUssQ0FBQyxRQUFRLGlCQUFpQixDQUFDLENBQUM7S0FDcEU7SUFDRCxJQUFJLENBQUMsc0JBQXNCLENBQUMsZ0JBQWdCLENBQUMsRUFBRTtRQUM3QyxPQUFPLElBQUksa0JBQVMsQ0FDbEIscUJBQXFCLGdCQUFnQix1QkFBdUIsaUNBQXFCLEVBQUUsQ0FDcEYsQ0FBQztLQUNIO0lBQ0QsSUFBSSxDQUFDLGVBQWUsQ0FBQyxLQUFLLENBQUMsVUFBVSxDQUFDLElBQUksQ0FBQyxrQkFBa0IsQ0FBQyxLQUFLLENBQUMsWUFBWSxDQUFDLEVBQUU7UUFDakYsT0FBTyxJQUFJLGtCQUFTLENBQUMsaUJBQWlCLEtBQUssQ0FBQyxZQUFZ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