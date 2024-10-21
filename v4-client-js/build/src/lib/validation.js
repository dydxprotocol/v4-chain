"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.validatePlaceOrderMessage = validatePlaceOrderMessage;
exports.validateCancelOrderMessage = validateCancelOrderMessage;
exports.validateTransferMessage = validateTransferMessage;
exports.verifyOrderFlags = verifyOrderFlags;
exports.isStatefulOrder = isStatefulOrder;
exports.isValidAddress = isValidAddress;
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
function isStatefulOrder(orderFlags) {
    return orderFlags === types_1.OrderFlags.LONG_TERM || orderFlags === types_1.OrderFlags.CONDITIONAL;
}
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoidmFsaWRhdGlvbi5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3NyYy9saWIvdmFsaWRhdGlvbi50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7OztBQWdCQSw4REEwQkM7QUFNRCxnRUEwQkM7QUFNRCwwREEyQkM7QUEwQkQsNENBR0M7QUFFRCwwQ0FFQztBQVlELHdDQUdDO0FBM0pELG1DQUFnQztBQUNoQyxnREFBd0I7QUFFeEIsb0RBQTBFO0FBQzFFLDRDQUswQjtBQUMxQixxQ0FBcUM7QUFFckM7OztHQUdHO0FBQ0gsU0FBZ0IseUJBQXlCLENBQ3ZDLGdCQUF3QixFQUN4QixLQUFrQjtJQUVsQixJQUFJLENBQUMsb0JBQW9CLENBQUMsS0FBSyxDQUFDLFFBQVEsQ0FBQyxFQUFFLENBQUM7UUFDMUMsT0FBTyxJQUFJLGtCQUFTLENBQUMsYUFBYSxLQUFLLENBQUMsUUFBUSx3QkFBd0IsQ0FBQyxDQUFDO0lBQzVFLENBQUM7SUFDRCxJQUFJLEtBQUssQ0FBQyxRQUFRLENBQUMsZUFBZSxDQUFDLGNBQUksQ0FBQyxJQUFJLENBQUMsRUFBRSxDQUFDO1FBQzlDLE9BQU8sSUFBSSxrQkFBUyxDQUFDLGFBQWEsS0FBSyxDQUFDLFFBQVEsaUJBQWlCLENBQUMsQ0FBQztJQUNyRSxDQUFDO0lBQ0QsSUFBSSxLQUFLLENBQUMsUUFBUSxDQUFDLGVBQWUsQ0FBQyxjQUFJLENBQUMsSUFBSSxDQUFDLEVBQUUsQ0FBQztRQUM5QyxPQUFPLElBQUksa0JBQVMsQ0FBQyxhQUFhLEtBQUssQ0FBQyxRQUFRLGlCQUFpQixDQUFDLENBQUM7SUFDckUsQ0FBQztJQUNELElBQUksQ0FBQyxzQkFBc0IsQ0FBQyxnQkFBZ0IsQ0FBQyxFQUFFLENBQUM7UUFDOUMsT0FBTyxJQUFJLGtCQUFTLENBQ2xCLHFCQUFxQixnQkFBZ0IsdUJBQXVCLGlDQUFxQixFQUFFLENBQ3BGLENBQUM7SUFDSixDQUFDO0lBQ0QsSUFBSSxDQUFDLGVBQWUsQ0FBQyxLQUFL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