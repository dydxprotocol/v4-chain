"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Composer = void 0;
const long_1 = __importDefault(require("long"));
const protobufjs_1 = __importDefault(require("protobufjs"));
const proto_includes_1 = require("./proto-includes");
protobufjs_1.default.util.Long = long_1.default;
protobufjs_1.default.configure();
class Composer {
    composeMsgPlaceOrder(address, subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock, goodTilBlockTime, side, quantums, subticks, timeInForce, reduceOnly, clientMetadata, conditionType = proto_includes_1.Order_ConditionType.CONDITION_TYPE_UNSPECIFIED, conditionalOrderTriggerSubticks = long_1.default.fromInt(0)) {
        this.validateGoodTilBlockAndTime(orderFlags, goodTilBlock, goodTilBlockTime);
        const subaccountId = {
            owner: address,
            number: subaccountNumber,
        };
        const orderId = {
            subaccountId,
            clientId,
            orderFlags,
            clobPairId,
        };
        const order = {
            orderId,
            side,
            quantums,
            subticks,
            goodTilBlock: goodTilBlock === 0 ? undefined : goodTilBlock,
            goodTilBlockTime: goodTilBlock === 0 ? goodTilBlockTime : undefined,
            timeInForce,
            reduceOnly,
            clientMetadata: clientMetadata !== null && clientMetadata !== void 0 ? clientMetadata : 0,
            conditionType,
            conditionalOrderTriggerSubticks,
        };
        const msg = {
            order,
        };
        return {
            typeUrl: '/dydxprotocol.clob.MsgPlaceOrder',
            value: msg,
        };
    }
    composeMsgCancelOrder(address, subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock, goodTilBlockTime) {
        this.validateGoodTilBlockAndTime(orderFlags, goodTilBlock, goodTilBlockTime);
        const subaccountId = {
            owner: address,
            number: subaccountNumber,
        };
        const orderId = {
            subaccountId,
            clientId,
            orderFlags,
            clobPairId,
        };
        const msg = {
            orderId,
            goodTilBlock: goodTilBlock === 0 ? undefined : goodTilBlock,
            goodTilBlockTime: goodTilBlock === 0 ? goodTilBlockTime : undefined,
        };
        return {
            typeUrl: '/dydxprotocol.clob.MsgCancelOrder',
            value: msg,
        };
    }
    composeMsgTransfer(address, subaccountNumber, recipientAddress, recipientSubaccountNumber, assetId, amount) {
        const sender = {
            owner: address,
            number: subaccountNumber,
        };
        const recipient = {
            owner: recipientAddress,
            number: recipientSubaccountNumber,
        };
        const transfer = {
            sender,
            recipient,
            assetId,
            amount,
        };
        const msg = {
            transfer,
        };
        return {
            typeUrl: '/dydxprotocol.sending.MsgCreateTransfer',
            value: msg,
        };
    }
    composeMsgDepositToSubaccount(address, subaccountNumber, assetId, quantums) {
        const recipient = {
            owner: address,
            number: subaccountNumber,
        };
        const msg = {
            sender: address,
            recipient,
            assetId,
            quantums,
        };
        return {
            typeUrl: '/dydxprotocol.sending.MsgDepositToSubaccount',
            value: msg,
        };
    }
    composeMsgWithdrawFromSubaccount(address, subaccountNumber, assetId, quantums, recipient = address) {
        const sender = {
            owner: address,
            number: subaccountNumber,
        };
        const msg = {
            sender,
            recipient,
            assetId,
            quantums,
        };
        return {
            typeUrl: '/dydxprotocol.sending.MsgWithdrawFromSubaccount',
            value: msg,
        };
    }
    composeMsgSendToken(address, recipient, coinDenom, quantums) {
        const coin = {
            denom: coinDenom,
            amount: quantums,
        };
        const msg = {
            fromAddress: address,
            toAddress: recipient,
            amount: [coin],
        };
        return {
            typeUrl: '/cosmos.bank.v1beta1.MsgSend',
            value: msg,
        };
    }
    validateGoodTilBlockAndTime(orderFlags, goodTilBlock, goodTilBlockTime) {
        if (orderFlags === 0 && goodTilBlock === 0) {
            throw new Error('goodTilBlock must be set if orderFlags is 0');
        }
        else if (orderFlags !== 0 && goodTilBlockTime === 0) {
            throw new Error('goodTilBlockTime must be set if orderFlags is not 0');
        }
    }
}
exports.Composer = Composer;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29tcG9zZXIuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvY2xpZW50cy9tb2R1bGVzL2NvbXBvc2VyLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7OztBQUdBLGdEQUF3QjtBQUN4Qiw0REFBa0M7QUFFbEMscURBYTBCO0FBRTFCLG9CQUFRLENBQUMsSUFBSSxDQUFDLElBQUksR0FBRyxjQUFJLENBQUM7QUFDMUIsb0JBQVEsQ0FBQyxTQUFTLEVBQUUsQ0FBQztBQUVyQixNQUFhLFFBQVE7SUFDWixvQkFBb0IsQ0FDekIsT0FBZSxFQUNmLGdCQUF3QixFQUN4QixRQUFnQixFQUNoQixVQUFrQixFQUNsQixVQUFrQixFQUNsQixZQUFvQixFQUNwQixnQkFBd0IsRUFDeEIsSUFBZ0IsRUFDaEIsUUFBYyxFQUNkLFFBQWMsRUFDZCxXQUE4QixFQUM5QixVQUFtQixFQUNuQixjQUFzQixFQUN0QixnQkFBcUMsb0NBQW1CLENBQUMsMEJBQTBCLEVBQ25GLGtDQUF3QyxjQUFJLENBQUMsT0FBTyxDQUFDLENBQUMsQ0FBQztRQUV2RCxJQUFJLENBQUMsMkJBQTJCLENBQUMsVUFBVSxFQUFFLFlBQVksRUFBRSxnQkFBZ0IsQ0FBQyxDQUFDO1FBRTdFLE1BQU0sWUFBWSxHQUFpQjtZQUNqQyxLQUFLLEVBQUUsT0FBTztZQUNkLE1BQU0sRUFBRSxnQkFBZ0I7U0FDekIsQ0FBQztRQUVGLE1BQU0sT0FBTyxHQUFZO1lBQ3ZCLFlBQVk7WUFDWixRQUFRO1lBQ1IsVUFBVTtZQUNWLFVBQVU7U0FDWCxDQUFDO1FBQ0YsTUFBTSxLQUFLLEdBQVU7WUFDbkIsT0FBTztZQUNQLElBQUk7WUFDSixRQUFRO1lBQ1IsUUFBUTtZQUNSLFlBQVksRUFBRSxZQUFZLEtBQUssQ0FBQyxDQUFDLENBQUMsQ0FBQyxTQUFTLENBQUMsQ0FBQyxDQUFDLFlBQVk7WUFDM0QsZ0JBQWdCLEVBQUUsWUFBWSxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMs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