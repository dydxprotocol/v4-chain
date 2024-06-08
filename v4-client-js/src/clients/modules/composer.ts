import { EncodeObject } from '@cosmjs/proto-signing';
import { MsgSend } from 'cosmjs-types/cosmos/bank/v1beta1/tx';
import { Coin } from 'cosmjs-types/cosmos/base/v1beta1/coin';
import Long from 'long';
import protobuf from 'protobufjs';

import {
  OrderId,
  Order,
  Order_ConditionType,
  Order_Side,
  Order_TimeInForce,
  MsgPlaceOrder,
  MsgCancelOrder,
  SubaccountId,
  MsgCreateTransfer,
  Transfer,
  MsgDepositToSubaccount,
  MsgWithdrawFromSubaccount,
} from './proto-includes';

protobuf.util.Long = Long;
protobuf.configure();

export class Composer {
  public composeMsgPlaceOrder(
    address: string,
    subaccountNumber: number,
    clientId: number,
    clobPairId: number,
    orderFlags: number,
    goodTilBlock: number,
    goodTilBlockTime: number,
    side: Order_Side,
    quantums: Long,
    subticks: Long,
    timeInForce: Order_TimeInForce,
    reduceOnly: boolean,
    clientMetadata: number,
    conditionType: Order_ConditionType = Order_ConditionType.CONDITION_TYPE_UNSPECIFIED,
    conditionalOrderTriggerSubticks: Long = Long.fromInt(0),
  ): EncodeObject {
    this.validateGoodTilBlockAndTime(orderFlags, goodTilBlock, goodTilBlockTime);

    const subaccountId: SubaccountId = {
      owner: address,
      number: subaccountNumber,
    };

    const orderId: OrderId = {
      subaccountId,
      clientId,
      orderFlags,
      clobPairId,
    };
    const order: Order = {
      orderId,
      side,
      quantums,
      subticks,
      goodTilBlock: goodTilBlock === 0 ? undefined : goodTilBlock,
      goodTilBlockTime: goodTilBlock === 0 ? goodTilBlockTime : undefined,
      timeInForce,
      reduceOnly,
      clientMetadata: clientMetadata ?? 0,
      conditionType,
      conditionalOrderTriggerSubticks,
    };
    const msg: MsgPlaceOrder = {
      order,
    };
    return {
      typeUrl: '/dydxprotocol.clob.MsgPlaceOrder',
      value: msg,
    };
  }

  public composeMsgCancelOrder(
    address: string,
    subaccountNumber: number,
    clientId: number,
    clobPairId: number,
    orderFlags: number,
    goodTilBlock: number,
    goodTilBlockTime: number,
  ): EncodeObject {
    this.validateGoodTilBlockAndTime(orderFlags, goodTilBlock, goodTilBlockTime);

    const subaccountId: SubaccountId = {
      owner: address,
      number: subaccountNumber,
    };

    const orderId: OrderId = {
      subaccountId,
      clientId,
      orderFlags,
      clobPairId,
    };

    const msg: MsgCancelOrder = {
      orderId,
      goodTilBlock: goodTilBlock === 0 ? undefined : goodTilBlock,
      goodTilBlockTime: goodTilBlock === 0 ? goodTilBlockTime : undefined,
    };

    return {
      typeUrl: '/dydxprotocol.clob.MsgCancelOrder',
      value: msg,
    };
  }

  public composeMsgTransfer(
    address: string,
    subaccountNumber: number,
    recipientAddress: string,
    recipientSubaccountNumber: number,
    assetId: number,
    amount: Long,
  ): EncodeObject {
    const sender: SubaccountId = {
      owner: address,
      number: subaccountNumber,
    };
    const recipient: SubaccountId = {
      owner: recipientAddress,
      number: recipientSubaccountNumber,
    };

    const transfer: Transfer = {
      sender,
      recipient,
      assetId,
      amount,
    };

    const msg: MsgCreateTransfer = {
      transfer,
    };

    return {
      typeUrl: '/dydxprotocol.sending.MsgCreateTransfer',
      value: msg,
    };
  }

  public composeMsgDepositToSubaccount(
    address: string,
    subaccountNumber: number,
    assetId: number,
    quantums: Long,
  ): EncodeObject {
    const recipient: SubaccountId = {
      owner: address,
      number: subaccountNumber,
    };

    const msg: MsgDepositToSubaccount = {
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

  public composeMsgWithdrawFromSubaccount(
    address: string,
    subaccountNumber: number,
    assetId: number,
    quantums: Long,
    recipient: string = address,
  ): EncodeObject {
    const sender: SubaccountId = {
      owner: address,
      number: subaccountNumber,
    };

    const msg: MsgWithdrawFromSubaccount = {
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

  public composeMsgSendToken(
    address: string,
    recipient: string,
    coinDenom: string,
    quantums: string,
  ): EncodeObject {
    const coin: Coin = {
      denom: coinDenom,
      amount: quantums,
    };

    const msg: MsgSend = {
      fromAddress: address,
      toAddress: recipient,
      amount: [coin],
    };

    return {
      typeUrl: '/cosmos.bank.v1beta1.MsgSend',
      value: msg,
    };
  }

  public validateGoodTilBlockAndTime(
    orderFlags: number,
    goodTilBlock: number,
    goodTilBlockTime: number,
  ): void {
    if (orderFlags === 0 && goodTilBlock === 0) {
      throw new Error('goodTilBlock must be set if orderFlags is 0');
    } else if (orderFlags !== 0 && goodTilBlockTime === 0) {
      throw new Error('goodTilBlockTime must be set if orderFlags is not 0');
    }
  }
}
