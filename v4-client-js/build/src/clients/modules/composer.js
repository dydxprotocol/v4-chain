"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Composer = void 0;
const long_1 = __importDefault(require("long"));
const protobufjs_1 = __importDefault(require("protobufjs"));
const proto_includes_1 = require("./proto-includes");
// import {
//   OrderId,
//   Order,
//   Order_ConditionType,
//   Order_Side,
//   Order_TimeInForce,
// } from '../../../../../indexer/packages/v4-protos/src/codegen/dydxprotocol/clob/order';
protobufjs_1.default.util.Long = long_1.default;
protobufjs_1.default.configure();
class Composer {
    composeMsgPlaceOrder(address, subaccountNumber, clientId, clobPairId, orderFlags, goodTilBlock, goodTilBlockTime, side, quantums, subticks, timeInForce, reduceOnly, clientMetadata, conditionType = proto_includes_1.Order_ConditionType.CONDITION_TYPE_UNSPECIFIED, conditionalOrderTriggerSubticks = long_1.default.fromInt(0), routerFeePpm, routerFeeSubaccountOwner, routerFeeSubaccountNumber) {
        this.validateGoodTilBlockAndTime(orderFlags, goodTilBlock, goodTilBlockTime);
        const subaccountId = {
            owner: address,
            number: subaccountNumber,
        };
        const routerFeeSubaccountId = {
            owner: routerFeeSubaccountOwner,
            number: routerFeeSubaccountNumber,
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
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiY29tcG9zZXIuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi9zcmMvY2xpZW50cy9tb2R1bGVzL2NvbXBvc2VyLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiI7Ozs7OztBQUdBLGdEQUF3QjtBQUN4Qiw0REFBa0M7QUFFbEMscURBYTBCO0FBRTFCLFdBQVc7QUFDWCxhQUFhO0FBQ2IsV0FBVztBQUNYLHlCQUF5QjtBQUN6QixnQkFBZ0I7QUFDaEIsdUJBQXVCO0FBQ3ZCLDBGQUEwRjtBQUUxRixvQkFBUSxDQUFDLElBQUksQ0FBQyxJQUFJLEdBQUcsY0FBSSxDQUFDO0FBQzFCLG9CQUFRLENBQUMsU0FBUyxFQUFFLENBQUM7QUFFckIsTUFBYSxRQUFRO0lBQ1osb0JBQW9CLENBQ3pCLE9BQWUsRUFDZixnQkFBd0IsRUFDeEIsUUFBZ0IsRUFDaEIsVUFBa0IsRUFDbEIsVUFBa0IsRUFDbEIsWUFBb0IsRUFDcEIsZ0JBQXdCLEVBQ3hCLElBQWdCLEVBQ2hCLFFBQWMsRUFDZCxRQUFjLEVBQ2QsV0FBOEIsRUFDOUIsVUFBbUIsRUFDbkIsY0FBc0IsRUFDdEIsZ0JBQXFDLG9DQUFtQixDQUFDLDBCQUEwQixFQUNuRixrQ0FBd0MsY0FBSSxDQUFDLE9BQU8sQ0FBQyxDQUFDLENBQUMsRUFDdkQsWUFBb0IsRUFDcEIsd0JBQWdDLEVBQ2hDLHlCQUFpQztRQUVqQyxJQUFJLENBQUMsMkJBQTJCLENBQUMsVUFBVSxFQUFFLFlBQVksRUFBRSxnQkFBZ0IsQ0FBQyxDQUFDO1FBRTdFLE1BQU0sWUFBWSxHQUFpQjtZQUNqQyxLQUFLLEVBQUUsT0FBTztZQUNkLE1BQU0sRUFBRSxnQkFBZ0I7U0FDekIsQ0FBQztRQUVGLE1BQU0scUJBQXFCLEdBQWlCO1lBQzFDLEtBQUssRUFBRSx3QkFBd0I7WUFDL0IsTUFBTSxFQUFFLHlCQUF5QjtTQUNsQyxDQUFDO1FBRUYsTUFBTSxPQUFPLEdBQVk7WUFDdkIsWUFBWTtZQUNaLFFBQVE7WUFDUixVQUFV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