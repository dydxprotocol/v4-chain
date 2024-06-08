import { EncodeObject } from '@cosmjs/proto-signing';
import Long from 'long';
import { Order_ConditionType, Order_Side, Order_TimeInForce } from './proto-includes';
export declare class Composer {
    composeMsgPlaceOrder(address: string, subaccountNumber: number, clientId: number, clobPairId: number, orderFlags: number, goodTilBlock: number, goodTilBlockTime: number, side: Order_Side, quantums: Long, subticks: Long, timeInForce: Order_TimeInForce, reduceOnly: boolean, clientMetadata: number, conditionType?: Order_ConditionType, conditionalOrderTriggerSubticks?: Long): EncodeObject;
    composeMsgCancelOrder(address: string, subaccountNumber: number, clientId: number, clobPairId: number, orderFlags: number, goodTilBlock: number, goodTilBlockTime: number): EncodeObject;
    composeMsgTransfer(address: string, subaccountNumber: number, recipientAddress: string, recipientSubaccountNumber: number, assetId: number, amount: Long): EncodeObject;
    composeMsgDepositToSubaccount(address: string, subaccountNumber: number, assetId: number, quantums: Long): EncodeObject;
    composeMsgWithdrawFromSubaccount(address: string, subaccountNumber: number, assetId: number, quantums: Long, recipient?: string): EncodeObject;
    composeMsgSendToken(address: string, recipient: string, coinDenom: string, quantums: string): EncodeObject;
    validateGoodTilBlockAndTime(orderFlags: number, goodTilBlock: number, goodTilBlockTime: number): void;
}
